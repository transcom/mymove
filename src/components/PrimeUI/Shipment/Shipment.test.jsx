import React from 'react';
import { render, screen } from '@testing-library/react';

import Shipment from './Shipment';

import { formatDateFromIso, formatPrimeAPIFullAddress } from 'shared/formatters';
import { MockProviders } from 'testUtils';

const shipmentId = 'ce01a5b8-9b44-4511-8a8d-edb60f2a4aee';
const moveId = '9c7b255c-2981-4bf8-839f-61c7458e2b4d';

const mockPush = jest.fn();

jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom'),
  useParams: jest.fn().mockReturnValue({
    moveCode: 'LR4T8V',
    moveCodeOrID: '9c7b255c-2981-4bf8-839f-61c7458e2b4d',
    shipmentId: 'ce01a5b8-9b44-4511-8a8d-edb60f2a4aee',
  }),
  useHistory: () => ({
    push: mockPush,
  }),
}));

jest.mock('services/primeApi', () => ({
  ...jest.requireActual('services/primeApi'),
  updatePrimeMTOShipment: jest.fn().mockImplementation(() => Promise.resolve()),
}));

jest.mock('hooks/queries', () => ({
  usePrimeSimulatorGetMove: jest.fn(),
}));
const approvedMoveTaskOrder = {
  moveTaskOrder: {
    id: '9c7b255c-2981-4bf8-839f-61c7458e2b4d',
    moveCode: 'LR4T8V',
    mtoShipments: [
      {
        actualPickupDate: '2020-03-17',
        agents: [],
        approvedDate: '2021-10-20',
        createdAt: '2021-10-21T18:24:41.377Z',
        customerRemarks: 'Please treat gently',
        destinationAddress: {
          city: 'Fairfield',
          country: 'US',
          eTag: 'MjAyMS0xMC0xOFQxODoyNDo0MS4zNzI3NDJa',
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
          country: 'US',
          eTag: 'MjAyMS0xMC0xOFQxODoyNDo0MS4zNjc3Mjda',
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
        shipmentType: 'HHG_LONGHAUL_DOMESTIC',
        status: 'APPROVED',
        updatedAt: '2021-10-22T18:24:41.377Z',
        mtoServiceItems: null,
        reweigh: {
          id: '1234',
          weight: 9000,
          requestedAt: '2021-10-23T18:24:41.377Z',
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
  it('renders the component headings and links without errors', async () => {
    render(mockedComponent);
    const shipmentLevelHeader = screen.getByRole('heading', { name: 'HHG shipment', level: 3 });
    expect(shipmentLevelHeader).toBeInTheDocument();

    const updateShipmentLink = screen.getByText(/Update Shipment/, { selector: 'a.usa-button' });
    expect(updateShipmentLink).toBeInTheDocument();
    expect(updateShipmentLink.getAttribute('href')).toBe(`/simulator/moves/${moveId}/shipments/${shipmentId}`);

    const addServiceItemLink = screen.getByText(/Add Service Item/, { selector: 'a.usa-button' });
    expect(addServiceItemLink).toBeInTheDocument();
    expect(addServiceItemLink.getAttribute('href')).toBe(`/shipments/${shipmentId}/service-items/new`);

    expect(screen.queryAllByRole('link', { name: 'Edit' })).toHaveLength(2);
  });

  it('renders the shipment address values', async () => {
    render(mockedComponent);
    const shipment = approvedMoveTaskOrder.moveTaskOrder.mtoShipments[0];

    expect(screen.getByText(formatPrimeAPIFullAddress(shipment.pickupAddress))).toBeInTheDocument();
    expect(screen.getByText(formatPrimeAPIFullAddress(shipment.destinationAddress))).toBeInTheDocument();
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
  it('renders the component with missing reweigh error', async () => {
    render(
      <MockProviders>
        <Shipment shipment={shipmentMissingReweighWeight} moveId={moveId} />
      </MockProviders>,
    );

    await expect(screen.getByText('Missing')).toBeInTheDocument();
    await expect(screen.getByText('Reweigh Weight:')).toBeInTheDocument();
    await expect(screen.getByText('Reweigh Requested Date:')).toBeInTheDocument();
  });

  // Reweigh isn't missing here, it was not requested and therefore should not be present
  // in shipment display table
  it('renders the component with no reweigh requested', async () => {
    render(
      <MockProviders>
        <Shipment shipment={shipmentNoReweighRequested} moveId={moveId} />
      </MockProviders>,
    );

    await expect(screen.queryByText('Reweigh Weight:')).not.toBeInTheDocument();
    await expect(screen.queryByText('Reweigh Requested Date:')).not.toBeInTheDocument();
  });
});
