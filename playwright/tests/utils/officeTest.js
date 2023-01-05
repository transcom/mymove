/**
 * office test fixture for playwright
 * Using javascript classes following the examples in
 * <https://playwright.dev/docs/test-fixtures>
 */
// @ts-check
const base = require('@playwright/test');

const { BaseTestPage } = require('./baseTest');

/**
 * office test fixture for playwright
 * See https://playwright.dev/docs/test-fixtures
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

  /**
   * Use devlocal auth to sign in as office user with qaecsr role
   */
  async signInAsNewQAECSRUser() {
    await Promise.all([
      // It is important to call waitForNavigation before click to set up waiting.
      this.page.waitForNavigation(),
      this.signIn.office.newQAECSRUser(),
    ]);
    await this.waitForLoading();
  }

  /**
   * search for and navigate to move
   * @param {string} moveLocator
   */
  async qaeCsrSearchForAndNavigateToMove(moveLocator) {
    await this.page.locator('input[name="searchText"]').type(moveLocator);
    await this.page.locator('input[name="searchText"]').blur();

    await this.page.getByRole('button', { name: 'Search' }).click();
    await this.waitForLoading();

    await base.expect(this.page.locator('tbody >> tr')).toHaveCount(1);
    base.expect(this.page.locator('tbody >> tr').first()).toContainText(moveLocator);

    // click result to navigate to move details page
    await this.page.locator('tbody > tr').first().click();
    await this.waitForLoading();

    base.expect(this.page.url()).toContain(`/moves/${moveLocator}/details`);
  }

  /**
   * TIO search for and navigate to move
   * @param {string} moveLocator
   */
  async tioNavigateToMove(moveLocator) {
    await this.page.locator('#locator').type(moveLocator);
    await this.page.locator('th[data-testid="locator"]').first().click();
    await this.page.locator('[data-testid="locator-0"]').click();
    await this.waitForLoading();
  }

  /**
   * TOO search for and navigate to move
   * @param {string} moveLocator
   */
  async tooNavigateToMove(moveLocator) {
    await this.page.locator('input[name="locator"]').type(moveLocator);
    await this.page.locator('input[name="locator"]').blur();

    // click result to navigate to move details page
    await this.page.locator('tbody > tr').first().click();
    await this.waitForLoading();
    base.expect(this.page.url()).toContain(`/moves/${moveLocator}/details`);
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

exports.OfficePage = OfficePage;

exports.test = base.test.extend(officeFixtures);

exports.expect = base.expect;
