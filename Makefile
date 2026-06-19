.PHONY: build test lint fmt docker-build docker-up docker-down

build:
	go build ./cmd/local-solar-time

test:
	go test ./...

lint:
	golangci-lint run

fmt:
	goimports -w .

docker-build:
	docker compose build

docker-up:
	docker compose up

docker-down:
	docker compose down
