export function customerChoosesAPPMMove() {
  cy.get('button[data-testid="shipment-selection-btn"]').click();
  cy.nextPage();

  cy.get('input[value="PPM"]').check({ force: true });
  cy.nextPage();
}

export function submitsDateAndLocation() {
  cy.get('input[name="pickupPostalCode"]').clear().type('90210').blur();

  cy.get('input[name="destinationPostalCode"]').clear().type('76127');
  cy.get('input[name="expectedDepartureDate"]').clear().type('01 Feb 2022').blur();

  cy.get('button').contains('Save & Continue').click();

  cy.location().should((loc) => {
    expect(loc.pathname).to.match(/^\/moves\/[^/]+\/shipments\/[^/]+\/estimated-weight/);
  });
}
