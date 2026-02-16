package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/georgi-georgiev/testmesh-agent/internal/executor"
	"github.com/georgi-georgiev/testmesh-agent/internal/heartbeat"
)

var (
	Version   = "dev"
	BuildTime = "unknown"
)

func main() {
	// Parse flags
	apiURL := flag.String("api-url", getEnv("TESTMESH_API_URL", "http://localhost:5016"), "TestMesh API URL")
	agentID := flag.String("agent-id", getEnv("TESTMESH_AGENT_ID", ""), "Agent ID (auto-generated if empty)")
	token := flag.String("token", getEnv("TESTMESH_AGENT_TOKEN", ""), "Agent authentication token")
	tags := flag.String("tags", getEnv("TESTMESH_AGENT_TAGS", ""), "Comma-separated agent tags")
	maxConcurrent := flag.Int("max-concurrent", 5, "Maximum concurrent executions")
	heartbeatInterval := flag.Duration("heartbeat-interval", 30*time.Second, "Heartbeat interval")
	version := flag.Bool("version", false, "Print version and exit")
	flag.Parse()

	if *version {
		fmt.Printf("TestMesh Agent %s (built %s)\n", Version, BuildTime)
		os.Exit(0)
	}

	// Create context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigChan
		fmt.Printf("\nReceived signal %v, shutting down...\n", sig)
		cancel()
	}()

	// Create agent configuration
	config := &AgentConfig{
		APIURL:            *apiURL,
		AgentID:           *agentID,
		Token:             *token,
		Tags:              *tags,
		MaxConcurrent:     *maxConcurrent,
		HeartbeatInterval: *heartbeatInterval,
	}

	// Run agent
	if err := runAgent(ctx, config); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

// AgentConfig holds agent configuration
type AgentConfig struct {
	APIURL            string
	AgentID           string
	Token             string
	Tags              string
	MaxConcurrent     int
	HeartbeatInterval time.Duration
}

func runAgent(ctx context.Context, config *AgentConfig) error {
	fmt.Println("🚀 TestMesh Agent starting...")
	fmt.Printf("   Version: %s\n", Version)
	fmt.Printf("   API URL: %s\n", config.APIURL)
	fmt.Printf("   Agent ID: %s\n", config.AgentID)
	fmt.Printf("   Max Concurrent: %d\n", config.MaxConcurrent)
	fmt.Println()

	// Create executor
	exec, err := executor.New(&executor.Config{
		APIURL:        config.APIURL,
		AgentID:       config.AgentID,
		Token:         config.Token,
		MaxConcurrent: config.MaxConcurrent,
	})
	if err != nil {
		return fmt.Errorf("failed to create executor: %w", err)
	}

	// Create heartbeat manager
	hb := heartbeat.New(&heartbeat.Config{
		APIURL:   config.APIURL,
		AgentID:  config.AgentID,
		Token:    config.Token,
		Tags:     config.Tags,
		Interval: config.HeartbeatInterval,
	})

	// Register agent
	if err := hb.Register(ctx); err != nil {
		return fmt.Errorf("failed to register agent: %w", err)
	}
	fmt.Println("✅ Agent registered successfully")

	// Start heartbeat
	go hb.Start(ctx)

	// Start executor
	go exec.Start(ctx)

	// Wait for shutdown
	<-ctx.Done()

	fmt.Println("👋 Agent shutting down...")

	// Deregister agent
	if err := hb.Deregister(context.Background()); err != nil {
		fmt.Printf("Warning: failed to deregister agent: %v\n", err)
	}

	return nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
