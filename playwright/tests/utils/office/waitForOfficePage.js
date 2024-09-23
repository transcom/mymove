// @ts-check
import * as base from '@playwright/test';

import WaitForPage from '../waitForPage';

/**
 * extension of WaitForPage that provides functions to wait for pages in the office app
 * @extends WaitForPage
 */
export class WaitForOfficePage extends WaitForPage {
  /**
   * @returns {Promise<void>}
   */
  async addNTSShipment() {
    await base.expect(this.page.getByRole('heading', { level: 1 })).toHaveText('Add shipment details');
    await base.expect(this.page.getByTestId('tag')).toHaveText('NTS');
  }

  /**
   * @returns {Promise<void>}
   */
  async addNTSReleaseShipment() {
    await base.expect(this.page.getByRole('heading', { level: 1 })).toHaveText('Add shipment details');
    await base.expect(this.page.getByTestId('tag')).toHaveText('NTS-release');
  }

  /**
   * @returns {Promise<void>}
   */
  async counselingQueue() {
    await base.expect(this.page.getByRole('link', { name: 'Counseling' })).toHaveClass('usa-current');
  }

  /**
   * @returns {Promise<void>}
   */
  async closeoutQueue() {
    await base.expect(this.page.getByRole('link', { name: 'PPM Closeout' })).toHaveClass('usa-current');
  }

  /**
   * @returns {Promise<void>}
   */
  async moveSearchTab() {
    await base.expect(this.page.getByRole('link', { name: 'Move Search' })).toHaveClass('usa-current');
  }

  /**
   * @returns {Promise<void>}
   */
  async moveSearchResults() {
    await base.expect(this.page.getByRole('heading', { level: 2 })).toHaveText('Results (1)');
  }

  /**
   * @returns {Promise<void>}
   */
  async editNTSShipment() {
    await base.expect(this.page.getByRole('heading', { level: 1 })).toHaveText('Edit shipment details');
    await base.expect(this.page.getByTestId('tag')).toHaveText('NTS');
  }

  /**
   * @returns {Promise<void>}
   */
  async editNTSReleaseShipment() {
    await base.expect(this.page.getByRole('heading', { level: 1 })).toHaveText('Edit shipment details');
    await base.expect(this.page.getByTestId('tag')).toHaveText('NTS-release');
  }

  /**
   * @returns {Promise<void>}
   */
  async editMobileHomeShipment() {
    await base.expect(this.page.getByRole('heading', { level: 1 })).toHaveText('Edit shipment details');
    await base.expect(this.page.getByTestId('tag')).toHaveText('Mobile Home');
  }

  /**
   * @returns {Promise<void>}
   */
  async moveDetails() {
    await base.expect(this.page.getByRole('heading', { level: 1 })).toHaveText('Move details');
  }

  /**
   * @returns {Promise<void>}
   */
  async moveOrders() {
    await base.expect(this.page.getByRole('heading', { level: 2, name: 'View orders' })).toBeVisible();
  }

  /**
   * @returns {Promise<void>}
   */
  async moveTaskOrder() {
    await base.expect(this.page.getByRole('heading', { level: 1 })).toHaveText('Move task order');
  }

  /**
   * @returns {Promise<void>}
   */
  async reviewDocumentsConfirmation() {
    await base.expect(this.page.getByRole('heading', { name: 'Send to customer?', level: 3 })).toBeVisible();
  }

  /**
   * @returns {Promise<void>}
   */
  async reviewProGear() {
    await base.expect(this.page.getByRole('heading', { name: 'Review pro-gear 1', level: 3 })).toBeVisible();
  }

  /**
   * @returns {Promise<void>}
   */
  async reviewReceipt() {
    await base.expect(this.page.getByRole('heading', { name: 'Review receipt 1', level: 3 })).toBeVisible();
  }

  /**
   * @returns {Promise<void>}
   */
  async reviewShipmentWeights() {
    await base.expect(this.page.getByRole('heading', { name: 'Review shipment weights', level: 1 })).toBeVisible();
  }

  /**
   * @returns {Promise<void>}
   */
  async reviewWeightTicket() {
    await base.expect(this.page.getByRole('heading', { name: /Review trip 1/i, level: 3 })).toBeVisible();
  }

  /**
   * @param {string} expense name of the expense
   * @param {number}receiptNumber the receipt index number
   * @param {number} expenseIndex the index of the expense
   * @example reveiewExpenseTicket('packing materials', 1, 1);
   * @returns {Promise<void>}
   */
  async reviewExpenseTicket(expense, receiptNumber, expenseIndex) {
    const receiptCheck = `Receipt ${receiptNumber}`;
    const receiptElement = this.page.getByRole('heading', { name: receiptCheck, level: 3, exact: false });

    const expenseCheck = `Review ${expense} #${expenseIndex}`;
    const expenseElement = this.page.getByRole('heading', { name: expenseCheck, level: 3, exact: false });
    await base.expect(receiptElement).toBeVisible({ timeout: 500 });
    await base.expect(expenseElement).toBeVisible({ timeout: 500 });
  }
}

export default WaitForOfficePage;
