import { userEntersDates } from '../../support/datesPanel';

/* global cy */
describe('TSP User Completes Dates Panel', function() {
  beforeEach(() => {
    cy.signIntoTSP();
  });
  it('tsp user completes dates panel', function() {
    tspUserGoesToDatesPanel('DATESP');
    userEntersDates();
  });
});

function tspUserGoesToDatesPanel(locator) {
  // Open new shipments queue
  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/new/);
  });

  // Find shipment and open it
  cy.selectQueueItemMoveLocator(locator);

  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/shipments\/[^/]+/);
  });

  cy
    .get('.editable-panel-header')
    .contains('Dates')
    .siblings()
    .click();
}
