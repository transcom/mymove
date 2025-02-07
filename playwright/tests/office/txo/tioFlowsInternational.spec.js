import { test, expect, OfficePage } from '../../utils/office/officeTest';

/**
 * TioFlowPage test fixture
 *
 * The logic in TioFlowPage is only used in this file, so keep the
 * playwright test fixture in this file.
 * @extends OfficePage
 */
class TioFlowPage extends OfficePage {
  /**
   * @param {OfficePage} officePage
   * @param {Object} move
   * @param {Boolean} usePaymentRequest
   * @override
   */
  constructor(officePage, move, usePaymentRequest) {
    super(officePage.page, officePage.request);
    this.move = move;
    this.moveLocator = move.locator;
    if (usePaymentRequest !== false) {
      this.paymentRequest = this.findPaymentRequestBySequenceNumber(1);
    }
  }

  /**
   * @param {number} sequenceNumber
   * @returns {object}
   */
  findPaymentRequestBySequenceNumber(sequenceNumber) {
    return this.move.PaymentRequests.find((pr) => pr.sequence_number === sequenceNumber);
  }

  /**
   * Complete the service item card
   * @param {import('@playwright/test').Locator} serviceItemCardLocator
   * @param {boolean} approve
   */
  async completeServiceItemCard(serviceItemCardLocator, approve = false) {
    // serviceItemAmount
    if (!approve) {
      const inputEl = serviceItemCardLocator.locator('input[data-testid="rejectRadio"]');
      const id = await inputEl.getAttribute('id');
      await this.page.locator(`label[for="${id}"]`).click();
      await this.page.locator('textarea[data-testid="rejectionReason"]').fill('This is not a valid request');
    } else {
      const inputEl = serviceItemCardLocator.locator('input[data-testid="approveRadio"]');
      const id = await inputEl.getAttribute('id');
      await this.page.locator(`label[for="${id}"]`).click();
    }
    await this.slowDown();
  }

  /**
   * approve the service item
   */
  async approveServiceItem() {
    await this.completeServiceItem(true);
  }

  /**
   * complete the service item, approving or rejecting it
   * @param {boolean} approved
   */
  async completeServiceItem(approved) {
    const cards = this.page.getByTestId('ServiceItemCard');
    const cardCount = await cards.count();
    expect(cardCount).toBeGreaterThan(0);
    for (let i = 0; i < cardCount; i += 1) {
      await this.completeServiceItemCard(cards.nth(i), approved);
    }
  }

  async slowDown() {
    await this.page.waitForLoadState('networkidle');
    // sleep for 500ms
    await new Promise((r) => {
      setTimeout(() => r(undefined), 500);
    });
  }
}

const alaskaEnabled = process.env.FEATURE_FLAG_ENABLE_ALASKA;

test.describe('TIO user', () => {
  /** @type {TioFlowPage} */
  let tioFlowPage;
  test.skip(alaskaEnabled === 'false', 'Skip if Alaska FF is disabled.');
  test('can review a payment request for a basic iHHG Alaska move', async ({ page, officePage }) => {
    test.slow();
    const move = await officePage.testHarness.buildnternationalHHGMoveWithServiceItemsandPaymentRequestsForTIO();
    await officePage.signInAsNewTIOUser();

    tioFlowPage = new TioFlowPage(officePage, move, true);
    await tioFlowPage.waitForLoading();
    await officePage.tioNavigateToMove(tioFlowPage.moveLocator);
    await officePage.page.getByRole('heading', { name: 'Payment Requests', exact: true }).waitFor();
    expect(page.url()).toContain('/payment-requests');
    await expect(page.getByTestId('MovePaymentRequests')).toBeVisible();

    const prNumber = tioFlowPage.paymentRequest.payment_request_number;
    const prHeading = page.getByRole('heading', { name: `Payment Request ${prNumber}` });
    await expect(prHeading).toBeVisible();
    await tioFlowPage.waitForLoading();

    await page.getByRole('button', { name: 'Review service items' }).click();

    await page.waitForURL(`**/payment-requests/${tioFlowPage.paymentRequest.id}`);
    await tioFlowPage.waitForLoading();

    // there should be four service items - let's approve all of them
    await expect(page.getByTestId('ReviewServiceItems')).toBeVisible();
    await expect(page.getByText('International Shipping & Linehaul')).toBeVisible();
    await page.getByText('Show calculations').click();
    await expect(page.locator('[data-testid="ServiceItemCalculations"]')).toContainText('Calculations');
    await expect(page.locator('[data-testid="ServiceItemCalculations"]')).toContainText('Billable weight (cwt)');
    await expect(page.locator('[data-testid="ServiceItemCalculations"]')).toContainText('ISLH price');
    await expect(page.locator('[data-testid="ServiceItemCalculations"]')).toContainText('Price escalation factor');
    // approve
    await tioFlowPage.approveServiceItem();
    await page.getByText('Next').click();
    await tioFlowPage.slowDown();

    await expect(page.getByText('International HHG Pack')).toBeVisible();
    await page.getByText('Show calculations').click();
    await expect(page.locator('[data-testid="ServiceItemCalculations"]')).toContainText('Calculations');
    await expect(page.locator('[data-testid="ServiceItemCalculations"]')).toContainText('Billable weight (cwt)');
    await expect(page.locator('[data-testid="ServiceItemCalculations"]')).toContainText('International Pack price');
    await expect(page.locator('[data-testid="ServiceItemCalculations"]')).toContainText('Price escalation factor');
    // approve
    await tioFlowPage.approveServiceItem();
    await page.getByText('Next').click();
    await tioFlowPage.slowDown();

    await expect(page.getByText('International HHG Unpack')).toBeVisible();
    await page.getByText('Show calculations').click();
    await expect(page.locator('[data-testid="ServiceItemCalculations"]')).toContainText('Calculations');
    await expect(page.locator('[data-testid="ServiceItemCalculations"]')).toContainText('Billable weight (cwt)');
    await expect(page.locator('[data-testid="ServiceItemCalculations"]')).toContainText('International Unpack price');
    await expect(page.locator('[data-testid="ServiceItemCalculations"]')).toContainText('Price escalation factor');
    // approve
    await tioFlowPage.approveServiceItem();
    await page.getByText('Next').click();
    await tioFlowPage.slowDown();

    await expect(page.getByText('International destination shuttle service')).toBeVisible();
    await page.getByText('Show calculations').click();
    await expect(page.locator('[data-testid="ServiceItemCalculations"]')).toContainText('Calculations');
    await expect(page.locator('[data-testid="ServiceItemCalculations"]')).toContainText('Billable weight (cwt)');
    await expect(page.locator('[data-testid="ServiceItemCalculations"]')).toContainText('Destination price');
    await expect(page.locator('[data-testid="ServiceItemCalculations"]')).toContainText('Price escalation factor');
    // approve
    await tioFlowPage.approveServiceItem();
    await page.getByText('Next').click();
    await tioFlowPage.slowDown();

    await expect(page.getByText('International POE Fuel Surcharge')).toBeVisible();
    await page.getByText('Show calculations').click();
    await expect(page.locator('[data-testid="ServiceItemCalculations"]')).toContainText('Calculations');
    await expect(page.locator('[data-testid="ServiceItemCalculations"]')).toContainText('Billable weight (cwt)');
    await expect(page.locator('[data-testid="ServiceItemCalculations"]')).toContainText('Mileage');
    await expect(page.locator('[data-testid="ServiceItemCalculations"]')).toContainText('ZIP 74133 to Port ZIP 98424');
    await expect(page.locator('[data-testid="ServiceItemCalculations"]')).toContainText('Mileage factor');
    // approve
    await tioFlowPage.approveServiceItem();
    await page.getByText('Next').click();
    await tioFlowPage.slowDown();

    // that should be it
    await expect(page.getByText('needs your review')).toHaveCount(0, { timeout: 10000 });
    await page.getByText('Complete request').click();

    await expect(page.locator('[data-testid="requested"]')).toContainText('$4,287.71');
    await expect(page.locator('[data-testid="accepted"]')).toContainText('$4,287.71');
    await expect(page.locator('[data-testid="rejected"]')).toContainText('$0.00');

    await page.getByText('Authorize payment').click();
    await tioFlowPage.waitForLoading();

    await tioFlowPage.slowDown();
    expect(page.url()).toContain('/payment-requests');

    await expect(page.getByTestId('tag')).toBeVisible();
    await expect(page.getByTestId('tag').getByText('Reviewed')).toHaveCount(1);

    // ensure the payment request we approved no longer has the "Review Service Items" button
    await expect(page.getByText('Review Service Items')).toHaveCount(0);

    // Go back to queue
    await page.locator('a[title="Home"]').click();
    await tioFlowPage.waitForLoading();

    // search for the moveLocator in case this move doesn't show up on the first page
    await page.locator('#locator').fill(tioFlowPage.moveLocator);
    await page.locator('#locator').blur();
    const paymentSection = page.locator(`[data-uuid="${tioFlowPage.paymentRequest.id}"]`);
    // the payment request that is now in the "Reviewed" status will no longer appear
    // in the TIO queue - only "Payment requested" moves will appear
    await expect(paymentSection.locator('td', { hasText: 'Reviewed' })).not.toBeVisible();
  });
});
