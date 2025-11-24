# Сервис назначения ревьюеров для Pull Request’ов

## Описание

Сервис, который назначает ревьюеров на PR из команды автора, позволяет выполнять переназначение ревьюверов и получать список PR’ов, назначенных конкретному пользователю, а также управлять командами и активностью пользователей. После merge PR изменение состава ревьюверов запрещено.

## Запуск

При development сервис запускается с помощью Air для hot reloading и также поднимается adminer

### С помощью Makefile

- Production: `make up ENV=prod`
- Development: `make up`

### С помощью Docker Compose

- Production: `docker compose up`
- Development: `docker compose -f docker-compose.dev.yaml up`
