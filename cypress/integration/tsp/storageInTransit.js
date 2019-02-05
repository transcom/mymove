/* global cy */
describe('TSP user interacts with storage in transit panel', function() {
  beforeEach(() => {
    cy.signIntoTSP();
  });
  it('TSP user creates storage in transit request', function() {
    tspUserCreatesSitRequest();
  });
});

function tspUserCreatesSitRequest() {
  // Open new shipments queue
  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/new/);
  });

  // Find shipment and open it
  cy
    .get('div')
    .contains('DATESP') // TODO: Need our own shipment to work with?
    .dblclick();

  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/shipments\/[^/]+/);
  });

  // Click on Request SIT and see SIT Request form
  cy
    .get('.storage-in-transit-panel .add-request')
    .contains('Request SIT')
    .click()
    .get('.storage-in-transit-request-form')
    .should($div => {
      const text = $div.text();
      expect(text).to.include('SIT Location');
      expect(text).to.include('Warehouse ID number');
      expect(text).to.include('Warehouse Name');
      expect(text).to.include('Address Line 1');
    })
    .get('.usa-button-secondary')
    .contains('Cancel')
    .click()
    .get('.storage-in-transit-panel .add-request')
    .should($div => {
      const text = $div.text();
      expect(text).to.not.include('Sit Location');
    });
}
