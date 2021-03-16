/* eslint-disable  react/jsx-props-no-spreading */
import React from 'react';
import { mount } from 'enzyme';
import { Formik } from 'formik';

import OrdersDetailForm from './OrdersDetailForm';

import { DEPARTMENT_INDICATOR_OPTIONS } from 'constants/departmentIndicators';
import { dropdownInputOptions } from 'shared/formatters';
import { ORDERS_TYPE_OPTIONS, ORDERS_TYPE_DETAILS_OPTIONS } from 'constants/orders';

const dutyStation = {
  address: {
    city: 'Scott Air Force Base',
    id: '9f8b0fad-afe1-4a44-bb28-296a335c1141',
    postal_code: '62225',
    state: 'IL',
    street_address_1: '',
  },
  address_id: '9f8b0fad-afe1-4a44-bb28-296a335c1141',
  affiliation: 'AIR_FORCE',
  created_at: '2018-10-04T22:54:46.589Z',
  id: '071f6286-8255-4e35-b8ac-0e7fe1d10aa4',
  name: 'Scott AFB',
  updated_at: '2018-10-04T22:54:46.589Z',
};

const initialValues = {
  currentDutyStation: dutyStation,
  newDutyStation: dutyStation,
  dateIssued: '2020-03-08',
  reportByDate: '2020-04-01',
  departmentIndicator: 'NAVY_AND_MARINES',
  ordersNumber: '999999999',
  ordersType: 'PERMANENT_CHANGE_OF_STATION',
  ordersTypeDetail: 'HHG_PERMITTED',
  tac: 'Tac',
  sac: 'Sac',
};

const deptOptions = dropdownInputOptions(DEPARTMENT_INDICATOR_OPTIONS);
const ordersTypeOptions = dropdownInputOptions(ORDERS_TYPE_OPTIONS);
const ordersTypeDetailOptions = dropdownInputOptions(ORDERS_TYPE_DETAILS_OPTIONS);
const defaultProps = {
  deptIndicatorOptions: deptOptions,
  ordersTypeOptions,
  ordersTypeDetailOptions,
  validateTac: jest.fn,
};

function mountOrdersDetailForm(props) {
  return mount(
    <Formik initialValues={initialValues}>
      <form>
        <OrdersDetailForm {...defaultProps} {...props} />
      </form>
    </Formik>,
  );
}

describe('OrdersDetailForm', () => {
  const wrapper = mountOrdersDetailForm();

  it('renders the Form', () => {
    expect(wrapper.find(OrdersDetailForm).exists()).toBe(true);
  });

  it('accepts deptIndicatorOptions prop', () => {
    expect(wrapper.find('DropdownInput[name="departmentIndicator"]').prop('options')).toBe(deptOptions);
  });

  it('accepts ordersTypeOptions prop', () => {
    expect(wrapper.find('DropdownInput[name="ordersType"]').prop('options')).toBe(ordersTypeOptions);
  });

  it('accepts ordersTypeDetailOptions prop', () => {
    expect(wrapper.find('DropdownInput[name="ordersTypeDetail"]').prop('options')).toBe(ordersTypeDetailOptions);
  });

  it('populates initial field values', () => {
    /*
    expect(wrapper.find('[name="originDutyStation"]').value).toBe(dutyStation);
    expect(wrapper.find('[name="newDutyStation"]').prop('value')).toBe(dutyStation);
    expect(wrapper.find('[name="issueDate"]').prop('value')).toBe('08 Mar 2020');
    expect(wrapper.find('[name="reportByDate"]').prop('value')).toBe('01 Apr 2020');
    expect(wrapper.find('[name="departmentIndicator"]').prop('value')).toBe('NAVY_AND_MARINES');
    expect(wrapper.find('[name="ordersNumber"]').prop('value')).toBe('999999999');
    expect(wrapper.find('[name="ordersType"]').prop('value')).toBe('PERMANENT_CHANGE_OF_STATION');
    expect(wrapper.find('[name="ordersTypeDetail"]').prop('value')).toBe('HHG_PERMITTED');
    expect(wrapper.find('[name="tac"]').prop('value')).toBe('Tac');
    expect(wrapper.find('[name="sac"]').prop('value')).toBe('Sac');
    */
  });

  it('shows the tac warning', () => {
    const tacWarning = 'You have been warned';
    const wrapperWarn = mountOrdersDetailForm({ tacWarning });

    expect(wrapperWarn.find('[data-testid="textInputWarning"]').exists()).toBe(true);
    expect(wrapperWarn.find('[data-testid="textInputWarning"]').text()).toEqual(tacWarning);
  });
});
