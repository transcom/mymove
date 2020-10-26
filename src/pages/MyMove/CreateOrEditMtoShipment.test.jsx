/*  react/jsx-props-no-spreading */
import React from 'react';
import { mount } from 'enzyme';

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
  },
  fetchCustomerData: jest.fn(),
  createMTOShipment: jest.fn(),
  updateMTOShipment: jest.fn(),
  selectedMoveType: '',
  mtoShipment: {},
  currentResidence: {},
  serviceMember: {
    weight_allotment: {
      total_weight_self: 5000,
    },
  },
};

const mockMtoShipment = {
  id: 'mock id',
  moveTaskOrderId: 'move123',
  customerRemarks: 'mock remarks',
  requestedPickupDate: '1 Mar 2020',
  requestedDeliveryDate: '30 Mar 2020',
  pickupAddress: {
    street_address_1: '812 S 129th St',
    city: 'San Antonio',
    state: 'TX',
    postal_code: '78234',
  },
  destinationAddress: {
    street_address_1: '441 SW Rio de la Plata Drive',
    city: 'Tacoma',
    state: 'WA',
    postal_code: '98421',
  },
};

const mountCreateOrEditMtoShipment = (props) => mount(<CreateOrEditMtoShipment {...defaultProps} {...props} />);

describe('CreateOrEditMtoShipment component', () => {
  it('fetches customer data on mount', () => {
    mountCreateOrEditMtoShipment({
      selectedMoveType: SHIPMENT_OPTIONS.NTSR,
    });
    expect(defaultProps.fetchCustomerData).toHaveBeenCalled();
  });

  describe('when creating a new shipment', () => {
    it('renders the MtoShipmentForm component right away', () => {
      const createWrapper = mountCreateOrEditMtoShipment({
        selectedMoveType: SHIPMENT_OPTIONS.HHG,
        isCreate: true,
      });
      expect(createWrapper.find('MtoShipmentForm').length).toBe(1);
      expect(createWrapper.find('LoadingPlaceholder').exists()).toBe(false);
    });
  });

  describe('when editing an existing shipment', () => {
    const editWrapper = mountCreateOrEditMtoShipment({
      selectedMoveType: SHIPMENT_OPTIONS.NTS,
      match: getMockMatchProp('/moves/:moveId/mto-shipments/:mtoShipmentId/edit'),
    });

    it('renders the loader right away', () => {
      expect(editWrapper.find('LoadingPlaceholder').exists()).toBe(true);
      expect(editWrapper.find('MtoShipmentForm').length).toBe(0);
    });

    it('renders the MtoShipmentForm after an MTO shipment has loaded', () => {
      editWrapper.setProps({
        mtoShipment: mockMtoShipment,
      });
      editWrapper.update();
      expect(editWrapper.find('LoadingPlaceholder').exists()).toBe(false);
      expect(editWrapper.find('MtoShipmentForm').length).toBe(1);
    });
  });
});
