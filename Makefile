.PHONY: all build test lint clean cross

# Definições Padrão
APP_NAME := crom-agente-cli
BIN_DIR := bin
GO_FILES := $(shell find . -name '*.go' -not -path "./vendor/*")
VERSION := dev

all: lint test build

build:
	@echo "==> Building $(APP_NAME)..."
	@mkdir -p $(BIN_DIR)
	go build -ldflags="-s -w -X main.Version=$(VERSION)" -v -o $(BIN_DIR)/$(APP_NAME) ./cmd/crom-agente-cli

# Compila para os principais sistemas operacionais e arquiteturas (Item 41)
cross:
	@echo "==> Cross-compiling $(APP_NAME)..."
	@mkdir -p $(BIN_DIR)
	@echo "Compilando para Linux (amd64)..."
	@GOOS=linux GOARCH=amd64 go build -ldflags="-s -w -X main.Version=$(VERSION)" -o $(BIN_DIR)/$(APP_NAME)-linux-amd64 ./cmd/crom-agente-cli
	@echo "Compilando para Linux (arm64)..."
	@GOOS=linux GOARCH=arm64 go build -ldflags="-s -w -X main.Version=$(VERSION)" -o $(BIN_DIR)/$(APP_NAME)-linux-arm64 ./cmd/crom-agente-cli
	@echo "Compilando para macOS (amd64)..."
	@GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w -X main.Version=$(VERSION)" -o $(BIN_DIR)/$(APP_NAME)-darwin-amd64 ./cmd/crom-agente-cli
	@echo "Compilando para macOS (arm64)..."
	@GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w -X main.Version=$(VERSION)" -o $(BIN_DIR)/$(APP_NAME)-darwin-arm64 ./cmd/crom-agente-cli
	@echo "Compilação cross-platform concluída com sucesso em $(BIN_DIR)/"

test:
	@echo "==> Running tests..."
	go test -v -race -cover ./...

lint:
	@echo "==> Running linter..."
	@if command -v golangci-lint >/dev/null; then \
		golangci-lint run; \
	else \
		echo "golangci-lint is not installed."; \
		exit 1; \
	fi

clean:
	@echo "==> Cleaning..."
	go clean
	rm -rf $(BIN_DIR)
