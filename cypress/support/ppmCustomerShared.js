import { signAgreement } from '../integration/mymove/utilities/customer';

export function setMobileViewport() {
  cy.viewport(479, 875);
}

export function customerStartsAddingAPPMShipment() {
  cy.get('button[data-testid="shipment-selection-btn"]').click();
  cy.nextPage();

  cy.get('input[value="PPM"]').check({ force: true });
  cy.nextPage();
}

export function signInAndNavigateFromHomePageToReviewPage(userId, isMoveSubmitted = false) {
  cy.apiSignInAsUser(userId);

  cy.wait('@getShipment');

  navigateFromHomePageToReviewPage(isMoveSubmitted);
}

export function signInAndNavigateToAboutPage(userId, selectAdvance) {
  cy.apiSignInAsUser(userId);

  cy.wait('@getShipment');
  cy.screenshot();
  cy.get('button[data-testid="button"]').contains('Upload PPM Documents').click();
  cy.location().should((loc) => {
    expect(loc.pathname).to.match(/^\/moves\/[^/]+\/shipments\/[^/]+\/about/);
  });

  fillOutAboutPage(selectAdvance);
}

export function navigateFromHomePageToReviewPage(isMoveSubmitted = false) {
  if (isMoveSubmitted) {
    cy.get('h3').contains('Next step: Your move gets approved');

    cy.get('button').contains('Review your request').click();
  } else {
    cy.get('h3').should('contain', 'Time to submit your move');

    cy.get('button').contains('Review and submit').click();
  }
}

export function fillOutAboutPage(selectAdvance) {
  cy.get('input[name="actualMoveDate"]').clear().type('01 Feb 2022').blur();
  cy.get('input[name="actualPickupPostalCode"]').clear().type('90210').blur();
  cy.get('input[name="actualDestinationPostalCode"]').clear().type('76127');
  if (selectAdvance) {
    cy.get('input[name="hasReceivedAdvance"][value="true"]').check({ force: true });
    cy.get('input[name="advanceAmountReceived"]').clear().type('5000');
  } else {
    cy.get('input[name="hasReceivedAdvance"][value="false"]').check({ force: true });
  }
  cy.get('button').contains('Save & Continue').should('be.enabled');
  navigateFromAboutPageToWeightTicketPage();
}

export function navigateFromAboutPageToWeightTicketPage() {
  cy.get('button').contains('Save & Continue').click();
  cy.location().should((loc) => {
    expect(loc.pathname).to.match(/^\/moves\/[^/]+\/shipments\/[^/]+\/weight-tickets/);
  });
}

export function signInAndNavigateToWeightTicketPage(userId, options) {
  signInAndNavigateToAboutPage(userId, true);
  fillOutAboutPage(true);
  navigateFromAboutPageToWeightTicketPage();
  fillOutWeightTicketPage(options);
}

export function fillOutWeightTicketPage(options) {
  const { useConstructedWeight, hasTrailer, ownTrailer } = options;
  cy.get('input[name="vehicleDescription"]').clear().type('Kia Forte').blur();

  if (useConstructedWeight) {
    cy.get('input[name="emptyWeight"]').clear().type('1000').blur();
    cy.get('input[name="missingEmptyWeightTicket"]').check({ force: true });
    cy.upload_file('.emptyWeightTicketsRef', 'constructedWeight.xls');
    cy.get('input[name="fullWeight"]').clear().type('3000');
    cy.get('input[name="missingFullWeightTicket"]').check({ force: true });
    cy.upload_file('.fullWeightTicketsRef', 'constructedWeight.xls');
  } else {
    cy.get('input[name="emptyWeight"]').clear().type('1000').blur();
    // TODO: add weight ticket image to /fixtures
    cy.upload_file('.emptyWeightTicketsRef', 'weighTicket.pdf');
    cy.get('input[name="fullWeight"]').clear().type('3000');
    // TODO: add file upload for full weight ticket
    cy.upload_file('.fullWeightTicketsRef', 'weighTicket.pdf');
  }

  cy.get('.tripWeightTotal').contains('Trip weight: 2,000 lbs');

  if (hasTrailer) {
    cy.get('input[name="hasOwnTrailer"][value="true"]').check({ force: true });
    if (ownTrailer) {
      cy.get('input[name="trailerMeetsCriteria"][value="true"]').check({ force: true });
      cy.upload_file('.trailerOwnershipDocsRef', 'trailerOwnership.pdf');
    } else {
      cy.get('input[name="trailerMeetsCriteria"][value="false"]').check({ force: true });
      cy.get('.doNotClaimTrailerWeight').contains('Do not claim the weight of this trailer for this trip.');
    }
  } else {
    cy.get('input[name="hasOwnTrailer"][value="false"]').check({ force: true });
  }

  cy.get('button').contains('Save & Continue').should('be.enabled');
  navigateFromWeightTicketPage();
}

export function navigateFromWeightTicketPage() {
  cy.get('button').contains('Save & Continue').click();
  cy.location().should((loc) => {
    // TODO: use correct path for subsequent page
    expect(loc.pathname).to.match(/^\/moves\/[^/]+\/shipments\/[^/]+\/weight-tickets/);
  });
}

export function signInAndNavigateFromHomePageToExistingPPMDateAndLocationPage(userId) {
  cy.apiSignInAsUser(userId);

  cy.wait('@getShipment');

  cy.get('h3').should('contain', 'Time to submit your move');

  cy.get('[data-testid="shipment-list-item-container"] button').contains('Edit').click();

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

  navigateFromDateAndLocationPageToEstimatedWeightsPage();
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
  cy.get('input[name="estimatedWeight"]').clear().type(4000).blur();
  cy.get('input[name="hasProGear"][value="true"]').check({ force: true });
  cy.get('input[name="proGearWeight"]').clear().type(500).blur();
  cy.get('input[name="spouseProGearWeight"]').clear().type(400).blur();
  cy.get('button').contains('Save & Continue').should('be.enabled');

  navigateFromEstimatedWeightsPageToEstimatedIncentivePage();
}

export function submitsEstimatedWeights() {
  cy.get('input[name="estimatedWeight"]').clear().type(4000).blur();
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
    cy.get('input[name="hasRequestedAdvance"][value="true"]').check({ force: true });

    cy.get('input[name="advanceAmountRequested"]').clear().type(4000).blur();

    cy.get('input[name="agreeToTerms"]').check({ force: true });
  } else {
    cy.get('input[name="hasRequestedAdvance"][value="false"]').check({ force: true });
  }

  navigateFromAdvancesPageToReviewPage(isMobile);
}

export function navigateFromAdvancesPageToReviewPage(isMobile = false) {
  // when navigating through an existing PPM that requested an advance, we must agree to the terms again to proceed
  // using cypress get or contains would result in an assertion failure for the case where advance requested is No
  cy.get('body').then(($body) => {
    if ($body.find('input[name="hasRequestedAdvance"][value="true"]:checked').length) {
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

  cy.get('.usa-alert__heading')
    .contains('Details saved')
    .next()
    .contains('Review your info and submit your move request now, or come back and finish later.');
}

export function navigateFromReviewPageToHomePage() {
  cy.get('button').contains('Return home').click();

  cy.location().should((loc) => {
    expect(loc.pathname).to.equal('/');
  });
}

export function navigateToAgreementAndSign() {
  cy.nextPage();
  signAgreement();
}

export function deleteShipment(selector, expectedLength) {
  cy.get(selector).contains('Delete').click();
  cy.get('[data-testid="modal"]').within(($modal) => {
    expect($modal).to.be.visible;
    cy.get('button').contains('Yes, Delete').click();
  });
  cy.wait(['@deleteShipment', '@getShipment']);
  if (expectedLength > 0) {
    cy.get(selector).should('have.length', expectedLength);
  } else {
    cy.get(selector).should('not.exist');
  }
  cy.get('[data-testid="alert"]').should('contain', 'The shipment was deleted.');
}
