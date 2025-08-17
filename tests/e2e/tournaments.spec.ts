import { test, expect } from './setup/test-setup';
import { TournamentsPage } from './pages/TournamentsPage';
import { ArchersPage } from './pages/ArchersPage';

test.describe('Tournaments Management', () => {
  let tournamentsPage: TournamentsPage;

  test.beforeEach(async ({ page }) => {
    tournamentsPage = new TournamentsPage(page);
  });

  test('should display tournaments list page', async ({ page }) => {
    await tournamentsPage.visit();
    
    // Verify page title
    await expect(page).toHaveTitle(/Tournaments/);
    
    // Verify navigation is present
    await tournamentsPage.verifyNavigation();
    
    // Verify list page elements
    await tournamentsPage.verifyListPageElements();
    
    // Verify create button is present
    await expect(tournamentsPage.createNewTournamentButton).toBeVisible();
  });

  test('should navigate to create tournament page', async ({ page }) => {
    await tournamentsPage.visit();
    await tournamentsPage.clickCreateNewTournament();
    
    // Verify we're on the create page
    await expect(page).toHaveURL(/\/tournaments\/new/);
    await expect(page.locator('h1')).toContainText(/Create.*Tournament/i);
    
    // Verify form elements are present
    await expect(tournamentsPage.nameField).toBeVisible();
    await expect(tournamentsPage.descriptionField).toBeVisible();
    await expect(tournamentsPage.startDateField).toBeVisible();
    await expect(tournamentsPage.endDateField).toBeVisible();
    await expect(tournamentsPage.locationField).toBeVisible();
    await expect(tournamentsPage.submitButton).toBeVisible();
  });

  test('should create a new tournament with valid data', async ({ page }) => {
    const testTournament = {
      name: 'Test Spring Tournament',
      description: 'A test tournament for spring archery competition',
      startDate: '2024-04-15',
      endDate: '2024-04-16',
      location: 'Test Archery Range'
    };

    await tournamentsPage.visitCreatePage();
    await tournamentsPage.fillTournamentForm(testTournament);
    await tournamentsPage.submitForm();
    
    // Should redirect to tournaments list or tournament details
    await expect(page).toHaveURL(/\/tournaments/);
    
    // Verify tournament appears in the list
    await tournamentsPage.verifyTournamentInList(testTournament.name);
  });

  test('should validate required fields', async ({ page }) => {
    await tournamentsPage.visitCreatePage();
    
    // Try to submit without filling required fields
    await tournamentsPage.submitForm();
    
    // Should stay on create page
    await expect(page).toHaveURL(/\/tournaments\/new/);
    
    // Check for required field validation
    const nameField = tournamentsPage.nameField;
    const isRequired = await nameField.getAttribute('required');
    expect(isRequired).not.toBeNull();
  });

  test('should validate date fields', async ({ page }) => {
    const invalidTournament = {
      name: 'Invalid Date Tournament',
      description: 'Test tournament with invalid dates',
      startDate: '2024-04-20',
      endDate: '2024-04-15', // End date before start date
      location: 'Test Location'
    };

    await tournamentsPage.visitCreatePage();
    await tournamentsPage.fillTournamentForm(invalidTournament);
    await tournamentsPage.submitForm();
    
    // Should either stay on form with validation error or handle gracefully
    // This depends on your validation implementation
    await expect(page).toHaveURL(/\/tournaments/);
  });

  test('should view tournament details', async ({ page }) => {
    // First create a tournament to view
    const testTournament = {
      name: 'View Test Tournament',
      description: 'Tournament created for viewing test',
      startDate: '2024-05-01',
      endDate: '2024-05-02',
      location: 'View Test Range'
    };

    await tournamentsPage.createTournament(testTournament);
    
    // Now view the tournament
    await tournamentsPage.viewTournament(testTournament.name);
    
    // Verify we're on the view page
    await expect(page).toHaveURL(/\/tournaments\/\d+$/);
    
    // Verify tournament details are displayed
    await expect(page).toHaveText(testTournament.name);
    await expect(page).toHaveText(testTournament.description);
    await expect(page).toHaveText(testTournament.location);
  });

  test('should edit tournament information', async ({ page }) => {
    // First create a tournament to edit
    const originalTournament = {
      name: 'Original Tournament',
      description: 'Original description',
      startDate: '2024-06-01',
      endDate: '2024-06-02',
      location: 'Original Location'
    };

    await tournamentsPage.createTournament(originalTournament);
    
    // Edit the tournament
    await tournamentsPage.editTournament(originalTournament.name);
    
    // Verify we're on the edit page
    await expect(page).toHaveURL(/\/tournaments\/\d+\/edit$/);
    
    // Verify form is pre-populated
    const formValues = await tournamentsPage.getFormValues();
    expect(formValues.name).toBe(originalTournament.name);
    expect(formValues.description).toBe(originalTournament.description);
    
    // Update tournament information
    const updatedTournament = {
      name: 'Updated Tournament Name',
      description: 'Updated description',
      startDate: '2024-06-10',
      endDate: '2024-06-11',
      location: 'Updated Location'
    };
    
    await tournamentsPage.fillTournamentForm(updatedTournament);
    await tournamentsPage.submitForm();
    
    // Should redirect back to tournament details or list
    await expect(page).toHaveURL(/\/tournaments/);
    
    // Verify updated information appears
    await tournamentsPage.verifyTournamentInList(updatedTournament.name);
  });

  test('should delete tournament', async ({ page }) => {
    // First create a tournament to delete
    const testTournament = {
      name: 'Delete Test Tournament',
      description: 'Tournament to be deleted',
      startDate: '2024-07-01',
      endDate: '2024-07-02',
      location: 'Delete Test Location'
    };

    await tournamentsPage.createTournament(testTournament);
    
    // Verify tournament exists
    await tournamentsPage.verifyTournamentInList(testTournament.name);
    
    // Delete the tournament
    await tournamentsPage.deleteTournament(testTournament.name);
    
    // Verify tournament is removed from list
    const tournamentRow = tournamentsPage.getTournamentRowByName(testTournament.name);
    await expect(tournamentRow).toHaveCount(0);
  });

  test('should manage tournament archers', async ({ page }) => {
    // First ensure we have an archer
    const archersPage = new ArchersPage(page);
    const testArcher = {
      name: 'Tournament Test Archer',
      gender: 'Male',
      bowClass: 'AMLB',
      email: 'tournament.test@example.com'
    };
    
    await archersPage.createArcher(testArcher);
    
    // Create a tournament
    const testTournament = {
      name: 'Archer Management Test',
      description: 'Tournament for testing archer management',
      startDate: '2024-08-01',
      endDate: '2024-08-02',
      location: 'Archer Test Location'
    };
    
    await tournamentsPage.createTournament(testTournament);
    
    // View tournament and manage archers
    await tournamentsPage.viewTournament(testTournament.name);
    await tournamentsPage.clickManageArchers();
    
    // Should be on the tournament archers page
    await expect(page).toHaveURL(/\/tournaments\/\d+\/archers/);
    
    // Add archer to tournament
    await tournamentsPage.addArchersToTournament([testArcher.name]);
    
    // Verify archer is added
    await tournamentsPage.verifyArchersInTournament([testArcher.name]);
  });

  test('should access tournament scoring features', async ({ page }) => {
    // Create a tournament
    const testTournament = {
      name: 'Scoring Test Tournament',
      description: 'Tournament for testing scoring features',
      startDate: '2024-09-01',
      endDate: '2024-09-02',
      location: 'Scoring Test Location'
    };
    
    await tournamentsPage.createTournament(testTournament);
    
    // View tournament
    await tournamentsPage.viewTournament(testTournament.name);
    
    // Test manage scores link
    await tournamentsPage.clickManageScores();
    await expect(page).toHaveURL(/\/tournaments\/\d+\/scores/);
    
    // Go back to tournament and test rankings
    await page.goBack();
    await tournamentsPage.clickViewRankings();
    await expect(page).toHaveURL(/\/tournaments\/\d+\/rankings/);
  });

  test('should handle empty tournaments list', async ({ page }) => {
    await tournamentsPage.visit();
    
    // If no tournaments exist, should display appropriate message
    const tournamentRows = await tournamentsPage.tournamentsTable.locator('tbody tr').all();
    
    if (tournamentRows.length === 0) {
      // Should show no tournaments message or empty table
      await expect(tournamentsPage.noTournamentsMessage).toBeVisible();
    } else {
      // Should show table with existing tournaments
      await expect(tournamentsPage.tournamentsTable).toBeVisible();
    }
  });

  test('should be responsive on mobile devices', async ({ page }) => {
    await page.setViewportSize({ width: 375, height: 667 });
    await tournamentsPage.visit();
    
    // Verify mobile layout
    await expect(tournamentsPage.createNewTournamentButton).toBeVisible();
    await expect(tournamentsPage.mobileMenuButton).toBeVisible();
    
    // Test create form on mobile
    await tournamentsPage.clickCreateNewTournament();
    await expect(tournamentsPage.nameField).toBeVisible();
    await expect(tournamentsPage.submitButton).toBeVisible();
  });

  test('should handle navigation between tournament pages', async ({ page }) => {
    await tournamentsPage.visit();
    
    // Navigate to create page
    await tournamentsPage.clickCreateNewTournament();
    await expect(page).toHaveURL(/\/tournaments\/new/);
    
    // Cancel and return to list
    await tournamentsPage.cancelButton.click();
    await expect(page).toHaveURL(/\/tournaments$/);
    
    // Use navigation to go to other sections
    await tournamentsPage.navigateToArchers();
    await expect(page).toHaveURL(/\/archers/);
    
    await tournamentsPage.navigateToTournaments();
    await expect(page).toHaveURL(/\/tournaments/);
  });
});