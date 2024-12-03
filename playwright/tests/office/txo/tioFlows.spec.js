/**
 * Semi-automated converted from a cypress test, and thus may contain
 * non best-practices, in particular: heavy use of `page.locator`
 * instead of `page.getBy*`.
 */

// @ts-check
import { test, expect, OfficePage } from '../../utils/office/officeTest';
import findOptionWithinOpenedDropdown from '../../utils/playwrightUtility';

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

const TIOTabsTitles = ['Payment Request Queue', 'Search'];

const SearchRBSelection = ['Move Code', 'DOD ID', 'Customer Name'];

const SearchTerms = ['SITEXT', '8796353598', 'Spacemen'];

const StatusFilterOptions = ['Draft', 'New Move', 'Needs Counseling', 'Service counseling completed', 'Move approved'];

test.describe('TIO user', () => {
  /** @type {TioFlowPage} */
  let tioFlowPage;
  let testMove;
  test.describe('with Payment Requests Queue', () => {
    test.beforeEach(async ({ officePage }) => {
      await officePage.signInAsNewTIOUser();
    });
  });
  test.describe('with Search Queue', () => {
    test.beforeEach(async ({ officePage }) => {
      testMove = await officePage.testHarness.buildHHGMoveWithServiceItemsandPaymentRequestsForTIO();
      await officePage.signInAsNewTIOUser();

      tioFlowPage = new TioFlowPage(officePage, testMove);

      const searchTab = officePage.page.getByTitle(TIOTabsTitles[1]);
      await searchTab.click();
    });

    test('can search for moves using Move Code', async ({ page }) => {
      const selectedRadio = page.getByRole('group').locator(`label:text("${SearchRBSelection[0]}")`);
      await selectedRadio.click();
      await page.getByTestId('searchText').fill(testMove.locator);
      await page.getByTestId('searchTextSubmit').click();

      await expect(page.getByText('Results')).toBeVisible();
      await expect(page.getByTestId('locator-0')).toContainText(testMove.locator);
    });
    test('can search for moves using DOD ID', async ({ page }) => {
      const selectedRadio = page.getByRole('group').locator(`label:text("${SearchRBSelection[1]}")`);
      await selectedRadio.click();
      await page.getByTestId('searchText').fill(testMove.Orders.ServiceMember.edipi);
      await page.getByTestId('searchTextSubmit').click();

      await expect(page.getByText('Results')).toBeVisible();
      await expect(page.getByTestId('edipi-0')).toContainText(testMove.Orders.ServiceMember.edipi);
    });
    test('can search for moves using Customer Name', async ({ page }) => {
      const CustomerName = `${testMove.Orders.ServiceMember.last_name}, ${testMove.Orders.ServiceMember.first_name}`;
      const selectedRadio = page.getByRole('group').locator(`label:text("${SearchRBSelection[2]}")`);
      await selectedRadio.click();
      await page.getByTestId('searchText').fill(CustomerName);
      await page.getByTestId('searchTextSubmit').click();

      await expect(page.getByText('Results')).toBeVisible();
      await expect(page.getByTestId('customerName-0')).toContainText(CustomerName);
    });
    test('Can filter status using Payment Request Status', async ({ page }) => {
      const selectedRadio = page.getByRole('group').locator(`label:text("${SearchRBSelection[0]}")`);
      await selectedRadio.click();
      await page.getByTestId('searchText').fill(SearchTerms[0]);
      await page.getByTestId('searchTextSubmit').click();

      // Check if Payment Request Status options are present
      const StatusFilter = page.getByTestId('MultiSelectCheckBoxFilter');
      await StatusFilter.click();

      for (const item of StatusFilterOptions) {
        const found = page
          .locator('[id^="react-select"][id*="listbox"]')
          .locator(`[id*="option"]:has(:text("${item}"))`);
        await expect(found).toBeVisible();
      }
    });
    test('Can select a filter status using Payment Request', async ({ page }) => {
      const selectedRadio = page.getByRole('group').locator(`label:text("${SearchRBSelection[0]}")`);
      await selectedRadio.click();
      await page.getByTestId('searchText').fill(testMove.locator);
      await page.getByTestId('searchTextSubmit').click();

      // Check if Payment Request Status options are present
      const StatusFilter = page.getByTestId('MultiSelectCheckBoxFilter');
      await StatusFilter.click();

      const found = findOptionWithinOpenedDropdown(page, StatusFilterOptions[1]);
      await found.click();
      await expect(page.getByText('Results')).toBeVisible();
    });
    test('cant search for empty move code', async ({ page }) => {
      const selectedRadio = page.getByRole('group').locator(`label:text("${SearchRBSelection[0]}")`);
      await selectedRadio.click();

      const SearchBox = page.getByTestId('searchText');
      await SearchBox.fill('');
      await SearchBox.blur();

      await expect(page.getByText('Move Code Must be exactly 6 characters')).toBeVisible();
      await expect(page.getByRole('table')).not.toBeVisible();
    });
    test('cant search for short move code', async ({ page }) => {
      const selectedRadio = page.getByRole('group').locator(`label:text("${SearchRBSelection[0]}")`);
      await selectedRadio.click();

      const SearchBox = page.getByTestId('searchText');
      await SearchBox.fill('MOVE');
      await SearchBox.blur();

      await expect(page.getByText('Move Code Must be exactly 6 characters')).toBeVisible();
      await expect(page.getByRole('table')).not.toBeVisible();
    });
    test('cant search for long move code', async ({ page }) => {
      const selectedRadio = page.getByRole('group').locator(`label:text("${SearchRBSelection[0]}")`);
      await selectedRadio.click();

      const SearchBox = page.getByTestId('searchText');
      await SearchBox.fill('ASUPERLONGMOVE');
      await SearchBox.blur();

      await expect(page.getByText('Move Code Must be exactly 6 characters')).toBeVisible();
      await expect(page.getByRole('table')).not.toBeVisible();
    });
    test('cant search for empty DOD ID', async ({ page }) => {
      const selectedRadio = page.getByRole('group').locator(`label:text("${SearchRBSelection[1]}")`);
      await selectedRadio.click();

      const SearchBox = page.getByTestId('searchText');
      await SearchBox.fill('');
      await SearchBox.blur();

      await expect(page.getByText('DOD ID must be exactly 10 characters')).toBeVisible();
      await expect(page.getByRole('table')).not.toBeVisible();
    });
    test('cant search for short DOD ID', async ({ page }) => {
      const selectedRadio = page.getByRole('group').locator(`label:text("${SearchRBSelection[1]}")`);
      await selectedRadio.click();

      const SearchBox = page.getByTestId('searchText');
      await SearchBox.fill('1234567');
      await SearchBox.blur();

      await expect(page.getByText('DOD ID must be exactly 10 characters')).toBeVisible();
      await expect(page.getByRole('table')).not.toBeVisible();
    });
    test('cant search for long DOD ID', async ({ page }) => {
      const selectedRadio = page.getByRole('group').locator(`label:text("${SearchRBSelection[1]}")`);
      await selectedRadio.click();

      const SearchBox = page.getByTestId('searchText');
      await SearchBox.fill('123456789011');
      await SearchBox.blur();

      await expect(page.getByText('DOD ID must be exactly 10 characters')).toBeVisible();
      await expect(page.getByRole('table')).not.toBeVisible();
    });
    test('cant search for empty Customer Name', async ({ page }) => {
      const selectedRadio = page.getByRole('group').locator(`label:text("${SearchRBSelection[2]}")`);
      await selectedRadio.click();

      const SearchBox = page.getByTestId('searchText');
      await SearchBox.fill('');
      await SearchBox.blur();

      await expect(page.getByText('Customer search must contain a value')).toBeVisible();
      await expect(page.getByRole('table')).not.toBeVisible();
    });
  });
  test.describe('with HHG Moves', () => {
    test.beforeEach(async ({ officePage }) => {
      const move = await officePage.testHarness.buildHHGMoveWithServiceItemsandPaymentRequestsForTIO();
      await officePage.signInAsNewTIOUser();

      tioFlowPage = new TioFlowPage(officePage, move);
      await tioFlowPage.waitForLoading();
      await officePage.tioNavigateToMove(tioFlowPage.moveLocator);
      await officePage.page.getByRole('heading', { name: 'Payment Requests', exact: true }).waitFor();
    });

    test('can use a payment request page to update orders and review a payment request', async ({ page }) => {
      // this is a long test
      test.slow();
      // Payment Requests page
      expect(page.url()).toContain('/payment-requests');
      await expect(page.getByTestId('MovePaymentRequests')).toBeVisible();

      // View Orders page for first payment request
      //
      // FLAKY: the payment requests can show up in any order, so the
      // first payment request on the page is not necessarily the
      // first payment request by sequence number
      //
      // I think this is a bug in the payments request page, but maybe
      // not. Regardless, to reduce flakiness, always pick the payment
      // request with the first sequence number -- ahobson 2023-01-01

      const prNumber = tioFlowPage.paymentRequest.payment_request_number;
      const prHeading = page.getByRole('heading', { name: `Payment Request ${prNumber}` });
      await expect(prHeading).toBeVisible();
      const prCard = prHeading.locator('../../..');

      await prCard.getByRole('link', { name: 'View orders' }).click();
      await tioFlowPage.waitForLoading();

      const form = page.locator('form');
      // Viewing allowances and then orders still navigates back to the payment request details page
      await form.getByRole('link', { name: 'View Allowances' }).click();
      await tioFlowPage.waitForLoading();

      await form.getByRole('link', { name: 'View orders' }).click();
      await tioFlowPage.waitForLoading();

      // Check for link that allows TIO to download the PDF for copy/paste functionality
      await expect(page.locator('p[class*="DocumentViewer_downloadLink"] > a > span')).toHaveText('Download file');

      // Edit orders page | Make edits
      await form.locator('input[name="tac"]').clear();
      await form.locator('input[name="tac"]').fill('E15A');
      await form.locator('input[name="sac"]').clear();
      await form.locator('input[name="sac"]').fill('4K988AS098F');
      await form.locator('input[name="sac"]').blur();
      // Edit orders page | Save
      await page.getByRole('button', { name: 'Save' }).click();
      await page.waitForURL('**/payment-requests');
      await tioFlowPage.waitForLoading();

      // confirm the move is set up as expected
      await expect(page.getByRole('heading', { name: '$1,130.21' })).toBeVisible();
      await expect(page.getByRole('heading', { name: '$805.55' })).toBeVisible();

      // again, find the first payment request by sequence number and
      // operate on that
      await prCard.getByRole('button', { name: 'Review service items' }).click();

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

      await expect(prCard.getByTestId('tag')).toBeVisible();

      await expect(prCard.getByTestId('tag').getByText('Reviewed')).toHaveCount(1);

      // ensure the payment request we approved no longer has the
      // "Review Service Items" button
      await expect(prCard.getByText('Review Service Items')).toHaveCount(0);

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
      //   await page.locator('input[data-testid="ntsTacInput"]').click().fill('E19A');
      //   await page.locator('input[data-testid="ntsSacInput"]').click().fill('3L988AS098F');
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

      await page.getByText('Review service items').first().click();

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

    test('is able to view Origin GBLOC', async ({ page }) => {
      // Check for Origin GBLOC label
      await expect(page.getByTestId('originGBLOC')).toHaveText('Origin GBLOC');
      await expect(page.getByTestId('infoBlock')).toContainText('KKFA');
    });
  });

  test.describe('with NTSR moves without service items', () => {
    test.beforeEach(async ({ officePage }) => {
      const move = await officePage.testHarness.buildNTSRMoveWithPaymentRequest();
      await officePage.signInAsNewTIOUser();

      tioFlowPage = new TioFlowPage(officePage, move);
      await tioFlowPage.waitForLoading();
      await officePage.tioNavigateToMove(tioFlowPage.moveLocator);
      await officePage.page.getByRole('heading', { name: 'Payment Requests', exact: true }).waitFor();
    });

    test('can review a NTS-R', async ({ page, officePage }) => {
      // Payment Requests page
      expect(page.url()).toContain('/payment-requests');
      await expect(page.getByTestId('MovePaymentRequests')).toBeVisible();

      // Verify weight info
      const weightSection = '#billable-weights';
      await expect(page.locator(weightSection)).toContainText('Billable weights');

      // Verify Actual Billable Weight info
      const totalBillableWeightParent = page
        .getByRole('heading', { name: 'Actual billable weight', exact: true })
        .locator('..');
      await expect(totalBillableWeightParent.getByRole('heading', { name: '2,000 lbs', exact: true })).toBeVisible();

      // Verify Maximum billable weight info
      const maximumBillableWeightParent = page
        .getByRole('heading', { name: 'Maximum billable weight', exact: true })
        .locator('..');
      await expect(maximumBillableWeightParent.getByRole('heading', { name: '8,000 lbs', exact: true })).toBeVisible();

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
      await expect(page.locator(prSection)).toContainText('HTC111-11-1-1112');
      await expect(page.locator(prSection)).toContainText('Non-temp storage release');
      await expect(page.getByTestId('pickup-to-destination')).toBeVisible();
      await expect(page.locator(prSection)).toContainText('1111 (HHG)');

      // Verify Service Item
      await expect(page.getByTestId('serviceItemName')).toContainText('Counseling');
      await expect(page.getByTestId('serviceItemAmount')).toContainText('$324.00');

      // Review Weights
      await page.locator(weightSection).getByText('Review shipment weights').click();
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
      const reviewWeightsBtn = page.locator('#billable-weights').getByText('Review shipment weights');
      await tioFlowPage.waitForLoading();
      await reviewWeightsBtn.click();

      await tioFlowPage.waitForLoading();
      await page.locator('[data-testid="closeSidebar"]').click();

      // Start reviewing service items
      await page.getByText('Review service items').click();
      await tioFlowPage.waitForLoading();

      await tioFlowPage.rejectServiceItem();
      await page.getByText('Next').click();
      await tioFlowPage.slowDown();

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
      await tioFlowPage.waitForLoading();
      await officePage.tioNavigateToMove(tioFlowPage.moveLocator);
      await officePage.page.getByRole('heading', { name: 'Payment Requests', exact: true }).waitFor();
    });

    test('can view calculation factors', async ({ page }) => {
      // Payment Requests page
      expect(page.url()).toContain('/payment-requests');
      await expect(page.getByTestId('MovePaymentRequests')).toBeVisible();

      // Review Weights
      await page.locator('#billable-weights').getByText('Review shipment weights').click();
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
