/**
 * Semi-automated converted from a cypress test, and thus may contain
 * non best-practices, in particular: heavy use of `page.locator`
 * instead of `page.getBy*`.
 */

// @ts-check
import { expect, test } from './customerPpmTestFixture';

const multiMoveEnabled = process.env.FEATURE_FLAG_MULTI_MOVE;

test.describe('PPM Onboarding - Add dates and location flow', () => {
  test.beforeEach(async ({ customerPpmPage }) => {
    const move = await customerPpmPage.testHarness.buildSpouseProGearMove();
    await customerPpmPage.signInForPPMWithMove(move);
    await customerPpmPage.customerStartsAddingAPPMShipment();
  });

  test('doesn’t allow SM to progress if form is in an invalid state', async ({ page }) => {
    await expect(page.getByText('PPM date & location')).toBeVisible();
    expect(page.url()).toContain('/new-shipment');

    // invalid date
    await page.locator('input[name="expectedDepartureDate"]').fill('01 ZZZ 20222');
    await page.locator('input[name="expectedDepartureDate"]').blur();
    const errorMessage = page.locator('[class="usa-error-message"]');
    await expect(errorMessage).toContainText('Enter a complete date in DD MMM YYYY format (day, month, year).');
    await page.locator('input[name="expectedDepartureDate"]').clear();
    const expectedDeparture = new Date(Date.now() + 24 * 60 * 60 * 1000).toLocaleDateString('en-US');
    await page.locator('input[name="expectedDepartureDate"]').fill(expectedDeparture);
    await page.locator('input[name="expectedDepartureDate"]').blur();
    await expect(errorMessage).not.toBeVisible();

    const pickupLocation = 'BEVERLY HILLS, CA 90210 (LOS ANGELES)';
    const secondaryPickupLocation = 'YUMA, AZ 85364 (YUMA)';

    await page.locator('input[name="pickupAddress.address.streetAddress1"]').fill('123 Street');
    await page.locator('input[name="pickupAddress.address.streetAddress1"]').clear();
    await page.locator('input[name="pickupAddress.address.streetAddress1"]').blur();
    await expect(errorMessage).toContainText('Required');
    await page.locator('input[name="pickupAddress.address.streetAddress1"]').fill('123 Street');
    await page.locator('input[id="pickupAddress.address-input"]').fill('90210');
    await expect(page.getByText(pickupLocation, { exact: true })).toBeVisible();
    await page.keyboard.press('Enter');
    await expect(errorMessage).not.toBeVisible();

    // missing secondary pickup street address
    await page.locator('label[for="yes-secondary-pickup-address"]').click();
    await page.locator('input[name="secondaryPickupAddress.address.streetAddress1"]').fill('123 Street');
    await page.locator('input[name="secondaryPickupAddress.address.streetAddress1"]').clear();
    await page.locator('input[name="secondaryPickupAddress.address.streetAddress1"]').blur();
    await page.locator('input[id="secondaryPickupAddress.address-input"]').fill('85364');
    await expect(page.getByText(secondaryPickupLocation, { exact: true })).toBeVisible();
    await page.keyboard.press('Enter');

    await expect(page.getByText('Save & Continue')).toBeDisabled();
  });

  test('can continue to next page', async ({ customerPpmPage }) => {
    await customerPpmPage.submitsDateAndLocation();
  });
});

test.describe('(MultiMove) PPM Onboarding - Add dates and location flow', () => {
  test.skip(multiMoveEnabled === 'false', 'Skip if MultiMove workflow is not enabled.');
  test.fail(multiMoveEnabled === 'true', 'Need to fix zipcode validation error messages.');
  test.beforeEach(async ({ customerPpmPage }) => {
    const move = await customerPpmPage.testHarness.buildSpouseProGearMove();
    await customerPpmPage.signInForPPMWithMove(move);
    await customerPpmPage.clickOnGoToMoveButton();
    await customerPpmPage.customerStartsAddingAPPMShipment();
  });

  test.skip('doesn’t allow SM to progress if form is in an invalid state', async ({ page }) => {
    await expect(page.getByText('PPM date & location')).toBeVisible();
    expect(page.url()).toContain('/new-shipment');

    // invalid date
    await page.locator('input[name="expectedDepartureDate"]').fill('01 ZZZ 20222');
    await page.locator('input[name="expectedDepartureDate"]').blur();
    const errorMessage = page.locator('[class="usa-error-message"]');
    await expect(errorMessage).toContainText('Enter a complete date in DD MMM YYYY format (day, month, year).');
    await page.locator('input[name="expectedDepartureDate"]').clear();
    const expectedDeparture = new Date(Date.now() + 24 * 60 * 60 * 1000).toLocaleDateString('en-US');
    await page.locator('input[name="expectedDepartureDate"]').fill(expectedDeparture);
    await page.locator('input[name="expectedDepartureDate"]').blur();
    await expect(errorMessage).not.toBeVisible();

    // invalid postal codes
    await page.locator('input[name="pickupAddress.address.postalCode"]').fill('00000');
    await page.locator('input[name="pickupAddress.address.postalCode"]').blur();
    await expect(errorMessage).toContainText(
      "We don't have rates for this ZIP code. Please verify that you have entered the correct one. Contact support if this problem persists.",
    );
    await page.locator('input[name="pickupAddress.address.postalCode"]').clear();
    await page.locator('input[name="pickupAddress.address.postalCode"]').fill('90210');
    await page.locator('input[name="pickupAddress.address.postalCode"]').blur();
    await expect(errorMessage).not.toBeVisible();

    await page.locator('input[name="pickupAddress.address.postalCode"]').clear();
    await page.locator('input[name="pickupAddress.address.postalCode"]').blur();
    await expect(errorMessage).toContainText('Required');
    await expect(page.locator('[class="usa-error-message"] + input')).toHaveAttribute(
      'id',
      'pickupAddress.address.postalCode',
    );
    await page.locator('input[name="pickupAddress.address.postalCode"]').fill('90210');
    await page.locator('input[name="pickupAddress.address.postalCode"]').blur();
    await expect(errorMessage).not.toBeVisible();

    // missing secondary pickup postal code
    await page.locator('label[for="yes-secondary-pickup-postal-code"]').click();
    await page.locator('input[name="secondaryPickupAddress.address.postalCode"]').clear();
    await page.locator('input[name="secondaryPickupAddress.address.postalCode"]').blur();
    await expect(errorMessage).toContainText('Required');
    await expect(page.locator('[class="usa-error-message"] + input')).toHaveAttribute(
      'id',
      'secondaryPickupAddress.address.postalCode',
    );
    await page.locator('input[name="secondaryPickupAddress.address.postalCode"]').fill('90210');
    await page.locator('input[name="secondaryPickupAddress.address.postalCode"]').blur();
    await expect(errorMessage).not.toBeVisible();

    // missing secondary destination postal code
    await page.locator('label[for="hasSecondaryDestinationAddressYes"]').click();
    await page.locator('input[name="secondaryDestinationAddress.address.postalCode"]').clear();
    await page.locator('input[name="secondaryDestinationAddress.address.postalCode"]').blur();
    await expect(errorMessage).toContainText('Required');
    await expect(page.locator('[class="usa-error-message"] + input')).toHaveAttribute(
      'id',
      'secondaryDestinationAddress.address.postalCode',
    );
    await page.locator('input[name="secondaryDestinationAddress.address.postalCode"]').fill('90210');
    await page.locator('input[name="secondaryDestinationAddress.address.postalCode"]').blur();
    await expect(errorMessage).not.toBeVisible();
  });

  test.skip('can continue to next page', async ({ customerPpmPage }) => {
    await customerPpmPage.submitsDateAndLocation();
  });
});
