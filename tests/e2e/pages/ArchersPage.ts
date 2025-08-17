import { Page, Locator } from '@playwright/test';
import { BasePage } from './BasePage';

/**
 * Archers Page Object Model
 * 
 * Represents the archers management pages including list, create, edit, and view.
 * Covers the full CRUD operations for archer management.
 */
export class ArchersPage extends BasePage {
  // List page elements
  readonly pageTitle: Locator;
  readonly createNewArcherButton: Locator;
  readonly addNewArcherButton: Locator;
  readonly archersTable: Locator;
  readonly noArchersMessage: Locator;
  
  // Form elements (for create/edit)
  readonly nameField: Locator;
  readonly genderSelect: Locator;
  readonly bowClassSelect: Locator;
  readonly emailField: Locator;
  readonly submitButton: Locator;
  readonly cancelButton: Locator;
  
  // Individual archer actions
  readonly viewButton: Locator;
  readonly editButton: Locator;
  readonly deleteButton: Locator;
  
  // Success/Error messages
  readonly successMessage: Locator;
  readonly errorMessage: Locator;

  constructor(page: Page) {
    super(page);
    
    // List page elements
    this.pageTitle = page.locator('h1');
    this.createNewArcherButton = page.locator('a', { hasText: 'Create New Archer' });
    this.addNewArcherButton = page.locator('a[href="/archers/new"]');
    this.archersTable = page.locator('table');
    this.noArchersMessage = page.locator('text=No archers registered yet');
    
    // Form elements (using IDs for more precise targeting)
    this.nameField = page.locator('#name');
    this.genderSelect = page.locator('#gender');
    this.bowClassSelect = page.locator('#bow_class');
    this.emailField = page.locator('#email');
    this.submitButton = page.locator('button[type="submit"]');
    this.cancelButton = page.locator('a[href="/archers"]');
    
    // Action buttons (these will be context-specific)
    this.viewButton = page.locator('a', { hasText: 'View' });
    this.editButton = page.locator('a', { hasText: 'Edit' });
    this.deleteButton = page.locator('button', { hasText: 'Delete' });
    
    // Messages
    this.successMessage = page.locator('.alert-success, .success, [class*="success"]');
    this.errorMessage = page.locator('.alert-error, .error, [class*="error"]');
  }

  /**
   * Navigate to the archers list page
   */
  async visit() {
    await this.page.goto('/archers');
    await this.waitForPageLoad();
  }

  /**
   * Navigate to the create new archer page
   */
  async visitCreatePage() {
    await this.page.goto('/archers/new');
    await this.waitForPageLoad();
  }

  /**
   * Navigate to edit archer page
   */
  async visitEditPage(archerId: string) {
    await this.page.goto(`/archers/${archerId}/edit`);
    await this.waitForPageLoad();
  }

  /**
   * Navigate to view archer page
   */
  async visitViewPage(archerId: string) {
    await this.page.goto(`/archers/${archerId}`);
    await this.waitForPageLoad();
  }

  /**
   * Verify the archers list page elements
   */
  async verifyListPageElements() {
    await this.pageTitle.waitFor({ state: 'visible' });
    await this.addNewArcherButton.waitFor({ state: 'visible' });
  }

  /**
   * Click the create new archer button
   */
  async clickCreateNewArcher() {
    await this.addNewArcherButton.click();
    await this.waitForPageLoad();
  }

  /**
   * Fill in the archer form
   */
  async fillArcherForm(archerData: {
    name: string;
    gender: string;
    bowClass: string;
    email: string;
  }) {
    await this.nameField.fill(archerData.name);
    await this.genderSelect.selectOption(archerData.gender);
    await this.bowClassSelect.selectOption(archerData.bowClass);
    await this.emailField.fill(archerData.email);
  }

  /**
   * Submit the archer form
   */
  async submitForm() {
    await this.submitButton.click();
    await this.waitForPageLoad();
  }

  /**
   * Create a new archer with provided data
   */
  async createArcher(archerData: {
    name: string;
    gender: string;
    bowClass: string;
    email: string;
  }) {
    await this.visitCreatePage();
    await this.fillArcherForm(archerData);
    await this.submitForm();
  }

  /**
   * Get all archer rows from the table
   */
  async getArcherRows() {
    return await this.archersTable.locator('tbody tr').all();
  }

  /**
   * Find an archer row by name
   */
  getArcherRowByName(name: string) {
    return this.archersTable.locator('tr', { hasText: name });
  }

  /**
   * Click view button for a specific archer
   */
  async viewArcher(archerName: string) {
    const row = this.getArcherRowByName(archerName);
    await row.locator('a', { hasText: 'View' }).click();
    await this.waitForPageLoad();
  }

  /**
   * Click edit button for a specific archer
   */
  async editArcher(archerName: string) {
    const row = this.getArcherRowByName(archerName);
    await row.locator('a', { hasText: 'Edit' }).click();
    await this.waitForPageLoad();
  }

  /**
   * Delete an archer (with confirmation)
   */
  async deleteArcher(archerName: string) {
    const row = this.getArcherRowByName(archerName);
    
    // Handle potential confirmation dialog
    this.page.on('dialog', async dialog => {
      await dialog.accept();
    });
    
    await row.locator('button', { hasText: 'Delete' }).click();
    await this.waitForPageLoad();
  }

  /**
   * Verify an archer appears in the list
   */
  async verifyArcherInList(archerName: string) {
    const row = this.getArcherRowByName(archerName);
    await row.waitFor({ state: 'visible' });
  }

  /**
   * Verify an archer does not appear in the list
   */
  async verifyArcherNotInList(archerName: string) {
    const row = this.getArcherRowByName(archerName);
    await row.waitFor({ state: 'hidden' });
  }

  /**
   * Get the values of form fields (useful for edit page testing)
   */
  async getFormValues() {
    return {
      name: await this.nameField.inputValue(),
      gender: await this.genderSelect.inputValue(),
      bowClass: await this.bowClassSelect.inputValue(),
      email: await this.emailField.inputValue()
    };
  }

  /**
   * Verify form validation errors
   */
  async verifyValidationError(message: string) {
    await this.errorMessage.waitFor({ state: 'visible' });
    await this.page.locator('text=' + message).waitFor({ state: 'visible' });
  }
}