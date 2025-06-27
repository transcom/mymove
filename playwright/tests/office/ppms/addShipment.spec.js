/**
 * Semi-automated converted from a cypress test, and thus may contain
 * non best-practices, in particular: heavy use of `page.locator`
 * instead of `page.getBy*`.
 */

// @ts-check
import { test, expect } from './ppmTestFixture';

/**
 * @param {string} dateString
 */
function formatDate(dateString) {
  const [month, day, year] = dateString.split('/').map(Number);
  const date = new Date(year, month - 1, day);

  const dayFormatted = String(date.getDate()).padStart(2, '0');
  const monthFormatted = date.toLocaleString('default', { month: 'short' });
  const yearFormatted = date.getFullYear();

  return `${dayFormatted} ${monthFormatted} ${yearFormatted}`;
}
const gunSafeEnabled = process.env.FEATURE_FLAG_GUN_SAFE;

test.describe('Services counselor user', () => {
  test.beforeEach(async ({ ppmPage }) => {
    const move = await ppmPage.testHarness.buildSubmittedMoveWithPPMShipmentForSC();
    await ppmPage.navigateToMove(move.locator);
  });

  test('is able to add a new PPM shipment', async ({ page, ppmPage }) => {
    // Delete existing shipment
    await page.locator('[data-testid="ShipmentContainer"] .usa-button').click();
    await page.locator('[data-testid="grid"] button').getByText('Delete shipment').click();
    await expect(page.locator('[data-testid="modal"]')).toBeVisible();

    await page.locator('[data-testid="modal"] button').getByText('Delete shipment').click();

    await expect(page.locator('[data-testid="ShipmentContainer"] .usa-button')).not.toBeVisible();

    // Click add shipment button and select PPM
    await page.locator('[data-testid="dropdown"]').first().selectOption({ label: 'PPM' });

    // Fill out page one
    await ppmPage.fillOutOriginInfo();
    await ppmPage.fillOutDestinationInfo();
    await ppmPage.fillOutWeight({ hasProGear: true });
    await ppmPage.selectDutyLocation('JPPSO NORTHWEST', 'closeoutOffice');

    await page.locator('[data-testid="submitForm"]').click();
    await ppmPage.waitForLoading();

    // Fill out page two
    const selector = page.getByText('Incentive & advance');
    await selector.waitFor({ state: 'visible' });
    await ppmPage.fillOutIncentiveAndAdvance({ advance: '5987' });
    await page.locator('[data-testid="counselor-remarks"]').fill('The requested advance amount has been added.');
    await page.locator('[data-testid="counselor-remarks"]').blur();

    await page.locator('[data-testid="submitForm"]').click();
    await ppmPage.waitForLoading();

    // Confirm new shipment is visible
    await expect(page.locator('[data-testid="ShipmentContainer"]')).toBeVisible();
    const shipmentContainer = page.locator('[data-testid="ShipmentContainer"]');
    // Verify unexpanded view
    const expectedDeparture = new Date(Date.now() + 24 * 60 * 60 * 1000).toLocaleDateString('en-US');
    const formattedDate = formatDate(expectedDeparture);
    await expect(shipmentContainer.locator('[data-testid="expectedDepartureDate"]')).toContainText(formattedDate);

    await expect(shipmentContainer.locator('[data-testid="pickupAddress"]')).toContainText('123 Street');
    await expect(shipmentContainer.locator('[data-testid="pickupAddress"]')).toContainText('BEVERLY HILLS');
    await expect(shipmentContainer.locator('[data-testid="pickupAddress"]')).toContainText('CA');
    await expect(shipmentContainer.locator('[data-testid="pickupAddress"]')).toContainText('90210');

    await expect(shipmentContainer.locator('[data-testid="destinationAddress"]')).toContainText('123 Street');
    await expect(shipmentContainer.locator('[data-testid="destinationAddress"]')).toContainText('FORT WORTH');
    await expect(shipmentContainer.locator('[data-testid="destinationAddress"]')).toContainText('TX');
    await expect(shipmentContainer.locator('[data-testid="destinationAddress"]')).toContainText('76127');

    await expect(shipmentContainer.locator('[data-testid="sitPlanned"]')).toContainText('No');
    await expect(shipmentContainer.locator('[data-testid="estimatedWeight"]')).toContainText('4,000 lbs');
    await expect(shipmentContainer.locator('[data-testid="hasRequestedAdvance"]')).toContainText('Yes, $5,987');
    await expect(shipmentContainer.locator('[data-testid="counselorRemarks"]')).toContainText(
      'The requested advance amount has been added.',
    );

    // Verify fields in the expanded view
    await shipmentContainer.locator('[data-prefix="fas"][data-icon="chevron-down"]').click();
    await expect(shipmentContainer.locator('[data-testid="proGearWeight"]')).toContainText('Yes, 1,000 lbs');
    await expect(shipmentContainer.locator('[data-testid="spouseProGear"]')).toContainText('Yes, 500 lbs');
  });

  test('is able to add a new PPM shipment with gun safe', async ({ page, ppmPage }) => {
    test.skip(gunSafeEnabled === 'false', 'Skip if Gun Safe FF is disabled.');
    // Delete existing shipment
    await page.locator('[data-testid="ShipmentContainer"] .usa-button').click();
    await page.locator('[data-testid="grid"] button').getByText('Delete shipment').click();
    await expect(page.locator('[data-testid="modal"]')).toBeVisible();

    await page.locator('[data-testid="modal"] button').getByText('Delete shipment').click();

    await expect(page.locator('[data-testid="ShipmentContainer"] .usa-button')).not.toBeVisible();

    // Click add shipment button and select PPM
    await page.locator('[data-testid="dropdown"]').first().selectOption({ label: 'PPM' });

    // Fill out page one
    await ppmPage.fillOutOriginInfo();
    await ppmPage.fillOutDestinationInfo();
    await ppmPage.fillOutWeight({ hasProGear: true, hasGunSafe: true });
    await ppmPage.selectDutyLocation('JPPSO NORTHWEST', 'closeoutOffice');

    await page.locator('[data-testid="submitForm"]').click();
    await ppmPage.waitForLoading();

    // Fill out page two
    await expect(page.getByText('Incentive & advance')).toBeVisible();
    await ppmPage.fillOutIncentiveAndAdvance({ advance: '5987' });
    await page.locator('[data-testid="counselor-remarks"]').fill('The requested advance amount has been added.');
    await page.locator('[data-testid="counselor-remarks"]').blur();

    await page.locator('[data-testid="submitForm"]').click();
    await ppmPage.waitForLoading();

    // Confirm new shipment is visible
    await expect(page.locator('[data-testid="ShipmentContainer"]')).toBeVisible();
    const shipmentContainer = page.locator('[data-testid="ShipmentContainer"]');
    // Verify unexpanded view
    await expect(shipmentContainer.locator('[data-testid="expectedDepartureDate"]')).toContainText('09 Jun 2025');

    await expect(shipmentContainer.locator('[data-testid="pickupAddress"]')).toContainText('123 Street');
    await expect(shipmentContainer.locator('[data-testid="pickupAddress"]')).toContainText('BEVERLY HILLS');
    await expect(shipmentContainer.locator('[data-testid="pickupAddress"]')).toContainText('CA');
    await expect(shipmentContainer.locator('[data-testid="pickupAddress"]')).toContainText('90210');

    await expect(shipmentContainer.locator('[data-testid="destinationAddress"]')).toContainText('123 Street');
    await expect(shipmentContainer.locator('[data-testid="destinationAddress"]')).toContainText('FORT WORTH');
    await expect(shipmentContainer.locator('[data-testid="destinationAddress"]')).toContainText('TX');
    await expect(shipmentContainer.locator('[data-testid="destinationAddress"]')).toContainText('76127');

    await expect(shipmentContainer.locator('[data-testid="sitPlanned"]')).toContainText('No');
    await expect(shipmentContainer.locator('[data-testid="estimatedWeight"]')).toContainText('4,000 lbs');
    await expect(shipmentContainer.locator('[data-testid="hasRequestedAdvance"]')).toContainText('Yes, $5,987');
    await expect(shipmentContainer.locator('[data-testid="counselorRemarks"]')).toContainText(
      'The requested advance amount has been added.',
    );

    // Verify fields in the expanded view
    await shipmentContainer.locator('[data-prefix="fas"][data-icon="chevron-down"]').click();
    await expect(shipmentContainer.locator('[data-testid="proGearWeight"]')).toContainText('Yes, 1,000 lbs');
    await expect(shipmentContainer.locator('[data-testid="spouseProGear"]')).toContainText('Yes, 500 lbs');
    await expect(shipmentContainer.locator('[data-testid="gunSafeWeight"]')).toContainText('Yes, 400 lbs');
  });
});
