.PHONY: all build lint fmt fmt-check vet install clean deps

APP_NAME=magicCli
SRC=application/magicCli/cmd/main.go

GIT_VERSION = $(shell git describe --tags --always --dirty 2>/dev/null || echo "dailyBuild")
GIT_COMMIT = $(shell git rev-parse HEAD 2>/dev/null || echo "unknown")
BUILD_DATE = $(shell date -u +"%Y-%m-%dT%H:%M:%SZ" 2>/dev/null || echo "unknown")

LDFLAGS = -X main.version=$(GIT_VERSION) -X main.gitCommit=$(GIT_COMMIT) -X main.buildDate=$(BUILD_DATE)

all: fmt vet build

fmt:
	@echo "Formatting code..."
	@find . -name "*.go" -not -path "./vendor/*" -exec gofmt -w {} +
	@if command -v goimports >/dev/null 2>&1; then \
		find . -name "*.go" -not -path "./vendor/*" -exec goimports -w {} +; \
		echo "goimports completed."; \
	else \
		echo "goimports not found, skipping..."; \
	fi

vet:
	@echo "Running go vet..."
	go vet ./...

lint:
	@echo "Running code quality checks..."
	@echo "1. Running go vet..."
	@go vet ./... 2>&1 | (grep -v "^vendor/" || true)
	@echo "2. Checking code format (excluding vendor)..."
	@if find . -name "*.go" -not -path "./vendor/*" -exec gofmt -d {} + 2>/dev/null | grep -q '^'; then \
		echo "Code is not formatted correctly. Run 'make fmt' to fix."; \
		find . -name "*.go" -not -path "./vendor/*" -exec gofmt -d {} + 2>/dev/null | head -50; \
		exit 1; \
	else \
		echo "Code is properly formatted."; \
	fi

fmt-check:
	@echo "Checking code format (excluding vendor)..."
	@if find . -name "*.go" -not -path "./vendor/*" -exec gofmt -d {} + 2>/dev/null | grep -q '^'; then \
		echo "Code is not formatted correctly. Run 'make fmt' to fix."; \
		find . -name "*.go" -not -path "./vendor/*" -exec gofmt -d {} + 2>/dev/null | head -50; \
		exit 1; \
	else \
		echo "Code is properly formatted."; \
	fi

build:
	go build -ldflags "$(LDFLAGS)" -o bin/$(APP_NAME) $(SRC)

install: build
	cp bin/$(APP_NAME) ~/.local/bin/

clean:
	rm -rf bin/

deps:
	@echo "Installing dependencies..."
	go mod tidy
	go mod download
