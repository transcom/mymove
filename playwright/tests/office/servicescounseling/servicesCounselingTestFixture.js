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
    await this.page.getByLabel('Requested pickup date').fill('16 Mar 2022');
    await this.page.getByLabel('Requested pickup date').blur();
    await this.page.getByText('Use current address').click();

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

    const storageAddress = this.page.getByRole('heading', { name: 'Storage facility address' }).locator('..');
    await storageAddress.getByLabel('Address 1').fill('148 S East St');
    await storageAddress.getByLabel('Address 2').fill('Suite 7A');
    await this.page.locator('input[id="deliveryAddress-location-input"]').fill('30301');
    await expect(storageAddress.getByText(StorageLocationLookup, { exact: true })).toBeVisible();
    await this.page.keyboard.press('Enter');
    await this.page.getByLabel('Lot number').fill('1111111');

    // Requested delivery date
    await this.page.getByLabel('Requested delivery date').fill('20 Mar 2022');
    await this.page.getByLabel('Requested delivery date').blur();

    // Delivery location
    const DeliveryLocationLookup = 'MONTGOMERY, AL 36101 (MONTGOMERY)';

    const deliveryLocation = this.page.getByRole('group', { name: 'Delivery location' });
    await deliveryLocation.getByLabel('Address 1').fill('448 Washington Blvd NE');
    await deliveryLocation.getByLabel('Address 2').fill('Apt D3');
    await this.page.locator('input[id="deliveryAddress-location-input"]').fill('36101');
    await expect(deliveryLocation.getByText(DeliveryLocationLookup, { exact: true })).toBeVisible();
    await this.page.keyboard.press('Enter');

    // Remarks
    await this.page.getByLabel('Counselor remarks').fill('NTS-release counselor remarks');

    // Save the shipment, progress back to the move details page, and verify it's been created
    await this.page.getByRole('button', { name: 'Save' }).click();
    await this.waitForPage.moveDetails();
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
