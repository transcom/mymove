/**
 * Semi-automated converted from a cypress test, and thus may contain
 * non best-practices, in particular: heavy use of `page.locator`
 * instead of `page.getBy*`.
 */

// @ts-check
const { expect, test, useMobileViewport, CustomerPpmPage } = require('./customerPpmTestFixture');

test.describe('PPM Onboarding - Estimated Incentive', () => {
  /** @type {CustomerPpmPage} */
  let customerPpmPage;

  test.beforeEach(async ({ customerPage }) => {
    const move = await customerPage.testHarness.buildUnSubmittedMoveWithPPMShipmentThroughEstimatedWeights();
    customerPpmPage = new CustomerPpmPage(customerPage, move);
    await customerPpmPage.signInAndNavigateFromHomePageToExistingPPMDateAndLocationPage();
    await customerPpmPage.navigateFromDateAndLocationPageToEstimatedWeightsPage();
    await customerPpmPage.navigateFromEstimatedWeightsPageToEstimatedIncentivePage();
    await expect(customerPage.page.locator('.container h2')).toContainText('$10,000');
  });

  //
  // https://playwright.dev/docs/test-parameterize
  //
  // use forEach to avoid
  // https://eslint.org/docs/latest/rules/no-loop-func
  [true, false].forEach((isMobile) => {
    const viewportName = isMobile ? 'mobile' : 'desktop';
    test.describe(`with ${viewportName} viewport`, async () => {
      if (isMobile) {
        useMobileViewport();
      }
      test('go to estimated incentives page', async () => {
        customerPpmPage.generalVerifyEstimatedIncentivePage({ isMobile });
      });
    });
  });
});
