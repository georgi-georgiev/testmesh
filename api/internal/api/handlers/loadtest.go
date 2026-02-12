package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/georgi-georgiev/testmesh/internal/loadtest"
	"github.com/georgi-georgiev/testmesh/internal/storage/models"
	"github.com/georgi-georgiev/testmesh/internal/storage/repository"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// LoadTestHandler handles load test requests
type LoadTestHandler struct {
	loadTester *loadtest.LoadTester
	flowRepo   *repository.FlowRepository
	logger     *zap.Logger
	// Store running tests
	runningTests map[uuid.UUID]*loadtest.LoadTestResult
}

// NewLoadTestHandler creates a new load test handler
func NewLoadTestHandler(loadTester *loadtest.LoadTester, flowRepo *repository.FlowRepository, logger *zap.Logger) *LoadTestHandler {
	return &LoadTestHandler{
		loadTester:   loadTester,
		flowRepo:     flowRepo,
		logger:       logger,
		runningTests: make(map[uuid.UUID]*loadtest.LoadTestResult),
	}
}

// StartLoadTestRequest represents a request to start a load test
type StartLoadTestRequest struct {
	FlowIDs       []string          `json:"flow_ids" binding:"required"`
	VirtualUsers  int               `json:"virtual_users" binding:"required,min=1,max=1000"`
	DurationSec   int               `json:"duration_sec" binding:"required,min=1,max=3600"`
	RampUpSec     int               `json:"ramp_up_sec"`
	RampDownSec   int               `json:"ramp_down_sec"`
	ThinkTimeMs   int               `json:"think_time_ms"`
	Variables     map[string]string `json:"variables"`
	Environment   string            `json:"environment"`
}

// Start handles POST /api/v1/load-tests
func (h *LoadTestHandler) Start(c *gin.Context) {
	var req StartLoadTestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Parse flow IDs
	flowIDs := make([]uuid.UUID, 0, len(req.FlowIDs))
	for _, idStr := range req.FlowIDs {
		id, err := uuid.Parse(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid flow ID: " + idStr})
			return
		}
		flowIDs = append(flowIDs, id)
	}

	// Fetch flows
	var flows []*models.Flow
	for _, flowID := range flowIDs {
		flow, err := h.flowRepo.GetByID(flowID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "flow not found: " + flowID.String()})
			return
		}
		flows = append(flows, flow)
	}

	// Build config
	config := &loadtest.LoadTestConfig{
		FlowIDs:      flowIDs,
		VirtualUsers: req.VirtualUsers,
		Duration:     time.Duration(req.DurationSec) * time.Second,
		RampUpTime:   time.Duration(req.RampUpSec) * time.Second,
		RampDownTime: time.Duration(req.RampDownSec) * time.Second,
		ThinkTime:    time.Duration(req.ThinkTimeMs) * time.Millisecond,
		Variables:    req.Variables,
		Environment:  req.Environment,
	}

	// Default ramp up time
	if config.RampUpTime == 0 {
		config.RampUpTime = 10 * time.Second
	}

	// Start load test in background
	go func() {
		ctx := c.Request.Context()
		result, err := h.loadTester.Run(ctx, config, flows, func(progress *loadtest.LoadTestResult) {
			h.runningTests[progress.ID] = progress
		})
		if err != nil {
			h.logger.Error("Load test failed", zap.Error(err))
		}
		h.runningTests[result.ID] = result
	}()

	// Return immediately with test ID
	testID := uuid.New()
	c.JSON(http.StatusAccepted, gin.H{
		"id":      testID,
		"status":  "starting",
		"message": "Load test started. Use GET /api/v1/load-tests/{id} to check status.",
	})
}

// Get handles GET /api/v1/load-tests/:id
func (h *LoadTestHandler) Get(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid test ID"})
		return
	}

	result, ok := h.runningTests[id]
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "load test not found"})
		return
	}

	c.JSON(http.StatusOK, result)
}

// List handles GET /api/v1/load-tests
func (h *LoadTestHandler) List(c *gin.Context) {
	tests := make([]*loadtest.LoadTestResult, 0)
	for _, test := range h.runningTests {
		tests = append(tests, test)
	}

	c.JSON(http.StatusOK, gin.H{
		"tests": tests,
		"total": len(tests),
	})
}

// Stop handles POST /api/v1/load-tests/:id/stop
func (h *LoadTestHandler) Stop(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid test ID"})
		return
	}

	result, ok := h.runningTests[id]
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "load test not found"})
		return
	}

	// Mark as cancelled (actual cancellation would need context management)
	result.Status = "cancelled"
	finishedAt := time.Now()
	result.FinishedAt = &finishedAt

	c.JSON(http.StatusOK, gin.H{
		"message": "Load test stop requested",
		"status":  result.Status,
	})
}

// GetMetrics handles GET /api/v1/load-tests/:id/metrics
func (h *LoadTestHandler) GetMetrics(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid test ID"})
		return
	}

	result, ok := h.runningTests[id]
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "load test not found"})
		return
	}

	c.JSON(http.StatusOK, result.Metrics)
}

// GetTimeline handles GET /api/v1/load-tests/:id/timeline
func (h *LoadTestHandler) GetTimeline(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid test ID"})
		return
	}

	result, ok := h.runningTests[id]
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "load test not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"timeline": result.Timeline,
	})
}
