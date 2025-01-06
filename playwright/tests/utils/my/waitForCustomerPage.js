// @ts-check
import * as base from '@playwright/test';

import WaitForPage from '../waitForPage';

/**
 * extension of WaitForPage that provides functions to wait for pages in the customer app
 * @extends WaitForPage
 */
export class WaitForCustomerPage extends WaitForPage {
  /**
   * @returns {Promise<void>}
   */
  async onboardingConus() {
    await this.runAccessibilityAudit();
    await base.expect(this.page.getByRole('heading', { name: 'Where are you moving?' })).toBeVisible();
    await this.runAccessibilityAudit();
  }

  /**
   * @returns {Promise<void>}
   */
  async onboardingDodId() {
    await this.runAccessibilityAudit();
    await base.expect(this.page.getByRole('heading', { name: 'Create your profile' })).toBeVisible();
    await this.runAccessibilityAudit();
  }

  /**
   * @returns {Promise<void>}
   */
  async onboardingName() {
    await this.runAccessibilityAudit();
    await base.expect(this.page.getByRole('heading', { name: 'Name' })).toBeVisible();
    await this.runAccessibilityAudit();
  }

  /**
   * @returns {Promise<void>}
   */
  async onboardingContactInfo() {
    await this.runAccessibilityAudit();
    await base.expect(this.page.getByRole('heading', { name: 'Your contact info' })).toBeVisible();
    await this.runAccessibilityAudit();
  }

  /**
   * @returns {Promise<void>}
   */
  async onboardingDutyLocation() {
    await this.runAccessibilityAudit();
    await base.expect(this.page.getByRole('heading', { name: 'Current duty location' })).toBeVisible();
    await this.runAccessibilityAudit();
  }

  /**
   * @returns {Promise<void>}
   */
  async onboardingCurrentAddress() {
    await this.runAccessibilityAudit();
    await base.expect(this.page.getByRole('heading', { name: 'Current Address' })).toBeVisible();
    await this.runAccessibilityAudit();
  }

  /**
   * @returns {Promise<void>}
   */
  async onboardingBackupAddress() {
    await this.runAccessibilityAudit();
    await base.expect(this.page.getByRole('heading', { name: 'Backup address' })).toBeVisible();
    await this.runAccessibilityAudit();
  }

  /**
   * @returns {Promise<void>}
   */
  async onboardingBackupContact() {
    await this.runAccessibilityAudit();
    await base.expect(this.page.getByRole('heading', { name: 'Backup contact' })).toBeVisible();
    await this.runAccessibilityAudit();
  }

  /**
   * @returns {Promise<void>}
   */
  async home() {
    await this.runAccessibilityAudit();
    await base.expect(this.page.getByTestId('stepContainer1').getByText('Profile complete')).toBeVisible();
    await this.runAccessibilityAudit();
  }

  /**
   * @returns {Promise<void>}
   */
  async multiMoveDashboard() {
    await this.runAccessibilityAudit();
    await base.expect(this.page.getByText('Welcome to MilMove!')).toBeVisible();
    await this.runAccessibilityAudit();
  }

  /**
   * @returns {Promise<void>}
   */
  async multiMoveLandingPage() {
    await this.runAccessibilityAudit();
    await base.expect(this.page.getByText('Welcome to MilMove!')).toBeVisible();
    await this.runAccessibilityAudit();
  }

  /**
   * @returns {Promise<void>}
   */
  async ordersDetails() {
    await this.runAccessibilityAudit();
    await base.expect(this.page.getByRole('heading', { level: 1 })).toHaveText('Tell us about your move orders');
    await this.runAccessibilityAudit();
  }

  /**
   * @returns {Promise<void>}
   */
  async ordersUpload() {
    await this.runAccessibilityAudit();
    await base.expect(this.page.getByRole('heading', { level: 1 })).toHaveText('Upload your orders');
    await this.runAccessibilityAudit();
  }

  async editOrders() {
    await this.runAccessibilityAudit();
    await base.expect(this.page.getByRole('heading', { level: 1 })).toHaveText('Orders');
    await this.runAccessibilityAudit();
  }

  /**
   * @returns {Promise<void>}
   */
  async aboutShipments() {
    await this.runAccessibilityAudit();
    await base
      .expect(this.page.getByRole('heading', { level: 1 }))
      .toHaveText('Things to know about selecting shipments');
    await this.runAccessibilityAudit();
  }

  /**
   * @returns {Promise<void>}
   */
  async selectShipmentType() {
    await this.runAccessibilityAudit();
    await base.expect(this.page.getByRole('heading', { level: 1 })).toHaveText('How should this shipment move?');
    await this.runAccessibilityAudit();
  }

  /**
   * @returns {Promise<void>}
   */
  async hhgShipment() {
    await this.runAccessibilityAudit();
    await base
      .expect(this.page.getByRole('heading', { level: 1 }))
      .toHaveText('Movers pack and transport this shipment');
    await this.runAccessibilityAudit();
  }

  /**
   * @returns {Promise<void>}
   */
  async ntsShipment() {
    await this.runAccessibilityAudit();
    await base
      .expect(this.page.getByRole('heading', { level: 1 }))
      .toHaveText('Where and when should the movers pick up your personal property going into storage?');
    await this.runAccessibilityAudit();
  }

  /**
   * @returns {Promise<void>}
   */
  async ntsReleaseShipment() {
    await this.runAccessibilityAudit();
    await base
      .expect(this.page.getByRole('heading', { level: 1 }))
      .toHaveText('Where and when should the movers deliver your personal property from storage?');
    await this.runAccessibilityAudit();
  }

  /**
   * @returns {Promise<void>}
   */
  async boatShipment() {
    await this.runAccessibilityAudit();
    await base.expect(this.page.getByRole('heading', { level: 1 })).toHaveText('Boat details and measurements');
    await this.runAccessibilityAudit();
  }

  /**
   * @returns {Promise<void>}
   */
  async mobileHomeShipment() {
    await this.runAccessibilityAudit();
    await base.expect(this.page.getByRole('heading', { level: 1 })).toHaveText('Mobile Home details and measurements');
    await this.runAccessibilityAudit();
  }

  /**
   * @returns {Promise<void>}
   */
  async reviewShipments() {
    await this.runAccessibilityAudit();
    await base.expect(this.page.getByRole('heading', { level: 1 })).toHaveText('Review your details');
    await this.runAccessibilityAudit();
  }
}

export default WaitForCustomerPage;
