describe('Office user looks at the invoice tab to view unbilled line items', () => {
  beforeEach(() => {
    cy.signIntoOffice();
  });

  it('there are no unbilled line items', checkNoUnbilledLineItems);

  it('there are unbilled line items', checkExistUnbilledLineItems);

  it('there are unbilled line items with an approve button', checkApproveButton);
});

function checkNoUnbilledLineItems() {
  // Open the shipments tab.
  cy.visit('/queues/new/moves/6eee3663-1973-40c5-b49e-e70e9325b895/hhg');

  // The invoice table should be empty.
  cy
    .get('.invoice-panel .basic-panel-content')
    .children()
    .first()
    .should('have.text', 'No line items');
}

function checkExistUnbilledLineItems() {
  // Open the shipments tab.
  cy.visit('/queues/new/moves/fb4105cf-f5a5-43be-845e-d59fdb34f31c/hhg');

  // The invoice table should display the unbilled line items.
  cy
    .get('.invoice-panel .basic-panel-content tbody')
    .children()
    // For each line item, I should see item code, description, etc.
    .each((row, index, lst) => {
      // Last line is the Totals line.
      if (index === lst.length - 1) {
        return;
      }

      cy
        .wrap(row)
        .children()
        .each(cell => {
          // Each cell should have a value present.
          cy.wrap(cell).should('not.have.text', '');
        });
    });
}

function checkApproveButton() {
  // Open the shipments tab.
  cy.visit('/queues/new/moves/fb4105cf-f5a5-43be-845e-d59fdb34f31c/hhg');

  // The invoice tab should have a button with the correct text.
  cy.get('.invoice-panel-header-cont button').should('have.text', 'Approve Payment');

  // Open shipments tab of move with no unbilled line items.
  cy.visit('/queues/new/moves/6eee3663-1973-40c5-b49e-e70e9325b895/hhg');

  // The invoice tab should not have a button.
  cy
    .get('.invoice-panel')
    .first()
    .within(() => {
      cy.get('button').should('not.exist');
    });
}
