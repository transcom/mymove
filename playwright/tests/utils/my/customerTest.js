/**
 * customer test fixture for playwright
 * Using javascript classes following the examples in
 * <https://playwright.dev/docs/test-fixtures>
 */
// @ts-check
import * as base from '@playwright/test';

import { BaseTestPage } from '../baseTest';

import WaitForCustomerPage from './waitForCustomerPage';

/**
 * devlocal auth user types
 */
export const milmoveUserType = 'milmove';

export const { expect } = base;

/**
 * CustomerPage
 * @extends BaseTestPage
 */
export class CustomerPage extends BaseTestPage {
  /**
   * Create an CustomerPage.
   * @param {import('@playwright/test').Page} page
   * @param {import('@playwright/test').APIRequestContext} request
   */
  constructor(page, request) {
    super(page, request);
    this.waitForPage = new WaitForCustomerPage(page);
  }

  /**
   * Sign in as new customer with devlocal
   *
   * returns {Promise<void>}
   */
  async signInAsNewCustomer() {
    await this.signInAsNewUser(milmoveUserType);
  }

  /**
   * Submits a validation code
   *
   * returns {Promise<void>}
   */
  async submitParameterValue() {
    const testCode = '123456';
    await this.page.locator('[name="code"]').fill(testCode);
    await expect(this.page.getByTestId('wizardNextButton')).toBeVisible();

    // Regex for the path of the validation code api call
    const pathRegex = /\/internal\/application_parameters$/;

    // Mock the api call and its response
    await this.page.route(pathRegex, async (route) => {
      await route.fulfill({
        status: 200,
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ parameterValue: '123456', parameterName: 'validation_code' }),
      });
    });

    // Click on the submit button
    await this.page.getByTestId('wizardNextButton').click();
  }

  /**
   * Sign in as existing customer with devlocal
   *
   * @param {string} userId
   * returns {Promise<void>}
   */
  async signInAsExistingCustomer(userId) {
    await this.signInAsUserWithId(userId);
    // ensure the home page has loaded
    await this.page.getByLabel('Home').waitFor();
  }

  /**
   * returns {Promise<void>}
   */
  async navigateFromMMDashboardToMove(move) {
    await expect(this.page.getByTestId('goToMoveBtn')).toBeVisible();

    await this.page.getByTestId('goToMoveBtn').click();

    await expect(this.page.getByTestId('customer-header').getByText(`${move.locator}`)).toBeVisible();

    const targetElements = await this.page.$$(`[data-testid="shipment-list-item-container"]`);

    for (const element of targetElements) {
      const matches = (await element.textContent()).match(/[0-9|A-Z]{6}-[0-9]{2}/);
      expect(matches).not.toBeNull();
    }
  }

  /**
   * returns {Promise<void>}
   */
  async navigateBack() {
    await this.page.getByTestId('wizardCancelButton').click();
  }

  /**
   * returns {Promise<void>}
   */
  async navigateForward() {
    await this.page.getByTestId('wizardNextButton').click();
  }

  /**
   * @param {Object} move
   * returns {Promise<void>}
   */
  async signInForPPMWithMove(move) {
    await this.signInAsExistingCustomer(move?.Orders?.service_member?.user_id);
  }

  /**
   * returns {Promise<void>}
   */
  async createMoveButtonClick() {
    await this.page.getByTestId('createMoveBtn').click();
  }
}

/**
 * set the viewport for mobile dimensions
 */
export function useMobileViewport() {
  // https://playwright.dev/docs/emulation#viewport
  base.test.use({ viewport: { width: 479, height: 875 } });
}

/**
 * @typedef {Object} ViewportCallbackProps
 * @property {string} viewportName
 * @property {boolean} isMobile
 */

/**
 * @callback forEachViewportCallback
 * @param {ViewportCallbackProps} props
 */

/**
 * @param {forEachViewportCallback} callbackfn
 */
export function forEachViewport(callbackfn) {
  const viewportsString = process.env.PLAYWRIGHT_VIEWPORTS || 'desktop mobile';
  const viewports = viewportsString.split(/\s+/);
  // use forEach to avoid
  // https://eslint.org/docs/latest/rules/no-loop-func
  viewports.forEach(async (viewportName) => {
    const isMobile = viewportName === 'mobile';
    //
    // https://playwright.dev/docs/test-parameterize
    //
    base.test.describe(`with ${viewportName} viewport`, async () => {
      if (isMobile) {
        useMobileViewport();
      }
      await callbackfn({ viewportName, isMobile });
    });
  });
}

/**
 * @typedef {object} CustomerPageTestArgs - customer page test args
 * @property {CustomerPage} customerPage  - customer page
 */

/** @type {base.Fixtures<CustomerPageTestArgs, {}, base.PlaywrightTestArgs, base.PlaywrightWorkerArgs>} */
const officeFixtures = {
  customerPage: async ({ page, request }, use) => {
    const customerPage = new CustomerPage(page, request);
    await use(customerPage);
  },
};

export const test = base.test.extend(officeFixtures);

export default test;
