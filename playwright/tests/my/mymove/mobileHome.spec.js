import { test, expect } from '../../utils/my/customerTest';

const multiMoveEnabled = process.env.FEATURE_FLAG_MULTI_MOVE;

test.describe('Mobile Home shipment', () => {
  test.skip(multiMoveEnabled === 'true', 'Skip if MultiMove workflow is enabled.');

  test('A customer can create a Mobile Home shipment', async ({ page, customerPage }) => {
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

    // Create an Mobile Home shipment
    await page.getByText('Move a Mobile Home').click();
    await customerPage.navigateForward();

    // Fill in form to create Mobile Home shipment
    await customerPage.waitForPage.mobileHomeShipment();
    await page.getByLabel('Year').fill('2022');
    await page.getByLabel('Make').fill('make');
    await page.getByLabel('Model').fill('model');
    await page.getByTestId('lengthFeet').fill('22');
    await page.getByTestId('widthFeet').fill('22');
    await page.getByTestId('heightFeet').fill('22');
    await page.getByRole('button', { name: 'Continue' }).click();

    await expect(page.getByTestId('tag')).toContainText('Mobile Home');

    await expect(page.getByText('Pickup info')).toBeVisible();
    await page.getByLabel('Preferred pickup date').fill('25 Dec 2022');
    await page.getByLabel('Preferred pickup date').blur();
    await page.getByText('Use my current address').click();
    await page.getByLabel('Preferred delivery date').fill('25 Dec 2022');
    await page.getByLabel('Preferred delivery date').blur();
    await page.getByRole('button', { name: 'Save & Continue' }).click();
    await customerPage.waitForPage.reviewShipments();
  });
});

test.describe('(MultiMove) Mobile Home shipment', () => {
  test.skip(multiMoveEnabled === 'false', 'Skip if MultiMove workflow is not enabled.');

  test('A customer can create a Mobile Home shipment', async ({ page, customerPage }) => {
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

    // Create an Mobile Home shipment
    await page.getByText('Move a mobile home').click();
    await customerPage.navigateForward();

    // Fill in form to create Mobile Home shipment
    await customerPage.waitForPage.mobileHomeShipment();
    await page.getByLabel('Year').fill('2022');
    await page.getByLabel('Make').fill('make');
    await page.getByLabel('Model').fill('model');
    await page.getByTestId('lengthFeet').fill('22');
    await page.getByTestId('widthFeet').fill('22');
    await page.getByTestId('heightFeet').fill('22');
    await page.getByRole('button', { name: 'Continue' }).click();

    await expect(page.getByTestId('tag')).toContainText('Mobile Home');

    await expect(page.getByText('Pickup info')).toBeVisible();
    await page.getByLabel('Preferred pickup date').fill('25 Dec 2022');
    await page.getByLabel('Preferred pickup date').blur();
    await page.getByText('Use my current address').click();
    await page.getByLabel('Preferred delivery date').fill('25 Dec 2022');
    await page.getByLabel('Preferred delivery date').blur();
    await page.getByRole('button', { name: 'Save & Continue' }).click();
    await customerPage.waitForPage.reviewShipments();
  });

  test('Is able to delete a Mobile Home shipment', async ({ page, customerPage }) => {
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

    // Create an Mobile Home shipment
    await page.getByText('Move a mobile home').click();
    await customerPage.navigateForward();

    // Fill in form to create Mobile Home shipment
    await customerPage.waitForPage.mobileHomeShipment();
    await page.getByLabel('Year').fill('2022');
    await page.getByLabel('Make').fill('make');
    await page.getByLabel('Model').fill('model');
    await page.getByTestId('lengthFeet').fill('22');
    await page.getByTestId('widthFeet').fill('22');
    await page.getByTestId('heightFeet').fill('22');
    await page.getByRole('button', { name: 'Continue' }).click();

    await expect(page.getByTestId('tag')).toContainText('Mobile Home');

    await expect(page.getByText('Pickup info')).toBeVisible();
    await page.getByLabel('Preferred pickup date').fill('25 Dec 2022');
    await page.getByLabel('Preferred pickup date').blur();
    await page.getByText('Use my current address').click();
    await page.getByLabel('Preferred delivery date').fill('25 Dec 2022');
    await page.getByLabel('Preferred delivery date').blur();
    await page.getByRole('button', { name: 'Save & Continue' }).click();
    await customerPage.waitForPage.reviewShipments();

    await expect(page.getByRole('heading', { name: 'Mobile Home 1' })).toBeVisible();
    await page.getByTestId('deleteShipmentButton').click();
    await expect(page.getByRole('heading', { name: 'Delete this?' })).toBeVisible();
    await page.getByText('Yes, Delete').click();

    await expect(page.getByRole('heading', { name: 'Mobile Home 1' })).not.toBeVisible();
  });
});
