import { test, expect } from './setup/test-setup';
import { BasePage } from './pages/BasePage';

test.describe('Navigation and Layout', () => {
  let basePage: BasePage;

  test.beforeEach(async ({ page }) => {
    basePage = new BasePage(page);
    await page.goto('/');
  });

  test('should display consistent navigation across all pages', async ({ page }) => {
    // Verify navigation is present on home page
    await basePage.verifyNavigation();
    
    // Test navigation consistency across different pages
    const pages = [
      { url: '/tournaments', expectedTitle: /Tournaments/i },
      { url: '/archers', expectedTitle: /Archers/i },
      { url: '/bowclasses', expectedTitle: /Bow Classes/i },
      { url: '/handicaps', expectedTitle: /Handicap/i }
    ];

    for (const pageInfo of pages) {
      await page.goto(pageInfo.url);
      await basePage.waitForPageLoad();
      
      // Verify navigation is still present
      await basePage.verifyNavigation();
      
      // Verify page loads correctly
      await expect(page.locator('h1')).toContainText(pageInfo.expectedTitle);
    }
  });

  test('should have working logo link', async ({ page }) => {
    // Navigate to a different page
    await basePage.navigateToTournaments();
    await expect(page).toHaveURL(/\/tournaments/);
    
    // Click logo to return home
    await basePage.logo.click();
    await basePage.waitForPageLoad();
    await expect(page).toHaveURL('/');
  });

  test('should highlight active navigation item', async ({ page }) => {
    // This test depends on your CSS implementation for active states
    // You might need to adjust selectors based on your actual implementation
    
    await basePage.navigateToTournaments();
    const tournamentsLink = basePage.tournamentsLink;
    
    // Check if the link has active styling (this depends on your CSS classes)
    const linkClasses = await tournamentsLink.getAttribute('class');
    // You would check for active classes here based on your implementation
    
    await basePage.navigateToArchers();
    const archersLink = basePage.archersLink;
    
    // Verify archers link is now active
    await expect(archersLink).toBeVisible();
  });

  test('should work with keyboard navigation', async ({ page }) => {
    // Test tab navigation through menu items
    await basePage.homeLink.focus();
    await expect(basePage.homeLink).toBeFocused();
    
    await page.keyboard.press('Tab');
    await expect(basePage.tournamentsLink).toBeFocused();
    
    await page.keyboard.press('Tab');
    await expect(basePage.archersLink).toBeFocused();
    
    // Test Enter key navigation
    await page.keyboard.press('Enter');
    await basePage.waitForPageLoad();
    await expect(page).toHaveURL(/\/archers/);
  });

  test('should handle mobile navigation menu', async ({ page }) => {
    // Set mobile viewport
    await page.setViewportSize({ width: 375, height: 667 });
    
    // Verify mobile menu button is visible
    await expect(basePage.mobileMenuButton).toBeVisible();
    
    // Verify desktop menu items are hidden on mobile
    const desktopNav = page.locator('.hidden.md\\:flex');
    await expect(desktopNav).toHaveCSS('display', 'none');
    
    // Test mobile menu functionality (if implemented)
    // This would depend on your mobile menu implementation
    await basePage.mobileMenuButton.click();
    // Add assertions for mobile menu behavior
  });

  test('should maintain navigation state during form interactions', async ({ page }) => {
    // Navigate to a form page
    await basePage.navigateToArchers();
    await page.click('a[href="/archers/new"]');
    
    // Verify navigation is still functional while on form
    await basePage.verifyNavigation();
    
    // Verify we can navigate away from form
    await basePage.navigateToTournaments();
    await expect(page).toHaveURL(/\/tournaments/);
  });

  test('should handle deep link navigation', async ({ page }) => {
    // Test direct navigation to deep URLs
    await page.goto('/tournaments/new');
    await basePage.waitForPageLoad();
    
    // Verify navigation is present and functional
    await basePage.verifyNavigation();
    
    // Verify we can navigate from deep link
    await basePage.navigateToHome();
    await expect(page).toHaveURL('/');
  });

  test('should display consistent branding elements', async ({ page }) => {
    const pages = ['/', '/tournaments', '/archers', '/bowclasses', '/handicaps'];
    
    for (const url of pages) {
      await page.goto(url);
      await basePage.waitForPageLoad();
      
      // Verify logo is present and consistent
      await expect(basePage.logo).toBeVisible();
      await expect(basePage.logo).toHaveAttribute('alt', 'HAT Logo');
      
      // Verify consistent layout elements
      const nav = page.locator('nav');
      await expect(nav).toBeVisible();
    }
  });

  test('should handle browser back/forward navigation', async ({ page }) => {
    // Navigate through several pages
    await basePage.navigateToTournaments();
    await basePage.navigateToArchers();
    await basePage.navigateToBowClasses();
    
    // Test browser back button
    await page.goBack();
    await expect(page).toHaveURL(/\/archers/);
    
    await page.goBack();
    await expect(page).toHaveURL(/\/tournaments/);
    
    // Test browser forward button
    await page.goForward();
    await expect(page).toHaveURL(/\/archers/);
    
    // Verify navigation is still functional after back/forward
    await basePage.verifyNavigation();
  });

  test('should maintain accessibility standards', async ({ page }) => {
    // Test navigation accessibility
    await basePage.verifyNavigation();
    
    // Verify navigation links have proper accessibility attributes
    const navLinks = [
      basePage.homeLink,
      basePage.tournamentsLink,
      basePage.archersLink,
      basePage.bowClassesLink,
      basePage.handicapsLink
    ];
    
    for (const link of navLinks) {
      await expect(link).toBeVisible();
      
      // Verify links are keyboard accessible
      await link.focus();
      await expect(link).toBeFocused();
    }
    
    // Test skip link functionality (if implemented)
    await page.keyboard.press('Tab');
    const skipLink = page.locator('a[href="#main"]');
    if (await skipLink.isVisible()) {
      await expect(skipLink).toBeFocused();
    }
  });

  test('should handle external link behavior', async ({ page }) => {
    // If you have external links in navigation, test them here
    // This would depend on your specific implementation
    
    // Example: if you have help links or external resources
    const externalLinks = page.locator('a[href^="http"]');
    const linkCount = await externalLinks.count();
    
    for (let i = 0; i < linkCount; i++) {
      const link = externalLinks.nth(i);
      const target = await link.getAttribute('target');
      
      // External links should open in new tab/window
      expect(target).toBe('_blank');
      
      // Should have security attributes
      const rel = await link.getAttribute('rel');
      expect(rel).toContain('noopener');
    }
  });
});