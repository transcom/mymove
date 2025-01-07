/**
 * Semi-automated converted from a cypress test, and thus may contain
 * non best-practices, in particular: heavy use of `page.locator`
 * instead of `page.getBy*`.
 */

// @ts-check
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
   * @constant TIOTabsTitles Search tabs available for TIO
   * @readonly
   */
  TIOTabsTitles = ['Payment Request Queue', 'Search'];

  /**
   * @constant SearchRBSelection Available search options within Search tab
   * @readonly
   */
  SearchRBSelection = ['Move Code', 'DOD ID', 'Customer Name'];

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
   * reject the service item
   */
  async rejectServiceItem() {
    await this.completeServiceItem(false);
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

  async approveCurrentServiceItem() {
    await this.completeServiceItem(true);
    await this.page.getByText('Next').click();
    await this.slowDown();
  }

  async validateDLCalcValues() {
    const siCalc = this.page.locator('[data-testid="ServiceItemCalculations"]');
    await expect(siCalc).toContainText('14 cwt');
    await expect(siCalc).toContainText('354');
    await expect(siCalc).toContainText('ZIP 80301 to ZIP 80501');
    await expect(siCalc).toContainText('21');
    await expect(siCalc).toContainText('Domestic non-peak');
    await expect(siCalc).toContainText('Origin service area: 144');
    await expect(siCalc).toContainText('1.01000');
    await expect(siCalc).toContainText('$800.00');
  }

  async validateFSCalcValues() {
    const siCalc = this.page.locator('[data-testid="ServiceItemCalculations"]');
    await expect(siCalc).toContainText('14 cwt');
    await expect(siCalc).toContainText('354');
    await expect(siCalc).toContainText('ZIP 80301 to ZIP 80501');
    await expect(siCalc).toContainText('0.1');
    await expect(siCalc).toContainText('EIA diesel: $2.81');
    await expect(siCalc).toContainText('Weight-based distance multiplier: 0.0004170');
    await expect(siCalc).toContainText('$107.00');
  }

  async validateDUCalcValues() {
    const siCalc = this.page.locator('[data-testid="ServiceItemCalculations"]');
    await expect(siCalc).toContainText('43 cwt');
    await expect(siCalc).toContainText('5.79');
    await expect(siCalc).toContainText('Destination service schedule: 1');
    await expect(siCalc).toContainText('Domestic non-peak');
    await expect(siCalc).toContainText('1.04071');
    await expect(siCalc).toContainText('Base year: DUPK Test Year');
    await expect(siCalc).toContainText('$459.00');
  }

  async validateDXCalcValues() {
    const siCalc = this.page.locator('[data-testid="ServiceItemCalculations"]');
    await expect(siCalc).toContainText('43 cwt');
    await expect(siCalc).toContainText('6.25');
    await expect(siCalc).toContainText('service area: 144');
    await expect(siCalc).toContainText('Domestic non-peak');
    await expect(siCalc).toContainText('1.04071');
    await expect(siCalc).toContainText('$150.00');
  }

  // ugh, needed for flaky tests
  // without this, local tests on a fast machine would fail maybe 1/3
  // of the time
  // pretty sure the client app is not updating its state correctly
  // when clicking through payment requests too quickly
  async slowDown() {
    // sigh, give the app time to catch up
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
    // this is a long test
    test.slow();
    const move = await officePage.testHarness.buildnternationalHHGMoveWithServiceItemsandPaymentRequestsForTIO();
    await officePage.signInAsNewTIOUser();

    tioFlowPage = new TioFlowPage(officePage, move, true);
    await tioFlowPage.waitForLoading();
    await officePage.tioNavigateToMove(tioFlowPage.moveLocator);
    await officePage.page.getByRole('heading', { name: 'Payment Requests', exact: true }).waitFor();
    // Payment Requests page
    expect(page.url()).toContain('/payment-requests');
    await expect(page.getByTestId('MovePaymentRequests')).toBeVisible();

    const prNumber = tioFlowPage.paymentRequest.payment_request_number;
    const prHeading = page.getByRole('heading', { name: `Payment Request ${prNumber}` });
    await expect(prHeading).toBeVisible();
    await tioFlowPage.waitForLoading();

    // // confirm the move is set up as expected
    // await expect(page.getByRole('heading', { name: '$1,130.21' })).toBeVisible();
    // await expect(page.getByRole('heading', { name: '$805.55' })).toBeVisible();

    // again, find the first payment request by sequence number and
    // operate on that
    await page.getByRole('button', { name: 'Review service items' }).click();

    // Payment Request detail page
    await page.waitForURL(`**/payment-requests/${tioFlowPage.paymentRequest.id}`);
    await tioFlowPage.waitForLoading();

    await expect(page.getByTestId('ReviewServiceItems')).toBeVisible();

    // Approve the first service item
    await tioFlowPage.approveServiceItem();
    await page.getByText('Next').click();
    await tioFlowPage.slowDown();

    // Approve the second service item
    await tioFlowPage.approveServiceItem();
    await page.getByText('Next').click();
    await tioFlowPage.slowDown();

    // Approve the shuttling service item

    // Confirm TIO can view the calculations
    await page.getByText('Show calculations').click();
    await expect(page.locator('[data-testid="ServiceItemCalculations"]')).toContainText('Calculations');
    await expect(page.locator('[data-testid="ServiceItemCalculations"]')).toContainText('Total:');
    await expect(page.locator('[data-testid="ServiceItemCalculations"]')).toContainText('Service schedule: 2');

    // Confirm TIO can hide the calculations. This ensures there's
    // no scrolling weirdness before the next action
    await page.getByText('Hide calculations').click();

    await tioFlowPage.approveServiceItem();
    await page.getByText('Next').click();
    await tioFlowPage.slowDown();

    // Approve the second service item
    await tioFlowPage.approveServiceItem();
    await page.getByText('Next').click();
    await tioFlowPage.slowDown();

    // Approve the crating service item

    // Confirm TIO can view the calculations
    await page.getByText('Show calculations').click();
    await expect(page.locator('[data-testid="ServiceItemCalculations"]')).toContainText('Calculations');
    await expect(page.locator('[data-testid="ServiceItemCalculations"]')).toContainText('Total:');
    await expect(page.locator('[data-testid="ServiceItemCalculations"]')).toContainText('Dimensions: 12x3x10 in');

    // Confirm TIO can hide the calculations. This ensures there's no scrolling weirdness before the next action
    await page.getByText('Hide calculations').click();

    await tioFlowPage.approveServiceItem();
    await page.getByText('Next').click();
    await tioFlowPage.slowDown();

    // Reject the last
    await tioFlowPage.rejectServiceItem();
    await page.getByText('Next').click();
    await tioFlowPage.slowDown();

    await expect(page.getByText('needs your review')).toHaveCount(0, { timeout: 10000 });
    // Complete Request
    await page.getByText('Complete request').click();

    await expect(page.locator('[data-testid="requested"]')).toContainText('$1,130.21');
    await expect(page.locator('[data-testid="accepted"]')).toContainText('$130.22');
    await expect(page.locator('[data-testid="rejected"]')).toContainText('$999.99');

    await page.getByText('Authorize payment').click();
    await tioFlowPage.waitForLoading();

    await tioFlowPage.slowDown();

    // Returns to payment requests overview for move
    expect(page.url()).toContain('/payment-requests');

    await expect(page.getByTestId('tag')).toBeVisible();

    await expect(page.getByTestId('tag').getByText('Reviewed')).toHaveCount(1);

    // ensure the payment request we approved no longer has the
    // "Review Service Items" button
    await expect(page.getByText('Review Service Items')).toHaveCount(0);

    // Go back to queue
    await page.locator('a[title="Home"]').click();
    await tioFlowPage.waitForLoading();

    // search for the moveLocator in case this move doesn't show up
    // on the first page
    await page.locator('#locator').fill(tioFlowPage.moveLocator);
    await page.locator('#locator').blur();
    const paymentSection = page.locator(`[data-uuid="${tioFlowPage.paymentRequest.id}"]`);
    // the payment request that is now in the "Reviewed" status will no longer appear
    // in the TIO queue - only "Payment requested" moves will appear
    await expect(paymentSection.locator('td', { hasText: 'Reviewed' })).not.toBeVisible();
  });

  // This is a stripped down version of the above to build tests on for MB-7936
  test('can review a payment request', async ({ page, officePage }) => {
    // The test starts on the move page

    // Payment Requests page
    await page.getByText('Review service items').last().click();
    await expect(page.getByRole('heading', { name: 'Review service items' })).toBeVisible();

    // Approve the first service item
    await tioFlowPage.approveServiceItem();
    await page.getByText('Next').click();
    await tioFlowPage.slowDown();

    // Reject the second
    await tioFlowPage.rejectServiceItem();
    await page.getByText('Next').click();
    await tioFlowPage.slowDown();

    // Complete Request
    await page.getByText('Complete request').click();

    await page.getByText('Authorize payment').click();
    await officePage.waitForLoading();

    // Returns to payment requests overview for move
    await expect(page.getByTestId('MovePaymentRequests')).toBeVisible();

    // Go back to queue
    await page.locator('a[title="Home"]').click();
    await expect(page.getByRole('heading', { name: 'Payment requests' })).toBeVisible();
  });

  test('can flag and unflag the move for review', async ({ page }) => {
    // Payment Requests page
    expect(page.url()).toContain('/payment-requests');
    await expect(page.getByTestId('MovePaymentRequests')).toBeVisible();

    // click to trigger financial review modal
    await page.getByText('Flag move for financial review').click();

    // Enter information in modal and submit
    await page.locator('label').getByText('Yes').click();
    await page.locator('textarea').fill('Something is rotten in the state of Denmark');

    // Click save on the modal
    await page.getByRole('button', { name: 'Save' }).click();

    // Verify sucess alert and tag
    await expect(page.getByText('Move flagged for financial review.')).toBeVisible();
    await expect(page.getByText('Flagged for financial review', { exact: true })).toBeVisible();

    // Payment Requests page
    expect(page.url()).toContain('/payment-requests');
    await expect(page.getByTestId('MovePaymentRequests')).toBeVisible();

    // click to trigger financial review modal
    await page.getByText('Edit').click();

    // Enter information in modal and submit
    await page.locator('label').getByText('No').click();

    // Click save on the modal
    await page.getByRole('button', { name: 'Save' }).click();

    // Verify sucess alert and tag
    await expect(page.getByText('Move unflagged for financial review.')).toBeVisible();
  });
});
