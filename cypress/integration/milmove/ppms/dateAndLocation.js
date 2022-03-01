describe('the PPM flow', function () {
  before(() => {
    cy.prepareCustomerApp();
  });

  beforeEach(() => {
    cy.logout();
  });

  it('can submit a PPM move', () => {
    // profile@comple.te
    const userId = '3b9360a3-3304-4c60-90f4-83d687884070';
    cy.apiSignInAsUser(userId);
    customerChoosesAPPMMove();
    submitsDateAndLocation();
  });

  it('doesn’t allow SM to progress if an invalid postal code is provided"', () => {
    // profile@co.mple.te
    const userId = '3b9360a3-3304-4c60-90f4-83d687884070';
    cy.apiSignInAsUser(userId);
    resumePPM();
    inputsInvalidPostalCodes();
  });

  it('doesn’t allow SM to progress if required date field is not filled out"', () => {
    // profile@co.mple.te
    const userId = '3b9360a3-3304-4c60-90f4-83d687884070';
    cy.apiSignInAsUser(userId);
    resumePPM();
    requiredDateFieldMissing();
  });

  it('doesn’t allow SM to progress if date field is invalid"', () => {
    // profile@co.mple.te
    const userId = '3b9360a3-3304-4c60-90f4-83d687884070';
    cy.apiSignInAsUser(userId);
    resumePPM();
    invalidDate();
  });

  it('doesn’t allow SM to progress if having a secondary zip is indicated but not provided"', () => {
    // profile@co.mple.te
    const userId = '3b9360a3-3304-4c60-90f4-83d687884070';
    cy.apiSignInAsUser(userId);
    resumePPM();
    missingSecondaryZip();
  });
});

function customerChoosesAPPMMove() {
  cy.get('button[data-testid="shipment-selection-btn"]').click();
  cy.nextPage();

  cy.get('input[type="radio"]').eq(1).check({ force: true });
  cy.nextPage();
}

function resumePPM() {
  cy.get('[data-testid="shipment-list-item-container').click();
}

function submitsDateAndLocation() {
  cy.contains('PPM date & location');
  cy.get('input[name="pickupPostalCode"]').clear().type('90210').blur();

  cy.get('input[name="destinationPostalCode"]').clear().type('76127');
  cy.get('input[name="expectedDepartureDate"]').clear().type('01 Feb 2022').blur();

  cy.get('button').contains('Save & Continue').click();

  cy.location().should((loc) => {
    expect(loc.pathname).to.match(/^\/moves\/[^/]+\/shipments\/[^/]+\/estimated-weight/);
  });
}

function inputsInvalidPostalCodes() {
  cy.get('[class="usa-error-message"]').should('not.exist');
  cy.get('input[name="pickupPostalCode"]').clear().type('00000').blur();
  cy.get('[class="usa-error-message"]').should('exist');

  // Fill in otherwise required fields
  cy.get('input[name="destinationPostalCode"]').clear().type('76127');
  cy.get('input[name="expectedDepartureDate"]').clear().type('01 Feb 2022').blur();

  cy.get('button').contains('Save & Continue').should('be.disabled');
}

function requiredDateFieldMissing() {
  // Fill in required fields other than the departure date field
  cy.get('input[name="pickupPostalCode"]').clear().type('90210').blur();
  cy.get('input[name="destinationPostalCode"]').clear().type('76127');

  cy.get('button').contains('Save & Continue').should('be.disabled');
}

function invalidDate() {
  cy.contains('PPM date & location');
  cy.get('input[name="pickupPostalCode"]').clear().type('90210').blur();

  cy.get('input[name="destinationPostalCode"]').clear().type('76127');
  cy.get('input[name="expectedDepartureDate"]').clear().type('01 ZZZ 20222').blur();
  cy.get('[class="usa-error-message"]').should('exist');

  cy.get('button').contains('Save & Continue').should('be.disabled');
}

function missingSecondaryZip() {
  cy.contains('PPM date & location');
  cy.get('input[name="pickupPostalCode"]').clear().type('90210').blur();
  cy.get('input[name="hasSecondaryPickupPostalCode"]').eq(0).check({ force: true });

  cy.get('input[name="destinationPostalCode"]').clear().type('76127');
  cy.get('input[name="expectedDepartureDate"]').clear().type('01 Feb 2022').blur();

  cy.get('button').contains('Save & Continue').should('be.disabled');
}
