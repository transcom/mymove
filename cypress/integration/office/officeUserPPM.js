/* global cy */
import { userCancelsStorageDetails, userSavesStorageDetails } from '../../support/storagePanel';

describe('office user finds the move', function() {
  beforeEach(() => {
    cy.signIntoOffice();
  });

  it('office user views moves in queue new moves', function() {
    officeUserViewsMoves();
  });
  it('office user verifies the orders tab', function() {
    officeUserVerifiesOrders('VGHEIS');
  });
  it('office user verifies the accounting tab', function() {
    officeUserVerifiesAccounting();
  });
  it('office user approves move, verifies and approves PPM', function() {
    officeUserApprovesMoveAndPPM('VGHEIS');
    officeUserVerifiesPPM();
  });
  it('office user approves move, approves PPM with no advance', function() {
    const moveLocator = 'NOADVC';
    officeUserVerifiesOrders(moveLocator);
    cy.patientVisit('/');
    officeUserApprovesMoveAndPPM(moveLocator);

    // Make sure the page hasn't exploded
    cy
      .get('button')
      .contains('Approve PPM')
      .should('have.class', 'btn__approve--green');
  });
  it('office user views actual move date', function() {
    officeUserGoesToPPMPanel('PAYMNT');
    cy.get('.actual_move_date').should($div => {
      const text = $div.text();
      expect(text).to.include('Departure date');
      expect(text).to.include('11-Nov-18');
    });
  });
  it('edits missing text when the actual move date is not set', function() {
    officeUserGoesToPPMPanel('FDXTIU');
    cy.get('.actual_move_date').should($div => {
      const text = $div.text();
      expect(text).to.include('Departure date');
      expect(text).to.include('missing');
    });

    const actualDate = '11/20/2018';

    officeUserEditsDatesAndLocationsPanel(actualDate);

    cy.get('.actual_move_date').should($div => {
      const text = $div.text();
      expect(text).to.include('Departure date');
      expect(text).to.include('20-Nov-18');
    });
  });

  it('office user edits ppm net weight', function() {
    officeUserEditsNetWeight();
  });

  it('edits pickup and destination zip codes in estimates panel and these values are reflected in the storage and incentive calculators', function() {
    officeUserGoesToPPMPanel('FDXTIU');
    officeUserEditsEstimatesPanel(60606, 72018, 6000);
  });

  it('office user completes storage panel', function() {
    officeUserGoesToStoragePanel('STORAG');
    userSavesStorageDetails();
  });

  it('office user cancels storage panel', function() {
    officeUserGoesToStoragePanel('NOSTRG');
    userCancelsStorageDetails();
  });
});

function officeUserEditsNetWeight() {
  cy.patientVisit('/queues/ppm');

  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/ppm/);
  });

  cy.selectQueueItemMoveLocator('VGHEIS');

  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/new\/moves\/[^/]+\/basics/);
  });

  cy
    .get('span')
    .contains('PPM')
    .click();

  cy.get('.net_weight').contains('missing');

  cy
    .get('.editable-panel-header')
    .contains('Weights')
    .siblings()
    .click();

  cy.get('input[name="net_weight"]').type('6000');

  cy
    .get('button')
    .contains('Save')
    .click();

  cy.get('.net_weight').contains('6,000');
}

function officeUserViewsMoves() {
  // Open new moves queue
  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/new/);
  });

  // Find move and open it
  cy.selectQueueItemMoveLocator('VGHEIS');

  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/new\/moves\/[^/]+\/basics/);
  });
}

function officeUserVerifiesOrders(moveLocator) {
  // Open new moves queue
  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/new/);
  });

  // Find move and open it
  cy.selectQueueItemMoveLocator(moveLocator);

  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/new\/moves\/[^/]+\/basics/);
  });

  // Click on Orders document link, check that link matches
  cy
    .get('.panel-field')
    .contains('Orders (')
    // .find('a')
    .should('have.attr', 'href')
    .and('match', /^\/moves\/[^/]+\/orders/);

  // Click on edit orders
  cy
    .get('.editable-panel-header')
    .contains('Orders')
    .siblings()
    .click();

  // Enter details in form and save orders
  cy.get('input[name="orders.orders_number"]').type('666666');
  cy.get('select[name="orders.orders_type_detail"]').select('DELAYED_APPROVAL');

  cy.get('input[name="orders.orders_issuing_agency"]').type('ISSUING AGENCY');
  cy.get('input[name="orders.paragraph_number"]').type('FP-TP');

  cy
    .get('button')
    .contains('Save')
    .should('be.enabled');

  cy
    .get('button')
    .contains('Save')
    .click();

  // Verify data has been saved in the UI
  cy.get('span').contains('666666');

  cy.get('span').contains('ISSUING AGENCY');
  cy.get('span').contains('FP-TP');

  // Refresh browser and make sure changes persist
  cy.patientReload();
  cy
    .get('button')
    .contains('Approve PPM')
    .should('be.disabled');

  cy.get('span').contains('666666');
  cy.get('span').contains('Delayed Approval 20 Weeks or More');

  cy.get('span').contains('ISSUING AGENCY');
  cy.get('span').contains('FP-TP');
}

function officeUserVerifiesAccounting() {
  // Open new moves queue
  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/new/);
  });

  // Find move and open it
  cy.selectQueueItemMoveLocator('VGHEIS');

  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/new\/moves\/[^/]+\/basics/);
  });

  // Enter details in form and save
  cy
    .get('.editable-panel-header')
    .contains('Accounting')
    .siblings()
    .click();

  cy.get('input[name="tac"]').type('6789');
  cy.get('select[name="department_indicator"]').select('AIR_FORCE');
  cy.get('input[name="sac"]').type('N002214CSW32Y9');

  cy
    .get('button')
    .contains('Save')
    .should('be.enabled');

  cy
    .get('button')
    .contains('Save')
    .click();

  // Verify data has been saved in the UI
  // added '.tac' as an additional selector since
  // '6789' is part of the DoD ID that was a false positive
  cy.get('.tac > span').contains('6789');
  cy.get('span').contains('N002214CSW32Y9');

  // Refresh browser and make sure changes persist
  cy.patientReload();
  cy
    .get('button')
    .contains('Approve PPM')
    .should('be.disabled');

  cy.get('span').contains('6789');
  cy.get('span').contains('57 - United States Air Force');
  cy.get('span').contains('N002214CSW32Y9');
}

function officeUserApprovesMoveAndPPM(moveLocator) {
  // Open new moves queue
  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/new/);
  });

  // Find move and open it
  cy.selectQueueItemMoveLocator(moveLocator);

  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/new\/moves\/[^/]+\/basics/);
  });

  // Approve the move
  cy
    .get('button')
    .contains('Approve PPM')
    .should('be.disabled');
  cy
    .get('button')
    .contains('Approve Basics')
    .click();

  cy
    .get('button')
    .contains('Approve PPM')
    .should('be.enabled');

  // Open new moves queue
  cy.patientVisit('/');
  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/new/);
  });

  // Open PPMs Queue
  cy
    .get('span')
    .contains('PPMs')
    .click();

  // Find move and open it
  cy.selectQueueItemMoveLocator(moveLocator);

  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/new\/moves\/[^/]+\/basics/);
  });

  // Click on PPM tab
  cy
    .get('span')
    .contains('PPM')
    .click();
  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/new\/moves\/[^/]+\/ppm/);
  });

  // Approve PPM
  cy
    .get('button')
    .contains('Approve PPM')
    .click();
}

function officeUserVerifiesPPM() {
  // Verify that the Estimates section contains expected data
  cy.get('span').contains('8,000');

  // Approve advance
  cy.get('.payment-table').within(() => {
    // Verify the status icon
    cy
      .get('td svg:first')
      .should('have.attr', 'title')
      .and('eq', 'Awaiting Review');
    // Verify the approve checkmark
    cy
      .get('td svg:last')
      .should('have.attr', 'title')
      .and('eq', 'Approve');

    // Approve advance and verify icon change
    cy.get('td svg:last').click();
    cy
      .get('td svg:first')
      .should('have.attr', 'title')
      .and('eq', 'Approved');
  });
}

function officeUserGoesToPPMPanel(locator) {
  // Open ppm queue
  cy.patientVisit('/queues/ppm');

  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/ppm/);
  });

  // Find shipment and open it
  cy.selectQueueItemMoveLocator(locator);

  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/new\/moves\/[^/]+\/basics/);
  });

  cy
    .get('.title')
    .contains('PPM')
    .click();

  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/new\/moves\/[^/]+\/ppm/);
  });
}

function officeUserEditsDatesAndLocationsPanel(date) {
  cy
    .get('.editable-panel-header')
    .contains('Dates & Locations')
    .siblings()
    .click();

  cy
    .get('input[name="actual_move_date"]')
    .first()
    .clear()
    .type(date)
    .blur();

  cy
    .get('button')
    .contains('Save')
    .should('be.enabled');

  cy
    .get('button')
    .contains('Save')
    .click();
}

function officeUserEditsEstimatesPanel(destinationPostalCode, pickupPostalCode, weightEstimate) {
  cy
    .get('.editable-panel-header')
    .contains('Estimates')
    .siblings()
    .click();

  cy.get('.estimates').within(() => {
    cy
      .get('input[name="PPMEstimate.destination_postal_code"]')
      .clear()
      .type(destinationPostalCode);

    cy
      .get('input[name="PPMEstimate.pickup_postal_code"]')
      .clear()
      .type(pickupPostalCode);

    cy
      .get('input[name="PPMEstimate.weight_estimate"]')
      .clear()
      .type(weightEstimate);
  });

  cy
    .get('button')
    .contains('Save')
    .should('be.enabled');

  cy
    .get('button')
    .contains('Save')
    .click();
}

function officeUserGoesToStoragePanel(locator) {
  // Open new moves queue
  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/new/);
  });

  // Find shipment and open it
  cy.selectQueueItemMoveLocator(locator);

  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/new\/moves\/[^/]+\/basics/);
  });

  cy
    .get('.title')
    .contains('PPM')
    .click();

  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/new\/moves\/[^/]+\/ppm/);
  });

  cy
    .get('.editable-panel-header')
    .contains('Storage')
    .siblings()
    .click();
}
