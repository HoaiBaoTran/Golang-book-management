DATABASE_URL ?= "postgresql://postgres:hoaibao1142@localhost:5432/book_management?sslmode=disable"

MIGRATION_PATH ?= pkg/database/migration

migrate-up:
	migrate -path=$(MIGRATION_PATH) -database=$(DATABASE_URL) -verbose up

migrate-down:
	migrate -path=$(MIGRATION_PATH) -database=$(DATABASE_URL) -verbose down

migrate-create:
	migrate create -ext=sql -dir=$(MIGRATION_PATH) $(name)

.PHONY: migrate-up migrate-down migrate-create