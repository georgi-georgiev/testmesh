package api

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/georgi-georgiev/testmesh/internal/ai"
	"github.com/georgi-georgiev/testmesh/internal/api/handlers"
	"github.com/georgi-georgiev/testmesh/internal/api/middleware"
	"github.com/georgi-georgiev/testmesh/internal/api/websocket"
	"github.com/georgi-georgiev/testmesh/internal/reporting"
	"github.com/georgi-georgiev/testmesh/internal/storage/repository"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// NewRouter creates and configures the API router
func NewRouter(db *gorm.DB, logger *zap.Logger, wsHub *websocket.Hub) *gin.Engine {
	// Set Gin mode
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()

	// Global middleware
	router.Use(middleware.Logger(logger))
	router.Use(middleware.Recovery(logger))
	router.Use(middleware.CORS())

	// Initialize repositories
	flowRepo := repository.NewFlowRepository(db)
	executionRepo := repository.NewExecutionRepository(db)
	mockRepo := repository.NewMockRepository(db)
	contractRepo := repository.NewContractRepository(db)
	reportingRepo := repository.NewReportingRepository(db)

	// Initialize reporting services
	reportOutputDir := filepath.Join(os.TempDir(), "testmesh", "reports")
	aggregator := reporting.NewAggregator(db, reportingRepo, executionRepo, flowRepo, logger)
	generator := reporting.NewGenerator(db, reportingRepo, executionRepo, flowRepo, logger, reportOutputDir)

	// Start scheduled aggregation
	if err := aggregator.ScheduleAggregation(); err != nil {
		logger.Error("Failed to schedule aggregation", zap.Error(err))
	}

	// Initialize AI services (using environment variables for API keys)
	aiConfig := ai.ProviderConfig{
		AnthropicAPIKey: os.Getenv("ANTHROPIC_API_KEY"),
		OpenAIAPIKey:    os.Getenv("OPENAI_API_KEY"),
		LocalEndpoint:   os.Getenv("LOCAL_LLM_ENDPOINT"),
	}
	aiProviders := ai.NewProviderManager(aiConfig, logger)
	aiGenerator := ai.NewGenerator(db, aiProviders, flowRepo, logger)
	aiAnalyzer := ai.NewAnalyzer(db, aiProviders, flowRepo, logger)
	aiSelfHealing := ai.NewSelfHealingEngine(db, aiProviders, flowRepo, executionRepo, logger)

	// Initialize handlers
	healthHandler := handlers.NewHealthHandler(db)
	flowHandler := handlers.NewFlowHandler(flowRepo, logger)
	executionHandler := handlers.NewExecutionHandler(executionRepo, flowRepo, mockRepo, contractRepo, logger, wsHub)
	mockHandler := handlers.NewMockHandler(mockRepo, logger)
	contractHandler := handlers.NewContractHandler(contractRepo, logger)
	reportingHandler := handlers.NewReportingHandler(reportingRepo, aggregator, generator, logger)
	aiHandler := handlers.NewAIHandler(db, aiGenerator, aiAnalyzer, aiSelfHealing, aiProviders, logger)
	wsHandler := websocket.NewHandler(wsHub, logger)

	// Health check
	router.GET("/health", healthHandler.Check)

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Flow routes
		flows := v1.Group("/flows")
		{
			flows.POST("", flowHandler.Create)
			flows.GET("", flowHandler.List)
			flows.GET("/:id", flowHandler.Get)
			flows.PUT("/:id", flowHandler.Update)
			flows.DELETE("/:id", flowHandler.Delete)
		}

		// Execution routes
		executions := v1.Group("/executions")
		{
			executions.POST("", executionHandler.Create)
			executions.GET("", executionHandler.List)
			executions.GET("/:id", executionHandler.Get)
			executions.POST("/:id/cancel", executionHandler.Cancel)
			executions.GET("/:id/logs", executionHandler.GetLogs)
			executions.GET("/:id/steps", executionHandler.GetSteps)
		}

		// Mock server routes
		mocks := v1.Group("/mock-servers")
		{
			mocks.GET("", mockHandler.ListServers)
			mocks.GET("/:id", mockHandler.GetServer)
			mocks.DELETE("/:id", mockHandler.DeleteServer)
			mocks.GET("/:id/endpoints", mockHandler.GetEndpoints)
			mocks.POST("/:id/endpoints", mockHandler.CreateEndpoint)
			mocks.PUT("/:id/endpoints/:endpoint_id", mockHandler.UpdateEndpoint)
			mocks.DELETE("/:id/endpoints/:endpoint_id", mockHandler.DeleteEndpoint)
			mocks.GET("/:id/requests", mockHandler.GetRequests)
			mocks.GET("/:id/state", mockHandler.GetStates)
			mocks.GET("/:id/state/:key", mockHandler.GetState)
		}

		// Contract testing routes
		contractsGroup := v1.Group("/contracts")
		{
			contractsGroup.GET("", contractHandler.ListContracts)
			contractsGroup.GET("/versions", contractHandler.GetContractVersions)
			contractsGroup.POST("/import", contractHandler.ImportPact)
			contractsGroup.POST("/breaking-changes", contractHandler.DetectBreakingChanges)
			contractsGroup.GET("/:id", contractHandler.GetContract)
			contractsGroup.DELETE("/:id", contractHandler.DeleteContract)
			contractsGroup.GET("/:id/pact", contractHandler.ExportPact)
			contractsGroup.GET("/:id/verifications", contractHandler.ListVerifications)
			contractsGroup.GET("/:id/breaking-changes", contractHandler.ListBreakingChanges)
		}

		// Verification routes
		verifications := v1.Group("/verifications")
		{
			verifications.GET("/:id", contractHandler.GetVerification)
		}

		// Report routes
		reports := v1.Group("/reports")
		{
			reports.POST("/generate", reportingHandler.GenerateReport)
			reports.GET("", reportingHandler.ListReports)
			reports.GET("/:id", reportingHandler.GetReport)
			reports.GET("/:id/download", reportingHandler.DownloadReport)
			reports.DELETE("/:id", reportingHandler.DeleteReport)
		}

		// Analytics routes
		analytics := v1.Group("/analytics")
		{
			analytics.GET("/metrics", reportingHandler.GetMetrics)
			analytics.GET("/flakiness", reportingHandler.GetFlakiness)
			analytics.GET("/trends", reportingHandler.GetTrends)
			analytics.GET("/steps", reportingHandler.GetStepPerformance)
			analytics.POST("/aggregate", reportingHandler.TriggerAggregation)
		}

		// AI routes
		aiRoutes := v1.Group("/ai")
		{
			aiRoutes.POST("/generate", aiHandler.Generate)
			aiRoutes.POST("/import/openapi", aiHandler.ImportOpenAPI)
			aiRoutes.POST("/import/postman", aiHandler.ImportPostman)
			aiRoutes.POST("/import/pact", aiHandler.ImportPact)
			aiRoutes.POST("/coverage/analyze", aiHandler.AnalyzeCoverage)
			aiRoutes.POST("/analyze/:execution_id", aiHandler.AnalyzeFailure)
			aiRoutes.GET("/suggestions", aiHandler.ListSuggestions)
			aiRoutes.GET("/suggestions/:id", aiHandler.GetSuggestion)
			aiRoutes.POST("/suggestions/:id/apply", aiHandler.ApplySuggestion)
			aiRoutes.POST("/suggestions/:id/accept", aiHandler.AcceptSuggestion)
			aiRoutes.POST("/suggestions/:id/reject", aiHandler.RejectSuggestion)
			aiRoutes.GET("/usage", aiHandler.GetUsage)
			aiRoutes.GET("/providers", aiHandler.GetProviders)
		}
	}

	// WebSocket routes
	router.GET("/ws/executions/:id", wsHandler.HandleConnection)

	// 404 handler
	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"error": "endpoint not found"})
	})

	return router
}
