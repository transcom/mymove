import React from 'react';
import { shallow } from 'enzyme';
import { Form } from 'components/form/Form';
import { OrdersDetailForm } from './OrdersDetailForm';

const mockHandleReset = jest.fn();
const mockHandleSubmit = jest.fn();
// mock out formik hook as we are not testing formik
// needs to be before first describe
jest.mock('formik', () => {
  return {
    ...jest.requireActual('formik'),
    useFormikContext: () => ({
      errors: { sampleField: 'Required' },
      touched: { sampleField: true },
      handleReset: mockHandleReset,
      handleSubmit: mockHandleSubmit,
    }),
  };
});

describe('OrdersDetailForm', () => {
  const deptOptions = [['key', 'value']];
  const ordersTypeOptions = [['key', 'value']];
  const ordersTypeDetailOptions = [['key', 'value']];
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
    dateIssued: '08 Mar 2020',
    reportByDate: '01 Apr 2020',
    departmentIndicator: 'NAVY_AND_MARINES',
    ordersNumber: '999999999',
    ordersType: 'PERMANENT_CHANGE_OF_STATION',
    ordersTypeDetail: 'HHG_PERMITTED',
    tac: 'Tac',
    sac: 'Sac',
  };
  const wrapper = shallow(
    <OrdersDetailForm
      onSubmit={mockHandleSubmit}
      onReset={mockHandleReset}
      deptIndicatorOptions={deptOptions}
      ordersTypeOptions={ordersTypeOptions}
      ordersTypeDetailOptions={ordersTypeDetailOptions}
      initialValues={initialValues}
    />,
  );

  it('renders the Form', () => {
    expect(wrapper.find(Form).length).toBe(1);
  });

  it('accepts onSubmit method', () => {
    expect(wrapper.prop('onSubmit')).toBe(mockHandleSubmit);
  });

  it('accepts onReset method', () => {
    expect(wrapper.prop('onReset')).toBe(mockHandleReset);
  });

  it('accepts initialValue prop', () => {
    expect(wrapper.prop('initialValues')).toBe(initialValues);
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

  it('call submit handler', () => {
    wrapper.simulate('submit');
    expect(mockHandleSubmit).toHaveBeenCalled();
    expect(mockHandleReset).not.toHaveBeenCalled();
  });

  it('call reset handler', () => {
    wrapper.simulate('reset');
    expect(mockHandleSubmit).not.toHaveBeenCalled();
    expect(mockHandleReset).toHaveBeenCalled();
  });

  afterEach(jest.resetAllMocks);
});
