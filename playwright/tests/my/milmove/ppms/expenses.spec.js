/**
 * Semi-automated converted from a cypress test, and thus may contain
 * non best-practices, in particular: heavy use of `page.locator`
 * instead of `page.getBy*`.
 */

// @ts-check
const { expect, test, forEachViewport } = require('./customerPpmTestFixture');

test.describe('Expenses', () => {
  forEachViewport(async () => {
    test.beforeEach(async ({ customerPpmPage }) => {
      const move = await customerPpmPage.testHarness.buildApprovedMoveWithPPMMovingExpense();
      await customerPpmPage.signInForPPMWithMove(move);
      await customerPpmPage.navigateToPPMReviewPage();
    });

    test(`new expense page loads`, async ({ customerPpmPage }) => {
      await customerPpmPage.navigateFromCloseoutReviewPageToExpensesPage();
      await customerPpmPage.submitExpensePage();
    });

    test(`edit expense page loads`, async ({ page }) => {
      // edit the first expense receipt
      const receipt1 = page.getByText('Receipt 1', { exact: true });
      await expect(receipt1).toBeVisible();
      await receipt1.locator('../..').getByRole('link', { name: 'Edit' }).click();
      await expect(page).toHaveURL(/\/moves\/[^/]+\/shipments\/[^/]+\/expenses\/[^/]+/);

      const expenseType = page.locator('select[name="expenseType"]');
      await expect(expenseType).toHaveValue('PACKING_MATERIALS');
      await expenseType.selectOption({ label: 'Tolls' });

      const expenseDescription = page.locator('input[name="description"]');
      await expect(expenseDescription).toHaveValue('Packing Peanuts');
      await expenseDescription.clear();
      await expenseDescription.type('PA Turnpike EZ-Pass');

      await expect(page.locator('label[for="yes-used-gtcc"]')).toBeChecked();
      await page.locator('label[for="no-did-not-use-gtcc"]').click();

      const expenseAmount = page.locator('input[name="amount"]');
      await expect(expenseAmount).toHaveValue('23.45');
      await expenseAmount.clear();
      await expenseAmount.type('54.32');

      const missingReceipt = page.locator('label[for="missingReceipt"]');
      await expect(missingReceipt).not.toBeChecked();
      await missingReceipt.click();

      await page.getByRole('button', { name: 'Save & Continue' }).click();
      await expect(page).toHaveURL(/\/moves\/[^/]+\/shipments\/[^/]+\/review/);

      await expect(page.getByText('PA Turnpike EZ-Pass')).toBeVisible();
    });
  });
});
