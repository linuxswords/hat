import { Page, Locator } from '@playwright/test';
import { BasePage } from './BasePage';

/**
 * Tournaments Page Object Model
 * 
 * Represents the tournaments management pages including list, create, edit, and view.
 * Also covers tournament-specific features like archer management and scoring.
 */
export class TournamentsPage extends BasePage {
  // List page elements
  readonly pageTitle: Locator;
  readonly createNewTournamentButton: Locator;
  readonly tournamentsTable: Locator;
  readonly noTournamentsMessage: Locator;
  readonly currentTournamentsSection: Locator;
  readonly upcomingTournamentsSection: Locator;
  
  // Form elements (for create/edit)
  readonly nameField: Locator;
  readonly descriptionField: Locator;
  readonly startDateField: Locator;
  readonly endDateField: Locator;
  readonly locationField: Locator;
  readonly handicapSetSelect: Locator;
  readonly submitButton: Locator;
  readonly cancelButton: Locator;
  
  // Tournament details page
  readonly tournamentName: Locator;
  readonly tournamentDescription: Locator;
  readonly tournamentDates: Locator;
  readonly manageArchersButton: Locator;
  readonly manageScoresButton: Locator;
  readonly viewRankingsButton: Locator;
  
  // Tournament actions
  readonly viewButton: Locator;
  readonly editButton: Locator;
  readonly deleteButton: Locator;
  
  // Archer management (within tournament)
  readonly addArchersButton: Locator;
  readonly archersInTournament: Locator;
  readonly removeArcherButton: Locator;
  
  // Success/Error messages
  readonly successMessage: Locator;
  readonly errorMessage: Locator;

  constructor(page: Page) {
    super(page);
    
    // List page elements
    this.pageTitle = page.locator('h1');
    this.createNewTournamentButton = page.locator('a[href="/tournaments/new"]');
    this.tournamentsTable = page.locator('table');
    this.noTournamentsMessage = page.locator('text=No tournaments created yet');
    this.currentTournamentsSection = page.locator('h2', { hasText: 'Current Tournaments' });
    this.upcomingTournamentsSection = page.locator('h2', { hasText: 'Upcoming Tournaments' });
    
    // Form elements
    this.nameField = page.locator('input[name="name"]');
    this.descriptionField = page.locator('textarea[name="description"]');
    this.startDateField = page.locator('input[name="start_date"]');
    this.endDateField = page.locator('input[name="end_date"]');
    this.locationField = page.locator('input[name="location"]');
    this.handicapSetSelect = page.locator('select[name="handicap_set"]');
    this.submitButton = page.locator('button[type="submit"]');
    this.cancelButton = page.locator('a', { hasText: 'Cancel' });
    
    // Tournament details
    this.tournamentName = page.locator('h1, h2').first();
    this.tournamentDescription = page.locator('p', { hasText: /.+/ }).first();
    this.tournamentDates = page.locator('text=/\d{4}-\d{2}-\d{2}/');
    this.manageArchersButton = page.locator('a', { hasText: 'Manage Archers' });
    this.manageScoresButton = page.locator('a', { hasText: 'Manage Scores' });
    this.viewRankingsButton = page.locator('a', { hasText: 'View Rankings' });
    
    // Action buttons
    this.viewButton = page.locator('a', { hasText: 'View' });
    this.editButton = page.locator('a', { hasText: 'Edit' });
    this.deleteButton = page.locator('button', { hasText: 'Delete' });
    
    // Archer management
    this.addArchersButton = page.locator('a', { hasText: 'Add Archers' });
    this.archersInTournament = page.locator('table').locator('tbody tr');
    this.removeArcherButton = page.locator('button', { hasText: 'Remove' });
    
    // Messages
    this.successMessage = page.locator('.alert-success, .success, [class*="success"]');
    this.errorMessage = page.locator('.alert-error, .error, [class*="error"]');
  }

  /**
   * Navigate to the tournaments list page
   */
  async visit() {
    await this.page.goto('/tournaments');
    await this.waitForPageLoad();
  }

  /**
   * Navigate to the create new tournament page
   */
  async visitCreatePage() {
    await this.page.goto('/tournaments/new');
    await this.waitForPageLoad();
  }

  /**
   * Navigate to view tournament page
   */
  async visitViewPage(tournamentId: string) {
    await this.page.goto(`/tournaments/${tournamentId}`);
    await this.waitForPageLoad();
  }

  /**
   * Navigate to edit tournament page
   */
  async visitEditPage(tournamentId: string) {
    await this.page.goto(`/tournaments/${tournamentId}/edit`);
    await this.waitForPageLoad();
  }

  /**
   * Navigate to tournament archers page
   */
  async visitArchersPage(tournamentId: string) {
    await this.page.goto(`/tournaments/${tournamentId}/archers`);
    await this.waitForPageLoad();
  }

  /**
   * Navigate to tournament scores page
   */
  async visitScoresPage(tournamentId: string) {
    await this.page.goto(`/tournaments/${tournamentId}/scores`);
    await this.waitForPageLoad();
  }

  /**
   * Verify the tournaments list page elements
   */
  async verifyListPageElements() {
    await this.pageTitle.waitFor({ state: 'visible' });
    await this.createNewTournamentButton.waitFor({ state: 'visible' });
  }

  /**
   * Click the create new tournament button
   */
  async clickCreateNewTournament() {
    await this.createNewTournamentButton.click();
    await this.waitForPageLoad();
  }

  /**
   * Fill in the tournament form
   */
  async fillTournamentForm(tournamentData: {
    name: string;
    description: string;
    startDate: string;
    endDate: string;
    location: string;
    handicapSet?: string;
  }) {
    await this.nameField.fill(tournamentData.name);
    await this.descriptionField.fill(tournamentData.description);
    await this.startDateField.fill(tournamentData.startDate);
    await this.endDateField.fill(tournamentData.endDate);
    await this.locationField.fill(tournamentData.location);
    
    if (tournamentData.handicapSet) {
      await this.handicapSetSelect.selectOption(tournamentData.handicapSet);
    }
  }

  /**
   * Submit the tournament form
   */
  async submitForm() {
    await this.submitButton.click();
    await this.waitForPageLoad();
  }

  /**
   * Create a new tournament with provided data
   */
  async createTournament(tournamentData: {
    name: string;
    description: string;
    startDate: string;
    endDate: string;
    location: string;
    handicapSet?: string;
  }) {
    await this.visitCreatePage();
    await this.fillTournamentForm(tournamentData);
    await this.submitForm();
  }

  /**
   * Find a tournament row by name
   */
  getTournamentRowByName(name: string) {
    return this.tournamentsTable.locator('tr', { hasText: name });
  }

  /**
   * Click view button for a specific tournament
   */
  async viewTournament(tournamentName: string) {
    const row = this.getTournamentRowByName(tournamentName);
    await row.locator('a', { hasText: 'View' }).click();
    await this.waitForPageLoad();
  }

  /**
   * Click edit button for a specific tournament
   */
  async editTournament(tournamentName: string) {
    const row = this.getTournamentRowByName(tournamentName);
    await row.locator('a', { hasText: 'Edit' }).click();
    await this.waitForPageLoad();
  }

  /**
   * Delete a tournament (with confirmation)
   */
  async deleteTournament(tournamentName: string) {
    const row = this.getTournamentRowByName(tournamentName);
    
    // Handle potential confirmation dialog
    this.page.on('dialog', async dialog => {
      await dialog.accept();
    });
    
    await row.locator('button', { hasText: 'Delete' }).click();
    await this.waitForPageLoad();
  }

  /**
   * Verify a tournament appears in the list
   */
  async verifyTournamentInList(tournamentName: string) {
    const row = this.getTournamentRowByName(tournamentName);
    await row.waitFor({ state: 'visible' });
  }

  /**
   * Click manage archers for current tournament
   */
  async clickManageArchers() {
    await this.manageArchersButton.click();
    await this.waitForPageLoad();
  }

  /**
   * Click manage scores for current tournament
   */
  async clickManageScores() {
    await this.manageScoresButton.click();
    await this.waitForPageLoad();
  }

  /**
   * Click view rankings for current tournament
   */
  async clickViewRankings() {
    await this.viewRankingsButton.click();
    await this.waitForPageLoad();
  }

  /**
   * Get the values of form fields
   */
  async getFormValues() {
    return {
      name: await this.nameField.inputValue(),
      description: await this.descriptionField.inputValue(),
      startDate: await this.startDateField.inputValue(),
      endDate: await this.endDateField.inputValue(),
      location: await this.locationField.inputValue(),
      handicapSet: await this.handicapSetSelect.inputValue()
    };
  }

  /**
   * Add archers to tournament
   */
  async addArchersToTournament(archerNames: string[]) {
    await this.addArchersButton.click();
    await this.waitForPageLoad();
    
    // Select archers (implementation depends on the UI)
    for (const archerName of archerNames) {
      const checkbox = this.page.locator(`input[type="checkbox"][value*="${archerName}"]`);
      await checkbox.check();
    }
    
    await this.submitButton.click();
    await this.waitForPageLoad();
  }

  /**
   * Verify archers are listed in tournament
   */
  async verifyArchersInTournament(archerNames: string[]) {
    for (const archerName of archerNames) {
      await this.page.locator('text=' + archerName).waitFor({ state: 'visible' });
    }
  }
}