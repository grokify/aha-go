.PHONY: all build build-cli test test-coverage lint generate clean tidy fmt vet check install completions help

# Binary name
BINARY_NAME=aha
BUILD_DIR=bin

all: build

# Build all packages
build:
	go build ./...

# Build the CLI binary
build-cli:
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/aha

# Run tests
test:
	go test -v -race ./...

# Run tests with coverage
test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Run linter
lint:
	golangci-lint run ./...

# Generate API client from OpenAPI spec
generate:
	./generate.sh

# Clean build artifacts
clean:
	rm -rf internal/api/*.go
	rm -rf $(BUILD_DIR)
	rm -f coverage.out coverage.html

# Tidy dependencies
tidy:
	go mod tidy

# Format code
fmt:
	gofmt -s -w .

# Run go vet
vet:
	go vet ./...

# Run all checks (fmt, vet, lint, test)
check: fmt vet lint test

# Install CLI to GOPATH/bin
install:
	go install ./cmd/aha

# Generate shell completions (requires build-cli first)
completions: build-cli
	@mkdir -p $(BUILD_DIR)/completions
	$(BUILD_DIR)/$(BINARY_NAME) completion bash > $(BUILD_DIR)/completions/aha.bash
	$(BUILD_DIR)/$(BINARY_NAME) completion zsh > $(BUILD_DIR)/completions/_aha
	$(BUILD_DIR)/$(BINARY_NAME) completion fish > $(BUILD_DIR)/completions/aha.fish

# Show help
help:
	@echo "Available targets:"
	@echo "  build          - Build all packages"
	@echo "  build-cli      - Build the CLI binary to bin/"
	@echo "  test           - Run tests with race detector"
	@echo "  test-coverage  - Run tests with coverage report"
	@echo "  lint           - Run golangci-lint"
	@echo "  generate       - Generate API client from OpenAPI spec"
	@echo "  clean          - Clean build artifacts"
	@echo "  tidy           - Tidy go.mod dependencies"
	@echo "  fmt            - Format code with gofmt"
	@echo "  vet            - Run go vet"
	@echo "  check          - Run all checks (fmt, vet, lint, test)"
	@echo "  install        - Install CLI to GOPATH/bin"
	@echo "  completions    - Generate shell completions"
	@echo "  help           - Show this help"

.DEFAULT_GOAL := build
