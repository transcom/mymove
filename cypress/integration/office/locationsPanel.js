/* global cy */
describe('Office User Checks Shipment Locations', function() {
  beforeEach(() => {
    cy.signIntoOffice();
  });

  it('office user primary pickup location', function() {
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

    officeUserViewsLocation({ shipmentId: 'BACON3', type: 'Pickup', expectation });
  });

  it('office user primary delivery location when delivery address exists', function() {
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

    officeUserViewsLocation({
      shipmentId: 'BACON3',
      type: 'Delivery',
      expectation,
    });
  });

  it('office user primary delivery location when delivery address does not exist', function() {
    const address = {
      city: 'Augusta',
      state: 'GA',
      postal_code: '30813',
    };
    const expectation = text => {
      expect(text).to.equal(`${address.city}, ${address.state} ${address.postal_code}`);
    };

    officeUserViewsLocation({
      shipmentId: 'DTYSTN',
      type: 'Delivery',
      expectation,
    });
  });
});

function officeUserViewsLocation({ shipmentId, type, expectation }) {
  // Open new shipments queue
  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/new/);
  });

  // Find a shipment and open it
  cy.selectQueueItemMoveLocator(shipmentId);

  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/new\/moves\/[^/]+\/basics/);
  });

  cy
    .get('a')
    .contains('HHG')
    .click(); // navtab

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

describe('office user Completes Locations Panel', function() {
  beforeEach(() => {
    cy.signIntoOffice();
  });
  it('office user completes locations panel', function() {
    officeUserEntersLocations();
  });
});

function officeUserEntersLocations() {
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
    postal_code: '66666-6666',
  };
  const newDutyStation = {
    city: 'Augusta',
    state: 'GA',
    postal_code: '30813',
  };

  // Open new shipments queue
  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/new/);
  });

  // Find shipment and open it
  cy.selectQueueItemMoveLocator('BACON3');

  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/new\/moves\/[^/]+\/basics/);
  });

  cy
    .get('a')
    .contains('HHG')
    .click(); // navtab

  cy
    .get('.editable-panel-header')
    .contains('Locations')
    .siblings()
    .click();

  // Enter details in form and save locations
  cy
    .get('input[name="pickup_address.street_address_1"]')
    .first()
    .clear()
    .type(pickupAddress.street_1)
    .blur();
  cy
    .get('input[name="pickup_address.city"]')
    .first()
    .clear()
    .type(pickupAddress.city)
    .blur();
  cy
    .get('input[name="pickup_address.state"]')
    .first()
    .clear()
    .type(pickupAddress.state)
    .blur();
  cy
    .get('input[name="pickup_address.postal_code"]')
    .first()
    .clear()
    .type('1002')
    .blur();
  // Shouldn't be able to save without 5 digit zip
  cy
    .get('button')
    .contains('Save')
    .should('be.disabled');
  cy
    .get('input[name="pickup_address.postal_code"]')
    .first()
    .clear()
    .type(pickupAddress.postal_code)
    .blur();

  // Make sure delivery address is not required
  cy
    .get('label[for="has_delivery_address"]')
    .parent()
    .find('[type="radio"][value="no"]')
    .check({ force: true });

  cy
    .get('button')
    .contains('Save')
    .should('be.enabled');

  // Set Secondary Pickup Address to required.
  cy
    .get('label[for="has_secondary_pickup_address"]')
    .parent()
    .find('[type="radio"][value="yes"]')
    .check({ force: true });
  cy
    .get('input[name="secondary_pickup_address.street_address_1"]')
    .first()
    .clear()
    .type(secondaryPickupAddress.street_1)
    .blur();
  cy
    .get('input[name="secondary_pickup_address.street_address_2"]')
    .first()
    .clear();
  cy
    .get('input[name="secondary_pickup_address.city"]')
    .first()
    .clear()
    .type(secondaryPickupAddress.city)
    .blur();
  cy
    .get('input[name="secondary_pickup_address.state"]')
    .first()
    .clear()
    .type(secondaryPickupAddress.state)
    .blur();
  cy
    .get('input[name="secondary_pickup_address.postal_code"]')
    .first()
    .clear()
    .type('1002')
    .blur();
  // Shouldn't be able to save without 5 digit zip
  cy
    .get('button')
    .contains('Save')
    .should('be.disabled');
  cy
    .get('input[name="secondary_pickup_address.postal_code"]')
    .first()
    .clear()
    .type(secondaryPickupAddress.postal_code)
    .blur();
  cy
    .get('button')
    .contains('Save')
    .should('be.enabled');

  // Allow secondary delivery address
  cy
    .get('label[for="has_delivery_address"]')
    .parent()
    .find('[type="radio"][value="yes"]')
    .check({ force: true });

  cy
    .get('input[name="delivery_address.street_address_1"]')
    .first()
    .clear()
    .type(deliveryAddress.street_1)
    .blur();
  cy
    .get('input[name="delivery_address.city"]')
    .first()
    .clear()
    .type(deliveryAddress.city)
    .blur();
  cy
    .get('input[name="delivery_address.state"]')
    .first()
    .clear()
    .type(deliveryAddress.state)
    .blur();
  cy
    .get('input[name="delivery_address.postal_code"]')
    .first()
    .clear()
    .type('1002')
    .blur();
  // Shouldn't be able to save without 5 digit zip
  cy
    .get('button')
    .contains('Save')
    .should('be.disabled');
  cy
    .get('input[name="delivery_address.postal_code"]')
    .first()
    .clear()
    .type(deliveryAddress.postal_code)
    .blur();
  cy
    .get('button')
    .contains('Save')
    .should('be.enabled');
  cy
    .get('button')
    .contains('Save')
    .should('be.enabled');

  cy
    .get('button')
    .contains('Save')
    .click();

  // Refresh browser and make sure changes persist
  cy.patientReload();

  cy
    .contains('Locations')
    .parents('.editable-panel')
    .within(() => {
      cy
        .contains('Delivery')
        .parent('.editable-panel-column')
        .children('.panel-field')
        .children('.field-value')
        .should($div => {
          const text = $div.text();
          expect(text).to.include(deliveryAddress.street_1);
          expect(text).to.include(`${deliveryAddress.city}, ${deliveryAddress.state} ${deliveryAddress.postal_code}`);
        });
    });

  cy
    .contains('Locations')
    .parents('.editable-panel')
    .within(() => {
      cy
        .contains('Pickup')
        .parent('.editable-panel-column')
        .children('.panel-field')
        .children('.field-value')
        .should($div => {
          const text = $div.text();
          expect(text).to.include(pickupAddress.street_1);
          expect(text).to.include(`${pickupAddress.city}, ${pickupAddress.state} ${pickupAddress.postal_code}`);
        });
    });

  cy
    .contains('Locations')
    .parents('.editable-panel')
    .within(() => {
      cy
        .contains('Pickup')
        .parent('.editable-panel-column')
        .children('.panel-field')
        .children('.field-value')
        .should($div => {
          const text = $div.text();
          expect(text).to.include(secondaryPickupAddress.street_1);
          expect(text).to.include(
            `${secondaryPickupAddress.city}, ${secondaryPickupAddress.state} ${secondaryPickupAddress.postal_code}`,
          );
        });
    });
  cy
    .get('.editable-panel-header')
    .contains('Locations')
    .siblings()
    .click();

  // Click every radio button, which means you'll end up with two 'No's selected
  cy
    .get('[type="radio"]')
    .eq(1)
    .check({ force: true });
  cy
    .get('[type="radio"]')
    .eq(3)
    .check({ force: true });

  cy
    .get('button')
    .contains('Save')
    .should('be.enabled');

  cy
    .get('button')
    .contains('Save')
    .click();

  // Refresh browser and make sure changes persist
  cy.patientReload();

  cy
    .contains('Locations')
    .parents('.editable-panel')
    .within(() => {
      cy
        .contains('Delivery')
        .parent('.editable-panel-column')
        .children('.panel-field')
        .children('.field-value')
        .should($div => {
          const text = $div.text();
          expect(text).to.include(`${newDutyStation.city}, ${newDutyStation.state} ${newDutyStation.postal_code}`);
        });
    });
}
