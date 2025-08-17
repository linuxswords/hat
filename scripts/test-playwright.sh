#!/bin/bash

# Playwright test runner script for HAT application
# This script helps run and debug Playwright tests

set -e

echo "🎯 HAT Playwright Test Runner"
echo "=============================="

# Check if dependencies are installed
if [ ! -d "node_modules" ]; then
    echo "📦 Installing npm dependencies..."
    npm install
fi

# Check if Playwright browsers are installed
if [ ! -d "node_modules/@playwright/test" ]; then
    echo "🎭 Installing Playwright..."
    npm install @playwright/test
fi

echo "🌐 Installing Playwright browsers..."
npx playwright install chromium

# Build CSS
echo "🎨 Building CSS..."
npm run build-css-once

# Build Go application
echo "🔨 Building Go application..."
make build

# Function to run tests
run_tests() {
    local test_type=$1
    case $test_type in
        "smoke")
            echo "💨 Running smoke tests..."
            npx playwright test basic-smoke.spec.ts --headed
            ;;
        "debug")
            echo "🐛 Running tests in debug mode..."
            npx playwright test basic-smoke.spec.ts --debug
            ;;
        "ui")
            echo "🖥️  Running tests in UI mode..."
            npx playwright test --ui
            ;;
        "all")
            echo "🧪 Running all tests..."
            npx playwright test
            ;;
        *)
            echo "🧪 Running basic smoke tests..."
            npx playwright test basic-smoke.spec.ts
            ;;
    esac
}

# Check command line arguments
if [ $# -eq 0 ]; then
    echo "Usage: $0 [smoke|debug|ui|all]"
    echo "  smoke - Run basic smoke tests"
    echo "  debug - Run tests in debug mode"
    echo "  ui    - Run tests in UI mode"
    echo "  all   - Run all tests"
    echo ""
    echo "Running default smoke tests..."
    run_tests "smoke"
else
    run_tests "$1"
fi

echo "✅ Test run completed!"