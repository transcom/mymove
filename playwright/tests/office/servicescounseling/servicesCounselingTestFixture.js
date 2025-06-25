// @ts-check
import { expect, test as officeTest, OfficePage } from '../../utils/office/officeTest';

/**
 * ServiceCounselorPage test fixture
 *
 * @extends OfficePage
 */
export class ServiceCounselorPage extends OfficePage {
  /**
   * @param {OfficePage} officePage
   * @override
   */
  constructor(officePage) {
    super(officePage.page, officePage.request);
  }

  /**
   * Verify that the user is in the correct move
   * @param {string} moveLocator
   */
  async verifyMoveByLocatorCode(moveLocator) {
    await expect(this.page.getByTestId('nameBlock').getByText(`#${moveLocator}`)).toHaveClass(/usa-tag/);
  }

  /**
   * Finds a shipment container on the move details page from it's shipment type
   * @param {string} shipmentType
   * @returns async {import('@playwright/test').Locator}
   */
  async getShipmentContainerByType(shipmentType) {
    const header = await this.page.getByRole('heading', { level: 3, name: shipmentType });
    const container = await header.locator('../../../..');
    return container;
  }

  /**
   * Service Counselor searches for a move that needs counseling without clicking on it
   * @param {string} moveLocator
   */
  async searchForMove(moveLocator) {
    await this.waitForPage.counselingQueue();

    // Type in move code/locator to search for.
    // (There's no accessible or testId way to find this textbox, so we need to use .locator)
    await this.page.locator('input[name="locator"]').fill(moveLocator);
    await this.page.locator('input[name="locator"]').blur();
  }

  /**
   * Service Counselor searches for a closeout move without clicking on it
   * @param {string} moveLocator
   */
  async searchForCloseoutMove(moveLocator) {
    await this.waitForPage.counselingQueue();

    // Navigate to "PPM Closeout" section first
    await this.page.getByRole('link', { name: 'PPM Closeout' }).click();
    await this.waitForPage.closeoutQueue();

    // Type in move code/locator to search for.
    // (There's no accessible or testId way to find this textbox, so we need to use .locator)
    await this.page.locator('input[name="locator"]').fill(moveLocator);
    await this.page.locator('input[name="locator"]').blur();
  }

  /**
   * Service Counselor navigate to closeout move
   * @param {string} moveLocator
   */
  async navigateToCloseoutMove(moveLocator) {
    await this.waitForPage.counselingQueue();

    // Navigate to "PPM Closeout" section first
    await this.page.getByRole('link', { name: 'PPM Closeout' }).click();
    await this.waitForPage.closeoutQueue();

    // Type in move code/locator to search for.
    // (There's no accessible or testId way to find this textbox, so we need to use .locator)
    await this.page.locator('input[name="locator"]').fill(moveLocator);
    await this.page.locator('input[name="locator"]').blur();

    // Click the first returned row
    await this.page.getByTestId('locator-0').click();
    await this.waitForPage.moveDetails();
    await this.verifyMoveByLocatorCode(moveLocator);
  }

  /**
   * Service Counselor navigate to move search tab
   * @param {string} moveLocator
   */
  async navigateToMoveUsingMoveSearch(moveLocator) {
    await this.waitForPage.counselingQueue();

    // Navigate to "Move Search" tab
    await this.page.getByRole('link', { name: 'Move Search' }).click();
    await this.waitForPage.moveSearchTab();

    // Type in move code/locator to search for.
    // (There's no accessible or testId way to find this textbox, so we need to use .locator)
    await this.page.locator('input[name="searchText"]').type(moveLocator);
    await this.page.locator('input[name="searchText"]').blur();

    await this.page.getByTestId('searchTextSubmit').click();
    await this.waitForPage.moveSearchResults();
    await this.page.getByTestId('locator-0').click();
    await this.waitForPage.moveDetails();
    await this.verifyMoveByLocatorCode(moveLocator);
  }

  /**
   * Service Counselor navigate to move
   * @param {string} moveLocator
   */
  async navigateToMove(moveLocator) {
    await this.waitForPage.counselingQueue();

    // Type in move code/locator to search for.
    // (There's no accessible or testId way to find this textbox, so we need to use .locator)
    await this.page.locator('input[name="locator"]').fill(moveLocator);
    await this.page.locator('input[name="locator"]').blur();

    // Click the first returned row
    await this.page.getByTestId('locator-0').click();
    await this.waitForPage.moveDetails();
    await this.verifyMoveByLocatorCode(moveLocator);
  }

  async addNTSShipment() {
    await this.page.getByTestId('dropdown').selectOption({ label: 'NTS' });

    await this.waitForPage.addNTSShipment();
    const pickupDate = new Date(Date.now() + 24 * 60 * 60 * 1000).toLocaleDateString('en-US');
    await this.page.getByLabel('Requested pickup date').fill(pickupDate);
    await this.page.getByLabel('Requested pickup date').blur();
    await this.page.getByText('Use pickup address').click();

    await this.page.getByLabel('Counselor remarks').fill('Sample counselor remarks');

    // Save the shipment, progress back to the move details page, and verify it's been created
    await this.page.getByRole('button', { name: 'Save' }).click();
    await this.waitForPage.moveDetails();
  }

  async addNTSReleaseShipment() {
    await this.page.getByTestId('dropdown').selectOption({ label: 'NTS-release' });

    await this.waitForPage.addNTSReleaseShipment();

    await this.page.getByLabel('Previously recorded weight (lbs)').fill('1300');

    // Storage facility info
    const storageInfo = await this.page.getByRole('heading', { name: 'Storage facility info' }).locator('..');
    await storageInfo.getByLabel('Facility name').fill('Sample Facility Name');
    await storageInfo.getByLabel('Phone').fill('999-999-9999');
    await storageInfo.getByLabel('Email').fill('sample@example.com');
    await storageInfo.getByLabel('Service order number').fill('999999');

    // Storage facility address
    const StorageLocationLookup = 'ATLANTA, GA 30301 (FULTON)';
    const countrySearch = 'UNITED STATES';

    const storageAddress = this.page.getByRole('heading', { name: 'Storage facility address' }).locator('..');
    await storageAddress.getByLabel('Address 1').fill('148 S East St');
    await storageAddress.getByLabel('Address 2').fill('Suite 7A');
    await this.page.locator('input[id="storageFacility.address-country-input"]').fill(countrySearch);
    let spanLocator = this.page.locator(`span:has(mark:has-text("${countrySearch}"))`);
    await expect(spanLocator).toBeVisible();
    await this.page.keyboard.press('Enter');
    const storageLocator = this.page.locator('input[id="storageFacility.address-input"]');
    await storageLocator.click({ timeout: 5000 });
    await storageLocator.fill('30301');
    await expect(storageAddress.getByText(StorageLocationLookup, { exact: true })).toBeVisible();
    await this.page.keyboard.press('Enter');
    await this.page.getByLabel('Lot number').fill('1111111');

    // Requested delivery date
    await this.page.getByLabel('Requested delivery date').fill('20 Mar 2022');
    await this.page.getByLabel('Requested delivery date').blur();

    // Delivery Address
    const DeliveryLocationLookup = 'MONTGOMERY, AL 36101 (MONTGOMERY)';
    const deliveryLocation = this.page.getByRole('group', { name: 'Delivery Address' });
    await deliveryLocation.getByLabel('Address 1').fill('448 Washington Blvd NE');
    await deliveryLocation.getByLabel('Address 2').fill('Apt D3');
    await this.page.locator('input[id="delivery.address-country-input"]').fill(countrySearch);
    spanLocator = this.page.locator(`span:has(mark:has-text("${countrySearch}"))`);
    await expect(spanLocator).toBeVisible();
    await this.page.keyboard.press('Enter');
    const deliveryLocator = this.page.locator('input[id="delivery.address-input"]');
    await deliveryLocator.click({ timeout: 5000 });
    await deliveryLocator.fill('36101');
    await expect(deliveryLocation.getByText(DeliveryLocationLookup, { exact: true })).toBeVisible();
    await this.page.keyboard.press('Enter');

    // Remarks
    await this.page.getByLabel('Counselor remarks').fill('NTS-release counselor remarks');

    // Save the shipment, progress back to the move details page, and verify it's been created
    await this.page.getByRole('button', { name: 'Save' }).click();
    await this.waitForPage.moveDetails();
  }

  /**
   * @param {Object} options
   * @param {boolean} [options.selectAdvance=false]
   * returns {Promise<void>}
   */
  async fillOutAboutPage(options = { selectAdvance: false }) {
    // editing this field with the keyboard instead of the date picker runs async validators for pre-filled postal codes
    // this helps debounce the API calls that would be triggered in quick succession
    await this.page.locator('input[name="actualMoveDate"]').fill('01 Feb 2022');

    const LocationLookup = 'YUMA, AZ 85364 (YUMA)';
    const countrySearch = 'UNITED STATES';

    await this.page.locator('input[name="pickupAddress.streetAddress1"]').fill('1819 S Cedar Street');
    await this.page.locator('input[id="pickupAddress-country-input"]').fill(countrySearch);
    let spanLocator = this.page.locator(`span:has(mark:has-text("${countrySearch}"))`);
    await expect(spanLocator).toBeVisible();
    await this.page.keyboard.press('Enter');
    const pickupLocator = this.page.locator('input[id="pickupAddress-input"]');
    await pickupLocator.click({ timeout: 5000 });
    await pickupLocator.fill('85364');
    await expect(this.page.getByText(LocationLookup, { exact: true })).toBeVisible();
    await this.page.keyboard.press('Enter');

    await this.page.locator('input[name="destinationAddress.streetAddress1"]').fill('1819 S Cedar Street');
    await this.page.locator('input[id="destinationAddress-country-input"]').fill(countrySearch);
    spanLocator = this.page.locator(`span:has(mark:has-text("${countrySearch}"))`);
    await expect(spanLocator).toBeVisible();
    await this.page.keyboard.press('Enter');
    const deliveryLocator = this.page.locator('input[id="destinationAddress-input"]');
    await deliveryLocator.click({ timeout: 5000 });
    await deliveryLocator.fill('85364');
    spanLocator = this.page.locator(`span:has(mark:has-text("85364"))`);
    await expect(spanLocator).toBeVisible();
    await this.page.keyboard.press('Enter');

    if (options?.selectAdvance) {
      await this.page.locator('label[for="yes-has-received-advance"]').click();
      await this.page.locator('input[name="advanceAmountReceived"]').fill('5000');
    } else {
      await this.page.locator('label[for="no-has-received-advance"]').click();
    }

    await this.page.locator('input[name="w2Address.streetAddress1"]').fill('1819 S Cedar Street');
    await this.page.locator('input[id="w2Address-country-input"]').fill(countrySearch);
    spanLocator = this.page.locator(`span:has(mark:has-text("${countrySearch}"))`);
    await expect(spanLocator).toBeVisible();
    await this.page.keyboard.press('Enter');
    const w2Locator = this.page.locator('input[id="w2Address-input"]');
    await w2Locator.click({ timeout: 5000 });
    await w2Locator.fill('85364');
    spanLocator = this.page.locator(`span:has(mark:has-text("85364"))`);
    await expect(spanLocator).toBeVisible();
    await this.page.keyboard.press('Enter');

    await this.page.getByRole('button', { name: 'Save & Continue' }).click();
  }

  /**
   * @param {Object} options
   * @param {boolean} [options.hasTrailer=false]
   * @param {boolean} [options.ownTrailer=false]
   * @param {boolean} [options.useConstructedWeight=false]
   * returns {Promise<void>}
   */
  async fillOutWeightTicketPage(options) {
    const { hasTrailer = false, ownTrailer = false, useConstructedWeight = false } = options;
    await this.page.locator('input[name="vehicleDescription"]').fill('Kia Forte');
    await this.page.locator('input[name="vehicleDescription"]').blur();

    await this.page.getByLabel('Empty weight').clear();
    await this.page.getByLabel('Empty weight').fill('1000');
    await this.page.getByLabel('Empty weight').blur();
    if (useConstructedWeight) {
      // this page has multiple labels with the same text, grab the
      // first one for Empty Weight. Not sure why getByLabel does not work
      await this.page.locator('label').getByText("I don't have this weight ticket").first().click();

      const emptyRental = this.page.locator('label').getByText(
        `Since you do not have a certified weight ticket, upload the registration or rental agreement for the vehicle used
        during the PPM`,
      );

      await expect(emptyRental).toBeVisible();
      let filepond = emptyRental.locator('../..').locator('.filepond--wrapper');
      await expect(filepond).toBeVisible();

      await this.uploadFileViaFilepond(filepond, 'sampleWeightTicket.jpg');

      // wait for the file to be visible in the uploads
      await expect(
        filepond
          .locator('../..')
          .locator('p')
          .getByText(/sampleWeightTicket-\d{14}\.jpg/, { exact: false }),
      ).toBeVisible();

      await this.page.getByLabel('Full weight').clear();
      await this.page.getByLabel('Full weight').fill('3000');

      // this page has multiple labels with the same text, grab the
      // second one for Full Weight. Not sure why getByLabel does not work
      await this.page.locator('label').getByText("I don't have this weight ticket").nth(1).click();

      const fullConstructed = this.page
        .locator('label')
        .getByText('Upload your completed constructed weight spreadsheet');

      await expect(fullConstructed).toBeVisible();
      filepond = fullConstructed.locator('../..').locator('.filepond--wrapper');
      await expect(filepond).toBeVisible();

      await this.uploadFileViaFilepond(filepond, 'constructedWeight.xlsx');

      // weight estimator file should be converted to .pdf so we verify it was
      const re = /constructedWeight.+-\d{14}\.pdf$/;

      // wait for the file to be visible in the uploads
      await expect(filepond.locator('../..').locator('p').getByText(re, { exact: false })).toBeVisible();
    } else {
      // find the label, then find the filepond wrapper. Not sure why
      // getByLabel doesn't work
      const emptyWeightLabel = this.page.locator('label').getByText('Upload empty weight ticket', { exact: true });
      await expect(emptyWeightLabel).toBeVisible();
      const emptyFilepond = emptyWeightLabel.locator('../..').locator('.filepond--wrapper');
      await expect(emptyFilepond).toBeVisible();

      await this.uploadFileViaFilepond(emptyFilepond, 'sampleWeightTicket.jpg');

      // wait for the file to be visible in the uploads
      await expect(
        emptyFilepond
          .locator('../..')
          .locator('p')
          .getByText(/sampleWeightTicket-\d{14}\.jpg/, { exact: false }),
      ).toBeVisible();

      await this.page.getByLabel('Full Weight').clear();
      await this.page.getByLabel('Full Weight').fill('3000');

      // find the label, then find the filepond wrapper. Not sure why
      // getByLabel doesn't work
      const fullWeightLabel = this.page.locator('label').getByText('Upload full weight ticket', { exact: true });
      await expect(fullWeightLabel).toBeVisible();
      const fullFilepond = fullWeightLabel.locator('../..').locator('.filepond--wrapper');
      await expect(fullFilepond).toBeVisible();

      await this.uploadFileViaFilepond(fullFilepond, 'sampleWeightTicket.jpg');
      // wait for the file to be visible in the uploads
      await expect(
        fullFilepond
          .locator('../..')
          .locator('p')
          .getByText(/sampleWeightTicket-\d{14}\.jpg/, { exact: false }),
      ).toBeVisible();
    }

    await expect(this.page.locator('.tripWeightTotal')).toContainText('Trip weight: 2,000 lbs');

    if (hasTrailer) {
      // the page design makes it hard to click without using a css locator
      await this.page.locator('label[for="yesOwnsTrailer"]').click();
      if (ownTrailer) {
        // the page design makes it hard to click without using a css locator
        await this.page.locator('label[for="yestrailerMeetsCriteria"]').click();

        // find the label, then find the filepond wrapper, not sure
        // why getByLabel does not work
        const ownershipLabel = this.page.locator('label').getByText('Upload proof of ownership', { exact: true });
        await expect(ownershipLabel).toBeVisible();
        const ownershipFilepond = ownershipLabel.locator('../..').locator('.filepond--wrapper');
        await expect(ownershipFilepond).toBeVisible();

        await this.uploadFileViaFilepond(ownershipFilepond, 'trailerOwnership.pdf');

        // wait for the file to be visible in the uploads
        await expect(
          ownershipFilepond
            .locator('../..')
            .locator('p')
            .getByText(/trailerOwnership-\d{14}\.pdf/, { exact: false }),
        ).toBeVisible();
      } else {
        // the page design makes it hard to click without using a css locator
        await this.page.locator('label[for="notrailerMeetsCriteria"]').click();
      }
    }
  }

  async fillOutWeightTicketWithIncorrectXlsx() {
    const emptyWeightLabel = this.page.locator('label').getByText('Upload empty weight ticket', { exact: true });
    await expect(emptyWeightLabel).toBeVisible();
    const emptyFilepond = emptyWeightLabel.locator('../..').locator('.filepond--wrapper');
    await expect(emptyFilepond).toBeVisible();

    await this.uploadFileViaFilepond(emptyFilepond, 'weightEstimatorExpectFailedUpload.xlsx');

    // await modal is visible and close modal
    await expect(
      this.page.getByText(
        'The only Excel file this uploader accepts is the Weight Estimator file. Please convert any other Excel file to PDF.',
      ),
    ).toBeVisible();
    await this.page.getByTestId('modalCloseButton').click();

    // wait for the an incorrect file to not be visible in the uploads
    await expect(this.page.getByRole('heading', { name: '1 FILES UPLOADED' })).not.toBeVisible();

    // find the label, then find the filepond wrapper.
    const fullWeightLabel = this.page.locator('label').getByText('Upload full weight ticket', { exact: true });
    await expect(fullWeightLabel).toBeVisible();
    const fullFilepond = fullWeightLabel.locator('../..').locator('.filepond--wrapper');
    await expect(fullFilepond).toBeVisible();

    await this.uploadFileViaFilepond(fullFilepond, 'weightEstimatorExpectFailedUpload.xlsx');
    // await modal is visible and close modal
    await expect(
      this.page.getByText(
        'The only Excel file this uploader accepts is the Weight Estimator file. Please convert any other Excel file to PDF.',
      ),
    ).toBeVisible();
    await this.page.getByTestId('modalCloseButton').click();

    // wait for the file not to be visible in the uploads
    await expect(this.page.getByRole('heading', { name: '1 FILES UPLOADED' })).not.toBeVisible();

    // add successful file upload and look for "1 FILES UPLOADED": weightEstimatorExpectSuccessfulUpload
    await this.uploadFileViaFilepond(fullFilepond, 'weightEstimatorExpectSuccessfulUpload.xlsx');
    // wait for the file to be visible in the uploads
    await expect(this.page.getByRole('heading', { name: '1 FILES UPLOADED' })).toBeVisible();
  }

  async fillOutProGearWithIncorrectXlsx() {
    // find the label, then find the filepond wrapper.
    const proGearWeightLabel = this.page.locator('label').getByText('Upload your');
    await expect(proGearWeightLabel).toBeVisible();
    const proGearFilepond = proGearWeightLabel.locator('../..').locator('.filepond--wrapper');
    await expect(proGearFilepond).toBeVisible();

    await this.uploadFileViaFilepond(proGearFilepond, 'weightEstimatorExpectFailedUpload.xlsx');

    // await modal is visible and close modal
    await expect(
      this.page.getByText(
        'The only Excel file this uploader accepts is the Weight Estimator file. Please convert any other Excel file to PDF.',
      ),
    ).toBeVisible();
    await this.page.getByTestId('modalCloseButton').click();

    // wait for the an incorrect file to not be visible in the uploads
    await expect(this.page.getByRole('heading', { name: '1 FILES UPLOADED' })).not.toBeVisible();

    // add successful file upload and look for "1 FILES UPLOADED": weightEstimatorExpectSuccessfulUpload
    await this.uploadFileViaFilepond(proGearFilepond, 'weightEstimatorExpectSuccessfulUpload.xlsx');
    // wait for the file to be visible in the uploads
    await expect(this.page.getByRole('heading', { name: '1 FILES UPLOADED' })).toBeVisible();
  }
}

/**
 * @typedef {object} ServiceCounselorPageTestArgs - services counselor page test args
 * @property {ServiceCounselorPage} scPage    - services counselor page
 */

/** @type {import('@playwright/test').Fixtures<ServiceCounselorPageTestArgs, {}, import('../../utils/office/officeTest').OfficePageTestArgs, import('@playwright/test').PlaywrightWorkerArgs>} */
const scFixtures = {
  scPage: async ({ officePage }, use) => {
    const scPage = new ServiceCounselorPage(officePage);
    await scPage.signInAsNewServicesCounselorUser();
    await use(scPage);
  },
};

export const test = officeTest.extend(scFixtures);

export { expect };

export default ServiceCounselorPage;
