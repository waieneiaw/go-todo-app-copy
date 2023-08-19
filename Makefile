.DEFAULT_GOAL := help

DOCKER_TAG := latest

.PHONY: build
build: ## Build docker image to deploy
	docker build \
		-t go-todo-app-copy:${DOCKER_TAG} \
		--target deploy ./

.PHONY: build-local
build-local: ## Build docker image to local development
	docker compose build --no-cache

.PHONY: up
up: ## Do docker compose up with hot reload
	docker compose up -d

.PHONY: down
down: ## Do docker compose down
	docker compose down

.PHONY: logs
logs: ## Tail docker compose logs
	docker compose logs -f

.PHONY: ps
ps: ## Check container status
	docker compose ps

.PHONY: test
test: ## Execute tests
	go test -race -shuffle=on ./...

.PHONY: help
help: ## Show options
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
