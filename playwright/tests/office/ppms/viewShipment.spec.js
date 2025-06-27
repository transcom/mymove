/**
 * Semi-automated converted from a cypress test, and thus may contain
 * non best-practices, in particular: heavy use of `page.locator`
 * instead of `page.getBy*`.
 */

// @ts-check
import { test, expect } from './ppmTestFixture';

test.describe('Services counselor user', () => {
  test.beforeEach(async ({ ppmPage }) => {
    const move = await ppmPage.testHarness.buildSubmittedMoveWithPPMShipmentForSC();
    await ppmPage.navigateToMove(move.locator);
  });

  test('is able to click on move and submit after using the move code filter', async ({ page }) => {
    // Shipment card
    await page.locator('[data-testid="ShipmentContainer"] ').click();

    await expect(page.locator('[data-testid="expectedDepartureDate"]')).toContainText('15 Mar 2025');

    await expect(page.locator('[data-testid="pickupAddress"]')).toContainText('987 New Street');
    await expect(page.locator('[data-testid="pickupAddress"]')).toContainText('P.O. Box 12345');
    await expect(page.locator('[data-testid="pickupAddress"]')).toContainText(/Des Moines/i);
    await expect(page.locator('[data-testid="pickupAddress"]')).toContainText('IA');
    await expect(page.locator('[data-testid="pickupAddress"]')).toContainText('50309');

    await expect(page.locator('[data-testid="destinationAddress"]')).toContainText('123 New Street');
    await expect(page.locator('[data-testid="destinationAddress"]')).toContainText('P.O. Box 12345');
    await expect(page.locator('[data-testid="destinationAddress"]')).toContainText(/Fort Eisenhower/i);
    await expect(page.locator('[data-testid="destinationAddress"]')).toContainText('GA');
    await expect(page.locator('[data-testid="destinationAddress"]')).toContainText('30813');

    await expect(page.locator('[data-testid="sitPlanned"]')).toContainText('No');
    await expect(page.locator('[data-testid="estimatedWeight"]')).toContainText('4,000 lbs');
    await expect(page.locator('[data-testid="hasRequestedAdvance"]')).toContainText('Yes, $5,987');
    await expect(page.locator('[data-testid="counselorRemarks"]')).toContainText('—');
  });
});
