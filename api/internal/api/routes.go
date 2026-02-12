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
	"github.com/georgi-georgiev/testmesh/internal/auth"
	"github.com/georgi-georgiev/testmesh/internal/loadtest"
	"github.com/georgi-georgiev/testmesh/internal/plugins"
	"github.com/georgi-georgiev/testmesh/internal/reporting"
	"github.com/georgi-georgiev/testmesh/internal/runner"
	"github.com/georgi-georgiev/testmesh/internal/scheduler"
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
	collectionRepo := repository.NewCollectionRepository(db)
	historyRepo := repository.NewHistoryRepository(db)

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
	aiRepo := repository.NewAIRepository(db)

	// Initialize services
	oauth2Service := auth.NewOAuth2Service(logger)

	// Initialize handlers
	healthHandler := handlers.NewHealthHandler(db)
	flowHandler := handlers.NewFlowHandler(flowRepo, logger)
	executionHandler := handlers.NewExecutionHandler(executionRepo, flowRepo, mockRepo, contractRepo, logger, wsHub)
	mockHandler := handlers.NewMockHandler(mockRepo, logger)
	contractHandler := handlers.NewContractHandler(contractRepo, logger)
	reportingHandler := handlers.NewReportingHandler(reportingRepo, aggregator, generator, logger)
	aiHandler := handlers.NewAIHandler(db, aiRepo, aiGenerator, aiAnalyzer, aiSelfHealing, aiProviders, logger)
	collectionHandler := handlers.NewCollectionHandler(collectionRepo, flowRepo, logger)
	oauth2Handler := handlers.NewOAuth2Handler(oauth2Service, logger)
	historyHandler := handlers.NewHistoryHandler(historyRepo, logger)
	wsHandler := websocket.NewHandler(wsHub, logger)

	// Initialize collection runner (executor created per-run to support parallel executions)
	executor := runner.NewExecutor(executionRepo, contractRepo, logger, wsHub, nil)
	collectionRunner := runner.NewCollectionRunner(executor, logger)
	runnerHandler := handlers.NewRunnerHandler(collectionRunner, flowRepo, logger)

	// Initialize workspace handler
	workspaceRepo := repository.NewWorkspaceRepository(db)
	workspaceHandler := handlers.NewWorkspaceHandler(workspaceRepo, logger)

	// Initialize bulk handler
	bulkHandler := handlers.NewBulkHandler(flowRepo, collectionRepo, logger)

	// Initialize import/export handler
	importExportHandler := handlers.NewImportExportHandler(flowRepo, logger)

	// Initialize load test handler
	loadTester := loadtest.NewLoadTester(logger)
	loadTestHandler := handlers.NewLoadTestHandler(loadTester, flowRepo, logger)

	// Initialize plugin registry
	pluginDir := filepath.Join(os.TempDir(), "testmesh", "plugins")
	pluginRegistry := plugins.NewRegistry(pluginDir, logger)
	pluginRegistry.Discover()
	pluginRegistry.LoadAll()
	pluginHandler := handlers.NewPluginHandler(pluginRegistry, logger)

	// Initialize scheduler
	scheduleRepo := repository.NewScheduleRepository(db)
	sched := scheduler.NewScheduler(scheduleRepo, logger)
	scheduleHandler := handlers.NewScheduleHandler(scheduleRepo, sched, logger)

	// Start the scheduler
	if err := sched.Start(); err != nil {
		logger.Error("Failed to start scheduler", zap.Error(err))
	}

	// Initialize collaboration handler
	collaborationRepo := repository.NewCollaborationRepository(db)
	collaborationHandler := handlers.NewCollaborationHandler(collaborationRepo, logger)

	// Health check
	router.GET("/health", healthHandler.Check)

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Collection routes
		collections := v1.Group("/collections")
		{
			collections.POST("", collectionHandler.Create)
			collections.GET("", collectionHandler.List)
			collections.GET("/tree", collectionHandler.GetTree)
			collections.GET("/search", collectionHandler.Search)
			collections.GET("/:id", collectionHandler.Get)
			collections.PUT("/:id", collectionHandler.Update)
			collections.DELETE("/:id", collectionHandler.Delete)
			collections.GET("/:id/children", collectionHandler.GetChildren)
			collections.GET("/:id/flows", collectionHandler.GetFlows)
			collections.POST("/:id/flows", collectionHandler.AddFlow)
			collections.DELETE("/:id/flows/:flow_id", collectionHandler.RemoveFlow)
			collections.GET("/:id/ancestors", collectionHandler.GetAncestors)
			collections.POST("/:id/move", collectionHandler.Move)
			collections.POST("/:id/duplicate", collectionHandler.Duplicate)
			collections.POST("/:id/reorder", collectionHandler.Reorder)
		}

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
			executions.GET("/:id/steps/:step_id", executionHandler.GetStep)
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
			mocks.GET("/:id/requests/:request_id", mockHandler.GetRequest)
			mocks.DELETE("/:id/requests", mockHandler.DeleteRequests)
			mocks.GET("/:id/state", mockHandler.GetStates)
			mocks.POST("/:id/state", mockHandler.CreateState)
			mocks.GET("/:id/state/:key", mockHandler.GetState)
			mocks.PUT("/:id/state/:key", mockHandler.UpdateState)
			mocks.DELETE("/:id/state/:key", mockHandler.DeleteState)
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
			contractsGroup.GET("/:id/interactions", contractHandler.ListInteractions)
			contractsGroup.GET("/:id/interactions/:interaction_id", contractHandler.GetInteraction)
			contractsGroup.DELETE("/:id/interactions/:interaction_id", contractHandler.DeleteInteraction)
		}

		// Verification routes
		verifications := v1.Group("/verifications")
		{
			verifications.POST("", contractHandler.CreateVerification)
			verifications.GET("/:id", contractHandler.GetVerification)
			verifications.PUT("/:id", contractHandler.UpdateVerification)
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
			aiRoutes.DELETE("/suggestions/:id", aiHandler.DeleteSuggestion)
			aiRoutes.GET("/usage", aiHandler.GetUsage)
			aiRoutes.GET("/providers", aiHandler.GetProviders)
			aiRoutes.GET("/generation-history", aiHandler.ListGenerationHistory)
			aiRoutes.GET("/generation-history/:id", aiHandler.GetGenerationHistory)
			aiRoutes.GET("/import-history", aiHandler.ListImportHistory)
			aiRoutes.GET("/import-history/:id", aiHandler.GetImportHistory)
			aiRoutes.GET("/coverage-analysis", aiHandler.ListCoverageAnalyses)
			aiRoutes.GET("/coverage-analysis/:id", aiHandler.GetCoverageAnalysis)
		}

		// OAuth2 routes
		oauth2 := v1.Group("/oauth2")
		{
			oauth2.GET("/providers", oauth2Handler.GetProviders)
			oauth2.GET("/providers/:name", oauth2Handler.GetProvider)
			oauth2.POST("/auth-url", oauth2Handler.GetAuthorizationURL)
			oauth2.POST("/token/code", oauth2Handler.ExchangeCode)
			oauth2.POST("/token/client-credentials", oauth2Handler.ClientCredentials)
			oauth2.POST("/token/password", oauth2Handler.PasswordGrant)
			oauth2.POST("/token/refresh", oauth2Handler.RefreshToken)
		}

		// Request history routes
		history := v1.Group("/history")
		{
			history.POST("", historyHandler.Create)
			history.GET("", historyHandler.List)
			history.GET("/stats", historyHandler.GetStats)
			history.GET("/:id", historyHandler.Get)
			history.DELETE("/:id", historyHandler.Delete)
			history.POST("/:id/save", historyHandler.Save)
			history.POST("/:id/unsave", historyHandler.Unsave)
			history.POST("/:id/tags", historyHandler.AddTag)
			history.DELETE("/:id/tags/:tag", historyHandler.RemoveTag)
			history.DELETE("", historyHandler.Clear)
		}

		// Collection runner routes (data-driven testing)
		runnerRoutes := v1.Group("/runner")
		{
			runnerRoutes.POST("/run", runnerHandler.Run)
			runnerRoutes.POST("/parse-data", runnerHandler.ParseData)
		}

		// Workspace routes
		workspaces := v1.Group("/workspaces")
		{
			workspaces.POST("", workspaceHandler.Create)
			workspaces.GET("", workspaceHandler.List)
			workspaces.GET("/personal", workspaceHandler.GetPersonal)
			workspaces.GET("/slug/:slug", workspaceHandler.GetBySlug)
			workspaces.GET("/:id", workspaceHandler.Get)
			workspaces.PUT("/:id", workspaceHandler.Update)
			workspaces.DELETE("/:id", workspaceHandler.Delete)
			workspaces.GET("/:id/role", workspaceHandler.GetUserRole)
			workspaces.GET("/:id/members", workspaceHandler.ListMembers)
			workspaces.POST("/:id/members", workspaceHandler.AddMember)
			workspaces.PUT("/:id/members/:user_id", workspaceHandler.UpdateMember)
			workspaces.DELETE("/:id/members/:user_id", workspaceHandler.RemoveMember)
			workspaces.GET("/:id/invitations", workspaceHandler.ListInvitations)
			workspaces.POST("/:id/invitations", workspaceHandler.InviteMember)
			workspaces.DELETE("/:id/invitations/:invitation_id", workspaceHandler.RevokeInvitation)
		}

		// Invitation acceptance (outside workspace context)
		v1.POST("/invitations/accept", workspaceHandler.AcceptInvitation)

		// Bulk operations routes
		bulk := v1.Group("/bulk/flows")
		{
			bulk.POST("/tags/add", bulkHandler.AddTags)
			bulk.POST("/tags/remove", bulkHandler.RemoveTags)
			bulk.POST("/move", bulkHandler.Move)
			bulk.POST("/delete", bulkHandler.Delete)
			bulk.POST("/duplicate", bulkHandler.Duplicate)
			bulk.POST("/export", bulkHandler.Export)
			bulk.POST("/find-replace", bulkHandler.FindReplace)
		}

		// Import/Export routes
		v1.POST("/import/parse", importExportHandler.Parse)
		v1.POST("/import", importExportHandler.Import)
		v1.POST("/export", importExportHandler.Export)
		v1.GET("/export/download", importExportHandler.ExportDownload)

		// Load testing routes
		loadTests := v1.Group("/load-tests")
		{
			loadTests.POST("", loadTestHandler.Start)
			loadTests.GET("", loadTestHandler.List)
			loadTests.GET("/:id", loadTestHandler.Get)
			loadTests.POST("/:id/stop", loadTestHandler.Stop)
			loadTests.GET("/:id/metrics", loadTestHandler.GetMetrics)
			loadTests.GET("/:id/timeline", loadTestHandler.GetTimeline)
		}

		// Plugin routes
		pluginsRoutes := v1.Group("/plugins")
		{
			pluginsRoutes.GET("", pluginHandler.List)
			pluginsRoutes.GET("/types", pluginHandler.GetTypes)
			pluginsRoutes.POST("/discover", pluginHandler.Discover)
			pluginsRoutes.POST("/install", pluginHandler.Install)
			pluginsRoutes.GET("/:id", pluginHandler.Get)
			pluginsRoutes.POST("/:id/enable", pluginHandler.Enable)
			pluginsRoutes.POST("/:id/disable", pluginHandler.Disable)
			pluginsRoutes.DELETE("/:id", pluginHandler.Uninstall)
		}

		// Schedule routes
		schedules := v1.Group("/schedules")
		{
			schedules.POST("", scheduleHandler.Create)
			schedules.GET("", scheduleHandler.List)
			schedules.GET("/presets", scheduleHandler.GetPresets)
			schedules.GET("/timezones", scheduleHandler.GetTimezones)
			schedules.POST("/validate-cron", scheduleHandler.ValidateCron)
			schedules.GET("/:id", scheduleHandler.Get)
			schedules.PUT("/:id", scheduleHandler.Update)
			schedules.DELETE("/:id", scheduleHandler.Delete)
			schedules.POST("/:id/pause", scheduleHandler.Pause)
			schedules.POST("/:id/resume", scheduleHandler.Resume)
			schedules.POST("/:id/trigger", scheduleHandler.Trigger)
			schedules.GET("/:id/runs", scheduleHandler.GetRuns)
			schedules.GET("/:id/stats", scheduleHandler.GetStats)
		}

		// Collaboration routes
		collaboration := v1.Group("/collaboration")
		{
			// Presence
			collaboration.POST("/presence", collaborationHandler.SetPresence)
			collaboration.DELETE("/presence", collaborationHandler.RemovePresence)
			collaboration.GET("/presence/:resource_type/:resource_id", collaborationHandler.GetPresence)

			// Comments
			collaboration.POST("/comments", collaborationHandler.CreateComment)
			collaboration.GET("/comments/:id", collaborationHandler.GetComment)
			collaboration.PUT("/comments/:id", collaborationHandler.UpdateComment)
			collaboration.DELETE("/comments/:id", collaborationHandler.DeleteComment)
			collaboration.POST("/comments/:id/resolve", collaborationHandler.ResolveComment)
			collaboration.POST("/comments/:id/unresolve", collaborationHandler.UnresolveComment)

			// Flow-specific comments
			collaboration.GET("/flows/:flow_id/comments", collaborationHandler.ListFlowComments)

			// Flow versions
			collaboration.GET("/flows/:flow_id/versions", collaborationHandler.ListFlowVersions)
			collaboration.GET("/flows/:flow_id/versions/compare", collaborationHandler.CompareVersions)
			collaboration.GET("/flows/:flow_id/versions/:version", collaborationHandler.GetFlowVersion)

			// Activity feed
			collaboration.GET("/activity", collaborationHandler.ListActivity)
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
