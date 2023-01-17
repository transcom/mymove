/**
 * Semi-automated converted from a cypress test, and thus may contain
 * non best-practices, in particular: heavy use of `page.locator`
 * instead of `page.getBy*`.
 */

// @ts-check
const { test, expect } = require('../customerPage');

/**
 * @param {import('@playwright/test').Page} page
 * @param {string} inputData
 * @param {string} fieldName
 * @param {string} classSelector
 */
async function genericSelect(page, inputData, fieldName, classSelector) {
  // fieldName is passed as a classname to the react-select component,
  // so select for it if provided
  const actualClassSelector = fieldName ? `${classSelector}.${fieldName}` : classSelector;
  await page.locator(`${actualClassSelector} input[type="text"]`).type(inputData);

  // Click on the first presented option
  await page.locator(classSelector).locator('div[class*="option"]').first().click();
}

/**
 * @param {import('@playwright/test').Page} page
 * @param {string} dutyLocationName
 * @param {string} fieldName
 */
async function selectDutyLocation(page, dutyLocationName, fieldName) {
  return genericSelect(page, dutyLocationName, fieldName, '.duty-input-box');
}

/**
 * @param {import('@playwright/test').Page} page
 * @param {string} flagVal
 */
async function setFeatureFlag(page, flagVal, url = '/queues/new') {
  await page.goto(`${url}?flag:${flagVal}`);
}

test('orders entry will accept orders information', async ({ page, customerPage }) => {
  const user = await customerPage.testHarness.buildNeedsOrdersUser();
  const userId = user.id;
  await customerPage.signIn.customer.existingCustomer(userId);

  await expect(page.getByText('Next step: Add your orders')).toBeVisible();
  await expect(page.getByText('Profile complete')).toBeVisible();
  await expect(page.getByText('Upload orders')).toBeVisible();

  await page.getByRole('button', { name: 'Add orders' }).click();
  expect(page.url()).toContain('/orders/info');

  await page.locator('select[name="orders_type"]').selectOption('SEPARATION');
  await page.locator('select[name="orders_type"]').selectOption('RETIREMENT');
  await page.locator('select[name="orders_type"]').selectOption('PERMANENT_CHANGE_OF_STATION');

  await page.locator('input[name="issue_date"]').click();
  await page.locator('input[name="issue_date"]').type('6/2/2018');
  await page.locator('input[name="issue_date"]').blur();

  await page.locator('input[name="report_by_date"]').type('8/9/2018');
  await page.locator('input[name="report_by_date"]').blur();

  // UGH
  // because of the styling of this input item, we cannot use a
  // css locator for the input item and then click it
  //
  // The styling is very similar to the issue described in
  //
  // https://github.com/microsoft/playwright/issues/3688
  //
  await page.locator('div:has(label:has-text("Are dependents")) >> div.usa-radio').getByText('No').click();

  // Choosing same current and destination duty location should block you from progressing and give an error
  await selectDutyLocation(page, 'Yuma AFB', 'new_duty_location');
  await expect(page.locator('.usa-error-message')).toContainText(
    'You entered the same duty location for your origin and destination. Please change one of them.',
  );
  await expect(page.locator('button[data-testid="wizardNextButton"]')).toBeDisabled();

  await selectDutyLocation(page, 'NAS Fort Worth JRB', 'new_duty_location');

  await page.getByRole('button', { name: 'Next' }).click();
  await expect(page.getByText('Upload your orders')).toBeVisible();

  expect(page.url()).toContain('/orders/upload');

  setFeatureFlag(page, 'ppmPaymentRequest=false', '/ppm');
  await expect(page.getByText('NAS Fort Worth (from Yuma AFB)')).toBeVisible();
  await expect(page.locator('[data-testid="move-header-weight-estimate"]')).toHaveText('5,000 lbs');

  await page.getByText('Continue Move Setup').click();
  expect(page.url()).toContain('/orders/upload');
});
