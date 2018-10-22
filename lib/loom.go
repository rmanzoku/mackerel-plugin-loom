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
	return nil
}

// MetricKeyPrefix interface for PluginWithPrefix
func (p LoomPlugin) MetricKeyPrefix() string {
	return p.Prefix
}

// Do do doo
func Do() {
	var (
		optPrefix   = flag.String("metric-key-prefix", "", "Metric key prefix")
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
