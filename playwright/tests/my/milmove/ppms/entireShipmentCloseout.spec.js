/**
 * Semi-automated converted from a cypress test, and thus may contain
 * non best-practices, in particular: heavy use of `page.locator`
 * instead of `page.getBy*`.
 */

// @ts-check
const { test, expect, forEachViewport } = require('./customerPpmTestFixture');

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
    test(`happy path with line item deletions`, async ({ customerPpmPage }) => {
      const move = await customerPpmPage.testHarness.buildMoveWithPPMShipmentReadyForFinalCloseout();

      await customerPpmPage.signInForPPMWithMove(move);
      await customerPpmPage.navigateToPPMReviewPage();
      // await customerPpmPage.deleteLineItem(customerPpmPage.page.getByText('Weight moved'), 'You are about to delete Trip 1. This cannot be undone.');
      // await customerPpmPage.deleteLineItem(customerPpmPage.page, 'You are about to delete Trip 1. This cannot be undone.');
      // await customerPpmPage.navigateFromCloseoutReviewPageToAboutPage();

      // First approach here is to click on each of the delete buttons separately by index
      // The indices shift so this is not reliable. Either need to not check the message and just do it 3 times, or have
      // a smarter selector
      const weightMoved = await customerPpmPage.page.getByRole('heading', { name: 'Weight moved' });
      await expect(weightMoved).toBeVisible();
      const foo = weightMoved.locator('../../..');
      await expect(foo).toBeVisible();
      await expect(foo.getByRole('button', { name: 'Delete' })).toBeVisible();
      await foo.getByRole('button', { name: 'Delete' }).click();
      await expect(customerPpmPage.page.getByText('You are about to delete')).toBeVisible();
      // await expect(customerPpmPage.page.getByText('You are about to delete Trip 1. This cannot be undone.')).toBeVisible();
      await customerPpmPage.page.getByRole('button', { name: 'Yes, Delete' }).click();
      await customerPpmPage.page.getByRole('button', { name: 'Delete' }).nth(2).click();
      // await customerPpmPage.page.getByTestId('modal').getByTestId('button').click();
      await customerPpmPage.page.getByText('You are about to delete Set 1. This cannot be undone.').click();
      await customerPpmPage.page.getByTestId('button').nth(1).click();
      await customerPpmPage.page.getByRole('button', { name: 'Delete' }).nth(3).click();
      await customerPpmPage.page.getByText('You are about to delete Receipt 1. This cannot be undone.').click();
    });
  });
});
