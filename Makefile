GOBASE=$(shell pwd)
GOBIN=$(GOBASE)/bin
GOPATH=$(shell go env GOPATH)

.PHONY: all build clean run test lint fmt deps

all: build

deps:
	@echo " > Installing dependencies..."
	go install github.com/wailsapp/wails/v2/cmd/wails@latest
	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.8.0
	go install mvdan.cc/gofumpt@latest
	go install golang.org/x/tools/cmd/goimports@latest
	go install github.com/evilmartians/lefthook@latest

	@echo " > Setting up git hooks..."
	lefthook install

	@echo " > Install Frontend dependencies..."
	cd frontend && pnpm install

dev:
	@echo " > Starting development server..."
	wails dev

build-window:
	@echo " > Building on window"
	wails build -platform windows/amd64 -nsis

format:
	@echo " > Formatting code..."
	gofumpt -w .
	goimports -w -local github.com/zenkiet/window-service-watcher .

lint:
	@echo " > Linting..."
	golangci-lint run ./...