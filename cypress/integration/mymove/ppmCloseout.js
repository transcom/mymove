// import { fileUploadTimeout } from '../../support/constants';

// describe('allows a SM to request a payment', function () {
//   before(() => {
//     cy.prepareCustomerApp();
//   });

//   beforeEach(() => {
//     cy.logout();
//   });

//   beforeEach(() => {
//     // TODO - commenting this out for now because cy.intercept has an open bug related to file upload endpoints
//     // https://github.com/cypress-io/cypress/issues/9534
//     // cy.intercept('POST', '**/internal/uploads').as('postUploadDocument');

//     cy.intercept('POST', '**/moves/**/weight_ticket').as('postWeightTicket');
//     cy.intercept('POST', '**/moves/**/moving_expense_documents').as('postMovingExpense');
//     cy.intercept('POST', '**/internal/personally_procured_move/**/request_payment').as('requestPayment');
//     cy.intercept('POST', '**/moves/**/signed_certifications').as('signedCertifications');
//   });

//   const moveID = 'f9f10492-587e-43b3-af2a-9f67d2ac8757';

//   it('service member goes through entire request payment flow', () => {
//     cy.apiSignInAsPpmUser('8e0d7e98-134e-4b28-bdd1-7d6b1ff34f9e');
//     serviceMemberStartsPPMPaymentRequest();
//     serviceMemberSubmitsWeightTicket('CAR', true);
//     serviceMemberChecksNumberOfWeightTickets('2nd');
//     serviceMemberSubmitsWeightTicket('BOX_TRUCK', false);
//     serviceMemberViewsExpensesLandingPage();
//     serviceMemberUploadsExpenses(false);
//     serviceMemberReviewsDocuments();
//     serviceMemberEditsPaymentRequest();
//     serviceMemberAddsWeightTicketSetWithMissingDocuments();
//     serviceMemberSubmitsPaymentRequestWithMissingDocuments();
//   });

//   it('service member reads introduction to ppm payment and goes back to homepage', () => {
//     cy.apiSignInAsPpmUser('745e0eba-4028-4c78-a262-818b00802748');
//     serviceMemberStartsPPMPaymentRequest();
//   });

//   it('service member can save a weight ticket for later', () => {
//     cy.apiSignInAsUser('745e0eba-4028-4c78-a262-818b00802748');

//     cy.visit(`/moves/${moveID}/ppm-weight-ticket`);
//     serviceMemberCanFinishWeightTicketLater('BOX_TRUCK');
//   });

//   it('service member submits weight tickets without any documents', () => {
//     cy.apiSignInAsUser('745e0eba-4028-4c78-a262-818b00802748');
//     cy.visit(`/moves/${moveID}/ppm-weight-ticket`);
//     serviceMemberSubmitsWeightsTicketsWithoutReceipts();
//   });

//   it('service member requests a car + trailer weight ticket payment', () => {
//     cy.apiSignInAsUser('745e0eba-4028-4c78-a262-818b00802748');
//     cy.visit(`/moves/${moveID}/ppm-weight-ticket`);
//     serviceMemberSubmitsCarTrailerWeightTicket();
//   });

//   it('service member starting at review page returns to review page after adding a weight ticket', () => {
//     cy.apiSignInAsUser('745e0eba-4028-4c78-a262-818b00802748');
//     cy.visit(`/moves/${moveID}/ppm-payment-review`);
//     cy.location().should((loc) => {
//       expect(loc.pathname).to.match(/^\/moves\/[^/]+\/ppm-payment-review/);
//     });
//     cy.get('[data-testid="weight-ticket-link"]').click();
//     cy.location().should((loc) => {
//       expect(loc.pathname).to.match(/^\/moves\/[^/]+\/ppm-weight-ticket/);
//     });
//     serviceMemberSubmitsWeightTicket('CAR', false);
//     cy.location().should((loc) => {
//       expect(loc.pathname).to.match(/^\/moves\/[^/]+\/ppm-payment-review/);
//     });
//   });

//   it('service member starting at review page returns to review page after adding an expense', () => {
//     cy.apiSignInAsUser('745e0eba-4028-4c78-a262-818b00802748');
//     cy.visit(`/moves/${moveID}/ppm-payment-review`);
//     cy.location().should((loc) => {
//       expect(loc.pathname).to.match(/^\/moves\/[^/]+\/ppm-payment-review/);
//     });
//     cy.get('[data-testid="expense-link"]').click();
//     cy.location().should((loc) => {
//       expect(loc.pathname).to.match(/^\/moves\/[^/]+\/ppm-expenses/);
//     });
//     serviceMemberUploadsExpenses(false);
//     cy.location().should((loc) => {
//       expect(loc.pathname).to.match(/^\/moves\/[^/]+\/ppm-payment-review/);
//     });
//   });

//   it('service member can skip weight tickets and expenses if already have one', () => {
//     cy.apiSignInAsUser('745e0eba-4028-4c78-a262-818b00802748');
//     cy.visit(`/moves/${moveID}/ppm-weight-ticket`);
//     serviceMemberSubmitsWeightTicket('CAR', true);
//     serviceMemberSkipsStep();
//     serviceMemberViewsExpensesLandingPage();
//     serviceMemberUploadsExpenses();
//     serviceMemberSkipsStep();
//   });

//   it('service member with old weight tickets can see and delete them', () => {
//     cy.apiSignInAsPpmUser('beccca28-6e15-40cc-8692-261cae0d4b14');
//     cy.get('[data-testid="edit-payment-request"]').contains('Edit Payment Request').should('exist').click();
//     cy.get('.ticket-item').first().should('not.contain', 'set');
//     cy.get('[data-testid="delete-ticket"]').first().click();
//     cy.get('[data-testid="delete-confirmation-button"]').click();
//     cy.get('.ticket-item').should('not.exist');
//   });
// });

// function serviceMemberSkipsStep() {
//   cy.get('[data-testid=skip]').contains('Skip').click();
// }

// function serviceMemberSubmitsPaymentRequestWithMissingDocuments() {
//   cy.get('.review-customer-agreement a').contains('Legal Agreement').click();
//   cy.location().should((loc) => {
//     expect(loc.pathname).to.match(/^\/ppm-customer-agreement/);
//   });
//   cy.get('.usa-button').contains('Back').click();
//   cy.location().should((loc) => {
//     expect(loc.pathname).to.match(/^\/moves\/[^/]+\/ppm-payment-review/);
//   });
//   cy.get('.missing-label').contains('Your estimated payment is unknown');

//   cy.get('input[id="agree-checkbox"]').check({ force: true });

//   cy.get('button').contains('Submit Request').should('be.enabled').click();
//   cy.wait('@signedCertifications');
//   cy.wait('@requestPayment');

//   cy.get('.usa-alert--warning').contains('Payment request is missing info');
//   cy.get('.usa-alert--warning').contains(
//     'You will need to contact your local PPPO office to resolve your missing weight ticket.',
//   );

//   cy.get('.title').contains('Next step: Contact the PPPO office');
//   //TODO add this back in when we have BVS scores
//   // cy.get('.missing-label').contains('Unknown');
// }

// function serviceMemberReviewsDocuments() {
//   cy.location().should((loc) => {
//     expect(loc.pathname).to.match(/^\/moves\/[^/]+\/ppm-payment-review/);
//   });
//   cy.get('.review-customer-agreement a').contains('Legal Agreement').click();
//   cy.location().should((loc) => {
//     expect(loc.pathname).to.match(/^\/ppm-customer-agreement/);
//   });
//   cy.get('[data-testid="back-button"]').contains('Back').click();
//   cy.location().should((loc) => {
//     expect(loc.pathname).to.match(/^\/moves\/[^/]+\/ppm-payment-review/);
//   });
//   cy.get('input[id="agree-checkbox"]').check({ force: true });
//   cy.contains(`You're requesting a payment of $`);
//   cy.get('button').contains('Submit Request').should('be.enabled').click();
//   cy.wait('@signedCertifications');
//   cy.wait('@requestPayment');
// }
// function serviceMemberDeletesDocuments() {
//   cy.location().should((loc) => {
//     expect(loc.pathname).to.match(/^\/moves\/[^/]+\/ppm-payment-review/);
//   });
//   cy.get('.ticket-item').should('have.length', 4);
//   cy.get('[data-testid="delete-ticket"]').first().click();
//   cy.get('[data-testid="delete-confirmation-button"]').click();
//   cy.get('.ticket-item').should('have.length', 3);
// }
// function serviceMemberEditsPaymentRequest() {
//   cy.get('.usa-alert--success').contains('Payment request submitted').should('exist');
//   cy.get('[data-testid="edit-payment-request"]').contains('Edit Payment Request').should('exist').click();
//   cy.get('[data-testid=weight-ticket-link]').should('exist').click();
//   serviceMemberSubmitsWeightTicket('CAR', false);
//   serviceMemberDeletesDocuments();
//   serviceMemberReviewsDocuments();
// }
// function serviceMemberAddsWeightTicketSetWithMissingDocuments(hasAnother = false) {
//   cy.get('[data-testid="edit-payment-request"]').contains('Edit Payment Request').should('exist').click();
//   cy.get('[data-testid=weight-ticket-link]').should('exist').click();

//   cy.get('select[name="weight_ticket_set_type"]').select('BOX_TRUCK');

//   cy.get('input[name="vehicle_nickname"]').type('Nickname');

//   cy.get('input[name="empty_weight"]').type('1000');
//   cy.get('input[name="missingEmptyWeightTicket"]').check({ force: true });

//   cy.get('input[name="full_weight"]').type('5000');
//   cy.upload_file('[data-testid=full-weight-upload] .filepond--root', 'sample-orders.png');
//   // cy.wait('@postUploadDocument').its('response.statusCode').should('eq', 201);
//   cy.get('[data-filepond-item-state="processing-complete"]', { timeout: fileUploadTimeout }).should('have.length', 1);

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

// function serviceMemberViewsExpensesLandingPage() {
//   cy.location().should((loc) => {
//     expect(loc.pathname).to.match(/^\/moves\/[^/]+\/ppm-expenses-intro/);
//   });

//   cy.get('[data-testid=documents-uploaded]').should('exist');
//   cy.get('button').contains('Continue').should('be.disabled');

//   cy.get('[type="radio"]').first().should('be.not.checked');
//   cy.get('[type="radio"]').last().should('be.not.checked');

//   cy.get('a')
//     .contains('More about expenses')
//     .should('have.attr', 'href')
//     .and('match', /^\/allowable-expenses/);

//   cy.get('input[name="hasExpenses"][value="Yes"]').should('not.be.checked');
//   cy.get('input[name="hasExpenses"][value="No"]').should('not.be.checked');
//   cy.get('input[name="hasExpenses"][value="Yes"]+label').click();

//   cy.get('button').contains('Continue').should('be.enabled').click();
// }

// function serviceMemberUploadsExpenses(hasAnother = true, expenseNumber = null) {
//   cy.location().should((loc) => {
//     expect(loc.pathname).to.match(/^\/moves\/[^/]+\/ppm-expenses/);
//   });

//   if (expenseNumber) {
//     cy.contains(`Expense ${expenseNumber}`);
//   }
//   cy.get('[data-testid=documents-uploaded]').should('exist');

//   cy.get('select[name="moving_expense_type"]').select('GAS');
//   cy.get('input[name="title"]').type('title');
//   cy.get('input[name="requested_amount_cents"]').type('1000');

//   cy.upload_file('.filepond--root:first', 'sample-orders.png');
//   // cy.wait('@postUploadDocument').its('response.statusCode').should('eq', 201);
//   cy.get('[data-filepond-item-state="processing-complete"]', { timeout: fileUploadTimeout }).should('have.length', 1);

//   cy.get('input[name="missingReceipt"]').should('not.be.checked');
//   cy.get('input[name="paymentMethod"][value="GTCC"]').should('be.checked');
//   cy.get('input[name="paymentMethod"][value="OTHER"]').should('not.be.checked');
//   cy.get('input[name="haveMoreExpenses"][value="Yes"]').should('not.be.checked');
//   cy.get('input[name="haveMoreExpenses"][value="No"]').should('be.checked');
//   cy.get('input[name="haveMoreExpenses"][value="Yes"]+label').click();
//   if (hasAnother) {
//     cy.get('button').contains('Save & Add Another').click();
//     cy.wait('@postMovingExpense').its('response.statusCode').should('eq', 200);
//     cy.get('[data-testid=documents-uploaded]').should('exist');
//   } else {
//     cy.get('input[name="haveMoreExpenses"][value="No"]+label').click();
//     cy.get('input[name="haveMoreExpenses"][value="No"]').should('be.checked');
//     cy.get('button').contains('Save & Continue').click();
//     cy.wait('@postMovingExpense').its('response.statusCode').should('eq', 200);
//   }
// }

// function serviceMemberSubmitsCarTrailerWeightTicket() {
//   cy.get('select[name="weight_ticket_set_type"]').select('CAR_TRAILER');

//   cy.get('input[name="vehicle_make"]').type('Make');
//   cy.get('input[name="vehicle_model"]').type('Model');

//   cy.contains('Do you own this trailer').children('a').should('have.attr', 'href', '/trailer-criteria');

//   cy.get('input[name="isValidTrailer"][value="Yes"]').should('not.be.checked');
//   cy.get('input[name="isValidTrailer"][value="No"]').should('be.checked');
//   cy.get('input[name="isValidTrailer"][value="Yes"]+label').click();

//   cy.upload_file('[data-testid=trailer-upload] .filepond--root', 'sample-orders.png');
//   // cy.wait('@postUploadDocument').its('response.statusCode').should('eq', 201);
//   cy.get('[data-filepond-item-state="processing-complete"]', { timeout: fileUploadTimeout }).should('have.length', 1);

//   cy.get('input[name="missingDocumentation"]').check({ force: true });

//   cy.get('input[name="empty_weight"]').type('1000');
//   cy.get('input[name="missingEmptyWeightTicket"]').check({ force: true });

//   cy.get('input[name="full_weight"]').type('5000');
//   cy.upload_file('.filepond--root:last', 'sample-orders.png');
//   // cy.wait('@postUploadDocument').its('response.statusCode').should('eq', 201);
//   cy.get('[data-filepond-item-state="processing-complete"]', { timeout: fileUploadTimeout }).should('have.length', 2);
//   cy.get('input[name="weight_ticket_date"]').type('6/2/2018{enter}').blur();

//   cy.get('input[name="additional_weight_ticket"][value="Yes"]').should('not.be.checked');
//   cy.get('input[name="additional_weight_ticket"][value="No"]').should('be.checked');
// }
// function serviceMemberCanFinishWeightTicketLater(vehicleType) {
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

//   cy.get('button').contains('Finish Later').click();

//   cy.get('button').contains('OK').click();

//   cy.location().should((loc) => {
//     expect(loc.pathname).to.match(/^\/ppm$/);
//   });
// }

// function serviceMemberSubmitsWeightsTicketsWithoutReceipts() {
//   cy.get('select[name="weight_ticket_set_type"]').select('CAR_TRAILER');
//   cy.get('input[name="vehicle_make"]').type('Make');
//   cy.get('input[name="vehicle_model"]').type('Model');
//   cy.get('input[name="empty_weight"]').type('1000');
//   cy.get('input[name="full_weight"]').type('2000');
//   cy.upload_file('[data-testid=full-weight-upload] .filepond--root', 'sample-orders.png');
//   // cy.wait('@postUploadDocument').its('response.statusCode').should('eq', 201);
//   cy.get('[data-filepond-item-state="processing-complete"]', { timeout: fileUploadTimeout }).should('have.length', 1);

//   cy.get('input[name="isValidTrailer"][value="Yes"]+label').click();
//   cy.get('input[name="missingDocumentation"]+label').click();
//   cy.get('[data-testid=trailer-warning]').contains(
//     'If your state does not provide a registration or bill of sale for your trailer, you may write and upload a signed and dated statement certifying that you or your spouse own the trailer and meets the trailer criteria. Upload your statement using the proof of ownership field.',
//   );
//   cy.get('input[name="missingDocumentation"]+label').click({ force: false });
//   cy.upload_file('[data-testid=trailer-upload] .filepond--root', 'sample-orders.png');
//   // cy.wait('@postUploadDocument').its('response.statusCode').should('eq', 201);
//   cy.get('[data-filepond-item-state="processing-complete"]', { timeout: fileUploadTimeout }).should('have.length', 2);

//   cy.get('input[name="missingEmptyWeightTicket"]+label').click();
//   cy.get('[data-testid=empty-warning]').contains(
//     'Contact your local Transportation Office (PPPO) to let them know you’re missing this weight ticket. For now, keep going and enter the info you do have.',
//   );

//   cy.get('input[name="weight_ticket_date"]').type('6/2/2018{enter}').blur();

//   cy.get('input[name="additional_weight_ticket"][value="Yes"]+label').click();
//   cy.get('input[name="additional_weight_ticket"][value="Yes"]').should('be.checked');
//   cy.contains('Save & Add Another').should('not.be.disabled').click();
//   cy.wait('@postWeightTicket').its('response.statusCode').should('eq', 200);
// }

// function serviceMemberStartsPPMPaymentRequest() {
//   cy.contains('Request Payment').click();
//   cy.get('input[name="actual_move_date"]').type('6/20/2018{enter}').blur();
//   cy.get('button').contains('Get Started').click();
// }

// function serviceMemberChecksNumberOfWeightTickets(ordinal) {
//   cy.contains(`Weight Tickets - ${ordinal} set`);
//   cy.get('[data-testid=documents-uploaded]').should('exist');
// }

// function serviceMemberSubmitsWeightTicket(vehicleType, hasAnother = true) {
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
