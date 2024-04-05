/**
 * Semi-automated converted from a cypress test, and thus may contain
 * non best-practices, in particular: heavy use of `page.locator`
 * instead of `page.getBy*`.
 */

// @ts-check
import { expect, test, forEachViewport } from './customerPpmTestFixture';

const multiMoveEnabled = process.env.FEATURE_FLAG_MULTI_MOVE;

test.describe('PPM Request Payment - Begin providing documents flow', () => {
  test.skip(multiMoveEnabled === 'true', 'Skip if MultiMove workflow is enabled.');
  forEachViewport(async () => {
    test.beforeEach(async ({ customerPpmPage }) => {
      const move = await customerPpmPage.testHarness.buildApprovedMoveWithPPM();
      await customerPpmPage.signInForPPMWithMove(move);
    });

    test('has upload documents button enabled', async ({ page }) => {
      await expect(page.getByRole('heading', { name: 'Your move is in progress.' })).toBeVisible();
      let stepContainer = page.getByTestId('stepContainer6');

      if (stepContainer == null) {
        stepContainer = page.getByTestId('stepContainer5');
      }

      await expect(stepContainer.locator('p').getByText('15 Apr 2022')).toBeVisible();
      await stepContainer.getByRole('button', { name: 'Upload PPM Documents' }).click();
      await expect(page).toHaveURL(/\/moves\/[^/]+\/shipments\/[^/]+\/about/);
    });
  });
});

test.describe('(MultiMove) PPM Request Payment - Begin providing documents flow', () => {
  test.skip(multiMoveEnabled === 'false', 'Skip if MultiMove workflow is not enabled.');

  forEachViewport(async () => {
    test.beforeEach(async ({ customerPpmPage }) => {
      const move = await customerPpmPage.testHarness.buildApprovedMoveWithPPM();
      await customerPpmPage.signInForPPMWithMove(move);
    });

    test('has upload documents button enabled', async ({ page }) => {
      await expect(page.getByRole('heading', { name: 'Your move is in progress.' })).toBeVisible();
      let stepContainer = page.getByTestId('stepContainer6');

      if (stepContainer == null) {
        stepContainer = page.getByTestId('stepContainer5');
      }

      await expect(stepContainer.locator('p').getByText('15 Apr 2022')).toBeVisible();
      await stepContainer.getByRole('button', { name: 'Upload PPM Documents' }).click();
      await expect(page).toHaveURL(/\/moves\/[^/]+\/shipments\/[^/]+\/about/);
    });
  });
});
