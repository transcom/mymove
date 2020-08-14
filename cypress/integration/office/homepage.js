/* global cy */
import { officeAppName, officeBaseURL } from '../../support/constants';

describe('Office Home Page', function () {
  beforeEach(() => {
    cy.setupBaseUrl(officeAppName);
  });
  it('successfully loads when not logged in', function () {
    cy.logout();
    officeUserIsOnSignInPage();
  });
  it('open accepted shipments queue and see moves', function () {
    cy.signInAsNewOfficeUser();
    officeAllMoves();
  });
  it('office user can use a single click to view move info', function () {
    cy.waitForReactTableLoad();

    cy.get('[data-testid=queueTableRow]:first').click();
    cy.url().should('include', '/moves/');
  });
});

describe('Office authorization', () => {
  describe('for a TOO user', () => {
    it('redirects TOO to TOO homepage', () => {
      cy.signInAsNewTOOUser();
    });
  });

  describe('for a TIO user', () => {
    it('redirects TIO to TIO homepage', () => {
      cy.signInAsNewTIOUser();
    });
  });

  describe('for a PPM user', () => {
    it('redirects PPM office user to old office queue', () => {
      cy.signInAsNewOfficeUser();
    });
  });

  describe('multiple role selection', () => {
    beforeEach(() => {
      cy.removeFetch();
      cy.server();
      cy.route('GET', '/ghc/v1/swagger.yaml').as('getGHCClient');
      cy.route('GET', '/ghc/v1/move-orders').as('getMoveOrders');
      cy.route('GET', '/ghc/v1/payment-requests').as('getPaymentRequests');
    });

    it('can switch between TOO & TIO roles', () => {
      cy.signInAsMultiRoleUser();
      cy.wait(['@getGHCClient', '@getMoveOrders']);
      cy.contains('All Customer Moves'); // TOO home

      cy.contains('Change user role').click();
      cy.url().should('contain', '/select-application');

      cy.contains('Select transportation_invoicing_officer').click();
      cy.url().should('eq', officeBaseURL + '/');
      cy.wait('@getPaymentRequests');
      cy.contains('Payment Requests');

      cy.contains('Change user role').click();
      cy.url().should('contain', '/select-application');
      cy.contains('Select transportation_ordering_officer').click();
      cy.wait('@getMoveOrders');
      cy.url().should('eq', officeBaseURL + '/');
      cy.contains('All Customer Moves');
    });
  });
});

describe('Queue staleness indicator', () => {
  it('displays the correct time ago text', () => {
    cy.clock();
    cy.setupBaseUrl(officeAppName);
    cy.signInAsNewOfficeUser();
    cy.patientVisit('/queues/all');

    cy.get('[data-testid=staleness-indicator]').should('have.text', 'Last updated a few seconds ago');

    cy.tick(120000);

    cy.get('[data-testid=staleness-indicator]').should('have.text', 'Last updated 2 mins ago');
  });
});

function officeUserIsOnSignInPage() {
  cy.contains('office.move.mil');
  cy.contains('Sign In');
}

function officeAllMoves() {
  cy.patientVisit('/queues/all');
  cy.location().should((loc) => {
    expect(loc.pathname).to.match(/^\/queues\/all/);
  });

  cy.get('[data-testid=locator]').contains('NOSHOW').should('not.exist');
}
