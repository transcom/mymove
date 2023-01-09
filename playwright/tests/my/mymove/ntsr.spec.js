// @ts-check
const { test, expect } = require('../../utils/customerTest');

test('A customer can create, edit, and delete an NTS-release shipment', async ({ page, customerPage }) => {
  const move = await customerPage.testHarness.buildMoveWithOrders();
  const userId = move.Orders.ServiceMember.user_id;
  await customerPage.signIn.customer.existingCustomer(userId);

  // Navigate to create a new shipment
  await page.getByTestId('shipment-selection-btn').click();
  await customerPage.navigateForward();
  await page.getByText('It was stored during a previous move (NTS-release)Movers pick up things you put ').click();
  await customerPage.navigateForward();

  // Fill in form to create NTS shipment
  await page.getByLabel('Preferred delivery date').fill('25 Dec 2022');
  await page.getByLabel('Preferred delivery date').blur();
  await page.getByLabel('Address 1').fill('7 Q St');
  await page.getByLabel('City').fill('Atco');
  await page.getByLabel('State').selectOption('NJ');
  await page.getByLabel('ZIP').fill('08004');

  await page.getByTestId('remarks').fill('Grandfather antique clock');
  await customerPage.navigateForward();

  // Verify that form submitted
  await customerPage.waitForPage.reviewShipments();
  await expect(page.getByText('Grandfather antique clock')).toBeVisible();

  // Navigate to edit shipment from the review page
  await page.getByTestId('edit-ntsr-shipment-btn').click();
  await customerPage.waitForPage.ntsReleaseShipment();

  // Update form (adding releasing agent)
  await page.getByLabel('First name').fill('Grace');
  await page.getByLabel('Last name').fill('Griffin');
  await page.getByLabel('Phone').fill('2025551234');
  await page.getByLabel('Email').fill('grace.griffin@example.com');
  await customerPage.navigateForward();

  // Verify that form submitted
  await customerPage.waitForPage.reviewShipments();
  await expect(page.getByText('Grace Griffin')).toBeVisible();

  // Navigate to homepage and delete shipment
  await customerPage.navigateBack();
  await customerPage.waitForPage.home();
  await page.getByRole('button', { name: 'Delete' }).click();
  await page.getByTestId('modal').getByTestId('button').click();

  await expect(page.getByText('The shipment was deleted.')).toBeVisible();
  await expect(page.getByTestId('stepContainer3').getByText('Set up shipments')).toBeVisible();
});
