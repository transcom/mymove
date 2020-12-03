describe('A customer following HHG Setup flow', function () {
  before(() => {
    cy.prepareCustomerApp();
  });

  beforeEach(() => {
    cy.intercept('POST', '/internal/service_members').as('createServiceMember');
    cy.intercept('PATCH', '**/internal/mto-shipments/**').as('patchShipment');
    cy.intercept('GET', '/internal/moves/**/mto_shipments').as('getMTOShipments');
    cy.intercept('GET', '/internal/users/logged_in').as('getLoggedInUser');
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
  cy.get(`[data-testid="mailingAddress1"]`).first().should('be.empty');
  cy.get(`[data-testid="city"]`).first().should('be.empty');
  cy.get(`[data-testid="state"]`).first().should('be.empty');
  cy.get(`[data-testid="zip"]`).first().should('be.empty');

  // should have expected "Required" error for required fields
  cy.get(`[data-testid="mailingAddress1"]`).first().focus().blur();
  cy.get('[class="usa-error-message"]').contains('Required');
  cy.get(`[data-testid="mailingAddress1"]`).first().type('Some address');
  cy.get('[class="usa-error-message"]').should('not.exist');

  cy.get(`[data-testid="city"]`).first().focus().blur();
  cy.get('[class="usa-error-message"]').contains('Required');
  cy.get(`[data-testid="city"]`).first().type('Some city');
  cy.get('[class="usa-error-message"]').should('not.exist');

  cy.get(`[data-testid="state"]`).first().focus().blur();
  cy.get('[class="usa-error-message"]').contains('Required');
  cy.get(`[data-testid="state"]`).first().type('CA');
  cy.get('[class="usa-error-message"]').should('not.exist');

  cy.get(`[data-testid="zip"]`).first().focus().blur();
  cy.get('[class="usa-error-message"]').contains('Required');
  cy.get(`[data-testid="zip"]`).first().type('9').blur();
  cy.get('[class="usa-error-message"]').contains('Must be valid zip code');
  cy.get(`[data-testid="zip"]`).first().type('1111').blur();
  cy.get('[class="usa-error-message"]').should('not.exist');

  // Next button disabled
  cy.get('button[data-testid="wizardNextButton"]').should('be.disabled');

  // overwrites data typed from above
  cy.get(`input[name="useCurrentResidence"]`).check({ force: true });

  // releasing agent
  cy.get(`[data-testid="firstName"]`).first().type('John');
  cy.get(`[data-testid="lastName"]`).first().type('Lee');
  cy.get(`[data-testid="phone"]`).first().type('9999999999');
  cy.get(`[data-testid="email"]`).first().type('ron').blur();
  cy.get('[class="usa-error-message"]').contains('Must be valid email');
  cy.get(`[data-testid="email"]`).first().type('@example.com');
  cy.get('[class="usa-error-message"]').should('not.exist');

  cy.get('button[data-testid="wizardNextButton"]').should('be.disabled');

  // requested delivery date
  cy.get('input[name="delivery.requestedDate"]').first().type('09/20/2020').blur();

  // checks has delivery address (default does not have delivery address)
  cy.get('input[type="radio"]').first().check({ force: true });

  // delivery location
  cy.get(`[data-testid="mailingAddress1"]`).last().type('412 Avenue M');
  cy.get(`[data-testid="mailingAddress2"]`).last().type('#3E');
  cy.get(`[data-testid="city"]`).last().type('Los Angeles');
  cy.get(`[data-testid="state"]`).last().type('CA');
  cy.get(`[data-testid="zip"]`).last().type('91111').blur();

  // releasing agent
  cy.get(`[data-testid="firstName"]`).last().type('John');
  cy.get(`[data-testid="lastName"]`).last().type('Lee');
  cy.get(`[data-testid="phone"]`).last().type('9999999999');
  cy.get(`[data-testid="email"]`).last().type('ron').blur();
  cy.get('[class="usa-error-message"]').contains('Must be valid email');
  cy.get(`[data-testid="email"]`).last().type('@example.com');
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
    expect(loc.pathname).to.match(/^\/moves\/[^/]+\/mto-shipments\/[^/]+\/edit-shipment/);
  });

  cy.get(`[data-testid="firstName"]`).last().type('Johnson').blur();

  // Ensure remarks is displayed in form
  cy.get(`[data-testid="remarks"]`).should('have.value', 'some customer remark');

  // Edit remarks and agent info
  cy.get(`[data-testid="remarks"]`).clear().type('some edited customer remark');
  cy.get(`[data-testid="email"]`).last().clear().type('John@example.com').blur();
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
  cy.get('h1').contains('Now for the official part...');
  cy.get('input[name="signature"]').type('Signature');
  cy.get('button').contains('Complete').click();
  cy.get('.usa-alert--success').within(() => {
    cy.contains('Congrats - your move is submitted!');
    cy.contains('Next, wait for approval. Once approved:');
    cy.get('a').contains('PPM info sheet').should('have.attr', 'href').and('include', '/downloads/ppm_info_sheet.pdf');
  });
}
