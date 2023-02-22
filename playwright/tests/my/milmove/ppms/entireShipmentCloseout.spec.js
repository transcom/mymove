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
      const move = await customerPpmPage.testHarness.buildApprovedMoveWithPPM();

      await customerPpmPage.signInForPPMWithMove(move);
      await customerPpmPage.navigateToAboutPage();
      await customerPpmPage.submitWeightTicketPage();
      await customerPpmPage.navigateFromCloseoutReviewPageToProGearPage();
      await customerPpmPage.submitProgearPage();
      await customerPpmPage.navigateFromCloseoutReviewPageToExpensesPage();
      await customerPpmPage.submitExpensePage();
      await customerPpmPage.navigateFromPPMReviewPageToFinalCloseoutPage();
      await customerPpmPage.submitFinalCloseout({
        totalNetWeight: '2,000 lbs',
        proGearWeight: '2,000 lbs',
        expensesClaimed: '675.99',
        finalIncentiveAmount: '$31,180.87',
      });
    });

    test(`happy path with edits and backs`, async ({ customerPpmPage }) => {
      const move = await customerPpmPage.testHarness.buildMoveWithPPMShipmentReadyForFinalCloseout();

      await customerPpmPage.signInForPPMWithMove(move);
      await customerPpmPage.navigateToPPMReviewPage();
      await customerPpmPage.navigateFromCloseoutReviewPageToAboutPage();
      await customerPpmPage.fillOutAboutPage();
      await customerPpmPage.navigateFromCloseoutReviewPageToEditWeightTicketPage();
      await customerPpmPage.submitWeightTicketPage();
      await customerPpmPage.navigateFromCloseoutReviewPageToEditProGearPage();
      await customerPpmPage.submitProgearPage({ belongsToSelf: false });
      await customerPpmPage.navigateFromCloseoutReviewPageToEditExpensePage();
      await customerPpmPage.submitExpensePage({ isEditExpense: true, amount: '833.41' });
      await customerPpmPage.navigateFromPPMReviewPageToFinalCloseoutPage();
      await customerPpmPage.submitFinalCloseout({
        totalNetWeight: '2,000 lbs',
        proGearWeight: '500 lbs',
        expensesClaimed: '833.41',
        finalIncentiveAmount: '$31,180.87',
      });
    });

    test(`delete complete and incomplete line items`, async ({ customerPpmPage }) => {
      const move = await customerPpmPage.testHarness.buildMoveWithPPMShipmentReadyForFinalCloseout();

      await customerPpmPage.signInForPPMWithMove(move);

      await customerPpmPage.navigateToPPMReviewPage();
      await customerPpmPage.verifySaveAndContinueEnabled();

      // Add incomplete weight ticket
      await customerPpmPage.navigateFromCloseoutReviewPageToAddWeightTicketPage();
      await customerPpmPage.cancelAddLineItemAndReturnToCloseoutReviewPage();

      // Add incomplete moving expense
      await customerPpmPage.navigateFromCloseoutReviewPageToAddExpensePage();
      await customerPpmPage.cancelAddLineItemAndReturnToCloseoutReviewPage();

      // Add incomplete pro-gear weight ticket
      await customerPpmPage.navigateFromCloseoutReviewPageToAddProGearPage();
      await customerPpmPage.cancelAddLineItemAndReturnToCloseoutReviewPage();

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
      const move = await customerPpmPage.testHarness.buildMoveWithPPMShipmentReadyForFinalCloseout();

      await customerPpmPage.signInForPPMWithMove(move);

      await customerPpmPage.navigateToPPMReviewPage();
      await customerPpmPage.navigateFromCloseoutReviewPageToAddWeightTicketPage();
      await customerPpmPage.submitWeightTicketPage();

      await customerPpmPage.navigateFromPPMReviewPageToFinalCloseoutPage();
      await customerPpmPage.verifyFinalIncentiveAndTotals({
        totalNetWeight: '6,000 lbs',
        proGearWeight: '1,500 lbs',
        expensesClaimed: '450.00',
        finalIncentiveAmount: '$128,268.35',
      });
      await customerPpmPage.page.getByRole('button', { name: 'Return to Homepage' }).click();
      await customerPpmPage.navigateToPPMReviewPage();

      await customerPpmPage.deleteWeightTicket(1, false);
      await customerPpmPage.navigateFromPPMReviewPageToFinalCloseoutPage();
      await customerPpmPage.submitFinalCloseout({
        totalNetWeight: '4,000 lbs',
        proGearWeight: '1,500 lbs',
        expensesClaimed: '450.00',
        finalIncentiveAmount: '$48,527.43',
      });
    });
  });
});
