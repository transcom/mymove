// @ts-check
import * as base from '@playwright/test';

export class WaitForPage {
  /**
   * Create a WaitForPage.
   * Use a class so we can catch typos or type errors
   *
   * @param {import('@playwright/test').Page} page
   */
  constructor(page) {
    this.page = page;
  }

  /**
   * @returns {Promise<void>}
   */
  async localLogin() {
    await base.expect(this.page.getByRole('heading', { name: 'Select an Existing User' })).toBeVisible();
  }

  /**
   * @returns {Promise<void>}
   */
  async onboardingConus() {
    await base.expect(this.page.getByRole('heading', { name: 'Where are you moving?' })).toBeVisible();
  }

  /**
   * @returns {Promise<void>}
   */
  async onboardingDodId() {
    await base.expect(this.page.getByRole('heading', { name: 'Create your profile' })).toBeVisible();
  }

  /**
   * @returns {Promise<void>}
   */
  async onboardingName() {
    await base.expect(this.page.getByRole('heading', { name: 'Name' })).toBeVisible();
  }

  /**
   * @returns {Promise<void>}
   */
  async onboardingContactInfo() {
    await base.expect(this.page.getByRole('heading', { name: 'Your contact info' })).toBeVisible();
  }

  /**
   * @returns {Promise<void>}
   */
  async onboardingDutyLocation() {
    await base.expect(this.page.getByRole('heading', { name: 'Current duty location' })).toBeVisible();
  }

  /**
   * @returns {Promise<void>}
   */
  async onboardingCurrentAddress() {
    await base.expect(this.page.getByRole('heading', { name: 'Current mailing address' })).toBeVisible();
  }

  /**
   * @returns {Promise<void>}
   */
  async onboardingBackupAddress() {
    await base.expect(this.page.getByRole('heading', { name: 'Backup address' })).toBeVisible();
  }

  /**
   * @returns {Promise<void>}
   */
  async onboardingBackupContact() {
    await base.expect(this.page.getByRole('heading', { name: 'Backup contact' })).toBeVisible();
  }

  /**
   * @returns {Promise<void>}
   */
  async home() {
    await base.expect(this.page.getByTestId('stepContainer1').getByText('Profile complete')).toBeVisible();
  }

  /**
   * @returns {Promise<void>}
   */
  async ordersDetails() {
    await base.expect(this.page.getByRole('heading', { level: 1 })).toHaveText('Tell us about your move orders');
  }

  /**
   * @returns {Promise<void>}
   */
  async ordersUpload() {
    await base.expect(this.page.getByRole('heading', { level: 1 })).toHaveText('Upload your orders');
  }

  /**
   * @returns {Promise<void>}
   */
  async aboutShipments() {
    await base
      .expect(this.page.getByRole('heading', { level: 1 }))
      .toHaveText('Things to know about selecting shipments');
  }

  /**
   * @returns {Promise<void>}
   */
  async selectShipmentType() {
    await base.expect(this.page.getByRole('heading', { level: 1 })).toHaveText('How should this shipment move?');
  }

  /**
   * @returns {Promise<void>}
   */
  async hhgShipment() {
    await base
      .expect(this.page.getByRole('heading', { level: 1 }))
      .toHaveText('Movers pack and transport this shipment');
  }

  /**
   * @returns {Promise<void>}
   */
  async ntsShipment() {
    await base
      .expect(this.page.getByRole('heading', { level: 1 }))
      .toHaveText('Where and when should the movers pick up your things going into storage?');
  }

  /**
   * @returns {Promise<void>}
   */
  async ntsReleaseShipment() {
    await base
      .expect(this.page.getByRole('heading', { level: 1 }))
      .toHaveText('Where and when should the movers deliver your things from storage?');
  }

  /**
   * @returns {Promise<void>}
   */
  async reviewShipments() {
    await base.expect(this.page.getByRole('heading', { level: 1 })).toHaveText('Review your details');
  }
}

export default WaitForPage;
