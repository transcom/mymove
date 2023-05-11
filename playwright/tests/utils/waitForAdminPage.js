// @ts-check
import * as base from '@playwright/test';

import WaitForPage from './waitForPage';

/**
 * extension of WaitForPage that provides functions to wait for pages in the admin app
 * @extends WaitForPage
 */
export class WaitForAdminPage extends WaitForPage {
  /**
   * Wait for the loading placeholder to go away
   * @returns {Promise<void>}
   */
  async waitForLoading() {
    await base.expect(this.page.locator('h2[data-name="loading-placeholder"]')).toHaveCount(0);
    await base.expect(this.page.locator('svg.MuiCircularProgress-svg')).toHaveCount(0);
  }

  /**
   * wait for the admin page to finish loading
   * @returns {Promise<void>}
   */
  async adminPage() {
    // ensure the admin page has fully loaded before moving on
    await base.expect(this.page.locator('a:has-text("Logout")')).toHaveCount(1, { timeout: 10000 });
    await this.waitForLoading();
    await base.expect(this.page.locator('.ReactTable').locator('.-loading.-active')).toHaveCount(0);
  }
}

export default WaitForAdminPage;
