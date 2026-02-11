# TestMesh v1.0 - Complete Technology Stack

> **Comprehensive review of all technologies, frameworks, and libraries**

**Version**: 1.0
**Date**: 2026-02-11
**Purpose**: Review and approval before implementation

---

## Overview

TestMesh v1.0 uses a **modular monolith** architecture with:
- **Backend**: Single Go binary (modular design)
- **Frontend**: Next.js 14 + React 18 + TypeScript
- **CLI**: Go with Cobra
- **Database**: PostgreSQL 14+ with TimescaleDB
- **Cache/Queue**: Redis 7+

---

## 1. Backend Stack (Go)

### 1.1 Core Language & Runtime

| Technology | Version | Purpose |
|------------|---------|---------|
| **Go** | 1.21+ | Primary backend language |
| **CGO** | Disabled | For static binary compilation |

**Why Go?**
- ✅ Compiled language → better performance
- ✅ Built-in concurrency (goroutines) → parallel test execution
- ✅ Static typing → catch errors at compile time
- ✅ Single binary deployment → no runtime dependencies
- ✅ Fast startup → critical for CLI responsiveness
- ✅ Low resource usage → cost-effective scaling

---

### 1.2 Web Framework & HTTP

| Library | Version | Purpose | Import Path |
|---------|---------|---------|-------------|
| **Gin** | Latest | HTTP web framework, routing, middleware | `github.com/gin-gonic/gin` |
| **Gorilla WebSocket** | Latest | WebSocket support for real-time updates | `github.com/gorilla/websocket` |

**Why Gin?**
- Fast HTTP router (most performant Go framework)
- Middleware support (auth, logging, rate limiting)
- JSON validation built-in
- Good documentation and community

---

### 1.3 Database & ORM

| Library | Version | Purpose | Import Path |
|---------|---------|---------|-------------|
| **GORM** | Latest | ORM for database operations | `gorm.io/gorm` |
| **GORM PostgreSQL Driver** | Latest | PostgreSQL driver for GORM | `gorm.io/driver/postgres` |
| **jackc/pgx** | v5 | PostgreSQL driver (alternative, lower-level) | `github.com/jackc/pgx/v5` |
| **lib/pq** | Latest | PostgreSQL driver (for array support) | `github.com/lib/pq` |
| **golang-migrate** | Latest | Database migrations | `github.com/golang-migrate/migrate/v4` |

**Database Choices**:
- **Primary ORM**: GORM (easier development, good for CRUD)
- **Performance-critical queries**: jackc/pgx (when GORM overhead is too much)
- **Migrations**: golang-migrate (standard tool)

---

### 1.4 Redis & Caching

| Library | Version | Purpose | Import Path |
|---------|---------|---------|-------------|
| **go-redis** | v9 | Redis client for caching & streams | `github.com/redis/go-redis/v9` |

**Redis Usage**:
- Caching (sessions, flow definitions)
- Distributed locking
- Redis Streams for message queue

**Why Redis Streams over RabbitMQ?**
- ✅ Already using Redis for caching → reuse infrastructure
- ✅ Simpler deployment (one less service)
- ✅ Redis Streams has persistence, consumer groups, acknowledgments
- ✅ Good enough for small-medium workloads (v1.0 scale)

---

### 1.5 Authentication & Security

| Library | Version | Purpose | Import Path |
|---------|---------|---------|-------------|
| **golang-jwt/jwt** | v5 | JWT token generation/validation | `github.com/golang-jwt/jwt/v5` |
| **bcrypt** | Latest | Password hashing | `golang.org/x/crypto/bcrypt` |
| **crypto/rand** | stdlib | Secure random generation | `crypto/rand` |

---

### 1.6 Configuration & CLI

| Library | Version | Purpose | Import Path |
|---------|---------|---------|-------------|
| **Cobra** | Latest | CLI framework | `github.com/spf13/cobra` |
| **Viper** | Latest | Configuration management | `github.com/spf13/viper` |

**Why Cobra + Viper?**
- Industry standard for Go CLI tools
- Used by kubectl, Hugo, Docker CLI
- Excellent flag parsing, subcommands, help generation

---

### 1.7 Logging & Observability

| Library | Version | Purpose | Import Path |
|---------|---------|---------|-------------|
| **Zap** | Latest | Structured logging | `go.uber.org/zap` |
| **Prometheus client** | Latest | Metrics collection | `github.com/prometheus/client_golang/prometheus` |
| **OpenTelemetry** | Latest | Distributed tracing | `go.opentelemetry.io/otel` |
| **Jaeger Exporter** | Latest | Trace export to Jaeger | `go.opentelemetry.io/otel/exporters/jaeger` |

**Logging Strategy**:
- Structured JSON logs (production)
- Human-readable logs (development)
- Log levels: debug, info, warn, error

**Metrics Strategy**:
- Prometheus for metrics collection
- Grafana for visualization
- Custom metrics for test executions, duration, success rate

---

### 1.8 Testing & Automation

| Library | Version | Purpose | Import Path |
|---------|---------|---------|-------------|
| **Playwright (Go)** | Latest | Browser automation | `github.com/playwright-community/playwright-go` |
| **Testify** | Latest | Test assertions & mocking | `github.com/stretchr/testify` |
| **gomock** | Latest | Mock generation | `github.com/golang/mock` |

**Why Playwright?**
- Cross-browser support (Chrome, Firefox, Safari)
- Headless/headful modes
- Network interception
- Screenshot/video recording
- Better than Selenium (more modern, faster, better API)

---

### 1.9 Utilities

| Library | Version | Purpose | Import Path |
|---------|---------|---------|-------------|
| **uuid** | Latest | UUID generation | `github.com/google/uuid` |
| **retry-go** | Latest | Retry logic with backoff | `github.com/avast/retry-go/v4` |
| **ojg** | Latest | JSONPath implementation | `github.com/ohler55/ojg` |
| **validator** | v10 | Struct validation | `github.com/go-playground/validator/v10` |
| **color** | Latest | Colored terminal output | `github.com/fatih/color` |
| **progressbar** | Latest | CLI progress bars | `github.com/schollz/progressbar/v3` |
| **tablewriter** | Latest | ASCII table rendering | `github.com/olekukonko/tablewriter` |

---

### 1.10 Protocol Support (Built-in Handlers)

| Library | Version | Purpose | Import Path |
|---------|---------|---------|-------------|
| **net/http** | stdlib | HTTP client | `net/http` |
| **grpc-go** | v1.60+ | gRPC client/server | `google.golang.org/grpc` |
| **sarama** | Latest | Kafka client (producer/consumer) | `github.com/IBM/sarama` |
| **gorilla/websocket** | Latest | WebSocket client | `github.com/gorilla/websocket` |

**Note**: All these are **built-in handlers**, not external plugins!

---

## 2. Frontend Stack (TypeScript/React)

### 2.1 Core Framework

| Technology | Version | Purpose |
|------------|---------|---------|
| **Next.js** | 14.1.0+ | React meta-framework (App Router) |
| **React** | 18.2.0+ | UI library |
| **TypeScript** | 5.3.3+ | Type safety |
| **Turbopack** | Latest | Fast bundler (Next.js 14) |

**Why Next.js 14?**
- ✅ App Router (new paradigm, better DX)
- ✅ Server Components (better performance)
- ✅ Turbopack (faster dev builds)
- ✅ Built-in API routes (optional)
- ✅ Static + SSR + ISR flexibility

---

### 2.2 State Management & Data Fetching

| Library | Version | Purpose | NPM Package |
|---------|---------|---------|-------------|
| **TanStack Query** | 5.17.0+ | Server state management, caching | `@tanstack/react-query` |
| **Zustand** | 4.5.0+ | Client state management | `zustand` |
| **Axios** | 1.6.5+ | HTTP client | `axios` |

**Why TanStack Query (React Query)?**
- Best-in-class data fetching/caching
- Automatic background refetching
- Optimistic updates
- Pagination/infinite scroll support

**Why Zustand?**
- Lightweight (vs Redux)
- Simple API
- No boilerplate
- Works great with React Query

---

### 2.3 UI Components & Styling

| Library | Version | Purpose | NPM Package |
|---------|---------|---------|-------------|
| **Tailwind CSS** | 3.4.0+ | Utility-first CSS | `tailwindcss` |
| **tailwindcss-animate** | 1.0.7+ | Animation utilities | `tailwindcss-animate` |
| **Radix UI** | Latest | Headless UI components | `@radix-ui/react-*` |
| **shadcn/ui** | Latest | Pre-built components (Radix + Tailwind) | Copy/paste components |
| **Lucide React** | 0.309.0+ | Icon library | `lucide-react` |
| **class-variance-authority** | 0.7.0+ | Component variants | `class-variance-authority` |
| **clsx** | 2.1.0+ | Conditional classnames | `clsx` |
| **tailwind-merge** | 2.2.0+ | Merge Tailwind classes | `tailwind-merge` |

**UI Strategy**:
- **shadcn/ui** for base components (Button, Input, Dialog, etc.)
- **Radix UI** for accessibility and keyboard navigation
- **Tailwind** for styling (no custom CSS files)

---

### 2.4 Specialized Components

| Library | Version | Purpose | NPM Package |
|---------|---------|---------|-------------|
| **React Flow** | 11.10.0+ | Visual flow builder (node editor) | `reactflow` |
| **Monaco Editor** | 0.45.0+ | Code editor (VS Code engine) | `monaco-editor` |
| **@monaco-editor/react** | 4.6.0+ | React wrapper for Monaco | `@monaco-editor/react` |
| **Recharts** | 2.10.0+ | Charts and graphs | `recharts` |
| **date-fns** | 3.2.0+ | Date manipulation | `date-fns` |
| **cmdk** | 0.2.0+ | Command palette (⌘K) | `cmdk` |

**Why React Flow?**
- Best visual node editor for React
- Drag-and-drop flows
- Custom nodes/edges
- Zoom/pan
- Mini-map support

**Why Monaco Editor?**
- Same editor as VS Code
- Syntax highlighting
- IntelliSense support
- Diff viewer

**Why Recharts?**
- Composable chart components
- Built on D3
- Responsive
- Good TypeScript support

---

### 2.5 Forms & Validation

| Library | Version | Purpose | NPM Package |
|---------|---------|---------|-------------|
| **React Hook Form** | 7.49.0+ | Form state management | `react-hook-form` |
| **Zod** | 3.22.0+ | Schema validation | `zod` |
| **@hookform/resolvers** | 3.3.0+ | Zod + React Hook Form integration | `@hookform/resolvers` |

**Why React Hook Form?**
- Best performance (uncontrolled inputs)
- Minimal re-renders
- Easy validation
- Great DX

**Why Zod?**
- TypeScript-first validation
- Type inference
- Composable schemas
- Runtime type safety

---

### 2.6 Real-Time Communication

| Library | Version | Purpose | NPM Package |
|---------|---------|---------|-------------|
| **Socket.io Client** | 4.6.0+ | WebSocket client (real-time updates) | `socket.io-client` |

**Alternative**: Native WebSocket API (stdlib)

**Usage**:
- Real-time execution updates
- Live log streaming
- Multi-user collaboration (future)

---

### 2.7 Development Tools

| Library | Version | Purpose | NPM Package |
|---------|---------|---------|-------------|
| **ESLint** | 8.56.0+ | Linting | `eslint` |
| **eslint-config-next** | 14.1.0+ | Next.js ESLint config | `eslint-config-next` |
| **Prettier** | 3.2.0+ | Code formatting | `prettier` |
| **prettier-plugin-tailwindcss** | 0.5.0+ | Tailwind class sorting | `prettier-plugin-tailwindcss` |
| **Vitest** | 1.2.0+ | Unit testing | `vitest` |
| **@testing-library/react** | 14.1.2+ | React component testing | `@testing-library/react` |
| **Playwright** | 1.41.0+ | E2E testing | `@playwright/test` |

---

## 3. Database & Storage

### 3.1 Primary Database

| Technology | Version | Purpose |
|------------|---------|---------|
| **PostgreSQL** | 14+ | Primary relational database |
| **TimescaleDB** | Latest | Time-series extension for metrics |

**Why PostgreSQL?**
- ✅ ACID compliant
- ✅ JSON/JSONB support (flexible schema)
- ✅ Array types (tags storage)
- ✅ Full-text search
- ✅ Extensions (TimescaleDB, pgvector for future AI)
- ✅ Proven at scale

**Why TimescaleDB?**
- ✅ Built on PostgreSQL (familiar SQL)
- ✅ Automatic partitioning for time-series data
- ✅ Compression (80%+ savings)
- ✅ Continuous aggregates
- ✅ Perfect for execution metrics

**Schema Organization**:
- Separate schemas per domain (`flows`, `executions`, `scheduler`, `agents`, `users`)
- Easy migration to microservices later

---

### 3.2 Caching Layer

| Technology | Version | Purpose |
|------------|---------|---------|
| **Redis** | 7+ | In-memory cache, sessions, distributed locks |

**Usage**:
- Session storage
- Flow definition caching
- Distributed locking (prevent duplicate executions)
- Rate limiting

---

### 3.3 Message Queue

| Technology | Version | Purpose |
|------------|---------|---------|
| **Redis Streams** | 7+ (built into Redis) | Message queue for async jobs |

**Why NOT RabbitMQ/Kafka for v1.0?**
- Redis Streams is "good enough" for v1.0 scale
- Reuse existing Redis infrastructure
- Simpler operations (one less service to manage)
- Has persistence, consumer groups, acknowledgments

**Can upgrade to RabbitMQ/Kafka in v2.0 if needed.**

---

### 3.4 Artifact Storage

| Technology | Version | Purpose |
|------------|---------|---------|
| **S3** | Latest API | Screenshots, logs, videos, HAR files |
| **MinIO** | Latest | Self-hosted S3-compatible storage (for on-prem) |

**Strategy**:
- Use S3-compatible API
- Works with AWS S3, MinIO, DigitalOcean Spaces, etc.
- Store references in PostgreSQL, files in S3

---

## 4. Infrastructure & DevOps

### 4.1 Containerization

| Technology | Version | Purpose |
|------------|---------|---------|
| **Docker** | Latest | Containerization |
| **Docker Compose** | Latest | Local development orchestration |

**Docker Images**:
- `golang:1.21-alpine` (build stage)
- `alpine:latest` (runtime stage)
- Multi-stage builds for small images (~20MB final)

---

### 4.2 Orchestration

| Technology | Version | Purpose |
|------------|---------|---------|
| **Kubernetes** | 1.28+ | Container orchestration (production) |
| **Helm** | 3+ | Kubernetes package manager |

**Kubernetes Components**:
- Deployments: `testmesh-server`, `testmesh-worker`, `testmesh-dashboard`
- StatefulSets: `postgres`, `redis`
- Services: LoadBalancer, ClusterIP
- ConfigMaps & Secrets
- HorizontalPodAutoscaler (worker scaling)
- PersistentVolumeClaims

---

### 4.3 Infrastructure as Code

| Technology | Version | Purpose |
|------------|---------|---------|
| **Terraform** | Latest | Cloud infrastructure provisioning |

**Cloud Providers**:
- AWS (primary)
- GCP (supported)
- Azure (supported)
- Self-hosted (supported)

---

### 4.4 CI/CD

**Supported Platforms**:
- GitHub Actions
- GitLab CI
- Jenkins
- CircleCI
- Azure Pipelines

**Pipeline Stages**:
1. Lint (golangci-lint, ESLint)
2. Test (go test, Vitest)
3. Build (Docker images)
4. Security scan (Trivy, npm audit)
5. Deploy (Helm)

---

## 5. Observability Stack

### 5.1 Metrics

| Technology | Version | Purpose |
|------------|---------|---------|
| **Prometheus** | Latest | Metrics collection & storage |
| **Grafana** | Latest | Metrics visualization |

**Metrics Exported**:
- Test executions (total, success rate)
- Test duration (histogram)
- Active tests (gauge)
- Queue depth
- HTTP request rate/latency

---

### 5.2 Logging

| Technology | Version | Purpose |
|------------|---------|---------|
| **Zap (Go)** | Latest | Structured logging |
| **Loki** | Latest (optional) | Log aggregation |

**Log Format**:
- JSON (production)
- Human-readable (development)

---

### 5.3 Tracing

| Technology | Version | Purpose |
|------------|---------|---------|
| **OpenTelemetry** | Latest | Distributed tracing standard |
| **Jaeger** | Latest | Trace visualization |

**Trace Context**:
- Test execution spans
- HTTP request spans
- Database query spans
- External API call spans

---

### 5.4 Error Tracking (Optional)

| Technology | Version | Purpose |
|------------|---------|---------|
| **Sentry** | Latest | Error tracking & monitoring |

**Note**: Built-in error handling sufficient for v1.0, Sentry optional.

---

## 6. Testing Tools

### 6.1 Backend Testing

| Tool | Purpose |
|------|---------|
| **go test** | Standard Go testing |
| **testify/assert** | Test assertions |
| **testify/mock** | Mocking |
| **gomock** | Mock generation |
| **httptest** | HTTP handler testing |

**Coverage Target**: >80%

---

### 6.2 Frontend Testing

| Tool | Purpose |
|------|---------|
| **Vitest** | Unit tests (faster than Jest) |
| **@testing-library/react** | Component testing |
| **Playwright** | E2E testing |

**Coverage Target**: >80%

---

## 7. Security Tools

| Tool | Purpose |
|------|---------|
| **golangci-lint** | Go static analysis |
| **gosec** | Go security scanner |
| **Trivy** | Container vulnerability scanning |
| **npm audit** | NPM dependency scanning |
| **OWASP Dependency-Check** | Dependency vulnerability scanning |

---

## 8. Documentation Tools

| Tool | Purpose |
|------|---------|
| **godoc** | Go documentation generation |
| **Swagger/OpenAPI** | REST API documentation |
| **Storybook** (optional) | UI component documentation |

---

## 9. Notable Exclusions & Alternatives Considered

### ❌ NOT Using

| Technology | Why NOT? |
|------------|----------|
| **RabbitMQ** | Redis Streams sufficient for v1.0 scale |
| **Kafka** | Overkill for v1.0, Redis Streams enough |
| **GraphQL** | REST API simpler for v1.0 |
| **Redux** | Zustand simpler, less boilerplate |
| **Sass/Less** | Tailwind CSS sufficient |
| **Emotion/Styled Components** | Tailwind CSS chosen |
| **ElasticSearch** | PostgreSQL full-text search sufficient |
| **MongoDB** | PostgreSQL with JSONB better fit |

---

## 10. Version Pinning Strategy

### Backend (Go)
```go
// go.mod
go 1.21

require (
    github.com/gin-gonic/gin v1.9.1
    gorm.io/gorm v1.25.5
    gorm.io/driver/postgres v1.5.4
    github.com/redis/go-redis/v9 v9.3.0
    go.uber.org/zap v1.26.0
    // ... (exact versions in go.mod)
)
```

### Frontend (NPM)
```json
{
  "dependencies": {
    "next": "^14.1.0",
    "react": "^18.2.0",
    "@tanstack/react-query": "^5.17.0",
    // ... (see package.json)
  }
}
```

**Strategy**:
- Pin major versions
- Allow minor/patch updates (^)
- Lock dependencies with lockfiles
- Monthly dependency updates
- Security patches immediately

---

## 11. Summary by Category

### Languages
- **Go** 1.21+ (backend, CLI)
- **TypeScript** 5.3+ (frontend)
- **SQL** (PostgreSQL 14+)
- **YAML** (configuration, test definitions)

### Frameworks
- **Gin** (Go web framework)
- **Next.js 14** (React framework)
- **Cobra** (CLI framework)

### Databases
- **PostgreSQL 14+** with TimescaleDB (primary)
- **Redis 7+** (cache, queue)

### Frontend Libraries (15 key libraries)
- React, Next.js, TypeScript
- TanStack Query, Zustand
- Tailwind, Radix UI, shadcn/ui
- React Flow, Monaco Editor, Recharts
- React Hook Form, Zod
- Socket.io, date-fns

### Backend Libraries (20+ key libraries)
- Gin, GORM, go-redis
- Cobra, Viper
- Zap, Prometheus, OpenTelemetry
- Playwright, Testify
- JWT, bcrypt
- UUID, retry-go, ojg

### Infrastructure
- Docker, Kubernetes, Helm
- Terraform
- Prometheus, Grafana, Jaeger

---

## 12. License Compatibility

**All dependencies are permissively licensed:**
- ✅ MIT License (most libraries)
- ✅ Apache 2.0 (some Go libraries)
- ✅ BSD 3-Clause (some libraries)

**No GPL/AGPL dependencies** (allows commercial use without restrictions)

---

## 13. Performance Expectations

### Backend (Go)
- **Startup time**: <1 second
- **Memory**: 50-100MB base, 200MB under load
- **Request latency**: <10ms (p95)
- **Throughput**: 100-200 flows/sec per server

### Frontend (Next.js)
- **Initial load**: <2 seconds (LCP)
- **Time to Interactive**: <3 seconds
- **Bundle size**: <300KB (gzipped)

### Database (PostgreSQL)
- **Query latency**: <5ms (p95)
- **Concurrent connections**: 100-200
- **Storage**: ~100MB per 10K executions

---

## 14. Questions for Review

### 1. Message Queue Decision
**Current**: Redis Streams
**Question**: Should we use RabbitMQ/Kafka from the start?

**Recommendation**: ✅ Stick with Redis Streams for v1.0, upgrade in v2.0 if needed.

---

### 2. WebSocket vs Socket.io
**Current**: Socket.io client
**Question**: Use native WebSocket API instead?

**Recommendation**: ✅ Socket.io has reconnection, room support, better DX. Keep it.

---

### 3. Monaco Editor vs CodeMirror
**Current**: Monaco Editor
**Question**: CodeMirror 6 is lighter, should we use it?

**Recommendation**: ✅ Monaco is VS Code engine, better features. Keep it.

---

### 4. Recharts vs Chart.js
**Current**: Recharts
**Question**: Chart.js has more features, should we switch?

**Recommendation**: ✅ Recharts is more React-friendly, composable. Keep it.

---

### 5. Playwright vs Selenium
**Current**: Playwright
**Question**: Selenium is more mature, should we use it?

**Recommendation**: ✅ Playwright is faster, modern API, better DX. Keep it.

---

## 15. Approval Checklist

✅ **ALL APPROVED - 2026-02-11**

- [x] **Go 1.21+** approved as backend language
- [x] **Next.js 14 App Router** approved as frontend framework
- [x] **PostgreSQL 14+** approved as primary database
- [x] **Redis Streams** approved as message queue (vs RabbitMQ/Kafka)
- [x] **Gin** approved as Go web framework
- [x] **GORM + selective raw SQL** approved as ORM strategy
- [x] **TanStack Query** approved for data fetching
- [x] **Zustand** approved for state management
- [x] **Tailwind + shadcn/ui + Radix UI** approved for UI
- [x] **React Flow** approved for visual editor
- [x] **Monaco Editor** approved for code editor
- [x] **Recharts** approved for charts
- [x] **Playwright** approved for browser automation
- [x] **Prometheus + Grafana** approved for observability
- [x] **OpenTelemetry + Jaeger** approved for distributed tracing
- [x] **Zap + JSON logs** approved for logging
- [x] **Cobra + Viper** approved for CLI framework
- [x] **Docker + Kubernetes** approved for deployment

---

## 16. Amendments

**Status**: ✅ No amendments - all recommendations approved as-is

### Approved Decisions (2026-02-11)

All 12 key technology decisions approved without amendments:

1. ✅ Redis Streams (vs RabbitMQ/Kafka)
2. ✅ GORM + selective raw SQL (hybrid approach)
3. ✅ Next.js 14 App Router (vs Pages Router)
4. ✅ shadcn/ui + Radix UI + Tailwind (vs Material UI)
5. ✅ Playwright (vs Selenium)
6. ✅ TanStack Query + Zustand (vs Redux)
7. ✅ Monaco Editor (vs CodeMirror)
8. ✅ Recharts (vs Chart.js)
9. ✅ Prometheus + Grafana (vs SaaS solutions)
10. ✅ Zap + JSON logs (vs ELK stack)
11. ✅ OpenTelemetry + Jaeger (included in v1.0)
12. ✅ Cobra + Viper (vs alternatives)

---

## Version History

| Version | Date | Changes |
|---------|------|---------|
| 1.0 | 2026-02-11 | Initial technology summary for review |

---

**Status**: ✅ APPROVED - 2026-02-11

This is now the **locked technology stack** for v1.0 implementation.

All changes to core technologies require explicit approval via amendment process.
