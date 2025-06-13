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

test.describe('Services counselor user', () => {
  test.beforeEach(async ({ ppmPage }) => {
    const move = await ppmPage.testHarness.buildSubmittedMoveWithPPMShipmentForSC();
    await ppmPage.navigateToMove(move.locator);
  });

  test('is able to edit a PPM shipment', async ({ page, ppmPage }) => {
    // View existing shipment
    await page.locator('[data-testid="ShipmentContainer"] .usa-button').click();

    await ppmPage.fillOutSitExpected();
    await ppmPage.selectDutyLocation('JPPSO NORTHWEST', 'closeoutOffice');

    // Submit page 1 of form
    await page.locator('[data-testid="submitForm"]').click();
    await ppmPage.waitForLoading();

    // Verify SIT info
    await expect(page.getByText(/Government constructed cost:/)).toBeVisible();
    await expect(page.getByText('1,000 lbs of destination SIT at 30813 for 31 days.')).toBeVisible();
    // Verify estimated incentive
    await expect(page.getByRole('heading', { name: 'Estimated incentive:' })).toBeVisible();

    // Update page 2
    await ppmPage.fillOutIncentiveAndAdvance();
    await page.locator('[data-testid="counselor-remarks"]').fill('Increased incentive to max');
    await page.locator('[data-testid="counselor-remarks"]').blur();

    // Submit page 2 of form
    await page.locator('[data-testid="submitForm"]').click();
    await ppmPage.waitForLoading();

    // Expand details and verify information
    await expect(page.getByText('Your changes were saved.')).toBeVisible();
    await expect(page.locator('[data-testid="ShipmentContainer"]')).toBeVisible();
    let shipmentContainer = page.locator('[data-testid="ShipmentContainer"]');
    await shipmentContainer.locator('[data-prefix="fas"][data-icon="chevron-down"]').click();
    await expect(shipmentContainer.locator('[data-testid="expectedDepartureDate"]')).toContainText('15 Mar 2025');

    await expect(shipmentContainer.locator('[data-testid="pickupAddress"]')).toContainText('987 New Street');
    await expect(shipmentContainer.locator('[data-testid="pickupAddress"]')).toContainText('P.O. Box 12345');
    await expect(shipmentContainer.locator('[data-testid="pickupAddress"]')).toContainText('Des Moines');
    await expect(shipmentContainer.locator('[data-testid="pickupAddress"]')).toContainText('IA');
    await expect(shipmentContainer.locator('[data-testid="pickupAddress"]')).toContainText('50309');

    await expect(shipmentContainer.locator('[data-testid="destinationAddress"]')).toContainText('123 New Street');
    await expect(shipmentContainer.locator('[data-testid="pickupAddress"]')).toContainText('P.O. Box 12345');
    await expect(shipmentContainer.locator('[data-testid="destinationAddress"]')).toContainText('Fort Eisenhower');
    await expect(shipmentContainer.locator('[data-testid="destinationAddress"]')).toContainText('GA');
    await expect(shipmentContainer.locator('[data-testid="destinationAddress"]')).toContainText('30813');

    await expect(shipmentContainer.locator('[data-testid="sitPlanned"]')).toContainText('Yes');
    await expect(shipmentContainer.locator('[data-testid="estimatedWeight"]')).toContainText('4,000 lbs');
    await expect(shipmentContainer.locator('[data-testid="proGearWeight"]')).toContainText('Yes, 1,987 lbs');
    await expect(shipmentContainer.locator('[data-testid="spouseProGear"]')).toContainText('Yes, 498 lbs');
    await expect(shipmentContainer.locator('[data-testid="hasRequestedAdvance"]')).toContainText('Yes, $6,000');
    await expect(shipmentContainer.locator('[data-testid="counselorRemarks"]')).toContainText(
      'Increased incentive to max',
    );

    // combined with test above
    // test 'is able to add a second PPM shipment'
    await page.locator('[data-testid="dropdown"]').selectOption({ label: 'PPM' });

    // Fill out page one
    await ppmPage.fillOutOriginInfo();
    await ppmPage.fillOutDestinationInfo();
    await ppmPage.fillOutSitExpected();
    await ppmPage.fillOutWeight({ hasProGear: true });

    await page.locator('[data-testid="submitForm"]').click();
    await ppmPage.waitForLoading();

    // Verify SIT info
    await expect(page.getByText(/Government constructed cost:/)).toBeVisible();
    await expect(page.getByText('1,000 lbs of destination SIT at 76127 for 31 days.')).toBeVisible();
    // Verify estimated incentive
    await expect(page.getByRole('heading', { name: 'Estimated incentive:' })).toBeVisible();

    // Fill out page two
    await ppmPage.fillOutIncentiveAndAdvance({ advance: '10000' });
    await page.locator('[data-testid="counselor-remarks"]').fill('Added correct incentive');
    await page.locator('[data-testid="counselor-remarks"]').blur();

    // Submit page two
    await page.locator('[data-testid="submitForm"]').click();
    await ppmPage.waitForLoading();

    // Expand details and verify information
    await expect(page.getByText('Your changes were saved.')).toBeVisible();
    shipmentContainer = page.locator('[data-testid="ShipmentContainer"]').last();
    await shipmentContainer.locator('[data-prefix="fas"][data-icon="chevron-down"]').click();
    const expectedDeparture = new Date(Date.now() + 24 * 60 * 60 * 1000).toLocaleDateString('en-US');
    const formattedDate = formatDate(expectedDeparture);
    await expect(shipmentContainer.locator('[data-testid="expectedDepartureDate"]')).toContainText(formattedDate);

    await expect(shipmentContainer.locator('[data-testid="pickupAddress"]')).toContainText('123 Street');
    await expect(shipmentContainer.locator('[data-testid="pickupAddress"]')).toContainText('BEVERLY HILLS');
    await expect(shipmentContainer.locator('[data-testid="pickupAddress"]')).toContainText('CA');
    await expect(shipmentContainer.locator('[data-testid="pickupAddress"]')).toContainText('90210');

    await expect(shipmentContainer.locator('[data-testid="destinationAddress"]')).toContainText('123 Street');
    await expect(shipmentContainer.locator('[data-testid="destinationAddress"]')).toContainText('FORT WORTH,');
    await expect(shipmentContainer.locator('[data-testid="destinationAddress"]')).toContainText('TX');
    await expect(shipmentContainer.locator('[data-testid="destinationAddress"]')).toContainText('76127');

    await expect(shipmentContainer.locator('[data-testid="sitPlanned"]')).toContainText('Yes');
    await expect(shipmentContainer.locator('[data-testid="estimatedWeight"]')).toContainText('4,000 lbs');
    await expect(shipmentContainer.locator('[data-testid="proGearWeight"]')).toContainText('Yes, 1,000 lbs');
    await expect(shipmentContainer.locator('[data-testid="spouseProGear"]')).toContainText('Yes, 500 lbs');
    const text = await shipmentContainer.locator('[data-testid="estimatedIncentive"]').textContent();
    expect(text).toMatch(/^\$\d/); // Check that it starts with a dollar sign and digit - this is to address flaky tests
    expect(text).not.toContain('undefined');
    await expect(shipmentContainer.locator('[data-testid="hasRequestedAdvance"]')).toContainText('Yes, $10,000');
    await expect(shipmentContainer.locator('[data-testid="counselorRemarks"]')).toContainText(
      'Added correct incentive',
    );
  });
});

test.describe('Services counselor user approval', () => {
  test.beforeEach(async ({ ppmPage }) => {
    const move = await ppmPage.testHarness.buildApprovedMoveWithSubmittedPPMShipmentForSC();
    await ppmPage.navigateToMoveUsingMoveSearch(move.locator);
  });

  test('is able to approve a submitted PPM shipment', async ({ page }) => {
    // View existing shipment
    await page.locator('[data-testid="sendPpmToCustomerButton"]').click();
    await page.locator('button:text("Yes, submit")').click();

    await expect(page.locator('[data-testid="ppmStatusTag"]')).toContainText('Waiting on customer');
  });
});
