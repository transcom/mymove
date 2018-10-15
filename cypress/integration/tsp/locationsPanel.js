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
      expect(text).to.include(`${address.city}, ${address.state} ${address.postal_code}`);
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
      expect(text).to.include(`${address.city}, ${address.state} ${address.postal_code}`);
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
      expect(text).to.equal(`${address.city}, ${address.state} ${address.postal_code}`);
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
  cy.get('div')
    .contains(shipmentId)
    .dblclick();

  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/shipments\/[^/]+/);
  });

  // Expect Customer Info to be loaded
  cy.contains('Locations')
    .parents('.editable-panel')
    .within(() => {
      cy.contains(type)
        .parent('.editable-panel-column')
        .children('.panel-field')
        .children('.field-value')
        .should($div => {
          expectation($div.text());
        });
    });
}

describe('TSP User Completes Locations Panel', function() {
  beforeEach(() => {
    cy.signIntoTSP();
  });
  it('tsp user completes locations panel', function() {
    tspUserEntersLocations();
  });
});

function tspUserEntersLocations() {
  const deliveryAddress = {
    street_1: '500 Something Avenue',
    city: 'Grandfather',
    state: 'ID',
    postal_code: '99999',
  };
  const pickupAddress = {
    street_1: '1 Main Street',
    city: 'Utopia',
    state: 'MT',
    postal_code: '11111',
  };
  const secondaryPickupAddress = {
    street_1: '666 Diagon Alley',
    city: 'London',
    state: 'NJ',
    postal_code: '66666',
  };

  // Open new shipments queue
  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/new/);
  });

  // Find shipment and open it
  cy.get('div')
    .contains('BACON1')
    .dblclick();

  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/shipments\/[^/]+/);
  });

  cy.get('.editable-panel-header')
    .contains('Locations')
    .siblings()
    .click();

  // Enter details in form and save dates
  cy.get('input[name="pickupAddress.street_address_1"]')
    .first()
    .clear()
    .type(pickupAddress.street_1)
    .blur();
  cy.get('input[name="pickupAddress.city"]')
    .first()
    .clear()
    .type(pickupAddress.city)
    .blur();
  cy.get('input[name="pickupAddress.state"]')
    .first()
    .clear()
    .type(pickupAddress.state)
    .blur();
  cy.get('input[name="pickupAddress.postal_code"]')
    .first()
    .clear()
    .type(pickupAddress.postal_code)
    .blur();

  // Secondary Pickup Address isn't actually required, so please remove later.
  cy.get('input[name="secondaryPickupAddress.street_address_1"]')
    .first()
    .clear()
    .type(secondaryPickupAddress.street_1)
    .blur();
  cy.get('input[name="secondaryPickupAddress.street_address_2"]')
    .first()
    .clear();
  cy.get('input[name="secondaryPickupAddress.city"]')
    .first()
    .clear()
    .type(secondaryPickupAddress.city)
    .blur();
  cy.get('input[name="secondaryPickupAddress.state"]')
    .first()
    .clear()
    .type(secondaryPickupAddress.state)
    .blur();
  cy.get('input[name="secondaryPickupAddress.postal_code"]')
    .first()
    .clear()
    .type(secondaryPickupAddress.postal_code)
    .blur();

  cy.get('input[name="deliveryAddress.street_address_1"]')
    .first()
    .clear()
    .type(deliveryAddress.street_1)
    .blur();
  cy.get('input[name="deliveryAddress.city"]')
    .first()
    .clear()
    .type(deliveryAddress.city)
    .blur();
  cy.get('input[name="deliveryAddress.state"]')
    .first()
    .clear()
    .type(deliveryAddress.state)
    .blur();
  cy.get('input[name="deliveryAddress.postal_code"]')
    .first()
    .clear()
    .type(deliveryAddress.postal_code)
    .blur();

  cy.get('button')
    .contains('Save')
    .should('be.enabled');

  cy.get('button')
    .contains('Save')
    .click();

  // Refresh browser and make sure changes persist
  cy.reload();

  cy.contains('Locations')
    .parents('.editable-panel')
    .within(() => {
      cy.contains('Delivery')
        .parent('.editable-panel-column')
        .children('.panel-field')
        .children('.field-value')
        .should($div => {
          const text = $div.text();
          expect(text).to.include(deliveryAddress.street_1);
          expect(text).to.include(
            `${deliveryAddress.city}, ${deliveryAddress.state} ${
              deliveryAddress.postal_code
            }`,
          );
        });
    });

  cy.contains('Locations')
    .parents('.editable-panel')
    .within(() => {
      cy.contains('Pickup')
        .parent('.editable-panel-column')
        .children('.panel-field')
        .children('.field-value')
        .should($div => {
          const text = $div.text();
          expect(text).to.include(pickupAddress.street_1);
          expect(text).to.include(
            `${pickupAddress.city}, ${pickupAddress.state} ${
              pickupAddress.postal_code
            }`,
          );
        });
    });

  cy.contains('Locations')
    .parents('.editable-panel')
    .within(() => {
      cy.contains('Pickup')
        .parent('.editable-panel-column')
        .children('.panel-field')
        .children('.field-value')
        .should($div => {
          const text = $div.text();
          expect(text).to.include(secondaryPickupAddress.street_1);
          expect(text).to.include(
            `${secondaryPickupAddress.city}, ${secondaryPickupAddress.state} ${
              secondaryPickupAddress.postal_code
            }`,
          );
        });
    });
}
