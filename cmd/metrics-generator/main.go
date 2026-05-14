package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/kube-burner/metrics-generator/pkg/metrics"
)

var profiles = []struct {
	filename string
	build    func(*metrics.Generator)
}{
	{"metrics.yml", metrics.BuildMetricsProfile},
	{"build-farm-metrics.yml", metrics.BuildBuildFarmProfile},
	{"kueue-metrics.yml", metrics.BuildKueueProfile},
	{"metrics-aggregated.yml", metrics.BuildMetricsAggregatedProfile},
	{"metrics-report-stackrox.yml", metrics.BuildMetricsReportStackroxProfile},
}

func main() {
	if len(os.Args) > 1 && (os.Args[1] == "-h" || os.Args[1] == "--help") {
		printHelp()
		return
	}

	outDir := "."
	if len(os.Args) > 1 {
		outDir = os.Args[1]
	}

	for _, p := range profiles {
		g := &metrics.Generator{}
		p.build(g)
		out, err := g.Generate()
		if err != nil {
			fmt.Fprintf(os.Stderr, "error generating %s: %v\n", p.filename, err)
			os.Exit(1)
		}
		path := filepath.Join(outDir, p.filename)
		if err := os.WriteFile(path, out, 0644); err != nil {
			fmt.Fprintf(os.Stderr, "error writing %s: %v\n", path, err)
			os.Exit(1)
		}
		fmt.Fprintf(os.Stderr, "wrote %d metrics → %s\n", g.Count(), path)
	}
}

func printHelp() {
	fmt.Print(`metrics-generator — generate kube-burner metrics YAML profiles

Usage:
  metrics-generator [output-dir]
  metrics-generator -h | --help

Arguments:
  output-dir   directory to write generated files (default: current directory)

Profiles generated:
  metrics.yml
  build-farm-metrics.yml
  kueue-metrics.yml
  metrics-aggregated.yml
  metrics-report-stackrox.yml
`)
}
