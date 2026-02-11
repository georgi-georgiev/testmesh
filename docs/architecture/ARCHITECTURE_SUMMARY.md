# Architecture Summary - Modular Monolith

## 📊 **Quick Overview**

**Architecture Pattern**: Modular Monolith → Microservices (when needed)
**Deployment**: Single Go binary + Worker process
**Communication**: In-process (direct calls) + Redis Streams (async jobs)
**Database**: Single PostgreSQL with domain schemas

---

## 🏗️ **System Architecture**

```
                              ┌─────────────────────────────────────┐
                              │      TestMesh Server (Go)           │
                              │      Single Binary                  │
                              ├─────────────────────────────────────┤
                              │                                     │
                              │  ┌───────────────────────────────┐ │
External                      │  │    API Domain                 │ │
Clients                       │  │  - REST API (port 8080)       │ │
   │                          │  │  - WebSocket (real-time)      │ │
   │                          │  │  - Auth & middleware          │ │
   ▼                          │  └──────────┬────────────────────┘ │
┌──────┐                      │             │                       │
│ CLI  │────────────┐         │  ┌──────────▼────────────────────┐ │
└──────┘            │         │  │  Scheduler Domain             │ │
                    │         │  │  - Cron scheduler             │ │
┌──────────┐        │         │  │  - Job queue                  │ │
│Dashboard │────────┼─────────┼─▶│  - Worker pool                │ │
└──────────┘        │         │  └──────────┬────────────────────┘ │
                    │         │             │                       │
┌──────────┐        │         │  ┌──────────▼────────────────────┐ │
│ Agents   │────────┘         │  │  Runner Domain                │ │
└──────────┘                  │  │  - Execution engine           │ │
                              │  │  - Action handlers            │ │
                              │  │  - Assertion engine           │ │
                              │  │  - Plugin system              │ │
                              │  └──────────┬────────────────────┘ │
                              │             │                       │
                              │  ┌──────────▼────────────────────┐ │
                              │  │  Storage Domain               │ │
                              │  │  - Flow repository            │ │
                              │  │  - Execution store            │ │
                              │  │  - Metrics store              │ │
                              │  └──────────┬────────────────────┘ │
                              │             │                       │
                              │  ┌──────────▼────────────────────┐ │
                              │  │  Shared Layer                 │ │
                              │  │  - DB, Redis, Queue clients   │ │
                              │  │  - Auth, Logging, Config      │ │
                              │  └───────────────────────────────┘ │
                              └─────────────┬───────────────────────┘
                                            │
                                            ▼
                              ┌─────────────────────────────────────┐
                              │  External Infrastructure            │
                              ├─────────────────────────────────────┤
                              │  ┌──────────┐  ┌───────┐  ┌──────┐│
                              │  │PostgreSQL│  │ Redis │  │Redis Streams││
                              │  │  (DB)    │  │(Cache)│  │(Queue)││
                              │  └──────────┘  └───────┘  └──────┘│
                              └─────────────────────────────────────┘
```

---

## 📁 **Project Structure**

```
testmesh/
├── server/                       # Backend monolith (Go)
│   ├── cmd/
│   │   ├── server/main.go       # HTTP server
│   │   ├── worker/main.go       # Background worker
│   │   └── migrate/main.go      # DB migrations
│   ├── internal/
│   │   ├── api/                 # API Domain
│   │   ├── runner/              # Runner Domain
│   │   ├── scheduler/           # Scheduler Domain
│   │   ├── storage/             # Storage Domain
│   │   ├── agent/               # Agent coordination
│   │   ├── cleanup/             # Data cleanup
│   │   └── shared/              # Shared utilities
│   └── go.mod
│
├── cli/                          # CLI tool (Go)
│   ├── cmd/
│   │   ├── root.go
│   │   ├── run.go
│   │   ├── watch.go
│   │   └── ...
│   └── go.mod
│
├── web/
│   └── dashboard/                # Next.js 14 dashboard
│       ├── app/                  # App Router
│       ├── components/
│       ├── lib/
│       └── package.json
│
├── docs/                         # Docusaurus documentation
├── infrastructure/               # Deployment configs
│   ├── docker/
│   │   └── docker-compose.yaml
│   ├── kubernetes/
│   │   └── helm/
│   └── terraform/
└── docker-compose.yaml           # Development setup
```

---

## 🔄 **Request Flow Examples**

### Example 1: Run Flow via API

```
1. User → Dashboard → API Domain
   POST /api/v1/executions

2. API Domain → Runner Domain (direct call)
   executor.Execute(ctx, flow)

3. Runner Domain → Storage Domain (direct call)
   storage.SaveExecution(ctx, result)

4. Runner Domain → API Domain
   return result

5. API Domain → User
   HTTP 201 Created + result
```

### Example 2: Scheduled Flow Execution

```
1. Cron → Scheduler Domain
   trigger scheduled job

2. Scheduler Domain → Redis Streams (queue)
   publish job message

3. Worker → Redis Streams (consume)
   consume job message

4. Worker → Runner Domain (direct call)
   executor.Execute(ctx, flow)

5. Runner Domain → Storage Domain (direct call)
   storage.SaveExecution(ctx, result)

6. Storage Domain → WebSocket (via API)
   broadcast execution complete
```

---

## 🗄️ **Database Schema Organization**

```sql
-- Separate schema per domain
CREATE SCHEMA flows;        -- Storage Domain
CREATE SCHEMA executions;   -- Storage Domain
CREATE SCHEMA scheduler;    -- Scheduler Domain
CREATE SCHEMA agents;       -- Agent Domain
CREATE SCHEMA users;        -- Storage Domain

-- Example table
CREATE TABLE flows.flows (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    definition JSONB NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE executions.executions (
    id UUID PRIMARY KEY,
    flow_id UUID REFERENCES flows.flows(id),
    status VARCHAR(50) NOT NULL,
    started_at TIMESTAMPTZ NOT NULL,
    completed_at TIMESTAMPTZ
);
```

**Benefits**:
- Clean separation
- Easy to migrate to separate databases later
- Clear ownership per domain

---

## 🚀 **Deployment Models**

### **Development (Docker Compose)**

```yaml
services:
  testmesh:
    build: ./server
    ports:
      - "8080:8080"
    command: ["server"]

  testmesh-worker:
    build: ./server
    command: ["worker"]

  postgres:
    image: postgres:14

  redis:
    image: redis:6

  rabbitmq:
    image: rabbitmq:3
```

**Single command**: `docker-compose up`

### **Production (Kubernetes)**

```yaml
# API Server (3 replicas)
apiVersion: apps/v1
kind: Deployment
metadata:
  name: testmesh-server
spec:
  replicas: 3

---

# Background Workers (5 replicas)
apiVersion: apps/v1
kind: Deployment
metadata:
  name: testmesh-worker
spec:
  replicas: 5
```

**Scale independently**: More workers for heavy load, fewer API servers.

---

## 🔀 **Domain Dependencies**

```
API Domain
    ↓ (calls directly)
Scheduler Domain
    ↓ (publishes jobs)
Redis Streams Queue
    ↓ (worker consumes)
Runner Domain
    ↓ (saves results)
Storage Domain
    ↓ (uses)
Shared Layer
```

**Rules**:
- ✅ Higher layers call lower layers
- ✅ Async via queue for long-running tasks
- ❌ No circular dependencies
- ❌ Lower layers don't call higher layers

---

## ⚡ **Performance Characteristics**

### **In-Process Communication**

```go
// Direct function call - ~1-10 microseconds
result := runner.Execute(ctx, flow)
```

**vs Microservices (HTTP)**

```go
// HTTP call - ~1-10 milliseconds (1000x slower)
result := http.Post("http://runner-service/execute", flow)
```

### **Throughput**

- **Single server**: 100-200 flows/sec
- **Scale horizontally**: Add more worker replicas
- **Database**: Bottleneck will be DB, not Go server

---

## 🔄 **Future Migration to Microservices**

### **When to Split** (v2.0+)

✅ **DO split when**:
- Test Runner needs 10x more capacity than API
- Different deployment schedules required
- Team structure demands it (separate teams)
- Observability shows clear boundaries

❌ **DON'T split when**:
- "It's best practice" (premature optimization)
- Traffic is still manageable
- Team is still small

### **How to Split**

**Phase 1**: Extract Storage
```
Before: API → Storage (in-process)
After:  API → Storage Service (HTTP/gRPC)
```

**Phase 2**: Extract Runner
```
Before: API → Runner (in-process)
After:  API → Runner Service (HTTP/gRPC)
```

**Phase 3**: Extract Scheduler
```
Before: Scheduler → Runner (in-process)
After:  Scheduler → Runner Service (queue)
```

**Cost**: Each split adds:
- Network latency
- Deployment complexity
- Operational overhead
- Distributed system complexity

**Benefit**: Independent scaling and deployment

---

## 📊 **Comparison**

| Aspect | Modular Monolith (v1.0) | Microservices (v2.0+) |
|--------|-------------------------|------------------------|
| **Deployment** | Single binary | 4+ services |
| **Latency** | Microseconds | Milliseconds |
| **Debugging** | Simple (one process) | Complex (distributed) |
| **Transactions** | Easy (single DB) | Hard (distributed) |
| **Scaling** | Horizontal (replicas) | Independent per service |
| **Ops Complexity** | Low | High |
| **Development Speed** | Fast | Slower |
| **Cost** | Lower (fewer resources) | Higher |

---

## ✅ **Decision Summary**

**v1.0**: Modular Monolith
- ✅ Faster to build
- ✅ Easier to debug
- ✅ Simpler to deploy
- ✅ Better performance
- ✅ Clean boundaries for future split

**v2.0+**: Microservices (if needed)
- ✅ Independent scaling
- ✅ Independent deployment
- ✅ Technology diversity
- ❌ More complexity

**This is the pragmatic path!** Start simple, scale when necessary. 🚀

---

## 📚 **Key Documents**

- **MODULAR_MONOLITH.md** - Detailed architecture
- **DECISIONS.md** - All technical decisions
- **PROJECT_STRUCTURE.md** - File organization
- **TECH_STACK.md** - Technology choices
- **IMPLEMENTATION_PLAN.md** - Development roadmap

---

## 🎯 **Next Steps**

1. ✅ Architecture decided: Modular Monolith
2. ✅ Domains defined: API, Runner, Scheduler, Storage
3. ✅ Tech stack finalized: Go + Next.js + PostgreSQL
4. ⏳ **Ready to start implementation**

**Awaiting your signal to begin!** 🚀
