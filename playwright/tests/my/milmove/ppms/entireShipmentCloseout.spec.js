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

      const weightMoved = await customerPpmPage.page.getByRole('heading', { name: 'Weight moved' });
      await expect(weightMoved).toBeVisible();
      const weightTicketsContainer = weightMoved.locator('../../..');
      await expect(weightTicketsContainer).toBeVisible();
      await expect(weightTicketsContainer.getByRole('button', { name: 'Delete' })).toBeVisible();
      await weightTicketsContainer.getByRole('button', { name: 'Delete' }).click();
      await expect(customerPpmPage.page.getByText('You are about to delete Trip 1')).toBeVisible();
      await customerPpmPage.page.getByRole('button', { name: 'Yes, Delete' }).click();
      await expect(
        customerPpmPage.page.getByText('No weight moved documented. At least one trip is required to continue.'),
      ).toBeVisible();
      await expect(
        customerPpmPage.page.getByText(
          'There are items below that are missing required information. Please select “Edit” to enter all required information or “Delete” to remove the item.',
        ),
      ).toBeVisible();
      await expect(customerPpmPage.page.getByText('Trip 1 successfully deleted.')).toBeVisible();

      const proGearExpense = await customerPpmPage.page.getByRole('heading', { name: 'Pro-gear' });
      await expect(proGearExpense).toBeVisible();
      const proGearExpenseContainer = proGearExpense.locator('../../..');
      await expect(proGearExpenseContainer).toBeVisible();
      await expect(proGearExpenseContainer.getByRole('button', { name: 'Delete' })).toBeVisible();
      await proGearExpenseContainer.getByRole('button', { name: 'Delete' }).click();
      await expect(customerPpmPage.page.getByText('You are about to delete Set 1')).toBeVisible();
      await customerPpmPage.page.getByRole('button', { name: 'Yes, Delete' }).click();
      await expect(customerPpmPage.page.getByText('No pro-gear weight documented.')).toBeVisible();
      await expect(customerPpmPage.page.getByText('Set 1 successfully deleted.')).toBeVisible();

      const moveExpense = await customerPpmPage.page.getByRole('heading', { name: 'Expenses' });
      await expect(moveExpense).toBeVisible();
      const moveExpensesContainer = moveExpense.locator('../../..');
      await expect(moveExpensesContainer).toBeVisible();
      await expect(moveExpensesContainer.getByRole('button', { name: 'Delete' })).toBeVisible();
      await moveExpensesContainer.getByRole('button', { name: 'Delete' }).click();
      await expect(customerPpmPage.page.getByText('You are about to delete Receipt 1')).toBeVisible();
      await customerPpmPage.page.getByRole('button', { name: 'Yes, Delete' }).click();
      await expect(customerPpmPage.page.getByText('No receipts uploaded.')).toBeVisible();
      await expect(customerPpmPage.page.getByText('Receipt 1 successfully deleted.')).toBeVisible();
    });
  });
});
