/* global cy */
import { officeAppName } from '../../support/constants';

describe('TOO page', function() {
  beforeEach(() => {
    cy.setupBaseUrl(officeAppName);
  });
  // TODO: create a TOO user when a TOO user exists
  it('creates new devlocal user', function() {
    cy.signInAsNewOfficeUser();
    tooVisitsTOOHomepage();
  });
});

function tooVisitsTOOHomepage() {
  cy.patientVisit('/ghc/too');
  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/ghc\/too/);
  });
}
