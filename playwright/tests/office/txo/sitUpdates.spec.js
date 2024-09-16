// @ts-check
import { test, expect } from '../../utils/office/officeTest';

import { TooFlowPage } from './tooTestFixture';

test.describe('TOO user', () => {
  /** @type {TooFlowPage} */
  let tooFlowPage;

  // Helper func to calculate days between 2 given dates
  // This is to support months of 30 and 31 dayss
  const calculateDaysDiff = (startDate, endDate) => {
    let days = 0;
    const currentDate = new Date(startDate);

    while (currentDate < endDate) {
      days += 1;
      currentDate.setDate(currentDate.getDate() + 1);
    }

    return days;
  };

  test.describe('previewing shipment with current SIT with past SIT', () => {
    test.beforeEach(async ({ officePage }) => {
      // build move in SIT with 90 days authorized and without pending extension requests
      // SIT entry date of 30 days ago, no departure date so it is current
      const move = await officePage.testHarness.buildHHGMoveInSIT();
      await officePage.signInAsNewTOOUser();
      tooFlowPage = new TooFlowPage(officePage, move);
      await tooFlowPage.waitForLoading();
      await officePage.tooNavigateToMove(tooFlowPage.moveLocator);
    });

    test('sum of days is correct between past Origin SIT and current Destination SIT', async ({ page }) => {
      // navigate to MTO tab
      await page.getByTestId('MoveTaskOrder-Tab').click();
      await tooFlowPage.waitForPage.moveTaskOrder();
      // assert that days authorization is 90
      await expect(page.getByTestId('sitStatusTable').getByText('90', { exact: true }).first()).toBeVisible();
      // get today
      const today = new Date();
      // get 1 month ago
      const oneMonthAgo = new Date(today);
      oneMonthAgo.setMonth(today.getMonth() - 1);
      // get 2 months ago
      const twoMonthsAgo = new Date(today);
      twoMonthsAgo.setMonth(today.getMonth() - 2);
      // get days between
      const daysBetweenTwoMonthsAgoAndOneMonthAgo = calculateDaysDiff(twoMonthsAgo, oneMonthAgo);
      const daysBetweenOneMonthAgoAndToday = calculateDaysDiff(oneMonthAgo, today);
      // get sums
      const totalDaysUsed = daysBetweenTwoMonthsAgoAndOneMonthAgo + daysBetweenOneMonthAgoAndToday;
      const remainingDays = 90 - totalDaysUsed;
      // assert that days used is the following sum
      // - past origin SIT (entry 2 months ago, departure 1 month ago)
      // - current destination SIT (entry 1 month ago, departure not given yet)
      await expect(
        page.getByTestId('sitStatusTable').getByText(`${totalDaysUsed}`, { exact: true }).first(),
      ).toBeVisible();
      // assert that days remaining is authorized minus totalDaysUsed
      await expect(
        page.getByTestId('sitStatusTable').getByText(`${remainingDays}`, { exact: true }).first(),
      ).toBeVisible();
      // assert that total days in destination sit is 1 month, inclusive of last day
      await expect(
        page.getByTestId('sitStartAndEndTable').getByText(`${daysBetweenOneMonthAgoAndToday}`, { exact: true }).first(),
      ).toBeVisible();

      // Get the SIT start date from the UI
      const sitStartDateText = await page.getByTestId('sitStartAndEndTable').locator('td').nth(0).innerText();
      const sitStartDate = new Date(sitStartDateText);

      // Calculate the authorized end date by adding 90 days to the start date and then subtracting total days used
      const authorizedEndDate = new Date(sitStartDate);
      // Use 89 because the last day is counted as a whole day
      // Subtract by "daysBetweenTwoMonthsAgoAndOneMonthAgo" (past SIT days) instead of total days used
      // This is because the authorized end date is based on the remaining days at the start of the current SIT
      authorizedEndDate.setDate(authorizedEndDate.getDate() + 89 - daysBetweenTwoMonthsAgoAndOneMonthAgo);
      // format
      const day = new Intl.DateTimeFormat('en', { day: '2-digit' }).format(authorizedEndDate);
      const month = new Intl.DateTimeFormat('en', { month: 'short' }).format(authorizedEndDate);
      const year = new Intl.DateTimeFormat('en', { year: 'numeric' }).format(authorizedEndDate);
      const expectedAuthorizedEndDate = `${day} ${month} ${year}`;
      // assert
      await expect(
        page.getByTestId('sitStartAndEndTable').getByText(`${expectedAuthorizedEndDate}`, { exact: true }).first(),
      ).toBeVisible();
    });
  });

  test.describe('previewing shipment with past origin and destination SIT', () => {
    test.beforeEach(async ({ officePage }) => {
      // build move in SIT with 90 days authorized and without pending extension requests
      // Origin sit had an entry date of four months ago, departure date of three months ago
      // Destination sit had an entry date of two months ago, departure date of one month ago
      const move = await officePage.testHarness.buildHHGMoveWithPastSITs();
      await officePage.signInAsNewTOOUser();
      tooFlowPage = new TooFlowPage(officePage, move);
      await tooFlowPage.waitForLoading();
      await officePage.tooNavigateToMove(tooFlowPage.moveLocator);
    });

    test('sum of days is correct between past Origin SIT and past Destination SIT', async ({ page }) => {
      // navigate to MTO tab
      await page.getByTestId('MoveTaskOrder-Tab').click();
      await tooFlowPage.waitForPage.moveTaskOrder();
      // assert that days authorization is 90
      await expect(page.getByTestId('sitStatusTable').getByText('90', { exact: true }).first()).toBeVisible();
      // get today
      const today = new Date();
      // get months
      const destinationDepartureDate = new Date(today);
      destinationDepartureDate.setMonth(today.getMonth() - 1);
      const destinationEntryDate = new Date(today);
      destinationEntryDate.setMonth(today.getMonth() - 2);
      const originDepartureDate = new Date(today);
      originDepartureDate.setMonth(today.getMonth() - 3);
      const originEntryDate = new Date(today);
      originEntryDate.setMonth(today.getMonth() - 4);
      // days between
      const totalDaysBetweenDestination = calculateDaysDiff(destinationEntryDate, destinationDepartureDate);
      const totalDaysBetweenOrigin = calculateDaysDiff(originEntryDate, originDepartureDate);
      // get sums
      const totalDaysUsed = totalDaysBetweenDestination + totalDaysBetweenOrigin;
      const remainingDays = 90 - totalDaysUsed;
      // assert sums
      await expect(
        page.getByTestId('sitStatusTable').getByText(`${totalDaysUsed}`, { exact: true }).first(),
      ).toBeVisible();
      // assert that days remaining is authorized minus totalDaysUsed
      await expect(
        page.getByTestId('sitStatusTable').getByText(`${remainingDays}`, { exact: true }).first(),
      ).toBeVisible();
      // assert previous sit days
      await expect(
        page
          .getByTestId('previouslyUsedSitTable')
          .getByText(`${totalDaysBetweenDestination} days at destination`, { exact: false })
          .first(),
      ).toBeVisible();
      await expect(
        page
          .getByTestId('previouslyUsedSitTable')
          .getByText(`${totalDaysBetweenOrigin} days at origin`, { exact: false })
          .first(),
      ).toBeVisible();
    });
  });

  test.describe('updating a move shipment in SIT', () => {
    test.beforeEach(async ({ officePage }) => {
      // build move in SIT with 90 days authorized and without pending extension requests
      const move = await officePage.testHarness.buildHHGMoveInSIT();
      await officePage.signInAsNewTOOUser();
      tooFlowPage = new TooFlowPage(officePage, move);
      await tooFlowPage.waitForLoading();
      await officePage.tooNavigateToMove(tooFlowPage.moveLocator);
    });

    test('is able to see the SIT Departure Date', async ({ page }) => {
      // navigate to MTO tab
      await page.getByTestId('MoveTaskOrder-Tab').click();
      await tooFlowPage.waitForPage.moveTaskOrder();

      const target = await page
        .getByTestId('currentSitDepartureDate')
        .locator('table[class="DataTable_dataTable__TGt9M table--data-point"]')
        .locator('tbody')
        .locator('td')
        .nth(0)
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
      expect(daysUsedCapture).toEqual('62'); // 31 days in past origin sit, 31 days in destination sit
      expect(daysLeftCapture).toEqual('28'); // of the 90 authorized, 62 have been used
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
      await officePage.tooNavigateToMove(tooFlowPage.moveLocator);
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
      await officePage.tooNavigateToMove(tooFlowPage.moveLocator);
    });

    test('is able to approve the SIT extension request', async ({ page }) => {
      // navigate to MTO tab
      await page.getByTestId('MoveTaskOrder-Tab').click();
      await tooFlowPage.waitForPage.moveTaskOrder();

      // assert that there is a pending SIT extension request
      await expect(page.getByText('SIT EXTENSION REQUESTED')).toBeVisible();

      // approve SIT extension with an adjusted approved days value of 100 days and change the extension reason
      await page.getByTestId('sitExtensions').getByTestId('button').click();
      await expect(page.getByRole('heading', { name: 'Review SIT Extension Request' })).toBeVisible();
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
      await expect(page.getByText('SIT EXTENSION REQUESTED')).toBeVisible();

      // deny SIT extension
      await page.getByTestId('sitExtensions').getByTestId('button').click();
      await expect(page.getByRole('heading', { name: 'Review SIT Extension Request' })).toBeVisible();
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
      await expect(page.getByText('SIT EXTENSION REQUESTED')).toBeVisible();

      // deny SIT extension
      await page.getByTestId('sitExtensions').getByTestId('button').click();
      await expect(page.getByRole('heading', { name: 'Review SIT Extension Request' })).toBeVisible();
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
      await expect(page.getByText('Total days in destination SIT')).toBeVisible();
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
