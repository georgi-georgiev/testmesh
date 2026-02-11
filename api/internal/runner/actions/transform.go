package actions

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/georgi-georgiev/testmesh/internal/storage/models"
	"github.com/tidwall/gjson"
	"go.uber.org/zap"
)

// TransformHandler handles data transformation actions
type TransformHandler struct {
	logger *zap.Logger
}

// NewTransformHandler creates a new transform handler
func NewTransformHandler(logger *zap.Logger) *TransformHandler {
	return &TransformHandler{
		logger: logger,
	}
}

// Execute transforms data using JSONPath extraction and manipulation
func (h *TransformHandler) Execute(ctx context.Context, config map[string]interface{}) (models.OutputData, error) {
	// Get input data
	input, ok := config["input"]
	if !ok {
		return nil, fmt.Errorf("input is required")
	}

	// Get transformations map
	transforms, ok := config["transforms"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("transforms is required and must be a map")
	}

	// Convert input to JSON for JSONPath processing
	inputJSON, err := json.Marshal(input)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal input: %w", err)
	}

	result := make(models.OutputData)

	// Apply each transformation
	for key, pathOrValue := range transforms {
		switch v := pathOrValue.(type) {
		case string:
			// If it starts with $, treat as JSONPath
			if len(v) > 0 && v[0] == '$' {
				extracted := gjson.GetBytes(inputJSON, v[1:]) // Remove $ prefix
				if !extracted.Exists() {
					h.logger.Warn("JSONPath did not match", zap.String("path", v), zap.String("key", key))
					result[key] = nil
					continue
				}
				result[key] = extracted.Value()
			} else {
				// Static string value
				result[key] = v
			}
		default:
			// Static value
			result[key] = v
		}
	}

	h.logger.Info("Data transformed", zap.Int("fields", len(result)))

	return result, nil
}
