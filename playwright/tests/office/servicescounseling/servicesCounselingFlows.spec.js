/**
 * Semi-automated converted from a cypress test, and thus may contain
 * non best-practices, in particular: heavy use of `page.locator`
 * instead of `page.getBy*`.
 */

// @ts-check
import { DEPARTMENT_INDICATOR_OPTIONS } from '../../utils/office/officeTest';

import { test, expect } from './servicesCounselingTestFixture';

const supportingDocsEnabled = process.env.FEATURE_FLAG_MANAGE_SUPPORTING_DOCS;
const LocationLookup = 'BEVERLY HILLS, CA 90210 (LOS ANGELES)';

test.describe('Services counselor user', () => {
  test.describe('with basic HHG move', () => {
    test.beforeEach(async ({ scPage }) => {
      const move = await scPage.testHarness.buildHHGMoveNeedsSC();
      await scPage.navigateToMove(move.locator);
    });

    test('is able to click on move and submit after using the move code filter', async ({ page }) => {
      /**
       * Move Details page
       */
      // click to trigger confirmation modal
      await page.getByText('Submit move details').click();

      // modal should pop up with text
      await expect(page.locator('h2').getByText('Are you sure?')).toBeVisible();
      await expect(page.locator('p').getByText('You can’t make changes after you submit the move.')).toBeVisible();

      // click submit
      await page.getByRole('button', { name: 'Yes, submit' }).click();

      // verify success alert
      await expect(page.getByText('Move submitted.')).toBeVisible();
    });

    test('is able to flag a move for financial review', async ({ page, scPage }) => {
      // click to trigger financial review modal
      await page.getByText('Flag move for financial review').click();

      // Enter information in modal and submit
      await page.locator('label').getByText('Yes').click();
      await page.locator('textarea').fill('Because I said so...');

      // Click save on the modal
      await page.getByRole('button', { name: 'Save' }).click();
      await scPage.waitForLoading();

      // Verify sucess alert and tag
      await expect(page.getByText('Move flagged for financial review.')).toBeVisible();
      await expect(page.getByText('Flagged for financial review', { exact: true })).toBeVisible();

      // test('is able to unflag a move for financial review')
      // combining test with above

      // click to trigger financial review modal
      await page.getByText('Edit', { exact: true }).click();

      // Enter information in modal and submit
      await page.locator('label').getByText('No').click();

      // Click save on the modal
      await page.getByRole('button', { name: 'Save' }).click();
      await scPage.waitForLoading();

      // Verify sucess alert and tag
      await expect(page.getByText('Move unflagged for financial review.', { exact: true })).toBeVisible();
    });

    test('is able to edit a shipment', async ({ page, scPage }) => {
      await page.locator('[data-testid="ShipmentContainer"] .usa-button').first().click();
      await page.locator('#requestedPickupDate').clear();
      await page.locator('#requestedPickupDate').fill('16 Mar 2022');
      await page.locator('#requestedPickupDate').blur();
      await page.getByText('Use current address').click();

      await page.locator('#requestedDeliveryDate').clear();
      await page.locator('#requestedDeliveryDate').fill('16 May 2022');
      await page.locator('#requestedDeliveryDate').blur();

      await page.getByRole('group', { name: 'Delivery location' }).getByText('Yes').nth(1).click();
      await page.locator('input[name="delivery.address.streetAddress1"]').clear();
      await page.locator('input[name="delivery.address.streetAddress1"]').fill('7 q st');
      await page.locator('input[id="delivery.address-location-input"]').fill('90210');
      await expect(page.getByText(LocationLookup, { exact: true })).toBeVisible();
      await page.keyboard.press('Enter');

      // Select that we do not know the destination address yet
      await page.getByRole('group', { name: 'Delivery location' }).getByText('No').nth(1).click();
      await expect(page.getByText('We can use the zip of their new duty location:')).toBeVisible();

      await page.locator('[data-testid="submitForm"]').click();
      await scPage.waitForLoading();

      await expect(page.locator('.usa-alert__text')).toContainText('Your changes were saved.');
    });
    test('is able to view Origin GBLOC', async ({ page }) => {
      // Check for Origin GBLOC label
      await expect(page.getByTestId('originGBLOC')).toHaveText('Origin GBLOC');
      await expect(page.getByTestId('infoBlock')).toContainText('KKFA');
    });
  });

  test.describe('with HHG Move with Marine Corps as BOS', () => {
    test.beforeEach(async ({ scPage }) => {
      const move = await scPage.testHarness.buildHHGMoveAsUSMCNeedsSC();
      await scPage.navigateToMove(move.locator);
    });

    test('is able to view USMC as Origin GBLOC', async ({ page }) => {
      // Check for Origin GBLOC label
      await expect(page.getByTestId('originGBLOC')).toHaveText('Origin GBLOC');
      await expect(page.getByTestId('infoBlock')).toContainText('KKFA / USMC');
    });
  });

  test.describe('with HHG Move with amended orders', () => {
    test.beforeEach(async ({ scPage }) => {
      const move = await scPage.testHarness.buildHHGWithAmendedOrders();
      await scPage.navigateToMove(move.locator);
    });

    test('is able to view orders and amended orders', async ({ page }) => {
      await page.getByRole('link', { name: 'View and edit orders' }).click();
      await page.getByTestId('openMenu').click();
      await expect(page.getByTestId('DocViewerMenu').getByTestId('button')).toHaveCount(3);

      // Check for link that allows counselor to download the PDF for copy/paste functionality
      await expect(page.locator('p[class*="DocumentViewer_downloadLink"] > a > span')).toHaveText('Download file');

      // Check for department indicators
      await page.getByLabel('Department indicator').selectOption(DEPARTMENT_INDICATOR_OPTIONS.AIR_AND_SPACE_FORCE);
      await page.getByLabel('Department indicator').selectOption(DEPARTMENT_INDICATOR_OPTIONS.ARMY);
      await page.getByLabel('Department indicator').selectOption(DEPARTMENT_INDICATOR_OPTIONS.ARMY_CORPS_OF_ENGINEERS);
      await page.getByLabel('Department indicator').selectOption(DEPARTMENT_INDICATOR_OPTIONS.COAST_GUARD);
      await page.getByLabel('Department indicator').selectOption(DEPARTMENT_INDICATOR_OPTIONS.NAVY_AND_MARINES);
      await page
        .getByLabel('Department indicator')
        .selectOption(DEPARTMENT_INDICATOR_OPTIONS.OFFICE_OF_SECRETARY_OF_DEFENSE);
    });

    test('is able to add and delete orders and amended orders', async ({ page, officePage }) => {
      await page.getByRole('link', { name: 'View and edit orders' }).click();

      // check initial quanity of files
      await page.getByTestId('openMenu').click();
      await expect(page.getByTestId('DocViewerMenu').getByTestId('button')).toHaveCount(3);
      await page.getByTestId('closeMenu').click();

      // add orders
      await page.getByRole('button', { name: 'Manage Orders' }).click();
      const filepondContainer = page.locator('.filepond--wrapper');
      await officePage.uploadFileViaFilepond(filepondContainer, 'AF Orders Sample.pdf');
      await expect(page.getByText('Uploading')).toBeVisible();
      await expect(page.getByText('Uploading')).not.toBeVisible();
      await expect(page.getByText('Upload complete')).not.toBeVisible();
      await expect(page.getByTestId('uploads-table').getByText('AF Orders Sample.pdf')).toBeVisible();
      await page.getByTestId('openMenu').click();
      await expect(page.getByTestId('DocViewerMenu').getByTestId('button')).toHaveCount(4);
      await page.getByTestId('closeMenu').click();

      // delete orders
      const firstDeleteButton = page.locator('text=Delete').nth(0);
      await expect(firstDeleteButton).toBeVisible();
      await firstDeleteButton.click();
      await page.getByTestId('confirm-delete').click();
      await expect(page.getByText('Yes, delete')).not.toBeVisible();
      await expect(page.getByTestId('uploads-table').getByText('AF Orders Sample.pdf')).not.toBeVisible();
      await page.getByTestId('openMenu').click();
      await expect(page.getByTestId('DocViewerMenu').getByTestId('button')).toHaveCount(3);
      await page.getByTestId('closeMenu').click();
      await page.getByRole('button', { name: 'Manage Orders' }).click();

      // add amended orders
      await page.getByRole('button', { name: 'Manage Amended Orders' }).click();
      const filepondContainer2 = page.locator('.filepond--wrapper');
      await officePage.uploadFileViaFilepond(filepondContainer2, 'AF Orders Sample.pdf');
      await expect(page.getByText('Uploading')).toBeVisible();
      await expect(page.getByText('Uploading')).not.toBeVisible();
      await expect(page.getByText('Upload complete')).not.toBeVisible();
      await expect(page.getByTestId('uploads-table').getByText('AF Orders Sample.pdf')).toBeVisible();
      await page.getByTestId('openMenu').click();
      await expect(page.getByTestId('DocViewerMenu').getByTestId('button')).toHaveCount(4);
      await page.getByTestId('closeMenu').click();

      // delete amended orders
      const firstDeleteButtonAmended = page.locator('text=Delete').nth(0);
      await expect(firstDeleteButtonAmended).toBeVisible();
      await firstDeleteButtonAmended.click();
      await page.getByTestId('confirm-delete').click();
      await expect(page.getByText('Yes, delete')).not.toBeVisible();
      await expect(page.getByTestId('uploads-table').getByText('AF Orders Sample.pdf')).not.toBeVisible();
      await page.getByTestId('openMenu').click();
      await expect(page.getByTestId('DocViewerMenu').getByTestId('button')).toHaveCount(3);
      await page.getByTestId('closeMenu').click();
    });

    test('is able to add and delete supporting documents', async ({ page, officePage }) => {
      test.skip(supportingDocsEnabled === 'false', 'Skip if Supporting Documents is not enabled.');
      await page.getByRole('link', { name: 'Supporting Documents' }).click();
      await expect(page.getByText('No supporting documents have been uploaded.')).toBeVisible();

      // add orders
      const filepondContainer = page.locator('.filepond--wrapper');
      await officePage.uploadFileViaFilepond(filepondContainer, 'AF Orders Sample.pdf');
      await expect(page.getByText('Uploading')).toBeVisible();
      await expect(page.getByText('Uploading')).not.toBeVisible();
      await expect(page.getByText('Upload complete')).not.toBeVisible();
      await expect(page.getByTestId('uploads-table').getByText('AF Orders Sample.pdf')).toBeVisible();
      await expect(page.getByText('No supporting documents have been uploaded.')).not.toBeVisible();
      await page.getByTestId('openMenu').click();
      await expect(page.getByTestId('DocViewerMenu').getByTestId('button')).toHaveCount(1);
      await page.getByTestId('closeMenu').click();

      // delete orders
      const firstDeleteButton = page.locator('text=Delete').nth(0);
      await expect(firstDeleteButton).toBeVisible();
      await firstDeleteButton.click();
      await page.getByTestId('confirm-delete').click();
      await expect(page.getByText('Yes, delete')).not.toBeVisible();
      await expect(page.getByTestId('uploads-table').getByText('AF Orders Sample.pdf')).not.toBeVisible();
      await expect(page.getByText('No supporting documents have been uploaded.')).toBeVisible();
    });
  });

  test.describe('with separation HHG move', () => {
    test.beforeEach(async ({ scPage }) => {
      const move = await scPage.testHarness.buildHHGMoveForSeparationNeedsSC();
      await scPage.navigateToMove(move.locator);
    });

    test('is able to add a shipment', async ({ page, scPage }) => {
      const deliveryDate = new Date().toLocaleDateString('en-US');
      await expect(page.locator('[data-testid="ShipmentContainer"] .usa-button')).toHaveCount(2);

      // add a shipment
      await page.locator('[data-testid="dropdown"]').first().selectOption({ label: 'HHG' });

      await page.locator('#requestedPickupDate').fill(deliveryDate);
      await page.locator('#requestedPickupDate').blur();
      await page.getByText('Use current address').click();
      await page.locator('#requestedDeliveryDate').fill('16 Mar 2022');
      await page.locator('#requestedDeliveryDate').blur();
      await page.getByRole('group', { name: 'Delivery location' }).getByText('Yes').click();
      await page.locator('input[name="delivery.address.streetAddress1"]').fill('7 q st');
      await page.locator('input[id="delivery.address-location-input"]').fill('90210');
      await expect(page.getByText(LocationLookup, { exact: true })).toBeVisible();
      await page.keyboard.press('Enter');
      await page.locator('select[name="destinationType"]').selectOption({ label: 'Home of record (HOR)' });
      await page.locator('[data-testid="submitForm"]').click();
      await scPage.waitForLoading();

      // expect new shipment to show up
      await expect(page.locator('[data-testid="ShipmentContainer"] .usa-button')).toHaveCount(3);
    });
  });

  test.describe('with separation HHG move', () => {
    test.beforeEach(async ({ scPage }) => {
      const move = await scPage.testHarness.buildHHGMoveForRetireeNeedsSC();
      await scPage.navigateToMove(move.locator);
    });

    /**
     * This test is being temporarily skipped until flakiness issues
     * can be resolved. It was skipped in cypress and is not part of
     * the initial playwright conversion. - ahobson 2023-01-12
     */
    test.skip('is able to edit allowances', () => {
      //   // TOO Moves queue
      //   cy.wait(['@getSortedMoves']);
      //   await expect(page.getByText(moveLocator).click()).toBeVisible();
      //   cy.url().should('include', `/moves/${moveLocator}/details`);
      //   // Move Details page
      //   cy.watest(['@getMoves', '@getOrders', '@getMTOShipments', async ({page}) => {
      //   // Navigate to Edit allowances page
      //   await expect(page.locator('[data-testid="edit-allowances"]')).toContainText('Edit allowances').click();
      //   // Toggle between Edit Allowances and Edit Orders page
      //   await page.locator('[data-testid="view-orders"]').click();
      //   cy.url().should('include', `/moves/${moveLocator}/orders`);
      //   await page.locator('[data-testid="view-allowances"]').click();
      //   cy.url().should('include', `/moves/${moveLocator}/allowances`);
      //   cy.watest(['@getMoves', async ({page}) => {
      //   await page.locator('form').within(($form) => {
      //     // Edit pro-gear, pro-gear spouse, RME, SIT, and OCIE fields
      //     await page.locator('input[name="proGearWeight"]').fill('1999');
      //     await page.locator('input[name="proGearWeightSpouse"]').fill('499');
      //     await page.locator('input[name="requiredMedicalEquipmentWeight"]').fill('999');
      //     await page.locator('input[name="storageInTransit"]').fill('199');
      //     await page.locator('input[name="organizationalClothingAndIndividualEquipment"]').siblings('label[for="ocieInput"]').click();
      //     // Edit grade and authorized weight
      //     await expect(page.locator('select[name=agency]')).toContainText('Army');
      //     await page.locator('select[name=agency]').selectOption({ label: 'Navy'});
      //     await expect(page.locator('select[name="grade"]')).toContainText('E-1');
      //     await page.locator('select[name="grade"]').selectOption({ label: 'W-2'});
      //     //Edit DependentsAuthorized
      //     await page.locator('input[name="dependentsAuthorized"]').siblings('label[for="dependentsAuthorizedInput"]').click();
      //     // Edit allowances page | Save
      //     await expect(page.locator('[data-testid="scAllowancesSave"]')).toBeEnabled().click();
      //   });
      //   cy.wait(['@patchAllowances']);
      //   // Verify edited values are saved
      //   cy.url().should('include', `/moves/${moveLocator}/details`);
      //   cy.watest(['@getMoves', '@getOrders', '@getMTOShipments', async ({page}) => {
      //   await expect(page.locator('[data-testid="progear"]')).toContainText('1,999');
      //   await expect(page.locator('[data-testid="spouseProgear"]')).toContainText('499');
      //   await expect(page.locator('[data-testid="rme"]')).toContainText('999');
      //   await expect(page.locator('[data-testid="storageInTransit"]')).toContainText('199');
      //   await expect(page.locator('[data-testid="ocie"]')).toContainText('Unauthorized');
      //   await expect(page.locator('[data-testid="branchGrade"]')).toContainText('Navy');
      //   await expect(page.locator('[data-testid="branchGrade"]')).toContainText('W-2');
      //   await expect(page.locator('[data-testid="dependents"]')).toContainText('Unauthorized');
      //   // Edit allowances page | Cancel
      //   await expect(page.locator('[data-testid="edit-allowances"]')).toContainText('Edit allowances').click();
      //   await expect(page.locator('button')).toContainText('Cancel').click();
      //   cy.url().should('include', `/moves/${moveLocator}/details`);
    });

    test('is able to see and use the left navigation', async ({ page }) => {
      await expect(page.locator('a[href*="#shipments"]')).toContainText('Shipments');
      await expect(page.locator('a[href*="#orders"]')).toContainText('Orders');
      await expect(page.locator('a[href*="#allowances"]')).toContainText('Allowances');
      await expect(page.locator('a[href*="#customer-info"]')).toContainText('Customer info');
      // one warning in red for the missing destinationType
      await expect(page.locator('[data-testid="shipment-missing-info-alert"]')).toContainText('1');

      // Assert that the window has scrolled after clicking a left nav item
      const origScrollY = await page.evaluate(() => window.scrollY);
      await page.locator('#customer-info').click();
      const newScrollY = await page.evaluate(() => window.scrollY);
      expect(newScrollY).toBeGreaterThan(origScrollY);
    });

    test('is able to edit a shipment', async ({ page, scPage }) => {
      await page.locator('[data-testid="ShipmentContainer"] .usa-button').first().click();
      await page.locator('#requestedPickupDate').clear();
      await page.locator('#requestedPickupDate').fill('16 Mar 2022');
      await page.locator('#requestedPickupDate').blur();
      await page.getByText('Use current address').click();

      await page.locator('#requestedDeliveryDate').clear();
      await page.locator('#requestedDeliveryDate').fill('16 May 2022');
      await page.locator('#requestedDeliveryDate').blur();
      await page.getByRole('group', { name: 'Delivery location' }).getByText('Yes').nth(1).click();
      await page.locator('input[name="delivery.address.streetAddress1"]').clear();
      await page.locator('input[name="delivery.address.streetAddress1"]').fill('7 q st');
      await page.locator('input[id="delivery.address-location-input"]').fill('90210');
      await expect(page.getByText(LocationLookup, { exact: true })).toBeVisible();
      await page.keyboard.press('Enter');
      await page.locator('select[name="destinationType"]').selectOption({ label: 'Home of selection (HOS)' });
      await page.locator('[data-testid="submitForm"]').click();
      await scPage.waitForLoading();

      await expect(page.locator('.usa-alert__text')).toContainText('Your changes were saved.');
    });

    test('is able to update destination type if destination address is unknown', async ({ page, scPage }) => {
      await page.locator('[data-testid="ShipmentContainer"] .usa-button').first().click();
      await page.locator('#requestedPickupDate').clear();
      await page.locator('#requestedPickupDate').fill('16 Mar 2022');
      await page.locator('#requestedPickupDate').blur();
      await page.getByText('Use current address').click();

      await page.locator('#requestedDeliveryDate').clear();
      await page.locator('#requestedDeliveryDate').fill('16 May 2022');
      await page.locator('#requestedDeliveryDate').blur();

      // Select that we do not know the destination address yet
      await page.getByRole('group', { name: 'Delivery location' }).getByText('No').nth(1).click();

      await expect(page.locator('select[name="destinationType"]')).toBeVisible();
      await expect(page.getByText('We can use the zip of their HOR, HOS or PLEAD:')).toBeVisible();
      await page.locator('select[name="destinationType"]').selectOption({ label: 'Home of selection (HOS)' });
      await page.locator('[data-testid="submitForm"]').click();
      await scPage.waitForLoading();

      await expect(page.locator('.usa-alert__text')).toContainText('Your changes were saved.');
    });

    test('is able to see that the tag next to shipment is updated', async ({ page, scPage }) => {
      // Verify that there's a tag on the left nav that flags missing information
      await expect(page.locator('[data-testid="shipment-missing-info-alert"]')).toContainText('1');

      // Edit the shipment so that the tag disappears
      await page.locator('[data-testid="ShipmentContainer"] .usa-button').last().click();
      await page.locator('select[name="destinationType"]').selectOption({ label: 'Home of selection (HOS)' });
      await page.locator('[data-testid="submitForm"]').click();
      await scPage.waitForLoading();

      await expect(page.locator('.usa-alert__text')).toContainText('Your changes were saved.');

      // Verify that the tag after the update is blank since missing information was filled
      await expect(page.locator('[data-testid="shipment-missing-info-alert"]')).toHaveCount(0);
    });
  });

  test('can complete review of PPM shipment documents and view documents after', async ({ page, scPage }) => {
    const move = await scPage.testHarness.buildApprovedMoveWithPPMAllDocTypesOffice();
    await scPage.navigateToCloseoutMove(move.locator);

    // Navigate to the "Review documents" page
    await expect(page.getByRole('button', { name: /Review documents/i })).toBeVisible();
    await page.getByRole('button', { name: 'Review documents' }).click();

    await scPage.waitForPage.reviewWeightTicket();
    await expect(page.getByLabel('Accept')).toBeVisible();
    await page.getByText('Accept').click();
    await page.getByRole('button', { name: 'Continue' }).click();

    await scPage.waitForPage.reviewProGear();
    await expect(page.getByLabel('Accept')).toBeVisible();
    await page.getByText('Accept').click();
    await page.getByRole('button', { name: 'Continue' }).click();

    await scPage.waitForPage.reviewExpenseTicket('Packing Materials', 1, 1);
    await expect(page.getByLabel('Accept')).toBeVisible();
    await page.getByText('Accept').click();
    await page.getByRole('button', { name: 'Continue' }).click();

    await scPage.waitForPage.reviewDocumentsConfirmation();

    await page.getByRole('button', { name: 'Confirm' }).click();
    await scPage.waitForPage.moveDetails();

    await expect(page.getByText('PACKET READY FOR DOWNLOAD')).toBeVisible();

    // Navigate to the "View documents" page
    await expect(page.getByRole('button', { name: /View documents/i })).toBeVisible();
    await page.getByRole('button', { name: 'View documents' }).click();

    await scPage.waitForPage.reviewWeightTicket();
    await expect(page.getByLabel('Accept')).toBeVisible();
    await page.getByLabel('Accept').isDisabled();
    await page.getByRole('button', { name: 'Continue' }).click();

    await scPage.waitForPage.reviewProGear();
    await expect(page.getByLabel('Accept')).toBeVisible();
    await page.getByLabel('Accept').isDisabled();
    await page.getByRole('button', { name: 'Continue' }).click();

    await scPage.waitForPage.reviewExpenseTicket('Packing Materials', 1, 1);
    await expect(page.getByLabel('Accept')).toBeVisible();
    await page.getByLabel('Accept').isDisabled();
  });

  test.describe('Edit shipment info and incentives', () => {
    test.beforeEach(async ({ scPage }) => {
      const move = await scPage.testHarness.buildApprovedMoveWithPPMAllDocTypesOffice();
      await scPage.navigateToCloseoutMove(move.locator);
    });

    test('is able to edit/save actual move start date', async ({ page, scPage }) => {
      // Navigate to the "Review documents" page
      await expect(page.getByRole('button', { name: /Review documents/i })).toBeVisible();
      await page.getByRole('button', { name: 'Review documents' }).click();

      await scPage.waitForPage.reviewWeightTicket();
      // Edit Actual Move Start Date
      await page.getByTestId('actualMoveDate').getByTestId('editTextButton').click();
      await page.waitForSelector('text="Edit Shipment Info"');
      await page.getByRole('button', { name: 'Save' }).click();
      await page.waitForSelector('text="Edit Shipment Info"', { state: 'hidden' });
      await expect(page.getByLabel('Accept')).toBeVisible();
      await page.getByLabel('Accept').dispatchEvent('click');
      await page.getByRole('button', { name: 'Continue' }).click();
    });

    test('is able to edit/save pickup address', async ({ page, scPage }) => {
      // Navigate to the "Review documents" page
      await expect(page.getByRole('button', { name: /Review documents/i })).toBeVisible();
      await page.getByRole('button', { name: 'Review documents' }).click();

      await scPage.waitForPage.reviewWeightTicket();
      // Edit Starting Address
      await page.getByTestId('pickupAddress').getByTestId('editTextButton').click();
      await page.waitForSelector('text="Edit Shipment Info"');
      await page.getByRole('button', { name: 'Save' }).click();
      await page.waitForSelector('text="Edit Shipment Info"', { state: 'hidden' });
      await expect(page.getByLabel('Accept')).toBeVisible();
      await page.getByLabel('Accept').dispatchEvent('click');
      await page.getByRole('button', { name: 'Continue' }).click();
    });

    test('is able to edit/save destination address', async ({ page, scPage }) => {
      // Navigate to the "Review documents" page
      await expect(page.getByRole('button', { name: /Review documents/i })).toBeVisible();
      await page.getByRole('button', { name: 'Review documents' }).click();

      await scPage.waitForPage.reviewWeightTicket();
      // Edit Ending Address
      await page.getByTestId('destinationAddress').getByTestId('editTextButton').click();
      await page.waitForSelector('text="Edit Shipment Info"');
      await page.getByRole('button', { name: 'Save' }).click();
      await page.waitForSelector('text="Edit Shipment Info"', { state: 'hidden' });
      await expect(page.getByLabel('Accept')).toBeVisible();
      await page.getByLabel('Accept').dispatchEvent('click');
      await page.getByRole('button', { name: 'Continue' }).click();
    });

    test('is able to edit/save advance received', async ({ page, scPage }) => {
      // Navigate to the "Review documents" page
      await expect(page.getByRole('button', { name: /Review documents/i })).toBeVisible();
      await page.getByRole('button', { name: 'Review documents' }).click();

      await scPage.waitForPage.reviewWeightTicket();
      await expect(page.getByLabel('Accept')).toBeVisible();
      await page.getByLabel('Accept').dispatchEvent('click');
      await page.getByRole('button', { name: 'Continue' }).click();

      await scPage.waitForPage.reviewProGear();
      await expect(page.getByLabel('Accept')).toBeVisible();
      await page.getByLabel('Accept').dispatchEvent('click');
      await page.getByRole('button', { name: 'Continue' }).click();

      await scPage.waitForPage.reviewExpenseTicket('Packing Materials', 1, 1);
      await expect(page.getByLabel('Accept')).toBeVisible();
      await page.getByLabel('Accept').dispatchEvent('click');
      await page.getByRole('button', { name: 'Continue' }).click();

      await scPage.waitForPage.reviewDocumentsConfirmation();
      await page.waitForSelector('text="Loading, Please Wait..."', { state: 'hidden' });
      await page.getByTestId('incentives').getByTestId('incentives-showRequestDetailsButton').click();
      await page.getByTestId('advanceReceived').getByTestId('editTextButton').click();
      await page.waitForSelector('text="Edit Incentives/Costs"');
      await page.getByRole('button', { name: 'Save' }).click();
      await page.waitForSelector('text="Edit Incentives/Costs"', { state: 'hidden' });
    });
  });

  test.describe('Checking for Partial/Full PPM functionality', () => {
    let partialPpmCloseoutLocator = '';
    let partialPpmCounselingLocator = '';
    let fullPpmMoveLocator = '';

    test('counselor can see partial PPM ready for closeout', async ({ page, scPage }) => {
      const partialPpmMoveCloseout = await scPage.testHarness.buildPartialPPMMoveReadyForCloseout();
      partialPpmCloseoutLocator = partialPpmMoveCloseout.locator;
      await scPage.searchForCloseoutMove(partialPpmCloseoutLocator);
      await expect(page.getByTestId('ppmType-0')).toContainText('Partial');
    });

    test('counselor can see partial PPM ready for counseling', async ({ page, scPage }) => {
      const partialPpmMoveCounseling = await scPage.testHarness.buildPartialPPMMoveReadyForCounseling();
      partialPpmCounselingLocator = partialPpmMoveCounseling.locator;
      await scPage.searchForMove(partialPpmCounselingLocator);
      await expect(page.getByTestId('locator-0')).toContainText(partialPpmCounselingLocator);
    });

    test('counselor can see full PPM ready for closeout', async ({ page, scPage }) => {
      const fullPpmMove = await scPage.testHarness.buildPPMMoveWithCloseout();
      fullPpmMoveLocator = fullPpmMove.locator;
      await scPage.searchForCloseoutMove(fullPpmMoveLocator);
      await expect(page.getByTestId('ppmType-0')).toContainText('Full');
    });
  });

  test.describe('Actual expense reimbursement tests', () => {
    test.describe('is able to view/edit actual expense reimbursement for non-civilian moves', () => {
      test('view/edit actual expense reimbursement - edit shipments page', async ({ page, scPage }) => {
        const move = await scPage.testHarness.buildSubmittedMoveWithPPMShipmentForSC();
        await scPage.navigateToMove(move.locator);

        await expect(page.getByTestId('payGrade')).toContainText('E-1');
        await expect(page.getByText('actual expense reimbursement')).not.toBeVisible();

        await page.getByRole('button', { name: 'Edit shipment' }).click();
        await expect(page.locator('h1').getByText('Edit shipment details')).toBeVisible();

        expect(await page.locator('[data-testid="actualExpenseReimbursementTag"]').count()).toBe(0);

        await page.getByText('Yes').first().click();
        await page.getByTestId('submitForm').click();
        await expect(page.getByTestId('actualExpenseReimbursementTag')).toContainText('Actual Expense Reimbursement');
        await page.getByText('Approve').click();
        await page.getByTestId('counselor-remarks').click();
        await page.getByTestId('counselor-remarks').fill('test');
        await page.getByTestId('submitForm').click();

        await expect(page.getByTestId('payGrade')).toContainText('E-1');
        await expect(page.getByTestId('ShipmentContainer').getByTestId('tag')).toContainText(
          'actual expense reimbursement',
        );

        await page.getByRole('button', { name: 'Edit shipment' }).click();
        await expect(page.locator('h1').getByText('Edit shipment details')).toBeVisible();
        await expect(page.getByTestId('actualExpenseReimbursementTag')).toContainText('Actual Expense Reimbursement');
      });

      test('view/edit actual expense reimbursement - PPM closeout review documents', async ({ page, scPage }) => {
        const move = await scPage.testHarness.buildApprovedMoveWithPPMProgearWeightTicketOffice();
        await scPage.navigateToMoveUsingMoveSearch(move.locator);

        await expect(page.getByTestId('payGrade')).toContainText('E-1');
        await expect(page.getByText('actual expense reimbursement')).not.toBeVisible();

        await page.getByText('Review documents').click();
        await expect(page.getByRole('heading', { name: 'View documents' })).toBeVisible();

        expect(await page.locator('[data-testid="tag"]').count()).toBe(0);
        await expect(page.locator('label').getByText('Actual Expense Reimbursement')).toBeVisible();
        await page.getByTestId('isActualExpenseReimbursement').getByTestId('editTextButton').click();

        await expect(page.getByText('Is this PPM an Actual Expense Reimbursement?')).toBeVisible();
        await page.getByTestId('modal').getByText('Yes').click();
        await page.getByTestId('modal').getByTestId('button').click();

        await expect(page.getByText('Is this PPM an Actual Expense Reimbursement?')).not.toBeVisible();
        await page.getByTestId('shipmentInfo-showRequestDetailsButton').click();
        expect(await page.locator('[data-testid="tag"]').count()).toBe(1);
        await page.getByText('Accept').click();
        await page.getByTestId('closeSidebar').click();
        await expect(page.getByRole('heading', { name: 'Move details' })).toBeVisible();
        await expect(page.getByText('actual expense reimbursement')).toBeVisible();
      });
    });

    test.describe('is unable to edit actual expense reimbursement for civilian moves', () => {
      test('cannot edit actual expense reimbursement - edit shipments page', async ({ page, scPage }) => {
        const move = await scPage.testHarness.buildSubmittedMoveWithPPMShipmentForSC();
        await scPage.navigateToMove(move.locator);

        await expect(page.getByText('actual expense reimbursement')).not.toBeVisible();
        await page.getByTestId('view-edit-orders').click();
        await page.getByTestId('payGradeInput').selectOption('AVIATION_CADET');
        await page.getByTestId('payGradeInput').selectOption('CIVILIAN_EMPLOYEE');
        await page.getByRole('button', { name: 'Save' }).click();

        await expect(page.getByTestId('payGrade')).toContainText('Civilian Employee');
        await expect(page.getByText('actual expense reimbursement')).toBeVisible();
        await page.getByRole('button', { name: 'Edit shipment' }).click();

        await expect(page.locator('h1').getByText('Edit shipment details')).toBeVisible();

        expect(await page.locator('[data-testid="isActualExpenseReimbursementYes"]').isDisabled()).toBe(true);
        expect(await page.locator('[data-testid="isActualExpenseReimbursementNo"]').isDisabled()).toBe(true);
      });

      test('cannot edit actual expense reimbursement - PPM closeout review documents', async ({ page, scPage }) => {
        const move = await scPage.testHarness.buildApprovedMoveWithPPMProgearWeightTicketOfficeCivilian();
        await scPage.navigateToMoveUsingMoveSearch(move.locator);

        await expect(page.getByTestId('payGrade')).toContainText('Civilian Employee');

        await page.getByText('Review documents').click();
        await expect(page.getByRole('heading', { name: 'View documents' })).toBeVisible();
        await expect(page.getByTestId('tag')).toContainText('actual expense reimbursement');

        await expect(page.locator('label').getByText('Actual Expense Reimbursement')).toBeVisible();
        expect(await page.getByTestId('isActualExpenseReimbursement').getByTestId('editTextButton').isDisabled()).toBe(
          true,
        );
      });
    });
  });

  test.describe('with approved HHG move sent to Prime', () => {
    test.beforeEach(async ({ scPage }) => {
      const move = await scPage.testHarness.buildHHGMoveInSIT();
      await scPage.navigateToMoveUsingMoveSearch(move.locator);
    });

    test('is unable to view/edit orders after MTO has been created(sent to prime)', async ({ page }) => {
      await expect(page.getByTestId('view-edit-orders')).toBeHidden();
      await expect(page.getByTestId('edit-allowances')).toBeHidden();
    });
  });
});
