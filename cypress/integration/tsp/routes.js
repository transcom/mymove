import { tspAppName } from '../../support/constants';

/* global cy */
describe('TSP User Navigating the App', function() {
  it('unauthorized user tries to access an authorized route', function() {
    unauthorizedTspUserGoesToAuthorizedRoute();
  });
  it('tsp user navigates to an invalid shipment', function() {
    tspUserViewsInvalidShipment();
  });
  it('tsp user enters an invalid url route', function() {
    tspUserNavigatesToInvalidRoute();
  });
});

function unauthorizedTspUserGoesToAuthorizedRoute() {
  cy.setupBaseUrl(tspAppName);
  cy.logout();
  cy.patientVisit('/queues/new');
  cy.contains('Welcome to tsp.move.mil');
  cy.contains('Sign In');
}

function tspUserViewsInvalidShipment() {
  cy.signIntoTSP();
  // Open new shipments queue
  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/new/);
  });

  // visit an invalid url
  cy.patientVisit('/shipments/some-invalid-uuid');

  // redirected to the queues page due to invalid shipment
  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/new/);
  });
}

function tspUserNavigatesToInvalidRoute() {
  cy.signIntoTSP();

  // Open new shipments queue
  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/new/);
  });

  // visits an invalid url
  cy.patientVisit('/i-do-not-exist');

  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/new/);
  });
}
