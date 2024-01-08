// @ts-nocheck
import { PrimeSimulatorUserType } from '../../utils/office/officeTest';
import { test, expect, AuditTestPage } from './auditTestSetup';
import { RUN_PRIME_SETS_UP_SHIPMENT } from './primeEntersShipmentDetailsSteps';
import { RUN_MOVER_OWNER_SETS_UP_HHG } from './customerHhgSetupSteps';
import './auditUtils';

const GetPrimeStepHelpers = ({ page, helpers }) => ({
  ...helpers,
  signInAsNewPrimeSimulatorUser: async () => {
    await page.goto('/devlocal-auth/login');
    await page.locator(`button[data-hook="new-user-login-${PrimeSimulatorUserType}"]`).click();
    await page.getByRole('heading', { name: 'Moves available to Prime' }).waitFor();
  },
});

export const getCustomerStepHelpers = ({ page, helpers: { stringHelpers, ...helpers } }) => ({
  ...helpers,
  stringHelpers,
  signInAsExisting: async (userId) => {
    const currentUrl = page.url();
    const baseUrl = new URL(currentUrl).origin;
    const loginPath = '/devlocal-auth/login';
    const loginUrl = stringHelpers.concatUrlSegments([baseUrl, loginPath]);
    await page.goto(loginUrl);

    const navigationWaiter = page.waitForURL(({ pathname, ...url }) => pathname === '/');

    await page.locator(`button[value="${userId}"]`).click();

    const test = await navigationWaiter;
    await page.locator('#title-announcer', { hasText: 'my.move.mil' }).waitFor();
  },
  signInAsNew: async () => {
    const currentUrl = page.url();
    const baseUrl = new URL(currentUrl).origin;
    const loginPath = '/devlocal-auth/login';
    const loginUrl = stinerHelpers.trimUrlOperation([baseUrl, loginPath]);
    await page.goto('http://milmovelocal:3000/devlocal-auth/login');
    await page.locator(`button[data-hook="new-user-login-${milmoveUserType}"]`).click();
  },
});

/** Array of test scenarios which also receive the test fixtures via pageInstance */
const SCENARIOS = {
  CREATING_SHIPMENT_CAUSES_SIT_CALCULATION: async ({ pageInstance }) => {
    //await RUN_PRIME_SETS_UP_SHIPMENT.run(pageInstance, GetPrimeStepHelpers);
    await RUN_MOVER_OWNER_SETS_UP_HHG.run(pageInstance, getCustomerStepHelpers);
  },
};

const TESTS = Object.entries(SCENARIOS).map(([key, testFunction]) => ({
  testLabel: key.split('_').join(' '),
  testFunction,
}));

/** Helper function to allow referencing the test info which is passed to the test framework function*/
const GetTest = (testInfo) => test(testInfo.testLabel, testInfo.testFunction);

test.describe('Prime Updates SIT', () => {
  GetTest(TESTS[0]);
});