import { test, expect } from './setup/test-setup';
import { HomePage } from './pages/HomePage';
import { ArchersPage } from './pages/ArchersPage';
import { TournamentsPage } from './pages/TournamentsPage';

test.describe('Complete User Journey', () => {
  let homePage: HomePage;
  let archersPage: ArchersPage;
  let tournamentsPage: TournamentsPage;

  test.beforeEach(async ({ page }) => {
    homePage = new HomePage(page);
    archersPage = new ArchersPage(page);
    tournamentsPage = new TournamentsPage(page);
  });

  test('should complete full tournament setup workflow', async ({ page }) => {
    /**
     * Complete user journey: Create archers, create tournament, add archers to tournament
     * This simulates a real tournament organizer workflow
     */
    
    // Start from dashboard
    await homePage.visit();
    await homePage.verifyDashboardElements();
    
    // Step 1: Create some archers first
    await homePage.clickManageArchers();
    
    const testArchers = [
      {
        name: 'Alice Johnson',
        gender: 'Female',
        bowClass: 'AFLB', // Adult Female Longbow
        email: 'alice.johnson@email.com'
      },
      {
        name: 'Bob Smith',
        gender: 'Male',
        bowClass: 'AMLB', // Adult Male Longbow
        email: 'bob.smith@email.com'
      },
      {
        name: 'Carol Williams',
        gender: 'Female',
        bowClass: 'AFTR', // Adult Female Traditional Recurve
        email: 'carol.williams@email.com'
      }
    ];

    // Create each archer
    for (const archer of testArchers) {
      await archersPage.clickCreateNewArcher();
      await archersPage.fillArcherForm(archer);
      await archersPage.submitForm();
      
      // Verify archer was created
      await archersPage.verifyArcherInList(archer.name);
    }

    // Step 2: Create a tournament
    await homePage.navigateToHome();
    await homePage.clickCreateTournament();
    
    const testTournament = {
      name: 'Spring Championship 2024',
      description: 'Annual spring archery championship featuring multiple bow classes',
      startDate: '2024-05-15',
      endDate: '2024-05-16',
      location: 'Central Archery Range'
    };

    await tournamentsPage.fillTournamentForm(testTournament);
    await tournamentsPage.submitForm();
    
    // Verify tournament was created
    await tournamentsPage.verifyTournamentInList(testTournament.name);

    // Step 3: View tournament details
    await tournamentsPage.viewTournament(testTournament.name);
    
    // Verify tournament details page
    await expect(page).toHaveText(testTournament.name);
    await expect(page).toHaveText(testTournament.description);

    // Step 4: Add archers to tournament
    await tournamentsPage.clickManageArchers();
    
    // Add all created archers to the tournament
    await tournamentsPage.addArchersToTournament(testArchers.map(a => a.name));
    
    // Verify archers are added to tournament
    await tournamentsPage.verifyArchersInTournament(testArchers.map(a => a.name));

    // Step 5: Navigate to scoring
    await tournamentsPage.clickManageScores();
    await expect(page).toHaveURL(/\/tournaments\/\d+\/scores/);

    // Step 6: View rankings
    await tournamentsPage.clickViewRankings();
    await expect(page).toHaveURL(/\/tournaments\/\d+\/rankings/);
    
    // Verify the complete workflow was successful
    await expect(page.locator('h1')).toContainText(/Rankings/i);
  });

  test('should handle archer management workflow', async ({ page }) => {
    /**
     * Archer management workflow: Create, view, edit, and manage archer lifecycle
     */
    
    await homePage.visit();
    await homePage.clickManageArchers();

    // Create archer
    const archer = {
      name: 'Test Workflow Archer',
      gender: 'Male',
      bowClass: 'AMTR', // Adult Male Traditional Recurve
      email: 'workflow.test@email.com'
    };

    await archersPage.createArcher(archer);

    // View archer details
    await archersPage.viewArcher(archer.name);
    await expect(page).toHaveText(archer.name);
    await expect(page).toHaveText(archer.email);

    // Edit archer
    await archersPage.editArcher(archer.name);
    
    const updatedArcher = {
      ...archer,
      name: 'Updated Workflow Archer',
      email: 'updated.workflow@email.com'
    };
    
    await archersPage.fillArcherForm(updatedArcher);
    await archersPage.submitForm();

    // Verify update
    await archersPage.verifyArcherInList(updatedArcher.name);

    // Clean up - delete archer
    await archersPage.deleteArcher(updatedArcher.name);
    await archersPage.verifyArcherNotInList(updatedArcher.name);
  });

  test('should handle tournament lifecycle workflow', async ({ page }) => {
    /**
     * Tournament lifecycle: Create, edit, view, and manage tournament
     */
    
    await homePage.visit();
    await homePage.clickCreateTournament();

    // Create tournament
    const tournament = {
      name: 'Test Lifecycle Tournament',
      description: 'Tournament for testing complete lifecycle',
      startDate: '2024-08-01',
      endDate: '2024-08-02',
      location: 'Test Lifecycle Range'
    };

    await tournamentsPage.fillTournamentForm(tournament);
    await tournamentsPage.submitForm();

    // View tournament
    await tournamentsPage.viewTournament(tournament.name);
    await expect(page).toHaveText(tournament.name);

    // Edit tournament
    await tournamentsPage.editTournament(tournament.name);
    
    const updatedTournament = {
      ...tournament,
      name: 'Updated Lifecycle Tournament',
      description: 'Updated description for lifecycle test'
    };
    
    await tournamentsPage.fillTournamentForm(updatedTournament);
    await tournamentsPage.submitForm();

    // Verify update
    await tournamentsPage.verifyTournamentInList(updatedTournament.name);

    // Clean up - delete tournament
    await tournamentsPage.deleteTournament(updatedTournament.name);
  });

  test('should handle error recovery scenarios', async ({ page }) => {
    /**
     * Test error handling and recovery in user workflows
     */
    
    await homePage.visit();

    // Test form validation recovery
    await homePage.clickCreateTournament();
    
    // Submit empty form to trigger validation
    await tournamentsPage.submitForm();
    
    // Should stay on form page
    await expect(page).toHaveURL(/\/tournaments\/new/);
    
    // Fill valid data and retry
    const validTournament = {
      name: 'Recovery Test Tournament',
      description: 'Testing error recovery',
      startDate: '2024-09-01',
      endDate: '2024-09-02',
      location: 'Recovery Test Location'
    };
    
    await tournamentsPage.fillTournamentForm(validTournament);
    await tournamentsPage.submitForm();
    
    // Should succeed this time
    await expect(page).toHaveURL(/\/tournaments/);
    await tournamentsPage.verifyTournamentInList(validTournament.name);
  });

  test('should maintain data consistency across navigation', async ({ page }) => {
    /**
     * Test that data remains consistent when navigating between pages
     */
    
    // Create an archer
    await homePage.visit();
    await homePage.clickManageArchers();
    
    const testArcher = {
      name: 'Consistency Test Archer',
      gender: 'Female',
      bowClass: 'AFBHR', // Adult Female Bowhunter Recurve
      email: 'consistency.test@email.com'
    };
    
    await archersPage.createArcher(testArcher);
    
    // Navigate away and back
    await homePage.navigateToTournaments();
    await homePage.navigateToArchers();
    
    // Verify archer still exists
    await archersPage.verifyArcherInList(testArcher.name);
    
    // Create tournament
    await homePage.navigateToTournaments();
    await tournamentsPage.clickCreateNewTournament();
    
    const testTournament = {
      name: 'Consistency Test Tournament',
      description: 'Testing data consistency',
      startDate: '2024-10-01',
      endDate: '2024-10-02',
      location: 'Consistency Test Location'
    };
    
    await tournamentsPage.fillTournamentForm(testTournament);
    await tournamentsPage.submitForm();
    
    // Navigate to archer management and back
    await homePage.navigateToArchers();
    await homePage.navigateToTournaments();
    
    // Verify tournament still exists
    await tournamentsPage.verifyTournamentInList(testTournament.name);
  });

  test('should support mobile user workflows', async ({ page }) => {
    /**
     * Test complete workflows on mobile devices
     */
    
    // Set mobile viewport
    await page.setViewportSize({ width: 375, height: 667 });
    
    await homePage.visit();
    
    // Verify mobile layout
    await expect(homePage.mobileMenuButton).toBeVisible();
    await homePage.verifyDashboardElements();
    
    // Create archer on mobile
    await homePage.clickManageArchers();
    await archersPage.clickCreateNewArcher();
    
    const mobileArcher = {
      name: 'Mobile Test Archer',
      gender: 'Male',
      bowClass: 'YAMBHR', // Young Adult Male Bowhunter Recurve
      email: 'mobile.test@email.com'
    };
    
    await archersPage.fillArcherForm(mobileArcher);
    await archersPage.submitForm();
    
    // Verify archer creation on mobile
    await archersPage.verifyArcherInList(mobileArcher.name);
    
    // Create tournament on mobile
    await homePage.navigateToTournaments();
    await tournamentsPage.clickCreateNewTournament();
    
    const mobileTournament = {
      name: 'Mobile Test Tournament',
      description: 'Tournament created on mobile',
      startDate: '2024-11-01',
      endDate: '2024-11-02',
      location: 'Mobile Test Location'
    };
    
    await tournamentsPage.fillTournamentForm(mobileTournament);
    await tournamentsPage.submitForm();
    
    // Verify tournament creation on mobile
    await tournamentsPage.verifyTournamentInList(mobileTournament.name);
  });
});