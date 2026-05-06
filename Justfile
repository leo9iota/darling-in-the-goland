exe_ext := if os() == "windows" { ".exe" } else { "" }

BINARY_NAME := "darling"
CMD_PATH := "./cmd/game"
BUILD_DIR := "./build"
VERSION := `git describe --tags --always --dirty 2>/dev/null || echo "dev"`
COMMIT := `git rev-parse --short HEAD 2>/dev/null || echo "unknown"`
BUILD_TIME := `date -u +"%Y-%m-%dT%H:%M:%SZ" 2>/dev/null || echo "unknown"`
LDFLAGS := "-s -w -X main.Version=" + VERSION + " -X main.Commit=" + COMMIT + " -X main.BuildTime=" + BUILD_TIME

# List all available tasks
default:
    @just --list

# ─── Development ───────────────────────────────────────────

# Run the game
run:
    go run {{CMD_PATH}}

# Build for current OS
build:
    go build -ldflags "{{LDFLAGS}}" -o {{BUILD_DIR}}/{{BINARY_NAME}}{{exe_ext}} {{CMD_PATH}}

# Remove all build artifacts
clean:
    rm -rf {{BUILD_DIR}}/windows/*.exe
    rm -rf {{BUILD_DIR}}/linux/{{BINARY_NAME}}
    rm -rf {{BUILD_DIR}}/macos/{{BINARY_NAME}}
    rm -rf {{BUILD_DIR}}/{{BINARY_NAME}}*

# ─── Cross-Platform Builds ────────────────────────────────
# NOTE: Ebitengine uses CGO (GLFW) on Linux/macOS.
# Cross-compiling from Windows requires a cross-compiler (e.g. zig cc).
# For CI, build each target on its native runner.

# Build for Windows (amd64)
build-windows:
    GOOS=windows GOARCH=amd64 go build -ldflags "{{LDFLAGS}}" -o {{BUILD_DIR}}/windows/{{BINARY_NAME}}.exe {{CMD_PATH}}

# Build for Linux (amd64) — requires Linux or cross-compiler
build-linux:
    GOOS=linux GOARCH=amd64 go build -ldflags "{{LDFLAGS}}" -o {{BUILD_DIR}}/linux/{{BINARY_NAME}} {{CMD_PATH}}

# Build for macOS (amd64) — requires macOS or cross-compiler
build-macos:
    GOOS=darwin GOARCH=amd64 go build -ldflags "{{LDFLAGS}}" -o {{BUILD_DIR}}/macos/{{BINARY_NAME}} {{CMD_PATH}}

# Build for macOS (Apple Silicon) — requires macOS or cross-compiler
build-macos-arm:
    GOOS=darwin GOARCH=arm64 go build -ldflags "{{LDFLAGS}}" -o {{BUILD_DIR}}/macos/{{BINARY_NAME}}-arm64 {{CMD_PATH}}

# Build for all platforms (use on CI with matrix builds)
build-all: build-windows build-linux build-macos build-macos-arm

# ─── Release ──────────────────────────────────────────────

# Full release pipeline (check + build all + package)
release: check build-all package

# Package builds into zip archives
package:
    cd {{BUILD_DIR}} && tar -czf windows/{{BINARY_NAME}}-windows-amd64.tar.gz -C windows {{BINARY_NAME}}.exe || true
    cd {{BUILD_DIR}} && tar -czf linux/{{BINARY_NAME}}-linux-amd64.tar.gz -C linux {{BINARY_NAME}} || true
    cd {{BUILD_DIR}} && tar -czf macos/{{BINARY_NAME}}-macos-amd64.tar.gz -C macos {{BINARY_NAME}} || true
    cd {{BUILD_DIR}} && tar -czf macos/{{BINARY_NAME}}-macos-arm64.tar.gz -C macos {{BINARY_NAME}}-arm64 || true

# ─── Code Quality ─────────────────────────────────────────

# Run go vet
vet:
    go vet ./...

# Format all Go source files
fmt:
    gofumpt -w .

# Run golangci-lint
lint:
    golangci-lint run ./...

# Run golangci-lint with auto-fix
lint-fix:
    golangci-lint run --fix ./...

# ─── Testing ──────────────────────────────────────────────

# Run all tests
test:
    go test ./...

# Run all tests (verbose)
test-v:
    go test -v ./...

# Run physics engine tests
test-physics:
    go test -v ./internal/physics/...

# Run benchmarks
bench:
    go test -bench=. ./...

# ─── Dependencies ─────────────────────────────────────────

# Tidy go modules
tidy:
    go mod tidy

# Download dependencies
deps:
    go mod download

# ─── Full Pipeline ────────────────────────────────────────

# Run all checks (fmt, vet, test, build)
check: fmt vet test build
