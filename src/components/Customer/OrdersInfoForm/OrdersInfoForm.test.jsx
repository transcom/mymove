/* react/jsx-props-no-spreading */
import React from 'react';
import { mount } from 'enzyme';
import { Formik } from 'formik';

import OrdersInfoForm from './OrdersInfoForm';

const testProps = {
  ordersTypeOptions: [
    { key: 'PERMANENT_CHANGE_OF_STATION', value: 'Permanent Change Of Station (PCS)' },
    { key: 'RETIREMENT', value: 'Retirement' },
    { key: 'SEPARATION', value: 'Separation' },
  ],
};

describe('OrdersInfoForm component', () => {
  describe('with no initial values', () => {
    const wrapper = mount(
      <Formik>
        <OrdersInfoForm {...testProps} />
      </Formik>,
    );

    it('renders without errors', () => {
      expect(wrapper.exists()).toBe(true);
    });

    it('renders the expected form inputs', () => {
      expect(wrapper.find('DropdownInput[name="orders_type"]').length).toBe(1);
      expect(wrapper.find('DatePickerInput[name="issue_date"]').length).toBe(1);
      expect(wrapper.find('DatePickerInput[name="report_by_date"]').length).toBe(1);
      expect(wrapper.find('input[name="has_dependents"][value="yes"]').length).toBe(1);
      expect(wrapper.find('input[name="has_dependents"][value="no"]').length).toBe(1);
      expect(wrapper.find('DutyStationInput[name="new_duty_station"]').length).toBe(1);
    });

    it('renders each option for orders type', () => {
      const ordersTypeDropdown = wrapper.find('DropdownInput[name="orders_type"]');
      const expectedOptions = [
        { key: 'PERMANENT_CHANGE_OF_STATION', value: 'Permanent Change Of Station (PCS)' },
        { key: 'RETIREMENT', value: 'Retirement' },
        { key: 'SEPARATION', value: 'Separation' },
      ];
      expect(ordersTypeDropdown.prop('options')).toEqual(expectedOptions);
    });
  });

  describe('with initial values', () => {
    const testInitialValues = {
      orders_type: 'PERMANENT_CHANGE_OF_STATION',
      issue_date: '2020-11-08',
      report_by_date: '2020-11-26',
      has_dependents: 'no',
      new_duty_station: {
        address: {
          city: 'Des Moines',
          country: 'US',
          id: 'a4b30b99-4e82-48a6-b736-01662b499d6a',
          postal_code: '50309',
          state: 'IA',
          street_address_1: '987 Other Avenue',
          street_address_2: 'P.O. Box 1234',
          street_address_3: 'c/o Another Person',
        },
        address_id: 'a4b30b99-4e82-48a6-b736-01662b499d6a',
        affiliation: 'AIR_FORCE',
        created_at: '2020-10-19T17:01:16.114Z',
        id: 'f9299768-16d2-4a13-ae39-7087a58b1f62',
        name: 'Yuma AFB',
        updated_at: '2020-10-19T17:01:16.114Z',
      },
    };
    const wrapper = mount(
      <Formik initialValues={testInitialValues}>
        <OrdersInfoForm {...testProps} />
      </Formik>,
    );

    it('renders without errors', () => {
      expect(wrapper.exists()).toBe(true);
    });

    it('pre-fills the inputs', () => {
      expect(wrapper.find('select[name="orders_type"]').prop('value')).toBe(testInitialValues.orders_type);
      expect(wrapper.find('input[name="issue_date"]').prop('value')).toBe('08 Nov 2020');
      expect(wrapper.find('input[name="report_by_date"]').prop('value')).toBe('26 Nov 2020');
      expect(wrapper.find('input[name="has_dependents"][value="yes"]').prop('checked')).toBe(false);
      expect(wrapper.find('input[name="has_dependents"][value="no"]').prop('checked')).toBe(true);
      expect(wrapper.find('input[name="new_duty_station"]').prop('value')).toBe(
        testInitialValues.new_duty_station.name,
      );
    });
  });
});
