# wonderlic-calc — monorepo orchestration
# Delegates into frontend/ (npm) and backend/ (go).

.DEFAULT_GOAL := help
.PHONY: help install build test test-frontend test-backend lint fmt \
        run-frontend run-backend dev clean all

FRONTEND_DIR := frontend
BACKEND_DIR  := backend

help: ## Show available targets
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
		| awk 'BEGIN{FS=":.*?## "}{printf "  \033[36m%-16s\033[0m %s\n", $$1, $$2}'

install: ## Install deps for both workspaces
	cd $(FRONTEND_DIR) && npm install
	cd $(BACKEND_DIR) && go mod download

build: ## Production build for both workspaces
	cd $(FRONTEND_DIR) && npm run build
	cd $(BACKEND_DIR) && go build -o bin/server ./cmd/server

test: test-frontend test-backend ## Run all tests

test-frontend: ## Run frontend tests
	cd $(FRONTEND_DIR) && npm run test -- --run

test-backend: ## Run backend tests
	cd $(BACKEND_DIR) && go test ./...

lint: ## Lint both workspaces
	cd $(FRONTEND_DIR) && npm run lint
	cd $(BACKEND_DIR) && go vet ./...

fmt: ## Format both workspaces
	cd $(FRONTEND_DIR) && npm run format
	cd $(BACKEND_DIR) && gofmt -w .

run-frontend: ## Run frontend dev server
	cd $(FRONTEND_DIR) && npm run dev

run-backend: ## Run backend server
	cd $(BACKEND_DIR) && go run ./cmd/server

dev: ## Run frontend and backend together
	$(MAKE) -j2 run-frontend run-backend

clean: ## Remove build artifacts
	rm -rf $(FRONTEND_DIR)/dist $(FRONTEND_DIR)/node_modules $(BACKEND_DIR)/bin

all: install lint test build ## Install, lint, test, build
