import { ServicesCounselorOfficeUserType } from '../../../support/constants';

describe('Services counselor user', () => {
  before(() => {
    cy.prepareOfficeApp();
  });

  beforeEach(() => {
    cy.intercept('**/ghc/v1/swagger.yaml').as('getGHCClient');
    cy.intercept('**/ghc/v1/queues/counseling?page=1&perPage=20&sort=submittedAt&order=asc').as('getSortedMoves');
    cy.intercept('**/ghc/v1/queues/counseling?page=1&perPage=20&sort=submittedAt&order=asc&locator=SCE4ET').as(
      'getFilterSortedMoves',
    );
    cy.intercept('**/ghc/v1/move/**').as('getMoves');
    cy.intercept('**/ghc/v1/orders/**').as('getOrders');
    cy.intercept('**/ghc/v1/move_task_orders/**/mto_shipments').as('getMTOShipments');
    cy.intercept('**/ghc/v1/move_task_orders/**/mto_service_items').as('getMTOServiceItems');
    cy.intercept('**/ghc/v1/move-task-orders/**/status/service-counseling-completed').as(
      'patchServiceCounselingCompleted',
    );
    cy.intercept('**/ghc/v1/moves/**/financial-review-flag').as('financialReviewFlagCompleted');
    cy.intercept('POST', '**/ghc/v1/mto-shipments').as('createShipment');
    cy.intercept('PATCH', '**/ghc/v1/move_task_orders/**/mto_shipments/**').as('patchShipment');
    cy.intercept('PATCH', '**/ghc/v1/counseling/orders/**/allowances').as('patchAllowances');

    const userId = 'a6c8663f-998f-4626-a978-ad60da2476ec';
    cy.apiSignInAsUser(userId, ServicesCounselorOfficeUserType);
  });

  it('is able to click on move and submit after using the move code filter', () => {
    const moveLocator = 'SCE4ET';

    /**
     * SC Moves queue
     */
    cy.wait(['@getSortedMoves']);
    cy.get('input[name="locator"]').as('moveCodeFilterInput');

    // type in move code/locator to filter
    cy.get('@moveCodeFilterInput').type(moveLocator).blur();
    cy.wait(['@getFilterSortedMoves']);

    // check if results appear, should be 1
    // and see if result have move code
    cy.get('tbody > tr').as('results');
    cy.get('@results').should('have.length', 1);
    cy.get('@results').first().contains(moveLocator);

    // click result to navigate to move details page
    cy.get('@results').first().click();
    cy.url().should('include', `/counseling/moves/${moveLocator}/details`);
    cy.wait(['@getMoves', '@getOrders', '@getMTOShipments', '@getMTOServiceItems']);

    /**
     * Move Details page
     */
    // click to trigger confirmation modal
    cy.contains('Submit move details').click();

    // modal should pop up with text
    cy.get('h2').contains('Are you sure?');
    cy.get('p').contains('You can’t make changes after you submit the move.');

    // click submit
    cy.get('button').contains('Yes, submit').click();
    cy.waitFor(['@patchServiceCounselingCompleted', '@getMoves']);

    // verify success alert
    cy.contains('Move submitted.');
  });

  it('is able to flag a move for financial review', () => {
    cy.wait(['@getSortedMoves']);
    // It doesn't matter which move we click on in the queue.
    cy.get('td').first().click();
    cy.url().should('include', `details`);
    cy.wait(['@getMoves', '@getOrders', '@getMTOShipments', '@getMTOServiceItems']);

    // click to trigger financial review modal
    cy.contains('Flag move for financial review').click();

    // Enter information in modal and submit
    cy.get('label').contains('Yes').click();
    cy.get('textarea').type('Because I said so...');

    // Click save on the modal
    cy.get('button').contains('Save').click();
    cy.wait(['@financialReviewFlagCompleted']);

    // Verify sucess alert and tag
    cy.contains('Move flagged for financial review.');
    cy.contains('Flagged for financial review');
  });

  it('is able to unflag a move for financial review', () => {
    cy.wait(['@getSortedMoves']);
    // It doesn't matter which move we click on in the queue.
    cy.get('td').first().click();
    cy.url().should('include', `details`);
    cy.wait(['@getMoves', '@getOrders', '@getMTOShipments', '@getMTOServiceItems']);

    // click to trigger financial review modal
    cy.contains('Edit').click();

    // Enter information in modal and submit
    cy.get('label').contains('No').click();

    // Click save on the modal
    cy.get('button').contains('Save').click();
    cy.wait(['@financialReviewFlagCompleted']);

    // Verify sucess alert and tag
    cy.contains('Move unflagged for financial review.');
  });

  it('is able to add a shipment', () => {
    const deliveryDate = new Date().toLocaleDateString('en-US');

    const moveLocator = 'S3PAR3';

    /**
     * SC Moves queue
     */
    cy.wait(['@getSortedMoves']);
    cy.get('input[name="locator"]').as('moveCodeFilterInput');
    cy.get('@moveCodeFilterInput').type(moveLocator).blur();
    cy.get('td').first().click();
    cy.url().should('include', `details`);
    cy.wait(['@getMoves', '@getOrders', '@getMTOShipments', '@getMTOServiceItems']);

    // add a shipment
    cy.get('[data-testid="dropdown"]').first().select('HHG');
    cy.get('#requestedPickupDate').clear().type(deliveryDate).blur();
    cy.get('[data-testid="useCurrentResidence"]').click({ force: true });
    cy.get('#requestedDeliveryDate').clear().type('16 Mar 2022').blur();
    cy.get('#has-delivery-address').click({ force: true });
    cy.get('input[name="delivery.address.streetAddress1"]').type('7 q st');
    cy.get('input[name="delivery.address.city"]').type('city');
    cy.get('select[name="delivery.address.state"]').select('OH');
    cy.get('input[name="delivery.address.postalCode"]').type('90210');
    cy.get('select[name="destinationType"]').select('Home of record (HOR)');
    cy.get('[data-testid="submitForm"]').click();
    // the shipment should be saved with the type
    cy.wait('@createShipment');
  });

  it.only('is able to edit allowances', () => {
    cy.wait(['@getSortedMoves']);
    // It doesn't matter which move we click on in the queue.
    cy.get('td').first().click();
    cy.url().should('include', `details`);
    cy.wait(['@getMoves', '@getOrders', '@getMTOShipments', '@getMTOServiceItems']);
    cy.get('[data-testid="edit-allowances"]').click();

    // the form
    cy.get('[data-testid="proGearWeightInput"]').focus().clear().type('1999').blur();
    cy.get('[data-testid="sitInput"]').focus().clear().type('199').blur();

    // Edit allowances page | Save
    cy.get('[data-testid="scAllowancesSave"]').click();

    cy.wait('@patchAllowances');

    cy.location().should((loc) => {
      expect(loc.pathname).to.include('/details');
    });

    // things should save and then load afterward with new data
    cy.wait(['@getMoves', '@getOrders', '@getMTOShipments', '@getMTOServiceItems']);
    cy.get('[data-testid="progear"]').contains('1,999');
    cy.get('[data-testid="storageInTransit"]').contains('199');
  });

  it('is able to see and use the left navigation', () => {
    const moveLocator = 'RET1RE';

    /**
     * SC Moves queue
     */
    cy.wait(['@getSortedMoves']);
    cy.get('input[name="locator"]').as('moveCodeFilterInput');
    cy.get('@moveCodeFilterInput').type(moveLocator).blur();
    cy.get('td').first().click();
    cy.url().should('include', `details`);
    cy.wait(['@getMoves', '@getOrders', '@getMTOShipments', '@getMTOServiceItems']);

    cy.get('a[href*="#shipments"]').contains('Shipments');
    cy.get('a[href*="#orders"]').contains('Orders');
    cy.get('a[href*="#allowances"]').contains('Allowances');
    cy.get('a[href*="#customer-info"]').contains('Customer info');

    cy.get('[data-testid="requestedShipmentsTag"]').contains('3');

    // Assert that the window has scrolled after clicking a left nav item
    cy.get('#customer-info').click().window().its('scrollY').should('not.equal', 0);
  });

  it('is able to edit a shipment', () => {
    const deliveryDate = new Date().toLocaleDateString('en-US');

    const moveLocator = 'RET1RE';

    /**
     * SC Moves queue
     */
    cy.wait(['@getSortedMoves']);
    cy.get('input[name="locator"]').as('moveCodeFilterInput');
    cy.get('@moveCodeFilterInput').type(moveLocator).blur();
    cy.get('td').first().click();
    cy.url().should('include', `details`);
    cy.wait(['@getMoves', '@getOrders', '@getMTOShipments', '@getMTOServiceItems']);

    // edit a shipment
    cy.get('[data-testid="ShipmentContainer"] .usa-button').first().click();
    cy.get('#requestedPickupDate').clear().type(deliveryDate).blur();
    cy.get('[data-testid="useCurrentResidence"]').click({ force: true });
    cy.get('#requestedDeliveryDate').clear().type('16 Mar 2022').blur();
    cy.get('#has-delivery-address').click({ force: true });
    cy.get('input[name="delivery.address.streetAddress1"]').clear().type('7 q st');
    cy.get('input[name="delivery.address.city"]').clear().type('city');
    cy.get('select[name="delivery.address.state"]').select('OH');
    cy.get('input[name="delivery.address.postalCode"]').clear().type('90210');
    cy.get('select[name="destinationType"]').select('Home of selection (HOS)');
    cy.get('[data-testid="submitForm"]').click();
    // the shipment should be saved with the type
    cy.wait('@patchShipment');
    cy.get('.usa-alert__text').contains('Your changes were saved.');
  });

  it('is able to see that the tag next to shipment is updated', () => {
    const moveLocator = 'RET1RE';

    /**
     * SC Moves queue
     */
    cy.wait(['@getSortedMoves']);
    cy.get('input[name="locator"]').as('moveCodeFilterInput');
    cy.get('@moveCodeFilterInput').type(moveLocator).blur();
    cy.get('td').first().click();
    cy.url().should('include', `details`);
    cy.wait(['@getMoves', '@getOrders', '@getMTOShipments', '@getMTOServiceItems']);

    cy.get('a[href*="#shipments"]').contains('Shipments');

    // Verify that there's a tag on the left nav that flags missing information
    cy.get('[data-testid="requestedShipmentsTag"]').contains('3');

    // Edit the shipment so that the tag disappears
    cy.get('[data-testid="ShipmentContainer"] .usa-button').last().click();
    cy.get('select[name="destinationType"]').select('Home of selection (HOS)');
    cy.get('[data-testid="submitForm"]').click();
    // the shipment should be saved with the type
    cy.wait('@patchShipment');
    cy.get('.usa-alert__text').contains('Your changes were saved.');

    // Verify that the tag after the update is a 2 since missing information was filled
    cy.get('[data-testid="requestedShipmentsTag"]').contains('2');
  });
});
