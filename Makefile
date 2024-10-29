# Makefile

.PHONY: help up down restart build logs ps clean prune reload

help: ## Show help information.
	@echo "Available commands:"
	@awk 'BEGIN {FS = ":.*?##"; printf "\nUsage:\n"} /^[a-zA-Z0-9_-]+:.*?##/ { printf "  make %-10s %s\n", $$1, $$2 }' $(MAKEFILE_LIST)

up: ## Starts the services.
	docker-compose up -d

down: ## Stops and removes the containers.
	docker-compose down

restart: down up ## Restarts the services.

build: ## Rebuilds the images.
	docker-compose build

logs: ## Follows the logs.
	docker-compose logs -f

ps: ## Lists the running services.
	docker-compose ps

reload: ## Applies configuration changes without restarting services.
	docker-compose up -d

clean: ## Cleans up all containers, networks, volumes, and images.
	docker-compose down --rmi all --volumes --remove-orphans

prune: ## Prunes the Docker system (WARNING: removes unused data).
	docker system prune -af --volumes
