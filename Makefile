include .envrc

run/api:
	go run ./cmd/api

psql:
	psql ${EHOME_DB_DSN}

migrations_up:
	migrate -path=./migrations -database=${EHOME_DB_DSN} up

migrations_down:
	migrate -path=./migrations -database=${EHOME_DB_DSN} down

	