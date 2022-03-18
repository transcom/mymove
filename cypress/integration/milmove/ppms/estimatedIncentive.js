describe('PPM Onboarding - Estimated Incentive', function () {
  before(() => {
    cy.prepareCustomerApp();
  });

  beforeEach(() => {
    // cy.intercept('GET', '**/internal/moves/**/mto_shipments').as('getShipment');
    cy.intercept('PATCH', '**/internal/mto-shipments/**').as('patchShipment');
    cy.logout();
    const userId = '4512dc8c-c777-444e-b6dc-7971e398f2dc';
    cy.apiSignInAsUser(userId);
  });

  it('go to estimated incentives page', () => {
    navigateToEstimatedIncentivePageAndVerifyAmount();
    cy.get('[data-testid="ppm-estimated-incentive-next"]').should('not.be.disabled');
  });

  it('mobile - go to estimated incentives page', () => {
    cy.viewport(425, 875);
    navigateToEstimatedIncentivePageAndVerifyAmount();
    cy.get('[data-testid="ppm-estimated-incentive-next"]').should('not.be.disabled').should('have.css', 'order', '1');
  });
});

function navigateToEstimatedIncentivePageAndVerifyAmount() {
  cy.get('[data-testid="shipment-list-item-container"]').click();

  // cy.url()
  //   .should('include', '/edit')
  //   .then((url) => {
  //     const regex = /(.*)\/edit/;
  //     const estimatedIncentivePath = `${url.match(regex)[1]}/estimated-incentive`;
  //     cy.visit(estimatedIncentivePath);
  //   });

  cy.get('[data-testid="ppm-date-and-location-submit"]').contains('Save & Continue').click();
  cy.wait('@patchShipment');

  cy.get('[data-testid="ppm-estimated-weights-submit"]').click();
  cy.wait('@patchShipment');

  cy.get('[data-testid="ppm-estimated-incentive-amount-sentence"]').contains(/\$\d{1,3}(?:,\d{3})*? is/);
}
