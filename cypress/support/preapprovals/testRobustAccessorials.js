import { createItemRequest } from '../../fixtures/preApprovals/requests/create105';

/* global cy */
export function addLegacyRequest({ code, quantity1 }) {
  return cy.location().then(loc => {
    // eslint-disable-next-line security/detect-non-literal-regexp
    const pattern = new RegExp(`^/shipments/(.*)`);
    expect(loc.pathname).to.match(/^\/shipments\/[^/]+/);
    const shipmentId = loc.pathname.match(pattern)[1];
    return cy.getCookie('masked_gorilla_csrf').then(token => {
      const csrfToken = token.value;
      return createLegacyRequest(shipmentId, csrfToken, code, quantity1);
    });
  });
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

function createLegacyRequest(shipmentId, csrfToken, code, quantity1) {
  const item = createItemRequest({
    shipmentId: shipmentId,
    csrfToken: csrfToken,
    code: code,
    quantity1: quantity1,
  });
  return cy.request(item);
}

export function add35A({ estimate_amount_dollars = 250, actual_amount_dollars }) {
  clickAddARequest();
  cy.selectTariff400ngItem('35A');
  cy.get('select[name="location"]').select('ORIGIN');
  cy.typeInTextarea({ name: 'description', value: `description description 35A` });
  cy.typeInTextarea({ name: 'reason', value: `reason reason 35A` });
  cy.typeInInput({ name: 'estimate_amount_cents', value: estimate_amount_dollars });
  if (actual_amount_dollars) {
    cy.typeInInput({ name: 'actual_amount_cents', value: actual_amount_dollars });
  }
  cy.typeInTextarea({ name: 'notes', value: 'notes notes' });
  clickSaveAndClose();
}
