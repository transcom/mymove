import {
  deleteShipment,
  navigateFromAdvancesPageToReviewPage,
  navigateFromDateAndLocationPageToEstimatedWeightsPage,
  navigateFromEstimatedIncentivePageToAdvancesPage,
  navigateFromEstimatedWeightsPageToEstimatedIncentivePage,
  navigateFromHomePageToReviewPage,
  navigateFromReviewPageToHomePage,
  navigateToAgreementAndSign,
  setMobileViewport,
  signInAndNavigateFromHomePageToExistingPPMDateAndLocationPage,
  signInAndNavigateFromHomePageToReviewPage,
} from '../../../support/ppmCustomerShared';

import { submitMove } from '../../mymove/utilities/customer';

const fullPPMShipmentFields = [
  ['Expected departure', '15 Mar 2020'],
  ['Origin ZIP', '90210'],
  ['Second origin ZIP', '90211'],
  ['Destination ZIP', '30813'],
  ['Second destination ZIP', '30814'],
  ['Storage expected? (SIT)', 'No'],
  ['Estimated weight', '4,000 lbs'],
  ['Pro-gear', 'Yes, 1,987 lbs'],
  ['Spouse pro-gear', 'Yes, 498 lbs'],
  ['Estimated incentive', '$10,000'],
  ['Advance requested?', 'Yes, $5,987'],
];

describe('PPM Onboarding - Review', function () {
  before(() => {
    cy.prepareCustomerApp();
  });

  beforeEach(() => {
    cy.intercept('GET', '**/internal/moves/**/mto_shipments').as('getShipment');
    cy.intercept('PATCH', '**/internal/mto-shipments/**').as('patchShipment');
    cy.intercept('DELETE', '**/internal/mto-shipments/**').as('deleteShipment');
    cy.intercept('GET', '**/internal/moves/**/signed_certifications').as('signedCertifications');
  });

  const viewportType = [
    { viewport: 'desktop', isMobile: false, userId: 'afcc7029-4810-4f19-999a-2b254c659e19' }, // multiComplete@ppm.unsubmitted
    { viewport: 'mobile', isMobile: true, userId: '836d8363-1a5a-45b7-aee0-996a97724c24' }, // multiComplete2@ppm.unsubmitted
  ];

  viewportType.forEach(({ viewport, isMobile, userId }) => {
    it(`navigates to the review page and deletes a shipment -  ${viewport}`, () => {
      if (isMobile) {
        setMobileViewport();
      }
      signInAndNavigateFromHomePageToReviewPage(userId);
      deleteShipmentFromCard();
    });
  });

  viewportType.forEach(({ viewport, isMobile }) => {
    it(`navigates to the review page after finishing editing the PPM shipment - ${viewport}`, () => {
      if (isMobile) {
        setMobileViewport();
      }

      // multiComplete@ppm.unsubmitted
      const userId = 'afcc7029-4810-4f19-999a-2b254c659e19';

      getToReviewPage(isMobile, userId);
      verifyPPMShipmentCard(fullPPMShipmentFields, true);
      navigateToAgreementAndSign(); // other tests submit the move otherwise we'd have an excessive number of moves
    });
  });

  viewportType.forEach(({ viewport, isMobile, userId }) => {
    it(`navigates to review page from home page and submits the move - ${viewport}`, () => {
      if (isMobile) {
        setMobileViewport();
      }

      signInAndNavigateFromHomePageToReviewPage(userId);
      verifyPPMShipmentCard(fullPPMShipmentFields, true);
      navigateToAgreementAndSign();
      submitMove('@signedCertifications');
      navigateFromHomePageToReviewPage(true);
      verifyPPMShipmentCard(fullPPMShipmentFields, false);
      navigateFromReviewPageToHomePage();
    });
  });
});

function getToReviewPage(isMobile = false, userId) {
  signInAndNavigateFromHomePageToExistingPPMDateAndLocationPage(userId);
  navigateFromDateAndLocationPageToEstimatedWeightsPage('@patchShipment');
  navigateFromEstimatedWeightsPageToEstimatedIncentivePage();
  navigateFromEstimatedIncentivePageToAdvancesPage();
  navigateFromAdvancesPageToReviewPage(isMobile);
}

function verifyPPMShipmentCard(shipmentCardFields, isEditable = false) {
  cy.get('h2').contains('Move setup').as('moveSetup');

  cy.get('@moveSetup')
    .next()
    .within(() => {
      cy.get('h4').contains('PPM 1');

      if (isEditable) {
        cy.get('button').contains('Edit');
        cy.get('button').contains('Delete');
      } else {
        cy.get('[data-testid="ShipmentContainer"]').find('button').should('not.exist');
      }

      cy.get('dt').should('have.length', shipmentCardFields.length).as('shipmentLabels');
      cy.get('dd').should('have.length', shipmentCardFields.length).as('shipmentValues');

      shipmentCardFields.forEach((shipmentField, index) => {
        cy.get('@shipmentLabels').eq(index).contains(shipmentField[0]);
        cy.get('@shipmentValues').eq(index).contains(shipmentField[1]);
      });
    });
}

function deleteShipmentFromCard() {
  cy.get('[data-testid="ShipmentContainer"]').last().as('shipmentContainer');
  deleteShipment('@shipmentContainer', 1);
}
