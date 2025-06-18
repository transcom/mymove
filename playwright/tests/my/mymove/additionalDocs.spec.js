import { test, expect } from '../../utils/my/customerTest';

test.describe('Additional Documents', () => {
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
    await expect(page.getByRole('link', { name: /trustedAgent-\d{14}.+/ })).toBeVisible();
  });
});
