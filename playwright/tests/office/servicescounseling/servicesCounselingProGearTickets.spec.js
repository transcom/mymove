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

test('The only Xlsx file that a service counselor can only upload is a weight estimator file template.', async ({
  page,
  scPage,
}) => {
  const move = await scPage.testHarness.buildApprovedMoveWithPPMWithAboutFormComplete();
  await scPage.signInAsNewServicesCounselorUser();
  await scPage.navigateToMoveUsingMoveSearch(move.locator);

  // click on "Complete PPM on behalf of customer"
  await page.getByRole('button', { name: 'Complete PPM on behalf of the Customer' }).click();

  // expect the new page to have heading "Review"
  await expect(page.getByRole('heading', { name: 'Review' })).toBeVisible();

  // click on Set 1 of Pro-gear section

  await expect(page.getByRole('heading', { name: 'Pro-gear' })).toBeVisible();

  await expect(page.getByText('Add Pro-gear Weight')).toBeVisible();
  await page.getByText('Add Pro-gear Weight').click();

  await expect(page.getByRole('heading', { name: 'Pro-gear' })).toBeVisible();
  const progearTypeSelector = `label[for="ownerOfProGearSelf"]`;
  await page.locator(progearTypeSelector).click();
  await expect(page.getByRole('heading', { name: 'Description' })).toBeVisible();

  await scPage.fillOutProGearWithIncorrectXlsx();
});

test('A service counselor can see the total for a progear weight ticket regular and spouse after update', async ({
  page,
  scPage,
}) => {
  // To solve timeout issues
  test.slow();
  // Create a move with TestHarness, and then navigate to the move details page for it
  const move = await scPage.testHarness.buildApprovedMoveWithPPMWithMultipleProgearWeightTicketsOffice();
  const reviewShipmentPage = scPage.waitForPage.reviewShipmentWeights();
  await scPage.navigateToCloseoutMove(move.locator);

  // Navigate to the "Review documents" page
  await page.getByRole('button', { name: 'Review shipment weights' }).click();
  await reviewShipmentPage;

  await expect(page.locator('[data-testid="proGear-0"]')).toHaveText('300 lbs');
  await expect(page.locator('[data-testid="spouseProGear-0"]')).toHaveText('75 lbs');
});

test('A service counselor can see the total for a progear weight ticket regular and spouse after delete', async ({
  page,
  scPage,
}) => {
  // To solve timeout issues
  test.slow();
  // Create a move with TestHarness, and then navigate to the move details page for it
  const move = await scPage.testHarness.buildApprovedMoveWithPPMWithMultipleProgearWeightTicketsOffice2();
  const reviewShipmentPage = scPage.waitForPage.reviewShipmentWeights();
  await scPage.navigateToCloseoutMove(move.locator);

  // Navigate to the "Review documents" page
  await page.getByRole('button', { name: 'Review shipment weights' }).click();
  await reviewShipmentPage;

  await expect(page.locator('[data-testid="proGear-0"]')).toHaveText('300 lbs');
  await expect(page.locator('[data-testid="spouseProGear-0"]')).toHaveText('50 lbs');
});
