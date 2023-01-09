// office test fixture for playwright
// See https://playwright.dev/docs/test-fixtures
// @ts-check
const base = require('@playwright/test');

const { BaseTestPage } = require('./baseTest');

/**
 * OfficePage
 * @extends BaseTestPage
 */
class OfficePage extends BaseTestPage {
  /**
   * Wait for the page to finish loading.
   */
  async waitForLoading() {
    // make sure we have signed in
    await base.expect(this.page.locator('button:has-text("Sign out")')).toHaveCount(1, { timeout: 10000 });
    await base.expect(this.page.locator('h2[data-name="loading-placeholder"]')).toHaveCount(0);
  }

  /**
   * Goto URL and wait for the page to finish loading.
   *
   * @param {string} url to visit
   */
  async gotoAndWaitForLoading(url) {
    await this.page.goto(url);
    await this.waitForLoading();
  }

  /**
   * Use devlocal auth to sign in as new PPM User
   *
   * @deprecated since the PPM office user is going away
   * @this {BaseOfficePage & SignInMixin}
   */
  async signInAsNewPPMUser() {
    await Promise.all([
      // It is important to call waitForNavigation before click to set up waiting.
      this.page.waitForNavigation(),
      this.signIntoOfficeAsNewPPMUser(),
    ]);
    await this.waitForLoading();
  }

  /**
   * Use devlocal auth to sign in as new Service Counselor
   */
  async signInAsNewServicesCounselorUser() {
    await Promise.all([
      // It is important to call waitForNavigation before click to set up waiting.
      this.page.waitForNavigation(),
      this.signIn.office.newServicesCounselorUser(),
    ]);
    await this.waitForLoading();
  }

  /**
   * Use devlocal auth to sign in as new TIO
   */
  async signInAsNewTIOUser() {
    await Promise.all([
      // It is important to call waitForNavigation before click to set up waiting.
      this.page.waitForNavigation(),
      this.signIn.office.newTIOUser(),
    ]);
    await this.waitForLoading();
  }

  /**
   * Use devlocal auth to sign in as new TOO
   */
  async signInAsNewTOOUser() {
    await Promise.all([
      // It is important to call waitForNavigation before click to set up waiting.
      this.page.waitForNavigation(),
      this.signIn.office.newTOOUser(),
    ]);
    await this.waitForLoading();
  }

  /**
   * Use devlocal auth to sign in as office user with both TOO and TIO roles
   */
  async signInAsNewTIOAndTOOUser() {
    const user = await this.testHarness.buildOfficeUserWithTOOAndTIO();
    await Promise.all([
      // It is important to call waitForNavigation before click to set up waiting.
      this.page.waitForNavigation(),
      this.signIn.office.existingOfficeUser(user.login_gov_email),
    ]);
    await this.waitForLoading();
  }

  /**
   * Use devlocal auth to sign in as office user with prime simulator role
   */
  async signInAsNewPrimeSimulatorUser() {
    await Promise.all([
      // It is important to call waitForNavigation before click to set up waiting.
      this.page.waitForNavigation(),
      this.signIn.office.newPrimeSimulatorUser(),
    ]);
    await this.waitForLoading();
  }
}

/**
 * @typedef {object} OfficePageTestArgs - office page test args
 * @property {OfficePage} officePage    - office page
 */

/** @type {base.Fixtures<OfficePageTestArgs, {}, base.PlaywrightTestArgs, base.PlaywrightWorkerArgs>} */
const officeFixtures = {
  officePage: async ({ page, request }, use) => {
    const officePage = new OfficePage(page, request);
    await use(officePage);
  },
};

exports.test = base.test.extend(officeFixtures);

exports.expect = base.expect;
