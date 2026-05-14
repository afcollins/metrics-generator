package main

import (
	"fmt"
	"os"
	"path/filepath"
)

var profiles = []struct {
	filename string
	build    func(*Generator)
}{
	{"metrics.yml", BuildMetricsProfile},
	{"build-farm-metrics.yml", BuildBuildFarmProfile},
	{"kueue-metrics.yml", BuildKueueProfile},
	{"metrics-aggregated.yml", BuildMetricsAggregatedProfile},
	{"metrics-report-stackrox.yml", BuildMetricsReportStackroxProfile},
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
		g := &Generator{}
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
