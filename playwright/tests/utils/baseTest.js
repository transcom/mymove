/**
 * base test fixture for playwright
 * Using javascript classes following the examples in
 * <https://playwright.dev/docs/test-fixtures>
 */
// @ts-check

const { newSignIn } = require('./signIn');
const { newTestHarness } = require('./testharness');

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
    this.signIn = newSignIn(page);
    this.testHarness = newTestHarness(request);
  }

  /**
   * @param {string} inputData
   * @param {string} fieldName
   * @param {string} classSelector
   * @param {number} nth
   */
  async genericSelect(inputData, fieldName, classSelector, nth) {
    // fieldName is passed as a classname to the react-select component,
    // so select for it if provided
    const actualClassSelector = fieldName ? `${classSelector}.${fieldName}` : classSelector;
    await this.page.locator(`${actualClassSelector} input[type="text"]`).type(inputData);

    // Click on the first presented option
    await this.page.locator(classSelector).locator('div[class*="option"]').nth(nth).click();
  }

  /**
   * @param {string} dutyLocationName
   * @param {string} fieldName
   * @param {number} nth
   */
  async selectDutyLocation(dutyLocationName, fieldName, nth = 0) {
    return this.genericSelect(dutyLocationName, fieldName, '.duty-input-box', nth);
  }
}

export default BaseTestPage;
