const base = require('@playwright/test');

const { buildPPMInProgressMove } = require('./testharness');

class OfficePage {
  /**
   * @param {import('@playwright/test').Page} page
   */
  constructor(page, request) {
    this.page = page;
    this.request = request;
  }

  async gotoAndWaitForLoading(url) {
    await this.page.goto(url);
    await base.expect(this.page.locator('h2[data-name="loading-placeholder"]')).toHaveCount(0);
  }

  async loginAsNewPPMOfficeUser() {
    await this.page.goto('/devlocal-auth/login');
    await this.page.locator('button[data-hook="new-user-login-PPM office"]').click();
  }

  async buildInProgressPPMMove() {
    return buildPPMInProgressMove(this.request);
  }
}

exports.test = base.test.extend({
  officePage: async ({ page, request }, use) => {
    const officePage = new OfficePage(page, request);
    await officePage.loginAsNewPPMOfficeUser();
    await use(officePage);
  },
});

exports.expect = base.expect;
