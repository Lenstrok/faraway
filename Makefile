run-server-app:
	go run cmd/server/main.go

run-parser-app:
	go run cmd/parser/main.go

run-migration-app:
	go run cmd/migration/main.go

run-linters:
	golangci-lint run ./...

run-unit-tests:
	go test ./internal/...

run-dev-docker-compose:
	docker compose -f docker-compose-dev.yml up -d

stop-dev-docker-compose:
	docker compose -f docker-compose-dev.yml down
