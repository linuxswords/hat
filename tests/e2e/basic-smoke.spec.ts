import { test, expect } from './setup/test-setup';

test.describe('Basic Smoke Tests', () => {
  test('should load the home page', async ({ page }) => {
    await page.goto('/');
    
    // Basic page load verification
    await expect(page).toHaveTitle(/HAT/);
    
    // Check for basic navigation elements
    await expect(page.locator('nav')).toBeVisible();
    await expect(page.locator('img[alt="HAT Logo"]')).toBeVisible();
    
    // Check for main content
    await expect(page.locator('h1')).toBeVisible();
  });

  test('should navigate to archers page', async ({ page }) => {
    await page.goto('/');
    
    // Click on archers link in navigation
    await page.click('a[href="/archers"]');
    
    // Verify we're on the archers page
    await expect(page).toHaveURL(/\/archers/);
    await expect(page.locator('h1')).toBeVisible();
  });

  test('should navigate to tournaments page', async ({ page }) => {
    await page.goto('/');
    
    // Click on tournaments link in navigation
    await page.click('a[href="/tournaments"]');
    
    // Verify we're on the tournaments page
    await expect(page).toHaveURL(/\/tournaments/);
    await expect(page.locator('h1')).toBeVisible();
  });

  test('should navigate to archer creation form', async ({ page }) => {
    await page.goto('/archers');
    
    // Look for the "Add New Archer" button and click it
    const addButton = page.locator('a[href="/archers/new"]');
    await expect(addButton).toBeVisible();
    await addButton.click();
    
    // Verify we're on the new archer page
    await expect(page).toHaveURL(/\/archers\/new/);
    
    // Check for form elements
    await expect(page.locator('#name')).toBeVisible();
    await expect(page.locator('#gender')).toBeVisible();
    await expect(page.locator('#bow_class')).toBeVisible();
    await expect(page.locator('#email')).toBeVisible();
  });

  test('should fill out archer form with valid data', async ({ page }) => {
    await page.goto('/archers/new');
    
    // Fill out the form
    await page.fill('#name', 'Test Archer');
    await page.selectOption('#gender', 'Male');
    
    // Check if bow class options are available
    const bowClassOptions = await page.locator('#bow_class option').count();
    if (bowClassOptions > 1) {
      // Select the first available option (not the placeholder)
      await page.selectOption('#bow_class', { index: 1 });
    }
    
    await page.fill('#email', 'test@example.com');
    
    // Submit the form
    await page.click('button[type="submit"]');
    
    // Should redirect somewhere (either back to list or show an error)
    await page.waitForURL(url => url.pathname !== '/archers/new');
  });

  test('should navigate to tournament creation form', async ({ page }) => {
    await page.goto('/tournaments');
    
    // Look for the "Create New Tournament" button and click it
    const createButton = page.locator('a[href="/tournaments/new"]');
    await expect(createButton).toBeVisible();
    await createButton.click();
    
    // Verify we're on the new tournament page
    await expect(page).toHaveURL(/\/tournaments\/new/);
    
    // Check for basic form structure
    await expect(page.locator('form')).toBeVisible();
  });

  test('should navigate to bow classes page', async ({ page }) => {
    await page.goto('/');
    
    // Click on bow classes link
    await page.click('a[href="/bowclasses"]');
    
    // Verify we're on the bow classes page
    await expect(page).toHaveURL(/\/bowclasses/);
    await expect(page.locator('h1')).toBeVisible();
  });

  test('should navigate to handicaps page', async ({ page }) => {
    await page.goto('/');
    
    // Click on handicaps link
    await page.click('a[href="/handicaps"]');
    
    // Verify we're on the handicaps page
    await expect(page).toHaveURL(/\/handicaps/);
    await expect(page.locator('h1')).toBeVisible();
  });

  test('should handle non-existent pages gracefully', async ({ page }) => {
    const response = await page.goto('/nonexistent-page');
    
    // Should get a 404 or redirect to a valid page
    expect([200, 404]).toContain(response?.status());
  });

  test('should be responsive on mobile viewports', async ({ page }) => {
    // Set mobile viewport
    await page.setViewportSize({ width: 375, height: 667 });
    
    await page.goto('/');
    
    // Check that the page loads and basic elements are visible
    await expect(page.locator('nav')).toBeVisible();
    await expect(page.locator('img[alt="HAT Logo"]')).toBeVisible();
    
    // Check for mobile menu button
    const mobileButton = page.locator('.md\\:hidden button');
    await expect(mobileButton).toBeVisible();
  });
});