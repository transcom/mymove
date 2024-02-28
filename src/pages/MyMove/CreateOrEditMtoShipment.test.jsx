/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import { screen } from '@testing-library/react';

import { CreateOrEditMtoShipment } from './CreateOrEditMtoShipment';

import { customerRoutes } from 'constants/routes';
import { SHIPMENT_OPTIONS } from 'shared/constants';
import { renderWithRouterProp } from 'testUtils';

const mockParams = { moveId: 'move123', mtoShipmentId: 'shipment123' };
const mockPath = customerRoutes.SHIPMENT_EDIT_PATH;

const defaultProps = {
  pageList: ['page1', 'anotherPage/:foo/:bar'],
  pageKey: 'page1',
  fetchCustomerData: jest.fn(),
  updateMTOShipment: jest.fn(),
  updateAllMoves: jest.fn(),
  mtoShipment: {},
  currentResidence: {},
  serviceMember: {
    id: '1234',
  },
  orders: {
    new_duty_location: {
      address: {
        postalCode: '20050',
      },
    },
    authorizedWeight: 5000,
    entitlement: {
      proGear: 2000,
      proGearSpouse: 500,
    },
  },
  serviceMemberMoves: {
    currentMove: [
      {
        createdAt: '2024-02-28T19:42:27.489Z',
        eTag: 'MjAyNC0wMi0yOFQxOTo0MjoyNy40ODkyMjha',
        id: 'move124',
        moveCode: '3W4PTF',
        orders: {},
        status: 'DRAFT',
        submittedAt: '0001-01-01T00:00:00.000Z',
        updatedAt: '0001-01-01T00:00:00.000Z',
      },
    ],
    previousMoves: [
      {
        createdAt: '2024-02-28T15:43:46.721Z',
        eTag: 'MjAyNC0wMi0yOFQxODo0OTowMi41NDk3MDda',
        id: '88604b7c-9fe7-4d2f-bbc8-659022e2457f',
        moveCode: '9VWGKW',
        mtoShipments: [
          {
            agents: [
              {
                agentType: 'RELEASING_AGENT',
                createdAt: '2024-02-28T17:08:28.150Z',
                email: 'djordan218@gmail.com',
                firstName: 'Daniel',
                id: '19f06ce5-da64-4082-8478-7358ac45f398',
                lastName: 'Jordan',
                mtoShipmentID: '4e160e31-d171-4792-afe2-85cf1c0fb7f5',
                phone: '555-555-5555',
                updatedAt: '2024-02-28T19:20:18.955Z',
              },
              {
                agentType: 'RECEIVING_AGENT',
                createdAt: '2024-02-28T17:08:28.152Z',
                email: 'brittany.bortnem@gmail.com',
                firstName: 'Brittany',
                id: '7440acb4-8dd2-4725-baba-8952c8ccb57b',
                lastName: 'Jordan',
                mtoShipmentID: '4e160e31-d171-4792-afe2-85cf1c0fb7f5',
                updatedAt: '2024-02-28T19:20:18.958Z',
              },
            ],
            createdAt: '2024-02-28T17:08:28.146Z',
            customerRemarks: '',
            eTag: 'MjAyNC0wMi0yOFQxOToyMDoxOC45NTk5NTFa',
            hasSecondaryDeliveryAddress: false,
            hasSecondaryPickupAddress: false,
            id: 'shipment123',
            moveTaskOrderID: '88604b7c-9fe7-4d2f-bbc8-659022e2457f',
            requestedDeliveryDate: '2024-02-29',
            requestedPickupDate: '2024-03-30',
            shipmentType: 'HHG',
            status: 'SUBMITTED',
            updatedAt: '2024-02-28T19:20:18.959Z',
          },
        ],
        orders: {},
        status: 'DRAFT',
        submittedAt: '0001-01-01T00:00:00.000Z',
        updatedAt: '0001-01-01T00:00:00.000Z',
      },
    ],
  },
};

const mockHHGShipment = {
  id: 'mock id',
  moveTaskOrderId: 'move123',
  customerRemarks: 'mock remarks',
  requestedPickupDate: '1 Mar 2020',
  requestedDeliveryDate: '30 Mar 2020',
  pickupAddress: {
    streetAddress1: '812 S 129th St',
    city: 'San Antonio',
    state: 'TX',
    postalCode: '78234',
  },
  destinationAddress: {
    streetAddress1: '441 SW Rio de la Plata Drive',
    city: 'Tacoma',
    state: 'WA',
    postalCode: '98421',
  },
  shipmentType: 'HHG',
};

const mockPPMShipment = {
  id: 'mock id',
  moveTaskOrderId: 'move123',
  shipmentType: 'PPM',
};

const renderComponent = (props, options) =>
  renderWithRouterProp(<CreateOrEditMtoShipment {...defaultProps} {...props} />, options);

describe('CreateOrEditMtoShipment component', () => {
  it('fetches customer data on mount', () => {
    renderComponent(
      {
        mtoShipment: { shipmentType: SHIPMENT_OPTIONS.NTSR },
      },
      { path: mockPath, params: mockParams },
    );
    expect(defaultProps.fetchCustomerData).toHaveBeenCalled();
  });

  describe('when creating a new shipment', () => {
    it('renders the PPM date and location page if the shipment type is PPM', async () => {
      renderComponent(
        {
          mtoShipment: { ...mockPPMShipment, shipmentType: SHIPMENT_OPTIONS.NTSR },
        },
        { path: mockPath, params: mockParams, search: `?type=${SHIPMENT_OPTIONS.PPM}`, includeProviders: true },
      );

      expect(screen.getByRole('heading', { level: 1 })).toHaveTextContent('PPM date & location');
    });

    it('renders the MtoShipmentForm component right away', async () => {
      renderComponent(
        {
          mtoShipment: { shipmentType: SHIPMENT_OPTIONS.HHG },
        },
        { path: mockPath, params: mockParams, search: `?type=${SHIPMENT_OPTIONS.HHG}` },
      );

      expect(await screen.findByRole('heading', { level: 1 })).toHaveTextContent(
        'Movers pack and transport this shipment',
      );
      expect(screen.queryByText('Loading, please wait...')).not.toBeInTheDocument();
    });
  });

  describe('when editing an existing shipment', () => {
    it('renders the loader right away', () => {
      renderComponent({}, { path: customerRoutes.SHIPMENT_EDIT_PATH, params: mockParams });

      expect(screen.getByText('Loading, please wait...')).toBeInTheDocument();
    });

    it('renders the MtoShipmentForm after an HHG shipment has loaded', async () => {
      renderComponent(
        { mtoShipment: mockHHGShipment },
        { path: customerRoutes.SHIPMENT_EDIT_PATH, params: mockParams },
      );

      expect(await screen.findByRole('heading', { level: 1 })).toHaveTextContent(
        'Movers pack and transport this shipment',
      );
      expect(screen.queryByText('Loading, please wait...')).not.toBeInTheDocument();
    });

    it('renders the PPM date and location page after a PPM shipment has loaded', async () => {
      renderComponent(
        { mtoShipment: mockPPMShipment },
        { path: customerRoutes.SHIPMENT_EDIT_PATH, params: mockParams, includeProviders: true },
      );

      expect(screen.getByRole('heading', { level: 1 })).toHaveTextContent('PPM date & location');
    });
  });
});
