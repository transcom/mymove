import { test, expect } from './servicesCounselingTestFixture';

test('A service counselor can edit and approve gun safe weight tickets', async ({ page, scPage }) => {
  // To solve timeout issues
  test.slow();
  // Create a move with TestHarness, and then navigate to the move details page for it
  const move = await scPage.testHarness.buildApprovedMoveWithPPMGunSafeWeightTicketOffice();
  await scPage.navigateToCloseoutMove(move.locator);

  // Navigate to the "Review documents" page
  await page.getByRole('button', { name: 'Review documents' }).click();
  await scPage.waitForPage.reviewWeightTicket();

  // Weight ticket
  await page.getByText('Accept').click();
  await page.getByRole('button', { name: 'Continue' }).click();

  // Gun safe weight ticket
  await expect(page.getByRole('heading', { name: 'Review gun safe 1' })).toBeVisible();
  await page.getByTestId('gunSafeWeightTicket').isChecked();
  await page.getByText('Constructed weight').click();
  await page.getByTestId('gunSafeWeight').clear();
  await page.getByTestId('gunSafeWeight').fill('300');
  await page.getByText('Accept').click();
  await page.getByRole('button', { name: 'Continue' }).click();
  await scPage.waitForPage.reviewDocumentsConfirmation();
  await expect(page.getByText('Gun Safe 1')).toBeVisible();
  await expect(page.getByTestId('gunSafeStatus')).toContainText('Accept');
  await expect(page.getByText('Missing Weight Ticket (Constructed)')).toBeVisible();
  await expect(page.getByTestId('gunSafeHasWeightTickets')).toContainText('Yes');
  await expect(page.getByText('Gun Safe Weight')).toBeVisible();
  await expect(page.getByTestId('gunSafeWeight')).toContainText('300');

  // Click "Confirm" on confirmation page, returning to move details page
  await page.getByRole('button', { name: 'PPM Review Complete' }).click();
  await scPage.waitForPage.moveDetails();
});

test('A service counselor can reject gun safe weight tickets', async ({ page, scPage }) => {
  // To solve timeout issues
  test.slow();
  // Create a move with TestHarness, and then navigate to the move details page for it
  const move = await scPage.testHarness.buildApprovedMoveWithPPMGunSafeWeightTicketOffice();
  await scPage.navigateToCloseoutMove(move.locator);

  // Navigate to the "Review documents" page
  await page.getByRole('button', { name: 'Review documents' }).click();
  await scPage.waitForPage.reviewWeightTicket();

  // Weight ticket
  await page.getByText('Accept').click();
  await page.getByRole('button', { name: 'Continue' }).click();

  // Gun safe weight ticket
  await expect(page.getByRole('heading', { name: 'Review gun safe 1' })).toBeVisible();
  await page.getByText('Reject').click();
  await page.getByLabel('Reason').fill('Reason for rejection');
  await page.getByRole('button', { name: 'Continue' }).click();
  await scPage.waitForPage.reviewDocumentsConfirmation();

  await expect(page.getByText('Gun Safe 1')).toBeVisible();
  await expect(page.getByTestId('gunSafeStatus')).toContainText('Reject');

  // Click "Confirm" on confirmation page, returning to move details page
  await page.getByRole('button', { name: 'PPM Review Complete' }).click();
  await scPage.waitForPage.moveDetails();
});
