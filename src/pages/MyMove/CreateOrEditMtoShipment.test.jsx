/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import { render, waitFor, screen } from '@testing-library/react';

import { CreateOrEditMtoShipment } from './CreateOrEditMtoShipment';

import { SHIPMENT_OPTIONS } from 'shared/constants';

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
  orders: {
    authorizedWeight: 5000,
  },
};

const mockMtoShipment = {
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

const renderComponent = (props) => render(<CreateOrEditMtoShipment {...defaultProps} {...props} />);

describe('CreateOrEditMtoShipment component', () => {
  it('fetches customer data on mount', () => {
    renderComponent({
      selectedMoveType: SHIPMENT_OPTIONS.NTSR,
    });
    expect(defaultProps.fetchCustomerData).toHaveBeenCalled();
  });

  describe('when creating a new shipment', () => {
    it('redirects to the PPM start page if selected shipment type is PPM', async () => {
      renderComponent({
        location: {
          search: `?type=${SHIPMENT_OPTIONS.PPM}`,
        },
      });

      await waitFor(() => {
        expect(defaultProps.history.replace).toHaveBeenCalledWith('/moves/move123/ppm-start');
      });
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

    it('renders the MtoShipmentForm after an MTO shipment has loaded', async () => {
      renderComponent({
        match: getMockMatchProp('/moves/:moveId/shipments/:mtoShipmentId/edit'),
        mtoShipment: mockMtoShipment,
      });

      expect(await screen.findByRole('heading', { level: 1 })).toHaveTextContent(
        'Movers pack and transport this shipment',
      );
      expect(screen.queryByText('Loading, please wait...')).not.toBeInTheDocument();
    });
  });
});
