# Phonic AI Calling Agent - Makefile
# Production-grade build automation for Go microservices

# Variables
PROJECT_NAME := phonic
GO_VERSION := 1.24
DOCKER_REGISTRY := ghcr.io/arbajansari19
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME := $(shell date -u '+%Y-%m-%d_%H:%M:%S')
GIT_COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")

# Go build variables
LDFLAGS := -X main.Version=$(VERSION) \
           -X main.BuildTime=$(BUILD_TIME) \
           -X main.GitCommit=$(GIT_COMMIT)

# Directories
BIN_DIR := ./bin
PROTO_DIR := ./proto
SERVICES_DIR := ./services
BUILD_DIR := ./build

# Services
SERVICES := gateway stt-client tts-client orchestrator session

# Colors for output
RED := \033[31m
GREEN := \033[32m
YELLOW := \033[33m
BLUE := \033[34m
RESET := \033[0m

.PHONY: help
help: ## Show this help message
	@echo "$(BLUE)🎵 Phonic AI Calling Agent - Available Commands$(RESET)"
	@echo ""
	@echo "$(YELLOW)📦 Build Commands:$(RESET)"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | grep -E "(build|proto|install)" | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(GREEN)%-20s$(RESET) %s\n", $$1, $$2}'
	@echo ""
	@echo "$(YELLOW)🚀 Development Commands:$(RESET)"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | grep -E "(dev|run|up|down|logs)" | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(GREEN)%-20s$(RESET) %s\n", $$1, $$2}'
	@echo ""
	@echo "$(YELLOW)🧪 Testing & Quality:$(RESET)"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | grep -E "(test|lint|fmt|check)" | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(GREEN)%-20s$(RESET) %s\n", $$1, $$2}'
	@echo ""
	@echo "$(YELLOW)🐳 Docker Commands:$(RESET)"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | grep -E "(docker|compose)" | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(GREEN)%-20s$(RESET) %s\n", $$1, $$2}'
	@echo ""
	@echo "$(YELLOW)🛠️ Utility Commands:$(RESET)"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | grep -v -E "(build|proto|install|dev|run|up|down|logs|test|lint|fmt|check|docker|compose)" | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(GREEN)%-20s$(RESET) %s\n", $$1, $$2}'

# =============================================================================
# Setup and Installation
# =============================================================================

.PHONY: setup
setup: ## Install all dependencies and setup development environment
	@echo "$(BLUE)🔧 Setting up Phonic development environment...$(RESET)"
	@command -v go >/dev/null 2>&1 || { echo "$(RED)❌ Go is not installed$(RESET)"; exit 1; }
	@command -v docker >/dev/null 2>&1 || { echo "$(RED)❌ Docker is not installed$(RESET)"; exit 1; }
	@command -v protoc >/dev/null 2>&1 || { echo "$(RED)❌ protoc is not installed$(RESET)"; exit 1; }
	@go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	@go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	@go mod download
	@mkdir -p $(BIN_DIR) $(BUILD_DIR)
	@echo "$(GREEN)✅ Development environment ready!$(RESET)"

.PHONY: deps
deps: ## Download and tidy Go dependencies
	@echo "$(BLUE)📦 Managing Go dependencies...$(RESET)"
	@go mod download
	@go mod tidy
	@echo "$(GREEN)✅ Dependencies updated$(RESET)"

# =============================================================================
# Code Generation
# =============================================================================

.PHONY: proto
proto: ## Generate Go code from protobuf definitions
	@echo "$(BLUE)🔄 Generating Go code from protobuf files...$(RESET)"
	@mkdir -p $(PROTO_DIR)/gen
	@find $(PROTO_DIR) -name "*.proto" -exec protoc \
		--go_out=. \
		--go_opt=paths=source_relative \
		--go-grpc_out=. \
		--go-grpc_opt=paths=source_relative \
		{} \;
	@echo "$(GREEN)✅ Protobuf code generated$(RESET)"

.PHONY: proto-clean
proto-clean: ## Clean generated protobuf files
	@echo "$(BLUE)🧹 Cleaning generated protobuf files...$(RESET)"
	@find . -name "*.pb.go" -delete
	@find . -name "*_grpc.pb.go" -delete
	@echo "$(GREEN)✅ Protobuf files cleaned$(RESET)"

# =============================================================================
# Build Commands
# =============================================================================

.PHONY: build
build: proto ## Build all services
	@echo "$(BLUE)🏗️ Building all services...$(RESET)"
	@mkdir -p $(BIN_DIR)
	@for service in $(SERVICES); do \
		echo "$(YELLOW)Building $$service...$(RESET)"; \
		go build -ldflags "$(LDFLAGS)" -o $(BIN_DIR)/$$service ./services/$$service/; \
	done
	@echo "$(GREEN)✅ All services built successfully$(RESET)"

.PHONY: build-gateway
build-gateway: proto ## Build gateway service
	@echo "$(BLUE)🏗️ Building gateway service...$(RESET)"
	@mkdir -p $(BIN_DIR)
	@go build -ldflags "$(LDFLAGS)" -o $(BIN_DIR)/gateway ./services/gateway/
	@echo "$(GREEN)✅ Gateway service built$(RESET)"

.PHONY: build-orchestrator
build-orchestrator: proto ## Build orchestrator service
	@echo "$(BLUE)🏗️ Building orchestrator service...$(RESET)"
	@mkdir -p $(BIN_DIR)
	@go build -ldflags "$(LDFLAGS)" -o $(BIN_DIR)/orchestrator ./services/orchestrator/
	@echo "$(GREEN)✅ Orchestrator service built$(RESET)"

.PHONY: install
install: build ## Install binaries to $GOPATH/bin
	@echo "$(BLUE)📦 Installing binaries...$(RESET)"
	@for service in $(SERVICES); do \
		cp $(BIN_DIR)/$$service $(GOPATH)/bin/; \
	done
	@echo "$(GREEN)✅ Binaries installed to $(GOPATH)/bin$(RESET)"

# =============================================================================
# Development
# =============================================================================

.PHONY: dev
dev: ## Start development environment with hot reload
	@echo "$(BLUE)🚀 Starting development environment...$(RESET)"
	@echo "$(YELLOW)Use Ctrl+C to stop all services$(RESET)"
	@docker-compose -f docker-compose.dev.yml up --build

.PHONY: dev-status
dev-status: ## Show development environment status
	@./scripts/dev-status.sh

.PHONY: run-gateway
run-gateway: build-gateway ## Run gateway service locally
	@echo "$(BLUE)🚀 Starting gateway service...$(RESET)"
	@$(BIN_DIR)/gateway

.PHONY: run-orchestrator
run-orchestrator: build-orchestrator ## Run orchestrator service locally
	@echo "$(BLUE)🚀 Starting orchestrator service...$(RESET)"
	@$(BIN_DIR)/orchestrator

# =============================================================================
# Docker Commands
# =============================================================================

.PHONY: docker-build
docker-build: ## Build all Docker images
	@echo "$(BLUE)🐳 Building Docker images...$(RESET)"
	@for service in $(SERVICES); do \
		echo "$(YELLOW)Building $$service Docker image...$(RESET)"; \
		docker build -f infra/docker/$$service.Dockerfile -t $(DOCKER_REGISTRY)/$(PROJECT_NAME)-$$service:$(VERSION) .; \
	done
	@echo "$(GREEN)✅ All Docker images built$(RESET)"

.PHONY: compose-up
compose-up: ## Start all services with docker-compose
	@echo "$(BLUE)🚀 Starting all services with Docker Compose...$(RESET)"
	@docker-compose up -d
	@echo "$(GREEN)✅ All services started$(RESET)"
	@echo "$(YELLOW)Run 'make logs' to see service logs$(RESET)"

.PHONY: compose-down
compose-down: ## Stop all services and remove containers
	@echo "$(BLUE)🛑 Stopping all services...$(RESET)"
	@docker-compose down
	@echo "$(GREEN)✅ All services stopped$(RESET)"

.PHONY: compose-logs
compose-logs: ## Show logs from all services
	@docker-compose logs -f

.PHONY: compose-restart
compose-restart: compose-down compose-up ## Restart all services

# =============================================================================
# Testing
# =============================================================================

.PHONY: test
test: ## Run all tests
	@echo "$(BLUE)🧪 Running all tests...$(RESET)"
	@go test -v ./...
	@echo "$(GREEN)✅ All tests passed$(RESET)"

.PHONY: test-coverage
test-coverage: ## Run tests with coverage report
	@echo "$(BLUE)🧪 Running tests with coverage...$(RESET)"
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "$(GREEN)✅ Coverage report generated: coverage.html$(RESET)"

.PHONY: test-race
test-race: ## Run tests with race detection
	@echo "$(BLUE)🧪 Running tests with race detection...$(RESET)"
	@go test -race -v ./...
	@echo "$(GREEN)✅ Race condition tests passed$(RESET)"

.PHONY: benchmark
benchmark: ## Run benchmarks
	@echo "$(BLUE)⚡ Running benchmarks...$(RESET)"
	@go test -bench=. -benchmem ./...

# =============================================================================
# Code Quality
# =============================================================================

.PHONY: fmt
fmt: ## Format Go code
	@echo "$(BLUE)✨ Formatting Go code...$(RESET)"
	@go fmt ./...
	@echo "$(GREEN)✅ Code formatted$(RESET)"

.PHONY: lint
lint: ## Run linters
	@echo "$(BLUE)🔍 Running linters...$(RESET)"
	@command -v golangci-lint >/dev/null 2>&1 || { echo "$(YELLOW)Installing golangci-lint...$(RESET)"; go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; }
	@golangci-lint run
	@echo "$(GREEN)✅ Linting passed$(RESET)"

.PHONY: vet
vet: ## Run go vet
	@echo "$(BLUE)🔍 Running go vet...$(RESET)"
	@go vet ./...
	@echo "$(GREEN)✅ Vet checks passed$(RESET)"

.PHONY: check
check: fmt vet lint test ## Run all quality checks

# =============================================================================
# Cleanup
# =============================================================================

.PHONY: clean
clean: ## Clean build artifacts
	@echo "$(BLUE)🧹 Cleaning build artifacts...$(RESET)"
	@rm -rf $(BIN_DIR) $(BUILD_DIR)
	@rm -f coverage.out coverage.html
	@go clean ./...
	@echo "$(GREEN)✅ Build artifacts cleaned$(RESET)"

.PHONY: clean-docker
clean-docker: ## Clean Docker images and containers
	@echo "$(BLUE)🧹 Cleaning Docker resources...$(RESET)"
	@docker-compose down --volumes --remove-orphans
	@docker system prune -f
	@echo "$(GREEN)✅ Docker resources cleaned$(RESET)"

.PHONY: clean-all
clean-all: clean clean-docker proto-clean ## Clean everything

# =============================================================================
# Utility Commands
# =============================================================================

.PHONY: version
version: ## Show version information
	@echo "$(BLUE)📋 Version Information:$(RESET)"
	@echo "Version: $(VERSION)"
	@echo "Build Time: $(BUILD_TIME)"
	@echo "Git Commit: $(GIT_COMMIT)"
	@echo "Go Version: $(shell go version)"

.PHONY: deps-update
deps-update: ## Update all dependencies to latest versions
	@echo "$(BLUE)📦 Updating dependencies...$(RESET)"
	@go get -u ./...
	@go mod tidy
	@echo "$(GREEN)✅ Dependencies updated$(RESET)"

.PHONY: security-check
security-check: ## Run security vulnerability checks
	@echo "$(BLUE)🔒 Running security checks...$(RESET)"
	@command -v gosec >/dev/null 2>&1 || { echo "$(YELLOW)Installing gosec...$(RESET)"; go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest; }
	@gosec ./...
	@echo "$(GREEN)✅ Security checks passed$(RESET)"

# Default target
.DEFAULT_GOAL := help
