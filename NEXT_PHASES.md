# TestMesh v1.0 - Remaining Implementation Phases

**Status**: Phases 1 & 2 Complete ✅
**Last Updated**: 2026-02-11
**Remaining Work**: Phases 3 & 4 (Advanced Reporting & AI Integration)

---

## Completed Work Summary

### ✅ Phase 1: Mock Server System (Complete)
- **Backend**: Full implementation with port pool management, request matching, stateful responses
- **Frontend**: List/detail pages with auto-refreshing request logs
- **Database**: `mocks` schema with 4 tables (mock_servers, mock_endpoints, mock_requests, mock_state)
- **Action Handlers**: `mock_server_start`, `mock_server_stop`, `mock_server_configure`
- **Examples**: `mock-server-example.yaml`, `contract-simple-example.yaml`
- **Critical Fix**: Added support for both `${VAR}` and `{{VAR}}` variable interpolation formats

### ✅ Phase 2: Contract Testing (Complete)
- **Backend**: Pact v4.0 compatible generator, verifier, and breaking change detection
- **Frontend**: Contract list/detail pages with interactions, verifications, and breaking changes tabs
- **Database**: `contracts` schema with 4 tables (contracts, interactions, verifications, breaking_changes)
- **Action Handlers**: `contract_generate`, `contract_verify`
- **Examples**: `contract-testing-example.yaml`, `contract-simple-example.yaml`
- **Export**: Pact JSON export/import functionality

### 🔧 Known Issues & Fixes
- **Variable Interpolation**: Updated `/api/internal/runner/interpolator.go` to support both `${VAR}` and `{{VAR}}` formats
- **YAML Syntax**: Fixed colon-in-string issues in example files (must quote strings with colons)
- **Navigation**: Added Mock Servers and Contracts cards to dashboard

---

## Phase 3: Advanced Reporting & Analytics

**Timeline**: 4 weeks
**Complexity**: Medium
**Priority**: High (Business Intelligence)

### Overview
Provide comprehensive test analytics with HTML reports, trend analysis, flakiness detection, and multiple export formats. This phase focuses on data aggregation and visualization without external dependencies.

### Database Schema

**New schema**: `reporting`

**Tables to create**:

```sql
-- Daily aggregated metrics
CREATE TABLE IF NOT EXISTS reporting.daily_metrics (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    date DATE NOT NULL,
    flow_id UUID REFERENCES flows.flows(id) ON DELETE CASCADE,
    suite VARCHAR(255),
    total_executions INTEGER NOT NULL DEFAULT 0,
    passed_executions INTEGER NOT NULL DEFAULT 0,
    failed_executions INTEGER NOT NULL DEFAULT 0,
    avg_duration_ms BIGINT,
    p50_duration_ms BIGINT,
    p95_duration_ms BIGINT,
    p99_duration_ms BIGINT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(date, flow_id)
);
CREATE INDEX idx_daily_metrics_date ON reporting.daily_metrics(date);
CREATE INDEX idx_daily_metrics_suite ON reporting.daily_metrics(suite);

-- Flakiness tracking
CREATE TABLE IF NOT EXISTS reporting.flakiness_metrics (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    flow_id UUID NOT NULL REFERENCES flows.flows(id) ON DELETE CASCADE,
    window_start TIMESTAMP WITH TIME ZONE NOT NULL,
    window_end TIMESTAMP WITH TIME ZONE NOT NULL,
    total_runs INTEGER NOT NULL,
    passed_runs INTEGER NOT NULL,
    failed_runs INTEGER NOT NULL,
    flakiness_score DECIMAL(5,2) NOT NULL, -- 0-100 scale
    consecutive_failures INTEGER NOT NULL DEFAULT 0,
    consecutive_successes INTEGER NOT NULL DEFAULT 0,
    calculated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(flow_id, window_start)
);
CREATE INDEX idx_flakiness_flow ON reporting.flakiness_metrics(flow_id);

-- Generated reports
CREATE TABLE IF NOT EXISTS reporting.reports (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    report_type VARCHAR(50) NOT NULL, -- 'html', 'junit', 'json'
    title VARCHAR(500) NOT NULL,
    filters JSONB,
    file_path VARCHAR(1000) NOT NULL,
    file_size_bytes BIGINT,
    generated_by VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP WITH TIME ZONE
);
CREATE INDEX idx_reports_created ON reporting.reports(created_at);
CREATE INDEX idx_reports_type ON reporting.reports(report_type);

-- Step-level performance tracking
CREATE TABLE IF NOT EXISTS reporting.step_performance (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    date DATE NOT NULL,
    step_action VARCHAR(100) NOT NULL,
    flow_id UUID REFERENCES flows.flows(id) ON DELETE CASCADE,
    avg_duration_ms BIGINT NOT NULL,
    min_duration_ms BIGINT NOT NULL,
    max_duration_ms BIGINT NOT NULL,
    total_executions INTEGER NOT NULL,
    error_count INTEGER NOT NULL DEFAULT 0,
    error_rate DECIMAL(5,2) NOT NULL DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(date, step_action, flow_id)
);
CREATE INDEX idx_step_perf_date ON reporting.step_performance(date);
CREATE INDEX idx_step_perf_action ON reporting.step_performance(step_action);
```

### Backend Implementation

#### 1. Models (`/api/internal/storage/models/reporting.go`)

```go
package models

import (
    "database/sql/driver"
    "encoding/json"
    "time"
    "github.com/google/uuid"
)

type DailyMetric struct {
    ID                uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
    Date              time.Time  `gorm:"type:date;not null;index" json:"date"`
    FlowID            *uuid.UUID `gorm:"type:uuid;index" json:"flow_id,omitempty"`
    Suite             string     `gorm:"type:varchar(255);index" json:"suite,omitempty"`
    TotalExecutions   int        `gorm:"not null;default:0" json:"total_executions"`
    PassedExecutions  int        `gorm:"not null;default:0" json:"passed_executions"`
    FailedExecutions  int        `gorm:"not null;default:0" json:"failed_executions"`
    AvgDurationMs     *int64     `json:"avg_duration_ms,omitempty"`
    P50DurationMs     *int64     `json:"p50_duration_ms,omitempty"`
    P95DurationMs     *int64     `json:"p95_duration_ms,omitempty"`
    P99DurationMs     *int64     `json:"p99_duration_ms,omitempty"`
    CreatedAt         time.Time  `json:"created_at"`
}

func (DailyMetric) TableName() string {
    return "reporting.daily_metrics"
}

type FlakinessMetric struct {
    ID                    uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
    FlowID                uuid.UUID `gorm:"type:uuid;not null;index" json:"flow_id"`
    WindowStart           time.Time `gorm:"not null" json:"window_start"`
    WindowEnd             time.Time `gorm:"not null" json:"window_end"`
    TotalRuns             int       `gorm:"not null" json:"total_runs"`
    PassedRuns            int       `gorm:"not null" json:"passed_runs"`
    FailedRuns            int       `gorm:"not null" json:"failed_runs"`
    FlakinessScore        float64   `gorm:"type:decimal(5,2);not null" json:"flakiness_score"`
    ConsecutiveFailures   int       `gorm:"not null;default:0" json:"consecutive_failures"`
    ConsecutiveSuccesses  int       `gorm:"not null;default:0" json:"consecutive_successes"`
    CalculatedAt          time.Time `json:"calculated_at"`
}

func (FlakinessMetric) TableName() string {
    return "reporting.flakiness_metrics"
}

type ReportType string

const (
    ReportTypeHTML  ReportType = "html"
    ReportTypeJUnit ReportType = "junit"
    ReportTypeJSON  ReportType = "json"
)

type Report struct {
    ID            uuid.UUID       `gorm:"type:uuid;primary_key" json:"id"`
    ReportType    ReportType      `gorm:"type:varchar(50);not null;index" json:"report_type"`
    Title         string          `gorm:"type:varchar(500);not null" json:"title"`
    Filters       ReportFilters   `gorm:"type:jsonb" json:"filters,omitempty"`
    FilePath      string          `gorm:"type:varchar(1000);not null" json:"file_path"`
    FileSizeBytes *int64          `json:"file_size_bytes,omitempty"`
    GeneratedBy   string          `gorm:"type:varchar(255)" json:"generated_by,omitempty"`
    CreatedAt     time.Time       `json:"created_at"`
    ExpiresAt     *time.Time      `json:"expires_at,omitempty"`
}

func (Report) TableName() string {
    return "reporting.reports"
}

type ReportFilters struct {
    StartDate   *time.Time `json:"start_date,omitempty"`
    EndDate     *time.Time `json:"end_date,omitempty"`
    Suite       string     `json:"suite,omitempty"`
    FlowIDs     []string   `json:"flow_ids,omitempty"`
    Status      string     `json:"status,omitempty"`
}

func (rf *ReportFilters) Scan(value interface{}) error {
    bytes, ok := value.([]byte)
    if !ok {
        return nil
    }
    return json.Unmarshal(bytes, rf)
}

func (rf ReportFilters) Value() (driver.Value, error) {
    return json.Marshal(rf)
}

type StepPerformance struct {
    ID              uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
    Date            time.Time  `gorm:"type:date;not null;index" json:"date"`
    StepAction      string     `gorm:"type:varchar(100);not null;index" json:"step_action"`
    FlowID          *uuid.UUID `gorm:"type:uuid" json:"flow_id,omitempty"`
    AvgDurationMs   int64      `gorm:"not null" json:"avg_duration_ms"`
    MinDurationMs   int64      `gorm:"not null" json:"min_duration_ms"`
    MaxDurationMs   int64      `gorm:"not null" json:"max_duration_ms"`
    TotalExecutions int        `gorm:"not null" json:"total_executions"`
    ErrorCount      int        `gorm:"not null;default:0" json:"error_count"`
    ErrorRate       float64    `gorm:"type:decimal(5,2);not null;default:0" json:"error_rate"`
    CreatedAt       time.Time  `json:"created_at"`
}

func (StepPerformance) TableName() string {
    return "reporting.step_performance"
}
```

#### 2. Repository (`/api/internal/storage/repository/reporting.go`)

```go
package repository

import (
    "time"
    "github.com/georgi-georgiev/testmesh/internal/storage/models"
    "github.com/google/uuid"
    "gorm.io/gorm"
)

type ReportingRepository struct {
    db *gorm.DB
}

func NewReportingRepository(db *gorm.DB) *ReportingRepository {
    return &ReportingRepository{db: db}
}

// Daily Metrics
func (r *ReportingRepository) CreateDailyMetric(metric *models.DailyMetric) error {
    return r.db.Create(metric).Error
}

func (r *ReportingRepository) GetDailyMetrics(startDate, endDate time.Time, suite string, flowID *uuid.UUID) ([]models.DailyMetric, error) {
    var metrics []models.DailyMetric
    query := r.db.Where("date BETWEEN ? AND ?", startDate, endDate)

    if suite != "" {
        query = query.Where("suite = ?", suite)
    }
    if flowID != nil {
        query = query.Where("flow_id = ?", *flowID)
    }

    err := query.Order("date DESC").Find(&metrics).Error
    return metrics, err
}

// Flakiness Metrics
func (r *ReportingRepository) CreateFlakinessMetric(metric *models.FlakinessMetric) error {
    return r.db.Create(metric).Error
}

func (r *ReportingRepository) GetFlakyFlows(threshold float64, limit int) ([]models.FlakinessMetric, error) {
    var metrics []models.FlakinessMetric
    err := r.db.Where("flakiness_score >= ?", threshold).
        Order("flakiness_score DESC").
        Limit(limit).
        Find(&metrics).Error
    return metrics, err
}

func (r *ReportingRepository) GetFlakinessHistory(flowID uuid.UUID, days int) ([]models.FlakinessMetric, error) {
    var metrics []models.FlakinessMetric
    startDate := time.Now().AddDate(0, 0, -days)
    err := r.db.Where("flow_id = ? AND window_start >= ?", flowID, startDate).
        Order("window_start DESC").
        Find(&metrics).Error
    return metrics, err
}

// Reports
func (r *ReportingRepository) CreateReport(report *models.Report) error {
    return r.db.Create(report).Error
}

func (r *ReportingRepository) GetReportByID(id uuid.UUID) (*models.Report, error) {
    var report models.Report
    err := r.db.Where("id = ?", id).First(&report).Error
    return &report, err
}

func (r *ReportingRepository) ListReports(reportType models.ReportType, limit, offset int) ([]models.Report, int64, error) {
    var reports []models.Report
    var total int64

    query := r.db.Model(&models.Report{})
    if reportType != "" {
        query = query.Where("report_type = ?", reportType)
    }

    query.Count(&total)
    err := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&reports).Error
    return reports, total, err
}

func (r *ReportingRepository) DeleteExpiredReports() (int64, error) {
    result := r.db.Where("expires_at IS NOT NULL AND expires_at < ?", time.Now()).Delete(&models.Report{})
    return result.RowsAffected, result.Error
}

// Step Performance
func (r *ReportingRepository) CreateStepPerformance(perf *models.StepPerformance) error {
    return r.db.Create(perf).Error
}

func (r *ReportingRepository) GetStepPerformance(date time.Time, action string) ([]models.StepPerformance, error) {
    var perfs []models.StepPerformance
    query := r.db.Where("date = ?", date)
    if action != "" {
        query = query.Where("step_action = ?", action)
    }
    err := query.Find(&perfs).Error
    return perfs, err
}
```

#### 3. Aggregation Service (`/api/internal/reporting/aggregator.go`)

Create this service to calculate daily metrics and flakiness scores:

```go
package reporting

import (
    "time"
    "github.com/georgi-georgiev/testmesh/internal/storage/models"
    "github.com/georgi-georgiev/testmesh/internal/storage/repository"
    "go.uber.org/zap"
)

type Aggregator struct {
    execRepo     *repository.ExecutionRepository
    reportingRepo *repository.ReportingRepository
    logger       *zap.Logger
}

func NewAggregator(execRepo *repository.ExecutionRepository, reportingRepo *repository.ReportingRepository, logger *zap.Logger) *Aggregator {
    return &Aggregator{
        execRepo:     execRepo,
        reportingRepo: reportingRepo,
        logger:       logger,
    }
}

// AggregateDailyMetrics calculates daily metrics for a given date
func (a *Aggregator) AggregateDailyMetrics(date time.Time) error {
    // TODO: Implement aggregation logic
    // 1. Query executions for the date
    // 2. Group by flow_id and suite
    // 3. Calculate totals, pass rate, duration percentiles
    // 4. Save to reporting.daily_metrics
    return nil
}

// CalculateFlakiness calculates flakiness scores for flows
func (a *Aggregator) CalculateFlakiness(windowDays int) error {
    // TODO: Implement flakiness calculation
    // 1. Query recent executions for each flow
    // 2. Analyze pass/fail patterns
    // 3. Calculate flakiness score (0-100)
    // 4. Save to reporting.flakiness_metrics
    return nil
}

// ScheduleAggregation runs aggregation on a schedule
func (a *Aggregator) ScheduleAggregation() {
    // TODO: Implement scheduled aggregation (daily at 2 AM)
}
```

#### 4. Report Generator (`/api/internal/reporting/generator.go`)

```go
package reporting

import (
    "github.com/georgi-georgiev/testmesh/internal/storage/models"
    "github.com/georgi-georgiev/testmesh/internal/storage/repository"
    "go.uber.org/zap"
)

type Generator struct {
    execRepo      *repository.ExecutionRepository
    reportingRepo *repository.ReportingRepository
    logger        *zap.Logger
}

func NewGenerator(execRepo *repository.ExecutionRepository, reportingRepo *repository.ReportingRepository, logger *zap.Logger) *Generator {
    return &Generator{
        execRepo:      execRepo,
        reportingRepo: reportingRepo,
        logger:        logger,
    }
}

// GenerateHTMLReport creates an HTML report
func (g *Generator) GenerateHTMLReport(title string, filters models.ReportFilters) (*models.Report, error) {
    // TODO: Implement HTML generation
    // 1. Query data based on filters
    // 2. Generate HTML with charts (using templates)
    // 3. Save to file
    // 4. Create report record
    return nil, nil
}

// GenerateJUnitXML creates JUnit XML report
func (g *Generator) GenerateJUnitXML(filters models.ReportFilters) (*models.Report, error) {
    // TODO: Implement JUnit XML generation
    return nil, nil
}

// GenerateJSON creates JSON report
func (g *Generator) GenerateJSON(filters models.ReportFilters) (*models.Report, error) {
    // TODO: Implement JSON report generation
    return nil, nil
}
```

#### 5. API Handlers (`/api/internal/api/handlers/reporting.go`)

```go
package handlers

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "github.com/georgi-georgiev/testmesh/internal/reporting"
    "github.com/georgi-georgiev/testmesh/internal/storage/repository"
    "go.uber.org/zap"
)

type ReportingHandler struct {
    repo      *repository.ReportingRepository
    generator *reporting.Generator
    logger    *zap.Logger
}

func NewReportingHandler(repo *repository.ReportingRepository, generator *reporting.Generator, logger *zap.Logger) *ReportingHandler {
    return &ReportingHandler{
        repo:      repo,
        generator: generator,
        logger:    logger,
    }
}

// POST /api/v1/reports/generate
func (h *ReportingHandler) Generate(c *gin.Context) {
    // TODO: Parse request, generate report, return report ID
}

// GET /api/v1/reports
func (h *ReportingHandler) List(c *gin.Context) {
    // TODO: List generated reports
}

// GET /api/v1/reports/:id
func (h *ReportingHandler) Get(c *gin.Context) {
    // TODO: Get report details
}

// GET /api/v1/reports/:id/download
func (h *ReportingHandler) Download(c *gin.Context) {
    // TODO: Serve report file
}

// GET /api/v1/analytics/metrics
func (h *ReportingHandler) GetMetrics(c *gin.Context) {
    // TODO: Get daily metrics with filters
}

// GET /api/v1/analytics/flakiness
func (h *ReportingHandler) GetFlakiness(c *gin.Context) {
    // TODO: Get flakiness metrics
}

// GET /api/v1/analytics/trends
func (h *ReportingHandler) GetTrends(c *gin.Context) {
    // TODO: Get trend data for charts
}
```

#### 6. Routes (`/api/internal/api/routes.go` - add these)

```go
// Reporting routes
reportingRepo := repository.NewReportingRepository(db)
reportingGenerator := reporting.NewGenerator(execRepo, reportingRepo, logger)
reportingHandler := handlers.NewReportingHandler(reportingRepo, reportingGenerator, logger)

v1.POST("/reports/generate", reportingHandler.Generate)
v1.GET("/reports", reportingHandler.List)
v1.GET("/reports/:id", reportingHandler.Get)
v1.GET("/reports/:id/download", reportingHandler.Download)
v1.GET("/analytics/metrics", reportingHandler.GetMetrics)
v1.GET("/analytics/flakiness", reportingHandler.GetFlakiness)
v1.GET("/analytics/trends", reportingHandler.GetTrends)
```

### Frontend Implementation

#### 1. Types (`/web/lib/api/types.ts` - add)

```typescript
export interface DailyMetric {
  id: string;
  date: string;
  flow_id?: string;
  suite?: string;
  total_executions: number;
  passed_executions: number;
  failed_executions: number;
  avg_duration_ms?: number;
  p50_duration_ms?: number;
  p95_duration_ms?: number;
  p99_duration_ms?: number;
}

export interface FlakinessMetric {
  id: string;
  flow_id: string;
  window_start: string;
  window_end: string;
  total_runs: number;
  passed_runs: number;
  failed_runs: number;
  flakiness_score: number;
  consecutive_failures: number;
  consecutive_successes: number;
}

export interface Report {
  id: string;
  report_type: 'html' | 'junit' | 'json';
  title: string;
  file_path: string;
  file_size_bytes?: number;
  created_at: string;
  expires_at?: string;
}
```

#### 2. Pages to Create

- `/web/app/analytics/page.tsx` - Analytics dashboard with charts
- `/web/app/analytics/trends/page.tsx` - Trend analysis
- `/web/app/analytics/flakiness/page.tsx` - Flaky tests view
- `/web/app/reports/page.tsx` - Reports list
- `/web/app/reports/generate/page.tsx` - Report generation form

#### 3. Components to Create

- `/web/components/analytics/TrendChart.tsx` - Line chart for trends
- `/web/components/analytics/FlakinessTable.tsx` - Flaky tests table
- `/web/components/analytics/MetricsCard.tsx` - Metric display card
- `/web/components/reports/ReportGenerator.tsx` - Report generation form

### Implementation Checklist

- [ ] Create database schema and run migration
- [ ] Implement models and repository
- [ ] Implement aggregation service
- [ ] Implement report generator
- [ ] Create API handlers and routes
- [ ] Implement frontend types and API client
- [ ] Create analytics pages
- [ ] Create report pages
- [ ] Add charts library (recharts)
- [ ] Implement background aggregation job
- [ ] Add report cleanup job
- [ ] Test end-to-end

---

## Phase 4: AI Integration

**Timeline**: 6 weeks
**Complexity**: High
**Priority**: Medium-High (Innovation)

### Overview
Enable AI-powered test creation, import from specifications, coverage analysis, and self-healing suggestions. Highest risk phase due to external API dependencies and cost management requirements.

### Database Schema

**New schema**: `ai`

**Tables to create**:

```sql
-- AI generation history
CREATE TABLE IF NOT EXISTS ai.generation_history (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    prompt TEXT NOT NULL,
    provider VARCHAR(50) NOT NULL, -- 'anthropic', 'openai', 'local'
    model VARCHAR(100) NOT NULL,
    generated_flow_id UUID REFERENCES flows.flows(id) ON DELETE SET NULL,
    tokens_used INTEGER,
    cost_cents INTEGER,
    generation_time_ms BIGINT,
    status VARCHAR(50) NOT NULL, -- 'success', 'failed', 'partial'
    error_message TEXT,
    created_by VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_generation_created ON ai.generation_history(created_at);
CREATE INDEX idx_generation_provider ON ai.generation_history(provider);

-- AI suggestions (self-healing)
CREATE TABLE IF NOT EXISTS ai.suggestions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    suggestion_type VARCHAR(50) NOT NULL, -- 'fix', 'optimization', 'coverage'
    execution_id UUID REFERENCES executions.executions(id) ON DELETE CASCADE,
    flow_id UUID REFERENCES flows.flows(id) ON DELETE CASCADE,
    step_id VARCHAR(255),
    title VARCHAR(500) NOT NULL,
    description TEXT NOT NULL,
    suggested_changes JSONB NOT NULL,
    confidence_score DECIMAL(5,2) NOT NULL, -- 0-100
    status VARCHAR(50) NOT NULL DEFAULT 'pending', -- 'pending', 'accepted', 'rejected', 'applied'
    applied_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_suggestions_execution ON ai.suggestions(execution_id);
CREATE INDEX idx_suggestions_flow ON ai.suggestions(flow_id);
CREATE INDEX idx_suggestions_status ON ai.suggestions(status);

-- Import history
CREATE TABLE IF NOT EXISTS ai.import_history (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    source_type VARCHAR(50) NOT NULL, -- 'openapi', 'postman', 'pact'
    source_file VARCHAR(1000) NOT NULL,
    flows_generated INTEGER NOT NULL DEFAULT 0,
    flows_created TEXT[], -- Array of flow IDs
    import_time_ms BIGINT,
    status VARCHAR(50) NOT NULL,
    error_message TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_import_created ON ai.import_history(created_at);

-- Coverage analysis
CREATE TABLE IF NOT EXISTS ai.coverage_analysis (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    analysis_type VARCHAR(50) NOT NULL, -- 'openapi', 'api_catalog'
    total_endpoints INTEGER NOT NULL,
    covered_endpoints INTEGER NOT NULL,
    missing_endpoints JSONB NOT NULL,
    coverage_percentage DECIMAL(5,2) NOT NULL,
    recommendations JSONB,
    analyzed_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_coverage_analyzed ON ai.coverage_analysis(analyzed_at);
```

### Backend Implementation

#### 1. AI Provider Interface (`/api/internal/ai/provider.go`)

```go
package ai

import (
    "context"
)

type Provider interface {
    GenerateFlow(ctx context.Context, prompt string) (*GeneratedFlow, error)
    AnalyzeCoverage(ctx context.Context, spec string) (*CoverageAnalysis, error)
    SuggestFix(ctx context.Context, execution *models.Execution) (*Suggestion, error)
    GetName() string
    GetModel() string
}

type GeneratedFlow struct {
    YAML         string
    TokensUsed   int
    CostCents    int
    GenerationMs int64
}

type CoverageAnalysis struct {
    TotalEndpoints   int
    CoveredEndpoints int
    MissingEndpoints []MissingEndpoint
    Recommendations  []string
}

type MissingEndpoint struct {
    Method string
    Path   string
    Tags   []string
}

type Suggestion struct {
    Title            string
    Description      string
    SuggestedChanges map[string]interface{}
    ConfidenceScore  float64
}

// Anthropic Claude provider
type AnthropicProvider struct {
    apiKey string
    model  string
}

func NewAnthropicProvider(apiKey, model string) *AnthropicProvider {
    return &AnthropicProvider{apiKey: apiKey, model: model}
}

func (p *AnthropicProvider) GenerateFlow(ctx context.Context, prompt string) (*GeneratedFlow, error) {
    // TODO: Implement Claude API call
    return nil, nil
}

// OpenAI provider
type OpenAIProvider struct {
    apiKey string
    model  string
}

func NewOpenAIProvider(apiKey, model string) *OpenAIProvider {
    return &OpenAIProvider{apiKey: apiKey, model: model}
}

func (p *OpenAIProvider) GenerateFlow(ctx context.Context, prompt string) (*GeneratedFlow, error) {
    // TODO: Implement OpenAI API call
    return nil, nil
}

// Local LLM provider (Ollama)
type LocalProvider struct {
    endpoint string
    model    string
}

func NewLocalProvider(endpoint, model string) *LocalProvider {
    return &LocalProvider{endpoint: endpoint, model: model}
}

func (p *LocalProvider) GenerateFlow(ctx context.Context, prompt string) (*GeneratedFlow, error) {
    // TODO: Implement local LLM call
    return nil, nil
}
```

#### 2. Flow Generator (`/api/internal/ai/generator.go`)

```go
package ai

import (
    "context"
    "github.com/georgi-georgiev/testmesh/internal/storage/models"
    "github.com/georgi-georgiev/testmesh/internal/storage/repository"
    "go.uber.org/zap"
)

type Generator struct {
    provider     Provider
    flowRepo     *repository.FlowRepository
    aiRepo       *repository.AIRepository
    logger       *zap.Logger
}

func NewGenerator(provider Provider, flowRepo *repository.FlowRepository, aiRepo *repository.AIRepository, logger *zap.Logger) *Generator {
    return &Generator{
        provider: provider,
        flowRepo: flowRepo,
        aiRepo:   aiRepo,
        logger:   logger,
    }
}

// GenerateFromPrompt generates a flow from natural language
func (g *Generator) GenerateFromPrompt(ctx context.Context, prompt string) (*models.Flow, error) {
    // TODO:
    // 1. Include YAML schema in system prompt
    // 2. Include example flows
    // 3. Call AI provider
    // 4. Parse and validate YAML
    // 5. Save flow
    // 6. Log generation history
    return nil, nil
}

// ImportFromOpenAPI generates flows from OpenAPI spec
func (g *Generator) ImportFromOpenAPI(ctx context.Context, openAPISpec string) ([]*models.Flow, error) {
    // TODO:
    // 1. Parse OpenAPI spec
    // 2. Generate flow for each endpoint
    // 3. Validate flows
    // 4. Save flows
    // 5. Log import history
    return nil, nil
}

// ImportFromPostman generates flows from Postman collection
func (g *Generator) ImportFromPostman(ctx context.Context, collectionJSON string) ([]*models.Flow, error) {
    // TODO: Similar to OpenAPI import
    return nil, nil
}

// ImportFromPact generates flows from Pact contracts
func (g *Generator) ImportFromPact(ctx context.Context, pactJSON string) ([]*models.Flow, error) {
    // TODO: Convert Pact interactions to flow steps
    return nil, nil
}
```

#### 3. Coverage Analyzer (`/api/internal/ai/analyzer.go`)

```go
package ai

import (
    "context"
    "github.com/georgi-georgiev/testmesh/internal/storage/repository"
    "go.uber.org/zap"
)

type Analyzer struct {
    provider   Provider
    flowRepo   *repository.FlowRepository
    aiRepo     *repository.AIRepository
    logger     *zap.Logger
}

func NewAnalyzer(provider Provider, flowRepo *repository.FlowRepository, aiRepo *repository.AIRepository, logger *zap.Logger) *Analyzer {
    return &Analyzer{
        provider: provider,
        flowRepo: flowRepo,
        aiRepo:   aiRepo,
        logger:   logger,
    }
}

// AnalyzeOpenAPICoverage compares flows against OpenAPI spec
func (a *Analyzer) AnalyzeOpenAPICoverage(ctx context.Context, openAPISpec string) (*CoverageAnalysis, error) {
    // TODO:
    // 1. Parse OpenAPI spec to get all endpoints
    // 2. Query existing flows
    // 3. Match flows to endpoints
    // 4. Calculate coverage percentage
    // 5. Identify missing endpoints
    // 6. Generate recommendations
    return nil, nil
}
```

#### 4. Self-Healing Engine (`/api/internal/ai/self_healing.go`)

```go
package ai

import (
    "context"
    "github.com/georgi-georgiev/testmesh/internal/storage/models"
    "github.com/georgi-georgiev/testmesh/internal/storage/repository"
    "go.uber.org/zap"
)

type SelfHealingEngine struct {
    provider     Provider
    execRepo     *repository.ExecutionRepository
    aiRepo       *repository.AIRepository
    logger       *zap.Logger
}

func NewSelfHealingEngine(provider Provider, execRepo *repository.ExecutionRepository, aiRepo *repository.AIRepository, logger *zap.Logger) *SelfHealingEngine {
    return &SelfHealingEngine{
        provider: provider,
        execRepo: execRepo,
        aiRepo:   aiRepo,
        logger:   logger,
    }
}

// AnalyzeFailure generates suggestions for failed executions
func (s *SelfHealingEngine) AnalyzeFailure(ctx context.Context, executionID string) (*Suggestion, error) {
    // TODO:
    // 1. Get execution details and error
    // 2. Get execution steps and logs
    // 3. Use AI to analyze failure
    // 4. Generate fix suggestion
    // 5. Save suggestion
    return nil, nil
}

// ApplySuggestion applies an accepted suggestion
func (s *SelfHealingEngine) ApplySuggestion(ctx context.Context, suggestionID string) error {
    // TODO:
    // 1. Get suggestion
    // 2. Apply changes to flow
    // 3. Mark as applied
    return nil
}
```

#### 5. API Handlers (`/api/internal/api/handlers/ai.go`)

```go
package handlers

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "github.com/georgi-georgiev/testmesh/internal/ai"
    "github.com/georgi-georgiev/testmesh/internal/storage/repository"
    "go.uber.org/zap"
)

type AIHandler struct {
    generator     *ai.Generator
    analyzer      *ai.Analyzer
    selfHealing   *ai.SelfHealingEngine
    aiRepo        *repository.AIRepository
    logger        *zap.Logger
}

func NewAIHandler(generator *ai.Generator, analyzer *ai.Analyzer, selfHealing *ai.SelfHealingEngine, aiRepo *repository.AIRepository, logger *zap.Logger) *AIHandler {
    return &AIHandler{
        generator:   generator,
        analyzer:    analyzer,
        selfHealing: selfHealing,
        aiRepo:      aiRepo,
        logger:      logger,
    }
}

// POST /api/v1/ai/generate
func (h *AIHandler) Generate(c *gin.Context) {
    // TODO: Generate flow from prompt
}

// POST /api/v1/ai/import/openapi
func (h *AIHandler) ImportOpenAPI(c *gin.Context) {
    // TODO: Import from OpenAPI
}

// POST /api/v1/ai/import/postman
func (h *AIHandler) ImportPostman(c *gin.Context) {
    // TODO: Import from Postman
}

// POST /api/v1/ai/import/pact
func (h *AIHandler) ImportPact(c *gin.Context) {
    // TODO: Import from Pact
}

// POST /api/v1/ai/coverage/analyze
func (h *AIHandler) AnalyzeCoverage(c *gin.Context) {
    // TODO: Analyze coverage
}

// GET /api/v1/ai/coverage
func (h *AIHandler) GetCoverage(c *gin.Context) {
    // TODO: Get coverage analysis results
}

// GET /api/v1/ai/suggestions
func (h *AIHandler) ListSuggestions(c *gin.Context) {
    // TODO: List self-healing suggestions
}

// POST /api/v1/ai/suggestions/:id/apply
func (h *AIHandler) ApplySuggestion(c *gin.Context) {
    // TODO: Apply suggestion
}

// POST /api/v1/ai/suggestions/:id/reject
func (h *AIHandler) RejectSuggestion(c *gin.Context) {
    // TODO: Reject suggestion
}

// GET /api/v1/ai/usage
func (h *AIHandler) GetUsage(c *gin.Context) {
    // TODO: Get token usage and costs
}
```

### Frontend Implementation

#### 1. Pages to Create

- `/web/app/ai/page.tsx` - AI hub dashboard
- `/web/app/ai/generate/page.tsx` - Interactive flow generator
- `/web/app/ai/import/page.tsx` - Import wizard
- `/web/app/ai/coverage/page.tsx` - Coverage dashboard
- `/web/app/ai/suggestions/page.tsx` - Self-healing suggestions list

#### 2. Components to Create

- `/web/components/ai/PromptInput.tsx` - Natural language input
- `/web/components/ai/GeneratedFlowPreview.tsx` - Preview generated YAML
- `/web/components/ai/ImportWizard.tsx` - Multi-step import flow
- `/web/components/ai/SuggestionCard.tsx` - Display AI suggestions
- `/web/components/ai/CoverageDashboard.tsx` - Coverage visualization

### Security & Cost Management

**CRITICAL Requirements**:

1. **API Key Management**:
   ```go
   // Store in environment variables, never in code
   ANTHROPIC_API_KEY=sk-...
   OPENAI_API_KEY=sk-...
   ```

2. **Rate Limiting**:
   ```go
   // Implement per-user rate limits
   rateLimiter := middleware.NewRateLimiter(10, time.Minute) // 10 requests/min
   ```

3. **Token Usage Tracking**:
   - Log every AI API call with token usage
   - Calculate costs based on provider pricing
   - Alert at 80% of monthly budget

4. **Cost Limits**:
   ```go
   const (
       MaxTokensPerUser   = 100000  // Per month
       MaxCostPerUser     = 5000    // $50.00 per month
       AlertThreshold     = 0.80    // 80%
   )
   ```

5. **Input Sanitization**:
   - Never include credentials in prompts
   - Sanitize user input to prevent prompt injection
   - Validate generated YAML before execution

### Implementation Checklist

- [ ] Create database schema and run migration
- [ ] Implement AI provider interface
- [ ] Implement Anthropic provider
- [ ] Implement OpenAI provider
- [ ] Implement local LLM provider (optional)
- [ ] Implement flow generator
- [ ] Implement OpenAPI/Postman/Pact importers
- [ ] Implement coverage analyzer
- [ ] Implement self-healing engine
- [ ] Create API handlers and routes
- [ ] Add rate limiting middleware
- [ ] Add token usage tracking
- [ ] Add cost monitoring
- [ ] Implement frontend pages
- [ ] Add prompt engineering (include YAML schema)
- [ ] Test with various prompts
- [ ] Security audit (prompt injection, API keys)
- [ ] Load testing
- [ ] Documentation

---

## References

- **Original Plan**: `/docs/features/PHASE_1_TO_4_IMPLEMENTATION.md` (or in plan file)
- **YAML Schema**: `/docs/features/YAML_SCHEMA.md`
- **Architecture**: `/docs/architecture/ARCHITECTURE.md`
- **Feature Specs**:
  - Mock Servers: `/docs/features/MOCK_SERVER.md`
  - Contract Testing: `/docs/features/CONTRACT_TESTING.md`
  - Advanced Reporting: `/docs/features/ADVANCED_REPORTING.md`
  - AI Integration: `/docs/features/AI_INTEGRATION.md`

## Next Steps

1. **Choose Phase**: Decide whether to implement Phase 3 or Phase 4 first
   - Recommendation: Phase 3 (lower risk, immediate business value)

2. **Create Migration**: Generate SQL migration for chosen phase

3. **Update main.go**: Initialize new repositories and services

4. **Update routes.go**: Mount new API handlers

5. **Implement incrementally**: Backend → API → Frontend → Testing

6. **Update Task List**: Create tasks for the chosen phase

---

**Note**: Both phases are independent and can be implemented in any order. Phase 3 (Reporting) has no external dependencies and lower risk. Phase 4 (AI) requires API keys and careful cost management.
