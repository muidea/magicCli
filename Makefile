.PHONY: build install clean

APP_NAME=magicCli
SRC=application/magicCli/cmd/main.go

GIT_VERSION = $(shell git describe --tags --always --dirty 2>/dev/null || echo "dailyBuild")
GIT_COMMIT = $(shell git rev-parse HEAD 2>/dev/null || echo "unknown")
BUILD_DATE = $(shell date -u +"%Y-%m-%dT%H:%M:%SZ" 2>/dev/null || echo "unknown")

LDFLAGS = -X main.version=$(GIT_VERSION) -X main.gitCommit=$(GIT_COMMIT) -X main.buildDate=$(BUILD_DATE)

build:
	go build -ldflags "$(LDFLAGS)" -o bin/$(APP_NAME) $(SRC)

install: build
	cp bin/$(APP_NAME) ~/.local/bin/

clean:
	rm -rf bin/
