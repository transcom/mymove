// @ts-check
import { test, expect } from '../../utils/my/customerTest';

const multiMoveEnabled = process.env.FEATURE_FLAG_MULTI_MOVE;

test.describe('Orders', () => {
  test.skip(multiMoveEnabled === 'true', 'Skip if MultiMove workflow is enabled.');
  test('Users can upload orders, and delete if the move is in draft status', async ({ page, customerPage }) => {
    // Generate a new onboarded user and log in
    const user = await customerPage.testHarness.buildNeedsOrdersUser();
    const userId = user.id;
    await customerPage.signInAsExistingCustomer(userId);

    // Navigate to add orders
    await customerPage.waitForPage.home();
    await page.getByRole('button', { name: 'Add orders' }).click();
    await customerPage.waitForPage.ordersDetails();

    // Fill in orders details
    await page.getByLabel('Orders type').selectOption({ label: 'Permanent Change Of Station (PCS)' });
    await page.getByLabel('Orders date').fill('6/2/2018');
    await page.getByLabel('Orders date').blur();
    await page.getByLabel('Report by date').fill('8/9/2018');
    await page.getByLabel('Report by date').blur();

    // UGH
    // because of the styling of this input item, we cannot use a
    // css locator for the input item and then click it
    //
    // The styling is very similar to the issue described in
    //
    // https://github.com/microsoft/playwright/issues/3688
    //
    await page.locator('div:has(label:has-text("Are dependents")) >> div.usa-radio').getByText('No').click();

    await customerPage.selectDutyLocation('Marine Corps AS Yuma, AZ 85369', 'new_duty_location');
    await page.keyboard.press('Backspace'); // tests if backspace clears the duty location field
    await expect(page.getByLabel('New duty location')).toBeEmpty();
    await customerPage.selectDutyLocation('Marine Corps AS Yuma, AZ 85369', 'new_duty_location');

    await customerPage.selectDutyLocation('Marine Corps AS Yuma, AZ 85369', 'origin_duty_location');
    await page.keyboard.press('Backspace'); // tests if backspace clears the duty location field
    await expect(page.getByLabel('Current duty location')).toBeEmpty();
    await customerPage.selectDutyLocation('Marine Corps AS Yuma, AZ 85369', 'origin_duty_location');

    await page.getByRole('combobox', { name: 'Pay grade' }).selectOption({ label: 'E-7' });
    await page
      .getByRole('combobox', { name: 'Counseling Office' })
      .selectOption({ label: 'PPPO DMO Camp Pendleton - USMC' });

    await customerPage.navigateForward();
    await customerPage.waitForPage.ordersUpload();

    // Upload an orders document, then submit
    // Annoyingly, there's no test IDs or labeling text for this control, so the only way to access it is .locator
    const filepondContainer = page.locator('.filepond--wrapper');
    await customerPage.uploadFileViaFilepond(filepondContainer, 'AF Orders Sample.pdf');
    await customerPage.navigateForward();

    // Verify that we're on the home page and that orders have been uploaded
    await customerPage.waitForPage.home();
    await expect(page.getByText('Orders uploaded')).toBeVisible();

    // Delete orders in draft status
    await page.getByTestId('stepContainer2').getByRole('button', { name: 'Edit' }).click();
    await customerPage.waitForPage.editOrders();
    await expect(page.getByText('AF Orders Sample.pdf')).toBeVisible();
    await page.getByRole('button', { name: 'Delete' }).click();
    await expect(page.getByText('AF Orders Sample.pdf')).not.toBeVisible();
  });
});

test.describe('(MultiMove) Orders', () => {
  test.skip(multiMoveEnabled === 'false', 'Skip if MultiMove workflow is not enabled.');

  test('Users can upload orders, and delete if the move is in draft status', async ({ page, customerPage }) => {
    // Generate a new onboarded user and log in
    const user = await customerPage.testHarness.buildNeedsOrdersUser();
    const userId = user.id;
    await customerPage.signInAsExistingCustomer(userId);
    await customerPage.createMoveButtonClick();

    // Fill in orders details
    await page.getByLabel('Orders type').selectOption({ label: 'Permanent Change Of Station (PCS)' });
    await page.getByLabel('Orders date').fill('6/2/2018');
    await page.getByLabel('Orders date').blur();
    await page.getByLabel('Report by date').fill('8/9/2018');
    await page.getByLabel('Report by date').blur();

    // UGH
    // because of the styling of this input item, we cannot use a
    // css locator for the input item and then click it
    //
    // The styling is very similar to the issue described in
    //
    // https://github.com/microsoft/playwright/issues/3688
    //
    await page.locator('div:has(label:has-text("Are dependents")) >> div.usa-radio').getByText('No').click();

    await customerPage.selectDutyLocation('Marine Corps AS Yuma, AZ 85369', 'new_duty_location');
    await page.keyboard.press('Backspace'); // tests if backspace clears the duty location field
    await expect(page.getByLabel('New duty location')).toBeEmpty();
    await customerPage.selectDutyLocation('Marine Corps AS Yuma, AZ 85369', 'new_duty_location');

    await customerPage.selectDutyLocation('Marine Corps AS Yuma, AZ 85369', 'origin_duty_location');
    await page.keyboard.press('Backspace'); // tests if backspace clears the duty location field
    await expect(page.getByLabel('Current duty location')).toBeEmpty();
    await customerPage.selectDutyLocation('Marine Corps AS Yuma, AZ 85369', 'origin_duty_location');

    await page.getByRole('combobox', { name: 'Pay grade' }).selectOption({ label: 'E-7' });
    await page
      .getByRole('combobox', { name: 'Counseling Office' })
      .selectOption({ label: 'PPPO DMO Camp Pendleton - USMC' });

    await customerPage.navigateForward();
    await customerPage.waitForPage.ordersUpload();

    // Upload an orders document, then submit
    // Annoyingly, there's no test IDs or labeling text for this control, so the only way to access it is .locator
    const filepondContainer = page.locator('.filepond--wrapper');
    await customerPage.uploadFileViaFilepond(filepondContainer, 'AF Orders Sample.pdf');
    await customerPage.navigateForward();

    // Verify that we're on the home page and that orders have been uploaded
    await customerPage.waitForPage.home();
    await expect(page.getByText('Orders uploaded')).toBeVisible();

    // Delete orders in draft status
    await page.getByTestId('stepContainer2').getByRole('button', { name: 'Edit' }).click();
    await customerPage.waitForPage.editOrders();
    await expect(page.getByText('AF Orders Sample.pdf')).toBeVisible();
    await page.getByRole('button', { name: 'Delete' }).click();
    await expect(page.getByText('AF Orders Sample.pdf')).not.toBeVisible();
  });
});

test.describe('Download Orders', () => {
  test('Users can download their orders for viewing', async ({ page, customerPage }) => {
    // Generate a move that has the status of SUBMITTED
    const move = await customerPage.testHarness.buildSubmittedMoveWithPPMShipmentForSC();
    const userId = move.Orders.ServiceMember.user_id;

    // Sign-in and navigate to move home page
    await customerPage.signInAsExistingCustomer(userId);
    if (multiMoveEnabled === 'true') {
      await customerPage.navigateFromMMDashboardToMove(move);
    }
    await customerPage.waitForPage.home();

    // Go to the Edit Orders page
    await page.getByTestId('review-and-submit-btn').click();
    await page.getByTestId('edit-orders-table').click();

    // Upload second set of orders
    const filepondContainer = page.locator('.filepond--wrapper');
    await customerPage.uploadFileViaFilepond(filepondContainer, 'secondOrders.pdf');

    // Verify filename is a downloadable link
    await expect(page.getByRole('link', { name: 'secondOrders.pdf' })).toBeVisible();
  });
});

test.describe('Download Amended Orders', () => {
  test('Users can download their amended orders for viewing', async ({ page, customerPage }) => {
    // Generate a move that has the status of SUBMITTED
    const move = await customerPage.testHarness.buildSubmittedMoveWithPPMShipmentForSC();
    const userId = move.Orders.ServiceMember.user_id;

    // Sign-in and navigate to move home page
    await customerPage.signInAsExistingCustomer(userId);
    if (multiMoveEnabled === 'true') {
      await customerPage.navigateFromMMDashboardToMove(move);
    }
    await customerPage.waitForPage.home();

    // Go to the Upload Amended Documents page
    await page.getByRole('button', { name: 'Upload documents' }).click();

    // Upload amended orders
    const filepondContainer = page.locator('.filepond--wrapper');
    await customerPage.uploadFileViaFilepond(filepondContainer, 'amendedOrders.pdf');

    // Verify filename is a downloadable link
    await expect(page.getByRole('link', { name: 'amendedOrders.pdf' })).toBeVisible();
  });
});
