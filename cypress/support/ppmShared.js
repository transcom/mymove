export function setMobileViewport() {
  cy.viewport(479, 875);
}

export function customerStartsAddingAPPMShipment() {
  cy.get('button[data-testid="shipment-selection-btn"]').click();
  cy.nextPage();

  cy.get('input[value="PPM"]').check({ force: true });
  cy.nextPage();
}

export function signInAndNavigateFromHomePageToReviewPage(userId) {
  cy.apiSignInAsUser(userId);

  cy.wait('@getShipment');

  cy.get('h3').should('contain', 'Time to submit your move');

  cy.get('button').contains('Review and submit').click();
}

export function signInAndNavigateFromHomePageToExistingPPMDateAndLocationPage(userId) {
  cy.apiSignInAsUser(userId);

  cy.wait('@getShipment');

  cy.get('h3').should('contain', 'Time to submit your move');

  cy.get('button').contains('PPM').click();

  cy.wait('@getShipment');

  cy.location().should((loc) => {
    expect(loc.pathname).to.match(/^\/moves\/[^/]+\/shipments\/[^/]+\/edit/);
  });

  cy.get('h1').should('contain', 'PPM date & location');
}

// used for creating a new shipment
export function submitsDateAndLocation() {
  cy.get('input[name="pickupPostalCode"]').clear().type('90210').blur();

  cy.get('input[name="destinationPostalCode"]').clear().type('76127');
  cy.get('input[name="expectedDepartureDate"]').clear().type('01 Feb 2022').blur();

  navigateFromDateAndLocationPageToEstimatedWeightsPage('@createShipment');
}

export function navigateFromDateAndLocationPageToEstimatedWeightsPage(actionToWaitOn) {
  cy.get('button').contains('Save & Continue').should('be.enabled').click();

  if (actionToWaitOn) {
    cy.wait(actionToWaitOn);
  }

  cy.location().should((loc) => {
    expect(loc.pathname).to.match(/^\/moves\/[^/]+\/shipments\/[^/]+\/estimated-weight/);
  });

  cy.get('h1').should('contain', 'Estimated weight');
}

export function submitsEstimatedWeightsAndProGear() {
  cy.get('input[name="estimatedWeight"]').clear().type(500).blur();
  cy.get('input[name="hasProGear"][value="true"]').check({ force: true });
  cy.get('input[name="proGearWeight"]').clear().type(500).blur();
  cy.get('button').contains('Save & Continue').should('be.enabled');

  navigateFromEstimatedWeightsPageToEstimatedIncentivePage();
}

export function submitsEstimatedWeights() {
  cy.get('input[name="estimatedWeight"]').clear().type(500).blur();
  cy.get('button').contains('Save & Continue').should('be.enabled');

  navigateFromEstimatedWeightsPageToEstimatedIncentivePage();
}

export function navigateFromEstimatedWeightsPageToEstimatedIncentivePage() {
  cy.get('button').contains('Save & Continue').should('be.enabled').click();

  cy.wait('@patchShipment');

  cy.location().should((loc) => {
    expect(loc.pathname).to.match(/^\/moves\/[^/]+\/shipments\/[^/]+\/estimated-incentive/);
  });

  cy.get('h1').should('contain', 'Estimated incentive');
}

export function generalVerifyEstimatedIncentivePage(isMobile = false) {
  cy.get('h1').should('contain', 'Estimated incentive');

  // checks the format of the incentive amount statement is `$<some comma-separated number without decimals> is`
  cy.get('.container h2').contains(/\$\d{1,3}(?:,\d{3})*? is/);

  if (!isMobile) {
    cy.get('button').contains('Next').should('not.be.disabled');
  } else {
    cy.get('button').contains('Next').should('not.be.disabled').should('have.css', 'order', '1');
  }

  navigateFromEstimatedIncentivePageToAdvancesPage();
}

export function navigateFromEstimatedIncentivePageToAdvancesPage() {
  cy.get('button').contains('Next').should('be.enabled').click();

  cy.location().should((loc) => {
    expect(loc.pathname).to.match(/^\/moves\/[^/]+\/shipments\/[^/]+\/advances/);
  });

  cy.get('h1').should('contain', 'Advances');
}

export function submitsAdvancePage(addAdvance = false, isMobile = false) {
  if (addAdvance) {
    cy.get('input[name="advanceRequested"][value="true"]').check({ force: true });

    cy.get('input[name="amountRequested"]').clear().type(4000).blur();

    cy.get('input[name="agreeToTerms"]').check({ force: true });
  } else {
    cy.get('input[name="advanceRequested"][value="false"]').check({ force: true });
  }

  navigateFromAdvancesPageToReviewPage(isMobile);
}

export function navigateFromAdvancesPageToReviewPage(isMobile = false) {
  // when navigating through an existing PPM that requested an advance, we must agree to the terms again to proceed
  // using cypress get or contains would result in an assertion failure for the case where advance requested is No
  cy.get('body').then(($body) => {
    if ($body.find('input[name="advanceRequested"][value="true"]:checked').length) {
      cy.get('input[name="agreeToTerms"]').check({ force: true });
    }
  });

  cy.get('button').contains('Save & Continue').as('saveButton');

  if (isMobile) {
    cy.get('@saveButton').should('have.css', 'order', '1');
  }

  cy.get('@saveButton').should('be.enabled').click();

  cy.wait('@patchShipment');

  cy.location().should((loc) => {
    expect(loc.pathname).to.match(/^\/moves\/[^/]+\/review/);
  });

  cy.get('h1').should('contain', 'Review your details');
}
