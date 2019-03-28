import { addLegacyRequest, clickSaveAndClose, clickAddARequest } from './testRobustAccessorials';

/* global cy */
export function add226A({ actual_amount_dollars = 250 }) {
  clickAddARequest();
  cy.selectTariff400ngItem('226A');
  cy.get('select[name="location"]').select('ORIGIN');
  cy.typeInTextarea({ name: 'description', value: `description description 226A` });
  cy.typeInTextarea({ name: 'reason', value: `reason reason 226A` });
  cy.typeInInput({ name: 'actual_amount_cents', value: actual_amount_dollars });
  cy.typeInTextarea({ name: 'notes', value: 'notes notes' });
  clickSaveAndClose();
}

export function test226ALegacy() {
  cy.selectQueueItemMoveLocator('DATESP');
  addLegacyRequest({ code: '226A', quantity1: 12 }).then(res => {
    expect(res.status).to.equal(201);
  });

  // must reload page because original 226A are added by cy.request()
  cy.reload();
  cy.get('td[details-cy="226A-default-details"]').should('contain', '12.0000 notes notes 226A');
}

export function test226A() {
  cy.selectQueueItemMoveLocator('DATESP');

  add226A({});
  cy.get('td[details-cy="226A-details"]').should('contain', 'description description 226A reason reason 226A $250.00');
}
