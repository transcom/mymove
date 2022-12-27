const base = require('@playwright/test');

class CustomerPage {
  constructor(page, request) {
    this.page = page;
    this.request = request;
  }

  waitForPage = {
    localLogin: async () => {
      await base.expect(await this.page.getByRole('heading', { name: 'Select an Existing User' })).toBeVisible();
    },
    onboardingConus: async () => {
      await base.expect(await this.page.getByRole('heading', { name: 'Where are you moving?' })).toBeVisible();
    },
    onboardingDodId: async () => {
      await base.expect(await this.page.getByRole('heading', { name: 'Create your profile' })).toBeVisible();
    },
    onboardingName: async () => {
      await base.expect(await this.page.getByRole('heading', { name: 'Name' })).toBeVisible();
    },
    onboardingContactInfo: async () => {
      await base.expect(await this.page.getByRole('heading', { name: 'Your contact info' })).toBeVisible();
    },
    onboardingDutyLocation: async () => {
      await base.expect(await this.page.getByRole('heading', { name: 'Current duty location' })).toBeVisible();
    },
    onboardingCurrentAddress: async () => {
      await base.expect(await this.page.getByRole('heading', { name: 'Current mailing address' })).toBeVisible();
    },
    onboardingBackupAddress: async () => {
      await base.expect(await this.page.getByRole('heading', { name: 'Backup mailing address' })).toBeVisible();
    },
    onboardingBackupContact: async () => {
      await base.expect(await this.page.getByRole('heading', { name: 'Backup contact' })).toBeVisible();
    },
    home: async () => {
      await base.expect(await this.page.getByRole('heading', { name: 'Leo Spacemen', level: 2 })).toBeVisible();
    },
    ntsShipment: async () => {
      await base
        .expect(await this.page.getByRole('heading', { level: 1 }))
        .toHaveText('Where and when should the movers pick up your things going into storage?');
    },
    ntsReleaseShipment: async () => {
      await base
        .expect(await this.page.getByRole('heading', { level: 1 }))
        .toHaveText('Where and when should the movers deliver your things from storage?');
    },
    reviewShipments: async () => {
      await base.expect(await this.page.getByRole('heading', { level: 1 })).toHaveText('Review your details');
    },
  };

  async signInAsNewCustomer() {
    await this.page.goto('/devlocal-auth/login');
    await this.waitForPage.localLogin();
    await this.page.getByRole('button', { name: 'Create a New milmove User' }).click();
  }

  async navigateBack() {
    await this.page.getByTestId('wizardCancelButton').click();
  }

  async navigateForward() {
    await this.page.getByTestId('wizardNextButton').click();
  }
}

exports.test = base.test.extend({
  customerPage: async ({ page, request }, use) => {
    const customerPage = new CustomerPage(page, request);
    await use(customerPage);
  },
});

exports.expect = base.expect;
