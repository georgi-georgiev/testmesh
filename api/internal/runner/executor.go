package runner

import (
	"context"
	"fmt"
	"time"

	"github.com/georgi-georgiev/testmesh/internal/runner/actions"
	"github.com/georgi-georgiev/testmesh/internal/runner/assertions"
	"github.com/georgi-georgiev/testmesh/internal/storage/models"
	"github.com/georgi-georgiev/testmesh/internal/storage/repository"
	"go.uber.org/zap"
)

// Executor orchestrates flow execution
type Executor struct {
	repo   *repository.ExecutionRepository
	logger *zap.Logger
}

// NewExecutor creates a new executor instance
func NewExecutor(repo *repository.ExecutionRepository, logger *zap.Logger) *Executor {
	return &Executor{
		repo:   repo,
		logger: logger,
	}
}

// Execute runs a flow definition
func (e *Executor) Execute(execution *models.Execution, definition *models.FlowDefinition, variables map[string]string) error {
	ctx := context.Background()

	// Create execution context
	execCtx := NewContext(variables, definition.Env)

	// Count total steps
	totalSteps := len(definition.Setup) + len(definition.Steps) + len(definition.Teardown)
	execution.TotalSteps = totalSteps

	// Execute setup steps
	if len(definition.Setup) > 0 {
		e.logger.Info("Executing setup steps", zap.Int("count", len(definition.Setup)))
		if err := e.executeSteps(ctx, execution, definition.Setup, execCtx, "setup"); err != nil {
			return fmt.Errorf("setup failed: %w", err)
		}
	}

	// Execute main steps
	e.logger.Info("Executing main steps", zap.Int("count", len(definition.Steps)))
	if err := e.executeSteps(ctx, execution, definition.Steps, execCtx, "main"); err != nil {
		// Run teardown even if main steps fail
		if len(definition.Teardown) > 0 {
			e.logger.Info("Executing teardown steps after failure")
			e.executeSteps(ctx, execution, definition.Teardown, execCtx, "teardown")
		}
		return fmt.Errorf("execution failed: %w", err)
	}

	// Execute teardown steps
	if len(definition.Teardown) > 0 {
		e.logger.Info("Executing teardown steps", zap.Int("count", len(definition.Teardown)))
		if err := e.executeSteps(ctx, execution, definition.Teardown, execCtx, "teardown"); err != nil {
			return fmt.Errorf("teardown failed: %w", err)
		}
	}

	return nil
}

// executeSteps executes a slice of steps
func (e *Executor) executeSteps(ctx context.Context, execution *models.Execution, steps []models.Step, execCtx *Context, phase string) error {
	for i, step := range steps {
		stepID := step.ID
		if stepID == "" {
			stepID = fmt.Sprintf("%s_%d", phase, i)
		}

		e.logger.Info("Executing step",
			zap.String("step_id", stepID),
			zap.String("action", step.Action),
			zap.String("phase", phase),
		)

		// Create step record
		execStep := &models.ExecutionStep{
			ExecutionID: execution.ID,
			StepID:      stepID,
			StepName:    step.Name,
			Action:      step.Action,
			Status:      models.StepStatusRunning,
		}
		now := time.Now()
		execStep.StartedAt = &now

		if err := e.repo.CreateStep(execStep); err != nil {
			return err
		}

		// Execute the step with retry logic
		result, err := e.executeStepWithRetry(ctx, &step, execStep, execCtx)

		// Update step record
		finishedAt := time.Now()
		execStep.FinishedAt = &finishedAt
		execStep.DurationMs = finishedAt.Sub(*execStep.StartedAt).Milliseconds()

		if err != nil {
			execStep.Status = models.StepStatusFailed
			execStep.ErrorMessage = err.Error()
			e.repo.UpdateStep(execStep)

			execution.FailedSteps++

			// Wrap error with execution context
			execErr := NewExecutionError(
				phase,
				stepID,
				step.Name,
				step.Action,
				err.Error(),
				err,
			)
			return execErr
		}

		execStep.Status = models.StepStatusCompleted
		execStep.Output = result
		e.repo.UpdateStep(execStep)

		execution.PassedSteps++

		// Store step output in context
		if step.Output != nil {
			for key, path := range step.Output {
				value := extractValue(result, path)
				execCtx.SetStepOutput(stepID, key, value)
			}
		}
	}

	return nil
}

// executeStepWithRetry executes a step with retry logic
func (e *Executor) executeStepWithRetry(ctx context.Context, step *models.Step, execStep *models.ExecutionStep, execCtx *Context) (models.OutputData, error) {
	maxAttempts := 1
	var delay time.Duration
	backoff := "fixed"

	// Get retry configuration
	if step.Retry != nil {
		maxAttempts = step.Retry.MaxAttempts
		if maxAttempts < 1 {
			maxAttempts = 1
		}

		// Parse delay
		if step.Retry.Delay != "" {
			if d, err := time.ParseDuration(step.Retry.Delay); err == nil {
				delay = d
			}
		}

		if step.Retry.Backoff != "" {
			backoff = step.Retry.Backoff
		}
	}

	var lastErr error
	currentDelay := delay

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		execStep.Attempt = attempt
		e.repo.UpdateStep(execStep)

		if attempt > 1 {
			e.logger.Info("Retrying step",
				zap.String("step_id", step.ID),
				zap.Int("attempt", attempt),
				zap.Int("max_attempts", maxAttempts),
			)

			// Wait before retry
			if currentDelay > 0 {
				time.Sleep(currentDelay)

				// Apply backoff
				if backoff == "exponential" {
					currentDelay *= 2
				}
			}
		}

		result, err := e.executeStep(ctx, step, execCtx)
		if err == nil {
			return result, nil
		}

		lastErr = err

		// Log retry failure
		if attempt < maxAttempts {
			e.logger.Warn("Step execution failed, will retry",
				zap.String("step_id", step.ID),
				zap.Error(err),
				zap.Int("attempt", attempt),
			)
		}
	}

	return nil, fmt.Errorf("failed after %d attempts: %w", maxAttempts, lastErr)
}

// executeStep executes a single step
func (e *Executor) executeStep(ctx context.Context, step *models.Step, execCtx *Context) (models.OutputData, error) {
	// Create interpolator
	interpolator := NewInterpolator(execCtx)

	// Interpolate variables in config
	config := interpolator.InterpolateMap(step.Config)

	// Get action handler
	handler, err := e.getActionHandler(step.Action)
	if err != nil {
		return nil, err
	}

	// Execute action
	result, err := handler.Execute(ctx, config)
	if err != nil {
		return nil, err
	}

	// Run assertions if any
	if len(step.Assert) > 0 {
		evaluator := assertions.NewEvaluator(result)
		if err := evaluator.Evaluate(step.Assert); err != nil {
			return result, fmt.Errorf("assertion failed: %w", err)
		}
		e.logger.Info("All assertions passed", zap.Int("count", len(step.Assert)))
	}

	return result, nil
}

// getActionHandler returns the appropriate action handler
func (e *Executor) getActionHandler(actionType string) (actions.Handler, error) {
	switch actionType {
	case "http_request":
		return actions.NewHTTPHandler(e.logger), nil
	case "database_query":
		return actions.NewDatabaseHandler(e.logger), nil
	default:
		return nil, fmt.Errorf("unknown action type: %s", actionType)
	}
}

// extractValue extracts a value from result using JSONPath
func extractValue(result models.OutputData, path string) interface{} {
	if path == "" || path == "$" {
		return result
	}

	// Use assertion evaluator for JSONPath extraction
	evaluator := assertions.NewEvaluator(result)
	value, err := evaluator.EvaluateJSONPath(path)
	if err != nil {
		// Fallback to direct field access
		return result[path]
	}

	return value
}
