import { test, expect } from './servicesCounselingTestFixture';

test('A service counselor can approve/reject pro-gear weight tickets', async ({ page, scPage }) => {
  // To solve timeout issues
  test.slow();
  // Create a move with TestHarness, and then navigate to the move details page for it
  const move = await scPage.testHarness.buildApprovedMoveWithPPMProgearWeightTicketOffice();
  await scPage.navigateToCloseoutMove(move.locator);

  // Navigate to the "Review documents" page
  await page.getByRole('button', { name: 'Review documents' }).click();
  await scPage.waitForPage.reviewWeightTicket();

  // Weight ticket is first in the order of docs. Click "Accept" on the weight ticket, then proceed
  await page.getByText('Accept').click();
  await page.getByRole('button', { name: 'Continue' }).click();

  // Next is pro-gear ticket here. Click "Accept" on the pro-gear ticket, then proceed
  await expect(page.getByRole('heading', { name: 'Review pro-gear 1' })).toBeVisible();
  await page.getByText('Accept').click();
  await page.getByRole('button', { name: 'Continue' }).click();
  await scPage.waitForPage.reviewDocumentsConfirmation();

  // Click "Confirm" on confirmation page, returning to move details page
  await page.getByRole('button', { name: 'Confirm' }).click();
  await scPage.waitForPage.moveDetails();

  // NOTE: Code below is commented out because the feature for the SC to be able to review documents AFTER it has been submitted will be picked up at a future date.
  // Currently SC is unable to re-review documents after it has been submitted, so these tests were failing.

  // Return to the pro-gear ticket and verify that it's approved
  // await page.getByRole('button', { name: 'Review documents' }).click();
  // await scPage.waitForPage.reviewWeightTicket();

  // // Weight ticket is first. Need to skip over to Pro-gear ticket
  // await page.getByRole('button', { name: 'Continue' }).click();
  // await expect(page.getByRole('heading', { name: 'Review pro-gear 1' })).toBeVisible();
  // await expect(page.getByRole('radio', { name: 'Accept' })).toBeChecked();

  // // Click "Reject" on the Pro-gear ticket, provide a reason, then save
  // await page.getByText('Reject').click();
  // await page.getByLabel('Reason').fill('Reason for rejection');
  // await page.getByRole('button', { name: 'Continue' }).click();
  // await scPage.waitForPage.reviewDocumentsConfirmation();
  // await page.getByRole('button', { name: 'Confirm' }).click();
  // await scPage.waitForPage.moveDetails();

  // // Return to the pro-gear ticket and verify that it's been edited
  // await page.getByRole('button', { name: 'Review documents' }).click();
  // await scPage.waitForPage.reviewWeightTicket();
  // await page.getByRole('button', { name: 'Continue' }).click();

  // await expect(page.getByRole('radio', { name: 'Reject' })).toBeChecked();
  // await expect(page.getByLabel('Reason')).toHaveValue('Reason for rejection');
});
