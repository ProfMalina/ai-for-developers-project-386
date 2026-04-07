.PHONY: help compile fmt fmt-check lint openapi clean install-dev dev backend-build backend-run backend-db-up backend-db-down backend-docker-build \
        docker-up docker-down docker-logs docker-build docker-rebuild frontend-test fronttest frontend-coverage frontend-lint \
        backend-test backend-test-coverage backend-test-verbose backend-lint backend-lint-fix backend-fmt backtest

TYPESPEC_DIR := typespec
BACKEND_DIR := backend
FRONTEND_DIR := frontend

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

frontend-e2e: ## Run Playwright E2E tests
	cd $(FRONTEND_DIR) && npx playwright test

frontend-e2e-ui: ## Run Playwright E2E tests with UI
	cd $(FRONTEND_DIR) && npx playwright test --ui

frontend-e2e-headed: ## Run Playwright E2E tests in headed mode
	cd $(FRONTEND_DIR) && npx playwright test --headed

frontend-e2e-debug: ## Run Playwright E2E tests in debug mode
	cd $(FRONTEND_DIR) && npx playwright test --debug

frontend-e2e-chromium: ## Run Playwright E2E tests on Chromium only
	cd $(FRONTEND_DIR) && npx playwright test --project=chromium

frontend-e2e-report: ## Open Playwright HTML report
	cd $(FRONTEND_DIR) && npx playwright show-report

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
	cd $(BACKEND_DIR) && go test ./... -coverprofile=coverage.out
	cd $(BACKEND_DIR) && go tool cover -html=coverage.out -o coverage.html

backend-test-verbose: ## Run backend tests with verbose output
	cd $(BACKEND_DIR) && go test ./... -v

backend-lint: ## Run backend linter (go vet + static analysis)
	cd $(BACKEND_DIR) && go vet ./...
	cd $(BACKEND_DIR) && go run honnef.co/go/tools/cmd/staticcheck@latest ./...

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
		--exclude='node_modules' \
		--exclude='tsp-output' \
		--exclude='*.md' \
		--exclude='Makefile' \
		-e ssh \
		./ $(SERVER):$(REMOTE_PATH)
	ssh $(SERVER) "cd $(REMOTE_PATH) && docker compose build --no-cache && docker compose up -d"

docker-ps: ## Show status of all services
	docker-compose ps
