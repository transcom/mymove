/**
 * customer test fixture for playwright
 * Using javascript classes following the examples in
 * <https://playwright.dev/docs/test-fixtures>
 */
// @ts-check
import * as base from '@playwright/test';

import { BaseTestPage } from './baseTest';
import { WaitForPage } from './waitForPage';

/**
 * devlocal auth user types
 */
export const milmoveUserType = 'milmove';

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
    this.injectAccessibilityTestsInWaitForPage();
    this.waitForPage = new WaitForPage(page);
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
  async navigateBack() {
    await this.page.getByTestId('wizardCancelButton').click();
  }

  /**
   * returns {Promise<void>}
   */
  async navigateForward() {
    await this.page.getByTestId('wizardNextButton').click();
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

export const { expect } = base;

export default test;
