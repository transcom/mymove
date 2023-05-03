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
  mtoShipment: {},
  currentResidence: {},
  serviceMember: {
    id: '1234',
    weight_allotment: {
      total_weight_self: 5000,
      total_weight_self_plus_dependents: 7500,
      pro_gear_weight: 2000,
      pro_gear_weight_spouse: 500,
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
