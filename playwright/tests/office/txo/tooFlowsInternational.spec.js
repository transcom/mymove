import { test, expect } from '../../utils/office/officeTest';

import { TooFlowPage } from './tooTestFixture';

const alaskaEnabled = process.env.FEATURE_FLAG_ENABLE_ALASKA;

test.describe('TOO user', () => {
  /** @type {TooFlowPage} */
  let tooFlowPage;

  test.describe('with HHG Moves', () => {
    test.skip(alaskaEnabled === 'false', 'Skip if Alaska FF is disabled.');
    test('is able to approve an AK iHHG shipment that generates 4 basic service items', async ({
      page,
      officePage,
    }) => {
      const move = await officePage.testHarness.buildInternationalAlaskaBasicHHGMoveForTOO();
      await officePage.signInAsNewTOOUser();
      tooFlowPage = new TooFlowPage(officePage, move);
      await tooFlowPage.waitForLoading();
      await officePage.tooNavigateToMove(tooFlowPage.moveLocator);
      await expect(page.locator('#approved-shipments')).not.toBeVisible();
      await expect(page.locator('#requested-shipments')).toBeVisible();
      await expect(page.getByText('Approve selected')).toBeDisabled();
      await expect(page.locator('#approvalConfirmationModal [data-testid="modal"]')).not.toBeVisible();

      await tooFlowPage.waitForLoading();
      await tooFlowPage.approveAllShipments();

      // Redirected to Move Task Order page - should have 4 basic iHHG service items
      expect(page.url()).toContain(`/moves/${tooFlowPage.moveLocator}/mto`);
      await expect(page.getByTestId('ShipmentContainer')).toBeVisible();
      await expect(page.locator('[data-testid="ApprovedServiceItemsTable"] h3')).toContainText(
        'Approved Service Items (4 items)',
      );

      // Navigate back to Move Details
      await page.getByTestId('MoveDetails-Tab').click();
      await tooFlowPage.waitForLoading();

      expect(page.url()).toContain(`/moves/${tooFlowPage.moveLocator}/details`);
      await expect(page.locator('#approvalConfirmationModal [data-testid="modal"]')).not.toBeVisible();
      await expect(page.locator('#approved-shipments')).toBeVisible();
      await expect(page.locator('#requested-shipments')).not.toBeVisible();
      await expect(page.getByText('Approve selected')).not.toBeVisible();
    });

    test.skip(alaskaEnabled === 'false', 'Skip if Alaska FF is disabled.');
    test('is able to approve and reject international crating/uncrating service items', async ({
      officePage,
      page,
    }) => {
      const move = await officePage.testHarness.buildHHGMoveWithIntlCratingServiceItemsTOO();
      await officePage.signInAsNewTOOUser();
      tooFlowPage = new TooFlowPage(officePage, move);
      await tooFlowPage.waitForLoading();
      await officePage.tooNavigateToMove(tooFlowPage.moveLocator);

      // Edit the shipment address to AK
      await page.locator('[data-testid="ShipmentContainer"] .usa-button').first().click();
      await page.locator('input[id="delivery.address-location-input"]').fill('99505');
      await page.keyboard.press('Enter');

      await page.getByRole('button', { name: 'Save' }).click();
      await tooFlowPage.waitForPage.moveDetails();

      await tooFlowPage.waitForLoading();
      await tooFlowPage.approveAllShipments();

      await page.getByTestId('MoveTaskOrder-Tab').click();
      await tooFlowPage.waitForLoading();
      expect(page.url()).toContain(`/moves/${tooFlowPage.moveLocator}/mto`);

      // Wait for page to load to deal with flakiness resulting from Service Item tables loading
      await tooFlowPage.page.waitForLoadState();

      // Move Task Order page
      await expect(page.getByTestId('ShipmentContainer')).toHaveCount(1);

      /**
       * @function
       * @description This test approves and rejects service items, which moves them from one table to another
       * and expects the counts of each table to increment/decrement by one item each time
       * This function gets the service items for a given table to help count them
       * @param {import("playwright-core").Locator} table
       * @returns {import("playwright-core").Locator}
       */
      const getServiceItemsInTable = (table) => {
        return table.getByRole('rowgroup').nth(1).getByRole('row');
      };

      const requestedServiceItemsTable = page.getByTestId('RequestedServiceItemsTable');
      let requestedServiceItemCount = await getServiceItemsInTable(requestedServiceItemsTable).count();
      const approvedServiceItemsTable = page.getByTestId('ApprovedServiceItemsTable');
      let approvedServiceItemCount = await getServiceItemsInTable(approvedServiceItemsTable).count();
      const rejectedServiceItemsTable = page.getByTestId('RejectedServiceItemsTable');
      let rejectedServiceItemCount = await getServiceItemsInTable(rejectedServiceItemsTable).count();

      await expect(page.getByText('Requested Service Items', { exact: false })).toBeVisible();
      await expect(getServiceItemsInTable(requestedServiceItemsTable).nth(1)).toBeVisible();

      await expect(page.getByTestId('modal')).not.toBeVisible();

      // Approve a requested service item
      expect((await getServiceItemsInTable(requestedServiceItemsTable).count()) > 0);
      // ICRT
      await requestedServiceItemsTable.getByRole('button', { name: 'Accept' }).first().click();
      await tooFlowPage.waitForLoading();

      await expect(getServiceItemsInTable(approvedServiceItemsTable)).toHaveCount(approvedServiceItemCount + 1);
      approvedServiceItemCount = await getServiceItemsInTable(approvedServiceItemsTable).count();

      await expect(getServiceItemsInTable(requestedServiceItemsTable)).toHaveCount(requestedServiceItemCount - 1);
      requestedServiceItemCount = await getServiceItemsInTable(requestedServiceItemsTable).count();

      // IUCRT
      await requestedServiceItemsTable.getByRole('button', { name: 'Accept' }).first().click();
      await tooFlowPage.waitForLoading();

      await expect(getServiceItemsInTable(approvedServiceItemsTable)).toHaveCount(approvedServiceItemCount + 1);
      approvedServiceItemCount = await getServiceItemsInTable(approvedServiceItemsTable).count();

      await expect(getServiceItemsInTable(requestedServiceItemsTable)).toHaveCount(requestedServiceItemCount - 1);
      requestedServiceItemCount = await getServiceItemsInTable(requestedServiceItemsTable).count();

      // Reject a requested service item
      await expect(page.getByText('Requested Service Items', { exact: false })).toBeVisible();
      expect((await getServiceItemsInTable(requestedServiceItemsTable).count()) > 0);
      // ICRT
      await requestedServiceItemsTable.getByRole('button', { name: 'Reject' }).first().click();

      await expect(page.getByTestId('modal')).toBeVisible();
      let modal = page.getByTestId('modal');

      await expect(modal.getByRole('button', { name: 'Submit' })).toBeDisabled();
      await modal.getByRole('textbox').fill('my very valid reason');
      await modal.getByRole('button', { name: 'Submit' }).click();

      await expect(page.getByTestId('modal')).not.toBeVisible();

      await expect(page.getByText('Rejected Service Items', { exact: false })).toBeVisible();
      await expect(getServiceItemsInTable(rejectedServiceItemsTable)).toHaveCount(rejectedServiceItemCount + 1);
      rejectedServiceItemCount = await getServiceItemsInTable(rejectedServiceItemsTable).count();

      await expect(getServiceItemsInTable(requestedServiceItemsTable)).toHaveCount(requestedServiceItemCount - 1);
      requestedServiceItemCount = await getServiceItemsInTable(requestedServiceItemsTable).count();

      // IUCRT
      await requestedServiceItemsTable.getByRole('button', { name: 'Reject' }).first().click();

      await expect(page.getByTestId('modal')).toBeVisible();
      modal = page.getByTestId('modal');

      await expect(modal.getByRole('button', { name: 'Submit' })).toBeDisabled();
      await modal.getByRole('textbox').fill('my very valid reason');
      await modal.getByRole('button', { name: 'Submit' }).click();

      await expect(page.getByTestId('modal')).not.toBeVisible();

      await expect(page.getByText('Rejected Service Items', { exact: false })).toBeVisible();
      await expect(getServiceItemsInTable(rejectedServiceItemsTable)).toHaveCount(rejectedServiceItemCount + 1);
      rejectedServiceItemCount = await getServiceItemsInTable(rejectedServiceItemsTable).count();

      await expect(getServiceItemsInTable(requestedServiceItemsTable)).toHaveCount(requestedServiceItemCount - 1);
    });
  });
});
