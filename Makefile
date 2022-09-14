include .env
export

run:
	go run main.go

.PHONY: mock
mock:
	go generate ./...

test:
	go test ./...

test/bench:
	go test ./... -bench=.

test/cov:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out

migrate:
	migrate -database ${MIGRATE_URL} -path db/migrations up

swag:
	swag init