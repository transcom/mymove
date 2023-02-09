/**
 * Semi-automated converted from a cypress test, and thus may contain
 * non best-practices, in particular: heavy use of `page.locator`
 * instead of `page.getBy*`.
 */

// @ts-check
const { test, expect } = require('./scPpmTestFixture');

test.describe('Services counselor user', () => {
  test.beforeEach(async ({ scPpmPage }) => {
    const move = await scPpmPage.testHarness.buildSubmittedMoveWithPPMShipmentForSC();
    await scPpmPage.navigateToMove(move.locator);
  });

  test('is able to click on move and submit after using the move code filter', async ({ page }) => {
    // Shipment card
    await page.locator('[data-testid="ShipmentContainer"] ').click();

    await expect(page.locator('[data-testid="expectedDepartureDate"]')).toContainText('15 Mar 2020');
    await expect(page.locator('[data-testid="originZIP"]')).toContainText('90210');
    await expect(page.locator('[data-testid="destinationZIP"]')).toContainText('30813');
    await expect(page.locator('[data-testid="sitPlanned"]')).toContainText('No');
    await expect(page.locator('[data-testid="estimatedWeight"]')).toContainText('4,000 lbs');
    await expect(page.locator('[data-testid="hasRequestedAdvance"]')).toContainText('Yes, $5,987');
    await expect(page.locator('[data-testid="counselorRemarks"]')).toContainText('—');
  });
});
