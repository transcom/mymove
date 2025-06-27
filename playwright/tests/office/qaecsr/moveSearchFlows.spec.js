/**
 * Semi-automated converted from a cypress test, and thus may contain
 * non best-practices, in particular: heavy use of `page.locator`
 * instead of `page.getBy*`.
 */

// @ts-check
import { test, expect } from '../../utils/office/officeTest';

test.describe('QAE Move Search', () => {
  test('is able to search by move code', async ({ page, officePage }) => {
    const move = await officePage.testHarness.buildHHGMoveWithNTSAndNeedsSC();
    const moveLocator = move.locator;

    await officePage.signInAsNewQAEUser();

    // Type move code into search bar (move code is default search type)
    await officePage.qaeSearchForAndNavigateToMove(moveLocator);
    await expect(page.locator('h1').getByText('Move Details', { exact: true })).toBeVisible();
  });

  test('is able to search by DOD ID', async ({ page, officePage }) => {
    const move = await officePage.testHarness.buildHHGMoveWithNTSAndNeedsSC();
    const moveLocator = move.locator;
    const { edipi } = move.Orders.service_member;

    await officePage.signInAsNewQAEUser();

    // Type dodID into search bar and select DOD ID as search type
    //
    // because of the styling of this input item, we cannot use a
    // css locator for the input item and then click it
    //
    // The styling is very similar to the issue described in
    //
    // https://github.com/microsoft/playwright/issues/3688
    //
    await page.getByText('DOD ID').click();
    await page.locator('input[name="searchText"]').fill(edipi);
    await page.locator('input[name="searchText"]').blur();

    // Search for moves
    await page.getByRole('button', { name: 'Search' }).click();
    await officePage.waitForLoading();

    // Verify results table contents
    await expect(page.locator('tbody >> tr')).toHaveCount(1);
    await expect(page.locator('tbody >> tr').first()).toContainText(moveLocator);
    await expect(page.locator('tbody >> tr').first()).toContainText(edipi);

    // Click result to navigate to move details page
    await page.locator('tbody > tr').first().click();
    await officePage.waitForLoading();

    expect(page.url()).toContain(`/moves/${moveLocator}/details`);
    await expect(page.locator('h1').getByText('Move Details', { exact: true })).toBeVisible();
  });

  test('is able to search by customer name', async ({ page, officePage }) => {
    const move = await officePage.testHarness.buildHHGMoveWithNTSAndNeedsSC();
    const moveLocator = move.locator;
    const lastName = move.Orders.service_member.last_name;

    await officePage.signInAsNewQAEUser();

    // Type name into search bar and select name as search type
    //
    // because of the styling of this input item, we cannot use a
    // css locator for the input item and then click it
    //
    // The styling is very similar to the issue described in
    //
    // https://github.com/microsoft/playwright/issues/3688
    //
    await page.getByText('Customer Name').click();
    await page.locator('input[name="searchText"]').fill(lastName);
    await page.locator('input[name="searchText"]').blur();

    // Search for moves
    await page.getByRole('button', { name: 'Search' }).click();
    await officePage.waitForLoading();

    // Verify results table contents
    await expect(page.locator('tbody >> tr')).toHaveCount(1);
    await expect(page.locator('tbody >> tr').first()).toContainText(lastName);

    // Click result to navigate to move details page
    await page.locator('tbody > tr').first().click();
    await officePage.waitForLoading();

    // Verify results table contents
    expect(page.url()).toContain(`/moves/${moveLocator}/details`);
    await expect(page.locator('h1').getByText('Move Details', { exact: true })).toBeVisible();
  });

  test('handles searches that do not return results', async ({ page, officePage }) => {
    await officePage.signInAsNewQAEUser();
    // Search for a bad move code
    await page.locator('input[name="searchText"]').fill('BAD_ID');
    await page.locator('input[name="searchText"]').blur();
    await page.getByRole('button', { name: 'Search' }).click();
    await officePage.waitForLoading();

    // Verify no results
    await expect(page.locator('[data-testid=table-queue] > h2')).toContainText('Results (0)');
    await expect(page.locator('[data-testid=table-queue] > h2')).toContainText('No results found.');
  });
});
