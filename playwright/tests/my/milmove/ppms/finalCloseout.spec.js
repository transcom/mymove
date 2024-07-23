/**
 * Semi-automated converted from a cypress test, and thus may contain
 * non best-practices, in particular: heavy use of `page.locator`
 * instead of `page.getBy*`.
 */

// @ts-check
import { test, forEachViewport } from './customerPpmTestFixture';

const multiMoveEnabled = process.env.FEATURE_FLAG_MULTI_MOVE;

test.describe('Final Closeout', () => {
  test.skip(multiMoveEnabled === 'true', 'Skip if MultiMove workflow is enabled.');
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

test.describe('(MultiMove) Final Closeout', () => {
  test.skip(multiMoveEnabled === 'false', 'Skip if MultiMove workflow is not enabled.');

  forEachViewport(async () => {
    test.beforeEach(async ({ customerPpmPage }) => {
      const move = await customerPpmPage.testHarness.buildMoveWithPPMShipmentReadyForFinalCloseout();
      await customerPpmPage.signInForPPMWithMove(move);
      await customerPpmPage.clickOnGoToMoveButton();
    });

    test('can see final closeout page with final estimated incentive and shipment totals', async ({
      customerPpmPage,
    }) => {
      await customerPpmPage.navigateToPPMReviewPageWithCompletePPM();
      await customerPpmPage.navigateFromPPMReviewPageToFinalCloseoutPage();

      await customerPpmPage.verifyFinalIncentiveAndTotals();
    });
  });
});
