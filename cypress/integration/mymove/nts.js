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
});

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
