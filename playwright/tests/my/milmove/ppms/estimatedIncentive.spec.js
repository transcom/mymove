/**
 * Semi-automated converted from a cypress test, and thus may contain
 * non best-practices, in particular: heavy use of `page.locator`
 * instead of `page.getBy*`.
 */

// @ts-check
import { test, expect, forEachViewport } from './customerPpmTestFixture';

const multiMoveEnabled = process.env.FEATURE_FLAG_MULTI_MOVE;

test.describe('PPM Onboarding - Estimated Incentive', () => {
  forEachViewport(async ({ isMobile }) => {
    test.beforeEach(async ({ customerPpmPage }) => {
      const move = await customerPpmPage.testHarness.buildUnSubmittedMoveWithPPMShipmentThroughEstimatedWeights();
      await customerPpmPage.signInForPPMWithMove(move);
      await customerPpmPage.navigateFromHomePageToExistingPPMDateAndLocationPage();
      await customerPpmPage.navigateFromDateAndLocationPageToEstimatedWeightsPage();
      await customerPpmPage.navigateFromEstimatedWeightsPageToEstimatedIncentivePage();
      await expect(customerPpmPage.page.locator('.container h2')).toContainText('is your estimated incentive');
    });

    test('go to estimated incentives page', async ({ customerPpmPage }) => {
      await customerPpmPage.generalVerifyEstimatedIncentivePage({ isMobile });
    });
  });
});

test.describe('(MultiMove) PPM Onboarding - Estimated Incentive', () => {
  test.skip(multiMoveEnabled === 'false', 'Skip if MultiMove workflow is not enabled.');

  forEachViewport(async ({ isMobile }) => {
    test.beforeEach(async ({ customerPpmPage }) => {
      const move = await customerPpmPage.testHarness.buildUnSubmittedMoveWithPPMShipmentThroughEstimatedWeights();
      await customerPpmPage.signInForPPMWithMove(move);
      await customerPpmPage.navigateFromHomePageToExistingPPMDateAndLocationPage();
      await customerPpmPage.navigateFromDateAndLocationPageToEstimatedWeightsPage();
      await customerPpmPage.navigateFromEstimatedWeightsPageToEstimatedIncentivePage();
      await expect(customerPpmPage.page.locator('.container h2')).toContainText('is your estimated incentive');
    });

    test('go to estimated incentives page', async ({ customerPpmPage }) => {
      await customerPpmPage.generalVerifyEstimatedIncentivePage({ isMobile });
    });
  });
});
