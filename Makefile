.PHONY: run dev build migrate-up migrate-down migrate-status

run:
	go run ./cmd/api

dev:
	air

build:
	go build -o bin/api ./cmd/api

DB_URL ?= postgres://postgres:postgres@localhost:5432/go_rest_api_db?sslmode=disable

migrate-up:
	goose -dir migrations postgres "$(DB_URL)" up

migrate-down:
	goose -dir migrations postgres "$(DB_URL)" down

migrate-status:
	goose -dir migrations postgres "$(DB_URL)" status