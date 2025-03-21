// @ts-check
import { test, expect } from '../../utils/office/officeTest';

import TooFlowPage from './tooTestFixture';

const waitForFilterInput = async (page, testId, value) => {
  await page.waitForFunction(
    (args) => {
      /** @type {HTMLInputElement} */
      const input = document.querySelector(`[data-testid="${args.testId}"] [data-testid="TextBoxFilter"]`);
      if (input && input.value === args.value) return true;
      if (input) {
        input.value = args.value;
        const keyupEvent = new KeyboardEvent('keyup', {
          code: 'Enter',
          key: 'Enter',
          bubbles: true,
          cancelable: true,
        });
        input.dispatchEvent(keyupEvent);
      }
      return false;
    },
    { testId, value },
    { polling: 500, timeout: 5000 },
  );
};

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

    // Add nonsense string to our filter (so now we're searching for 'abcde')
    await waitForFilterInput(page, 'assignedTo', 'abcde');

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
    await waitForFilterInput(page, 'locator', testMove.locator);

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

    // search for moves with assigned user filter only
    await waitForFilterInput(page, 'locator', '');

    await expect(page.getByTestId('locator').getByTestId('TextBoxFilter')).toBeEmpty();
    await expect(page.getByRole('heading', { level: 1 })).toContainText(allMoves);

    await waitForFilterInput(page, 'assignedTo', assigned);

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

test.describe('TOO user queue filters - Destination Requests Queue', async () => {
  let testMove;
  let testMoveNotAssigned;
  let tooFlowPage;

  test.beforeEach(async ({ officePage }) => {
    testMove = await officePage.testHarness.buildHHGMoveInSITWithPendingExtension();
    testMoveNotAssigned = await officePage.testHarness.buildHHGMoveInSITWithPendingExtension();
    await officePage.signInAsNewTOOUser();
    tooFlowPage = new TooFlowPage(officePage, testMove);
    await tooFlowPage.waitForLoading();
  });

  test('filters out all moves with nonsense assigned user', async ({ page }) => {
    // We should still see all moves
    await Promise.all([
      page.waitForResponse((res) => res.url().includes('/ghc/v1/queues/destination-requests')),
      page.getByTestId('destination-requests-tab-link').click(),
    ]);
    await expect(page.getByRole('heading', { level: 1 })).not.toContainText('Destination requests (0)');

    // Add nonsense string to our filter (so now we're searching for 'abcde')
    await waitForFilterInput(page, 'assignedTo', 'abcde');

    // Now we shouldn't see any results
    await expect(page.getByRole('heading', { level: 1 })).toContainText('Destination requests (0)');
    await expect(page.getByRole('row').getByText(testMove.locator)).not.toBeVisible();
    await expect(page.getByRole('row').getByText(testMoveNotAssigned.locator)).not.toBeVisible();
  });

  test('filters all moves EXCEPT with assigned user, restores all moves when filter removed', async ({ page }) => {
    // We should still see all moves
    await Promise.all([
      page.waitForResponse((res) => res.url().includes('/ghc/v1/queues/destination-requests')),
      page.getByTestId('destination-requests-tab-link').click(),
    ]);

    await expect(page.getByRole('heading', { level: 1 })).not.toContainText('Destination requests (0)');
    const destinationRequests = await page.getByRole('heading', { level: 1 }).innerHTML();

    // assign user to test move
    await waitForFilterInput(page, 'locator', testMove.locator);

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

    // search for moves with assigned user filter only
    await waitForFilterInput(page, 'locator', '');

    await expect(page.getByTestId('locator').getByTestId('TextBoxFilter')).toBeEmpty();
    await expect(page.getByRole('heading', { level: 1 })).toContainText(destinationRequests);

    await waitForFilterInput(page, 'assignedTo', assigned);

    // assigned moves for user show in queue
    await expect(page.getByRole('table').getByText(testMove.locator)).toBeVisible();
    await expect(page.getByRole('heading', { level: 1 })).not.toContainText(destinationRequests);
    await expect(page.getByRole('heading', { level: 1 })).not.toContainText('Destination requests (0)');
    await expect(page.getByRole('row').getByText(testMove.locator)).toBeVisible();
    await expect(page.getByRole('row').getByText(testMoveNotAssigned.locator)).not.toBeVisible();

    // Now, remove filter and ensure retores all moves in queue
    const removeFilterButton = page.getByTestId('remove-filters-assignedTo');
    await expect(removeFilterButton).toBeVisible();
    await removeFilterButton.click();
    await expect(page.getByRole('heading', { level: 1 })).toContainText(destinationRequests);
  });
});
