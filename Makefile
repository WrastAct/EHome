include .envrc

.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

.PHONY: confirm
confirm:
	@echo -n 'Are you sure? [y/N] ' && read ans && [ $${ans:-N} = y ]

## run/api: run the cmd/api application
.PHONY: run/api
run/api:
	@go run ./cmd/api -db-dsn=${EHOME_DB_DSN}

## test/api: run the cmd/api application without rate limiting
.PHONY: test/api
test/api:
	@go run ./cmd/api -db-dsn=${EHOME_DB_DSN} -limiter-enabled=false

## db/psql: connect to database
.PHONY: db/psql
db/psql:
	@psql ${EHOME_DB_DSN}

## db/migrations/new name=$1: create a new database migration
.PHONY: db/migrations/new
db/migrations/new:
	migrate create -seq -ext=.sql -dir=./migrations ${name}

## db/migrations/up: apply all up database migrations
.PHONY: db/migrations/up
db/migrations/up: confirm
	@echo 'Migrations up'
	@migrate -path=./migrations -database=${EHOME_DB_DSN} up

## db/migrations/down: apply all down database migrations
.PHONY: db/migrations/down
db/migrations/down: confirm
	@echo 'Migrations down'
	@migrate -path=./migrations -database=${EHOME_DB_DSN} down

	