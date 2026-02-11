package actions

import (
	"context"
	"fmt"

	"github.com/georgi-georgiev/testmesh/internal/runner/assertions"
	"github.com/georgi-georgiev/testmesh/internal/storage/models"
	"go.uber.org/zap"
)

// AssertHandler handles standalone assertion actions
type AssertHandler struct {
	logger *zap.Logger
}

// NewAssertHandler creates a new assert handler
func NewAssertHandler(logger *zap.Logger) *AssertHandler {
	return &AssertHandler{
		logger: logger,
	}
}

// Execute runs assertions against provided data
func (h *AssertHandler) Execute(ctx context.Context, config map[string]interface{}) (models.OutputData, error) {
	// Get data to assert against
	data, ok := config["data"]
	if !ok {
		return nil, fmt.Errorf("data is required")
	}

	// Convert to OutputData
	var outputData models.OutputData
	switch v := data.(type) {
	case models.OutputData:
		outputData = v
	case map[string]interface{}:
		outputData = models.OutputData(v)
	default:
		// Wrap in a map
		outputData = models.OutputData{"value": data}
	}

	h.logger.Debug("Assert action data",
		zap.Any("data", outputData),
		zap.Int("fields", len(outputData)),
	)

	// Get assertions
	assertionsRaw, ok := config["assertions"]
	if !ok {
		return nil, fmt.Errorf("assertions is required")
	}

	var assertionList []string
	switch v := assertionsRaw.(type) {
	case []string:
		assertionList = v
	case []interface{}:
		for _, a := range v {
			if s, ok := a.(string); ok {
				assertionList = append(assertionList, s)
			}
		}
	default:
		return nil, fmt.Errorf("assertions must be an array of strings")
	}

	// Run assertions
	evaluator := assertions.NewEvaluator(outputData)
	if err := evaluator.Evaluate(assertionList); err != nil {
		return nil, fmt.Errorf("assertion failed: %w", err)
	}

	h.logger.Info("All assertions passed", zap.Int("count", len(assertionList)))

	return models.OutputData{
		"assertions_count": len(assertionList),
		"passed":           true,
		"data":             outputData,
	}, nil
}
