import { test, expect } from './servicesCounselingTestFixture';

test('A service counselor can approve/reject pro-gear weight tickets', async ({ page, scPage }) => {
  // To solve timeout issues
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
  await page.getByRole('button', { name: 'Preview PPM Payment Packet' }).click();
  await expect(page.getByTestId('loading-spinner')).not.toBeVisible();
  await page.getByRole('button', { name: 'Complete PPM Review' }).click();
  await page.getByRole('button', { name: 'Yes' }).click();
  await scPage.waitForPage.moveDetails();
});
