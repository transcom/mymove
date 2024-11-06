// @ts-check
import { test, expect } from './servicesCounselingTestFixture';

const today = new Date();
const pickupDate = today;
const pickupDateString = pickupDate.toLocaleDateString('en-US');
const deliveryDate = new Date(new Date().setDate(today.getDate() + 14));
const deliveryDateString = deliveryDate.toLocaleDateString('en-US');

const pickupAddress = {
  Address1: '7 Q St',
  City: 'Atco',
  State: 'NJ',
  ZIP: '08004',
};

const secondaryPickupAddress = {
  ...pickupAddress,
};
secondaryPickupAddress.Address1 = '8 Q St';

const addressToString = (address) => {
  return `${address.Address1},  ${address.Address2 ? `${address.Address2}, ` : ''}${
    address.Address3 ? `${address.Address3}, ` : ''
  }${address.City}, ${address.State} ${address.ZIP}`;
};

const deliveryAddress = {
  Address1: '9 W 2nd Ave',
  Address2: 'P.O. Box 456',
  City: 'Hollywood',
  State: 'MD',
  ZIP: '20636',
};

const secondaryDeliveryAddress = {
  Address1: '9 Q St',
  City: 'Atco',
  State: 'NJ',
  ZIP: '08004',
};

const releasingAgent = {
  firstName: 'Grace',
  lastName: 'Griffin',
  phone: '2025551234',
  email: 'grace.griffin@example.com',
};

const receivingAgent = {
  firstName: 'Leo',
  lastName: 'Spacemen',
  phone: '2025552345',
  email: 'leo.spaceman@example.com',
};

const formatPhone = (phone) => {
  return `${phone.slice(0, 3)}-${phone.slice(3, 6)}-${phone.slice(6)}`;
};

const agentToString = (agent) => {
  return `${agent.firstName} ${agent.lastName}${formatPhone(agent.phone)}${agent.email}`;
};

const formatDate = (date) => {
  const formattedDay = date.toLocaleDateString(undefined, { day: '2-digit' });
  const formattedMonth = date.toLocaleDateString(undefined, {
    month: 'short',
  });
  const formattedYear = date.toLocaleDateString(undefined, {
    year: 'numeric',
  });

  return `${formattedDay} ${formattedMonth} ${formattedYear}`;
};

test.describe('Services counselor user', () => {
  test.beforeEach(async ({ scPage }) => {
    const move = await scPage.testHarness.buildMobileHomeMoveNeedsSC();
    await scPage.navigateToMove(move.locator);
  });

  test('Services Counselor can create a mobile home shipment and view shipment card info', async ({ page, scPage }) => {
    await page.getByTestId('dropdown').selectOption({ label: 'Mobile Home' });

    await expect(page.getByText('Mobile Home Information')).toBeVisible();

    await expect(page.getByRole('heading', { level: 1 })).toHaveText('Add shipment details');
    await expect(page.getByTestId('tag')).toHaveText('Mobile Home');

    await page.getByLabel('Year').fill('2022');
    await page.getByLabel('Make').fill('make');
    await page.getByLabel('Model').fill('model');
    await page.getByTestId('lengthFeet').fill('22');
    await page.getByTestId('lengthInches').fill('0');
    await page.getByTestId('widthFeet').fill('22');
    await page.getByTestId('widthInches').fill('0');
    await page.getByTestId('heightFeet').fill('22');
    await page.getByTestId('heightInches').fill('0');

    await page.locator('#requestedPickupDate').fill(pickupDateString);
    await page.locator('#requestedPickupDate').blur();
    await page.getByText('Use pickup address').click();
    await page.locator('#requestedDeliveryDate').fill('16 Mar 2022');
    await page.locator('#requestedDeliveryDate').blur();

    await page.getByLabel('Counselor remarks').fill('Sample counselor remarks');

    await page.locator('#requestedDeliveryDate').fill(deliveryDateString);
    await page.locator('#requestedDeliveryDate').blur();

    // Save the shipment
    await page.getByRole('button', { name: 'Save' }).click();
    await scPage.waitForPage.moveDetails();

    await expect(page.getByTestId('ShipmentContainer')).toHaveCount(2);

    await expect(page.getByText('Mobile home year').last()).toBeVisible();
    await expect(page.getByTestId('year').last()).toHaveText('2022');
    await expect(page.getByText('Mobile home make').last()).toBeVisible();
    await expect(page.getByTestId('make').last()).toHaveText('make');
    await expect(page.getByText('Mobile home model').last()).toBeVisible();
    await expect(page.getByTestId('model').last()).toHaveText('model');
    await expect(page.getByText('Dimensions').last()).toBeVisible();
    await expect(page.getByTestId('dimensions').last()).toHaveText("22' L x 22' W x 22' H");
  });

  test('Services Counselor can delete an existing Mobile Home shipment', async ({ page, scPage }) => {
    await expect(page.getByText('Edit Shipment')).toHaveCount(1);
    // Choose a shipment and store it's shipment ID
    const editShipmentButton = await page.getByRole('button', { name: 'Edit Shipment' });
    process.stdout.write(await editShipmentButton.evaluate((el) => el.outerHTML));

    await editShipmentButton.click();
    await scPage.waitForLoading();
    await scPage.waitForPage.editMobileHomeShipment();

    // Delete that shipment
    await page.getByRole('button', { name: 'Delete shipment' }).click();
    await expect(page.getByTestId('modalCloseButton')).toBeVisible();
    await page.getByTestId('modal').getByRole('button', { name: 'Delete shipment' }).click();
    await scPage.waitForPage.moveDetails();

    // Verify that the shipment has been deleted
    await expect(page.getByTestId('ShipmentContainer')).toHaveCount(0);
  });

  test('Services Counselor can edit an existing Mobile Home shipment', async ({ page, scPage }) => {
    await expect(page.getByText('Edit Shipment')).toHaveCount(1);

    // Choose a shipment, store it's container, and click the edit button
    const shipmentContainer = await page.getByTestId('ShipmentContainer');
    await shipmentContainer.getByRole('button').click();
    await scPage.waitForLoading();
    await scPage.waitForPage.editMobileHomeShipment();

    // Fill in all of the form fields with new data
    await page.getByLabel('Year').fill('2024');
    await page.getByLabel('Make').fill('Test Make');
    await page.getByLabel('Model').fill('Test Model');

    await page.getByTestId('lengthFeet').fill('20');
    await page.getByTestId('lengthInches').fill('6');

    await page.getByTestId('widthFeet').fill('15');
    await page.getByTestId('widthInches').fill('1');

    await page.getByTestId('heightFeet').fill('10');
    await page.getByTestId('heightInches').fill('0');

    await page.locator('#requestedPickupDate').fill(pickupDateString);
    await page.locator('#requestedPickupDate').blur();
    await page.locator('#requestedDeliveryDate').fill(deliveryDateString);
    await page.locator('#requestedDeliveryDate').blur();

    // Update form (adding pickup and delivery address)
    const pickupAddressGroup = await page.getByRole('group', { name: 'Pickup location' });
    await pickupAddressGroup.getByText('Yes').click();
    await pickupAddressGroup.getByLabel('Address 1').nth(0).fill(pickupAddress.Address1);
    await pickupAddressGroup.getByLabel('Address 2').nth(0).clear();
    await pickupAddressGroup.getByLabel('Address 3').nth(0).clear();
    await pickupAddressGroup.getByLabel('City').nth(0).fill(pickupAddress.City);
    await pickupAddressGroup.getByLabel('State').nth(0).selectOption({ label: pickupAddress.State });
    await pickupAddressGroup.getByLabel('ZIP').nth(0).fill(pickupAddress.ZIP);

    // Secondary pickup address
    await pickupAddressGroup.getByText('Yes').click();
    await pickupAddressGroup.getByLabel('Address 1').nth(1).fill(secondaryPickupAddress.Address1);
    await pickupAddressGroup.getByLabel('Address 2').nth(1).clear();
    await pickupAddressGroup.getByLabel('Address 3').nth(1).clear();
    await pickupAddressGroup.getByLabel('City').nth(1).fill(secondaryPickupAddress.City);
    await pickupAddressGroup.getByLabel('State').nth(1).selectOption({ label: secondaryPickupAddress.State });
    await pickupAddressGroup.getByLabel('ZIP').nth(1).fill(secondaryPickupAddress.ZIP);

    // Releasing agent
    await page.locator(`[name='pickup.agent.firstName']`).fill(releasingAgent.firstName);
    await page.locator(`[name='pickup.agent.lastName']`).fill(releasingAgent.lastName);
    await page.locator(`[name='pickup.agent.phone']`).fill(releasingAgent.phone);
    await page.locator(`[name='pickup.agent.email']`).fill(releasingAgent.email);

    const deliveryAddressGroup = await page.getByRole('group', { name: 'Delivery location' });
    await deliveryAddressGroup.getByText('Yes').nth(0).click();
    await deliveryAddressGroup.getByLabel('Address 1').nth(0).fill(deliveryAddress.Address1);
    await deliveryAddressGroup.getByLabel('Address 2').nth(0).fill(deliveryAddress.Address2);
    await deliveryAddressGroup.getByLabel('Address 3').nth(0).clear();
    await deliveryAddressGroup.getByLabel('City').nth(0).fill(deliveryAddress.City);
    await deliveryAddressGroup.getByLabel('State').nth(0).selectOption({ label: deliveryAddress.State });
    await deliveryAddressGroup.getByLabel('ZIP').nth(0).fill(deliveryAddress.ZIP);

    // Secondary delivery address
    await deliveryAddressGroup.getByText('Yes').nth(1).click();
    await deliveryAddressGroup.getByLabel('Address 1').nth(1).fill(secondaryDeliveryAddress.Address1);
    await deliveryAddressGroup.getByLabel('Address 2').nth(1).clear();
    await deliveryAddressGroup.getByLabel('Address 3').nth(1).clear();
    await deliveryAddressGroup.getByLabel('City').nth(1).fill(secondaryDeliveryAddress.City);
    await deliveryAddressGroup.getByLabel('State').nth(1).selectOption({ label: secondaryDeliveryAddress.State });
    await deliveryAddressGroup.getByLabel('ZIP').nth(1).fill(secondaryDeliveryAddress.ZIP);

    // Receiving agent
    await page.locator(`[name='delivery.agent.firstName']`).fill(receivingAgent.firstName);
    await page.locator(`[name='delivery.agent.lastName']`).fill(receivingAgent.lastName);
    await page.locator(`[name='delivery.agent.phone']`).fill(receivingAgent.phone);
    await page.locator(`[name='delivery.agent.email']`).fill(receivingAgent.email);

    await page.getByLabel('Counselor remarks').fill('Sample counselor remarks');

    // Submit edits
    await page.getByTestId('submitForm').click();
    await scPage.waitForLoading();
    await expect(page.locator('.usa-alert__text')).toContainText('Your changes were saved.');

    // Check that the data in the shipment card now matches what we just submitted
    await shipmentContainer.locator('[data-prefix="fas"][data-icon="chevron-down"]').click();
    await expect(shipmentContainer.getByTestId('requestedPickupDate')).toHaveText(formatDate(pickupDate));
    await expect(shipmentContainer.getByTestId('pickupAddress')).toHaveText(addressToString(pickupAddress));
    await expect(shipmentContainer.getByTestId('secondaryPickupAddress')).toHaveText(
      addressToString(secondaryPickupAddress),
    );

    await expect(shipmentContainer.getByTestId('RELEASING_AGENT')).toHaveText(agentToString(releasingAgent));

    await expect(shipmentContainer.getByTestId('requestedDeliveryDate')).toHaveText(formatDate(deliveryDate));
    await expect(shipmentContainer.getByTestId('destinationAddress')).toHaveText(addressToString(deliveryAddress));
    await expect(shipmentContainer.getByTestId('secondaryDeliveryAddress')).toHaveText(
      addressToString(secondaryDeliveryAddress),
    );

    await expect(shipmentContainer.getByTestId('RECEIVING_AGENT')).toHaveText(agentToString(receivingAgent));

    await expect(shipmentContainer.getByTestId('year')).toHaveText('2024');
    await expect(shipmentContainer.getByTestId('make')).toHaveText('Test Make');
    await expect(shipmentContainer.getByTestId('model')).toHaveText('Test Model');

    await expect(shipmentContainer.getByTestId('dimensions')).toHaveText(`20' 6" L x 15' 1" W x 10' H`);

    await expect(shipmentContainer.getByTestId('counselorRemarks')).toHaveText('Sample counselor remarks');
  });
});
