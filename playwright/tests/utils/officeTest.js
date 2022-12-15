const base = require('@playwright/test');

const {
  signIntoOfficeAsNewPPMUser,
  signIntoOfficeAsNewTIOUser,
  signIntoOfficeAsNewTOOUser,
  signIntoOfficeAsNewServicesCounselorUser,
} = require('./signIn');
const { buildPPMInProgressMove } = require('./testharness');

class OfficePage {
  /**
   * @param {import('@playwright/test').Page} page
   */
  constructor(page, request) {
    this.page = page;
    this.request = request;
  }

  async waitForLoading() {
    await base.expect(this.page.locator('h2[data-name="loading-placeholder"]')).toHaveCount(0);
  }

  async gotoAndWaitForLoading(url) {
    await this.page.goto(url);
    await this.waitForLoading();
  }

  async signInAsNewPPMUser() {
    await signIntoOfficeAsNewPPMUser(this.page);
    await this.waitForLoading();
  }

  async signInAsNewServicesCounselorUser() {
    await signIntoOfficeAsNewServicesCounselorUser(this.page);
    await this.waitForLoading();
  }

  async signInAsNewTIOUser() {
    await signIntoOfficeAsNewTIOUser(this.page);
    await this.waitForLoading();
  }

  async signInAsNewTOOUser() {
    await signIntoOfficeAsNewTOOUser(this.page);
    await this.waitForLoading();
  }

  async buildInProgressPPMMove() {
    return buildPPMInProgressMove(this.request);
  }
}

exports.test = base.test.extend({
  officePage: async ({ page, request }, use) => {
    const officePage = new OfficePage(page, request);
    await use(officePage);
  },
});

exports.expect = base.expect;
