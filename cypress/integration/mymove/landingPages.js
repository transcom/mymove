/* global cy */

import { milmoveAppName } from '../../support/constants';

describe('testing landing pages', function () {
  // Submitted draft move/orders but no move type yet.
  it('tests pre-move type', function () {
    // sm_no_move_type@example.com
    draftMove('9ceb8321-6a82-4f6d-8bb3-a1d85922a202');
  });

  // PPM: SUBMITTED
  it('tests submitted PPM', function () {
    // ppm@incomple.te
    ppmSubmitted('e10d5964-c070-49cb-9bd1-eaf9f7348eb6');
  });

  // PPM: APPROVED
  it('tests approved PPM', function () {
    // ppm@approv.ed
    ppmApproved('70665111-7bbb-4876-a53d-18bb125c943e');
  });

  // PPM: PAYMENT_REQUESTED
  it('tests PPM that has requested payment', function () {
    // ppmpayment@request.ed
    ppmPaymentRequested('beccca28-6e15-40cc-8692-261cae0d4b14');
  });

  // PPM: COMPLETED
  // Not seeing a path to a COMPLETED PPM move at this time.

  // PPM: CANCELED
  it('tests canceled PPM', function () {
    // ppm-canceled@example.com
    canceledMove('20102768-4d45-449c-a585-81bc386204b1');
  });
});

function draftMove(userId) {
  cy.signInAsUserPostRequest(milmoveAppName, userId);
  cy.contains('Next Step: Finish setting up your move');
  cy.logout();
}

function ppmSubmitted(userId) {
  cy.signInAsUserPostRequest(milmoveAppName, userId);
  cy.contains('Handle your own move (PPM)');
  cy.contains('Next Step: Wait for approval');
  cy.should('not.contain', 'Add PPM (DITY) Move');
  cy.logout();
}

function ppmApproved(userId) {
  cy.signInAsUserPostRequest(milmoveAppName, userId);
  cy.contains('Handle your own move (PPM)');
  cy.contains('Next Step: Request payment');
  cy.logout();
}

function ppmPaymentRequested(userId) {
  cy.signInAsUserPostRequest(milmoveAppName, userId);
  cy.setFeatureFlag('ppmPaymentRequest=false', '/');
  cy.contains('Handle your own move (PPM)');
  cy.contains('Edit Payment Request');
  cy.contains('Estimated');
  cy.contains('Submitted');
  cy.logout();
}

function canceledMove(userId) {
  cy.signInAsUserPostRequest(milmoveAppName, userId);
  cy.contains('New move');
  cy.contains('Start here');
  cy.logout();
}
