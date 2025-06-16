// @ts-check
import { test, expect } from './servicesCounselingTestFixture';

const completePPMCloseoutForCustomerEnabled = process.env.FEATURE_FLAG_COMPLETE_PPM_CLOSEOUT_FOR_CUSTOMER;

test('A service counselor can approve/reject weight tickets', async ({ page, scPage }) => {
  // Create a move with TestHarness, and then navigate to the move details page for it
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
  const move = await scPage.testHarness.buildApprovedMoveWithPPMWeightTicketOffice();
  await scPage.navigateToCloseoutMove(move.locator);

  // Navigate to the "Review documents" page
  await page.getByRole('button', { name: 'Review documents' }).click();
  await scPage.waitForPage.reviewWeightTicket();

  // Click "Accept" on the weight ticket, then proceed
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

test('A services counselor can reduce PPM weights for a move with excess weight', async ({ page, scPage }) => {
  const move = await scPage.testHarness.buildApprovedMoveWithPPMShipmentAndExcessWeight();
  await scPage.navigateToCloseoutMove(move.locator);

  // navigate to review-shipment-weights page and verify page components are rendered
  await page.getByRole('button', { name: 'Review shipment weights' }).click();
  await scPage.waitForPage.reviewShipmentWeights();

  // verify that the excess weight alert is visible, since the move has excess weight
  await expect(
    page.getByText('This move has excess weight. Review PPM weight ticket documents to resolve.'),
  ).toBeVisible();

  // navigate to review-documents page and decrease ppm shipment weight below threshold
  await page.getByRole('link', { name: 'Review Documents' }).click();
  await scPage.waitForPage.reviewWeightTicket();

  await page.getByTestId('fullWeight').clear();
  await page.getByTestId('fullWeight').fill('8000');
  await page.getByTestId('fullWeight').blur();

  await page.getByText('Accept').click();
  await page.getByRole('button', { name: 'Continue' }).click();
  await page.getByTestId('closeSidebar').click();

  // navigate to review-shipment-weights page and verify excess weight alert is no longer visible
  await page.getByRole('button', { name: 'Review shipment weights' }).click();
  await scPage.waitForPage.reviewShipmentWeights();

  await expect(
    page.getByText('This move has excess weight. Review PPM weight ticket documents to resolve.'),
  ).toHaveCount(0);
});

test('A services counselor can edit allowable weight', async ({ page, scPage }) => {
  const move = await scPage.testHarness.buildApprovedMoveWithPPMShipmentAndExcessWeight();
  await scPage.navigateToCloseoutMove(move.locator);

  await page.getByRole('button', { name: 'Review shipment weights' }).click();
  await scPage.waitForPage.reviewShipmentWeights();

  await page.getByRole('link', { name: 'Review Documents' }).click();
  await scPage.waitForPage.reviewWeightTicket();

  await page.getByTestId('editAllowableWeightButton').click();
  await page.getByText('Cancel').click();

  await page.getByTestId('editAllowableWeightButton').click();
  await page.getByTestId('editAllowableWeightInput').focus();
  await page.getByTestId('editAllowableWeightInput').fill('8000');
  await page.getByTestId('editAllowableWeightInput').blur();
  await page.getByText('Save').click();
  await expect(page.getByText('8,000 lbs')).toBeVisible();

  await page.getByTestId('closeSidebar').click();

  // Ensure change appears in audit history
  await page.getByText('Move History').click();
  await expect(page.getByText('Allowable Weight: 8,000 lbs')).toBeVisible();
});

test('A service counselor cannot upload an other excel file except a weight estimator file into a weight ticket file upload', async ({
  page,
  scPage,
}) => {
  test.skip(completePPMCloseoutForCustomerEnabled === 'false', 'Skip if FF is disabled.');
  const move = await scPage.testHarness.buildApprovedMoveWithPPMWithAboutFormComplete();
  await scPage.signInAsNewServicesCounselorUser();
  await scPage.navigateToMoveUsingMoveSearch(move.locator);

  // click on "Complete PPM on behalf of the customer" button
  await page.getByRole('button', { name: 'Complete PPM on behalf of the Customer' }).click();

  // expect the new page to have heading "Review"
  await expect(page.getByRole('heading', { name: 'Review' })).toBeVisible();

  // click the edit button in the "Trip 1" card of "Weight Moved" section
  await expect(page.getByText('Edit')).toBeVisible();
  await page.getByText('Add More Weight').click();

  // expect the weight tickets page to have loaded
  await expect(page.getByRole('heading', { name: 'Weight Tickets' })).toBeVisible();

  // can use fillOutWeightTicketPage as base
  await scPage.fillOutWeightTicketWithIncorrectXlsx();
});

test('A service counselor can see HHG weights when reviewing weight tickets', async ({ page, scPage }) => {
  // Create a move with TestHarness, and then navigate to the move details page for it
  const move = await scPage.testHarness.buildApprovedMoveWithPPMWeightTicketOfficeWithHHG();
  await scPage.navigateToCloseoutMove(move.locator);

  // Navigate to the "Review documents" page
  await page.getByRole('button', { name: 'Review documents' }).click();
  await scPage.waitForPage.reviewWeightTicket();

  // Verify the heading "HHG 1" is present
  await expect(page.getByRole('heading', { level: 3, name: /HHG 1/ })).toBeVisible();

  // Verify that the labels for Estimated Weight and Actual Weight are present
  await expect(page.getByText('Estimated Weight')).toBeVisible();
  await expect(page.getByText('Actual Weight')).toBeVisible();

  // Verify that the specific weights are displayed as expected
  await expect(page.getByText('1,400 lbs')).toBeVisible();
  await expect(page.getByText('2,000 lbs')).toBeVisible();
});
