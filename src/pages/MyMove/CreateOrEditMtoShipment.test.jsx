/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import { render, screen } from '@testing-library/react';

import { CreateOrEditMtoShipment } from './CreateOrEditMtoShipment';

import { SHIPMENT_OPTIONS } from 'shared/constants';
import { MockProviders } from 'testUtils';

jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom'),
  useHistory: () => ({
    push: jest.fn(),
  }),
  useParams: () => ({
    moveCode: 'move123',
  }),
}));

function getMockMatchProp(path = '') {
  return {
    path,
    isExact: false,
    url: '',
    params: { moveId: 'move123' },
  };
}

const defaultProps = {
  location: {
    search: '',
  },
  match: {
    path: '',
    isExact: false,
    url: '',
    params: { moveId: 'move123' },
  },
  pageList: ['page1', 'anotherPage/:foo/:bar'],
  pageKey: 'page1',
  history: {
    goBack: jest.fn(),
    push: jest.fn(),
    replace: jest.fn(),
  },
  fetchCustomerData: jest.fn(),
  updateMTOShipment: jest.fn(),
  selectedMoveType: '',
  mtoShipment: {},
  currentResidence: {},
  serviceMember: {
    id: '1234',
    weight_allotment: {
      total_weight_self: 5000,
    },
  },
  orders: {
    new_duty_location: {
      address: {
        postalCode: '20050',
      },
    },
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

const renderComponent = (props, options) => render(<CreateOrEditMtoShipment {...defaultProps} {...props} />, options);

describe('CreateOrEditMtoShipment component', () => {
  it('fetches customer data on mount', () => {
    renderComponent({
      selectedMoveType: SHIPMENT_OPTIONS.NTSR,
    });
    expect(defaultProps.fetchCustomerData).toHaveBeenCalled();
  });

  describe('when creating a new shipment', () => {
    it('renders the PPM date and location page if the shipment type is PPM', async () => {
      renderComponent(
        {
          location: {
            search: `?type=${SHIPMENT_OPTIONS.PPM}`,
          },
        },
        { wrapper: MockProviders },
      );

      expect(screen.getByRole('heading', { level: 1 })).toHaveTextContent('PPM date & location');
    });

    it('renders the MtoShipmentForm component right away', async () => {
      renderComponent({
        location: {
          search: `?type=${SHIPMENT_OPTIONS.HHG}`,
        },
      });

      expect(await screen.findByRole('heading', { level: 1 })).toHaveTextContent(
        'Movers pack and transport this shipment',
      );
      expect(screen.queryByText('Loading, please wait...')).not.toBeInTheDocument();
    });
  });

  describe('when editing an existing shipment', () => {
    it('renders the loader right away', () => {
      renderComponent({
        match: getMockMatchProp('/moves/:moveId/shipments/:mtoShipmentId/edit'),
      });

      expect(screen.getByText('Loading, please wait...')).toBeInTheDocument();
    });

    it('renders the MtoShipmentForm after an HHG shipment has loaded', async () => {
      renderComponent({
        match: getMockMatchProp('/moves/:moveId/shipments/:mtoShipmentId/edit'),
        mtoShipment: mockHHGShipment,
      });

      expect(await screen.findByRole('heading', { level: 1 })).toHaveTextContent(
        'Movers pack and transport this shipment',
      );
      expect(screen.queryByText('Loading, please wait...')).not.toBeInTheDocument();
    });

    it('renders the PPM date and location page after a PPM shipment has loaded', async () => {
      renderComponent(
        {
          match: getMockMatchProp('/moves/:moveId/shipments/:mtoShipmentId/edit'),
          mtoShipment: mockPPMShipment,
        },
        { wrapper: MockProviders },
      );

      expect(screen.getByRole('heading', { level: 1 })).toHaveTextContent('PPM date & location');
    });
  });
});
