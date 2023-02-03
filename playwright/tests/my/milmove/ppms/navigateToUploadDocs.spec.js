/**
 * Semi-automated converted from a cypress test, and thus may contain
 * non best-practices, in particular: heavy use of `page.locator`
 * instead of `page.getBy*`.
 */

// @ts-check
const { expect, test, forEachViewport } = require('./customerPpmTestFixture');

test.describe('PPM Request Payment - Begin providing documents flow', () => {
  forEachViewport(async () => {
    test.beforeEach(async ({ customerPpmPage }) => {
      const move = await customerPpmPage.testHarness.buildApprovedMoveWithPPM();
      await customerPpmPage.signInForPPMWithMove(move);
    });

    test('has upload documents button enabled', async ({ page }) => {
      await expect(page.getByRole('heading', { name: 'Your move is in progress.' })).toBeVisible();
      const stepContainer5 = page.getByTestId('stepContainer5');
      await expect(stepContainer5.locator('p').getByText('15 Apr 2022')).toBeVisible();
      await stepContainer5.getByRole('button', { name: 'Upload PPM Documents' }).click();
      await expect(page).toHaveURL(/\/moves\/[^/]+\/shipments\/[^/]+\/about/);
    });
  });
});
