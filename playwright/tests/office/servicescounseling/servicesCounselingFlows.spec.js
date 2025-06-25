/**
 * Semi-automated converted from a cypress test, and thus may contain
 * non best-practices, in particular: heavy use of `page.locator`
 * instead of `page.getBy*`.
 */

// @ts-check
import { DEPARTMENT_INDICATOR_OPTIONS } from '../../utils/office/officeTest';

import { test, expect } from './servicesCounselingTestFixture';

const completePPMCloseoutForCustomerEnabled = process.env.FEATURE_FLAG_COMPLETE_PPM_CLOSEOUT_FOR_CUSTOMER;
const LocationLookup = 'BEVERLY HILLS, CA 90210 (LOS ANGELES)';

test.describe('Services counselor user', () => {
  test.describe('GBLOC tests', () => {
    test.describe('Origin Duty Location', () => {
      let moveLocatorKKFA = '';
      let moveLocatorCNNQ = '';
      test.beforeEach(async ({ scPage }) => {
        const moveKKFA = await scPage.testHarness.buildHHGMoveNeedsSC();
        moveLocatorKKFA = moveKKFA.locator;
        const moveCNNQ = await scPage.testHarness.buildHHGMoveNeedsSC();
        moveLocatorCNNQ = moveCNNQ.locator;
      });

      test('when origin duty location GBLOC matches services counselor GBLOC', async ({ page }) => {
        const locatorFilter = await page.getByTestId('locator').getByTestId('TextBoxFilter');
        await locatorFilter.fill(moveLocatorKKFA);
        await locatorFilter.blur();
        await expect(page.getByTestId('locator-0')).toBeVisible();
      });

      test('when origin duty location GBLOC does not match services counselor GBLOC', async ({ page }) => {
        const locatorFilter = await page.getByTestId('locator').getByTestId('TextBoxFilter');
        await locatorFilter.fill(moveLocatorCNNQ);
        await locatorFilter.blur();
        await expect(page.getByTestId('locator-0')).not.toBeVisible();
      });
    });
  });

  test.describe('with basic HHG move', () => {
    test.beforeEach(async ({ scPage }) => {
      const move = await scPage.testHarness.buildHHGMoveNeedsSC();
      await scPage.navigateToMove(move.locator);
    });

    test('is able to click on move and submit after using the move code filter', async ({ page }) => {
      test.slow();
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
      test.slow();
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
      test.slow();
      await page.locator('[data-testid="ShipmentContainer"] .usa-button').first().click();
      await page.locator('#requestedPickupDate').clear();
      const pickupDate = new Date(Date.now() + 24 * 60 * 60 * 1000).toLocaleDateString('en-US');
      await page.getByLabel('Requested pickup date').fill(pickupDate);
      await page.locator('#requestedPickupDate').blur();
      await page.getByText('Use pickup address').click();

      await page.locator('#requestedDeliveryDate').clear();
      await page.locator('#requestedDeliveryDate').fill('16 May 2022');
      await page.locator('#requestedDeliveryDate').blur();

      await page.getByRole('group', { name: 'Delivery Address' }).getByText('Yes').nth(1).click();
      await page.locator('input[name="delivery.address.streetAddress1"]').clear();
      await page.locator('input[name="delivery.address.streetAddress1"]').fill('7 q st');
      const countrySearch = 'UNITED STATES';
      await page.locator('input[id="delivery.address-country-input"]').fill(countrySearch);
      let spanLocator = page.locator(`span:has(mark:has-text("${countrySearch}"))`);
      await expect(spanLocator).toBeVisible();
      await page.keyboard.press('Enter');
      const deliveryLocator = page.locator('input[id="delivery.address-input"]');
      await deliveryLocator.click({ timeout: 5000 });
      await deliveryLocator.fill('90210');
      await expect(page.getByText(LocationLookup, { exact: true })).toBeVisible();
      await page.keyboard.press('Enter');

      // Select that we do not know the delivery address yet
      await page.getByRole('group', { name: 'Delivery Address' }).getByText('No').nth(1).click();
      await expect(page.getByText('We can use the zip of their new duty location:')).toBeVisible();

      await page.locator('[data-testid="submitForm"]').click();
      await scPage.waitForLoading();

      await expect(page.locator('.usa-alert__text')).toContainText('Your changes were saved.');
    });
    test('is able to view Origin GBLOC', async ({ page }) => {
      test.slow();
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
      test.slow();
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
      test.slow();
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
      test.slow();
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
      test.slow();
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
      test.slow();
      const pickupDate = new Date(Date.now() + 24 * 60 * 60 * 1000).toLocaleDateString('en-US');
      const deliveryDate = new Date().toLocaleDateString('en-US');
      await expect(page.locator('[data-testid="ShipmentContainer"] .usa-button')).toHaveCount(2);

      // add a shipment
      await page.locator('[data-testid="dropdown"]').first().selectOption({ label: 'HHG' });

      await page.locator('#requestedPickupDate').fill(pickupDate);
      await page.locator('#requestedPickupDate').blur();
      await page.getByText('Use pickup address').click();
      await page.locator('#requestedDeliveryDate').fill(deliveryDate);
      await page.locator('#requestedDeliveryDate').blur();
      await page.getByRole('group', { name: 'Delivery Address' }).getByText('Yes').click();
      await page.locator('input[name="delivery.address.streetAddress1"]').fill('7 q st');
      const countrySearch = 'UNITED STATES';
      await page.locator('input[id="delivery.address-country-input"]').fill(countrySearch);
      let spanLocator = page.locator(`span:has(mark:has-text("${countrySearch}"))`);
      await expect(spanLocator).toBeVisible();
      await page.keyboard.press('Enter');
      const deliveryLocator = page.locator('input[id="delivery.address-input"]');
      await deliveryLocator.click({ timeout: 5000 });
      await deliveryLocator.fill('90210');
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

    test('is able to see and use the left navigation', async ({ page }) => {
      test.slow();
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
      test.slow();
      await page.locator('[data-testid="ShipmentContainer"] .usa-button').first().click();
      await page.locator('#requestedPickupDate').clear();
      const pickupDate = new Date(Date.now() + 24 * 60 * 60 * 1000).toLocaleDateString('en-US');
      await page.getByLabel('Requested pickup date').fill(pickupDate);
      await page.locator('#requestedPickupDate').blur();
      await page.getByText('Use pickup address').click();

      await page.locator('#requestedDeliveryDate').clear();
      await page.locator('#requestedDeliveryDate').fill('16 May 2022');
      await page.locator('#requestedDeliveryDate').blur();
      await page.getByRole('group', { name: 'Delivery Address' }).getByText('Yes').nth(1).click();
      await page.locator('input[name="delivery.address.streetAddress1"]').clear();
      await page.locator('input[name="delivery.address.streetAddress1"]').fill('7 q st');
      const countrySearch = 'UNITED STATES';
      await page.locator('input[id="delivery.address-country-input"]').fill(countrySearch);
      let spanLocator = page.locator(`span:has(mark:has-text("${countrySearch}"))`);
      await expect(spanLocator).toBeVisible();
      await page.keyboard.press('Enter');
      const deliveryLocator = page.locator('input[id="delivery.address-input"]');
      await deliveryLocator.click({ timeout: 5000 });
      await deliveryLocator.fill('90210');
      await expect(page.getByText(LocationLookup, { exact: true })).toBeVisible();
      await page.keyboard.press('Enter');
      await page.locator('select[name="destinationType"]').selectOption({ label: 'Home of selection (HOS)' });
      await page.getByLabel('Requested pickup date').fill(pickupDate);

      await page.locator('[data-testid="submitForm"]').click();
      await scPage.waitForLoading();

      await expect(page.locator('.usa-alert__text')).toContainText('Your changes were saved.');
    });

    test('is able to update destination type if delivery address is unknown', async ({ page, scPage }) => {
      test.slow();
      await page.locator('[data-testid="ShipmentContainer"] .usa-button').first().click();
      await page.locator('#requestedPickupDate').clear();
      const pickupDate = new Date(Date.now() + 24 * 60 * 60 * 1000).toLocaleDateString('en-US');
      await page.getByLabel('Requested pickup date').fill(pickupDate);
      await page.locator('#requestedPickupDate').blur();
      await page.getByText('Use pickup address').click();

      await page.locator('#requestedDeliveryDate').clear();
      await page.locator('#requestedDeliveryDate').fill('16 May 2022');
      await page.locator('#requestedDeliveryDate').blur();

      // Select that we do not know the delivery address yet
      await page.getByRole('group', { name: 'Delivery Address' }).getByText('No').nth(1).click();

      await expect(page.locator('select[name="destinationType"]')).toBeVisible();
      await expect(page.getByText('We can use the zip of their HOR, HOS or PLEAD:')).toBeVisible();
      await page.locator('select[name="destinationType"]').selectOption({ label: 'Home of selection (HOS)' });
      await page.locator('[data-testid="submitForm"]').click();
      await scPage.waitForLoading();

      await expect(page.locator('.usa-alert__text')).toContainText('Your changes were saved.');
    });

    test('is able to see that the tag next to shipment is updated', async ({ page, scPage }) => {
      test.slow();
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
    test.slow();
    await page.route('**/ghc/v1/ppm-shipments/*/payment-packet', async (route) => {
      // mocked blob
      const fakePdfBlob = new Blob(['%PDF-1.4 foo'], { type: 'application/pdf' });
      const arrayBuffer = await fakePdfBlob.arrayBuffer();
      // playwright route mocks only want a Buffer or string
      const bodyBuffer = Buffer.from(arrayBuffer);
      await route.fulfill({
        status: 200,
        contentType: 'application/pdf',
        body: bodyBuffer,
      });
    });
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

    await page.getByRole('button', { name: 'Preview PPM Payment Packet' }).click();
    await expect(page.getByTestId('loading-spinner')).not.toBeVisible();
    await page.getByRole('button', { name: 'Complete PPM Review' }).click();
    await page.getByRole('button', { name: 'Yes' }).click();
    await scPage.waitForPage.moveDetails();

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
      test.slow();
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
      test.slow();
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

    test('is able to edit/save delivery address', async ({ page, scPage }) => {
      test.slow();
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
      test.slow();
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
      test.slow();
      const partialPpmMoveCloseout = await scPage.testHarness.buildPartialPPMMoveReadyForCloseout();
      partialPpmCloseoutLocator = partialPpmMoveCloseout.locator;
      await scPage.searchForCloseoutMove(partialPpmCloseoutLocator);
      await expect(page.getByTestId('ppmType-0')).toContainText('Partial');
    });

    test('counselor can see partial PPM ready for counseling', async ({ page, scPage }) => {
      test.slow();
      const partialPpmMoveCounseling = await scPage.testHarness.buildPartialPPMMoveReadyForCounseling();
      partialPpmCounselingLocator = partialPpmMoveCounseling.locator;
      await scPage.searchForMove(partialPpmCounselingLocator);
      await expect(page.getByTestId('locator-0')).toContainText(partialPpmCounselingLocator);
    });

    test('counselor can see full PPM ready for closeout', async ({ page, scPage }) => {
      test.slow();
      const fullPpmMove = await scPage.testHarness.buildPPMMoveWithCloseout();
      fullPpmMoveLocator = fullPpmMove.locator;
      await scPage.searchForCloseoutMove(fullPpmMoveLocator);
      await expect(page.getByTestId('ppmType-0')).toContainText('Full');
    });

    test.describe('Complete PPM closeout on behalf of customer', () => {
      test.skip(completePPMCloseoutForCustomerEnabled === 'false', 'Skip if FF is disabled.');
      test('can complete PPM About page', async ({ page, scPage }) => {
        const move = await scPage.testHarness.buildApprovedMoveWithPPM();
        await scPage.navigateToMoveUsingMoveSearch(move.locator);

        await expect(page.getByRole('button', { name: /Complete PPM on behalf of the Customer/i })).toBeVisible();
        await page.getByRole('button', { name: 'Complete PPM on behalf of the Customer' }).click();

        // fill out About PPM page
        await expect(page.getByRole('heading', { name: 'About your PPM' })).toBeVisible();
        await expect(page.getByRole('heading', { name: 'How to complete your PPM' })).toBeVisible();
        await scPage.fillOutAboutPage();

        await expect(page.getByRole('heading', { name: 'How to complete your PPM' })).not.toBeVisible();
      });

      test('can navigate to PPM review page and edit About PPM page', async ({ page, scPage }) => {
        const move = await scPage.testHarness.buildMoveWithPPMShipmentReadyForFinalCloseout();
        await scPage.navigateToMoveUsingMoveSearch(move.locator);

        await expect(page.getByRole('button', { name: /Complete PPM on behalf of the Customer/i })).toBeVisible();
        await page.getByRole('button', { name: 'Complete PPM on behalf of the Customer' }).click();

        // Review page
        await expect(page.getByRole('heading', { name: 'Review' })).toBeVisible();
        await expect(page.getByRole('heading', { name: 'Documents' })).toBeVisible();
        await expect(page.getByRole('heading', { name: 'Weight moved' })).toBeVisible();
        await expect(page.getByRole('heading', { name: 'Pro-gear' })).toBeVisible();
        await expect(page.getByRole('heading', { name: 'Expenses' })).toBeVisible();

        // Edit About PPM page
        await page.locator('[data-testid="aboutYourPPM"] a').getByText('Edit').click();
        await expect(page.getByRole('heading', { name: 'About your PPM' })).toBeVisible();
        await page.getByRole('button', { name: 'Save & Continue' }).click();
        await expect(page.getByRole('heading', { name: 'Review' })).toBeVisible();
      });

      test('can add, edit, and delete Weight moved', async ({ page, scPage }) => {
        const move = await scPage.testHarness.buildApprovedMoveWithPPMWithAboutFormComplete();
        await scPage.navigateToMoveUsingMoveSearch(move.locator);

        await expect(page.getByRole('button', { name: /Complete PPM on behalf of the Customer/i })).toBeVisible();
        await page.getByRole('button', { name: 'Complete PPM on behalf of the Customer' }).click();

        // Add Weight Ticket
        await expect(page.getByRole('heading', { name: 'Review' })).toBeVisible();
        await page.getByText('Add More Weight').click();
        await expect(page.getByRole('heading', { name: 'Weight Tickets' })).toBeVisible();
        await scPage.fillOutWeightTicketPage({ hasTrailer: true, ownTrailer: true });
        await page.getByRole('button', { name: 'Save & Continue' }).click();

        // Edit
        await expect(page.getByRole('heading', { name: 'Review' })).toBeVisible();
        await page.getByTestId('weightMoved-1').click();
        await expect(page.getByRole('heading', { name: 'Weight Tickets' })).toBeVisible();
        await page.getByRole('button', { name: 'Save & Continue' }).click();

        // Delete
        await expect(page.getByRole('heading', { name: 'Review' })).toBeVisible();
        await page.getByTestId('weightMovedDelete-1').click();
        await expect(page.getByRole('heading', { name: 'Delete this?' })).toBeVisible();
        await page.getByRole('button', { name: 'Yes, Delete' }).click();
        await expect(page.getByText('Trip 1 successfully deleted.')).toBeVisible();
        await expect(page.getByTestId('weightMovedDelete-1')).not.toBeVisible();
      });
    });
  });

  test.describe('Actual expense reimbursement tests', () => {
    test.describe('is able to view/edit actual expense reimbursement for non-civilian moves', () => {
      test('view/edit actual expense reimbursement - edit shipments page', async ({ page, scPage }) => {
        test.slow();
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
        await expect(page.getByTestId('ShipmentContainer').getByTestId('actualReimbursementTag')).toContainText(
          'actual expense reimbursement',
        );

        await page.getByRole('button', { name: 'Edit shipment' }).click();
        await expect(page.locator('h1').getByText('Edit shipment details')).toBeVisible();
        await expect(page.getByTestId('actualExpenseReimbursementTag')).toContainText('Actual Expense Reimbursement');
      });

      test('view/edit actual expense reimbursement - PPM closeout review documents', async ({ page, scPage }) => {
        test.slow();
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
        await expect(page.getByRole('heading', { name: 'Move Details' })).toBeVisible();
        await expect(page.getByText('actual expense reimbursement')).toBeVisible();
      });
    });

    test.describe('is unable to edit actual expense reimbursement for civilian moves', () => {
      test('cannot edit actual expense reimbursement - edit shipments page', async ({ page, scPage }) => {
        test.slow();
        const move = await scPage.testHarness.buildSubmittedMoveWithPPMShipmentForSC();
        await scPage.navigateToMove(move.locator);

        await expect(page.getByText('actual expense reimbursement')).not.toBeVisible();
        await page.getByTestId('view-edit-orders').click();
        await page.getByTestId('payGradeInput').selectOption('CIVILIAN_EMPLOYEE');
        await page.getByRole('button', { name: 'Save' }).click();

        await expect(page.getByTestId('payGrade')).toContainText('CIVILIAN_EMPLOYEE');
        await expect(page.getByText('actual expense reimbursement')).toBeVisible();
        await page.getByRole('button', { name: 'Edit shipment' }).click();

        await expect(page.locator('h1').getByText('Edit shipment details')).toBeVisible();

        expect(await page.locator('[data-testid="isActualExpenseReimbursementYes"]').isDisabled()).toBe(true);
        expect(await page.locator('[data-testid="isActualExpenseReimbursementNo"]').isDisabled()).toBe(true);
      });

      test('cannot edit actual expense reimbursement - PPM closeout review documents', async ({ page, scPage }) => {
        test.slow();
        const move = await scPage.testHarness.buildApprovedMoveWithPPMProgearWeightTicketOfficeCivilian();
        await scPage.navigateToMoveUsingMoveSearch(move.locator);

        await expect(page.getByTestId('payGrade')).toContainText('CIVILIAN EMPLOYEE');

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
      test.slow();
      await expect(page.getByTestId('view-edit-orders')).toBeHidden();
      await expect(page.getByTestId('edit-allowances')).toBeHidden();
    });
  });
});
