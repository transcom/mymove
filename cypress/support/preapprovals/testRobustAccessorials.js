import { createItemRequest } from '../../fixtures/preApprovals/requests/createRobustAccessorial';

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

export function clickAddARequest() {
  cy
    .get('.add-request')
    .contains('Add a request')
    .click();
}

export function clickSaveAndClose() {
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
