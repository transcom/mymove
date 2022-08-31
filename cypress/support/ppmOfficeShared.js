export function navigateToShipmentDetails(moveLocator) {
  /**
   * SC Moves queue
   */
  cy.wait(['@getSortedMoves']);
  cy.get('input[name="locator"]').as('moveCodeFilterInput');

  // type in move code/locator to filter
  cy.get('@moveCodeFilterInput').type(moveLocator).blur();
  cy.wait(['@getFilterSortedMoves']);

  // check if results appear, should be 1
  // and see if result have move code
  cy.get('tbody > tr').as('results');
  cy.get('@results').should('have.length', 1);
  cy.get('@results').first().contains(moveLocator);

  // click result to navigate to move details page
  cy.get('@results').first().click();
  cy.url().should('include', `/counseling/moves/${moveLocator}/details`);
  cy.wait(['@getMoves', '@getOrders', '@getMTOShipments', '@getMTOServiceItems']);
}

export const fillOutWeight = (options = {}) => {
  const { estimatedWeight = '4000', hasProGear = false, proGearWeight = '1000', spouseProGearWeight = '500' } = options;

  cy.get('input[name="estimatedWeight"]').clear().type(estimatedWeight).blur();
  if (hasProGear) {
    cy.get('input[name="hasProGear"][value="yes"]').check({ force: true });
    cy.get('input[name="proGearWeight"]').type(proGearWeight).blur();
    cy.get('input[name="spouseProGearWeight"]').type(spouseProGearWeight).blur();
  } else {
    cy.get('input[name="hasProGear"][value="no"]').check({ force: true });
  }
};

export const fillOutSitExpected = (options = {}) => {
  const {
    hasSit = true,
    sitEstimatedWeight = '1000',
    sitEstimatedEntryDate = '01 Mar 2020',
    sitEstimatedDepartureDate = '31 Mar 2020',
    sitLocation = 'DESTINATION', // takes in values of either 'DESTINATION' or 'ORIGIN'
  } = options;

  if (hasSit) {
    cy.get('input[name="sitExpected"][value="yes"]').check({ force: true });
    cy.get(`input[name="sitLocation"][value="${sitLocation}"]`).check({ force: true });
    cy.get('input[name="sitEstimatedWeight"]').clear().type(sitEstimatedWeight).blur();
    cy.get('input[name="sitEstimatedEntryDate"]').clear().type(sitEstimatedEntryDate).blur();
    cy.get('input[name="sitEstimatedDepartureDate"]').clear().type(sitEstimatedDepartureDate).blur();
  } else {
    cy.get('input[name="sitExpected"][value="no"]').check({ force: true });
  }
};

export const fillOutOriginInfo = (options = {}) => {
  const {
    expectedDepartureDate = '09 Jun 2022',
    pickupPostalCode = '90210',
    secondPickupPostalCode = '07003',
  } = options;

  cy.get('input[name="expectedDepartureDate"]').clear().type(expectedDepartureDate).blur();
  cy.get('input[name="pickupPostalCode"]').clear().type(pickupPostalCode).blur();
  if (secondPickupPostalCode) {
    cy.get('input[name="secondPickupPostalCode"]').clear().type(secondPickupPostalCode).blur();
  }
};

export const fillOutDestinationInfo = (options = {}) => {
  const { destinationPostalCode = '76127', secondDestinationPostalCode = '08540' } = options;

  cy.get('input[name="destinationPostalCode"]').clear().type(destinationPostalCode).blur();
  if (secondDestinationPostalCode) {
    cy.get('input[name="secondDestinationPostalCode"]').clear().type(secondDestinationPostalCode).blur();
  }
};

export const fillOutIncentiveAndAdvance = (options = {}) => {
  const { hasAdvance = true, advance = '6000' } = options;

  if (hasAdvance) {
    cy.get('input[name="advanceRequested"][value="Yes"]').check({ force: true });
    cy.get('input[name="advance"]').clear().type(advance).blur();
  } else {
    cy.get('input[name="advanceRequested"][value="No"]').check({ force: true });
  }
};
