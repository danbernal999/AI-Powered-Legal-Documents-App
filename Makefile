.DEFAULT_GOAL := help

MIGRATE_IMAGE ?= migrate/migrate:v4.15.2
MIGRATE ?= migrate
DATABASE_URL ?= $(shell grep -E '^DATABASE_URL=' .env 2>/dev/null | cut -d'=' -f2- || echo "postgres://user:password@localhost:5432/kiradoc?sslmode=disable")

help:
	@echo "Makefile commands:"
	@echo "  make db-migrate     - apply up migrations (uses local migrate or docker image)"
	@echo "  make db-down        - apply down migrations (uses local migrate or docker image)"
	@echo "  make db-create      - start postgres and apply migrations"
	@echo "  make up             - docker compose up (builds images)"
	@echo "  make down           - docker compose down"

db-migrate:
	@if command -v $(MIGRATE) >/dev/null 2>&1; then \
		$(MIGRATE) -path ./backend/migrations -database "$(DATABASE_URL)" up; \
	else \
		docker run --rm -v $(PWD)/backend/migrations:/migrations -e DATABASE_URL="$(DATABASE_URL)" $(MIGRATE_IMAGE) -path=/migrations -database "$(DATABASE_URL)" up; \
	fi

db-down:
	@if command -v $(MIGRATE) >/dev/null 2>&1; then \
		$(MIGRATE) -path ./backend/migrations -database "$(DATABASE_URL)" down; \
	else \
		docker run --rm -v $(PWD)/backend/migrations:/migrations -e DATABASE_URL="$(DATABASE_URL)" $(MIGRATE_IMAGE) -path=/migrations -database "$(DATABASE_URL)" down; \
	fi

db-create:
	docker compose up -d postgres
	@sleep 2
	$(MAKE) db-migrate

up:
	docker compose up -d --build

down:
	docker compose down
