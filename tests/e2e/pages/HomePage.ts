import { Page, Locator } from '@playwright/test';
import { BasePage } from './BasePage';

/**
 * Home Page Object Model
 * 
 * Represents the dashboard/home page of the HAT application
 * with quick action cards and recent activity section.
 */
export class HomePage extends BasePage {
  // Main content elements
  readonly pageTitle: Locator;
  readonly subtitle: Locator;
  
  // Quick action cards
  readonly createTournamentCard: Locator;
  readonly manageArchersCard: Locator;
  readonly bowClassesCard: Locator;
  readonly handicapSetsCard: Locator;
  readonly viewReportsCard: Locator;
  
  // Quick action buttons
  readonly createTournamentButton: Locator;
  readonly manageArchersButton: Locator;
  readonly bowClassesButton: Locator;
  readonly handicapSetsButton: Locator;
  readonly viewReportsButton: Locator;
  
  // Recent activity section
  readonly recentActivitySection: Locator;
  readonly viewAllTournamentsLink: Locator;

  constructor(page: Page) {
    super(page);
    
    // Main content
    this.pageTitle = page.locator('h1').first();
    this.subtitle = page.locator('p', { hasText: 'Tournament Management Dashboard' });
    
    // Quick action cards (using href selectors for precision)
    this.createTournamentCard = page.locator('a[href="/tournaments/new"]');
    this.manageArchersCard = page.locator('a[href="/archers"]');
    this.bowClassesCard = page.locator('a[href="/bowclasses"]');
    this.handicapSetsCard = page.locator('a[href="/handicaps"]');
    this.viewReportsCard = page.locator('a[href="/reports"]');
    
    // Quick action buttons within cards
    this.createTournamentButton = this.createTournamentCard.locator('span', { hasText: 'Create' });
    this.manageArchersButton = this.manageArchersCard.locator('span', { hasText: 'Manage' });
    this.bowClassesButton = this.bowClassesCard.locator('span', { hasText: 'Configure' });
    this.handicapSetsButton = this.handicapSetsCard.locator('span', { hasText: 'Configure' });
    this.viewReportsButton = this.viewReportsCard.locator('span', { hasText: 'View' });
    
    // Recent activity
    this.recentActivitySection = page.locator('div', { hasText: 'Recent Activity' });
    this.viewAllTournamentsLink = page.locator('a', { hasText: 'View all tournaments' });
  }

  /**
   * Navigate to the home page
   */
  async visit() {
    await this.page.goto('/');
    await this.waitForPageLoad();
  }

  /**
   * Verify all dashboard elements are visible
   */
  async verifyDashboardElements() {
    await this.pageTitle.waitFor({ state: 'visible' });
    await this.subtitle.waitFor({ state: 'visible' });
    
    // Verify all quick action cards are visible
    await this.createTournamentCard.waitFor({ state: 'visible' });
    await this.manageArchersCard.waitFor({ state: 'visible' });
    await this.bowClassesCard.waitFor({ state: 'visible' });
    await this.handicapSetsCard.waitFor({ state: 'visible' });
    await this.viewReportsCard.waitFor({ state: 'visible' });
    
    // Verify recent activity section
    await this.recentActivitySection.waitFor({ state: 'visible' });
  }

  /**
   * Click on Create Tournament quick action
   */
  async clickCreateTournament() {
    await this.createTournamentCard.click();
    await this.waitForPageLoad();
  }

  /**
   * Click on Manage Archers quick action
   */
  async clickManageArchers() {
    await this.manageArchersCard.click();
    await this.waitForPageLoad();
  }

  /**
   * Click on Bow Classes quick action
   */
  async clickBowClasses() {
    await this.bowClassesCard.click();
    await this.waitForPageLoad();
  }

  /**
   * Click on Handicap Sets quick action
   */
  async clickHandicapSets() {
    await this.handicapSetsCard.click();
    await this.waitForPageLoad();
  }

  /**
   * Click on View Reports quick action
   */
  async clickViewReports() {
    await this.viewReportsCard.click();
    await this.waitForPageLoad();
  }

  /**
   * Verify the hover effects on quick action cards
   */
  async verifyCardHoverEffects() {
    // Hover over create tournament card and verify visual changes
    await this.createTournamentCard.hover();
    // Could add specific visual verification here if needed
    
    // Test other cards similarly
    await this.manageArchersCard.hover();
    await this.bowClassesCard.hover();
    await this.handicapSetsCard.hover();
    await this.viewReportsCard.hover();
  }
}