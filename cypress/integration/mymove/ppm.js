/* global cy */

describe('completing the ppm flow', function() {
  beforeEach(() => {
    //profile@comple.te
    // cy.resetDb();
    cy.signInAsUser('13F3949D-0D53-4BE4-B1B1-AE4314793F34');
  });

  //tear down currently means doing this:
  //update moves set status='DRAFT';
  //delete from personally_procured_moves
  it('progresses thru forms', function() {
    cy.contains('Yuma AFB (from Yuma AFB)');
    cy.get('.whole_box > :nth-child(3) > span').contains('10,500 lbs');
    cy.contains('Continue Move Setup').click();

    cy.location().should(loc => {
      expect(loc.pathname).to.match(/^\/moves\/[^/]+\/ppm-start/);
    });
    cy
      .get('input[name="planned_move_date"]')
      .first()
      .type('9/2/2018{enter}')
      .blur();
    cy
      .get('input[name="pickup_postal_code"]')
      .clear()
      .type('80913');
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

    cy.get('.rangeslider__handle').click();

    cy.get('.incentive').contains('$');

    cy.get('input[type="radio"]').check('yes', { force: true });
    cy.get('input[name="requested_amount"]').type('1,333.91');
    cy.get('select[name="method_of_receipt"]').select('MilPay');
    cy.nextPage();

    cy.location().should(loc => {
      expect(loc.pathname).to.match(/^\/moves\/[^/]+\/review/);
    });

    // //todo: should probably have test suite for review and edit screens
    cy.contains('$1,333.91'); // Verify that the advance matches what was input

    cy.nextPage();

    cy.location().should(loc => {
      expect(loc.pathname).to.match(/^\/moves\/[^/]+\/agreement/);
    });

    cy.get('input[name="signature"]').type('Jane Doe');

    cy.nextPage();
    cy.contains('Success');
    cy.contains('Next Step: Awaiting approval');
    cy.contains('Advance Requested: $1,333.91');
  });
});
