run-server-app:
	go run cmd/server/main.go

run-client-app:
	go run cmd/client/main.go

run-linters:
	golangci-lint run ./...

run-unit-tests:
	go test ./internal/...
