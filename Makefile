.PHONY: debug fe server sync_db dump_db middleware web down clean python help

# 定义脚本路径
SCRIPTS_DIR := ./scripts
BUILD_FE_SCRIPT := $(SCRIPTS_DIR)/build_fe.sh
BUILD_SERVER_SCRIPT := $(SCRIPTS_DIR)/setup/server.sh
SETUP_DOCKER_SCRIPT := $(SCRIPTS_DIR)/setup/docker.sh
SETUP_PYTHON_SCRIPT := $(SCRIPTS_DIR)/setup/python.sh
COMPOSE_FILE := docker/docker-compose-debug.yml
MYSQL_SCHEMA := ./docker/volumes/mysql/schema.sql
ENV_FILE := ./docker/.env.debug
WEB_ENV_FILE := ./docker/.env
STATIC_DIR := ./bin/resources/static

debug: env middleware python fe server

env:
	@if [ ! -f "$(ENV_FILE)" ]; then \
		echo "Env file '$(ENV_FILE)' not found, using example env..."; \
		cp ./docker/.env.debug.example $(ENV_FILE); \
	fi

fe:
	@echo "Building frontend..."
	@bash $(BUILD_FE_SCRIPT)

server: env
	@if [ ! -d "$(STATIC_DIR)" ]; then \
		echo "Static directory '$(STATIC_DIR)' not found, building frontend..."; \
		$(MAKE) fe; \
	fi
	@echo "Building and run server..."
	@APP_ENV=debug bash $(BUILD_SERVER_SCRIPT) -start


build_server:
	@echo "Building server..."
	@bash $(BUILD_SERVER_SCRIPT)

sync_db: env
	@echo "Syncing database..."
	@docker compose -f $(COMPOSE_FILE) --env-file $(ENV_FILE) --profile mysql-setup up -d

sql_init:
	@echo "Init sql data..."
	@docker compose -f $(COMPOSE_FILE) --env-file $(ENV_FILE) --profile mysql-setup up -d

middleware:
	@echo "Start middleware docker environment for opencoze app"
	@docker compose -f $(COMPOSE_FILE) --env-file $(ENV_FILE) --profile middleware up -d --wait

build_docker:
	@echo "Build docker image"
	@docker compose -f $(COMPOSE_FILE) --profile build-server build

web_env:
	@if [ ! -f "$(WEB_ENV_FILE)" ]; then \
		echo "Env file '$(WEB_ENV_FILE)' not found, using example env..."; \
		cp ./docker/.env.example $(WEB_ENV_FILE); \
	fi

web: web_env
	@echo "Start web server in docker"
	@docker compose -f docker/docker-compose.yml --env-file $(WEB_ENV_FILE) up -d --build

down_web:
	@echo "Stop web server in docker"
	@docker compose -f docker/docker-compose.yml --env-file $(WEB_ENV_FILE) down

down: env
	@echo "Stop all docker containers"
	@docker compose -f $(COMPOSE_FILE) --profile '*' down

clean: down
	@echo "Remove docker containers and volumes data"
	@rm -rf ./docker/data

python:
	@echo "Setting up Python..."
	@bash $(SETUP_PYTHON_SCRIPT)

sync_helm_schema:
	@cp $(MYSQL_SCHEMA) ./helm/charts/opencoze/files/mysql/schema.sql

help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@echo "  debug            - Start the debug environment."
	@echo "  env              - Setup env file."
	@echo "  fe               - Build the frontend."
	@echo "  server           - Build and run the server binary."
	@echo "  build_server     - Build the server binary."
	@echo "  sync_db          - Create missing lightweight tables without altering existing tables."
	@echo "  sql_init         - Init sql data..."
	@echo "  middleware       - Setup middlewares docker environment, but exclude the server app."
	@echo "  web              - Setup web docker environment, include middlewares docker."
	@echo "  down             - Stop the docker containers."
	@echo "  down_web         - Stop the web docker containers."
	@echo "  clean            - Stop the docker containers and clean volumes."
	@echo "  python           - Setup python environment."
	@echo ""
	@echo ""
	@echo "  help             - Show this help message."
