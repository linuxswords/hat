import { defineConfig, devices } from '@playwright/test';

/**
 * Playwright configuration for HAT (Handicap Archery Tournament) frontend testing
 * 
 * This configuration sets up cross-browser testing for the archery tournament
 * management application, covering all major user journeys and interactions.
 */
export default defineConfig({
  // Test directory
  testDir: './tests/e2e',
  
  // Test output directory
  outputDir: './test-results',
  
  // Run tests in files in parallel
  fullyParallel: true,
  
  // Fail the build on CI if you accidentally left test.only in the source code
  forbidOnly: !!process.env.CI,
  
  // Retry on CI only
  retries: process.env.CI ? 2 : 0,
  
  // Opt out of parallel tests on CI
  workers: process.env.CI ? 1 : undefined,
  
  // Reporter to use
  reporter: process.env.CI ? 'github' : [
    ['html', { outputFolder: 'playwright-report' }],
    ['json', { outputFile: 'test-results.json' }],
    ['list']
  ],
  
  // Shared settings for all the projects below
  use: {
    // Base URL for the application
    baseURL: process.env.BASE_URL || 'http://localhost:8080',
    
    // Collect trace when retrying the failed test
    trace: 'on-first-retry',
    
    // Record video only when retrying with failures
    video: 'retain-on-failure',
    
    // Take screenshot only when retrying with failures  
    screenshot: 'only-on-failure',
    
    // Global timeout for each action (increased for development)
    actionTimeout: 15000,
    
    // Global timeout for navigation (increased for development)
    navigationTimeout: 45000,
    
    // Ignore HTTPS errors for local development
    ignoreHTTPSErrors: true,
  },

  // Configure projects for major browsers
  projects: [
    {
      name: 'chromium',
      use: { ...devices['Desktop Chrome'] },
    },
    
    // Uncomment these when basic tests are working
    // {
    //   name: 'firefox',
    //   use: { ...devices['Desktop Firefox'] },
    // },
    
    // {
    //   name: 'webkit',
    //   use: { ...devices['Desktop Safari'] },
    // },
    
    // Mobile testing
    // {
    //   name: 'Mobile Chrome',
    //   use: { ...devices['Pixel 5'] },
    // },
    
    // {
    //   name: 'Mobile Safari',
    //   use: { ...devices['iPhone 12'] },
    // },
  ],

  // Configure the development server to auto-start
  webServer: {
    command: 'make run',
    url: 'http://localhost:8080',
    reuseExistingServer: !process.env.CI,
    timeout: 120 * 1000, // 2 minutes
    stdout: 'ignore',
    stderr: 'pipe',
  },
});