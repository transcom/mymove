import { userEntersDates, userEntersAndRemovesDates } from '../../support/datesPanel';

/* global cy */
describe('office user interacts with dates panel', function() {
  beforeEach(() => {
    cy.signIntoOffice();
  });
  it('office user completes dates panel', function() {
    officeUserGoesToDatesPanel('ODATES');
    userEntersDates();
  });
  it('office user completes dates panel and zeroes it out', function() {
    officeUserGoesToDatesPanel('ODATE0');
    userEntersAndRemovesDates();
  });
});

function officeUserGoesToDatesPanel(locator) {
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
    .contains('HHG')
    .click();

  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/new\/moves\/[^/]+\/hhg/);
  });

  cy
    .get('.editable-panel-header')
    .contains('Dates')
    .siblings()
    .click();
}
