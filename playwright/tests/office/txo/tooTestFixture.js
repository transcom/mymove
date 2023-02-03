/**
 * Semi-automated converted from a cypress test, and thus may contain
 * non best-practices, in particular: heavy use of `page.locator`
 * instead of `page.getBy*`.
 */

// @ts-check
const { expect, OfficePage } = require('../../utils/officeTest');
/**
 * TooFlowPage test fixture
 *
 * The logic in TooFlowPage is only used in this file, so keep the
 * playwright test fixture in this file.
 * @extends OfficePage
 */
export class TooFlowPage extends OfficePage {
  /**
   * @param {OfficePage} officePage
   * @param {Object} move
   * @override
   */
  constructor(officePage, move) {
    super(officePage.page, officePage.request);
    this.move = move;
    this.moveLocator = move.locator;
  }

  /**
   * select and approve all service items on the page
   */
  async selectAndApproveAllServiceItems() {
    // Select & approve items
    const checkboxes = this.page.locator('input[data-testid="shipment-display-checkbox"]');
    const boxCount = await checkboxes.count();
    expect(boxCount).toBeGreaterThan(0);
    for (let i = 0; i < boxCount; i += 1) {
      const id = await checkboxes.nth(i).getAttribute('id');
      const label = this.page.locator(`label[for="${id}"]`);
      await expect(label).toHaveCount(1);
      await label.click();
    }
  }

  /**
   * approve all shipments on the move
   * @param {Object} options
   * @param {boolean} [options.usesExternalVendor=false]
   * @param {boolean} [options.hasServiceItems=true]
   */
  async approveAllShipments(options = {}) {
    const { usesExternalVendor = false, hasServiceItems = true } = options;
    const shipmentCount = await this.page.locator('input[data-testid="shipment-display-checkbox"]').count();
    if (hasServiceItems) {
      // Select & approve items
      await this.selectAndApproveAllServiceItems();
    }
    // Select additional service items
    await this.page.locator('label[for="shipmentManagementFee"]').click();
    await this.page.locator('label[for="counselingFee"]').click();
    // Open modal
    await this.page.getByRole('button', { name: 'Approve selected' }).click();

    const modal = this.page.locator('#approvalConfirmationModal [data-testid="modal"]');
    await expect(modal).toBeVisible();
    // Verify modal content
    await expect(modal.getByText('Preview and post move task order')).toBeVisible();
    await expect(this.page.locator('#approvalConfirmationModal [data-testid="ShipmentContainer"]')).toHaveCount(
      shipmentCount,
    );
    expect(modal.getByText('Approved service items for this move')).toBeVisible();
    const siTable = modal.getByText('Approved service items for this move').locator('..').locator('table');
    await expect(siTable).toContainText('Move management');
    await expect(siTable).toContainText('Counseling');

    if (usesExternalVendor) {
      const shipment = this.page.locator('[data-testid="modal"] [data-testid="ShipmentContainer"]').last();
      await expect(shipment.locator('[data-testid="usesExternalVendor"]')).toBeVisible();
      await expect(shipment.locator('[data-testid="tacType"]')).toBeVisible();
      await expect(shipment.locator('[data-testid="sacType"]')).toBeVisible();
    }

    // Click approve
    await modal.getByRole('button', { name: 'Approve and send' }).click();
    // In CircleCI it can take more than 5 seconds for this modal to disappear
    await expect(modal).not.toBeVisible({ timeout: 10000 });
    await this.waitForLoading();
  }

  /**
   * edit tac sac for NTS
   */
  async editTacSac() {
    await this.page.locator('[data-testid="ShipmentContainer"] .usa-button').last().click();
    await this.waitForLoading();

    await this.page.locator('[data-testid="grid"] button').getByText('Add or edit codes').click();

    const form = this.page.locator('form');
    await form.locator('select[name="departmentIndicator"]').selectOption({ label: '21 Army' });
    await form.locator('input[name="ordersNumber"]').type('ORDER66');
    await form.locator('select[name="ordersType"]').selectOption({ label: 'Permanent Change Of Station (PCS)' });
    await form.locator('select[name="ordersTypeDetail"]').selectOption({ label: 'Shipment of HHG Permitted' });

    await form.locator('[data-testid="hhgTacInput"]').type('E15A');
    await form.locator('[data-testid="hhgSacInput"]').type('4K988AS098F');
    await form.locator('[data-testid="ntsTacInput"]').type('F123');
    await form.locator('[data-testid="ntsSacInput"]').type('3L988AS098F');
    // Edit orders page | Save
    await form.getByRole('button', { name: 'Save' }).click();
    await this.waitForLoading();

    await expect(this.page.locator('[data-testid="tacMDC"]')).toContainText('E15A');
    await expect(this.page.locator('[data-testid="sacSDN"]')).toContainText('4K988AS098F');
    await expect(this.page.locator('[data-testid="NTStac"]')).toContainText('F123');
    await expect(this.page.locator('[data-testid="NTSsac"]')).toContainText('3L988AS098F');
  }
}

export default TooFlowPage;
