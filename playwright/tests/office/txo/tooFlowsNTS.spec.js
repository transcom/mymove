/**
 * Semi-automated converted from a cypress test, and thus may contain
 * non best-practices, in particular: heavy use of `page.locator`
 * instead of `page.getBy*`.
 */

// @ts-check
import { test, expect } from '../../utils/office/officeTest';

import { TooFlowPage } from './tooTestFixture';

const TOOTabsTitles = ['Move Queue', 'Search'];
const SearchRBSelection = ['Move Code', 'DOD ID', 'Customer Name'];

test.describe('TOO user', () => {
  /** @type {TooFlowPage} */
  let tooFlowPage;
  let tac;
  test.describe('with unapproved HHG + NTS Move', () => {
    test.beforeEach(async ({ officePage }) => {
      const move = await officePage.testHarness.buildHHGMoveWithNTSShipmentsForTOO();
      tac = await officePage.testHarness.buildGoodTACAndLoaCombination();
      await officePage.signInAsNewTOOUser();
      tooFlowPage = new TooFlowPage(officePage, move);
      await tooFlowPage.waitForLoading();
      await officePage.tooNavigateToMove(tooFlowPage.moveLocator);
    });

    // TODO FOR NTS-RELEASE
    // Enter/edit shipment weight if needed -- user will look up in TOPS and manually enter
    // TOO has to enter Service Order Number (SON) for NTS-RELEASE shipments prior to posting to the MTO

    // This test covers editing the NTS shipment and prepares it for approval
    test('TOO can edit a request for Domestic NTS Shipment handled by the Prime', async ({ page }) => {
      // This test is almost exactly a duplicate of the test in
      // tooFlowsNTSR.
      await expect(page.getByText('Approve selected')).toBeDisabled();
      let lastShipment = page.locator('[data-testid="ShipmentContainer"]').last();
      await expect(lastShipment.locator('[data-testid="shipment-display-checkbox"]')).toBeDisabled();
      await lastShipment.locator('[data-icon="chevron-down"]').click();

      await expect(
        lastShipment.locator('div[class*="missingInfoError"] [data-testid="storageFacilityName"]'),
      ).toBeVisible();
      await expect(lastShipment.locator('[data-testid="storageFacilityName"]')).toContainText('Missing');
      await expect(
        lastShipment.locator('div[class*="missingInfoError"] [data-testid="serviceOrderNumber"]'),
      ).toBeVisible();
      await expect(
        lastShipment.locator('div[class*="missingInfoError"] [data-testid="storageFacilityAddress"]'),
      ).toBeVisible();
      await expect(lastShipment.locator('[data-testid="storageFacilityAddress"]')).toContainText('Missing');
      await expect(page.locator('div[class*="missingInfoError"] [data-testid="tacType"]')).toBeVisible();

      await tooFlowPage.editTacSac();

      // Edit shipments to enter missing info
      await page.locator('[data-testid="ShipmentContainer"] .usa-button').last().click();
      // Basic info
      await page.locator('#requestedPickupDate').clear();
      await page.locator('#requestedPickupDate').fill('16 Mar 2022');
      await page.getByText('Use current address').click();

      // Storage facility info
      await page.locator('#facilityName').fill('Sample Facility Name');
      await page.locator('#facilityName').blur();
      await page.locator('#facilityPhone').fill('999-999-9999');
      await page.locator('#facilityPhone').blur();
      await page.locator('#facilityEmail').fill('sample@example.com');
      await page.locator('#facilityEmail').blur();
      await page.locator('#facilityServiceOrderNumber').fill('999999');
      await page.locator('#facilityServiceOrderNumber').blur();

      // Storage facility address
      const StorageLocationLookup = 'ATLANTA, GA 30301 (FULTON)';
      await page.locator('input[name="storageFacility.address.streetAddress1"]').fill('148 S East St');
      await page.locator('input[name="storageFacility.address.streetAddress1"]').blur();
      await page.locator('input[name="storageFacility.address.streetAddress2"]').fill('Suite 7A');
      await page.locator('input[name="storageFacility.address.streetAddress2"]').blur();
      await page.locator('input[id="storageFacility.address-location-input"]').fill('30301');
      await expect(page.getByText(StorageLocationLookup, { exact: true })).toBeVisible();
      await page.keyboard.press('Enter');
      await page.locator('#facilityLotNumber').fill('1111111');
      await page.locator('#facilityLotNumber').blur();

      // Delivery info
      await page.locator('#requestedDeliveryDate').clear();
      await page.locator('#requestedDeliveryDate').fill('16 Mar 2022');

      // TAC and SAC
      await page.locator('[data-testid="radio"] [for="tacType-NTS"]').click();
      await page.locator('[data-testid="radio"] [for="sacType-HHG"]').click();

      await page.locator('[data-testid="submitForm"]').click();
      await tooFlowPage.waitForLoading();

      // edit the NTS shipment to be handled by an external vendor
      await page.locator('[data-testid="ShipmentContainer"] .usa-button').last().click();
      await page.locator('label[for="vendorExternal"]').click();
      await page.locator('[data-testid="submitForm"]').click();
      await tooFlowPage.waitForLoading();

      await expect(page.locator('[data-testid="ShipmentContainer"] [data-testid="tag"]')).toContainText(
        'external vendor',
      );

      // edit the NTS shipment back to being handled by the GHC Prime contractor
      await page.locator('[data-testid="ShipmentContainer"] .usa-button').last().click();
      await expect(page.locator('[data-testid="alert"]')).toContainText(
        'The GHC prime contractor is not handling the shipment.',
      );

      await page.locator('label[for="vendorPrime"]').click();
      await page.locator('[data-testid="submitForm"]').click();
      await tooFlowPage.waitForLoading();

      lastShipment = page.locator('[data-testid="ShipmentContainer"]').last();

      await expect(lastShipment.locator('[data-testid="shipment-display-checkbox"]')).toBeEnabled();
      await expect(lastShipment.locator('[data-testid="tag"]')).not.toBeVisible();
      await lastShipment.locator('[data-icon="chevron-down"]').click();

      await expect(
        lastShipment.locator('div[class*="missingInfoError"] [data-testid="storageFacilityName"]'),
      ).not.toBeVisible();
      await expect(
        lastShipment.locator('div[class*="missingInfoError"] [data-testid="serviceOrderNumber"]'),
      ).not.toBeVisible();
      await expect(
        lastShipment.locator('div[class*="missingInfoError"] [data-testid="storageFacilityAddress"]'),
      ).not.toBeVisible();
      await expect(lastShipment.locator('div[class*="missingInfoError"] [data-testid="tacType"]')).not.toBeVisible();

      // Combine this test with the one above as it needs the state
      // after the test above

      // test('TOO can approve an NTS shipment handled by the Prime')

      // Make sure that it shows all the relevant information on the approve page
      // captures the information about the NTS Facility and relevant storage information
      // verifies that all the information is shown for NTS shipments handled by the GHC Contractor
      await expect(page.locator('#approved-shipments')).not.toBeVisible();
      await expect(page.locator('#requested-shipments')).toBeVisible();
      await expect(page.getByText('Approve selected')).toBeDisabled();

      // Select & approve items
      await expect(page.locator('input[data-testid="shipment-display-checkbox"]')).toHaveCount(2);
      await tooFlowPage.approveAllShipments();

      // Redirected to Move Task Order page
      expect(page.url()).toContain(`/moves/${tooFlowPage.moveLocator}/mto`);

      // confirm the changes made above
      lastShipment = page.locator('[data-testid="ShipmentContainer"]').last();
      // pickup address header
      await expect(lastShipment.locator('[class*="ShipmentAddresses_mtoShipmentAddresses"]')).toContainText(
        'Pickup address',
      );
      // facility address header
      await expect(lastShipment.locator('[class*="ShipmentAddresses_mtoShipmentAddresses"]')).toContainText(
        'Facility address',
      );

      // Confirming expected data
      // facility info and address
      // service order number
      // accounting codes
      const details = lastShipment.locator('[class*="ShipmentDetailsSidebar"]').first();
      await expect(details).toBeVisible();
      await expect(details.locator('section header').first().getByText('Facility info and address')).toBeVisible();
      await expect(details.locator('section').first()).toContainText('148 S East St');

      await expect(details.locator('section header').nth(1).getByText('Service order number')).toBeVisible();
      await expect(details.locator('section').nth(1)).toContainText('999999');

      await expect(details.locator('section header').last().getByText('Accounting codes')).toBeVisible();
      await expect(details.locator('section').last()).toContainText('F123');
    });

    test('NTS LOA populates on orders page with good TGET data for TOO', async ({ page }) => {
      await page.getByTestId('edit-orders').click();

      // Fill out the orders details
      await page.getByLabel('Orders number').fill('1234');
      await page.getByLabel('Date issued').fill('1234');
      await page.getByLabel('Department indicator').selectOption({ label: '21 Army' });
      await page
        .getByLabel('Orders type', { exact: true })
        .selectOption({ label: 'Permanent Change Of Station (PCS)' });
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
      // "GOOD" is a hard-set GOOD TAC by the e2e seed data
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

  test.describe('with approved HHG + NTS Move', () => {
    let move;

    test.beforeEach(async ({ officePage }) => {
      move = await officePage.testHarness.buildHHGMoveWithApprovedNTSShipmentsForTOO();
      await officePage.signInAsNewTOOUser();
      tooFlowPage = new TooFlowPage(officePage, move);
      await tooFlowPage.waitForLoading();
      const searchTab = officePage.page.getByTitle(TOOTabsTitles[1]);
      await searchTab.click();
    });

    test('TOO can view and edit Domestic NTS Shipments handled by the Prime on the MTO page', async ({ page }) => {
      // This test is almost exactly a duplicate of the test in
      // tooFlowsNTSR.
      const selectedRadio = page.getByRole('group').locator(`label:text("${SearchRBSelection[0]}")`);
      await selectedRadio.click();
      await page.getByTestId('searchText').fill(move.locator);
      await page.getByTestId('searchTextSubmit').click();

      await expect(page.getByText('Results (1)')).toBeVisible();
      await expect(page.getByTestId('locator-0')).toContainText(move.locator);

      // await page.getByTestId('MoveTaskOrder-Tab').click();
      await page.getByTestId('locator-0').click();

      await page.getByTestId('MoveTaskOrder-Tab').click();
      await tooFlowPage.waitForLoading();
      await expect(
        page.locator('[id="move-weights"] div').getByText('1 shipment not moved by GHC prime.'),
      ).not.toBeVisible();

      await expect(page.locator('[data-testid="ShipmentContainer"]')).toHaveCount(2);

      let lastShipment = page.locator('[data-testid="ShipmentContainer"]').last();

      // non-temp storage header
      await expect(lastShipment.locator('h2')).toContainText('Non-temp storage');
      // pickup address header
      await expect(lastShipment.locator('[class*="ShipmentAddresses_mtoShipmentAddresses"]')).toContainText(
        'Pickup address',
      );
      // facility address header
      await expect(lastShipment.locator('[class*="ShipmentAddresses_mtoShipmentAddresses"]')).toContainText(
        'Facility address',
      );
      // edit facility info and address
      await page.locator('[data-testid="edit-facility-info-modal-open"]').click();

      await expect(page.getByTestId('modal')).toBeVisible();
      let modal = page.getByTestId('modal');
      // Storage facility info
      await modal.locator('#facilityName').clear();
      await modal.locator('#facilityName').fill('New Facility Name');
      await modal.locator('#facilityPhone').clear();
      await modal.locator('#facilityPhone').fill('999-999-9999');
      await modal.locator('#facilityEmail').clear();
      await modal.locator('#facilityEmail').fill('new@example.com');
      await modal.locator('#facilityServiceOrderNumber').clear();
      await modal.locator('#facilityServiceOrderNumber').fill('098098');

      // Storage facility address
      await modal.locator('input[name="storageFacility.address.streetAddress1"]').clear();
      await modal.locator('input[name="storageFacility.address.streetAddress1"]').fill('265 S East St');
      await modal.locator('#facilityLotNumber').clear();
      await modal.locator('#facilityLotNumber').fill('1111111');

      await modal.locator('button[type="submit"]').click();
      await expect(modal).not.toBeVisible();

      lastShipment = page.locator('[data-testid="ShipmentContainer"]').last();
      let sidebar = lastShipment.locator('[class*="ShipmentDetailsSidebar"]');
      await expect(sidebar.locator('section header').first()).toContainText('Facility info and address');
      await expect(sidebar.locator('section').first()).toContainText('Storage R Us');
      await expect(sidebar.locator('section').first()).toContainText('265 S East St');
      await expect(sidebar.locator('section').first()).toContainText('Lot 1111111');

      // edit service order number
      await lastShipment.locator('[data-testid="service-order-number-modal-open"]').click();

      await expect(page.getByTestId('modal')).toBeVisible();
      modal = page.getByTestId('modal');

      await modal.locator('[data-testid="textInput"]').clear();
      await modal.locator('[data-testid="textInput"]').fill('ORDER456');

      await modal.locator('button[type="submit"]').click();
      await expect(modal).not.toBeVisible();

      lastShipment = page.locator('[data-testid="ShipmentContainer"]').last();
      await expect(lastShipment.locator('[class*="ShipmentDetailsSidebar"] section').nth(1)).toContainText('ORDER456');

      // edit accounting codes
      await lastShipment.locator('[data-testid="edit-accounting-codes-modal-open"]').click();

      await expect(page.getByTestId('modal')).toBeVisible();
      modal = page.getByTestId('modal');
      await modal.locator('[data-testid="radio"] [for="tacType-HHG"]').click();
      await modal.locator('[data-testid="radio"] [for="sacType-NTS"]').click();

      await modal.locator('button[type="submit"]').click();
      await expect(modal).not.toBeVisible();

      lastShipment = page.locator('[data-testid="ShipmentContainer"]').last();
      sidebar = lastShipment.locator('[class*="ShipmentDetailsSidebar"]');
      await expect(sidebar.locator('section').last()).toContainText('F123');
      await expect(sidebar.locator('section').last()).toContainText('4K988AS098F');

      await expect(lastShipment.locator('[data-testid="ApprovedServiceItemsTable"] h3').last()).toContainText(
        'Approved Service Items (5 items)',
      );
    });
  });

  test.describe('with HHG Move plus NTS Shipment handled by external vendor', () => {
    test.beforeEach(async ({ officePage }) => {
      const move = await officePage.testHarness.buildHHGMoveWithExternalNTSShipmentsForTOO();
      await officePage.signInAsNewTOOUser();
      tooFlowPage = new TooFlowPage(officePage, move);
      await tooFlowPage.waitForLoading();
      await officePage.tooNavigateToMove(tooFlowPage.moveLocator);
    });

    test('can approve an shipment', async ({ page }) => {
      // This test is almost exactly a duplicate of the test in
      // tooFlowsNTSR.
      await expect(page.locator('#approved-shipments')).not.toBeVisible();
      await expect(page.locator('#requested-shipments')).toBeVisible();
      await expect(page.getByText('Approve selected')).toBeDisabled();

      const lastShipment = page.locator('[data-testid="ShipmentContainer"]').last();

      await lastShipment.locator('[data-icon="chevron-down"]').click();

      await expect(
        lastShipment.locator('div[class*="missingInfoError"] [data-testid="storageFacilityName"]'),
      ).toBeVisible();
      await expect(lastShipment.getByTestId('storageFacilityName')).toContainText('Missing');
      await expect(
        lastShipment.locator('div[class*="missingInfoError"] [data-testid="serviceOrderNumber"]'),
      ).toBeVisible();
      await expect(
        lastShipment.locator('div[class*="missingInfoError"] [data-testid="storageFacilityAddress"]'),
      ).toBeVisible();
      await expect(lastShipment.getByTestId('storageFacilityAddress')).toContainText('Missing');
      await expect(lastShipment.locator('div[class*="missingInfoError"] [data-testid="tacType"]')).toBeVisible();

      await tooFlowPage.editTacSac();

      await tooFlowPage.approveAllShipments();

      // Redirected to Move Task Order page
      expect(page.url()).toContain(`/moves/${tooFlowPage.moveLocator}/mto`);

      // Confirm estimated weight shows expected extra shipment detail link
      await expect(
        page.locator('[id="move-weights"] div').getByText('1 shipment not moved by GHC prime.'),
      ).toBeVisible();
    });
  });

  test.describe('with NTS-only Move handled by external vendor', () => {
    test.beforeEach(async ({ officePage }) => {
      const move = await officePage.testHarness.buildMoveWithNTSShipmentsForTOO();
      await officePage.signInAsNewTOOUser();
      tooFlowPage = new TooFlowPage(officePage, move);
      await tooFlowPage.waitForLoading();
      await officePage.tooNavigateToMove(tooFlowPage.moveLocator);
    });

    test('can submit service items', async ({ page }) => {
      // This test is almost exactly a duplicate of the test in
      // tooFlowsNTSR.
      await expect(page.locator('#approved-shipments')).not.toBeVisible();
      await expect(page.locator('#requested-shipments')).toBeVisible();
      await expect(page.getByText('Approve selected')).toBeEnabled();

      const lastShipment = page.locator('[data-testid="ShipmentContainer"]').last();

      await expect(lastShipment.locator('[data-testid="shipment-display-checkbox"]')).not.toBeVisible();
      await lastShipment.locator('[data-icon="chevron-down"]').click();

      await expect(
        lastShipment.locator('div[class*="missingInfoError"] [data-testid="storageFacilityName"]'),
      ).toBeVisible();
      await expect(lastShipment.getByTestId('storageFacilityName')).toContainText('Missing');
      await expect(
        lastShipment.locator('div[class*="missingInfoError"] [data-testid="serviceOrderNumber"]'),
      ).toBeVisible();
      await expect(
        lastShipment.locator('div[class*="missingInfoError"] [data-testid="storageFacilityAddress"]'),
      ).toBeVisible();
      await expect(lastShipment.getByTestId('storageFacilityAddress')).toContainText('Missing');
      await expect(lastShipment.locator('div[class*="missingInfoError"] [data-testid="tacType"]')).toBeVisible();

      await tooFlowPage.editTacSac();

      await tooFlowPage.approveAllShipments({ hasServiceItems: false });

      // Redirected to Move Task Order page
      expect(page.url()).toContain(`/moves/${tooFlowPage.moveLocator}/mto`);

      // Confirm estimated weight shows expected extra shipment detail link
      await expect(
        page.getByRole('main').getByText('This move does not have any approved shipments yet.'),
      ).toBeVisible();
    });
  });
});
