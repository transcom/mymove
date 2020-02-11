/* global cy */

import { milmoveAppName } from '../../support/constants';

describe('completing the ppm flow', function() {
  it('progresses thru forms', function() {
    //profile@comple.te
    cy.signInAsUserPostRequest(milmoveAppName, '13f3949d-0d53-4be4-b1b1-ae4314793f34');
    cy.contains('Fort Gordon (from Yuma AFB)');
    cy.get('[data-cy="move-header-weight-estimate"]').contains('8,000 lbs');
    cy.contains('Continue Move Setup').click();

    cy.location().should(loc => {
      expect(loc.pathname).to.match(/^\/moves\/[^/]+\/ppm-start/);
    });
    cy.get('.wizard-header').should('not.exist');
    cy.get('input[name="original_move_date"]')
      .first()
      .type('9/2/2018{enter}')
      .blur();
    cy.get('input[name="pickup_postal_code"]')
      .clear()
      .type('80913');

    cy.get('input[name="destination_postal_code"]')
      .clear()
      .type('76127');

    cy.get('input[type="radio"][value="yes"]')
      .eq(1)
      .check('yes', { force: true });
    cy.get('input[name="days_in_storage"]')
      .clear()
      .type('30');

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

    cy.nextPage();

    cy.location().should(loc => {
      expect(loc.pathname).to.match(/^\/moves\/[^/]+\/review/);
    });
    cy.get('.wizard-header').should('not.exist');

    // todo: should probably have test suite for review and edit screens
    cy.get('[data-cy="sit-display"]')
      .contains('30 days')
      .contains('$2328.64');

    cy.get('[data-cy="edit-ppm-dates"]').click();

    cy.location().should(loc => {
      expect(loc.pathname).to.match(/^\/moves\/[^/]+\/review\/edit-date-and-location/);
    });

    cy.get('[data-cy="storage-estimate"]').contains('$2328.64');

    cy.get('input[name="days_in_storage"]')
      .clear()
      .type('35');

    cy.get('[data-cy="storage-estimate"]').contains('$2328.64');

    cy.get('button')
      .contains('Save')
      .click();

    cy.location().should(loc => {
      expect(loc.pathname).to.match(/^\/moves\/[^/]+\/review/);
    });

    cy.get('[data-cy="sit-display"]')
      .contains('35 days')
      .contains('$2405.14');

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

    cy.get('.usa-alert--success').within(() => {
      cy.contains('Congrats - your move is submitted!');
      cy.contains('Next, wait for approval. Once approved:');
      cy.get('a')
        .contains('PPM info sheet')
        .should('have.attr', 'href')
        .and('include', '/downloads/ppm_info_sheet.pdf');
    });
  });
});

describe('completing the ppm flow with a move date that we currently do not have rates for', function() {
  it('complete a PPM move', function() {
    //profile@complete.draft
    cy.signInAsUserPostRequest(milmoveAppName, '3b9360a3-3304-4c60-90f4-83d687884070');
    cy.contains('Fort Gordon (from Yuma AFB)');
    cy.get('[data-cy="move-header-weight-estimate"]').contains('8,000 lbs');
    cy.contains('Continue Move Setup').click();

    cy.location().should(loc => {
      expect(loc.pathname).to.match(/^\/moves\/[^/]+\/ppm-start/);
    });
    cy.get('.wizard-header').should('not.exist');
    cy.get('input[name="original_move_date"]')
      .first()
      .type('9/2/2030{enter}')
      .blur();
    cy.get('input[name="pickup_postal_code"]')
      .clear()
      .type('80913');

    cy.get('input[name="destination_postal_code"]')
      .clear()
      .type('76127');

    cy.get('input[type="radio"][value="yes"]')
      .eq(1)
      .check('yes', { force: true });
    cy.get('input[name="days_in_storage"]')
      .clear()
      .type('30');

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

    cy.get('.incentive').contains('Not ready yet');
    cy.get('[data-icon="question-circle"]').click();
    cy.get('[data-cy="tooltip"]').contains(
      'We expect to receive rate data covering your move dates by the end of this month. Check back then to see your estimated incentive.',
    );
    cy.nextPage();

    cy.location().should(loc => {
      expect(loc.pathname).to.match(/^\/moves\/[^/]+\/review/);
    });
    cy.get('.wizard-header').should('not.exist');
    cy.get('td').contains('Not ready yet');
    cy.get('[data-icon="question-circle"]').click();
    cy.get('[data-cy="tooltip"]').contains(
      'We expect to receive rate data covering your move dates by the end of this month. Check back then to see your estimated incentive.',
    );

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

    cy.get('.usa-alert--success').within(() => {
      cy.contains('Congrats - your move is submitted!');
      cy.contains('Next, wait for approval. Once approved:');
      cy.get('a')
        .contains('PPM info sheet')
        .should('have.attr', 'href')
        .and('include', '/downloads/ppm_info_sheet.pdf');
    });

    cy.contains('Incentive (est.): Not ready yet');
    cy.get('[data-icon="question-circle"]').click();
    cy.get('[data-cy="tooltip"]').contains(
      'We expect to receive rate data covering your move dates by the end of this month. Check back then to see your estimated incentive.',
    );

    cy.get('[data-cy="edit-move"]')
      .contains('Edit Move')
      .click();

    cy.get('td').contains('Not ready yet');
    cy.get('[data-icon="question-circle"]').click();
    cy.get('[data-cy="tooltip"]').contains(
      'We expect to receive rate data covering your move dates by the end of this month. Check back then to see your estimated incentive.',
    );
  });
});

describe('check invalid ppm inputs', () => {
  it('doesnt allow same origin and destination zip', function() {
    cy.signInAsUserPostRequest(milmoveAppName, '99360a51-8cfa-4e25-ae57-24e66077305f');
    cy.contains('Continue Move Setup').click();
    cy.location().should(loc => {
      expect(loc.pathname).to.match(/^\/moves\/[^/]+\/ppm-start/);
    });
    cy.get('.wizard-header').should('not.exist');
    cy.get('input[name="original_move_date"]')
      .first()
      .type('9/2/2018{enter}')
      .blur();
    cy.get('input[name="pickup_postal_code"]')
      .clear()
      .type('80913');
    cy.get('input[name="destination_postal_code"]')
      .type('80913')
      .blur();

    cy.get('#destination_postal_code-error').should('exist');
  });

  it('doesnt allow SM to progress if dont have rate data for zips"', function() {
    cy.signInAsUserPostRequest(milmoveAppName, '99360a51-8cfa-4e25-ae57-24e66077305f');

    cy.contains('Continue Move Setup').click();
    cy.location().should(loc => {
      expect(loc.pathname).to.match(/^\/moves\/[^/]+\/ppm-start/);
    });
    cy.get('.wizard-header').should('not.exist');
    cy.get('input[name="original_move_date"]')
      .type('6/3/2100')
      .blur();
    // test an invalid pickup zip code
    cy.get('input[name="pickup_postal_code"]')
      .clear()
      .type('00000')
      .blur();
    cy.get('#pickup_postal_code-error').should('exist');

    cy.get('input[name="pickup_postal_code"]')
      .clear()
      .type('80913');

    // test an invalid destination zip code
    cy.get('input[name="destination_postal_code"]')
      .clear()
      .type('00000')
      .blur();
    cy.get('#destination_postal_code-error').should('exist');

    cy.get('input[name="destination_postal_code"]')
      .clear()
      .type('30813');
    cy.nextPage();

    cy.location().should(loc => {
      expect(loc.pathname).to.match(/^\/moves\/[^/]+\/ppm-start/);
    });
  });
});

describe('editing ppm only move', () => {
  it('sees only details relevant to PPM only move', () => {
    cy.signInAsUserPostRequest(milmoveAppName, 'e10d5964-c070-49cb-9bd1-eaf9f7348eb6');
    cy.get('.sidebar button')
      .contains('Edit Move')
      .click();

    cy.get('[data-cy="ppm-summary"]').should(ppmContainer => {
      expect(ppmContainer).to.have.length(1);
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

    cy.get('button')
      .contains('Finish Later')
      .click();

    cy.get('button')
      .contains('OK')
      .click();

    cy.location().should(loc => {
      expect(loc.pathname).to.match(/^\/$/);
    });

    cy.get('a')
      .contains('Continue Requesting Payment')
      .click();
    cy.location().should(loc => {
      expect(loc.pathname).to.match(/^\/moves\/[^/]+\/ppm-payment-review/);
    });
  });
});
function serviceMemberStartsPPMPaymentRequest() {
  cy.contains('Request Payment').click();
  cy.get('input[name="actual_move_date"]')
    .type('6/20/2018{enter}')
    .blur();
  cy.get('button')
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

  cy.get('select[name="weight_ticket_set_type"]').select(vehicleType);

  cy.get('input[name="vehicle_nickname"]').type('Nickname');

  cy.get('input[name="empty_weight"]').type('1000');

  cy.upload_file('[data-cy=empty-weight-upload] .filepond--root', 'top-secret.png');
  cy.wait('@postUploadDocument');
  cy.get('[data-filepond-item-state="processing-complete"]').should('have.length', 1);

  cy.get('input[name="full_weight"]').type('5000');
  cy.upload_file('[data-cy=full-weight-upload] .filepond--root', 'top-secret.png');
  cy.wait('@postUploadDocument');
  cy.get('[data-filepond-item-state="processing-complete"]').should('have.length', 2);
  cy.get('input[name="weight_ticket_date"]')
    .type('6/2/2018{enter}')
    .blur();
  cy.get('input[name="additional_weight_ticket"][value="Yes"]').should('not.be.checked');
  cy.get('input[name="additional_weight_ticket"][value="No"]').should('be.checked');
  if (hasAnother) {
    cy.get('input[name="additional_weight_ticket"][value="Yes"]+label').click();
    cy.get('input[name="additional_weight_ticket"][value="Yes"]').should('be.checked');
    cy.get('button')
      .contains('Save & Add Another')
      .click();
    cy.wait('@postWeightTicket')
      .its('status')
      .should('eq', 200);
    cy.get('[data-cy=documents-uploaded]').should('exist');
  } else {
    cy.get('button')
      .contains('Save & Continue')
      .click();
    cy.wait('@postWeightTicket')
      .its('status')
      .should('eq', 200);
  }
}
