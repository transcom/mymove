/**
 * Semi-automated converted from a cypress test, and thus may contain
 * non best-practices, in particular: heavy use of `page.locator`
 * instead of `page.getBy*`.
 */

// @ts-check
const { expect, test, CustomerPpmPage } = require('./customerPpmTestFixture');

test.describe('PPM Onboarding - Add dates and location flow', () => {
  /** @type {CustomerPpmPage} */
  let customerPpmPage;

  test.beforeEach(async ({ customerPage }) => {
    const move = await customerPage.testHarness.buildSpouseProGearMove();
    customerPpmPage = new CustomerPpmPage(customerPage, move);
    await customerPpmPage.signInForPPM();
    await customerPpmPage.customerStartsAddingAPPMShipment();
  });

  test('doesn’t allow SM to progress if form is in an invalid state', async ({ page }) => {
    await expect(page.getByText('PPM date & location')).toBeVisible();
    expect(page.url()).toContain('/new-shipment');

    // invalid date
    await page.locator('input[name="expectedDepartureDate"]').type('01 ZZZ 20222');
    await page.locator('input[name="expectedDepartureDate"]').blur();
    const errorMessage = page.locator('[class="usa-error-message"]');
    await expect(errorMessage).toContainText('Enter a complete date in DD MMM YYYY format (day, month, year).');
    await page.locator('input[name="expectedDepartureDate"]').clear();
    await page.locator('input[name="expectedDepartureDate"]').type('01 Feb 2022');
    await page.locator('input[name="expectedDepartureDate"]').blur();
    await expect(errorMessage).not.toBeVisible();

    // invalid postal codes
    await page.locator('input[name="pickupPostalCode"]').type('00000');
    await page.locator('input[name="pickupPostalCode"]').blur();
    await expect(errorMessage).toContainText(
      "We don't have rates for this ZIP code. Please verify that you have entered the correct one. Contact support if this problem persists.",
    );
    await page.locator('input[name="pickupPostalCode"]').clear();
    await page.locator('input[name="pickupPostalCode"]').type('90210');
    await page.locator('input[name="pickupPostalCode"]').blur();
    await expect(errorMessage).not.toBeVisible();

    await page.locator('input[name="pickupPostalCode"]').clear();
    await page.locator('input[name="pickupPostalCode"]').blur();
    await expect(errorMessage).toContainText('Required');
    await expect(page.locator('[class="usa-error-message"] + input')).toHaveAttribute('id', 'pickupPostalCode');
    await page.locator('input[name="pickupPostalCode"]').type('90210');
    await page.locator('input[name="pickupPostalCode"]').blur();
    await expect(errorMessage).not.toBeVisible();

    // missing secondary pickup postal code
    await page.locator('label[for="yes-secondary-pickup-postal-code"]').click();
    await page.locator('input[name="secondaryPickupPostalCode"]').clear();
    await page.locator('input[name="secondaryPickupPostalCode"]').blur();
    await expect(errorMessage).toContainText('Required');
    await expect(page.locator('[class="usa-error-message"] + input')).toHaveAttribute(
      'id',
      'secondaryPickupPostalCode',
    );
    await page.locator('input[name="secondaryPickupPostalCode"]').type('90210');
    await page.locator('input[name="secondaryPickupPostalCode"]').blur();
    await expect(errorMessage).not.toBeVisible();

    // missing secondary destination postal code
    await page.locator('label[for="hasSecondaryDestinationPostalCodeYes"]').click();
    await page.locator('input[name="secondaryDestinationPostalCode"]').clear();
    await page.locator('input[name="secondaryDestinationPostalCode"]').blur();
    await expect(errorMessage).toContainText('Required');
    await expect(page.locator('[class="usa-error-message"] + input')).toHaveAttribute(
      'id',
      'secondaryDestinationPostalCode',
    );
    await page.locator('input[name="secondaryDestinationPostalCode"]').type('90210');
    await page.locator('input[name="secondaryDestinationPostalCode"]').blur();
    await expect(errorMessage).not.toBeVisible();
  });

  test('can continue to next page', async () => {
    await customerPpmPage.submitsDateAndLocation();
  });
});
