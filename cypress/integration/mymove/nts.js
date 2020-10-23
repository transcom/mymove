describe('Customer NTS(r) Setup flow', function () {
  // profile@comple.te
  const profileCompleteUser = '3b9360a3-3304-4c60-90f4-83d687884077';
  // nts@ntsr.unsubmitted
  const ntsUser = '583cfbe1-cb34-4381-9e1f-54f68200da1b';

  before(() => {
    cy.prepareCustomerApp();
  });

  it('Sets up an NTS shipment', function () {
    cy.apiSignInAsUser(profileCompleteUser);
    customerCreatesAnNTSShipment();
    customerReviewsNTSMoveDetails();
  });

  it('Sets up an NTSr shipment', function () {
    cy.apiSignInAsUser(profileCompleteUser);
    customerCreatesAnNTSRShipment();
    customerReviewsNTSRMoveDetails();
  });

  it('Edits an NTS shipment', function () {
    cy.apiSignInAsUser(ntsUser);
    customerVisitsReviewPage();
    customerEditsNTSShipment();
  });

  it('Edits an NTSr shipment', function () {
    cy.apiSignInAsUser(ntsUser);
    customerVisitsReviewPage();
    customerEditsNTSRShipment();
  });
});

function customerReviewsNTSRMoveDetails() {
  cy.get('[data-testid="review-move-header"]').contains('Review your details');

  // Requested delivery date
  cy.get('[data-testid="ntsr-summary"]').last().contains('02 Jan 2020');

  // Destination
  cy.get('[data-testid="ntsr-summary"]').last().contains('30813');

  // Receiving agent
  cy.get('[data-testid="ntsr-summary"]').last().contains('James Bond');
  cy.get('[data-testid="ntsr-summary"]').last().contains('777-777-7777');
  cy.get('[data-testid="ntsr-summary"]').last().contains('007@example.com');

  // Remarks
  cy.get('[data-testid="ntsr-summary"]').last().contains('some other customer remark');
}

function customerReviewsNTSMoveDetails() {
  cy.get('[data-testid="review-move-header"]').contains('Review your details');
  cy.get('[data-testid="nts-summary"]').contains('NTS');

  // Requested pickup date
  cy.get('[data-testid="nts-summary"]').contains('02 Aug 2020');

  // Pickup location
  cy.get('[data-testid="nts-summary"]').contains('123 Any Street P.O. Box 12345');
  cy.get('[data-testid="nts-summary"]').contains('Beverly Hills, CA 90210');

  // Releasing agent
  cy.get('[data-testid="nts-summary"]').contains('John Lee');
  cy.get('[data-testid="nts-summary"]').contains('999-999-9999');
  cy.get('[data-testid="nts-summary"]').contains('john@example.com');

  // Remarks
  cy.get('[data-testid="nts-summary"]').contains('some customer remark');
}

function customerEditsNTSRShipment() {
  cy.get('button[data-testid="edit-ntsr-shipment-btn"]').contains('Edit').click();
  cy.get('input[name="delivery.requestedDate"]').clear().type('01/01/2022').blur();
  cy.get('[data-testid="mailingAddress1"]').clear().type('123 Maple street');
  cy.get('input[data-testid="firstName"]').clear().type('Ketchum').blur();
  cy.get('input[data-testid="remarks"]').clear().type('Warning: fragile').blur();
  cy.get('button').contains('Save').click();
  cy.location().should((loc) => {
    expect(loc.pathname).to.match(/^\/moves\/[^/]+\/review/);
  });
  cy.get('[data-testid="ntsr-summary"]').contains('01 Jan 2022');
  cy.get('[data-testid="ntsr-summary"]').contains('123 Maple street');
  cy.get('[data-testid="ntsr-summary"]').contains('Ketchum Ash');
  cy.get('[data-testid="ntsr-summary"]').contains('Warning: fragile');
}

function customerEditsNTSShipment() {
  cy.get('button[data-testid="edit-nts-shipment-btn"]').contains('Edit').click();
  cy.get('input[name="pickup.requestedDate"]').clear().type('12/25/2020').blur();
  cy.get('input[data-testid="lastName"]').clear().type('Bourne').blur();
  cy.get('input[data-testid="remarks"]').clear().type('Handle with care').blur();
  cy.get('button').contains('Save').click();
  cy.location().should((loc) => {
    expect(loc.pathname).to.match(/^\/moves\/[^/]+\/review/);
  });
  cy.get('[data-testid="nts-summary"]').contains('25 Dec 2020');
  cy.get('[data-testid="nts-summary"]').contains('Jason Bourne');
  cy.get('[data-testid="nts-summary"]').contains('Handle with care');
}

function customerVisitsReviewPage() {
  cy.get('button[data-testid="review-and-submit-btn"]').contains('Review and submit').click();
  cy.get('h2[data-testid="review-move-header"]').contains('Review your details');
}

function customerCreatesAnNTSShipment() {
  cy.get('[data-testid="shipment-selection-btn"]').contains('Plan your shipments').click();
  cy.get('h1').contains('Figure out your shipments');
  cy.nextPage();
  cy.get('h1').contains('How do you want to move your belongings?');

  cy.get('input[type="radio"]').eq(2).check({ force: true });
  cy.nextPage();
  cy.location().should((loc) => {
    expect(loc.pathname).to.match(/^\/moves\/[^/]+\/nts-start/);
  });

  // pickup date
  cy.get('input[name="pickup.requestedDate"]').type('08/02/2020').blur();

  // pickup location
  cy.get(`input[name="useCurrentResidence"]`).check({ force: true });

  // releasing agent
  cy.get(`[data-testid="firstName"]`).type('John');
  cy.get(`[data-testid="lastName"]`).type('Lee');
  cy.get(`[data-testid="phone"]`).type('9999999999');
  cy.get(`[data-testid="email"]`).type('john@example.com').blur();

  // remarks
  cy.get(`[data-testid="remarks"]`).first().type('some customer remark');

  cy.nextPage();
}

function customerCreatesAnNTSRShipment() {
  cy.get('[data-testid="shipment-selection-btn"]').contains('Add another shipment').click();
  cy.get('h1').contains('How do you want this group of things moved?');

  cy.get('input[type="radio"]').eq(3).check({ force: true });
  cy.nextPage();
  cy.location().should((loc) => {
    expect(loc.pathname).to.match(/^\/moves\/[^/]+\/ntsr-start/);
  });

  // pickup date
  cy.get('input[name="delivery.requestedDate"]').type('01/02/2020').blur();

  // no delivery location

  // receiving agent
  cy.get(`[data-testid="firstName"]`).type('James');
  cy.get(`[data-testid="lastName"]`).type('Bond');
  cy.get(`[data-testid="phone"]`).type('7777777777');
  cy.get(`[data-testid="email"]`).type('007@example.com').blur();

  // remarks
  cy.get(`[data-testid="remarks"]`).first().type('some other customer remark');

  cy.nextPage();
}
