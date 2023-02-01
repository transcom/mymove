/**
 * Semi-automated converted from a cypress test, and thus may contain
 * non best-practices, in particular: heavy use of `page.locator`
 * instead of `page.getBy*`.
 */

// @ts-check
const { test, useMobileViewport, CustomerPpmPage } = require('./customerPpmTestFixture');

test.describe('Final Closeout', () => {
  /** @type {CustomerPpmPage} */
  let customerPpmPage;

  test.beforeEach(async ({ customerPage }) => {
    const move = await customerPage.testHarness.buildMoveWithPPMShipmentReadyForFinalCloseout();
    customerPpmPage = new CustomerPpmPage(customerPage, move);
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
      test('can see final closeout page with final estimated incentive and shipment totals', async () => {
        await customerPpmPage.signInAndNavigateToFinalCloseoutPage();

        await customerPpmPage.verifyFinalIncentiveAndTotals();
      });
    });
  });
});
