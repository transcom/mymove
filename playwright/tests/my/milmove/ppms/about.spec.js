/**
 * Semi-automated converted from a cypress test, and thus may contain
 * non best-practices, in particular: heavy use of `page.locator`
 * instead of `page.getBy*`.
 */

// @ts-check
import { test, forEachViewport } from './customerPpmTestFixture';

const multiMoveEnabled = process.env.FEATURE_FLAG_MULTI_MOVE;

test.describe('About Your PPM', () => {
  forEachViewport(async () => {
    test.beforeEach(async ({ customerPpmPage }) => {
      const move = await customerPpmPage.testHarness.buildApprovedMoveWithPPM();
      await customerPpmPage.signInForPPMWithMove(move);
    });

    [true, false].forEach((selectAdvance) => {
      const advanceText = selectAdvance ? 'with' : 'without';
      test(`can submit actual PPM shipment info ${advanceText} an advance`, async ({ customerPpmPage }) => {
        if (multiMoveEnabled) {
          // Need to have move show on multiMoveLanding Page. Click Go to Move, then ...
          // await customerPpmPage.navigateToAboutPage({ selectAdvance });
        } else {
          await customerPpmPage.navigateToAboutPage({ selectAdvance });
        }
      });
    });
  });
});
