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

  // Refresh browser and make sure changes persist
  cy.patientReload();

  cy.get('.storage-in-transit').should($div => {
    const text = $div.text();
    expect(text).to.include('Dates');
    expect(text).to.include('24-Oct-2018');
    expect(text).to.include('Warehouse');
    expect(text).to.include('warehouse haus');
    expect(text).to.include('Warehouse ID');
    expect(text).to.include('SIT123456SIT');
    expect(text).to.include('Contact info');
    expect(text).to.include('123 Anystreet St.');
    expect(text).to.include('Citycitycity');
    expect(text).to.include('State');
    expect(text).to.include('94703');
  });
}
