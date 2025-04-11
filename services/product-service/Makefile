# Makefile

# Variables
APP_NAME := product-service
BUILD_DIR := build
DOCKER_REGISTRY := your-registry
VERSION := $(shell git describe --tags --always --dirty)
GOOS ?= linux
GOARCH ?= amd64

.PHONY: all build clean test lint docker-build docker-push run-local swagger generate help

all: clean lint test build

# Build the application
build:
	@echo "Building $(APP_NAME)..."
	@mkdir -p $(BUILD_DIR)
	@CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) go build -a -installsuffix cgo -ldflags "-X main.version=$(VERSION)" -o $(BUILD_DIR)/$(APP_NAME) ./cmd/server/main.go
	@echo "Build completed: $(BUILD_DIR)/$(APP_NAME)"

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)
	@go clean
	@echo "Clean completed"

# Run tests
test:
	@echo "Running tests..."
	@go test -v ./... -cover
	@echo "Tests completed"

# Run linter
lint:
	@echo "Running linter..."
	@golangci-lint run ./...
	@echo "Lint completed"

# Build docker image
docker-build:
	@echo "Building Docker image..."
	@docker build -t $(DOCKER_REGISTRY)/$(APP_NAME):$(VERSION) -f deployments/docker/Dockerfile .
	@docker tag $(DOCKER_REGISTRY)/$(APP_NAME):$(VERSION) $(DOCKER_REGISTRY)/$(APP_NAME):latest
	@echo "Docker build completed"

# Push docker image
docker-push:
	@echo "Pushing Docker image..."
	@docker push $(DOCKER_REGISTRY)/$(APP_NAME):$(VERSION)
	@docker push $(DOCKER_REGISTRY)/$(APP_NAME):latest
	@echo "Docker push completed"

# Run locally with docker-compose
run-local:
	@echo "Running locally with docker-compose..."
	@docker-compose -f deployments/docker/docker-compose.yml up -d
	@echo "Application is running"

# Generate swagger docs
swagger:
	@echo "Generating Swagger documentation..."
	@swag init -g cmd/server/main.go
	@echo "Swagger generation completed"

# Generate code
generate:
	@echo "Running code generation..."
	@go generate ./...
	@echo "Code generation completed"

# Help
help:
	@echo "Makefile for $(APP_NAME)"
	@echo ""
	@echo "Usage:"
	@echo "  make all          - Run clean, lint, test, and build"
	@echo "  make build        - Build the application"
	@echo "  make clean        - Clean build artifacts"
	@echo "  make test         - Run tests"
	@echo "  make lint         - Run linter"
	@echo "  make docker-build - Build Docker image"
	@echo "  make docker-push  - Push Docker image to registry"
	@echo "  make run-local    - Run locally with docker-compose"
	@echo "  make swagger      - Generate Swagger documentation"
	@echo "  make generate     - Run code generation"
	@echo "  make help         - Show this help"
