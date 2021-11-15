import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';

import { formatDateFromIso } from '../../../shared/formatters';

import Shipment from './Shipment';

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
        agents: [
          {
            agentType: 'RELEASING_AGENT',
            createdAt: '2021-10-18T18:24:41.521Z',
            eTag: 'MjAyMS0xMC0xOFQxODoyNDo0MS41MjE4NzNa',
            email: 'test@test.email.com',
            firstName: 'Test',
            id: 'f2619e1b-7729-4b97-845d-6ae1ebe299f2',
            lastName: 'Agent',
            mtoShipmentID: 'ce01a5b8-9b44-4511-8a8d-edb60f2a4aee',
            phone: '202-555-9301',
            updatedAt: '2021-10-19T18:24:41.521Z',
          },
        ],
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
        secondaryPickupAddress: {
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

const formatAddress = (address) => {
  const { streetAddress1, streetAddress2, city, state, postalCode } = address;
  return `${streetAddress1}, ${streetAddress2}, ${city}, ${state} ${postalCode}`;
};

describe('Shipment details component', () => {
  it('renders the component without errors', async () => {
    render(mockedComponent);
    const shipmentLevelHeader = screen.getByRole('heading', { name: 'HHG shipment', level: 3 });
    expect(shipmentLevelHeader).toBeInTheDocument();

    const updateShipmentLink = screen.getByText(/Update Shipment/, { selector: 'a.usa-button' });
    expect(updateShipmentLink).toBeInTheDocument();
    expect(updateShipmentLink.getAttribute('href')).toBe(`/simulator/moves/${moveId}/shipments/${shipmentId}`);

    const addServiceItemLink = screen.getByText(/Add Service Item/, { selector: 'a.usa-button' });
    expect(addServiceItemLink).toBeInTheDocument();
    expect(addServiceItemLink.getAttribute('href')).toBe(`/shipments/${shipmentId}/service-items/new`);

    const shipmentFields = [
      'Status:',
      'Shipment ID:',
      'Shipment eTag:',
      'Requested Pickup Date:',
      'Scheduled Pickup Date:',
      'Actual Pickup Date:',
      'Actual Weight:',
      'Estimated Weight:',
      'Reweigh Weight:',
      'Reweigh Requested Date:',
      'Destination Address:',
      'Pickup Address:',
      'Created at:',
      'Approved at:',
    ];
    await waitFor(() => {
      shipmentFields.forEach((shipmentField) => {
        expect(screen.getByText(shipmentField)).toBeVisible();
      });
    });
    expect(screen.queryAllByRole('link', { name: 'Edit' })).toHaveLength(2);
  });

  it('renders the shipment details values', async () => {
    render(mockedComponent);
    const shipment = approvedMoveTaskOrder.moveTaskOrder.mtoShipments[0];
    const shipmentFieldValues = [
      shipment.status,
      shipment.id,
      shipment.eTag,
      shipment.requestedPickupDate,
      shipment.scheduledPickupDate,
      shipment.actualPickupDate,
      shipment.primeActualWeight,
      shipment.primeEstimatedWeight,
      shipment.reweigh.weight,
      formatDateFromIso(shipment.reweigh.requestedAt, 'YYYY-MM-DD'),
      formatDateFromIso(shipment.createdAt, 'YYYY-MM-DD'),
      shipment.approvedDate,
    ];

    await waitFor(() => {
      shipmentFieldValues.forEach((shipmentFieldValue) => {
        // console.log(shipmentFieldValue);
        // console.log(new RegExp(shipmentFieldValue));
        // expect(screen.findByText(new RegExp(shipmentFieldValue))).toBeVisible();
        expect(screen.getByText(shipmentFieldValue)).toBeVisible();
      });
    });

    expect(screen.getByText(formatAddress(shipment.pickupAddress))).toBeInTheDocument();
    expect(screen.getByText(formatAddress(shipment.destinationAddress))).toBeInTheDocument();
  });

  it('renders the shipment addresses', async () => {
    // const { getByText } = render(mockedComponent);
    render(mockedComponent);

    // These won't match
    // getByText("Hello world");
    // getByText(/Hello world/);
    /*
    getByText((content, node) => {

      console.log('node.textContent');
      console.log(node.textContent);
      const hasText = node => node.textContent === content;
      const nodeHasText = hasText(node);
      const childrenDontHaveText = Array.from(node.children).every(
        child => !hasText(child)
      );

      return nodeHasText && childrenDontHaveText;
    });

     */

    /*
    const shipment = approvedMoveTaskOrder.moveTaskOrder.mtoShipments[0];
const shipmentAddressFieldValues = [
shipment.destinationAddress.streetAddress1,
`/${shipment.destinationAddress.streetAddress2}/`,
shipment.destinationAddress.city,
shipment.destinationAddress.zip,
shipment.pickupAddress.streetAddress1,
shipment.pickupAddress.streetAddress1,
shipment.pickupAddress.streetAddress2,
shipment.pickupAddress.city,
shipment.pickupAddress.zip,
];
 */

    /*
    await waitFor(() => {
      shipmentAddressFieldValues.forEach((shipmentAddressFieldValue)=>{
        console.log(shipmentAddressFieldValue);
        console.log(new RegExp(shipmentAddressFieldValue));
        //expect(screen.findByText(new RegExp(shipmentAddressFieldValue))).toBeVisible();
        expect(screen.getByText(shipmentAddressFieldValue)).toBeVisible();
      });
    });

     */

    /*
    await waitFor(() => {
      console.log(shipment.status);
      console.log(new RegExp(shipment.status));
      expect(screen.getByText(shipment.status)).toBeVisible();
      expect(screen.getByText(new RegExp(shipment.status))).toBeInTheDocument();
      expect(screen.findByText(shipment.id)).toBeVisible();
      expect(screen.findByText(shipment.eTag)).toBeVisible();
      expect(screen.findByText(shipment.requestedPickupDate)).toBeVisible();
      expect(screen.findByText(shipment.scheduledPickupDate)).toBeVisible();
      expect(screen.findByText(shipment.actualPickupDate)).toBeVisible();
      expect(screen.findByText(shipment.primeActualWeight)).toBeVisible();
      expect(screen.findByText(shipment.primeEstimatedWeight)).toBeVisible();
      expect(screen.findByText(shipment.reweigh.weight)).toBeVisible();
      expect(screen.findByText(formatDateFromIso(shipment.reweigh.requestedAt, 'YYYY-MM-DD'))).toBeVisible();
      expect(screen.findByText(formatPrimeAPIShipmentAddress(shipment.destinationAddress))).toBeVisible();
      expect(screen.findByText(formatPrimeAPIShipmentAddress(shipment.pickupAddress))).toBeVisible();
      expect(screen.findByText(formatDateFromIso(shipment.createdAt, 'YYYY-MM-DD'))).toBeVisible();
      expect(screen.findByText(shipment.approvedDate)).toBeVisible();
    });
     */
  });
});
