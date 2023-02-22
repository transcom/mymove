// @ts-check
import { test, expect } from '../../utils/customerTest';

test('A customer can create, edit, and delete an HHG shipment', async ({ page, customerPage }) => {
  // Generate a new onboarded user with orders and log in
  const move = await customerPage.testHarness.buildMoveWithOrders();
  const userId = move.Orders.ServiceMember.user_id;
  await customerPage.signInAsExistingCustomer(userId);

  // Navigate to create a new shipment
  await customerPage.waitForPage.home();
  await page.getByTestId('shipment-selection-btn').click();
  await customerPage.waitForPage.aboutShipments();
  await customerPage.navigateForward();
  await customerPage.waitForPage.selectShipmentType();

  // Create an HHG shipment
  await page.getByText('Movers pack and ship it, paid by the government').click();
  await customerPage.navigateForward();

  // Fill in form to create HHG shipment
  await customerPage.waitForPage.hhgShipment();
  await page.getByLabel('Preferred pickup date').fill('25 Dec 2022');
  await page.getByLabel('Preferred pickup date').blur();
  await page.getByText('Use my current address').click();
  await page.getByLabel('Preferred delivery date').fill('25 Dec 2022');
  await page.getByLabel('Preferred delivery date').blur();
  await page.getByTestId('remarks').fill('Grandfather antique clock');
  await customerPage.navigateForward();

  // Verify that form submitted
  await customerPage.waitForPage.reviewShipments();
  await expect(page.getByText('Grandfather antique clock')).toBeVisible();
  await expect(page.getByTestId('ShipmentContainer').getByText('123 Any Street')).toBeVisible();

  // Navigate to edit shipment from the review page
  await page.getByTestId('edit-shipment-btn').click();
  await customerPage.waitForPage.hhgShipment();

  // Update form (adding pickup and delivery address)
  const pickupAddress = await page.getByRole('group', { name: 'Pickup location' });
  await pickupAddress.getByLabel('Address 1').fill('7 Q St');
  await pickupAddress.getByLabel('Address 2').clear();
  await pickupAddress.getByLabel('City').fill('Atco');
  await pickupAddress.getByLabel('State').selectOption({ label: 'NJ' });
  await pickupAddress.getByLabel('ZIP').fill('08004');

  const deliveryAddress = await page.getByRole('group', { name: 'Delivery location' });
  await deliveryAddress.getByText('Yes').click();
  await deliveryAddress.getByLabel('Address 1').fill('9 W 2nd Ave');
  await deliveryAddress.getByLabel('Address 2').fill('P.O. Box 456');
  await deliveryAddress.getByLabel('City').fill('Hollywood');
  await deliveryAddress.getByLabel('State').selectOption({ label: 'MD' });
  await deliveryAddress.getByLabel('ZIP').fill('20636');
  await customerPage.navigateForward();

  // Verify that shipment updated
  await customerPage.waitForPage.reviewShipments();
  await expect(page.getByTestId('ShipmentContainer').getByText('7 Q St')).toBeVisible();

  // Navigate to homepage and delete shipment
  await customerPage.navigateBack();
  await customerPage.waitForPage.home();
  await page.getByRole('button', { name: 'Delete' }).click();
  await page.getByTestId('modal').getByTestId('button').click();

  await expect(page.getByText('The shipment was deleted.')).toBeVisible();
  await expect(page.getByTestId('stepContainer3').getByText('Set up shipments')).toBeVisible();
});
