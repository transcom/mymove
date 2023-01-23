/**
 * Semi-automated converted from a cypress test, and thus may contain
 * non best-practices, in particular: heavy use of `page.locator`
 * instead of `page.getBy*`.
 */

// @ts-check
const { test, expect, OfficePage } = require('../../utils/officeTest');

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
   * @override
   */
  constructor(officePage, move) {
    super(officePage.page, officePage.request);
    this.move = move;
    this.moveLocator = move.locator;
    this.paymentRequest = this.findPaymentRequestBySequenceNumber(1);
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
      await this.page.locator('textarea[data-testid="textarea"]').type('This is not a valid request');
      await serviceItemCardLocator.getByRole('button', { name: 'Save' }).click();
    } else {
      const inputEl = serviceItemCardLocator.locator('input[data-testid="approveRadio"]');
      const id = await inputEl.getAttribute('id');
      await this.page.locator(`label[for="${id}"]`).click();
    }
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
    await expect(siCalc).toContainText('0.15');
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
}

test.describe('TIO user', () => {
  /** @type {TioFlowPage} */
  let tioFlowPage;

  test.describe('with HHG Moves', () => {
    test.beforeEach(async ({ officePage }) => {
      const move = await officePage.testHarness.buildHHGMoveWithServiceItemsandPaymentRequestsForTIO();
      await officePage.signInAsNewTIOUser();

      tioFlowPage = new TioFlowPage(officePage, move);
      await officePage.tioNavigateToMove(move.locator);
    });

    test('can use a payment request page to update orders and review a payment request', async ({ page }) => {
      // this is a long test
      test.slow();
      // Payment Requests page
      expect(page.url()).toContain('/payment-requests');
      await expect(page.getByTestId('MovePaymentRequests')).toBeVisible();

      // View Orders page for first payment request
      await page.getByText('View orders').first().click();
      await tioFlowPage.waitForLoading();

      const form = page.locator('form');
      await form.locator('input[name="tac"]').clear();
      await form.locator('input[name="tac"]').type('E15A');
      await form.locator('input[name="sac"]').clear();
      await form.locator('input[name="sac"]').type('4K988AS098F');
      await form.locator('input[name="sac"]').blur();
      // Edit orders page | Save
      await Promise.all([page.waitForNavigation(), page.getByRole('button', { name: 'Save' }).click()]);
      await tioFlowPage.waitForLoading();

      expect(page.url()).toContain('/details');

      await Promise.all([
        page.waitForNavigation(),
        page.getByRole('link', { name: 'Payment requests', exact: true }).click(),
      ]);
      await tioFlowPage.waitForLoading();

      // Payment Requests page
      expect(page.url()).toContain('/payment-requests');

      // confirm the move is set up as expected
      await expect(page.getByRole('heading', { name: '$1,130.21' })).toBeVisible();
      await expect(page.getByRole('heading', { name: '$805.55' })).toBeVisible();

      // first payment request
      await Promise.all([page.waitForNavigation(), page.getByText('Review service items').first().click()]);
      await tioFlowPage.waitForLoading();

      // // Payment Request detail page
      expect(page.url()).toContain(`/payment-requests/${tioFlowPage.paymentRequest.id}`);
      await expect(page.getByTestId('ReviewServiceItems')).toBeVisible();

      // Approve the first service item
      await tioFlowPage.approveServiceItem();
      await page.getByText('Next').click();

      // Approve the second service item
      await tioFlowPage.approveServiceItem();
      await page.getByText('Next').click();
      await tioFlowPage.waitForLoading();

      // Approve the shuttling service item

      // Confirm TIO can view the calculations
      await page.getByText('Show calculations').click();
      await expect(page.locator('[data-testid="ServiceItemCalculations"]')).toContainText('Calculations');
      await expect(page.locator('[data-testid="ServiceItemCalculations"]')).toContainText('Total amount requested');
      await expect(page.locator('[data-testid="ServiceItemCalculations"]')).toContainText('Service schedule: 2');

      // Confirm TIO can hide the calculations. This ensures there's
      // no scrolling weirdness before the next action
      await page.getByText('Hide calculations').click();

      await tioFlowPage.approveServiceItem();
      await page.getByText('Next').click();

      // Approve the second service item
      await tioFlowPage.approveServiceItem();
      await page.getByText('Next').click();

      // Approve the crating service item

      // Confirm TIO can view the calculations
      await page.getByText('Show calculations').click();
      await expect(page.locator('[data-testid="ServiceItemCalculations"]')).toContainText('Calculations');
      await expect(page.locator('[data-testid="ServiceItemCalculations"]')).toContainText('Total amount requested');
      await expect(page.locator('[data-testid="ServiceItemCalculations"]')).toContainText('Dimensions: 12x3x10 in');

      // Confirm TIO can hide the calculations. This ensures there's no scrolling weirdness before the next action
      await page.getByText('Hide calculations').click();

      await tioFlowPage.approveServiceItem();
      await page.getByText('Next').click();

      // Reject the last
      await tioFlowPage.rejectServiceItem();
      await page.getByText('Next').click();

      await expect(page.getByText('item still needs your review')).not.toBeVisible();
      // Complete Request
      await page.getByText('Complete request').click();

      await expect(page.locator('[data-testid="requested"]')).toContainText('$1,130.21');
      await expect(page.locator('[data-testid="accepted"]')).toContainText('$130.22');
      await expect(page.locator('[data-testid="rejected"]')).toContainText('$999.99');

      await page.getByText('Authorize payment').click();
      await tioFlowPage.waitForLoading();

      await page.waitForLoadState('networkidle', { timeout: 10000 });

      // Returns to payment requests overview for move
      expect(page.url()).toContain('/payment-requests');
      await expect(page.getByTestId('MovePaymentRequests')).toBeVisible();
      await expect(page.locator('[data-testid="MovePaymentRequests"] [data-testid="tag"]').first()).toBeVisible();

      await expect(
        page.locator('[data-testid="MovePaymentRequests"] [data-testid="tag"]').getByText('Reviewed'),
      ).toHaveCount(1);

      // ensure the payment request we approved no longer has the
      // "Review Service Items" button
      const paymentRequestCard = page.getByText(tioFlowPage.paymentRequest.payment_request_number).locator('../..');
      await expect(paymentRequestCard.getByText('Review Service Items')).toHaveCount(0);

      // Go back to queue
      await page.locator('a[title="Home"]').click();
      await tioFlowPage.waitForLoading();

      // search for the moveLocator in case this move doesn't show up
      // on the first page
      await page.locator('#locator').type(tioFlowPage.moveLocator);
      await page.locator('#locator').blur();
      const paymentSection = page.locator(`[data-uuid="${tioFlowPage.paymentRequest.id}"]`);
      await expect(paymentSection).toHaveCount(1);
      await expect(paymentSection.locator('td', { hasText: 'Reviewed' })).toBeVisible();
    });

    // This is a stripped down version of the above to build tests on for MB-7936
    test('can review a payment request', async ({ page, officePage }) => {
      // The test starts on the move page

      // Payment Requests page
      const secondPaymentRequest = tioFlowPage.findPaymentRequestBySequenceNumber(2);
      const paymentRequestCard = page.getByText(secondPaymentRequest.payment_request_number).locator('../..');
      await Promise.all([page.waitForNavigation(), paymentRequestCard.getByText('Review service items').click()]);
      await officePage.waitForLoading();

      // Approve the first service item
      await tioFlowPage.approveServiceItem();
      await page.getByText('Next').click();

      // Reject the second
      await tioFlowPage.rejectServiceItem();
      await page.getByText('Next').click();

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
      await page.locator('textarea').type('Something is rotten in the state of Denmark');

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

    /**
     * This test is being temporarily skipped until flakiness issues
     * can be resolved. It was skipped in cypress and is not part of
     * the initial playwright conversion. - ahobson 2023-01-05
     */
    test.skip('can add/edit TAC/SAC', async ({ page }) => {
      // Payment Requests page
      expect(page.url()).toContain('/payment-requests');
      await expect(page.getByTestId('MovePaymentRequests')).toBeVisible();

      // await expect(page.locator('button')).toContainText('Edit').click();
      // await expect(page.locator('button')).toContainText('Add or edit codes').click();
      // cy.url().should('include', `/moves/NTSTIO/orders`);

      // await page.locator('form').within(() => {
      //   await page.locator('input[data-testid="ntsTacInput"]').click().type('E19A');
      //   await page.locator('input[data-testid="ntsSacInput"]').click().type('3L988AS098F');
      //   // Edit orders page | Save
      //   await expect(page.locator('button')).toContainText('Save').click();
      // });
      // cy.url().should('include', `/moves/NTSTIO/details`);
      // await expect(page.getByText('Payment requests').click()).toBeVisible();
      // cy.url().should('include', `/payment-requests`);
      // await expect(page.locator('button')).toContainText('Edit').click();

      // await page.locator('input#tacType-NTS').click({ force: true });
      // await page.locator('input#sacType-NTS').click({ force: true });
      // await page.locator('button[type="submit"]').click();

      // await expect(page.locator('[data-testid="tac"]')).toContainText('E19A (NTS)');
      // await expect(page.locator('[data-testid="sac"]')).toContainText('3L988AS098F (NTS)');
    });

    // ahobson - 2023-01-05 skipping this test as it is a subset of
    // the test called 'can use a payment request page to update
    // orders and review a payment request'
    test.skip('can view and approve service items', async ({ page }) => {
      // Payment Requests page
      expect(page.url()).toContain('/payment-requests');
      await expect(page.getByTestId('MovePaymentRequests')).toBeVisible();

      await Promise.all([page.waitForNavigation(), page.getByText('Review service items').first().click()]);

      // await expect(page.locator('[data-testid="serviceItemName"]')).toContainText('Move management');
      // await page.locator('[data-testid="approveRadio"]').click({ force: true });
      // cy.wait('@patchPaymentServiceItemStatus');
      // await expect(page.locator('button')).toContainText('Next').click();

      // await expect(page.locator('[data-testid="serviceItemName"]')).toContainText('Domestic origin shuttle service');
      // await page.locator('[data-testid="approveRadio"]').click({ force: true });
      // cy.wait('@patchPaymentServiceItemStatus');
      // await expect(page.locator('button')).toContainText('Next').click();

      // await expect(page.locator('[data-testid="serviceItemName"]')).toContainText('Domestic origin shuttle service');
      // await page.locator('[data-testid="approveRadio"]').click({ force: true });
      // cy.wait('@patchPaymentServiceItemStatus');
      // await expect(page.locator('button')).toContainText('Next').click();

      // await expect(page.locator('[data-testid="serviceItemName"]')).toContainText('Domestic crating');
      // await page.locator('[data-testid="approveRadio"]').click({ force: true });
      // cy.wait('@patchPaymentServiceItemStatus');
      // await expect(page.locator('button')).toContainText('Next').click();

      // await expect(page.locator('[data-testid="serviceItemName"]')).toContainText('Domestic crating');
      // await page.locator('[data-testid="approveRadio"]').click({ force: true });
      // cy.wait('@patchPaymentServiceItemStatus');
      // await expect(page.locator('button')).toContainText('Next').click();

      // await expect(page.locator('[data-testid="serviceItemName"]')).toContainText('Domestic linehaul');
      // await page.locator('[data-testid="approveRadio"]').click({ force: true });
      // cy.wait('@patchPaymentServiceItemStatus');
      // await expect(page.locator('button')).toContainText('Next').click();

      // await expect(page.locator('[data-testid="accepted"]')).toContainText('$1,130.21');
      // await expect(page.locator('button')).toContainText('Authorize payment').click();
      // cy.wait(['@getMovePaymentRequests']);

      // await expect(page.locator('[data-testid="tag"]')).toContainText('Reviewed');
    });
  });

  test.describe('with NTSR moves without service items', () => {
    test.beforeEach(async ({ officePage }) => {
      const move = await officePage.testHarness.buildNTSRMoveWithPaymentRequest();
      await officePage.signInAsNewTIOUser();

      tioFlowPage = new TioFlowPage(officePage, move);
      await officePage.tioNavigateToMove(move.locator);
    });

    test('can review a NTS-R', async ({ page, officePage }) => {
      // Payment Requests page
      expect(page.url()).toContain('/payment-requests');
      await expect(page.getByTestId('MovePaymentRequests')).toBeVisible();

      // Verify weight info
      const weightSection = '#billable-weights';
      await expect(page.locator(weightSection)).toContainText('Billable weights');
      await expect(page.locator(weightSection)).toContainText('8,000 lbs');
      await expect(page.locator(weightSection)).toContainText('2,000 lbs');

      // Verify External Shipment shown
      await expect(page.locator(weightSection)).toContainText('1 other shipment:');
      await expect(page.locator(weightSection)).toContainText('0 lbs');
      await expect(page.locator(weightSection).locator('a', { hasText: 'View move details' })).toBeVisible();

      // Verify relevant payment request info
      const prSection = '#payment-requests';
      await expect(page.locator(prSection)).toContainText('Needs review');
      await expect(page.locator(prSection)).not.toContainText('Reviewed');
      await expect(page.locator(prSection)).not.toContainText('Rejected');

      await expect(page.locator(prSection)).toContainText('$324.00');
      await expect(page.locator(prSection)).toContainText('HTC111-11-1-1111');
      await expect(page.locator(prSection)).toContainText('Non-temp storage release');
      await expect(page.getByTestId('pickup-to-destination')).toBeVisible();
      await expect(page.locator(prSection)).toContainText('1111 (HHG)');

      // Verify Service Item
      await expect(page.getByTestId('serviceItemName')).toContainText('Counseling');
      await expect(page.getByTestId('serviceItemAmount')).toContainText('$324.00');

      // Review Weights
      await page.locator(weightSection).getByText('Review weights').click();
      await officePage.waitForLoading();
      await expect(page.getByText('Edit max billable weight')).toBeVisible();
      await page.locator('[data-testid="closeSidebar"]').click();

      // Review service items
      await page.getByText('Review service items').click();
      await officePage.waitForLoading();

      // Approve the service item
      await tioFlowPage.approveCurrentServiceItem();

      // Complete Request
      await expect(page.getByText('Complete request')).toBeVisible();

      await page.getByText('Authorize payment').click();
      await officePage.waitForLoading();

      // Returns to payment requests overview for move
      await expect(page.getByTestId('MovePaymentRequests')).toBeVisible();
      await expect(page.getByTestId('serviceItemStatus').getByText('Accepted')).toBeVisible();

      // Should now have 'Reviewed' Tag
      await expect(page.locator(prSection)).toContainText('Reviewed');
      await expect(page.locator(prSection)).not.toContainText('Needs Review');
      await expect(page.locator(prSection)).not.toContainText('Rejected');

      // Go back home
      await page.locator('a[title="Home"]').click();
      await expect(page.getByRole('heading', { name: 'Payment requests' })).toBeVisible();
    });

    test('can reject a NTS-R', async ({ page }) => {
      // Payment Requests page
      expect(page.url()).toContain('/payment-requests');
      await expect(page.getByTestId('MovePaymentRequests')).toBeVisible();

      // Verify payment request status
      const prSection = '#payment-requests';
      await expect(page.locator(prSection)).toContainText('Needs review');
      await expect(page.locator(prSection)).not.toContainText('Reviewed');
      await expect(page.locator(prSection)).not.toContainText('Rejected');

      // Review Weights
      await page.locator('#billable-weights').getByText('Review weights').click();
      await tioFlowPage.waitForLoading();
      await page.locator('[data-testid="closeSidebar"]').click();

      // Start reviewing service items
      await page.getByText('Review service items').click();
      await tioFlowPage.waitForLoading();

      await tioFlowPage.rejectServiceItem();
      await page.getByText('Next').click();

      // Reject the Request
      await expect(page.getByText('Review details')).toBeVisible();
      await page.getByText('Reject request').click();
      await tioFlowPage.waitForLoading();

      // Returns to payment requests overview for move
      await expect(page.getByTestId('MovePaymentRequests')).toBeVisible();

      // Should now have 'Rejected' Tag
      await expect(page.locator(prSection)).toContainText('Rejected');
      await expect(page.locator(prSection)).not.toContainText('Needs Review');
      await expect(page.locator(prSection)).not.toContainText('Reviewed');

      // Go back home
      await page.locator('a[title="Home"]').click();
      await expect(page.getByRole('heading', { name: 'Payment requests' })).toBeVisible();
    });
  });

  test.describe('with NTSR moves with service items and payment request', () => {
    test.beforeEach(async ({ officePage }) => {
      const move = await officePage.testHarness.buildNTSRMoveWithServiceItemsAndPaymentRequest();
      await officePage.signInAsNewTIOUser();

      tioFlowPage = new TioFlowPage(officePage, move);
      await officePage.tioNavigateToMove(move.locator);
    });

    test('can view calculation factors', async ({ page }) => {
      // Payment Requests page
      expect(page.url()).toContain('/payment-requests');
      await expect(page.getByTestId('MovePaymentRequests')).toBeVisible();

      // Review Weights
      await page.locator('#billable-weights').getByText('Review weights').click();
      await tioFlowPage.waitForLoading();
      await page.locator('[data-testid="closeSidebar"]').click();

      // Review service items
      await page.getByText('Review service items').click();
      await tioFlowPage.waitForLoading();

      // Verify at domestic linehaul calculations
      await page.locator('[data-testid=toggleCalculations]').click();
      await tioFlowPage.validateDLCalcValues();
      await tioFlowPage.approveCurrentServiceItem();

      // Verify at fuel surcharge calculations
      await page.locator('[data-testid=toggleCalculations]').click();
      await tioFlowPage.validateFSCalcValues();
      await tioFlowPage.approveCurrentServiceItem();

      // Verify at domestic origin calculations
      await page.locator('[data-testid=toggleCalculations]').click();
      await tioFlowPage.validateDXCalcValues();
      await tioFlowPage.approveCurrentServiceItem();

      // Verify at domestic destination calculations
      await page.locator('[data-testid=toggleCalculations]').click();
      await tioFlowPage.validateDXCalcValues();
      await tioFlowPage.approveCurrentServiceItem();

      // Verify at domestic unpacking calculations
      await page.locator('[data-testid=toggleCalculations]').click();
      await tioFlowPage.validateDUCalcValues();
      await tioFlowPage.approveCurrentServiceItem();

      // Complete Request
      await expect(page.getByText('Complete request')).toBeVisible();
      await page.getByText('Authorize payment').click();
      await tioFlowPage.waitForLoading();

      // Returns to payment requests overview for move
      await expect(page.getByTestId('MovePaymentRequests')).toBeVisible();

      // Verify Service Item Calcs now that approved and visable
      await expect(page.getByTestId('serviceItemStatus').getByText('Accepted')).toHaveCount(5);

      // Verify domestic linehaul cacls on payment request
      await page.getByText('Domestic linehaul').click();
      await tioFlowPage.validateDLCalcValues();
      await page.getByText('Domestic linehaul').click();

      // Verify fuel surcharge calcs on payment request
      await page.getByText('Fuel surcharge', { exact: true }).click();
      await tioFlowPage.validateFSCalcValues();
      await page.getByText('Fuel surcharge', { exact: true }).click();

      // Verify Domestic origin price cacls on payment request
      await page.getByText('Domestic origin price').click();
      await tioFlowPage.validateDXCalcValues();
      await page.getByText('Domestic origin price').click();

      // Verify domestic destination price cacls on payment request
      await page.getByText('Domestic destination price').click();
      await tioFlowPage.validateDXCalcValues();
      await page.getByText('Domestic destination price').click();

      // Verify domestic unpacking cacls on payment request
      await page.getByText('Domestic unpacking').click();
      await tioFlowPage.validateDUCalcValues();
      await page.getByText('Domestic unpacking').click();

      // calcs are good, go home
      await page.locator('a[title="Home"]').click();
      await expect(page.getByRole('heading', { name: 'Payment requests' })).toBeVisible();
    });
  });
});
