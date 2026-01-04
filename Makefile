.PHONY: run dev build docs migrate-up migrate-down migrate-status lint lint-fix test vet fmt tidy

run:
	go run ./cmd/api

dev:
	air

build:
	go build -o bin/api ./cmd/api

docs:
	swag init -g cmd/api/main.go

DB_URL ?= postgres://postgres:postgres@localhost:5432/go_rest_api_db?sslmode=disable

migrate-up:
	goose -dir migrations postgres "$(DB_URL)" up

migrate-down:
	goose -dir migrations postgres "$(DB_URL)" down

migrate-status:
	goose -dir migrations postgres "$(DB_URL)" status

# Linting
lint:
	@echo "Running golangci-lint..."
	@if ! command -v golangci-lint > /dev/null; then \
		echo "golangci-lint not found. Installing..."; \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin v1.60.1; \
	fi
	golangci-lint run

lint-fix:
	@echo "Running golangci-lint with auto-fix..."
	@if ! command -v golangci-lint > /dev/null; then \
		echo "golangci-lint not found. Installing..."; \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin v1.60.1; \
	fi
	golangci-lint run --fix

# Code quality checks
test:
	go test -v -race -coverprofile=coverage.out ./...

test-coverage:
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

vet:
	go vet ./...

fmt:
	@echo "Checking code formatting..."
	@if [ $$(gofmt -l . | wc -l) -ne 0 ]; then \
		echo "Code is not formatted. Run 'make fmt-fix' to fix."; \
		gofmt -d .; \
		exit 1; \
	fi
	@echo "Code is properly formatted."

fmt-fix:
	@echo "Fixing code formatting..."
	gofmt -w .
	goimports -w -local github.com/mrhpn/go-rest-api .

tidy:
	go mod tidy
	go mod verify

# Run all quality checks
check: fmt vet lint test
	@echo "All quality checks passed!"

# CI-friendly check (fails on any issue)
ci-check: fmt vet lint
	@echo "CI checks completed."