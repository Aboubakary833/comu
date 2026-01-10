APP_ADDR ?= ":4000"
DB_DRIVER ?= "mysql"
DB_SOURCE ?= "abubakr:root@/comu_db?parseTime=true"
MIGRATION_PATH="./migrations"

help:
	@echo "Comu makefile commands\n"

	@echo "List of available commands:\n"
	@echo " - \"dev\"\n\tRun the development serve without building the app.\n"
	@echo " - \"make migrations-status\"\n\tCheck the database migrations status\n"
	@echo " - \"make db-migrate-up\"\n\tMigrate all the modules migration schema definitions\n"
	@echo " - \"make db-migrate-down\"\n\tRollback all the modules migrations\n"

dev:
	@APP_ADDR=$(APP_ADDR) DB_DRIVER=$(DB_DRIVER) DB_SOURCE=$(DB_SOURCE) go run ./cmd/web


# Database related commands

migrations-status:
	@goose $(DB_DRIVER) $(DB_SOURCE) -dir=$(MIGRATION_PATH) status

migrations-up:
	@goose $(DB_DRIVER) $(DB_SOURCE) -dir=$(MIGRATION_PATH) up

migrations-down:
	@goose $(DB_DRIVER) $(DB_SOURCE) -dir=$(MIGRATION_PATH) down
