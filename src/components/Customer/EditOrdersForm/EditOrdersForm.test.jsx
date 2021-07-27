import React from 'react';
import { screen, render, waitFor, fireEvent } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import selectEvent from 'react-select-event';

import EditOrdersForm from './EditOrdersForm';

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
  onCancel: jest.fn(),
  onUploadComplete: jest.fn(),
  createUpload: jest.fn(),
  onDelete: jest.fn(),
  existingUploads: [],
  filePond: {},
  ordersTypeOptions: [
    { key: 'PERMANENT_CHANGE_OF_STATION', value: 'Permanent Change Of Station (PCS)' },
    { key: 'RETIREMENT', value: 'Retirement' },
    { key: 'SEPARATION', value: 'Separation' },
  ],
  currentStation: {},
};

describe('EditOrdersForm component', () => {
  it('renders the form inputs', async () => {
    render(<EditOrdersForm {...testProps} />);

    expect(await screen.findByLabelText('Orders type')).toBeInstanceOf(HTMLSelectElement);
    expect(await screen.findByLabelText('Orders type')).toBeRequired();
    expect(await screen.findByLabelText('Orders date')).toBeInstanceOf(HTMLInputElement);
    expect(await screen.findByLabelText('Orders date')).toBeRequired();
    expect(await screen.findByLabelText('Report-by date')).toBeInstanceOf(HTMLInputElement);
    expect(await screen.findByLabelText('Report-by date')).toBeRequired();
    expect(await screen.findByLabelText('Yes')).toBeInstanceOf(HTMLInputElement);
    expect(await screen.findByLabelText('No')).toBeInstanceOf(HTMLInputElement);
    expect(await screen.findByLabelText('New duty station')).toBeInstanceOf(HTMLInputElement);
  });

  it.each([
    ['PERMANENT_CHANGE_OF_STATION', 'PERMANENT_CHANGE_OF_STATION'],
    ['RETIREMENT', 'RETIREMENT'],
    ['SEPARATION', 'SEPARATION'],
  ])('renders the %s option for the orders type field', async (selectionOption, expectedValue) => {
    render(<EditOrdersForm {...testProps} />);

    const ordersTypeDropdown = await screen.findByLabelText('Orders type');
    expect(ordersTypeDropdown).toBeInstanceOf(HTMLSelectElement);

    userEvent.selectOptions(ordersTypeDropdown, selectionOption);
    await waitFor(() => {
      expect(ordersTypeDropdown).toHaveValue(expectedValue);
    });
  });

  it('validates the new duty station against the current duty station', async () => {
    render(<EditOrdersForm {...testProps} currentStation={{ name: 'Luke AFB' }} />);

    userEvent.selectOptions(screen.getByLabelText('Orders type'), 'PERMANENT_CHANGE_OF_STATION');
    userEvent.type(screen.getByLabelText('Orders date'), '08 Nov 2020');
    userEvent.type(screen.getByLabelText('Report-by date'), '26 Nov 2020');
    userEvent.click(screen.getByLabelText('No'));

    // Test Duty Station Search Box interaction
    fireEvent.change(screen.getByLabelText('New duty station'), { target: { value: 'AFB' } });
    await selectEvent.select(await screen.findByLabelText('New duty station'), /Luke/);
    await waitFor(() => {
      expect(screen.getByRole('form')).toHaveFormValues({
        new_duty_station: 'Luke AFB',
      });
    });
    const submitButton = screen.getByRole('button', { name: 'Save' });

    await waitFor(() => {
      expect(submitButton).toBeDisabled();
      expect(
        screen.queryByText(
          'You entered the same duty station for your origin and destination. Please change one of them.',
        ),
      ).toBeInTheDocument();
    });
  });

  it('shows an error message if the form is invalid', async () => {
    const initialValues = {
      orders_type: 'PERMANENT_CHANGE_OF_STATION',
      issue_date: '2020-11-08',
      report_by_date: '2020-11-26',
      has_dependents: 'No',
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

    const existingUploads = [
      {
        id: '123',
        created_at: '2020-11-08',
        bytes: 1,
        url: 'url',
        filename: 'Test Upload',
      },
    ];

    render(<EditOrdersForm {...testProps} initialValues={initialValues} existingUploads={existingUploads} />);
    const submitButton = screen.getByRole('button', { name: 'Save' });

    await waitFor(() => {
      expect(submitButton).toBeEnabled();
    });

    const ordersTypeDropdown = screen.getByLabelText('Orders type');
    userEvent.selectOptions(ordersTypeDropdown, '');
    userEvent.tab();

    await waitFor(() => {
      expect(submitButton).toBeDisabled();
    });

    const required = screen.getByText('Required');
    expect(required).toBeInTheDocument();
  });

  /*
  it('submits the form when its valid', async () => {
    render(<EditOrdersForm {...testProps} />);

    userEvent.selectOptions(screen.getByLabelText('Orders type'), 'PERMANENT_CHANGE_OF_STATION');
    userEvent.type(screen.getByLabelText('Orders date'), '08 Nov 2020');
    userEvent.type(screen.getByLabelText('Report-by date'), '26 Nov 2020');
    userEvent.click(screen.getByLabelText('No'));

    // Test Duty Station Search Box interaction
    fireEvent.change(screen.getByLabelText('New duty station'), { target: { value: 'AFB' } });
    await selectEvent.select(screen.getByLabelText('New duty station'), /Luke/);
    expect(screen.getByRole('form')).toHaveFormValues({
      new_duty_station: 'Luke AFB',
    });

    const submitBtn = screen.getByRole('button', { name: 'Save' });
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
  */

  it('implements the onCancel handler when the Cancel button is clicked', async () => {
    render(<EditOrdersForm {...testProps} />);
    const cancelButton = screen.getByRole('button', { name: 'Cancel' });

    userEvent.click(cancelButton);

    await waitFor(() => {
      expect(testProps.onCancel).toHaveBeenCalled();
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
      render(<EditOrdersForm {...testProps} initialValues={testInitialValues} />);

      expect(await screen.findByRole('form')).toHaveFormValues({
        new_duty_station: 'Yuma AFB',
      });

      expect(await screen.findByLabelText('Orders type')).toHaveValue(testInitialValues.orders_type);
      expect(await screen.findByLabelText('Orders date')).toHaveValue('08 Nov 2020');
      expect(await screen.findByLabelText('Report-by date')).toHaveValue('26 Nov 2020');
      expect(await screen.findByLabelText('Yes')).not.toBeChecked();
      expect(await screen.findByLabelText('No')).toBeChecked();
      expect(await screen.findByText('Yuma AFB')).toBeInTheDocument();
    });
  });

  describe('disables the save button', () => {
    it('when no orders type is selected', async () => {
      render(<EditOrdersForm {...testProps} />);

      const save = screen.getByRole('button', { name: 'Save' });
      await waitFor(() => {
        expect(save).toBeInTheDocument();
      });

      expect(save).toBeDisabled();
    });

    it('when no orders date is selected', async () => {
      render(<EditOrdersForm {...testProps} />);

      const save = screen.getByRole('button', { name: 'Save' });
      await waitFor(() => {
        expect(save).toBeInTheDocument();
      });

      expect(save).toBeDisabled();
    });

    it('when no report by date is selected', async () => {
      render(<EditOrdersForm {...testProps} />);

      const save = screen.getByRole('button', { name: 'Save' });
      await waitFor(() => {
        expect(save).toBeInTheDocument();
      });

      expect(save).toBeDisabled();
    });

    it('when no duty station is selected', async () => {
      render(<EditOrdersForm {...testProps} />);

      const save = screen.getByRole('button', { name: 'Save' });
      await waitFor(() => {
        expect(save).toBeInTheDocument();
      });

      expect(save).toBeDisabled();
    });

    it('when no orders are uploaded', async () => {
      render(<EditOrdersForm {...testProps} />);

      const save = screen.getByRole('button', { name: 'Save' });
      await waitFor(() => {
        expect(save).toBeInTheDocument();
      });

      expect(save).toBeDisabled();
    });
  });

  afterEach(jest.restoreAllMocks);
});
