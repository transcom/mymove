/* global cy */

describe('completing the ppm flow', function() {
  it('progresses thru forms', function() {
    //profile@comple.te
    cy.signInAsUser('13f3949d-0d53-4be4-b1b1-ae4314793f34');
    cy.contains('Fort Gordon (from Yuma AFB)');
    cy.get('.whole_box > div > :nth-child(3) > span').contains('10,500 lbs');
    cy.contains('Continue Move Setup').click();

    cy.location().should(loc => {
      expect(loc.pathname).to.match(/^\/moves\/[^/]+\/ppm-start/);
    });
    cy.get('.wizard-header').should('not.exist');
    cy
      .get('input[name="original_move_date"]')
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

    cy.get('.wizard-header').should('not.exist');
    //todo verify entitlement
    cy.contains('moving truck').click();

    cy.nextPage();

    cy.location().should(loc => {
      expect(loc.pathname).to.match(/^\/moves\/[^/]+\/ppm-incentive/);
    });

    cy.get('.wizard-header').should('not.exist');
    cy.get('.rangeslider__handle').click();

    cy.get('.incentive').contains('$');

    cy.get('input[type="radio"]').check('yes', { force: true });
    cy.get('input[name="requested_amount"]').type('1,333.91');
    cy.get('select[name="method_of_receipt"]').select('MilPay');
    cy.nextPage();

    cy.location().should(loc => {
      expect(loc.pathname).to.match(/^\/moves\/[^/]+\/review/);
    });
    cy.get('.wizard-header').should('not.exist');

    // //todo: should probably have test suite for review and edit screens
    cy.contains('$1,333.91'); // Verify that the advance matches what was input
    cy.contains('Storage: Not requested'); // Verify SIT on the ppm review page since it's optional on HHG_PPM

    cy.nextPage();

    cy.location().should(loc => {
      expect(loc.pathname).to.match(/^\/moves\/[^/]+\/agreement/);
    });
    cy.get('.wizard-header').should('not.exist');

    cy.get('input[name="signature"]').type('Jane Doe');

    cy.nextPage();

    cy.location().should(loc => {
      expect(loc.pathname).to.match(/^\/$/);
    });

    cy.contains('Success');
    cy.contains('Next Step: Wait for approval');
    cy.contains('Advance Requested: $1,333.91');
  });

  it('allows a SM to request payment', function() {
    cy.logout();
    //profile@comple.te
    cy.signInAsUser('8e0d7e98-134e-4b28-bdd1-7d6b1ff34f9e');
    cy.contains('Fort Gordon (from Yuma AFB)');
    cy.contains('Request Payment').click();

    cy.location().should(loc => {
      expect(loc.pathname).to.match(/^\/moves\/[^/]+\/request-payment/);
    });
    cy.get('input[type="checkbox"]').should('not.be.checked');
    const stub = cy.stub();
    cy.on('window:alert', stub);

    cy
      .contains('Legal Agreement / Privacy Act')
      .click()
      .then(() => {
        expect(stub.getCall(0)).to.be.calledWithMatch('LEGAL AGREEMENT / PRIVACY ACT');
      });
    cy.get('input[type="checkbox"]').should('not.be.checked');
    cy.get('input[type="checkbox"]').check({ force: true });
  });
});

describe('editing ppm only move', () => {
  it('sees only details relevant to PPM only move', () => {
    cy.signInAsUser('e10d5964-c070-49cb-9bd1-eaf9f7348eb6');
    cy
      .get('.sidebar button')
      .contains('Edit Move')
      .click();

    cy.get('.ppm-container').should(ppmContainer => {
      expect(ppmContainer).to.have.length(1);
      expect(ppmContainer).to.not.have.class('hhg-shipment-summary');
    });
  });
});
