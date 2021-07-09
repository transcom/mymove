describe('orders entry', function () {
  before(() => {
    cy.prepareCustomerApp();
  });

  beforeEach(() => {
    // cy.logout();
    cy.intercept('**/internal/orders/**').as('uploadAmendedOrders');
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

  it('will allow amended orders upload', function () {
    const userId = '6016e423-f8d5-44ca-98a8-af03c8445c94';
    cy.apiSignInAsUser(userId);
    cy.intercept('PATCH', '**/internal/orders/**').as('uploadAmendedOrder');
    cy.intercept('GET', '**/internal/service_members/**').as('currentOrders');

    cy.contains('Upload documents').click();

    cy.location().should((loc) => {
      expect(loc.pathname).to.eq('/orders/amend');
    });

    cy.upload_file('.filepond--root', 'top-secret.png');
    cy.get('[data-filepond-item-state="processing-complete"]').should('have.length', 1);

    cy.wait(['@currentOrders', '@uploadAmendedOrder']);

    cy.nextPage();

    cy.location().should((loc) => {
      expect(loc.pathname).to.eq('/');
    });

    cy.get('.usa-alert--success').within(() => {
      cy.contains('The transportation office will review your new documents');
    });
  });
});
