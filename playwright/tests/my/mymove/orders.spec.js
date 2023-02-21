// @ts-check
import { test, expect } from '../../utils/customerTest';

test('Users can upload orders', async ({ page, customerPage }) => {
  // Generate a new onboarded user and log in
  const user = await customerPage.testHarness.buildNeedsOrdersUser();
  const userId = user.id;
  await customerPage.signInAsExistingCustomer(userId);

  // Navigate to add orders
  await customerPage.waitForPage.home();
  await page.getByRole('button', { name: 'Add orders' }).click();
  await customerPage.waitForPage.ordersDetails();

  // Fill in orders details
  await page.getByTestId('dropdown').selectOption('PERMANENT_CHANGE_OF_STATION');
  await page.getByLabel('Orders date').fill('6/2/2018');
  await page.getByLabel('Orders date').blur();
  await page.getByLabel('Report by date').fill('8/9/2018');
  await page.getByLabel('Report by date').blur();

  // UGH
  // because of the styling of this input item, we cannot use a
  // css locator for the input item and then click it
  //
  // The styling is very similar to the issue described in
  //
  // https://github.com/microsoft/playwright/issues/3688
  //
  await page.locator('div:has(label:has-text("Are dependents")) >> div.usa-radio').getByText('No').click();

  // Verify that a user can't use the same duty location
  await customerPage.selectDutyLocation('Yuma AFB', 'new_duty_location');
  await expect(
    page.getByText('You entered the same duty location for your origin and destination. Please change one of them.'),
  ).toBeVisible();

  // Change to a different destination duty location, then proceed to next page
  await customerPage.selectDutyLocation('NAS Fort Worth JRB', 'new_duty_location');
  await customerPage.navigateForward();
  await customerPage.waitForPage.ordersUpload();

  // Upload an orders document, then submit
  // Annoyingly, there's no test IDs or labeling text for this control, so the only way to access it is .locator
  const filepondContainer = page.locator('.filepond--wrapper');
  await customerPage.uploadFileViaFilepond(filepondContainer, 'sampleWeightTicket.jpg');
  await customerPage.navigateForward();

  // Verify that we're on the home page and that orders have been uploaded
  await customerPage.waitForPage.home();
  await expect(page.getByText('Orders uploaded')).toBeVisible();
});
