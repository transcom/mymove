/* global cy */
describe('TSP User Checks Shipment Locations', function() {
  beforeEach(() => {
    cy.signIntoTSP();
  });
  it('tsp user primary pickup location', function() {
    const address = {
      street_1: '123 Any Street',
      street_2: 'P.O. Box 12345',
      street_3: 'c/o Some Person',
      city: 'Beverly Hills',
      state: 'CA',
      postal_code: '90210',
    };
    tspUserViewsLocation({ shipmentId: 'BACON1', type: 'Pickup', address });
  });
  it('tsp user primary delivery location when delivery address exists', function() {
    const address = {
      street_1: '987 Any Avenue',
      street_2: 'P.O. Box 9876',
      street_3: 'c/o Some Person',
      city: 'Fairfield',
      state: 'CA',
      postal_code: '94535',
    };
    tspUserViewsLocation({ shipmentId: 'BACON1', type: 'Delivery', address });
  });
});

function tspUserViewsLocation({ shipmentId, type, address }) {
  // Open new shipments queue
  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/new/);
  });

  // Find a shipment and open it
  cy
    .get('div')
    .contains(shipmentId)
    .dblclick();

  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/shipments\/[^/]+/);
  });

  // Expect Customer Info to be loaded
  cy
    .contains('Locations')
    .parents('.editable-panel')
    .within(() => {
      cy
        .contains(type)
        .parent('.editable-panel-column')
        .should($div => {
          const text = $div.text();
          expect(text).to.include(address.street_1);
          expect(text).to.include(address.street_2);
          expect(text).to.include(address.street_3);
          expect(text).to.include(
            `${address.city}, ${address.state} ${address.postal_code}`,
          );
        });
    });
}
