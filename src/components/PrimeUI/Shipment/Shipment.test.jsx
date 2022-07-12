import React from 'react';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import Shipment from './Shipment';

import {
  formatCents,
  formatDateFromIso,
  formatPrimeAPIFullAddress,
  formatYesNoInputValue,
  toDollarString,
} from 'utils/formatters';
import { MockProviders } from 'testUtils';

const shipmentId = 'ce01a5b8-9b44-4511-8a8d-edb60f2a4aee';
const moveId = '9c7b255c-2981-4bf8-839f-61c7458e2b4d';

const approvedMoveTaskOrder = {
  moveTaskOrder: {
    id: '9c7b255c-2981-4bf8-839f-61c7458e2b4d',
    moveCode: 'LR4T8V',
    mtoShipments: [
      {
        actualPickupDate: '2020-03-17',
        agents: [],
        approvedDate: '2021-10-20',
        counselorRemarks: 'These are counselor remarks for an HHG.',
        createdAt: '2021-10-21',
        customerRemarks: 'Please treat gently',
        destinationAddress: {
          city: 'Fairfield',
          id: 'bfe61147-5fd7-426e-b473-54ccf77bde35',
          postalCode: '94535',
          state: 'CA',
          streetAddress1: '987 Any Avenue',
          streetAddress2: 'P.O. Box 9876',
          streetAddress3: 'c/o Some Person',
        },
        eTag: 'MjAyMS0xMC0xOFQxODoyNDo0MS4zNzc5Nzha',
        firstAvailableDeliveryDate: null,
        id: 'ce01a5b8-9b44-4511-8a8d-edb60f2a4aee',
        moveTaskOrderID: '9c7b255c-2981-4bf8-839f-61c7458e2b4d',
        pickupAddress: {
          city: 'Beverly Hills',
          id: 'cf159eca-162c-4131-84a0-795e684416a6',
          postalCode: '90210',
          state: 'CA',
          streetAddress1: '123 Any Street',
          streetAddress2: 'P.O. Box 12345',
          streetAddress3: 'c/o Some Person',
        },
        primeActualWeight: 2000,
        primeEstimatedWeight: 1400,
        primeEstimatedWeightRecordedDate: null,
        requestedPickupDate: '2020-03-15',
        requiredDeliveryDate: null,
        scheduledPickupDate: '2020-03-16',
        secondaryDeliveryAddress: {
          city: null,
          postalCode: null,
          state: null,
          streetAddress1: null,
        },
        shipmentType: 'HHG',
        status: 'APPROVED',
        updatedAt: '2021-10-22',
        mtoServiceItems: null,
        reweigh: {
          id: '1234',
          weight: 9000,
          verificationReason: 'Reweigh requested.',
          requestedAt: '2021-10-23',
        },
      },
    ],
  },
};

const mockedComponent = (
  <MockProviders>
    <Shipment shipment={approvedMoveTaskOrder.moveTaskOrder.mtoShipments[0]} moveId={moveId} />
  </MockProviders>
);

describe('Shipment details component', () => {
  it('renders the component headings and links without errors', () => {
    render(mockedComponent);
    const shipmentLevelHeader = screen.getByRole('heading', { name: 'HHG shipment', level: 3 });
    expect(shipmentLevelHeader).toBeInTheDocument();

    const updateShipmentLink = screen.getByText(/Update Shipment/, { selector: 'a.usa-button' });
    expect(updateShipmentLink).toBeInTheDocument();
    expect(updateShipmentLink.getAttribute('href')).toBe(`/simulator/moves/${moveId}/shipments/${shipmentId}`);

    const addServiceItemLink = screen.getByText(/Add Service Item/, { selector: 'a.usa-button' });
    expect(addServiceItemLink).toBeInTheDocument();
    expect(addServiceItemLink.getAttribute('href')).toBe(`/shipments/${shipmentId}/service-items/new`);

    expect(screen.queryAllByRole('link', { name: 'Edit' })).toHaveLength(3);
  });

  it('renders the shipment address values', async () => {
    render(mockedComponent);
    const shipment = approvedMoveTaskOrder.moveTaskOrder.mtoShipments[0];

    expect(screen.getByText(formatPrimeAPIFullAddress(shipment.pickupAddress))).toBeInTheDocument();
    expect(screen.getByText(formatPrimeAPIFullAddress(shipment.destinationAddress))).toBeInTheDocument();
  });

  it('renders the shipment info', () => {
    render(mockedComponent);
    const shipment = approvedMoveTaskOrder.moveTaskOrder.mtoShipments[0];

    // shipment text values
    let field = screen.getByText('Status:');
    expect(field).toBeInTheDocument();
    expect(field.nextElementSibling.textContent).toBe(shipment.status);

    field = screen.getByText('Shipment ID:');
    expect(field).toBeInTheDocument();
    expect(field.nextElementSibling.textContent).toBe(shipment.id);

    field = screen.getByText('Shipment eTag:');
    expect(field).toBeInTheDocument();
    expect(field.nextElementSibling.textContent).toBe(shipment.eTag);

    field = screen.getByText('Requested Pickup Date:');
    expect(field).toBeInTheDocument();
    expect(field.nextElementSibling.textContent).toBe(shipment.requestedPickupDate);

    field = screen.getByText('Scheduled Pickup Date:');
    expect(field).toBeInTheDocument();
    expect(field.nextElementSibling.textContent).toBe(shipment.scheduledPickupDate);

    field = screen.getByText('Actual Pickup Date:');
    expect(field).toBeInTheDocument();
    expect(field.nextElementSibling.textContent).toBe(shipment.actualPickupDate);

    field = screen.getByText('Estimated Weight:');
    expect(field).toBeInTheDocument();
    expect(field.nextElementSibling.textContent).toBe(shipment.primeEstimatedWeight.toString());

    field = screen.getByText('Actual Weight:');
    expect(field).toBeInTheDocument();
    expect(field.nextElementSibling.textContent).toBe(shipment.primeActualWeight.toString());

    field = screen.getByText('Reweigh Weight:');
    expect(field).toBeInTheDocument();
    expect(field.nextElementSibling.textContent).toBe(shipment.reweigh.weight.toString());

    field = screen.getByText('Reweigh Requested Date:');
    expect(field).toBeInTheDocument();
    expect(field.nextElementSibling.textContent).toBe(shipment.reweigh.requestedAt);

    field = screen.getByText('Pickup Address:');
    expect(field).toBeInTheDocument();
    expect(field.nextElementSibling.textContent).toContain(shipment.pickupAddress.city);
    expect(field.nextElementSibling.textContent).toContain(shipment.pickupAddress.state);
    expect(field.nextElementSibling.textContent).toContain(shipment.pickupAddress.streetAddress1);
    expect(field.nextElementSibling.textContent).toContain(shipment.pickupAddress.streetAddress2);
    expect(field.nextElementSibling.textContent).toContain(shipment.pickupAddress.postalCode);

    field = screen.getByText('Destination Address:');
    expect(field).toBeInTheDocument();
    expect(field.nextElementSibling.textContent).toContain(shipment.destinationAddress.city);
    expect(field.nextElementSibling.textContent).toContain(shipment.destinationAddress.state);
    expect(field.nextElementSibling.textContent).toContain(shipment.destinationAddress.streetAddress1);
    expect(field.nextElementSibling.textContent).toContain(shipment.destinationAddress.streetAddress2);
    expect(field.nextElementSibling.textContent).toContain(shipment.destinationAddress.postalCode);

    field = screen.getByText('Created at:');
    expect(field).toBeInTheDocument();
    expect(field.nextElementSibling.textContent).toBe(shipment.createdAt);

    field = screen.getByText('Approved at:');
    expect(field).toBeInTheDocument();
    expect(field.nextElementSibling.textContent).toBe(shipment.approvedDate);

    field = screen.getByText('Counselor Remarks:');
    expect(field).toBeInTheDocument();
    expect(field.nextElementSibling.textContent).toBe(shipment.counselorRemarks);

    // This is an HHG, so make sure elements that are specific to PPMs are not visible.
    const deleteShipmentButton = screen.queryByText(/Delete Shipment/, { selector: 'button' });
    expect(deleteShipmentButton).not.toBeInTheDocument();

    field = screen.queryByText('PPM Status:');
    expect(field).not.toBeInTheDocument();
  });
});

describe('Shipment details component fields and values are present', () => {
  const shipment = approvedMoveTaskOrder.moveTaskOrder.mtoShipments[0];
  it.each([
    ['Status:', shipment.status],
    ['Shipment ID:', shipment.id],
    ['Shipment eTag:', shipment.eTag],
    ['Requested Pickup Date:', shipment.requestedPickupDate],
    ['Scheduled Pickup Date:', shipment.scheduledPickupDate],
    ['Actual Pickup Date:', shipment.actualPickupDate],
    ['Actual Weight:', shipment.primeActualWeight],
    ['Estimated Weight:', shipment.primeEstimatedWeight],
    ['Reweigh Weight:', shipment.reweigh.weight],
    ['Reweigh Remarks:', shipment.reweigh.verificationReason],
    ['Reweigh Requested Date:', formatDateFromIso(shipment.reweigh.requestedAt, 'YYYY-MM-DD')],
    ['Created at:', formatDateFromIso(shipment.createdAt, 'YYYY-MM-DD')],
    ['Approved at:', shipment.approvedDate],
  ])('Verify PrimeUI Move Shipment field %s with value %s is present', async (shipmentField, shipmentFieldValue) => {
    render(mockedComponent);
    await expect(screen.getByText(shipmentField)).toBeVisible();
    await expect(screen.getByText(shipmentFieldValue)).toBeVisible();
  });
});

const shipmentMissingReweighWeight = {
  ...approvedMoveTaskOrder.moveTaskOrder.mtoShipments[0],
  reweigh: {
    id: '1234',
    requestedAt: '2021-10-23T18:24:41.377Z',
  },
};

const shipmentNoReweighRequested = {
  ...approvedMoveTaskOrder.moveTaskOrder.mtoShipments[0],
  reweigh: null,
};

describe('Shipment has missing reweigh', () => {
  it('renders the component with missing reweigh error', () => {
    render(
      <MockProviders>
        <Shipment shipment={shipmentMissingReweighWeight} moveId={moveId} />
      </MockProviders>,
    );

    expect(screen.getByText('Missing')).toBeInTheDocument();
    expect(screen.getByText('Reweigh Weight:')).toBeInTheDocument();
    expect(screen.getByText('Reweigh Requested Date:')).toBeInTheDocument();
  });

  // Reweigh isn't missing here, it was not requested and therefore should not be present
  // in shipment display table
  it('renders the component with no reweigh requested', () => {
    render(
      <MockProviders>
        <Shipment shipment={shipmentNoReweighRequested} moveId={moveId} />
      </MockProviders>,
    );

    expect(screen.queryByText('Reweigh Weight:')).not.toBeInTheDocument();
    expect(screen.queryByText('Reweigh Remarks:')).not.toBeInTheDocument();
    expect(screen.queryByText('Reweigh Requested Date:')).not.toBeInTheDocument();
  });
});

const ppmShipment = {
  actualPickupDate: null,
  approvedDate: null,
  counselorRemarks: 'These are counselor remarks for a PPM.',
  createdAt: '2022-07-01T13:41:33.261Z',
  destinationAddress: {
    city: null,
    postalCode: null,
    state: null,
    streetAddress1: null,
  },
  eTag: 'MjAyMi0wNy0wMVQxNDoyMzoxOS43MzgzODla',
  firstAvailableDeliveryDate: null,
  id: '1b695b60-c3ed-401b-b2e3-808d095eb8cc',
  moveTaskOrderID: '7024c8c5-52ca-4639-bf69-dd8238308c98',
  pickupAddress: {
    city: null,
    postalCode: null,
    state: null,
    streetAddress1: null,
  },
  ppmShipment: {
    actualDestinationPostalCode: '30814',
    actualMoveDate: '2022-07-13',
    actualPickupPostalCode: '90212',
    advanceAmountReceived: 598600,
    advanceAmountRequested: 598700,
    approvedAt: '2022-07-03T14:20:21.620Z',
    createdAt: '2022-06-30T13:41:33.265Z',
    destinationPostalCode: '30813',
    eTag: 'MjAyMi0wNy0wMVQxNDoyMzoxOS43ODA1Mlo=',
    estimatedIncentive: 1000000,
    estimatedWeight: 4000,
    expectedDepartureDate: '2020-03-15',
    hasProGear: true,
    hasReceivedAdvance: true,
    hasRequestedAdvance: true,
    id: 'd733fe2f-b08d-434a-ad8d-551f4d597b03',
    netWeight: 3900,
    pickupPostalCode: '90210',
    proGearWeight: 1987,
    reviewedAt: '2022-07-02T14:20:14.636Z',
    secondaryDestinationPostalCode: '30814',
    secondaryPickupPostalCode: '90211',
    shipmentId: '1b695b60-c3ed-401b-b2e3-808d095eb8cc',
    sitEstimatedCost: 123456,
    sitEstimatedDepartureDate: '2022-07-13',
    sitEstimatedEntryDate: '2022-07-05',
    sitEstimatedWeight: 1100,
    sitExpected: true,
    sitLocation: 'DESTINATION',
    spouseProGearWeight: 498,
    status: 'SUBMITTED',
    submittedAt: '2022-07-01T13:41:33.252Z',
    updatedAt: '2022-07-01T14:23:19.780Z',
  },
  primeEstimatedWeightRecordedDate: null,
  requestedPickupDate: null,
  requiredDeliveryDate: null,
  scheduledPickupDate: null,
  secondaryDeliveryAddress: {
    city: null,
    postalCode: null,
    state: null,
    streetAddress1: null,
  },
  secondaryPickupAddress: {
    city: null,
    postalCode: null,
    state: null,
    streetAddress1: null,
  },
  shipmentType: 'PPM',
  status: 'APPROVED',
  updatedAt: '2022-07-01T14:23:19.738Z',
  mtoServiceItems: [],
};

const ppmShipmentWaitingOnCustomer = {
  ...ppmShipment,
  ppmShipment: {
    ...ppmShipment.ppmShipment,
    status: 'WAITING_ON_CUSTOMER',
  },
};

const ppmShipmentMissingObject = {
  ...ppmShipment,
  ppmShipment: null,
};

describe('PPM shipments are handled', () => {
  it('PPM fields header is present', () => {
    render(
      <MockProviders>
        <Shipment shipment={ppmShipment} moveId={moveId} />
      </MockProviders>,
    );

    const ppmFieldsHeader = screen.getByRole('heading', { name: 'PPM-specific fields', level: 4 });
    expect(ppmFieldsHeader).toBeInTheDocument();
  });

  const ppmShipmentFields = ppmShipment.ppmShipment;

  it.each([
    ['PPM Status:', ppmShipmentFields.status],
    ['PPM Shipment ID:', ppmShipmentFields.id],
    ['PPM Shipment eTag:', ppmShipmentFields.eTag],
    ['PPM Created at:', formatDateFromIso(ppmShipmentFields.createdAt, 'YYYY-MM-DD')],
    ['PPM Updated at:', formatDateFromIso(ppmShipmentFields.updatedAt, 'YYYY-MM-DD')],
    ['PPM Expected Departure Date:', formatDateFromIso(ppmShipmentFields.expectedDepartureDate, 'YYYY-MM-DD')],
    ['PPM Actual Move Date:', formatDateFromIso(ppmShipmentFields.actualMoveDate, 'YYYY-MM-DD')],
    ['PPM Submitted at:', formatDateFromIso(ppmShipmentFields.submittedAt, 'YYYY-MM-DD')],
    ['PPM Reviewed at:', formatDateFromIso(ppmShipmentFields.reviewedAt, 'YYYY-MM-DD')],
    ['PPM Approved at:', formatDateFromIso(ppmShipmentFields.approvedAt, 'YYYY-MM-DD')],
    ['PPM Pickup Postal Code:', ppmShipmentFields.pickupPostalCode],
    ['PPM Secondary Pickup Postal Code:', ppmShipmentFields.secondaryPickupPostalCode],
    ['PPM Destination Postal Code:', ppmShipmentFields.destinationPostalCode],
    ['PPM Secondary Destination Postal Code:', ppmShipmentFields.secondaryDestinationPostalCode],
    ['PPM SIT Expected:', formatYesNoInputValue(ppmShipmentFields.sitExpected)],
    ['PPM Estimated Weight:', ppmShipmentFields.estimatedWeight.toString()],
    ['PPM Net Weight:', ppmShipmentFields.netWeight.toString()],
    ['PPM Has Pro Gear:', formatYesNoInputValue(ppmShipmentFields.hasProGear)],
    ['PPM Pro Gear Weight:', ppmShipmentFields.proGearWeight.toString()],
    ['PPM Spouse Pro Gear Weight:', ppmShipmentFields.spouseProGearWeight.toString()],
    ['PPM Estimated Incentive:', toDollarString(formatCents(ppmShipmentFields.estimatedIncentive))],
    ['PPM SIT Location:', ppmShipmentFields.sitLocation],
    ['PPM SIT Estimated Weight:', ppmShipmentFields.sitEstimatedWeight.toString()],
    ['PPM SIT Estimated Entry Date:', formatDateFromIso(ppmShipmentFields.sitEstimatedEntryDate, 'YYYY-MM-DD')],
    ['PPM SIT Estimated Departure Date:', formatDateFromIso(ppmShipmentFields.sitEstimatedDepartureDate, 'YYYY-MM-DD')],
    ['PPM SIT Estimated Cost:', toDollarString(formatCents(ppmShipmentFields.sitEstimatedCost))],
    ['PPM Actual Pickup Postal Code:', ppmShipmentFields.actualPickupPostalCode],
    ['PPM Actual Destination Postal Code:', ppmShipmentFields.actualDestinationPostalCode],
    ['PPM Has Requested Advance:', formatYesNoInputValue(ppmShipmentFields.hasRequestedAdvance)],
    ['PPM Advance Amount Requested:', toDollarString(formatCents(ppmShipmentFields.advanceAmountRequested))],
    ['PPM Has Received Advance:', formatYesNoInputValue(ppmShipmentFields.hasReceivedAdvance)],
    ['PPM Advance Amount Received:', toDollarString(formatCents(ppmShipmentFields.advanceAmountReceived))],
  ])('PPM shipment field %s with value %s is present', async (ppmShipmentField, ppmShipmentFieldValue) => {
    render(
      <MockProviders>
        <Shipment shipment={ppmShipment} moveId={moveId} />
      </MockProviders>,
    );

    const field = screen.getByText(ppmShipmentField);
    await expect(field).toBeInTheDocument();
    await expect(field.nextElementSibling.textContent).toBe(ppmShipmentFieldValue);
  });

  it('PPM can be deleted', () => {
    const onDelete = jest.fn();

    render(
      <MockProviders>
        <Shipment shipment={ppmShipment} moveId={moveId} onDelete={onDelete} />
      </MockProviders>,
    );

    const deleteShipmentButton = screen.queryByText(/Delete Shipment/, { selector: 'button' });
    expect(deleteShipmentButton).toBeInTheDocument();

    userEvent.click(deleteShipmentButton);
    let modalTitle = screen.getByText('Are you sure?');
    expect(modalTitle).toBeInTheDocument();

    const modalDeleteButton = screen.getByText('Delete shipment', { selector: 'button.usa-button--destructive' });
    userEvent.click(modalDeleteButton);
    expect(onDelete).toHaveBeenCalledTimes(1);

    modalTitle = screen.queryByText('Are you sure?');
    expect(modalTitle).not.toBeInTheDocument();
  });

  it('PPM status does not allow deletion', () => {
    render(
      <MockProviders>
        <Shipment shipment={ppmShipmentWaitingOnCustomer} moveId={moveId} />
      </MockProviders>,
    );

    const deleteShipmentButton = screen.queryByText(/Delete Shipment/, { selector: 'button' });
    expect(deleteShipmentButton).not.toBeInTheDocument();
  });

  it('PPM shipment is missing ppmShipment object', () => {
    render(
      <MockProviders>
        <Shipment shipment={ppmShipmentMissingObject} moveId={moveId} />
      </MockProviders>,
    );

    const deleteShipmentButton = screen.queryByText(/Delete Shipment/, { selector: 'button' });
    expect(deleteShipmentButton).not.toBeInTheDocument();
  });
});
