// @ts-check
const { test, expect } = require('../customerPage');

test('A customer can create, edit, and delete an HHG shipment', async ({ page, customerPage }) => {
  const move = await customerPage.testHarness.buildMoveWithOrders();
  const userId = move.Orders.ServiceMember.user_id;
  await customerPage.signIn.customer.existingCustomer(userId);

  // Navigate to create a new shipment
  await page.getByTestId('shipment-selection-btn').click();
  await customerPage.navigateForward();
  await page.getByText('Movers pack and ship it, paid by the government (HHG)').click();
  await customerPage.navigateForward();

  // Fill in form to create HHG shipment
  await customerPage.waitForPage.hhgShipment();
  await page.getByLabel('Preferred pickup date').fill('25 Dec 2022');
  await page.getByLabel('Preferred pickup date').blur();
  await page.getByText('Use my current address').click();
  await page.getByLabel('Preferred delivery date').fill('29 Dec 2022');
  await page.getByTestId('remarks').fill('Grandfather antique clock');
  await customerPage.navigateForward();

  // Verify that form submitted
  await customerPage.waitForPage.reviewShipments();
  await expect(page.getByText('Grandfather antique clock')).toBeVisible();

  // Navigate to edit shipment from the review page
  await page.getByTestId('edit-shipment-btn').click();
  await customerPage.waitForPage.hhgShipment();

  // Update form (adding releasing agent)
  await page.getByRole('group', { name: 'Releasing agent Optional' }).getByLabel('First name').fill('Grace');
  await page.getByRole('group', { name: 'Releasing agent Optional' }).getByLabel('Last name').fill('Griffin');
  await page.getByRole('group', { name: 'Releasing agent Optional' }).getByLabel('Phone').fill('2025551234');
  await page
    .getByRole('group', { name: 'Releasing agent Optional' })
    .getByLabel('Email')
    .fill('grace.griffin@example.com');
  await page.getByTestId('wizardNextButton').click();

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
