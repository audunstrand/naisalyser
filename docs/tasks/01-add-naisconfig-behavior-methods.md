# Task 01: Add Behavior Methods to NaisConfig

## Deliverable
NaisConfig struct has helper methods that encapsulate business logic for querying its state.

## Context
Currently `NaisConfig` is an anemic data container. Business logic like "is this app stateful?" is scattered across files like `graph/builder.go` (lines 51-56). This violates encapsulation and causes code duplication.

**Problem code in `internal/graph/builder.go:51-56`:**
```go
stateful := false
if analysis.NaisConfig != nil && analysis.NaisConfig.Spec.GCP != nil {
    if len(analysis.NaisConfig.Spec.GCP.SQLInstances) > 0 || len(analysis.NaisConfig.Spec.GCP.Buckets) > 0 {
        stateful = true
    }
}
```

This logic belongs on `NaisConfig` itself.

## Key Decisions and Principles
- Methods return sensible defaults when nested fields are nil (no nil pointer panics)
- Methods are simple predicates or accessors, not complex logic
- Keep methods on the struct they query (Tell Don't Ask principle)

## Delivers
Four new methods on `NaisConfig` that encapsulate common queries about app configuration.

## Acceptance Criteria
- `HasDatabase() bool` → returns true if `Spec.GCP.SQLInstances` has entries
- `HasBuckets() bool` → returns true if `Spec.GCP.Buckets` has entries  
- `IsStateful() bool` → returns true if app has database OR buckets
- `GetInboundApps() []string` → returns list of app names from inbound access policy
- `GetOutboundApps() []string` → returns list of app names from outbound access policy
- All methods handle nil pointers safely (return false/empty slice)
- Existing code continues to compile

## Dependencies
None - this is a foundational task.

## Related Code

**File to modify:** `internal/analyzer/analyzer.go`

**Add methods after line 127** (after `NaisBucket` struct definition, before `Dependency` struct):

```go
// NaisConfig methods - add these after line 127

// HasDatabase returns true if the app has Cloud SQL instances configured
func (n *NaisConfig) HasDatabase() bool {
    if n == nil || n.Spec.GCP == nil {
        return false
    }
    return len(n.Spec.GCP.SQLInstances) > 0
}

// HasBuckets returns true if the app has GCS buckets configured
func (n *NaisConfig) HasBuckets() bool {
    if n == nil || n.Spec.GCP == nil {
        return false
    }
    return len(n.Spec.GCP.Buckets) > 0
}

// IsStateful returns true if the app has persistent storage (database or buckets)
func (n *NaisConfig) IsStateful() bool {
    return n.HasDatabase() || n.HasBuckets()
}

// GetInboundApps returns the list of application names allowed to call this app
func (n *NaisConfig) GetInboundApps() []string {
    if n == nil || n.Spec.AccessPolicy == nil || n.Spec.AccessPolicy.Inbound == nil {
        return nil
    }
    apps := make([]string, 0, len(n.Spec.AccessPolicy.Inbound.Rules))
    for _, rule := range n.Spec.AccessPolicy.Inbound.Rules {
        apps = append(apps, rule.Application)
    }
    return apps
}

// GetOutboundApps returns the list of application names this app can call
func (n *NaisConfig) GetOutboundApps() []string {
    if n == nil || n.Spec.AccessPolicy == nil || n.Spec.AccessPolicy.Outbound == nil {
        return nil
    }
    apps := make([]string, 0, len(n.Spec.AccessPolicy.Outbound.Rules))
    for _, rule := range n.Spec.AccessPolicy.Outbound.Rules {
        apps = append(apps, rule.Application)
    }
    return apps
}
```

## Verification

```bash
# Build should succeed
go build ./...

# Test the methods work (create a simple test)
cat > internal/analyzer/naisconfig_test.go << 'EOF'
package analyzer

import "testing"

func TestNaisConfig_IsStateful_WithDatabase(t *testing.T) {
    cfg := &NaisConfig{
        Spec: NaisSpec{
            GCP: &NaisGCP{
                SQLInstances: []NaisSQLInstance{{Type: "POSTGRES_14"}},
            },
        },
    }
    if !cfg.IsStateful() {
        t.Error("expected IsStateful() to return true when database exists")
    }
}

func TestNaisConfig_IsStateful_Nil(t *testing.T) {
    var cfg *NaisConfig
    if cfg.IsStateful() {
        t.Error("expected IsStateful() to return false for nil config")
    }
}

func TestNaisConfig_GetInboundApps(t *testing.T) {
    cfg := &NaisConfig{
        Spec: NaisSpec{
            AccessPolicy: &NaisAccessPolicy{
                Inbound: &NaisAccessRules{
                    Rules: []NaisAccessRule{
                        {Application: "app-a"},
                        {Application: "app-b"},
                    },
                },
            },
        },
    }
    apps := cfg.GetInboundApps()
    if len(apps) != 2 || apps[0] != "app-a" || apps[1] != "app-b" {
        t.Errorf("unexpected inbound apps: %v", apps)
    }
}
EOF

# Run the test
go test ./internal/analyzer/... -v
```

## Files Changed
- `internal/analyzer/analyzer.go` (add ~40 lines of methods)
- `internal/analyzer/naisconfig_test.go` (new file, ~45 lines)
