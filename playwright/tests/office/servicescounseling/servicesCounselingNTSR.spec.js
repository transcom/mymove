// @ts-check
import { test, expect } from './servicesCounselingTestFixture';

test.describe('Services counselor user', () => {
  test.beforeEach(async ({ scPage }) => {
    const move = await scPage.testHarness.buildHHGMoveWithNTSAndNeedsSC();
    await scPage.navigateToMove(move.locator);
  });

  test('Services Counselor can delete/remove an NTS-release shipment request', async ({ page, scPage }) => {
    // this test is almost identical to the NTS test
    await scPage.addNTSReleaseShipment();

    // Verify that there are two shipments -- a single HHG plus added NTS
    await expect(page.getByTestId('ShipmentContainer')).toHaveCount(2);

    // Find the card for the NTS shipment, then click the edit button
    const container = await scPage.getShipmentContainerByType('NTS');
    await container.getByRole('button', { name: 'Edit Shipment' }).click();
    await scPage.waitForPage.editNTSReleaseShipment();

    // Click the "Delete Shipment" button to bring up the confirmation modal
    await page.getByRole('button', { name: 'Delete shipment' }).click();
    await expect(page.getByTestId('modal')).toBeVisible();

    // Confirm deletion to be taken back to the move details page
    await page.getByTestId('modal').getByRole('button', { name: 'Delete shipment' }).click();
    await scPage.waitForPage.moveDetails();

    // Verify that there's only 1 shipment displayed now
    await expect(page.getByTestId('ShipmentContainer')).toHaveCount(1);
  });

  test('Services Counselor can enter accounting codes and submit shipment', async ({ page, scPage }) => {
    // this test is almost identical to the NTSR test
    await scPage.addNTSReleaseShipment();

    // Find the card for the NTS-release shipment, then click the edit button
    const container = await scPage.getShipmentContainerByType('NTS-release');
    await container.getByRole('button', { name: 'Edit Shipment' }).click();
    await scPage.waitForPage.editNTSReleaseShipment();

    // Click the "Add or edit codes" button to be taken to the orders form
    await page.getByRole('button', { name: 'Add or edit codes' }).click();
    await scPage.waitForPage.moveOrders();

    // Fill out the HHG and NTS accounting codes
    await page.getByTestId('hhgTacInput').fill('E15A');
    await page.getByTestId('hhgSacInput').fill('4K988AS098F');
    await page.getByTestId('ntsTacInput').fill('F123');
    await page.getByTestId('ntsSacInput').fill('3L988AS098F');

    // Fill out the orders details
    await page.getByLabel('Orders number').fill('1234');
    await page.getByLabel('Department indicator').selectOption({ label: '21 Army' });
    await page
      .getByLabel('Orders type *', { exact: true })
      .selectOption({ label: 'Permanent Change Of Station (PCS)' });
    await page.getByLabel('Orders type detail').selectOption({ label: 'Shipment of HHG Permitted' });

    // Click the Save button and return to the move details page
    await page.getByRole('button', { name: 'Save' }).click();
    await scPage.waitForPage.moveDetails();

    // Verify that orders information saved correctly
    await expect(page.getByTestId('tacMDC')).toHaveText('E15A');

    // Return to editing the NTS-release shipment
    await container.getByRole('button', { name: 'Edit Shipment' }).click();
    await scPage.waitForPage.editNTSReleaseShipment();

    // Select TAC and SAC codes for this shipment, then save
    await page.getByText('F123 (NTS)').click();
    await page.getByText('4K988AS098F (HHG)').click();
    await page.getByRole('button', { name: 'Save' }).click();
    await scPage.waitForPage.moveDetails();

    // Verify that shipment has been been updated
    await expect(page.getByText('Your changes were saved.')).toBeVisible();

    // Click the submit move button and bring up the confirmation modal
    await page.getByRole('button', { name: 'Submit move details' }).click();
    await expect(page.getByTestId('modal')).toBeVisible();

    // Submit the move and verify the success alert displays
    await page.getByRole('button', { name: 'Yes, submit' }).click();
    await expect(page.getByText('Move submitted.')).toBeVisible();
  });
});
