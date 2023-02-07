/**
 * Semi-automated converted from a cypress test, and thus may contain
 * non best-practices, in particular: heavy use of `page.locator`
 * instead of `page.getBy*`.
 */

// @ts-check
const { expect, test, forEachViewport } = require('./customerPpmTestFixture');

test.describe('PPM Onboarding - Estimated Incentive', () => {
  forEachViewport(async ({ isMobile }) => {
    test.beforeEach(async ({ customerPpmPage }) => {
      const move = await customerPpmPage.testHarness.buildUnSubmittedMoveWithPPMShipmentThroughEstimatedWeights();
      await customerPpmPage.signInForPPMWithMove(move);
      await customerPpmPage.navigateFromHomePageToExistingPPMDateAndLocationPage();
      await customerPpmPage.navigateFromDateAndLocationPageToEstimatedWeightsPage();
      await customerPpmPage.navigateFromEstimatedWeightsPageToEstimatedIncentivePage();
      await expect(customerPpmPage.page.locator('.container h2')).toContainText('$10,000');
    });

    test('go to estimated incentives page', async ({ customerPpmPage }) => {
      await customerPpmPage.generalVerifyEstimatedIncentivePage({ isMobile });
    });
  });
});
