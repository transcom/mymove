/* global cy */
export function addOriginal105({ code, quantity1 }) {
  clickAddARequest();
  cy.selectTariff400ngItem(code);
  cy.typeInInput({ name: 'quantity_1', value: quantity1 });
  cy.typeInTextarea({ name: 'notes', value: `notes notes ${code}` });
  clickSaveAndClose();
}

export function add105({ code, itemSize = 25, crateSize = 25 }) {
  clickAddARequest();
  cy.selectTariff400ngItem(code);
  cy.typeInTextarea({ name: 'description', value: `description description ${code}` });
  addDimensions('item', itemSize);
  addDimensions('crate', crateSize);
  cy.typeInTextarea({ name: 'notes', value: 'notes notes' });
  clickSaveAndClose();
}

function clickAddARequest() {
  cy
    .get('.add-request')
    .contains('Add a request')
    .click();
}

function addDimensions(name, value = 25) {
  cy.typeInInput({ name: `${name}_dimensions.length`, value: value });
  cy.typeInInput({ name: `${name}_dimensions.width`, value: value });
  cy.typeInInput({ name: `${name}_dimensions.height`, value: value });
}

function clickSaveAndClose() {
  cy
    .get('button')
    .contains('Save & Close')
    .click();
}
