.PHONY: help \
        dev stop down restart logs \
        migrate migrate-down migrate-status migrate-create \
        test lint \
        prod build \
        moderator shell-backend shell-frontend shell-db \
        clean

# ─── Config ───────────────────────────────────────────────────────────────────
DOCKER_DB_URL := postgres://goodcity:goodcity@postgres:5432/goodcity?sslmode=disable
GOOSE         := GOFLAGS=-mod=mod go run github.com/pressly/goose/v3/cmd/goose@latest
MIGRATIONS    := /migrations
DC            := docker compose
DC_RUN        := $(DC) run --rm

# ─── Help (default) ───────────────────────────────────────────────────────────
.DEFAULT_GOAL := help

help:
	@echo ""
	@echo "  GoodCity — все команды через Docker"
	@echo ""
	@echo "  Разработка"
	@echo "    make dev                       — поднять весь стек (postgres + backend + frontend)"
	@echo "    make stop                      — остановить контейнеры (сохранить данные)"
	@echo "    make down                      — остановить и удалить контейнеры"
	@echo "    make restart                   — перезапустить все сервисы"
	@echo "    make restart service=backend   — перезапустить конкретный сервис"
	@echo "    make logs                      — логи всех сервисов"
	@echo "    make logs service=backend      — логи конкретного сервиса"
	@echo ""
	@echo "  Миграции"
	@echo "    make migrate                   — применить все миграции"
	@echo "    make migrate-down              — откатить последнюю миграцию"
	@echo "    make migrate-status            — статус миграций"
	@echo "    make migrate-create NAME=foo   — создать новый файл миграции"
	@echo ""
	@echo "  Тесты и линтер"
	@echo "    make test                      — go test ./..."
	@echo "    make lint                      — go vet + eslint"
	@echo ""
	@echo "  Прод"
	@echo "    make prod                      — запустить прод-сборку (build + up)"
	@echo "    make build                     — пересобрать прод-образы"
	@echo ""
	@echo "  Утилиты"
	@echo "    make moderator EMAIL=u@ex.com  — выдать роль модератора"
	@echo "    make shell-backend             — sh в контейнер backend"
	@echo "    make shell-frontend            — sh в контейнер frontend"
	@echo "    make shell-db                  — psql в базу данных"
	@echo "    make clean                     — удалить контейнеры, тома, кеши"
	@echo ""

# ─── Development ──────────────────────────────────────────────────────────────
dev:
	$(DC) up

stop:
	$(DC) stop

down:
	$(DC) down

restart:
ifdef service
	$(DC) restart $(service)
else
	$(DC) restart
endif

logs:
ifdef service
	$(DC) logs -f $(service)
else
	$(DC) logs -f
endif

# ─── Migrations ───────────────────────────────────────────────────────────────
# Postgres запускается автоматически как зависимость backend-сервиса.
# Для ручного управления миграциями используем отдельный контейнер.

migrate:
	$(DC) up -d postgres
	@echo "→ Ожидание postgres..."
	@until $(DC) exec postgres pg_isready -U goodcity > /dev/null 2>&1; do printf '.'; sleep 1; done
	@echo " готов"
	$(DC_RUN) backend sh -c '$(GOOSE) -dir $(MIGRATIONS) postgres "$(DOCKER_DB_URL)" up'

migrate-down:
	$(DC_RUN) backend sh -c '$(GOOSE) -dir $(MIGRATIONS) postgres "$(DOCKER_DB_URL)" down'

migrate-status:
	$(DC_RUN) backend sh -c '$(GOOSE) -dir $(MIGRATIONS) postgres "$(DOCKER_DB_URL)" status'

migrate-create:
	@[ -n "$(NAME)" ] || (echo "Укажите NAME=<название>. Пример: make migrate-create NAME=add_tags" && exit 1)
	$(DC_RUN) backend sh -c '$(GOOSE) -dir $(MIGRATIONS) create $(NAME) sql'
	@echo "→ Файл миграции создан в migrations/"

# ─── Tests & Lint ─────────────────────────────────────────────────────────────
test:
	$(DC_RUN) backend go test ./...

lint:
	$(DC_RUN) backend go vet ./...
	$(DC_RUN) frontend sh -c "npm install --silent && npm run lint"

# ─── Production ───────────────────────────────────────────────────────────────
PROD_DC := $(DC) -f docker-compose.yml -f docker-compose.prod.yml

prod: build
	$(PROD_DC) up -d

build:
	$(PROD_DC) build --build-arg VITE_DADATA_TOKEN=$(VITE_DADATA_TOKEN)

# ─── Utilities ────────────────────────────────────────────────────────────────
moderator:
	@[ -n "$(EMAIL)" ] || (echo "Укажите EMAIL=user@example.com" && exit 1)
	$(DC) exec postgres psql -U goodcity -d goodcity \
		-c "UPDATE users SET role='moderator' WHERE email='$(EMAIL)' RETURNING id, email, role;"

shell-backend:
	$(DC) exec backend sh

shell-frontend:
	$(DC) exec frontend sh

shell-db:
	$(DC) exec postgres psql -U goodcity -d goodcity

# ─── Clean ────────────────────────────────────────────────────────────────────
clean:
	$(DC) down -v --remove-orphans
	rm -rf frontend/dist/ bin/
	@echo "→ Контейнеры, тома и артефакты удалены"
