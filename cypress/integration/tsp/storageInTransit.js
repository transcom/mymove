import { fillAndSaveStorageInTransit } from '../../support/testCreateStorageInTransit';

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
    .get('.storage-in-transit-request-form')
    .should($div => {
      const text = $div.text();
      expect(text).to.include('SIT location');
      expect(text).to.include('Warehouse ID number');
      expect(text).to.include('Warehouse name');
      expect(text).to.include('Address Line 1');
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
