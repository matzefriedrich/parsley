# Variables
BINARY_NAME=parsley-cli
CMD_DIR=./cmd/$(BINARY_NAME)
BUILD_DIR=build
VERSION=$(shell git describe --tags --abbrev=0 2>/dev/null || echo "0.0.0")
LDFLAGS=-ldflags "-X github.com/matzefriedrich/parsley/internal/utils.VersionString=$(VERSION)"
SKILLS_DIR=.agents/skills

# Copies the shared agent skills into .<agent>/skills, replacing any previous copy
define install_skills
	rm -rf .$(1)/skills
	mkdir -p .$(1)
	cp -R $(SKILLS_DIR) .$(1)/skills
endef

# Default target
all: build

# Phony targets
.PHONY: all build install test test-coverage lint lint-fix help clean \
	install-agents install-claude install-junie install-opencode

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

install-claude: ## Install shared agent skills into .claude/skills
	$(call install_skills,claude)

install-junie: ## Install shared agent skills into .junie/skills
	$(call install_skills,junie)

install-opencode: ## Install shared agent skills into .opencode/skills
	$(call install_skills,opencode)

install-agents: install-claude install-junie install-opencode ## Install shared agent skills for all agents

clean: ## Clean build artifacts
	rm -rf $(BUILD_DIR)
	rm -f coverage.out

help: ## Show this help message
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'
