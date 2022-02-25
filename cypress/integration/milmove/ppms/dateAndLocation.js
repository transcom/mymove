describe('the PPM flow', function () {
  before(() => {
    cy.prepareCustomerApp();
  });

  beforeEach(() => {
    cy.logout();
  });

  it('can submit a PPM move', () => {
    // profile@comple.te
    const userId = '3b9360a3-3304-4c60-90f4-83d687884077';
    cy.apiSignInAsUser(userId);
    customerChoosesAPPMMove();
    submitsDateAndLocation();
  });
});

function customerChoosesAPPMMove() {
  cy.get('button[data-testid="shipment-selection-btn"]').click();
  cy.nextPage();

  cy.get('input[type="radio"]').eq(1).check({ force: true });
  cy.nextPage();
}
//   it('doesn’t allow SM to progress if don’t have rate data for zips"', () => {
//     // profile@co.mple.te
//     const userId = 'f154929c-5f07-41f5-b90c-d90b83d5773d';
//     cy.apiSignInAsPpmUser(userId);
//     SMInputsInvalidPostalCodes();
//   });

//   it('doesn’t allow SM to progress if required fields are not filled out"', () => {
//     // profile@co.mple.te
//     const userId = 'f154929c-5f07-41f5-b90c-d90b83d5773d';
//     cy.apiSignInAsPpmUser(userId);
//     SMInputsInvalidPostalCodes();
//   });

//   it('doesn’t allow SM to progress if date fields are invalid"', () => {
//     // profile@co.mple.te
//     const userId = 'f154929c-5f07-41f5-b90c-d90b83d5773d';
//     cy.apiSignInAsPpmUser(userId);
//     SMInputsInvalidPostalCodes();
//   });

//   it('form is pre populated if SM has shipment data"', () => {
//     // profile@co.mple.te
//     const userId = 'f154929c-5f07-41f5-b90c-d90b83d5773d';
//     cy.apiSignInAsPpmUser(userId);
//     SMInputsInvalidPostalCodes();
//   });

function submitsDateAndLocation() {
  cy.contains('PPM date & location');
  cy.get('input[name="pickupPostalCode"]').first().type('90210{enter}').blur();

  cy.get('input[name="destinationPostalCode"]').clear().type('76127');
  cy.get('input[name="expectedDepartureDate"]').first().type('01 Feb 2022').blur();

  cy.contains('Save & Continue').click();

  cy.location().should((loc) => {
    expect(loc.pathname).to.match(/^\/moves\/[^/]+\/shipments\/[^/]+\/estimated-weight/);
  });
}

// function SMInputsInvalidPostalCodes() {
//   cy.contains('Continue Move Setup').click();
//   cy.location().should((loc) => {
//     expect(loc.pathname).to.match(/^\/moves\/[^/]+\/ppm-start/);
//   });
//   cy.get('.wizard-header').should('not.exist');
//   cy.get('input[name="original_move_date"]').type('6/3/2100').blur();
//   // test an invalid pickup zip code
//   cy.get('input[name="pickup_postal_code"]').clear().type('00000').blur();
//   cy.get('#pickup_postal_code-error').should('exist');

//   cy.get('input[name="pickup_postal_code"]').clear().type('80913');

//   // test an invalid destination zip code
//   cy.get('input[name="destination_postal_code"]').clear().type('00000').blur();
//   cy.get('#destination_postal_code-error').should('exist');

//   cy.get('input[name="destination_postal_code"]').clear().type('30813');
//   cy.nextPage();

//   cy.location().should((loc) => {
//     expect(loc.pathname).to.match(/^\/moves\/[^/]+\/ppm-start/);
//   });
// }
