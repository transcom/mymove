// import { fileUploadTimeout } from '../../support/constants';

// describe('the PPM flow', function () {
//   before(() => {
//     cy.prepareCustomerApp();
//   });

//   beforeEach(() => {
//     cy.logout();
//   });

//   it('can submit a PPM move', () => {
//     // profile2@complete.draft
//     const userId = '4635b5a7-0f57-4557-8ba4-bbbb760c300a';
//     cy.apiSignInAsPpmUser(userId);
//     SMSubmitsMove();
//   });

//   it('can complete the PPM flow with a move date that we currently do not have rates for', () => {
//     // profile@complete.draft
//     const userId = '3b9360a3-3304-4c60-90f4-83d687884070';
//     cy.apiSignInAsPpmUser(userId);
//     SMCompletesMove();
//   });

//   it('doesn’t allow a user to enter the same origin and destination zip', () => {
//     // profile@co.mple.te
//     const userId = '324dec0a-850c-41c8-976b-068e27121b84';
//     cy.apiSignInAsPpmUser(userId);
//     SMInputsSamePostalCodes();
//   });

//   it('doesn’t allow SM to progress if don’t have rate data for zips"', () => {
//     // profile@co.mple.te
//     const userId = 'f154929c-5f07-41f5-b90c-d90b83d5773d';
//     cy.apiSignInAsPpmUser(userId);
//     SMInputsInvalidPostalCodes();
//   });

//   it('when editing PPM only move, sees only details relevant to PPM only move', () => {
//     // ppm@incomple.te
//     const userId = 'e10d5964-c070-49cb-9bd1-eaf9f7348eb6';
//     cy.apiSignInAsPpmUser(userId);
//     SMSeesMoveDetails();
//   });

//   it('service member should be able to continue requesting payment', () => {
//     // ppm@continue.requestingpayment
//     const userId = '4ebc03b7-c801-4c0d-806c-a95aed242102';

//     // TODO - commenting this out for now because cy.intercept has an open bug related to file upload endpoints
//     // https://github.com/cypress-io/cypress/issues/9534
//     // cy.intercept('POST', '**/internal/uploads').as('postUploadDocument');

//     cy.intercept('POST', '**/moves/**/weight_ticket').as('postWeightTicket');
//     cy.intercept('POST', '**/moves/**/moving_expense_documents').as('postMovingExpense');
//     cy.intercept('POST', '**/internal/personally_procured_move/**/request_payment').as('requestPayment');
//     cy.intercept('POST', '**/moves/**/signed_certifications').as('signedCertifications');
//     cy.apiSignInAsPpmUser(userId);
//     SMContinueRequestPayment();
//   });
// });

// function SMSubmitsMove() {
//   cy.contains('Fort Gordon (from Yuma AFB)');
//   cy.get('[data-testid="move-header-weight-estimate"]').contains('5,000 lbs');
//   cy.contains('Continue Move Setup').click();

//   cy.location().should((loc) => {
//     expect(loc.pathname).to.match(/^\/moves\/[^/]+\/ppm-start/);
//   });

//   cy.get('.wizard-header').should('not.exist');
//   cy.get('input[name="original_move_date"]').first().type('9/2/2018{enter}').blur();
//   cy.get('input[name="pickup_postal_code"]').clear().type('80913');

//   cy.get('input[name="destination_postal_code"]').clear().type('76127');

//   cy.get('input[type="radio"][value="yes"]').eq(1).check('yes', { force: true });
//   cy.get('input[name="days_in_storage"]').clear().type('30');

//   cy.nextPageAndCheckLocation(
//     'weight-page-title',
//     "How much do you think you'll move?",
//     '^\\/moves\\/[^/]+\\/ppm-incentive',
//   );

//   cy.get('.wizard-header').should('not.exist');
//   cy.get('#incentive-estimation-slider').click();

//   cy.get('[data-testid="incentive-range-values"]').contains('$');

//   cy.nextPage();

//   cy.location().should((loc) => {
//     expect(loc.pathname).to.match(/^\/moves\/[^/]+\/review/);
//   });
//   cy.get('.wizard-header').should('not.exist');

//   // todo: should probably have test suite for review and edit screens
//   cy.get('[data-testid="sit-display"]').contains('30 days');

//   cy.get('[data-testid="edit-ppm-dates"]').click();

//   cy.location().should((loc) => {
//     expect(loc.pathname).to.match(/^\/moves\/[^/]+\/review\/edit-date-and-location/);
//   });

//   cy.get('input[name="days_in_storage"]').clear().type('35');

//   cy.get('[data-testid="storage-estimate"]').contains('$736.91');

//   cy.get('button').contains('Save').click();

//   cy.location().should((loc) => {
//     expect(loc.pathname).to.match(/^\/moves\/[^/]+\/review/);
//   });

//   cy.nextPage();

//   cy.location().should((loc) => {
//     expect(loc.pathname).to.match(/^\/moves\/[^/]+\/agreement/);
//   });
//   cy.get('.wizard-header').should('not.exist');

//   cy.get('input[name="signature"]').type('Jane Doe');

//   cy.completeFlow();

//   cy.location().should((loc) => {
//     expect(loc.pathname).to.match(/^\/$/);
//   });

//   cy.get('.usa-alert--success').within(() => {
//     cy.contains('You’ve submitted your move request.');
//   });
// }

// function SMCompletesMove() {
//   cy.contains('Fort Gordon (from Yuma AFB)');
//   cy.get('[data-testid="move-header-weight-estimate"]').contains('8,000 lbs');
//   cy.contains('Continue Move Setup').click();

//   cy.location().should((loc) => {
//     expect(loc.pathname).to.match(/^\/moves\/[^/]+\/ppm-start/);
//   });

//   cy.get('.wizard-header').should('not.exist');
//   cy.get('input[name="original_move_date"]').first().type('9/2/2030{enter}').blur();
//   cy.get('input[name="pickup_postal_code"]').clear().type('80913');

//   cy.get('input[name="destination_postal_code"]').clear().type('76127');

//   cy.get('input[type="radio"][value="yes"]').eq(1).check('yes', { force: true });
//   cy.get('input[name="days_in_storage"]').clear().type('30');

//   cy.nextPageAndCheckLocation(
//     'weight-page-title',
//     "How much do you think you'll move?",
//     '^\\/moves\\/[^/]+\\/ppm-incentive',
//   );

//   cy.get('.wizard-header').should('not.exist');
//   cy.get('#incentive-estimation-slider').click();

//   // we are hardcoding PPM rates so this will never come with this message
//   // cy.get('[data-icon="question-circle"]').click();
//   // cy.get('[data-testid="tooltip"]').contains(
//   //   'We expect to receive rate data covering your move dates by the end of this month. Check back then to see your estimated incentive.',
//   // );
//   cy.nextPage();

//   cy.location().should((loc) => {
//     expect(loc.pathname).to.match(/^\/moves\/[^/]+\/review/);
//   });
//   cy.get('.wizard-header').should('not.exist');
//   cy.get('dd').contains('Rate info unavailable');

//   cy.nextPage();

//   cy.location().should((loc) => {
//     expect(loc.pathname).to.match(/^\/moves\/[^/]+\/agreement/);
//   });
//   cy.get('.wizard-header').should('not.exist');

//   cy.get('input[name="signature"]').type('Jane Doe');

//   cy.completeFlow();

//   cy.location().should((loc) => {
//     expect(loc.pathname).to.match(/^\/$/);
//   });

//   cy.get('.usa-alert--success').within(() => {
//     cy.contains('You’ve submitted your move request.');
//   });

//   cy.visit('/ppm');
//   cy.contains('Payment: Not ready yet');
//   cy.get('[data-icon="question-circle"]').click();
//   cy.get('[data-testid="tooltip"]').contains(
//     'We expect to receive rate data covering your move dates by the end of this month. Check back then to see your estimated incentive.',
//   );

//   // cy.get('[data-testid="edit-move"]').contains('Edit Move').click();

//   // cy.get('td').contains('Not ready yet');
//   // cy.get('[data-icon="question-circle"]').click();
//   // cy.get('[data-testid="tooltip"]').contains(
//   //   'We expect to receive rate data covering your move dates by the end of this month. Check back then to see your estimated incentive.',
//   // );
// }

// function SMInputsSamePostalCodes() {
//   cy.contains('Continue Move Setup').click();
//   cy.location().should((loc) => {
//     expect(loc.pathname).to.match(/^\/moves\/[^/]+\/ppm-start/);
//   });

//   cy.get('.wizard-header').should('not.exist');
//   cy.get('input[name="original_move_date"]').first().type('9/2/2018{enter}').blur();
//   cy.get('input[name="pickup_postal_code"]').clear().type('80913');
//   cy.get('input[name="destination_postal_code"]').type('80913').blur();

//   cy.get('#destination_postal_code-error').should('exist');
// }

// function SMInputsInvalidPostalCodes() {
//   cy.contains('Continue Move Setup').click();
//   cy.location().should((loc) => {
//     expect(loc.pathname).to.match(/^\/moves\/[^/]+\/ppm-start/);
//   });
//   cy.get('.wizard-header').should('not.exist');
//   cy.get('input[name="original_move_date"]').type('6/3/2100').blur();
//   // test an invalid pickup zip code
//   cy.get('input[name="pickup_postal_code"]').clear().type('00000').blur();
//   cy.get('#pickup_postal_code-error').should('exist');

//   cy.get('input[name="pickup_postal_code"]').clear().type('80913');

//   // test an invalid destination zip code
//   cy.get('input[name="destination_postal_code"]').clear().type('00000').blur();
//   cy.get('#destination_postal_code-error').should('exist');

//   cy.get('input[name="destination_postal_code"]').clear().type('30813');
//   cy.nextPage();

//   cy.location().should((loc) => {
//     expect(loc.pathname).to.match(/^\/moves\/[^/]+\/ppm-start/);
//   });
// }

// function SMSeesMoveDetails() {
//   cy.get('.sidebar button').contains('Edit Move').click();

//   cy.get('[data-testid="ppm-summary"]').should((ppmContainer) => {
//     expect(ppmContainer).to.have.length(1);
//   });
// }

// function SMContinueRequestPayment() {
//   serviceMemberStartsPPMPaymentRequest();
//   serviceMemberSubmitsWeightTicket('CAR', true);

//   cy.get('button').contains('Finish Later').click();

//   cy.get('button').contains('OK').click();

//   cy.location().should((loc) => {
//     expect(loc.pathname).to.match(/^\/ppm$/);
//   });

//   cy.get('a').contains('Continue Requesting Payment').click();
//   cy.location().should((loc) => {
//     expect(loc.pathname).to.match(/^\/moves\/[^/]+\/ppm-payment-review/);
//   });
// }

// function serviceMemberStartsPPMPaymentRequest() {
//   cy.contains('Request Payment').click();
//   cy.get('input[name="actual_move_date"]').type('6/20/2018{enter}').blur();
//   cy.get('button').contains('Get Started').click();
// }

// function serviceMemberSubmitsWeightTicket(vehicleType, hasAnother = true, ordinal = null) {
//   if (ordinal) {
//     cy.contains(`Weight Tickets - ${ordinal} set`);
//     if (ordinal === '1st') {
//       cy.get('[data-testid=documents-uploaded]').should('not.exist');
//     } else {
//       cy.get('[data-testid=documents-uploaded]').should('exist');
//     }
//   }

//   cy.get('select[name="weight_ticket_set_type"]').select(vehicleType);

//   if (vehicleType === 'BOX_TRUCK' || vehicleType === 'PRO_GEAR') {
//     cy.get('input[name="vehicle_nickname"]').type('Nickname');
//   } else if (vehicleType === 'CAR' || vehicleType === 'CAR_TRAILER') {
//     cy.get('input[name="vehicle_make"]').type('Make');
//     cy.get('input[name="vehicle_model"]').type('Model');
//   }

//   cy.get('input[name="empty_weight"]').type('1000');

//   cy.upload_file('[data-testid=empty-weight-upload] .filepond--root', 'sample-orders.png');
//   // cy.wait('@postUploadDocument').its('response.statusCode').should('eq', 201);
//   cy.get('[data-filepond-item-state="processing-complete"]', { timeout: fileUploadTimeout }).should('have.length', 1);

//   cy.get('input[name="full_weight"]').type('5000');
//   cy.upload_file('[data-testid=full-weight-upload] .filepond--root', 'sample-orders.png');
//   // cy.wait('@postUploadDocument').its('response.statusCode').should('eq', 201);
//   cy.get('[data-filepond-item-state="processing-complete"]', { timeout: fileUploadTimeout }).should('have.length', 2);
//   cy.get('input[name="weight_ticket_date"]').type('6/2/2018{enter}').blur();
//   cy.get('input[name="additional_weight_ticket"][value="Yes"]').should('not.be.checked');
//   cy.get('input[name="additional_weight_ticket"][value="No"]').should('be.checked');
//   if (hasAnother) {
//     cy.get('input[name="additional_weight_ticket"][value="Yes"]+label').click();
//     cy.get('input[name="additional_weight_ticket"][value="Yes"]').should('be.checked');
//     cy.get('button').contains('Save & Add Another').click();
//     cy.wait('@postWeightTicket').its('response.statusCode').should('eq', 200);
//     cy.get('[data-testid=documents-uploaded]').should('exist');
//   } else {
//     cy.get('button').contains('Save & Continue').click();
//     cy.wait('@postWeightTicket').its('response.statusCode').should('eq', 200);
//   }
// }
