/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import { screen } from '@testing-library/react';
import { Formik } from 'formik';
import userEvent from '@testing-library/user-event';

import OrdersDetailForm from './OrdersDetailForm';

import { DEPARTMENT_INDICATOR_OPTIONS } from 'constants/departmentIndicators';
import { dropdownInputOptions } from 'utils/formatters';
import { ORDERS_TYPE_OPTIONS, ORDERS_TYPE_DETAILS_OPTIONS, ORDERS_PAY_GRADE_TYPE } from 'constants/orders';
import { renderWithRouter } from 'testUtils';

const dutyLocation = {
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
  currentDutyLocation: dutyLocation,
  newDutyLocation: dutyLocation,
  dateIssued: '2020-03-08',
  reportByDate: '2020-04-01',
  departmentIndicator: 'NAVY_AND_MARINES',
  ordersNumber: '999999999',
  ordersType: 'PERMANENT_CHANGE_OF_STATION',
  ordersTypeDetail: 'HHG_PERMITTED',
  tac: 'Tac',
  sac: 'Sac',
  ordersAcknowledgement: true,
  dependentsAuthorized: true,
};

const deptOptions = dropdownInputOptions(DEPARTMENT_INDICATOR_OPTIONS);
const ordersTypeOptions = dropdownInputOptions(ORDERS_TYPE_OPTIONS);
const ordersTypeDetailOptions = dropdownInputOptions(ORDERS_TYPE_DETAILS_OPTIONS);
const payGradeOptions = dropdownInputOptions(ORDERS_PAY_GRADE_TYPE);

const defaultProps = {
  deptIndicatorOptions: deptOptions,
  ordersTypeOptions,
  ordersTypeDetailOptions,
  showOrdersAcknowledgement: true,
  validateTac: jest.fn,
  ordersType: 'PERMANENT_CHANGE_OF_STATION',
  setFieldValue: jest.fn,
  payGradeOptions,
  hhgLongLineOfAccounting: 'Long line of accounting is present',
  ntsLongLineOfAccounting: 'NTS Long line of accounting is present',
};

function renderOrdersDetailForm(props) {
  renderWithRouter(
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
    expect(await screen.findByText('Current duty location')).toBeInTheDocument();

    // hidden fields are default visible
    expect(screen.getByLabelText('Department indicator *')).toBeInTheDocument();
    expect(screen.getByLabelText('Orders number *')).toBeInTheDocument();
    expect(screen.getByLabelText('Orders type detail *')).toBeInTheDocument();
    expect(screen.getByLabelText('TAC *')).toBeInTheDocument();
    expect(screen.queryAllByLabelText('SAC').length).toBe(2);
    expect(screen.getByLabelText('I have read the new orders')).toBeInTheDocument();
    expect(screen.getByTestId('reqAsteriskMsg')).toBeInTheDocument();
    expect(screen.getByTestId('reqAsteriskMsg').textContent).toContain('Fields marked with * are required.');
  });

  it('renders the Form disabled with all information if flag is passed', async () => {
    renderOrdersDetailForm({ formIsDisabled: true });
    const currentDutyLocationInput = screen.getByLabelText(/Current duty location/);
    expect(currentDutyLocationInput).toBeInTheDocument();
    expect(currentDutyLocationInput).toBeDisabled();
    const newDutyLocationInput = screen.getByLabelText(/New duty location/);
    expect(newDutyLocationInput).toBeInTheDocument();
    expect(newDutyLocationInput).toBeDisabled();
    const payGradeInput = screen.getByLabelText(/Pay grade/);
    expect(payGradeInput).toBeInTheDocument();
    expect(payGradeInput).toBeDisabled();
    const dateIssuedInput = screen.getByLabelText(/Date issued/);
    expect(dateIssuedInput).toBeInTheDocument();
    expect(dateIssuedInput).toBeDisabled();
    const reportByDateInput = screen.getByLabelText(/Report by date/);
    expect(reportByDateInput).toBeInTheDocument();
    expect(reportByDateInput).toBeDisabled();
    const departmentIndicatorInput = screen.getByLabelText(/Department indicator/);
    expect(departmentIndicatorInput).toBeInTheDocument();
    expect(departmentIndicatorInput).toBeDisabled();
    const ordersNumberInput = screen.getByLabelText(/Orders number/);
    expect(ordersNumberInput).toBeInTheDocument();
    expect(ordersNumberInput).toBeDisabled();
    const ordersTypeInputs = screen.getAllByLabelText(/Orders type/);
    ordersTypeInputs.forEach((input) => {
      expect(input).toBeDisabled();
    });
    const ordersTypeDetailInput = screen.getByLabelText(/Orders type detail/);
    expect(ordersTypeDetailInput).toBeInTheDocument();
    expect(ordersTypeDetailInput).toBeDisabled();
    const dependentsAuthorizedInput = screen.getByLabelText('Dependents authorized');
    expect(dependentsAuthorizedInput).toBeInTheDocument();
    expect(dependentsAuthorizedInput).toBeDisabled();
    const tacInput = screen.getByLabelText('TAC');
    expect(tacInput).toBeInTheDocument();
    expect(tacInput).toBeDisabled();
    const tacInputRequired = screen.getByLabelText('TAC *');
    expect(tacInputRequired).toBeInTheDocument();
    expect(tacInputRequired).toBeDisabled();
    const sacInputs = screen.queryAllByLabelText('SAC');
    const tacInputs = screen.queryAllByLabelText('TAC');
    expect(tacInputs.length).toBe(1);
    expect(sacInputs.length).toBe(2);
    expect(tacInputs[0]).toBeDisabled();
    expect(sacInputs[0]).toBeDisabled();
    expect(sacInputs[1]).toBeDisabled();
  });

  it('accepts deptIndicatorOptions prop', async () => {
    renderOrdersDetailForm();
    expect(await screen.findByLabelText('Department indicator *')).toBeInTheDocument();
  });

  it('accepts ordersTypeOptions prop', async () => {
    renderOrdersDetailForm();
    expect(await screen.findByLabelText('Orders type *')).toBeInTheDocument();
  });

  it('accepts ordersTypeDetailOptions prop', async () => {
    renderOrdersDetailForm();
    expect(await screen.findByLabelText('Orders type detail *')).toBeInTheDocument();
  });

  it('accepts showOrdersAcknowledgement prop', async () => {
    renderOrdersDetailForm();
    expect(await screen.findByLabelText('I have read the new orders')).toBeInTheDocument();
  });

  it('accepts hhg longLineOfAccounting prop', async () => {
    renderOrdersDetailForm();
    const loaTextField = screen.getByTestId('hhgLoaTextField');
    expect(loaTextField).toHaveValue('Long line of accounting is present');
  });

  it('accepts nts longLineOfAccounting prop', async () => {
    renderOrdersDetailForm();
    const loaTextField = screen.getByTestId('ntsLoaTextField');
    expect(loaTextField).toHaveValue('NTS Long line of accounting is present');
  });

  it('shows the tac warning', async () => {
    renderOrdersDetailForm({ hhgTacWarning: 'Test warning' });
    expect(await screen.findByText('Test warning')).toBeInTheDocument();
  });

  it('shows the loa warning', async () => {
    renderOrdersDetailForm({ hhgLoaWarning: 'Test LOA warning' });
    expect(await screen.findByText('Test LOA warning')).toBeInTheDocument();
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
      showLoa: false,
    });

    // fields are visible
    expect(await screen.findByLabelText('Current duty location *')).toBeInTheDocument();

    // fields are hidden
    expect(screen.queryByLabelText('Department indicator *')).not.toBeInTheDocument();
    expect(screen.queryByLabelText('Orders number *')).not.toBeInTheDocument();
    expect(screen.queryByLabelText('Orders type detail *')).not.toBeInTheDocument();
    expect(screen.queryByLabelText('TAC *')).not.toBeInTheDocument();
    expect(screen.queryByLabelText('SAC')).not.toBeInTheDocument();
    expect(screen.queryByLabelText('I have read the new orders')).not.toBeInTheDocument();
    expect(screen.queryByLabelText('Line of Accounting Preview')).not.toBeInTheDocument();
    expect(screen.queryByLabelText('LOA')).not.toBeInTheDocument();
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
    expect(await screen.findByLabelText('Date of retirement *')).toBeInTheDocument();
    expect(await screen.findByLabelText('HOR, HOS or PLEAD *')).toBeInTheDocument();
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
    expect(await screen.findByLabelText('Date of separation *')).toBeInTheDocument();
    expect(await screen.findByLabelText('HOR, HOS or PLEAD *')).toBeInTheDocument();
  });

  it('has orders type dropdown disabled if safety move', async () => {
    renderOrdersDetailForm({
      showDepartmentIndicator: false,
      showOrdersNumber: false,
      showOrdersTypeDetail: false,
      showHHGTac: false,
      showHHGSac: false,
      showNTSTac: false,
      showNTSSac: false,
      showOrdersAcknowledgement: false,
      ordersType: 'SAFETY',
    });

    // correct labels are visible
    expect(await screen.findByLabelText('Orders type *')).toBeDisabled();
  });
  it('has orders type dropdown disabled if bluebark move', async () => {
    renderOrdersDetailForm({
      showDepartmentIndicator: false,
      showOrdersNumber: false,
      showOrdersTypeDetail: false,
      showHHGTac: false,
      showHHGSac: false,
      showNTSTac: false,
      showNTSSac: false,
      showOrdersAcknowledgement: false,
      ordersType: 'BLUEBARK',
    });

    // correct labels are visible
    expect(await screen.findByLabelText('Orders type *')).toBeDisabled();
  });

  it('renders dependents authorized checkbox field', async () => {
    renderOrdersDetailForm();
    expect(await screen.findByTestId('dependentsAuthorizedInput')).toBeInTheDocument();
  });

  it('allows typing more than 4 characters into a SAC field', async () => {
    renderOrdersDetailForm();

    // there are two SAC fields (HHG SAC and NTS SAC)
    const sacInputs = screen.getAllByLabelText('SAC');
    expect(sacInputs.length).toBeGreaterThanOrEqual(1);

    const firstSacInput = sacInputs[0];
    await userEvent.type(firstSacInput, 'ABCDE');
    // Sac is already in the initial values, so we can confirm we can append to that
    expect(firstSacInput).toHaveValue('SacABCDE');

    // NTS SAC is not in the initial values, so we can check for exactly what we put in
    const secondSacInput = sacInputs[1];
    await userEvent.type(secondSacInput, 'FGHIJ123');
    expect(secondSacInput).toHaveValue('FGHIJ123');
  });
});
