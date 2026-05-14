# metrics-generator

Go library and CLI for generating [kube-burner](https://github.com/kube-burner/kube-burner) metrics YAML profiles from type-safe PromQL builders.

## Layout

```
pkg/metrics/          # importable library
  query.go            # fluent PromQL query builder (Q, Rate, Agg, …)
  filters.go          # label filter helpers (NSExact, NSIn, Filters, …)
  types.go            # typed enums: Metric, GroupBy, NodeRole, AggFunc, …
  generator.go        # Generator — collects queries, emits YAML
  profile_*.go        # pre-built profiles

cmd/metrics-generator/ # CLI binary
  main.go
```

## CLI

```bash
make build
./metrics-generator                  # write all profiles to current dir
./metrics-generator /path/to/output  # write to specific directory
./metrics-generator --help
```

## Library usage

```go
import "github.com/kube-burner/metrics-generator/pkg/metrics"

g := &metrics.Generator{}

// simple query
g.AddQuery("containerCPU",
    metrics.Q(metrics.MetricContainerCPU, metrics.Filters(
        metrics.NSIn("openshift-.*", "cilium"),
        `name!=""`,
    )).IRate(metrics.Rate2m).Multiply("100").
        Agg(metrics.AggSum, metrics.GroupByPod, metrics.GroupByNamespace, metrics.GroupByNode),
)

// node-role join
g.AddQuery("workerNodeLoad",
    metrics.Q(metrics.MetricNodeLoad1, "").
        MultiplyOnGroupLeft(
            []metrics.GroupBy{metrics.GroupByInstance},
            metrics.NodeRoleLabelReplace(metrics.RoleWorker),
        ),
)

// use a pre-built profile
metrics.BuildMetricsProfile(g)

out, _ := g.Generate() // returns []byte YAML
```

## Pre-built profiles

| Profile function                    | Output file                    |
|-------------------------------------|--------------------------------|
| `BuildMetricsProfile`               | `metrics.yml`                  |
| `BuildBuildFarmProfile`             | `build-farm-metrics.yml`       |
| `BuildKueueProfile`                 | `kueue-metrics.yml`            |
| `BuildMetricsAggregatedProfile`     | `metrics-aggregated.yml`       |
| `BuildMetricsReportStackroxProfile` | `metrics-report-stackrox.yml`  |

## Development

```bash
make test        # run all tests
make build       # build binary to ./metrics-generator
```

Golden files live in `pkg/metrics/testdata/`. Update them when query output intentionally changes.
