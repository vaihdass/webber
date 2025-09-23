###############################################################################
### WEBBER — Go web development library & toolchain                         ###
###############################################################################

###############################################################################
# VARIABLES
###############################################################################
SHELL                := /usr/bin/env bash
.DEFAULT_GOAL        := help

#--- Project configuration
MAIN_BRANCH          ?= master
PROJECT_MODULE       ?= github.com/vaihdass/webber

#--- Test exclusions (folders to skip during testing)
# Examples: scripts|tools|vendor|examples
# Leave empty to test ALL packages
TEST_EXCLUDE_DIRS    := vendor

#--- Tool versions
GOLANGCI_LINT_VERSION?= v2.5.0

#--- Directory paths
BIN_DIR              := $(CURDIR)/bin
COVERAGE_DIR         := $(CURDIR)/coverage

#--- File paths
COVERAGE_FILE        := $(COVERAGE_DIR)/coverage.out
COVERAGE_HTML        := $(COVERAGE_DIR)/coverage.html

###############################################################################
# PUBLIC COMMANDS (visible in help - sorted alphabetically)
###############################################################################

#--- Help command (default target)

help: ## Show all available tasks (default target)
	@echo -e "\n\033[1mAvailable tasks\033[0m"; \
	grep -E '^[a-zA-Z0-9_.-]+:.*##' $(MAKEFILE_LIST) \
		| grep -vE '^\.' \
		| awk 'BEGIN {FS=":.*##"} {printf "  \033[36m%-30s\033[0m %s\n", $$1, $$2}' ; \
	echo ""

#--- Compilation & work in progress commands

build: ## Build all packages in the module
	go build ./...

run: ## Run cmd/test/main.go for local drafting
	go run ./cmd/test/main.go

pre: build format lint test ## Prepare new version (build, format, lint & test)

#--- Code formatting

format: deps ## Run all code formatters
	gofmt -s -w .
	goimports -w .
	go mod tidy

#--- Linting

lint: deps ## Run all linters (diff-only with origin/main)
	$(BIN_DIR)/golangci-lint run --new-from-rev=origin/$(MAIN_BRANCH)

lint-full: deps ## Run all linters (full codebase)
	$(BIN_DIR)/golangci-lint run

#--- Testing

test: ## Run Go tests with race & coverage, output HTML
	mkdir -p $(COVERAGE_DIR)
	@if [ -z "$(TEST_EXCLUDE_DIRS)" ]; then \
		echo "Running tests on ALL packages..."; \
		go test -race -count=1 -coverprofile='$(COVERAGE_FILE)' -covermode=atomic ./...; \
	else \
		echo "Running tests with exclusions: $(TEST_EXCLUDE_DIRS)"; \
		go test -race -count=1 -coverprofile='$(COVERAGE_FILE)' -covermode=atomic \
			$$(go list ./... | grep -vE '($(TEST_EXCLUDE_DIRS))'); \
	fi
	go tool cover -func='$(COVERAGE_FILE)' | grep ^total | tr -s '\t'
	go tool cover -html=$(COVERAGE_FILE) -o $(COVERAGE_HTML)

#--- Setup and dependencies

deps: ## Install binary and Go dependencies (skips if already installed)
	@if [ ! -d "$(BIN_DIR)" ] || [ ! -f "$(BIN_DIR)/golangci-lint" ]; then \
		echo "🔧 Installing tools..."; \
		$(MAKE) .deps; \
	fi

force-deps: .deps ## Force reinstall all binary dependencies

tidy: ## go mod tidy & vendor
	go mod tidy
	go mod vendor

clean: .clean_bin ## Remove built artefacts, binary dependencies & coverage files
	go clean ./...
	rm -rf $(COVERAGE_FILE) $(COVERAGE_HTML)
	go mod tidy

###############################################################################
# INTERNAL TASKS (hidden from `make help` – all names start with ".")
###############################################################################

.deps: .clean_bin
	go mod tidy
	go mod download
	go mod vendor
	@mkdir -p $(BIN_DIR)
	GOBIN=$(BIN_DIR) go install golang.org/x/tools/cmd/goimports@latest
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(BIN_DIR) $(GOLANGCI_LINT_VERSION)

.clean_bin:
	rm -rf $(BIN_DIR)

#--- PHONY DECLARATIONS
.PHONY: help build run pre format lint lint-full test deps force-deps clean .deps .clean_bin
