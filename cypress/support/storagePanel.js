/* global cy */

export function userSavesStorageDetails() {
  cy
    .get('.storage')
    .find('input[name="total_sit_cost"]')
    .first()
    .type('600');

  cy
    .get('.storage')
    .find('input[name="days_in_storage"]')
    .first()
    .type('60');

  cy
    .get('button')
    .contains('Save')
    .should('be.enabled');

  cy
    .get('button')
    .contains('Save')
    .click();

  cy.patientReload();

  cy.get('div.panel-field.total_sit_cost').contains('$600.00');
  cy.get('div.panel-field.days_in_storage').contains('60');
}

export function userCancelsStorageDetails() {
  cy
    .get('.storage')
    .find('input[name="total_sit_cost"]')
    .first()
    .type('600');

  cy
    .get('.storage')
    .find('input[name="days_in_storage"]')
    .first()
    .type('60');

  cy
    .get('button')
    .contains('Cancel')
    .should('be.enabled');

  cy
    .get('button')
    .contains('Cancel')
    .click();

  cy.get('.total_sit_cost .field-value').contains('$0');
  cy.get('.days_in_storage .field-value').contains('0');
}
