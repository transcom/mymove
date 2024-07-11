// @ts-check

import { test, expect } from '../../utils/office/officeTest';

import { TooFlowPage } from './tooTestFixture';

test.describe('TOO user', () => {
  /** @type {TooFlowPage} */
  let tooFlowPage;

  test.describe('updating a move shipment in SIT', () => {
    test.beforeEach(async ({ officePage }) => {
      // build move in SIT with 90 days authorized and without pending extension requests
      const move = await officePage.testHarness.buildHHGMoveInSIT();
      await officePage.signInAsNewTOOUser();
      tooFlowPage = new TooFlowPage(officePage, move);
      await tooFlowPage.waitForLoading();
      await officePage.tooNavigateToMove(move.locator);
    });

    test('is able to see the SIT Departure Date', async ({ page }) => {
      // navigate to MTO tab
      await page.getByTestId('MoveTaskOrder-Tab').click();
      await tooFlowPage.waitForPage.moveTaskOrder();

      const target = await page
        .getByTestId('sitDaysAtCurrentLocation')
        .locator('table[class="DataTable_dataTable__TGt9M table--data-point"]')
        .locator('tbody')
        .locator('td')
        .nth(1)
        .locator('div')
        .locator('span')
        .textContent();
      const pattern = /(—|\d{2} \w{3} \d{4})/;

      expect(pattern.test(target)).toBeTruthy();
    });

    test('is able to increase a SIT authorization', async ({ page }) => {
      // navigate to MTO tab
      await page.getByTestId('MoveTaskOrder-Tab').click();
      await tooFlowPage.waitForPage.moveTaskOrder();

      // increase SIT authorization to 100 days
      await page.getByTestId('sitExtensions').getByRole('button', { name: 'Edit' }).click();
      await expect(page.getByRole('heading', { name: 'Edit SIT authorization' })).toBeVisible();
      await page.getByTestId('daysApproved').clear();
      await page.getByTestId('daysApproved').fill('100');
      await page.getByTestId('reasonDropdown').selectOption('AWAITING_COMPLETION_OF_RESIDENCE');
      await page.getByTestId('officeRemarks').fill('residence under construction');
      await expect(page.getByTestId('form').getByTestId('button')).toBeEnabled();
      await page.getByTestId('form').getByTestId('button').click();

      // assert that days authorization is now 100
      await expect(page.getByTestId('sitStatusTable').getByText('100', { exact: true }).first()).toBeVisible();
    });

    test('is able to decrease a SIT authorization', async ({ page }) => {
      // navigate to MTO tab
      await page.getByTestId('MoveTaskOrder-Tab').click();
      await tooFlowPage.waitForPage.moveTaskOrder();

      // decrease SIT authorization to 80 days
      await page.getByTestId('sitExtensions').getByRole('button', { name: 'Edit' }).click();
      await expect(page.getByRole('heading', { name: 'Edit SIT authorization' })).toBeVisible();
      await page.getByTestId('daysApproved').clear();
      await page.getByTestId('daysApproved').fill('80');
      await page.getByTestId('reasonDropdown').selectOption('AWAITING_COMPLETION_OF_RESIDENCE');
      await page.getByTestId('officeRemarks').fill('residence under construction');
      await page.getByTestId('form').getByTestId('button').click();

      // assert that days authorization is now 80
      await expect(page.getByTestId('sitStatusTable').getByText('80', { exact: true }).first()).toBeVisible();
    });

    test('is able to see appropriate results for 90 days approved and 90 days used', async ({ page }) => {
      await page.getByTestId('MoveTaskOrder-Tab').click();
      await tooFlowPage.waitForPage.moveTaskOrder();
      const daysApprovedCapture = await page
        .getByTestId('sitStatusTable')
        .locator('[class="DataTable_dataTable__TGt9M table--data-point"]')
        .locator('tbody')
        .locator('tr')
        .locator('span')
        .nth(0)
        .textContent();
      const daysUsedCapture = await page
        .getByTestId('sitStatusTable')
        .locator('[class="DataTable_dataTable__TGt9M table--data-point"]')
        .locator('tbody')
        .locator('tr')
        .locator('span')
        .nth(1)
        .textContent();
      const daysLeftCapture = await page
        .getByTestId('sitStatusTable')
        .locator('[class="DataTable_dataTable__TGt9M table--data-point"]')
        .locator('tbody')
        .locator('tr')
        .locator('span')
        .nth(2)
        .textContent();

      expect(daysApprovedCapture).toEqual('90');
      expect(daysUsedCapture).toEqual('90');
      expect(daysLeftCapture).toEqual('Expired');
    });

    test('is unable to decrease the SIT authorization below the number of days already used', async ({ page }) => {
      // navigate to MTO tab
      await page.getByTestId('MoveTaskOrder-Tab').click();
      await tooFlowPage.waitForPage.moveTaskOrder();

      // try to decrease SIT authorization to 1 day
      await page.getByTestId('sitExtensions').getByRole('button', { name: 'Edit' }).click();
      await expect(page.getByRole('heading', { name: 'Edit SIT authorization' })).toBeVisible();
      await page.getByTestId('daysApproved').clear();
      await page.getByTestId('daysApproved').fill('1');
      await page.getByTestId('reasonDropdown').selectOption('AWAITING_COMPLETION_OF_RESIDENCE');
      await page.getByTestId('officeRemarks').fill('residence under construction');

      // assert that save button is disabled and error messages are present
      await expect(page.getByTestId('form').getByTestId('button')).not.toBeEnabled();
      await expect(page.getByTestId('form').getByTestId('sitStatusTable').getByTestId('errorMessage')).toBeVisible();
      await expect(
        page.getByText('The end date must occur after the start date. Please select a new date.', { exact: true }),
      ).toBeVisible();
    });
  });

  test.describe('converting a SIT to customer expense using convert to customer expense button', () => {
    test.beforeEach(async ({ officePage }) => {
      // build move in SIT that ends today
      const move = await officePage.testHarness.buildHHGMoveInSITEndsToday();
      await officePage.signInAsNewTOOUser();
      tooFlowPage = new TooFlowPage(officePage, move);
      await tooFlowPage.waitForLoading();
      await officePage.tooNavigateToMove(move.locator);
    });

    test('is able to convert a SIT to customer expense', async ({ page }) => {
      // navigate to MTO tab
      await page.getByTestId('MoveTaskOrder-Tab').click();
      await tooFlowPage.waitForPage.moveTaskOrder();

      // assert that there is a Convert to customer expense button
      await expect(page.getByText('Convert to customer expense')).toBeVisible();

      // Convert SIT to customer expense through the modal
      await page.getByRole('button', { name: 'Convert to customer expense' }).click();
      await expect(page.getByRole('heading', { name: 'Convert SIT To Customer Expense' })).toBeVisible();
      await page.getByTestId('remarks').fill('testing convert to customer expense');
      await page.getByTestId('form').getByTestId('button').click();

      // assert that there is a Converted To Customer Expense Tag
      await expect(page.getByTestId('tag').getByText('Converted to customer expense')).toBeVisible();

      // assert that there is no Convert to customer expense button showing
      await expect(page.getByText('Convert to customer expense')).toBeHidden();
    });
  });
  test.describe('updating a move shipment in SIT with a SIT extension request', () => {
    test.beforeEach(async ({ officePage }) => {
      // build move in SIT with 90 days authorized and with one pending extension request
      const move = await officePage.testHarness.buildHHGMoveInSITWithPendingExtension();
      await officePage.signInAsNewTOOUser();
      tooFlowPage = new TooFlowPage(officePage, move);
      await tooFlowPage.waitForLoading();
      await officePage.tooNavigateToMove(move.locator);
    });

    test('is able to approve the SIT extension request', async ({ page }) => {
      // navigate to MTO tab
      await page.getByTestId('MoveTaskOrder-Tab').click();
      await tooFlowPage.waitForPage.moveTaskOrder();

      // assert that there is a pending SIT extension request
      await expect(page.getByText('Additional days requested')).toBeVisible();

      // approve SIT extension with an adjusted approved days value of 100 days and change the extension reason
      await page.getByTestId('sitExtensions').getByTestId('button').click();
      await expect(page.getByRole('heading', { name: 'Review additional days requested' })).toBeVisible();
      await page.getByTestId('daysApproved').clear();
      await page.getByTestId('daysApproved').fill('100');
      await page.getByText('Yes', { exact: true }).click();
      await page.getByTestId('reasonDropdown').selectOption('OTHER');
      await page.getByTestId('officeRemarks').fill('allowance increased by 20 days instead of the requested 45 days');
      await page.getByTestId('form').getByTestId('button').click();

      // assert that there is no pending SIT extension request and the days authorization is now 100
      await expect(page.getByText('Additional days requested')).toBeHidden();
      await expect(page.getByTestId('sitStatusTable').getByText('100', { exact: true }).first()).toBeVisible();
    });

    test('is able to deny the SIT extension request', async ({ page }) => {
      // navigate to MTO tab
      await page.getByTestId('MoveTaskOrder-Tab').click();
      await tooFlowPage.waitForPage.moveTaskOrder();

      // assert that there is a pending SIT extension request
      await expect(page.getByText('Additional days requested')).toBeVisible();

      // deny SIT extension
      await page.getByTestId('sitExtensions').getByTestId('button').click();
      await expect(page.getByRole('heading', { name: 'Review additional days requested' })).toBeVisible();
      await page.getByText('No', { exact: true }).click();
      await page.getByTestId('officeRemarks').fill('extension request denied');
      await page.getByTestId('convertToCustomerExpense');
      await page.getByTestId('form').getByTestId('button').click();

      // assert that there is no pending SIT extension request and the days authorization is still 90
      await expect(page.getByText('Additional days requested')).toBeHidden();
      await expect(page.getByTestId('sitStatusTable').getByText('90', { exact: true }).first()).toBeVisible();
    });

    test('is able to deny the SIT extension request AND convert to customer expense', async ({ page }) => {
      // navigate to MTO tab
      await page.getByTestId('MoveTaskOrder-Tab').click();
      await tooFlowPage.waitForPage.moveTaskOrder();

      // assert that there is a pending SIT extension request
      await expect(page.getByText('Additional days requested')).toBeVisible();

      // deny SIT extension
      await page.getByTestId('sitExtensions').getByTestId('button').click();
      await expect(page.getByRole('heading', { name: 'Review additional days requested' })).toBeVisible();
      await page.getByText('No', { exact: true }).click();
      await page.getByTestId('officeRemarks').fill('extension request denied');
      await page.getByTestId('convertToCustomerExpense').click();
      await page.getByTestId('convertToCustomerExpenseConfirmationYes').click();
      await page.getByTestId('form').getByTestId('button').click();

      // assert that there is a Converted To Customer Expense Tag
      await expect(page.getByTestId('tag').getByText('Converted to customer expense')).toBeVisible();

      // assert that there is no pending SIT extension request and the days authorization is still 90
      await expect(page.getByText('Additional days requested')).toBeHidden();
      await expect(page.getByTestId('sitStatusTable').getByText('90', { exact: true }).first()).toBeVisible();
    });
    test('is showing correct labels', async ({ page }) => {
      // navigate to MTO tab
      await page.getByTestId('MoveTaskOrder-Tab').click();
      await tooFlowPage.waitForPage.moveTaskOrder();

      await expect(page.getByText('Total days of SIT approved')).toBeVisible();
      await expect(page.getByText('Total days used')).toBeVisible();
      await expect(page.getByText('Total days remaining')).toBeVisible();
      await expect(page.getByText('SIT start date').nth(0)).toBeVisible();
      await expect(page.getByText('SIT authorized end date')).toBeVisible();
      await expect(page.getByText('Calculated total SIT days')).toBeVisible();
    });
    test('is showing the SIT Departure Date section', async ({ page }) => {
      // navigate to MTO tab
      await page.getByTestId('MoveTaskOrder-Tab').click();
      await tooFlowPage.waitForPage.moveTaskOrder();

      await expect(
        page.locator('table[class="DataTable_dataTable__TGt9M table--data-point"]').getByText('SIT departure date'),
      ).toBeVisible();
    });
  });
});
