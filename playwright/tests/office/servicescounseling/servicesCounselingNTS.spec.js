/**
 * Semi-automated converted from a cypress test, and thus may contain
 * non best-practices, in particular: heavy use of `page.locator`
 * instead of `page.getBy*`.
 */

// @ts-check
const { expect, test } = require('./servicesCounselingTestFixture');

test.describe('Services counselor user', () => {
  test.beforeEach(async ({ scPage }) => {
    const move = await scPage.testHarness.buildHHGMoveWithNTSAndNeedsSC();
    await scPage.navigateToMove(move.locator);
  });

  test('Services Counselor can delete/remove an NTS shipment request', async ({ page, scPage }) => {
    // this test is almost identical to the NTSR test
    await scPage.addNTSShipment();

    // single HHG plus added NTS
    await expect(page.locator('[data-testid="ShipmentContainer"] .usa-button')).toHaveCount(2);

    await page.locator('[data-testid="ShipmentContainer"] .usa-button').last().click();
    await scPage.waitForLoading();

    // click to trigger confirmation modal
    await page.locator('[data-testid="grid"] button').getByText('Delete shipment').click();

    await expect(page.getByTestId('modal')).toBeVisible();

    await page.getByTestId('modal').getByRole('button', { name: 'Delete shipment' }).click();
    await scPage.waitForLoading();

    await expect(page.locator('[data-testid="ShipmentContainer"] .usa-button')).toHaveCount(1);
  });

  test('Services Counselor can enter accounting codes and submit shipment', async ({ page, scPage }) => {
    // this test is almost identical to the NTSR test
    await scPage.addNTSShipment();
    // edit the newly added NTS shipment
    await page.locator('[data-testid="ShipmentContainer"] .usa-button').last().click();
    await scPage.waitForLoading();

    await page.locator('[data-testid="grid"]').getByRole('button', { name: 'Add or edit codes' }).click();

    const form = page.locator('form');
    await form.locator('input[name="tac"]').type('E15A');
    await form.locator('input[name="sac"]').type('4K988AS098F');
    await form.locator('input[name="ntsTac"]').type('F123');
    await form.locator('input[name="ntsSac"]').type('3L988AS098F');
    await form.locator('input[name="ordersNumber"]').type('1234');
    await form.locator('select[name="departmentIndicator"]').selectOption({ label: '21 Army' });
    await form.locator('select[name="ordersType"]').selectOption({ label: 'Permanent Change Of Station (PCS)' });
    await form.locator('select[name="ordersTypeDetail"]').selectOption({ label: 'Shipment of HHG Permitted' });
    // Edit orders page | Save
    await form.getByRole('button', { name: 'Save' }).click();
    await scPage.waitForLoading();

    await expect(page.locator('[data-testid="tacMDC"]')).toContainText('E15A');
    await expect(page.locator('[data-testid="sacSDN"]')).toContainText('4K988AS098F');
    await expect(page.locator('[data-testid="NTStac"]')).toContainText('F123');
    await expect(page.locator('[data-testid="NTSsac"]')).toContainText('3L988AS098F');

    // test 'Services Counselor can assign accounting code(s) to a shipment'
    // combining this test with the one above

    await page.locator('[data-testid="ShipmentContainer"] .usa-button').last().click();
    await scPage.waitForLoading();

    await page.locator('[data-testid="radio"] [for="tacType-NTS"]').click();
    await page.locator('[data-testid="radio"] [for="sacType-HHG"]').click();

    await page.locator('[data-testid="submitForm"]').click();
    await scPage.waitForLoading();

    await expect(page.locator('.usa-alert__text')).toContainText('Your changes were saved.');

    // test 'Services Counselor can submit a move with an NTS shipment'
    // combining test

    // click to trigger confirmation modal
    await page.getByText('Submit move details').click();

    await expect(page.locator('[data-testid="modal"]')).toBeVisible();

    await page.getByRole('button', { name: 'Yes, submit' }).click();

    // verify success alert
    await expect(page.getByText('Move submitted.')).toBeVisible();
  });
});
