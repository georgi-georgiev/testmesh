# TestMesh Documentation

Welcome to the TestMesh documentation. This directory contains all planning, architecture, feature design, and process documentation.

## Documentation Structure

### 📋 Planning (`/planning/`)
Core planning documents that define what we're building and when.

- **[V1_SCOPE.md](./planning/V1_SCOPE.md)** - Complete v1.0 scope definition
- **[FEATURES.md](./planning/FEATURES.md)** - All 27 features with detailed specifications
- **[IMPLEMENTATION_PLAN.md](./planning/IMPLEMENTATION_PLAN.md)** - 7-phase roadmap (10-13 months)
- **[SUMMARY.md](./planning/SUMMARY.md)** - Executive summary
- **[QUICKSTART.md](./planning/QUICKSTART.md)** - Getting started guide
- **[V1_CONSOLIDATION_SUMMARY.md](./planning/V1_CONSOLIDATION_SUMMARY.md)** - Version consolidation notes
- **[V1_NEW_FEATURES_SUMMARY.md](./planning/V1_NEW_FEATURES_SUMMARY.md)** - New features added to v1.0
- **[IMPLEMENTATION_PLAN_UPDATE_SUMMARY.md](./planning/IMPLEMENTATION_PLAN_UPDATE_SUMMARY.md)** - Plan change log
- **[AI_ROADMAP_SUMMARY.md](./planning/AI_ROADMAP_SUMMARY.md)** - AI features roadmap

### 🏗️ Architecture (`/architecture/`)
System architecture, technical stack, and structural decisions.

- **[ARCHITECTURE.md](./architecture/ARCHITECTURE.md)** - Complete system architecture
- **[MODULAR_MONOLITH.md](./architecture/MODULAR_MONOLITH.md)** - Architectural approach and rationale
- **[TECH_STACK.md](./architecture/TECH_STACK.md)** - Technology choices with examples
- **[PROJECT_STRUCTURE.md](./architecture/PROJECT_STRUCTURE.md)** - Code organization and folder layout
- **[ARCHITECTURE_SUMMARY.md](./architecture/ARCHITECTURE_SUMMARY.md)** - Quick architecture reference
- **[TECHNOLOGY_SUMMARY.md](./architecture/TECHNOLOGY_SUMMARY.md)** - Tech stack review and approval
- **[DECISIONS.md](./architecture/DECISIONS.md)** - Key design decisions and trade-offs

### ✨ Features (`/features/`)
Detailed design documents for each major feature area.

**Core Features:**
- **[FLOW_DESIGN.md](./features/FLOW_DESIGN.md)** - Flow-based testing model
- **[VISUAL_EDITOR_DESIGN.md](./features/VISUAL_EDITOR_DESIGN.md)** - Visual editor UI/UX
- **[YAML_SCHEMA.md](./features/YAML_SCHEMA.md)** - Flow definition specification

**Testing Features:**
- **[TAGGING_SYSTEM.md](./features/TAGGING_SYSTEM.md)** - Tagging and filtering
- **[TEST_DATA_MANAGEMENT.md](./features/TEST_DATA_MANAGEMENT.md)** - Test data tracking
- **[DATA_GENERATION.md](./features/DATA_GENERATION.md)** - Faker, templates, data generation
- **[ASYNC_PATTERNS.md](./features/ASYNC_PATTERNS.md)** - Eventual consistency testing
- **[JSON_SCHEMA_VALIDATION.md](./features/JSON_SCHEMA_VALIDATION.md)** - Response validation
- **[CONTRACT_TESTING.md](./features/CONTRACT_TESTING.md)** - Consumer-driven contracts

**Integration & Extensibility:**
- **[MCP_INTEGRATION.md](./features/MCP_INTEGRATION.md)** - AI agent integration (Claude)
- **[AI_INTEGRATION.md](./features/AI_INTEGRATION.md)** - AI-powered testing features
- **[MOCK_SERVER.md](./features/MOCK_SERVER.md)** - Mock server system

**User Experience:**
- **[POSTMAN_INSPIRED_FEATURES.md](./features/POSTMAN_INSPIRED_FEATURES.md)** - Postman-like UX features
- **[CLI_DASHBOARD_PARITY.md](./features/CLI_DASHBOARD_PARITY.md)** - Feature parity between CLI and UI
- **[DASHBOARD_UI_SPECIFICATION.md](./features/DASHBOARD_UI_SPECIFICATION.md)** - Dashboard UI design
- **[DASHBOARD_COVERAGE_ANALYSIS.md](./features/DASHBOARD_COVERAGE_ANALYSIS.md)** - Coverage analysis UI
- **[DASHBOARD_GAP_CLOSURE_SUMMARY.md](./features/DASHBOARD_GAP_CLOSURE_SUMMARY.md)** - Dashboard gap analysis

**Operations:**
- **[OBSERVABILITY.md](./features/OBSERVABILITY.md)** - Logging, metrics, tracing
- **[ADVANCED_REPORTING.md](./features/ADVANCED_REPORTING.md)** - HTML reports, trends, analytics
- **[LOCAL_DEVELOPMENT.md](./features/LOCAL_DEVELOPMENT.md)** - Local development setup
- **[CLOUD_EXECUTION.md](./features/CLOUD_EXECUTION.md)** - Distributed agent architecture

### 📖 Process (`/process/`)
Development workflows, coding standards, and guidelines.

- **[DEVELOPMENT_WORKFLOW.md](./process/DEVELOPMENT_WORKFLOW.md)** - Git workflow, TDD, permissions
- **[SECURITY_GUIDELINES.md](./process/SECURITY_GUIDELINES.md)** - Security rules and OWASP compliance
- **[AGENT_CONTRACT.md](./process/AGENT_CONTRACT.md)** - Meta-instructions for AI agents
- **[CODING_STANDARDS.md](./process/CODING_STANDARDS.md)** - Go/TypeScript standards with examples
- **[CODE_REVIEW_CHECKLIST.md](./process/CODE_REVIEW_CHECKLIST.md)** - Comprehensive review checklist
- **[PLUGIN_DEVELOPMENT.md](./process/PLUGIN_DEVELOPMENT.md)** - Plugin system guide
- **[RECOMMENDED_TOOLING.md](./process/RECOMMENDED_TOOLING.md)** - Development tools and IDE setup

### ✅ Assessment (`/assessment/`)
Readiness verification and pre-implementation checks.

- **[IMPLEMENTATION_READINESS.md](./assessment/IMPLEMENTATION_READINESS.md)** - Final go/no-go decision (✅ READY)
- **[PRE_IMPLEMENTATION_ASSESSMENT.md](./assessment/PRE_IMPLEMENTATION_ASSESSMENT.md)** - 12-point readiness check (12/12 ✅)
- **[PRE_IMPLEMENTATION_CHECKLIST.md](./assessment/PRE_IMPLEMENTATION_CHECKLIST.md)** - Implementation checklist

---

## Quick Navigation

### I want to...

**Understand the project:**
- Start here: [../README.md](../README.md)
- Executive summary: [planning/SUMMARY.md](./planning/SUMMARY.md)
- What's in v1.0: [planning/V1_SCOPE.md](./planning/V1_SCOPE.md)

**Understand the architecture:**
- System design: [architecture/ARCHITECTURE.md](./architecture/ARCHITECTURE.md)
- Tech choices: [architecture/TECH_STACK.md](./architecture/TECH_STACK.md)
- Why modular monolith: [architecture/MODULAR_MONOLITH.md](./architecture/MODULAR_MONOLITH.md)

**Start implementing:**
- Implementation roadmap: [planning/IMPLEMENTATION_PLAN.md](./planning/IMPLEMENTATION_PLAN.md)
- Development workflow: [process/DEVELOPMENT_WORKFLOW.md](./process/DEVELOPMENT_WORKFLOW.md)
- Coding standards: [process/CODING_STANDARDS.md](./process/CODING_STANDARDS.md)

**Design a feature:**
- All features: [planning/FEATURES.md](./planning/FEATURES.md)
- Feature designs: [features/](./features/)

**Review code:**
- Code review checklist: [process/CODE_REVIEW_CHECKLIST.md](./process/CODE_REVIEW_CHECKLIST.md)
- Security guidelines: [process/SECURITY_GUIDELINES.md](./process/SECURITY_GUIDELINES.md)

---

## Document Status

**Last Updated:** 2026-02-11

**Status:** ✅ **READY FOR IMPLEMENTATION**

All documentation is complete and approved. Implementation can begin.

See [assessment/IMPLEMENTATION_READINESS.md](./assessment/IMPLEMENTATION_READINESS.md) for final verification.
