import { test, expect } from './setup/test-setup';
import { ArchersPage } from './pages/ArchersPage';

test.describe('Archers Management', () => {
  let archersPage: ArchersPage;

  test.beforeEach(async ({ page }) => {
    archersPage = new ArchersPage(page);
  });

  test('should display archers list page', async ({ page }) => {
    await archersPage.visit();
    
    // Verify page title
    await expect(page).toHaveTitle(/Archers/);
    
    // Verify navigation is present
    await archersPage.verifyNavigation();
    
    // Verify list page elements
    await archersPage.verifyListPageElements();
    
    // Verify create button is present
    await expect(archersPage.addNewArcherButton).toBeVisible();
  });

  test('should navigate to create archer page', async ({ page }) => {
    await archersPage.visit();
    await archersPage.clickCreateNewArcher();
    
    // Verify we're on the create page
    await expect(page).toHaveURL(/\/archers\/new/);
    await expect(page.locator('h1')).toContainText(/Add New Archer|Create.*Archer/i);
    
    // Verify form elements are present
    await expect(archersPage.nameField).toBeVisible();
    await expect(archersPage.genderSelect).toBeVisible();
    await expect(archersPage.bowClassSelect).toBeVisible();
    await expect(archersPage.emailField).toBeVisible();
    await expect(archersPage.submitButton).toBeVisible();
  });

  test('should create a new archer with valid data', async ({ page }) => {
    const testArcher = {
      name: 'John Test Archer',
      gender: 'Male',
      bowClass: 'AMLB', // Adult Male Longbow
      email: 'john.test@example.com'
    };

    await archersPage.visitCreatePage();
    await archersPage.fillArcherForm(testArcher);
    await archersPage.submitForm();
    
    // Should redirect to archers list
    await expect(page).toHaveURL(/\/archers$/);
    
    // Verify archer appears in the list
    await archersPage.verifyArcherInList(testArcher.name);
  });

  test('should validate required fields', async ({ page }) => {
    await archersPage.visitCreatePage();
    
    // Try to submit without filling required fields
    await archersPage.submitForm();
    
    // Should stay on create page and show validation errors
    await expect(page).toHaveURL(/\/archers\/new/);
    
    // Check for validation error (exact message depends on implementation)
    // This might be a browser validation or server-side validation
    const nameField = archersPage.nameField;
    const isRequired = await nameField.getAttribute('required');
    expect(isRequired).not.toBeNull();
  });

  test('should handle invalid email format', async ({ page }) => {
    const invalidArcher = {
      name: 'Test Archer',
      gender: 'Male',
      bowClass: 'AMLB',
      email: 'invalid-email'
    };

    await archersPage.visitCreatePage();
    await archersPage.fillArcherForm(invalidArcher);
    await archersPage.submitForm();
    
    // Should either stay on form with validation error or handle gracefully
    // This depends on your validation implementation
    const emailField = archersPage.emailField;
    const emailType = await emailField.getAttribute('type');
    expect(emailType).toBe('email');
  });

  test('should view archer details', async ({ page }) => {
    // First create an archer to view
    const testArcher = {
      name: 'View Test Archer',
      gender: 'Female',
      bowClass: 'AFLB', // Adult Female Longbow
      email: 'view.test@example.com'
    };

    await archersPage.createArcher(testArcher);
    
    // Now view the archer
    await archersPage.viewArcher(testArcher.name);
    
    // Verify we're on the view page
    await expect(page).toHaveURL(/\/archers\/\d+$/);
    
    // Verify archer details are displayed
    await expect(page.locator('h1')).toContainText(/Archer Details/i);
    await expect(page).toHaveText(testArcher.name);
    await expect(page).toHaveText(testArcher.email);
  });

  test('should edit archer information', async ({ page }) => {
    // First create an archer to edit
    const originalArcher = {
      name: 'Original Archer',
      gender: 'Male',
      bowClass: 'AMLB',
      email: 'original@example.com'
    };

    await archersPage.createArcher(originalArcher);
    
    // Edit the archer
    await archersPage.editArcher(originalArcher.name);
    
    // Verify we're on the edit page
    await expect(page).toHaveURL(/\/archers\/\d+\/edit$/);
    
    // Verify form is pre-populated
    const formValues = await archersPage.getFormValues();
    expect(formValues.name).toBe(originalArcher.name);
    expect(formValues.email).toBe(originalArcher.email);
    
    // Update archer information
    const updatedArcher = {
      name: 'Updated Archer Name',
      gender: 'Male',
      bowClass: 'AMLB',
      email: 'updated@example.com'
    };
    
    await archersPage.fillArcherForm(updatedArcher);
    await archersPage.submitForm();
    
    // Should redirect back to archer details or list
    await expect(page).toHaveURL(/\/archers/);
    
    // Verify updated information appears
    await archersPage.verifyArcherInList(updatedArcher.name);
  });

  test('should delete archer', async ({ page }) => {
    // First create an archer to delete
    const testArcher = {
      name: 'Delete Test Archer',
      gender: 'Female',
      bowClass: 'AFLB',
      email: 'delete.test@example.com'
    };

    await archersPage.createArcher(testArcher);
    
    // Verify archer exists
    await archersPage.verifyArcherInList(testArcher.name);
    
    // Delete the archer
    await archersPage.deleteArcher(testArcher.name);
    
    // Verify archer is removed from list
    await archersPage.verifyArcherNotInList(testArcher.name);
  });

  test('should handle navigation between pages', async ({ page }) => {
    await archersPage.visit();
    
    // Navigate to create page
    await archersPage.clickCreateNewArcher();
    await expect(page).toHaveURL(/\/archers\/new/);
    
    // Cancel and return to list
    await archersPage.cancelButton.click();
    await expect(page).toHaveURL(/\/archers$/);
    
    // Use navigation to go to other sections
    await archersPage.navigateToTournaments();
    await expect(page).toHaveURL(/\/tournaments/);
    
    await archersPage.navigateToArchers();
    await expect(page).toHaveURL(/\/archers/);
  });

  test('should handle empty archers list', async ({ page }) => {
    await archersPage.visit();
    
    // If no archers exist, should display appropriate message
    const archerRows = await archersPage.getArcherRows();
    
    if (archerRows.length === 0) {
      // Should show no archers message or empty table
      await expect(archersPage.noArchersMessage).toBeVisible();
    } else {
      // Should show table with existing archers
      await expect(archersPage.archersTable).toBeVisible();
    }
  });

  test('should be responsive on mobile devices', async ({ page }) => {
    await page.setViewportSize({ width: 375, height: 667 });
    await archersPage.visit();
    
    // Verify mobile layout
    await expect(archersPage.createNewArcherButton).toBeVisible();
    await expect(archersPage.mobileMenuButton).toBeVisible();
    
    // Test create form on mobile
    await archersPage.clickCreateNewArcher();
    await expect(archersPage.nameField).toBeVisible();
    await expect(archersPage.submitButton).toBeVisible();
  });
});