.DEFAULT_GOAL := help

NAME = $(shell basename $(PWD))
AWAIT_DB_SCRIPT = build/wait-for-it.sh

.PHONY: help
help:
	@echo "------------------------------------------------------------------------"
	@echo "${NAME}"
	@echo "------------------------------------------------------------------------"
	@grep -E '^[a-zA-Z0-9_/%\-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: proto
proto: ## Generate gRPC code from proto files
	@protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		proto/users/v1/*.proto

.PHONY: build
build: ## Build the application
	@GOOS=linux go build -o $(NAME) main.go

.PHONY: run
run: build ## Run the application on a Docker container (requires Docker)
	@[ -f $(AWAIT_DB_SCRIPT) ] || curl -o $(AWAIT_DB_SCRIPT) "https://raw.githubusercontent.com/vishnubob/wait-for-it/master/wait-for-it.sh"
	@chmod +x $(AWAIT_DB_SCRIPT)
	@docker-compose -f build/docker-compose.yml up db usrsvc --force-recreate --build

.PHONY: lint
lint: ## Run go fmt and go vet
	@go fmt ./...
	@go vet ./...

.PHONY: test-unit
test-unit: ## Run unit tests
	@go test -v -race -vet=all -count=1 -timeout 60s ./app/... ./internal/...

.PHONY: test-it
test-it: ## Run integration tests (requires Docker)
	@docker-compose -f build/docker-compose.yml up -d db
	@sleep 3 # wait for db to be ready
	@go test -v -race -tags=integration -vet=all -count=1 -timeout 60s ./app/... ./internal/...
	@docker-compose -f build/docker-compose.yml down

.PHONY: test-e2e
test-e2e: ## Run end-to-end tests (requires Docker)
	@docker-compose -f build/docker-compose.yml up -d db
	@sleep 3 # wait for db to be ready
	@go test -v -tags=e2e -race -vet=all -count=1 -timeout 60s ./tests/...
	@docker-compose -f build/docker-compose.yml down

.PHONY: test ## Run all tests (lint, unit, integration, and end-to-end)
test: lint test-unit test-it test-e2e ## Run all tests
