// @ts-check
import { test, expect } from './servicesCounselingTestFixture';

test.describe('Services counselor user', () => {
  let tac;
  test.beforeEach(async ({ scPage }) => {
    const move = await scPage.testHarness.buildHHGMoveWithNTSAndNeedsSC();
    tac = await scPage.testHarness.buildGoodTACAndLoaCombination();
    await scPage.navigateToMove(move.locator);
  });

  test('Services Counselor can delete/remove an NTS shipment request', async ({ page, scPage }) => {
    // this test is almost identical to the NTSR test
    await scPage.addNTSShipment();

    // Verify that there are two shipments -- a single HHG plus added NTS
    await expect(page.getByTestId('ShipmentContainer')).toHaveCount(2);

    // Find the card for the NTS shipment, then click the edit button
    const container = await scPage.getShipmentContainerByType('NTS');
    await container.getByRole('button', { name: 'Edit Shipment' }).click();
    await scPage.waitForPage.editNTSShipment();

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
    await scPage.addNTSShipment();

    // Find the card for the NTS shipment, then click the edit button
    const container = await scPage.getShipmentContainerByType('NTS');
    await container.getByRole('button', { name: 'Edit Shipment' }).click();
    await scPage.waitForPage.editNTSShipment();

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
    await page.getByLabel('Orders type', { exact: true }).selectOption({ label: 'Permanent Change Of Station (PCS)' });
    await page.getByLabel('Orders type detail').selectOption({ label: 'Shipment of HHG Permitted' });

    // Click the Save button and return to the move details page
    await page.getByRole('button', { name: 'Save' }).click();
    await scPage.waitForPage.moveDetails();

    // Verify that orders information saved correctly
    await expect(page.getByTestId('tacMDC')).toHaveText('E15A');

    // Return to editing the NTS shipment
    await container.getByRole('button', { name: 'Edit Shipment' }).click();
    await scPage.waitForPage.editNTSShipment();

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

  test('Service Counselor can identify an active NTS LOA based on the date in which they are submitting the move', async ({
    page,
    scPage,
  }) => {
    // this test is almost identical to the NTSR test
    await scPage.addNTSShipment();

    // Find the card for the NTS shipment, then click the edit button
    const container = await scPage.getShipmentContainerByType('NTS');
    await container.getByRole('button', { name: 'Edit Shipment' }).click();
    await scPage.waitForPage.editNTSShipment();

    // Click the "Add or edit codes" button to be taken to the orders form
    await page.getByRole('button', { name: 'Add or edit codes' }).click();
    await scPage.waitForPage.moveOrders();

    // Fill out the orders details
    await page.getByLabel('Orders number').fill('1234');
    await page.getByLabel('Date issued').fill('1234');
    await page.getByLabel('Department indicator').selectOption({ label: '21 Army' });
    await page.getByLabel('Orders type', { exact: true }).selectOption({ label: 'Permanent Change Of Station (PCS)' });
    await page.getByLabel('Orders type detail').selectOption({ label: 'Shipment of HHG Permitted' });

    // Fill out the HHG and NTS accounting codes
    await page.getByTestId('hhgTacInput').fill(tac.tac);
    const today = new Date();
    const day = new Intl.DateTimeFormat('en', { day: '2-digit' }).format(today);
    const month = new Intl.DateTimeFormat('en', { month: 'short' }).format(today);
    const year = new Intl.DateTimeFormat('en', { year: 'numeric' }).format(today);
    const formattedDate = `${day} ${month} ${year}`;

    await page.locator('input[name="issueDate"]').fill(formattedDate);

    await page.getByTestId('hhgSacInput').fill('4K988AS098F');
    // Today's date will fall valid under the TAC and LOA and the NTS LOA should then populate
    await page.getByTestId('ntsTacInput').fill(tac.tac);
    const ntsLoaTextField = await page.getByTestId('ntsLoaTextField');
    await expect(ntsLoaTextField).toHaveValue('1*1*20232025*1*1*1*1*1*1*1*1*1*1*1*1*1*1*1*1*1*1*1*1*1*1*1*1*1');

    const loaMissingErrorMessage = page.getByText('Unable to find a LOA based on the provided details');
    const loaInvalidErrorMessage = page.getByText(
      'The LOA identified based on the provided details appears to be invalid',
    );

    await expect(loaMissingErrorMessage).not.toBeVisible();
    await expect(loaInvalidErrorMessage).not.toBeVisible();
  });
});
