import { test, expect } from '../../utils/my/customerTest';

const multiMoveEnabled = process.env.FEATURE_FLAG_MULTI_MOVE;

test.describe('Uploads', () => {
  test.skip(multiMoveEnabled === 'true', 'Skip if MultiMove workflow is enabled.');
  test('Users can upload but cannot delete orders once move has been submitted', async ({ page, customerPage }) => {
    // Generate a move that has the status of SUBMITTED
    const move = await customerPage.testHarness.buildSubmittedMoveWithPPMShipmentForSC();
    const userId = move?.Orders?.service_member?.user_id;
    await customerPage.signInAsExistingCustomer(userId);
    await customerPage.waitForPage.home();
    await page.getByRole('button', { name: 'Review your request' }).click();
    await page.getByTestId('edit-orders-table').click();
    await expect(page.getByText('Delete')).not.toBeVisible();
  });
});

test.describe('(MultiMove) Uploads', () => {
  test.skip(multiMoveEnabled === 'false', 'Skip if MultiMove workflow is not enabled.');

  test('Users can upload but cannot delete orders once move has been submitted', async ({ page, customerPage }) => {
    // Generate a move that has the status of SUBMITTED
    const move = await customerPage.testHarness.buildSubmittedMoveWithPPMShipmentForSC();
    const userId = move?.Orders?.service_member?.user_id;
    await customerPage.signInAsExistingCustomer(userId);
    await customerPage.navigateFromMMDashboardToMove(move);
    await customerPage.waitForPage.home();
    await page.getByRole('button', { name: 'Review your request' }).click();
    await page.getByTestId('edit-orders-table').click();
    await expect(page.getByText('Delete')).not.toBeVisible();
  });
});
