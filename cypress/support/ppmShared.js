export function customerChoosesAPPMMove() {
  cy.get('button[data-testid="shipment-selection-btn"]').click();
  cy.nextPage();

  cy.get('input[value="PPM"]').check({ force: true });
  cy.nextPage();
}

// used for creating a new shipment
export function submitsDateAndLocation() {
  cy.get('input[name="pickupPostalCode"]').clear().type('90210').blur();

  cy.get('input[name="destinationPostalCode"]').clear().type('76127');
  cy.get('input[name="expectedDepartureDate"]').clear().type('01 Feb 2022').blur();

  cy.get('button').contains('Save & Continue').click();
  cy.wait('@createShipment');

  cy.location().should((loc) => {
    expect(loc.pathname).to.match(/^\/moves\/[^/]+\/shipments\/[^/]+\/estimated-weight/);
  });
}

export function submitsEstimatedWeightsAndProgear() {
  cy.get('input[name="estimatedWeight"]').clear().type(500).blur();
  cy.get('input[name="hasProGear"][value="true"]').check({ force: true });
  cy.get('input[name="proGearWeight"]').clear().type(500).blur();
  cy.get('button').contains('Save & Continue').should('be.enabled');

  cy.get('button').contains('Save & Continue').click();
  cy.wait('@patchShipment');

  cy.location().should((loc) => {
    expect(loc.pathname).to.match(/^\/moves\/[^/]+\/shipments\/[^/]+\/estimated-incentive/);
  });
}

export function submitsEstimatedWeights() {
  cy.get('input[name="estimatedWeight"]').clear().type(500).blur();
  cy.get('button').contains('Save & Continue').should('be.enabled');

  cy.get('button').contains('Save & Continue').click();
  cy.wait('@patchShipment');

  cy.location().should((loc) => {
    expect(loc.pathname).to.match(/^\/moves\/[^/]+\/shipments\/[^/]+\/estimated-incentive/);
  });
}

export function generalVerifyEstimatedIncentivePage(isMobile = false) {
  cy.get('h1').should('contain', 'Estimated incentive');

  // checks the format of the incentive amount statment is `$<some comma-separated number without decimals> is`
  cy.get('.container h2').contains(/\$\d{1,3}(?:,\d{3})*? is/);

  if (!isMobile) {
    cy.get('button').contains('Next').should('not.be.disabled');
  } else {
    cy.get('button').contains('Next').should('not.be.disabled').should('have.css', 'order', '1');
  }

  cy.get('button').contains('Next').should('be.enabled').click();
  cy.location().should((loc) => {
    expect(loc.pathname).to.match(/^\/moves\/[^/]+\/shipments\/[^/]+\/advances/);
  });
}
