import { test, expect } from '../../utils/my/customerTest';

const multiMoveEnabled = process.env.FEATURE_FLAG_MULTI_MOVE;
const manageSupportDocsEnabled = process.env.FEATURE_FLAG_MANAGE_SUPPORTING_DOCS;

test.describe('Additional Documents', () => {
  test.skip(multiMoveEnabled === 'false', 'Skip if MultiMove workflow is not enabled.');
  test.skip(manageSupportDocsEnabled === 'false', 'Skip if manage supporting docs workflow is not enabled.');

  test('Users can download documents uploaded to Additional Documents', async ({ page, customerPage }) => {
    // Generate a move that has the status of SUBMITTED
    const move = await customerPage.testHarness.buildSubmittedMoveWithPPMShipmentForSC();
    const userId = move?.Orders?.service_member?.user_id;

    // Sign-in and navigate to move home page
    await customerPage.signInAsExistingCustomer(userId);
    await customerPage.navigateFromMMDashboardToMove(move);
    await customerPage.waitForPage.home();

    // Go to the Upload Additional Documents page
    await page.getByRole('button', { name: 'Upload/Manage Additional Documentation' }).click();

    // Upload document
    const filepondContainer = page.locator('.filepond--wrapper');
    await customerPage.uploadFileViaFilepond(filepondContainer, 'trustedAgent.pdf');

    // Verify filename is a downloadable link
    await expect(page.getByRole('link', { name: 'trustedAgent.pdf' })).toBeVisible();
  });
});
