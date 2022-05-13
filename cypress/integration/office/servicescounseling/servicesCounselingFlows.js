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
    cy.wait(['@patchServiceCounselingCompleted', '@getMoves']);

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

  it('is able to edit allowances', () => {
    const moveLocator = 'RET1RE';

    // TOO Moves queue
    cy.wait(['@getSortedMoves']);
    cy.contains(moveLocator).click();
    cy.url().should('include', `/moves/${moveLocator}/details`);

    // Move Details page
    cy.wait(['@getMoves', '@getOrders', '@getMTOShipments', '@getMTOServiceItems']);

    // Navigate to Edit allowances page
    cy.get('[data-testid="edit-allowances"]').contains('Edit allowances').click();

    // Toggle between Edit Allowances and Edit Orders page
    cy.get('[data-testid="view-orders"]').click();
    cy.url().should('include', `/moves/${moveLocator}/orders`);
    cy.get('[data-testid="view-allowances"]').click();
    cy.url().should('include', `/moves/${moveLocator}/allowances`);

    cy.wait(['@getMoves', '@getOrders']);

    cy.get('form').within(($form) => {
      // Edit pro-gear, pro-gear spouse, RME, SIT, and OCIE fields
      cy.get('input[name="proGearWeight"]').clear().type('1999');
      cy.get('input[name="proGearWeightSpouse"]').clear().type('499');
      cy.get('input[name="requiredMedicalEquipmentWeight"]').clear().type('999');
      cy.get('input[name="storageInTransit"]').clear().type('199');
      cy.get('input[name="organizationalClothingAndIndividualEquipment"]').siblings('label[for="ocieInput"]').click();

      // Edit grade and authorized weight
      cy.get('select[name=agency]').contains('Army');
      cy.get('select[name=agency]').select('Navy');
      cy.get('select[name="grade"]').contains('E-1');
      cy.get('select[name="grade"]').select('W-2');

      //Edit DependentsAuthorized
      cy.get('input[name="dependentsAuthorized"]').siblings('label[for="dependentsAuthorizedInput"]').click();

      // Edit allowances page | Save
      cy.get('button').contains('Save').should('be.enabled').click().should('be.disabled');
    });

    cy.wait(['@patchAllowances']);

    // Verify edited values are saved
    cy.url().should('include', `/moves/${moveLocator}/details`);

    cy.wait(['@getMoves', '@getOrders', '@getMTOShipments', '@getMTOServiceItems']);

    cy.get('[data-testid="progear"]').contains('1,999');
    cy.get('[data-testid="spouseProgear"]').contains('499');
    cy.get('[data-testid="rme"]').contains('999');
    cy.get('[data-testid="storageInTransit"]').contains('199');
    cy.get('[data-testid="ocie"]').contains('Unauthorized');

    cy.get('[data-testid="branchRank"]').contains('Navy');
    cy.get('[data-testid="branchRank"]').contains('W-2');
    cy.get('[data-testid="dependents"]').contains('Unauthorized');

    // Edit allowances page | Cancel
    cy.get('[data-testid="edit-allowances"]').contains('Edit allowances').click();
    cy.get('button').contains('Cancel').click();
    cy.url().should('include', `/moves/${moveLocator}/details`);
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
