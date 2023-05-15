// @ts-check
import { test, expect } from '../../utils/office/officeTest';

import { TooFlowPage } from './tooTestFixture';

test.describe('TOO user', () => {
  /** @type {TooFlowPage} */
  let tooFlowPage;

  test.describe('updating a move shipment in SIT', () => {
    test.beforeEach(async ({ officePage }) => {
      // build move in SIT with 200 days authorized and without pending extension requests
      const move = await officePage.testHarness.buildHHGMoveIn200DaysSIT();
      await officePage.signInAsNewTOOUser();
      tooFlowPage = new TooFlowPage(officePage, move);
      await officePage.tooNavigateToMove(move.locator);
    });

    test('is able to increase a SIT authorization', async ({ page }) => {
      // navigate to MTO tab
      await page.getByTestId('MoveTaskOrder-Tab').click();
      await tooFlowPage.waitForPage.moveTaskOrder();

      // increase SIT authorization to 220 days
      await page.getByTestId('sitExtensions').getByTestId('button').click();
      await expect(page.getByRole('heading', { name: 'Edit SIT authorization' })).toBeVisible();
      await page.getByTestId('daysApproved').clear();
      await page.getByTestId('daysApproved').fill('220');
      await page.getByTestId('dropdown').selectOption('AWAITING_COMPLETION_OF_RESIDENCE');
      await page.getByTestId('officeRemarks').fill('residence under construction');
      await expect(page.getByTestId('form').getByTestId('button')).toBeEnabled();
      await page.getByTestId('form').getByTestId('button').click();

      // assert that days authorization is now 220
      await expect(page.getByTestId('sitStatusTable').getByText('220', { exact: true }).first()).toBeVisible();
    });

    test('is able to decrease a SIT authorization', async ({ page }) => {
      // navigate to MTO tab
      await page.getByTestId('MoveTaskOrder-Tab').click();
      await tooFlowPage.waitForPage.moveTaskOrder();

      // decrease SIT authorization to 190 days
      await page.getByTestId('sitExtensions').getByTestId('button').click();
      await expect(page.getByRole('heading', { name: 'Edit SIT authorization' })).toBeVisible();
      await page.getByTestId('daysApproved').clear();
      await page.getByTestId('daysApproved').fill('190');
      await page.getByTestId('dropdown').selectOption('AWAITING_COMPLETION_OF_RESIDENCE');
      await page.getByTestId('officeRemarks').fill('residence under construction');
      await page.getByTestId('form').getByTestId('button').click();

      // assert that days authorization is now 190
      await expect(page.getByTestId('sitStatusTable').getByText('190', { exact: true }).first()).toBeVisible();
    });

    test('is unable to decrease the SIT authorization below the number of days already used', async ({ page }) => {
      // navigate to MTO tab
      await page.getByTestId('MoveTaskOrder-Tab').click();
      await tooFlowPage.waitForPage.moveTaskOrder();

      // try to decrease SIT authorization to 1 day
      await page.getByTestId('sitExtensions').getByTestId('button').click();
      await expect(page.getByRole('heading', { name: 'Edit SIT authorization' })).toBeVisible();
      await page.getByTestId('daysApproved').clear();
      await page.getByTestId('daysApproved').fill('1');
      await page.getByTestId('dropdown').selectOption('AWAITING_COMPLETION_OF_RESIDENCE');
      await page.getByTestId('officeRemarks').fill('residence under construction');

      // assert that save button is disabled and error messages are present
      await expect(page.getByTestId('form').getByTestId('button')).not.toBeEnabled();
      await expect(page.getByTestId('form').getByTestId('sitStatusTable').getByTestId('errorMessage')).toBeVisible();
      await expect(
        page.getByText('The end date must occur after the start date. Please select a new date.', { exact: true }),
      ).toBeVisible();
    });
  });

  test.describe('updating a move shipment in SIT with a SIT extension request', () => {
    test.beforeEach(async ({ officePage }) => {
      // build move in SIT with 200 days authorized and with one pending extension request
      const move = await officePage.testHarness.buildHHGMoveIn200DaysSITWithPendingExtension();
      await officePage.signInAsNewTOOUser();
      tooFlowPage = new TooFlowPage(officePage, move);
      await officePage.tooNavigateToMove(move.locator);
    });

    test('is able to approve the SIT extension request', async ({ page }) => {
      // navigate to MTO tab
      await page.getByTestId('MoveTaskOrder-Tab').click();
      await tooFlowPage.waitForPage.moveTaskOrder();

      // assert that there is a pending SIT extension request
      await expect(page.getByText('Additional days requested')).toBeVisible();

      // approve SIT extension with an adjusted approved days value of 220 days and change the extension reason
      await page.getByTestId('sitExtensions').getByTestId('button').click();
      await expect(page.getByRole('heading', { name: 'Review additional days requested' })).toBeVisible();
      await page.getByTestId('daysApproved').clear();
      await page.getByTestId('daysApproved').fill('220');
      await page.getByText('Yes', { exact: true }).click();
      await page.getByTestId('dropdown').selectOption('OTHER');
      await page.getByTestId('officeRemarks').fill('allowance increased by 20 days instead of the requested 45 days');
      await page.getByTestId('form').getByTestId('button').click();

      // assert that there is no pending SIT extension request and the days authorization is now 220
      await expect(page.getByText('Additional days requested')).toBeHidden();
      await expect(page.getByTestId('sitStatusTable').getByText('220', { exact: true }).first()).toBeVisible();
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
      await page.getByTestId('form').getByTestId('button').click();

      // assert that there is no pending SIT extension request and the days authorization is still 200
      await expect(page.getByText('Additional days requested')).toBeHidden();
      await expect(page.getByTestId('sitStatusTable').getByText('200', { exact: true }).first()).toBeVisible();
    });
  });
});
