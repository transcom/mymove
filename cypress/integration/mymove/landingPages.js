/* global cy */

describe('testing landing pages', function() {
  // Submitted draft move/orders but no move type yet.
  it('tests pre-move type', function() {
    // sm_no_move_type@example.com
    draftMove('9ceb8321-6a82-4f6d-8bb3-a1d85922a202');
  });

  // PPM: SUBMITTED
  it('tests submitted PPM', function() {
    // ppm@incomple.te
    ppmSubmitted('e10d5964-c070-49cb-9bd1-eaf9f7348eb6');
  });

  // PPM: APPROVED
  it('tests approved PPM', function() {
    // ppm@approv.ed
    ppmApproved('1842091b-b9a0-4d4a-ba22-1e2f38f26317');
  });

  // PPM: PAYMENT_REQUESTED
  it('tests PPM that has requested payment', function() {
    // ppmpayment@request.ed
    ppmPaymentRequested('beccca28-6e15-40cc-8692-261cae0d4b14');
  });

  // PPM: COMPLETED
  // Not seeing a path to a COMPLETED PPM move at this time.

  // PPM: CANCELED
  it('tests canceled PPM', function() {
    // ppm-canceled@example.com
    canceledMove('20102768-4d45-449c-a585-81bc386204b1');
  });

  // HHG: SUBMITTED
  it('tests submitted HHG', function() {
    // hhg@incomplete.serviceagent
    hhgMoveSummary('412e76e0-bb34-47d4-ba37-ff13e2dd40b9');
  });

  // HHG: AWARDED
  it('tests awarded HHG', function() {
    // hhg@reject.ing
    hhgMoveSummary('76bdcff3-ade4-41ff-bf09-0b2474cec751');
  });

  // HHG: ACCEPTED
  it('tests accepted HHG', function() {
    // hhg@accept.ed
    hhgMoveSummary('6a39dd2a-a23f-4967-a035-3bc9987c6848');
  });

  // HHG: APPROVED
  it('tests approved HHG', function() {
    // hhg@approv.ed
    hhgMoveSummary('68461d67-5385-4780-9cb6-417075343b0e');
  });

  // HHG: IN_TRANSIT
  it('tests in-transit HHG', function() {
    // hhg@in.transit
    hhgMoveSummary('1239dd2a-a23f-4967-a035-3bc9987c6848');
  });

  // HHG: DELIVERED
  it('tests delivered HHG', function() {
    // hhg@de.livered
    hhgDeliveredOrCompletedMoveSummary('3339dd2a-a23f-4967-a035-3bc9987c6848');
  });

  // HHG: COMPLETED
  it('tests completed HHG', function() {
    // hhg@com.pleted
    hhgDeliveredOrCompletedMoveSummary('4449dd2a-a23f-4967-a035-3bc9987c6848');
  });

  // HHG: CANCELED
  it('tests canceled HHG', function() {
    // hhg@cancel.ed
    canceledMove('05ea5bc3-fd77-4f42-bdc5-a984a81b3829');
  });
});

function draftMove(userId) {
  cy.signInAsUser(userId);
  cy.contains('Move to be scheduled');
  cy.contains('Next Step: Finish setting up your move');
  cy.logout();
}

function ppmSubmitted(userId) {
  cy.signInAsUser(userId);
  cy.contains('Move your own stuff (PPM)');
  cy.contains('Next Step: Wait for approval');
  cy.should('not.contain', 'Add PPM (DITY) Move');
  cy.logout();
}

function ppmApproved(userId) {
  cy.signInAsUser(userId);
  cy.contains('Move your own stuff (PPM)');
  cy.contains('Next Step: Get ready to move');
  cy.contains('Next Step: Request payment');
  cy.logout();
}

function ppmPaymentRequested(userId) {
  cy.signInAsUser(userId);
  cy.contains('Move your own stuff (PPM)');
  cy.contains('Your payment is in review');
  cy.contains('You will receive a notification from your destination PPPO office when it has been reviewed.');
  cy.logout();
}

function hhgMoveSummary(userId) {
  cy.signInAsUser(userId);
  cy.contains('Government Movers and Packers (HHG)');
  cy.contains('Next Step: Prepare for move');
  cy.contains('Add PPM (DITY) Move');
  cy.logout();
}

function hhgDeliveredOrCompletedMoveSummary(userId) {
  cy.signInAsUser(userId);
  cy.contains('Government Movers and Packers (HHG)');
  cy.contains('Next Step: Survey');
  cy.logout();
}

function canceledMove(userId) {
  cy.signInAsUser(userId);
  cy.contains('New move');
  cy.contains('Start here');
  cy.logout();
}
