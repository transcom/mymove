// @ts-check
const { test, expect } = require('../../utils/officeTest');

test.describe('TIO user', () => {
  // This test performs a mutation so it can only succeed on a fresh DB.
  test('can update max billable weight', async ({ page, officePage }) => {
    const move = await officePage.testHarness.buildNTSRMoveWithPaymentRequest();
    /** @type {string} */
    const moveLocator = move.locator;
    await officePage.signInAsNewTIOUser();

    // Go to known NTS-R move
    await page.locator('#locator').type(moveLocator);
    await page.locator('th[data-testid="locator"]').first().click();
    await page.locator('[data-testid="locator-0"]').click();
    await officePage.waitForLoading();

    // Verify we are on the Payment Requests page
    expect(page.url()).toContain('/payment-requests');
    await expect(page.getByTestId('MovePaymentRequests')).toBeVisible();

    await Promise.all([
      page.waitForNavigation(),
      page.locator('#billable-weights').getByText('Review weights').click(),
    ]);

    await page.getByRole('button', { name: 'Edit' }).click();
    await officePage.waitForLoading();

    const fieldset = page.locator('fieldset');
    await fieldset.locator('input#billableWeight').clear();
    await fieldset.locator('input#billableWeight').type('7400');
    await fieldset.locator('textarea#billableWeightJustification').type('Some basic remarks.');

    await page.getByRole('button', { name: 'Save changes' }).click();
    await officePage.waitForLoading();

    await expect(page.locator('[data-testid="billableWeightValue"]')).toContainText('7,400 lbs');
    await expect(page.locator('[data-testid="billableWeightRemarks"]')).toContainText('Some basic remarks.');
  });
});
