import { PrimeSimulatorUserType } from '../../../support/constants';

// I was trying to avoid importing moment.js, these can definitely be improved
const formatRelativeDate = (daysInterval) => {
  const dayFormat = { day: '2-digit' };
  const monthFormat = { month: 'short' };
  const yearFormat = { year: 'numeric' };

  const relativeDate = new Date();
  relativeDate.setDate(relativeDate.getDate() + daysInterval);

  return [
    relativeDate,
    `${relativeDate.toLocaleDateString(undefined, dayFormat)} ${relativeDate.toLocaleDateString(
      undefined,
      monthFormat,
    )} ${relativeDate.toLocaleDateString(undefined, yearFormat)}`,
  ];
};

const formatNumericDate = (date) => {
  const dayFormat = { day: '2-digit' };
  const monthFormat = { month: '2-digit' };
  const yearFormat = { year: 'numeric' };

  return `${date.toLocaleDateString(undefined, yearFormat)}-${date.toLocaleDateString(
    undefined,
    monthFormat,
  )}-${date.toLocaleDateString(undefined, dayFormat)}`;
};

describe('Prime simulator user', () => {
  before(() => {
    cy.prepareOfficeApp();
  });

  beforeEach(() => {
    cy.intercept('GET', '**/prime/v1/swagger.yaml').as('getPrimeClient');
    cy.intercept('GET', '**/prime/v1/moves').as('listMoves');
    cy.intercept('GET', '**/prime/v1/move-task-orders/**').as('getMove');
    cy.intercept('PATCH', '**/prime/v1/mto-shipments/**').as('updateMTOShipment');
    cy.intercept('POST', '**/prime/v1/payment-requests').as('createPaymentRequest');

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

    // there must be sufficient time prior to the pickup dates to update the estimated weight
    const [scheduledPickupDate, formattedScheduledPickupDate] = formatRelativeDate(11);
    cy.get('input[name="scheduledPickupDate"]').type(formattedScheduledPickupDate).blur();

    const [actualPickupDate, formattedActualPickupDate] = formatRelativeDate(12);
    cy.get('input[name="actualPickupDate"]').type(formattedActualPickupDate).blur();

    // update shipment does not require these fields but we need actual weight to create a payment request, we could
    // perform multiple updates.
    cy.get('input[name="estimatedWeight"]').type('{backspace}7500');
    cy.get('input[name="actualWeight"]').type('{backspace}8000');

    cy.get('input[name="destinationAddress.streetAddress1"]').type('142 E Barrel Hoop Circle');
    cy.get('input[name="destinationAddress.city"]').type('Joshua Tree');
    cy.get('select[name="destinationAddress.state"]').select('CA');
    cy.get('input[name="destinationAddress.postalCode"]').type('92252');

    cy.contains('Save').click();
    cy.wait(['@updateMTOShipment']);

    cy.url().should('include', `/simulator/moves/${moveID}/details`);
    cy.wait(['@getMove']);

    // If you added another shipment to the move you would want to scope these with within()
    cy.contains('Scheduled Pickup Date').siblings().contains(formatNumericDate(scheduledPickupDate));
    cy.contains('Actual Pickup Date').siblings().contains(formatNumericDate(actualPickupDate));

    cy.contains('Estimated Weight').siblings().contains('7500');
    cy.contains('Actual Weight').siblings().contains('8000');

    cy.contains('Destination Address').siblings().contains('142 E Barrel Hoop Circle, Joshua Tree, CA 92252');
  });

  it('is able to create a payment request', () => {
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
    cy.contains('Create Payment Request').click();

    // waits for the update shipment page to load
    cy.url().should('include', `/simulator/moves/${moveID}/payment-requests/new`);
    cy.wait(['@getMove']);

    // select all of the service items from the move and shipment, force is required because of USWDS positioning tricks
    cy.get('input[name="serviceItems"]').click({ multiple: true, force: true });

    cy.contains('Submit Payment Request').click();
    cy.wait(['@createPaymentRequest']);

    cy.url().should('include', `/simulator/moves/${moveID}/details`);
    cy.wait(['@getMove']);

    cy.contains('Successfully created payment request');

    // could also check for a payment request number but we won't know the value ahead of time
  });
});
