/* global cy */
export function fillAndSaveStorageInTransit() {
  // Select the location
  cy.get('select[name="location"]').select('ORIGIN');

  // Enter details in form and create the Storage In Transit request
  cy

    .get('input[name="estimated_start_date"]')

    .type('10/24/2018')

    .blur();

  cy
    .get('textarea[name="notes"]')
    .first()
    .type('notes notes', { force: true, delay: 150 });

  cy
    .get('input[name="warehouse_id"]')
    .first()
    .type('SIT123456SIT', { force: true, delay: 150 });

  cy
    .get('input[name="warehouse_name"]')
    .first()
    .type('warehouse haus', { force: true, delay: 150 });

  cy
    .get('input[name="warehouse_address.street_address_1"]')
    .first()
    .type('123 Anystreet St.', { force: true, delay: 150 });

  cy
    .get('input[name="warehouse_address.city"]')
    .first()
    .type('Citycitycity', { force: true, delay: 150 });

  cy
    .get('input[name="warehouse_address.state"]')
    .first()
    .type('State', { force: true, delay: 150 });

  cy
    .get('input[name="warehouse_address.postal_code"]')
    .first()
    .type('94703', { force: true, delay: 150 });

  cy
    .get('button')
    .contains('Send Request')
    .should('be.enabled');

  cy
    .get('button')
    .contains('Send Request')
    .click();
}
