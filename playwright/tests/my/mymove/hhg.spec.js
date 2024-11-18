// @ts-check
import { test, expect } from '../../utils/my/customerTest';

const multiMoveEnabled = process.env.FEATURE_FLAG_MULTI_MOVE;

test.describe('HHG', () => {
  test.skip(multiMoveEnabled === 'true', 'Skip if MultiMove workflow is enabled.');
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
    const pickupLocation = 'BEVERLY HILLS, CA 90212 (LOS ANGELES)';
    const pickupAddress = page.getByRole('group', { name: 'Pickup location' });
    await pickupAddress.getByLabel('Address 1').fill('7 Q St');
    await pickupAddress.getByLabel('Address 2').clear();
    await page.locator('input[id="pickupAddress-location-input"]').fill('90212');
    await expect(page.getByText(pickupLocation, { exact: true })).toBeVisible();
    await page.keyboard.press('Enter');

    // Secondary pickup address
    await pickupAddress.getByText('Yes').click();
    await pickupAddress.getByLabel('Address 1').nth(1).fill('8 Q St');
    await pickupAddress.getByLabel('Address 2').nth(1).clear();
    await page.locator('input[id="secondaryPickupAddress-location-input"]').fill('90212');
    await expect(page.getByText(pickupLocation, { exact: true })).toBeVisible();
    await page.keyboard.press('Enter');

    const deliveryLocation = 'HOLLYWOOD, MD 20636 (SAINT MARYS)';
    const deliveryAddress = page.getByRole('group', { name: 'Delivery location' });
    await deliveryAddress.getByText('Yes').nth(0).click();
    await deliveryAddress.getByLabel('Address 1').nth(0).fill('9 W 2nd Ave');
    await deliveryAddress.getByLabel('Address 2').nth(0).fill('P.O. Box 456');
    await page.locator('input[id="deliveryAddress-location-input"]').fill('20636');
    await expect(page.getByText(deliveryLocation, { exact: true })).toBeVisible();
    await page.keyboard.press('Enter');

    // Secondary delivery address
    const secondaryDeliveryLocation = 'ATCO, NJ 08004 (CAMDEN)';
    await deliveryAddress.getByText('Yes').nth(1).click();
    await deliveryAddress.getByLabel('Address 1').nth(1).fill('9 Q St');
    await deliveryAddress.getByLabel('Address 2').nth(1).clear();
    await page.locator('input[id="secondaryDeliveryAddress-location-input"]').fill('08004');
    await expect(page.getByText(secondaryDeliveryLocation, { exact: true })).toBeVisible();
    await page.keyboard.press('Enter');
    await customerPage.navigateForward();

    // Verify that shipment updated
    await customerPage.waitForPage.reviewShipments();
    await expect(page.getByTestId('ShipmentContainer').getByText('7 Q St')).toBeVisible();
    await expect(page.getByTestId('ShipmentContainer').getByText('8 Q St')).toBeVisible();
    await expect(page.getByTestId('ShipmentContainer').getByText('9 W 2nd Ave')).toBeVisible();
    await expect(page.getByTestId('ShipmentContainer').getByText('9 Q St')).toBeVisible();

    // Navigate to homepage and delete shipment
    await customerPage.navigateBack();
    await customerPage.waitForPage.home();
    // Remove secondary pickup and delivery addresses
    await page.getByTestId('shipment-list-item-container').getByRole('button', { name: 'Edit' }).click();
    await customerPage.waitForPage.hhgShipment();
    await pickupAddress.getByText('No').click();
    await deliveryAddress.getByText('No', { exact: true }).nth(1).click();
    await customerPage.navigateForward();

    await customerPage.waitForPage.reviewShipments();
    await expect(page.getByTestId('ShipmentContainer').getByText('7 Q St')).toBeVisible();
    // Make sure secondary pickup and delivery addresses are gone now
    await expect(page.getByTestId('ShipmentContainer').getByText('8 Q St')).toBeHidden();
    await expect(page.getByTestId('ShipmentContainer').getByText('9 Q St')).toBeHidden();
    await customerPage.navigateBack();
    await customerPage.waitForPage.home();
    await page.getByRole('button', { name: 'Delete' }).click();
    await page.getByTestId('modal').getByTestId('button').click();

    await expect(page.getByText('The shipment was deleted.')).toBeVisible();
    await expect(page.getByTestId('stepContainer3').getByText('Set up shipments')).toBeVisible();
  });
});

test.describe('(MultiMove) HHG', () => {
  test.skip(multiMoveEnabled === 'false', 'Skip if MultiMove workflow is not enabled.');

  test('A customer can create, edit, and delete an HHG shipment', async ({ page, customerPage }) => {
    // Generate a new onboarded user with orders and log in
    const move = await customerPage.testHarness.buildMoveWithOrders();
    const userId = move.Orders.ServiceMember.user_id;
    await customerPage.signInAsExistingCustomer(userId);

    // Navigate from MM Dashboard to Move
    await customerPage.navigateFromMMDashboardToMove(move);

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
    const location = 'ATCO, NJ 08004 (CAMDEN)';
    const pickupAddress = page.getByRole('group', { name: 'Pickup location' });
    await pickupAddress.getByLabel('Address 1').fill('7 Q St');
    await pickupAddress.getByLabel('Address 2').clear();
    await page.locator('input[id="pickup.address-location-input"]').fill('08004');
    await expect(page.getByText(location, { exact: true })).toBeVisible();
    await page.keyboard.press('Enter');

    // Secondary pickup address
    await pickupAddress.getByText('Yes').click();
    await pickupAddress.getByLabel('Address 1').nth(1).fill('8 Q St');
    await pickupAddress.getByLabel('Address 2').nth(1).clear();
    await page.locator('input[id="secondaryPickup.address-location-input"]').fill('08004');
    await expect(page.getByText(location, { exact: true })).toBeVisible();
    await page.keyboard.press('Enter');

    // Delivery address
    const deliveryLocation = 'HOLLYWOOD, MD 20636 (SAINT MARYS)';
    const deliveryAddress = page.getByRole('group', { name: 'Delivery location' });
    await deliveryAddress.getByText('Yes').nth(0).click();
    await deliveryAddress.getByLabel('Address 1').nth(0).fill('9 W 2nd Ave');
    await deliveryAddress.getByLabel('Address 2').nth(0).fill('P.O. Box 456');
    await page.locator('input[id="delivery.address-location-input"]').fill('20636');
    await expect(page.getByText(deliveryLocation, { exact: true })).toBeVisible();
    await page.keyboard.press('Enter');

    // Secondary delivery address
    await deliveryAddress.getByText('Yes').nth(1).click();
    await deliveryAddress.getByLabel('Address 1').nth(1).fill('9 Q St');
    await deliveryAddress.getByLabel('Address 2').nth(1).clear();
    await page.locator('input[id="secondaryDelivery.address-location-input"]').fill('08004');
    await expect(page.getByText(location, { exact: true })).toBeVisible();
    await page.keyboard.press('Enter');
    await customerPage.navigateForward();

    // Verify that shipment updated
    await customerPage.waitForPage.reviewShipments();
    await expect(page.getByTestId('ShipmentContainer').getByText('7 Q St')).toBeVisible();
    await expect(page.getByTestId('ShipmentContainer').getByText('8 Q St')).toBeVisible();
    await expect(page.getByTestId('ShipmentContainer').getByText('9 W 2nd Ave')).toBeVisible();
    await expect(page.getByTestId('ShipmentContainer').getByText('9 Q St')).toBeVisible();

    // Navigate to homepage and delete shipment
    await customerPage.navigateBack();
    await customerPage.waitForPage.home();
    // Remove secondary pickup and delivery addresses
    await page.getByTestId('shipment-list-item-container').getByRole('button', { name: 'Edit' }).click();
    await customerPage.waitForPage.hhgShipment();
    await pickupAddress.getByText('No').click();
    await deliveryAddress.getByText('No', { exact: true }).nth(1).click();
    await customerPage.navigateForward();

    await customerPage.waitForPage.reviewShipments();
    await expect(page.getByTestId('ShipmentContainer').getByText('7 Q St')).toBeVisible();
    // Make sure secondary pickup and delivery addresses are gone now
    await expect(page.getByTestId('ShipmentContainer').getByText('8 Q St')).toBeHidden();
    await expect(page.getByTestId('ShipmentContainer').getByText('9 Q St')).toBeHidden();
    await customerPage.navigateBack();
    await customerPage.waitForPage.home();
    await page.getByRole('button', { name: 'Delete' }).click();
    await page.getByTestId('modal').getByTestId('button').click();

    await expect(page.getByText('The shipment was deleted.')).toBeVisible();
    await expect(page.getByTestId('stepContainer3').getByText('Set up shipments')).toBeVisible();
  });
});
