.PHONY: help build run test clean install dev

# Default target
help:
	@echo "Available targets:"
	@echo "  make build     - Build the application"
	@echo "  make run       - Run the application"
	@echo "  make dev       - Run in development mode"
	@echo "  make test      - Run tests"
	@echo "  make clean     - Clean build artifacts"
	@echo "  make install   - Install dependencies"
	@echo "  make fmt       - Format code"
	@echo "  make lint      - Run linter"

# Install dependencies
install:
	go mod download
	go mod verify

# Build the application
build:
	go build -o spotify-sorter cmd/server/main.go

# Run the application
run:
	go run cmd/server/main.go

# Run in development mode (requires air)
dev:
	@which air > /dev/null || (echo "Installing air..." && go install github.com/air-verse/air@latest)
	air

# Run tests
test:
	go test -v ./...

# Run tests with coverage
test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Format code
fmt:
	go fmt ./...

# Run linter (requires golangci-lint)
lint:
	@which golangci-lint > /dev/null || (echo "Installing golangci-lint..." && go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest)
	golangci-lint run

# Clean build artifacts
clean:
	rm -f spotify-sorter
	rm -f coverage.out coverage.html
	go clean

# Generate .env file if it doesn't exist
setup-env:
	@if [ ! -f .env ]; then \
		cp .env.example .env; \
		echo ".env file created. Please edit it with your credentials."; \
	else \
		echo ".env file already exists."; \
	fi

# Full setup
setup: install setup-env
	@echo "Setup complete! Edit .env with your Spotify credentials, then run 'make run'"

# Docker targets (optional)
docker-build:
	docker build -t spotify-playlist-sorter .

docker-run:
	docker run -p 8080:8080 --env-file .env spotify-playlist-sorter
