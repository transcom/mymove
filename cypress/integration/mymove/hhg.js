describe('A customer following HHG Setup flow', function () {
  before(() => {
    cy.prepareCustomerApp();
  });

  beforeEach(() => {
    cy.intercept('POST', '**/internal/service_members').as('createServiceMember');
    cy.intercept('PATCH', '**/internal/mto-shipments/**').as('patchShipment');
    cy.intercept('**/internal/moves/**/mto_shipments').as('getMTOShipments');
    cy.intercept('**/internal/users/logged_in').as('getLoggedInUser');
  });

  it('can create an HHG shipment, review and edit details, and submit their move', function () {
    // profile@comple.te
    const userId = '3b9360a3-3304-4c60-90f4-83d687884077';
    cy.apiSignInAsUser(userId);
    customerChoosesAnHHGMove();
    customerSetsUpAnHHGMove();
    customerReviewsMoveDetailsAndEditsHHG();
    customerSubmitsMove();
  });
});

function customerChoosesAnHHGMove() {
  cy.get('button[data-testid="shipment-selection-btn"]').click();
  cy.nextPage();
  cy.get('h2').contains('Choose 1 shipment at a time.');

  cy.get('input[type="radio"]').eq(1).check({ force: true });
  cy.nextPage();
}

function customerSetsUpAnHHGMove() {
  cy.get('button[data-testid="wizardNextButton"]').should('be.disabled');

  cy.get('input[name="pickup.requestedDate"]').focus().blur();
  cy.get('[class="usa-error-message"]').contains('Required');
  cy.get('input[name="pickup.requestedDate"]').type('08/02/2020').blur();
  cy.get('[class="usa-error-message"]').should('not.exist');

  cy.get('button[data-testid="wizardNextButton"]').should('be.disabled');

  // should be empty before using "Use current residence" checkbox
  cy.get(`input[name="pickup.address.street_address_1"]`).should('be.empty');
  cy.get(`input[name="pickup.address.city"]`).should('be.empty');
  cy.get(`input[name="pickup.address.postal_code"]`).should('be.empty');

  // should have expected "Required" error for required fields
  cy.get(`input[name="pickup.address.street_address_1"]`).focus().blur();
  cy.get('[class="usa-error-message"]').contains('Required');
  cy.get(`input[name="pickup.address.street_address_1"]`).type('Some address');
  cy.get('[class="usa-error-message"]').should('not.exist');

  cy.get(`input[name="pickup.address.city"]`).focus().blur();
  cy.get('[class="usa-error-message"]').contains('Required');
  cy.get(`input[name="pickup.address.city"]`).type('Some city');
  cy.get('[class="usa-error-message"]').should('not.exist');

  cy.get(`select[name="pickup.address.state"]`).focus().blur();
  cy.get('[class="usa-error-message"]').contains('Required');
  cy.get(`select[name="pickup.address.state"]`).select('CA');
  cy.get('[class="usa-error-message"]').should('not.exist');

  cy.get(`input[name="pickup.address.postal_code"]`).focus().blur();
  cy.get('[class="usa-error-message"]').contains('Required');
  cy.get(`input[name="pickup.address.postal_code"]`).type('9').blur();
  cy.get('[class="usa-error-message"]').contains('Must be valid zip code');
  cy.get(`input[name="pickup.address.postal_code"]`).type('1111').blur();
  cy.get('[class="usa-error-message"]').should('not.exist');

  // Next button disabled
  cy.get('button[data-testid="wizardNextButton"]').should('be.disabled');

  // overwrites data typed from above
  cy.get(`input[name="useCurrentResidence"]`).check({ force: true });

  // releasing agent
  cy.get(`input[name="pickup.agent.firstName"]`).type('John');
  cy.get(`input[name="pickup.agent.lastName"]`).type('Lee');
  cy.get(`input[name="pickup.agent.phone"]`).type('9999999999');
  cy.get(`input[name="pickup.agent.email"]`).type('ron').blur();
  cy.get('[class="usa-error-message"]').contains('Must be valid email');
  cy.get(`input[name="pickup.agent.email"]`).type('@example.com');
  cy.get('[class="usa-error-message"]').should('not.exist');

  cy.get('button[data-testid="wizardNextButton"]').should('be.disabled');

  // requested delivery date
  cy.get('input[name="delivery.requestedDate"]').first().type('09/20/2020').blur();

  // checks has delivery address (default does not have delivery address)
  cy.get('input[type="radio"]').first().check({ force: true });

  // delivery location
  cy.get(`input[name="delivery.address.street_address_1"]`).type('412 Avenue M');
  cy.get(`input[name="delivery.address.street_address_2"]`).type('#3E');
  cy.get(`input[name="delivery.address.city"]`).type('Los Angeles');
  cy.get(`select[name="delivery.address.state"]`).select('CA');
  cy.get(`input[name="delivery.address.postal_code"]`).type('91111').blur();

  // releasing agent
  cy.get(`input[name="delivery.agent.firstName"]`).type('John');
  cy.get(`input[name="delivery.agent.lastName"]`).type('Lee');
  cy.get(`input[name="delivery.agent.phone"]`).type('9999999999');
  cy.get(`input[name="delivery.agent.email"]`).type('ron').blur();
  cy.get('[class="usa-error-message"]').contains('Must be valid email');
  cy.get(`input[name="delivery.agent.email"]`).type('@example.com');
  cy.get('[class="usa-error-message"]').should('not.exist');

  // customer remarks
  cy.get(`[data-testid="remarks"]`).first().type('some customer remark');
  cy.nextPage();
}

function customerReviewsMoveDetailsAndEditsHHG() {
  cy.get('[data-testid="review-move-header"]').contains('Review your details');

  cy.get('[data-testid="ShipmentContainer"]').contains('HHG 1');

  cy.get('[data-testid="edit-shipment-btn"]').contains('Edit').click();

  cy.location().should((loc) => {
    expect(loc.pathname).to.match(/^\/moves\/[^/]+\/shipments\/[^/]+\/edit/);
  });

  cy.get(`input[name="delivery.agent.firstName"]`).type('Johnson').blur();

  // Ensure remarks is displayed in form
  cy.get(`[data-testid="remarks"]`).should('have.value', 'some customer remark');

  // Edit remarks and agent info
  cy.get(`[data-testid="remarks"]`).clear().type('some edited customer remark');
  cy.get(`input[name="delivery.agent.email"]`).clear().type('John@example.com').blur();
  cy.get('button').contains('Save').click();

  cy.wait('@patchShipment');

  cy.location().should((loc) => {
    expect(loc.pathname).to.match(/^\/moves\/[^/]+\/review/);
  });

  cy.get('[data-testid="hhg-summary"]').find('dl').contains('some edited customer remark');
  cy.get('[data-testid="hhg-summary"]').find('dl').contains('JohnJohnson Lee');

  cy.get('button').contains('Finish later').click();
  cy.get('h3').contains('Time to submit your move');
  cy.get('button').contains('Review and submit').click();

  cy.nextPage();
}

function customerSubmitsMove() {
  cy.get('h1').contains('Now for the official part');
  cy.get('input[name="signature"]').type('Signature');
  cy.get('button').contains('Complete').click();
  cy.get('.usa-alert--success').within(() => {
    cy.contains('You’ve submitted your move request.');
  });
}
