/**
 * Semi-automated converted from a cypress test, and thus may contain
 * non best-practices, in particular: heavy use of `page.locator`
 * instead of `page.getBy*`.
 */

// @ts-check
import { expect, test } from './customerPpmTestFixture';

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
    await page.locator('input[name="expectedDepartureDate"]').type('01 ZZZ 20222');
    await page.locator('input[name="expectedDepartureDate"]').blur();
    const errorMessage = page.locator('[class="usa-error-message"]');
    await expect(errorMessage).toContainText('Enter a complete date in DD MMM YYYY format (day, month, year).');
    await page.locator('input[name="expectedDepartureDate"]').clear();
    await page.locator('input[name="expectedDepartureDate"]').type('01 Feb 2022');
    await page.locator('input[name="expectedDepartureDate"]').blur();
    await expect(errorMessage).not.toBeVisible();

    // missing pickup postal code
    await page.locator('input[name="pickupAddress.address.streetAddress1"]').type('123 Street');
    await page.locator('input[name="pickupAddress.address.city"]').type('SomeCity - Secondary');
    await page.locator('select[name="pickupAddress.address.state"]').selectOption({ label: 'CA' });
    await page.locator('input[name="pickupAddress.address.postalCode"]').clear();
    await page.locator('input[name="pickupAddress.address.postalCode"]').blur();
    await expect(errorMessage).toContainText('Required');
    await page.locator('input[name="pickupAddress.address.postalCode"]').type('90210');
    await page.locator('input[name="pickupAddress.address.postalCode"]').blur();
    await expect(errorMessage).not.toBeVisible();

    // missing secondary pickup postal code
    await page.locator('label[for="yes-secondary-pickup-address"]').click();
    await page.locator('input[name="secondaryPickupAddress.address.streetAddress1"]').type('123 Street');
    await page.locator('input[name="secondaryPickupAddress.address.city"]').type('SomeCity - Secondary');
    await page.locator('select[name="secondaryPickupAddress.address.state"]').selectOption({ label: 'CA' });
    await page.locator('input[name="secondaryPickupAddress.address.postalCode"]').clear();
    await page.locator('input[name="secondaryPickupAddress.address.postalCode"]').blur();
    await expect(errorMessage).toContainText('Required');
    await page.locator('input[name="secondaryPickupAddress.address.postalCode"]').type('90210');
    await page.locator('input[name="secondaryPickupAddress.address.postalCode"]').blur();
    await expect(errorMessage).not.toBeVisible();

    // missing secondary destination postal code
    await page.locator('label[for="hasSecondaryDestinationAddressYes"]').click();
    await page.locator('input[name="secondaryDestinationAddress.address.streetAddress1"]').type('123 Street');
    await page.locator('input[name="secondaryDestinationAddress.address.city"]').type('SomeCity - Secondary');
    await page.locator('select[name="secondaryDestinationAddress.address.state"]').selectOption({ label: 'CA' });
    await page.locator('input[name="secondaryDestinationAddress.address.postalCode"]').clear();
    await page.locator('input[name="secondaryDestinationAddress.address.postalCode"]').blur();
    await expect(errorMessage).toContainText('Required');
    await page.locator('input[name="secondaryDestinationAddress.address.postalCode"]').type('90210');
    await page.locator('input[name="secondaryDestinationAddress.address.postalCode"]').blur();
    await expect(errorMessage).not.toBeVisible();
  });

  test('can continue to next page', async ({ customerPpmPage }) => {
    await customerPpmPage.submitsDateAndLocation();
  });
});
