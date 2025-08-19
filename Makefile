.PHONY: build clean run dev css test test-verbose test-coverage test-coverage-html test-ci test-handlers test-e2e test-e2e-ui test-e2e-debug test-all db-up db-down db-setup db-migrate db-seed

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

# Database Operations
db-up:
	@echo "Starting PostgreSQL containers..."
	docker-compose up -d postgres postgres_test
	@echo "Waiting for PostgreSQL to be ready..."
	sleep 10

db-down:
	@echo "Stopping PostgreSQL containers..."
	docker-compose down

db-setup: db-up
	@echo "Setting up database..."
	@echo "PostgreSQL is running on:"
	@echo "  Development: localhost:5432"
	@echo "  Test: localhost:5433"
	@echo ""
	@echo "Database credentials:"
	@echo "  User: hat_user"
	@echo "  Password: hat_password"
	@echo "  Dev DB: hat_development"
	@echo "  Test DB: hat_test"

db-migrate:
	@echo "Running database migrations..."
	@echo "Migrations will run automatically when the application starts"

db-seed:
	@echo "Seeding database..."
	@echo "Seed data will be inserted automatically when the application starts"

# Environment variables for development
dev-env:
	@echo "Setting up development environment variables..."
	@echo "export DB_HOST=localhost"
	@echo "export DB_PORT=5432"
	@echo "export DB_USER=hat_user"
	@echo "export DB_PASSWORD=hat_password"
	@echo "export DB_NAME=hat_development"
	@echo "export TEST_DB_HOST=localhost"
	@echo "export TEST_DB_PORT=5433"
	@echo "export TEST_DB_USER=hat_user"
	@echo "export TEST_DB_PASSWORD=hat_password"
	@echo "export TEST_DB_NAME=hat_test"
