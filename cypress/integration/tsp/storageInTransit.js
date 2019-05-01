import { fillAndSaveStorageInTransit, editAndSaveStorageInTransit } from '../../support/testCreateStorageInTransit';

/* global cy */
describe('TSP user interacts with storage in transit panel', function() {
  beforeEach(() => {
    cy.signIntoTSP();
  });
  it('TSP user creates storage in transit request', function() {
    tspUserCreatesSitRequest();
  });
  it('TSP user starts and then cancels storage in transit request', function() {
    tspUserStartsAndCancelsSitRequest();
  });
  it('TSP user edits and saves storage in transit request', function() {
    tspUserEditsSitRequest();
  });
  it('TSP user starts and then cancels Place into SIT form', function() {
    tspUserStartsAndCancelsSitPlaceInSit();
  });
});

// need to simulate a form submit

function tspUserCreatesSitRequest() {
  // Open accepted shipments queue
  cy.patientVisit('/queues/accepted');

  // Open new shipments queue
  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/accepted/);
  });

  // Find shipment and open it
  cy.selectQueueItemMoveLocator('SITPAN');

  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/shipments\/[^/]+/);
  });

  // Click on Request SIT and see SIT Request form
  cy
    .get('.storage-in-transit-panel .add-request')
    .contains('Request SIT')
    .click()
    .get('.storage-in-transit-form')
    .should($div => {
      const text = $div.text();
      expect(text).to.include('SIT location');
      expect(text).to.include('Warehouse ID number');
      expect(text).to.include('Warehouse name');
      expect(text).to.include('Address line 1');
    });

  // fill out and submit the form
  fillAndSaveStorageInTransit();
}

function tspUserStartsAndCancelsSitRequest() {
  // Open accepted shipments queue
  cy.patientVisit('/queues/accepted');

  // Open new shipments queue
  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/accepted/);
  });

  // Find shipment and open it
  cy.selectQueueItemMoveLocator('SITPAN');

  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/shipments\/[^/]+/);
  });

  cy
    .get('.storage-in-transit-panel .add-request')
    .contains('Request SIT')
    .click();

  cy
    .get('.usa-button-secondary')
    .contains('Cancel')
    .click()
    .get('.storage-in-transit-panel .add-request')
    .should($div => {
      const text = $div.text();
      expect(text).to.not.include('Sit Location');
    });
}

function tspUserEditsSitRequest() {
  // Open accepted shipments queue
  cy.patientVisit('/queues/accepted');

  // Open new shipments queue
  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/accepted/);
  });

  // Find shipment and open it
  cy.selectQueueItemMoveLocator('SITPAN');

  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/shipments\/[^/]+/);
  });

  cy
    .get('.storage-in-transit-panel .add-request')
    .contains('Request SIT')
    .click();

  fillAndSaveStorageInTransit();

  // click edit
  cy
    .get('.sit-edit a')
    .first()
    .click();

  editAndSaveStorageInTransit();
}

function tspUserStartsAndCancelsSitPlaceInSit() {
  // Open accepted shipments queue
  cy.patientVisit('/queues/in_transit');

  // Open new shipments queue
  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/in_transit/);
  });

  // Find shipment and open it
  cy.selectQueueItemMoveLocator('SITAPR');

  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/shipments\/[^/]+/);
  });

  cy
    .get('[data-cy=storage-in-transit-panel] [data-cy=place-in-sit-link]')
    .contains('Place into SIT')
    .click();

  cy.get('input[name=actual_start_date').should('have.value', '3/26/2019');
  cy
    .get('.usa-button-secondary')
    .contains('Cancel')
    .click()
    .get('[data-cy=storage-in-transit-panel] [data-cy=place-in-sit-link]')
    .should($div => {
      const text = $div.text();
      expect(text).to.not.include('Actual start date');
    });
}
