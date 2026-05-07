# Variables
BINARY_NAME=parsley-cli
CMD_DIR=./cmd/$(BINARY_NAME)
BUILD_DIR=build
VERSION=$(shell git describe --tags --abbrev=0 2>/dev/null || echo "0.0.0")
LDFLAGS=-ldflags "-X github.com/matzefriedrich/parsley/internal/utils.VersionString=$(VERSION)"

# Default target
all: build

# Phony targets
.PHONY: all build install test test-coverage lint lint-fix help clean

build: ## Build the parsley-cli binary
	mkdir -p $(BUILD_DIR)
	go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(CMD_DIR)

install: ## Install the parsley-cli from source
	go install $(LDFLAGS) $(CMD_DIR)

test: ## Run all tests
	go test ./...

test-coverage: ## Run tests and compute coverage
	go test -coverpkg=./... ./... -coverprofile=coverage.out
	go tool cover -func=coverage.out

lint: ## Run golangci-lint
	golangci-lint run

lint-fix: ## Run golangci-lint and apply fixes
	golangci-lint run --fix

clean: ## Clean build artifacts
	rm -rf $(BUILD_DIR)
	rm -f coverage.out

help: ## Show this help message
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'
