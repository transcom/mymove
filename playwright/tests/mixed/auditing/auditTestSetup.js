import TestHarness from '../../utils/testharness';
import * as base from '@playwright/test';
import { stringHelpers, formatRelativeDate, formatNumericDate, textWithNoTrailingNumbers,  } from './auditUtils';
export const { expect } = base;
/**
 * audit test fixture for playwright
 * See https://playwright.dev/docs/test-fixtures
 */
export class AuditTestPage {
  /**
   * Create an BaseTestPage.
   * @param {import('@playwright/test').Page} page
   * @param {import('@playwright/test').APIRequestContext} request
   */
  constructor(page, request) {
    this.page = page;
    this.request = request;

    const baseMilmoveLocal = 'http://milmovelocal:3000';
    const baseOfficeLocal = 'http://officelocal:3000';
    const baseAdminLocal = 'http://adminlocal:3000';

    this.baseURLS = {
      my: process.env.PLAYWRIGHT_MY_URL || baseMilmoveLocal,
      office: process.env.PLAYWRIGHT_OFFICE_URL || baseOfficeLocal,
      admin: process.env.PLAYWRIGHT_ADMIN_URL || baseAdminLocal,
    };

    this.testHarness = new TestHarness(request);

    this.helpers = {
      waitForLoading: async () => {
        await expect(page.locator('button:has-text("Sign out")')).toHaveCount(1, { timeout: 10000 });
        await expect(page.locator('h2[data-name="loading-placeholder"]')).toHaveCount(0);
      },
      signOut: async () => {
        await this.page.getByText('Sign out').click();
      },
      stringHelpers,
      utils: {
        textWithNoTrailingNumbers,
      },
      fast: {
        focusSelector: async (selector) => await page.locator(selector).focus(),
        clickTextAsync: async (selector) => await page.getByText(selector).click(),
        typeInto: async (selector, value) =>
          (async (locator, value) => locator.focus().then(() => locator.type(value)))(page.locator(selector), value),
        typeAndBlur: async (selector, value) => {
          const locator = page.locator(selector);
          await locator.type(value);
          await locator.blur();
        },
        selectValue: async (selector, selectPayload) => await page.locator(selector).selectOption(selectPayload),
      },
    };
  }
}

/**
 * @typedef {object} AuditTestArgs
 * @property {AuditTestPage} pageInstance
 */

/** @type {base.Fixtures<AuditTestArgs, {}, base.PlaywrightTestArgs & base.PlaywrightTestOptions, base.PlaywrightWorkerArgs & base.PlaywrightWorkerOptions>} */
const auditFixtures = {
  pageInstance: async ({ page, request }, use) => {
    const thePage = new AuditTestPage(page, request);
    page.goto('about_blank');
    await use(thePage);
  },
};

export const test = base.test.extend(auditFixtures);
