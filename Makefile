DATABASE_URL ?= "postgresql://postgres:hoaibao1142@localhost:5432/book_management?sslmode=disable"

MIGRATION_PATH ?= pkg/database/migration

migrate-v1-force:
	migrate -path=$(MIGRATION_PATH) -database=$(DATABASE_URL) force 1

migrate-up-v1:
	migrate -path=$(MIGRATION_PATH) -database=$(DATABASE_URL) -verbose up 1

migrate-down-v1:
	migrate -path=$(MIGRATION_PATH) -database=$(DATABASE_URL) -verbose down 1

migrate-v2-force:
	migrate -path=$(MIGRATION_PATH) -database=$(DATABASE_URL) force 2

migrate-up-v2:
	migrate -path=$(MIGRATION_PATH) -database=$(DATABASE_URL) -verbose up 2

migrate-down-v2:
	migrate -path=$(MIGRATION_PATH) -database=$(DATABASE_URL) -verbose down 2

migrate-create:
	migrate create -ext=sql -dir=$(MIGRATION_PATH) $(name)

.PHONY: migrate-up-v1 migrate-down-v1 migrate-up-v2 migrate-create