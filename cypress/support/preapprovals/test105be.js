import { addLegacyRequest, clickSaveAndClose, clickAddARequest } from './testRobustAccessorials';

/* global cy */
function addDimensions(name, value = 25) {
  cy.typeInInput({ name: `${name}_dimensions.length`, value: value });
  cy.typeInInput({ name: `${name}_dimensions.width`, value: value });
  cy.typeInInput({ name: `${name}_dimensions.height`, value: value });
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

export function test105beLegacy() {
  cy.selectQueueItemMoveLocator('DATESP');

  addLegacyRequest({ code: '105B', quantity1: 12 }).then(res => {
    expect(res.status).to.equal(201);
  });

  addLegacyRequest({ code: '105E', quantity1: 90 }).then(res => {
    expect(res.status).to.equal(201);
  });

  // must reload page because original 105B/E are added by cy.request()
  cy.reload();
  cy.get('td[details-cy="105B-default-details"]').should('contain', '12.0000 notes notes 105B');
  cy.get('td[details-cy="105E-default-details"]').should('contain', '90.0000 notes notes 105E');
}

export function test105be() {
  cy.selectQueueItemMoveLocator('DATESP');

  add105({ code: '105B', itemSize: 30, crateSize: 25 });
  cy
    .get('td[details-cy="105B-details"]')
    .should(
      'contain',
      'description description 105B Crate: 25" x 25" x 25" (9.04 cu ft) Item: 30" x 30" x 30" notes notes',
    );

  add105({ code: '105E', itemSize: 40, crateSize: 50 });
  cy
    .get('td[details-cy="105E-details"]')
    .should(
      'contain',
      'description description 105E Crate: 50" x 50" x 50" (72.33 cu ft) Item: 40" x 40" x 40" notes notes',
    );
}
