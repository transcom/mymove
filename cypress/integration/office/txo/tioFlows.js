import { TIOOfficeUserType } from '../../../support/constants';

const completeServiceItemCard = ($serviceItem, approve = false) => {
  // serviceItemAmount
  if (!approve) {
    const inputEl = $serviceItem.find('input[data-testid="rejectRadio"]');
    const id = inputEl.attr('id');
    cy.get(`label[for="${id}"]`).click();

    cy.get('textarea[data-testid="textarea"]').type('This is not a valid request');

    cy.wrap($serviceItem).contains('Save').click();
  } else {
    const inputEl = $serviceItem.find('input[data-testid="approveRadio"]');
    const id = inputEl.attr('id');
    cy.get(`label[for="${id}"]`).click();
  }
};

const approveCurrentServiceItem = () => {
  // Approve the service item
  cy.get('[data-testid="ServiceItemCard"]').each((el) => {
    completeServiceItemCard(el, true);
  });

  // Go to next service item
  cy.get('[data-testid=nextServiceItem]').click();
};

const validateDLCalcValues = () => {
  cy.get('[data-testid="ServiceItemCalculations"]')
    .children()
    .should('contain', '14 cwt')
    .and('contain', '354')
    .and('contain', 'ZIP 803 to ZIP 805')
    .and('contain', '21')
    .and('contain', 'Domestic non-peak')
    .and('contain', 'Origin service area: 144')
    .and('contain', '1.01000')
    .and('contain', '$800.00');
};
const validateFSCalcValues = () => {
  cy.get('[data-testid="ServiceItemCalculations"]')
    .children()
    .should('contain', '14 cwt')
    .and('contain', '354')
    .and('contain', 'ZIP 803 to ZIP 805')
    .and('contain', '0.15')
    .and('contain', 'EIA diesel: $2.81')
    .and('contain', 'Weight-based distance multiplier: 0.0004170')
    .and('contain', '$107.00');
};
const validateDUCalcValues = () => {
  cy.get('[data-testid="ServiceItemCalculations"]')
    .children()
    .should('contain', '43 cwt')
    .and('contain', '5.79')
    .and('contain', 'Destination service schedule: 1')
    .and('contain', 'Domestic non-peak')
    .and('contain', '1.04071')
    .and('contain', 'Base year: DUPK Test Year')
    .and('contain', '$459.00');
};
const validateDXCalcValues = () => {
  cy.get('[data-testid="ServiceItemCalculations"]')
    .children()
    .should('contain', '43 cwt')
    .and('contain', '6.25')
    .and('contain', 'service area: 144')
    .and('contain', 'Domestic non-peak')
    .and('contain', '1.04071')
    .and('contain', '$150.00');
};

describe('TIO user', () => {
  before(() => {
    cy.prepareOfficeApp();
  });

  beforeEach(() => {
    cy.intercept('**/ghc/v1/swagger.yaml').as('getGHCClient');
    cy.intercept('**/ghc/v1/queues/payment-requests?**').as('getPaymentRequests');
    cy.intercept('**/ghc/v1/queues/payment-requests?sort=age&order=desc&page=1&perPage=20').as(
      'getSortedPaymentRequests',
    );
    cy.intercept('**/ghc/v1/moves/**/payment-requests').as('getMovePaymentRequests');
    cy.intercept('**/ghc/v1/payment-requests/**').as('getPaymentRequest');
    cy.intercept('**/ghc/v1/move/**').as('getMoves');
    cy.intercept('**/ghc/v1/orders/**').as('getOrders');
    cy.intercept('**/ghc/v1/documents/**').as('getDocuments');
    cy.intercept('**/ghc/v1/move_task_orders/**/mto_shipments').as('getMtoShipments');

    cy.intercept('PATCH', '**/ghc/v1/move-task-orders/**/payment-service-items/**/status').as(
      'patchPaymentServiceItemStatus',
    );
    cy.intercept('PATCH', '**/ghc/v1/payment-requests/**/status').as('patchPaymentRequestStatus');
    cy.intercept('PATCH', '**/ghc/v1/move_task_orders/**/mto_shipments/**').as('patchMtoShipment');
    cy.intercept('**/ghc/v1/moves/**/financial-review-flag').as('financialReviewFlagCompleted');

    const userId = '3b2cc1b0-31a2-4d1b-874f-0591f9127374';
    cy.apiSignInAsUser(userId, TIOOfficeUserType);
  });

  // This test performs a mutation so it can only succeed on a fresh DB.
  it('can use a payment request page to update orders and review a payment request', () => {
    const paymentRequestId = 'ea945ab7-099a-4819-82de-6968efe131dc';

    // TIO Payment Requests queue
    cy.wait(['@getGHCClient', '@getPaymentRequests', '@getSortedPaymentRequests']);
    cy.get('#locator').type('TIOFLO');
    cy.get('th[data-testid="locator"]').first().click();
    cy.get('[data-testid="locator-0"]').click();

    // Payment Requests page
    cy.url().should('include', `/payment-requests`);
    cy.wait(['@getMovePaymentRequests']);
    cy.get('[data-testid="MovePaymentRequests"]');

    // View Orders page
    cy.contains('View orders').click();

    cy.wait(['@getMoves', '@getOrders', '@getDocuments']);
    cy.get('form').within(($form) => {
      cy.get('input[name="tac"]').click().clear().type('E15A');
      cy.get('input[name="sac"]').click().clear().type('4K988AS098F');
      // Edit orders page | Save
      cy.get('button').contains('Save').click();
    });

    cy.url().should('include', `/details`);

    cy.contains('Payment requests').click();

    // Payment Requests page
    cy.url().should('include', `/payment-requests`);
    cy.contains('Review service items').click();

    // Payment Request detail page
    cy.url().should('include', `/payment-requests/${paymentRequestId}`);
    cy.wait(['@getPaymentRequest']);
    cy.get('[data-testid="ReviewServiceItems"]');

    // Approve the first service item
    cy.get('[data-testid="ServiceItemCard"]').each((el) => {
      completeServiceItemCard(el, true);
    });
    cy.wait('@patchPaymentServiceItemStatus');
    cy.contains('Next').click();

    // Approve the second service item
    cy.get('[data-testid="ServiceItemCard"]').each((el) => {
      completeServiceItemCard(el, true);
    });
    cy.wait('@patchPaymentServiceItemStatus');
    cy.contains('Next').click();

    // Approve the shuttling service item

    // Confirm TIO can view the calculations
    cy.contains('Show calculations').click();
    cy.get('[data-testid="ServiceItemCalculations"]').contains('Calculations');
    cy.get('[data-testid="ServiceItemCalculations"]').contains('Total amount requested');
    cy.get('[data-testid="ServiceItemCalculations"]').contains('Service schedule: 2');

    // Confirm TIO can hide the calculations. This ensures there's no scrolling weirdness before the next action
    cy.contains('Hide calculations').click();

    cy.get('[data-testid="ServiceItemCard"]').each((el) => {
      completeServiceItemCard(el, true);
    });
    cy.wait('@patchPaymentServiceItemStatus');
    cy.contains('Next').click();

    // Approve the second service item
    cy.get('[data-testid="ServiceItemCard"]').each((el) => {
      completeServiceItemCard(el, true);
    });
    cy.wait('@patchPaymentServiceItemStatus');
    cy.contains('Next').click();

    // Approve the crating service item

    // Confirm TIO can view the calculations
    cy.contains('Show calculations').click();
    cy.get('[data-testid="ServiceItemCalculations"]').contains('Calculations');
    cy.get('[data-testid="ServiceItemCalculations"]').contains('Total amount requested');
    cy.get('[data-testid="ServiceItemCalculations"]').contains('Dimensions: 12x3x10 in');

    // Confirm TIO can hide the calculations. This ensures there's no scrolling weirdness before the next action
    cy.contains('Hide calculations').click();

    cy.get('[data-testid="ServiceItemCard"]').each((el) => {
      completeServiceItemCard(el, true);
    });
    cy.wait('@patchPaymentServiceItemStatus');
    cy.contains('Next').click();

    // Reject the last
    cy.get('[data-testid="ServiceItemCard"]').each((el) => {
      completeServiceItemCard(el, false);
    });
    cy.wait('@patchPaymentServiceItemStatus');
    cy.contains('Next').click();

    // Complete Request
    cy.contains('Complete request');

    cy.get('[data-testid="requested"]').contains('$1,130.21');
    cy.get('[data-testid="accepted"]').contains('$130.22');
    cy.get('[data-testid="rejected"]').contains('$999.99');

    cy.contains('Authorize payment').click();
    cy.wait('@patchPaymentRequestStatus');

    // Returns to payment requests overview for move
    cy.url().should('include', `/payment-requests`);
    cy.wait(['@getMovePaymentRequests']);
    cy.get('[data-testid="MovePaymentRequests"]');
    cy.get('[data-testid="MovePaymentRequests"] [data-testid="tag"]').contains('Reviewed');
    cy.contains('Review Service Items').should('not.exist');

    // Go back to queue
    cy.get('a[title="Home"]').click();

    cy.contains('Payment requests', { matchCase: false });
    cy.contains('Reviewed', { matchCase: false });
    cy.get('[data-uuid="' + paymentRequestId + '"]').within(() => {
      cy.get('td').eq(2).contains('Reviewed');
    });
  });

  // This test performs a mutation so it can only succeed on a fresh DB.
  // This is a stripped down version of the above to build tests on for MB-7936
  it('can review a payment request', () => {
    const paymentRequestId = 'ea945ab7-099a-4819-82de-6968efe131dc';

    // TIO Payment Requests queue
    cy.wait(['@getGHCClient', '@getPaymentRequests', '@getSortedPaymentRequests']);
    cy.get('[data-uuid="' + paymentRequestId + '"]').click();

    // Payment Requests page
    cy.wait(['@getMovePaymentRequests']);
    cy.contains('Review service items').click();

    // Payment Request detail page
    cy.wait(['@getPaymentRequest']);

    // Approve the first service item
    cy.get('[data-testid="ServiceItemCard"]').each((el) => {
      completeServiceItemCard(el, true);
    });
    cy.wait('@patchPaymentServiceItemStatus');
    cy.contains('Next').click();

    // Reject the second
    cy.get('[data-testid="ServiceItemCard"]').each((el) => {
      completeServiceItemCard(el, false);
    });
    cy.wait('@patchPaymentServiceItemStatus');
    cy.contains('Next').click();

    // Complete Request
    cy.contains('Complete request');

    cy.contains('Authorize payment').click();
    cy.wait('@patchPaymentRequestStatus');

    // Returns to payment requests overview for move
    cy.wait(['@getMovePaymentRequests']);
    cy.get('[data-testid="MovePaymentRequests"]');

    // Go back to queue
    cy.get('a[title="Home"]').click();

    cy.contains('Payment requests', { matchCase: false });
  });

  // This test performs a mutation so it can only succeed on a fresh DB.
  it('can flag the move for review', () => {
    const paymentRequestId = 'ea945ab7-099a-4819-82de-6968efe131dc';

    // TIO Payment Requests queue
    cy.wait(['@getGHCClient', '@getPaymentRequests', '@getSortedPaymentRequests']);
    cy.get('#locator').type('TIOFLO');
    cy.get('th[data-testid="locator"]').first().click();
    cy.get('[data-testid="locator-0"]').click();

    // Payment Requests page
    cy.url().should('include', `/payment-requests`);
    cy.wait(['@getMovePaymentRequests']);
    cy.get('[data-testid="MovePaymentRequests"]');

    // click to trigger financial review modal
    cy.contains('Flag move for financial review').click();

    // Enter information in modal and submit
    cy.get('label').contains('Yes').click();
    cy.get('textarea').type('Something is rotten in the state of Denmark');

    // Click save on the modal
    cy.get('button').contains('Save').click();

    // Verify sucess alert and tag
    cy.contains('Move flagged for financial review.');
    cy.contains('Flagged for financial review');
  });

  // This test performs a mutation so it can only succeed on a fresh DB.
  it('can unflag the move for review', () => {
    const paymentRequestId = 'ea945ab7-099a-4819-82de-6968efe131dc';

    // TIO Payment Requests queue
    cy.wait(['@getGHCClient', '@getPaymentRequests', '@getSortedPaymentRequests']);
    cy.get('#locator').type('TIOFLO');
    cy.get('th[data-testid="locator"]').first().click();
    cy.get('[data-testid="locator-0"]').click();

    // Payment Requests page
    cy.url().should('include', `/payment-requests`);
    cy.wait(['@getMovePaymentRequests']);
    cy.get('[data-testid="MovePaymentRequests"]');

    // click to trigger financial review modal
    cy.contains('Edit').click();

    // Enter information in modal and submit
    cy.get('label').contains('No').click();

    // Click save on the modal
    cy.get('button').contains('Save').click();

    // Verify sucess alert and tag
    cy.contains('Move unflagged for financial review.');
  });

  it('can review a NTS-R', () => {
    cy.wait(['@getGHCClient', '@getPaymentRequests']);

    // Go to known NTS-R move
    cy.get('#locator').type('NTSRT2');
    cy.get('th[data-testid="locator"]').first().click();
    cy.get('[data-testid="locator-0"]').click();

    // Verify we are on the Payment Requests page
    cy.url().should('include', `/payment-requests`);
    cy.wait(['@getMovePaymentRequests', '@getMoves', '@getOrders']);
    cy.get('[data-testid="MovePaymentRequests"]');

    // Verify weight info
    const weightSection = '#billable-weights';
    cy.get(weightSection).contains('Billable weights');
    cy.get(weightSection).contains('8,000 lbs');
    cy.get(weightSection).contains('2,000 lbs');

    // Verify External Shipment shown
    cy.get(weightSection).contains('1 other shipment:');
    cy.get(weightSection).contains('0 lbs');
    cy.get(weightSection).contains('View move details').should('have.attr', 'href');

    // Verify relevant payment request info
    const prSection = '#payment-requests';
    cy.get(prSection).contains('Needs review');
    cy.get(prSection).contains('Reviewed').should('not.exist');
    cy.get(prSection).contains('Rejected').should('not.exist');

    cy.get(prSection).contains('$324.00');
    cy.get(prSection).contains('HTC111-11-1-1111');
    cy.get(prSection).contains('Non-temp storage release');
    cy.get('[data-testid="pickup-to-destination').should('exist');
    cy.get(prSection).contains('1111 (HHG)');

    // Verify Service Item
    cy.get('[data-testid="serviceItemName"]').contains('Counseling');
    cy.get('[data-testid="serviceItemAmount"]').contains('$324.00');

    // Review Weights
    cy.get(weightSection).contains('Review weights').click();
    cy.contains('Edit max billable weight');
    cy.get('[data-testid="closeSidebar"]').click();
    cy.wait(['@getMovePaymentRequests', '@getMoves', '@getOrders']);

    // Review service items
    cy.contains('Review service items').click();
    cy.wait(['@getPaymentRequest']);

    // Approve the service item
    approveCurrentServiceItem();

    // Complete Request
    cy.contains('Complete request');

    cy.contains('Authorize payment').click();
    cy.wait(['@getMovePaymentRequests', '@getMoves', '@getOrders']);
    cy.get('[data-testid="MovePaymentRequests"]');
    cy.contains('Accepted');

    // Should now have 'Reviewed' Tag
    cy.get(prSection).contains('Reviewed');
    cy.get(prSection).contains('Needs Review').should('not.exist');
    cy.get(prSection).contains('Rejected').should('not.exist');

    // Go back home
    cy.get('a[title="Home"]').click();
    cy.contains('Payment requests', { matchCase: false });
  });

  it('can reject a NTS-R', () => {
    cy.wait(['@getGHCClient', '@getPaymentRequests']);

    // Go to known NTS-R move
    cy.get('#locator').type('NTSRT3');
    cy.get('th[data-testid="locator"]').first().click();
    cy.get('[data-testid="locator-0"]').click();

    // Verify we are on the Payment Requests page
    cy.url().should('include', `/payment-requests`);
    cy.wait(['@getMovePaymentRequests', '@getMoves', '@getOrders']);
    cy.get('[data-testid="MovePaymentRequests"]');

    // Verify payment request status
    const prSection = '#payment-requests';
    cy.get(prSection).contains('Needs review');
    cy.get(prSection).contains('Reviewed').should('not.exist');
    cy.get(prSection).contains('Rejected').should('not.exist');

    // Review Weights
    cy.get('#billable-weights').contains('Review weights').click();
    cy.get('[data-testid="closeSidebar"]').click();
    cy.wait(['@getMovePaymentRequests', '@getMoves', '@getOrders']);

    // Start reviewing service items
    cy.contains('Review service items').click();
    cy.wait(['@getPaymentRequest']);

    // Reject the service item
    cy.get('[data-testid="ServiceItemCard"]').each((el) => {
      completeServiceItemCard(el, false);
    });
    cy.wait('@patchPaymentServiceItemStatus');
    cy.get('[data-testid=nextServiceItem]').click();

    // Reject the Request
    cy.contains('Review details');
    cy.contains('Reject request').click();
    cy.wait(['@getMovePaymentRequests', '@getMoves', '@getOrders']);
    cy.get('[data-testid="MovePaymentRequests"]');

    // Should now have 'Rejected' Tag
    cy.get(prSection).contains('Rejected');
    cy.get(prSection).contains('Needs Review').should('not.exist');
    cy.get(prSection).contains('Reviewed').should('not.exist');

    // Go back home
    cy.get('a[title="Home"]').click();
    cy.contains('Payment requests', { matchCase: false });
  });

  it('can view calculation factors', () => {
    cy.wait(['@getGHCClient', '@getPaymentRequests']);

    // Go to known NTS-R move as TIO
    cy.get('#locator').type('NTSRT1');
    cy.get('th[data-testid="locator"]').first().click();
    cy.get('[data-testid="locator-0"]').click();

    // Verify we are on the Payment Requests page
    cy.url().should('include', `/payment-requests`);
    cy.wait(['@getMovePaymentRequests', '@getMoves', '@getOrders']);
    cy.get('[data-testid="MovePaymentRequests"]');

    // Review Weights
    cy.get('#billable-weights').contains('Review weights').click();
    cy.get('[data-testid="closeSidebar"]').click();
    cy.wait(['@getMovePaymentRequests', '@getMoves', '@getOrders']);

    // Review service items
    cy.contains('Review service items').click();
    cy.wait(['@getPaymentRequest', '@getMoves', '@getOrders']);

    // Verify at domestic linehaul calculations
    cy.get('[data-testid=toggleCalculations]').click();
    validateDLCalcValues();
    approveCurrentServiceItem();

    // Verify at fuel surcharge calculations
    cy.get('[data-testid=toggleCalculations]').click();
    validateFSCalcValues();
    approveCurrentServiceItem();

    // Verify at domestic origin calculations
    cy.get('[data-testid=toggleCalculations]').click();
    validateDXCalcValues();
    approveCurrentServiceItem();

    // Verify at domestic destination calculations
    cy.get('[data-testid=toggleCalculations]').click();
    validateDXCalcValues();
    approveCurrentServiceItem();

    // Verify at domestic unpacking calculations
    cy.get('[data-testid=toggleCalculations]').click();
    validateDUCalcValues();
    approveCurrentServiceItem();

    // Complete Request
    cy.contains('Complete request');
    cy.contains('Authorize payment').click();

    // Return to payment requests overview for move
    cy.wait(['@getMovePaymentRequests', '@getMoves', '@getOrders']);
    cy.get('[data-testid="MovePaymentRequests"]');

    // Verify Service Item Calcs now that approved and visable
    cy.contains('Accepted');

    // Verify domestic linehaul cacls on payment request
    cy.contains('Domestic linehaul').click();
    validateDLCalcValues();
    cy.contains('Domestic linehaul').click();

    // Verify fuel surcharge calcs on payment request
    cy.contains('Fuel surcharge').click();
    validateFSCalcValues();
    cy.contains('Fuel surcharge').click();

    // Verify Domestic origin price cacls on payment request
    cy.contains('Domestic origin price').click();
    validateDXCalcValues();
    cy.contains('Domestic origin price').click();

    // Verify domestic destination price cacls on payment request
    cy.contains('Domestic destination price').click();
    validateDXCalcValues();
    cy.contains('Domestic destination price').click();

    // Verify domestic unpacking cacls on payment request
    cy.contains('Domestic unpacking').click();
    validateDUCalcValues();
    cy.contains('Domestic unpacking').click();

    // calcs are good, go home
    cy.get('a[title="Home"]').click();
    cy.contains('Payment requests', { matchCase: false });
  });

  it('can add/edit TAC/SAC', () => {
    // TIO Payment Requests queue
    cy.wait(['@getGHCClient', '@getPaymentRequests', '@getSortedPaymentRequests']);
    cy.get('#locator').type('NTSTIO');
    cy.get('th[data-testid="locator"]').first().click();
    cy.get('[data-testid="locator-0"]').click();

    cy.wait(['@getMovePaymentRequests', '@getMoves', '@getOrders']);
    cy.get('[data-testid="MovePaymentRequests"]');

    cy.get('button').contains('Edit').click();
    cy.get('button').contains('Add or edit codes').click();
    cy.url().should('include', `/moves/NTSTIO/orders`);

    cy.get('form').within(() => {
      cy.get('input[data-testid="ntsTacInput"]').click().clear().type('E19A');
      cy.get('input[data-testid="ntsSacInput"]').click().clear().type('3L988AS098F');
      // Edit orders page | Save
      cy.get('button').contains('Save').click();
    });
    cy.url().should('include', `/moves/NTSTIO/details`);
    cy.contains('Payment requests').click();
    cy.url().should('include', `/payment-requests`);
    cy.get('button').contains('Edit').click();

    cy.get('input#tacType-NTS').click({ force: true });
    cy.get('input#sacType-NTS').click({ force: true });
    cy.get('button[type="submit"]').click();
    cy.wait(['@patchMtoShipment', '@getMtoShipments']);

    cy.get('[data-testid="tac"]').contains('E19A (NTS)');
    cy.get('[data-testid="sac"]').contains('3L988AS098F (NTS)');
  });

  it('can view and approve service items', () => {
    // TIO Payment Requests queue
    cy.wait(['@getGHCClient', '@getPaymentRequests', '@getSortedPaymentRequests']);
    cy.get('#locator').type('NTSTIO');
    cy.get('th[data-testid="locator"]').first().click();
    cy.get('[data-testid="locator-0"]').click();

    cy.wait(['@getMovePaymentRequests', '@getMoves', '@getOrders']);
    cy.get('[data-testid="MovePaymentRequests"]');

    cy.get('[test-dataid="reviewBtn"]').click();

    cy.get('[data-testid="serviceItemName"]').contains('Move management');
    cy.get('[data-testid="approveRadio"]').click({ force: true });
    cy.wait('@patchPaymentServiceItemStatus');
    cy.get('button').contains('Next').click();

    cy.get('[data-testid="serviceItemName"]').contains('Domestic origin shuttle service');
    cy.get('[data-testid="approveRadio"]').click({ force: true });
    cy.wait('@patchPaymentServiceItemStatus');
    cy.get('button').contains('Next').click();

    cy.get('[data-testid="serviceItemName"]').contains('Domestic origin shuttle service');
    cy.get('[data-testid="approveRadio"]').click({ force: true });
    cy.wait('@patchPaymentServiceItemStatus');
    cy.get('button').contains('Next').click();

    cy.get('[data-testid="serviceItemName"]').contains('Domestic crating');
    cy.get('[data-testid="approveRadio"]').click({ force: true });
    cy.wait('@patchPaymentServiceItemStatus');
    cy.get('button').contains('Next').click();

    cy.get('[data-testid="serviceItemName"]').contains('Domestic crating');
    cy.get('[data-testid="approveRadio"]').click({ force: true });
    cy.wait('@patchPaymentServiceItemStatus');
    cy.get('button').contains('Next').click();

    cy.get('[data-testid="serviceItemName"]').contains('Domestic linehaul');
    cy.get('[data-testid="approveRadio"]').click({ force: true });
    cy.wait('@patchPaymentServiceItemStatus');
    cy.get('button').contains('Next').click();

    cy.get('[data-testid="accepted"]').contains('$1,130.21');
    cy.get('button').contains('Authorize payment').click();
    cy.wait(['@getMovePaymentRequests']);

    cy.get('[data-testid="tag"]').contains('Reviewed');
  });
});
