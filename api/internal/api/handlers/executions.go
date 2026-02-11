package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/georgi-georgiev/testmesh/internal/runner"
	"github.com/georgi-georgiev/testmesh/internal/storage/models"
	"github.com/georgi-georgiev/testmesh/internal/storage/repository"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// ExecutionHandler handles execution-related requests
type ExecutionHandler struct {
	execRepo *repository.ExecutionRepository
	flowRepo *repository.FlowRepository
	logger   *zap.Logger
	wsHub    runner.WSHub
}

// NewExecutionHandler creates a new execution handler
func NewExecutionHandler(execRepo *repository.ExecutionRepository, flowRepo *repository.FlowRepository, logger *zap.Logger, wsHub runner.WSHub) *ExecutionHandler {
	return &ExecutionHandler{
		execRepo: execRepo,
		flowRepo: flowRepo,
		logger:   logger,
		wsHub:    wsHub,
	}
}

// Create handles POST /api/v1/executions
func (h *ExecutionHandler) Create(c *gin.Context) {
	var req struct {
		FlowID      string            `json:"flow_id" binding:"required"`
		Environment string            `json:"environment"`
		Variables   map[string]string `json:"variables"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	flowID, err := uuid.Parse(req.FlowID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid flow ID"})
		return
	}

	// Get flow
	flow, err := h.flowRepo.GetByID(flowID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "flow not found"})
		return
	}

	// Create execution record
	execution := &models.Execution{
		FlowID:      flowID,
		Status:      models.ExecutionStatusPending,
		Environment: req.Environment,
	}

	if err := h.execRepo.Create(execution); err != nil {
		h.logger.Error("Failed to create execution", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create execution"})
		return
	}

	// Start execution in background (we'll implement this properly later)
	go h.executeFlow(execution, flow, req.Variables)

	c.JSON(http.StatusCreated, execution)
}

// executeFlow runs the flow execution (placeholder for now)
func (h *ExecutionHandler) executeFlow(execution *models.Execution, flow *models.Flow, variables map[string]string) {
	// Update status to running
	execution.Status = models.ExecutionStatusRunning
	now := time.Now()
	execution.StartedAt = &now
	h.execRepo.Update(execution)

	// Execute flow using the runner
	executor := runner.NewExecutor(h.execRepo, h.logger, h.wsHub)
	err := executor.Execute(execution, &flow.Definition, variables)

	// Update execution status
	finishedAt := time.Now()
	execution.FinishedAt = &finishedAt
	execution.DurationMs = finishedAt.Sub(*execution.StartedAt).Milliseconds()

	if err != nil {
		execution.Status = models.ExecutionStatusFailed
		execution.Error = err.Error()

		// Broadcast execution failed
		if h.wsHub != nil {
			h.wsHub.BroadcastExecutionFailed(execution.ID, map[string]interface{}{
				"error":       execution.Error,
				"duration_ms": execution.DurationMs,
			})
		}
	} else {
		execution.Status = models.ExecutionStatusCompleted

		// Broadcast execution completed
		if h.wsHub != nil {
			h.wsHub.BroadcastExecutionCompleted(execution.ID, map[string]interface{}{
				"passed_steps": execution.PassedSteps,
				"failed_steps": execution.FailedSteps,
				"total_steps":  execution.TotalSteps,
				"duration_ms":  execution.DurationMs,
			})
		}
	}

	h.execRepo.Update(execution)
}

// List handles GET /api/v1/executions
func (h *ExecutionHandler) List(c *gin.Context) {
	var flowID *uuid.UUID
	if flowIDStr := c.Query("flow_id"); flowIDStr != "" {
		id, err := uuid.Parse(flowIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid flow ID"})
			return
		}
		flowID = &id
	}

	status := models.ExecutionStatus(c.Query("status"))
	limit := 20
	offset := 0

	executions, total, err := h.execRepo.List(flowID, status, limit, offset)
	if err != nil {
		h.logger.Error("Failed to list executions", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list executions"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"executions": executions,
		"total":      total,
		"limit":      limit,
		"offset":     offset,
	})
}

// Get handles GET /api/v1/executions/:id
func (h *ExecutionHandler) Get(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid execution ID"})
		return
	}

	execution, err := h.execRepo.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "execution not found"})
		return
	}

	c.JSON(http.StatusOK, execution)
}

// Cancel handles POST /api/v1/executions/:id/cancel
func (h *ExecutionHandler) Cancel(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid execution ID"})
		return
	}

	execution, err := h.execRepo.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "execution not found"})
		return
	}

	if execution.Status != models.ExecutionStatusRunning && execution.Status != models.ExecutionStatusPending {
		c.JSON(http.StatusBadRequest, gin.H{"error": "execution cannot be cancelled"})
		return
	}

	execution.Status = models.ExecutionStatusCancelled
	if err := h.execRepo.Update(execution); err != nil {
		h.logger.Error("Failed to cancel execution", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to cancel execution"})
		return
	}

	c.JSON(http.StatusOK, execution)
}

// GetLogs handles GET /api/v1/executions/:id/logs
func (h *ExecutionHandler) GetLogs(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid execution ID"})
		return
	}

	steps, err := h.execRepo.GetSteps(id)
	if err != nil {
		h.logger.Error("Failed to get execution logs", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get logs"})
		return
	}

	// Format logs from steps
	var logs []string
	for _, step := range steps {
		logs = append(logs, step.StepName+": "+string(step.Status))
	}

	c.JSON(http.StatusOK, gin.H{"logs": logs})
}

// GetSteps handles GET /api/v1/executions/:id/steps
func (h *ExecutionHandler) GetSteps(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid execution ID"})
		return
	}

	steps, err := h.execRepo.GetSteps(id)
	if err != nil {
		h.logger.Error("Failed to get execution steps", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get steps"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"steps": steps})
}
