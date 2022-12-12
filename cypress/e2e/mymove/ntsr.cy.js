describe('Customer NTSr Setup flow', function () {
  // profile@co.mple.te
  const profileCompleteUser = '99360a51-8cfa-4e25-ae57-24e66077305f';
  // nts@ntsr.unsubmitted
  const ntsUser = '80da86f3-9dac-4298-8b03-b753b443668e';

  before(() => {
    cy.prepareCustomerApp();
  });

  it('Sets up an NTSr shipment', function () {
    cy.apiSignInAsUser(profileCompleteUser);
    customerCreatesAnNTSRShipment();
    customerReviewsNTSRMoveDetails();
  });

  it('Edits an NTSr shipment from review page', function () {
    cy.apiSignInAsUser(ntsUser);
    customerVisitsReviewPage();
    customerReviewsShipmentInfoOnReviewPage();
    customerEditsNTSRShipmentFromReviewPage();
  });

  it('Edits an NTSr shipment from home page', function () {
    cy.apiSignInAsUser(ntsUser);
    customerEditsNTSRShipmentFromHomePage();
  });
});

function customerEditsNTSRShipmentFromHomePage() {
  cy.get('[data-testid="shipment-list-item-container"]').eq(1).children().contains('Edit').click();
  cy.get('textarea[data-testid="remarks"]').clear().type('Warning: glass').blur();

  cy.get('button').contains('Save').click();
  cy.location().should((loc) => {
    expect(loc.pathname).to.match(/^\/moves\/[^/]+\/review/);
  });
}

function customerReviewsNTSRMoveDetails() {
  cy.get('[data-testid="review-move-header"]').contains('Review your details');

  // Requested delivery date
  cy.get('[data-testid="ntsr-summary"]').last().contains('02 Jan 2020');

  // Destination
  cy.get('[data-testid="ntsr-summary"]').last().contains('412 Avenue M');
  cy.get('[data-testid="ntsr-summary"]').last().contains('91111');

  // Receiving agent
  cy.get('[data-testid="ntsr-summary"]').last().contains('James Bond');
  cy.get('[data-testid="ntsr-summary"]').last().contains('777-777-7777');
  cy.get('[data-testid="ntsr-summary"]').last().contains('007@example.com');

  // Remarks
  cy.get('[data-testid="ntsr-summary"]').last().contains('some other customer remark');
}

function customerEditsNTSRShipmentFromReviewPage() {
  cy.get('button[data-testid="edit-ntsr-shipment-btn"]').contains('Edit').click();
  cy.get('input[name="delivery.requestedDate"]').clear().type('01/01/2022').blur();
  cy.get('input[name="delivery.address.streetAddress1"]').clear().type('123 Maple street');
  cy.get('input[data-testid="has-secondary-delivery"]').check({ force: true });
  cy.get('input[name="secondaryDelivery.address.streetAddress1"]').clear().type('123 Oak Street');
  cy.get('input[name="delivery.agent.firstName"]').clear().type('Ketchum').blur();
  cy.get('textarea[data-testid="remarks"]').clear().type('Warning: fragile').blur();
  cy.get('button').contains('Save').click();
  cy.location().should((loc) => {
    expect(loc.pathname).to.match(/^\/moves\/[^/]+\/review/);
  });
  cy.get('[data-testid="ntsr-summary"]').contains('01 Jan 2022');
  cy.get('[data-testid="ntsr-summary"]').contains('123 Maple street');
  cy.get('[data-testid="ntsr-summary"]').contains('123 Oak Street');
  cy.get('[data-testid="ntsr-summary"]').contains('Ketchum Ash');
  cy.get('[data-testid="ntsr-summary"]').contains('Warning: fragile');
}

function customerReviewsShipmentInfoOnReviewPage() {
  cy.get('[data-testid="review-move-header"]').contains('Review your details');

  // Requested delivery date
  cy.get('[data-testid="ntsr-summary"]').last().contains('15 Mar 2020');

  // Destination
  cy.get('[data-testid="ntsr-summary"]').last().contains('987 Any Avenue P.O. Box 9876');
  cy.get('[data-testid="ntsr-summary"]').last().contains('Fairfield, CA 94535');

  // Second destination
  cy.get('[data-testid="ntsr-summary"]').last().contains('123 Any Street P.O. Box 12345');
  cy.get('[data-testid="ntsr-summary"]').last().contains('Beverly Hills, CA 90210');

  // Receiving agent
  cy.get('[data-testid="ntsr-summary"]').last().contains('Jason Ash');
  cy.get('[data-testid="ntsr-summary"]').last().contains('202-555-9301');
  cy.get('[data-testid="ntsr-summary"]').last().contains('jason.ash@example.com');

  // Remarks
  cy.get('[data-testid="ntsr-summary"]').last().contains('Please treat gently');
}

function customerVisitsReviewPage() {
  cy.get('button[data-testid="review-and-submit-btn"]').click();
  cy.get('[data-testid="review-move-header"]').contains('Review your details');
}

function customerCreatesAnNTSRShipment() {
  cy.get('[data-testid="shipment-selection-btn"]').click();
  cy.nextPage();
  cy.get('input[type="radio"]').eq(3).check({ force: true });
  cy.nextPage();
  cy.location().should((loc) => {
    expect(loc.pathname).to.match(/^\/moves\/[^/]+\/new-shipment/);
  });

  // pickup date
  cy.get('input[name="delivery.requestedDate"]').type('01/02/2020').blur();

  // delivery location
  cy.get(`input[name="delivery.address.streetAddress1"]`).type('412 Avenue M');
  cy.get(`input[name="delivery.address.streetAddress2"]`).type('#3E');
  cy.get(`input[name="delivery.address.city"]`).type('Los Angeles');
  cy.get(`select[name="delivery.address.state"]`).select('CA');
  cy.get(`input[name="delivery.address.postalCode"]`).type('91111').blur();

  // receiving agent
  cy.get(`input[name="delivery.agent.firstName"]`).type('James');
  cy.get(`input[name="delivery.agent.lastName"]`).type('Bond');
  cy.get(`input[name="delivery.agent.phone"]`).type('7777777777');
  cy.get(`input[name="delivery.agent.email"]`).type('007@example.com').blur();

  // remarks
  cy.get(`[data-testid="remarks"]`).first().type('some other customer remark');

  cy.nextPage();
}
