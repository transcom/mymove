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
}

export default BaseTestPage;
