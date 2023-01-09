import { officeBaseURL } from '../../support/constants';

describe('Office Home Page', function () {
  before(() => {
    cy.prepareOfficeApp();
  });

  it('successfully loads when not logged in', function () {
    cy.logout();
    cy.contains('office.move.mil');
    cy.contains('Sign in');
  });

  it('office user can use a single click to view move info', function () {
    cy.signIntoOffice();
    cy.get('tr[data-uuid]').first().click();
    cy.url().should('include', '/moves/');
  });
});

describe('Office authorization', () => {
  before(() => {
    cy.prepareOfficeApp();
  });

  beforeEach(() => {
    cy.clearAllCookies();
    cy.intercept('**/ghc/v1/queues/counseling?page=1&perPage=20&sort=submittedAt&order=asc&needsPPMCloseout=false').as(
      'getCounselingSortedOrders',
    );
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

  it('redirects Services Counselor to Services Counselor homepage', () => {
    cy.signInAsNewServicesCounselorUser();
    cy.wait('@getCounselingSortedOrders');

    cy.contains('Moves');
    cy.contains('Needs counseling');
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
      cy.contains('KKFA moves');

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
});
