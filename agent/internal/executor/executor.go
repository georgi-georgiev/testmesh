package executor

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	"gopkg.in/yaml.v3"
)

// Config holds executor configuration
type Config struct {
	APIURL        string
	AgentID       string
	Token         string
	MaxConcurrent int
	PollInterval  time.Duration
}

// Executor handles flow execution on the agent
type Executor struct {
	config     *Config
	client     *http.Client
	sem        chan struct{}
	wg         sync.WaitGroup
	mu         sync.RWMutex
	running    map[string]*Execution
}

// Execution represents a running execution
type Execution struct {
	ID        string
	FlowID    string
	StartTime time.Time
	Cancel    context.CancelFunc
}

// Job represents a job from the server
type Job struct {
	ID          string                 `json:"id"`
	FlowID      string                 `json:"flow_id"`
	FlowYAML    string                 `json:"flow_yaml"`
	Environment map[string]string      `json:"environment"`
	Variables   map[string]interface{} `json:"variables"`
}

// New creates a new executor
func New(config *Config) (*Executor, error) {
	if config.PollInterval == 0 {
		config.PollInterval = 5 * time.Second
	}
	if config.MaxConcurrent == 0 {
		config.MaxConcurrent = 5
	}

	return &Executor{
		config:  config,
		client:  &http.Client{Timeout: 30 * time.Second},
		sem:     make(chan struct{}, config.MaxConcurrent),
		running: make(map[string]*Execution),
	}, nil
}

// Start starts the executor loop
func (e *Executor) Start(ctx context.Context) {
	fmt.Println("📋 Executor started, polling for jobs...")

	ticker := time.NewTicker(e.config.PollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			e.shutdown()
			return
		case <-ticker.C:
			e.pollForJobs(ctx)
		}
	}
}

func (e *Executor) pollForJobs(ctx context.Context) {
	// Check if we have capacity
	select {
	case e.sem <- struct{}{}:
		// Have capacity, try to get a job
	default:
		// At capacity, skip this poll
		return
	}

	job, err := e.fetchJob(ctx)
	if err != nil {
		<-e.sem // Release semaphore
		if err != errNoJobs {
			fmt.Printf("Error fetching job: %v\n", err)
		}
		return
	}

	// Execute job in background
	e.wg.Add(1)
	go func() {
		defer e.wg.Done()
		defer func() { <-e.sem }()
		e.executeJob(ctx, job)
	}()
}

var errNoJobs = fmt.Errorf("no jobs available")

func (e *Executor) fetchJob(ctx context.Context) (*Job, error) {
	url := fmt.Sprintf("%s/api/v1/agents/%s/jobs/next", e.config.APIURL, e.config.AgentID)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+e.config.Token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := e.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNoContent {
		return nil, errNoJobs
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("server error: %s", string(body))
	}

	var job Job
	if err := json.NewDecoder(resp.Body).Decode(&job); err != nil {
		return nil, err
	}

	return &job, nil
}

func (e *Executor) executeJob(ctx context.Context, job *Job) {
	execCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	exec := &Execution{
		ID:        uuid.New().String(),
		FlowID:    job.FlowID,
		StartTime: time.Now(),
		Cancel:    cancel,
	}

	e.mu.Lock()
	e.running[exec.ID] = exec
	e.mu.Unlock()

	defer func() {
		e.mu.Lock()
		delete(e.running, exec.ID)
		e.mu.Unlock()
	}()

	fmt.Printf("▶️  Executing job %s (flow: %s)\n", job.ID, job.FlowID)

	// Parse flow
	var flow map[string]interface{}
	if err := yaml.Unmarshal([]byte(job.FlowYAML), &flow); err != nil {
		e.reportResult(ctx, job.ID, "failed", err.Error(), nil)
		return
	}

	// Execute steps
	results := make([]StepResult, 0)
	steps, _ := flow["steps"].([]interface{})

	for i, step := range steps {
		select {
		case <-execCtx.Done():
			e.reportResult(ctx, job.ID, "cancelled", "execution cancelled", results)
			return
		default:
		}

		stepMap, ok := step.(map[string]interface{})
		if !ok {
			continue
		}

		stepName, _ := stepMap["name"].(string)
		if stepName == "" {
			stepName = fmt.Sprintf("Step %d", i+1)
		}

		result := e.executeStep(execCtx, stepMap, job.Environment, job.Variables)
		result.Name = stepName
		results = append(results, result)

		if result.Status == "failed" {
			e.reportResult(ctx, job.ID, "failed", result.Error, results)
			return
		}
	}

	e.reportResult(ctx, job.ID, "passed", "", results)
	fmt.Printf("✅ Job %s completed\n", job.ID)
}

// StepResult holds the result of a step execution
type StepResult struct {
	Name     string                 `json:"name"`
	Status   string                 `json:"status"`
	Duration int64                  `json:"duration_ms"`
	Output   map[string]interface{} `json:"output,omitempty"`
	Error    string                 `json:"error,omitempty"`
}

func (e *Executor) executeStep(ctx context.Context, step map[string]interface{}, env map[string]string, vars map[string]interface{}) StepResult {
	start := time.Now()
	result := StepResult{
		Status: "passed",
		Output: make(map[string]interface{}),
	}

	defer func() {
		result.Duration = time.Since(start).Milliseconds()
	}()

	action, ok := step["action"].(map[string]interface{})
	if !ok {
		result.Status = "failed"
		result.Error = "missing action"
		return result
	}

	actionType, _ := action["type"].(string)

	switch actionType {
	case "http":
		return e.executeHTTP(ctx, action, env, vars)
	case "delay", "sleep":
		return e.executeDelay(ctx, action)
	case "log":
		return e.executeLog(action)
	default:
		result.Status = "failed"
		result.Error = fmt.Sprintf("unsupported action type: %s", actionType)
	}

	return result
}

func (e *Executor) executeHTTP(ctx context.Context, action map[string]interface{}, env map[string]string, vars map[string]interface{}) StepResult {
	start := time.Now()
	result := StepResult{
		Status: "passed",
		Output: make(map[string]interface{}),
	}

	http_, ok := action["http"].(map[string]interface{})
	if !ok {
		result.Status = "failed"
		result.Error = "missing http configuration"
		return result
	}

	url, _ := http_["url"].(string)
	method, _ := http_["method"].(string)
	if method == "" {
		method = "GET"
	}

	var bodyReader io.Reader
	if body := http_["body"]; body != nil {
		bodyBytes, _ := json.Marshal(body)
		bodyReader = bytes.NewReader(bodyBytes)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
	if err != nil {
		result.Status = "failed"
		result.Error = err.Error()
		return result
	}

	// Add headers
	if headers, ok := http_["headers"].(map[string]interface{}); ok {
		for k, v := range headers {
			req.Header.Set(k, fmt.Sprintf("%v", v))
		}
	}

	resp, err := e.client.Do(req)
	if err != nil {
		result.Status = "failed"
		result.Error = err.Error()
		return result
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	result.Duration = time.Since(start).Milliseconds()
	result.Output["status_code"] = resp.StatusCode
	result.Output["body"] = string(respBody)

	// Parse JSON response
	var jsonBody interface{}
	if json.Unmarshal(respBody, &jsonBody) == nil {
		result.Output["json"] = jsonBody
	}

	return result
}

func (e *Executor) executeDelay(ctx context.Context, action map[string]interface{}) StepResult {
	result := StepResult{Status: "passed"}

	delay, _ := action["delay"].(string)
	if delay == "" {
		delay = "1s"
	}

	duration, err := time.ParseDuration(delay)
	if err != nil {
		result.Status = "failed"
		result.Error = fmt.Sprintf("invalid delay: %s", delay)
		return result
	}

	select {
	case <-time.After(duration):
	case <-ctx.Done():
		result.Status = "cancelled"
	}

	result.Duration = int64(duration.Milliseconds())
	return result
}

func (e *Executor) executeLog(action map[string]interface{}) StepResult {
	message, _ := action["message"].(string)
	fmt.Printf("   LOG: %s\n", message)
	return StepResult{Status: "passed"}
}

func (e *Executor) reportResult(ctx context.Context, jobID, status, errorMsg string, steps []StepResult) {
	url := fmt.Sprintf("%s/api/v1/agents/%s/jobs/%s/result", e.config.APIURL, e.config.AgentID, jobID)

	report := map[string]interface{}{
		"status": status,
		"steps":  steps,
	}
	if errorMsg != "" {
		report["error"] = errorMsg
	}

	body, _ := json.Marshal(report)
	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+e.config.Token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := e.client.Do(req)
	if err != nil {
		fmt.Printf("Warning: failed to report result: %v\n", err)
		return
	}
	resp.Body.Close()
}

func (e *Executor) shutdown() {
	fmt.Println("🛑 Stopping executor...")

	// Cancel all running executions
	e.mu.RLock()
	for _, exec := range e.running {
		exec.Cancel()
	}
	e.mu.RUnlock()

	// Wait for all to complete
	e.wg.Wait()
	fmt.Println("✅ Executor stopped")
}

// RunningCount returns the number of currently running executions
func (e *Executor) RunningCount() int {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return len(e.running)
}
