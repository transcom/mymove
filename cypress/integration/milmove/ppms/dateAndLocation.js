// import { UnsupportedZipCodePPMErrorMsg, InvalidZIPTypeError } from 'utils/validation';
const UnsupportedZipCodePPMErrorMsg =
  "We don't have rates for this ZIP code. Please verify that you have entered the correct one.  Contact support if this problem persists.";

const InvalidZIPTypeError = 'Enter a 5-digit ZIP code';

describe('the PPM flow', function () {
  before(() => {
    cy.prepareCustomerApp();
  });

  beforeEach(() => {
    cy.logout();
  });

  // it('can submit a PPM move', () => {
  //   // profile@comple.te
  //   const userId = '3b9360a3-3304-4c60-90f4-83d687884077';
  //   cy.apiSignInAsUser(userId);
  //   customerChoosesAPPMMove();
  //   submitsDateAndLocation();
  // });

  it('doesn’t allow SM to progress if invalid postal codes are provided"', () => {
    // profile@co.mple.te
    const userId = '3b9360a3-3304-4c60-90f4-83d687884077';
    cy.apiSignInAsPpmUser(userId);

    customerChoosesAPPMMove();
    inputsInvalidPostalCodes();
  });
});

function customerChoosesAPPMMove() {
  cy.get('button[data-testid="shipment-selection-btn"]').click();
  cy.nextPage();

  cy.get('input[type="radio"]').eq(1).check({ force: true });
  cy.nextPage();
}

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

function submitsDateAndLocation() {
  cy.contains('PPM date & location');
  cy.get('input[name="pickupPostalCode"]').first().type('90210{enter}').blur();

  cy.get('input[name="destinationPostalCode"]').clear().type('76127');
  cy.get('input[name="expectedDepartureDate"]').first().type('01 Feb 2022').blur();

  cy.get('button').contains('Save & Continue').click();

  cy.location().should((loc) => {
    expect(loc.pathname).to.match(/^\/moves\/[^/]+\/shipments\/[^/]+\/estimated-weight/);
  });
}

function inputsInvalidPostalCodes() {
  cy.get('input[name="pickupPostalCode"]').first().type('00000{enter}').blur();
  cy.contains(UnsupportedZipCodePPMErrorMsg);

  cy.get('input[name="destinationPostalCode"]').clear().type('761272343');
  cy.contains(InvalidZIPTypeError);
  cy.get('button').contains('Save & Continue').should('be.disabled');
}
