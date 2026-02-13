# TestMesh Quick Start Guide

This guide will help you get TestMesh up and running locally to test all features including real-time WebSocket updates.

## Prerequisites

- Go 1.21 or higher
- Node.js 18+ and pnpm
- PostgreSQL running (via Docker or local)

## Step 1: Start PostgreSQL

If using Docker Compose:

```bash
docker-compose up -d postgres redis
```

The database should be accessible at `localhost:5432` with:
- Database: `testmesh`
- User: `testmesh`
- Password: `testmesh_dev`

## Step 2: Start the API Server

Open a terminal and run:

```bash
cd api
go run main.go
```

You should see:
```
INFO    Starting TestMesh API server    {"port": 5016, "environment": "development"}
```

Verify it's running:
```bash
curl http://localhost:5016/health
# Should return: {"status":"ok"}
```

## Step 3: Start the Web App

Open a **new terminal** and run:

```bash
cd web
pnpm dev
```

The web app should start at http://localhost:3000

## Step 4: Test the Real-Time Updates

### Create a Flow

1. Navigate to http://localhost:3000/flows/new
2. Use this example flow:

```yaml
flow:
  name: "WebSocket Test Flow"
  description: "Test real-time updates"

steps:
  - id: step1
    name: "First step"
    action: http_request
    config:
      method: GET
      url: "https://jsonplaceholder.typicode.com/users/1"
    assert:
      - status == 200

  - id: step2
    name: "Second step"
    action: http_request
    config:
      method: GET
      url: "https://jsonplaceholder.typicode.com/posts/1"
    assert:
      - status == 200

  - id: step3
    name: "Third step"
    action: http_request
    config:
      method: GET
      url: "https://jsonplaceholder.typicode.com/comments/1"
    assert:
      - status == 200
```

3. Click "Create Flow"

### Run and Watch Live Updates

1. Click "Run" on the flow detail page
2. Click "View" to go to the execution detail page
3. You should see:
   - ✅ A green "Live" badge indicating WebSocket connection
   - ✅ Steps updating in real-time as they execute
   - ✅ Status changing from "pending" → "running" → "completed"
   - ✅ No page refresh needed!

## Troubleshooting

### WebSocket connection error

**Symptom:** Console shows "WebSocket connection error"

**Solution:** Make sure the API server is running on port 5016:
```bash
curl http://localhost:5016/health
```

If not running, start it:
```bash
cd api && go run main.go
```

### "All statuses" Select error

**Symptom:** Error about empty Select.Item value

**Solution:** This was fixed - make sure you have the latest code:
```bash
cd web && git pull
```

### Database connection failed

**Symptom:** API fails to start with database error

**Solution:** Ensure PostgreSQL is running and the database exists:
```bash
# Connect to postgres
docker exec -it testmesh-postgres-1 psql -U postgres

# Create database if needed
CREATE DATABASE testmesh;
CREATE USER testmesh WITH PASSWORD 'testmesh_dev';
GRANT ALL PRIVILEGES ON DATABASE testmesh TO testmesh;
```

### Port already in use

**Symptom:** "bind: address already in use"

**Solution:** Kill the process using the port:
```bash
# For API (port 5016)
lsof -ti:5016 | xargs kill -9

# For Web (port 3000)
lsof -ti:3000 | xargs kill -9
```

## Features to Test

### ✅ Phase 1-3 Features
- [x] Create flows with YAML
- [x] Run flows with HTTP requests
- [x] Database queries
- [x] Assertions (status, JSONPath, expressions)
- [x] Variable interpolation (${UUID}, ${TIMESTAMP}, etc.)
- [x] Setup/teardown hooks
- [x] Retry logic with exponential backoff
- [x] Execution history
- [x] Step results with outputs

### ✅ Phase 4 Features (NEW)
- [x] WebSocket real-time updates
- [x] Live execution status
- [x] Live step updates
- [x] Connection status indicator
- [x] Auto-reconnect on disconnect

## API Endpoints

All endpoints are available at `http://localhost:5016`:

**REST:**
- `GET /health` - Health check
- `POST /api/v1/flows` - Create flow
- `GET /api/v1/flows` - List flows
- `GET /api/v1/flows/:id` - Get flow
- `PUT /api/v1/flows/:id` - Update flow
- `DELETE /api/v1/flows/:id` - Delete flow
- `POST /api/v1/executions` - Run flow
- `GET /api/v1/executions` - List executions
- `GET /api/v1/executions/:id` - Get execution
- `GET /api/v1/executions/:id/steps` - Get execution steps

**WebSocket:**
- `WS /ws/executions/:id` - Real-time execution updates

## Next Steps

- Explore the example flows in `/examples`
- Read the full implementation plan in `IMPLEMENTATION_PLAN.md`
- Check phase completion docs: `PHASE3_COMPLETE.md`, `PHASE4_COMPLETE.md`
- Continue to Phase 5 for control flow actions (condition, for_each, etc.)
