import { test as base, expect } from '@playwright/test';

/**
 * Extended test setup for HAT application
 * 
 * This setup provides common utilities and fixtures for testing
 * the archery tournament management system.
 */

// Custom test with additional fixtures
export const test = base.extend({
  // Add any custom fixtures here if needed in the future
});

// Re-export expect for convenience
export { expect };

// Common test utilities
export class TestUtils {
  /**
   * Navigate to a page and wait for it to load
   */
  static async navigateAndWaitForLoad(page: any, path: string) {
    await page.goto(path);
    await page.waitForLoadState('networkidle');
  }

  /**
   * Fill a form field and wait for it to be updated
   */
  static async fillField(page: any, selector: string, value: string) {
    await page.fill(selector, value);
    await page.waitForTimeout(100); // Small delay for form updates
  }

  /**
   * Click a button and wait for navigation or response
   */
  static async clickAndWait(page: any, selector: string, waitForNavigation = true) {
    if (waitForNavigation) {
      await Promise.all([
        page.waitForNavigation(),
        page.click(selector)
      ]);
    } else {
      await page.click(selector);
    }
  }

  /**
   * Check if an element is visible and contains expected text
   */
  static async verifyElementText(page: any, selector: string, expectedText: string) {
    const element = page.locator(selector);
    await expect(element).toBeVisible();
    await expect(element).toContainText(expectedText);
  }
}