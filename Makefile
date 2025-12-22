.PHONY: run dev build

run:
	go run ./cmd/api

dev:
	air

build:
	go build -o bin/api ./cmd/api