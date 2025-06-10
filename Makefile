.PHONY: install start docker-build

install:
	go install github.com/air-verse/air@latest
	go mod tidy

start:
	go run cmd/hat/main.go
docker-build:
	docker build -t hat-app -f ops/docker/Dockerfile .
