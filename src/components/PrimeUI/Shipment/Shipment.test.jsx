import React from 'react';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import Shipment from './Shipment';

import { formatDateFromIso, formatPrimeAPIFullAddress } from 'utils/formatters';
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
  approvedDate: '2022-05-24',
  createdAt: '2022-05-24T21:06:35.888Z',
  eTag: 'MjAyMi0wNS0yNFQyMTowNzoyMS4wNjc0MzJa',
  id: '88ececed-eaf1-42e2-b060-cd90d11ad080',
  moveTaskOrderID: '9c7b255c-2981-4bf8-839f-61c7458e2b4d',
  ppmShipment: {
    advance: 598700,
    advanceRequested: true,
    createdAt: '2022-05-24T21:06:35.901Z',
    destinationPostalCode: '30813',
    eTag: 'MjAyMi0wNS0yNFQyMTowNjozNS45MDEwMjNa',
    estimatedIncentive: 1000000,
    estimatedWeight: 4000,
    expectedDepartureDate: '2020-03-15',
    hasProGear: false,
    id: '5b21b808-6933-43ea-8f6f-02fc0a639835',
    pickupPostalCode: '90210',
    shipmentId: '88ececed-eaf1-42e2-b060-cd90d11ad080',
    status: 'SUBMITTED',
    submittedAt: '2022-05-24T21:06:35.890Z',
    updatedAt: '2022-05-24T21:06:35.901Z',
  },
  shipmentType: 'PPM',
  status: 'APPROVED',
  updatedAt: '2022-05-24T21:07:21.067Z',
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
  it('renders the component when shipment is a PPM', () => {
    render(
      <MockProviders>
        <Shipment shipment={ppmShipment} moveId={moveId} />
      </MockProviders>,
    );

    const field = screen.getByText('PPM Status:');
    expect(field).toBeInTheDocument();
    expect(field.nextElementSibling.textContent).toBe(ppmShipment.ppmShipment.status);
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
