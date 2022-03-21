describe('PPM Onboarding - Estimated Incentive', function () {
  before(() => {
    cy.prepareCustomerApp();
  });

  beforeEach(() => {
    cy.intercept('GET', '**/internal/moves/**/mto_shipments').as('getShipment');
    cy.intercept('PATCH', '**/internal/mto-shipments/**').as('patchShipment');
  });

  it('go to estimated incentives page', () => {
    navigateToEstimatedIncentivePageAndVerifyFields(false);
  });

  it('mobile - go to estimated incentives page', () => {
    cy.viewport(479, 875);
    navigateToEstimatedIncentivePageAndVerifyFields(true);
  });
});

function navigateToEstimatedIncentivePageAndVerifyFields(isMobile = false) {
  // estimated_weights@ppm.unsubmitted
  const userId = '4512dc8c-c777-444e-b6dc-7971e398f2dc';
  cy.apiSignInAsUser(userId);
  cy.wait('@getShipment');

  cy.get('[data-testid="shipment-list-item-container"]').click();
  cy.wait('@getShipment');

  cy.get('[data-testid="ppm-date-and-location-submit"]').should('not.be.disabled').click();
  cy.wait('@patchShipment');

  cy.get('[data-testid="ppm-estimated-weights-submit"]').should('not.be.disabled').click();
  cy.wait('@patchShipment');

  // checks the format of the incentive amount statment is `$<some comma-separated number without decimals> is`
  cy.get('[data-testid="ppm-estimated-incentive-amount-sentence"]').contains(/\$\d{1,3}(?:,\d{3})*? is/);

  if (!isMobile) {
    cy.get('[data-testid="ppm-estimated-incentive-next"]').should('not.be.disabled');
  } else {
    cy.get('[data-testid="ppm-estimated-incentive-next"]').should('not.be.disabled').should('have.css', 'order', '1');
  }

  cy.get('[data-testid="ppm-estimated-incentive-next"]').should('be.enabled').click();
  cy.location().should((loc) => {
    expect(loc.pathname).to.match(/^\/moves\/[^/]+\/shipments\/[^/]+\/advances/);
  });
}
