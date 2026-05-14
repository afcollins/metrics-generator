## Repo purpose
Generate kube-burner metrics YAML profiles from type-safe Go PromQL builders.

## Layout
```
pkg/metrics/      — importable library (package metrics)
  query.go        — Query builder: Q(), Rate(), Agg(), MultiplyOnGroupLeft(), …
  filters.go      — Label helpers: NSExact(), NSIn(), NSRegex(), Filters()
  types.go        — Metric/GroupBy/NodeRole/AggFunc/… constants
  generator.go    — Generator struct: AddQuery(), Generate() → []byte YAML
  profile_*.go    — Pre-built profiles: BuildMetricsProfile(), BuildKueueProfile(), …
  testdata/       — Golden YAML files; update when query output intentionally changes
cmd/metrics-generator/main.go  — CLI entry point (package main)
```

## Key patterns

### Query building
```go
Q(MetricNodeCPU, Filters(NSIn("openshift-.*", "cilium"), `name!=""`)).
    IRate(Rate2m).Multiply("100").
    Agg(AggSum, GroupByNode)
```

### Node-role joins
```go
Q(MetricNodeLoad1, "").MultiplyOnGroupLeft(
    []GroupBy{GroupByInstance}, NodeRoleLabelReplace(RoleWorker))
```

### Subquery rate (Grafana $interval pattern)
```go
Q(MetricNodeCPU, `mode!="idle"`).
    MultiplyOnGroupLeft(...).Paren().RateSubquery("$interval")
```

### Namespace filters
- `NSExact("stackrox")` → `namespace="stackrox"`
- `NSIn("openshift-.*", "cilium")` → `namespace=~"openshift-.*|cilium"`
- `NSRegex("kueue-scale-.+")` → `namespace=~"kueue-scale-.+"`
- `Filters(NSIn(...), `name!=""`)` → comma-joined

## Tests & golden files
- `go test ./...` — runs all tests
- Profile tests compare against `pkg/metrics/testdata/*.yml`
- When a query changes intentionally: update the matching golden file

## CLI
```
make build
./metrics-generator [output-dir]   # writes all profiles; default: .
./metrics-generator --help
```

## Adding a new profile
1. Create `pkg/metrics/profile_<name>.go` with `func Build<Name>Profile(g *Generator)`
2. Add a golden file to `pkg/metrics/testdata/`
3. Add a test in `profiles_test.go` using `compareProfiles()`
4. Register in `cmd/metrics-generator/main.go` profiles slice
