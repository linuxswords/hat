.PHONY: build clean run dev css test test-verbose test-coverage test-coverage-html test-ci test-handlers test-e2e test-e2e-ui test-e2e-debug test-all

# Build the application
build: css
	go build -o bin/hat cmd/hat/main.go

# Build CSS with Tailwind
css:
	npx tailwindcss -i ./static/css/styles.css -o ./static/css/output.css

# Clean build artifacts
clean:
	rm -rf bin/
	rm -f static/css/output.css

# Run the application
run: build
	./bin/hat

# Development mode with CSS watching
dev:
	npx tailwindcss -i ./static/css/styles.css -o ./static/css/output.css --watch &
	air

# Install dependencies
deps:
	go mod tidy
	npm install

# Run tests
test:
	go test ./...

# Run tests with verbose output  
test-verbose:
	go test -v ./...

# Run tests with coverage
test-coverage:
	go test -cover ./...

# Run tests with coverage report
test-coverage-html:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run tests for CI
test-ci:
	@echo "Running full test suite for CI"
	go test ./...

# Run only handler tests
test-handlers:
	go test ./internal/handlers/...

# Frontend/E2E Tests
test-e2e:
	@echo "Running Playwright E2E tests"
	npm run test:e2e

# Run E2E tests with UI mode
test-e2e-ui:
	@echo "Running Playwright E2E tests in UI mode"
	npm run test:e2e:ui

# Debug E2E tests
test-e2e-debug:
	@echo "Running Playwright E2E tests in debug mode"
	npm run test:e2e:debug

# Run all tests (backend + frontend)
test-all: test test-e2e
	@echo "All tests completed successfully"

# Setup Playwright for first time
test-e2e-setup:
	@echo "Setting up Playwright for E2E testing"
	npm install
	npm run playwright:install
