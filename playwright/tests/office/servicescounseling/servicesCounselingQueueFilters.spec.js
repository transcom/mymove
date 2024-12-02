/**
 * Semi-automated converted from a cypress test, and thus may contain
 * non best-practices, in particular: heavy use of `page.locator`
 * instead of `page.getBy*`.
 */

// @ts-check
import { test, expect } from './servicesCounselingTestFixture';

test.describe('Services counselor user', () => {
  let moveLocator = '';
  let moveWithNeedsCloseoutLocator = '';
  test.slow();
  test.describe('with PPM shipment ready for closeout', () => {
    let dutyLocationName = '';
    test.beforeEach(async ({ scPage }) => {
      const move = await scPage.testHarness.buildPartialPPMMoveReadyForCloseout();
      const moveWithNeedsCloseout = await scPage.testHarness.buildApprovedMoveWithPPMAllDocTypesOffice();
      moveLocator = move.locator;
      moveWithNeedsCloseoutLocator = moveWithNeedsCloseout.locator;
      dutyLocationName = move.Orders.NewDutyLocation.name;
      await scPage.page.locator('[data-testid="closeout-tab-link"]').click();
    });

    test('is able to filter partial vs full moves based on ppm type', async ({ page }) => {
      // closeout tab

      // Created a single Partial PPM move, so when we search for
      // Partial, we should see it in the results
      await page.locator('th[data-testid="locator"] > div > input').fill(moveLocator);
      await page.locator('th[data-testid="locator"] > div > input').blur();

      await page.locator('th[data-testid="ppmType"] > div > select').selectOption({ label: 'Partial' });

      await expect(page.locator('td').getByText(moveLocator)).toBeVisible();

      // When we search for Full PPM moves, partial move should not come up
      await page.locator('th[data-testid="ppmType"] > div > select').selectOption({ label: 'Full' });
      await expect(page.locator('h1').getByText('Moves (0)')).toBeVisible();
    });

    test('is able to filter moves based on PPM status', async ({ page }) => {
      // Check for Needs closeout filter
      await page.locator('th[data-testid="locator"] > div > input').clear();
      await page.locator('th[data-testid="locator"] > div > input').fill(moveWithNeedsCloseoutLocator);
      await page.locator('th[data-testid="locator"] > div > input').blur();
      await page.locator('th[data-testid="ppmStatus"] > div > select').selectOption({ label: 'Needs closeout' });
      await expect(page.locator('td').getByText(moveWithNeedsCloseoutLocator)).toBeVisible();
    });

    test('is able to filter moves based on destination duty location', async ({ page }) => {
      // add filter for move code (PPM closeout that has Fort Gordon as
      // its destination duty location)

      await page.locator('th[data-testid="locator"] > div > input').fill(moveLocator);
      await page.locator('th[data-testid="locator"] > div > input').blur();

      const dutyLocationPrefix = dutyLocationName.substring(0, 4);

      // Add destination duty location filter for the first part of
      // the name
      await page.locator('th[data-testid="destinationDutyLocation"] > div > input').fill(dutyLocationPrefix);
      await page.locator('th[data-testid="destinationDutyLocation"] > div > input').blur();
      // We should still see our move
      await expect(page.locator('td').getByText(moveLocator)).toBeVisible();

      // Add nonsense string to our filter (so now we're searching for 'fortzzzz')
      await page.locator('th[data-testid="destinationDutyLocation"] > div > input').fill('zzzz');
      await page.locator('th[data-testid="destinationDutyLocation"] > div > input').blur();
      // Now we shouldn't see any results
      await expect(page.locator('h1')).toContainText('Moves (0)');
    });
  });

  test.describe('with PPM move with closeout', () => {
    test.beforeEach(async ({ scPage }) => {
      const move = await scPage.testHarness.buildPPMMoveWithCloseout();
      moveLocator = move.locator;
      await scPage.page.locator('[data-testid="closeout-tab-link"]').click();
    });

    test('is able to filter moves based on PPM Closeout initiated', async ({ page }) => {
      const closeoutDate = new Date().toLocaleDateString('en-US');

      // first test with bogus date and no moves are found
      await page.locator('th[data-testid="closeoutInitiated"] > div > div > input').fill('11 Dec 2020');
      await expect(page.locator('h1')).toContainText('Moves (0)');

      // test with the closeout date of the created move and that our
      // move is found
      await page.getByTestId('remove-filters-closeoutInitiated').click();
      await page.locator('th[data-testid="closeoutInitiated"] > div > div > input').fill(closeoutDate);
      await page.getByRole('cell', { name: 'Previous Month Next Month' }).getByRole('textbox').fill(closeoutDate);
      await expect(page.locator('h1')).not.toContainText('Moves (0)');
      await page.getByLabel('rows per page').click();
      await page.getByLabel('rows per page').selectOption('50');
      await expect(page.getByText(moveLocator)).toBeVisible();
    });
  });

  test.describe('with PPM move with closeout office', () => {
    let closeoutOffice = '';
    test.beforeEach(async ({ scPage }) => {
      const move = await scPage.testHarness.buildPPMMoveWithCloseoutOffice();
      moveLocator = move.locator;
      closeoutOffice = move.CloseoutOffice.name;
      await scPage.page.locator('[data-testid="closeout-tab-link"]').click();
    });

    test('is able to filter moves based on PPM Closeout location', async ({ page }) => {
      await page.locator('th[data-testid="locator"] > div > input').fill(moveLocator);
      await page.locator('th[data-testid="locator"] > div > input').blur();
      // add another filter for the closeout office column checking
      // it's not case sensitive
      await page.locator('th[data-testid="closeoutLocation"] > div > input').fill(closeoutOffice.toUpperCase());
      await page.locator('th[data-testid="closeoutLocation"] > div > input').blur();

      await expect(page.locator('td').getByText(moveLocator)).toBeVisible();
      // Add some nonsense z text to our filter
      await page.locator('th[data-testid="closeoutLocation"] > div > input').fill('z');
      await page.locator('th[data-testid="closeoutLocation"] > div > input').blur();
      // now we should get no results
      await expect(page.locator('h1')).toContainText('Moves (0)');
    });
  });
});
