/**
 * Semi-automated converted from a cypress test, and thus may contain
 * non best-practices, in particular: heavy use of `page.locator`
 * instead of `page.getBy*`.
 */

// @ts-check
const { expect, test, forEachViewport, CustomerPpmPage } = require('./customerPpmTestFixture');

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

  forEachViewport(async ({ isMobile }) => {
    test('go to estimated incentives page', async () => {
      customerPpmPage.generalVerifyEstimatedIncentivePage({ isMobile });
    });
  });
});
