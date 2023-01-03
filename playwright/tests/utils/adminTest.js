// admin test fixture for playwright
// See https://playwright.dev/docs/test-fixtures
// @ts-check
const base = require('@playwright/test');

const { signIntoAdminAsNewAdminUser } = require('./signIn');
const {
  buildDefaultAdminUser,
  buildDefaultMove,
  buildOfficeUserWithTOOAndTIO,
  buildWebhookSubscription,
} = require('./testharness');

class AdminPage {
  /**
   * Create an AdminPage.
   * @param {base.Page} page
   * @param {base.APIRequestContext} request
   */
  constructor(page, request) {
    this.page = page;
    this.request = request;
  }

  /**
   * Wait for the loading placeholder to go away
   */
  async waitForLoading() {
    await base.expect(this.page.locator('h2[data-name="loading-placeholder"]')).toHaveCount(0);
    await base.expect(this.page.locator('svg.MuiCircularProgress-svg')).toHaveCount(0);
  }

  /**
   * wait for the admin page to finish loading
   */

  async waitForAdminPageToLoad() {
    // ensure the admin page has fully loaded before moving on
    await base.expect(this.page.locator('a:has-text("Logout")')).toHaveCount(1, { timeout: 10000 });
    await this.waitForLoading();
    await base.expect(this.page.locator('.ReactTable').locator('.-loading.-active')).toHaveCount(0);
  }

  /**
   * Create a new admin user and sign in as them
   */
  async signInAsNewAdminUser() {
    await signIntoAdminAsNewAdminUser(this.page);
    await this.waitForAdminPageToLoad();
  }

  /**
   */
  async buildDefaultAdminUser() {
    return buildDefaultAdminUser(this.request);
  }

  /**
   */
  async buildOfficeUserWithTOOAndTIO() {
    return buildOfficeUserWithTOOAndTIO(this.request);
  }

  /**
   */
  async buildDefaultMove() {
    return buildDefaultMove(this.request);
  }

  /**
   */
  async buildWebhookSubscription() {
    return buildWebhookSubscription(this.request);
  }

  /**
   */
  async expectRoleLabelsByText(role, labels, options = { exact: true }) {
    for (const label of labels) {
      await base.expect(this.page.getByRole(role).getByText(label, options)).toBeVisible();
    }
  }

  /**
   */
  async expectLocatorLabelsByText(locator, labels, options = { exact: true }) {
    for (const label of labels) {
      await base.expect(this.page.locator(locator).getByText(label, options)).toBeVisible();
    }
  }
}

/**
 * @typedef {object} AdminPageTestArgs - admin page test args
 * @property {AdminPage} adminPage     - admin page
 */

/** @type {base.Fixtures<AdminPageTestArgs, {}, base.PlaywrightTestArgs, base.PlaywrightWorkerArgs>} */
const adminFixtures = {
  adminPage: async ({ page, request }, use) => {
    const adminPage = new AdminPage(page, request);
    await use(adminPage);
  },
};

exports.test = base.test.extend(adminFixtures);
exports.expect = base.expect;
