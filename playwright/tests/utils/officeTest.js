// office test fixture for playwright
// See https://playwright.dev/docs/test-fixtures
const base = require('@playwright/test');

const {
  signIntoOfficeAsNewPPMUser,
  signIntoOfficeAsNewTIOUser,
  signIntoOfficeAsNewTOOUser,
  signIntoOfficeAsNewServicesCounselorUser,
  signIntoOfficeAsNewPrimeSimulatorUser,
  signIntoOfficeAsNewQAECSRUser,
  signInAsExistingOfficeUser,
} = require('./signIn');
const {
  buildDefaultMove,
  buildWithShipmentMove,
  buildPPMInProgressMove,
  buildOfficeUserWithTOOAndTIO,
  buildHHGMoveWithServiceItemsAndPaymentRequestsAndFiles,
  buildPrimeSimulatorMoveNeedsShipmentUpdate,
  buildHHGMoveWithNTSAndNeedsSC,
} = require('./testharness');

class OfficePage {
  /**
   * @param {base.Page} page
   * @param {base.APIRequestContext} request
   */
  constructor(page, request) {
    this.page = page;
    this.request = request;
  }

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
   */
  async signInAsNewPPMUser() {
    await Promise.all([
      // It is important to call waitForNavigation before click to set up waiting.
      this.page.waitForNavigation(),
      signIntoOfficeAsNewPPMUser(this.page),
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
      signIntoOfficeAsNewServicesCounselorUser(this.page),
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
      signIntoOfficeAsNewTIOUser(this.page),
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
      signIntoOfficeAsNewTOOUser(this.page),
    ]);
    await this.waitForLoading();
  }

  /**
   * Use devlocal auth to sign in as office user with both TOO and TIO roles
   */
  async signInAsNewTIOAndTOOUser() {
    const user = await this.buildOfficeUserWithTOOAndTIO();
    await Promise.all([
      // It is important to call waitForNavigation before click to set up waiting.
      this.page.waitForNavigation(),
      signInAsExistingOfficeUser(this.page, user.login_gov_email),
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
      signIntoOfficeAsNewPrimeSimulatorUser(this.page),
    ]);
    await this.waitForLoading();
  }

  /**
   * Use devlocal auth to sign in as the office user with SC role
   */
  async signInAsNewQAECSRUser() {
    await Promise.all([
      // It is important to call waitForNavigation before click to set up waiting.
      this.page.waitForNavigation(),
      signIntoOfficeAsNewQAECSRUser(this.page),
    ]);
    await this.waitForLoading();
  }

  /**
   * Use testharness to build office user with both TOO and TIO roles
   */
  async buildOfficeUserWithTOOAndTIO() {
    return buildOfficeUserWithTOOAndTIO(this.request);
  }

  /**
   * search for and navigate to move
   * @param {string} moveCode
   */
  async searchForAndNavigateToMove(moveCode) {
    await this.page.locator('input[name="searchText"]').type(moveCode);
    await this.page.locator('input[name="searchText"]').blur();

    await this.page.getByRole('button', { name: 'Search' }).click();
    await this.waitForLoading();

    await base.expect(this.page.locator('tbody >> tr')).toHaveCount(1);
    base.expect(this.page.locator('tbody >> tr').first()).toContainText(moveCode);

    // click result to navigate to move details page
    await this.page.locator('tbody > tr').first().click();
    await this.waitForLoading();

    base.expect(this.page.url()).toContain(`/moves/${moveCode}/details`);
  }

  /**
   * Use testharness to build in progress PPM Move
   */
  async buildInProgressPPMMove() {
    return buildPPMInProgressMove(this.request);
  }

  /**
   * Use testharness to build default PPM Move
   */
  async buildDefaultMove() {
    return buildDefaultMove(this.request);
  }

  /**
   * Use testharness to build move with shipment
   */
  async buildWithShipmentMove() {
    return buildWithShipmentMove(this.request);
  }

  /**
   * Use testharness to build complicated move that will be visible to
   * TOO and TIO
   */
  async buildHHGMoveWithServiceItemsAndPaymentRequestsAndFiles() {
    return buildHHGMoveWithServiceItemsAndPaymentRequestsAndFiles(this.request);
  }

  async buildHHGMoveWithNTSAndNeedsSC() {
    return buildHHGMoveWithNTSAndNeedsSC(this.request);
  }

  /**
   * Use testharness to build complicated move that will be visible to
   * prime simulator
   */
  async buildPrimeSimulatorMoveNeedsShipmentUpdate() {
    return buildPrimeSimulatorMoveNeedsShipmentUpdate(this.request);
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
