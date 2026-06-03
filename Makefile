# wonderlic-calc — monorepo orchestration
# Delegates into frontend/ (npm) and backend/ (go).

.DEFAULT_GOAL := help
.PHONY: help install build test test-frontend test-backend cover-backend lint fmt \
        run-frontend run-backend dev clean all

FRONTEND_DIR := frontend
BACKEND_DIR  := backend
COVERAGE_MIN := 95.0

# Use bash so `nvm` (a shell function) works inside recipes.
SHELL := /bin/bash

# Proper Go toolchain on PATH for every recipe, regardless of caller's shell.
export PATH := $(HOME)/.local/go/bin:$(HOME)/go/bin:/usr/local/go/bin:$(PATH)
export GOPATH ?= $(HOME)/go

# Load nvm and activate the Node version pinned in ./.nvmrc. Prefix any
# node/npm command with $(use_node) so it runs under the pinned version.
NVM_DIR ?= $(HOME)/.nvm
use_node := export NVM_DIR="$(NVM_DIR)"; if [ -s "$$NVM_DIR/nvm.sh" ]; then . "$$NVM_DIR/nvm.sh" --no-use; nvm use --silent >/dev/null; fi;

help: ## Show available targets
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
		| awk 'BEGIN{FS=":.*?## "}{printf "  \033[36m%-16s\033[0m %s\n", $$1, $$2}'

install: ## Install deps for both workspaces
	$(use_node) cd $(FRONTEND_DIR) && npm install
	cd $(BACKEND_DIR) && go mod download

build: ## Production build for both workspaces
	$(use_node) cd $(FRONTEND_DIR) && npm run build
	cd $(BACKEND_DIR) && go build -o bin/server ./cmd/server

test: test-frontend test-backend ## Run all tests

test-frontend: ## Run frontend tests
	$(use_node) cd $(FRONTEND_DIR) && npm run test -- --run

test-backend: ## Run backend tests + coverage gate (>= $(COVERAGE_MIN)% over internal/)
	cd $(BACKEND_DIR) && go test ./... -coverprofile=coverage.out -covermode=atomic
	@cd $(BACKEND_DIR) && go test ./internal/... -coverprofile=coverage.internal.out -covermode=atomic >/dev/null
	@cd $(BACKEND_DIR) && go tool cover -func=coverage.internal.out | awk -v min=$(COVERAGE_MIN) \
		'/^total:/ {gsub("%","",$$3); printf "backend internal coverage: %.1f%% (min %.1f%%)\n", $$3, min; \
		if ($$3+0 < min+0) {print "FAIL: coverage below threshold"; exit 1}}'

cover-backend: ## Open backend coverage HTML report
	cd $(BACKEND_DIR) && go test ./internal/... -coverprofile=coverage.out -covermode=atomic
	cd $(BACKEND_DIR) && go tool cover -html=coverage.out

lint: ## Lint both workspaces
	$(use_node) cd $(FRONTEND_DIR) && npm run lint
	cd $(BACKEND_DIR) && go vet ./...

fmt: ## Format both workspaces
	$(use_node) cd $(FRONTEND_DIR) && npm run format
	cd $(BACKEND_DIR) && gofmt -w .

run-frontend: ## Run frontend dev server
	$(use_node) cd $(FRONTEND_DIR) && npm run dev

run-backend: ## Run backend server
	cd $(BACKEND_DIR) && go run ./cmd/server

dev: ## Run frontend and backend together
	$(MAKE) -j2 run-frontend run-backend

clean: ## Remove build artifacts
	rm -rf $(FRONTEND_DIR)/dist $(FRONTEND_DIR)/node_modules $(BACKEND_DIR)/bin

all: install lint test build ## Install, lint, test, build
