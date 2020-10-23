import { customerFillsInProfileInformation, customerFillsOutOrdersInformation } from './utilities/customer';

describe('NTS Setup flow', function () {
  before(() => {
    cy.prepareCustomerApp();
  });

  beforeEach(() => {
    cy.removeFetch();
    cy.server();
    cy.route('POST', '/internal/service_members').as('createServiceMember');
    cy.route('PATCH', '**/internal/mto-shipments/**').as('patchShipment');
    cy.route('GET', '/internal/moves/**/mto_shipments').as('getMTOShipments');
    cy.route('GET', '/internal/users/logged_in').as('getLoggedInUser');
  });

  it('Sets up an NTS shipment', function () {
    cy.signInAsNewMilMoveUser();
    customerFillsInProfileInformation();
    customerFillsOutOrdersInformation();
    customerChoosesAnNTSMove();
  });

  it('Edits an NTS shipment', function () {
    const ntsUser = '583cfbe1-cb34-4381-9e1f-54f68200da1b';
    cy.apiSignInAsPpmUser(ntsUser);
    customerVisitsReviewPage();
    customerEditsNTSShipment();
  });

  it('Edits an NTSr shipment', function () {
    const ntsUser = '583cfbe1-cb34-4381-9e1f-54f68200da1b';
    cy.apiSignInAsPpmUser(ntsUser);
    customerVisitsReviewPage();
    customerEditsNTSRShipment();
  });
});

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
  cy.get('dd').contains('01 Jan 2022');
  cy.get('dd').contains('123 Maple street');
  cy.get('dd').contains('Ketchum Ash');
  cy.get('dd').contains('Warning: fragile');
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
  cy.get('dd').contains('25 Dec 2020');
  cy.get('dd').contains('Jason Bourne');
  cy.get('dd').contains('Handle with care');
}

function customerVisitsReviewPage() {
  cy.visit('/');
  cy.get('button[data-testid="review-and-submit-btn"]').contains('Review and submit').click();
  cy.get('h2[data-testid="review-move-header"]').contains('Review your details');
}

function customerChoosesAnNTSMove() {
  cy.get('h1').contains('Figure out your shipments');
  cy.nextPage();
  cy.get('h1').contains('How do you want to move your belongings?');

  cy.get('input[type="radio"]').eq(2).check({ force: true });
  cy.nextPage();
  cy.location().should((loc) => {
    expect(loc.pathname).to.match(/^\/moves\/[^/]+\/nts-start/);
  });
}
