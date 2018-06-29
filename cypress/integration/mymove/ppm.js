/* global cy */

//this is skipped until we can figure out how to upload orders or reset specified user to needing to select move
describe.skip('setting up service member profile', function() {
  beforeEach(() => {
    //20180629092519@example.com
    cy.signInAsUser('f775c2b4-c411-4860-873d-ce9a89786830');
  });

  //tear down currently means doing this:
  //update moves set status='DRAFT';
  //delete from personally_procured_moves
  it('progresses thru forms', function() {
    cy.contains('NAS Fort Worth from Ft Carson');
    cy.get('.whole_box > :nth-child(3) > span').contains('15,000 lbs');
    cy.contains('Continue Move Setup').click();

    cy.location().should(loc => {
      expect(loc.pathname).to.match(/^\/moves\/[^/]+\/ppm-start/);
    });
    cy
      .get('input[placeholder="Date"]')
      .first()
      .type('2018-9-2{enter}')
      .blur();

    cy.get('input[name="destination_postal_code"]').type('76127');

    cy.nextPage();

    cy.location().should(loc => {
      expect(loc.pathname).to.match(/^\/moves\/[^/]+\/ppm-size/);
    });

    //todo verify entitlement
    cy.contains('moving truck').click();

    cy.nextPage();

    cy.location().should(loc => {
      expect(loc.pathname).to.match(/^\/moves\/[^/]+\/ppm-incentive/);
    });

    cy.nextPage();

    cy.location().should(loc => {
      expect(loc.pathname).to.match(/^\/moves\/[^/]+\/review/);
    });

    //todo: should probably have test suite for review and edit screens

    cy.nextPage();

    cy.location().should(loc => {
      expect(loc.pathname).to.match(/^\/moves\/[^/]+\/agreement/);
    });

    cy.get('input[name="signature"]').type('Jane Doe');

    cy.nextPage();
    cy.contains('Success');
    cy.contains('Next Step: Awaiting approval');
    cy.contains('Weight (est.): 8250 lbs');
  });
});
