/**
 * Semi-automated converted from a cypress test, and thus may contain
 * non best-practices, in particular: heavy use of `page.locator`
 * instead of `page.getBy*`.
 */

// @ts-check
import { test, expect } from '../../utils/office/officeTest';

test.describe('TIO user', () => {
  // This test performs a mutation so it can only succeed on a fresh DB.
  test('can update max billable weight', async ({ page, officePage }) => {
    const move = await officePage.testHarness.buildNTSRMoveWithPaymentRequest();
    /** @type {string} */
    const moveLocator = move.locator;
    await officePage.signInAsNewTIOUser();

    // Go to known NTS-R move
    await officePage.tioNavigateToMove(moveLocator);
    await officePage.waitForLoading();

    // Verify we are on the Payment Requests page
    expect(page.url()).toContain('/payment-requests');
    await expect(page.getByTestId('MovePaymentRequests')).toBeVisible();

    // Review Weights
    const reviewWeightsBtn = page.locator('#billable-weights').getByText('Review shipment weights');
    await reviewWeightsBtn.click();

    await page.getByRole('heading', { name: 'Review weights' }).waitFor({ state: 'visible' });

    await page.getByRole('button', { name: 'Edit' }).click();
    await officePage.waitForLoading();

    const fieldset = page.locator('fieldset');
    await fieldset.locator('input#billableWeight').clear();
    await fieldset.locator('input#billableWeight').type('2200');
    await fieldset.locator('textarea#billableWeightJustification').type('Some basic remarks.');

    await page.getByRole('button', { name: 'Save changes' }).click();
    await officePage.waitForLoading();

    await expect(page.locator('[data-testid="billableWeightValue"]')).toContainText('2,200 lbs');
    await expect(page.locator('[data-testid="billableWeightRemarks"]')).toContainText('Some basic remarks.');
  });
});
