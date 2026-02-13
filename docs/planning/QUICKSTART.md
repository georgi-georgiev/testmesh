# TestMesh Quick Start Guide

## What is TestMesh?

TestMesh is a production-ready platform for writing and running end-to-end integration tests. It provides:

- **Simple YAML-based test definitions** - Write tests in human-readable format
- **Multi-protocol support** - HTTP, databases, message queues, gRPC, WebSockets, browser automation
- **Excellent observability** - Detailed logs, artifacts, metrics, and real-time dashboards
- **Easy extensibility** - Plugin system for custom actions and integrations
- **Local development** - Run tests locally before deploying
- **Production-ready** - Built for reliability, security, and scale

## Installation

### Prerequisites

- Docker and Docker Compose (for local development)
- Node.js 18+ or Go 1.21+ (for CLI tool)
- PostgreSQL, Redis, Redis Streams (provided by Docker Compose)

### Install CLI

**Using npm:**
```bash
npm install -g @testmesh/cli
```

**Using Go:**
```bash
go install github.com/testmesh/cli@latest
```

**Using Homebrew (macOS):**
```bash
brew install testmesh
```

## Getting Started

### 1. Initialize a new project

```bash
mkdir my-tests
cd my-tests
testmesh init
```

This creates:
```
my-tests/
├── .testmesh.yaml          # Configuration
├── tests/                  # Test definitions
│   └── example.yaml
├── fixtures/               # Test data
└── plugins/                # Custom plugins
```

### 2. Write your first test

Create `tests/api-test.yaml`:

```yaml
test:
  name: "API Health Check"
  suite: "smoke"
  tags: ["critical", "api"]

  steps:
    - name: "Check API health"
      action: http_request
      config:
        method: GET
        url: "https://api.example.com/health"
        assertions:
          - status_code: 200
          - json_path: "$.status == 'ok'"
```

### 3. Run the test locally

```bash
testmesh run tests/api-test.yaml
```

Output:
```
✓ API Health Check (1.2s)
  ✓ Check API health (1.2s)

1 test passed, 0 failed
```

### 4. Set up local server (optional)

For the full experience with dashboard and scheduling:

```bash
# Start all services
testmesh server start

# Access dashboard at http://localhost:3000
```

### 5. Push tests to server

```bash
testmesh push --env local
```

## Example Tests

### HTTP API Test with Data Extraction

```yaml
test:
  name: "Create and Retrieve User"
  suite: "users"

  steps:
    - name: "Create new user"
      action: http_request
      config:
        method: POST
        url: "${API_URL}/users"
        headers:
          Authorization: "Bearer ${API_TOKEN}"
        body:
          name: "John Doe"
          email: "john@example.com"
        assertions:
          - status_code: 201
          - json_path: "$.id exists"
      save:
        user_id: "$.id"

    - name: "Retrieve created user"
      action: http_request
      config:
        method: GET
        url: "${API_URL}/users/${user_id}"
        headers:
          Authorization: "Bearer ${API_TOKEN}"
        assertions:
          - status_code: 200
          - json_path: "$.name == 'John Doe'"
          - json_path: "$.email == 'john@example.com'"
```

### Database Test

```yaml
test:
  name: "Verify User in Database"
  suite: "database"

  setup:
    - action: database_query
      config:
        query: "DELETE FROM users WHERE email = 'test@example.com'"

  steps:
    - name: "Insert user"
      action: database_query
      config:
        query: "INSERT INTO users (name, email) VALUES (?, ?)"
        params: ["Test User", "test@example.com"]
        assertions:
          - affected_rows: 1

    - name: "Verify user exists"
      action: database_query
      config:
        query: "SELECT * FROM users WHERE email = ?"
        params: ["test@example.com"]
        assertions:
          - row_count: 1
          - column: "name == 'Test User'"

  teardown:
    - action: database_query
      config:
        query: "DELETE FROM users WHERE email = 'test@example.com'"
```

### Message Queue Test

```yaml
test:
  name: "Kafka Message Flow"
  suite: "messaging"

  steps:
    - name: "Publish message to Kafka"
      action: kafka_publish
      config:
        topic: "user-events"
        message:
          event_type: "user.created"
          user_id: "12345"
          timestamp: "${NOW}"

    - name: "Consume and verify message"
      action: kafka_consume
      config:
        topic: "user-events"
        timeout: 5s
        assertions:
          - json_path: "$.event_type == 'user.created'"
          - json_path: "$.user_id == '12345'"
```

### Browser Automation Test

```yaml
test:
  name: "Login Flow"
  suite: "e2e"

  steps:
    - name: "Navigate to login page"
      action: browser_navigate
      config:
        url: "https://app.example.com/login"

    - name: "Fill login form"
      action: browser_fill
      config:
        selector: "#email"
        value: "user@example.com"

    - name: "Fill password"
      action: browser_fill
      config:
        selector: "#password"
        value: "password123"

    - name: "Click login button"
      action: browser_click
      config:
        selector: "button[type='submit']"

    - name: "Wait for dashboard"
      action: browser_wait_for
      config:
        selector: ".dashboard"
        timeout: 5s

    - name: "Take screenshot"
      action: browser_screenshot
      config:
        path: "dashboard.png"
```

## CLI Commands

### Run Tests

```bash
# Run all tests
testmesh run

# Run specific test
testmesh run tests/api-test.yaml

# Run by suite
testmesh run --suite smoke

# Run by tag
testmesh run --tag critical

# Run with specific environment
testmesh run --env staging

# Watch mode (re-run on changes)
testmesh watch
```

### Manage Tests

```bash
# Validate test syntax
testmesh validate tests/**/*.yaml

# List all tests
testmesh list

# Push tests to server
testmesh push --env staging

# Pull tests from server
testmesh pull --env production
```

### View Results

```bash
# View recent results
testmesh results

# View last 10 results
testmesh results --last 10

# View failed tests only
testmesh results --failed-only

# View specific execution
testmesh results --execution-id abc123

# View logs
testmesh logs --execution-id abc123
```

### Generate Tests

```bash
# Generate API test from OpenAPI spec
testmesh generate api-test \
  --openapi ./swagger.yaml \
  --endpoint "/api/users"

# Generate database test
testmesh generate db-test \
  --table users \
  --operation crud

# Generate browser test (with recording)
testmesh generate browser-test \
  --url "https://app.example.com" \
  --record
```

### Debug Tests

```bash
# Debug specific test
testmesh debug tests/api-test.yaml

# Debug with breakpoints
testmesh debug tests/api-test.yaml --step 3

# Interactive mode
testmesh interactive
```

## Configuration

### `.testmesh.yaml`

```yaml
version: 1
project: my-tests

environments:
  local:
    api_url: "http://localhost:3000"
    database_url: "postgresql://localhost:5432/testdb"

  staging:
    api_url: "https://api.staging.example.com"
    database_url: "${STAGING_DB_URL}"

  production:
    api_url: "https://api.example.com"
    database_url: "${PROD_DB_URL}"

defaults:
  environment: local
  timeout: 30s
  retry: 3
  retry_delay: 1s

plugins:
  - name: slack-notifier
    config:
      webhook_url: "${SLACK_WEBHOOK_URL}"
```

### Environment Variables

Create `.env` file:

```bash
# API Configuration
API_TOKEN=your_api_token_here

# Database
DATABASE_URL=postgresql://localhost:5432/testdb

# Kafka
KAFKA_BROKERS=localhost:9092

# Slack
SLACK_WEBHOOK_URL=https://hooks.slack.com/...
```

## Deployment

### Docker Compose (Local/Development)

```bash
# Start all services (server + worker + infrastructure)
docker-compose up -d

# View logs
docker-compose logs -f testmesh testmesh-worker

# Stop services
docker-compose down
```

**Architecture**: Single Docker image runs in two modes:
- `testmesh` service: HTTP server (port 5016)
- `testmesh-worker` service: Background job processor

### Kubernetes (Production)

```bash
# Install using Helm
helm repo add testmesh https://charts.testmesh.io
helm install testmesh testmesh/testmesh \
  --set postgres.password=yourpassword \
  --set redis.password=yourpassword

# Check status
kubectl get pods -n testmesh

# Access dashboard
kubectl port-forward -n testmesh svc/testmesh-server 5016:5016
kubectl port-forward -n testmesh svc/testmesh-dashboard 3000:3000
```

**Deployments**:
- `testmesh-server`: HTTP API (3 replicas)
- `testmesh-worker`: Job processor (5-20 replicas, auto-scaling)
- `testmesh-dashboard`: Next.js UI (2 replicas)

### Terraform (Infrastructure)

```bash
cd infrastructure/terraform
terraform init
terraform plan
terraform apply
```

## CI/CD Integration

### GitHub Actions

```yaml
name: Run Tests

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3

      - name: Install TestMesh CLI
        run: npm install -g @testmesh/cli

      - name: Run smoke tests
        run: testmesh run --suite smoke --env staging
        env:
          TESTMESH_API_TOKEN: ${{ secrets.TESTMESH_API_TOKEN }}
```

### GitLab CI

```yaml
test:
  stage: test
  image: testmesh/cli:latest
  script:
    - testmesh run --suite smoke --env $CI_ENVIRONMENT_NAME
  only:
    - merge_requests
    - main
```

## Web Dashboard

Access the dashboard at `http://localhost:3000` (local) or your deployed URL.

### Features:
- **Dashboard**: Overview of test executions and metrics
- **Tests**: Browse and manage test definitions
- **Executions**: View execution history and details
- **Results**: Detailed test results with artifacts
- **Schedules**: Manage scheduled test runs
- **Analytics**: Trends, metrics, and insights
- **Settings**: Configuration and user management

## Plugins

### Using Plugins

```yaml
# .testmesh.yaml
plugins:
  - name: slack-notifier
    version: 1.0.0
    config:
      webhook_url: "${SLACK_WEBHOOK_URL}"
      notify_on_failure: true
```

### Creating Custom Plugins

```typescript
// plugins/my-plugin/index.ts
import { TestMeshPlugin, ActionHandler } from '@testmesh/sdk';

const myAction: ActionHandler = {
  name: 'my_custom_action',

  async execute(config, context) {
    // Your custom logic here
    return {
      success: true,
      output: { result: 'custom action executed' }
    };
  }
};

export default {
  metadata: {
    name: 'my-plugin',
    version: '1.0.0',
    description: 'My custom plugin',
    author: 'Your Name'
  },
  actions: [myAction]
} as TestMeshPlugin;
```

Load plugin:
```bash
testmesh plugin install ./plugins/my-plugin
```

## Best Practices

1. **Organize tests by suite and tags**
   - Use suites for logical grouping (e.g., "auth", "payments")
   - Use tags for cross-cutting concerns (e.g., "critical", "slow")

2. **Use setup/teardown for test isolation**
   - Clean up test data after each test
   - Create fresh fixtures in setup

3. **Extract reusable data into variables**
   - Use environment variables for configuration
   - Save response data for reuse in later steps

4. **Name tests and steps descriptively**
   - Make it clear what's being tested
   - Make failures easy to understand

5. **Keep tests independent**
   - Each test should run standalone
   - Don't rely on execution order

6. **Use assertions liberally**
   - Verify all important aspects
   - Fail fast with clear error messages

7. **Capture artifacts on failure**
   - Enable screenshots for browser tests
   - Log request/response bodies

8. **Monitor flaky tests**
   - Use the dashboard to identify flaky tests
   - Fix or remove unreliable tests

## Troubleshooting

### Test fails with timeout

Increase timeout in test definition:
```yaml
test:
  timeout: 60s  # Increase from default 30s
```

### Can't connect to database

Check database URL and credentials:
```bash
testmesh config get database_url
```

### Plugin not found

Install plugin:
```bash
testmesh plugin install <plugin-name>
```

### Logs not showing

Enable debug logging:
```bash
testmesh run --log-level debug
```

## Next Steps

- Read the [User Guide](./docs/USER_GUIDE.md) for detailed documentation
- Explore [Example Tests](./examples/) for more inspiration
- Join our [Community](https://discord.gg/testmesh) for support
- Check out the [Plugin Development Guide](./docs/PLUGIN_DEVELOPMENT.md)
- Read the [API Documentation](./docs/API.md)

## Support

- **Documentation**: https://docs.testmesh.io
- **Discord**: https://discord.gg/testmesh
- **GitHub Issues**: https://github.com/testmesh/testmesh/issues
- **Email**: support@testmesh.io

## License

TestMesh is open source software licensed under the MIT license.
