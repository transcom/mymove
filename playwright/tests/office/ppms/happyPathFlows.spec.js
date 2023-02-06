/**
 * Semi-automated converted from a cypress test, and thus may contain
 * non best-practices, in particular: heavy use of `page.locator`
 * instead of `page.getBy*`.
 */

// @ts-check
const { test, expect, ScPpmPage } = require('./scPpmTestFixture');

test.describe('Services counselor user', () => {
  /** @type {ScPpmPage} */
  let scPpmPage;

  test.beforeEach(async ({ officePage }) => {
    const move = await officePage.testHarness.buildSubmittedMoveWithPPMShipmentForSC();
    await officePage.signInAsNewServicesCounselorUser();
    scPpmPage = new ScPpmPage(officePage, move);
    await scPpmPage.navigateToMove();
  });

  test('is able to edit a PPM shipment', async ({ page }) => {
    // View existing shipment
    await page.locator('[data-testid="ShipmentContainer"] .usa-button').click();

    await scPpmPage.fillOutSitExpected();
    await scPpmPage.selectDutyLocation('JPPSO NORTHWEST', 'closeoutOffice');

    // Submit page 1 of form
    await page.locator('[data-testid="submitForm"]').click();
    await scPpmPage.waitForLoading();

    // Verify SIT info
    await expect(page.getByText('Government constructed cost: $326')).toBeVisible();
    await expect(page.getByText('1,000 lbs of destination SIT at 30813 for 31 days.')).toBeVisible();
    // Verify estimated incentive
    await expect(page.getByRole('heading', { name: 'Estimated incentive: $10,000' })).toBeVisible();

    // Update page 2
    scPpmPage.fillOutIncentiveAndAdvance();
    await expect(page.locator('[data-testid="errorMessage"]')).toContainText('Required');
    await page.locator('[data-testid="counselor-remarks"]').type('Increased incentive to max');
    await page.locator('[data-testid="counselor-remarks"]').blur();

    // Submit page 2 of form
    await page.locator('[data-testid="submitForm"]').click();
    await scPpmPage.waitForLoading();

    // Expand details and verify information
    await expect(page.getByText('Your changes were saved.')).toBeVisible();
    expect(page.locator('[data-testid="ShipmentContainer"]')).toBeVisible();
    let shipmentContainer = page.locator('[data-testid="ShipmentContainer"]');
    await shipmentContainer.locator('[data-prefix="fas"][data-icon="chevron-down"]').click();
    await expect(shipmentContainer.locator('[data-testid="expectedDepartureDate"]')).toContainText('15 Mar 2020');
    await expect(shipmentContainer.locator('[data-testid="originZIP"]')).toContainText('90210');
    await expect(shipmentContainer.locator('[data-testid="secondOriginZIP"]')).toContainText('90211');
    await expect(shipmentContainer.locator('[data-testid="destinationZIP"]')).toContainText('30813');
    await expect(shipmentContainer.locator('[data-testid="secondDestinationZIP"]')).toContainText('30814');
    await expect(shipmentContainer.locator('[data-testid="sitPlanned"]')).toContainText('Yes');
    await expect(shipmentContainer.locator('[data-testid="estimatedWeight"]')).toContainText('4,000 lbs');
    await expect(shipmentContainer.locator('[data-testid="proGearWeight"]')).toContainText('Yes, 1,987 lbs');
    await expect(shipmentContainer.locator('[data-testid="spouseProGear"]')).toContainText('Yes, 498 lbs');
    await expect(shipmentContainer.locator('[data-testid="estimatedIncentive"]')).toContainText('$10,000');
    await expect(shipmentContainer.locator('[data-testid="hasRequestedAdvance"]')).toContainText('Yes, $6,000');
    await expect(shipmentContainer.locator('[data-testid="counselorRemarks"]')).toContainText(
      'Increased incentive to max',
    );

    // combined with test above
    // test 'is able to add a second PPM shipment'
    await page.locator('[data-testid="dropdown"]').selectOption({ label: 'PPM' });

    // Fill out page one
    await scPpmPage.fillOutOriginInfo();
    await scPpmPage.fillOutDestinationInfo();
    await scPpmPage.fillOutSitExpected();
    await scPpmPage.fillOutWeight({ hasProGear: true });

    await page.locator('[data-testid="submitForm"]').click();
    await scPpmPage.waitForLoading();

    // Verify SIT info
    await expect(page.getByText('Government constructed cost: $379')).toBeVisible();
    await expect(page.getByText('1,000 lbs of destination SIT at 76127 for 31 days.')).toBeVisible();
    // Verify estimated incentive
    await expect(page.getByRole('heading', { name: 'Estimated incentive: $67,692' })).toBeVisible();

    // Fill out page two
    await scPpmPage.fillOutIncentiveAndAdvance({ advance: '10000' });
    await expect(page.locator('[data-testid="errorMessage"]')).toContainText('Required');
    await page.locator('[data-testid="counselor-remarks"]').type('Added correct incentive');
    await page.locator('[data-testid="counselor-remarks"]').blur();

    // Submit page two
    await page.locator('[data-testid="submitForm"]').click();
    await scPpmPage.waitForLoading();

    // Expand details and verify information
    await expect(page.getByText('Your changes were saved.')).toBeVisible();
    expect(page.locator('[data-testid="ShipmentContainer"]')).toHaveCount(2);
    shipmentContainer = page.locator('[data-testid="ShipmentContainer"]').last();
    await shipmentContainer.locator('[data-prefix="fas"][data-icon="chevron-down"]').click();
    await expect(shipmentContainer.locator('[data-testid="expectedDepartureDate"]')).toContainText('09 Jun 2022');
    await expect(shipmentContainer.locator('[data-testid="originZIP"]')).toContainText('90210');
    await expect(shipmentContainer.locator('[data-testid="secondOriginZIP"]')).toContainText('07003');
    await expect(shipmentContainer.locator('[data-testid="destinationZIP"]')).toContainText('76127');
    await expect(shipmentContainer.locator('[data-testid="secondDestinationZIP"]')).toContainText('08540');
    await expect(shipmentContainer.locator('[data-testid="sitPlanned"]')).toContainText('Yes');
    await expect(shipmentContainer.locator('[data-testid="estimatedWeight"]')).toContainText('4,000 lbs');
    await expect(shipmentContainer.locator('[data-testid="proGearWeight"]')).toContainText('Yes, 1,000 lbs');
    await expect(shipmentContainer.locator('[data-testid="spouseProGear"]')).toContainText('Yes, 500 lbs');
    await expect(shipmentContainer.locator('[data-testid="estimatedIncentive"]')).toContainText('$67,692');
    await expect(shipmentContainer.locator('[data-testid="hasRequestedAdvance"]')).toContainText('Yes, $10,000');
    await expect(shipmentContainer.locator('[data-testid="counselorRemarks"]')).toContainText(
      'Added correct incentive',
    );
  });
});
