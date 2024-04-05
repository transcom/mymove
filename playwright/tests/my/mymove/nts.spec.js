// @ts-check
import { test, expect } from '../../utils/my/customerTest';

const multiMoveEnabled = process.env.FEATURE_FLAG_MULTI_MOVE;

test.describe('NTS', () => {
  test.skip(multiMoveEnabled === 'true', 'Skip if MultiMove workflow is enabled.');
  test('A customer can create, edit, and delete an NTS shipment', async ({ page, customerPage }) => {
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

    // Create an NTS shipment
    await page.getByText('It is going into storage for months or years (NTS)').click();
    await customerPage.navigateForward();

    // Fill in form to create NTS shipment
    await customerPage.waitForPage.ntsShipment();
    await page.getByLabel('Preferred pickup date').fill('25 Dec 2022');
    await page.getByLabel('Preferred pickup date').blur();
    await page.getByText('Use my current address').click();
    await page.getByTestId('remarks').fill('Grandfather antique clock');
    await customerPage.navigateForward();

    // Verify that form submitted
    await customerPage.waitForPage.reviewShipments();
    await expect(page.getByText('Grandfather antique clock')).toBeVisible();

    // Navigate to edit shipment from the review page
    await page.getByTestId('edit-nts-shipment-btn').click();
    await customerPage.waitForPage.ntsShipment();

    // Update form (adding releasing agent)
    await page.getByLabel('First name').fill('Grace');
    await page.getByLabel('Last name').fill('Griffin');
    await page.getByLabel('Phone').fill('2025551234');
    await page.getByLabel('Email').fill('grace.griffin@example.com');
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
});

test.describe('(MultiMove) NTS', () => {
  test.skip(multiMoveEnabled === 'false', 'Skip if MultiMove workflow is not enabled.');

  test('A customer can create, edit, and delete an NTS shipment', async ({ page, customerPage }) => {
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

    // Create an NTS shipment
    await page.getByText('It is going into storage for months or years (NTS)').click();
    await customerPage.navigateForward();

    // Fill in form to create NTS shipment
    await customerPage.waitForPage.ntsShipment();
    await page.getByLabel('Preferred pickup date').fill('25 Dec 2022');
    await page.getByLabel('Preferred pickup date').blur();
    await page.getByText('Use my current address').click();
    await page.getByTestId('remarks').fill('Grandfather antique clock');
    await customerPage.navigateForward();

    // Verify that form submitted
    await customerPage.waitForPage.reviewShipments();
    await expect(page.getByText('Grandfather antique clock')).toBeVisible();

    // Navigate to edit shipment from the review page
    await page.getByTestId('edit-nts-shipment-btn').click();
    await customerPage.waitForPage.ntsShipment();

    // Update form (adding releasing agent)
    await page.getByLabel('First name').fill('Grace');
    await page.getByLabel('Last name').fill('Griffin');
    await page.getByLabel('Phone').fill('2025551234');
    await page.getByLabel('Email').fill('grace.griffin@example.com');
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
});
