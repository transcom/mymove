// @ts-check
const { test, expect } = require('../../utils/officeTest');

test.describe('QAE/CSR Move Search', () => {
  test('is able to search by move code', async ({ page, officePage }) => {
    const move = await officePage.buildHHGMoveWithNTSAndNeedsSC();
    const moveCode = move.locator;

    await officePage.signInAsNewQAECSRUser();

    // Type move code into search bar (move code is default search type)
    await officePage.qaeCsrSearchForAndNavigateToMove(moveCode);
    await expect(page.locator('h1').getByText('Move details', { exact: true })).toBeVisible();
  });

  test('is able to search by DOD ID', async ({ page, officePage }) => {
    const move = await officePage.buildHHGMoveWithNTSAndNeedsSC();
    const moveCode = move.locator;
    const { edipi } = move.Orders.ServiceMember;

    await officePage.signInAsNewQAECSRUser();

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
    await page.locator('input[name="searchText"]').type(edipi);
    await page.locator('input[name="searchText"]').blur();

    // Search for moves
    await page.getByRole('button', { name: 'Search' }).click();
    await officePage.waitForLoading();

    // Verify results table contents
    await expect(page.locator('tbody >> tr')).toHaveCount(1);
    expect(page.locator('tbody >> tr').first()).toContainText(moveCode);
    expect(page.locator('tbody >> tr').first()).toContainText(edipi);

    // Click result to navigate to move details page
    await page.locator('tbody > tr').first().click();
    await officePage.waitForLoading();

    expect(page.url()).toContain(`/moves/${moveCode}/details`);
    await expect(page.locator('h1').getByText('Move details', { exact: true })).toBeVisible();
  });

  test('is able to search by customer name', async ({ page, officePage }) => {
    const move = await officePage.buildHHGMoveWithNTSAndNeedsSC();
    const moveCode = move.locator;
    const lastName = move.Orders.ServiceMember.last_name;

    await officePage.signInAsNewQAECSRUser();

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
    await page.locator('input[name="searchText"]').type(lastName);
    await page.locator('input[name="searchText"]').blur();

    // Search for moves
    await page.getByRole('button', { name: 'Search' }).click();
    await officePage.waitForLoading();

    // Verify results table contents
    await expect(page.locator('tbody >> tr')).toHaveCount(1);
    expect(page.locator('tbody >> tr').first()).toContainText(lastName);

    // Click result to navigate to move details page
    await page.locator('tbody > tr').first().click();
    await officePage.waitForLoading();

    // Verify results table contents
    expect(page.url()).toContain(`/moves/${moveCode}/details`);
    await expect(page.locator('h1').getByText('Move details', { exact: true })).toBeVisible();
  });

  test('handles searches that do not return results', async ({ page, officePage }) => {
    await officePage.signInAsNewQAECSRUser();
    // Search for a bad move code
    await page.locator('input[name="searchText"]').type('BAD_ID');
    await page.locator('input[name="searchText"]').blur();
    await page.getByRole('button', { name: 'Search' }).click();
    await officePage.waitForLoading();

    // Verify no results
    await expect(page.locator('[data-testid=table-queue] > h2')).toContainText('Results (0)');
    await expect(page.locator('[data-testid=table-queue] > p')).toContainText('No results found.');
  });
});
