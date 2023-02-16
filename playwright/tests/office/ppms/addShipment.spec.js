/**
 * Semi-automated converted from a cypress test, and thus may contain
 * non best-practices, in particular: heavy use of `page.locator`
 * instead of `page.getBy*`.
 */

// @ts-check
import { test, expect } from './scPpmTestFixture';

test.describe('Services counselor user', () => {
  test.beforeEach(async ({ scPpmPage }) => {
    const move = await scPpmPage.testHarness.buildSubmittedMoveWithPPMShipmentForSC();
    await scPpmPage.navigateToMove(move.locator);
  });

  test('is able to add a new PPM shipment', async ({ page, scPpmPage }) => {
    // Delete existing shipment
    await page.locator('[data-testid="ShipmentContainer"] .usa-button').click();
    await page.locator('[data-testid="grid"] button').getByText('Delete shipment').click();
    await expect(page.locator('[data-testid="modal"]')).toBeVisible();

    await page.locator('[data-testid="modal"] button').getByText('Delete shipment').click();

    await expect(page.locator('[data-testid="ShipmentContainer"] .usa-button')).not.toBeVisible();

    // Click add shipment button and select PPM
    await page.locator('[data-testid="dropdown"]').first().selectOption({ label: 'PPM' });

    // Fill out page one
    await scPpmPage.fillOutOriginInfo();
    await scPpmPage.fillOutDestinationInfo();
    await scPpmPage.fillOutWeight({ hasProGear: true });
    await scPpmPage.selectDutyLocation('JPPSO NORTHWEST', 'closeoutOffice');

    await page.locator('[data-testid="submitForm"]').click();
    await scPpmPage.waitForLoading();

    // Fill out page two
    await expect(page.getByText('Incentive & advance')).toBeVisible();
    await scPpmPage.fillOutIncentiveAndAdvance({ advance: '5987' });
    await page.locator('[data-testid="counselor-remarks"]').type('The requested advance amount has been added.');
    await page.locator('[data-testid="counselor-remarks"]').blur();

    await page.locator('[data-testid="submitForm"]').click();
    await scPpmPage.waitForLoading();

    // Confirm new shipment is visible
    await expect(page.locator('[data-testid="ShipmentContainer"]')).toBeVisible();
    const shipmentContainer = page.locator('[data-testid="ShipmentContainer"]');
    // Verify unexpanded view
    await expect(shipmentContainer.locator('[data-testid="expectedDepartureDate"]')).toContainText('09 Jun 2022');
    await expect(shipmentContainer.locator('[data-testid="originZIP"]')).toContainText('90210');
    await expect(shipmentContainer.locator('[data-testid="destinationZIP"]')).toContainText('76127');
    await expect(shipmentContainer.locator('[data-testid="sitPlanned"]')).toContainText('No');
    await expect(shipmentContainer.locator('[data-testid="estimatedWeight"]')).toContainText('4,000 lbs');
    await expect(shipmentContainer.locator('[data-testid="hasRequestedAdvance"]')).toContainText('Yes, $5,987');
    await expect(shipmentContainer.locator('[data-testid="counselorRemarks"]')).toContainText(
      'The requested advance amount has been added.',
    );

    // Verify fields in the expanded view
    await shipmentContainer.locator('[data-prefix="fas"][data-icon="chevron-down"]').click();
    await expect(shipmentContainer.locator('[data-testid="secondOriginZIP"]')).toContainText('07003');
    await expect(shipmentContainer.locator('[data-testid="secondDestinationZIP"]')).toContainText('08540');
    await expect(shipmentContainer.locator('[data-testid="proGearWeight"]')).toContainText('Yes, 1,000 lbs');
    await expect(shipmentContainer.locator('[data-testid="spouseProGear"]')).toContainText('Yes, 500 lbs');
    await expect(shipmentContainer.locator('[data-testid="estimatedIncentive"]')).toContainText('$67,692');
  });
});
