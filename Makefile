.PHONY: help run sqlc-generate docker-build docker-up docker-down docker-logs test test-coverage swagger

help:
	@echo "Available commands:"
	@echo "  make run            - Run the application locally"
	@echo "  make sqlc-generate  - Generate code from SQL queries"
	@echo "  make swagger        - Generate Swagger documentation"
	@echo "  make test           - Run tests"
	@echo "  make test-coverage  - Run tests with coverage"
	@echo "  make docker-build   - Build Docker image"
	@echo "  make docker-up      - Start services with docker-compose"
	@echo "  make docker-down    - Stop services"
	@echo "  make docker-logs    - View logs"

run:
	go run cmd/api/main.go

sqlc-generate:
	~/go/bin/sqlc generate

swagger:
	~/go/bin/swag init -g cmd/api/main.go -o docs

test:
	go test -v ./...

test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

docker-build:
	docker-compose build

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

docker-logs:
	docker-compose logs -f app
