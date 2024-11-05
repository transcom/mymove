import React from 'react';
import { render, waitFor, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import OrdersInfoForm from './OrdersInfoForm';

import { showCounselingOffices } from 'services/internalApi';
import { ORDERS_TYPE } from 'constants/orders';

jest.setTimeout(60000);

jest.setTimeout(60000);

jest.mock('services/internalApi', () => ({
  ...jest.requireActual('services/internalApi'),
  showCounselingOffices: jest.fn().mockImplementation(() =>
    Promise.resolve({
      body: [
        {
          id: '3e937c1f-5539-4919-954d-017989130584',
          name: 'Albuquerque AFB',
        },
        {
          id: 'fa51dab0-4553-4732-b843-1f33407f77bc',
          name: 'Glendale Luke AFB',
        },
      ],
    }),
  ),
}));

jest.mock('components/LocationSearchBox/api', () => ({
  ShowAddress: jest.fn().mockImplementation(() =>
    Promise.resolve({
      city: 'Glendale Luke AFB',
      country: 'United States',
      id: 'fa51dab0-4553-4732-b843-1f33407f77bc',
      postalCode: '85309',
      state: 'AZ',
      streetAddress1: 'n/a',
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
          city: '',
          id: '00000000-0000-0000-0000-000000000000',
          postalCode: '',
          state: '',
          streetAddress1: '',
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
          city: 'Glendale Luke AFB',
          country: 'United States',
          id: 'fa51dab0-4553-4732-b843-1f33407f77bc',
          postalCode: '85309',
          state: 'AZ',
          streetAddress1: 'n/a',
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
          postalCode: '',
          state: '',
          streetAddress1: '',
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
          postalCode: '',
          state: '',
          streetAddress1: '',
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
          postalCode: '',
          state: '',
          streetAddress1: '',
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
          postalCode: '',
          state: '',
          streetAddress1: '',
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
  initialValues: {
    orders_type: '',
    issue_date: '',
    report_by_date: '',
    has_dependents: '',
    new_duty_location: {},
    grade: '',
    origin_duty_location: {},
  },
  onBack: jest.fn(),
  ordersTypeOptions: [
    { key: 'PERMANENT_CHANGE_OF_STATION', value: 'Permanent Change Of Station (PCS)' },
    { key: 'LOCAL_MOVE', value: 'Local Move' },
    { key: 'RETIREMENT', value: 'Retirement' },
    { key: 'SEPARATION', value: 'Separation' },
    { key: 'TEMPORARY_DUTY', value: 'Temporary Duty (TDY)' },
  ],
};

describe('OrdersInfoForm component', () => {
  it('renders the form inputs', async () => {
    const { getByLabelText } = render(<OrdersInfoForm {...testProps} />);

    await waitFor(() => {
      expect(getByLabelText(/Orders type/)).toBeInstanceOf(HTMLSelectElement);
      expect(getByLabelText(/Orders type/)).toBeRequired();
      expect(getByLabelText(/Orders date/)).toBeInstanceOf(HTMLInputElement);
      expect(getByLabelText(/Orders date/)).toBeRequired();
      expect(getByLabelText(/Report by date/)).toBeInstanceOf(HTMLInputElement);
      expect(getByLabelText(/Report by date/)).toBeRequired();
      expect(getByLabelText('Yes')).toBeInstanceOf(HTMLInputElement);
      expect(getByLabelText('No')).toBeInstanceOf(HTMLInputElement);
      expect(getByLabelText(/New duty location/)).toBeInstanceOf(HTMLInputElement);
      expect(getByLabelText(/Pay grade/)).toBeInstanceOf(HTMLSelectElement);
      expect(getByLabelText(/Current duty location/)).toBeInstanceOf(HTMLInputElement);
    });
  });

  it('renders each option for orders type', async () => {
    showCounselingOffices.mockImplementation(() => Promise.resolve({}));
    const { getByLabelText } = render(<OrdersInfoForm {...testProps} />);

    const ordersTypeDropdown = getByLabelText(/Orders type/);
    expect(ordersTypeDropdown).toBeInstanceOf(HTMLSelectElement);

    await userEvent.selectOptions(ordersTypeDropdown, ORDERS_TYPE.PERMANENT_CHANGE_OF_STATION);
    expect(ordersTypeDropdown).toHaveValue(ORDERS_TYPE.PERMANENT_CHANGE_OF_STATION);

    await userEvent.selectOptions(ordersTypeDropdown, ORDERS_TYPE.LOCAL_MOVE);
    expect(ordersTypeDropdown).toHaveValue(ORDERS_TYPE.LOCAL_MOVE);

    await userEvent.selectOptions(ordersTypeDropdown, ORDERS_TYPE.RETIREMENT);
    expect(ordersTypeDropdown).toHaveValue(ORDERS_TYPE.RETIREMENT);

    await userEvent.selectOptions(ordersTypeDropdown, ORDERS_TYPE.SEPARATION);
    expect(ordersTypeDropdown).toHaveValue(ORDERS_TYPE.SEPARATION);

    await userEvent.selectOptions(ordersTypeDropdown, ORDERS_TYPE.TEMPORARY_DUTY);
    expect(ordersTypeDropdown).toHaveValue(ORDERS_TYPE.TEMPORARY_DUTY);
  });

  it('allows new and current duty location to be the same', async () => {
    render(<OrdersInfoForm {...testProps} />);

    await userEvent.selectOptions(screen.getByLabelText(/Orders type/), ORDERS_TYPE.PERMANENT_CHANGE_OF_STATION);
    await userEvent.type(screen.getByLabelText(/Orders date/), '08 Nov 2020');
    await userEvent.type(screen.getByLabelText(/Report by date/), '26 Nov 2020');
    await userEvent.click(screen.getByLabelText('No'));
    await userEvent.selectOptions(screen.getByLabelText(/Pay grade/), ['E_5']);

    // Test Current Duty Location Search Box interaction
    await userEvent.type(screen.getByLabelText(/Current duty location/), 'AFB', { delay: 100 });
    const selectedOptionCurrent = await screen.findByText(/Altus/);
    await userEvent.click(selectedOptionCurrent);

    // Test New Duty Location Search Box interaction
    await userEvent.type(screen.getByLabelText(/New duty location/), 'AFB', { delay: 100 });
    const selectedOptionNew = await screen.findByText(/Luke/);
    await userEvent.click(selectedOptionNew);

    await waitFor(() => {
      expect(screen.getByRole('form')).toHaveFormValues({
        new_duty_location: 'Luke AFB',
        origin_duty_location: 'Altus AFB',
      });
    });

    expect(screen.getByRole('button', { name: 'Next' })).not.toHaveAttribute('disabled');
  });

  it('shows an error message if trying to submit an invalid form', async () => {
    const { getByRole, getAllByTestId } = render(<OrdersInfoForm {...testProps} />);

    // Touch required fields to show validation errors
    await userEvent.click(screen.getByLabelText(/Orders type/));
    await userEvent.click(screen.getByLabelText(/Orders date/));
    await userEvent.click(screen.getByLabelText(/Report by date/));
    await userEvent.click(screen.getByLabelText(/Pay grade/));

    const submitBtn = getByRole('button', { name: 'Next' });
    await userEvent.click(submitBtn);

    await waitFor(() => {
      expect(getAllByTestId('errorMessage').length).toBe(4);
    });
    expect(testProps.onSubmit).not.toHaveBeenCalled();
  });

  it('submits the form when its valid and has a counseling office selected', async () => {
    const testPropsWithCounselingOffice = {
      onSubmit: jest.fn().mockImplementation(() => Promise.resolve()),
      initialValues: {
        orders_type: '',
        issue_date: '',
        report_by_date: '',
        has_dependents: '',
        new_duty_location: {},
        grade: '',
        origin_duty_location: {},
        counseling_office_id: '3be2381f-f9ed-4902-bbdc-69c69e43eb86',
      },
      onBack: jest.fn(),
      ordersTypeOptions: [
        { key: 'PERMANENT_CHANGE_OF_STATION', value: 'Permanent Change Of Station (PCS)' },
        { key: 'LOCAL_MOVE', value: 'Local Move' },
        { key: 'RETIREMENT', value: 'Retirement' },
        { key: 'SEPARATION', value: 'Separation' },
        { key: 'TEMPORARY_DUTY', value: 'Temporary Duty (TDY)' },
      ],
    };

    render(<OrdersInfoForm {...testPropsWithCounselingOffice} />);

    await userEvent.selectOptions(screen.getByLabelText(/Orders type/), ORDERS_TYPE.PERMANENT_CHANGE_OF_STATION);
    await userEvent.type(screen.getByLabelText(/Orders date/), '08 Nov 2020');
    await userEvent.type(screen.getByLabelText(/Report by date/), '26 Nov 2020');
    await userEvent.click(screen.getByLabelText('No'));
    await userEvent.selectOptions(screen.getByLabelText(/Pay grade/), ['E_5']);

    // Test Current Duty Location Search Box interaction
    await userEvent.type(screen.getByLabelText(/Current duty location/), 'AFB', { delay: 100 });
    const selectedOptionCurrent = await screen.findByText(/Altus/);
    await userEvent.click(selectedOptionCurrent);

    // Test New Duty Location Search Box interaction
    await userEvent.type(screen.getByLabelText(/New duty location/), 'AFB', { delay: 100 });
    const selectedOptionNew = await screen.findByText(/Luke/);
    await userEvent.click(selectedOptionNew);

    await waitFor(() => {
      expect(screen.getByRole('form')).toHaveFormValues({
        new_duty_location: 'Luke AFB',
        origin_duty_location: 'Altus AFB',
      });
    });

    const submitBtn = screen.getByRole('button', { name: 'Next' });
    await userEvent.click(submitBtn);

    await waitFor(() => {
      expect(testPropsWithCounselingOffice.onSubmit).toHaveBeenCalledWith(
        expect.objectContaining({
          orders_type: ORDERS_TYPE.PERMANENT_CHANGE_OF_STATION,
          counseling_office_id: '3be2381f-f9ed-4902-bbdc-69c69e43eb86',
          has_dependents: 'no',
          issue_date: '08 Nov 2020',
          report_by_date: '26 Nov 2020',
          new_duty_location: {
            address: {
              city: 'Glendale Luke AFB',
              country: 'United States',
              id: 'fa51dab0-4553-4732-b843-1f33407f77bc',
              postalCode: '85309',
              state: 'AZ',
              streetAddress1: 'n/a',
            },
            address_id: '25be4d12-fe93-47f1-bbec-1db386dfa67f',
            affiliation: 'AIR_FORCE',
            created_at: '2021-02-11T16:48:04.117Z',
            id: 'a8d6b33c-8370-4e92-8df2-356b8c9d0c1a',
            name: 'Luke AFB',
            updated_at: '2021-02-11T16:48:04.117Z',
          },
          grade: 'E_5',
          origin_duty_location: {
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
        }),
        expect.anything(),
      );
    });
  });

  it('submits the form when its valid', async () => {
    render(<OrdersInfoForm {...testProps} />);

    await userEvent.selectOptions(screen.getByLabelText(/Orders type/), ORDERS_TYPE.PERMANENT_CHANGE_OF_STATION);
    await userEvent.type(screen.getByLabelText(/Orders date/), '08 Nov 2020');
    await userEvent.type(screen.getByLabelText(/Report by date/), '26 Nov 2020');
    await userEvent.click(screen.getByLabelText('No'));
    await userEvent.selectOptions(screen.getByLabelText(/Pay grade/), ['E_5']);

    // Test Current Duty Location Search Box interaction
    await userEvent.type(screen.getByLabelText(/Current duty location/), 'AFB', { delay: 100 });
    const selectedOptionCurrent = await screen.findByText(/Altus/);
    await userEvent.click(selectedOptionCurrent);

    // Test New Duty Location Search Box interaction
    await userEvent.type(screen.getByLabelText(/New duty location/), 'AFB', { delay: 100 });
    const selectedOptionNew = await screen.findByText(/Luke/);
    await userEvent.click(selectedOptionNew);

    await waitFor(() => {
      expect(screen.getByRole('form')).toHaveFormValues({
        new_duty_location: 'Luke AFB',
        origin_duty_location: 'Altus AFB',
      });
    });

    const submitBtn = screen.getByRole('button', { name: 'Next' });
    await userEvent.click(submitBtn);

    await waitFor(() => {
      expect(testProps.onSubmit).toHaveBeenCalledWith(
        expect.objectContaining({
          orders_type: ORDERS_TYPE.PERMANENT_CHANGE_OF_STATION,
          has_dependents: 'no',
          issue_date: '08 Nov 2020',
          report_by_date: '26 Nov 2020',
          new_duty_location: {
            address: {
              city: 'Glendale Luke AFB',
              country: 'United States',
              id: 'fa51dab0-4553-4732-b843-1f33407f77bc',
              postalCode: '85309',
              state: 'AZ',
              streetAddress1: 'n/a',
            },
            address_id: '25be4d12-fe93-47f1-bbec-1db386dfa67f',
            affiliation: 'AIR_FORCE',
            created_at: '2021-02-11T16:48:04.117Z',
            id: 'a8d6b33c-8370-4e92-8df2-356b8c9d0c1a',
            name: 'Luke AFB',
            updated_at: '2021-02-11T16:48:04.117Z',
          },
          grade: 'E_5',
          origin_duty_location: {
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
        }),
        expect.anything(),
      );
    });
  });

  it('submits the form when temporary duty orders type is selected', async () => {
    render(<OrdersInfoForm {...testProps} />);

    await userEvent.selectOptions(screen.getByLabelText(/Orders type/), ORDERS_TYPE.TEMPORARY_DUTY);
    await userEvent.type(screen.getByLabelText(/Orders date/), '28 Oct 2024');
    await userEvent.type(screen.getByLabelText(/Report by date/), '28 Oct 2024');
    await userEvent.click(screen.getByLabelText('No'));
    await userEvent.selectOptions(screen.getByLabelText(/Pay grade/), ['E_7']);

    // Test Current Duty Location Search Box interaction
    await userEvent.type(screen.getByLabelText(/Current duty location/), 'AFB', { delay: 100 });
    const selectedOptionCurrent = await screen.findByText(/Altus/);
    await userEvent.click(selectedOptionCurrent);

    // Test New Duty Location Search Box interaction
    await userEvent.type(screen.getByLabelText(/New duty location/), 'AFB', { delay: 100 });
    const selectedOptionNew = await screen.findByText(/Luke/);
    await userEvent.click(selectedOptionNew);

    const submitBtn = screen.getByRole('button', { name: 'Next' });
    await userEvent.click(submitBtn);

    await waitFor(() => {
      expect(testProps.onSubmit).toHaveBeenCalledWith(
        expect.objectContaining({
          orders_type: ORDERS_TYPE.TEMPORARY_DUTY,
          has_dependents: 'no',
          issue_date: '28 Oct 2024',
          report_by_date: '28 Oct 2024',
          new_duty_location: {
            address: {
              city: 'Glendale Luke AFB',
              country: 'United States',
              id: 'fa51dab0-4553-4732-b843-1f33407f77bc',
              postalCode: '85309',
              state: 'AZ',
              streetAddress1: 'n/a',
            },
            address_id: '25be4d12-fe93-47f1-bbec-1db386dfa67f',
            affiliation: 'AIR_FORCE',
            created_at: '2021-02-11T16:48:04.117Z',
            id: 'a8d6b33c-8370-4e92-8df2-356b8c9d0c1a',
            name: 'Luke AFB',
            updated_at: '2021-02-11T16:48:04.117Z',
          },
          grade: 'E_7',
          origin_duty_location: {
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
        }),
        expect.anything(),
      );
    });
  });

  it('implements the onBack handler when the Back button is clicked', async () => {
    const { getByRole } = render(<OrdersInfoForm {...testProps} />);
    const backBtn = getByRole('button', { name: 'Back' });

    await userEvent.click(backBtn);

    await waitFor(() => {
      expect(testProps.onBack).toHaveBeenCalled();
    });
  });

  describe('with initial values', () => {
    const testInitialValues = {
      orders_type: ORDERS_TYPE.PERMANENT_CHANGE_OF_STATION,
      issue_date: '2020-11-08',
      report_by_date: '2020-11-26',
      has_dependents: 'no',
      new_duty_location: {
        address: {
          city: 'Des Moines',
          country: 'US',
          id: 'a4b30b99-4e82-48a6-b736-01662b499d6a',
          postalCode: '50309',
          state: 'IA',
          streetAddress1: '987 Other Avenue',
          streetAddress2: 'P.O. Box 1234',
          streetAddress3: 'c/o Another Person',
        },
        address_id: 'a4b30b99-4e82-48a6-b736-01662b499d6a',
        affiliation: 'AIR_FORCE',
        created_at: '2020-10-19T17:01:16.114Z',
        id: 'f9299768-16d2-4a13-ae39-7087a58b1f62',
        name: 'Yuma AFB',
        updated_at: '2020-10-19T17:01:16.114Z',
      },
      grade: 'E_1',
      origin_duty_location: {
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
    };

    it('pre-fills the inputs', async () => {
      const { getByRole, queryByText, getByLabelText } = render(
        <OrdersInfoForm {...testProps} initialValues={testInitialValues} />,
      );

      await waitFor(() => {
        expect(getByRole('form')).toHaveFormValues({
          new_duty_location: 'Yuma AFB',
          origin_duty_location: 'Altus AFB',
        });

        expect(getByLabelText(/Orders type/)).toHaveValue(testInitialValues.orders_type);
        expect(getByLabelText(/Orders date/)).toHaveValue('08 Nov 2020');
        expect(getByLabelText(/Report by date/)).toHaveValue('26 Nov 2020');
        expect(getByLabelText('Yes')).not.toBeChecked();
        expect(getByLabelText('No')).toBeChecked();
        expect(queryByText('Yuma AFB')).toBeInTheDocument();
        expect(getByLabelText(/Pay grade/)).toHaveValue(testInitialValues.grade);
        expect(queryByText('Altus AFB')).toBeInTheDocument();
      });
    });
  });

  afterEach(jest.restoreAllMocks);
});
