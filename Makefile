# Variables
BINARY_NAME=ghm
CMD_PATH=./cmd/ghm
MAIN_PACKAGE=github.com/Cassin01/ghm
VERSION?=$(shell git describe --tags --always --dirty)
COMMIT?=$(shell git rev-parse --short HEAD)
BUILD_DATE?=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS=-ldflags "-s -w -X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.buildDate=$(BUILD_DATE)"

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOLINT=golangci-lint

# Build targets
.PHONY: all build clean test coverage lint fmt vet install uninstall help
.PHONY: build-linux build-darwin build-all
.PHONY: release-linux release-darwin release-all

# Default target
all: clean lint test build

# Build the binary
build:
	$(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME) $(CMD_PATH)

# Clean build artifacts
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_NAME)-*
	rm -f *.tar.gz
	rm -f *.zip
	rm -f checksums.txt
	rm -f coverage.out
	rm -f coverage.html

# Run tests
test:
	$(GOTEST) -v -race ./...

# Run tests with coverage
coverage:
	$(GOTEST) -v -race -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html

# Run linter
lint:
	$(GOLINT) run --timeout=5m

# Format code
fmt:
	$(GOCMD) fmt ./...

# Run vet
vet:
	$(GOCMD) vet ./...

# Install dependencies
deps:
	$(GOMOD) download
	$(GOMOD) tidy

# Install binary to GOPATH/bin
install:
	$(GOCMD) install $(LDFLAGS) $(CMD_PATH)

# Uninstall binary from GOPATH/bin
uninstall:
	rm -f $(GOPATH)/bin/$(BINARY_NAME)

# Cross-compilation builds
build-linux:
	GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME)-linux-amd64 $(CMD_PATH)
	GOOS=linux GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME)-linux-arm64 $(CMD_PATH)

build-darwin:
	GOOS=darwin GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME)-darwin-amd64 $(CMD_PATH)
	GOOS=darwin GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME)-darwin-arm64 $(CMD_PATH)

build-all: build-linux build-darwin

# Release targets with archives
release-linux: build-linux
	tar -czf $(BINARY_NAME)-linux-amd64.tar.gz $(BINARY_NAME)-linux-amd64
	tar -czf $(BINARY_NAME)-linux-arm64.tar.gz $(BINARY_NAME)-linux-arm64

release-darwin: build-darwin
	tar -czf $(BINARY_NAME)-darwin-amd64.tar.gz $(BINARY_NAME)-darwin-amd64
	tar -czf $(BINARY_NAME)-darwin-arm64.tar.gz $(BINARY_NAME)-darwin-arm64

release-all: release-linux release-darwin
	sha256sum $(BINARY_NAME)-*.tar.gz > checksums.txt

# Development helpers
dev: clean fmt vet lint test build

# Run the binary
run: build
	./$(BINARY_NAME)

# Show help
help:
	@echo "Available targets:"
	@echo "  build         - Build the binary"
	@echo "  clean         - Clean build artifacts"
	@echo "  test          - Run tests"
	@echo "  coverage      - Run tests with coverage"
	@echo "  lint          - Run linter"
	@echo "  fmt           - Format code"
	@echo "  vet           - Run vet"
	@echo "  deps          - Install dependencies"
	@echo "  install       - Install binary to GOPATH/bin"
	@echo "  uninstall     - Remove binary from GOPATH/bin"
	@echo "  build-all     - Build for all platforms"
	@echo "  release-all   - Build release archives for all platforms"
	@echo "  dev           - Run development workflow (clean, fmt, vet, lint, test, build)"
	@echo "  run           - Build and run the binary"
	@echo "  help          - Show this help message"