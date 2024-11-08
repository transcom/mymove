/**
 * Semi-automated converted from a cypress test, and thus may contain
 * non best-practices, in particular: heavy use of `page.locator`
 * instead of `page.getBy*`.
 */

// @ts-check
import { expect, test as scTest, ServiceCounselorPage } from '../servicescounseling/servicesCounselingTestFixture';

/**
 * PpmPage test fixture
 *
 * @extends ServiceCounselorPage
 */
export class PpmPage extends ServiceCounselorPage {
  /**
   * @param {Object} options
   * @param {string} [options.estimatedWeight=4000]
   * @param {boolean} [options.hasProGear=false]
   * @param {string} [options.proGearWeight=1000]
   * @param {string} [options.spouseProGearWeight=500]
   *
   * @returns Promise<void>
   */
  async fillOutWeight(options = {}) {
    const {
      estimatedWeight = '4000',
      hasProGear = false,
      proGearWeight = '1000',
      spouseProGearWeight = '500',
    } = options;

    await this.page.locator('input[name="estimatedWeight"]').clear();
    await this.page.locator('input[name="estimatedWeight"]').fill(estimatedWeight);

    if (hasProGear) {
      await this.page.locator('label[for="hasProGearYes"]').click();
      await this.page.locator('input[name="proGearWeight"]').fill(proGearWeight);
      await this.page.locator('input[name="spouseProGearWeight"]').fill(spouseProGearWeight);
    } else {
      await this.page.locator('label[for="hasProGearNo"]').click();
    }
  }

  /**
   * @param {Object} options
   * @param {boolean} [options.hasSit=true]
   * @param {string} [options.sitEstimatedWeight=1000]
   * @param {string} [options.sitEstimatedEntryDate='01 Mar 2020']
   * @param {string} [options.sitEstimatedDepartureDate='31 Mar 2020']
   * @param {string} [options.sitLocation='Destination']
   *
   * @returns Promise<void>
   */
  async fillOutSitExpected(options = {}) {
    const {
      hasSit = true,
      sitEstimatedWeight = '1000',
      sitEstimatedEntryDate = '01 Mar 2020',
      sitEstimatedDepartureDate = '31 Mar 2020',
      sitLocation = 'Destination', // takes in values of either 'Destination' or 'Origin'
    } = options;

    if (hasSit) {
      await this.page.locator('label[for="sitExpectedYes"]').click();
      await this.page.locator(`label[for="sitLocation${sitLocation}"]`).click();
      await this.page.locator('input[name="sitEstimatedWeight"]').clear();
      await this.page.locator('input[name="sitEstimatedWeight"]').fill(sitEstimatedWeight);
      await this.page.locator('input[name="sitEstimatedEntryDate"]').clear();
      await this.page.locator('input[name="sitEstimatedEntryDate"]').fill(sitEstimatedEntryDate);
      await this.page.locator('input[name="sitEstimatedDepartureDate"]').clear();
      await this.page.locator('input[name="sitEstimatedDepartureDate"]').fill(sitEstimatedDepartureDate);
    } else {
      await this.page.locator('label[for="sitExpectedNo"]').click();
    }
  }

  /**
   * @param {Object} options
   * @param {string} [options.expectedDepartureDate='09 Jun 2022']
   * @param {string} [options.pickupPostalCode=90210]
   * @param {string} [options.secondPickupPostalCode='07003']
   *
   * @returns Promise<void>
   */
  async fillOutOriginInfo() {
    const LocationLookup = 'BEVERLY HILLS, CA 90210 (LOS ANGELES)';

    await this.page.locator('input[name="expectedDepartureDate"]').fill('09 Jun 2022');

    await this.page.locator('input[name="pickup.address.streetAddress1"]').fill('123 Street');
    await this.page.getByLabel('Location Lookup').nth(0).fill('90210');
    await expect(this.page.getByText(LocationLookup, { exact: true })).toBeVisible();
    await this.page.keyboard.press('Enter');
  }

  /**
   * @param {Object} options
   * @param {string} [options.destinationPostalCode=76127]
   * @param {string} [options.secondDestinatixonPostalCode=98444]
   *
   * @returns Promise<void>
   */
  async fillOutDestinationInfo() {
    const LocationLookup = 'FORT WORTH, TX 76127 (TARRANT)';

    await this.page.locator('input[name="destination.address.streetAddress1"]').fill('123 Street');
    await this.page.getByLabel('Location Lookup').nth(1).fill('76127');
    await expect(this.page.getByText(LocationLookup, { exact: true })).toBeVisible();
    await this.page.keyboard.press('Enter');
  }

  /**
   * @param {Object} options
   * @param {boolean} [options.hasAdvance=true]
   * @param {string} [options.advance=6000]
   *
   * @returns Promise<void>
   */
  async fillOutIncentiveAndAdvance(options = {}) {
    const { hasAdvance = true, advance = '6000' } = options;

    if (hasAdvance) {
      await this.page.locator('label[for="hasRequestedAdvanceYes"]').click();
      await this.page.locator('input[name="advance"]').clear();
      await this.page.locator('input[name="advance"]').fill(advance);
      await this.page.locator('input[name="advance"]').blur();
      await this.page.locator('label[for="approveAdvanceRequest"]').click();
    } else {
      await this.page.locator('label[for="hasRequestedAdvanceNo"]').click();
    }
  }
}

/**
 * @typedef {object} PpmPageTestArgs - ppm page test args
 * @property {PpmPage} ppmPage    -  ppm page
 */

/** @type {import('@playwright/test').Fixtures<PpmPageTestArgs, {}, import('../../utils/office/officeTest').OfficePageTestArgs, import('@playwright/test').PlaywrightWorkerArgs>} */
const ppmFixtures = {
  ppmPage: async ({ officePage }, use) => {
    const ppmPage = new PpmPage(officePage);
    await ppmPage.signInAsNewServicesCounselorUser();
    await use(ppmPage);
  },
};

export const test = scTest.extend(ppmFixtures);

export { expect };

export default PpmPage;
