package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/georgi-georgiev/testmesh/internal/runner"
	"github.com/georgi-georgiev/testmesh/internal/shared/logger"
	"github.com/georgi-georgiev/testmesh/internal/storage/models"
	"github.com/google/uuid"
)

var (
	environment string
	verbose     bool
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run <flow.yaml>",
	Short: "Execute a flow locally",
	Long: `Execute a test flow from a YAML file locally without connecting to a server.

This command parses the flow definition and executes it immediately,
showing results in the terminal.

Example:
  testmesh run examples/simple-http.yaml
  testmesh run my-flow.yaml --env staging --verbose`,
	Args: cobra.ExactArgs(1),
	RunE: runFlow,
}

func init() {
	rootCmd.AddCommand(runCmd)

	runCmd.Flags().StringVarP(&environment, "env", "e", "development", "Environment to run the flow in")
	runCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Verbose output (show all logs)")
}

func runFlow(cmd *cobra.Command, args []string) error {
	flowFile := args[0]

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
		return fmt.Errorf("failed to parse flow YAML: %w", err)
	}

	definition := flowWrapper.Flow

	// Print flow info
	fmt.Printf("🚀 Running flow: %s\n", definition.Name)
	if definition.Description != "" {
		fmt.Printf("   %s\n", definition.Description)
	}
	fmt.Printf("   Environment: %s\n", environment)
	fmt.Println()

	// Initialize logger
	log := logger.New()
	// Note: Logger will output at configured level
	// Verbose flag could be used to adjust log level in future

	// Create in-memory execution record
	execution := &models.Execution{
		ID:          uuid.New(),
		Status:      models.ExecutionStatusRunning,
		Environment: environment,
	}
	now := time.Now()
	execution.StartedAt = &now

	// Create executor (without database and websocket for local execution)
	executor := runner.NewExecutor(nil, log, nil)

	// Execute flow
	startTime := time.Now()
	err = executor.Execute(execution, &definition, nil)
	duration := time.Since(startTime)

	fmt.Println()
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println()

	if err != nil {
		fmt.Printf("❌ Flow failed after %v\n", duration.Round(time.Millisecond))
		fmt.Printf("   Error: %v\n", err)
		fmt.Println()
		return fmt.Errorf("flow execution failed")
	}

	// Print success summary
	fmt.Printf("✅ Flow completed successfully in %v\n", duration.Round(time.Millisecond))
	fmt.Printf("   Total steps: %d\n", execution.TotalSteps)
	fmt.Printf("   Passed: %d\n", execution.PassedSteps)
	fmt.Printf("   Failed: %d\n", execution.FailedSteps)
	fmt.Println()

	return nil
}
