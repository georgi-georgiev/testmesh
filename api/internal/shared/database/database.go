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
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			deleted_at TIMESTAMP WITH TIME ZONE
		);
		CREATE INDEX IF NOT EXISTS idx_flows_name ON flows.flows(name);
		CREATE INDEX IF NOT EXISTS idx_flows_suite ON flows.flows(suite);
		CREATE INDEX IF NOT EXISTS idx_flows_deleted_at ON flows.flows(deleted_at);
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

	return nil
}
