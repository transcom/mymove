/**
 * Semi-automated converted from a cypress test, and thus may contain
 * non best-practices, in particular: heavy use of `page.locator`
 * instead of `page.getBy*`.
 */

// @ts-check
import { test, forEachViewport } from './customerPpmTestFixture';

test.describe('About Your PPM', () => {
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

test.describe('About You PPM with incorrect Xlsx Upload', () => {
  forEachViewport(async () => {
    test.beforeEach(async ({ customerPpmPage }) => {
      const move = await customerPpmPage.testHarness.buildApprovedMoveWithPPM();
      await customerPpmPage.signInForPPMWithMove(move);
      await customerPpmPage.clickOnGoToMoveButton();
      await customerPpmPage.navigateToAboutPage({ selectAdvance: false });
    });

    test('Fill out weight ticket with incorrect Xlsx document', async ({ customerPpmPage }) => {
      await customerPpmPage.submitIncorrectXlsxFileForWeightTicketPage();
    });
  });
});

test.describe('(MultiMove) About Your PPM', () => {
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
