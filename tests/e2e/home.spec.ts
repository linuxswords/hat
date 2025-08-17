import { test, expect } from './setup/test-setup';
import { HomePage } from './pages/HomePage';

test.describe('Home Page', () => {
  let homePage: HomePage;

  test.beforeEach(async ({ page }) => {
    homePage = new HomePage(page);
    await homePage.visit();
  });

  test('should display the dashboard with all elements', async ({ page }) => {
    // Verify page title
    await expect(page).toHaveTitle(/HAT/);
    
    // Verify navigation is present
    await homePage.verifyNavigation();
    
    // Verify main dashboard elements exist (but be flexible about exact content)
    await expect(homePage.pageTitle).toBeVisible();
    
    // Check for the presence of quick action cards (be flexible about exact text)
    await expect(page.locator('a[href="/tournaments/new"]')).toBeVisible();
    await expect(page.locator('a[href="/archers"]')).toBeVisible();
    await expect(page.locator('a[href="/bowclasses"]')).toBeVisible();
    await expect(page.locator('a[href="/handicaps"]')).toBeVisible();
  });

  test('should navigate to create tournament page', async ({ page }) => {
    await homePage.clickCreateTournament();
    
    // Verify we're on the create tournament page
    await expect(page).toHaveURL(/\/tournaments\/new/);
    await expect(page.locator('h1')).toContainText(/Create.*Tournament/i);
  });

  test('should navigate to manage archers page', async ({ page }) => {
    await homePage.clickManageArchers();
    
    // Verify we're on the archers page
    await expect(page).toHaveURL(/\/archers/);
    await expect(page.locator('h1')).toContainText(/Archers/i);
  });

  test('should navigate to bow classes page', async ({ page }) => {
    await homePage.clickBowClasses();
    
    // Verify we're on the bow classes page
    await expect(page).toHaveURL(/\/bowclasses/);
    await expect(page.locator('h1')).toContainText(/Bow Classes/i);
  });

  test('should navigate to handicaps page', async ({ page }) => {
    await homePage.clickHandicapSets();
    
    // Verify we're on the handicaps page
    await expect(page).toHaveURL(/\/handicaps/);
    await expect(page.locator('h1')).toContainText(/Handicap/i);
  });

  test('should display recent activity section', async () => {
    await expect(homePage.recentActivitySection).toBeVisible();
    await expect(homePage.viewAllTournamentsLink).toBeVisible();
  });

  test('should have working navigation links', async ({ page }) => {
    // Test main navigation links
    await homePage.navigateToTournaments();
    await expect(page).toHaveURL(/\/tournaments/);
    
    await homePage.navigateToArchers();
    await expect(page).toHaveURL(/\/archers/);
    
    await homePage.navigateToBowClasses();
    await expect(page).toHaveURL(/\/bowclasses/);
    
    await homePage.navigateToHandicaps();
    await expect(page).toHaveURL(/\/handicaps/);
    
    // Navigate back to home
    await homePage.navigateToHome();
    await expect(page).toHaveURL('/');
  });

  test('should have responsive design elements', async ({ page }) => {
    // Test on mobile viewport
    await page.setViewportSize({ width: 375, height: 667 });
    
    // Verify mobile menu button is visible
    await expect(homePage.mobileMenuButton).toBeVisible();
    
    // Verify cards are still visible and clickable
    await expect(homePage.createTournamentCard).toBeVisible();
    await expect(homePage.manageArchersCard).toBeVisible();
    
    // Test tablet viewport
    await page.setViewportSize({ width: 768, height: 1024 });
    
    // Verify layout adjusts appropriately
    await expect(homePage.createTournamentCard).toBeVisible();
    await expect(homePage.manageArchersCard).toBeVisible();
  });

  test('should have accessible elements', async () => {
    // Verify important elements have proper alt text or labels
    await expect(homePage.logo).toHaveAttribute('alt', 'HAT Logo');
    
    // Verify buttons are keyboard accessible
    await homePage.createTournamentCard.focus();
    await expect(homePage.createTournamentCard).toBeFocused();
    
    // Test keyboard navigation
    await homePage.createTournamentCard.press('Tab');
    await expect(homePage.manageArchersCard).toBeFocused();
  });

  test('should handle hover effects on quick action cards', async () => {
    // Verify hover effects work
    await homePage.verifyCardHoverEffects();
    
    // Verify cards return to normal state after hover
    await homePage.createTournamentCard.hover();
    await page.mouse.move(0, 0); // Move mouse away
    
    // Cards should still be visible and functional
    await expect(homePage.createTournamentCard).toBeVisible();
  });
});