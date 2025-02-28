// @ts-check
import { test, expect } from '../../utils/office/officeTest';

import TooFlowPage from './tooTestFixture';

test.describe('TOO user queue filters - Move Queue', async () => {
  let testMove;
  let testMoveNotAssigned;
  let tooFlowPage;

  test.beforeEach(async ({ officePage }) => {
    testMove = await officePage.testHarness.buildMoveWithNTSShipmentsForTOO();
    testMoveNotAssigned = await officePage.testHarness.buildMoveWithNTSShipmentsForTOO();
    await officePage.signInAsNewTOOUser();
    tooFlowPage = new TooFlowPage(officePage, testMove);
    await tooFlowPage.waitForLoading();
  });

  test('filters out all moves with nonsense assigned user', async ({ page }) => {
    // We should still see all moves
    await expect(page.getByRole('heading', { level: 1 })).not.toContainText('All moves (0)');

    // Add nonsense string to our filter (so now we're searching for 'zzzz')
    await page.waitForTimeout(100);
    await page.getByTestId('assignedTo').getByTestId('TextBoxFilter').fill('zzzz');
    await page.waitForTimeout(100);
    await Promise.all([
      page.waitForResponse((res) => res.url().includes('/ghc/v1/queues/moves')),
      page.getByTestId('assignedTo').getByTestId('TextBoxFilter').blur(),
    ]);
    // Now we shouldn't see any results
    await expect(page.getByRole('heading', { level: 1 })).toContainText('All moves (0)');
    await expect(page.getByRole('row').getByText(testMove.locator)).not.toBeVisible();
    await expect(page.getByRole('row').getByText(testMoveNotAssigned.locator)).not.toBeVisible();
  });

  test('filters all moves EXCEPT with assigned user, restores all moves when filter removed', async ({ page }) => {
    await Promise.all([
      page.waitForResponse((res) => res.url().includes('/ghc/v1/queues/moves')),
      page.getByText('KKFA moves').click(),
    ]);

    await expect(page.getByRole('heading', { level: 1 })).not.toContainText('All moves (0)');
    const allMoves = await page.getByRole('heading', { level: 1 }).innerHTML();

    // assign user to test move
    await page.waitForTimeout(100);
    await page.getByTestId('locator').getByTestId('TextBoxFilter').fill(testMove.locator);
    await Promise.all([
      page.waitForResponse((res) => res.url().includes('/ghc/v1/queues/moves')),
      page.getByTestId('locator').getByTestId('TextBoxFilter').blur(),
    ]);
    await expect(page.getByTestId('locator').getByTestId('TextBoxFilter')).toHaveValue(testMove.locator);
    await expect(page.getByRole('row').getByText(testMove.locator)).toBeVisible();

    expect(
      await page
        .getByTestId('assignedTo-0')
        .getByTestId('dropdown')
        .getByRole('option', { selected: true })
        .textContent(),
    ).toEqual('—');

    await Promise.all([
      page.waitForResponse((res) => res.url().includes('assignOfficeUser')),
      page.getByTestId('assignedTo-0').getByTestId('dropdown').selectOption({ index: 1 }),
    ]);

    const assigned = (
      await page
        .getByTestId('assignedTo-0')
        .getByTestId('dropdown')
        .getByRole('option', { selected: true })
        .textContent()
    )
      .split(',')
      .at(0);

    // search for moves with assigned user only
    await page.waitForTimeout(100);
    await page.getByTestId('locator').getByTestId('TextBoxFilter').fill('');
    await Promise.all([
      page.waitForResponse((res) => res.url().includes('/ghc/v1/queues/moves')),
      page.getByTestId('locator').getByTestId('TextBoxFilter').blur(),
    ]);

    await expect(page.getByTestId('locator').getByTestId('TextBoxFilter')).toBeEmpty();
    await expect(page.getByRole('heading', { level: 1 })).toContainText(allMoves);

    await page.waitForTimeout(100);
    await page.getByTestId('assignedTo').getByTestId('TextBoxFilter').fill(assigned);
    await Promise.all([
      page.waitForResponse((res) => res.url().includes('/ghc/v1/queues/moves')),
      page.getByTestId('assignedTo').getByTestId('TextBoxFilter').blur(),
    ]);
    await expect(page.getByTestId('remove-filters-assignedTo')).toBeVisible();

    // assigned moves for user show in queue
    await expect(page.getByRole('table').getByText(testMove.locator)).toBeVisible();
    await expect(page.getByRole('heading', { level: 1 })).not.toContainText(allMoves);
    await expect(page.getByRole('heading', { level: 1 })).not.toContainText('All moves (0)');
    await expect(page.getByRole('row').getByText(testMove.locator)).toBeVisible();
    await expect(page.getByRole('row').getByText(testMoveNotAssigned.locator)).not.toBeVisible();

    // Now, remove filter and ensure retores all moves in queue
    const removeFilterButton = page.getByTestId('remove-filters-assignedTo');
    await expect(removeFilterButton).toBeVisible();
    await removeFilterButton.click();
    await expect(page.getByRole('heading', { level: 1 })).toContainText(allMoves);
  });
});
