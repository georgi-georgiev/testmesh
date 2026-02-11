package runner

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Interpolator handles variable interpolation in strings
type Interpolator struct {
	context *Context
}

// NewInterpolator creates a new interpolator
func NewInterpolator(ctx *Context) *Interpolator {
	return &Interpolator{context: ctx}
}

// Interpolate replaces all variables in a string
// Supports:
//   - ${VAR} - environment/context variables
//   - ${RANDOM_ID} - random UUID v4
//   - ${UUID} - random UUID v4 (alias)
//   - ${TIMESTAMP} - Unix timestamp
//   - ${ISO_TIMESTAMP} - ISO 8601 timestamp
//   - ${DATE} - Current date (YYYY-MM-DD)
//   - ${TIME} - Current time (HH:MM:SS)
//   - ${DATETIME} - Current datetime (YYYY-MM-DD HH:MM:SS)
//   - ${step_id.output_key} - step output reference
func (i *Interpolator) Interpolate(input string) string {
	if !strings.Contains(input, "${") {
		return input
	}

	result := input

	// Replace built-in functions
	result = i.replaceBuiltInVariables(result)

	// Replace step outputs
	result = i.replaceStepOutputs(result)

	// Replace context variables
	result = i.replaceContextVariables(result)

	return result
}

// replaceBuiltInVariables replaces built-in variable functions
func (i *Interpolator) replaceBuiltInVariables(input string) string {
	now := time.Now()

	replacements := map[string]string{
		"${RANDOM_ID}":      uuid.New().String(),
		"${UUID}":           uuid.New().String(),
		"${TIMESTAMP}":      fmt.Sprintf("%d", now.Unix()),
		"${ISO_TIMESTAMP}":  now.Format(time.RFC3339),
		"${DATE}":           now.Format("2006-01-02"),
		"${TIME}":           now.Format("15:04:05"),
		"${DATETIME}":       now.Format("2006-01-02 15:04:05"),
		"${YEAR}":           fmt.Sprintf("%d", now.Year()),
		"${MONTH}":          fmt.Sprintf("%02d", now.Month()),
		"${DAY}":            fmt.Sprintf("%02d", now.Day()),
		"${HOUR}":           fmt.Sprintf("%02d", now.Hour()),
		"${MINUTE}":         fmt.Sprintf("%02d", now.Minute()),
		"${SECOND}":         fmt.Sprintf("%02d", now.Second()),
	}

	result := input
	for placeholder, value := range replacements {
		result = strings.ReplaceAll(result, placeholder, value)
	}

	return result
}

// replaceStepOutputs replaces step output references
// Format: ${step_id.output_key} or ${step_id.nested.path}
func (i *Interpolator) replaceStepOutputs(input string) string {
	// Pattern to match ${step.output} or ${step.nested.path}
	pattern := regexp.MustCompile(`\$\{([a-zA-Z_][a-zA-Z0-9_]*)\.([\w.]+)\}`)

	result := pattern.ReplaceAllStringFunc(input, func(match string) string {
		// Extract step_id and path from ${step_id.path}
		parts := pattern.FindStringSubmatch(match)
		if len(parts) != 3 {
			return match // Return original if pattern doesn't match
		}

		stepID := parts[1]
		path := parts[2]

		// Get step output
		if value, exists := i.context.GetStepOutput(stepID, path); exists {
			return fmt.Sprintf("%v", value)
		}

		// If not found, try nested path
		pathParts := strings.Split(path, ".")
		if len(pathParts) > 1 {
			// Try getting the first part
			if value, exists := i.context.GetStepOutput(stepID, pathParts[0]); exists {
				// Try to navigate nested structure
				if nestedValue := i.navigateNestedValue(value, pathParts[1:]); nestedValue != nil {
					return fmt.Sprintf("%v", nestedValue)
				}
			}
		}

		// Return original if not found
		return match
	})

	return result
}

// navigateNestedValue navigates through nested map/slice structures
func (i *Interpolator) navigateNestedValue(value interface{}, path []string) interface{} {
	current := value

	for _, key := range path {
		switch v := current.(type) {
		case map[string]interface{}:
			current = v[key]
		case map[interface{}]interface{}:
			current = v[key]
		default:
			return nil
		}

		if current == nil {
			return nil
		}
	}

	return current
}

// replaceContextVariables replaces context/environment variables
// Format: ${VAR_NAME}
func (i *Interpolator) replaceContextVariables(input string) string {
	// Pattern to match ${VAR_NAME}
	pattern := regexp.MustCompile(`\$\{([A-Z_][A-Z0-9_]*)\}`)

	result := pattern.ReplaceAllStringFunc(input, func(match string) string {
		// Extract variable name from ${VAR}
		parts := pattern.FindStringSubmatch(match)
		if len(parts) != 2 {
			return match
		}

		varName := parts[1]

		// Get from context
		if value, exists := i.context.Get(varName); exists {
			return value
		}

		// Return original if not found
		return match
	})

	return result
}

// InterpolateMap recursively interpolates all string values in a map
func (i *Interpolator) InterpolateMap(input map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})

	for key, value := range input {
		result[key] = i.InterpolateValue(value)
	}

	return result
}

// InterpolateValue interpolates a single value (handles strings, maps, slices)
func (i *Interpolator) InterpolateValue(value interface{}) interface{} {
	switch v := value.(type) {
	case string:
		return i.Interpolate(v)
	case map[string]interface{}:
		return i.InterpolateMap(v)
	case []interface{}:
		result := make([]interface{}, len(v))
		for idx, item := range v {
			result[idx] = i.InterpolateValue(item)
		}
		return result
	default:
		return value
	}
}
