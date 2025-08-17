import { Page, Locator } from '@playwright/test';

/**
 * Base Page Object Model for HAT application
 * 
 * Contains common elements and methods shared across all pages
 * like navigation, header, and common UI patterns.
 */
export class BasePage {
  readonly page: Page;
  
  // Common navigation elements
  readonly logo: Locator;
  readonly homeLink: Locator;
  readonly tournamentsLink: Locator;
  readonly archersLink: Locator;
  readonly bowClassesLink: Locator;
  readonly handicapsLink: Locator;
  readonly reportsLink: Locator;
  readonly mobileMenuButton: Locator;

  constructor(page: Page) {
    this.page = page;
    
    // Navigation elements (more specific selectors)
    this.logo = page.locator('img[alt="HAT Logo"]');
    this.homeLink = page.locator('nav a[href="/"]');
    this.tournamentsLink = page.locator('nav a[href="/tournaments"]');
    this.archersLink = page.locator('nav a[href="/archers"]');
    this.bowClassesLink = page.locator('nav a[href="/bowclasses"]');
    this.handicapsLink = page.locator('nav a[href="/handicaps"]');
    this.reportsLink = page.locator('nav a[href="/reports"]');
    this.mobileMenuButton = page.locator('.md\\:hidden button');
  }

  /**
   * Navigate to the home page
   */
  async navigateToHome() {
    await this.homeLink.click();
    await this.page.waitForLoadState('networkidle');
  }

  /**
   * Navigate to tournaments page
   */
  async navigateToTournaments() {
    await this.tournamentsLink.click();
    await this.page.waitForLoadState('networkidle');
  }

  /**
   * Navigate to archers page
   */
  async navigateToArchers() {
    await this.archersLink.click();
    await this.page.waitForLoadState('networkidle');
  }

  /**
   * Navigate to bow classes page
   */
  async navigateToBowClasses() {
    await this.bowClassesLink.click();
    await this.page.waitForLoadState('networkidle');
  }

  /**
   * Navigate to handicaps page
   */
  async navigateToHandicaps() {
    await this.handicapsLink.click();
    await this.page.waitForLoadState('networkidle');
  }

  /**
   * Navigate to reports page
   */
  async navigateToReports() {
    await this.reportsLink.click();
    await this.page.waitForLoadState('networkidle');
  }

  /**
   * Verify the navigation bar is visible and functional
   */
  async verifyNavigation() {
    await this.logo.waitFor({ state: 'visible' });
    await this.homeLink.waitFor({ state: 'visible' });
    await this.tournamentsLink.waitFor({ state: 'visible' });
    await this.archersLink.waitFor({ state: 'visible' });
    await this.bowClassesLink.waitFor({ state: 'visible' });
    await this.handicapsLink.waitFor({ state: 'visible' });
  }

  /**
   * Get the current page title
   */
  async getPageTitle(): Promise<string> {
    return await this.page.title();
  }

  /**
   * Wait for the page to be fully loaded
   */
  async waitForPageLoad() {
    await this.page.waitForLoadState('networkidle');
    await this.page.waitForLoadState('domcontentloaded');
  }
}