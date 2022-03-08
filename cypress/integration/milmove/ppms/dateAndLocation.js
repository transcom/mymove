import { submitsDateAndLocation, customerChoosesAPPMMove } from './shared';

describe('PPM Onboarding - Add dates and location flow', function () {
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
    invalidInputs();
  });

  it('can continue to next page', () => {
    cy.apiSignInAsUser(userId);
    customerChoosesAPPMMove();
    submitsDateAndLocation();
  });
});

function invalidInputs() {
  cy.contains('PPM date & location');
  cy.url().should('include', '/new-shipment');

  // invalid date
  cy.get('input[name="expectedDepartureDate"]').clear().type('01 ZZZ 20222').blur();
  cy.get('[class="usa-error-message"]').as('errorMessage');
  cy.get('@errorMessage').contains('Enter a complete date in DD MMM YYYY format (day, month, year).');
  cy.get('input[name="expectedDepartureDate"]').clear().type('01 Feb 2022').blur();
  cy.get('@errorMessage').should('not.exist');

  // invalid postal codes
  cy.get('input[name="pickupPostalCode"]').clear().type('00000').blur();
  cy.get('@errorMessage').contains(
    "We don't have rates for this ZIP code. Please verify that you have entered the correct one. Contact support if this problem persists.",
  );
  cy.get('input[name="pickupPostalCode"]').clear().type('90210').blur();
  cy.get('@errorMessage').should('not.exist');
  cy.get('input[name="pickupPostalCode"]').clear().blur();
  cy.get('@errorMessage').contains('Required');
  cy.get('@errorMessage').next('input').should('have.id', 'pickupPostalCode');
  cy.get('input[name="pickupPostalCode"]').clear().type('90210').blur();
  cy.get('@errorMessage').should('not.exist');

  // missing secondary pickup postal code
  cy.get('input[name="hasSecondaryPickupPostalCode"]').eq(0).check({ force: true });
  cy.get('input[name="secondaryPickupPostalCode"]').clear().blur();
  cy.get('@errorMessage').contains('Required');
  cy.get('@errorMessage').next('input').should('have.id', 'secondaryPickupPostalCode');
  cy.get('input[name="secondaryPickupPostalCode"]').clear().type('90210').blur();
  cy.get('@errorMessage').should('not.exist');

  // missing secondary destination postal code
  cy.get('input[name="hasSecondaryDestinationPostalCode"]').eq(0).check({ force: true });
  cy.get('input[name="secondaryDestinationPostalCode"]').clear().blur();
  cy.get('@errorMessage').contains('Required');
  cy.get('@errorMessage').next('input').should('have.id', 'secondaryDestinationPostalCode');
  cy.get('input[name="secondaryDestinationPostalCode"]').clear().type('90210').blur();
  cy.get('@errorMessage').should('not.exist');
}
