---
name: feature-request-multiply-on-group-left
description: Feature request for metrics-generator to add MultiplyOnGroupLeft and related binary op join methods
metadata: 
  node_type: memory
  type: project
  originSessionId: c4e1f31d-24ba-4351-a263-fd83902e6380
---

metrics-generator (`github.com/kube-burner/metrics-generator`) needs binary operator join methods for `* on (labels) group_left` patterns.

**Why:** performance-dashboards uses this pattern heavily for node-role joins:
```promql
node_cpu_seconds_total{mode!="idle"} * on (instance) group_left label_replace(kube_node_role{role="worker"}, "instance", "$1", "node", "(.*)")
```
and simpler variant:
```promql
container_threads{container!=""} * on (node) group_left kube_node_role{role="worker"}
```
Without these methods, ~30+ queries in ocp_performance.go must stay as `Raw()`.

**How to apply:** Add to `query.go`:
- `MultiplyOnGroupLeft(labels []GroupBy, right *Query) *Query` — produces `expr * on (labels) group_left right`
- `OnGroupLeft(op string, labels []GroupBy, right *Query) *Query` — generic version for any binary op

Existing `NodeRoleLabelReplace()` already builds the right-hand side. Combined usage would be:
```go
Q(MetricNodeCPU, `mode!="idle"`).
    MultiplyOnGroupLeft([]GroupBy{GroupByInstance}, NodeRoleLabelReplace(RoleWorker))
```

Related: [[metrics-generator-import]]
