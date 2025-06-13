import { test, expect } from './servicesCounselingTestFixture';

test('A service counselor can approve/reject moving expenses', async ({ page, scPage }) => {
  test.slow();
  await page.route('**/ghc/v1/ppm-shipments/*/payment-packet', async (route) => {
    // mocked blob
    const fakePdfBlob = new Blob(['%PDF-1.4 foo'], { type: 'application/pdf' });
    const arrayBuffer = await fakePdfBlob.arrayBuffer();
    // playwright route mocks only want a Buffer or string
    const bodyBuffer = Buffer.from(arrayBuffer);
    await route.fulfill({
      status: 200,
      contentType: 'application/pdf',
      body: bodyBuffer,
    });
  });
  // Create a move with TestHarness, and then navigate to the move details page for it
  const move = await scPage.testHarness.buildApprovedMoveWithPPMMovingExpenseOffice();
  await scPage.navigateToCloseoutMove(move.locator);

  // Navigate to the "Review documents" page
  await page.getByRole('button', { name: /Review documents/i }).click();
  await scPage.waitForPage.reviewWeightTicket();
  // Weight ticket is first in the order of docs. Click "Accept" on the weight ticket, then proceed
  await page.getByText('Accept').click();
  await page.getByRole('button', { name: 'Continue' }).click();

  // Next is packing materials expense ticket here. Click "Accept" on the expense, then proceed
  await scPage.waitForPage.reviewExpenseTicket('Packing Materials', 1, 1);
  await page.getByText('Accept').click();
  await page.getByRole('button', { name: 'Continue' }).click();

  // Next is storage expense ticket. Click "Accept", then proceed
  await scPage.waitForPage.reviewExpenseTicket('Storage', 2, 1);
  await page.getByText('Accept').click();
  await page.getByRole('button', { name: 'Continue' }).click();
  await scPage.waitForPage.reviewDocumentsConfirmation();

  // Click "Confirm" on confirmation page, returning to move details page
  await page.getByRole('button', { name: 'Preview PPM Payment Packet' }).click();
  await expect(page.getByTestId('loading-spinner')).not.toBeVisible();
  await page.getByRole('button', { name: 'Complete PPM Review' }).click();
  await page.getByRole('button', { name: 'Yes' }).click();
  await scPage.waitForPage.moveDetails();
});

test('Review documents page displays correct value for Total days in SIT', async ({ page, scPage }) => {
  test.slow();
  test.setTimeout(300000); // This one has been a headache forever. Shoehorn fix to go way above default "slow" timeout
  // Create a move with TestHarness, and then navigate to the move details page for it
  const move = await scPage.testHarness.buildApprovedMoveWithPPMMovingExpenseOffice();
  await scPage.navigateToCloseoutMove(move.locator);

  // Navigate to the "Review documents" page
  await page.getByRole('button', { name: /Review documents/i }).click();

  await scPage.waitForPage.reviewWeightTicket();
  await page.getByText('Accept').click();
  await page.getByRole('button', { name: 'Continue' }).click();

  await scPage.waitForPage.reviewExpenseTicket('Packing Materials', 1, 1);
  await page.getByText('Accept').click();
  await page.getByRole('button', { name: 'Continue' }).click();

  // Storage expense ticket.
  await scPage.waitForPage.reviewExpenseTicket('Storage', 2, 1);

  // The SIT Days should be counting the "Start date" and "End date" (i.e. Start date of 15 Apr 24 and End date of 19 Apr 24 is 5 days of Storage.
  await page.locator('[name="sitStartDate"]').clear();
  await page.locator('[name="sitStartDate"]').fill('15 Apr 2024');
  await page.locator('[name="sitEndDate"]').clear();
  await page.locator('[name="sitEndDate"]').fill('19 Apr 2024');
  await page.locator('[name="sitEndDate"]').press('Tab'); // Exit out of datepicker view

  await expect(page.locator('[data-testid="days-in-sit"]')).toContainText('5');

  // Final review page in Review documents
  await page.getByText('Accept').click();
  await page.getByRole('button', { name: 'Continue' }).click();
  await scPage.waitForPage.reviewDocumentsConfirmation();

  await expect(page.locator('[data-testid="days-in-sit"]')).toContainText('5');
});
