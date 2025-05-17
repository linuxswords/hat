.PHONY: install start docker-build

install:
	go install github.com/air-verse/air@latest
	go mod tidy

dev:
	air -c air.toml

start:
	go run cmd/hat/main.go
docker-build:
	docker build -t hat-app -f ops/docker/Dockerfile .
