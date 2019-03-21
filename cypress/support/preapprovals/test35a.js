import { addLegacyRequest, clickSaveAndClose, clickAddARequest } from './testRobustAccessorials';

/* global cy */
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

export function test35ALegacy() {
  cy.selectQueueItemMoveLocator('DATESP');
  addLegacyRequest({ code: '35A', quantity1: 12 }).then(res => {
    expect(res.status).to.equal(201);
  });

  // must reload page because original 35A are added by cy.request()
  cy.reload();
  cy.get('td[details-cy="35A-default-details"]').should('contain', '12.0000 notes notes 35A');
}

export function test35A() {
  cy.selectQueueItemMoveLocator('DATESP');

  add35A({});
  cy
    .get('td[details-cy="35A-details"]')
    .should('contain', 'description description 35A reason reason 35A Est. not to exceed: $250.00 Actual amount: --');
}
