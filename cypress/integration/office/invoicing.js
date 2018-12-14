describe(
  'As an office user, I want to keep track of HHG shipment costs before they are billed, ' +
    'so if there are no pre-approved, unbilled charges, I should see none',
  () => {
    beforeEach(() => {
      cy.signIntoOffice();
    });
    it('opens the shipment tab', () => {
      cy.visit('/queues/new/moves/fb4105cf-f5a5-43be-845e-d59fdb34f31c/hhg');

      // The invoice table should be empty.
      cy
        .get('.invoice-panel .basic-panel-content')
        .children()
        .first()
        .should('have.class', 'empty-content');
    });
  },
);
