import { DEPARTMENT_INDICATOR_OPTIONS } from '../../utils/office/officeTest';
import { appendTimestampToFilenamePrefix, formatDate } from '../../utils/playwrightUtility';

import { test, expect } from './servicesCounselingTestFixture';

const createCustomerFF = process.env.FEATURE_FLAG_COUNSELOR_MOVE_CREATE;
const alaskaFF = process.env.FEATURE_FLAG_ENABLE_ALASKA;
const LocationLookup1 = 'BEVERLY HILLS, CA 90210 (LOS ANGELES)';
const LocationLookup2 = 'BEVERLY HILLS, CA 90212 (LOS ANGELES)';

test.describe('Services counselor user', () => {
  test.describe('Can create a customer with an international Alaska move', () => {
    let pickupDate;
    let pickupDateStr;
    let deliveryDate;
    let deliveryDateStr;

    test.beforeEach(async ({ scPage }) => {
      await scPage.signInAsNewServicesCounselorUser();
    });

    test.skip(
      createCustomerFF === 'false' || alaskaFF === 'false',
      'Skip if the create customer & AK FFs are not enabled.',
    );

    test('create a customer and add a basic iHHG shipment with Alaska address', async ({ page, officePage }) => {
      // make sure we see the queue
      await expect(page.getByText('Moves')).toBeVisible();
      await expect(page.getByRole('link', { name: 'Counseling' })).toBeVisible();
      await expect(page.getByRole('link', { name: 'Customer Search' })).toBeVisible();

      // we need to search before we have access to the create customer button
      await page.getByRole('link', { name: 'Customer Search' }).click();
      await page.getByText('Customer Name').click();
      await page.getByLabel('Search').fill('Test');
      await page.getByRole('button', { name: 'Search' }).click();
      await expect(page.getByRole('button', { name: 'Add Customer' })).toBeEnabled();
      await page.getByRole('button', { name: 'Add Customer' }).click();

      // fill out the customer form
      await page.getByRole('combobox', { name: 'Branch of service' }).selectOption({ label: 'Army' });
      await page.getByLabel('DoD ID number').fill('1234567890');
      await page.getByLabel('First name').nth(0).fill('Mister');
      await page.getByLabel('Last name').nth(0).fill('Alaska');
      await page.getByLabel('Best contact phone').fill('555-555-5555');
      await page.getByLabel('Personal email').fill('alaskaBoi@mail.mil');
      await page.getByText('Phone', { exact: true }).nth(0).click();
      await page.getByLabel('Address 1').nth(0).fill('1234 Pickup St.');
      await page.getByLabel('Location Lookup').nth(0).fill('90210');
      await expect(page.getByText(LocationLookup1, { exact: true })).toBeVisible();
      await page.keyboard.press('Enter');
      await page.getByLabel('Address 1').nth(1).fill('1234 Backup St.');
      await page.getByLabel('Location Lookup').nth(1).fill('90212');
      await expect(page.getByText(LocationLookup2, { exact: true })).toBeVisible();
      await page.keyboard.press('Enter');
      await page.getByLabel('First Name').nth(1).fill('Backup');
      await page.getByLabel('Last Name').nth(1).fill('Friend');
      await page.getByLabel('Email', { exact: true }).nth(1).fill('backupFriend@mail.mil');
      await page.getByLabel('Phone', { exact: true }).nth(1).fill('555-867-5309');
      await page.locator('label[for="noCreateOktaAccount"]').click();
      await page.locator('label[for="yesCacUser"]').click();
      await page.keyboard.press('Tab');
      await expect(page.getByRole('button', { name: 'Save' })).toBeEnabled();
      await page.getByRole('button', { name: 'Save' }).click();

      // fill out the orders form
      await page.getByLabel('Orders type').selectOption({ label: 'Permanent Change Of Station (PCS)' });
      await page.getByLabel('Orders date').fill('12/25/2024');
      await page.getByLabel('Orders date').blur();
      await page.getByLabel('Report by date').fill('1/25/2025');
      await page.getByLabel('Report by date').blur();
      const originLocation = 'Tinker AFB, OK 73145';
      await page.getByLabel('Current duty location').fill(originLocation);
      await expect(page.getByText(originLocation, { exact: true })).toBeVisible();
      await page.keyboard.press('Enter');
      await page.getByLabel('Counseling office').selectOption({ label: 'PPPO Tinker AFB - USAF' });
      const counselingOffice = page.locator('#counselingOfficeId');
      await counselingOffice.selectOption('PPPO Tinker AFB - USAF');
      const dutyLocation = 'Elmendorf AFB, AK 99506';
      await page.getByLabel('New duty location').fill('ELMENDORF AFB');
      await expect(page.getByText(dutyLocation, { exact: true })).toBeVisible();
      await page.keyboard.press('Enter');
      await page.locator('label[for="hasDependentsNo"]').click();
      await page.getByLabel('Pay grade').selectOption({ label: 'E-7' });
      await expect(page.getByRole('button', { name: 'Next' })).toBeEnabled();
      await page.getByRole('button', { name: 'Next' }).click();

      // now we need to add order data
      await page.getByRole('link', { name: 'View and edit orders' }).click();
      const filepondContainer = page.locator('.filepond--wrapper');
      await officePage.uploadFileViaFilepond(filepondContainer, 'AF Orders Sample.pdf');
      await expect(page.getByTestId('documentAlertMessage')).toContainText('Uploading');
      await expect(page.getByTestId('documentAlertMessage')).not.toBeVisible();
      await expect(page.getByText('Upload complete')).not.toBeVisible();
      await expect(
        page.getByTestId('uploads-table').getByText(appendTimestampToFilenamePrefix('AF Orders Sample')),
      ).toBeVisible();
      await page.getByRole('button', { name: 'Done' }).click();
      await page.getByLabel('Department indicator').selectOption(DEPARTMENT_INDICATOR_OPTIONS.ARMY);
      await page.getByLabel('Orders number').fill('123456');
      await page.getByLabel('Orders type detail').selectOption('Shipment of HHG Permitted');
      await page.getByTestId('hhgTacInput').fill('TEST');
      await expect(page.getByRole('button', { name: 'Save' })).toBeEnabled();
      await page.getByRole('button', { name: 'Save' }).click();

      // adding an HHG shipment
      await page.getByLabel('Add a new shipment').selectOption('HHG');
      await expect(page.getByText('Add shipment details')).toBeVisible();
      await expect(page.getByText('Weight allowance: 11,000 lbs')).toBeVisible();
      pickupDate = new Date();
      pickupDate.setDate(pickupDate.getDate() + 5);
      pickupDateStr = formatDate(pickupDate);
      deliveryDate = new Date();
      deliveryDate.setDate(pickupDate.getDate() + 10);
      deliveryDateStr = formatDate(deliveryDate);
      await page.getByLabel('Requested pickup date').fill(pickupDateStr);
      await page.getByLabel('Requested pickup date').blur();
      await page.getByText('Use pickup address').click();
      await page.getByLabel('Requested delivery date').fill(deliveryDateStr);
      await page.getByLabel('Requested delivery date').blur();
      await expect(page.getByRole('button', { name: 'Save' })).toBeEnabled();
      await page.getByRole('button', { name: 'Save' }).click();

      // verify we can see the iHHG shipment, submit it to the TOO
      await expect(page.getByText('iHHG')).toBeVisible();
      await expect(page.getByRole('button', { name: 'Submit move details' })).toBeEnabled();
      await page.getByRole('button', { name: 'Submit move details' }).click();
      await expect(page.getByText('Are you sure?')).toBeVisible();
      await page.getByRole('button', { name: 'Yes, submit' }).click();
      await expect(page.getByText('Move submitted.')).toBeVisible();
    });

    test('can create PPM with existing iHHG shipment on move', async ({ page, scPage }) => {
      const move = await scPage.testHarness.buildIntlHHGMoveNeedsSC();
      await scPage.navigateToMove(move.locator);

      // adding a PPM shipment
      await page.getByLabel('Add a new shipment').selectOption('PPM');
      await expect(page.getByText('Add shipment details')).toBeVisible();
      await expect(page.getByText('Weight allowance: 8,000 lbs')).toBeVisible();
      await page.getByLabel('Planned Departure Date').fill('01 May 2025');
      await page.getByLabel('Planned Departure Date').blur();
      await page.getByText('Use pickup address').click();

      const deliveryLocation = 'FAIRBANKS, AK 99702 (FAIRBANKS NORTH STAR)';
      const deliveryAddress = page.getByRole('group', { name: 'Delivery Address' });
      await deliveryAddress.getByLabel('Address 1').nth(0).fill('123 Cold St.');
      await page.locator('input[id="destination.address-location-input"]').fill('99702');
      await expect(page.getByText(deliveryLocation, { exact: true })).toBeVisible();
      await page.keyboard.press('Enter');

      await scPage.selectDutyLocation('JPPSO NORTHWEST', 'closeoutOffice');

      await page.getByTestId('estimatedWeight').fill('2,000');

      await expect(page.getByRole('button', { name: 'Save and Continue' })).toBeEnabled();
      await page.getByRole('button', { name: 'Save and Continue' }).click();

      await expect(page.getByText('Incentive & advance')).toBeVisible();
      await expect(page.getByText('Estimated incentive')).toBeVisible();

      await page.getByTestId('counselor-remarks').fill('remarks');

      await expect(page.getByRole('button', { name: 'Save and Continue' })).toBeEnabled();
      await page.getByRole('button', { name: 'Save and Continue' }).click();

      // verify we can see the PPM shipment, submit it to the TOO
      await expect(page.getByText('iPPM')).toBeVisible();
      await expect(page.getByRole('button', { name: 'Submit move details' })).toBeEnabled();
      await page.getByRole('button', { name: 'Submit move details' }).click();
      await expect(page.getByText('Are you sure?')).toBeVisible();
      await page.getByRole('button', { name: 'Yes, submit' }).click();
      await expect(page.getByText('Move submitted.')).toBeVisible();
    });
  });
});
