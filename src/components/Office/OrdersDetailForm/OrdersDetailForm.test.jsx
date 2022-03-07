/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import { render, screen } from '@testing-library/react';
import { Formik } from 'formik';

import OrdersDetailForm from './OrdersDetailForm';

import { DEPARTMENT_INDICATOR_OPTIONS } from 'constants/departmentIndicators';
import { dropdownInputOptions } from 'shared/formatters';
import { ORDERS_TYPE_OPTIONS, ORDERS_TYPE_DETAILS_OPTIONS } from 'constants/orders';

const dutyStation = {
  address: {
    city: 'Scott Air Force Base',
    id: '9f8b0fad-afe1-4a44-bb28-296a335c1141',
    postalCode: '62225',
    state: 'IL',
    streetAddress1: '',
  },
  address_id: '9f8b0fad-afe1-4a44-bb28-296a335c1141',
  affiliation: 'AIR_FORCE',
  created_at: '2018-10-04T22:54:46.589Z',
  id: '071f6286-8255-4e35-b8ac-0e7fe1d10aa4',
  name: 'Scott AFB',
  updated_at: '2018-10-04T22:54:46.589Z',
};

const initialValues = {
  currentDutyLocation: dutyStation,
  newDutyLocation: dutyStation,
  dateIssued: '2020-03-08',
  reportByDate: '2020-04-01',
  departmentIndicator: 'NAVY_AND_MARINES',
  ordersNumber: '999999999',
  ordersType: 'PERMANENT_CHANGE_OF_STATION',
  ordersTypeDetail: 'HHG_PERMITTED',
  tac: 'Tac',
  sac: 'Sac',
  ordersAcknowledgement: true,
};

const deptOptions = dropdownInputOptions(DEPARTMENT_INDICATOR_OPTIONS);
const ordersTypeOptions = dropdownInputOptions(ORDERS_TYPE_OPTIONS);
const ordersTypeDetailOptions = dropdownInputOptions(ORDERS_TYPE_DETAILS_OPTIONS);
const defaultProps = {
  deptIndicatorOptions: deptOptions,
  ordersTypeOptions,
  ordersTypeDetailOptions,
  showOrdersAcknowledgement: true,
  validateTac: jest.fn,
  ordersType: 'PERMANENT_CHANGE_OF_STATION',
  setFieldValue: jest.fn,
};

function renderOrdersDetailForm(props) {
  render(
    <Formik initialValues={initialValues}>
      <form>
        <OrdersDetailForm {...defaultProps} {...props} />
      </form>
    </Formik>,
  );
}

describe('OrdersDetailForm', () => {
  it('renders the Form', async () => {
    renderOrdersDetailForm();
    expect(await screen.findByLabelText('Current duty location')).toBeInTheDocument();

    // hidden fields are default visible
    expect(screen.getByLabelText('Department indicator')).toBeInTheDocument();
    expect(screen.getByLabelText('Orders number')).toBeInTheDocument();
    expect(screen.getByLabelText('Orders type detail')).toBeInTheDocument();
    expect(screen.queryAllByLabelText('TAC').length).toBe(2);
    expect(screen.queryAllByLabelText('SAC').length).toBe(2);
    // expect(screen.getByLabelText('SAC')).toBeInTheDocument();
    expect(screen.getByLabelText('I have read the new orders')).toBeInTheDocument();
  });

  it('accepts deptIndicatorOptions prop', async () => {
    renderOrdersDetailForm();
    expect(await screen.findByLabelText('Department indicator')).toBeInTheDocument();
  });

  it('accepts ordersTypeOptions prop', async () => {
    renderOrdersDetailForm();
    expect(await screen.findByLabelText('Orders type')).toBeInTheDocument();
  });

  it('accepts ordersTypeDetailOptions prop', async () => {
    renderOrdersDetailForm();
    expect(await screen.findByLabelText('Orders type detail')).toBeInTheDocument();
  });

  it('accepts showOrdersAcknowledgement prop', async () => {
    renderOrdersDetailForm();
    expect(await screen.findByLabelText('I have read the new orders')).toBeInTheDocument();
  });

  it('shows the tac warning', async () => {
    renderOrdersDetailForm({ hhgTacWarning: 'Test warning' });
    expect(await screen.findByText('Test warning')).toBeInTheDocument();
  });

  it('hides hideable fields', async () => {
    renderOrdersDetailForm({
      showDepartmentIndicator: false,
      showOrdersNumber: false,
      showOrdersTypeDetail: false,
      showHHGTac: false,
      showHHGSac: false,
      showNTSTac: false,
      showNTSSac: false,
      showOrdersAcknowledgement: false,
    });

    // fields are visible
    expect(await screen.findByLabelText('Current duty location')).toBeInTheDocument();

    // fields are hidden
    expect(screen.queryByLabelText('Department indicator')).not.toBeInTheDocument();
    expect(screen.queryByLabelText('Orders number')).not.toBeInTheDocument();
    expect(screen.queryByLabelText('Orders type detail')).not.toBeInTheDocument();
    expect(screen.queryByLabelText('TAC')).not.toBeInTheDocument();
    expect(screen.queryByLabelText('SAC')).not.toBeInTheDocument();
    expect(screen.queryByLabelText('I have read the new orders')).not.toBeInTheDocument();
  });

  it('has the right labels for a retiree', async () => {
    renderOrdersDetailForm({
      showDepartmentIndicator: false,
      showOrdersNumber: false,
      showOrdersTypeDetail: false,
      showHHGTac: false,
      showHHGSac: false,
      showNTSTac: false,
      showNTSSac: false,
      showOrdersAcknowledgement: false,
      ordersType: 'RETIREMENT',
    });

    // correct labels are visible
    expect(await screen.findByLabelText('Date of retirement')).toBeInTheDocument();
    expect(await screen.findByLabelText('HOR, HOS or PLEAD')).toBeInTheDocument();
  });

  it('has the right labels for a separatee', async () => {
    renderOrdersDetailForm({
      showDepartmentIndicator: false,
      showOrdersNumber: false,
      showOrdersTypeDetail: false,
      showHHGTac: false,
      showHHGSac: false,
      showNTSTac: false,
      showNTSSac: false,
      showOrdersAcknowledgement: false,
      ordersType: 'SEPARATION',
    });

    // correct labels are visible
    expect(await screen.findByLabelText('Date of separation')).toBeInTheDocument();
    expect(await screen.findByLabelText('HOR, HOS or PLEAD')).toBeInTheDocument();
  });
});
