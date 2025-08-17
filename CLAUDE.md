# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

HAT (Handicap Archery Tournament) is a web application for managing handicap archery tournaments. It enables tournament organizers to manage archers, bow classes, and calculate handicap-adjusted scores for fair competition across different bow types.

## Architecture

### Tech Stack
- **Backend**: Go 1.24.2+ with Gin-Gonic web framework
- **Templates**: Go's standard `html/template` with custom layout system
- **Frontend**: Tailwind CSS for styling and responsive design
- **Data Storage**: Currently JSON files (database integration planned)
- **Deployment**: Docker containers on Google Cloud Platform (GKE)

### Project Structure
- `cmd/hat/main.go` - Application entry point with route definitions
- `internal/` - Core application logic
  - `handlers/` - HTTP request handlers and template rendering
  - `models/` - Domain models (Archer, Tournament, BowClass, Handicap, Score)
  - `repositories/` - Data access layer (currently JSON file-based)
- `templates/` - HTML templates with layouts and partials
- `static/` - CSS, images, and other static assets
- `doc/` - Comprehensive project documentation
- `ops/` - Infrastructure and deployment configurations

### Key Domain Concepts
- **Archer**: Athletes with name, gender, bow class, and contact info
- **Bow Class**: Categories like "Adult Male Longbow" that determine handicap factors
- **Handicap**: Multipliers to normalize scores across different bow classes
- **Tournament**: Competitions with multiple archers using handicap scoring
- **Raw Score vs Adjusted Score**: Raw scores multiplied by bow class handicaps

## Development Commands

### Building and Running
```bash
make build       # Build the application binary
make run         # Build and run the application on :8080
make dev         # Development mode with CSS watching (requires air)
make css         # Build Tailwind CSS once
make clean       # Remove build artifacts
make deps        # Install Go and npm dependencies
```

### Testing
```bash
make test                # Run all tests
make test-verbose        # Run tests with verbose output
make test-coverage       # Run tests with coverage report
make test-coverage-html  # Generate HTML coverage report
make test-ci             # Run tests for CI environments
make test-handlers       # Run only handler tests
```

**Test Coverage**: The test suite provides comprehensive coverage of handler business logic including:
- CRUD operations for archers, tournaments, and scores
- Input validation and error handling
- Template rendering and form processing
- Repository integration with mocks
- HTTP status codes and redirects

### CSS Development
```bash
npx tailwindcss -i ./static/css/styles.css -o ./static/css/output.css --watch
```

### Development Tools
- Use `air` for hot reloading during development
- Tailwind config targets `./templates/**/*.html` for CSS purging

## Template System

The application uses a custom layout system in `internal/handlers/renderer.go`:
- Base layout: `templates/layouts/base.html`
- Shared partials: `templates/partials/`
- Page templates: `templates/{section}/{page}.html`
- Template functions include `contains` for string matching

## Data Layer

Currently uses JSON file persistence via repositories:
- Archer data stored in `doc/data/archers/archers.json`
- Auto-incrementing IDs and atomic file operations
- Repository pattern abstracts data access
- Database integration planned for production

## Infrastructure and Deployment

### Docker
- Multi-stage build with Go 1.24.2
- Debian slim runtime with SQLite support
- Exposes port 8080
- Copies static assets and templates

### GCP Deployment (ops/infra/)
- Terraform configuration for GKE cluster
- Uses e2-medium preemptible nodes
- Located in us-central1
- Service account authentication required

## Project Status

This is an active development project with the following completed:
- Basic archer CRUD operations
- Bow class management
- Template rendering system
- Docker containerization
- Terraform infrastructure setup

Planned features:
- Database integration (likely PostgreSQL)
- Tournament management
- Score tracking and handicap calculations
- User authentication
- CI/CD pipeline
