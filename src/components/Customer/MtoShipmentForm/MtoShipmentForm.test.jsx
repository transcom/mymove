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
    const createHhgWrapper = mountMtoShipmentForm({ selectedMoveType: SHIPMENT_OPTIONS.HHG });

    it('renders expected child components', () => {
      expect(createHhgWrapper.find('MtoShipmentForm').length).toBe(1);
      expect(createHhgWrapper.find('DatePickerInput').length).toBe(2);
      expect(createHhgWrapper.find('AddressFields').length).toBe(1);
      expect(createHhgWrapper.find('ContactInfoFields').length).toBe(2);
      expect(createHhgWrapper.find('Field[name="customerRemarks"]').length).toBe(1);
    });

    it('does not render special NTS What to expect section', () => {
      expect(createHhgWrapper.find('div[data-testid="nts-what-to-expect"]').length).toBe(0);
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
      expect(wrapper.find('Field[name="customerRemarks"]').length).toBe(1);
      expect(wrapper.find('Field[name="customerRemarks"]').prop('value')).toEqual(mockMtoShipment.customerRemarks);
      expect(wrapper.find('Field[name="delivery.address.street_address_1"]').prop('value')).toContain(
        mockMtoShipment.destinationAddress.street_address_1,
      );
    });
  });

  describe('creating a new NTS shipment', () => {
    const createNtsWrapper = mountMtoShipmentForm({ selectedMoveType: SHIPMENT_OPTIONS.NTS });

    it('renders expected child components', () => {
      expect(createNtsWrapper.find('MtoShipmentForm').length).toBe(1);
      expect(createNtsWrapper.find('DatePickerInput').length).toBe(1);
      expect(createNtsWrapper.find('AddressFields').length).toBe(1);
      expect(createNtsWrapper.find('ContactInfoFields').length).toBe(1);
      expect(createNtsWrapper.find('Field[name="customerRemarks"]').length).toBe(1);
    });

    it('renders special NTS What to expect section', () => {
      expect(createNtsWrapper.find('div[data-testid="nts-what-to-expect"]').length).toBe(1);
    });
  });

  describe('creating a new NTS-R shipment', () => {
    const createNtsrWrapper = mountMtoShipmentForm({ selectedMoveType: SHIPMENT_OPTIONS.NTSR });

    it('renders expected child components', () => {
      expect(createNtsrWrapper.find('MtoShipmentForm').length).toBe(1);
      expect(createNtsrWrapper.find('DatePickerInput').length).toBe(1);
      expect(createNtsrWrapper.find('AddressFields').length).toBe(0);
      expect(createNtsrWrapper.find('ContactInfoFields').length).toBe(1);
      expect(createNtsrWrapper.find('Field[name="customerRemarks"]').length).toBe(1);
    });

    it('does not render special NTS What to expect section', () => {
      expect(createNtsrWrapper.find('div[data-testid="nts-what-to-expect"]').length).toBe(0);
    });
  });
});
