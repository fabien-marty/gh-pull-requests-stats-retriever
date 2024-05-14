SHELL:=/bin/bash

FIX=1
CGO_ENABLED=1
COMMON_TEST_OPTIONS=-race -cover -covermode=atomic
CMDS=cmd/gh-pr-stats-retriever/gh-pr-stats-retriever
BUILDARGS=

default: help

.PHONY: build
build: $(CMDS) ## Build Go binaries

cmd/gh-pr-stats-retriever/gh-pr-stats-retriever: $(shell find cmd/gh-pr-stats-retriever internal -type f -name '*.go')
	cd `dirname $@` && go build $(BUILDARGS) -o `basename $@` *.go

.PHONY: gofmt
gofmt:
	@if test "$(FIX)" = "1"; then \
		set -x ; gofmt -s -w . ;\
	else \
		set -x ; gofmt -s -d . ;\
	fi

.PHONY: golangcilint
golangcilint: tmp/bin/golangci-lint
	@if test "$(FIX)" = "1"; then \
		set -x ; $< run --fix --timeout 10m;\
	else \
		set -x ; $< run --timeout 10m;\
	fi

tmp/venv/bin/activate:
	@mkdir -p tmp
	python3 -m venv tmp/venv
	source tmp/venv/bin/activate && pip install -r requirements.txt

.PHONY: govet
govet:
	go vet ./...

.PHONY: _unit_test
_unit_test: ## Execute all unit tests
	@rm -Rf covdatafiles/unit
	@mkdir -p covdatafiles/unit 
	go test $(COMMON_TEST_OPTIONS) ./... -args -test.gocoverdir=$$(pwd)/covdatafiles/unit

.PHONY: _prepare_coverage
_prepare_coverage:
	@mkdir -p covdatafiles/unit

.PHONY: _merge_coverage
_merge_coverage:
	go tool covdata textfmt -i=./covdatafiles/unit -o coverage.out
	rm -Rf covdatafiles

.PHONY: test
test: _prepare_coverage _unit_test _merge_coverage ## Execute all tests 

.PHONY: html-coverage
html-coverage: test ## Build html coverage
	go tool cover -html coverage.out -o cover.html

.PHONY: lint
lint: govet gofmt golangcilint ## Lint the code (also fix the code if FIX=1, default)

tmp/bin/golangci-lint:
	@mkdir -p tmp/bin
	cd tmp/bin && curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b . v1.56.2 && chmod +x `basename $@`

.PHONY: clean
clean: _cmd_clean ## Clean the repo
	rm -f coverage.out
	rm -Rf covdatafiles
	rm -Rf tmp
	rm -Rf build
	rm -f cover.html

.PHONY: _cmd_clean
_cmd_clean:
	rm -f $(CMDS)

.PHONY: doc
doc: build tmp/venv/bin/activate ## Generate documentation
	source tmp/venv/bin/activate && jinja-tree .

.PHONY: help
help::
	@# See https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
