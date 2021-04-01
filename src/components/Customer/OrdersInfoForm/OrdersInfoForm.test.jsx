import React from 'react';
import { render, waitFor, fireEvent } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import selectEvent from 'react-select-event';

import OrdersInfoForm from './OrdersInfoForm';

jest.mock('scenes/ServiceMembers/api.js', () => ({
  ShowAddress: jest.fn().mockImplementation(() => Promise.resolve()),
  SearchDutyStations: jest.fn().mockImplementation(() =>
    Promise.resolve([
      {
        address: {
          city: '',
          id: '00000000-0000-0000-0000-000000000000',
          postal_code: '',
          state: '',
          street_address_1: '',
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
          city: '',
          id: '00000000-0000-0000-0000-000000000000',
          postal_code: '',
          state: '',
          street_address_1: '',
        },
        address_id: '2d7e17f6-1b8a-4727-8949-007c80961a62',
        affiliation: 'AIR_FORCE',
        created_at: '2021-02-11T16:48:04.117Z',
        id: '7d123884-7c1b-4611-92ae-e8d43ca03ad9',
        name: 'Hill AFB',
        updated_at: '2021-02-11T16:48:04.117Z',
      },
      {
        address: {
          city: '',
          id: '00000000-0000-0000-0000-000000000000',
          postal_code: '',
          state: '',
          street_address_1: '',
        },
        address_id: '25be4d12-fe93-47f1-bbec-1db386dfa67f',
        affiliation: 'AIR_FORCE',
        created_at: '2021-02-11T16:48:04.117Z',
        id: 'a8d6b33c-8370-4e92-8df2-356b8c9d0c1a',
        name: 'Luke AFB',
        updated_at: '2021-02-11T16:48:04.117Z',
      },
      {
        address: {
          city: '',
          id: '00000000-0000-0000-0000-000000000000',
          postal_code: '',
          state: '',
          street_address_1: '',
        },
        address_id: '3dbf1fc7-3289-4c6e-90aa-01b530a7c3c3',
        affiliation: 'AIR_FORCE',
        created_at: '2021-02-11T16:48:20.225Z',
        id: 'd01bd2a4-6695-4d69-8f2f-69e88dff58f8',
        name: 'Shaw AFB',
        updated_at: '2021-02-11T16:48:20.225Z',
      },
      {
        address: {
          city: '',
          id: '00000000-0000-0000-0000-000000000000',
          postal_code: '',
          state: '',
          street_address_1: '',
        },
        address_id: '1af8f0f3-f75f-46d3-8dc8-c67c2feeb9f0',
        affiliation: 'AIR_FORCE',
        created_at: '2021-02-11T16:49:14.322Z',
        id: 'b1f9a535-96d4-4cc3-adf1-b76505ce0765',
        name: 'Yuma AFB',
        updated_at: '2021-02-11T16:49:14.322Z',
      },
      {
        address: {
          city: '',
          id: '00000000-0000-0000-0000-000000000000',
          postal_code: '',
          state: '',
          street_address_1: '',
        },
        address_id: 'f2adfebc-7703-4d06-9b49-c6ca8f7968f1',
        affiliation: 'AIR_FORCE',
        created_at: '2021-02-11T16:48:20.225Z',
        id: 'a268b48f-0ad1-4a58-b9d6-6de10fd63d96',
        name: 'Los Angeles AFB',
        updated_at: '2021-02-11T16:48:20.225Z',
      },
      {
        address: {
          city: '',
          id: '00000000-0000-0000-0000-000000000000',
          postal_code: '',
          state: '',
          street_address_1: '',
        },
        address_id: '13eb2cab-cd68-4f43-9532-7a71996d3296',
        affiliation: 'AIR_FORCE',
        created_at: '2021-02-11T16:48:20.225Z',
        id: 'a48fda70-8124-4e90-be0d-bf8119a98717',
        name: 'Wright-Patterson AFB',
        updated_at: '2021-02-11T16:48:20.225Z',
      },
    ]),
  ),
}));

const testProps = {
  onSubmit: jest.fn().mockImplementation(() => Promise.resolve()),
  initialValues: { orders_type: '', issue_date: '', report_by_date: '', has_dependents: '', new_duty_station: {} },
  onBack: jest.fn(),
  ordersTypeOptions: [
    { key: 'PERMANENT_CHANGE_OF_STATION', value: 'Permanent Change Of Station (PCS)' },
    { key: 'RETIREMENT', value: 'Retirement' },
    { key: 'SEPARATION', value: 'Separation' },
  ],
  currentStation: {},
};

describe('OrdersInfoForm component', () => {
  it('renders the form inputs', async () => {
    const { getByLabelText } = render(<OrdersInfoForm {...testProps} />);

    await waitFor(() => {
      expect(getByLabelText('Orders type')).toBeInstanceOf(HTMLSelectElement);
      expect(getByLabelText('Orders type')).toBeRequired();
      expect(getByLabelText('Orders date')).toBeInstanceOf(HTMLInputElement);
      expect(getByLabelText('Orders date')).toBeRequired();
      expect(getByLabelText('Report-by date')).toBeInstanceOf(HTMLInputElement);
      expect(getByLabelText('Report-by date')).toBeRequired();
      expect(getByLabelText('Yes')).toBeInstanceOf(HTMLInputElement);
      expect(getByLabelText('No')).toBeInstanceOf(HTMLInputElement);
      expect(getByLabelText('New duty station')).toBeInstanceOf(HTMLInputElement);
    });
  });

  it('renders each option for orders type', async () => {
    const { getByLabelText } = render(<OrdersInfoForm {...testProps} />);

    await waitFor(() => {
      const ordersTypeDropdown = getByLabelText('Orders type');
      expect(ordersTypeDropdown).toBeInstanceOf(HTMLSelectElement);

      userEvent.selectOptions(ordersTypeDropdown, 'PERMANENT_CHANGE_OF_STATION');
      expect(ordersTypeDropdown).toHaveValue('PERMANENT_CHANGE_OF_STATION');

      userEvent.selectOptions(ordersTypeDropdown, 'RETIREMENT');
      expect(ordersTypeDropdown).toHaveValue('RETIREMENT');

      userEvent.selectOptions(ordersTypeDropdown, 'SEPARATION');
      expect(ordersTypeDropdown).toHaveValue('SEPARATION');
    });
  });

  it('validates the new duty station against the current duty station', async () => {
    const { queryByText, getByRole, getByLabelText } = render(
      <OrdersInfoForm {...testProps} currentStation={{ name: 'Luke AFB' }} />,
    );

    userEvent.selectOptions(getByLabelText('Orders type'), 'PERMANENT_CHANGE_OF_STATION');
    userEvent.type(getByLabelText('Orders date'), '08 Nov 2020');
    userEvent.type(getByLabelText('Report-by date'), '26 Nov 2020');
    userEvent.click(getByLabelText('No'));

    // Test Duty Station Search Box interaction
    fireEvent.change(getByLabelText('New duty station'), { target: { value: 'AFB' } });
    await selectEvent.select(getByLabelText('New duty station'), /Luke/);
    expect(getByRole('form')).toHaveFormValues({
      new_duty_station: 'Luke AFB',
    });

    await waitFor(() => {
      expect(getByRole('button', { name: 'Next' })).toHaveAttribute('disabled');
      expect(
        queryByText('You entered the same duty station for your origin and destination. Please change one of them.'),
      ).toBeInTheDocument();
    });
  });

  it('shows an error message if trying to submit an invalid form', async () => {
    const { getByRole, getAllByText } = render(<OrdersInfoForm {...testProps} />);
    const submitBtn = getByRole('button', { name: 'Next' });

    userEvent.click(submitBtn);

    await waitFor(() => {
      expect(getAllByText('Required').length).toBe(3);
    });
    expect(testProps.onSubmit).not.toHaveBeenCalled();
  });

  it('submits the form when its valid', async () => {
    const { getByRole, getByLabelText } = render(<OrdersInfoForm {...testProps} />);

    userEvent.selectOptions(getByLabelText('Orders type'), 'PERMANENT_CHANGE_OF_STATION');
    userEvent.type(getByLabelText('Orders date'), '08 Nov 2020');
    userEvent.type(getByLabelText('Report-by date'), '26 Nov 2020');
    userEvent.click(getByLabelText('No'));

    // Test Duty Station Search Box interaction
    fireEvent.change(getByLabelText('New duty station'), { target: { value: 'AFB' } });
    await selectEvent.select(getByLabelText('New duty station'), /Luke/);
    expect(getByRole('form')).toHaveFormValues({
      new_duty_station: 'Luke AFB',
    });

    const submitBtn = getByRole('button', { name: 'Next' });
    userEvent.click(submitBtn);

    await waitFor(() => {
      expect(testProps.onSubmit).toHaveBeenCalledWith(
        expect.objectContaining({
          orders_type: 'PERMANENT_CHANGE_OF_STATION',
          has_dependents: 'no',
          issue_date: '08 Nov 2020',
          report_by_date: '26 Nov 2020',
          new_duty_station: {
            address: undefined,
            address_id: '25be4d12-fe93-47f1-bbec-1db386dfa67f',
            affiliation: 'AIR_FORCE',
            created_at: '2021-02-11T16:48:04.117Z',
            id: 'a8d6b33c-8370-4e92-8df2-356b8c9d0c1a',
            name: 'Luke AFB',
            updated_at: '2021-02-11T16:48:04.117Z',
          },
        }),
        expect.anything(),
      );
    });
  });

  it('implements the onBack handler when the Back button is clicked', async () => {
    const { getByRole } = render(<OrdersInfoForm {...testProps} />);
    const backBtn = getByRole('button', { name: 'Back' });

    userEvent.click(backBtn);

    await waitFor(() => {
      expect(testProps.onBack).toHaveBeenCalled();
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

    it('pre-fills the inputs', async () => {
      const { getByRole, queryByText, getByLabelText } = render(
        <OrdersInfoForm {...testProps} initialValues={testInitialValues} />,
      );

      await waitFor(() => {
        expect(getByRole('form')).toHaveFormValues({
          new_duty_station: 'Yuma AFB',
        });

        expect(getByLabelText('Orders type')).toHaveValue(testInitialValues.orders_type);
        expect(getByLabelText('Orders date')).toHaveValue('08 Nov 2020');
        expect(getByLabelText('Report-by date')).toHaveValue('26 Nov 2020');
        expect(getByLabelText('Yes')).not.toBeChecked();
        expect(getByLabelText('No')).toBeChecked();
        expect(queryByText('Yuma AFB')).toBeInTheDocument();
      });
    });
  });

  afterEach(jest.clearAllMocks);
});
