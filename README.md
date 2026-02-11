# TestMesh

> **A production-ready platform for writing and running e2e integration tests**

TestMesh makes it easy to write, manage, and execute end-to-end integration tests across multiple protocols and services. With excellent observability, extensibility, and developer experience, it's designed for modern distributed systems.

## Repository Structure

This is a monorepo containing all TestMesh components:

```
testmesh/
├── web/          # Next.js frontend dashboard
├── api/          # Go backend service
├── docs/         # Comprehensive documentation
└── examples/     # Example test flows and use cases
```

## Why TestMesh?

### 🎯 Simple & Powerful
Write tests in human-readable YAML. No complex coding required. Yet powerful enough to handle complex scenarios with data extraction, assertions, and flow control.

### 🔌 Multi-Protocol Support
- HTTP/REST APIs
- Databases (PostgreSQL, MySQL, MongoDB)
- Message Queues (Kafka, Redis Streams, NATS)
- gRPC services
- WebSockets
- Browser automation

### 👀 Excellent Observability
- Real-time execution dashboard
- Detailed logs with context
- Automatic artifact capture (screenshots, network traces)
- Metrics and analytics
- Flaky test detection

### 🔧 Easily Extensible
- Plugin architecture for custom actions
- Custom assertions
- Integration with external systems
- Community marketplace

### 💻 Great Developer Experience
- CLI tool for local development
- Hot reload in watch mode
- Test generators from OpenAPI/Swagger
- IDE integrations (VS Code)
- Interactive debugging mode

### 🚀 Production Ready
- Horizontal scaling
- High availability
- Security hardened
- Kubernetes native
- CI/CD integration
- Scheduled execution

### 🤖 AI-Powered Testing (NEW!)
Transform testing from tedious manual work into strategic, high-level specifications:
- **Natural Language Generation**: Describe what you want to test, AI generates the YAML flow
- **Smart Import**: Convert OpenAPI specs, Postman collections, or Pact contracts into tests automatically
- **Coverage Analysis**: AI detects gaps and generates missing tests
- **Self-Healing Tests**: AI analyzes failures and suggests fixes
- **Interactive Chat**: Conversational test creation and debugging
- **Local AI Support**: Works offline with privacy-first local models

```bash
# Generate a complete test from natural language
$ testmesh generate "Test daily fare cap with payment gateway mock"

🤖 Generating test flow...
✓ Created flows/daily-fare-cap-test.yaml
✓ Created mocks/payment-gateway-mock.yaml
✓ Created data/fare-test-data.json

Run test? [y/n]
```

**See [AI Integration Guide](./docs/AI_INTEGRATION.md) for full capabilities.**

## Quick Start

```bash
# Install CLI
npm install -g @testmesh/cli

# Initialize project
mkdir my-tests && cd my-tests
testmesh init

# Create your first test
cat > tests/api-test.yaml <<EOF
test:
  name: "API Health Check"
  steps:
    - name: "Check API health"
      action: http_request
      config:
        method: GET
        url: "https://api.example.com/health"
        assertions:
          - status_code: 200
EOF

# Run it!
testmesh run tests/api-test.yaml
```

See [QUICKSTART.md](./QUICKSTART.md) for detailed getting started guide.

## Features

### Test Definition
```yaml
test:
  name: "User Registration Flow"
  suite: "authentication"
  tags: ["critical"]

  steps:
    - name: "Create user"
      action: http_request
      config:
        method: POST
        url: "${API_URL}/users"
        body:
          email: "user@example.com"
          password: "SecurePass123!"
        assertions:
          - status_code: 201
          - json_path: "$.user.id exists"
      save:
        user_id: "$.user.id"

    - name: "Verify in database"
      action: database_query
      config:
        query: "SELECT * FROM users WHERE id = ?"
        params: ["${user_id}"]
        assertions:
          - row_count: 1
```

### Execution Modes
- **Local**: Run tests on your machine
- **Parallel**: Execute multiple tests concurrently
- **Scheduled**: Cron-based execution
- **On-demand**: Trigger via API or webhook

### Real-Time Dashboard
- Live execution monitoring
- Historical results and trends
- Success rate analytics
- Resource utilization
- Detailed test reports with artifacts

### CLI Tool
```bash
testmesh run --suite smoke --env staging
testmesh watch --filter "api-*"
testmesh debug user-registration --step 3
testmesh generate api-test --openapi swagger.yaml
```

### Extensibility
```typescript
// Create custom action plugin
const myAction: ActionHandler = {
  name: 'my_custom_action',
  async execute(config, context) {
    // Your logic here
    return { success: true, output: {...} };
  }
};
```

## Documentation

All documentation has been organized for easy navigation:

- **[📚 Documentation Index](./docs/README.md)** - Start here for complete documentation
- **[📋 Planning](./docs/planning/)** - Project scope, features, and roadmap
  - [V1_SCOPE.md](./docs/planning/V1_SCOPE.md) - v1.0 scope (all 27 features)
  - [FEATURES.md](./docs/planning/FEATURES.md) - Complete feature specifications
  - [IMPLEMENTATION_PLAN.md](./docs/planning/IMPLEMENTATION_PLAN.md) - 7-phase roadmap (11-15 months)
  - [QUICKSTART.md](./docs/planning/QUICKSTART.md) - Get started in 5 minutes
- **[🏗️ Architecture](./docs/architecture/)** - System design and tech stack
  - [ARCHITECTURE.md](./docs/architecture/ARCHITECTURE.md) - Complete architecture
  - [TECH_STACK.md](./docs/architecture/TECH_STACK.md) - Technology decisions
  - [MODULAR_MONOLITH.md](./docs/architecture/MODULAR_MONOLITH.md) - Architectural approach
- **[✨ Features](./docs/features/)** - Detailed feature designs (20 docs)
  - [FLOW_DESIGN.md](./docs/features/FLOW_DESIGN.md) - Flow-based testing
  - [VISUAL_EDITOR_DESIGN.md](./docs/features/VISUAL_EDITOR_DESIGN.md) - Visual editor UI
  - [YAML_SCHEMA.md](./docs/features/YAML_SCHEMA.md) - Flow definition spec
  - [AI_INTEGRATION.md](./docs/features/AI_INTEGRATION.md) - 🤖 AI-powered testing
  - [CONTRACT_TESTING.md](./docs/features/CONTRACT_TESTING.md) - Consumer-driven contracts
  - [ADVANCED_REPORTING.md](./docs/features/ADVANCED_REPORTING.md) - Reports & analytics
  - ...and 14 more feature docs
- **[📖 Process](./docs/process/)** - Development workflows and standards
  - [DEVELOPMENT_WORKFLOW.md](./docs/process/DEVELOPMENT_WORKFLOW.md) - Git workflow, TDD
  - [SECURITY_GUIDELINES.md](./docs/process/SECURITY_GUIDELINES.md) - Security rules
  - [CODING_STANDARDS.md](./docs/process/CODING_STANDARDS.md) - Code standards
- **[✅ Assessment](./docs/assessment/)** - Implementation readiness
  - [IMPLEMENTATION_READINESS.md](./docs/assessment/IMPLEMENTATION_READINESS.md) - ✅ READY

See [docs/DOCUMENTATION_UPDATE.md](./docs/DOCUMENTATION_UPDATE.md) for recent documentation improvements.

## Architecture

TestMesh is built as a **modular monolith** - a single Go service with clear domain boundaries:

```
┌─────────────┬───────────────┬───────────────┐
│  CLI Tool   │  Web Dashboard│   API Client  │
└──────┬──────┴───────┬───────┴───────┬───────┘
       │              │               │
       └──────────────┴───────────────┘
                      │
       ┌──────────────┴──────────────────┐
       │   TestMesh Server (Go Binary)   │
       ├─────────────────────────────────┤
       │  ┌────────────────────────────┐ │
       │  │ API Domain (HTTP + WS)     │ │
       │  └──────────┬─────────────────┘ │
       │             │                    │
       │  ┌──────────▼─────────────────┐ │
       │  │ Scheduler Domain           │ │
       │  └──────────┬─────────────────┘ │
       │             │                    │
       │  ┌──────────▼─────────────────┐ │
       │  │ Runner Domain              │ │
       │  └──────────┬─────────────────┘ │
       │             │                    │
       │  ┌──────────▼─────────────────┐ │
       │  │ Storage Domain             │ │
       │  └────────────────────────────┘ │
       └──────────────┬──────────────────┘
                      │
       ┌──────────────┴──────────────────┐
       │  PostgreSQL + Redis + Redis Streams  │
       └─────────────────────────────────┘
```

**Benefits**: Faster development, easier debugging, better performance, simpler deployment.

See [docs/architecture/ARCHITECTURE.md](./docs/architecture/ARCHITECTURE.md) or [docs/architecture/MODULAR_MONOLITH.md](./docs/architecture/MODULAR_MONOLITH.md) for details.

## Tech Stack

### Backend (Modular Monolith)
- **Go** - Single service with domain modules
- **Gin** - HTTP framework
- **GORM** - Database ORM
- **PostgreSQL + TimescaleDB** - Data storage with schemas per domain
- **Redis** - Caching and distributed locking
- **Redis Streams** - Async job processing

### Frontend
- **Next.js 14** - React framework with App Router
- **TypeScript** - Type safety
- **React Flow** - Visual flow editor
- **shadcn/ui + Radix UI** - UI components
- **Tailwind CSS** - Styling
- **Socket.io** - Real-time updates

### CLI
- **Go** with Cobra framework
- Cross-platform binary (macOS, Linux, Windows)

### Deployment
- **Docker** - Single image for server + worker
- **Kubernetes + Helm** - Production orchestration
- **Terraform** - Infrastructure as code

### Observability
- **Prometheus + Grafana** - Metrics and dashboards
- **OpenTelemetry + Jaeger** - Distributed tracing
- **Structured logging** - JSON logs

## Deployment

### Local Development
```bash
docker-compose up -d
```

### Kubernetes
```bash
helm install testmesh testmesh/testmesh \
  --set postgres.password=yourpassword
```

### Cloud Providers
- AWS (ECS, EKS, Fargate)
- Google Cloud (GKE, Cloud Run)
- Azure (AKS, Container Instances)

See deployment guides in [infrastructure/](./infrastructure/) directory.

## CI/CD Integration

### GitHub Actions
```yaml
- name: Run TestMesh Tests
  uses: testmesh/github-action@v1
  with:
    suite: smoke
    environment: staging
```

### GitLab CI
```yaml
test:
  image: testmesh/cli:latest
  script:
    - testmesh run --suite critical
```

### Jenkins
```groovy
stage('Integration Tests') {
  steps {
    sh 'testmesh run --suite smoke'
  }
}
```

## Examples

### HTTP API Test
```yaml
test:
  name: "User CRUD Operations"
  steps:
    - name: "Create user"
      action: http_request
      config:
        method: POST
        url: "${API_URL}/users"
        body: { name: "John Doe", email: "john@example.com" }
        assertions:
          - status_code: 201
```

### Database Test
```yaml
test:
  name: "Database Migration"
  steps:
    - name: "Check table exists"
      action: database_query
      config:
        query: "SELECT * FROM information_schema.tables WHERE table_name = 'users'"
        assertions:
          - row_count: 1
```

### Message Queue Test
```yaml
test:
  name: "Event Processing"
  steps:
    - name: "Publish event"
      action: kafka_publish
      config:
        topic: "events"
        message: { type: "user.created", id: "123" }
```

### Browser Test
```yaml
test:
  name: "Login Flow"
  steps:
    - action: browser_navigate
      config: { url: "https://app.example.com/login" }
    - action: browser_fill
      config: { selector: "#email", value: "user@example.com" }
    - action: browser_click
      config: { selector: "button[type=submit]" }
```

More examples in [examples/](./examples/) directory.

## Roadmap

### v1.0 - Comprehensive Launch (Current Focus)

**Core Testing Features:**
- ✅ Flow Definition & Execution (YAML + Visual)
- ✅ Multi-Protocol Support (HTTP, Database, Kafka, gRPC, WebSocket, Browser, MCP)
- ✅ Assertions & Validation (including JSON Schema)
- ✅ Visual Flow Editor (Drag & Drop)
- ✅ Tagging System
- ✅ Plugin System
- ✅ Local Development (CLI tool)
- ✅ Cloud Execution (Kubernetes deployment)
- ✅ Observability & Debugging

**Postman-Inspired Features:**
- ✅ Request Builder UI
- ✅ Response Visualization
- ✅ Collections & Folders
- ✅ Request History
- ✅ Environment Switcher
- ✅ Variable Autocomplete
- ✅ Advanced Auth Helpers (OAuth 2.0, JWT, AWS Signature, etc.)
- ✅ Mock Servers
- ✅ Import/Export (OpenAPI, Swagger, Postman, HAR, cURL, GraphQL)
- ✅ Workspaces
- ✅ Bulk Operations
- ✅ Data-Driven Testing
- ✅ Load Testing

**Advanced Features:**
- ✅ Contract Testing (Pact-compatible)
- ✅ Advanced Reporting & Analytics
- ✅ AI-Powered Testing (Natural Language, Smart Import, Coverage Analysis, Self-Healing)

**Timeline:** 10-13 months with parallel development

See [docs/planning/V1_SCOPE.md](./docs/planning/V1_SCOPE.md) for complete v1.0 scope and [docs/planning/IMPLEMENTATION_PLAN.md](./docs/planning/IMPLEMENTATION_PLAN.md) for detailed roadmap.

### Post-v1.0 Future Enhancements

**v1.1 (Minor Enhancements)**
- Performance optimization
- Additional protocol support (SOAP, MQTT, AMQP)
- Enhanced plugin marketplace

**v2.0 (Major Evolution)**
- Advanced AI capabilities (predictive analytics, automated optimization)
- Multi-region distributed execution
- Enhanced security features
- Enterprise-scale features

## Contributing

We welcome contributions! Please see [CONTRIBUTING.md](./CONTRIBUTING.md) for guidelines.

### Development Setup
```bash
# Clone repository
git clone https://github.com/testmesh/testmesh.git
cd testmesh

# Start development environment
docker-compose up -d

# Install dependencies
npm install

# Run tests
npm test

# Start services in dev mode
npm run dev
```

## Community

- **Discord**: [Join our community](https://discord.gg/testmesh)
- **GitHub Discussions**: [Ask questions and share ideas](https://github.com/testmesh/testmesh/discussions)
- **Twitter**: [@testmesh](https://twitter.com/testmesh)
- **Blog**: [blog.testmesh.io](https://blog.testmesh.io)

## Support

- **Documentation**: [docs.testmesh.io](https://docs.testmesh.io)
- **GitHub Issues**: [Report bugs](https://github.com/testmesh/testmesh/issues)
- **Email**: support@testmesh.io
- **Enterprise Support**: Available for production deployments

## License

TestMesh is open source software licensed under the [MIT License](./LICENSE).

## Acknowledgments

Built with ❤️ by the TestMesh team and contributors.

Special thanks to:
- The open source community
- All our contributors and early adopters
- Projects that inspired TestMesh: Postman, Playwright, K6, and many others

---

**Ready to get started?** Check out the [Quick Start Guide](./docs/planning/QUICKSTART.md) or [join our Discord](https://discord.gg/testmesh) for help!

**Questions?** Open an issue or start a discussion on GitHub.

**Want to contribute?** See [CONTRIBUTING.md](./CONTRIBUTING.md) to get involved!
