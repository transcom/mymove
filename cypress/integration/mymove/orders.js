/* global cy*/
describe('orders entry', function() {
  it('will accept orders information', function() {
    cy.signInAsUser('feac0e92-66ec-4cab-ad29-538129bf918e');
    cy.contains('New move (from Yuma AFB)');
    cy.contains('No details');
    cy.contains('No documents');
    cy.contains('Continue Move Setup').click();

    cy.location().should(loc => {
      expect(loc.pathname).to.eq('/orders/');
    });

    cy.get('select[name="orders_type"]').select('Permanent Change Of Station');

    cy
      .get('input[name="issue_date"]')
      .first()
      .click();

    cy
      .get('input[name="issue_date"]')
      .first()
      .type('6/2/2018{enter}')
      .blur();

    cy
      .get('input[name="report_by_date"]')
      .last()
      .type('8/9/2018{enter}')
      .blur();

    cy.selectDutyStation('NAS Fort Worth', 'new_duty_station');

    cy.nextPage();

    cy.location().should(loc => {
      expect(loc.pathname).to.eq('/orders/upload');
    });

    cy.visit('/');
    cy.contains('NAS Fort Worth (from Yuma AFB)');
    cy.get('.whole_box > div > :nth-child(3) > span').contains('7,000 lbs');
    cy.contains('Continue Move Setup').click();
    cy.location().should(loc => {
      expect(loc.pathname).to.eq('/orders/upload');
    });
  });
});
