/**
 * Semi-automated converted from a cypress test, and thus may contain
 * non best-practices, in particular: heavy use of `page.locator`
 * instead of `page.getBy*`.
 */

// @ts-check
const { expect, test } = require('./servicesCounselingTestFixture');

test.describe('Services counselor user', () => {
  let moveLocator = '';

  test.describe('with PPM shipment ready for closeout', () => {
    let dutyLocationName = '';
    test.beforeEach(async ({ scPage }) => {
      const move = await scPage.testHarness.buildMoveWithPPMShipmentReadyForFinalCloseout();
      moveLocator = move.locator;
      dutyLocationName = move.Orders.NewDutyLocation.name;
      await scPage.page.locator('[data-testid="closeout-tab-link"]').click();
    });

    test('is able to filter partial vs full moves based on ppm type', async ({ page }) => {
      // closeout tab

      // Created a single Partial PPM move, so when we search for
      // Partial, we should see it in the results
      await page.locator('th[data-testid="locator"] > div > input').type(moveLocator);
      await page.locator('th[data-testid="locator"] > div > input').blur();

      await page.locator('th[data-testid="ppmType"] > div > select').selectOption({ label: 'Partial' });

      await expect(page.locator('td').getByText(moveLocator)).toBeVisible();

      // When we search for Full PPM moves, partial move should not come up
      await page.locator('th[data-testid="ppmType"] > div > select').selectOption({ label: 'Full' });
      await expect(page.locator('h1').getByText('Moves (0)')).toBeVisible();
    });

    test('is able to filter moves based on destination duty location', async ({ page }) => {
      // add filter for move code (PPM closeout that has Fort Gordon as
      // its destination duty location)

      await page.locator('th[data-testid="locator"] > div > input').type(moveLocator);
      await page.locator('th[data-testid="locator"] > div > input').blur();

      const dutyLocationPrefix = dutyLocationName.substring(0, 4);

      // Add destination duty location filter for the first part of
      // the name
      await page.locator('th[data-testid="destinationDutyLocation"] > div > input').type(dutyLocationPrefix);
      await page.locator('th[data-testid="destinationDutyLocation"] > div > input').blur();
      // We should still see our move
      await expect(page.locator('td').getByText(moveLocator)).toBeVisible();

      // Add nonsense string to our filter (so now we're searching for 'fortzzzz')
      await page.locator('th[data-testid="destinationDutyLocation"] > div > input').type('zzzz');
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
      await page.locator('th[data-testid="closeoutInitiated"] > div > div > input').type('11 Dec 2020');
      await expect(page.locator('h1')).toContainText('Moves (0)');

      // test with the closeout date of the created move and that our
      // move is found
      await page.locator('th[data-testid="closeoutInitiated"] > div > div > input').clear();
      await page.locator('th[data-testid="closeoutInitiated"] > div > div > input').type(closeoutDate);
      await expect(page.locator('h1')).not.toContainText('Moves (0)');
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
      await page.locator('th[data-testid="locator"] > div > input').type(moveLocator);
      await page.locator('th[data-testid="locator"] > div > input').blur();
      // add another filter for the closeout office column checking
      // it's not case sensitive
      await page.locator('th[data-testid="closeoutLocation"] > div > input').type(closeoutOffice.toUpperCase());
      await page.locator('th[data-testid="closeoutLocation"] > div > input').blur();

      await expect(page.locator('td').getByText(moveLocator)).toBeVisible();
      // Add some nonsense z text to our filter
      await page.locator('th[data-testid="closeoutLocation"] > div > input').type('z');
      await page.locator('th[data-testid="closeoutLocation"] > div > input').blur();
      // now we should get no results
      await expect(page.locator('h1')).toContainText('Moves (0)');
    });
  });
});
