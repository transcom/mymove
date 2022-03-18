describe('PPM Onboarding - Estimated Incentive', function () {
  before(() => {
    cy.prepareCustomerApp();
  });

  beforeEach(() => {
    cy.intercept('GET', '**/internal/moves/**/mto_shipments').as('getShipment');
    cy.intercept('PATCH', '**/internal/mto-shipments/**').as('patchShipment');

    const userId = '4512dc8c-c777-444e-b6dc-7971e398f2dc';
    cy.apiSignInAsUser(userId);
    cy.wait('@getShipment');
  });

  it('go to estimated incentives page', () => {
    navigateToEstimatedIncentivePageAndVerifyAmount();
    cy.get('[data-testid="ppm-estimated-incentive-next"]').should('not.be.disabled');
  });

  it('mobile - go to estimated incentives page', () => {
    cy.viewport(479, 875);
    navigateToEstimatedIncentivePageAndVerifyAmount();
    cy.get('[data-testid="ppm-estimated-incentive-next"]').should('not.be.disabled').should('have.css', 'order', '1');
  });
});

function navigateToEstimatedIncentivePageAndVerifyAmount() {
  cy.get('[data-testid="shipment-list-item-container"]').click();
  cy.wait('@getShipment');

  cy.get('[data-testid="ppm-date-and-location-submit"]').should('not.be.disabled').click();
  cy.wait('@patchShipment');

  cy.get('[data-testid="ppm-estimated-weights-submit"]').should('not.be.disabled').click();
  cy.wait('@patchShipment');

  cy.get('[data-testid="ppm-estimated-incentive-amount-sentence"]').contains(/\$\d{1,3}(?:,\d{3})*? is/);
}
