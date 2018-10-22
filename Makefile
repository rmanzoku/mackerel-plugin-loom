LATEST_TAG=$(shell git describe --abbrev=0 --tags))
NAME="mackerel-plugin-loom"

.PHONY: run help meta build

run: # Run trial
	@go run main.go -url http://172.31.32.97:9999/metrics

build: # Build
	@go build

meta: # Show meta
	@MACKEREL_AGENT_PLUGIN_META=1 go run main.go

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
