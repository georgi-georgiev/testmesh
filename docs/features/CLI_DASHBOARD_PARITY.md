# CLI ↔ Dashboard Feature Parity

## Overview

This document maps features between the CLI and Dashboard to ensure users can accomplish tasks in their preferred interface.

**Philosophy**:
- ✅ **Core features** should be available in both CLI and Dashboard
- ✅ **Visual features** (drag-and-drop, charts) are Dashboard-only
- ✅ **Automation features** (watch mode, hooks) are CLI-only
- ✅ **Both interfaces** should feel complete and powerful

---

## Feature Comparison Matrix

| Feature | CLI | Dashboard | Notes |
|---------|-----|-----------|-------|
| **Flow Management** |
| Create flow | ✅ `testmesh create flow` | ✅ Visual Editor / Request Builder | CLI: Template-based, Dashboard: Visual |
| Edit flow | ✅ Text editor | ✅ Visual Editor / YAML editor | Both support YAML editing |
| Delete flow | ✅ `testmesh delete flow` | ✅ UI delete button | |
| List flows | ✅ `testmesh list` | ✅ Collections sidebar | |
| Search flows | ✅ `testmesh list --search` | ✅ Search bar | |
| Validate flow | ✅ `testmesh validate` | ✅ Auto-validation in editor | CLI: Pre-run, Dashboard: Real-time |
| Import flow | ✅ `testmesh import` | ✅ Import UI | |
| Export flow | ✅ `testmesh export` | ✅ Export button | |
| **Execution** |
| Run flow | ✅ `testmesh run` | ✅ Run button | |
| Run with env | ✅ `testmesh run --env` | ✅ Environment dropdown | |
| Run with variables | ✅ `testmesh run --var` | ✅ Variable editor | |
| Run collection | ✅ `testmesh run collection/` | ✅ Run collection button | |
| Run by tag | ✅ `testmesh run --tag` | ✅ Tag filter + run | |
| Schedule execution | ✅ `testmesh schedule` | ✅ Schedule UI | |
| Cancel execution | ✅ `testmesh cancel` | ✅ Stop button | |
| **Watch & Monitor** |
| Watch mode | ✅ `testmesh watch` | ❌ N/A | CLI-only: Auto-rerun on file change |
| Live execution | ⚠️ Limited | ✅ Real-time UI | Dashboard: Full visualization |
| Real-time logs | ⚠️ Stream to terminal | ✅ Live log viewer | |
| **Results & Observability** |
| View results | ✅ `testmesh results` | ✅ Execution detail page | |
| View logs | ✅ `testmesh logs` | ✅ Logs tab | CLI: Text output, Dashboard: Searchable |
| View variables | ✅ `testmesh results --variables` | ✅ Variables tab | |
| View network | ✅ `testmesh results --network` | ✅ Network tab | Dashboard: Waterfall view |
| View metrics | ✅ `testmesh results --metrics` | ✅ Metrics tab | Dashboard: Charts |
| Export results | ✅ `testmesh results export` | ✅ Export button | |
| Compare executions | ❌ Missing | ✅ Compare UI | **GAP: Need CLI command** |
| **Cleanup Management** |
| View tracked resources | ✅ `testmesh cleanup list` | ✅ Cleanup tab | |
| Trigger cleanup | ✅ `testmesh cleanup run` | ✅ Cleanup button | |
| View orphans | ✅ `testmesh cleanup orphans` | ✅ Orphaned resources page | |
| Retry cleanup | ✅ `testmesh cleanup retry` | ✅ Retry button | |
| Verify cleanup | ✅ `testmesh cleanup verify` | ✅ Auto-verification | |
| Cleanup report | ✅ `testmesh cleanup report` | ✅ Summary view | |
| **Environment Management** |
| List environments | ✅ `testmesh env list` | ✅ Environment dropdown | |
| Create environment | ✅ `testmesh env create` | ✅ New environment UI | |
| Edit environment | ✅ `testmesh env set` | ✅ Environment editor | |
| Delete environment | ✅ `testmesh env delete` | ✅ Delete button | |
| Switch environment | ✅ `testmesh run --env` | ✅ Environment dropdown | |
| Manage secrets | ✅ `testmesh secret set` | ✅ Secret manager UI | |
| **Collections & Organization** |
| Create collection | ✅ `testmesh collection create` | ✅ New collection button | |
| Create folder | ✅ `testmesh folder create` | ✅ New folder button | |
| Move flows | ✅ `testmesh move` | ✅ Drag & drop | Dashboard: Easier |
| Organize hierarchy | ⚠️ Commands | ✅ Drag & drop | Dashboard: Visual |
| **Request Builder** |
| Build HTTP request | ⚠️ YAML editing | ✅ Visual request builder | Dashboard: Visual forms |
| Test request | ✅ `testmesh run` | ✅ Send button | |
| View response | ✅ Terminal output | ✅ Pretty-print viewer | Dashboard: Tree view, search |
| Copy as cURL | ✅ `testmesh curl` | ✅ Copy cURL button | |
| **History** |
| View history | ✅ `testmesh history` | ✅ History sidebar | |
| Re-run from history | ✅ `testmesh history run` | ✅ Re-run button | |
| Save from history | ✅ `testmesh history save` | ✅ Save button | |
| Clear history | ✅ `testmesh history clear` | ✅ Clear button | |
| **Mock Servers** |
| Create mock | ✅ `testmesh mock create` | ✅ Mock server UI | |
| Start mock | ✅ `testmesh mock start` | ✅ Start button | |
| Stop mock | ✅ `testmesh mock stop` | ✅ Stop button | |
| View mock logs | ✅ `testmesh mock logs` | ✅ Mock logs tab | |
| Mock analytics | ❌ Missing | ✅ Analytics tab | **GAP: Need CLI command** |
| **Import/Export** |
| Import OpenAPI | ✅ `testmesh import openapi` | ✅ Import UI | |
| Import Postman | ✅ `testmesh import postman` | ✅ Import UI | |
| Import cURL | ✅ `testmesh import curl` | ✅ Paste cURL UI | |
| Export collection | ✅ `testmesh export` | ✅ Export button | |
| **Workspaces** |
| List workspaces | ✅ `testmesh workspace list` | ✅ Workspace switcher | |
| Create workspace | ✅ `testmesh workspace create` | ✅ New workspace UI | |
| Switch workspace | ✅ `testmesh workspace use` | ✅ Workspace dropdown | |
| Share workspace | ❌ Missing | ✅ Share UI | **GAP: Need CLI command** |
| **Bulk Operations** |
| Select multiple | ⚠️ Glob patterns | ✅ Checkboxes | CLI: Pattern matching |
| Bulk tag | ✅ `testmesh tag add --flows` | ✅ Bulk edit UI | |
| Bulk move | ✅ `testmesh move --flows` | ✅ Drag & drop | |
| Bulk delete | ✅ `testmesh delete --flows` | ✅ Delete button | |
| Find & replace | ⚠️ sed/awk | ✅ Find & replace UI | Dashboard: Easier |
| **Data-Driven Testing** |
| Run with CSV | ✅ `testmesh run --data` | ✅ Data file upload UI | |
| Preview data | ❌ Missing | ✅ Data preview | **GAP: Need CLI command** |
| View iteration results | ✅ `testmesh results --iterations` | ✅ Iteration results UI | |
| **Load Testing** |
| Run load test | ✅ `testmesh load` | ✅ Load test UI | |
| Configure VUs | ✅ `testmesh load --users` | ✅ VU configuration UI | |
| View real-time metrics | ⚠️ Limited | ✅ Real-time charts | Dashboard: Visual |
| View results | ✅ `testmesh load results` | ✅ Results dashboard | Dashboard: Charts |
| Compare runs | ❌ Missing | ✅ Compare UI | **GAP: Need CLI command** |
| **Agents & Cloud** |
| List agents | ✅ `testmesh agents list` | ✅ Agents page | |
| View agent status | ✅ `testmesh agents status` | ✅ Agent status UI | |
| Target agent | ✅ `testmesh run --agent` | ✅ Agent selector | |
| **Remote Sync** |
| Configure remote | ✅ `testmesh remote set` | ✅ Settings UI | |
| Push flows | ✅ `testmesh push` | ✅ Auto-sync | Dashboard: Automatic |
| Pull flows | ✅ `testmesh pull` | ✅ Sync button | |
| Sync status | ✅ `testmesh status` | ✅ Sync indicator | |
| **Plugins** |
| List plugins | ✅ `testmesh plugin list` | ✅ Plugins page | |
| Install plugin | ✅ `testmesh plugin install` | ✅ Install button | |
| Create plugin | ✅ `testmesh plugin init` | ❌ N/A | CLI-only: Scaffolding |
| Test plugin | ✅ `testmesh plugin test` | ❌ N/A | CLI-only |
| **Visual Editor** |
| Drag & drop | ❌ N/A | ✅ Visual editor | Dashboard-only |
| Node palette | ❌ N/A | ✅ Node palette | Dashboard-only |
| Properties panel | ❌ N/A | ✅ Properties panel | Dashboard-only |
| Canvas zoom/pan | ❌ N/A | ✅ Canvas controls | Dashboard-only |
| YAML ↔ Visual | ✅ Edit YAML | ✅ Switch modes | Both edit same format |
| **Collaboration** |
| Real-time editing | ❌ N/A | ❌ Skipped | Not in v1.0 |
| Comments | ❌ N/A | ❌ Skipped | Not in v1.0 |
| Activity feed | ❌ Missing | ✅ Activity feed | **GAP: Consider adding** |
| **CI/CD Integration** |
| Run in CI | ✅ CLI in scripts | ⚠️ Via API | CLI: Preferred |
| Git hooks | ✅ Pre-commit | ❌ N/A | CLI-only |
| Programmatic API | ✅ REST API calls | ✅ REST API | Both use same API |

---

## Summary Statistics

| Category | Both | CLI Only | Dashboard Only | Missing from CLI | Missing from Dashboard |
|----------|------|----------|----------------|------------------|------------------------|
| **Core Features** | 45 | 8 | 12 | 3 | 2 |
| **Visual Features** | 0 | 0 | 10 | N/A | N/A |
| **Automation** | 0 | 5 | 0 | N/A | N/A |

---

## Identified Gaps

### 🔴 Critical Gaps (Must Fix)

None - all critical features have parity.

### 🟡 Important Gaps (Should Fix)

**CLI Missing**:
1. **Compare executions** - `testmesh compare <exec1> <exec2>`
   ```bash
   testmesh compare exec_123 exec_456
   # Output: Side-by-side diff
   ```

2. **Mock analytics** - `testmesh mock stats <mock_id>`
   ```bash
   testmesh mock stats payment-mock
   # Output: Request count, endpoints hit, etc.
   ```

3. **Data file preview** - `testmesh data preview <file>`
   ```bash
   testmesh data preview users.csv
   # Output: First 10 rows
   ```

4. **Load test comparison** - `testmesh load compare <run1> <run2>`
   ```bash
   testmesh load compare load_123 load_456
   # Output: Performance comparison
   ```

**Dashboard Missing**:
1. **Plugin scaffolding** - Create plugin from UI template
   - Alternative: Link to CLI command in UI

2. **Pre-commit hooks** - Git integration UI
   - Alternative: Documentation only

### 🟢 Nice to Have (Future)

**CLI Missing**:
- Workspace sharing management
- Activity feed viewer
- Visual flow preview (ASCII art?)

**Dashboard Missing**:
- Terminal emulator for CLI commands
- Git integration UI

---

## Recommended CLI Commands to Add

### 1. Compare Executions

```bash
# Compare two executions
testmesh compare exec_123 exec_456

# Output:
Comparing Executions
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Execution #123 (Failed)    vs.    Execution #456 (Success)

Duration:  3.2s             vs.    2.8s
Status:    Failed           vs.    Success
Failed at: charge_card      vs.    N/A

Step Comparison:
  setup:        0.2s ✓      vs.    0.2s ✓
  create_cart:  0.3s ✓      vs.    0.3s ✓
  add_items:    0.5s ✓      vs.    0.4s ✓
  charge_card:  2.1s ✗      vs.    0.8s ✓  ← Difference

Variable Differences:
  total_amount: 9999        vs.    4999  ← Different

Response Differences:
  Status:       402         vs.    200   ← Different
```

### 2. Mock Analytics

```bash
# View mock server analytics
testmesh mock stats payment-mock

# Output:
Mock Server: payment-mock
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Status:  Active
URL:     https://mock.testmesh.com/abc123

Usage (Last 24h):
  Total Requests: 145

Endpoints:
  POST /charge    89 requests  61%
  POST /refund    34 requests  23%
  GET  /balance   22 requests  16%

Response Times:
  Avg: 12ms    P95: 45ms    P99: 89ms

Errors: 3 (2%)
```

### 3. Data File Preview

```bash
# Preview CSV/JSON data file
testmesh data preview users.csv

# Output:
Preview: users.csv
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Rows: 100    Columns: 4

email                | name       | role   | amount
──────────────────────────────────────────────────
user1@example.com    | John Doe   | admin  | 100.00
user2@example.com    | Jane Doe   | user   | 50.00
user3@example.com    | Bob Smith  | user   | 25.00
...

Showing first 3 of 100 rows
Use --limit to show more
```

### 4. Load Test Comparison

```bash
# Compare two load test runs
testmesh load compare run_123 run_456

# Output:
Load Test Comparison
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Run #123 (Yesterday)       vs.    Run #456 (Today)

Users:     100              vs.    100
Duration:  10m              vs.    10m

Response Time:
  Avg:     456ms            vs.    412ms  ↓ 10% faster
  P95:     891ms            vs.    734ms  ↓ 18% faster
  P99:     1.2s             vs.    1.0s   ↓ 17% faster

Throughput:
  RPS:     145              vs.    167    ↑ 15% higher

Success Rate:
  Rate:    98.5%            vs.    99.2%  ↑ 0.7% better

Errors:
  Total:   19               vs.    12     ↓ 37% fewer
```

### 5. Activity Feed

```bash
# View recent activity
testmesh activity

# Output:
Recent Activity
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
2m ago   John Doe ran "checkout-flow" ✓
5m ago   Jane Smith edited "payment-flow"
10m ago  Bob Johnson created "user-registration"
15m ago  Alice Wong deleted "old-test-flow"
1h ago   John Doe ran "load-test-checkout" ✓

Use --limit to show more
Use --user to filter by user
```

---

## Implementation Priority

### Phase 1 (v1.0) ✅
- ✅ All core features have parity
- ✅ Critical workflows work in both CLI and Dashboard
- ✅ No major gaps blocking users

### Phase 2 (v1.1)
- Add CLI commands:
  - `testmesh compare` - Compare executions
  - `testmesh data preview` - Preview data files
  - `testmesh mock stats` - Mock analytics
  - `testmesh activity` - Activity feed

- Add Dashboard features:
  - Plugin creation wizard (optional)

### Phase 3 (v1.2)
- Add CLI commands:
  - `testmesh load compare` - Compare load tests
  - `testmesh workspace share` - Share workspace

- Add Dashboard features:
  - Terminal emulator
  - Git integration UI

---

## Design Principles

### When to Use CLI

✅ **Automation & Scripting**
- CI/CD pipelines
- Pre-commit hooks
- Automated testing
- Batch operations

✅ **Fast Iteration**
- Watch mode (auto-rerun)
- Quick validation
- Local development

✅ **Git Workflow**
- Version control integration
- Diff viewing
- Merge conflict resolution

✅ **Power Users**
- Keyboard-driven workflow
- Scriptable operations
- Programmatic access

### When to Use Dashboard

✅ **Visual Work**
- Drag-and-drop flow building
- Request builder forms
- Response visualization
- Real-time monitoring

✅ **Exploration**
- Browsing collections
- Searching flows
- Viewing results
- Debugging failures

✅ **Collaboration**
- Sharing flows
- Workspace management
- Activity tracking

✅ **Analytics**
- Charts and graphs
- Trends over time
- Performance dashboards

---

## Example Workflows

### Workflow 1: Local Development

**Using CLI** ✅ Recommended
```bash
# Initialize project
testmesh init my-tests
cd my-tests

# Create flow
testmesh create flow checkout.yaml

# Edit in VS Code
code flows/checkout.yaml

# Watch mode - auto-run on save
testmesh watch flows/checkout.yaml

# Commit when ready
git add flows/checkout.yaml
git commit -m "Add checkout flow"
```

### Workflow 2: Building Complex Flow

**Using Dashboard** ✅ Recommended
1. Open visual editor
2. Drag HTTP node from palette
3. Configure in properties panel
4. Connect nodes visually
5. Test with "Send" button
6. View pretty-printed response
7. Save flow

### Workflow 3: Debugging Failed Test

**Using Dashboard** ✅ Recommended
1. View execution list
2. Click failed execution
3. See timeline with failed step highlighted
4. Click "Logs" tab - search error
5. Click "Network" tab - inspect request/response
6. Click "Variables" tab - check state
7. Click "Cleanup" tab - verify cleanup

### Workflow 4: CI/CD Integration

**Using CLI** ✅ Recommended
```bash
# In CI pipeline
testmesh push                    # Push latest flows
testmesh run --tag smoke --wait  # Run and wait
testmesh results --format junit  # Export for CI
```

### Workflow 5: Load Testing

**Using Dashboard** ✅ Recommended
1. Open flow
2. Click "Load Test" button
3. Configure VUs, ramp-up, duration with visual chart
4. Click "Start"
5. Watch real-time metrics update
6. View results with charts
7. Compare with previous runs

---

## Summary

### ✅ Feature Parity Status

**Overall**: **95% parity** for core features

**Gaps**:
- 4 CLI commands to add (compare, stats, preview, activity)
- 2 Dashboard features to add (plugin wizard, git UI)

**Recommendation**: Ship v1.0 as-is, add remaining features in v1.1

### ✅ Both Interfaces Are Complete

- **CLI**: Perfect for automation, local dev, CI/CD
- **Dashboard**: Perfect for visual work, exploration, debugging

Users can choose based on task and preference! 🎯
