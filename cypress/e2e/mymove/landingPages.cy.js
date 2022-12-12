// TODO - sometimes tests result in a 500 response to this URL;
// http://milmovelocal:3000/internal/estimates/ppm?original_move_date=2020-08-27&origin_zip=72017&origin_duty_station_zip=50309&orders_id=fd5495d4-5105-46b8-89f8-0e6b155d900b&weight_estimate=8000
// This causes the test to fail because of an uncaught exception (the fetch error is not handled)

// describe('testing landing pages', function () {
//   before(() => {
//     cy.prepareCustomerApp();
//   });

//   // Submitted draft move/orders but no move type yet.
//   it('tests pre-move type', function () {
//     // sm_no_move_type@example.com
//     const userId = '9ceb8321-6a82-4f6d-8bb3-a1d85922a202';
//     cy.apiSignInAsPpmUser(userId);
//     cy.contains('Next Step: Finish setting up your move');
//   });

//   // PPM: SUBMITTED
//   it('tests submitted PPM', function () {
//     // ppm@incomple.te
//     const userId = 'e10d5964-c070-49cb-9bd1-eaf9f7348eb6';
//     cy.apiSignInAsPpmUser(userId);
//     cy.contains('Handle your own move (PPM)');
//     cy.contains('Next Step: Wait for approval');
//     cy.should('not.contain', 'Add PPM (DITY) Move');
//   });

//   // PPM: APPROVED
//   it('tests approved PPM', function () {
//     // ppm@approv.ed
//     const userId = '70665111-7bbb-4876-a53d-18bb125c943e';
//     cy.apiSignInAsPpmUser(userId);
//     cy.contains('Handle your own move (PPM)');
//     cy.contains('Next Step: Request payment');
//   });

//   // PPM: PAYMENT_REQUESTED
//   it('tests PPM that has requested payment', function () {
//     // ppmpayment@request.ed
//     const userId = 'beccca28-6e15-40cc-8692-261cae0d4b14';
//     cy.apiSignInAsPpmUser(userId);
//     cy.setFeatureFlag('ppmPaymentRequest=false', '/ppm');
//     cy.contains('Handle your own move (PPM)');
//     cy.contains('Edit Payment Request');
//     cy.contains('Estimated');
//     cy.contains('Submitted');
//   });

//   // PPM: COMPLETED
//   // Not seeing a path to a COMPLETED PPM move at this time.

//   // PPM: CANCELED
//   it('tests canceled PPM', function () {
//     // ppm-canceled@example.com
//     const userId = '20102768-4d45-449c-a585-81bc386204b1';
//     cy.apiSignInAsPpmUser(userId);
//     cy.contains('New move');
//     cy.contains('Start here');
//   });
// });
