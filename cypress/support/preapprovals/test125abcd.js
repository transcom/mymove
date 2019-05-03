import { addLegacyRequest, clickSaveAndClose, clickAddARequest } from './testRobustAccessorials';

/* global cy */
export function add125({ code }) {
  clickAddARequest();
  cy.selectTariff400ngItem(code);
  cy.get('select[name="location"]').select('ORIGIN');
  cy.typeInTextarea({ name: 'reason', value: `reason reason ${code}` });
  cy
    .get('input[name="date"]')
    .last()
    .type('4/29/2019{enter}')
    .blur();
  cy.typeInInput({ name: 'time', value: '0400J' });
  cy.typeInTextarea({ name: 'notes', value: `notes notes ${code}` });
  cy.typeInInput({ name: 'address.street_address_1', value: 'street address 1' });
  cy.typeInInput({ name: 'address.street_address_2', value: 'street address 2' });
  cy.typeInInput({ name: 'address.city', value: 'city' });
  cy.get('select[name="address.state"]').select('CA');
  cy.typeInInput({ name: 'address.postal_code', value: '90210' });
  clickSaveAndClose();
}

export function test125ABCDLegacy() {
  cy.selectQueueItemMoveLocator('DATESP');
  addLegacyRequest({ code: '125A', quantity1: 12 }).then(res => {
    expect(res.status).to.equal(201);
  });
  addLegacyRequest({ code: '125B', quantity1: 12 }).then(res => {
    expect(res.status).to.equal(201);
  });
  addLegacyRequest({ code: '125C', quantity1: 12 }).then(res => {
    expect(res.status).to.equal(201);
  });
  addLegacyRequest({ code: '125D', quantity1: 12 }).then(res => {
    expect(res.status).to.equal(201);
  });

  // must reload page because original 125 are added by cy.request()
  cy.reload();
  cy.get('td[details-cy="125A-default-details"]').should('contain', '12.0000 notes notes 125A');
  cy.get('td[details-cy="125B-default-details"]').should('contain', '12.0000 notes notes 125B');
  cy.get('td[details-cy="125C-default-details"]').should('contain', '12.0000 notes notes 125C');
  cy.get('td[details-cy="125D-default-details"]').should('contain', '12.0000 notes notes 125D');
}

export function test125ABCD() {
  cy.selectQueueItemMoveLocator('DATESP');

  add125({ code: '125A' });
  cy
    .get('td[data-cy="125A-details"]')
    .should(
      'contain',
      'reason reason 125A Date of service: 29-Apr-19 Time of service: 0400J street address 1 street address 2 city, CA 90210',
    );

  add125({ code: '125B' });
  cy
    .get('td[data-cy="125B-details"]')
    .should(
      'contain',
      'reason reason 125B Date of service: 29-Apr-19 Time of service: 0400J street address 1 street address 2 city, CA 90210',
    );

  add125({ code: '125C' });
  cy
    .get('td[data-cy="125C-details"]')
    .should(
      'contain',
      'reason reason 125C Date of service: 29-Apr-19 Time of service: 0400J street address 1 street address 2 city, CA 90210',
    );

  add125({ code: '125D' });
  cy
    .get('td[data-cy="125D-details"]')
    .should(
      'contain',
      'reason reason 125D Date of service: 29-Apr-19 Time of service: 0400J street address 1 street address 2 city, CA 90210',
    );
}
