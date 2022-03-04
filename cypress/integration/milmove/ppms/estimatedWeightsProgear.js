describe('the PPM flow', function () {
  before(() => {
    cy.prepareCustomerApp();
  });

  beforeEach(() => {
    cy.logout();
  });

  // profile@comple.te
  const userId = '3b9360a3-3304-4c60-90f4-83d687884070';
  it('doesn’t allow SM to progress if form is in an invalid state', () => {
    cy.apiSignInAsUser(userId);
    customerChoosesAPPMMove();
    submitsDateAndLocation();
    invalidInputs();
  });

  it('can submit a PPM move', () => {
    cy.apiSignInAsUser(userId);
    customerChoosesAPPMMove();
    submitsDateAndLocation();
  });
});

function customerChoosesAPPMMove() {
  cy.get('button[data-testid="shipment-selection-btn"]').click();
  cy.nextPage();

  cy.get('input[value="PPM"]').check({ force: true });
  cy.nextPage();
}

function submitsDateAndLocation() {
  cy.get('input[name="estimatedWeight"]').clear().type(900).blur();

  cy.get('button').contains('Save & Continue').click();

  cy.location().should((loc) => {
    expect(loc.pathname).to.match(/^\/moves\/[^/]+\/shipments\/[^/]+\/estimated-weight/);
  });
}

function invalidInputs() {
  cy.contains('Estimated Weight');
  cy.url().should('include', '/estimated-weight');

  // missing required weight
  cy.get('input[name="estimatedWeight"]').clear().blur();
  cy.get('[class="usa-error-message"]').as('errorMessage');
  cy.get('@errorMessage').contains('Required');
  cy.get('@errorMessage').next('input').should('have.id', 'estimatedWeight');

  // estimated weight violates min
  cy.get('input[name="estimatedWeight"]').type(0).blur();
  cy.get('@errorMessage').contains('Enter a weight greater than 0 lbs');
  cy.get('@errorMessage').next('input').should('have.id', 'estimatedWeight');
  cy.get('input[name="estimatedWeight"]').clear().type(500).blur();
  cy.get('@errorMessage').should('not.exist');

  // a warning is displayed when estimated weight is greater than the SM's weight allowance
  cy.get('input[name="estimatedWeight"]').clear().type(17000).blur();
  cy.get('[class="warning"]').as('warningMessage');
  cy.get('@warningMessage').contains(
    'This weight is more than your weight allowance. Talk to your counselor about what that could mean for your move.',
  );
  cy.get('@warningMessage').next('input').should('have.id', 'estimatedWeight');
  cy.get('input[name="estimatedWeight"]').clear().type(500).blur();
  cy.get('@warningMessage').should('not.exist');

  // pro gear violates max
  cy.get('input[name="hasProGear][value="true"]').check({ force: true });
  cy.get('input[name="proGearWeight"]').type(5000).blur();
  cy.get('@errorMessage').contains('Enter a weight less than 2,000 lbs');
  cy.get('@errorMessage').next('input').should('have.id', 'proGearWeight');
  cy.get('input[name="proGearWeight"]').clear().type(500).blur();
  cy.get('@errorMessage').should('not.exist');

  // When hasProGear is true show error if either personal or spouse pro gear isn't specified
  cy.get('input[name="proGearWeight"]').clear().blur();
  cy.get('@errorMessage').contains(
    "Enter a weight into at least one pro-gear field. If you won't have pro-gear, select No above.",
  );
  cy.get('@errorMessage').next('input').should('have.id', 'proGearWeight');
  cy.get('input[name="proGearWeight"]').clear().type(500).blur();
  cy.get('@errorMessage').should('not.exist');

  // spouse pro gear max violation
  cy.get('input[name="spouseProGearWeight"]').type(1000).blur();
  cy.get('@errorMessage').contains('Enter a weight less than 500 lbs');
  cy.get('@errorMessage').next('input').should('have.id', 'spouseProGearWeight');
  cy.get('input[name="spouseProGearWeight"]').clear().type(100).blur();
  cy.get('@errorMessage').should('not.exist');
}
