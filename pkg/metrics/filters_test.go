package metrics

import "testing"

func TestNSFilters(t *testing.T) {
	tests := []struct {
		name     string
		got      string
		expected string
	}{
		{"NSExact", NSExact("stackrox"), `namespace="stackrox"`},
		{"NSRegex", NSRegex("openshift-.*"), `namespace=~"openshift-.*"`},
		{"NSPrefix", NSPrefix("openshift-"), `namespace=~"openshift-.*"`},
		{"NSIn single", NSIn("stackrox"), `namespace="stackrox"`},
		{"NSIn multi", NSIn("openshift-monitoring", "stackrox", "cilium"), `namespace=~"openshift-monitoring|stackrox|cilium"`},
		{
			"Filters compose",
			Filters(NSIn("openshift-.*", "cilium"), `name!=""`, `container!="POD"`),
			`namespace=~"openshift-.*|cilium",name!="",container!="POD"`,
		},
		{
			"Q with NSExact",
			Q(MetricContainerMemoryRSS, Filters(NSExact("openshift-etcd"), `name!=""`)).String(),
			`container_memory_rss{namespace="openshift-etcd",name!=""}`,
		},
		{
			"Q with NSIn multi",
			Q(MetricContainerCPU, NSIn("openshift-monitoring", "openshift-ovn-kubernetes")).String(),
			`container_cpu_usage_seconds_total{namespace=~"openshift-monitoring|openshift-ovn-kubernetes"}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.expected {
				t.Errorf("got  %q\nwant %q", tt.got, tt.expected)
			}
		})
	}
}
