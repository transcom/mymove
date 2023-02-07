/**
 * base test fixture for playwright
 * Using javascript classes following the examples in
 * <https://playwright.dev/docs/test-fixtures>
 */
// @ts-check

const path = require('path');

const { expect } = require('@playwright/test');

const { TestHarness } = require('./testharness');

/**
 * base test fixture for playwright
 * See https://playwright.dev/docs/test-fixtures
 */
export class BaseTestPage {
  /**
   * Create an BaseTestPage.
   * @param {import('@playwright/test').Page} page
   * @param {import('@playwright/test').APIRequestContext} request
   */
  constructor(page, request) {
    this.page = page;
    this.request = request;
    this.testHarness = new TestHarness(request);
  }

  /**
   * @param {string} dutyLocationName
   * @param {string} fieldName
   * @param {number} nth
   */
  async selectDutyLocation(dutyLocationName, fieldName, nth = 0) {
    // fieldName is passed as a classname to the react-select component,
    // so select for it if provided
    const classSelector = '.duty-input-box';
    const actualClassSelector = fieldName ? `${classSelector}.${fieldName}` : classSelector;
    await this.page.locator(`${actualClassSelector} input[type="text"]`).type(dutyLocationName);

    // Click on the first presented option
    await this.page.locator(classSelector).locator('div[class*="option"]').nth(nth).click();
  }

  /**
   * Sign in as a new user with devlocal
   *
   * @param {string} userType
   */
  async signInAsNewUser(userType) {
    await this.page.goto('/devlocal-auth/login');
    await this.page.locator(`button[data-hook="new-user-login-${userType}"]`).click();
  }

  /**
   * Sign in as existing user with devlocal
   *
   * @param {string} userId
   */
  async signInAsUserWithId(userId) {
    await this.page.goto('/devlocal-auth/login');
    await this.page.locator(`button[value="${userId}"]`).click();
  }

  /**
   * Use fileChooser to upload files
   *
   * @param {import('@playwright/test').Locator} locator
   * @param {string} relativeFilePath path relative to playwright/fixtures
   */
  async uploadFileViaFilepond(locator, relativeFilePath) {
    const filePath = path.join('playwright/fixtures', relativeFilePath);
    const chooser = locator.getByText('choose from folder');
    await expect(chooser).toBeVisible();
    const fileChooserPromise = this.page.waitForEvent('filechooser');
    await chooser.click();
    const fileChooser = await fileChooserPromise;
    await fileChooser.setFiles(filePath);
  }
}

export default BaseTestPage;
