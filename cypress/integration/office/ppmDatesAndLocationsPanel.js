/* global cy */
describe('office user interacts with ppm dates and locations panel', function() {
  beforeEach(() => {
    cy.signIntoOffice();
  });
  it('office user sees actual move date', function() {
    officeUserGoesToDatesAndLocationsPanel('PAYMNT');
    cy.get('.actual_move_date').should($div => {
      const text = $div.text();
      expect(text).to.include('Departure date');
      expect(text).to.include('11-Nov-18');
    });
  });

  it('edits missing text when the actual move date is not set', function() {
    officeUserGoesToDatesAndLocationsPanel('FDXTIU');
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
});

function officeUserGoesToDatesAndLocationsPanel(locator) {
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
