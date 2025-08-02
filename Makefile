# Gateway Worker Service Makefile

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
BINARY_NAME=worker
BINARY_PATH=bin/$(BINARY_NAME)
MAIN_PATH=cmd/worker/main.go

# Build parameters
LDFLAGS=-ldflags "-s -w"
BUILD_FLAGS=-v $(LDFLAGS)

# Docker parameters
DOCKER_IMAGE=gateway-worker
DOCKER_TAG=latest

.PHONY: all build clean test coverage deps lint proto help

all: deps proto build test

## Development
build: ## Build the application
	@echo "Building $(BINARY_NAME)..."
	mkdir -p bin
	$(GOBUILD) $(BUILD_FLAGS) -o $(BINARY_PATH) $(MAIN_PATH)

clean: ## Clean build artifacts
	@echo "Cleaning..."
	$(GOCLEAN)
	rm -rf bin/
	rm -rf pkg/proto/*.pb.go

test: ## Run tests
	@echo "Running tests..."
	$(GOTEST) -v ./...

test-coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html

deps: ## Install dependencies
	@echo "Installing dependencies..."
	$(GOMOD) download
	$(GOMOD) tidy

## Code Quality
lint: ## Run linters
	@echo "Running linters..."
	@which golangci-lint > /dev/null || (echo "Installing golangci-lint..." && go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest)
	golangci-lint run

fmt: ## Format code
	@echo "Formatting code..."
	$(GOCMD) fmt ./...

vet: ## Run go vet
	@echo "Running go vet..."
	$(GOCMD) vet ./...

## Protocol Buffers
proto: ## Generate protobuf code
	@echo "Generating protobuf code..."
	@which protoc > /dev/null || (echo "protoc is required. Please install Protocol Buffers compiler." && exit 1)
	@which protoc-gen-go > /dev/null || (echo "Installing protoc-gen-go..." && go install google.golang.org/protobuf/cmd/protoc-gen-go@latest)
	@which protoc-gen-go-grpc > /dev/null || (echo "Installing protoc-gen-go-grpc..." && go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest)
	mkdir -p pkg/proto
	protoc --go_out=pkg/proto --go_opt=paths=source_relative \
		--go-grpc_out=pkg/proto --go-grpc_opt=paths=source_relative \
		api/*.proto

## Docker
docker-build: ## Build Docker image
	@echo "Building Docker image..."
	docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) .

docker-run: ## Run Docker container
	@echo "Running Docker container..."
	docker run -p 8080:8080 -p 8081:8081 $(DOCKER_IMAGE):$(DOCKER_TAG)

## Development Server
run: build ## Build and run the application
	@echo "Running $(BINARY_NAME)..."
	./$(BINARY_PATH)

dev: ## Run in development mode with live reload (requires air)
	@which air > /dev/null || (echo "Installing air..." && go install github.com/air-verse/air@latest)
	air

## Utilities
check: fmt vet lint test ## Run all checks

install-tools: ## Install development tools
	@echo "Installing development tools..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/air-verse/air@latest

version: ## Show version information
	@echo "Go version: $(shell $(GOCMD) version)"
	@echo "Worker version: 1.0.0"

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Default target
.DEFAULT_GOAL := help