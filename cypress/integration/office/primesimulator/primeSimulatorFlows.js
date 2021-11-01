import { PrimeSimulatorUserType } from '../../../support/constants';

describe('Prime simulator user', () => {
  before(() => {
    cy.prepareOfficeApp();
  });

  beforeEach(() => {
    cy.intercept('GET', '**/prime/v1/swagger.yaml').as('getPrimeClient');
    cy.intercept('GET', '**/prime/v1/moves').as('listMoves');
    cy.intercept('GET', '**/prime/v1/move-task-orders/**').as('getMove');
    cy.intercept('PATCH', '**/prime/v1/mto-shipments/**').as('updateMTOShipment');

    const userId = 'cf5609e9-b88f-4a98-9eda-9d028bc4a515';
    cy.apiSignInAsUser(userId, PrimeSimulatorUserType);
  });

  it('is able to update a shipment', () => {
    const moveLocator = 'PRMUPD';
    const moveID = 'ef4a2b75-ceb3-4620-96a8-5ccf26dddb16';
    const shipmentID = '5375f237-430c-406d-9ec8-5a27244d563a';

    // wait for the the available moves page to load
    cy.wait(['@getPrimeClient', '@listMoves']);

    // select the PRMUPD move from the list
    cy.contains(moveLocator).click();
    cy.url().should('include', `/simulator/moves/${moveID}/details`);

    // waits for the move details page to load
    cy.wait(['@getMove']);
    cy.contains('Update Shipment').click();

    // waits for the update shipment page to load
    cy.url().should('include', `/simulator/moves/${moveID}/shipments/${shipmentID}`);
    cy.wait(['@getMove']);

    cy.get('input[name="scheduledPickupDate"]').type('01 Nov 2021').blur();
    cy.get('input[name="actualPickupDate"]').type('02 Nov 2021').blur();

    cy.get('input[name="destinationAddress.street_address_1"]').type('142 E Barrel Hoop Circle');
    cy.get('input[name="destinationAddress.city"]').type('Joshua Tree');
    cy.get('select[name="destinationAddress.state"]').select('CA');
    cy.get('input[name="destinationAddress.postal_code"]').type('92252');

    cy.contains('Save').click();
    cy.wait(['@updateMTOShipment']);

    cy.url().should('include', `/simulator/moves/${moveID}/details`);
    cy.wait(['@getMove']);
  });
});
