/**
 * admin test fixture for playwright
 * Using javascript classes following the examples in
 * <https://playwright.dev/docs/test-fixtures>
 */
// @ts-check
import * as base from '@playwright/test';

import { BaseTestPage } from '../baseTest';

import WaitForAdminPage from './waitForAdminPage';

/**
 * devlocal auth user types
 */
export const AdminUserType = 'admin';

/**
 * AdminPage
 * @extends BaseTestPage
 */
class AdminPage extends BaseTestPage {
  /**
   * Create an AdminPage.
   * @param {import('@playwright/test').Page} page
   * @param {import('@playwright/test').APIRequestContext} request
   */
  constructor(page, request) {
    super(page, request);
    this.waitForPage = new WaitForAdminPage(page);
  }

  /**
   * Create a new admin user and sign in as them
   * @returns {Promise<void>}
   */
  async signInAsNewAdminUser() {
    await this.signInAsNewUser(AdminUserType);
    await this.waitForPage.adminPage();
  }

  /**
   * Create a new admin user and sign in as them
   * @param {string} userId
   * @returns {Promise<void>}
   */
  async signInAsExistingAdminUser(userId) {
    await this.signInAsUserWithId(userId);
    await this.waitForPage.adminPage();
  }

  /**
   * @param {import('aria-query').ARIARole} role
   * @param {Array<string>} labels
   * @param {Object} options
   * @param {boolean} [options.exact=true]
   * @returns {Promise<void>}
   */
  async expectRoleLabelsByText(role, labels, options = { exact: true }) {
    for (const label of labels) {
      await base.expect(this.page.getByRole(role).getByText(label, options)).toBeVisible();
    }
  }

  /**
   * @param {Array<string>} labels
   * @returns {Promise<void>}
   */
  async expectLabels(labels) {
    for (const label of labels) {
      await base
        .expect(this.page.getByRole('paragraph').filter({ has: this.page.locator(`text="${label}"`) }))
        .toBeVisible();
    }
  }
}

/**
 * @typedef {object} AdminPageTestArgs - admin page test args
 * @property {AdminPage} adminPage     - admin page
 */

/** @type {base.Fixtures<AdminPageTestArgs, {}, base.PlaywrightTestArgs, base.PlaywrightWorkerArgs>} */
const adminFixtures = {
  adminPage: async ({ page, request }, use) => {
    const adminPage = new AdminPage(page, request);
    await use(adminPage);
  },
};

export const test = base.test.extend(adminFixtures);

export const { expect } = base;

export default test;
