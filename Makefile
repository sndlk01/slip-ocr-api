# Makefile for OCR API

.PHONY: help build run test clean docker-build docker-up docker-down docker-logs install dev

# Default target
help:
	@echo "OCR API - Available Commands"
	@echo "============================"
	@echo ""
	@echo "Docker Commands (Recommended):"
	@echo "  make docker-up      - Start the application with Docker"
	@echo "  make docker-down    - Stop the application"
	@echo "  make docker-logs    - View application logs"
	@echo "  make docker-build   - Rebuild Docker image"
	@echo "  make docker-restart - Restart the application"
	@echo "  make docker-clean   - Remove containers and volumes"
	@echo ""
	@echo "Local Development:"
	@echo "  make install        - Install dependencies"
	@echo "  make build          - Build the application"
	@echo "  make run            - Run the application"
	@echo "  make dev            - Run in development mode"
	@echo "  make test           - Run tests"
	@echo "  make clean          - Clean build artifacts"
	@echo ""

# Docker commands
docker-up:
	@echo "🚀 Starting OCR API with Docker..."
	docker-compose up -d
	@echo "✅ Application started at http://localhost:8080"
	@echo "   Health check: http://localhost:8080/health"

docker-down:
	@echo "🛑 Stopping OCR API..."
	docker-compose down

docker-logs:
	docker-compose logs -f

docker-build:
	@echo "🔨 Building Docker image..."
	docker-compose build --no-cache

docker-restart:
	@echo "🔄 Restarting OCR API..."
	docker-compose restart

docker-clean:
	@echo "🧹 Cleaning Docker resources..."
	docker-compose down -v
	docker system prune -f

docker-shell:
	@echo "🐚 Opening shell in container..."
	docker-compose exec ocr-api sh

# Local development commands
install:
	@echo "📦 Installing dependencies..."
	go mod download
	go mod tidy

build:
	@echo "🔨 Building application..."
	go build -o ocr-api.exe

run: build
	@echo "🚀 Running application..."
	./ocr-api.exe

dev:
	@echo "🔧 Running in development mode..."
	go run main.go

test:
	@echo "🧪 Running tests..."
	go test -v ./...

test-coverage:
	@echo "📊 Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

clean:
	@echo "🧹 Cleaning build artifacts..."
	rm -f ocr-api ocr-api.exe
	rm -f coverage.out coverage.html
	rm -rf uploads/*
	@echo "✅ Clean complete"

# Database commands
db-reset:
	@echo "⚠️  Resetting database..."
	rm -f db.sqlite
	@echo "✅ Database reset"

# Quick test commands
health-check:
	@echo "🏥 Checking API health..."
	curl -s http://localhost:8080/health | jq

test-upload:
	@echo "📤 Testing upload endpoint..."
	@echo "Note: Make sure you have a test_slip.jpg file"
	curl -X POST http://localhost:8080/api/v1/upload \
		-F "slip=@test_slip.jpg" \
		-F "type=income" | jq

test-list:
	@echo "📋 Listing all transactions..."
	curl -s http://localhost:8080/api/v1/transactions | jq

# Linting and formatting
fmt:
	@echo "🎨 Formatting code..."
	go fmt ./...

lint:
	@echo "🔍 Linting code..."
	golangci-lint run

# Documentation
docs:
	@echo "📚 Available documentation:"
	@echo "  README.md         - Main documentation"
	@echo "  DOCKER.md         - Docker guide"
	@echo "  TESTING.md        - API testing examples"
	@echo "  DOCUMENTATION.md  - Technical documentation"
