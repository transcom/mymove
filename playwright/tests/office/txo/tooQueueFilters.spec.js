/**
 * Semi-automated converted from a cypress test, and thus may contain
 * non best-practices, in particular: heavy use of `page.locator`
 * instead of `page.getBy*`.
 */

// @ts-check
import { test, expect } from '../../utils/office/officeTest';

import TooFlowPage from './tooTestFixture';

test.describe('TOO user queue filters - Move Queue', () => {
  let testMove;
  let tooFlowPage;

  test.beforeEach(async ({ officePage }) => {
    testMove = await officePage.testHarness.buildMoveWithNTSShipmentsForTOO();
    await officePage.signInAsNewTOOUser();
    tooFlowPage = new TooFlowPage(officePage, testMove);
    await tooFlowPage.waitForLoading();
  });

  test('filters out all moves with nonsense assigned user', async ({ page }) => {
    test.slow();
    await page.locator('th[data-testid="assignedTo"] > div > input').fill('');
    await page.locator('th[data-testid="assignedTo"] > div > input').blur();
    // We should still see all moves
    await expect(page.locator('h1')).not.toContainText('All moves (0)');

    // Add nonsense string to our filter (so now we're searching for 'zzzz')
    await page.locator('th[data-testid="assignedTo"] > div > input').fill('zzzz');
    await page.locator('th[data-testid="assignedTo"] > div > input').blur();
    // Now we shouldn't see any results
    await expect(page.locator('h1')).toContainText('All moves (0)');
  });

  test('filters all moves EXCEPT with assigned user', async ({ page }) => {
    test.slow();

    await page.getByText('KKFA moves').click();

    await page.locator('th[data-testid="locator"] > div > input').fill(testMove.locator);
    await page.locator('th[data-testid="locator"] > div > input').blur();
    await expect(page.locator('td').getByText(testMove.locator)).toBeVisible();

    await page.locator('td[data-testid="assignedTo-0"] > div > select').selectOption({ index: 1 });
    const assigned = (await page.locator('td[data-testid="assignedTo-0"] > div > select option:checked').textContent())
      .split(',')
      .at(0);

    // Add nonsense string to our filter (so now we're searching for 'zzzz')
    await page.locator('th[data-testid="assignedTo"] > div > input').fill(assigned);
    await page.locator('th[data-testid="assignedTo"] > div > input').blur();

    await expect(page.locator('td').getByText(testMove.locator)).toBeVisible();
    await expect(page.locator('h1')).toContainText('All moves (1)');
  });
});
