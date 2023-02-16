// @ts-check
import { test, expect } from './servicesCounselingTestFixture';

test('A service counselor can approve/reject weight tickets', async ({ page, scPage }) => {
  // Create a move with TestHarness, and then navigate to the move details page for it
  const move = await scPage.testHarness.buildApprovedMoveWithPPMWeightTicketOffice();
  await scPage.navigateToCloseoutMove(move.locator);

  // Navigate to the "Review documents" page
  await page.getByRole('button', { name: 'Review documents' }).click();
  await scPage.waitForPage.reviewDocuments();

  // Click "Accept" on the weight ticket, then save (and navigate back to the move details page)
  await page.getByText('Accept').click();
  await page.getByRole('button', { name: 'Continue' }).click();
  await scPage.waitForPage.moveDetails();

  // Return to the weight ticket and verify that it's approved
  await page.getByRole('button', { name: 'Review documents' }).click();
  await scPage.waitForPage.reviewDocuments();
  await expect(page.getByRole('radio', { name: 'Accept' })).toBeChecked();

  // Click "Reject" on the weight ticket, provide a reason, then save
  await page.getByText('Reject').click();
  await page.getByLabel('Reason').fill('Justification for rejection');
  await page.getByRole('button', { name: 'Continue' }).click();
  await scPage.waitForPage.moveDetails();

  // Return to the weight ticket and verify that it's been edited
  await page.getByRole('button', { name: 'Review documents' }).click();
  await scPage.waitForPage.reviewDocuments();
  await expect(page.getByRole('radio', { name: 'Reject' })).toBeChecked();
  await expect(page.getByLabel('Reason')).toHaveValue('Justification for rejection');
});
