package metrics

import (
	"fmt"
	"strings"
)

// ---------------------------------------------------------------------------
// Label filter builders — return filter strings for use inside Q(metric, filter)
// ---------------------------------------------------------------------------

// NSExact returns namespace="name".
func NSExact(name string) string {
	return fmt.Sprintf(`namespace="%s"`, name)
}

// NSRegex returns namespace=~"pattern".
func NSRegex(pattern string) string {
	return fmt.Sprintf(`namespace=~"%s"`, pattern)
}

// NSNotRegex returns namespace!~"pattern".
func NSNotRegex(pattern string) string {
	return fmt.Sprintf(`namespace!~"%s"`, pattern)
}

// NSPrefix returns namespace=~"prefix.*".
func NSPrefix(prefix string) string {
	return NSRegex(prefix + ".*")
}

// NSIn returns namespace="name" for one namespace, or namespace=~"a|b|c" for many.
func NSIn(namespaces ...string) string {
	if len(namespaces) == 1 {
		return NSExact(namespaces[0])
	}
	return NSRegex(strings.Join(namespaces, "|"))
}

// Filters joins multiple filter clauses with commas.
// Use to compose namespace filters with other label matchers:
//
//	Q(MetricContainerCPU, Filters(NSIn("openshift-.*", "cilium"), `name!=""`))
func Filters(parts ...string) string {
	return strings.Join(parts, ",")
}
