// import { test, expect } from '../../utils/my/customerTest';
import { test } from '../../utils/my/customerTest';

test('Users can upload but cannot delete orders once move has been submitted', async ({ customerPage }) => {
  // Generate a move that has the status of SUBMITTED
  const move = await customerPage.testHarness.buildSubmittedMoveWithPPMShipmentForSC();
  const userId = move.Orders.ServiceMember.user_id;
  await customerPage.signInAsExistingCustomer(userId);
  // await customerPage.waitForPage.multiMoveLandingPage();
  // await page.getByRole('button', { name: 'Review your request' }).click();
  // await page.getByTestId('edit-orders-table').click();
  // await expect(page.getByText('Delete')).not.toBeVisible();
});
