describe('Customer NTS Setup flow', function () {
  // profile2@complete.draft
  const profileCompleteUser = '2194daed-3589-408f-b988-e9889c9f120e';
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

  it('Edits an NTS shipment', function () {
    cy.apiSignInAsUser(ntsUser);
    customerVisitsReviewPage();
    customerReviewsShipmentInfoOnReviewPageBeforeEdit();
    customerEditsNTSShipment();
  });

  it('Edits an NTS shipment from homepage', function () {
    cy.apiSignInAsUser(ntsUser);
    customerEditsNTSShipmentFromHomePage();
  });

  it('Submits an NTS shipment from homepage', function () {
    cy.apiSignInAsUser(ntsUser);
    customerVisitsReviewPage();
    customerReviewsShipmentInfoOnReviewPageAfterEdit();
    customerSubmitsNTSShipmentMoveFromHomePage();
  });
});

function customerSubmitsNTSShipmentMoveFromHomePage() {
  cy.nextPage();
  cy.location().should((loc) => {
    expect(loc.pathname).to.match(/^\/moves\/[^/]+\/agreement/);
  });
  cy.get('.wizard-header').should('not.exist');

  cy.get('input[name="signature"]').type('Jane Doe');
  cy.completeFlow();
}

function customerEditsNTSShipmentFromHomePage() {
  cy.get('[data-testid="shipment-list-item-container"]').contains('NTS').click();
  cy.get('textarea[data-testid="remarks"]').clear().type('Warning: glass').blur();
  cy.get('input[name="secondaryPickup.address.city"]').clear().type('The Hills');
  cy.get('button').contains('Save').click();
  cy.location().should((loc) => {
    expect(loc.pathname).to.match(/^\/moves\/[^/]+\/review/);
  });
  cy.get('[data-testid="nts-summary"]').contains('The Hills, CA 90210');
}

function customerReviewsNTSMoveDetails() {
  cy.get('[data-testid="review-move-header"]').contains('Review your details');
  cy.get('[data-testid="nts-summary"]').contains('NTS');

  // Requested pickup date
  cy.get('[data-testid="nts-summary"]').contains('02 Aug 2020');

  // Pickup location
  cy.get('[data-testid="nts-summary"]').contains('123 Any Street P.O. Box 12345');
  cy.get('[data-testid="nts-summary"]').contains('Beverly Hills, CA 90210');

  // Secondary pickup address
  cy.get('[data-testid="nts-summary"]').contains('123 Some address');
  cy.get('[data-testid="nts-summary"]').contains('Some city, CA 90210');

  // Releasing agent
  cy.get('[data-testid="nts-summary"]').contains('John Lee');
  cy.get('[data-testid="nts-summary"]').contains('999-999-9999');
  cy.get('[data-testid="nts-summary"]').contains('john@example.com');

  // Remarks
  cy.get('[data-testid="nts-summary"]').contains('some customer remark');
}

function customerEditsNTSShipment() {
  cy.get('button[data-testid="edit-nts-shipment-btn"]').contains('Edit').click();
  cy.get('input[name="pickup.requestedDate"]').clear().type('12/25/2020').blur();
  cy.get('input[name="pickup.agent.lastName"]').clear().type('Bourne').blur();
  cy.get('textarea[data-testid="remarks"]').clear().type('Handle with care').blur();
  cy.get('input[name="secondaryPickup.address.streetAddress1"]').clear().type('999 Any Street');
  cy.get('button').contains('Save').click();
  cy.location().should((loc) => {
    expect(loc.pathname).to.match(/^\/moves\/[^/]+\/review/);
  });
  cy.get('[data-testid="nts-summary"]').contains('25 Dec 2020');
  cy.get('[data-testid="nts-summary"]').contains('Jason Bourne');
  cy.get('[data-testid="nts-summary"]').contains('Handle with care');
  cy.get('[data-testid="nts-summary"]').contains('999 Any Street');
}

function customerVisitsReviewPage() {
  cy.get('button[data-testid="review-and-submit-btn"]').contains('Review and submit').click();
  cy.get('[data-testid="review-move-header"]').contains('Review your details');
}

function customerReviewsShipmentInfoOnReviewPageBeforeEdit() {
  cy.get('[data-testid="nts-summary"]').contains('NTS');

  // Requested pickup date
  cy.get('[data-testid="nts-summary"]').contains('15 Mar 2020');

  // Pickup location
  cy.get('[data-testid="nts-summary"]').contains('123 Any Street P.O. Box 12345');
  cy.get('[data-testid="nts-summary"]').contains('Beverly Hills, CA 90210');

  // Secondary pickup address
  cy.get('[data-testid="nts-summary"]').contains('123 Any Street P.O. Box 12345');
  cy.get('[data-testid="nts-summary"]').contains('Beverly Hills, CA 90210');

  // Releasing agent
  cy.get('[data-testid="nts-summary"]').contains('Jason Ash');
  cy.get('[data-testid="nts-summary"]').contains('202-555-9301');
  cy.get('[data-testid="nts-summary"]').contains('jason.ash@example.com');

  // Remarks
  cy.get('[data-testid="nts-summary"]').contains('Please treat gently');
}

function customerReviewsShipmentInfoOnReviewPageAfterEdit() {
  cy.get('[data-testid="nts-summary"]').contains('NTS');

  // Requested pickup date
  cy.get('[data-testid="nts-summary"]').contains('25 Dec 2020');

  // Pickup location
  cy.get('[data-testid="nts-summary"]').contains('123 Any Street P.O. Box 12345');
  cy.get('[data-testid="nts-summary"]').contains('Beverly Hills, CA 90210');

  // Secondary pickup address
  cy.get('[data-testid="nts-summary"]').contains('999 Any Street');
  cy.get('[data-testid="nts-summary"]').contains('The Hills, CA 90210');

  // Releasing agent
  cy.get('[data-testid="nts-summary"]').contains('Jason Bourne');
  cy.get('[data-testid="nts-summary"]').contains('202-555-9301');
  cy.get('[data-testid="nts-summary"]').contains('jason.ash@example.com');

  // Remarks
  cy.get('[data-testid="nts-summary"]').contains('Warning: glass');
}

function customerCreatesAnNTSShipment() {
  cy.get('[data-testid="shipment-selection-btn"]').contains('Set up your shipments').click();
  cy.nextPage();
  cy.get('input[type="radio"]').eq(2).check({ force: true });
  cy.nextPage();
  cy.location().should((loc) => {
    expect(loc.pathname).to.match(/^\/moves\/[^/]+\/new-shipment/);
  });

  // pickup date
  cy.get('input[name="pickup.requestedDate"]').type('08/02/2020').blur();

  // pickup location
  cy.get(`input[name="useCurrentResidence"]`).check({ force: true });

  // secondary pickup location
  cy.get(`input[data-testid="has-secondary-pickup"]`).check({ force: true });
  cy.get(`input[name="secondaryPickup.address.streetAddress1"]`).type('123 Some address');
  cy.get(`input[name="secondaryPickup.address.city"]`).type('Some city');
  cy.get(`select[name="secondaryPickup.address.state"]`).select('CA');
  cy.get(`input[name="secondaryPickup.address.postalCode"]`).type('90210').blur();

  // releasing agent
  cy.get('input[name="pickup.agent.firstName"]').type('John');
  cy.get('input[name="pickup.agent.lastName"]').type('Lee');
  cy.get('input[name="pickup.agent.phone"]').type('9999999999');
  cy.get('input[name="pickup.agent.email"]').type('john@example.com').blur();

  // remarks
  cy.get(`[data-testid="remarks"]`).first().type('some customer remark');

  cy.nextPage();
}
