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

  //TODO: remove when done with the new flow to request payment
  it('allows a SM to request payment', function() {
    cy.removeFetch();
    cy.server();
    cy.route('POST', '**/internal/uploads').as('postUploadDocument');
    const stub = cy.stub();
    cy.on('window:alert', stub);

    cy.logout();
    //profile@comple.te
    cy.signInAsUserPostRequest(milmoveAppName, '8e0d7e98-134e-4b28-bdd1-7d6b1ff34f9e');
    cy.setFeatureFlag('ppmPaymentRequest', '/');
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
      .first()
      .type('6/3/2100{enter}');
    cy
      .get('input[name="pickup_postal_code"]')
      .clear()
      .type('80913');
    cy.get('input[name="destination_postal_code"]').type('76127');
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

it('allows a SM to request ppm payment', function() {
  cy.removeFetch();
  cy.server();
  cy.route('POST', '**/internal/uploads').as('postUploadDocument');

  cy.signInAsUserPostRequest(milmoveAppName, '8e0d7e98-134e-4b28-bdd1-7d6b1ff34f9e');
  cy.contains('Fort Gordon (from Yuma AFB)');
  cy.get('.submitted .status_dates').should('exist');
  cy.get('.ppm_approved .status_dates').should('exist');
  cy.get('.in_progress .status_dates').should('exist');
  serviceMemberCanCancel();
  serviceMemberVisitsIntroToPPMPaymentRequest();
  serviceMemberFillsOutWeightTicket('Car');
  // TODO: remove when we are doing something with the data
  cy.reload();
  serviceMemberFillsOutWeightTicket('Box truck');
});

function serviceMemberFillsOutWeightTicket(vehicleType) {
  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/moves\/[^/]+\/ppm-weight-ticket/);
  });

  cy.get('select[name="vehicle_options"]').select(vehicleType);

  cy.get('input[name="vehicle_nickname"]').type('Nickname');

  cy.get('input[name="empty_weight"]').type('1000');
  cy.upload_file('.filepond--root:first', 'top-secret.png');
  cy.wait('@postUploadDocument');

  cy.get('input[name="full_weight"]').type('5000');
  cy.upload_file('.filepond--root:last', 'top-secret.png');
  cy.wait('@postUploadDocument');
  cy
    .get('input[name="weight_ticket_date"]')
    .type('6/2/2018{enter}')
    .blur();

  cy
    .get('[type="radio"]')
    .first()
    .should('be.checked');
}

function serviceMemberCanCancel() {
  cy.contains('Request Payment').click();

  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/moves\/[^/]+\/ppm-payment-request-intro/);
  });
  cy
    .get('button')
    .contains('Cancel')
    .click();

  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\//);
  });
}

function serviceMemberVisitsIntroToPPMPaymentRequest() {
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
    .contains('What expenses are allowed?')
    .click();

  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/allowable-expenses/);
  });

  cy.get('h3').contains('Allowable expenses');

  cy
    .get('button')
    .contains('Back')
    .click();

  cy
    .get('button')
    .contains('Get Started')
    .click();
}
