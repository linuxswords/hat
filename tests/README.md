# HAT Frontend Testing with Playwright

This directory contains end-to-end (E2E) frontend tests for the HAT (Handicap Archery Tournament) application using [Playwright](https://playwright.dev/).

## Overview

Our Playwright test suite provides comprehensive coverage of:
- **User Interface Testing**: Verifying all UI components render correctly
- **User Journey Testing**: Complete workflows from archer creation to tournament management
- **Cross-Browser Testing**: Chrome, Firefox, Safari, and mobile browsers
- **Responsive Design Testing**: Desktop, tablet, and mobile viewports
- **Accessibility Testing**: Keyboard navigation and screen reader compatibility

## Test Structure

```
tests/e2e/
├── pages/                  # Page Object Models
│   ├── BasePage.ts         # Common navigation and layout
│   ├── HomePage.ts         # Dashboard/home page
│   ├── ArchersPage.ts      # Archer management pages
│   └── TournamentsPage.ts  # Tournament management pages
├── setup/                  # Test configuration and utilities
│   └── test-setup.ts       # Extended test fixtures and utilities
├── home.spec.ts           # Home page and dashboard tests
├── archers.spec.ts        # Archer CRUD operations tests
├── tournaments.spec.ts    # Tournament management tests
├── navigation.spec.ts     # Navigation and layout tests
├── user-journey.spec.ts   # Complete user workflow tests
└── README.md              # This file
```

## Getting Started

### Prerequisites

- Node.js 18 or higher
- Go 1.24+ (for backend)
- HAT application dependencies

### Installation

1. Install dependencies:
   ```bash
   npm install
   ```

2. Install Playwright browsers:
   ```bash
   npm run playwright:install
   ```

### Running Tests

#### Local Development

```bash
# Run all E2E tests
npm run test:e2e

# Run tests with UI mode (interactive)
npm run test:e2e:ui

# Run tests in headed mode (see browser)
npm run test:e2e:headed

# Debug specific test
npm run test:e2e:debug -- home.spec.ts

# Run specific test file
npx playwright test archers.spec.ts

# Run tests on specific browser
npx playwright test --project=chromium
```

#### Test Reports

After running tests, view the HTML report:
```bash
npx playwright show-report
```

## Test Configuration

The test configuration is defined in `playwright.config.ts`:

- **Base URL**: `http://localhost:8080` (configurable via `BASE_URL` env var)
- **Browsers**: Chrome, Firefox, Safari, Mobile Chrome, Mobile Safari
- **Retries**: 2 retries on CI, 0 locally
- **Screenshots**: On failure only
- **Videos**: On failure only
- **Traces**: On first retry

## Page Object Models

We use the Page Object Model pattern to maintain clean, reusable test code:

### BasePage
Common functionality shared across all pages:
- Navigation elements (menu, logo, links)
- Common methods (waitForPageLoad, verifyNavigation)
- Browser utilities

### HomePage
Dashboard-specific functionality:
- Quick action cards
- Navigation to different sections
- Recent activity display

### ArchersPage
Archer management functionality:
- CRUD operations (Create, Read, Update, Delete)
- Form validation
- List management

### TournamentsPage
Tournament management functionality:
- Tournament lifecycle management
- Archer assignment to tournaments
- Scoring and rankings access

## Test Categories

### 1. Component Tests
- Individual page element verification
- Form field validation
- Button and link functionality

### 2. Integration Tests
- Multi-page workflows
- Data consistency across navigation
- Form submission and redirect handling

### 3. User Journey Tests
- Complete tournament setup workflow
- Archer management lifecycle
- Error recovery scenarios

### 4. Cross-Browser Tests
- Browser compatibility verification
- Mobile responsive design
- Touch interaction support

### 5. Accessibility Tests
- Keyboard navigation
- Screen reader compatibility
- Focus management

## Writing New Tests

### Basic Test Structure

```typescript
import { test, expect } from './setup/test-setup';
import { HomePage } from './pages/HomePage';

test.describe('Feature Name', () => {
  let homePage: HomePage;

  test.beforeEach(async ({ page }) => {
    homePage = new HomePage(page);
    await homePage.visit();
  });

  test('should do something', async ({ page }) => {
    // Test implementation
    await expect(page).toHaveTitle(/Expected Title/);
  });
});
```

### Best Practices

1. **Use Page Object Models**: Keep selectors and page logic in page classes
2. **Wait for Elements**: Use proper waits instead of fixed timeouts
3. **Descriptive Test Names**: Test names should describe behavior, not implementation
4. **Independent Tests**: Each test should be able to run independently
5. **Clean Up**: Clean up test data when possible
6. **Error Handling**: Test both success and error scenarios

### Locator Strategies

Prefer these locator strategies in order:
1. **Semantic selectors**: `page.getByRole('button', { name: 'Submit' })`
2. **Test IDs**: `page.getByTestId('submit-button')`
3. **Text content**: `page.getByText('Submit')`
4. **CSS selectors**: `page.locator('.submit-btn')` (last resort)

## Environment Configuration

### Local Development
- Application runs on `http://localhost:8080`
- Uses test data from JSON files
- No external dependencies required

### CI/CD
- GitHub Actions workflow in `.github/workflows/e2d-frontend-tests.yml`
- Automatic browser installation
- Parallel test execution
- Artifact upload for test reports

### Environment Variables

```bash
# Override base URL for testing
BASE_URL=http://localhost:3000 npm run test:e2e

# Run in CI mode
CI=true npm run test:e2e

# Debug mode
DEBUG=pw:api npm run test:e2e
```

## Troubleshooting

### Common Issues

1. **Test Timeout**: Increase timeout in `playwright.config.ts` or use `test.setTimeout()`
2. **Element Not Found**: Use proper waits: `waitFor({ state: 'visible' })`
3. **Flaky Tests**: Add proper waits and avoid fixed timeouts
4. **Browser Launch Issues**: Run `npx playwright install` to reinstall browsers
5. **Selector Mismatches**: Use the browser's developer tools to inspect actual HTML and update selectors
6. **Application Not Running**: Ensure the Go backend is running on port 8080

### Debug Mode

```bash
# Quick script for testing
./scripts/test-playwright.sh debug

# Run specific test in debug mode
npx playwright test --debug basic-smoke.spec.ts

# Run with headed browser for visual debugging
npx playwright test --headed basic-smoke.spec.ts

# Run just smoke tests
npx playwright test basic-smoke.spec.ts
```

### Getting Started - Quick Test

Start with the basic smoke tests to verify everything is working:

```bash
# Make sure application is built and running
make build
make run

# In another terminal, run basic tests
./scripts/test-playwright.sh smoke
```

### Selector Debugging

If tests fail due to element not found:

1. Run tests in headed mode to see the actual page
2. Use browser DevTools to inspect elements
3. Update selectors in page object models
4. Start with basic-smoke.spec.ts which has more robust selectors

### Test Data Management

- Tests create their own test data
- Use unique identifiers to avoid conflicts
- Clean up test data when possible
- Independent test execution

## Contributing

When adding new tests:

1. **Follow the existing patterns** in page objects and test structure
2. **Add new page objects** for new pages or major UI sections
3. **Update this README** if adding new test categories or patterns
4. **Test locally** across different browsers before submitting
5. **Include both positive and negative test cases**

## Integration with Backend

The E2E tests are designed to work with the Go backend:
- Tests use the actual HTTP endpoints
- Data persistence through JSON files
- Real template rendering and form handling
- Authentic user experience validation

This provides confidence that the entire application stack works correctly together.