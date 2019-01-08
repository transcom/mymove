/* global cy */
describe('TSP user looks at the invoice panel to view unbilled line items', () => {
  beforeEach(() => {
    cy.signIntoTSP();
  });

  it('there are no unbilled line items', checkNoUnbilledLineItems);

  it('there are unbilled line items', checkExistUnbilledLineItems);
});

function checkNoUnbilledLineItems() {
  // Open the shipment with no unbilled line items.
  cy.patientVisit('/shipments/0851706a-997f-46fb-84e4-2525a444ade0');

  // The invoice table should be empty.
  cy
    .get('.invoice-panel .basic-panel-content')
    .find('span.empty-content')
    .should('have.text', 'No line items');
}

function checkExistUnbilledLineItems() {
  // Open the shipment with unbilled line items.
  cy.patientVisit('shipments/67a3cbe7-4ae3-4f6a-9f9a-4f312e7458b9');

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
