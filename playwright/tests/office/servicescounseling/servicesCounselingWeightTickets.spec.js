// @ts-check
import { test, expect } from './servicesCounselingTestFixture';

test('A service counselor can approve/reject weight tickets', async ({ page, scPage }) => {
  // Create a move with TestHarness, and then navigate to the move details page for it
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
  await page.getByRole('button', { name: 'Confirm' }).click();
  await scPage.waitForPage.moveDetails();

  // NOTE: Code below is commented out because the feature for the SC to be able to review documents AFTER it has been submitted will be picked up at a future date.
  // Currently SC is unable to re-review documents after it has been submitted, so these tests were failing.

  // Return to the weight ticket and verify that it's approved
  // await page.getByRole('button', { name: 'Review documents' }).click();
  // await scPage.waitForPage.reviewWeightTicket();
  // await expect(page.getByRole('radio', { name: 'Accept' })).toBeChecked();

  // // Click "Reject" on the weight ticket, provide a reason, then save
  // await page.getByText('Reject').click();
  // await page.getByLabel('Reason').fill('Justification for rejection');
  // await page.getByRole('button', { name: 'Continue' }).click();
  // await scPage.waitForPage.reviewDocumentsConfirmation();
  // await page.getByRole('button', { name: 'Confirm' }).click();
  // await scPage.waitForPage.moveDetails();

  // // Return to the weight ticket and verify that it's been edited
  // await page.getByRole('button', { name: 'Review documents' }).click();
  // await scPage.waitForPage.reviewWeightTicket();
  // await expect(page.getByRole('radio', { name: 'Reject' })).toBeChecked();
  // await expect(page.getByLabel('Reason')).toHaveValue('Justification for rejection');
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
