/**
 * customer test fixture for playwright
 * Using javascript classes following the examples in
 * <https://playwright.dev/docs/test-fixtures>
 */
// @ts-check
const base = require('@playwright/test');

const { BaseTestPage } = require('../../utils/baseTest');

/**
 * CustomerPage
 * @extends BaseTestPage
 */
class CustomerPage extends BaseTestPage {
  waitForPage = {
    localLogin: async () => {
      await base.expect(this.page.getByRole('heading', { name: 'Select an Existing User' })).toBeVisible();
    },
    onboardingConus: async () => {
      await base
        .expect(this.page.getByRole('heading', { name: 'Where are you moving?' }))
        .toBeVisible({ timeout: 10000 });
    },
    onboardingDodId: async () => {
      await base.expect(this.page.getByRole('heading', { name: 'Create your profile' })).toBeVisible();
    },
    onboardingName: async () => {
      await base.expect(this.page.getByRole('heading', { name: 'Name' })).toBeVisible();
    },
    onboardingContactInfo: async () => {
      await base.expect(this.page.getByRole('heading', { name: 'Your contact info' })).toBeVisible();
    },
    onboardingDutyLocation: async () => {
      await base.expect(this.page.getByRole('heading', { name: 'Current duty location' })).toBeVisible();
    },
    onboardingCurrentAddress: async () => {
      await base.expect(this.page.getByRole('heading', { name: 'Current mailing address' })).toBeVisible();
    },
    onboardingBackupAddress: async () => {
      await base.expect(this.page.getByRole('heading', { name: 'Backup mailing address' })).toBeVisible();
    },
    onboardingBackupContact: async () => {
      await base.expect(this.page.getByRole('heading', { name: 'Backup contact' })).toBeVisible();
    },
    home: async () => {
      await base.expect(this.page.getByRole('heading', { name: 'Leo Spacemen', level: 2 })).toBeVisible();
    },
    orders: async () => {
      await base.expect(this.page.getByRole('heading', { level: 1 })).toHaveText('Tell us about your move orders');
    },
    ordersUpload: async () => {
      await base.expect(this.page.getByRole('heading', { level: 1 })).toHaveText('Upload your orders');
    },
    ntsShipment: async () => {
      await base
        .expect(this.page.getByRole('heading', { level: 1 }))
        .toHaveText('Where and when should the movers pick up your things going into storage?');
    },
    ntsReleaseShipment: async () => {
      await base
        .expect(this.page.getByRole('heading', { level: 1 }))
        .toHaveText('Where and when should the movers deliver your things from storage?');
    },
    hhgShipment: async () => {
      await base
        .expect(this.page.getByRole('heading', { level: 1 }))
        .toHaveText('Movers pack and transport this shipment');
    },
    reviewShipments: async () => {
      await base.expect(this.page.getByRole('heading', { level: 1 })).toHaveText('Review your details');
    },
  };

  async navigateBack() {
    await this.page.getByTestId('wizardCancelButton').click();
  }

  async navigateForward() {
    await this.page.getByTestId('wizardNextButton').click();
  }
}

/**
 * @typedef {object} CustomerPageTestArgs - customer page test args
 * @property {CustomerPage} customerPage  - customer page
 */

/** @type {base.Fixtures<CustomerPageTestArgs, {}, base.PlaywrightTestArgs, base.PlaywrightWorkerArgs>} */
const customerFixtures = {
  customerPage: async ({ page, request }, use) => {
    const customerPage = new CustomerPage(page, request);
    await use(customerPage);
  },
};

exports.test = base.test.extend(customerFixtures);

exports.expect = base.expect;
