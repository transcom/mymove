/**
 * Semi-automated converted from a cypress test, and thus may contain
 * non best-practices, in particular: heavy use of `page.locator`
 * instead of `page.getBy*`.
 */

// @ts-check
import { expect, test, ServiceCounselorPage } from '../servicescounseling/servicesCounselingTestFixture';

/**
 * PpmPage test fixture
 *
 * @extends ServiceCounselorPage
 */
export class PpmPage extends ServiceCounselorPage {
  async fillOutWeight(options = {}) {
    const {
      estimatedWeight = '4000',
      hasProGear = false,
      proGearWeight = '1000',
      spouseProGearWeight = '500',
    } = options;

    await this.page.locator('input[name="estimatedWeight"]').clear();
    await this.page.locator('input[name="estimatedWeight"]').type(estimatedWeight);

    if (hasProGear) {
      await this.page.locator('label[for="hasProGearYes"]').click();
      await this.page.locator('input[name="proGearWeight"]').type(proGearWeight);
      await this.page.locator('input[name="spouseProGearWeight"]').type(spouseProGearWeight);
    } else {
      await this.page.locator('label[for="hasProGearNo"]').click();
    }
  }

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
      await this.page.locator('input[name="sitEstimatedWeight"]').type(sitEstimatedWeight);
      await this.page.locator('input[name="sitEstimatedEntryDate"]').clear();
      await this.page.locator('input[name="sitEstimatedEntryDate"]').type(sitEstimatedEntryDate);
      await this.page.locator('input[name="sitEstimatedDepartureDate"]').clear();
      await this.page.locator('input[name="sitEstimatedDepartureDate"]').type(sitEstimatedDepartureDate);
    } else {
      await this.page.locator('label[for="sitExpectedNo"]').click();
    }
  }

  async fillOutOriginInfo(options = {}) {
    const {
      expectedDepartureDate = '09 Jun 2022',
      pickupPostalCode = '90210',
      secondPickupPostalCode = '07003',
    } = options;

    await this.page.locator('input[name="expectedDepartureDate"]').clear();
    await this.page.locator('input[name="expectedDepartureDate"]').type(expectedDepartureDate);
    await this.page.locator('input[name="pickupPostalCode"]').clear();
    await this.page.locator('input[name="pickupPostalCode"]').type(pickupPostalCode);
    if (secondPickupPostalCode) {
      await this.page.locator('input[name="secondPickupPostalCode"]').clear();
      await this.page.locator('input[name="secondPickupPostalCode"]').type(secondPickupPostalCode);
    }
  }

  async fillOutDestinationInfo(options = {}) {
    const { destinationPostalCode = '76127', secondDestinationPostalCode = '08540' } = options;

    await this.page.locator('input[name="destinationPostalCode"]').clear();
    await this.page.locator('input[name="destinationPostalCode"]').type(destinationPostalCode);
    if (secondDestinationPostalCode) {
      await this.page.locator('input[name="secondDestinationPostalCode"]').clear();
      await this.page.locator('input[name="secondDestinationPostalCode"]').type(secondDestinationPostalCode);
    }
  }

  async fillOutIncentiveAndAdvance(options = {}) {
    const { hasAdvance = true, advance = '6000' } = options;

    if (hasAdvance) {
      await this.page.locator('label[for="hasRequestedAdvanceYes"]').click();
      await this.page.locator('input[name="advance"]').clear();
      await this.page.locator('input[name="advance"]').type(advance);
    } else {
      await this.page.locator('label[for="hasRequestedAdvanceNo"]').click();
    }
  }
}

export { expect, test };

export default PpmPage;
