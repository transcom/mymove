/* global cy */
export function addOrigional105b() {
  cy
    .get('.add-request')
    .contains('Add a request')
    .click();

  cy.selectTariff400ngItem('Pack Reg Crate');

  cy
    .get('input[name="quantity_1"]')
    .first()
    .type('12', { force: true });

  cy
    .get('textarea[name="notes"]')
    .first()
    .type('notes notes', { force: true });

  cy
    .get('button')
    .contains('Save & Close')
    .should('be.enabled');

  cy
    .get('button')
    .contains('Save & Close')
    .click();
}

export function add105b() {
  cy
    .get('.add-request')
    .contains('Add a request')
    .click();

  cy.selectTariff400ngItem('105B');

  cy
    .get('textarea[name="description"]')
    .first()
    .type('description txtfield', { force: true });

  addDimensions('item', 30);
  addDimensions('crate');

  cy
    .get('textarea[name="notes"]')
    .first()
    .type('notes notes', { force: true });

  cy
    .get('button')
    .contains('Save & Close')
    .should('be.enabled');

  cy
    .get('button')
    .contains('Save & Close')
    .click();
}

function addDimensions(type, value = 25) {
  cy
    .get(`input[name="${type}_dimensions.length"]`)
    .first()
    .type(value, { force: true });

  cy
    .get(`input[name="${type}_dimensions.width"]`)
    .first()
    .type(value, { force: true });

  cy
    .get(`input[name="${type}_dimensions.height"]`)
    .first()
    .type(value, { force: true });
}
