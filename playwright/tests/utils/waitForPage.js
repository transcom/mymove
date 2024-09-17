// @ts-check
import * as base from '@playwright/test';
import { checkA11y } from 'axe-playwright';

export class WaitForPage {
  /**
   * Create a WaitForPage.
   * Use a class so we can catch typos or type errors
   *
   * @param {import('@playwright/test').Page} page
   */
  constructor(page) {
    this.page = page;
  }

  /**
   * Run an accessibility audit against the page in its current state
   *
   */

  async runAccessibilityAudit() {
    await checkA11y(
      this.page,
      undefined,
      {
        detailedReport: true,
        detailedReportOptions: { html: true },
      },
      // skip failures
      false,
      'html',
    );
  }

  /**
   * @returns {Promise<void>}
   */
  async localLogin() {
    await this.runAccessibilityAudit();
    await base.expect(this.page.getByRole('heading', { name: 'Select an Existing User' })).toBeVisible();
    await this.runAccessibilityAudit();
  }
}

export default WaitForPage;
