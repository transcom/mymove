import React from 'react';
import { render, waitFor, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import AddOrdersForm from './AddOrdersForm';

import { dropdownInputOptions } from 'utils/formatters';
import { ORDERS_PAY_GRADE_OPTIONS } from 'constants/orders';

import { isBooleanFlagEnabled } from 'utils/featureFlags';

jest.mock('components/LocationSearchBox/api', () => ({
  SearchDutyLocations: jest.fn().mockImplementation(() =>
    Promise.resolve([
      {
        address: {
          city: '',
          id: '00000000-0000-0000-0000-000000000000',
          postalCode: '',
          state: '',
          streetAddress1: '',
        },
        address_id: '46c4640b-c35e-4293-a2f1-36c7b629f903',
        affiliation: 'AIR_FORCE',
        created_at: '2021-02-11T16:48:04.117Z',
        id: '93f0755f-6f35-478b-9a75-35a69211da1c',
        name: 'Altus AFB',
        updated_at: '2021-02-11T16:48:04.117Z',
      },
      {
        address: {
          city: 'Elmendorf AFB',
          country: 'US',
          id: 'fa51dab0-4553-4732-b843-1f33407f11bc',
          postalCode: '78112',
          state: 'AK',
          streetAddress1: 'n/a',
          isOconus: true,
        },
        address_id: 'fa51dab0-4553-4732-b843-1f33407f11bc',
        affiliation: 'AIR_FORCE',
        created_at: '2021-02-11T16:48:04.117Z',
        id: 'a8d6b33c-8370-4e92-8df2-356b8c9d0c1a',
        name: 'Elmendorf AFB',
        updated_at: '2021-02-11T16:48:04.117Z',
      },
    ]),
  ),
}));

jest.mock('utils/featureFlags', () => ({
  ...jest.requireActual('utils/featureFlags'),
  isBooleanFlagEnabled: jest.fn().mockImplementation(() => Promise.resolve(false)),
}));

describe('CreateMoveCustomerInfo Component', () => {
  const initialValues = {
    ordersType: '',
    issueDate: '',
    reportByDate: '',
    hasDependents: '',
    newDutyLocation: '',
    grade: '',
    currentDutyLocation: {},
    newDutyLocation: {},
  };
  const testProps = {
    initialValues,
    ordersTypeOptions: dropdownInputOptions(ORDERS_PAY_GRADE_OPTIONS),
    onSubmit: jest.fn(),
    onBack: jest.fn(),
  };

  it('renders the form inputs', async () => {
    render(<AddOrdersForm {...testProps} />);

    await waitFor(() => {
      expect(screen.getByText('Tell us about the orders')).toBeInTheDocument();
      expect(screen.getByLabelText('Orders type')).toBeInTheDocument();
      expect(screen.getByLabelText('Orders date')).toBeInTheDocument();
      expect(screen.getByLabelText('Report by date')).toBeInTheDocument();
      expect(screen.getByText('Are dependents included in the orders?')).toBeInTheDocument();
      expect(screen.getByTestId('hasDependentsYes')).toBeInTheDocument();
      expect(screen.getByTestId('hasDependentsNo')).toBeInTheDocument();
      expect(screen.getByLabelText('Current duty location')).toBeInTheDocument();
      expect(screen.getByLabelText('New duty location')).toBeInTheDocument();
      expect(screen.getByLabelText('Pay grade')).toBeInTheDocument();
    });
  });

  it('shows an error message if trying to submit an invalid form', async () => {
    const { getByRole, findAllByRole, getByLabelText } = render(<AddOrdersForm {...testProps} />);
    await userEvent.click(getByLabelText('Orders type'));
    await userEvent.click(getByLabelText('Orders date'));
    await userEvent.click(getByLabelText('Report by date'));
    await userEvent.click(getByLabelText('Current duty location'));
    await userEvent.click(getByLabelText('New duty location'));
    await userEvent.click(getByLabelText('Pay grade'));

    const submitBtn = getByRole('button', { name: 'Next' });
    await userEvent.click(submitBtn);

    const alerts = await findAllByRole('alert');
    expect(alerts.length).toBe(4);

    alerts.forEach((alert) => {
      expect(alert).toHaveTextContent('Required');
    });

    expect(testProps.onSubmit).not.toHaveBeenCalled();
  });

  it('only renders dependents age groupings and accompanied tour if dependents are present', async () => {
    isBooleanFlagEnabled.mockImplementation(() => Promise.resolve(true));
    render(<AddOrdersForm {...testProps} />);

    // Select a CONUS current duty location
    await userEvent.type(screen.getByLabelText(/Current duty location/), 'AFB', { delay: 100 });
    const selectedOptionCurrent = await screen.findByText(/Altus/);
    await userEvent.click(selectedOptionCurrent);
    // Select an OCONUS new duty location
    await userEvent.type(screen.getByLabelText(/New duty location/), 'AFB', { delay: 100 });
    const selectedOptionNew = await screen.findByText(/Elmendorf/);
    await userEvent.click(selectedOptionNew);
    // Select that dependents are present
    await userEvent.click(screen.getByTestId('hasDependentsYes'));
    // With one of the duty locations being OCONUS, the number of dependents input boxes should be present
    expect(screen.queryByLabelText(/Number of dependents under the age of 12/)).toBeInTheDocument();
    expect(screen.queryByLabelText(/Number of dependents of the age 12 or over/)).toBeInTheDocument();
    // expect(screen.queryByLabelText(/Is this an accompanied tour?/)).toBeInTheDocument();
    expect(screen.findByTestId('isAnAccompaniedTourYes').toBeInTheDocument());

    // Select that dependents are not present, the fields should go away
    expect(screen.queryByLabelText(/Number of dependents under the age of 12/)).not.toBeInTheDocument();
    expect(screen.queryByLabelText(/Number of dependents of the age 12 or over/)).not.toBeInTheDocument();
    expect(screen.queryByLabelText(/Is this an accompanied tour?/)).not.toBeInTheDocument();
  });
});
