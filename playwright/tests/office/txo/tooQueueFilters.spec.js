/**
 * Semi-automated converted from a cypress test, and thus may contain
 * non best-practices, in particular: heavy use of `page.locator`
 * instead of `page.getBy*`.
 */

// @ts-check
import { test, expect } from '../../utils/office/officeTest';

import TooFlowPage from './tooTestFixture';

test.describe('TOO user queue filters - Move Queue', async () => {
  let testMove;
  let testMoveNotAssigned;
  let tooFlowPage;

  test.beforeEach(async ({ officePage }) => {
    if (!testMove) {
      testMove = await officePage.testHarness.buildMoveWithNTSShipmentsForTOO();
    }
    if (!testMoveNotAssigned) {
      testMoveNotAssigned = await officePage.testHarness.buildMoveWithNTSShipmentsForTOO();
    }
    await officePage.signInAsNewTOOUser();
    tooFlowPage = new TooFlowPage(officePage, testMove);
    await tooFlowPage.waitForLoading();
  });

  test('filters out all moves with nonsense assigned user', async ({ page }) => {
    test.slow();
    await page.getByTestId('assignedTo').getByRole('textbox').fill('');
    await page.getByTestId('assignedTo').getByRole('textbox').blur();
    // We should still see all moves
    await expect(page.getByRole('heading', { level: 1 })).not.toContainText('All moves (0)');

    // Add nonsense string to our filter (so now we're searching for 'zzzz')
    await page.getByTestId('assignedTo').getByRole('textbox').fill('zzzz');
    await page.getByTestId('assignedTo').getByRole('textbox').blur();
    // Now we shouldn't see any results
    await expect(page.getByRole('heading', { level: 1 })).toContainText('All moves (0)');
  });

  test('filters all moves EXCEPT with assigned user, restores all moves when filter removed', async ({ page }) => {
    test.slow();

    await page.getByText('KKFA moves').click();

    const allMoves = await page.getByRole('heading', { level: 1 }).innerHTML();

    await page.getByTestId('locator').getByRole('textbox').fill(testMove.locator);
    await page.getByTestId('locator').getByRole('textbox').blur();
    await expect(page.locator('td').getByText(testMove.locator)).toBeVisible();

    await page.getByTestId('assignedTo-0').getByTestId('dropdown').selectOption({ index: 1 });

    const assigned = (
      await page
        .getByTestId('assignedTo-0')
        .getByTestId('dropdown')
        .getByRole('option', { selected: true })
        .textContent()
    )
      .split(',')
      .at(0);

    await page.getByTestId('locator').getByRole('textbox').fill('');
    await page.getByTestId('locator').getByRole('textbox').blur();
    await expect(page.locator('remove-filters-locator').getByText(testMove.locator)).not.toBeVisible();

    await page.getByTestId('assignedTo').getByRole('textbox').fill(assigned);
    await page.getByTestId('assignedTo').getByRole('textbox').blur();

    await expect(page.getByRole('table').getByText(testMove.locator)).toBeVisible();
    await expect(page.getByRole('heading', { level: 1 })).not.toContainText(allMoves);
    await expect(page.getByRole('heading', { level: 1 })).not.toContainText('All moves(0)');

    // Now, remove filter and ensure retores all moves
    const removeFilterButton = page.getByTestId('remove-filters-assignedTo');
    await expect(removeFilterButton).toBeVisible();
    await removeFilterButton.click();
    await expect(page.getByRole('heading', { level: 1 })).toContainText(allMoves);
  });
});
