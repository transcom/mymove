/**
 * Semi-automated converted from a cypress test, and thus may contain
 * non best-practices, in particular: heavy use of `page.locator`
 * instead of `page.getBy*`.
 */

// @ts-check
import { test, forEachViewport } from './customerPpmTestFixture';

test.describe('Entire PPM closeout flow', () => {
  forEachViewport(async () => {
    test(`flows through happy path for existing shipment`, async ({ customerPpmPage }) => {
      test.slow();
      test.setTimeout(300000); // This one has been a headache forever. Shoehorn fix to go way above default "slow" timeout
      const move = await customerPpmPage.testHarness.buildApprovedMoveWithPPM();

      await customerPpmPage.signInForPPMWithMove(move);
      await customerPpmPage.navigateFromMMDashboardToMove(move);
      await customerPpmPage.navigateToAboutPage();
      await customerPpmPage.submitWeightTicketPage();
      await customerPpmPage.navigateFromCloseoutReviewPageToProGearPage();
      await customerPpmPage.submitProgearPage();
      await customerPpmPage.navigateFromCloseoutReviewPageToExpensesPage();
      await customerPpmPage.submitExpensePage();
      await customerPpmPage.navigateFromPPMReviewPageToFinalCloseoutPage();
      await customerPpmPage.submitFinalCloseout(move.id, {
        totalNetWeight: '2,000 lbs',
        proGearWeight: '2,000 lbs',
        expensesClaimed: '675.99',
        finalIncentiveAmount: '$92,921.12',
      });
    });

    test(`happy path with edits and backs`, async ({ customerPpmPage }) => {
      test.slow();
      test.setTimeout(300000); // This one has been a headache forever. Shoehorn fix to go way above default "slow" timeout
      const move = await customerPpmPage.testHarness.buildMoveWithPPMShipmentReadyForFinalCloseout();

      await customerPpmPage.signInForPPMWithMove(move);
      await customerPpmPage.navigateFromMMDashboardToMove(move);
      await customerPpmPage.navigateToPPMReviewPageWithCompletePPM();
      await customerPpmPage.navigateFromPPMReviewPageToFinalCloseoutPage();
      await customerPpmPage.returnToMoveHome();
      await customerPpmPage.navigateToPPMReviewPageWithCompletePPM();
      await customerPpmPage.navigateFromCloseoutReviewPageToAboutPage();
      await customerPpmPage.fillOutAboutPage();
      await customerPpmPage.navigateFromCloseoutReviewPageToEditWeightTicketPage();
      await customerPpmPage.submitWeightTicketPage();
      await customerPpmPage.navigateFromCloseoutReviewPageToEditProGearPage();
      await customerPpmPage.submitProgearPage({ belongsToSelf: false });
      await customerPpmPage.navigateFromCloseoutReviewPageToEditExpensePage();
      await customerPpmPage.submitExpensePage({ isEditExpense: true, amount: '833.41' });
      await customerPpmPage.navigateFromPPMReviewPageToFinalCloseoutPage();
      await customerPpmPage.submitFinalCloseout(move.id, {
        totalNetWeight: '2,000 lbs',
        proGearWeight: '500 lbs',
        expensesClaimed: '833.41',
        finalIncentiveAmount: '$92,921.12',
      });
    });

    test(`delete complete and incomplete line items`, async ({ customerPpmPage }) => {
      test.slow();
      const move = await customerPpmPage.testHarness.buildMoveWithPPMShipmentReadyForFinalCloseout();

      await customerPpmPage.signInForPPMWithMove(move);

      await customerPpmPage.navigateFromMMDashboardToMove(move);
      await customerPpmPage.navigateToPPMReviewPageWithCompletePPM();
      await customerPpmPage.verifySaveAndContinueEnabled();

      // Add incomplete weight ticket and exit
      await customerPpmPage.navigateFromCloseoutReviewPageToAddWeightTicketPage();
      await customerPpmPage.cancelAddLineItemAndReturnToCloseoutReviewPage();
      await customerPpmPage.returnToMoveHome();

      await customerPpmPage.navigateToPPMReviewPageWithCompletePPM();

      // Add incomplete moving expense and exit
      await customerPpmPage.navigateFromCloseoutReviewPageToAddExpensePage();
      await customerPpmPage.cancelAddLineItemAndReturnToCloseoutReviewPage();
      await customerPpmPage.returnToMoveHome();

      await customerPpmPage.navigateToPPMReviewPageWithCompletePPM();

      // Add incomplete pro-gear weight ticket and exit
      await customerPpmPage.navigateFromCloseoutReviewPageToAddProGearPage();
      await customerPpmPage.cancelAddLineItemAndReturnToCloseoutReviewPage();
      await customerPpmPage.returnToMoveHome();

      await customerPpmPage.navigateToPPMReviewPageWithCompletePPM();

      // Now that we have incomplete line items, we cannot submit the PPM
      await customerPpmPage.verifySaveAndContinueDisabled();

      // Delete incomplete line items
      await customerPpmPage.deleteWeightTicket(1, false);
      await customerPpmPage.deleteProGearExpense(1, false);
      await customerPpmPage.deleteMovingExpense(1, false);

      // Incomplete items have been deleted, so we should be allowed to submit now
      await customerPpmPage.verifySaveAndContinueEnabled();

      // Delete complete line items
      await customerPpmPage.deleteWeightTicket(0, true);
      await customerPpmPage.deleteProGearExpense(0, true);
      await customerPpmPage.deleteMovingExpense(0, true);

      // All the line items are gone, so we cannot submit the PPM
      await customerPpmPage.verifySaveAndContinueDisabled();
    });

    test(`deleting weight tickets updates final incentive`, async ({ customerPpmPage }) => {
      test.slow();
      const move = await customerPpmPage.testHarness.buildMoveWithPPMShipmentReadyForFinalCloseout();

      await customerPpmPage.signInForPPMWithMove(move);

      await customerPpmPage.navigateFromMMDashboardToMove(move);
      await customerPpmPage.navigateToPPMReviewPageWithCompletePPM();
      await customerPpmPage.navigateFromCloseoutReviewPageToAddWeightTicketPage();
      await customerPpmPage.submitWeightTicketPage();

      await customerPpmPage.navigateFromPPMReviewPageToFinalCloseoutPage();
      await customerPpmPage.verifyFinalIncentiveAndTotals({
        totalNetWeight: '6,000 lbs',
        proGearWeight: '1,500 lbs',
        expensesClaimed: '450.00',
        finalIncentiveAmount: '$172,795.18',
      });
      await customerPpmPage.page.getByRole('button', { name: 'Return to Homepage' }).click();
      await customerPpmPage.navigateToPPMReviewPageWithCompletePPM();

      await customerPpmPage.deleteWeightTicket(1, false);
      await customerPpmPage.navigateFromPPMReviewPageToFinalCloseoutPage();
      await customerPpmPage.submitFinalCloseout(move.id, {
        totalNetWeight: '4,000 lbs',
        proGearWeight: '1,500 lbs',
        expensesClaimed: '450.00',
        finalIncentiveAmount: '$69,498.74',
      });
    });
  });
});
