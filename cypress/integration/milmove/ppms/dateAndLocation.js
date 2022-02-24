describe('the PPM flow', function () {
  before(() => {
    cy.prepareCustomerApp();
  });

  beforeEach(() => {
    cy.logout();
  });

  it('can submit a PPM move', () => {
    // profile2@complete.draft
    const userId = '4635b5a7-0f57-4557-8ba4-bbbb760c300a';
    cy.apiSignInAsUser(userId);
    SMSubmitsDateAndLocation();
  });
});
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

function SMSubmitsDateAndLocation() {
  cy.contains('PPM date & location');
  cy.get('input[name="pickupPostalCode"]').first().type('90210{enter}').blur();
  cy.get('input[type="radio"][id="no-secondary-pickup-postal-code"][value="false"]').eq(0).check('no', { force: true });

  cy.get('input[name="destinationPostalCode"]').clear().type('76127');
  cy.get('input[type="radio"][id="hasSecondaryDestinationPostalCodeNo"][value="false"]')
    .eq(0)
    .check('no', { force: true });

  cy.get('input[type="radio"][id="sitExpectedNo"][value="false"]').eq(0).check('no', { force: true });

  cy.nextPage();

  //   cy.location().should((loc) => {
  //     expect(loc.pathname).to.match(/^\/moves\/[^/]+\/review/);
  //   });
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
