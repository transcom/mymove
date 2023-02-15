/**
 * Semi-automated converted from a cypress test, and thus may contain
 * non best-practices, in particular: heavy use of `page.locator`
 * instead of `page.getBy*`.
 */

// @ts-check
import { test, expect } from '../../utils/officeTest';

import { TooFlowPage } from './tooTestFixture';

test.describe('TOO user', () => {
  /** @type {TooFlowPage} */
  let tooFlowPage;
  test.describe('with unapproved HHG + NTS Move', () => {
    test.beforeEach(async ({ officePage }) => {
      const move = await officePage.testHarness.buildHHGMoveWithNTSShipmentsForTOO();
      await officePage.signInAsNewTOOUser();
      tooFlowPage = new TooFlowPage(officePage, move);
      await officePage.tooNavigateToMove(move.locator);
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
      await page.locator('#requestedPickupDate').type('16 Mar 2022');
      await page.getByText('Use current address').click();

      // Storage facility info
      await page.locator('#facilityName').type('Sample Facility Name');
      await page.locator('#facilityName').blur();
      await page.locator('#facilityPhone').type('999-999-9999');
      await page.locator('#facilityPhone').blur();
      await page.locator('#facilityEmail').type('sample@example.com');
      await page.locator('#facilityEmail').blur();
      await page.locator('#facilityServiceOrderNumber').type('999999');
      await page.locator('#facilityServiceOrderNumber').blur();

      // Storage facility address
      await page.locator('input[name="storageFacility.address.streetAddress1"]').type('148 S East St');
      await page.locator('input[name="storageFacility.address.streetAddress1"]').blur();
      await page.locator('input[name="storageFacility.address.streetAddress2"]').type('Suite 7A');
      await page.locator('input[name="storageFacility.address.streetAddress2"]').blur();
      await page.locator('input[name="storageFacility.address.city"]').type('Sample City');
      await page.locator('input[name="storageFacility.address.city"]').blur();
      await page.locator('select[name="storageFacility.address.state"]').selectOption({ label: 'GA' });
      await page.locator('input[name="storageFacility.address.postalCode"]').type('30301');
      await page.locator('input[name="storageFacility.address.postalCode"]').blur();
      await page.locator('#facilityLotNumber').type('1111111');
      await page.locator('#facilityLotNumber').blur();

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
  });

  test.describe('with approved HHG + NTS Move', () => {
    test.beforeEach(async ({ officePage }) => {
      const move = await officePage.testHarness.buildHHGMoveWithApprovedNTSShipmentsForTOO();
      await officePage.signInAsNewTOOUser();
      tooFlowPage = new TooFlowPage(officePage, move);
      await officePage.tooNavigateToMove(move.locator);
    });

    test('TOO can view and edit Domestic NTS Shipments handled by the Prime on the MTO page', async ({ page }) => {
      // This test is almost exactly a duplicate of the test in
      // tooFlowsNTSR.
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
      await modal.locator('#facilityName').type('New Facility Name');
      await modal.locator('#facilityPhone').clear();
      await modal.locator('#facilityPhone').type('999-999-9999');
      await modal.locator('#facilityEmail').clear();
      await modal.locator('#facilityEmail').type('new@example.com');
      await modal.locator('#facilityServiceOrderNumber').clear();
      await modal.locator('#facilityServiceOrderNumber').type('098098');

      // Storage facility address
      await modal.locator('input[name="storageFacility.address.streetAddress1"]').clear();
      await modal.locator('input[name="storageFacility.address.streetAddress1"]').type('265 S East St');
      await modal.locator('#facilityLotNumber').clear();
      await modal.locator('#facilityLotNumber').type('1111111');

      await modal.locator('button[type="submit"]').click();
      await expect(modal).not.toBeVisible();

      lastShipment = page.locator('[data-testid="ShipmentContainer"]').last();
      let sidebar = lastShipment.locator('[class*="ShipmentDetailsSidebar"]');
      await expect(sidebar.locator('section header').first()).toContainText('Facility info and address');
      await expect(sidebar.locator('section').first()).toContainText('New Facility Name');
      await expect(sidebar.locator('section').first()).toContainText('265 S East St');
      await expect(sidebar.locator('section').first()).toContainText('Lot 1111111');

      // edit service order number
      await lastShipment.locator('[data-testid="service-order-number-modal-open"]').click();

      await expect(page.getByTestId('modal')).toBeVisible();
      modal = page.getByTestId('modal');

      await modal.locator('[data-testid="textInput"]').clear();
      await modal.locator('[data-testid="textInput"]').type('ORDER456');

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
        'Approved service items (5 items)',
      );
    });
  });

  test.describe('with HHG Move plus NTS Shipment handled by external vendor', () => {
    test.beforeEach(async ({ officePage }) => {
      const move = await officePage.testHarness.buildHHGMoveWithExternalNTSShipmentsForTOO();
      await officePage.signInAsNewTOOUser();
      tooFlowPage = new TooFlowPage(officePage, move);
      await officePage.tooNavigateToMove(move.locator);
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
      await officePage.tooNavigateToMove(move.locator);
    });

    test('can submit service items', async ({ page }) => {
      // This test is almost exactly a duplicate of the test in
      // tooFlowsNTSR.
      await expect(page.locator('#approved-shipments')).not.toBeVisible();
      await expect(page.locator('#requested-shipments')).toBeVisible();
      await expect(page.getByText('Approve selected')).toBeDisabled();

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
