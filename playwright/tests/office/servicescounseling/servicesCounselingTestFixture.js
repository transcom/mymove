/**
 * Semi-automated converted from a cypress test, and thus may contain
 * non best-practices, in particular: heavy use of `page.locator`
 * instead of `page.getBy*`.
 */

// @ts-check
const { expect, test: officeTest, OfficePage } = require('../../utils/officeTest');
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
   * Service Counselor navigate to move
   * @param {string} moveLocator
   */
  async navigateToMove(moveLocator) {
    // type in move code/locator to filter
    await this.page.locator('input[name="locator"]').type(moveLocator);
    await this.page.locator('input[name="locator"]').blur();

    await this.page.locator('tbody > tr').first().click();
    await this.waitForLoading();
    expect(this.page.url()).toContain(`/counseling/moves/${moveLocator}/details`);
  }

  async addNTSShipment() {
    await this.page.locator('[data-testid="dropdown"]').first().selectOption({ label: 'NTS' });

    await this.page.locator('#requestedPickupDate').clear();
    await this.page.locator('#requestedPickupDate').type('16 Mar 2022');
    await this.page.locator('#requestedPickupDate').blur();
    await this.page.getByText('Use current address').click();
    await expect(this.page.locator(`[data-testid="remarks"]`)).not.toBeVisible();
    await this.page.locator('[name="counselorRemarks"]').type('Sample counselor remarks');

    await this.page.locator('[data-testid="submitForm"]').click();
    await this.waitForLoading();

    // the new shipment is visible on the Move details page
    await expect(this.page.locator('[data-testid="ShipmentContainer"]').last()).toContainText(
      'Sample counselor remarks',
    );
  }

  async addNTSReleaseShipment() {
    await this.page.locator('[data-testid="dropdown"]').first().selectOption({ label: 'NTS-release' });

    // Previously recorded weight
    await this.page.locator('#ntsRecordedWeight').type('1300');

    // Storage facility info
    await this.page.locator('#facilityName').type('Sample Facility Name');
    await this.page.locator('#facilityPhone').type('999-999-9999');
    await this.page.locator('#facilityEmail').type('sample@example.com');
    await this.page.locator('#facilityServiceOrderNumber').type('999999');

    // Storage facility address
    await this.page.locator('input[name="storageFacility.address.streetAddress1"]').type('148 S East St');
    await this.page.locator('input[name="storageFacility.address.streetAddress2"]').type('Suite 7A');
    await this.page.locator('input[name="storageFacility.address.city"]').type('Sample City');
    await this.page.locator('select[name="storageFacility.address.state"]').selectOption({ label: 'GA' });
    await this.page.locator('input[name="storageFacility.address.postalCode"]').type('30301');
    await this.page.locator('#facilityLotNumber').type('1111111');

    // Requested delivery date
    await this.page.locator('#requestedDeliveryDate').type('20 Mar 2022');
    await this.page.locator('#requestedDeliveryDate').blur();

    // Delivery location
    await this.page.locator('input[name="delivery.address.streetAddress1"]').type('448 Washington Blvd NE');
    await this.page.locator('input[name="delivery.address.streetAddress2"]').type('Apt D3');
    await this.page.locator('input[name="delivery.address.city"]').type('Another City');
    await this.page.locator('select[name="delivery.address.state"]').selectOption({ label: 'AL' });
    await this.page.locator('input[name="delivery.address.postalCode"]').type('36101');
    await this.page.locator('#destinationType').selectOption({ label: 'Home of record (HOR)' });

    // Receiving agent
    await this.page.locator('input[name="delivery.agent.firstName"]').type('Skyler');
    await this.page.locator('input[name="delivery.agent.lastName"]').type('Hunt');
    await this.page.locator('input[name="delivery.agent.phone"]').type('999-999-9999');
    await this.page.locator('input[name="delivery.agent.email"]').type('skyler.hunt@example.com');

    // Remarks
    await expect(this.page.locator(`[data-testid="remarks"]`)).not.toBeVisible();
    await this.page.locator('[data-testid="counselor-remarks"]').type('NTS-release counselor remarks');

    await this.page.locator('[data-testid="submitForm"]').click();
    await this.waitForLoading();

    // the new shipment is visible on the Move details page
    const lastShipment = this.page.locator('[data-testid="ShipmentContainer"]').last();

    await lastShipment.locator('[data-icon="chevron-down"]').click();

    await expect(lastShipment.locator('[data-testid="ntsRecordedWeight"]')).toContainText('1,300');
    await expect(lastShipment.locator('[data-testid="storageFacilityName"]')).toContainText('Sample Facility Name');
    await expect(lastShipment.locator('[data-testid="serviceOrderNumber"]')).toContainText('999999');
    await expect(lastShipment.locator('[data-testid="storageFacilityAddress"]')).toContainText(
      '148 S East St, Suite 7A, Sample City, GA 30301',
    );
    await expect(lastShipment.locator('[data-testid="storageFacilityAddress"]')).toContainText('1111111');
    await expect(lastShipment.locator('[data-testid="requestedDeliveryDate"]')).toContainText('20 Mar 2022');
    await expect(lastShipment.locator('[data-testid="destinationAddress"]')).toContainText(
      '448 Washington Blvd NE, Apt D3, Another City, AL 36101',
    );
    await expect(lastShipment.locator('[data-testid="secondaryDeliveryAddress"]')).toContainText('—');
    await expect(lastShipment.locator('[data-testid="customerRemarks"]')).toContainText('—');
    await expect(lastShipment.locator('[data-testid="counselorRemarks"]')).toContainText(
      'NTS-release counselor remarks',
    );
    await expect(lastShipment.locator('[data-testid="tacType"]')).toContainText('—');
    await expect(lastShipment.locator('[data-testid="sacType"]')).toContainText('—');
  }

  async editNTSReleaseShipment() {
    await this.page.getByRole('button', { name: 'Edit Shipment' }).click();

    // Previously recorded weight
    await this.page.locator('#ntsRecordedWeight').type('1100');

    // Storage facility info
    await this.page.locator('#facilityName').type('AAA Facility Name');
    await this.page.locator('#facilityPhone').type('999-999-9999');
    await this.page.locator('#facilityEmail').type('aaa@example.com');
    await this.page.locator('#facilityServiceOrderNumber').type('123456');

    // Storage facility address
    await this.page.locator('input[name="storageFacility.address.streetAddress1"]').type('9 W 2nd Ave');
    await this.page.locator('input[name="storageFacility.address.streetAddress2"]').type('Bldg 3');
    await this.page.locator('input[name="storageFacility.address.city"]').type('Big City');
    await this.page.locator('select[name="storageFacility.address.state"]').selectOption({ label: 'SC' });
    await this.page.locator('input[name="storageFacility.address.postalCode"]').type('29201');
    await this.page.locator('#facilityLotNumber').type('2222222');

    // Requested delivery date
    await this.page.locator('#requestedDeliveryDate').clear();
    await this.page.locator('#requestedDeliveryDate').type('21 Mar 2022');
    await this.page.locator('#requestedDeliveryDate').blur();

    // Delivery location
    await this.page.locator('input[name="delivery.address.streetAddress1"]').clear();
    await this.page.locator('input[name="delivery.address.streetAddress1"]').type('4124 Apache Dr');
    await this.page.locator('input[name="delivery.address.streetAddress2"]').clear();
    await this.page.locator('input[name="delivery.address.streetAddress2"]').type('Apt 18C');
    await this.page.locator('input[name="delivery.address.city"]').clear();
    await this.page.locator('input[name="delivery.address.city"]').type('Little City');
    await this.page.locator('select[name="delivery.address.state"]').selectOption({ label: 'GA' });
    await this.page.locator('input[name="delivery.address.postalCode"]').clear();
    await this.page.locator('input[name="delivery.address.postalCode"]').type('30901');
    await this.page.locator('#destinationType').selectOption({ label: 'Home of record (HOR)' });

    // Receiving agent
    await this.page.locator('input[name="delivery.agent.firstName"]').type('Jody');
    await this.page.locator('input[name="delivery.agent.lastName"]').type('Pitkin');
    await this.page.locator('input[name="delivery.agent.phone"]').type('999-111-1111');
    await this.page.locator('input[name="delivery.agent.email"]').type('jody.pitkin@example.com');

    // Remarks
    await expect(this.page.locator(`[data-testid="remarks"]`)).not.toBeVisible();
    await this.page.locator('[data-testid="counselor-remarks"]').type('NTS-release edited counselor remarks');

    await this.page.locator('[data-testid="submitForm"]').click();
    // the shipment should be saved with the type
    await this.waitForLoading();

    // the new shipment is visible on the Move details page
    const lastShipment = this.page.locator('[data-testid="ShipmentContainer"]').last();
    await lastShipment.locator('[data-icon="chevron-down"]').click();

    await expect(lastShipment.locator('[data-testid="ntsRecordedWeight"]')).toContainText('1,100');
    await expect(lastShipment.locator('[data-testid="storageFacilityName"]')).toContainText('AAA Facility Name');
    await expect(lastShipment.locator('[data-testid="serviceOrderNumber"]')).toContainText('123456');
    await expect(lastShipment.locator('[data-testid="storageFacilityAddress"]')).toContainText(
      '9 W 2nd Ave, Bldg 3, Big City, SC 29201',
    );
    await expect(lastShipment.locator('[data-testid="storageFacilityAddress"]')).toContainText('2222222');
    await expect(lastShipment.locator('[data-testid="requestedDeliveryDate"]')).toContainText('21 Mar 2022');
    await expect(lastShipment.locator('[data-testid="destinationAddress"]')).toContainText(
      '4124 Apache Dr, Apt 18C, Little City, GA 30901',
    );
    await expect(lastShipment.locator('[data-testid="secondaryDeliveryAddress"]')).toContainText('—');
    await expect(lastShipment.locator('[data-testid="customerRemarks"]')).toContainText('—');
    await expect(lastShipment.locator('[data-testid="counselorRemarks"]')).toContainText(
      'NTS-release edited counselor remarks',
    );
    await expect(lastShipment.locator('[data-testid="tacType"]')).toContainText('—');
    await expect(lastShipment.locator('[data-testid="sacType"]')).toContainText('—');
  }
}

/**
 * @typedef {object} ServiceCounselorPageTestArgs - services counselor page test args
 * @property {ServiceCounselorPage} scPage    - services counselor page
 */

/** @type {import('@playwright/test').Fixtures<ServiceCounselorPageTestArgs, {}, import('../../utils/officeTest').OfficePageTestArgs, import('@playwright/test').PlaywrightWorkerArgs>} */
const scFixtures = {
  scPage: async ({ officePage }, use) => {
    const scPage = new ServiceCounselorPage(officePage);
    await scPage.signInAsNewServicesCounselorUser();
    use(scPage);
  },
};

export const test = officeTest.extend(scFixtures);

export { expect };

export default ServiceCounselorPage;
