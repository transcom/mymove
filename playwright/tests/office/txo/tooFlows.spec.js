/**
 * Semi-automated converted from a cypress test, and thus may contain
 * non best-practices, in particular: heavy use of `page.locator`
 * instead of `page.getBy*`.
 */

// @ts-check
import { test, expect, DEPARTMENT_INDICATOR_OPTIONS } from '../../utils/office/officeTest';
import findOptionWithinOpenedDropdown from '../../utils/playwrightUtility';

import { TooFlowPage } from './tooTestFixture';

const TOOTabsTitles = ['Move Queue', 'Search'];

const SearchRBSelection = ['Move Code', 'DOD ID', 'Customer Name'];

const SearchTerms = ['SITEXT', '8796353598', 'Spacemen'];

const StatusFilterOptions = ['Draft', 'New Move', 'Needs Counseling', 'Service counseling completed', 'Move approved'];

test.describe('TOO user', () => {
  /** @type {TooFlowPage} */
  let tooFlowPage;
  let testMove;

  test.describe('with Search Queue', () => {
    test.beforeEach(async ({ officePage }) => {
      testMove = await officePage.testHarness.buildHHGMoveWithServiceItemsAndPaymentRequestsAndFilesForTOO();
      await officePage.signInAsNewTOOUser();
      tooFlowPage = new TooFlowPage(officePage, testMove);
      await tooFlowPage.waitForLoading();
      const searchTab = officePage.page.getByTitle(TOOTabsTitles[1]);
      await searchTab.click();
    });

    test('can search for moves using Move Code', async ({ page }) => {
      const selectedRadio = page.getByRole('group').locator(`label:text("${SearchRBSelection[0]}")`);
      await selectedRadio.click();
      await page.getByTestId('searchText').fill(testMove.locator);
      await page.getByTestId('searchTextSubmit').click();

      await expect(page.getByText('Results (1)')).toBeVisible();
      await expect(page.getByTestId('locator-0')).toContainText(testMove.locator);
    });
    test('can search for moves using DOD ID', async ({ page }) => {
      const selectedRadio = page.getByRole('group').locator(`label:text("${SearchRBSelection[1]}")`);
      await selectedRadio.click();
      await page.getByTestId('searchText').fill(testMove.Orders.ServiceMember.edipi);
      await page.getByTestId('searchTextSubmit').click();

      await expect(page.getByText('Results (1)')).toBeVisible();
      await expect(page.getByTestId('edipi-0')).toContainText(testMove.Orders.ServiceMember.edipi);
    });
    test('can search for moves using Customer Name', async ({ page }) => {
      const CustomerName = `${testMove.Orders.ServiceMember.last_name}, ${testMove.Orders.ServiceMember.first_name}`;
      const selectedRadio = page.getByRole('group').locator(`label:text("${SearchRBSelection[2]}")`);
      await selectedRadio.click();
      await page.getByTestId('searchText').fill(CustomerName);
      await page.getByTestId('searchTextSubmit').click();

      await expect(page.getByText('Results')).toBeVisible();
      await expect(page.getByTestId('customerName-0')).toContainText(CustomerName);
    });
    test('Can filter status using Move Status', async ({ page }) => {
      const selectedRadio = page.getByRole('group').locator(`label:text("${SearchRBSelection[0]}")`);
      await selectedRadio.click();
      await page.getByTestId('searchText').fill(SearchTerms[0]);
      await page.getByTestId('searchTextSubmit').click();

      const StatusFilter = page.getByTestId('MultiSelectCheckBoxFilter');
      await StatusFilter.click();

      for (const item of StatusFilterOptions) {
        const found = page
          .locator('[id^="react-select"][id*="listbox"]')
          .locator(`[id*="option"]:has(:text("${item}"))`);
        await expect(found).toBeVisible();
      }
    });
    test('Can select a filter status using Payment Request', async ({ page }) => {
      const selectedRadio = page.getByRole('group').locator(`label:text("${SearchRBSelection[0]}")`);
      await selectedRadio.click();
      await page.getByTestId('searchText').fill(testMove.locator);
      await page.getByTestId('searchTextSubmit').click();

      // Check if Payment Request Status options are present
      const StatusFilter = page.getByTestId('MultiSelectCheckBoxFilter');
      await StatusFilter.click();

      const found = findOptionWithinOpenedDropdown(page, StatusFilterOptions[1]);
      await found.click();
      await expect(page.getByText('Results')).toBeVisible();
    });
    test('cant search for empty move code', async ({ page }) => {
      const selectedRadio = page.getByRole('group').locator(`label:text("${SearchRBSelection[0]}")`);
      await selectedRadio.click();

      const SearchBox = page.getByTestId('searchText');
      await SearchBox.fill('');
      await SearchBox.blur();

      await expect(page.getByText('Move Code Must be exactly 6 characters')).toBeVisible();
      await expect(page.getByRole('table')).not.toBeVisible();
    });
    test('cant search for short move code', async ({ page }) => {
      const selectedRadio = page.getByRole('group').locator(`label:text("${SearchRBSelection[0]}")`);
      await selectedRadio.click();

      const SearchBox = page.getByTestId('searchText');
      await SearchBox.fill('MOVE');
      await SearchBox.blur();

      await expect(page.getByText('Move Code Must be exactly 6 characters')).toBeVisible();
      await expect(page.getByRole('table')).not.toBeVisible();
    });
    test('cant search for long move code', async ({ page }) => {
      const selectedRadio = page.getByRole('group').locator(`label:text("${SearchRBSelection[0]}")`);
      await selectedRadio.click();

      const SearchBox = page.getByTestId('searchText');
      await SearchBox.fill('ASUPERLONGMOVE');
      await SearchBox.blur();

      await expect(page.getByText('Move Code Must be exactly 6 characters')).toBeVisible();
      await expect(page.getByRole('table')).not.toBeVisible();
    });
    test('cant search for empty DOD ID', async ({ page }) => {
      const selectedRadio = page.getByRole('group').locator(`label:text("${SearchRBSelection[1]}")`);
      await selectedRadio.click();

      const SearchBox = page.getByTestId('searchText');
      await SearchBox.fill('');
      await SearchBox.blur();

      await expect(page.getByText('DOD ID must be exactly 10 characters')).toBeVisible();
      await expect(page.getByRole('table')).not.toBeVisible();
    });
    test('cant search for short DOD ID', async ({ page }) => {
      const selectedRadio = page.getByRole('group').locator(`label:text("${SearchRBSelection[1]}")`);
      await selectedRadio.click();

      const SearchBox = page.getByTestId('searchText');
      await SearchBox.fill('1234567');
      await SearchBox.blur();

      await expect(page.getByText('DOD ID must be exactly 10 characters')).toBeVisible();
      await expect(page.getByRole('table')).not.toBeVisible();
    });
    test('cant search for long DOD ID', async ({ page }) => {
      const selectedRadio = page.getByRole('group').locator(`label:text("${SearchRBSelection[1]}")`);
      await selectedRadio.click();

      const SearchBox = page.getByTestId('searchText');
      await SearchBox.fill('123456789011');
      await SearchBox.blur();

      await expect(page.getByText('DOD ID must be exactly 10 characters')).toBeVisible();
      await expect(page.getByRole('table')).not.toBeVisible();
    });
    test('cant search for empty Customer Name', async ({ page }) => {
      const selectedRadio = page.getByRole('group').locator(`label:text("${SearchRBSelection[2]}")`);
      await selectedRadio.click();

      const SearchBox = page.getByTestId('searchText');
      await SearchBox.fill('');
      await SearchBox.blur();

      await expect(page.getByText('Customer search must contain a value')).toBeVisible();
      await expect(page.getByRole('table')).not.toBeVisible();
    });
  });

  test.describe('with HHG Moves', () => {
    test.beforeEach(async ({ officePage }) => {
      const move = await officePage.testHarness.buildHHGMoveWithServiceItemsAndPaymentRequestsAndFilesForTOO();
      await officePage.signInAsNewTOOUser();
      tooFlowPage = new TooFlowPage(officePage, move);
      await tooFlowPage.waitForLoading();
      await officePage.tooNavigateToMove(tooFlowPage.moveLocator);
    });

    test('is able to approve a shipment', async ({ page }) => {
      await expect(page.locator('#approved-shipments')).not.toBeVisible();
      await expect(page.locator('#requested-shipments')).toBeVisible();
      await expect(page.getByText('Approve selected')).toBeDisabled();
      await expect(page.locator('#approvalConfirmationModal [data-testid="modal"]')).not.toBeVisible();

      await tooFlowPage.waitForLoading();
      await tooFlowPage.approveAllShipments();

      // Redirected to Move Task Order page

      expect(page.url()).toContain(`/moves/${tooFlowPage.moveLocator}/mto`);
      await expect(page.getByTestId('ShipmentContainer')).toBeVisible();
      await expect(page.locator('[data-testid="ApprovedServiceItemsTable"] h3')).toContainText(
        'Approved Service Items (14 items)',
      );

      // MTO compliance information is visible
      await expect(
        page.getByText('Payment will be made using the Third-Party Payment System (TPPS) Automated Payment System'),
      ).toBeVisible();
      await expect(
        page.getByText(
          'Packaging, packing, and shipping instructions as identified in the Conformed Copy of HTC111-11-1-1112 Attachment 1 Performance Work Statement',
        ),
      ).toBeVisible();
      await expect(page.getByText('NAICS: 488510 - FREIGHT TRANSPORTATION ARRANGEMENT')).toBeVisible();
      await expect(page.getByText('Contract #HTC111-11-1-1112')).toBeVisible();

      // Navigate back to Move Details
      await page.getByTestId('MoveDetails-Tab').click();
      await tooFlowPage.waitForLoading();

      expect(page.url()).toContain(`/moves/${tooFlowPage.moveLocator}/details`);
      await expect(page.locator('#approvalConfirmationModal [data-testid="modal"]')).not.toBeVisible();
      await expect(page.locator('#approved-shipments')).toBeVisible();
      await expect(page.locator('#requested-shipments')).not.toBeVisible();
      await expect(page.getByText('Approve selected')).not.toBeVisible();
    });

    test('is able to flag and unflag a move for financial review', async ({ page }) => {
      expect(page.url()).toContain('details');

      // click to trigger financial review modal
      await page.getByText('Flag move for financial review').click();

      // Enter information in modal and submit
      await page.locator('label').getByText('Yes', { exact: true }).click();
      await page.locator('textarea').fill('Something is rotten in the state of Denmark');

      // Click save on the modal
      await page.getByRole('button', { name: 'Save' }).click();

      // Verify sucess alert and tag
      await expect(page.getByText('Move flagged for financial review.')).toBeVisible();
      await expect(page.getByText('Flagged for financial review', { exact: true })).toBeVisible();

      // now test unflag
      expect(page.url()).toContain('details');

      // click to trigger financial review modal
      await page.getByText('Edit', { exact: true }).click();

      // Enter information in modal and submit
      await page.locator('label').getByText('No', { exact: true }).click();

      // Click save on the modal
      await page.getByRole('button', { name: 'Save' }).click();

      // Verify success alert and tag
      await expect(page.getByText('Move unflagged for financial review.')).toBeVisible();
    });

    test('is able to approve and reject mto service items', async ({ page }) => {
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

      // This test requires at least two requested service items
      await expect(page.getByText('Requested Service Items', { exact: false })).toBeVisible();
      await expect(getServiceItemsInTable(requestedServiceItemsTable).nth(1)).toBeVisible();

      await expect(page.getByTestId('modal')).not.toBeVisible();

      // Approve a requested service item
      expect((await getServiceItemsInTable(requestedServiceItemsTable).count()) > 0);
      await requestedServiceItemsTable.getByRole('button', { name: 'Accept' }).first().click();
      await tooFlowPage.waitForLoading();

      await expect(page.locator('[data-testid="ApprovedServiceItemsTable"] h3')).toContainText(
        'Approved Service Items (15 items)',
      );
      await expect(getServiceItemsInTable(approvedServiceItemsTable)).toHaveCount(approvedServiceItemCount + 1);
      approvedServiceItemCount = await getServiceItemsInTable(approvedServiceItemsTable).count();

      await expect(getServiceItemsInTable(requestedServiceItemsTable)).toHaveCount(requestedServiceItemCount - 1);
      requestedServiceItemCount = await getServiceItemsInTable(requestedServiceItemsTable).count();

      // Reject a requested service item
      await expect(page.getByText('Requested Service Items', { exact: false })).toBeVisible();
      expect((await getServiceItemsInTable(requestedServiceItemsTable).count()) > 0);
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

      // Accept a previously rejected service item
      await rejectedServiceItemsTable.getByRole('button').first().click();

      await expect(page.locator('[data-testid="ApprovedServiceItemsTable"] h3')).toContainText(
        'Approved Service Items (15 items)',
      );
      await expect(getServiceItemsInTable(approvedServiceItemsTable)).toHaveCount(approvedServiceItemCount + 1);
      approvedServiceItemCount = await getServiceItemsInTable(approvedServiceItemsTable).count();

      await expect(getServiceItemsInTable(rejectedServiceItemsTable)).toHaveCount(rejectedServiceItemCount - 1);
      rejectedServiceItemCount = await getServiceItemsInTable(rejectedServiceItemsTable).count();

      // Reject a previously accepted service item
      await approvedServiceItemsTable.getByRole('button').first().click();

      await expect(page.getByTestId('modal')).toBeVisible();
      modal = page.getByTestId('modal');
      await expect(modal.getByRole('button', { name: 'Submit' })).toBeDisabled();
      await modal.getByTestId('textInput').fill('changed my mind about this one');
      await modal.getByRole('button', { name: 'Submit' }).click();

      await expect(page.getByTestId('modal')).not.toBeVisible();

      await expect(page.getByText('Rejected Service Items', { exact: false })).toBeVisible();
      await expect(getServiceItemsInTable(rejectedServiceItemsTable)).toHaveCount(rejectedServiceItemCount + 1);
      rejectedServiceItemCount = await getServiceItemsInTable(rejectedServiceItemsTable).count();

      await expect(page.locator('[data-testid="ApprovedServiceItemsTable"] h3')).toContainText(
        'Approved Service Items (15 items)',
      );
      await expect(getServiceItemsInTable(approvedServiceItemsTable)).toHaveCount(approvedServiceItemCount - 1);
      approvedServiceItemCount = await getServiceItemsInTable(approvedServiceItemsTable).count();
    });

    test('is able to edit orders', async ({ page }) => {
      // Navigate to Edit orders page
      await expect(page.getByTestId('edit-orders')).toContainText('Edit orders');
      await page.getByText('Edit orders').click();
      await tooFlowPage.waitForLoading();

      // Check for department indicators
      await page.getByLabel('Department indicator').selectOption(DEPARTMENT_INDICATOR_OPTIONS.AIR_AND_SPACE_FORCE);
      await page.getByLabel('Department indicator').selectOption(DEPARTMENT_INDICATOR_OPTIONS.ARMY);
      await page.getByLabel('Department indicator').selectOption(DEPARTMENT_INDICATOR_OPTIONS.ARMY_CORPS_OF_ENGINEERS);
      await page.getByLabel('Department indicator').selectOption(DEPARTMENT_INDICATOR_OPTIONS.COAST_GUARD);
      await page.getByLabel('Department indicator').selectOption(DEPARTMENT_INDICATOR_OPTIONS.NAVY_AND_MARINES);
      await page
        .getByLabel('Department indicator')
        .selectOption(DEPARTMENT_INDICATOR_OPTIONS.OFFICE_OF_SECRETARY_OF_DEFENSE);

      // Toggle between Edit Allowances and Edit Orders page
      await page.getByTestId('view-allowances').click();
      await tooFlowPage.waitForLoading();
      expect(page.url()).toContain(`/moves/${tooFlowPage.moveLocator}/allowances`);
      await page.getByTestId('view-orders').click();
      await tooFlowPage.waitForLoading();
      expect(page.url()).toContain(`/moves/${tooFlowPage.moveLocator}/orders`);

      // Check for link that allows TIO to download the PDF for copy/paste functionality
      await expect(page.locator('p[class*="DocumentViewer_downloadLink"] > a > span')).toHaveText('Download file');

      // Edit orders fields

      await tooFlowPage.selectDutyLocation('Fort Irwin', 'originDutyLocation');
      // select the 5th option in the dropdown
      await tooFlowPage.selectDutyLocation('JB McGuire-Dix-Lakehurst', 'newDutyLocation', 5);

      await page.locator('input[name="issueDate"]').clear();
      await page.locator('input[name="issueDate"]').fill('16 Mar 2018');
      await page.locator('input[name="reportByDate"]').clear();
      await page.locator('input[name="reportByDate"]').fill('22 Mar 2018');
      await page.locator('select[name="departmentIndicator"]').selectOption({ label: '21 Army' });
      await page.locator('input[name="ordersNumber"]').clear();
      await page.locator('input[name="ordersNumber"]').fill('ORDER66');
      await page.locator('select[name="ordersType"]').selectOption({ label: 'Permanent Change Of Station (PCS)' });
      await page.locator('select[name="ordersTypeDetail"]').selectOption({ label: 'Shipment of HHG Permitted' });
      await page.locator('input[name="tac"]').clear();
      await page.locator('input[name="tac"]').fill('F123');
      await page.locator('input[name="sac"]').clear();
      await page.locator('input[name="sac"]').fill('4K988AS098F');

      // Edit orders page | Save
      await page.getByRole('button', { name: 'Save' }).click();
      await page.getByRole('heading', { name: 'Move details' }).waitFor();

      // Verify edited values are saved
      expect(page.url()).toContain(`/moves/${tooFlowPage.moveLocator}/details`);

      await expect(page.locator('[data-testid="currentDutyLocation"]')).toContainText('Fort Irwin');
      await expect(page.locator('[data-testid="newDutyLocation"]')).toContainText(
        'JB Langley-Eustis (Fort Eustis), VA 23604',
      );
      await expect(page.locator('[data-testid="issuedDate"]')).toContainText('16 Mar 2018');
      await expect(page.locator('[data-testid="reportByDate"]')).toContainText('22 Mar 2018');
      await expect(page.locator('[data-testid="departmentIndicator"]')).toContainText('Army');
      await expect(page.locator('[data-testid="ordersNumber"]')).toContainText('ORDER66');
      await expect(page.locator('[data-testid="ordersType"]')).toContainText('Permanent Change Of Station (PCS)');
      await expect(page.locator('[data-testid="ordersTypeDetail"]')).toContainText('Shipment of HHG Permitted');
      await expect(page.locator('[data-testid="tacMDC"]')).toContainText('F123');
      await expect(page.locator('[data-testid="sacSDN"]')).toContainText('4K988AS098F');

      // Edit orders page | Cancel
      // Navigate to Edit orders page
      await expect(page.getByTestId('edit-orders')).toContainText('Edit orders');
      await page.getByText('Edit orders').click();
      await tooFlowPage.waitForLoading();
      await page.getByRole('button', { name: 'Cancel' }).click();
      await tooFlowPage.waitForLoading();

      expect(page.url()).toContain(`/moves/${tooFlowPage.moveLocator}/details`);
    });

    test('is able to request cancellation for a shipment', async ({ page }) => {
      await tooFlowPage.waitForLoading();
      await tooFlowPage.approveAllShipments();

      await page.getByTestId('MoveTaskOrder-Tab').click();
      expect(page.url()).toContain(`/moves/${tooFlowPage.moveLocator}/mto`);

      // Move Task Order page
      await expect(page.getByTestId('ShipmentContainer')).toHaveCount(1);

      // Click requestCancellation button and display modal
      await page.locator('.shipment-heading').locator('button').getByText('Request Cancellation').click();

      await expect(page.getByTestId('modal')).toBeVisible();
      const modal = page.getByTestId('modal');

      await modal.locator('button[type="submit"]').click();

      // After updating, the button is disabeld and an alert is shown
      await expect(page.getByTestId('modal')).not.toBeVisible();
      await expect(page.locator('.shipment-heading')).toContainText('Cancellation Requested');
      await expect(
        page
          .locator('[data-testid="alert"]')
          .getByText('The request to cancel that shipment has been sent to the movers.'),
      ).toBeVisible();

      // Alert should disappear if focus changes
      await page.locator('[data-testid="rejectTextButton"]').first().click();
      await page.locator('[data-testid="closeRejectServiceItem"]').click();
      await expect(page.locator('[data-testid="alert"]')).not.toBeVisible();
    });

    /**
     * This test is being temporarily skipped until flakiness issues
     * can be resolved. It was skipped in cypress and is not part of
     * the initial playwright conversion. - ahobson 2023-01-10
     */
    test.skip('is able to edit allowances', async ({ page }) => {
      // Navigate to Edit allowances page
      await expect(page.getByTestId('edit-allowances')).toContainText('Edit allowances');
      await page.getByText('Edit allowances').click();

      // // Toggle between Edit Allowances and Edit Orders page
      // await page.locator('[data-testid="view-orders"]').click();
      // cy.url().should('include', `/moves/${moveLocator}/orders`);
      // await page.locator('[data-testid="view-allowances"]').click();
      // cy.url().should('include', `/moves/${moveLocator}/allowances`);

      // await page.locator('form').within(($form) => {
      //   // Edit pro-gear, pro-gear spouse, RME, SIT, and OCIE fields
      //   await page.locator('input[name="proGearWeight"]').fill('1999');
      //   await page.locator('input[name="proGearWeightSpouse"]').fill('499');
      //   await page.locator('input[name="requiredMedicalEquipmentWeight"]').fill('999');
      //   await page.locator('input[name="storageInTransit"]').fill('199');
      //   await page.locator('input[name="organizationalClothingAndIndividualEquipment"]').siblings('label[for="ocieInput"]').click();

      //   // Edit grade and authorized weight
      //   await expect(page.locator('select[name=agency]')).toContainText('Army');
      //   await page.locator('select[name=agency]').selectOption({ label: 'Navy'});
      //   await expect(page.locator('select[name="grade"]')).toContainText('E-1');
      //   await page.locator('select[name="grade"]').selectOption({ label: 'W-2'});
      //   await page.locator('input[name="authorizedWeight"]').fill('11111');

      //   //Edit DependentsAuthorized
      //   await page.locator('input[name="dependentsAuthorized"]').siblings('label[for="dependentsAuthorizedInput"]').click();

      //   // Edit allowances page | Save
      //   await expect(page.locator('button').contains('Save')).toBeEnabled().click();

      // cy.wait(['@patchAllowances']);

      // // Verify edited values are saved
      // cy.url().should('include', `/moves/${moveLocator}/details`);

      // await expect(page.locator('[data-testid="progear"]')).toContainText('1,999');
      // await expect(page.locator('[data-testid="spouseProgear"]')).toContainText('499');
      // await expect(page.locator('[data-testid="rme"]')).toContainText('999');
      // await expect(page.locator('[data-testid="storageInTransit"]')).toContainText('199');
      // await expect(page.locator('[data-testid="ocie"]')).toContainText('Unauthorized');

      // await expect(page.locator('[data-testid="authorizedWeight"]')).toContainText('11,111');
      // await expect(page.locator('[data-testid="branchGrade"]')).toContainText('Navy');
      // await expect(page.locator('[data-testid="branchGrade"]')).toContainText('W-2');
      // await expect(page.locator('[data-testid="dependents"]')).toContainText('Unauthorized');

      // // Edit allowances page | Cancel
      // await expect(page.locator('[data-testid="edit-allowances"]')).toContainText('Edit allowances').click();
      // await expect(page.locator('button')).toContainText('Cancel').click();
      // cy.url().should('include', `/moves/${moveLocator}/details`);
    });

    test('is able to edit shipment', async ({ page }) => {
      const deliveryDate = new Date().toLocaleDateString('en-US');

      // Edit the shipment
      await page.locator('[data-testid="ShipmentContainer"] .usa-button').first().click();
      // fill out some changes on the form
      await page.locator('#requestedDeliveryDate').clear();
      await page.locator('#requestedDeliveryDate').fill(deliveryDate);
      await page.locator('#requestedDeliveryDate').blur();
      await page.locator('input[name="delivery.address.streetAddress1"]').clear();
      await page.locator('input[name="delivery.address.streetAddress1"]').fill('7 q st');
      await page.locator('input[name="delivery.address.city"]').clear();
      await page.locator('input[name="delivery.address.city"]').fill('city');
      await page.locator('select[name="delivery.address.state"]').selectOption({ label: 'OH' });
      await page.locator('input[name="delivery.address.postalCode"]').clear();
      await page.locator('input[name="delivery.address.postalCode"]').fill('90210');
      await page.locator('[data-testid="submitForm"]').click();
      await expect(page.locator('[data-testid="submitForm"]')).not.toBeEnabled();

      await tooFlowPage.waitForPage.moveDetails();
    });

    // Test that the TOO is blocked from doing QAE actions
    test('is unable to see create report buttons', async ({ page }) => {
      await page.getByText('Quality assurance').click();
      await tooFlowPage.waitForLoading();
      await expect(page.getByText('Quality assurance reports')).toBeVisible();

      // Make sure there are no create report buttons on the page
      await expect(page.getByText('Create report')).not.toBeVisible();
    });

    test('cannot load evaluation report creation form', async ({ page }) => {
      // Attempt to visit edit page for an evaluation report (report ID doesn't matter since
      // we should get stopped before looking it up)
      await page.goto(`/moves/${tooFlowPage.moveLocator}/evaluation-reports/11111111-1111-1111-1111-111111111111`);
      await expect(page.getByText("Sorry, you can't access this page")).toBeVisible();
      await page.getByText('Go to move details').click();
      await tooFlowPage.waitForLoading();

      // Make sure we go to move details page
      expect(page.url()).toContain(`/moves/${tooFlowPage.moveLocator}/details`);
    });

    test('is able to view Origin GBLOC', async ({ page }) => {
      // Check for Origin GBLOC label
      await expect(page.getByTestId('originGBLOC')).toHaveText('Origin GBLOC');
      await expect(page.getByTestId('infoBlock')).toContainText('KKFA');
    });
  });

  test.describe('with HHG Moves after actual pickup date', () => {
    test.beforeEach(async ({ officePage }) => {
      const move = await officePage.testHarness.buildHHGMoveForTOOAfterActualPickupDate();
      await officePage.signInAsNewTOOUser();
      tooFlowPage = new TooFlowPage(officePage, move);
      await tooFlowPage.waitForLoading();
      await officePage.tooNavigateToMove(tooFlowPage.moveLocator);
    });

    test('is able to request diversion for a shipment and receive alert msg', async ({ page }) => {
      await tooFlowPage.waitForLoading();
      await tooFlowPage.approveAllShipments();

      await page.getByTestId('MoveTaskOrder-Tab').click();
      expect(page.url()).toContain(`/moves/${tooFlowPage.moveLocator}/mto`);

      // Move Task Order page
      await expect(page.getByTestId('ShipmentContainer')).toHaveCount(1);

      await page.locator('button').getByText('Request diversion').click();

      // Check modal title includes shipment locator
      const modalTitleText = await page.locator('div[data-testid="modal"] h3').textContent();
      const modalTitlePattern = /^Request Shipment Diversion for #([A-Za-z0-9]{6}-\d{2})$/;
      const hasValidModalTitle = modalTitlePattern.test(modalTitleText);
      expect(hasValidModalTitle).toBeTruthy();

      // Submit the diversion request
      await page.locator('input[name="diversionReason"]').fill('reasonable reason');
      await page.locator('button[data-testid="modalSubmitButton"]').click();
      await expect(page.locator('.shipment-heading')).toContainText('diversion requested');

      // Check the alert message with shipment locator
      const alertText = await page.locator('[data-testid="alert"]').textContent();
      const shipmentNumberPattern = /^Diversion successfully requested for Shipment #([A-Za-z0-9]{6}-\d{2})$/;
      const hasValidShipmentNumber = shipmentNumberPattern.test(alertText);
      expect(hasValidShipmentNumber).toBeTruthy();

      // Alert should disappear if focus changes
      await page.locator('[data-testid="rejectTextButton"]').first().click();
      await page.locator('[data-testid="closeRejectServiceItem"]').click();
      await expect(page.locator('[data-testid="alert"]')).not.toBeVisible();
    });
  });

  test.describe('with retiree moves', () => {
    test.beforeEach(async ({ officePage }) => {
      const move = await officePage.testHarness.buildHHGMoveWithRetireeForTOO();
      await officePage.signInAsNewTOOUser();

      tooFlowPage = new TooFlowPage(officePage, move);
      await tooFlowPage.waitForLoading();
      await officePage.tooNavigateToMove(move.locator);
    });

    test('is able to edit shipment for retiree', async ({ page }) => {
      const deliveryDate = new Date().toLocaleDateString('en-US');

      // Edit the shipment
      await page.locator('[data-testid="ShipmentContainer"] .usa-button').first().click();
      // fill out some changes on the form
      await page.locator('#requestedDeliveryDate').clear();
      await page.locator('#requestedDeliveryDate').fill(deliveryDate);
      await page.locator('#requestedDeliveryDate').blur();

      await page.locator('input[name="delivery.address.streetAddress1"]').clear();
      await page.locator('input[name="delivery.address.streetAddress1"]').fill('7 q st');
      await page.locator('input[name="delivery.address.city"]').clear();
      await page.locator('input[name="delivery.address.city"]').fill('city');
      await page.locator('select[name="delivery.address.state"]').selectOption({ label: 'OH' });
      await page.locator('input[name="delivery.address.postalCode"]').clear();
      await page.locator('input[name="delivery.address.postalCode"]').fill('90210');
      await page.locator('select[name="destinationType"]').selectOption({ label: 'Home of selection (HOS)' });

      await page.locator('[data-testid="submitForm"]').click();

      await tooFlowPage.waitForPage.moveDetails();
    });
  });

  let moveLoc;
  test.describe('with payment requests', () => {
    test.beforeEach(async ({ officePage, page }) => {
      const move = await officePage.testHarness.buildHHGMoveInSITEndsToday();
      moveLoc = move.locator;
      await officePage.signInAsNewMultiRoleUser();

      await page.getByRole('link', { name: 'Change user role' }).click();
      await page.getByRole('button', { name: 'Select prime_simulator' }).click();
      await page.locator('#moveCode').click();
      await page.locator('#moveCode').fill(moveLoc);
      await page.locator('#moveCode').press('Enter');
      await page.getByTestId('moveCode-0').click();
      await page.getByRole('link', { name: 'Create Payment Request' }).click();
      await page.waitForSelector('h3:has-text("Domestic origin SIT fuel surcharge")');
      const serviceItemID = await page.$eval(
        `//h3[text()='Domestic origin SIT fuel surcharge']/following-sibling::div[contains(@class, 'descriptionList_row__TsTvp')]//dt[text()='ID:']/following-sibling::dd[1]`,
        (ddElement) => ddElement.textContent.trim(),
      );
      await page.locator(`label[for="${serviceItemID}"]`).nth(0).check();
      await page.locator(`input[name="params\\.${serviceItemID}\\.WeightBilled"]`).fill('10000');
      await page.locator(`input[name="params\\.${serviceItemID}\\.WeightBilled"]`).blur();
      await page.getByTestId('form').getByTestId('button').click();
      await page.getByRole('link', { name: 'Change user role' }).click();
      await page.getByRole('button', { name: 'Select task_ordering_officer' }).click();
    });
    test('weight-based multiplier prioritizes billed weight', async ({ page }) => {
      await page.getByRole('row', { name: 'Select...' }).getByTestId('locator').getByTestId('TextBoxFilter').click();
      await page
        .getByRole('row', { name: 'Select...' })
        .getByTestId('locator')
        .getByTestId('TextBoxFilter')
        .fill(moveLoc);
      await page
        .getByRole('row', { name: 'Select...' })
        .getByTestId('locator')
        .getByTestId('TextBoxFilter')
        .press('Enter');
      await page.getByTestId('locator-0').click();
      await page.getByRole('link', { name: 'Payment requests' }).click();
      await page.getByRole('button', { name: 'Review shipment weights' }).click();
      await page.getByRole('button', { name: 'Review shipment weights' }).click();
      await page.getByRole('button', { name: 'Back' }).click();
      await page.getByRole('link', { name: 'Payment requests' }).click();
      await page.getByTestId('reviewBtn').click();
      await page.getByTestId('toggleCalculations').click();
      await expect(page.getByText('Weight-based distance multiplier: 0.0006255')).toBeVisible();
    });
  });

  test('approves a delivery address change request for an HHG shipment', async ({ officePage, page }) => {
    const shipmentAddressUpdate = await officePage.testHarness.bulidHHGMoveWithAddressChangeRequest();
    await officePage.signInAsNewTOOUser();
    tooFlowPage = new TooFlowPage(officePage, shipmentAddressUpdate.Shipment.MoveTaskOrder);
    await tooFlowPage.waitForLoading();
    await officePage.tooNavigateToMove(shipmentAddressUpdate.Shipment.MoveTaskOrder.locator);

    await expect(page.getByText('Review required')).toBeVisible();

    // Edit the shipment
    await page.getByRole('button', { name: 'Edit shipment' }).click();

    await expect(
      page.getByTestId('alert').getByText('Request needs review. See delivery location to proceed.'),
    ).toBeVisible();
    await expect(
      page
        .getByTestId('alert')
        .getByText('Pending delivery location change request needs review. Review request to proceed.'),
    ).toBeVisible();

    // click to trigger review modal
    await page.getByRole('button', { name: 'Review request' }).click();

    // Enter information in modal and submit
    await page.getByTestId('modal').getByTestId('radio').getByText('Yes').click();
    await page.getByTestId('modal').locator('textarea').fill('The delivery address change looks good. ');

    // Click save on the modal
    await page.getByTestId('modal').getByRole('button', { name: 'Save' }).click();
    await expect(page.getByTestId('modal')).not.toBeVisible();

    await expect(page.getByText('Changes sent to contractor.')).toBeVisible();

    const destinationAddress = page.getByRole('group', { name: 'Delivery location' });
    await expect(destinationAddress.getByLabel('Address 1')).toHaveValue('123 Any Street');
    await expect(destinationAddress.getByLabel('Address 2')).toHaveValue('P.O. Box 12345');
    await expect(destinationAddress.getByLabel('City')).toHaveValue('Beverly Hills');
    await expect(destinationAddress.getByLabel('State')).toHaveValue('CA');
    await expect(destinationAddress.getByLabel('ZIP')).toHaveValue('90210');

    // Click save on the page
    await page.getByRole('button', { name: 'Save' }).click();

    await expect(page.getByText('Update request details')).not.toBeVisible();
    await expect(page.getByText('Review required')).not.toBeVisible();
    await expect(page.getByTestId('destinationAddress')).toContainText(
      '123 Any Street, P.O. Box 12345, c/o Some Person, Beverly Hills, CA 90210',
    );

    await page.getByText('KKFA moves').click();

    await page.locator('input[name="locator"]').fill(shipmentAddressUpdate.Shipment.MoveTaskOrder.locator);
    await page.locator('input[name="locator"]').blur();
    // once the move is in the Move approved status, it will no longer show up in the TOO queue
    await expect(page.getByText('Move approved')).not.toBeVisible();
    await expect(page.getByText('Approvals requested')).not.toBeVisible();
  });

  test('approves a delivery address change request for a NTSr shipment', async ({ officePage, page }) => {
    const shipmentAddressUpdate = await officePage.testHarness.buildNTSRMoveWithAddressChangeRequest();
    await officePage.signInAsNewTOOUser();
    tooFlowPage = new TooFlowPage(officePage, shipmentAddressUpdate.Shipment.MoveTaskOrder);
    await tooFlowPage.waitForLoading();
    await officePage.tooNavigateToMove(shipmentAddressUpdate.Shipment.MoveTaskOrder.locator);

    await expect(page.getByText('Review required')).toBeVisible();
    await page.getByRole('button', { name: 'Edit shipment' }).click();

    await expect(
      page.getByTestId('alert').getByText('Request needs review. See delivery location to proceed.'),
    ).toBeVisible();
    await expect(
      page
        .getByTestId('alert')
        .getByText('Pending delivery location change request needs review. Review request to proceed.'),
    ).toBeVisible();
    await page.getByRole('button', { name: 'Review request' }).click();

    await page.getByTestId('modal').getByTestId('radio').getByText('Yes').click();
    await page.getByTestId('modal').locator('textarea').fill('The delivery address change looks good. ');
    await page.getByTestId('modal').getByRole('button', { name: 'Save' }).click();
    await expect(page.getByTestId('modal')).not.toBeVisible();
    await expect(page.getByText('Changes sent to contractor.')).toBeVisible();

    const destinationAddress = page.getByRole('group', { name: 'Delivery location' });
    await expect(destinationAddress.getByLabel('Address 1')).toHaveValue('123 Any Street');
    await expect(destinationAddress.getByLabel('Address 2')).toHaveValue('P.O. Box 12345');
    await expect(destinationAddress.getByLabel('City')).toHaveValue('Beverly Hills');
    await expect(destinationAddress.getByLabel('State')).toHaveValue('CA');
    await expect(destinationAddress.getByLabel('ZIP')).toHaveValue('90210');

    // Save the approved delivery address change
    await page.getByRole('button', { name: 'Save' }).click();

    await expect(page.getByText('Update request details')).not.toBeVisible();
    await expect(page.getByText('Review required')).not.toBeVisible();
    await expect(page.getByTestId('destinationAddress')).toContainText(
      '123 Any Street, P.O. Box 12345, c/o Some Person, Beverly Hills, CA 90210',
    );

    // go back and make sure the move is in approved status (won't be viewable in TOO queue)
    await page.getByText('KKFA moves').click();
    await page.locator('input[name="locator"]').fill(shipmentAddressUpdate.Shipment.MoveTaskOrder.locator);
    await page.locator('input[name="locator"]').blur();
    await expect(page.getByText('Move approved')).not.toBeVisible();
    await expect(page.getByText('Approvals requested')).not.toBeVisible();
  });
});
