import { officeBaseURL } from '../../support/constants';

describe('Office Home Page', function () {
  before(() => {
    cy.prepareOfficeApp();
  });

  it('successfully loads when not logged in', function () {
    cy.logout();
    cy.contains('office.move.mil');
    cy.contains('Sign In');
  });

  it('open accepted shipments queue and see moves', function () {
    cy.signInAsNewPPMOfficeUser();
    cy.patientVisit('/queues/all');
    cy.location().should((loc) => {
      expect(loc.pathname).to.match(/^\/queues\/all/);
    });
    cy.get('[data-testid=locator]').contains('NOSHOW').should('not.exist');
  });

  it('office user can use a single click to view move info', function () {
    cy.waitForReactTableLoad();
    cy.get('[data-testid=queueTableRow]:first').click();
    cy.url().should('include', '/moves/');
  });
});

describe('Office authorization', () => {
  before(() => {
    cy.prepareOfficeApp();
  });

  beforeEach(() => {
    cy.clearAllCookies();
  });

  it('redirects TOO to TOO homepage', () => {
    cy.signInAsNewTOOUser();
    cy.contains('All moves');
    cy.url().should('eq', officeBaseURL + '/');
  });

  it('redirects TIO to TIO homepage', () => {
    cy.signInAsNewTIOUser();
    cy.contains('Payment requests');
    cy.url().should('eq', officeBaseURL + '/');
  });

  it('redirects PPM office user to old office queue', () => {
    cy.signInAsNewPPMOfficeUser();
    cy.contains('New moves');
    cy.url().should('eq', officeBaseURL + '/');
  });

  describe('multiple role selection', () => {
    beforeEach(() => {
      cy.intercept('**/ghc/v1/swagger.yaml').as('getGHCClient');
      cy.intercept('**/ghc/v1/queues/moves?**').as('getOrders');
      cy.intercept('**/ghc/v1/queues/payment-requests?**').as('getPaymentRequests');
    });

    it('can switch between TOO & TIO roles', () => {
      cy.signInAsMultiRoleOfficeUser();
      cy.wait(['@getGHCClient', '@getOrders']);
      cy.contains('All moves'); // TOO home
      cy.contains('Spaceman, Leo');
      cy.contains('LKNQ moves');

      cy.contains('Change user role').click();
      cy.url().should('contain', '/select-application');

      cy.contains('Select transportation_invoicing_officer').click();
      cy.url().should('eq', officeBaseURL + '/');
      cy.wait('@getPaymentRequests');
      cy.contains('Payment requests');

      cy.contains('Change user role').click();
      cy.url().should('contain', '/select-application');
      cy.contains('Select transportation_ordering_officer').click();
      cy.wait('@getOrders');
      cy.url().should('eq', officeBaseURL + '/');
      cy.contains('All moves');
    });
  });
});

describe('Queue staleness indicator', () => {
  before(() => {
    cy.prepareOfficeApp();
    cy.clearAllCookies();
  });

  it('displays the correct time ago text', () => {
    cy.clock();
    cy.signInAsNewPPMOfficeUser();
    cy.patientVisit('/queues/all');

    cy.get('[data-testid=staleness-indicator]').contains('Last updated a few seconds ago');

    cy.tick(120000);

    cy.get('[data-testid=staleness-indicator]').contains('Last updated 2 mins ago');
  });
});
