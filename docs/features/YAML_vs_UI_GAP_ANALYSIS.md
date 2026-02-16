# TestMesh: YAML vs UI Flow Support - Gap Analysis

> **Comprehensive comparison to achieve feature parity between YAML and UI flow management**

**Date**: 2026-02-16
**Status**: In Progress
**Goal**: Enable 100% of YAML features to be manageable through the visual UI

---

## Table of Contents
1. [Executive Summary](#executive-summary)
2. [Action Types Comparison](#action-types-comparison)
3. [Control Flow Comparison](#control-flow-comparison)
4. [Configuration Features](#configuration-features)
5. [Advanced Features](#advanced-features)
6. [UI Component Gaps](#ui-component-gaps)
7. [Priority Recommendations](#priority-recommendations)
8. [Implementation Roadmap](#implementation-roadmap)

---

## Executive Summary

### Current State

**YAML Support**: ⭐⭐⭐⭐⭐ (100% - Comprehensive)
- 13+ action types
- Full control flow (conditional, loops, parallel)
- Advanced features (retry, error handling, polling, assertions)
- Complete variable system
- Mock server capabilities

**UI Support**: ⭐⭐⭐ (60% - Basic to Intermediate)
- 12 action types (missing 8+ critical ones)
- Basic control flow (condition, for_each)
- Limited configuration options
- Simple validation
- Missing advanced features

### Gap Summary

| Category | YAML Features | UI Features | Gap % |
|----------|--------------|-------------|-------|
| **Action Types** | 21 | 12 | **43%** |
| **Control Flow** | 6 patterns | 2 patterns | **67%** |
| **HTTP Config** | 15 options | 5 options | **67%** |
| **Assertions** | 20+ operators | 5 operators | **75%** |
| **Variables** | 4 types | 1 type | **75%** |
| **Error Handling** | 5 mechanisms | 1 mechanism | **80%** |
| **Output Extraction** | JSONPath + filters | Basic paths | **60%** |
| **Visual Features** | N/A | Limited | **50%** |

**Overall Feature Parity**: **~40%** ✅ implemented, **~60%** 🔴 missing

---

## Action Types Comparison

### ✅ Fully Supported in Both

| Action | YAML | UI | Notes |
|--------|------|----|----|
| HTTP Request | ✅ | ✅ | Basic support in UI |
| Database Query | ✅ | ✅ | Missing polling in UI |
| Log Message | ✅ | ✅ | Full support |
| Delay/Wait | ✅ | ✅ | Full support |
| Assert | ✅ | ✅ | Limited operators in UI |
| Transform | ✅ | ✅ | Basic support in UI |
| Condition (If/Else) | ✅ | ✅ | Full support |
| For Each (Loop) | ✅ | ✅ | Full support |
| Mock Server Start | ✅ | ✅ | Full support |
| Mock Server Stop | ✅ | ✅ | Full support |
| Contract Generate | ✅ | ✅ | Full support |
| Contract Verify | ✅ | ✅ | Full support |

### 🔴 Missing in UI

| Action | YAML Support | UI Support | Priority |
|--------|-------------|-----------|----------|
| **Kafka Publish** | ✅ Full | ❌ None | 🔴 HIGH |
| **Kafka Consume** | ✅ Full | ❌ None | 🔴 HIGH |
| **gRPC Call** | ✅ Full | ❌ None | 🟡 MEDIUM |
| **gRPC Stream** | ✅ Full | ❌ None | 🟡 MEDIUM |
| **WebSocket** | ✅ Full | ❌ None | 🟡 MEDIUM |
| **Browser Navigate** | ✅ Full | ❌ None | 🟡 MEDIUM |
| **Browser Actions** | ✅ 20+ actions | ❌ None | 🟡 MEDIUM |
| **Wait Until (Poll)** | ✅ Full | ❌ None | 🔴 HIGH |
| **Sub-flow** | ✅ Full | ❌ None | 🔴 HIGH |
| **Parallel** | ✅ Full | ❌ None | 🔴 HIGH |
| **Mock Server Verify** | ✅ Full | ❌ None | 🟢 LOW |
| **Mock Server Update** | ✅ Full | ❌ None | 🟢 LOW |
| **Mock Server Reset** | ✅ Full | ❌ None | 🟢 LOW |

---

## Control Flow Comparison

### YAML Capabilities

| Pattern | YAML | Example |
|---------|------|---------|
| **Sequential** | ✅ | Steps run in order |
| **Conditional (If/Else)** | ✅ | `when:` clause or `condition` action |
| **Loops (For Each)** | ✅ | `for_each` with items/range/glob |
| **Parallel** | ✅ | `parallel` action with concurrent steps |
| **Try/Catch** | ✅ | `on_error: continue` + conditional handling |
| **Switch/Case** | ✅ | Multiple conditional steps |

### UI Capabilities

| Pattern | UI | Status | Notes |
|---------|----|----|-------|
| **Sequential** | ✅ | Full | Connect nodes linearly |
| **Conditional (If/Else)** | ✅ | Full | Condition node with branches |
| **Loops (For Each)** | ✅ | Full | ForEach node |
| **Parallel** | ❌ | **Missing** | No parallel node yet |
| **Try/Catch** | ❌ | **Missing** | No error handling UI |
| **Switch/Case** | ⚠️ | Partial | Multiple conditions manually |

### 🔴 Missing Control Flow Features in UI

1. **Parallel Execution Node**
   - Execute multiple steps concurrently
   - Wait for all or continue on first success
   - Max concurrency control

2. **Error Handling**
   - `error_steps` visual representation
   - `on_error` dropdown (fail/continue/retry)
   - `on_timeout` handlers

3. **Advanced Conditionals**
   - Multi-branch conditions (switch/case)
   - Complex boolean expressions builder

---

## Configuration Features

### HTTP Request Configuration

| Feature | YAML | UI | Priority |
|---------|------|----|----|
| **Method** | ✅ GET/POST/PUT/PATCH/DELETE/HEAD/OPTIONS | ✅ All methods | ✅ Done |
| **URL** | ✅ With variables | ✅ With variables | ✅ Done |
| **Headers** | ✅ Key-value map | ✅ Key-value editor | ✅ Done |
| **Query Params** | ✅ `params` object | ❌ None | 🔴 HIGH |
| **Body** | ✅ JSON/String/Raw | ⚠️ Basic | 🟡 MEDIUM |
| **Auth** | ✅ Basic/Bearer/API Key/OAuth2 | ❌ None | 🔴 HIGH |
| **Follow Redirects** | ✅ Yes/No + max | ❌ None | 🟢 LOW |
| **SSL/TLS** | ✅ verify_ssl, client cert/key | ❌ None | 🟡 MEDIUM |
| **Cookies** | ✅ Key-value map | ❌ None | 🟢 LOW |
| **Timeout** | ✅ Duration | ⚠️ Basic | 🟡 MEDIUM |

### Database Query Configuration

| Feature | YAML | UI | Priority |
|---------|------|----|----|
| **Connection** | ✅ Type, host, port, database, user, password | ⚠️ Basic | 🔴 HIGH |
| **Query** | ✅ SQL/MongoDB query | ✅ Text area | ✅ Done |
| **Params** | ✅ Positional or named | ❌ None | 🔴 HIGH |
| **Transaction** | ✅ Boolean | ❌ None | 🟡 MEDIUM |
| **Polling** | ✅ enabled, timeout, interval | ❌ None | 🔴 HIGH |

### Kafka Configuration

| Feature | YAML | UI | Priority |
|---------|------|----|----|
| **Brokers** | ✅ Array of broker URLs | ❌ None | 🔴 HIGH |
| **Topic** | ✅ String | ❌ None | 🔴 HIGH |
| **Message Key/Value** | ✅ String/Object | ❌ None | 🔴 HIGH |
| **Headers** | ✅ Key-value map | ❌ None | 🟡 MEDIUM |
| **Partition** | ✅ Number | ❌ None | 🟢 LOW |
| **Compression** | ✅ none/gzip/snappy/lz4 | ❌ None | 🟢 LOW |
| **SASL Auth** | ✅ mechanism, username, password | ❌ None | 🟡 MEDIUM |
| **Consume Options** | ✅ timeout, max_messages, from_beginning | ❌ None | 🔴 HIGH |
| **Match/Filter** | ✅ key, json_path, header filters | ❌ None | 🔴 HIGH |

### gRPC Configuration

| Feature | YAML | UI | Priority |
|---------|------|----|----|
| **Service/Method** | ✅ Service, method names | ❌ None | 🟡 MEDIUM |
| **Request** | ✅ Message fields | ❌ None | 🟡 MEDIUM |
| **Metadata** | ✅ Headers | ❌ None | 🟡 MEDIUM |
| **TLS** | ✅ cert, key, ca | ❌ None | 🟢 LOW |
| **Streaming** | ✅ client/server/bidirectional | ❌ None | 🟡 MEDIUM |

### Browser Automation Configuration

| Feature | YAML | UI | Priority |
|---------|------|----|----|
| **Browser Type** | ✅ chromium/firefox/webkit | ❌ None | 🟡 MEDIUM |
| **Headless** | ✅ Boolean | ❌ None | 🟡 MEDIUM |
| **Viewport** | ✅ width, height | ❌ None | 🟢 LOW |
| **Device Emulation** | ✅ Predefined devices | ❌ None | 🟢 LOW |
| **Actions** | ✅ 20+ action types | ❌ None | 🟡 MEDIUM |
| **Network Intercept** | ✅ enabled, patterns | ❌ None | 🟢 LOW |

---

## Advanced Features

### Assertions

| Feature | YAML | UI | Priority |
|---------|------|----|----|
| **Equality** | ✅ `==`, `!=` | ✅ Basic | ✅ Done |
| **Comparison** | ✅ `>`, `<`, `>=`, `<=` | ⚠️ Limited | 🟡 MEDIUM |
| **Existence** | ✅ `exists`, `is null`, `is not null` | ❌ None | 🔴 HIGH |
| **Type Checks** | ✅ `is number`, `is string`, `is boolean`, etc. | ❌ None | 🟡 MEDIUM |
| **String Ops** | ✅ `contains`, `starts_with`, `ends_with`, `matches` | ⚠️ Limited | 🔴 HIGH |
| **Array Ops** | ✅ `length`, `contains` | ❌ None | 🟡 MEDIUM |
| **Boolean Logic** | ✅ `&&`, `||`, `!` | ❌ None | 🔴 HIGH |
| **Custom Messages** | ✅ `expression` + `message` | ❌ None | 🟡 MEDIUM |
| **Assert Modes** | ✅ `assert`, `assert_any`, `assert_none` | ❌ None | 🟡 MEDIUM |
| **Visual Builder** | N/A | ❌ None | 🔴 HIGH |

### Output Extraction

| Feature | YAML | UI | Priority |
|---------|------|----|----|
| **Basic Paths** | ✅ `response.body.id` | ✅ Text input | ✅ Done |
| **JSONPath** | ✅ Full syntax | ⚠️ No validation | 🔴 HIGH |
| **Array Access** | ✅ `[0]`, `[-1]`, `[0:3]` | ⚠️ No UI helper | 🟡 MEDIUM |
| **Filters** | ✅ `[?(@.status == 'active')]` | ❌ None | 🟡 MEDIUM |
| **Functions** | ✅ `length`, `sum()`, `avg()`, `min()`, `max()` | ❌ None | 🟡 MEDIUM |
| **Preview** | N/A | ⚠️ Basic | 🟡 MEDIUM |
| **Path Builder** | N/A | ❌ None | 🔴 HIGH |

### Variables

| Feature | YAML | UI | Priority |
|---------|------|----|----|
| **Environment Vars** | ✅ `env:` section | ⚠️ External only | 🔴 HIGH |
| **System Vars** | ✅ FLOW_ID, EXECUTION_ID, TIMESTAMP, etc. | ❌ None | 🟡 MEDIUM |
| **Faker Vars** | ✅ 100+ faker functions | ❌ None | 🟡 MEDIUM |
| **Step Outputs** | ✅ `${step_id.output}` | ✅ Basic | ✅ Done |
| **Default Values** | ✅ `${VAR:default}` | ❌ None | 🟢 LOW |
| **Ternary** | ✅ `${condition ? true : false}` | ❌ None | 🟢 LOW |
| **Autocomplete** | N/A | ⚠️ Basic | 🔴 HIGH |

### Error Handling

| Feature | YAML | UI | Priority |
|---------|------|----|----|
| **on_error** | ✅ continue/fail/retry | ❌ None | 🔴 HIGH |
| **error_steps** | ✅ Steps to run on error | ❌ None | 🔴 HIGH |
| **Retry Config** | ✅ max_attempts, delay, backoff | ⚠️ Basic | 🔴 HIGH |
| **retry_on** | ✅ Conditions for retry | ❌ None | 🟡 MEDIUM |
| **on_timeout** | ✅ Steps to run on timeout | ❌ None | 🟡 MEDIUM |
| **Visual Error Paths** | N/A | ❌ None | 🔴 HIGH |

### Flow-Level Features

| Feature | YAML | UI | Priority |
|---------|------|----|----|
| **Setup/Teardown** | ✅ Separate step sections | ✅ Visual sections | ✅ Done |
| **Flow Config** | ✅ timeout, fail_fast, retry | ❌ None | 🟡 MEDIUM |
| **Flow-level on_error** | ✅ Error handler steps | ❌ None | 🟡 MEDIUM |
| **Tags** | ✅ Array of strings | ✅ Basic | ✅ Done |
| **Suite** | ✅ String | ✅ Basic | ✅ Done |
| **Author** | ✅ String | ❌ None | 🟢 LOW |

---

## UI Component Gaps

### Properties Panel - Missing Components

1. **🔴 Authentication Builder**
   - Dropdown: Basic/Bearer/API Key/OAuth2
   - Dynamic fields based on type
   - Secure credential storage reference

2. **🔴 Query Parameters Editor**
   - Key-value pairs
   - Variable interpolation
   - URL preview with params

3. **🔴 Assertion Builder**
   - Visual builder with dropdowns
   - Field selector (with autocomplete)
   - Operator selector
   - Value input
   - Preview/test functionality

4. **🔴 JSONPath Builder**
   - Path input with autocomplete
   - Preview with sample data
   - Common patterns library
   - Filter builder

5. **🔴 Retry Configuration Panel**
   - Enable/disable toggle
   - Max attempts slider
   - Backoff strategy dropdown
   - Retry conditions builder

6. **🔴 Error Handling Panel**
   - on_error dropdown
   - error_steps mini-flow
   - on_timeout handler

7. **🟡 Database Connection Editor**
   - Connection type dropdown
   - Connection string builder
   - Test connection button

8. **🟡 Variable Picker**
   - Dropdown/modal with all available variables
   - Categories: System, Faker, Environment, Step outputs
   - Insert at cursor position

### Node Palette - Missing Nodes

1. **🔴 Kafka Nodes**
   - Kafka Publish node
   - Kafka Consume node

2. **🟡 gRPC Nodes**
   - gRPC Call node
   - gRPC Stream node

3. **🟡 WebSocket Node**
   - WebSocket connection/send/receive node

4. **🟡 Browser Nodes**
   - Browser Navigate node
   - Browser Actions node group

5. **🔴 Wait/Poll Node**
   - Wait Until node (polling)

6. **🔴 Sub-flow Node**
   - Run Flow node

7. **🔴 Parallel Node**
   - Parallel execution container

### Canvas Features - Missing

1. **🔴 Execution Visualization**
   - Live node status updates
   - Animated data flow on edges
   - Step timing display
   - Error highlighting

2. **🔴 Conditional Branch Visualization**
   - True/false path labels
   - Visual branch indicators

3. **🔴 Loop Visualization**
   - Loop-back arrow
   - Iteration counter

4. **🔴 Parallel Container**
   - Visual parallel group
   - Multiple output handles

5. **🟡 Mini-map**
   - Overview of entire flow
   - Viewport indicator
   - Click to navigate

6. **🟡 Grid & Snap**
   - Toggle grid visibility
   - Snap to grid
   - Alignment guides

### Toolbar - Missing Actions

1. **🟡 Auto-Layout**
   - Automatic node arrangement
   - Hierarchical layout
   - Align/distribute tools

2. **🟡 Export Options**
   - Export as PNG/SVG
   - Export as YAML
   - Export as JSON

3. **🔴 Validation**
   - Validate flow button
   - Show validation errors
   - Quick fix suggestions

---

## Priority Recommendations

### 🔴 Phase 1: Critical Gaps (Must Have)

**Goal**: Enable essential workflows that are currently impossible in the UI

1. **Kafka Support** (2-3 weeks)
   - Kafka Publish node
   - Kafka Consume node with filters
   - Broker/topic configuration UI

2. **Wait/Poll Node** (1 week)
   - Wait Until action
   - Polling configuration
   - Timeout handling

3. **Sub-flow Support** (2 weeks)
   - Sub-flow node
   - Flow selector
   - Input/output mapping

4. **Parallel Execution** (2 weeks)
   - Parallel node container
   - Multiple output handles
   - Wait for all/any configuration

5. **Enhanced HTTP Configuration** (2 weeks)
   - Query parameters editor
   - Authentication builder
   - Advanced options (SSL, cookies, redirects)

6. **Enhanced Database Configuration** (1 week)
   - Parameterized queries UI
   - Polling configuration
   - Connection editor

7. **Assertion Builder** (2 weeks)
   - Visual assertion builder
   - All operators
   - Custom messages

8. **Output/JSONPath Builder** (1-2 weeks)
   - JSONPath autocomplete
   - Preview with sample data
   - Common patterns

9. **Error Handling UI** (2 weeks)
   - on_error dropdown
   - error_steps visual editor
   - Retry configuration panel

10. **Execution Visualization** (2-3 weeks)
    - Live status updates
    - Animated flow
    - Step details popover

**Total: ~18-24 weeks (4.5-6 months)**

### 🟡 Phase 2: Important Enhancements (Should Have)

**Goal**: Improve usability and enable advanced scenarios

1. **gRPC Support** (2 weeks)
   - gRPC Call node
   - gRPC Stream node

2. **WebSocket Support** (1-2 weeks)
   - WebSocket node with actions

3. **Browser Automation** (3-4 weeks)
   - Browser nodes (20+ actions)
   - Action sequence builder

4. **Variable System** (2 weeks)
   - System variables
   - Faker variables
   - Variable picker UI

5. **Advanced Assertions** (1 week)
   - Type checks
   - Array operations
   - Boolean logic builder

6. **Flow-Level Configuration** (1 week)
   - Flow config editor
   - Flow-level error handlers

7. **Transform Operations** (2 weeks)
   - Map/filter/join UI
   - Transform operation builder

8. **Mini-map & Navigation** (1 week)
   - Canvas mini-map
   - Zoom controls
   - Overview panel

9. **Auto-Layout** (1-2 weeks)
   - Automatic arrangement
   - Alignment tools

**Total: ~15-19 weeks (3.5-4.5 months)**

### 🟢 Phase 3: Nice to Have (Could Have)

**Goal**: Polish and convenience features

1. **Mock Server Advanced** (1 week)
   - Verify/update/reset UI
   - State management

2. **Advanced SSL/TLS** (1 week)
   - Client certificates UI

3. **Cookies Editor** (1 week)
   - Cookie management UI

4. **Device Emulation** (Browser) (1 week)
   - Device picker

5. **Network Interception** (Browser) (1 week)
   - Intercept patterns UI

6. **Flow Templates** (2 weeks)
   - Template library
   - Template creation

7. **Export Options** (1 week)
   - PNG/SVG export
   - Pretty-print options

8. **Collaboration Features** (4-6 weeks)
   - Real-time editing
   - Comments
   - Activity feed

**Total: ~12-15 weeks (3-3.75 months)**

---

## Implementation Roadmap

### Quarter 1 (Q1) - Foundation

**Weeks 1-4**: Core Missing Actions
- ✅ Kafka Publish/Consume nodes
- ✅ Wait/Poll node
- ✅ Sub-flow node

**Weeks 5-8**: Essential Configurations
- ✅ HTTP: Query params, Auth builder
- ✅ Database: Params, Polling
- ✅ Assertion builder (basic)

**Weeks 9-12**: Control Flow & Error Handling
- ✅ Parallel execution node
- ✅ Error handling UI
- ✅ Retry configuration

### Quarter 2 (Q2) - Enhancement

**Weeks 13-16**: Advanced Actions
- ✅ gRPC support
- ✅ WebSocket support
- ✅ Output/JSONPath builder

**Weeks 17-20**: Visual Improvements
- ✅ Execution visualization
- ✅ Mini-map & navigation
- ✅ Auto-layout

**Weeks 21-24**: Variables & Transforms
- ✅ Variable system UI
- ✅ Transform operations
- ✅ Advanced assertions

### Quarter 3 (Q3) - Polish & Completion

**Weeks 25-28**: Browser Automation
- ✅ Browser nodes
- ✅ Action builder

**Weeks 29-32**: Flow-Level Features
- ✅ Flow configuration
- ✅ Flow-level error handlers
- ✅ Export options

**Weeks 33-36**: Final Polish
- ✅ Templates
- ✅ Testing & QA
- ✅ Documentation

---

## Success Metrics

### Feature Parity

| Milestone | Target | Timeline |
|-----------|--------|----------|
| **Phase 1 Complete** | 70% parity | Q1 End (Week 12) |
| **Phase 2 Complete** | 90% parity | Q2 End (Week 24) |
| **Phase 3 Complete** | 95%+ parity | Q3 End (Week 36) |

### User Adoption

- **Goal**: 80% of flows created via UI (vs YAML) by end of Q2
- **Metric**: Track creation method per flow
- **Validation**: User surveys and feedback

### Performance

- **Goal**: Visual editor handles 100+ node flows smoothly
- **Metric**: Canvas rendering performance, interaction latency
- **Target**: <100ms response time for all interactions

---

## Appendix: Quick Reference

### YAML → UI Mapping Table

| YAML Feature | UI Component | Status | Priority |
|--------------|--------------|--------|----------|
| `action: http_request` | HTTP Request Node | ✅ Partial | 🔴 Complete |
| `action: database_query` | Database Query Node | ✅ Partial | 🔴 Complete |
| `action: kafka_publish` | Kafka Publish Node | ❌ Missing | 🔴 High |
| `action: kafka_consume` | Kafka Consume Node | ❌ Missing | 🔴 High |
| `action: grpc_call` | gRPC Call Node | ❌ Missing | 🟡 Medium |
| `action: websocket` | WebSocket Node | ❌ Missing | 🟡 Medium |
| `action: browser` | Browser Node | ❌ Missing | 🟡 Medium |
| `action: wait_until` | Wait/Poll Node | ❌ Missing | 🔴 High |
| `action: run_flow` | Sub-flow Node | ❌ Missing | 🔴 High |
| `action: parallel` | Parallel Node | ❌ Missing | 🔴 High |
| `action: condition` | Condition Node | ✅ Done | ✅ Done |
| `action: for_each` | ForEach Node | ✅ Done | ✅ Done |
| `auth:` | Auth Builder | ❌ Missing | 🔴 High |
| `params:` | Query Params Editor | ❌ Missing | 🔴 High |
| `assert:` | Assertion Builder | ⚠️ Basic | 🔴 Complete |
| `output:` | Output Mapping | ⚠️ Basic | 🔴 Complete |
| `retry:` | Retry Config Panel | ⚠️ Basic | 🔴 Complete |
| `on_error:` | Error Handling UI | ❌ Missing | 🔴 High |
| `error_steps:` | Error Steps Editor | ❌ Missing | 🔴 High |
| `config.poll:` | Polling Config | ❌ Missing | 🔴 High |

---

## Conclusion

To achieve feature parity and make the UI a true replacement for YAML editing, focus on:

1. **Phase 1 (Critical)**: Kafka, Wait/Poll, Sub-flows, Parallel, Enhanced HTTP/DB config, Assertions, Error handling
2. **Phase 2 (Important)**: gRPC, WebSocket, Browser, Variables, Advanced features
3. **Phase 3 (Polish)**: Templates, Collaboration, Export options

**Estimated Total Effort**: 45-58 weeks (~11-14 months) for 95%+ parity

**Recommended Approach**:
- Start with Phase 1 (Q1) to enable most critical workflows
- Gather user feedback during Phase 1
- Adjust Phase 2 priorities based on actual usage patterns
- Phase 3 can be incremental based on user requests

By end of Q2, you'll have ~90% feature parity, enabling most users to work exclusively in the visual UI while maintaining YAML as a power-user option.
