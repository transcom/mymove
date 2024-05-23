/**
 * Semi-automated converted from a cypress test, and thus may contain
 * non best-practices, in particular: heavy use of `page.locator`
 * instead of `page.getBy*`.
 */

// @ts-check
import { test, forEachViewport } from './customerPpmTestFixture';

const multiMoveEnabled = process.env.FEATURE_FLAG_MULTI_MOVE;

test.describe('About Your PPM', () => {
  test.skip(multiMoveEnabled === 'true', 'Skip if MultiMove workflow is enabled.');
  forEachViewport(async () => {
    test.beforeEach(async ({ customerPpmPage }) => {
      const move = await customerPpmPage.testHarness.buildApprovedMoveWithPPMWithAboutFormComplete();
      await customerPpmPage.signInForPPMWithMove(move);
      await customerPpmPage.navigateToWeightTicketPage();
    });

    test('proceed with weight ticket documents', async ({ customerPpmPage }) => {
      await customerPpmPage.submitWeightTicketPage();
    });

    test('proceed with claiming trailer', async ({ customerPpmPage }) => {
      await customerPpmPage.submitWeightTicketPage({ hasTrailer: true, ownTrailer: true });
    });

    test('proceed without claiming trailer', async ({ customerPpmPage }) => {
      await customerPpmPage.submitWeightTicketPage({ hasTrailer: true, ownTrailer: false });
    });

    test('proceed with constructed weight ticket documents', async ({ customerPpmPage }) => {
      await customerPpmPage.submitWeightTicketPage({ useConstructedWeight: true });
    });
  });
});

test.describe('(MultiMove) About Your PPM', () => {
  test.skip(multiMoveEnabled === 'false', 'Skip if MultiMove workflow is not enabled.');

  forEachViewport(async () => {
    test.beforeEach(async ({ customerPpmPage }) => {
      const move = await customerPpmPage.testHarness.buildApprovedMoveWithPPMWithAboutFormComplete();
      await customerPpmPage.signInForPPMWithMove(move);
      await customerPpmPage.clickOnGoToMoveButton();
      await customerPpmPage.navigateToWeightTicketPage();
    });

    test('proceed with weight ticket documents', async ({ customerPpmPage }) => {
      await customerPpmPage.submitWeightTicketPage();
    });

    test('proceed with claiming trailer', async ({ customerPpmPage }) => {
      await customerPpmPage.submitWeightTicketPage({ hasTrailer: true, ownTrailer: true });
    });

    test('proceed without claiming trailer', async ({ customerPpmPage }) => {
      await customerPpmPage.submitWeightTicketPage({ hasTrailer: true, ownTrailer: false });
    });

    test('proceed with constructed weight ticket documents', async ({ customerPpmPage }) => {
      await customerPpmPage.submitWeightTicketPage({ useConstructedWeight: true });
    });
  });
});
