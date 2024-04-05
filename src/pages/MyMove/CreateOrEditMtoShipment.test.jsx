/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import { screen } from '@testing-library/react';

import { CreateOrEditMtoShipment } from './CreateOrEditMtoShipment';

import { customerRoutes } from 'constants/routes';
import { SHIPMENT_OPTIONS } from 'shared/constants';
import { renderWithRouterProp } from 'testUtils';
import { selectCurrentMoveFromAllMoves, selectCurrentShipmentFromMove } from 'store/entities/selectors';

const mockParams = { moveId: 'move123', mtoShipmentId: 'shipment123' };
const mockPath = customerRoutes.SHIPMENT_EDIT_PATH;

jest.mock('store/entities/selectors', () => ({
  ...jest.requireActual('store/entities/selectors'),
  selectServiceMemberFromLoggedInUser: jest.fn(),
  selectCurrentMoveFromAllMoves: jest.fn(),
  selectCurrentShipmentFromMove: jest.fn(),
}));

jest.mock('services/internalApi', () => ({
  ...jest.requireActual('services/internalApi'),
  getAllMoves: jest.fn().mockImplementation(() => Promise.resolve()),
}));

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
        id: 'move123',
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

const testMove = {
  createdAt: '2024-02-27T19:16:32.850Z',
  eTag: 'MjAyNC0wMi0yN1QxOToxNjozMi44NTAwNTda',
  id: 'move123',
  moveCode: 'WWYFP6',
  mtoShipments: [
    {
      createdAt: '2024-02-27T19:27:39.150Z',
      customerRemarks: '',
      destinationAddress: {
        city: 'Flagstaff',
        country: 'United States',
        id: '112e0d7b-90eb-44c4-80b1-5c1214fca0a7',
        postalCode: '86003',
        state: 'AZ',
        streetAddress1: 'N/A',
      },
      eTag: 'MjAyNC0wMi0yN1QxOToyNzozOS4xNTA3MjRa',
      hasSecondaryDeliveryAddress: false,
      hasSecondaryPickupAddress: false,
      id: 'f0082986-8e2f-443b-8411-191b3796e572',
      moveTaskOrderID: 'e23d629e-2a73-4b42-886b-fa60cb3db957',
      pickupAddress: {
        city: 'Tulsa',
        id: 'dac5e64d-1a69-4e83-a154-5fca04384544',
        postalCode: '74133',
        state: 'OK',
        streetAddress1: '8711 S 76th E Ave',
        streetAddress2: '',
      },
      requestedDeliveryDate: '2024-02-29',
      requestedPickupDate: '2024-03-01',
      shipmentType: 'HHG',
      status: 'SUBMITTED',
      updatedAt: '2024-02-27T19:27:39.150Z',
    },
  ],
  orders: {
    authorizedWeight: 8000,
    created_at: '2024-02-27T19:16:32.845Z',
    entitlement: {
      proGear: 2000,
      proGearSpouse: 500,
    },
    grade: 'E_6',
    has_dependents: false,
    id: '786e60ec-1bbd-48dd-bc12-b6ffcf212c54',
    issue_date: '2024-02-29',
    new_duty_location: {
      address: {
        city: 'Flagstaff',
        country: 'United States',
        id: 'cd51f4db-6195-473a-86cd-c3e5e07640e4',
        postalCode: '86003',
        state: 'AZ',
        streetAddress1: 'n/a',
      },
      address_id: 'cd51f4db-6195-473a-86cd-c3e5e07640e4',
      affiliation: null,
      created_at: '2024-02-27T18:22:12.471Z',
      id: '6ea57f62-2995-4b0c-a0a8-f11782cc8a3b',
      name: 'Flagstaff, AZ 86003',
      updated_at: '2024-02-27T18:22:12.471Z',
    },
    orders_type: 'PERMANENT_CHANGE_OF_STATION',
    originDutyLocationGbloc: 'BGAC',
    origin_duty_location: {
      address: {
        city: 'Aberdeen Proving Ground',
        country: 'United States',
        id: 'b6ca003e-1528-4e7c-a43e-830222ca7fb3',
        postalCode: '21005',
        state: 'MD',
        streetAddress1: 'n/a',
      },
      address_id: 'b6ca003e-1528-4e7c-a43e-830222ca7fb3',
      affiliation: 'ARMY',
      created_at: '2024-02-27T18:22:12.471Z',
      id: '61e092c4-575c-458a-9c3f-b93ad373c454',
      name: 'Aberdeen Proving Ground, MD 21005',
      transportation_office: {
        address: {
          city: 'Aberdeen Proving Ground',
          country: 'United States',
          id: 'ac4dbfa5-3068-4f8f-99d1-3cd850412bb9',
          postalCode: '21005',
          state: 'MD',
          streetAddress1: '4305 Susquehanna Ave',
          streetAddress2: 'Room 307',
        },
        created_at: '2018-05-28T14:27:41.772Z',
        gbloc: 'BGAC',
        id: '6a27dfbd-2a49-485f-86dd-49475d5facef',
        name: 'PPPO Aberdeen Proving Ground - USA',
        phone_lines: [],
        updated_at: '2018-05-28T14:27:41.772Z',
      },
      transportation_office_id: '6a27dfbd-2a49-485f-86dd-49475d5facef',
      updated_at: '2024-02-27T18:22:12.471Z',
    },
    report_by_date: '2024-02-29',
    service_member_id: 'c95824c7-425f-47e1-8264-bd6e55a2a2e4',
    spouse_has_pro_gear: false,
    status: 'DRAFT',
    updated_at: '2024-02-27T19:16:32.845Z',
    uploaded_orders: {
      id: 'b392f96f-20a0-43d2-9ca3-643cfd3b4182',
      service_member_id: 'c95824c7-425f-47e1-8264-bd6e55a2a2e4',
      uploads: [
        {
          bytes: 1137126,
          contentType: 'image/png',
          createdAt: '2024-02-27T19:16:38.998Z',
          filename: 'Screenshot 2024-02-15 at 12.22.53 PM.png',
          id: 'bc6c0e2d-fbef-4c32-8487-92c14b613040',
          status: 'PROCESSING',
          updatedAt: '2024-02-27T19:16:38.998Z',
          url: '/storage/user/f94c8fa6-89de-4594-be6a-ca6f1b4e22d0/uploads/bc6c0e2d-fbef-4c32-8487-92c14b613040?contentType=image%2Fpng',
        },
      ],
    },
  },
  status: 'DRAFT',
  submittedAt: '0001-01-01T00:00:00.000Z',
  updatedAt: '0001-01-01T00:00:00.000Z',
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
  selectCurrentMoveFromAllMoves.mockImplementation(() => testMove);
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
    selectCurrentMoveFromAllMoves.mockImplementation(() => testMove);

    it('renders the MtoShipmentForm after an HHG shipment has loaded', async () => {
      selectCurrentShipmentFromMove.mockImplementation(() => mockHHGShipment);
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
      selectCurrentShipmentFromMove.mockImplementation(() => mockPPMShipment);
      renderComponent(
        { mtoShipment: mockPPMShipment },
        { path: customerRoutes.SHIPMENT_EDIT_PATH, params: mockParams, includeProviders: true },
      );

      expect(screen.getByRole('heading', { level: 1 })).toHaveTextContent('PPM date & location');
    });
  });
});
