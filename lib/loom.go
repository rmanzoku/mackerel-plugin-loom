package mploom

import (
	"flag"
	"fmt"
	"strconv"

	mp "github.com/mackerelio/go-mackerel-plugin"
	dto "github.com/prometheus/client_model/go"
	"github.com/prometheus/prom2json"
)

// LoomPlugin is plugin for Loom SDK
type LoomPlugin struct {
	Tempfile string
	URL      string
	Prefix   string
}

func fetchFamilies(url string) ([]*prom2json.Family, error) {

	mfChan := make(chan *dto.MetricFamily, 1024)

	go func() {
		err := prom2json.FetchMetricFamilies(url, mfChan, "", "", true)
		if err != nil {
			panic(err)
		}
	}()

	result := []*prom2json.Family{}
	for mf := range mfChan {
		result = append(result, prom2json.NewFamily(mf))
	}

	return result, nil
}

// FetchMetrics interface for mackerelplugin
func (p LoomPlugin) FetchMetrics() (map[string]float64, error) {
	ret := map[string]float64{}
	families, err := fetchFamilies(p.URL)
	if err != nil {
		return nil, err
	}

	for _, f := range families {
		fmt.Println(f.Name)
		switch f.Name {
		case "loomchain_application_block_count":
			ret["block_count"], err = strconv.ParseFloat(f.Metrics[0].(prom2json.Metric).Value, 64)

		case "loomchain_application_checktx_latency_microseconds":
			ret["latency_check_tx"], err = strconv.ParseFloat(f.Metrics[0].(prom2json.Summary).Sum, 64)

		case "loomchain_application_commit_block_latency_microseconds":
			ret["latency_commit_block"], err = strconv.ParseFloat(f.Metrics[0].(prom2json.Summary).Sum, 64)

		case "loomchain_application_delivertx_latency_microseconds":
			ret["latency_deliver_tx"], err = strconv.ParseFloat(f.Metrics[0].(prom2json.Summary).Sum, 64)

		case "loomchain_application_evm_transaction_count":
			ret["evm_tx_count"], err = strconv.ParseFloat(f.Metrics[0].(prom2json.Metric).Value, 64)

		case "loomchain_application_evm_tx_gas_cost":
			ret["evm_gas_cost"], err = strconv.ParseFloat(f.Metrics[0].(prom2json.Summary).Sum, 64)

		case "loomchain_application_evmtx_latency_microseconds":
			ret["latency_evm_tx"], err = strconv.ParseFloat(f.Metrics[0].(prom2json.Summary).Sum, 64)

		case "loomchain_application_request_count":
			ret["req_application"], err = strconv.ParseFloat(f.Metrics[0].(prom2json.Metric).Value, 64)

		}
	}

	return ret, err
}

// GraphDefinition interface for mackerelplugin
func (p LoomPlugin) GraphDefinition() map[string]mp.Graphs {
	return map[string]mp.Graphs{
		"block": {
			Label: "Block count",
			Unit:  "integer",
			Metrics: []mp.Metrics{
				{Name: "block_count", Label: "Block Count", Diff: false},
			},
		},
		"evm.tx": {
			Label: "EVM transaction count",
			Unit:  "integer",
			Metrics: []mp.Metrics{
				{Name: "evm_tx_count", Label: "EVM transaction count", Diff: true},
			},
		},
		"evm.gas": {
			Label: "EVM Gas cost",
			Unit:  "integer",
			Metrics: []mp.Metrics{
				{Name: "evm_gas_cost", Label: "Gas cost", Diff: true},
			},
		},
		"request": {
			Label: "Loomchain request count",
			Unit:  "integer",
			Metrics: []mp.Metrics{
				{Name: "req_application", Label: "Application", Diff: false},
				{Name: "req_query", Label: "Query", Diff: false},
				{Name: "req_tx", Label: "Transaction", Diff: false},
				{Name: "req_evm_tx", Label: "EVM transaction", Diff: false},
			},
		},
		"latency": {
			Label: "Loomchain latency (microseconds)",
			Unit:  "integer",
			Metrics: []mp.Metrics{
				{Name: "latency_check_tx", Label: "Check tx", Diff: true},
				{Name: "latency_commit_block", Label: "Commit block", Diff: false},
				{Name: "latency_deliver_tx", Label: "Deliver tx", Diff: false},
				{Name: "latency_evm_tx", Label: "EVM tx", Diff: false},
				{Name: "latency_query", Label: "Query request", Diff: false},
				{Name: "latency_tx", Label: "Tx request", Diff: false},
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
