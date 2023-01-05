// @ts-check
const { test, expect, OfficePage } = require('../../utils/officeTest');

/**
 * TooFlowPage test fixture
 *
 * The logic in TooFlowPage is only used in this file, so keep the
 * playwright test fixture in this file.
 * @extends OfficePage
 */
class TooFlowPage extends OfficePage {
  /**
   * @param {OfficePage} officePage
   * @param {Object} move
   * @override
   */
  constructor(officePage, move) {
    super(officePage.page, officePage.request);
    this.move = move;
    this.moveLocator = move.locator;
  }

  /**
   * select and approve all service items on the page
   */
  async selectAndApproveAllServiceItems() {
    // Select & approve items
    const checkboxes = this.page.locator('input[data-testid="shipment-display-checkbox"]');
    const boxCount = await checkboxes.count();
    expect(boxCount).toBeGreaterThan(0);
    for (let i = 0; i < boxCount; i += 1) {
      const id = await checkboxes.nth(i).getAttribute('id');
      const label = this.page.locator(`label[for="${id}"]`);
      await expect(label).toHaveCount(1);
      await label.click();
    }
  }

  /**
   * approve all shipments on the move
   */
  async approveAllShipments() {
    const shipmentCount = await this.page.locator('input[data-testid="shipment-display-checkbox"]').count();
    // Select & approve items
    this.selectAndApproveAllServiceItems();
    // Select additional service items
    await this.page.locator('label[for="shipmentManagementFee"]').click();
    await this.page.locator('label[for="counselingFee"]').click();
    // Open modal
    await this.page.getByText('Approve selected').click();

    const modal = this.page.locator('#approvalConfirmationModal [data-testid="modal"]');
    await expect(modal).toBeVisible();
    // Verify modal content
    await expect(modal.getByText('Preview and post move task order')).toBeVisible();
    await expect(this.page.locator('#approvalConfirmationModal [data-testid="ShipmentContainer"]')).toHaveCount(
      shipmentCount,
    );
    expect(modal.getByText('Approved service items for this move')).toBeVisible();
    const siTable = modal.getByText('Approved service items for this move').locator('..').locator('table');
    await expect(siTable).toContainText('Move management');
    await expect(siTable).toContainText('Counseling');

    // Click approve
    await modal.getByText('Approve and send').click();
    await expect(this.page.locator('#approvalConfirmationModal [data-testid="modal"]')).not.toBeVisible();
    await this.waitForLoading();
  }
}

test.describe('TOO user', () => {
  /** @type {TooFlowPage} */
  let tooFlowPage;
  test.describe('with HHG Moves', () => {
    test.beforeEach(async ({ officePage }) => {
      const move = await officePage.testHarness.buildHHGMoveWithServiceItemsAndPaymentRequestsAndFilesForTOO();
      await officePage.signInAsNewTOOUser();
      tooFlowPage = new TooFlowPage(officePage, move);
      await officePage.tooNavigateToMove(move.locator);
    });

    test('is able to approve a shipment', async ({ page }) => {
      await expect(page.locator('#approved-shipments')).not.toBeVisible();
      await expect(page.locator('#requested-shipments')).toBeVisible();
      await expect(page.getByText('Approve selected')).toBeDisabled();
      await expect(page.locator('#approvalConfirmationModal [data-testid="modal"]')).not.toBeVisible();

      await tooFlowPage.approveAllShipments();

      // Redirected to Move Task Order page

      expect(page.url()).toContain(`/moves/${tooFlowPage.moveLocator}/mto`);
      await expect(page.getByTestId('ShipmentContainer')).toBeVisible();
      await expect(page.locator('[data-testid="ApprovedServiceItemsTable"] h3')).toContainText(
        'Approved service items (12 items)',
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

    test('is able to flag and unflag a move for financial review', async ({ page }) => {
      expect(page.url()).toContain('details');

      // click to trigger financial review modal
      await page.getByText('Flag move for financial review').click();

      // Enter information in modal and submit
      await page.locator('label').getByText('Yes', { exact: true }).click();
      await page.locator('textarea').type('Something is rotten in the state of Denmark');

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
      await tooFlowPage.approveAllShipments();

      await page.getByTestId('MoveTaskOrder-Tab').click();
      await tooFlowPage.waitForLoading();
      expect(page.url()).toContain(`/moves/${tooFlowPage.moveLocator}/mto`);

      // Move Task Order page
      await expect(page.getByTestId('ShipmentContainer')).toHaveCount(1);

      await expect(page.getByText('Approved service items (12 items)')).toBeVisible();
      await expect(page.getByText('Rejected service items')).not.toBeVisible();

      await expect(page.getByTestId('modal')).not.toBeVisible();

      // Approve a requested service item
      let serviceItemsTable = page.getByTestId('RequestedServiceItemsTable');
      await expect(serviceItemsTable.locator('tbody tr')).toHaveCount(2);
      await serviceItemsTable.locator('.acceptButton').first().click();
      await tooFlowPage.waitForLoading();

      await expect(page.getByText('Approved service items (12 items)')).toBeVisible();
      await expect(page.locator('[data-testid="ApprovedServiceItemsTable"] tbody tr')).toHaveCount(13);

      // Reject a requested service item
      await expect(page.getByText('Requested service items (1 item)')).toBeVisible();
      serviceItemsTable = page.getByTestId('RequestedServiceItemsTable');
      await expect(serviceItemsTable.locator('tbody tr')).toHaveCount(1);
      await serviceItemsTable.locator('.rejectButton').first().click();

      expect(page.getByTestId('modal')).toBeVisible();
      let modal = page.getByTestId('modal');

      await expect(modal.locator('button[type="submit"]')).toBeDisabled();
      await modal.locator('[data-testid="textInput"]').type('my very valid reason');
      await modal.locator('button[type="submit"]').click();

      expect(page.getByTestId('modal')).not.toBeVisible();

      await expect(page.getByText('Rejected service items (1 item)')).toBeVisible();
      await expect(page.locator('[data-testid="RejectedServiceItemsTable"] tbody tr')).toHaveCount(1);

      // Accept a previously rejected service item
      await page.locator('[data-testid="RejectedServiceItemsTable"] button').click();

      await expect(page.getByText('Approved service items (13 items)')).toBeVisible();
      await expect(page.locator('[data-testid="ApprovedServiceItemsTable"] tbody tr')).toHaveCount(13);
      await expect(page.getByText('Rejected service items (1 item)')).not.toBeVisible();

      // Reject a previously accpeted service item
      await page.locator('[data-testid="ApprovedServiceItemsTable"] button').first().click();

      expect(page.getByTestId('modal')).toBeVisible();
      modal = page.getByTestId('modal');
      await expect(modal.locator('button[type="submit"]')).toBeDisabled();
      await modal.getByTestId('textInput').type('changed my mind about this one');
      await modal.locator('button[type="submit"]').click();

      expect(page.getByTestId('modal')).not.toBeVisible();

      await expect(page.getByText('Rejected service items (1 item)')).toBeVisible();
      await expect(page.locator('[data-testid="RejectedServiceItemsTable"] tbody tr')).toHaveCount(1);

      await expect(page.getByText('Requested service items')).not.toBeVisible();
      await expect(page.getByText('Approved service items (13 items)')).toBeVisible();
      await expect(page.locator('[data-testid="ApprovedServiceItemsTable"] tbody tr')).toHaveCount(13);
    });

    test('is able to edit orders', async ({ page }) => {
      // Navigate to Edit orders page
      await expect(page.getByTestId('edit-orders')).toContainText('Edit orders');
      await page.getByText('Edit orders').click();
      await tooFlowPage.waitForLoading();

      // Toggle between Edit Allowances and Edit Orders page
      await page.getByTestId('view-allowances').click();
      await tooFlowPage.waitForLoading();
      expect(page.url()).toContain(`/moves/${tooFlowPage.moveLocator}/allowances`);
      await page.getByTestId('view-orders').click();
      await tooFlowPage.waitForLoading();
      expect(page.url()).toContain(`/moves/${tooFlowPage.moveLocator}/orders`);

      // Edit orders fields

      await tooFlowPage.selectDutyLocation('Fort Irwin', 'originDutyLocation');
      // select the 5th option in the dropdown
      await tooFlowPage.selectDutyLocation('JB McGuire-Dix-Lakehurst', 'newDutyLocation', 5);

      await page.locator('input[name="issueDate"]').clear();
      await page.locator('input[name="issueDate"]').type('16 Mar 2018');
      await page.locator('input[name="reportByDate"]').clear();
      await page.locator('input[name="reportByDate"]').type('22 Mar 2018');
      await page.locator('select[name="departmentIndicator"]').selectOption({ label: '21 Army' });
      await page.locator('input[name="ordersNumber"]').clear();
      await page.locator('input[name="ordersNumber"]').type('ORDER66');
      await page.locator('select[name="ordersType"]').selectOption({ label: 'Permanent Change Of Station (PCS)' });
      await page.locator('select[name="ordersTypeDetail"]').selectOption({ label: 'Shipment of HHG Permitted' });
      await page.locator('input[name="tac"]').clear();
      await page.locator('input[name="tac"]').type('F123');
      await page.locator('input[name="sac"]').clear();
      await page.locator('input[name="sac"]').type('4K988AS098F');

      // Edit orders page | Save
      await Promise.all([page.waitForNavigation(), page.getByRole('button', { name: 'Save' }).click()]);
      await tooFlowPage.waitForLoading();

      // Verify edited values are saved
      expect(page.url()).toContain(`/moves/${tooFlowPage.moveLocator}/details`);

      await expect(page.locator('[data-testid="currentDutyLocation"]')).toContainText('Fort Irwin');
      await expect(page.locator('[data-testid="newDutyLocation"]')).toContainText(
        'Joint Base Lewis-McChord (McChord AFB)',
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
      await tooFlowPage.approveAllShipments();

      await page.getByTestId('MoveTaskOrder-Tab').click();
      expect(page.url()).toContain(`/moves/${tooFlowPage.moveLocator}/mto`);

      // Move Task Order page
      await expect(page.getByTestId('ShipmentContainer')).toHaveCount(1);

      // Click requestCancellation button and display modal
      await page.locator('.shipment-heading').locator('button').getByText('Request Cancellation').click();

      expect(page.getByTestId('modal')).toBeVisible();
      const modal = page.getByTestId('modal');

      await modal.locator('button[type="submit"]').click();

      // After updating, the button is disabeld and an alert is shown
      expect(page.getByTestId('modal')).not.toBeVisible();
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
      //   await page.locator('input[name="proGearWeight"]').type('1999');
      //   await page.locator('input[name="proGearWeightSpouse"]').type('499');
      //   await page.locator('input[name="requiredMedicalEquipmentWeight"]').type('999');
      //   await page.locator('input[name="storageInTransit"]').type('199');
      //   await page.locator('input[name="organizationalClothingAndIndividualEquipment"]').siblings('label[for="ocieInput"]').click();

      //   // Edit grade and authorized weight
      //   await expect(page.locator('select[name=agency]')).toContainText('Army');
      //   await page.locator('select[name=agency]').selectOption({ label: 'Navy'});
      //   await expect(page.locator('select[name="grade"]')).toContainText('E-1');
      //   await page.locator('select[name="grade"]').selectOption({ label: 'W-2'});
      //   await page.locator('input[name="authorizedWeight"]').type('11111');

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
      // await expect(page.locator('[data-testid="branchRank"]')).toContainText('Navy');
      // await expect(page.locator('[data-testid="branchRank"]')).toContainText('W-2');
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
      await page.locator('#requestedDeliveryDate').type(deliveryDate);
      await page.locator('#requestedDeliveryDate').blur();
      await page.locator('input[name="delivery.address.streetAddress1"]').clear();
      await page.locator('input[name="delivery.address.streetAddress1"]').type('7 q st');
      await page.locator('input[name="delivery.address.city"]').clear();
      await page.locator('input[name="delivery.address.city"]').type('city');
      await page.locator('select[name="delivery.address.state"]').selectOption({ label: 'OH' });
      await page.locator('input[name="delivery.address.postalCode"]').clear();
      await page.locator('input[name="delivery.address.postalCode"]').type('90210');
      await page.locator('[data-testid="submitForm"]').click();
      await tooFlowPage.waitForLoading();

      await expect(page.locator('[data-testid="submitForm"]')).not.toBeEnabled();
    });
    /**
     * This test is being temporarily skipped until flakiness issues
     * can be resolved. It was skipped in cypress and is not part of
     * the initial playwright conversion. - ahobson 2023-01-10
     */
    test.skip('is able to view SIT and create and edit SIT extensions', async ({ page }) => {
      await page.getByTestId('MoveTaskOrder-Tab').click();
      // cy.url().should('include', `/moves/${moveLocator}/mto`);

      // // View SIT display
      // await page.locator('[data-testid="sitExtensions"]');

      // // Total SIT
      // await expect(page.getByText('379 authorized')).toBeVisible();
      // // await expect(page.getByText('60 used')).toBeVisible();
      // // await expect(page.getByText('210 remaining')).toBeVisible();
      // // await expect(page.getByText('Ends 26 Apr 2022')).toBeVisible();

      // // Current SIT
      // await expect(page.getByText('Current location: destination')).toBeVisible();
      // await expect(page.getByText('60')).toBeVisible();
      // // await expect(page.getByText('29 Aug 2021')).toBeVisible();
      // // await expect(page.getByText('Ends 26 Apr 2022')).toBeVisible();

      // // Previous SIT
      // // await expect(page.getByText('30 days at origin (30 Jul 2021 - 29 Aug 2021)')).toBeVisible();
      // await expect(page.getByText('30 days at origin')).toBeVisible();

      // // SIT extensions
      // await expect(page.getByText('90 days added')).toBeVisible();
      // // await expect(page.getByText('on 28 Sep 2021')).toBeVisible();
      // await expect(page.getByText('Serious illness of the member')).toBeVisible();
      // await expect(page.getByText('The customer requested an extension.')).toBeVisible();
      // await expect(page.getByText('The service member is unable to move into their new home at the expected time')).toBeVisible();
    });

    // Test that the TOO is blocked from doing QAECSR actions
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
  });

  test.describe('with retiree moves', () => {
    test.beforeEach(async ({ officePage }) => {
      const move = await officePage.testHarness.buildHHGMoveWithRetireeForTOO();
      await officePage.signInAsNewTOOUser();

      tooFlowPage = new TooFlowPage(officePage, move);
      await officePage.tioNavigateToMove(move.locator);
    });

    test('is able to edit shipment for retiree', async ({ page }) => {
      const deliveryDate = new Date().toLocaleDateString('en-US');

      // Edit the shipment
      await page.locator('[data-testid="ShipmentContainer"] .usa-button').first().click();
      // fill out some changes on the form
      await page.locator('#requestedDeliveryDate').clear();
      await page.locator('#requestedDeliveryDate').type(deliveryDate);
      await page.locator('#requestedDeliveryDate').blur();

      await page.locator('input[name="delivery.address.streetAddress1"]').clear();
      await page.locator('input[name="delivery.address.streetAddress1"]').type('7 q st');
      await page.locator('input[name="delivery.address.city"]').clear();
      await page.locator('input[name="delivery.address.city"]').type('city');
      await page.locator('select[name="delivery.address.state"]').selectOption({ label: 'OH' });
      await page.locator('input[name="delivery.address.postalCode"]').clear();
      await page.locator('input[name="delivery.address.postalCode"]').type('90210');
      await page.locator('select[name="destinationType"]').selectOption({ label: 'Home of selection (HOS)' });
      await page.locator('[data-testid="submitForm"]').click();
      await tooFlowPage.waitForLoading();

      await expect(page.locator('[data-testid="submitForm"]')).not.toBeEnabled();
    });
  });
});
