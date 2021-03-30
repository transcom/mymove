describe('orders entry', function () {
  before(() => {
    cy.prepareCustomerApp();
  });

  it('will accept orders information', function () {
    // needs@orde.rs
    cy.apiSignInAsPpmUser('feac0e92-66ec-4cab-ad29-538129bf918e');
    cy.contains('New move (from Yuma AFB)');
    cy.contains('No details');
    cy.contains('No documents');
    cy.contains('Continue Move Setup').click();

    cy.location().should((loc) => {
      expect(loc.pathname).to.eq('/orders/info');
    });

    cy.get('select[name="orders_type"]').select('Separation');
    cy.get('select[name="orders_type"]').select('Retirement');
    cy.get('select[name="orders_type"]').select('Permanent Change Of Station (PCS)');

    cy.get('input[name="issue_date"]').first().click();

    cy.get('input[name="issue_date"]').first().type('6/2/2018{enter}').blur();

    cy.get('input[name="report_by_date"]').last().type('8/9/2018{enter}').blur();

    cy.get('label[for="hasDependentsNo"]').first().click();

    // Choosing same current and destination duty station should block you from progressing and give an error
    cy.selectDutyStation('Yuma AFB', 'new_duty_station');
    cy.get('.usa-error-message').contains(
      'You entered the same duty station for your origin and destination. Please change one of them.',
    );
    cy.get('button[data-testid="wizardNextButton"]').should('be.disabled');

    cy.selectDutyStation('NAS Fort Worth JRB', 'new_duty_station');

    cy.nextPage();

    cy.location().should((loc) => {
      expect(loc.pathname).to.eq('/orders/upload');
    });

    cy.setFeatureFlag('ppmPaymentRequest=false', '/ppm');
    cy.contains('NAS Fort Worth JRB (from Yuma AFB)');
    cy.get('[data-testid="move-header-weight-estimate"]').contains('5,000 lbs');
    cy.contains('Continue Move Setup').click();
    cy.location().should((loc) => {
      expect(loc.pathname).to.eq('/orders/upload');
    });
  });
});
