# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
BINARY_NAME=request-management-system
BINARY_UNIX=$(BINARY_NAME)_unix

# Docker parameters
DOCKER_IMAGE=request-management-system
DOCKER_TAG=latest

.PHONY: all build clean test deps run docker-build docker-run docker-stop help

# Default target
all: test build

# Build the application
build:
	$(GOBUILD) -o $(BINARY_NAME) -v ./cmd/server

# Build for Linux
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v ./cmd/server

# Clean build artifacts
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)

# Run tests
test:
	$(GOTEST) -v ./...

# Run tests with coverage
test-coverage:
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out

# Download dependencies
deps:
	$(GOMOD) download
	$(GOMOD) tidy

# Run the application locally
run:
	$(GOBUILD) -o $(BINARY_NAME) -v ./cmd/server
	./$(BINARY_NAME)

# Run with live reload (requires air: go install github.com/cosmtrek/air@latest)
dev:
	air

# Database operations
db-setup:
	docker-compose up -d postgres
	sleep 10
	docker-compose exec postgres psql -U postgres -d request_management_system -f /docker-entrypoint-initdb.d/01-schema.sql

db-reset:
	docker-compose down -v
	docker-compose up -d postgres
	sleep 10

# Docker operations
docker-build:
	docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) .

docker-run:
	docker-compose up -d

docker-stop:
	docker-compose down

docker-logs:
	docker-compose logs -f app

docker-clean:
	docker-compose down -v
	docker rmi $(DOCKER_IMAGE):$(DOCKER_TAG) || true

# Development environment
dev-up:
	docker-compose up -d postgres
	sleep 10
	$(MAKE) run

dev-down:
	docker-compose down

# Production deployment
deploy:
	docker-compose -f docker-compose.yml -f docker-compose.prod.yml up -d

# Linting and formatting
lint:
	golangci-lint run

fmt:
	$(GOCMD) fmt ./...

# Security scan
security:
	gosec ./...

# Generate API documentation (requires swag: go install github.com/swaggo/swag/cmd/swag@latest)
docs:
	swag init -g cmd/server/main.go

# Install development tools
install-tools:
	$(GOGET) github.com/cosmtrek/air@latest
	$(GOGET) github.com/swaggo/swag/cmd/swag@latest
	$(GOGET) github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	$(GOGET) github.com/securecodewarrior/gosec/v2/cmd/gosec@latest

# Help
help:
	@echo "Available commands:"
	@echo "  build         - Build the application"
	@echo "  build-linux   - Build for Linux"
	@echo "  clean         - Clean build artifacts"
	@echo "  test          - Run tests"
	@echo "  test-coverage - Run tests with coverage"
	@echo "  deps          - Download dependencies"
	@echo "  run           - Run the application locally"
	@echo "  dev           - Run with live reload"
	@echo "  db-setup      - Set up database"
	@echo "  db-reset      - Reset database"
	@echo "  docker-build  - Build Docker image"
	@echo "  docker-run    - Run with Docker Compose"
	@echo "  docker-stop   - Stop Docker containers"
	@echo "  docker-logs   - View application logs"
	@echo "  docker-clean  - Clean Docker resources"
	@echo "  dev-up        - Start development environment"
	@echo "  dev-down      - Stop development environment"
	@echo "  deploy        - Deploy to production"
	@echo "  lint          - Run linter"
	@echo "  fmt           - Format code"
	@echo "  security      - Run security scan"
	@echo "  docs          - Generate API documentation"
	@echo "  install-tools - Install development tools"
	@echo "  help          - Show this help message"

