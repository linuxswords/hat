.PHONY: install start

install:
	go install github.com/air-verse/air@latest
	go mod tidy

dev:
	air -c air.toml

start:
	go run cmd/main/main.go
