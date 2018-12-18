describe(
  'As an office user, I want to keep track of HHG shipment costs before they are billed, ' +
    'so if there are no pre-approved, unbilled charges, I should see none',
  () => {
    beforeEach(() => {
      cy.signIntoOffice();
    });
    it('opens the shipment tab', () => {
      cy.visit('/queues/new/moves/6eee3663-1973-40c5-b49e-e70e9325b895/hhg');

      // The invoice table should be empty.
      cy
        .get('.invoice-panel .basic-panel-content')
        .children()
        .first()
        .should('have.text', 'No line items');
    });
  },
);
