# Development Commands for RuneScape 2 Build 317 Renderer

# Build targets
build:
	@echo "Building Go renderer..."
	go build -o bin/render3d ./cmd/render3d
	@echo "Build complete: bin/render3d"

# Run the renderer
run:
	@echo "Running 3D renderer..."
	go run ./cmd/render3d

# Run with Ebitengine
run-ebit:
	@echo "Running with Ebitengine..."
	go run ./cmd/render3d --ebit

# Build for WebAssembly
build-wasm:
	@echo "Building WebAssembly..."
	GOOS=js GOARCH=wasm go build -o bin/render3d.wasm ./cmd/render3d

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf bin/
	go clean

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Run linter
lint:
	@echo "Running linter..."
	golangci-lint run ./...

# Vet code
vet:
	@echo "Vetting code..."
	go vet ./...

# Check for issues
check: fmt vet lint

# Install dependencies
deps:
	@echo "Installing dependencies..."
	go mod download
	go mod tidy

# Update dependencies
update-deps:
	@echo "Updating dependencies..."
	go get -u ./...
	go mod tidy

# Generate mocks for testing
mocks:
	@echo "Generating mocks..."
	mockgen -source=pkg/renderer/interfaces.go -destination=pkg/renderer/mocks/

# Benchmark performance
bench:
	@echo "Running benchmarks..."
	go test -bench=. -benchmem ./...

# Profile CPU
profile-cpu:
	@echo "Running CPU profiling..."
	go test -cpuprofile=cpu.prof -bench=. ./...
	go tool pprof -text cpu.prof

# Profile memory
profile-mem:
	@echo "Running memory profiling..."
	go test -memprofile=mem.prof -bench=. ./...
	go tool pprof -text mem.prof

# Help
help:
	@echo "Available commands:"
	@echo "  make build          - Build the renderer"
	@echo "  make run            - Run the renderer"
	@echo "  make run-ebit       - Run with Ebitengine"
	@echo "  make build-wasm     - Build for WebAssembly"
	@echo "  make clean          - Clean build artifacts"
	@echo "  make test           - Run tests"
	@echo "  make test-coverage  - Run tests with coverage"
	@echo "  make fmt            - Format code"
	@echo "  make lint           - Run linter"
	@echo "  make vet            - Vet code"
	@echo "  make check          - Run all checks (fmt, vet, lint)"
	@echo "  make deps           - Install dependencies"
	@echo "  make update-deps    - Update dependencies"
	@echo "  make mocks          - Generate test mocks"
	@echo "  make bench          - Run benchmarks"
	@echo "  make profile-cpu    - Profile CPU usage"
	@echo "  make profile-mem    - Profile memory usage"
	@echo "  make help           - Show this help message"

.PHONY: build run run-ebit build-wasm clean test test-coverage fmt vet lint check deps update-deps mocks bench profile-cpu profile-mem help
