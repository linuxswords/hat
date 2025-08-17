.PHONY: build clean run dev css test test-verbose test-coverage test-coverage-html test-ci test-handlers

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
