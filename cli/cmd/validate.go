package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/georgi-georgiev/testmesh/internal/storage/models"
)

// validateCmd represents the validate command
var validateCmd = &cobra.Command{
	Use:   "validate <flow.yaml>",
	Short: "Validate flow YAML syntax",
	Long: `Validate the syntax and structure of a flow YAML file without executing it.

This command checks:
- YAML syntax is valid
- Required fields are present
- Flow structure is correct
- Action types are recognized

Example:
  testmesh validate examples/simple-http.yaml
  testmesh validate my-flow.yaml`,
	Args: cobra.ExactArgs(1),
	RunE: validateFlow,
}

func init() {
	rootCmd.AddCommand(validateCmd)
}

func validateFlow(cmd *cobra.Command, args []string) error {
	flowFile := args[0]

	fmt.Printf("🔍 Validating: %s\n", flowFile)
	fmt.Println()

	// Read flow file
	data, err := os.ReadFile(flowFile)
	if err != nil {
		return fmt.Errorf("failed to read flow file: %w", err)
	}

	// Parse flow definition
	var flowWrapper struct {
		Flow models.FlowDefinition `yaml:"flow"`
	}
	if err := yaml.Unmarshal(data, &flowWrapper); err != nil {
		fmt.Printf("❌ YAML parsing failed\n")
		fmt.Printf("   %v\n", err)
		return fmt.Errorf("invalid YAML syntax")
	}

	definition := flowWrapper.Flow

	// Validate required fields
	errors := []string{}

	if definition.Name == "" {
		errors = append(errors, "flow.name is required")
	}

	if len(definition.Steps) == 0 {
		errors = append(errors, "flow.steps must have at least one step")
	}

	// Validate each step
	validActions := map[string]bool{
		"http_request":    true,
		"database_query":  true,
		"log":             true,
		"delay":           true,
		"transform":       true,
		"assert":          true,
		"condition":       true,
		"for_each":        true,
	}

	for i, step := range definition.Steps {
		if step.Action == "" {
			errors = append(errors, fmt.Sprintf("step %d: action is required", i+1))
		} else if !validActions[step.Action] {
			errors = append(errors, fmt.Sprintf("step %d: unknown action type '%s'", i+1, step.Action))
		}

		if step.Name == "" && step.ID == "" {
			errors = append(errors, fmt.Sprintf("step %d: either name or id is required", i+1))
		}
	}

	// Report results
	if len(errors) > 0 {
		fmt.Printf("❌ Validation failed with %d error(s):\n", len(errors))
		for _, err := range errors {
			fmt.Printf("   • %s\n", err)
		}
		fmt.Println()
		return fmt.Errorf("validation failed")
	}

	// Print success summary
	fmt.Printf("✅ Flow is valid\n")
	fmt.Printf("   Name: %s\n", definition.Name)
	if definition.Description != "" {
		fmt.Printf("   Description: %s\n", definition.Description)
	}
	if definition.Suite != "" {
		fmt.Printf("   Suite: %s\n", definition.Suite)
	}

	totalSteps := len(definition.Setup) + len(definition.Steps) + len(definition.Teardown)
	fmt.Printf("   Total steps: %d", totalSteps)
	if len(definition.Setup) > 0 {
		fmt.Printf(" (%d setup", len(definition.Setup))
	}
	if len(definition.Steps) > 0 {
		if len(definition.Setup) > 0 {
			fmt.Printf(", %d main", len(definition.Steps))
		} else {
			fmt.Printf(" (%d main", len(definition.Steps))
		}
	}
	if len(definition.Teardown) > 0 {
		fmt.Printf(", %d teardown)", len(definition.Teardown))
	} else if len(definition.Setup) > 0 || len(definition.Steps) > 0 {
		fmt.Printf(")")
	}
	fmt.Println()

	// Show step summary
	if len(definition.Steps) > 0 {
		fmt.Println()
		fmt.Println("   Steps:")
		for i, step := range definition.Steps {
			stepName := step.Name
			if stepName == "" {
				stepName = step.ID
			}
			if stepName == "" {
				stepName = fmt.Sprintf("step_%d", i+1)
			}
			fmt.Printf("   %d. %s (%s)\n", i+1, stepName, step.Action)
		}
	}

	fmt.Println()
	return nil
}
