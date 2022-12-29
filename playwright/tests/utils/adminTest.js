// admin test fixture for playwright
// See https://playwright.dev/docs/test-fixtures
const base = require('@playwright/test');

const { signIntoAdminAsNewAdminUser } = require('./signIn');
const { buildDefaultAdminUser, buildDefaultMove, buildOfficeUserWithTOOAndTIO } = require('./testharness');

class AdminPage {
  /**
   * @param {import('@playwright/test').Page} page
   */
  constructor(page, request) {
    this.page = page;
    this.request = request;
  }

  /**
   */
  async waitForLoading() {
    await base.expect(this.page.locator('h2[data-name="loading-placeholder"]')).toHaveCount(0);
    await base.expect(this.page.locator('svg.MuiCircularProgress-svg')).toHaveCount(0);
  }

  async waitForAdminPageToLoad() {
    // ensure the admin page has fully loaded before moving on
    await base.expect(this.page.locator('a:has-text("Logout")')).toHaveCount(1, { timeout: 10000 });
    await this.waitForLoading();
    await base.expect(this.page.locator('.ReactTable').locator('.-loading.-active')).toHaveCount(0);
  }

  /**
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
}

exports.test = base.test.extend({
  adminPage: async ({ page, request }, use) => {
    const adminPage = new AdminPage(page, request);
    await use(adminPage);
  },
});

exports.expect = base.expect;
