/* global cy */

// Sign in user
//  choose hhg move type
// Fill in the form (including optional addresses)
// click next (which submits data)
// verify next page is Review page
// Update hhg. . .
describe('completing the hhg flow', function() {
  beforeEach(() => {
    // sm_hhg@example.com
    cy.signInAsUser('4b389406-9258-4695-a091-0bf97b5a132f');
  });

  it('selects hhg and progresses thru form', function() {
    cy.visit('localhost:4000/moves/8718c8ac-e0c6-423b-bdc6-af971ee05b9a');

    cy.contains('Household Goods Move').click();
    cy.contains('Next').click({ force: true });
    cy.location().should(loc => {
      expect(loc.pathname).to.match(/^\/moves\/[^/]+\/hhg-start/);
    });
  });
});
