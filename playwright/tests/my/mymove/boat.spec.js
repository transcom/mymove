// @ts-check
import { test, expect } from '../../utils/my/customerTest';

const multiMoveEnabled = process.env.FEATURE_FLAG_MULTI_MOVE;

test.describe('Boat shipment', () => {
  test.skip(multiMoveEnabled === 'true', 'Skip if MultiMove workflow is enabled.');
  test('A customer can create a Boat shipment - Tow-Away', async ({ page, customerPage }) => {
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

    // Create an Boat shipment
    await page.getByText('Move a boat').click();
    await customerPage.navigateForward();

    // Fill in form to create Boat shipment
    await customerPage.waitForPage.boatShipment();
    await page.getByLabel('Year').fill('2022');
    await page.getByLabel('Make').fill('make');
    await page.getByLabel('Model').fill('model');
    await page.getByTestId('lengthFeet').fill('22');
    await page.getByTestId('widthFeet').fill('22');
    await page.getByTestId('heightFeet').fill('22');
    // await page.getByLabel('hasTrailerYes').click();
    await page.locator('label[for="hasTrailerYes"]').click();
    await page.locator('label[for="isRoadworthyYes"]').click();
    await page.getByTestId('remarks').fill('remarks test');
    await page.getByRole('button', { name: 'Continue' }).click();

    await expect(page.getByRole('heading', { name: 'Boat Tow-Away (BTA)' })).toBeVisible();
    await page.getByTestId('boatConfirmationContinue').click();

    await expect(page.getByText('Pickup info')).toBeVisible();
    await page.getByLabel('Preferred pickup date').fill('25 Dec 2022');
    await page.getByLabel('Preferred pickup date').blur();
    await page.getByText('Use my current address').click();
    await page.getByLabel('Preferred delivery date').fill('25 Dec 2022');
    await page.getByLabel('Preferred delivery date').blur();
    await page.getByRole('button', { name: 'Save & Continue' }).click();
    await customerPage.waitForPage.reviewShipments();
  });

  test('A customer can create a Boat shipment - Haul-Away', async ({ page, customerPage }) => {
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

    // Create an Boat shipment
    await page.getByText('Move a boat').click();
    await customerPage.navigateForward();

    // Fill in form to create Boat shipment
    await customerPage.waitForPage.boatShipment();
    await page.getByLabel('Year').fill('2022');
    await page.getByLabel('Make').fill('make');
    await page.getByLabel('Model').fill('model');
    await page.getByTestId('lengthFeet').fill('22');
    await page.getByTestId('widthFeet').fill('22');
    await page.getByTestId('heightFeet').fill('22');
    await page.locator('label[for="hasTrailerNo"]').click();
    await page.getByTestId('remarks').fill('remarks test');
    await page.getByRole('button', { name: 'Continue' }).click();

    await expect(page.getByRole('heading', { name: 'Boat Haul-Away (BHA)' })).toBeVisible();
    await page.getByTestId('boatConfirmationContinue').click();

    await expect(page.getByText('Pickup info')).toBeVisible();
    await page.getByLabel('Preferred pickup date').fill('25 Dec 2022');
    await page.getByLabel('Preferred pickup date').blur();
    await page.getByText('Use my current address').click();
    await page.getByLabel('Preferred delivery date').fill('25 Dec 2022');
    await page.getByLabel('Preferred delivery date').blur();
    await page.getByRole('button', { name: 'Save & Continue' }).click();
    await customerPage.waitForPage.reviewShipments();
  });

  test('A customer is redirected to HHG if dimension requirement is not met', async ({ page, customerPage }) => {
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

    // Create an Boat shipment
    await page.getByText('Move a boat').click();
    await customerPage.navigateForward();

    // Fill in form to create Boat shipment
    await customerPage.waitForPage.boatShipment();
    await page.getByLabel('Year').fill('2022');
    await page.getByLabel('Make').fill('make');
    await page.getByLabel('Model').fill('model');
    await page.getByTestId('lengthFeet').fill('2');
    await page.getByTestId('widthFeet').fill('2');
    await page.getByTestId('heightFeet').fill('2');
    await page.locator('label[for="hasTrailerNo"]').click();
    await page.getByTestId('remarks').fill('remarks test');
    await page.getByRole('button', { name: 'Continue' }).click();

    await expect(
      page.getByRole('heading', { name: 'Movers pack and ship it, paid by the government (HHG)' }),
    ).toBeVisible();
    await page.getByTestId('boatConfirmationContinue').click();

    await expect(page.getByText('HHG')).toBeVisible();
  });

  test('Is able to delete a boat shipment', async ({ page, customerPage }) => {
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

    // Create an Boat shipment
    await page.getByText('Move a boat').click();
    await customerPage.navigateForward();

    // Fill in form to create Boat shipment
    await customerPage.waitForPage.boatShipment();
    await page.getByLabel('Year').fill('2022');
    await page.getByLabel('Make').fill('make');
    await page.getByLabel('Model').fill('model');
    await page.getByTestId('lengthFeet').fill('22');
    await page.getByTestId('widthFeet').fill('22');
    await page.getByTestId('heightFeet').fill('22');
    await page.locator('label[for="hasTrailerYes"]').click();
    await page.locator('label[for="isRoadworthyYes"]').click();
    await page.getByTestId('remarks').fill('remarks test');
    await page.getByRole('button', { name: 'Continue' }).click();

    await expect(page.getByRole('heading', { name: 'Boat Tow-Away (BTA)' })).toBeVisible();
    await page.getByTestId('boatConfirmationContinue').click();

    await expect(page.getByText('Pickup info')).toBeVisible();
    await page.getByLabel('Preferred pickup date').fill('25 Dec 2022');
    await page.getByLabel('Preferred pickup date').blur();
    await page.getByText('Use my current address').click();
    await page.getByLabel('Preferred delivery date').fill('25 Dec 2022');
    await page.getByLabel('Preferred delivery date').blur();
    await page.getByRole('button', { name: 'Save & Continue' }).click();
    await customerPage.waitForPage.reviewShipments();

    await expect(page.getByRole('heading', { name: 'Boat 1' })).toBeVisible();
    await page.getByTestId('deleteShipmentButton').click();
    await expect(page.getByRole('heading', { name: 'Delete this?' })).toBeVisible();
    await page.getByText('Yes, Delete').click();

    await expect(page.getByRole('heading', { name: 'Boat 1' })).not.toBeVisible();
  });

  test('Deletes existing boat shipment and is redirected to HHG after edit', async ({ page, customerPage }) => {
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

    // Create an Boat shipment
    await page.getByText('Move a boat').click();
    await customerPage.navigateForward();

    // Fill in form to create Boat shipment
    await customerPage.waitForPage.boatShipment();
    await page.getByLabel('Year').fill('2022');
    await page.getByLabel('Make').fill('make');
    await page.getByLabel('Model').fill('model');
    await page.getByTestId('lengthFeet').fill('22');
    await page.getByTestId('widthFeet').fill('22');
    await page.getByTestId('heightFeet').fill('22');
    await page.locator('label[for="hasTrailerYes"]').click();
    await page.locator('label[for="isRoadworthyYes"]').click();
    await page.getByTestId('remarks').fill('remarks test');
    await page.getByRole('button', { name: 'Continue' }).click();

    await expect(page.getByRole('heading', { name: 'Boat Tow-Away (BTA)' })).toBeVisible();
    await page.getByTestId('boatConfirmationContinue').click();

    await expect(page.getByText('Pickup info')).toBeVisible();
    await page.getByLabel('Preferred pickup date').fill('25 Dec 2022');
    await page.getByLabel('Preferred pickup date').blur();
    await page.getByText('Use my current address').click();
    await page.getByLabel('Preferred delivery date').fill('25 Dec 2022');
    await page.getByLabel('Preferred delivery date').blur();
    await page.getByRole('button', { name: 'Save & Continue' }).click();
    await customerPage.waitForPage.reviewShipments();

    await expect(page.getByRole('heading', { name: 'Boat 1' })).toBeVisible();
    await page.getByTestId('editShipmentButton').click();

    await customerPage.waitForPage.boatShipment();
    await page.getByTestId('lengthFeet').fill('2');
    await page.getByTestId('widthFeet').fill('2');
    await page.getByTestId('heightFeet').fill('2');
    await page.getByRole('button', { name: 'Continue' }).click();

    await expect(
      page.getByRole('heading', { name: 'Movers pack and ship it, paid by the government (HHG)' }),
    ).toBeVisible();
    await page.getByTestId('boatConfirmationContinue').click();
    await expect(
      page.getByRole('heading', { name: 'Movers pack and ship it, paid by the government (HHG)' }),
    ).not.toBeVisible();
    await expect(page.getByText('HHG')).toBeVisible();
    await expect(page.getByText('Movers pack and transport this shipment')).toBeVisible();
    await page.getByTestId('wizardNextButton').click();
    await customerPage.waitForPage.reviewShipments();
    await expect(page.getByRole('heading', { name: 'Boat 1' })).not.toBeVisible();
    await expect(page.getByRole('heading', { name: 'HHG 1' })).toBeVisible();
  });

  test('A customer is unable to sign if boat shipment is incomplete', async ({ page, customerPage }) => {
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

    // Create an Boat shipment
    await page.getByText('Move a boat').click();
    await customerPage.navigateForward();

    // Fill in form to create Boat shipment
    await customerPage.waitForPage.boatShipment();
    await page.getByLabel('Year').fill('2022');
    await page.getByLabel('Make').fill('make');
    await page.getByLabel('Model').fill('model');
    await page.getByTestId('lengthFeet').fill('22');
    await page.getByTestId('widthFeet').fill('22');
    await page.getByTestId('heightFeet').fill('22');
    await page.locator('label[for="hasTrailerYes"]').click();
    await page.locator('label[for="isRoadworthyYes"]').click();
    await page.getByTestId('remarks').fill('remarks test');
    await page.getByRole('button', { name: 'Continue' }).click();

    await expect(page.getByRole('heading', { name: 'Boat Tow-Away (BTA)' })).toBeVisible();
    await page.getByTestId('boatConfirmationContinue').click();

    // Back out to mark as incomplete
    await expect(page.getByText('Pickup info')).toBeVisible();
    await page.getByRole('button', { name: 'Back' }).click();
    await customerPage.waitForPage.boatShipment();
    await page.getByRole('button', { name: 'Back' }).click();

    await expect(page.getByText('Time to submit your move')).toBeVisible();
    await expect(page.getByTestId('review-and-submit-btn')).toBeDisabled();

    // Create another shipment to test continue button is disabled from the review page
    await page.getByTestId('shipment-selection-btn').click();
    await customerPage.waitForPage.selectShipmentType();

    // Create an Boat shipment
    await page.getByText('Move a boat').click();
    await customerPage.navigateForward();

    // Fill in form to create Boat shipment
    await customerPage.waitForPage.boatShipment();
    await page.getByLabel('Year').fill('2022');
    await page.getByLabel('Make').fill('make');
    await page.getByLabel('Model').fill('model');
    await page.getByTestId('lengthFeet').fill('22');
    await page.getByTestId('widthFeet').fill('22');
    await page.getByTestId('heightFeet').fill('22');
    await page.locator('label[for="hasTrailerNo"]').click();
    await page.getByTestId('remarks').fill('remarks test');
    await page.getByRole('button', { name: 'Continue' }).click();

    await expect(page.getByRole('heading', { name: 'Boat Haul-Away (BHA)' })).toBeVisible();
    await page.getByTestId('boatConfirmationContinue').click();

    await expect(page.getByText('Pickup info')).toBeVisible();
    await page.getByLabel('Preferred pickup date').fill('25 Dec 2022');
    await page.getByLabel('Preferred pickup date').blur();
    await page.getByText('Use my current address').click();
    await page.getByLabel('Preferred delivery date').fill('25 Dec 2022');
    await page.getByLabel('Preferred delivery date').blur();
    await page.getByRole('button', { name: 'Save & Continue' }).click();
    await customerPage.waitForPage.reviewShipments();
    await expect(page.getByRole('button', { name: 'Next' })).toBeDisabled();
  });
});

test.describe('(MultiMove) Boat shipment', () => {
  test.skip(multiMoveEnabled === 'false', 'Skip if MultiMove workflow is not enabled.');

  test('A customer can create a Boat shipment - Tow-Away', async ({ page, customerPage }) => {
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

    // Create an Boat shipment
    await page.getByText('Move a boat').click();
    await customerPage.navigateForward();

    // Fill in form to create Boat shipment
    await customerPage.waitForPage.boatShipment();
    await page.getByLabel('Year').fill('2022');
    await page.getByLabel('Make').fill('make');
    await page.getByLabel('Model').fill('model');
    await page.getByTestId('lengthFeet').fill('22');
    await page.getByTestId('widthFeet').fill('22');
    await page.getByTestId('heightFeet').fill('22');
    await page.locator('label[for="hasTrailerYes"]').click();
    await page.locator('label[for="isRoadworthyYes"]').click();
    await page.getByTestId('remarks').fill('remarks test');
    await page.getByRole('button', { name: 'Continue' }).click();

    await expect(page.getByRole('heading', { name: 'Boat Tow-Away (BTA)' })).toBeVisible();
    await page.getByTestId('boatConfirmationContinue').click();

    await expect(page.getByText('Pickup info')).toBeVisible();
    await page.getByLabel('Preferred pickup date').fill('25 Dec 2022');
    await page.getByLabel('Preferred pickup date').blur();
    await page.getByText('Use my current address').click();
    await page.getByLabel('Preferred delivery date').fill('25 Dec 2022');
    await page.getByLabel('Preferred delivery date').blur();
    await page.getByRole('button', { name: 'Save & Continue' }).click();
    await customerPage.waitForPage.reviewShipments();
  });

  test('A customer can create a Boat shipment - Haul-Away', async ({ page, customerPage }) => {
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

    // Create an Boat shipment
    await page.getByText('Move a boat').click();
    await customerPage.navigateForward();

    // Fill in form to create Boat shipment
    await customerPage.waitForPage.boatShipment();
    await page.getByLabel('Year').fill('2022');
    await page.getByLabel('Make').fill('make');
    await page.getByLabel('Model').fill('model');
    await page.getByTestId('lengthFeet').fill('22');
    await page.getByTestId('widthFeet').fill('22');
    await page.getByTestId('heightFeet').fill('22');
    await page.locator('label[for="hasTrailerNo"]').click();
    await page.getByTestId('remarks').fill('remarks test');
    await page.getByRole('button', { name: 'Continue' }).click();

    await expect(page.getByRole('heading', { name: 'Boat Haul-Away (BHA)' })).toBeVisible();
    await page.getByTestId('boatConfirmationContinue').click();

    await expect(page.getByText('Pickup info')).toBeVisible();
    await page.getByLabel('Preferred pickup date').fill('25 Dec 2022');
    await page.getByLabel('Preferred pickup date').blur();
    await page.getByText('Use my current address').click();
    await page.getByLabel('Preferred delivery date').fill('25 Dec 2022');
    await page.getByLabel('Preferred delivery date').blur();
    await page.getByRole('button', { name: 'Save & Continue' }).click();
    await customerPage.waitForPage.reviewShipments();
  });

  test('A customer is redirected to HHG if dimension requirement is not met', async ({ page, customerPage }) => {
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

    // Create an Boat shipment
    await page.getByText('Move a boat').click();
    await customerPage.navigateForward();

    // Fill in form to create Boat shipment
    await customerPage.waitForPage.boatShipment();
    await page.getByLabel('Year').fill('2022');
    await page.getByLabel('Make').fill('make');
    await page.getByLabel('Model').fill('model');
    await page.getByTestId('lengthFeet').fill('2');
    await page.getByTestId('widthFeet').fill('2');
    await page.getByTestId('heightFeet').fill('2');
    await page.locator('label[for="hasTrailerNo"]').click();
    await page.getByTestId('remarks').fill('remarks test');
    await page.getByRole('button', { name: 'Continue' }).click();

    await expect(
      page.getByRole('heading', { name: 'Movers pack and ship it, paid by the government (HHG)' }),
    ).toBeVisible();
    await page.getByTestId('boatConfirmationContinue').click();

    await expect(page.getByText('HHG')).toBeVisible();
  });

  test('Is able to delete a boat shipment', async ({ page, customerPage }) => {
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

    // Create an Boat shipment
    await page.getByText('Move a boat').click();
    await customerPage.navigateForward();

    // Fill in form to create Boat shipment
    await customerPage.waitForPage.boatShipment();
    await page.getByLabel('Year').fill('2022');
    await page.getByLabel('Make').fill('make');
    await page.getByLabel('Model').fill('model');
    await page.getByTestId('lengthFeet').fill('22');
    await page.getByTestId('widthFeet').fill('22');
    await page.getByTestId('heightFeet').fill('22');
    await page.locator('label[for="hasTrailerYes"]').click();
    await page.locator('label[for="isRoadworthyYes"]').click();
    await page.getByTestId('remarks').fill('remarks test');
    await page.getByRole('button', { name: 'Continue' }).click();

    await expect(page.getByRole('heading', { name: 'Boat Tow-Away (BTA)' })).toBeVisible();
    await page.getByTestId('boatConfirmationContinue').click();

    await expect(page.getByText('Pickup info')).toBeVisible();
    await page.getByLabel('Preferred pickup date').fill('25 Dec 2022');
    await page.getByLabel('Preferred pickup date').blur();
    await page.getByText('Use my current address').click();
    await page.getByLabel('Preferred delivery date').fill('25 Dec 2022');
    await page.getByLabel('Preferred delivery date').blur();
    await page.getByRole('button', { name: 'Save & Continue' }).click();
    await customerPage.waitForPage.reviewShipments();

    await expect(page.getByRole('heading', { name: 'Boat 1' })).toBeVisible();
    await page.getByTestId('deleteShipmentButton').click();
    await expect(page.getByRole('heading', { name: 'Delete this?' })).toBeVisible();
    await page.getByText('Yes, Delete').click();

    await expect(page.getByRole('heading', { name: 'Boat 1' })).not.toBeVisible();
  });

  test('Deletes existing boat shipment and is redirected to HHG after edit', async ({ page, customerPage }) => {
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

    // Create an Boat shipment
    await page.getByText('Move a boat').click();
    await customerPage.navigateForward();

    // Fill in form to create Boat shipment
    await customerPage.waitForPage.boatShipment();
    await page.getByLabel('Year').fill('2022');
    await page.getByLabel('Make').fill('make');
    await page.getByLabel('Model').fill('model');
    await page.getByTestId('lengthFeet').fill('22');
    await page.getByTestId('widthFeet').fill('22');
    await page.getByTestId('heightFeet').fill('22');
    await page.locator('label[for="hasTrailerYes"]').click();
    await page.locator('label[for="isRoadworthyYes"]').click();
    await page.getByTestId('remarks').fill('remarks test');
    await page.getByRole('button', { name: 'Continue' }).click();

    await expect(page.getByRole('heading', { name: 'Boat Tow-Away (BTA)' })).toBeVisible();
    await page.getByTestId('boatConfirmationContinue').click();

    await expect(page.getByText('Pickup info')).toBeVisible();
    await page.getByLabel('Preferred pickup date').fill('25 Dec 2022');
    await page.getByLabel('Preferred pickup date').blur();
    await page.getByText('Use my current address').click();
    await page.getByLabel('Preferred delivery date').fill('25 Dec 2022');
    await page.getByLabel('Preferred delivery date').blur();
    await page.getByRole('button', { name: 'Save & Continue' }).click();
    await customerPage.waitForPage.reviewShipments();

    await expect(page.getByRole('heading', { name: 'Boat 1' })).toBeVisible();
    await page.getByTestId('editShipmentButton').click();

    await customerPage.waitForPage.boatShipment();
    await page.getByTestId('lengthFeet').fill('2');
    await page.getByTestId('widthFeet').fill('2');
    await page.getByTestId('heightFeet').fill('2');
    await page.getByRole('button', { name: 'Continue' }).click();

    await expect(
      page.getByRole('heading', { name: 'Movers pack and ship it, paid by the government (HHG)' }),
    ).toBeVisible();
    await page.getByTestId('boatConfirmationContinue').click();
    await expect(
      page.getByRole('heading', { name: 'Movers pack and ship it, paid by the government (HHG)' }),
    ).not.toBeVisible();
    await expect(page.getByText('HHG')).toBeVisible();
    await expect(page.getByText('Movers pack and transport this shipment')).toBeVisible();
    await page.getByTestId('wizardNextButton').click();
    await customerPage.waitForPage.reviewShipments();
    await expect(page.getByRole('heading', { name: 'Boat 1' })).not.toBeVisible();
    await expect(page.getByRole('heading', { name: 'HHG 1' })).toBeVisible();
  });

  test('A customer is unable to sign if boat shipment is incomplete', async ({ page, customerPage }) => {
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

    // Create an Boat shipment
    await page.getByText('Move a boat').click();
    await customerPage.navigateForward();

    // Fill in form to create Boat shipment
    await customerPage.waitForPage.boatShipment();
    await page.getByLabel('Year').fill('2022');
    await page.getByLabel('Make').fill('make');
    await page.getByLabel('Model').fill('model');
    await page.getByTestId('lengthFeet').fill('22');
    await page.getByTestId('widthFeet').fill('22');
    await page.getByTestId('heightFeet').fill('22');
    await page.locator('label[for="hasTrailerYes"]').click();
    await page.locator('label[for="isRoadworthyYes"]').click();
    await page.getByTestId('remarks').fill('remarks test');
    await page.getByRole('button', { name: 'Continue' }).click();

    await expect(page.getByRole('heading', { name: 'Boat Tow-Away (BTA)' })).toBeVisible();
    await page.getByTestId('boatConfirmationContinue').click();

    // Back out to mark as incomplete
    await expect(page.getByText('Pickup info')).toBeVisible();
    await page.getByRole('button', { name: 'Back' }).click();
    await customerPage.waitForPage.boatShipment();
    await page.getByRole('button', { name: 'Back' }).click();

    await expect(page.getByText('Time to submit your move')).toBeVisible();
    await expect(page.getByTestId('review-and-submit-btn')).toBeDisabled();

    // Create another shipment to test continue button is disabled from the review page
    await page.getByTestId('shipment-selection-btn').click();
    await customerPage.waitForPage.selectShipmentType();

    // Create an Boat shipment
    await page.getByText('Move a boat').click();
    await customerPage.navigateForward();

    // Fill in form to create Boat shipment
    await customerPage.waitForPage.boatShipment();
    await page.getByLabel('Year').fill('2022');
    await page.getByLabel('Make').fill('make');
    await page.getByLabel('Model').fill('model');
    await page.getByTestId('lengthFeet').fill('22');
    await page.getByTestId('widthFeet').fill('22');
    await page.getByTestId('heightFeet').fill('22');
    await page.locator('label[for="hasTrailerNo"]').click();
    await page.getByTestId('remarks').fill('remarks test');
    await page.getByRole('button', { name: 'Continue' }).click();

    await expect(page.getByRole('heading', { name: 'Boat Haul-Away (BHA)' })).toBeVisible();
    await page.getByTestId('boatConfirmationContinue').click();

    await expect(page.getByText('Pickup info')).toBeVisible();
    await page.getByLabel('Preferred pickup date').fill('25 Dec 2022');
    await page.getByLabel('Preferred pickup date').blur();
    await page.getByText('Use my current address').click();
    await page.getByLabel('Preferred delivery date').fill('25 Dec 2022');
    await page.getByLabel('Preferred delivery date').blur();
    await page.getByRole('button', { name: 'Save & Continue' }).click();
    await customerPage.waitForPage.reviewShipments();
    await expect(page.getByRole('button', { name: 'Next' })).toBeDisabled();
  });
});
