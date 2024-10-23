/**
 * office test fixture for playwright
 * Using javascript classes following the examples in
 * <https://playwright.dev/docs/test-fixtures>
 */
// @ts-check
import * as base from '@playwright/test';

import { BaseTestPage } from '../baseTest';

import WaitForOfficePage from './waitForOfficePage';

export const { expect } = base;

/**
 * devlocal auth user types
 */
export const TOOOfficeUserType = 'TOO office';
export const TIOOfficeUserType = 'TIO office';
export const QAEOfficeUserType = 'QAE office';
export const CustomerServiceRepresentativeOfficeUserType = 'CSR office';
export const ServicesCounselorOfficeUserType = 'Services Counselor office';
export const GSROfficeUserType = 'GSR office';
export const PrimeSimulatorUserType = 'Prime Simulator office';
export const MultiRoleOfficeUserType = 'Multi role office';

export const DEPARTMENT_INDICATOR_OPTIONS = {
  AIR_AND_SPACE_FORCE: '57 Air Force and Space Force',
  ARMY: '21 Army',
  ARMY_CORPS_OF_ENGINEERS: '96 Army Corps of Engineers',
  COAST_GUARD: '70 Coast Guard',
  NAVY_AND_MARINES: '17 Navy and Marine Corps',
  OFFICE_OF_SECRETARY_OF_DEFENSE: '97 Office of the Secretary of Defense',
};

/**
 * office test fixture for playwright
 * See https://playwright.dev/docs/test-fixtures
 * @extends BaseTestPage
 */
export class OfficePage extends BaseTestPage {
  /**
   * Create an OfficePage.
   * @param {import('@playwright/test').Page} page
   * @param {import('@playwright/test').APIRequestContext} request
   */
  constructor(page, request) {
    super(page, request);
    this.waitForPage = new WaitForOfficePage(page);
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
   * Sign in as existing office user with devlocal
   *
   * @param {string} email
   */
  async signInAsExistingOfficeUser(email) {
    await this.page.goto('/devlocal-auth/login');
    await this.page.locator('input[name=email]').fill(email);
    await this.page.locator('p', { hasText: 'User Email' }).locator('button').click();
  }

  /**
   * Use devlocal auth to sign in as new Service Counselor
   */
  async signInAsNewServicesCounselorUser() {
    await this.signInAsNewUser(ServicesCounselorOfficeUserType);
    await this.page.getByRole('heading', { name: 'Moves' }).waitFor();
  }

  /**
   * Use devlocal auth to sign in as new Multi-role User
   */
  async signInAsNewMultiRoleUser() {
    await this.signInAsNewUser(MultiRoleOfficeUserType);
    await this.page.getByText('Change user role').waitFor();
  }

  /**
   * Use devlocal auth to sign in as new TIO
   */
  async signInAsNewTIOUser() {
    await this.signInAsNewUser(TIOOfficeUserType);
    await this.page.getByRole('heading', { name: 'Payment requests' }).waitFor();
  }

  /**
   * Use devlocal auth to sign in as new TOO
   */
  async signInAsNewTOOUser() {
    await this.signInAsNewUser(TOOOfficeUserType);
    await this.page.getByRole('heading', { name: 'All moves' }).waitFor();
  }

  /**
   * Use devlocal auth to sign in as office user with both TOO and TIO roles
   */
  async signInAsNewTIOAndTOOUser() {
    const user = await this.testHarness.buildOfficeUserWithTOOAndTIO();
    await this.signInAsExistingOfficeUser(user.okta_email);
  }

  /**
   * Use devlocal auth to sign in as office user with prime simulator role
   */
  async signInAsNewPrimeSimulatorUser() {
    await this.signInAsNewUser(PrimeSimulatorUserType);
    await this.page.getByRole('heading', { name: 'Moves available to Prime' }).waitFor();
  }

  /**
   * Use devlocal auth to sign in as office user with qae role
   */
  async signInAsNewQAEUser() {
    await this.signInAsNewUser(QAEOfficeUserType);
    await this.page.getByRole('heading', { name: 'Search for a move' }).waitFor();
  }

  /**
   * Use devlocal auth to sign in as new Multi-role User
   */
  async signInAsNewGSRUser() {
    await this.signInAsNewUser(GSROfficeUserType);
    await this.page.getByRole('heading', { name: 'Search for a move' }).waitFor();
  }

  /**
   * Use devlocal auth to sign in as an office user with the customer service representative role
   */
  async signInAsNewCustomerServiceRepresentativeUser() {
    await this.signInAsNewUser(CustomerServiceRepresentativeOfficeUserType);
    await this.page.getByRole('heading', { name: 'Search for a move' }).waitFor();
  }

  /**
   * search for and navigate to move (Prime Simulator role)
   * @param {string} moveLocator
   */
  async primeSimulatorNavigateToMove(moveLocator) {
    await this.page.locator('input[name="moveCode"]').fill(moveLocator);
    await this.page.locator('input[name="moveCode"]').blur();

    // Click the first returned row
    await this.page.getByTestId('locator-0').click();
    await this.waitForPage.moveDetails();
    await expect(this.page.getByText(moveLocator)).toBeVisible();
  }

  /**
   * search for and navigate to move for QAE
   * @param {string} moveLocator
   */
  async qaeSearchForAndNavigateToMove(moveLocator) {
    await this.page.locator('input[name="searchText"]').fill(moveLocator);
    await this.page.locator('input[name="searchText"]').blur();

    await this.page.getByRole('button', { name: 'Search' }).click();
    await this.page.getByRole('heading', { name: 'Results' }).waitFor();

    await base.expect(this.page.locator('tbody >> tr')).toHaveCount(1);
    await base.expect(this.page.locator('tbody >> tr').first()).toContainText(moveLocator);

    // click result to navigate to move details page
    await this.page.locator('tbody > tr').first().click();
    await this.waitForLoading();

    base.expect(this.page.url()).toContain(`/moves/${moveLocator}/details`);
  }

  /**
   * search for and navigate to move for CSR
   * @param {string} moveLocator
   */
  async csrSearchForAndNavigateToMove(moveLocator) {
    await this.page.locator('input[name="searchText"]').fill(moveLocator);
    await this.page.locator('input[name="searchText"]').blur();

    await this.page.getByRole('button', { name: 'Search' }).click();
    await this.page.getByRole('heading', { name: 'Results' }).waitFor();

    await base.expect(this.page.locator('tbody >> tr')).toHaveCount(1);
    await base.expect(this.page.locator('tbody >> tr').first()).toContainText(moveLocator);

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
    await this.page.goto(`/moves/${moveLocator}/payment-requests/`);
    await this.page.waitForLoadState('networkidle', { timeout: 30000 });
    await this.page.getByRole('heading', { name: 'Payment requests' }).first().waitFor();
  }

  /**
   * TOO search for and navigate to move
   * @param {string} moveLocator
   */
  async tooNavigateToMove(moveLocator) {
    await this.page.goto(`/moves/${moveLocator}/details/`);
    await this.page.waitForLoadState('networkidle', { timeout: 30000 });
    await this.page.getByRole('heading', { name: 'Move details' }).waitFor();
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

export const test = base.test.extend(officeFixtures);

export default test;
