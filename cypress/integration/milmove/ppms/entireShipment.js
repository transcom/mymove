const {
  customerChoosesAPPMMove,
  submitsDateAndLocation,
  submitsEstimatedWeightsAndProgear,
  verifyEstimatedIncentivePage,
} = require('../../../support/ppmShared');

describe('Entire PPM onboarding flow', function () {
  before(() => {
    cy.prepareCustomerApp();
  });

  beforeEach(() => {
    cy.intercept('GET', '**/internal/moves/**/mto_shipments').as('getShipment');
    cy.intercept('POST', '**/internal/mto_shipments').as('createShipment');
    cy.intercept('PATCH', '**/internal/mto-shipments/**').as('patchShipment');
  });

  // TODO: need to change id
  const userId = '3b9360a3-3304-4c60-90f4-83d687884070';
  it('happy path with new shipment', () => {
    cy.apiSignInAsUser(userId);
    cy.wait('@getShipment');
    customerChoosesAPPMMove();
    submitsDateAndLocation();
    submitsEstimatedWeightsAndProgear();
    verifyEstimatedIncentivePage();
  });
});
