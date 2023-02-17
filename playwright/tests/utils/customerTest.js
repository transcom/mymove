/**
 * customer test fixture for playwright
 * Using javascript classes following the examples in
 * <https://playwright.dev/docs/test-fixtures>
 */
// @ts-check
import * as base from '@playwright/test';

import { BaseTestPage } from './baseTest';

/**
 * devlocal auth user types
 */
export const milmoveUserType = 'milmove';

/**
 * CustomerPage
 * @extends BaseTestPage
 */
export class CustomerPage extends BaseTestPage {
  waitForPage = {
    localLogin: async () => {
      await base.expect(this.page.getByRole('heading', { name: 'Select an Existing User' })).toBeVisible();
    },
    onboardingConus: async () => {
      await base.expect(this.page.getByRole('heading', { name: 'Where are you moving?' })).toBeVisible();
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
    aboutShipments: async () => {
      await base
        .expect(this.page.getByRole('heading', { level: 1 }))
        .toHaveText('Things to know about selecting shipments');
    },
    selectShipmentType: async () => {
      await base.expect(this.page.getByRole('heading', { level: 1 })).toHaveText('How should this shipment move?');
    },
    hhgShipment: async () => {
      await base
        .expect(this.page.getByRole('heading', { level: 1 }))
        .toHaveText('Movers pack and transport this shipment');
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
    reviewShipments: async () => {
      await base.expect(this.page.getByRole('heading', { level: 1 })).toHaveText('Review your details');
    },
  };

  /**
   * Sign in as new customer with devlocal
   *
   */
  async signInAsNewCustomer() {
    await this.signInAsNewUser(milmoveUserType);
  }

  /**
   * Sign in as existing customer with devlocal
   *
   * @param {string} userId
   */
  async signInAsExistingCustomer(userId) {
    await this.signInAsUserWithId(userId);
  }

  async navigateBack() {
    await this.page.getByTestId('wizardCancelButton').click();
  }

  async navigateForward() {
    await this.page.getByTestId('wizardNextButton').click();
  }
}

/**
 * set the viewport for mobile dimensions
 */
export function useMobileViewport() {
  // https://playwright.dev/docs/emulation#viewport
  base.test.use({ viewport: { width: 479, height: 875 } });
}

/**
 * @typedef {Object} ViewportCallbackProps
 * @property {string} viewportName
 * @property {boolean} isMobile
 */

/**
 * @callback forEachViewportCallback
 * @param {ViewportCallbackProps} props
 */

/**
 * @param {forEachViewportCallback} callbackfn
 */
export function forEachViewport(callbackfn) {
  const viewportsString = process.env.PLAYWRIGHT_VIEWPORTS || 'desktop mobile';
  const viewports = viewportsString.split(/\s+/);
  // use forEach to avoid
  // https://eslint.org/docs/latest/rules/no-loop-func
  viewports.forEach(async (viewportName) => {
    const isMobile = viewportName === 'mobile';
    //
    // https://playwright.dev/docs/test-parameterize
    //
    base.test.describe(`with ${viewportName} viewport`, async () => {
      if (isMobile) {
        useMobileViewport();
      }
      await callbackfn({ viewportName, isMobile });
    });
  });
}

/**
 * @typedef {object} CustomerPageTestArgs - customer page test args
 * @property {CustomerPage} customerPage  - customer page
 */

/** @type {base.Fixtures<CustomerPageTestArgs, {}, base.PlaywrightTestArgs, base.PlaywrightWorkerArgs>} */
const officeFixtures = {
  customerPage: async ({ page, request }, use) => {
    const customerPage = new CustomerPage(page, request);
    await use(customerPage);
  },
};

export const test = base.test.extend(officeFixtures);

export const { expect } = base;

export default test;
