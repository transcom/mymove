/**
 * Semi-automated converted from a cypress test, and thus may contain
 * non best-practices, in particular: heavy use of `page.locator`
 * instead of `page.getBy*`.
 */

// @ts-check
const { test, forEachViewport, CustomerPpmPage } = require('./customerPpmTestFixture');

test.describe('Final Closeout', () => {
  /** @type {CustomerPpmPage} */
  let customerPpmPage;

  test.beforeEach(async ({ customerPage }) => {
    const move = await customerPage.testHarness.buildMoveWithPPMShipmentReadyForFinalCloseout();
    customerPpmPage = new CustomerPpmPage(customerPage, move);
  });

  forEachViewport(async () => {
    test('can see final closeout page with final estimated incentive and shipment totals', async () => {
      await customerPpmPage.signInAndNavigateToFinalCloseoutPage();

      await customerPpmPage.verifyFinalIncentiveAndTotals();
    });
  });
});
