package heartbeat

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Config holds heartbeat configuration
type Config struct {
	APIURL   string
	AgentID  string
	Token    string
	Tags     string
	Interval time.Duration
}

// Manager handles agent heartbeats
type Manager struct {
	config   *Config
	client   *http.Client
	agentID  string
	hostname string
}

// New creates a new heartbeat manager
func New(config *Config) *Manager {
	hostname, _ := os.Hostname()

	agentID := config.AgentID
	if agentID == "" {
		agentID = uuid.New().String()
	}

	return &Manager{
		config:   config,
		client:   &http.Client{Timeout: 10 * time.Second},
		agentID:  agentID,
		hostname: hostname,
	}
}

// AgentInfo holds agent information for registration
type AgentInfo struct {
	ID       string            `json:"id"`
	Hostname string            `json:"hostname"`
	Version  string            `json:"version"`
	Platform string            `json:"platform"`
	Arch     string            `json:"arch"`
	Tags     []string          `json:"tags"`
	Metadata map[string]string `json:"metadata"`
}

// Register registers the agent with the server
func (m *Manager) Register(ctx context.Context) error {
	info := AgentInfo{
		ID:       m.agentID,
		Hostname: m.hostname,
		Version:  "1.0.0",
		Platform: runtime.GOOS,
		Arch:     runtime.GOARCH,
		Tags:     m.parseTags(),
		Metadata: map[string]string{
			"go_version": runtime.Version(),
			"num_cpu":    fmt.Sprintf("%d", runtime.NumCPU()),
		},
	}

	body, err := json.Marshal(info)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s/api/v1/agents/register", m.config.APIURL)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+m.config.Token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := m.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to register: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("registration failed: %s", string(respBody))
	}

	return nil
}

// Deregister deregisters the agent from the server
func (m *Manager) Deregister(ctx context.Context) error {
	url := fmt.Sprintf("%s/api/v1/agents/%s/deregister", m.config.APIURL, m.agentID)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+m.config.Token)

	resp, err := m.client.Do(req)
	if err != nil {
		return err
	}
	resp.Body.Close()

	return nil
}

// Start starts the heartbeat loop
func (m *Manager) Start(ctx context.Context) {
	ticker := time.NewTicker(m.config.Interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := m.sendHeartbeat(ctx); err != nil {
				fmt.Printf("Warning: heartbeat failed: %v\n", err)
			}
		}
	}
}

// HeartbeatPayload holds heartbeat data
type HeartbeatPayload struct {
	AgentID       string  `json:"agent_id"`
	Timestamp     int64   `json:"timestamp"`
	Status        string  `json:"status"`
	RunningJobs   int     `json:"running_jobs"`
	CPUUsage      float64 `json:"cpu_usage"`
	MemoryUsage   float64 `json:"memory_usage"`
	UptimeSeconds int64   `json:"uptime_seconds"`
}

var startTime = time.Now()

func (m *Manager) sendHeartbeat(ctx context.Context) error {
	payload := HeartbeatPayload{
		AgentID:       m.agentID,
		Timestamp:     time.Now().Unix(),
		Status:        "healthy",
		RunningJobs:   0, // Would be updated by executor
		CPUUsage:      m.getCPUUsage(),
		MemoryUsage:   m.getMemoryUsage(),
		UptimeSeconds: int64(time.Since(startTime).Seconds()),
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s/api/v1/agents/%s/heartbeat", m.config.APIURL, m.agentID)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+m.config.Token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := m.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("heartbeat failed: %s", string(respBody))
	}

	return nil
}

func (m *Manager) parseTags() []string {
	if m.config.Tags == "" {
		return []string{}
	}
	tags := strings.Split(m.config.Tags, ",")
	for i, tag := range tags {
		tags[i] = strings.TrimSpace(tag)
	}
	return tags
}

func (m *Manager) getCPUUsage() float64 {
	// Simplified CPU usage - in production would use proper system metrics
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	return float64(runtime.NumGoroutine()) / float64(runtime.NumCPU()) * 10
}

func (m *Manager) getMemoryUsage() float64 {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	return float64(memStats.Alloc) / float64(memStats.Sys) * 100
}

// AgentID returns the agent ID
func (m *Manager) AgentID() string {
	return m.agentID
}
