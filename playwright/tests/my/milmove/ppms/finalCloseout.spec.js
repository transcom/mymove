/**
 * Semi-automated converted from a cypress test, and thus may contain
 * non best-practices, in particular: heavy use of `page.locator`
 * instead of `page.getBy*`.
 */

// @ts-check
import { test, forEachViewport } from './customerPpmTestFixture';

test.describe('Final Closeout', () => {
  forEachViewport(async () => {
    test.beforeEach(async ({ customerPpmPage }) => {
      const move = await customerPpmPage.testHarness.buildMoveWithPPMShipmentReadyForFinalCloseout();
      await customerPpmPage.signInForPPMWithMove(move);
    });

    test('can see final closeout page with final estimated incentive and shipment totals', async ({
      customerPpmPage,
    }) => {
      await customerPpmPage.navigateToFinalCloseoutPage();

      await customerPpmPage.verifyFinalIncentiveAndTotals();
    });
  });
});
