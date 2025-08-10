# HAT - Handicap Archery Tournament

## Description

HAT (Handicap Archery Tournament) is a web application for managing handicap archery tournaments built with Go and the Gin web framework. It includes user stories, personas, and navigation documentation to guide development.

## Tech Stack

- **Backend**: Go with [Gin-Gonic](https://gin-gonic.com/) web framework
- **Templating**: Go's standard `html/template` package
- **Frontend**: [Tailwind CSS](https://tailwindcss.com/) for styling and UI components

## Documentation

For detailed documentation, please see the [documentation index](doc/index.md).

## Getting Started

### Prerequisites

- Go 1.21 or later
- Node.js and npm (for Tailwind CSS)

### Installation

1. Clone the repository and navigate to the project directory
2. Initialize Go module and install dependencies:
   ```bash
   go mod tidy
   ```
3. Install Tailwind CSS dependencies:
   ```bash
   npm install -D tailwindcss
   npx tailwindcss init
   ```
4. Run the application:
   ```bash
   go run main.go
   ```

## Contributing

[Add contribution guidelines here]
