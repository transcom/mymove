/* global cy */
export function addOrigional105be() {
  clickAddPAR();
  cy.selectTariff400ngItem('105B');
  enterInput('quantity_1', 12);
  enterTextarea('notes', 'notes notes 105B');
  clickSaveButtonAndWait();

  clickAddPAR();
  cy.selectTariff400ngItem('105E');
  enterInput('quantity_1', 90);
  enterTextarea('notes', 'notes notes 105E');
  clickSaveButtonAndWait();
}

export function add105be() {
  clickAddPAR();
  cy.selectTariff400ngItem('105B');
  enterTextarea('description', 'description description 105B');
  addDimensions('item', 30);
  addDimensions('crate');
  enterTextarea('notes', 'notes notes');
  clickSaveButtonAndWait();

  clickAddPAR();
  cy.selectTariff400ngItem('105E');
  enterTextarea('description', 'description description 105E');
  addDimensions('item', 40);
  addDimensions('crate', 50);
  enterTextarea('notes', 'notes notes');
  clickSaveButtonAndWait();
}

function clickAddPAR() {
  cy
    .get('.add-request')
    .contains('Add a request')
    .click();
}

function addDimensions(type, value = 25) {
  enterInput(`${type}_dimensions.length`, value);
  enterInput(`${type}_dimensions.width`, value);
  enterInput(`${type}_dimensions.height`, value);
}

function enterTextarea(type, value) {
  cy.get(`textarea[name="${type}"]`).type(value, { force: true });
}

function enterInput(type, value) {
  cy.get(`input[name="${type}"]`).type(value, { force: true });
}

function clickSaveButtonAndWait() {
  cy
    .get('button')
    .contains('Save & Close')
    .click();
  cy.wait('@accessorialsCheck');
}
