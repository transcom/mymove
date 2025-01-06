/**
 * Semi-automated converted from a cypress test, and thus may contain
 * non best-practices, in particular: heavy use of `page.locator`
 * instead of `page.getBy*`.
 */

// @ts-check
import { test, forEachViewport } from './customerPpmTestFixture';

const multiMoveEnabled = process.env.FEATURE_FLAG_MULTI_MOVE;

const fullPPMShipmentFields = [
  ['Expected departure', '15 Mar 2020'],
  ['Origin ZIP', '90210'],
  ['Second origin ZIP', '90211'],
  ['Delivery Address ZIP', '30813'],
  ['Second Delivery Address ZIP', '30814'],
  ['Closeout office', 'Creech AFB'],
  ['Storage expected? (SIT)', 'No'],
  ['Estimated weight', '4,000 lbs'],
  ['Pro-gear', 'Yes, 1,987 lbs'],
  ['Spouse pro-gear', 'Yes, 498 lbs'],
  ['Estimated incentive', '$10,000'],
  ['Advance requested?', 'Yes, $5,987'],
];

test.describe('PPM Onboarding - Review', () => {
  test.skip(multiMoveEnabled === 'true', 'Skip if MultiMove workflow is enabled.');
  forEachViewport(async ({ isMobile }) => {
    test.beforeEach(async ({ customerPpmPage }) => {
      const move = await customerPpmPage.testHarness.buildUnsubmittedMoveWithMultipleFullPPMShipmentComplete();
      await customerPpmPage.signInForPPMWithMove(move);
      await customerPpmPage.navigateFromHomePageToReviewPage();
    });

    test(`navigates to the review page, deletes and edit shipment`, async ({ customerPpmPage }) => {
      test.skip(true, 'This test fail due to navigateFromDateAndLocationPageToEstimatedWeightsPage()');
      const shipmentContainer = customerPpmPage.page.locator('[data-testid="ShipmentContainer"]').last();
      await customerPpmPage.deleteShipment(shipmentContainer, 1);

      // combining test for
      // navigates to the review page after finishing editing the PPM
      // shipment
      await customerPpmPage.navigateToHomePage();
      await customerPpmPage.navigateFromHomePageToExistingPPMDateAndLocationPage();
      await customerPpmPage.navigateFromDateAndLocationPageToEstimatedWeightsPage();
      await customerPpmPage.navigateFromEstimatedWeightsPageToEstimatedIncentivePage();
      await customerPpmPage.navigateFromEstimatedIncentivePageToAdvancesPage();
      await customerPpmPage.navigateFromAdvancesPageToReviewPage({ isMobile });
      await customerPpmPage.verifyPPMShipmentCard(fullPPMShipmentFields, { isEditable: true });
      // other tests submit the move otherwise we'd have an excessive number of moves
      await customerPpmPage.navigateToAgreementAndSign();
    });

    test('navigates to review page from home page and submits the move', async ({ customerPpmPage }) => {
      await customerPpmPage.verifyPPMShipmentCard(fullPPMShipmentFields, { isEditable: true });
      await customerPpmPage.navigateToAgreementAndSign();
      await customerPpmPage.submitMove();
      await customerPpmPage.navigateFromHomePageToReviewPage({ isMoveSubmitted: true });
      await customerPpmPage.verifyPPMShipmentCard(fullPPMShipmentFields, { isEditable: false });
      await customerPpmPage.navigateFromReviewPageToHomePage();
    });
  });
});

test.describe('(MultiMove) PPM Onboarding - Review', () => {
  test.skip(multiMoveEnabled === 'false', 'Skip if MultiMove workflow is not enabled.');
  let testMove;

  forEachViewport(async ({ isMobile }) => {
    test.beforeEach(async ({ customerPpmPage }) => {
      const move = await customerPpmPage.testHarness.buildUnsubmittedMoveWithMultipleFullPPMShipmentComplete();
      testMove = move;
      await customerPpmPage.signInForPPMWithMove(move);
      await customerPpmPage.clickOnGoToMoveButton();
      await customerPpmPage.navigateFromHomePageToReviewPage();
    });

    test(`navigates to the review page, deletes and edit shipment`, async ({ customerPpmPage }) => {
      const shipmentContainer = customerPpmPage.page.locator('[data-testid="ShipmentContainer"]').last();
      await customerPpmPage.deleteShipment(shipmentContainer, 1);

      // combining test for
      // navigates to the review page after finishing editing the PPM
      // shipment
      await customerPpmPage.navigateToHomePage();
      await customerPpmPage.navigateFromHomePageToExistingPPMDateAndLocationPage();
      await customerPpmPage.navigateFromDateAndLocationPageToEstimatedWeightsPage();
      await customerPpmPage.navigateFromEstimatedWeightsPageToEstimatedIncentivePage();
      await customerPpmPage.navigateFromEstimatedIncentivePageToAdvancesPage();
      await customerPpmPage.navigateFromAdvancesPageToReviewPage({ isMobile });
      await customerPpmPage.verifyPPMShipmentCard(fullPPMShipmentFields, { isEditable: true });
      // other tests submit the move otherwise we'd have an excessive number of moves
      await customerPpmPage.navigateToAgreementAndSign();
    });

    test('navigates to review page from home page and submits the move', async ({ customerPpmPage }) => {
      await customerPpmPage.verifyPPMShipmentCard(fullPPMShipmentFields, { isEditable: true });
      await customerPpmPage.navigateToAgreementAndSign();
      await customerPpmPage.submitMove();
      await customerPpmPage.navigateFromHomePageToReviewPage({ isMoveSubmitted: true });
      await customerPpmPage.verifyPPMShipmentCard(fullPPMShipmentFields, { isEditable: false });
      await customerPpmPage.navigateFromReviewPageToHomePageMM(testMove);
    });
  });
});
