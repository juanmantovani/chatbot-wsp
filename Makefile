# WhatsApp Chatbot Makefile

.PHONY: help build run test clean docker-build docker-run deploy

# Default target
help:
	@echo "Available targets:"
	@echo "  build        - Build the application"
	@echo "  run          - Run the application locally"
	@echo "  test         - Run tests"
	@echo "  clean        - Clean build artifacts"
	@echo "  docker-build - Build Docker image"
	@echo "  docker-run   - Run Docker container"
	@echo "  deploy       - Deploy to AWS"

# Build the application
build:
	@echo "Building application..."
	go build -o bin/chatbot-wsp ./cmd/main.go

# Run the application locally
run:
	@echo "Running application..."
	go run ./cmd/main.go

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf bin/
	go clean

# Build Docker image
docker-build:
	@echo "Building Docker image..."
	docker build -t chatbot-wsp:latest .

# Run Docker container
docker-run:
	@echo "Running Docker container..."
	docker run -p 8080:8080 --env-file .env chatbot-wsp:latest

# Deploy to AWS
deploy:
	@echo "Deploying to AWS..."
	./scripts/deploy.sh

# Install dependencies
deps:
	@echo "Installing dependencies..."
	go mod download
	go mod tidy

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Lint code
lint:
	@echo "Linting code..."
	golangci-lint run

# Generate mocks (if using mockgen)
mocks:
	@echo "Generating mocks..."
	mockgen -source=internal/domain/repository/chatbot_repository.go -destination=internal/mocks/chatbot_repository_mock.go
	mockgen -source=internal/domain/service/chatbot_service.go -destination=internal/mocks/chatbot_service_mock.go
