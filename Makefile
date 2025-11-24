ENV ?= dev

COMPOSE_FILE = docker-compose.dev.yaml
ifeq ($(ENV),prod)
	COMPOSE_FILE = docker-compose.yaml
endif

DC = docker compose -f $(COMPOSE_FILE)

up:
	$(DC) up

build:
	$(DC) build

down:
	$(DC) down
