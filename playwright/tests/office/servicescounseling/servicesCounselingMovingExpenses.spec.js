import { test, expect } from './servicesCounselingTestFixture';

test('A service counselor can approve/reject moving expenses', async ({ page, scPage }) => {
  // Create a move with TestHarness, and then navigate to the move details page for it
  const move = await scPage.testHarness.buildApprovedMoveWithPPMMovingExpense();
  await scPage.navigateToCloseoutMove(move.locator);

  // Navigate to the "Review documents" page
  await page.getByRole('button', { name: 'Review documents' }).click();
  await scPage.waitForPage.reviewDocuments();

  // Weight ticket is first in the order of docs. Click "Accept" on the weight ticket, then proceed
  await page.getByText('Accept').click();
  await page.getByRole('button', { name: 'Continue' }).click();

  // Next is expense ticket here. Click "Accept" on the  expense, then proceed
  await expect(page.getByRole('heading', { name: 'Review receipt 1' })).toBeVisible();
  await page.getByText('Accept').click();
  await page.getByRole('button', { name: 'Continue' }).click();

  // Next is storage expense ticket. Click "Accept", then proceed
  await expect(page.getByRole('heading', { name: 'Review storage 2' })).toBeVisible();
  await page.getByText('Accept').click();
  await page.getByRole('button', { name: 'Continue' }).click();
  await scPage.waitForPage.reviewDocumentsConfirmation();

  // Click "Confirm" on confirmation page, returning to move details page
  await page.getByRole('button', { name: 'Confirm' }).click();
  await scPage.waitForPage.moveDetails();

  // Return to the expenses ticket and verify that it's approved
  await page.getByRole('button', { name: 'Review documents' }).click();
  await scPage.waitForPage.reviewDocuments();

  // Weight ticket is first. Need to skip over to expense ticket
  await page.getByRole('button', { name: 'Continue' }).click();

  await expect(page.getByRole('radio', { name: 'Accept' })).toBeChecked();

  // Click "Reject" on the expense ticket, provide a reason, then save
  await page.getByText('Reject').click();
  await page.getByLabel('Reason').fill('Reason for expense rejection');
  await page.getByRole('button', { name: 'Continue' }).click();

  // Verify storage ticket is approved. Then, edit to "Reject" on the storage ticket, provide a reason, then save
  await expect(page.getByRole('heading', { name: 'Review storage 2' })).toBeVisible();
  await expect(page.getByRole('radio', { name: 'Accept' })).toBeChecked();
  await page.getByText('Reject').click();
  await page.getByLabel('Reason').fill('Reason for storage rejection');
  await page.getByRole('button', { name: 'Continue' }).click();

  await scPage.waitForPage.reviewDocumentsConfirmation();
  await page.getByRole('button', { name: 'Confirm' }).click();
  await scPage.waitForPage.moveDetails();

  // Return to the expense and verify that it's been rejected
  await page.getByRole('button', { name: 'Review documents' }).click();
  await scPage.waitForPage.reviewDocuments();
  await page.getByRole('button', { name: 'Continue' }).click();

  await expect(page.getByRole('heading', { name: 'Review receipt 1' })).toBeVisible();
  await expect(page.getByRole('radio', { name: 'Reject' })).toBeChecked();
  await expect(page.getByLabel('Reason')).toHaveValue('Reason for expense rejection');

  // Edit expense ticket to "Exclude", provide a reason, then save
  await page.getByText('Exclude').click();
  await page.getByLabel('Reason').fill('Reason for excluding expense');
  await page.getByRole('button', { name: 'Continue' }).click();

  // Verify storage ticket, Edit to "Exclude", provide a reason, then save
  await page.getByText('Exclude').click();
  await page.getByLabel('Reason').fill('Reason for excluding storage');
  await page.getByRole('button', { name: 'Continue' }).click();

  await scPage.waitForPage.reviewDocumentsConfirmation();
  await page.getByRole('button', { name: 'Confirm' }).click();
  await scPage.waitForPage.moveDetails();

  // Return to the expense and verify that it's been excluded
  await page.getByRole('button', { name: 'Review documents' }).click();
  await scPage.waitForPage.reviewDocuments();
  await page.getByRole('button', { name: 'Continue' }).click();

  await expect(page.getByRole('heading', { name: 'Review receipt 1' })).toBeVisible();
  await expect(page.getByRole('radio', { name: 'Exclude' })).toBeChecked();
  await expect(page.getByLabel('Reason')).toHaveValue('Reason for excluding expense');

  // Return to storage and verify that it's been excluded
  await page.getByRole('button', { name: 'Continue' }).click();
  await expect(page.getByRole('radio', { name: 'Exclude' })).toBeChecked();
  await expect(page.getByLabel('Reason')).toHaveValue('Reason for excluding storage');
});
