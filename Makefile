include .envrc

run/api:
	go run ./cmd/api

psql:
	psql ${EHOME_DB_DSN}