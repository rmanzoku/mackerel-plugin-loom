package mploom

import (
	"flag"

	mp "github.com/mackerelio/go-mackerel-plugin"
)

// LoomPlugin is plugin for Loom SDK
type LoomPlugin struct {
	Tempfile string
	URL      string
	Prefix   string
}

// FetchMetrics interface for mackerelplugin
func (p LoomPlugin) FetchMetrics() (map[string]float64, error) {
	return nil, nil
}

// GraphDefinition interface for mackerelplugin
func (p LoomPlugin) GraphDefinition() map[string]mp.Graphs {
	return map[string]mp.Graphs{
		"block": {
			Label: "Block count",
			Unit:  "integer",
			Metrics: []mp.Metrics{
				{Name: "count", Label: "Count", Diff: false},
			},
		},
		"evm.gas": {
			Label: "EVM Gas cost",
			Unit:  "integer",
			Metrics: []mp.Metrics{
				{Name: "cost", Label: "Gas cost", Diff: false},
			},
		},
		"request": {
			Label: "Loomchain request count",
			Unit:  "integer",
			Metrics: []mp.Metrics{
				{Name: "application", Label: "Application", Diff: false},
				{Name: "query", Label: "Query", Diff: false},
				{Name: "tx", Label: "Transaction", Diff: false},
				{Name: "evm_tx", Label: "EVM transaction", Diff: false},
			},
		},
		"latency": {
			Label: "Loomchain latency (microseconds)",
			Unit:  "integer",
			Metrics: []mp.Metrics{
				{Name: "check_tx", Label: "Check tx", Diff: false},
				{Name: "commit_block", Label: "Commit block", Diff: false},
				{Name: "deliver_tx", Label: "Deliver tx", Diff: false},
				{Name: "evm_tx", Label: "EVM tx", Diff: false},
				{Name: "query", Label: "Query request", Diff: false},
				{Name: "tx", Label: "Tx request", Diff: false},
			},
		},
	}
}

// MetricKeyPrefix interface for PluginWithPrefix
func (p LoomPlugin) MetricKeyPrefix() string {
	return p.Prefix
}

// Do do doo
func Do() {
	var (
		optPrefix   = flag.String("metric-key-prefix", "loomchain", "Metric key prefix")
		optURL      = flag.String("url", "http://127.0.0.1:46658/metrics", "Loom RPC Bind Address")
		optTempfile = flag.String("tempfile", "", "Temp file name")
	)
	flag.Parse()

	var l LoomPlugin
	l.Prefix = *optPrefix
	l.URL = *optURL

	helper := mp.NewMackerelPlugin(l)
	helper.Tempfile = *optTempfile
	helper.Run()
}
