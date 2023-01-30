/**
 * Semi-automated converted from a cypress test, and thus may contain
 * non best-practices, in particular: heavy use of `page.locator`
 * instead of `page.getBy*`.
 */

// @ts-check
const { test, useMobileViewport, CustomerPpmPage } = require('./customerPpmTestFixture');

const fullPPMShipmentFields = [
  ['Expected departure', '15 Mar 2020'],
  ['Origin ZIP', '90210'],
  ['Second origin ZIP', '90211'],
  ['Destination ZIP', '30813'],
  ['Second destination ZIP', '30814'],
  ['Closeout office', 'Base Ketchikan'],
  ['Storage expected? (SIT)', 'No'],
  ['Estimated weight', '4,000 lbs'],
  ['Pro-gear', 'Yes, 1,987 lbs'],
  ['Spouse pro-gear', 'Yes, 498 lbs'],
  ['Estimated incentive', '$10,000'],
  ['Advance requested?', 'Yes, $5,987'],
];

test.describe('PPM Onboarding - Review', () => {
  /** @type {CustomerPpmPage} */
  let customerPpmPage;

  test.beforeEach(async ({ customerPage }) => {
    const move = await customerPage.testHarness.buildUnsubmittedMoveWithMultipleFullPPMShipmentComplete();
    customerPpmPage = new CustomerPpmPage(customerPage, move);
    await customerPpmPage.signInAndNavigateFromHomePageToReviewPage();
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

      test(`navigates to the review page, deletes and edit shipment`, async () => {
        const shipmentContainer = customerPpmPage.page.locator('[data-testid="ShipmentContainer"]').last();
        await customerPpmPage.deleteShipment(shipmentContainer, 1);

        // combining test for
        // navigates to the review page after finishing editing the PPM shipment
        await customerPpmPage.signInAndNavigateFromHomePageToExistingPPMDateAndLocationPage();
        await customerPpmPage.navigateFromDateAndLocationPageToEstimatedWeightsPage();
        await customerPpmPage.navigateFromEstimatedWeightsPageToEstimatedIncentivePage();
        await customerPpmPage.navigateFromEstimatedIncentivePageToAdvancesPage();
        await customerPpmPage.navigateFromAdvancesPageToReviewPage({ isMobile });
        await customerPpmPage.verifyPPMShipmentCard(fullPPMShipmentFields, { isEditable: true });
        // other tests submit the move otherwise we'd have an excessive number of moves
        await customerPpmPage.navigateToAgreementAndSign();
      });

      test('navigates to review page from home page and submits the move', async () => {
        await customerPpmPage.verifyPPMShipmentCard(fullPPMShipmentFields, { isEditable: true });
        await customerPpmPage.navigateToAgreementAndSign();
        await customerPpmPage.submitMove();
        await customerPpmPage.navigateFromHomePageToReviewPage({ isMoveSubmitted: true });
        await customerPpmPage.verifyPPMShipmentCard(fullPPMShipmentFields, { isEditable: false });
        await customerPpmPage.navigateFromReviewPageToHomePage();
      });
    });
  });
});
