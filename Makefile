.PHONY: help compile fmt fmt-check lint openapi clean

TYPESPEC_DIR := typespec

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
