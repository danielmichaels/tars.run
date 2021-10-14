include .env
default: help

# ==================================================================================== #
# HELPERS
# ==================================================================================== #

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

.PHONY: confirm
confirm:
	@echo -n 'Are you sure? [y/N] ' && read ans && [ $${ans:-N} = y ]

# ==================================================================================== #
# DEVELOPMENT
# ==================================================================================== #

## run/api: run the cmd/api application
.PHONY: run/api
run/api:
	go run ./cmd/api


## run/api/development: run the API server with hot-reloading (uses `air`)
.PHONY: run/api/development
run/api/development:
	@air

## db/migrations/up: Run migration up (apply migrations)
.PHONY: db/migration/up
db/migration/up: confirm
	@echo "setting up the database extensions..."
	goose -dir ./migrations sqlite3 ./${DB_NAME} up

## db/migration/down/to version=$1: Roll back migration to $1
.PHONY: db/migration/down/to
db/migration/down/to: confirm
	@echo "rolling back migrations to ${version}"
	goose -dir ./migrations sqlite3 ./${DB_NAME} down-to ${version}

## db/migration/down: drop all migrations
.PHONY: db/migration/down
db/migration/down: confirm
	@echo "rolling back all migrations"
	goose -dir ./migrations sqlite3 ./${DB_NAME} down

## db/migration/redo: rollback latest migration, then reapply
.PHONY: db/migration/redo
db/migration/redo: confirm
	@echo "rolling back latest migration, then re-applying it; redo"
	goose -dir ./migrations sqlite3 ./${DB_NAME} redo

## db/migration/create name=$1: Create new migrations with name of $1
.PHONY: db/migration/create
db/migration/create: confirm
	@echo "creating migration files for ${name} in ${DB_NAME}"
	goose -dir ./migrations sqlite3 ./${DB_NAME} create ${name} sql

## db/migration/status: Check database migration status
.PHONY: db/migration/status
db/migration/status:
	@echo "checking migration status for ${DB_NAME}"
	goose -dir ./migrations sqlite3 ./${DB_NAME} status
# ==================================================================================== #
# QUALITY CONTROL
# ==================================================================================== #

## audit: tidy and vendor dependencies and format, vet and test all code
.PHONY: audit
audit: vendor
	@echo 'Formatting code...'
	go fmt ./...
	@echo 'Vetting code...'
	go vet ./...
	staticcheck ./...
	@echo 'Running tests...'
	go test -race -vet=off ./...

## vendor: tidy and vendor dependencies
.PHONY: vendor
vendor:
	@echo 'Tidying and verifying module dependencies...'
	go mod tidy
	go mod verify
	@echo 'Vendoring dependencies...'
	go mod vendor

# ==================================================================================== #
# BUILD
# ==================================================================================== #

current_time = $(shell date --iso-8601=seconds)
git_description = $(shell git describe --always --dirty --tags --long)
linker_flags = '-s -X main.buildTime=${current_time} -X main.version=${git_description}'

## build/api: build the cmd/api application
.PHONY: build/api
build/api:
	@echo 'Building cmd/api...'
	go build -ldflags=${linker_flags} -o=./bin/api ./cmd/api
	GOOS=linux GOARCH=amd64 go build -ldflags=${linker_flags} -o=./bin/linux_amd64/api ./cmd/api

# ==================================================================================== #
# PRODUCTION
# ====================================================================================

# build a docker container?
