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
    const expectation = text => {
      expect(text).to.include(address.street_1);
      expect(text).to.include(address.street_2);
      expect(text).to.include(address.street_3);
      expect(text).to.include(
        `${address.city}, ${address.state} ${address.postal_code}`,
      );
    };
    tspUserViewsLocation({ shipmentId: 'BACON1', type: 'Pickup', expectation });
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
    const expectation = text => {
      expect(text).to.include(address.street_1);
      expect(text).to.include(address.street_2);
      expect(text).to.include(address.street_3);
      expect(text).to.include(
        `${address.city}, ${address.state} ${address.postal_code}`,
      );
    };
    tspUserViewsLocation({
      shipmentId: 'BACON1',
      type: 'Delivery',
      expectation,
    });
  });
  it('tsp user primary delivery location when delivery address does not exist', function() {
    const address = {
      city: 'Beverly Hills',
      state: 'CA',
      postal_code: '90210',
    };
    const expectation = text => {
      expect(text).to.equal(
        `${address.city}, ${address.state} ${address.postal_code}`,
      );
    };
    tspUserViewsLocation({
      shipmentId: 'DTYSTN',
      type: 'Delivery',
      expectation,
    });
  });
});

function tspUserViewsLocation({ shipmentId, type, expectation }) {
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
        .children('.panel-field')
        .children('.field-value')
        .should($div => {
          expectation($div.text());
        });
    });
}
