/* global cy */

import { milmoveAppName } from '../../support/constants';

describe('completing the ppm flow', function() {
  it('progresses thru forms', function() {
    //profile@comple.te
    cy.signInAsUserPostRequest(milmoveAppName, '13f3949d-0d53-4be4-b1b1-ae4314793f34');
    cy.contains('Fort Gordon (from Yuma AFB)');
    cy.get('.whole_box > div > :nth-child(3) > span').contains('10,500 lbs');
    cy.contains('Continue Move Setup').click();

    cy.location().should(loc => {
      expect(loc.pathname).to.match(/^\/moves\/[^/]+\/ppm-start/);
    });
    cy.get('.wizard-header').should('not.exist');
    cy
      .get('input[name="original_move_date"]')
      .first()
      .type('9/2/2018{enter}')
      .blur();
    cy
      .get('input[name="pickup_postal_code"]')
      .clear()
      .type('80913');

    cy
      .get('input[name="destination_postal_code"]')
      .clear()
      .type('76127');

    cy.nextPage();

    cy.location().should(loc => {
      expect(loc.pathname).to.match(/^\/moves\/[^/]+\/ppm-size/);
    });

    cy.get('.wizard-header').should('not.exist');
    //todo verify entitlement
    cy.contains('moving truck').click();

    cy.nextPage();

    cy.location().should(loc => {
      expect(loc.pathname).to.match(/^\/moves\/[^/]+\/ppm-incentive/);
    });

    cy.get('.wizard-header').should('not.exist');
    cy.get('.rangeslider__handle').click();

    cy.get('.incentive').contains('$');

    cy.get('input[type="radio"]').check('yes', { force: true });
    cy.get('input[name="requested_amount"]').type('1,333.91');
    cy.get('select[name="method_of_receipt"]').select('MilPay');
    cy.nextPage();

    cy.location().should(loc => {
      expect(loc.pathname).to.match(/^\/moves\/[^/]+\/review/);
    });
    cy.get('.wizard-header').should('not.exist');

    // //todo: should probably have test suite for review and edit screens
    cy.contains('$1,333.91'); // Verify that the advance matches what was input
    cy.contains('Storage: Not requested'); // Verify SIT on the ppm review page since it's optional on HHG_PPM

    cy.nextPage();

    cy.location().should(loc => {
      expect(loc.pathname).to.match(/^\/moves\/[^/]+\/agreement/);
    });
    cy.get('.wizard-header').should('not.exist');

    cy.get('input[name="signature"]').type('Jane Doe');

    cy.nextPage();

    cy.location().should(loc => {
      expect(loc.pathname).to.match(/^\/$/);
    });

    cy.get('.usa-alert-success').within(() => {
      cy.contains('Congrats - your move is submitted!');
      cy.contains('Next, wait for approval. Once approved:');
      cy
        .get('a')
        .contains('PPM info sheet')
        .should('have.attr', 'href')
        .and('include', '/downloads/ppm_info_sheet.pdf');
    });

    cy.get('.usa-width-three-fourths').within(() => {
      cy.contains('Next Step: Wait for approval');
      cy
        .contains('Go to weight scales')
        .children('a')
        .should('have.attr', 'href', 'https://move.mil/resources/locator-maps');
      cy.contains('Advance Requested: $1,333.91');
    });
  });
});

describe('check invalid ppm inputs', () => {
  it('doesnt allow SM to progress if dont have rate data for move dates + zips"', function() {
    cy.signInAsUserPostRequest(milmoveAppName, '99360a51-8cfa-4e25-ae57-24e66077305f');

    cy.contains('Continue Move Setup').click();
    cy.location().should(loc => {
      expect(loc.pathname).to.match(/^\/moves\/[^/]+\/ppm-start/);
    });
    cy.get('.wizard-header').should('not.exist');
    cy
      .get('input[name="original_move_date"]')
      .type('6/3/2100')
      .blur();
    // test an invalid pickup zip code
    cy
      .get('input[name="pickup_postal_code"]')
      .clear()
      .type('00000')
      .blur();
    cy.get('#pickup_postal_code-error').should('exist');

    cy
      .get('input[name="pickup_postal_code"]')
      .clear()
      .type('80913');

    // test an invalid destination zip code
    cy
      .get('input[name="destination_postal_code"]')
      .clear()
      .type('00000')
      .blur();
    cy.get('#destination_postal_code-error').should('exist');

    cy
      .get('input[name="destination_postal_code"]')
      .clear()
      .type('30813');
    cy.nextPage();

    cy.location().should(loc => {
      expect(loc.pathname).to.match(/^\/moves\/[^/]+\/ppm-start/);
    });

    cy.get('#original_move_date-error').should('exist');
  });

  it('doesnt allow same origin and destination zip', function() {
    cy.signInAsUserPostRequest(milmoveAppName, '99360a51-8cfa-4e25-ae57-24e66077305f');
    cy.contains('Continue Move Setup').click();
    cy.location().should(loc => {
      expect(loc.pathname).to.match(/^\/moves\/[^/]+\/ppm-start/);
    });
    cy.get('.wizard-header').should('not.exist');
    cy
      .get('input[name="original_move_date"]')
      .first()
      .type('9/2/2018{enter}')
      .blur();
    cy
      .get('input[name="pickup_postal_code"]')
      .clear()
      .type('80913');
    cy
      .get('input[name="destination_postal_code"]')
      .type('80913')
      .blur();

    cy.get('#destination_postal_code-error').should('exist');
  });
});

describe('editing ppm only move', () => {
  it('sees only details relevant to PPM only move', () => {
    cy.signInAsUserPostRequest(milmoveAppName, 'e10d5964-c070-49cb-9bd1-eaf9f7348eb6');
    cy
      .get('.sidebar button')
      .contains('Edit Move')
      .click();

    cy.get('.ppm-container').should(ppmContainer => {
      expect(ppmContainer).to.have.length(1);
      expect(ppmContainer).to.not.have.class('hhg-shipment-summary');
    });
  });
});

describe('allows a SM to continue requesting a payment', function() {
  const smId = '4ebc03b7-c801-4c0d-806c-a95aed242102';
  beforeEach(() => {
    cy.removeFetch();
    cy.server();
    cy.route('POST', '**/internal/uploads').as('postUploadDocument');
    cy.route('POST', '**/moves/**/weight_ticket').as('postWeightTicket');
    cy.route('POST', '**/moves/**/moving_expense_documents').as('postMovingExpense');
    cy.route('POST', '**/internal/personally_procured_move/**/request_payment').as('requestPayment');
    cy.route('POST', '**/moves/**/signed_certifications').as('signedCertifications');
    cy.signInAsUserPostRequest(milmoveAppName, smId);
  });

  it('service should be able to continue requesting payment', () => {
    serviceMemberStartsPPMPaymentRequest();
    serviceMemberSubmitsWeightTicket('CAR', true);

    cy
      .get('button')
      .contains('Cancel')
      .click();
    cy
      .get('.usa-button-secondary')
      .contains('Continue Requesting Payment')
      .click();
    cy.location().should(loc => {
      expect(loc.pathname).to.match(/^\/moves\/[^/]+\/ppm-payment-review/);
    });
  });
});
describe('allows a SM to request a payment', function() {
  const smId = '745e0eba-4028-4c78-a262-818b00802748';
  const moveID = 'f9f10492-587e-43b3-af2a-9f67d2ac8757';
  beforeEach(() => {
    cy.removeFetch();
    cy.server();
    cy.route('POST', '**/internal/uploads').as('postUploadDocument');
    cy.route('POST', '**/moves/**/weight_ticket').as('postWeightTicket');
    cy.route('POST', '**/moves/**/moving_expense_documents').as('postMovingExpense');
    cy.route('POST', '**/internal/personally_procured_move/**/request_payment').as('requestPayment');
    cy.route('POST', '**/moves/**/signed_certifications').as('signedCertifications');
    cy.signInAsUserPostRequest(milmoveAppName, smId);
  });

  it('service member reads introduction to ppm payment and cancels to go back to homepage', () => {
    serviceMemberStartsPPMPaymentRequestWithAssertions();
  });

  it('service member goes through entire request payment flow', () => {
    serviceMemberStartsPPMPaymentRequest();
    serviceMemberSubmitsWeightTicket('CAR', false);
    serviceMemberViewsExpensesLandingPage();
    serviceMemberUploadsExpenses(false);
    serviceMemberReviewsDocuments();
    serviceMemberEditsPaymentRequest();
  });

  it('service member can save a weight ticket for later', () => {
    cy.visit(`/moves/${moveID}/ppm-weight-ticket`);
    serviceMemberSavesWeightTicketForLater('BOX_TRUCK');
  });

  it('service member submits weight tickets without any documents', () => {
    cy.visit(`/moves/${moveID}/ppm-weight-ticket`);
    serviceMemberSubmitsWeightsTicketsWithoutReceipts();
  });

  it('service member requests a car + trailer weight ticket payment', () => {
    cy.visit(`/moves/${moveID}/ppm-weight-ticket`);
    serviceMemberSubmitsCarTrailerWeightTicket();
  });

  it('makes missing weight ticket fields optional when missing is checked', () => {
    // Always required fields
    cy.visit(`/moves/${moveID}/ppm-weight-ticket`);
    cy.get('select[name="vehicle_options"]').select('CAR');
    cy.get('input[name="vehicle_nickname"]').type('Nickname');

    // only required when missing not checked
    cy.get('input[name="missingEmptyWeightTicket"]+label').click();

    // only required when missing is not checked
    cy.get('input[name="full_weight"]').type('5000');
    cy.upload_file('[data-cy=full-weight-upload] .filepond--root', 'top-secret.png');
    cy.wait('@postUploadDocument');
    cy.get('[data-filepond-item-state="processing-complete"]').should('have.length', 1);
    cy
      .get('input[name="weight_ticket_date"]')
      .type('6/2/2018{enter}')
      .blur();

    cy.get('input[name="additional_weight_ticket"][value="Yes"]').should('not.be.checked');
    cy.get('input[name="additional_weight_ticket"][value="No"]').should('be.checked');
    cy.get('input[name="additional_weight_ticket"][value="Yes"]+label').click();
    cy
      .get('button')
      .contains('Save & Add Another')
      .should('be.enabled');
  });

  it('makes full weight ticket fields optional when missing is checked', () => {
    // Always required fields
    cy.visit(`/moves/${moveID}/ppm-weight-ticket`);
    cy.get('select[name="vehicle_options"]').select('CAR');
    cy.get('input[name="vehicle_nickname"]').type('Nickname');

    // only required when missing not checked
    cy.get('input[name="empty_weight"]').type('1000');
    cy.upload_file('[data-cy=empty-weight-upload] .filepond--root', 'top-secret.png');
    cy.wait('@postUploadDocument');
    cy.get('[data-filepond-item-state="processing-complete"]').should('have.length', 1);

    // only required when missing is not checked
    cy.get('input[name="missingFullWeightTicket"]+label').click();

    cy.get('input[name="additional_weight_ticket"][value="Yes"]').should('not.be.checked');
    cy.get('input[name="additional_weight_ticket"][value="No"]').should('be.checked');
    cy.get('input[name="additional_weight_ticket"][value="Yes"]+label').click();
    cy
      .get('button')
      .contains('Save & Add Another')
      .should('be.enabled');
  });

  it('service member starting at review page returns to review page after adding a weight ticket', () => {
    cy.visit(`/moves/${moveID}/ppm-payment-review`);
    cy.location().should(loc => {
      expect(loc.pathname).to.match(/^\/moves\/[^/]+\/ppm-payment-review/);
    });
    cy.get('[data-cy="weight-ticket-link"]').click();
    cy.location().should(loc => {
      expect(loc.pathname).to.match(/^\/moves\/[^/]+\/ppm-weight-ticket/);
    });
    serviceMemberSubmitsWeightTicket('CAR', false);
    cy.location().should(loc => {
      expect(loc.pathname).to.match(/^\/moves\/[^/]+\/ppm-payment-review/);
    });
  });

  it('service member starting at review page returns to review page after adding an expense', () => {
    cy.visit(`/moves/${moveID}/ppm-payment-review`);
    cy.location().should(loc => {
      expect(loc.pathname).to.match(/^\/moves\/[^/]+\/ppm-payment-review/);
    });
    cy.get('[data-cy="expense-link"]').click();
    cy.location().should(loc => {
      expect(loc.pathname).to.match(/^\/moves\/[^/]+\/ppm-expenses/);
    });
    serviceMemberUploadsExpenses(false);
    cy.location().should(loc => {
      expect(loc.pathname).to.match(/^\/moves\/[^/]+\/ppm-payment-review/);
    });
  });

  it('service member can skip weight tickets and expenses if already have one', () => {
    cy.visit(`/moves/${moveID}/ppm-weight-ticket`);

    serviceMemberSubmitsWeightTicket('CAR', true);
    cy
      .get('[data-cy=skip]')
      .contains('Skip')
      .click();
    serviceMemberViewsExpensesLandingPage();
    serviceMemberUploadsExpenses();
    cy
      .get('[data-cy=skip]')
      .contains('Skip')
      .click();
  });

  //TODO: remove when done with the new flow to request payment
  it('service member submits request for payment', function() {
    cy.removeFetch();
    cy.server();
    cy.route('POST', '**/internal/uploads').as('postUploadDocument');
    const stub = cy.stub();
    cy.on('window:alert', stub);

    cy.logout();
    //profile@comple.te
    cy.signInAsUserPostRequest(milmoveAppName, '8e0d7e98-134e-4b28-bdd1-7d6b1ff34f9e');
    cy.setFeatureFlag('ppmPaymentRequest=false', '/');
    cy.contains('Fort Gordon (from Yuma AFB)');
    cy.contains('Request Payment').click();

    cy.location().should(loc => {
      expect(loc.pathname).to.match(/^\/moves\/[^/]+\/request-payment/);
    });

    cy.get('input[type="checkbox"]').should('not.be.checked');

    cy
      .contains('Legal Agreement / Privacy Act')
      .click()
      .then(() => {
        expect(stub.getCall(0)).to.be.calledWithMatch('LEGAL AGREEMENT / PRIVACY ACT');
      });
    cy.get('input[type="checkbox"]').should('not.be.checked');
    cy.get('select[name="move_document_type"]').select('WEIGHT_TICKET');
    cy.get('input[name="title"]').type('WEIGHT_TICKET');
    cy.upload_file('.filepond--root', 'top-secret.png');
    cy.wait('@postUploadDocument');
    cy
      .get('button')
      .contains('Save')
      .click();
    cy.get('input[id="agree-checkbox"]').check({ force: true });
    cy
      .get('button')
      .contains('Submit Payment')
      .click();
    cy.location().should(loc => {
      expect(loc.pathname).to.match(/^\/$/);
    });
  });
});

function serviceMemberReviewsDocuments() {
  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/moves\/[^/]+\/ppm-payment-review/);
  });
  cy
    .get('.review-customer-agreement a')
    .contains('Legal Agreement')
    .click();
  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/ppm-customer-agreement/);
  });
  cy
    .get('.usa-button-secondary')
    .contains('Back')
    .click();
  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/moves\/[^/]+\/ppm-payment-review/);
  });
  cy.get('input[id="agree-checkbox"]').check({ force: true });
  cy
    .get('button')
    .contains('Submit Request')
    .should('be.enabled')
    .click();
  cy.wait('@signedCertifications');
  cy.wait('@requestPayment');
}
function serviceMemberEditsPaymentRequest() {
  cy
    .get('.usa-alert-success')
    .contains('Payment request submitted')
    .should('exist');
  cy
    .get('.usa-button-secondary')
    .contains('Edit Payment Request')
    .should('exist')
    .click();
  cy
    .get('[data-cy=weight-ticket-link]')
    .should('exist')
    .click();
  serviceMemberSubmitsWeightTicket('CAR', false);
  serviceMemberReviewsDocuments();
}
function serviceMemberViewsExpensesLandingPage() {
  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/moves\/[^/]+\/ppm-expenses-intro/);
  });

  cy.get('[data-cy=documents-uploaded]').should('exist');
  cy
    .get('button')
    .contains('Continue')
    .should('be.disabled');

  cy
    .get('[type="radio"]')
    .first()
    .should('be.not.checked');
  cy
    .get('[type="radio"]')
    .last()
    .should('be.not.checked');

  cy
    .get('a')
    .contains('More about expenses')
    .should('have.attr', 'href')
    .and('match', /^\/allowable-expenses/);

  cy.get('input[name="hasExpenses"][value="Yes"]').should('not.be.checked');
  cy.get('input[name="hasExpenses"][value="No"]').should('not.be.checked');
  cy.get('input[name="hasExpenses"][value="Yes"]+label').click();

  cy
    .get('button')
    .contains('Continue')
    .should('be.enabled')
    .click();
}

function serviceMemberUploadsExpenses(hasAnother = true, expenseNumber = null) {
  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/moves\/[^/]+\/ppm-expenses/);
  });

  if (expenseNumber) {
    cy.contains(`Expense ${expenseNumber}`);
  }
  cy.get('[data-cy=documents-uploaded]').should('exist');

  cy.get('select[name="moving_expense_type"]').select('GAS');
  cy.get('input[name="title"]').type('title');
  cy.get('input[name="requested_amount_cents"]').type('1000');

  cy.upload_file('.filepond--root:first', 'top-secret.png');
  cy.wait('@postUploadDocument');
  cy.get('[data-filepond-item-state="processing-complete"]').should('have.length', 1);

  cy.get('input[name="missingReceipt"]').should('not.be.checked');
  cy.get('input[name="paymentMethod"][value="GTCC"]').should('be.checked');
  cy.get('input[name="paymentMethod"][value="OTHER"]').should('not.be.checked');
  cy.get('input[name="haveMoreExpenses"][value="Yes"]').should('not.be.checked');
  cy.get('input[name="haveMoreExpenses"][value="No"]').should('be.checked');
  cy.get('input[name="haveMoreExpenses"][value="Yes"]+label').click();
  if (hasAnother) {
    cy
      .get('button')
      .contains('Save & Add Another')
      .click();
    cy
      .wait('@postMovingExpense')
      .its('status')
      .should('eq', 200);
    cy.get('[data-cy=documents-uploaded]').should('exist');
  } else {
    cy.get('input[name="haveMoreExpenses"][value="No"]+label').click();
    cy.get('input[name="haveMoreExpenses"][value="No"]').should('be.checked');
    cy
      .get('button')
      .contains('Save & Continue')
      .click();
    cy
      .wait('@postMovingExpense')
      .its('status')
      .should('eq', 200);
  }
}

function serviceMemberSubmitsCarTrailerWeightTicket() {
  cy.get('select[name="vehicle_options"]').select('CAR_TRAILER');

  cy.get('input[name="vehicle_nickname"]').type('Nickname');

  cy
    .contains('Do you own this trailer')
    .children('a')
    .should('have.attr', 'href', '/trailer-criteria');

  cy.get('input[name="isValidTrailer"][value="Yes"]').should('not.be.checked');
  cy.get('input[name="isValidTrailer"][value="No"]').should('be.checked');
  cy.get('input[name="isValidTrailer"][value="Yes"]+label').click();

  cy.upload_file('[data-cy=trailer-upload] .filepond--root', 'top-secret.png');
  cy.wait('@postUploadDocument');
  cy.get('[data-filepond-item-state="processing-complete"]').should('have.length', 1);

  cy.get('input[name="missingDocumentation"]').check({ force: true });

  cy.get('input[name="empty_weight"]').type('1000');
  cy.get('input[name="missingEmptyWeightTicket"]').check({ force: true });

  cy.get('input[name="full_weight"]').type('5000');
  cy.upload_file('.filepond--root:last', 'top-secret.png');
  cy.wait('@postUploadDocument');
  cy.get('[data-filepond-item-state="processing-complete"]').should('have.length', 2);
  cy
    .get('input[name="weight_ticket_date"]')
    .type('6/2/2018{enter}')
    .blur();

  cy.get('input[name="additional_weight_ticket"][value="Yes"]').should('not.be.checked');
  cy.get('input[name="additional_weight_ticket"][value="No"]').should('be.checked');
}
function serviceMemberSavesWeightTicketForLater(vehicleType) {
  cy.get('select[name="vehicle_options"]').select(vehicleType);

  cy.get('input[name="vehicle_nickname"]').type('Nickname');

  cy.get('input[name="empty_weight"]').type('1000');
  cy.upload_file('[data-cy=empty-weight-upload] .filepond--root', 'top-secret.png');
  cy.wait('@postUploadDocument');
  cy.get('[data-filepond-item-state="processing-complete"]').should('have.length', 1);

  cy.get('input[name="full_weight"]').type('5000');
  cy.upload_file('[data-cy=full-weight-upload] .filepond--root', 'top-secret.png');
  cy.wait('@postUploadDocument');
  cy.get('[data-filepond-item-state="processing-complete"]').should('have.length', 2);
  cy
    .get('input[name="weight_ticket_date"]')
    .type('6/2/2018{enter}')
    .blur();

  cy
    .get('button')
    .contains('Save For Later')
    .click();
  cy
    .wait('@postWeightTicket')
    .its('status')
    .should('eq', 200);
  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/$/);
  });
}

function serviceMemberSubmitsWeightsTicketsWithoutReceipts() {
  cy.get('select[name="vehicle_options"]').select('CAR_TRAILER');
  cy.get('input[name="vehicle_nickname"]').type('Nickname');
  cy.get('input[name="empty_weight"]').type('1000');
  cy.get('input[name="full_weight"]').type('2000');
  cy.get('input[name="isValidTrailer"][value="Yes"]+label').click();
  cy.get('input[name="missingDocumentation"]+label').click();
  cy
    .get('[data-cy=trailer-warning]')
    .contains(
      'If your state does not provide a registration or bill of sale for your trailer, you may write and upload a signed and dated statement certifying that you or your spouse own the trailer and meets the trailer criteria. Upload your statement using the proof of ownership field.',
    );
  cy.get('input[name="missingEmptyWeightTicket"]+label').click();
  cy
    .get('[data-cy=empty-warning]')
    .contains(
      'Contact your local Transportation Office (PPPO) to let them know you’re missing this weight ticket. For now, keep going and enter the info you do have.',
    );
  cy.get('input[name="missingFullWeightTicket"]+label').click();
  cy
    .get('[data-cy=full-warning]')
    .contains(
      'Contact your local Transportation Office (PPPO) to let them know you’re missing this weight ticket. For now, keep going and enter the info you do have.',
    );
  cy
    .get('input[name="weight_ticket_date"]')
    .type('6/2/2018{enter}')
    .blur();

  cy.get('input[name="additional_weight_ticket"][value="Yes"]+label').click();
  cy.get('input[name="additional_weight_ticket"][value="Yes"]').should('be.checked');
  cy
    .get('button')
    .contains('Save & Add Another')
    .click();
  cy.wait('@postWeightTicket');
}

function serviceMemberStartsPPMPaymentRequest() {
  cy.contains('Request Payment').click();
  cy
    .get('button')
    .contains('Get Started')
    .click();
}

function serviceMemberSubmitsWeightTicket(vehicleType, hasAnother = true, ordinal = null) {
  if (ordinal) {
    cy.contains(`Weight Tickets - ${ordinal} set`);
    if (ordinal === '1st') {
      cy.get('[data-cy=documents-uploaded]').should('not.exist');
    } else {
      cy.get('[data-cy=documents-uploaded]').should('exist');
    }
  }

  cy.get('select[name="vehicle_options"]').select(vehicleType);

  cy.get('input[name="vehicle_nickname"]').type('Nickname');

  cy.get('input[name="empty_weight"]').type('1000');

  cy.upload_file('[data-cy=empty-weight-upload] .filepond--root', 'top-secret.png');
  cy.wait('@postUploadDocument');
  cy.get('[data-filepond-item-state="processing-complete"]').should('have.length', 1);

  cy.get('input[name="full_weight"]').type('5000');
  cy.upload_file('[data-cy=full-weight-upload] .filepond--root', 'top-secret.png');
  cy.wait('@postUploadDocument');
  cy.get('[data-filepond-item-state="processing-complete"]').should('have.length', 2);
  cy
    .get('input[name="weight_ticket_date"]')
    .type('6/2/2018{enter}')
    .blur();
  cy.get('input[name="additional_weight_ticket"][value="Yes"]').should('not.be.checked');
  cy.get('input[name="additional_weight_ticket"][value="No"]').should('be.checked');
  if (hasAnother) {
    cy.get('input[name="additional_weight_ticket"][value="Yes"]+label').click();
    cy.get('input[name="additional_weight_ticket"][value="Yes"]').should('be.checked');
    cy
      .get('button')
      .contains('Save & Add Another')
      .click();
    cy
      .wait('@postWeightTicket')
      .its('status')
      .should('eq', 200);
    cy.get('[data-cy=documents-uploaded]').should('exist');
  } else {
    cy
      .get('button')
      .contains('Save & Continue')
      .click();
    cy
      .wait('@postWeightTicket')
      .its('status')
      .should('eq', 200);
  }
}

function serviceMemberStartsPPMPaymentRequestWithAssertions() {
  cy.contains('Request Payment').click();

  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/moves\/[^/]+\/ppm-payment-request-intro/);
  });

  cy.get('h3').contains('Request PPM Payment');

  cy.get('.weight-ticket-examples-link').click();

  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/weight-ticket-examples/);
  });

  cy.get('h3').contains('Example weight ticket scenarios');

  cy
    .get('button')
    .contains('Back')
    .click();

  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/moves\/[^/]+\/ppm-payment-request-intro/);
  });

  cy
    .get('a')
    .contains('More about expenses')
    .click();

  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/allowable-expenses/);
  });

  cy.get('h3').contains('Storage & Moving Expenses');

  cy
    .get('button')
    .contains('Back')
    .click();

  cy
    .get('button')
    .contains('Get Started')
    .click();
}
