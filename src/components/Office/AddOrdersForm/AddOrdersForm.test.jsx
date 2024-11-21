import React from 'react';
import { render, waitFor, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { Provider } from 'react-redux';

import AddOrdersForm from './AddOrdersForm';

import { dropdownInputOptions } from 'utils/formatters';
import { ORDERS_TYPE_OPTIONS } from 'constants/orders';
import { configureStore } from 'shared/store';
import { isBooleanFlagEnabled } from 'utils/featureFlags';

jest.setTimeout(60000);

jest.mock('components/LocationSearchBox/api', () => ({
  ShowAddress: jest.fn().mockImplementation(() =>
    Promise.resolve({
      city: 'Luke AFB',
      country: 'United States',
      id: 'fa51dab0-4553-4732-b843-1f33407f77bc',
      postalCode: '85309',
      state: 'AZ',
      streetAddress1: 'n/a',
      isOconus: true,
    }),
  ),
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
      {
        address: {
          city: 'Glendale Luke AFB',
          country: 'United States',
          id: 'fa51dab0-4553-4732-b843-1f33407f77bc',
          postalCode: '85309',
          state: 'AZ',
          streetAddress1: 'n/a',
          isOconus: true,
        },
        address_id: '25be4d12-fe93-47f1-bbec-1db386dfa67f',
        affiliation: 'AIR_FORCE',
        created_at: '2021-02-11T16:48:04.117Z',
        id: 'a8d6b33c-8370-4e92-8df2-356b8c9d0c1a',
        name: 'Luke AFB',
        updated_at: '2021-02-11T16:48:04.117Z',
      },
    ]),
  ),
}));

jest.mock('utils/featureFlags', () => ({
  ...jest.requireActual('utils/featureFlags'),
  isBooleanFlagEnabled: jest.fn().mockImplementation(() => Promise.resolve(false)),
}));

const mockStore = configureStore({});
const initialValues = {
  ordersType: '',
  issueDate: '',
  reportByDate: '',
  hasDependents: '',
  newDutyLocation: '',
  grade: '',
  originDutyLocation: '',
  accompaniedTour: '',
  dependentsUnderTwelve: '',
  dependentsTwelveAndOver: '',
};
const testProps = {
  initialValues,
  ordersTypeOptions: dropdownInputOptions(ORDERS_TYPE_OPTIONS),
  onSubmit: jest.fn(),
  onBack: jest.fn(),
};

describe('CreateMoveCustomerInfo Component', () => {
  it('renders the form inputs', async () => {
    render(
      <Provider store={mockStore.store}>
        <AddOrdersForm {...testProps} />
      </Provider>,
    );

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
    const { getByRole, findAllByRole, getByLabelText } = render(
      <Provider store={mockStore.store}>
        <AddOrdersForm {...testProps} />
      </Provider>,
    );
    await userEvent.click(getByLabelText('Orders type'));
    await userEvent.click(getByLabelText('Orders date'));
    await userEvent.click(getByLabelText('Report by date'));
    await userEvent.click(getByLabelText('Current duty location'));
    await userEvent.click(getByLabelText('New duty location'));
    await userEvent.click(getByLabelText('Pay grade'));

    const submitBtn = getByRole('button', { name: 'Next' });
    await userEvent.click(submitBtn);

    const alerts = await findAllByRole('alert');
    expect(alerts.length).toBe(5);

    alerts.forEach((alert) => {
      expect(alert).toHaveTextContent('Required');
    });

    expect(testProps.onSubmit).not.toHaveBeenCalled();
  });
});

describe('AddOrdersForm - OCONUS and Accompanied Tour Test', () => {
  it('submits the form with OCONUS values and accompanied tour selection', async () => {
    isBooleanFlagEnabled.mockResolvedValue(true);

    render(
      <Provider store={mockStore.store}>
        <AddOrdersForm {...testProps} />
      </Provider>,
    );

    await userEvent.selectOptions(await screen.findByLabelText(/Orders type/), 'PERMANENT_CHANGE_OF_STATION');
    await userEvent.type(screen.getByLabelText(/Orders date/), '08 Nov 2020');
    await userEvent.type(screen.getByLabelText(/Report by date/), '26 Nov 2020');
    await userEvent.click(screen.getByLabelText('No'));
    await userEvent.selectOptions(screen.getByLabelText(/Pay grade/), ['E_5']);

    await userEvent.type(screen.getByLabelText(/Current duty location/), 'AFB');
    await userEvent.click(await screen.findByText(/Elmendorf/));

    await userEvent.type(screen.getByLabelText(/New duty location/), 'AFB');
    await userEvent.click(await screen.findByText(/Luke/));

    await userEvent.click(screen.getByTestId('hasDependentsYes'));
    await userEvent.click(screen.getByTestId('isAnAccompaniedTourYes'));
    await userEvent.type(screen.getByTestId('dependentsUnderTwelve'), '2');
    await userEvent.type(screen.getByTestId('dependentsTwelveAndOver'), '1');

    const nextBtn = screen.getByRole('button', { name: 'Next' });
    await userEvent.click(nextBtn);

    await waitFor(() => {
      expect(testProps.onSubmit).toHaveBeenCalled();
    });
  });
});
