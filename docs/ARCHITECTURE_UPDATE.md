# Architecture Update Summary

**Date**: 2026-02-11
**Status**: ✅ Complete

## Changes Made

### 1. Multi-Repository Architecture ✅

Changed from monorepo to **multi-repository architecture** with separate repos for Go and TypeScript:

**Before**: Single monorepo with all code
```
testmesh/ (monorepo)
├── server/ (Go)
├── web/dashboard/ (TypeScript)
├── cli/ (Go)
└── ...
```

**After**: Separate repositories
```
testmesh-server/ (Go repository)
├── cmd/
├── internal/
├── plugins/
└── ...

testmesh-dashboard/ (TypeScript repository)
├── app/ (Next.js 16)
├── components/
├── lib/
└── ...
```

### Benefits

✅ **Clear Separation**: Backend and frontend are completely independent
✅ **Independent CI/CD**: Each repo has its own pipeline optimized for its language
✅ **Faster Builds**: Changes to one repo don't trigger builds in the other
✅ **No Tool Conflicts**: No conflicts between Go and Node.js tooling
✅ **Independent Versioning**: Backend and frontend can be versioned separately
✅ **Team Scaling**: Different teams can own different repos
✅ **Flexible Deployment**: Deploy backend and frontend to different platforms

---

### 2. Next.js 16 Upgrade ✅

**Before**: Next.js 14
**After**: Next.js 16 (latest)

**Key Features of Next.js 16**:
- Latest App Router improvements
- Enhanced performance optimizations
- Better TypeScript support
- Improved Server Components
- Latest React 19 support

---

### 3. Repository Structure

#### Backend Repository: `testmesh-server`

**Language**: Go 1.22+
**Purpose**: API server, CLI, plugins, workers

```
testmesh-server/
├── cmd/
│   ├── server/          # HTTP API server
│   ├── worker/          # Background worker
│   ├── cli/             # CLI tool
│   └── migrate/         # Database migrations
├── internal/            # Private application code
│   ├── api/             # API Domain
│   ├── runner/          # Runner Domain
│   ├── scheduler/       # Scheduler Domain
│   ├── storage/         # Storage Domain
│   └── shared/          # Shared utilities
├── pkg/                 # Public libraries
├── plugins/             # Built-in plugins
├── tests/               # Tests
├── deployments/         # Docker & K8s configs
│   ├── docker/
│   ├── kubernetes/
│   └── helm/
├── go.mod
├── go.sum
└── Makefile
```

**Tech Stack**:
- Go 1.22+
- Gin (web framework)
- GORM (ORM)
- Cobra (CLI framework)
- Redis + Redis Streams
- PostgreSQL + TimescaleDB

#### Frontend Repository: `testmesh-dashboard`

**Language**: TypeScript 5.6+
**Framework**: Next.js 16

```
testmesh-dashboard/
├── app/                 # Next.js 16 App Router
│   ├── (auth)/          # Auth routes
│   ├── (dashboard)/     # Dashboard routes
│   ├── api/             # API routes
│   ├── layout.tsx       # Root layout
│   └── page.tsx         # Home page
├── components/          # React components
│   ├── ui/              # shadcn/ui components
│   ├── flow/            # Flow editor
│   ├── request/         # Request builder
│   └── collections/     # Collections
├── lib/                 # Utilities
│   ├── api/             # API client
│   ├── hooks/           # React hooks
│   ├── store/           # State (Zustand)
│   └── types/           # TypeScript types
├── public/              # Static assets
├── tests/               # Tests
├── next.config.ts       # Next.js config
├── tailwind.config.ts
├── tsconfig.json
├── package.json
└── pnpm-lock.yaml       # Using pnpm
```

**Tech Stack**:
- Next.js 16 + React 19
- TypeScript 5.6+
- pnpm (package manager)
- Tailwind CSS
- shadcn/ui (UI components)
- React Flow (visual editor)
- Monaco Editor (code editor)
- TanStack Query (server state)
- Zustand (client state)
- Socket.io-client (WebSocket)
- Recharts (charts)

---

### 4. Development Workflow

#### Backend Development

```bash
# Clone backend repo
git clone https://github.com/testmesh/testmesh-server.git
cd testmesh-server

# Install dependencies
go mod download

# Start dependencies (PostgreSQL, Redis)
docker-compose up -d

# Run migrations
make migrate-up

# Run tests
make test

# Start server
make dev-server

# Start worker
make dev-worker

# Build CLI
make build-cli
```

#### Frontend Development

```bash
# Clone frontend repo
git clone https://github.com/testmesh/testmesh-dashboard.git
cd testmesh-dashboard

# Install dependencies (using pnpm)
pnpm install

# Set up environment
cp .env.example .env.local

# Run development server
pnpm dev

# Run tests
pnpm test

# Build for production
pnpm build
```

---

### 5. Inter-Repository Communication

**API Communication**:
- Frontend calls backend via REST API (HTTP)
- Real-time updates via WebSocket
- Log streaming via Server-Sent Events

**Example API Client** (Frontend):
```typescript
// lib/api/client.ts
import axios from 'axios';

const apiClient = axios.create({
  baseURL: process.env.NEXT_PUBLIC_API_URL, // http://localhost:8080
  headers: {
    'Content-Type': 'application/json',
  },
});

// Add auth token interceptor
apiClient.interceptors.request.use((config) => {
  const token = localStorage.getItem('auth_token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

export default apiClient;
```

---

### 6. Deployment

#### Separate Deployments

Each repository has its own deployment:

**Backend Deployment**:
- Docker image: `testmesh/server:latest`
- Kubernetes with HPA (Horizontal Pod Autoscaler)
- Exposes API at `:8080`

**Frontend Deployment**:
- Docker image: `testmesh/dashboard:latest`
- Can deploy to:
  - Vercel (recommended for Next.js 16)
  - Netlify
  - Kubernetes
  - CDN with static export
- Configures `NEXT_PUBLIC_API_URL` to point to backend

#### Local Development (Docker Compose)

```yaml
# docker-compose.yml
version: '3.8'

services:
  postgres:
    image: postgres:14
    ports: ["5432:5432"]

  redis:
    image: redis:7
    ports: ["6379:6379"]

  server:
    build: ./testmesh-server
    ports: ["8080:8080"]
    environment:
      DATABASE_URL: postgresql://testmesh:testmesh@postgres:5432/testmesh
      REDIS_URL: redis://redis:6379/0

  worker:
    build: ./testmesh-server
    command: worker
    depends_on: [postgres, redis]

  dashboard:
    build: ./testmesh-dashboard
    ports: ["3000:3000"]
    environment:
      NEXT_PUBLIC_API_URL: http://localhost:8080
```

---

### 7. CI/CD Pipelines

#### Backend Pipeline (GitHub Actions)

```yaml
# .github/workflows/backend.yml
name: Backend CI/CD

on:
  push:
    branches: [main]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: '1.22'
      - run: go test ./...

  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - run: docker build -t testmesh/server:latest .
      - run: docker push testmesh/server:latest
```

#### Frontend Pipeline (GitHub Actions)

```yaml
# .github/workflows/frontend.yml
name: Frontend CI/CD

on:
  push:
    branches: [main]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: pnpm/action-setup@v2
      - uses: actions/setup-node@v3
        with:
          node-version: '20'
          cache: 'pnpm'
      - run: pnpm install
      - run: pnpm test
      - run: pnpm build

  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: vercel/deploy@v1  # Deploy to Vercel
```

---

### 8. Version Management

**Semantic Versioning** for both repos:

- Backend: `testmesh-server v1.0.0`
- Frontend: `testmesh-dashboard v1.0.0`

**Version Compatibility Matrix**:
- Document which backend versions work with which frontend versions
- Use OpenAPI spec to document API contract
- Use feature flags for gradual rollouts

---

## Updated Documentation

The following files have been updated:

1. ✅ **PROJECT_STRUCTURE.md** - Complete rewrite for multi-repo architecture
2. ✅ **TECH_STACK.md** - Updated Next.js 14 → 16, clarified multi-repo
3. ⏭️ **ARCHITECTURE.md** - Needs update (separate architecture per repo)
4. ⏭️ **IMPLEMENTATION_PLAN.md** - Needs update (Phase 1 project setup)

---

## Migration Path (If Starting Fresh)

### Option 1: Start with Multi-Repo (Recommended)

1. Create `testmesh-server` repository
2. Create `testmesh-dashboard` repository
3. Set up each independently
4. Use docker-compose for local development

### Option 2: Start Monorepo, Split Later

1. Start with monorepo structure
2. Develop both in parallel
3. Split into separate repos when needed
4. Use Git subtree/filter-branch to preserve history

---

## Trade-offs & Considerations

### Pros of Multi-Repo

✅ Independent development
✅ Clear ownership
✅ Simpler CI/CD per repo
✅ No tool conflicts
✅ Flexible deployment

### Cons of Multi-Repo

⚠️ Need to coordinate API changes
⚠️ Version compatibility management
⚠️ Shared types require separate package

### Mitigation Strategies

1. **OpenAPI Specification**: Document API contract clearly
2. **Shared Types Package**: Publish `@testmesh/types` NPM package
3. **Semantic Versioning**: Backend and frontend both use semver
4. **Feature Flags**: Gradual rollout of breaking changes
5. **E2E Tests**: Test backend + frontend integration

---

## Summary

TestMesh now uses a **multi-repository architecture** with:

1. **testmesh-server** (Go) - Backend, CLI, plugins
2. **testmesh-dashboard** (Next.js 16) - Web UI
3. **(Optional) testmesh-docs** - Shared documentation

This structure provides:
- ✅ Clear separation of concerns
- ✅ Independent development & deployment
- ✅ Latest Next.js 16 features
- ✅ Optimized CI/CD pipelines
- ✅ Better team scalability

**Status**: Ready for Phase 1 implementation! 🚀
