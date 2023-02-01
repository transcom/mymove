/**
 * Semi-automated converted from a cypress test, and thus may contain
 * non best-practices, in particular: heavy use of `page.locator`
 * instead of `page.getBy*`.
 */

// @ts-check
const { expect, test, useMobileViewport, CustomerPpmPage } = require('./customerPpmTestFixture');

test.describe('PPM Request Payment - Begin providing documents flow', () => {
  test.beforeEach(async ({ customerPage }) => {
    const move = await customerPage.testHarness.buildApprovedMoveWithPPM();
    const customerPpmPage = new CustomerPpmPage(customerPage, move);
    await customerPpmPage.signInForPPM();
  });

  //
  // https://playwright.dev/docs/test-parameterize
  //
  // use forEach to avoid
  // https://eslint.org/docs/latest/rules/no-loop-func
  [true, false].forEach((isMobile) => {
    const viewportName = isMobile ? 'mobile' : 'desktop';
    test.describe(`with ${viewportName} viewport`, async () => {
      if (isMobile) {
        useMobileViewport();
      }
      test('has upload documents button enabled', async ({ page }) => {
        await expect(page.getByRole('heading', { name: 'Your move is in progress.' })).toBeVisible();
        const stepContainer5 = page.getByTestId('stepContainer5');
        await expect(stepContainer5.locator('p').getByText('15 Apr 2022')).toBeVisible();
        await stepContainer5.getByRole('button', { name: 'Upload PPM Documents' }).click();
        await expect(page).toHaveURL(/\/moves\/[^/]+\/shipments\/[^/]+\/about/);
      });
    });
  });
});
