.PHONY: build clean run dev css

# Build the application
build: css
	go build -o bin/hat main.go

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
