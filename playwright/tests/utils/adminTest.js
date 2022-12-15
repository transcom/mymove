const base = require('@playwright/test');

const { signIntoAdminAsNewAdminUser } = require('./signIn');
const { buildDefaultAdminUser, buildDefaultMove } = require('./testharness');

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

  /**
   */
  async signInAsNewAdminUser() {
    await signIntoAdminAsNewAdminUser(this.page);
  }

  /**
   */
  async buildDefaultAdminUser() {
    return buildDefaultAdminUser(this.request);
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
