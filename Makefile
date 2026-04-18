.PHONY: help compile fmt fmt-check lint openapi clean install-dev dev backend-build backend-run backend-db-up backend-db-down backend-docker-build \
        docker-up docker-down docker-logs docker-build docker-rebuild frontend-test fronttest frontend-coverage frontend-lint \
        backend-test backend-test-coverage backend-test-verbose backend-lint backend-lint-strict backend-lint-fix backend-fmt backtest

TYPESPEC_DIR := typespec
BACKEND_DIR := backend
FRONTEND_DIR := frontend
GOLANGCI_LINT_VERSION := v2.11.4

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

compile: ## Compile TypeSpec specification
	cd $(TYPESPEC_DIR) && npx tsp compile .

fmt-check: ## Check TypeSpec formatting
	cd $(TYPESPEC_DIR) && npx prettier --check "**/*.tsp"

fmt-fix: ## Fix TypeSpec formatting
	cd $(TYPESPEC_DIR) && npx prettier --write "**/*.tsp"

lint: fmt-check compile ## Run TypeSpec linter (format check + compile)

openapi: ## Generate OpenAPI 3.0 specification
	cd $(TYPESPEC_DIR) && npx tsp compile . --emit @typespec/openapi3

clean: ## Remove generated files
	rm -rf $(TYPESPEC_DIR)/tsp-output

install-dev: ## Install frontend dependencies
	cd $(FRONTEND_DIR) && npm install

install-e2e: ## Install Playwright browsers for E2E tests
	cd $(FRONTEND_DIR) && npx playwright install --with-deps

dev: ## Start frontend dev server
	cd $(FRONTEND_DIR) && npm run dev

alltest: fronttest backtest

frontend-test: ## Run frontend tests
	cd $(FRONTEND_DIR) && npm test -- --run

fronttest:
	cd $(FRONTEND_DIR) && npm test -- --run
	cd $(FRONTEND_DIR) && npm run lint

frontend-coverage: ## Run frontend tests with coverage report
	cd $(FRONTEND_DIR) && npm run test:coverage -- --run

frontend-lint: ## Run frontend ESLint
	cd $(FRONTEND_DIR) && npm run lint

frontend-e2e: ## Run Playwright E2E tests (starts dev server automatically)
	cd $(FRONTEND_DIR) && npx playwright test --workers=1

frontend-e2e-ui: ## Run Playwright E2E tests with UI
	cd $(FRONTEND_DIR) && npx playwright test --ui

frontend-e2e-headed: ## Run Playwright E2E tests in headed mode
	cd $(FRONTEND_DIR) && npx playwright test --headed

frontend-e2e-debug: ## Run Playwright E2E tests in debug mode
	cd $(FRONTEND_DIR) && npx playwright test --debug

frontend-e2e-chromium: ## Run Playwright E2E tests on Chromium only
	cd $(FRONTEND_DIR) && npx playwright test --project=chromium --workers=1

frontend-e2e-report: ## Open Playwright HTML report
	cd $(FRONTEND_DIR) && npx playwright show-report

frontend-e2e-with-server: ## Run Playwright E2E tests with dev server
	cd $(FRONTEND_DIR) && npm run dev & sleep 3 && npx playwright test --workers=1 && kill %1 2>/dev/null

# Backend targets
backend-build: ## Build the Go backend
	cd $(BACKEND_DIR) && go build -o server ./cmd/server

backend-run: ## Run the Go backend (requires running PostgreSQL)
	cd $(BACKEND_DIR) && go run ./cmd/server

backend-db-up: ## Start PostgreSQL in Docker
	cd $(BACKEND_DIR) && docker-compose up -d

backend-db-down: ## Stop PostgreSQL Docker
	cd $(BACKEND_DIR) && docker-compose down

backend-docker-build: ## Build backend Docker image
	cd $(BACKEND_DIR) && docker build -t booking-backend .

backend-test: ## Run backend tests
	cd $(BACKEND_DIR) && go test ./...

backend-test-coverage: ## Run backend tests with coverage report
	cd $(BACKEND_DIR) && go test ./... -covermode=atomic -coverpkg=./... -coverprofile=coverage.out
	cd $(BACKEND_DIR) && go tool cover -html=coverage.out -o coverage.html

backend-test-verbose: ## Run backend tests with verbose output
	cd $(BACKEND_DIR) && go test ./... -v

backend-lint: ## Run quick backend lint (go vet + optional golangci-lint)
	cd $(BACKEND_DIR) && go vet ./...
	cd $(BACKEND_DIR) && if command -v golangci-lint >/dev/null 2>&1; then \
		if golangci-lint version | grep -Eq 'built with go1\.(25|26|27|28)'; then \
			golangci-lint run --config .golangci.yml ./...; \
		else \
			echo "Skipping local golangci-lint: installed binary targets older Go than this module. CI enforces $(GOLANGCI_LINT_VERSION)."; \
		fi; \
	else \
		echo "Skipping local golangci-lint: install golangci-lint $(GOLANGCI_LINT_VERSION) locally or rely on CI."; \
	fi

backend-lint-strict: ## Run backend lint with required compatible golangci-lint
	cd $(BACKEND_DIR) && go vet ./...
	cd $(BACKEND_DIR) && command -v golangci-lint >/dev/null 2>&1 || (echo "golangci-lint $(GOLANGCI_LINT_VERSION) is required for strict lint" && exit 1)
	cd $(BACKEND_DIR) && golangci-lint version | grep -Eq 'built with go1\.(25|26|27|28)' || (echo "Installed golangci-lint is built with older Go; install $(GOLANGCI_LINT_VERSION) or newer built with Go 1.25+" && exit 1)
	cd $(BACKEND_DIR) && golangci-lint run --config .golangci.yml ./...

backend-lint-fix: ## Auto-fix lint issues (format, organize imports)
	cd $(BACKEND_DIR) && gofmt -w .
	cd $(BACKEND_DIR) && go mod tidy

backend-fmt: ## Check Go formatting
	cd $(BACKEND_DIR) && gofmt -l .
	cd $(BACKEND_DIR) && test -z "$$(gofmt -l .)" || (echo "Files need formatting:" && gofmt -l . && exit 1)

backtest: backend-test backend-lint ## Run backend tests and linter

# Docker Compose targets
docker-up: ## Start all services (frontend + backend + database)
	docker-compose up -d

docker-down: ## Stop all services
	docker-compose down

docker-down-v: ## Stop all services and remove volumes
	docker-compose down -v

docker-logs: ## View logs from all services
	docker-compose logs -f

docker-logs-backend: ## View backend logs
	docker-compose logs -f backend

docker-logs-frontend: ## View frontend logs
	docker-compose logs -f frontend

docker-logs-db: ## View database logs
	docker-compose logs -f postgres

docker-build: ## Rebuild all Docker images
	docker-compose build --no-cache

docker-build-backend: ## Rebuild backend Docker image
	docker-compose build backend

docker-build-frontend: ## Rebuild frontend Docker image
	docker-compose build frontend

docker-rebuild: ## Stop, rebuild, and start all services
	docker-compose down
	docker-compose build --no-cache
	docker-compose up -d

# Deployment
SERVER ?= admn@192.168.2.53
REMOTE_PATH ?= /home/admn/qwen-ai

devdeploy: ## Sync files to server and rebuild docker containers
		rsync -avz --delete \
			--exclude='.git' \
			--exclude='.github' \
			--exclude='.worktrees' \
			--exclude='node_modules' \
			--exclude='tsp-output' \
			--exclude='frontend/dist' \
			--exclude='frontend/coverage' \
			--exclude='frontend/playwright-report' \
			--exclude='frontend/test-results' \
			--exclude='backend/coverage.out' \
			--exclude='backend/coverage.html' \
			--exclude='*.md' \
			--exclude='Makefile' \
			-e ssh \
		./ $(SERVER):$(REMOTE_PATH)
	ssh $(SERVER) "cd $(REMOTE_PATH) && docker compose build --no-cache && docker compose up -d"

docker-ps: ## Show status of all services
	docker-compose ps
