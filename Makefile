.PHONY: install start

install:
	go mod tidy

start:
	go run cmd/main/main.go
