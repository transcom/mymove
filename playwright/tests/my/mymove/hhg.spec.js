// @ts-check
import { test, expect } from '../../utils/my/customerTest';

const multiMoveEnabled = process.env.FEATURE_FLAG_MULTI_MOVE;
const alaskaFF = process.env.FEATURE_FLAG_ENABLE_ALASKA;

test.describe('HHG', () => {
  test.skip(multiMoveEnabled === 'true', 'Skip if MultiMove workflow is enabled.');
  test('A customer can create, edit, and delete an HHG shipment', async ({ page, customerPage }) => {
    // Generate a new onboarded user with orders and log in
    const move = await customerPage.testHarness.buildMoveWithOrders();
    const userId = move?.Orders?.service_member?.user_id;
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
    const pickupAddress = page.getByRole('group', { name: 'Pickup Address' });
    await pickupAddress.getByLabel('Address 1').fill('7 Q St');
    await pickupAddress.getByLabel('Address 2').clear();
    await page.locator('input[id="pickupAddress-input"]').fill('90212');
    await expect(page.getByText(pickupLocation, { exact: true })).toBeVisible();
    await page.keyboard.press('Enter');

    // Secondary pickup address
    await pickupAddress.getByText('Yes').click();
    await pickupAddress.getByLabel('Address 1').nth(1).fill('8 Q St');
    await pickupAddress.getByLabel('Address 2').nth(1).clear();
    await page.locator('input[id="secondaryPickupAddress-input"]').fill('90212');
    await expect(page.getByText(pickupLocation, { exact: true })).toBeVisible();
    await page.keyboard.press('Enter');

    const deliveryLocation = 'HOLLYWOOD, MD 20636 (SAINT MARYS)';
    const deliveryAddress = page.getByRole('group', { name: 'Delivery Address' });
    await deliveryAddress.getByText('Yes').nth(0).click();
    await deliveryAddress.getByLabel('Address 1').nth(0).fill('9 W 2nd Ave');
    await deliveryAddress.getByLabel('Address 2').nth(0).fill('P.O. Box 456');
    await page.locator('input[id="deliveryAddress-input"]').fill('20636');
    await expect(page.getByText(deliveryLocation, { exact: true })).toBeVisible();
    await page.keyboard.press('Enter');

    // Secondary delivery address
    const secondaryDeliveryLocation = 'ATCO, NJ 08004 (CAMDEN)';
    await deliveryAddress.getByText('Yes').nth(1).click();
    await deliveryAddress.getByLabel('Address 1').nth(1).fill('9 Q St');
    await deliveryAddress.getByLabel('Address 2').nth(1).clear();
    await page.locator('input[id="secondaryDeliveryAddress-input"]').fill('08004');
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
    const userId = move?.Orders?.service_member?.user_id;
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
    const pickupLocation = 'BEVERLY HILLS, CA 90210 (LOS ANGELES)';
    const secondaryPickupLocation = 'YUMA, AZ 85364 (YUMA)';
    const deliveryLocation = 'YUMA, AZ 85367 (YUMA)';
    const secondaryDeliveryLocation = 'YUMA, AZ 85366 (YUMA)';

    const pickupAddress = page.getByRole('group', { name: 'Pickup Address' });
    await pickupAddress.getByLabel('Address 1').fill('7 Q St');
    await pickupAddress.getByLabel('Address 2').clear();
    await page.locator('input[id="pickup.address-input"]').fill('90210');
    await expect(page.getByText(pickupLocation, { exact: true })).toBeVisible();
    await page.keyboard.press('Enter');

    // Secondary pickup address
    await pickupAddress.getByText('Yes').click();
    await pickupAddress.getByLabel('Address 1').nth(1).fill('8 Q St');
    await pickupAddress.getByLabel('Address 2').nth(1).clear();
    await page.locator('input[id="secondaryPickup.address-input"]').fill('85364');
    await expect(page.getByText(secondaryPickupLocation, { exact: true })).toBeVisible();
    await page.keyboard.press('Enter');

    // Delivery address
    const deliveryAddress = page.getByRole('group', { name: 'Delivery Address' });
    await deliveryAddress.getByText('Yes').nth(0).click();
    await deliveryAddress.getByLabel('Address 1').nth(0).fill('9 W 2nd Ave');
    await deliveryAddress.getByLabel('Address 2').nth(0).fill('P.O. Box 456');
    await page.locator('input[id="delivery.address-input"]').fill('85367');
    await expect(page.getByText(deliveryLocation, { exact: true })).toBeVisible();
    await page.keyboard.press('Enter');

    // Secondary delivery address
    await deliveryAddress.getByText('Yes').nth(1).click();
    await deliveryAddress.getByLabel('Address 1').nth(1).fill('9 Q St');
    await deliveryAddress.getByLabel('Address 2').nth(1).clear();
    await page.locator('input[id="secondaryDelivery.address-input"]').fill('85366');
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

  test.skip(alaskaFF === 'false', 'Skip if the create customer & AK FFs are not enabled.');
  test('A customer can create, edit, and submit an international Alaska HHG shipment', async ({
    page,
    customerPage,
  }) => {
    // Generate a new onboarded user with orders and log in
    const move = await customerPage.testHarness.buildMoveWithOrders();
    const userId = move?.Orders?.service_member?.user_id;
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
    await page.getByLabel('Preferred pickup date').fill('29 Dec 2025');
    await page.getByLabel('Preferred pickup date').blur();
    await page.getByText('Use my current address').click();
    await page.getByLabel('Preferred delivery date').fill('29 Dec 2025');
    await page.getByLabel('Preferred delivery date').blur();
    await page.getByTestId('remarks').fill('Going to Alaska');
    await customerPage.navigateForward();

    // Verify that form submitted, initial setup has it being a domestic HHG shipment (dHHG)
    await customerPage.waitForPage.reviewShipments();
    await expect(page.getByText('dHHG')).toBeVisible();
    await expect(page.getByText('Going to Alaska')).toBeVisible();
    await expect(page.getByTestId('ShipmentContainer').getByText('123 Any Street')).toBeVisible();

    // Navigate to edit shipment from the review page
    await page.getByTestId('edit-shipment-btn').click();
    await customerPage.waitForPage.hhgShipment();

    // Update form (adding pickup and delivery address)
    const pickupLocation = 'LAWTON, OK 73505 (COMANCHE)';
    const pickupAddress = page.getByRole('group', { name: 'Pickup Address' });
    await pickupAddress.getByLabel('Address 1').fill('123 Warm St.');
    await page.locator('input[id="pickup.address-input"]').fill('73505');
    await expect(page.getByText(pickupLocation, { exact: true })).toBeVisible();
    await page.keyboard.press('Enter');

    // Delivery address
    const deliveryLocation = 'JBER, AK 99505 (ANCHORAGE)';
    const deliveryAddress = page.getByRole('group', { name: 'Delivery Address' });
    await deliveryAddress.getByText('Yes').nth(0).click();
    await deliveryAddress.getByLabel('Address 1').nth(0).fill('123 Cold Ave.');
    await page.locator('input[id="delivery.address-input"]').fill('99505');
    await expect(page.getByText(deliveryLocation, { exact: true })).toBeVisible();
    await page.keyboard.press('Enter');
    await customerPage.navigateForward();

    // Verify that shipment updated - should now be an iHHG shipment
    await customerPage.waitForPage.reviewShipments();
    await expect(page.getByText('iHHG')).toBeVisible();
    await expect(page.getByTestId('ShipmentContainer').getByText('123 Warm St.')).toBeVisible();
    await expect(page.getByTestId('ShipmentContainer').getByText('123 Cold Ave.')).toBeVisible();

    await page.getByRole('button', { name: 'Next' }).click();
    await expect(page).toHaveURL(/\/moves\/[^/]+\/agreement/);
    await expect(page.getByRole('heading', { name: 'Now for the official part…' })).toBeVisible();

    const scrollBox = page.locator('[data-testid="certificationTextScrollBox"]');
    const signatureBox = page.getByRole('textbox', { name: 'signature' });
    // Make sure it's visible
    await expect(scrollBox).toBeVisible();

    // Gradual scroll to bottom to trigger the React onScroll logic
    await scrollBox.evaluate(async (el) => {
      const scrWindow = el;
      const delay = (ms) =>
        new Promise((res) => {
          setTimeout(res, ms);
        });
      for (let i = 0; i <= scrWindow.scrollHeight; i += 100) {
        scrWindow.scrollTop = i;
        await delay(50);
      }
    });

    const checkbox = page.locator('[data-testid="acknowledgementCheckbox"]');
    await expect(checkbox).toBeEnabled();

    // Click it to acknowledge
    await checkbox.click();

    await expect(signatureBox).toBeEnabled();

    await page.locator('input[name="signature"]').fill('Leo Spacemen');
    await expect(page.getByRole('button', { name: 'Complete' })).toBeEnabled();
    await page.getByRole('button', { name: 'Complete' }).click();
    await expect(page.getByText('submitted your move request.')).toBeVisible();
  });
});
