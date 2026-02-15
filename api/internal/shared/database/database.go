package database

import (
	"fmt"
	"time"

	"github.com/georgi-georgiev/testmesh/internal/shared/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// New creates a new database connection
func New(cfg config.DatabaseConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	// Set connection pool settings
	sqlDB.SetMaxOpenConns(cfg.MaxConns)
	sqlDB.SetMaxIdleConns(cfg.MaxIdle)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return db, nil
}

// AutoMigrate runs database migrations
func AutoMigrate(db *gorm.DB) error {
	// Create schemas
	if err := db.Exec("CREATE SCHEMA IF NOT EXISTS flows").Error; err != nil {
		return err
	}
	if err := db.Exec("CREATE SCHEMA IF NOT EXISTS executions").Error; err != nil {
		return err
	}
	if err := db.Exec("CREATE SCHEMA IF NOT EXISTS scheduler").Error; err != nil {
		return err
	}
	if err := db.Exec("CREATE SCHEMA IF NOT EXISTS mocks").Error; err != nil {
		return err
	}
	if err := db.Exec("CREATE SCHEMA IF NOT EXISTS contracts").Error; err != nil {
		return err
	}
	if err := db.Exec("CREATE SCHEMA IF NOT EXISTS reporting").Error; err != nil {
		return err
	}
	if err := db.Exec("CREATE SCHEMA IF NOT EXISTS ai").Error; err != nil {
		return err
	}

	// Import and migrate models
	// We'll do this manually here to avoid circular dependencies
	// In production, you might want to use a proper migration tool like golang-migrate

	// Create environments table
	db.Exec(`
		CREATE TABLE IF NOT EXISTS flows.environments (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			name VARCHAR(255) NOT NULL,
			description TEXT,
			color VARCHAR(20),
			is_default BOOLEAN DEFAULT false,
			variables JSONB DEFAULT '[]',
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			deleted_at TIMESTAMP WITH TIME ZONE
		);
		CREATE INDEX IF NOT EXISTS idx_environments_name ON flows.environments(name);
		CREATE INDEX IF NOT EXISTS idx_environments_deleted_at ON flows.environments(deleted_at);
	`)

	// Create collections table
	db.Exec(`
		CREATE TABLE IF NOT EXISTS flows.collections (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			name VARCHAR(255) NOT NULL,
			description TEXT,
			icon VARCHAR(100),
			color VARCHAR(20),
			parent_id UUID REFERENCES flows.collections(id),
			sort_order INTEGER DEFAULT 0,
			variables JSONB DEFAULT '{}',
			auth JSONB DEFAULT '{}',
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			deleted_at TIMESTAMP WITH TIME ZONE
		);
		CREATE INDEX IF NOT EXISTS idx_collections_name ON flows.collections(name);
		CREATE INDEX IF NOT EXISTS idx_collections_parent_id ON flows.collections(parent_id);
		CREATE INDEX IF NOT EXISTS idx_collections_deleted_at ON flows.collections(deleted_at);
	`)

	// Create flows table
	db.Exec(`
		CREATE TABLE IF NOT EXISTS flows.flows (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			name VARCHAR(255) NOT NULL,
			description TEXT,
			suite VARCHAR(255),
			tags TEXT[],
			definition JSONB NOT NULL,
			environment VARCHAR(50) DEFAULT 'default',
			collection_id UUID,
			sort_order INTEGER DEFAULT 0,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			deleted_at TIMESTAMP WITH TIME ZONE
		);
		CREATE INDEX IF NOT EXISTS idx_flows_name ON flows.flows(name);
		CREATE INDEX IF NOT EXISTS idx_flows_suite ON flows.flows(suite);
		CREATE INDEX IF NOT EXISTS idx_flows_deleted_at ON flows.flows(deleted_at);
		CREATE INDEX IF NOT EXISTS idx_flows_collection_id ON flows.flows(collection_id);
	`)

	// Add missing columns to existing flows table (idempotent migration)
	db.Exec(`
		ALTER TABLE flows.flows ADD COLUMN IF NOT EXISTS collection_id UUID;
		ALTER TABLE flows.flows ADD COLUMN IF NOT EXISTS sort_order INTEGER DEFAULT 0;
	`)

	// Create executions table
	db.Exec(`
		CREATE TABLE IF NOT EXISTS executions.executions (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			flow_id UUID NOT NULL REFERENCES flows.flows(id),
			status VARCHAR(20) NOT NULL DEFAULT 'pending',
			environment VARCHAR(50) DEFAULT 'default',
			started_at TIMESTAMP WITH TIME ZONE,
			finished_at TIMESTAMP WITH TIME ZONE,
			duration_ms BIGINT,
			total_steps INTEGER DEFAULT 0,
			passed_steps INTEGER DEFAULT 0,
			failed_steps INTEGER DEFAULT 0,
			error TEXT,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		);
		CREATE INDEX IF NOT EXISTS idx_executions_flow_id ON executions.executions(flow_id);
		CREATE INDEX IF NOT EXISTS idx_executions_status ON executions.executions(status);
	`)

	// Create execution_steps table
	db.Exec(`
		CREATE TABLE IF NOT EXISTS executions.execution_steps (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			execution_id UUID NOT NULL REFERENCES executions.executions(id) ON DELETE CASCADE,
			step_id VARCHAR(255) NOT NULL,
			step_name VARCHAR(255),
			action VARCHAR(100) NOT NULL,
			status VARCHAR(20) NOT NULL DEFAULT 'pending',
			started_at TIMESTAMP WITH TIME ZONE,
			finished_at TIMESTAMP WITH TIME ZONE,
			duration_ms BIGINT,
			output JSONB,
			error_message TEXT,
			attempt INTEGER DEFAULT 1,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		);
		CREATE INDEX IF NOT EXISTS idx_execution_steps_execution_id ON executions.execution_steps(execution_id);
	`)

	// Create mock_servers table
	db.Exec(`
		CREATE TABLE IF NOT EXISTS mocks.mock_servers (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			execution_id UUID REFERENCES executions.executions(id),
			name VARCHAR(255) NOT NULL,
			port INTEGER NOT NULL,
			base_url VARCHAR(500) NOT NULL,
			status VARCHAR(20) NOT NULL DEFAULT 'starting',
			started_at TIMESTAMP WITH TIME ZONE,
			stopped_at TIMESTAMP WITH TIME ZONE,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		);
		CREATE INDEX IF NOT EXISTS idx_mock_servers_execution_id ON mocks.mock_servers(execution_id);
		CREATE INDEX IF NOT EXISTS idx_mock_servers_port ON mocks.mock_servers(port);
	`)

	// Create mock_endpoints table
	db.Exec(`
		CREATE TABLE IF NOT EXISTS mocks.mock_endpoints (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			mock_server_id UUID NOT NULL REFERENCES mocks.mock_servers(id) ON DELETE CASCADE,
			path VARCHAR(500) NOT NULL,
			method VARCHAR(10) NOT NULL,
			match_config JSONB,
			response_config JSONB NOT NULL,
			state_config JSONB,
			priority INTEGER DEFAULT 0,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		);
		CREATE INDEX IF NOT EXISTS idx_mock_endpoints_server_id ON mocks.mock_endpoints(mock_server_id);
		CREATE INDEX IF NOT EXISTS idx_mock_endpoints_path_method ON mocks.mock_endpoints(path, method);
	`)

	// Create mock_requests table
	db.Exec(`
		CREATE TABLE IF NOT EXISTS mocks.mock_requests (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			mock_server_id UUID NOT NULL REFERENCES mocks.mock_servers(id) ON DELETE CASCADE,
			endpoint_id UUID REFERENCES mocks.mock_endpoints(id),
			method VARCHAR(10) NOT NULL,
			path VARCHAR(500) NOT NULL,
			headers JSONB,
			query_params JSONB,
			body TEXT,
			matched BOOLEAN DEFAULT false,
			response_code INTEGER,
			received_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		);
		CREATE INDEX IF NOT EXISTS idx_mock_requests_server_id ON mocks.mock_requests(mock_server_id);
		CREATE INDEX IF NOT EXISTS idx_mock_requests_endpoint_id ON mocks.mock_requests(endpoint_id);
	`)

	// Create mock_state table
	db.Exec(`
		CREATE TABLE IF NOT EXISTS mocks.mock_state (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			mock_server_id UUID NOT NULL REFERENCES mocks.mock_servers(id) ON DELETE CASCADE,
			state_key VARCHAR(255) NOT NULL,
			state_value JSONB NOT NULL,
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			UNIQUE(mock_server_id, state_key)
		);
		CREATE INDEX IF NOT EXISTS idx_mock_state_server_id ON mocks.mock_state(mock_server_id);
		CREATE INDEX IF NOT EXISTS idx_mock_state_key ON mocks.mock_state(state_key);
	`)

	// Create contracts table
	db.Exec(`
		CREATE TABLE IF NOT EXISTS contracts.contracts (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			consumer VARCHAR(255) NOT NULL,
			provider VARCHAR(255) NOT NULL,
			version VARCHAR(100) NOT NULL,
			pact_version VARCHAR(10) DEFAULT '4.0',
			contract_data JSONB NOT NULL,
			flow_id UUID REFERENCES flows.flows(id),
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			UNIQUE(consumer, provider, version)
		);
		CREATE INDEX IF NOT EXISTS idx_contracts_consumer ON contracts.contracts(consumer);
		CREATE INDEX IF NOT EXISTS idx_contracts_provider ON contracts.contracts(provider);
		CREATE INDEX IF NOT EXISTS idx_contracts_flow_id ON contracts.contracts(flow_id);
	`)

	// Create interactions table
	db.Exec(`
		CREATE TABLE IF NOT EXISTS contracts.interactions (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			contract_id UUID NOT NULL REFERENCES contracts.contracts(id) ON DELETE CASCADE,
			description TEXT NOT NULL,
			provider_state VARCHAR(500),
			request JSONB NOT NULL,
			response JSONB NOT NULL,
			interaction_type VARCHAR(50) DEFAULT 'http',
			metadata JSONB,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		);
		CREATE INDEX IF NOT EXISTS idx_interactions_contract_id ON contracts.interactions(contract_id);
	`)

	// Create verifications table
	db.Exec(`
		CREATE TABLE IF NOT EXISTS contracts.verifications (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			contract_id UUID NOT NULL REFERENCES contracts.contracts(id) ON DELETE CASCADE,
			provider_version VARCHAR(100) NOT NULL,
			status VARCHAR(20) NOT NULL DEFAULT 'pending',
			verified_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			results JSONB NOT NULL,
			execution_id UUID REFERENCES executions.executions(id),
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		);
		CREATE INDEX IF NOT EXISTS idx_verifications_contract_id ON contracts.verifications(contract_id);
		CREATE INDEX IF NOT EXISTS idx_verifications_execution_id ON contracts.verifications(execution_id);
		CREATE INDEX IF NOT EXISTS idx_verifications_status ON contracts.verifications(status);
	`)

	// Create breaking_changes table
	db.Exec(`
		CREATE TABLE IF NOT EXISTS contracts.breaking_changes (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			old_contract_id UUID NOT NULL REFERENCES contracts.contracts(id),
			new_contract_id UUID NOT NULL REFERENCES contracts.contracts(id),
			change_type VARCHAR(100) NOT NULL,
			severity VARCHAR(20) NOT NULL,
			description TEXT NOT NULL,
			details JSONB,
			detected_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		);
		CREATE INDEX IF NOT EXISTS idx_breaking_changes_old_contract ON contracts.breaking_changes(old_contract_id);
		CREATE INDEX IF NOT EXISTS idx_breaking_changes_new_contract ON contracts.breaking_changes(new_contract_id);
		CREATE INDEX IF NOT EXISTS idx_breaking_changes_severity ON contracts.breaking_changes(severity);
	`)

	// Create daily_metrics table for reporting
	db.Exec(`
		CREATE TABLE IF NOT EXISTS reporting.daily_metrics (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			date DATE NOT NULL,
			environment VARCHAR(50) NOT NULL,
			total_flows INTEGER DEFAULT 0,
			total_execs INTEGER DEFAULT 0,
			passed_execs INTEGER DEFAULT 0,
			failed_execs INTEGER DEFAULT 0,
			pass_rate DECIMAL(5,2) DEFAULT 0,
			avg_duration_ms BIGINT DEFAULT 0,
			p50_duration_ms BIGINT DEFAULT 0,
			p95_duration_ms BIGINT DEFAULT 0,
			p99_duration_ms BIGINT DEFAULT 0,
			total_steps INTEGER DEFAULT 0,
			passed_steps INTEGER DEFAULT 0,
			failed_steps INTEGER DEFAULT 0,
			by_flow_metrics JSONB,
			by_suite_metrics JSONB,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			UNIQUE(date, environment)
		);
		CREATE INDEX IF NOT EXISTS idx_daily_metrics_date ON reporting.daily_metrics(date);
		CREATE INDEX IF NOT EXISTS idx_daily_metrics_environment ON reporting.daily_metrics(environment);
	`)

	// Create flakiness_metrics table
	db.Exec(`
		CREATE TABLE IF NOT EXISTS reporting.flakiness_metrics (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			flow_id UUID NOT NULL REFERENCES flows.flows(id),
			window_start_date DATE NOT NULL,
			window_end_date DATE NOT NULL,
			window_days INTEGER NOT NULL,
			total_execs INTEGER DEFAULT 0,
			passed_execs INTEGER DEFAULT 0,
			failed_execs INTEGER DEFAULT 0,
			transitions INTEGER DEFAULT 0,
			flakiness_score DECIMAL(5,4) DEFAULT 0,
			is_flaky BOOLEAN DEFAULT false,
			failure_patterns TEXT[],
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		);
		CREATE INDEX IF NOT EXISTS idx_flakiness_metrics_flow_id ON reporting.flakiness_metrics(flow_id);
		CREATE INDEX IF NOT EXISTS idx_flakiness_metrics_is_flaky ON reporting.flakiness_metrics(is_flaky);
		CREATE INDEX IF NOT EXISTS idx_flakiness_metrics_window ON reporting.flakiness_metrics(window_start_date, window_end_date);
	`)

	// Create reports table
	db.Exec(`
		CREATE TABLE IF NOT EXISTS reporting.reports (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			name VARCHAR(255) NOT NULL,
			format VARCHAR(20) NOT NULL,
			status VARCHAR(20) NOT NULL DEFAULT 'pending',
			filters JSONB,
			start_date DATE,
			end_date DATE,
			file_path VARCHAR(500),
			file_size BIGINT DEFAULT 0,
			generated_at TIMESTAMP WITH TIME ZONE,
			expires_at TIMESTAMP WITH TIME ZONE,
			error TEXT,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		);
		CREATE INDEX IF NOT EXISTS idx_reports_status ON reporting.reports(status);
		CREATE INDEX IF NOT EXISTS idx_reports_format ON reporting.reports(format);
		CREATE INDEX IF NOT EXISTS idx_reports_expires_at ON reporting.reports(expires_at);
	`)

	// Create step_performance table
	db.Exec(`
		CREATE TABLE IF NOT EXISTS reporting.step_performance (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			flow_id UUID NOT NULL REFERENCES flows.flows(id),
			step_id VARCHAR(255) NOT NULL,
			step_name VARCHAR(255),
			action VARCHAR(100) NOT NULL,
			date DATE NOT NULL,
			execution_count INTEGER DEFAULT 0,
			passed_count INTEGER DEFAULT 0,
			failed_count INTEGER DEFAULT 0,
			pass_rate DECIMAL(5,2) DEFAULT 0,
			avg_duration_ms BIGINT DEFAULT 0,
			min_duration_ms BIGINT DEFAULT 0,
			max_duration_ms BIGINT DEFAULT 0,
			p50_duration_ms BIGINT DEFAULT 0,
			p95_duration_ms BIGINT DEFAULT 0,
			p99_duration_ms BIGINT DEFAULT 0,
			common_errors TEXT[],
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		);
		CREATE INDEX IF NOT EXISTS idx_step_performance_flow_id ON reporting.step_performance(flow_id);
		CREATE INDEX IF NOT EXISTS idx_step_performance_step_id ON reporting.step_performance(step_id);
		CREATE INDEX IF NOT EXISTS idx_step_performance_date ON reporting.step_performance(date);
		CREATE INDEX IF NOT EXISTS idx_step_performance_action ON reporting.step_performance(action);
	`)

	// Create AI generation_history table
	db.Exec(`
		CREATE TABLE IF NOT EXISTS ai.generation_history (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			provider VARCHAR(50) NOT NULL,
			model VARCHAR(100) NOT NULL,
			prompt TEXT NOT NULL,
			status VARCHAR(20) NOT NULL DEFAULT 'pending',
			generated_yaml TEXT,
			flow_id UUID REFERENCES flows.flows(id),
			tokens_used INTEGER DEFAULT 0,
			latency_ms BIGINT DEFAULT 0,
			error TEXT,
			metadata JSONB,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		);
		CREATE INDEX IF NOT EXISTS idx_generation_history_provider ON ai.generation_history(provider);
		CREATE INDEX IF NOT EXISTS idx_generation_history_status ON ai.generation_history(status);
		CREATE INDEX IF NOT EXISTS idx_generation_history_flow_id ON ai.generation_history(flow_id);
	`)

	// Create AI suggestions table
	db.Exec(`
		CREATE TABLE IF NOT EXISTS ai.suggestions (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			flow_id UUID NOT NULL REFERENCES flows.flows(id),
			execution_id UUID REFERENCES executions.executions(id),
			type VARCHAR(50) NOT NULL,
			status VARCHAR(20) NOT NULL DEFAULT 'pending',
			title VARCHAR(500) NOT NULL,
			description TEXT,
			original_yaml TEXT,
			suggested_yaml TEXT,
			diff_patch TEXT,
			confidence DECIMAL(5,4) DEFAULT 0,
			reasoning TEXT,
			applied_at TIMESTAMP WITH TIME ZONE,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		);
		CREATE INDEX IF NOT EXISTS idx_suggestions_flow_id ON ai.suggestions(flow_id);
		CREATE INDEX IF NOT EXISTS idx_suggestions_execution_id ON ai.suggestions(execution_id);
		CREATE INDEX IF NOT EXISTS idx_suggestions_type ON ai.suggestions(type);
		CREATE INDEX IF NOT EXISTS idx_suggestions_status ON ai.suggestions(status);
	`)

	// Create AI import_history table
	db.Exec(`
		CREATE TABLE IF NOT EXISTS ai.import_history (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			source_type VARCHAR(50) NOT NULL,
			source_name VARCHAR(255) NOT NULL,
			source_content TEXT,
			source_url VARCHAR(500),
			status VARCHAR(20) NOT NULL DEFAULT 'pending',
			flows_generated INTEGER DEFAULT 0,
			flow_ids TEXT[],
			error TEXT,
			metadata JSONB,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		);
		CREATE INDEX IF NOT EXISTS idx_import_history_source_type ON ai.import_history(source_type);
		CREATE INDEX IF NOT EXISTS idx_import_history_status ON ai.import_history(status);
	`)

	// Create AI coverage_analysis table
	db.Exec(`
		CREATE TABLE IF NOT EXISTS ai.coverage_analysis (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			spec_type VARCHAR(50) NOT NULL,
			spec_name VARCHAR(255) NOT NULL,
			spec_content TEXT,
			spec_url VARCHAR(500),
			status VARCHAR(20) NOT NULL DEFAULT 'pending',
			total_endpoints INTEGER DEFAULT 0,
			covered_endpoints INTEGER DEFAULT 0,
			coverage_percent DECIMAL(5,2) DEFAULT 0,
			results JSONB,
			error TEXT,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		);
		CREATE INDEX IF NOT EXISTS idx_coverage_analysis_spec_type ON ai.coverage_analysis(spec_type);
		CREATE INDEX IF NOT EXISTS idx_coverage_analysis_status ON ai.coverage_analysis(status);
	`)

	// Create AI usage_stats table
	db.Exec(`
		CREATE TABLE IF NOT EXISTS ai.usage_stats (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			provider VARCHAR(50) NOT NULL,
			model VARCHAR(100) NOT NULL,
			date DATE NOT NULL,
			total_requests INTEGER DEFAULT 0,
			total_tokens INTEGER DEFAULT 0,
			success_count INTEGER DEFAULT 0,
			failure_count INTEGER DEFAULT 0,
			avg_latency_ms BIGINT DEFAULT 0,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			UNIQUE(provider, model, date)
		);
		CREATE INDEX IF NOT EXISTS idx_usage_stats_provider ON ai.usage_stats(provider);
		CREATE INDEX IF NOT EXISTS idx_usage_stats_date ON ai.usage_stats(date);
	`)

	// Create schedules table
	db.Exec(`
		CREATE TABLE IF NOT EXISTS schedules (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			name VARCHAR(255) NOT NULL,
			description TEXT,
			flow_id UUID NOT NULL REFERENCES flows.flows(id),
			cron_expr VARCHAR(100) NOT NULL,
			timezone VARCHAR(50) DEFAULT 'UTC',
			status VARCHAR(20) NOT NULL DEFAULT 'active',
			environment JSONB DEFAULT '{}',
			notify_on_failure BOOLEAN DEFAULT false,
			notify_on_success BOOLEAN DEFAULT false,
			notify_emails TEXT[] DEFAULT '{}',
			max_retries INTEGER DEFAULT 0,
			retry_delay VARCHAR(20) DEFAULT '1m',
			allow_overlap BOOLEAN DEFAULT false,
			next_run_at TIMESTAMP WITH TIME ZONE,
			last_run_at TIMESTAMP WITH TIME ZONE,
			last_run_id UUID,
			last_run_result VARCHAR(20),
			tags TEXT[] DEFAULT '{}',
			created_by UUID,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		);
		CREATE INDEX IF NOT EXISTS idx_schedules_flow_id ON schedules(flow_id);
		CREATE INDEX IF NOT EXISTS idx_schedules_status ON schedules(status);
		CREATE INDEX IF NOT EXISTS idx_schedules_next_run_at ON schedules(next_run_at);
	`)

	// Migrate existing jsonb columns to text[] (idempotent)
	db.Exec(`
		DO $$
		BEGIN
			IF EXISTS (
				SELECT 1 FROM information_schema.columns
				WHERE table_name = 'schedules' AND column_name = 'notify_emails' AND data_type = 'jsonb'
			) THEN
				ALTER TABLE schedules DROP COLUMN notify_emails;
				ALTER TABLE schedules ADD COLUMN notify_emails TEXT[] DEFAULT '{}';
			END IF;
			IF EXISTS (
				SELECT 1 FROM information_schema.columns
				WHERE table_name = 'schedules' AND column_name = 'tags' AND data_type = 'jsonb'
			) THEN
				ALTER TABLE schedules DROP COLUMN tags;
				ALTER TABLE schedules ADD COLUMN tags TEXT[] DEFAULT '{}';
			END IF;
		END $$;
	`)

	// Create schedule_runs table
	db.Exec(`
		CREATE TABLE IF NOT EXISTS schedule_runs (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			schedule_id UUID NOT NULL REFERENCES schedules(id) ON DELETE CASCADE,
			execution_id UUID,
			status VARCHAR(20) NOT NULL DEFAULT 'pending',
			result VARCHAR(20),
			error TEXT,
			retry_count INTEGER DEFAULT 0,
			scheduled_at TIMESTAMP WITH TIME ZONE NOT NULL,
			started_at TIMESTAMP WITH TIME ZONE,
			completed_at TIMESTAMP WITH TIME ZONE,
			duration BIGINT DEFAULT 0,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		);
		CREATE INDEX IF NOT EXISTS idx_schedule_runs_schedule_id ON schedule_runs(schedule_id);
		CREATE INDEX IF NOT EXISTS idx_schedule_runs_status ON schedule_runs(status);
		CREATE INDEX IF NOT EXISTS idx_schedule_runs_scheduled_at ON schedule_runs(scheduled_at);
	`)

	// Create request_history table
	db.Exec(`
		CREATE TABLE IF NOT EXISTS flows.request_history (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			flow_id UUID REFERENCES flows.flows(id),
			collection_id UUID REFERENCES flows.collections(id),
			method VARCHAR(10) NOT NULL,
			url VARCHAR(2048) NOT NULL,
			request JSONB NOT NULL,
			response JSONB,
			status_code INTEGER,
			duration_ms BIGINT,
			size_bytes BIGINT,
			error TEXT,
			tags TEXT[],
			saved_at TIMESTAMP WITH TIME ZONE,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		);
		CREATE INDEX IF NOT EXISTS idx_request_history_flow_id ON flows.request_history(flow_id);
		CREATE INDEX IF NOT EXISTS idx_request_history_collection_id ON flows.request_history(collection_id);
		CREATE INDEX IF NOT EXISTS idx_request_history_url ON flows.request_history(url);
		CREATE INDEX IF NOT EXISTS idx_request_history_created_at ON flows.request_history(created_at);
	`)

	// Add collection_items table for organizing flows within collections
	db.Exec(`
		CREATE TABLE IF NOT EXISTS flows.collection_items (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			collection_id UUID NOT NULL REFERENCES flows.collections(id) ON DELETE CASCADE,
			flow_id UUID NOT NULL REFERENCES flows.flows(id) ON DELETE CASCADE,
			sort_order INTEGER DEFAULT 0,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			UNIQUE(collection_id, flow_id)
		);
		CREATE INDEX IF NOT EXISTS idx_collection_items_collection_id ON flows.collection_items(collection_id);
		CREATE INDEX IF NOT EXISTS idx_collection_items_flow_id ON flows.collection_items(flow_id);
	`)

	// Create workspaces table
	db.Exec(`
		CREATE TABLE IF NOT EXISTS workspaces (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			name VARCHAR(255) NOT NULL,
			slug VARCHAR(255) NOT NULL UNIQUE,
			description TEXT,
			type VARCHAR(20) NOT NULL DEFAULT 'personal',
			owner_id UUID,
			settings JSONB DEFAULT '{}',
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			deleted_at TIMESTAMP WITH TIME ZONE
		);
		CREATE INDEX IF NOT EXISTS idx_workspaces_name ON workspaces(name);
		CREATE INDEX IF NOT EXISTS idx_workspaces_slug ON workspaces(slug);
		CREATE INDEX IF NOT EXISTS idx_workspaces_owner_id ON workspaces(owner_id);
		CREATE INDEX IF NOT EXISTS idx_workspaces_deleted_at ON workspaces(deleted_at);
	`)

	// Create workspace_members table
	db.Exec(`
		CREATE TABLE IF NOT EXISTS workspace_members (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			workspace_id UUID NOT NULL REFERENCES workspaces(id) ON DELETE CASCADE,
			user_id UUID NOT NULL,
			email VARCHAR(255),
			name VARCHAR(255),
			role VARCHAR(20) NOT NULL DEFAULT 'viewer',
			invited_by UUID,
			invited_at TIMESTAMP WITH TIME ZONE,
			joined_at TIMESTAMP WITH TIME ZONE,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		);
		CREATE INDEX IF NOT EXISTS idx_workspace_members_workspace_id ON workspace_members(workspace_id);
		CREATE INDEX IF NOT EXISTS idx_workspace_members_user_id ON workspace_members(user_id);
		CREATE INDEX IF NOT EXISTS idx_workspace_members_email ON workspace_members(email);
	`)

	// Create workspace_invitations table
	db.Exec(`
		CREATE TABLE IF NOT EXISTS workspace_invitations (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			workspace_id UUID NOT NULL REFERENCES workspaces(id) ON DELETE CASCADE,
			email VARCHAR(255) NOT NULL,
			role VARCHAR(20) NOT NULL,
			token VARCHAR(255) NOT NULL UNIQUE,
			invited_by UUID NOT NULL,
			expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		);
		CREATE INDEX IF NOT EXISTS idx_workspace_invitations_workspace_id ON workspace_invitations(workspace_id);
		CREATE INDEX IF NOT EXISTS idx_workspace_invitations_email ON workspace_invitations(email);
		CREATE INDEX IF NOT EXISTS idx_workspace_invitations_token ON workspace_invitations(token);
	`)

	// Create collaboration schema
	db.Exec(`CREATE SCHEMA IF NOT EXISTS collaboration`)

	// Create user_presences table
	db.Exec(`
		CREATE TABLE IF NOT EXISTS user_presences (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			user_id UUID NOT NULL,
			user_name VARCHAR(255) NOT NULL,
			user_email VARCHAR(255),
			user_avatar VARCHAR(500),
			color VARCHAR(20) NOT NULL,
			resource_type VARCHAR(50) NOT NULL,
			resource_id UUID NOT NULL,
			status VARCHAR(20) NOT NULL DEFAULT 'viewing',
			cursor_data JSONB,
			last_active_at TIMESTAMP WITH TIME ZONE NOT NULL,
			connected_at TIMESTAMP WITH TIME ZONE NOT NULL,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		);
		CREATE INDEX IF NOT EXISTS idx_user_presences_user_id ON user_presences(user_id);
		CREATE INDEX IF NOT EXISTS idx_user_presences_resource ON user_presences(resource_type, resource_id);
		CREATE INDEX IF NOT EXISTS idx_user_presences_last_active ON user_presences(last_active_at);
	`)

	// Create flow_comments table
	db.Exec(`
		CREATE TABLE IF NOT EXISTS flow_comments (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			flow_id UUID NOT NULL,
			step_id VARCHAR(255),
			parent_id UUID REFERENCES flow_comments(id),
			author_id UUID NOT NULL,
			author_name VARCHAR(255) NOT NULL,
			author_avatar VARCHAR(500),
			content TEXT NOT NULL,
			resolved BOOLEAN DEFAULT false,
			position JSONB,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			deleted_at TIMESTAMP WITH TIME ZONE
		);
		CREATE INDEX IF NOT EXISTS idx_flow_comments_flow_id ON flow_comments(flow_id);
		CREATE INDEX IF NOT EXISTS idx_flow_comments_step_id ON flow_comments(step_id);
		CREATE INDEX IF NOT EXISTS idx_flow_comments_parent_id ON flow_comments(parent_id);
		CREATE INDEX IF NOT EXISTS idx_flow_comments_deleted_at ON flow_comments(deleted_at);
	`)

	// Create activity_events table
	db.Exec(`
		CREATE TABLE IF NOT EXISTS activity_events (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			actor_id UUID,
			actor_name VARCHAR(255) NOT NULL,
			actor_avatar VARCHAR(500),
			event_type VARCHAR(100) NOT NULL,
			resource_type VARCHAR(50) NOT NULL,
			resource_id UUID NOT NULL,
			resource_name VARCHAR(255),
			description TEXT,
			changes JSONB,
			metadata JSONB,
			workspace_id UUID,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		);
		CREATE INDEX IF NOT EXISTS idx_activity_events_actor_id ON activity_events(actor_id);
		CREATE INDEX IF NOT EXISTS idx_activity_events_event_type ON activity_events(event_type);
		CREATE INDEX IF NOT EXISTS idx_activity_events_resource ON activity_events(resource_type, resource_id);
		CREATE INDEX IF NOT EXISTS idx_activity_events_workspace_id ON activity_events(workspace_id);
		CREATE INDEX IF NOT EXISTS idx_activity_events_created_at ON activity_events(created_at);
	`)

	// Create flow_versions table
	db.Exec(`
		CREATE TABLE IF NOT EXISTS flow_versions (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			flow_id UUID NOT NULL,
			version INTEGER NOT NULL,
			content TEXT NOT NULL,
			author_id UUID,
			author_name VARCHAR(255),
			message VARCHAR(500),
			description TEXT,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		);
		CREATE INDEX IF NOT EXISTS idx_flow_versions_flow_id ON flow_versions(flow_id);
		CREATE UNIQUE INDEX IF NOT EXISTS idx_flow_versions_flow_version ON flow_versions(flow_id, version);
	`)

	return nil
}
