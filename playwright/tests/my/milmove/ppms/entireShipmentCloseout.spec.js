/**
 * Semi-automated converted from a cypress test, and thus may contain
 * non best-practices, in particular: heavy use of `page.locator`
 * instead of `page.getBy*`.
 */

// @ts-check
const { test, forEachViewport, CustomerPpmPage } = require('./customerPpmTestFixture');

test.describe('Entire PPM closeout flow', () => {
  forEachViewport(async () => {
    test(`flows through happy path for existing shipment`, async ({ customerPage }) => {
      const move = await customerPage.testHarness.buildApprovedMoveWithPPM();
      const customerPpmPage = new CustomerPpmPage(customerPage, move);

      await customerPpmPage.signInAndNavigateToAboutPage();
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

    test(`happy path with edits and backs`, async ({ customerPage }) => {
      const move = await customerPage.testHarness.buildMoveWithPPMShipmentReadyForFinalCloseout();
      const customerPpmPage = new CustomerPpmPage(customerPage, move);

      await customerPpmPage.signInAndNavigateToPPMReviewPage();
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
  });
});
