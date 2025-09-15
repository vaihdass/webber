###############################################################################
# WEBBER – Go Library Makefile (Library + Go toolchain)
###############################################################################
SHELL                := /usr/bin/env bash
.DEFAULT_GOAL        := help

###############################################################################
# VARIABLES
###############################################################################

# ---------------------------------------------------------------------------
# Project configuration
# ---------------------------------------------------------------------------
MAIN_BRANCH          ?= main
PROJECT_MODULE       ?= github.com/vaihdass/webber

# ---------------------------------------------------------------------------
# Directory paths
# ---------------------------------------------------------------------------
BIN_DIR              := $(CURDIR)/bin
COVERAGE_DIR         := $(CURDIR)/coverage

# ---------------------------------------------------------------------------
# File paths
# ---------------------------------------------------------------------------
COVERAGE_FILE        := $(COVERAGE_DIR)/coverage.out
COVERAGE_HTML        := $(COVERAGE_DIR)/coverage.html

# ---------------------------------------------------------------------------
# Test exclusions (folders to skip during testing)
# Examples: scripts|tools|vendor|examples
# Leave empty to test ALL packages
# ---------------------------------------------------------------------------
TEST_EXCLUDE_DIRS    := scripts|tools|vendor|examples|testdata

# ---------------------------------------------------------------------------
# Tool versions
# ---------------------------------------------------------------------------
GOLANGCI_LINT_VERSION?= v2.1.6

###############################################################################
# PUBLIC COMMANDS (visible in help - sorted alphabetically)
###############################################################################

# ---------------------------------------------------------------------------
# Help command (default target)
# ---------------------------------------------------------------------------
help: ## Show all available tasks (default target)
	@echo -e "\n\033[1mAvailable tasks\033[0m"; \
	grep -E '^[a-zA-Z0-9_.-]+:.*##' $(MAKEFILE_LIST) \
		| grep -vE '^\.' \
		| awk 'BEGIN {FS=":.*##"} {printf "  \033[36m%-30s\033[0m %s\n", $$1, $$2}' ; \
	echo ""

# ---------------------------------------------------------------------------
# Build and compilation commands
# ---------------------------------------------------------------------------
build: ## Build all packages in the module
	go build -v ./...

# ---------------------------------------------------------------------------
# Clean up commands
# ---------------------------------------------------------------------------
clean: ## Remove built artefacts & coverage files
	go clean ./...
	rm -rf $(BIN_DIR) $(COVERAGE_FILE) $(COVERAGE_HTML)
	go mod tidy

# ---------------------------------------------------------------------------
# Code formatting commands
# ---------------------------------------------------------------------------
format: ## Run all code formatters
	gofmt -s -w .
	goimports -w .
	go mod tidy

# ---------------------------------------------------------------------------
# Linting commands
# ---------------------------------------------------------------------------
lint: .lint-go ## Run all linters

lint-full: .lint-go-full ## Run all linters (full)

# ---------------------------------------------------------------------------
# Release preparation commands
# ---------------------------------------------------------------------------
prepare-release: ## Prepare release (format, lint, test)
	$(MAKE) format
	$(MAKE) lint
	$(MAKE) test

# ---------------------------------------------------------------------------
# Setup and dependency commands
# ---------------------------------------------------------------------------
deps: .deps ## Install Go modules and development tools

force-deps: ## Force reinstall all dependencies
	go mod tidy
	go mod download
	$(MAKE) .deps_tools

# ---------------------------------------------------------------------------
# Testing commands
# ---------------------------------------------------------------------------
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

###############################################################################
# INTERNAL TASKS (hidden from `make help` – all names start with ".")
###############################################################################

# ---------------------------------------------------------------------------
# Dependency management commands
# ---------------------------------------------------------------------------
.deps: ## Install Go modules and development tools (optimized: skips if already installed)
	go mod tidy
	go mod download
	@if [ ! -d "$(BIN_DIR)" ] || [ ! -f "$(BIN_DIR)/golangci-lint" ]; then \
		echo "🔧 Installing tools..."; \
		$(MAKE) .deps_tools; \
	else \
		echo "✅ Tools already installed"; \
	fi

.deps_tools:
	$(MAKE) .clean_tools
	@mkdir -p $(BIN_DIR)

	# -----------------------------------------------------------------------
	# Go-based toolchain
	# -----------------------------------------------------------------------
	GOBIN=$(BIN_DIR) go install golang.org/x/tools/cmd/goimports@latest
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(BIN_DIR) $(GOLANGCI_LINT_VERSION)

.clean_tools:
	rm -rf $(BIN_DIR)

# ---------------------------------------------------------------------------
# Linting sub-commands
# ---------------------------------------------------------------------------
.lint-go: ## Run golangci-lint (diff-only vs origin/main)
	@if [ -f "$(BIN_DIR)/golangci-lint" ]; then \
		$(BIN_DIR)/golangci-lint run --new-from-rev=origin/$(MAIN_BRANCH); \
	else \
		golangci-lint run --new-from-rev=origin/$(MAIN_BRANCH); \
	fi

.lint-go-full: ## Run golangci-lint (full run)
	@if [ -f "$(BIN_DIR)/golangci-lint" ]; then \
		$(BIN_DIR)/golangci-lint run; \
	else \
		golangci-lint run; \
	fi

###############################################################################
# PHONY DECLARATIONS
###############################################################################
.PHONY: help build clean format lint lint-full test \
		prepare-release deps force-deps \
		.deps .deps_tools .clean_tools .lint-go .lint-go-full
