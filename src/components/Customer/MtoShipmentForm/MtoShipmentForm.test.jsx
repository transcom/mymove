/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import { mount } from 'enzyme';

import MtoShipmentForm from './MtoShipmentForm';

import { SHIPMENT_OPTIONS } from 'shared/constants';

const defaultProps = {
  isCreatePage: true,
  pageList: ['page1', 'anotherPage/:foo/:bar'],
  pageKey: 'page1',
  match: { isExact: false, path: '', url: '', params: { moveId: '' } },
  history: {
    goBack: jest.fn(),
    push: jest.fn(),
  },
  showLoggedInUser: jest.fn(),
  createMTOShipment: jest.fn(),
  updateMTOShipment: jest.fn(),
  newDutyStationAddress: {
    city: 'Fort Benning',
    state: 'GA',
    postal_code: '31905',
  },
  currentResidence: {
    city: 'Fort Benning',
    state: 'GA',
    postal_code: '31905',
    street_address_1: '123 Main',
  },
  serviceMember: {
    weight_allotment: {
      total_weight_self: 5000,
    },
  },
};

const mockMtoShipment = {
  id: 'mock id',
  moveTaskOrderId: 'mock move id',
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

const mountMtoShipmentForm = (props) => mount(<MtoShipmentForm {...defaultProps} {...props} />);

describe('MtoShipmentForm component', () => {
  describe('creating a new HHG shipment', () => {
    it('renders expected child components', () => {
      const wrapper = mountMtoShipmentForm({ selectedMoveType: SHIPMENT_OPTIONS.HHG });
      expect(wrapper.find('MtoShipmentForm').length).toBe(1);
      expect(wrapper.find('DatePickerInput').length).toBe(2);
      expect(wrapper.find('AddressFields').length).toBe(1);
      expect(wrapper.find('ContactInfoFields').length).toBe(2);
      expect(wrapper.find('input[name="customerRemarks"]').length).toBe(1);
    });

    // TODO - Formik & Enzyme don't play well together :( - https://github.com/formium/formik/issues/937
    // Displaying Delivery address fields is just tested in Edit mode for now with existing values
    it.skip('renders second address field when has delivery address', () => {
      /*
      const wrapper = mount(<MtoShipmentForm {...defaultProps} selectedMoveType={SHIPMENT_OPTIONS.HHG} />);
      const checkbox = wrapper.find('input[name="hasDeliveryAddress"][value="yes"]');
      expect(checkbox.length).toBe(1);

      act(() => {
        checkbox.simulate('change', { target: { name: 'hasDeliveryAddress', value: 'yes', checked: true } });
      });

      wrapper.update();
      expect(wrapper.find('AddressFields').length).toBe(2);
      */
    });
  });

  describe('editing an already existing HHG shipment', () => {
    it('renders the pre-filled MtoShipmentForm', () => {
      const wrapper = mountMtoShipmentForm({
        isCreatePage: false,
        selectedMoveType: SHIPMENT_OPTIONS.HHG,
        mtoShipment: mockMtoShipment,
      });
      expect(wrapper.find('MtoShipmentForm').length).toBe(1);
      expect(wrapper.find('DatePickerInput').length).toBe(2);
      expect(wrapper.find('AddressFields').length).toBe(2);
      expect(wrapper.find('ContactInfoFields').length).toBe(2);
      expect(wrapper.find('input[name="customerRemarks"]').length).toBe(1);
      expect(wrapper.find('TextInput[name="customerRemarks"]').prop('value')).toEqual(mockMtoShipment.customerRemarks);
      expect(wrapper.find('Field[name="delivery.address.street_address_1"]').prop('value')).toContain(
        mockMtoShipment.destinationAddress.street_address_1,
      );
    });
  });

  describe('creating a new NTS shipment', () => {
    it('renders expected child components', () => {
      const wrapper = mountMtoShipmentForm({ selectedMoveType: SHIPMENT_OPTIONS.NTS });
      expect(wrapper.find('MtoShipmentForm').length).toBe(1);
      expect(wrapper.find('DatePickerInput').length).toBe(1);
      expect(wrapper.find('AddressFields').length).toBe(1);
      expect(wrapper.find('ContactInfoFields').length).toBe(1);
      expect(wrapper.find('input[name="customerRemarks"]').length).toBe(1);
    });
  });

  describe('creating a new NTS-R shipment', () => {
    it('renders expected child components', () => {
      const wrapper = mountMtoShipmentForm({ selectedMoveType: SHIPMENT_OPTIONS.NTSR });
      expect(wrapper.find('MtoShipmentForm').length).toBe(1);
      expect(wrapper.find('DatePickerInput').length).toBe(1);
      expect(wrapper.find('AddressFields').length).toBe(0);
      expect(wrapper.find('ContactInfoFields').length).toBe(1);
      expect(wrapper.find('input[name="customerRemarks"]').length).toBe(1);
    });
  });
});
