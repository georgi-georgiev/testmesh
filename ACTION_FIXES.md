# Action Fixes: Transform & Assert

## Issues Fixed

Both `transform` and `assert` actions had problems handling interpolated data from previous steps. They've been completely rewritten to be more robust.

---

## Transform Action (`actions/transform.go`)

### What Was Wrong
- Couldn't properly handle data from step outputs
- JSONPath extraction was too strict
- No support for direct field access

### What Was Fixed

**1. Better Input Type Handling**
```go
// Now handles:
- map[string]interface{}
- models.OutputData
- Any type (via JSON marshal/unmarshal)
```

**2. Two Path Formats**
```yaml
transforms:
  # JSONPath format (for nested data)
  city: "$.body.address.city"

  # Direct field access (for step outputs)
  user_id: "$user_id"
  name: "$user_name"

  # Static values
  status: "active"
```

**3. Better Error Messages**
- Shows available input fields when path fails
- Logs input/output field counts

### Example Usage

```yaml
- id: fetch_user
  action: http_request
  config:
    method: GET
    url: "https://api.example.com/users/1"
  output:
    user_id: "$.id"
    user_name: "$.name"

- id: transform_user
  action: transform
  config:
    input: ${fetch_user}
    transforms:
      # Direct field from output
      id: "$user_id"
      name: "$user_name"

      # JSONPath from body
      email: "$.body.email"
      city: "$.body.address.city"

      # Static value
      status: "active"
```

---

## Assert Action (`actions/assert.go`)

### What Was Wrong
- Couldn't handle interpolated data from other steps
- Limited type support
- Poor error messages

### What Was Fixed

**1. Robust Type Conversion**
```go
// Now handles:
- models.OutputData
- map[string]interface{}
- JSON strings (auto-parses)
- Any type (via JSON marshal/unmarshal)
```

**2. Better Debugging**
- Logs available fields before running assertions
- Shows available fields in error messages
- Validates assertion list isn't empty

**3. Improved Error Messages**
```
assertion failed: id != nil: assertion failed (available fields: [id name email status])
```

### Example Usage

```yaml
- id: transform_user
  action: transform
  config:
    input: ${previous_step}
    transforms:
      id: "$user_id"
      name: "$user_name"
      status: "active"

- id: validate
  action: assert
  config:
    data: ${transform_user}
    assertions:
      - id == 1
      - name != ""
      - status == "active"
```

---

## Testing

### Test Example

Created `examples/transform-assert-demo.yaml` that demonstrates:

1. ✅ Fetch data from API
2. ✅ Transform with direct field access (`$field`)
3. ✅ Assert on transformed data
4. ✅ Transform with JSONPath (`$.path.to.field`)
5. ✅ Assert on nested data
6. ✅ Static values in transforms

### Run the Test

**Via CLI:**
```bash
cd api/cmd/testmesh
./testmesh validate ../../../examples/transform-assert-demo.yaml
./testmesh run ../../../examples/transform-assert-demo.yaml
```

**Via API:**
1. Start API server: `cd api && go run main.go`
2. Start web UI: `cd web && pnpm dev`
3. Visit http://localhost:3000
4. Create flow from `examples/transform-assert-demo.yaml`
5. Run it and watch it succeed! ✅

---

## Breaking Changes

### Transform Action

**Old format (broken):**
```yaml
transforms:
  id: "$.user_id"  # This assumed input was nested
```

**New format (working):**
```yaml
transforms:
  # For step output fields
  id: "$user_id"   # Direct access

  # For nested data
  id: "$.body.user.id"  # JSONPath
```

**Migration:** If your paths start with `$.`, they now extract from nested JSON. For direct field access from step outputs, use `$field` instead.

### Assert Action

No breaking changes - it now handles more input types automatically!

---

## Summary

| Action | Status Before | Status After | Breaking Changes |
|--------|---------------|--------------|------------------|
| `transform` | ⚠️ Broken | ✅ Working | Yes - path syntax |
| `assert` | ⚠️ Broken | ✅ Working | No |

Both actions now:
- ✅ Handle interpolated data correctly
- ✅ Support multiple input types
- ✅ Provide better error messages
- ✅ Log debug information
- ✅ Work with step outputs

**All control flow actions are now functional!** 🎉
