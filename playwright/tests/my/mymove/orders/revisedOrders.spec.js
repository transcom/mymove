// @ts-check
const { test, expect } = require('../customerPage');

test('A customer can upload their orders', async ({ page, customerPage }) => {
  const user = await customerPage.testHarness.buildNeedsOrdersUser();
  const userId = user.id;
  await customerPage.signIn.customer.existingCustomer(userId);

  await page.getByRole('button', { name: 'Add orders' }).click();
  await customerPage.waitForPage.orders();

  // Fill out orders detail form
  await page.getByTestId('dropdown').selectOption('PERMANENT_CHANGE_OF_STATION');
  await page.getByLabel('Orders date').fill('25 Dec 2022');
  await page.getByLabel('Orders date').blur();
  await page.getByLabel('Report by date').fill('29 Dec 2022');
  await page.getByLabel('Report by date').blur();

  // Because of the styling for this item, we can't use a standard locator for it and click it
  // Similar to this existing issue: https://github.com/microsoft/playwright/issues/3688
  await page.locator('div:has(label:has-text("Are dependents")) >> div.usa-radio').getByText('No').click();

  // Select the user's current duty location as their new duty location, and verify that an error displays
  await page.getByLabel('New duty location').fill('Yuma AFB');
  // 'mark' is not yet supported by react testing library
  // https://github.com/testing-library/dom-testing-library/issues/1150
  // @ts-ignore:next-line
  await page.getByRole('mark').click();
  await expect(page.getByText('You entered the same duty location for your origin and destination.')).toBeVisible();

  // Enter another duty location and proceed to next page
  await page.getByLabel('New duty location').fill('Scott AFB');
  // 'mark' is not yet supported by react testing library
  // https://github.com/testing-library/dom-testing-library/issues/1150
  // @ts-ignore:next-line
  await page.getByRole('mark').click();
  await customerPage.navigateForward();

  // Upload orders document
  await customerPage.waitForPage.ordersUpload();
  await page.getByLabel('Drag & drop or click to upload orders').setInputFiles('../fixtures/sample-orders.png');
  await customerPage.navigateForward();
});
