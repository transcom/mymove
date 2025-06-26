import React from 'react';
import { render, waitFor, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { Provider } from 'react-redux';

import { isBooleanFlagEnabled } from '../../../utils/featureFlags';

import OrdersInfoForm from './OrdersInfoForm';

import { showCounselingOffices } from 'services/internalApi';
import { ORDERS_BRANCH_OPTIONS, ORDERS_PAY_GRADE_TYPE, ORDERS_TYPE, ORDERS_TYPE_OPTIONS } from 'constants/orders';
import { configureStore } from 'shared/store';
import { MockProviders } from 'testUtils';

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
  getRankOptions: jest.fn().mockImplementation(() => {
    return Promise.resolve([
      {
        id: 'cb0ee2b8-e852-40fe-b972-2730b53860c7',
        paygradeId: '5f871c82-f259-43cc-9245-a6e18975dde0',
        rankAbbv: 'SSgt',
        rankOrder: 24,
      },
    ]);
  }),
  getPayGradeOptions: jest.fn().mockImplementation(() => {
    const MOCKED__ORDERS_PAY_GRADE_TYPE = {
      E_5: 'E-5',
      E_6: 'E-6',
      CIVILIAN_EMPLOYEE: 'CIVILIAN_EMPLOYEE',
    };

    return Promise.resolve({
      body: [
        {
          grade: MOCKED__ORDERS_PAY_GRADE_TYPE.E_5,
          description: MOCKED__ORDERS_PAY_GRADE_TYPE.E_5,
        },
        {
          grade: MOCKED__ORDERS_PAY_GRADE_TYPE.E_6,
          description: MOCKED__ORDERS_PAY_GRADE_TYPE.E_6,
        },
        {
          description: MOCKED__ORDERS_PAY_GRADE_TYPE.CIVILIAN_EMPLOYEE,
          grade: MOCKED__ORDERS_PAY_GRADE_TYPE.CIVILIAN_EMPLOYEE,
        },
      ],
    });
  }),
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
      {
        address: {
          city: '',
          id: '1111111111',
          postalCode: '',
          state: '',
          streetAddress1: '',
        },
        address_id: '4334640b-c35e-4293-a2f1-36c7b629f903',
        affiliation: 'AIR_FORCE',
        created_at: '2021-02-11T16:48:04.117Z',
        id: '22f0755f-6f35-478b-9a75-35a69211da1c',
        name: 'Scott AFB',
        updated_at: '2021-02-11T16:48:04.117Z',
        provides_services_counseling: true,
      },
    ]),
  ),
}));

jest.mock('../../../utils/featureFlags', () => ({
  isBooleanFlagEnabled: jest.fn(),
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
    { key: ORDERS_TYPE.EARLY_RETURN_OF_DEPENDENTS, value: ORDERS_TYPE_OPTIONS.EARLY_RETURN_OF_DEPENDENTS },
    { key: ORDERS_TYPE.STUDENT_TRAVEL, value: ORDERS_TYPE_OPTIONS.STUDENT_TRAVEL },
  ],
  affiliation: ORDERS_BRANCH_OPTIONS.AIR_FORCE,
};

const civilianTDYTestProps = {
  onSubmit: jest.fn().mockImplementation(() => Promise.resolve()),
  initialValues: {
    orders_type: ORDERS_TYPE_OPTIONS.TEMPORARY_DUTY,
    issue_date: '',
    report_by_date: '',
    has_dependents: '',
    uploaded_orders: [],
    grade: ORDERS_PAY_GRADE_TYPE.CIVILIAN_EMPLOYEE,
    origin_duty_location: { name: 'Luke AFB', address: { isOconus: false } },
    new_duty_location: { name: 'Luke AFB', provides_services_counseling: false, address: { isOconus: true } },
  },
  onCancel: jest.fn(),
  onUploadComplete: jest.fn(),
  createUpload: jest.fn(),
  onDelete: jest.fn(),
  filePond: {},
  ordersTypeOptions: [
    { key: 'PERMANENT_CHANGE_OF_STATION', value: 'Permanent Change Of Station (PCS)' },
    { key: 'LOCAL_MOVE', value: 'Local Move' },
    { key: 'RETIREMENT', value: 'Retirement' },
    { key: 'SEPARATION', value: 'Separation' },
    { key: 'TEMPORARY_DUTY', value: 'Temporary Duty (TDY)' },
    { key: ORDERS_TYPE.EARLY_RETURN_OF_DEPENDENTS, value: ORDERS_TYPE_OPTIONS.EARLY_RETURN_OF_DEPENDENTS },
    { key: ORDERS_TYPE.STUDENT_TRAVEL, value: ORDERS_TYPE_OPTIONS.STUDENT_TRAVEL },
  ],
  currentDutyLocation: { name: 'Luke AFB', address: { isOconus: false } },
  grade: ORDERS_PAY_GRADE_TYPE.CIVILIAN_EMPLOYEE,
  orders_type: ORDERS_TYPE.TEMPORARY_DUTY,
};

const mockStore = configureStore({});

describe('OrdersInfoForm component', () => {
  it('renders the form inputs', async () => {
    const { getByLabelText } = render(
      <Provider store={mockStore.store}>
        <OrdersInfoForm {...testProps} />
      </Provider>,
    );

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

      expect(screen.getByTestId('reqAsteriskMsg')).toBeInTheDocument();

      // check for asterisks on required fields
      const formGroups = screen.getAllByTestId('formGroup');

      formGroups.forEach((group) => {
        const hasRequiredField = group.querySelector('[required]') !== null;

        if (hasRequiredField) {
          expect(group.textContent).toContain('*');
        }
      });
    });
  });

  it('renders each option for orders type', async () => {
    isBooleanFlagEnabled.mockImplementation(() => Promise.resolve(true));

    showCounselingOffices.mockImplementation(() => Promise.resolve({}));
    const { getByLabelText } = render(
      <Provider store={mockStore.store}>
        <OrdersInfoForm {...testProps} />
      </Provider>,
    );

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

    await userEvent.selectOptions(ordersTypeDropdown, ORDERS_TYPE.EARLY_RETURN_OF_DEPENDENTS);
    expect(ordersTypeDropdown).toHaveValue(ORDERS_TYPE.EARLY_RETURN_OF_DEPENDENTS);

    await userEvent.selectOptions(ordersTypeDropdown, ORDERS_TYPE.STUDENT_TRAVEL);
    expect(ordersTypeDropdown).toHaveValue(ORDERS_TYPE.STUDENT_TRAVEL);
  });

  it('allows new and current duty location to be the same', async () => {
    render(
      <Provider store={mockStore.store}>
        <OrdersInfoForm {...testProps} />
      </Provider>,
    );

    await userEvent.selectOptions(screen.getByLabelText(/Orders type/), ORDERS_TYPE.PERMANENT_CHANGE_OF_STATION);
    await userEvent.type(screen.getByLabelText(/Orders date/), '08 Nov 2020');
    await userEvent.type(screen.getByLabelText(/Report by date/), '26 Nov 2020');
    await userEvent.click(screen.getByLabelText('No'));
    await userEvent.selectOptions(screen.getByLabelText(/Pay grade/), [ORDERS_PAY_GRADE_TYPE.E_5]);

    // Test Current Duty Location Search Box interaction
    await userEvent.type(screen.getByLabelText(/Current duty location/), 'AFB', { delay: 100 });
    const selectedOptionCurrent = await screen.findByText('Altus');
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
    const { getByRole, getAllByTestId } = render(
      <Provider store={mockStore.store}>
        <OrdersInfoForm {...testProps} />
      </Provider>,
    );

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

  it('renders the counseling office if current duty location provides services counseling', async () => {
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
        counseling_office_id: '',
      },
      onBack: jest.fn(),
      ordersTypeOptions: [
        { key: 'PERMANENT_CHANGE_OF_STATION', value: 'Permanent Change Of Station (PCS)' },
        { key: 'LOCAL_MOVE', value: 'Local Move' },
        { key: 'RETIREMENT', value: 'Retirement' },
        { key: 'SEPARATION', value: 'Separation' },
        { key: 'TEMPORARY_DUTY', value: 'Temporary Duty (TDY)' },
        { key: ORDERS_TYPE.EARLY_RETURN_OF_DEPENDENTS, value: ORDERS_TYPE_OPTIONS.EARLY_RETURN_OF_DEPENDENTS },
        { key: ORDERS_TYPE.STUDENT_TRAVEL, value: ORDERS_TYPE_OPTIONS.STUDENT_TRAVEL },
      ],
      affiliation: ORDERS_BRANCH_OPTIONS.AIR_FORCE,
    };

    render(
      <Provider store={mockStore.store}>
        <OrdersInfoForm {...testPropsWithCounselingOffice} />
      </Provider>,
    );
    await userEvent.selectOptions(screen.getByLabelText(/Orders type/), ORDERS_TYPE.PERMANENT_CHANGE_OF_STATION);
    await userEvent.type(screen.getByLabelText(/Orders date/), '08 Nov 2020');
    await userEvent.type(screen.getByLabelText(/Report by date/), '26 Nov 2020');
    await userEvent.click(screen.getByLabelText('No'));
    await userEvent.selectOptions(screen.getByLabelText(/Pay grade/), [ORDERS_PAY_GRADE_TYPE.E_5]);

    await userEvent.type(screen.getByLabelText(/Current duty location/), 'AFB', { delay: 100 });
    const selectedOptionCurrent = await screen.findByText(/Scott/);
    await userEvent.click(selectedOptionCurrent);

    await userEvent.type(screen.getByLabelText(/New duty location/), 'AFB', { delay: 100 });
    const selectedOptionNew = await screen.findByText(/Luke/);
    await userEvent.click(selectedOptionNew);

    await waitFor(() => {
      expect(screen.getByLabelText(/Counseling office/));
    });
  });

  it('does not render the counseling office if current duty location does not provides services counseling', async () => {
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
        counseling_office_id: '',
      },
      onBack: jest.fn(),
      ordersTypeOptions: [
        { key: 'PERMANENT_CHANGE_OF_STATION', value: 'Permanent Change Of Station (PCS)' },
        { key: 'LOCAL_MOVE', value: 'Local Move' },
        { key: 'RETIREMENT', value: 'Retirement' },
        { key: 'SEPARATION', value: 'Separation' },
        { key: 'TEMPORARY_DUTY', value: 'Temporary Duty (TDY)' },
      ],
      affiliation: ORDERS_BRANCH_OPTIONS.AIR_FORCE,
    };

    render(
      <Provider store={mockStore.store}>
        <OrdersInfoForm {...testPropsWithCounselingOffice} />
      </Provider>,
    );
    await userEvent.selectOptions(screen.getByLabelText(/Orders type/), ORDERS_TYPE.PERMANENT_CHANGE_OF_STATION);
    await userEvent.type(screen.getByLabelText(/Orders date/), '08 Nov 2020');
    await userEvent.type(screen.getByLabelText(/Report by date/), '26 Nov 2020');
    await userEvent.click(screen.getByLabelText('No'));
    await userEvent.selectOptions(screen.getByLabelText(/Pay grade/), [ORDERS_PAY_GRADE_TYPE.E_5]);

    await userEvent.type(screen.getByLabelText(/Current duty location/), 'AFB', { delay: 100 });
    const selectedOptionCurrent = await screen.findByText(/Altus AFB/);
    await userEvent.click(selectedOptionCurrent);

    await userEvent.type(screen.getByLabelText(/New duty location/), 'AFB', { delay: 100 });
    const selectedOptionNew = await screen.findByText(/Luke/);
    await userEvent.click(selectedOptionNew);

    expect(screen.queryByText(/Counseling office/)).not.toBeInTheDocument();
  });

  it('submits the form when its valid', async () => {
    render(
      <Provider store={mockStore.store}>
        <OrdersInfoForm {...testProps} />
      </Provider>,
    );

    await userEvent.selectOptions(screen.getByLabelText(/Orders type/), ORDERS_TYPE.PERMANENT_CHANGE_OF_STATION);
    await userEvent.type(screen.getByLabelText(/Orders date/), '08 Nov 2020');
    await userEvent.type(screen.getByLabelText(/Report by date/), '26 Nov 2020');
    await userEvent.click(screen.getByLabelText('No'));
    await userEvent.selectOptions(screen.getByLabelText(/Pay grade/), [ORDERS_PAY_GRADE_TYPE.E_5]);

    // Test Current Duty Location Search Box interaction
    await userEvent.type(screen.getByLabelText(/Current duty location/, { exact: false }), 'AFB', { delay: 100 });
    const selectedOptionCurrent = await screen.findByText('Altus');
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
          grade: ORDERS_PAY_GRADE_TYPE.E_5,
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
    render(
      <Provider store={mockStore.store}>
        <OrdersInfoForm {...testProps} />
      </Provider>,
    );
    await userEvent.selectOptions(screen.getByLabelText(/Orders type/), ORDERS_TYPE.TEMPORARY_DUTY);
    await userEvent.type(screen.getByLabelText(/Orders date/), '28 Oct 2024');
    await userEvent.type(screen.getByLabelText(/Report by date/), '28 Oct 2024');
    await userEvent.click(screen.getByLabelText('No'));
    await userEvent.selectOptions(screen.getByLabelText(/Pay grade/), [ORDERS_PAY_GRADE_TYPE.E_5]);

    // Test Current Duty Location Search Box interaction
    await userEvent.type(screen.getByLabelText(/Current duty location/), 'AFB', { delay: 100 });
    const selectedOptionCurrent = await screen.findByText('Altus');
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
          grade: ORDERS_PAY_GRADE_TYPE.E_5,
          rank: 'cb0ee2b8-e852-40fe-b972-2730b53860c7',
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
    const { getByRole } = render(
      <Provider store={mockStore.store}>
        <OrdersInfoForm {...testProps} />
      </Provider>,
    );
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
      grade: ORDERS_PAY_GRADE_TYPE.E_5,
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
        <Provider store={mockStore.store}>
          <OrdersInfoForm {...testProps} initialValues={testInitialValues} />
        </Provider>,
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
        expect(getByLabelText(/Pay grade/)).toHaveTextContent(ORDERS_PAY_GRADE_TYPE.E_5);
        expect(queryByText('Altus AFB')).toBeInTheDocument();
      });
    });
  });

  it('has dependents is yes and disabled when order type is student travel', async () => {
    render(
      <Provider store={mockStore.store}>
        <OrdersInfoForm {...testProps} />
      </Provider>,
    );

    await userEvent.selectOptions(screen.getByLabelText(/Orders type/), ORDERS_TYPE.STUDENT_TRAVEL);

    const hasDependentsYes = screen.getByLabelText('Yes');
    const hasDependentsNo = screen.getByLabelText('No');

    await waitFor(() => {
      expect(hasDependentsYes).toBeChecked();
      expect(hasDependentsYes).toBeDisabled();
      expect(hasDependentsNo).toBeDisabled();
    });
  });

  it('has dependents is yes and disabled when order type is early return', async () => {
    render(
      <Provider store={mockStore.store}>
        <OrdersInfoForm {...testProps} />
      </Provider>,
    );

    await userEvent.selectOptions(screen.getByLabelText(/Orders type/), ORDERS_TYPE.EARLY_RETURN_OF_DEPENDENTS);

    const hasDependentsYes = screen.getByLabelText('Yes');
    const hasDependentsNo = screen.getByLabelText('No');

    await waitFor(() => {
      expect(hasDependentsYes).toBeChecked();
      expect(hasDependentsYes).toBeDisabled();
      expect(hasDependentsNo).toBeDisabled();
    });
  });

  it('has dependents becomes disabled and then re-enabled for order type student travel', async () => {
    render(
      <Provider store={mockStore.store}>
        <OrdersInfoForm {...testProps} />
      </Provider>,
    );
    // set order type to perm change and verify the "has dependents" state
    await userEvent.selectOptions(screen.getByLabelText(/Orders type/), 'PERMANENT_CHANGE_OF_STATION');

    const hasDependentsYesPermChg = screen.getByLabelText('Yes');
    const hasDependentsNoPermChg = screen.getByLabelText('No');

    await waitFor(() => {
      expect(hasDependentsYesPermChg).not.toBeChecked();
      expect(hasDependentsYesPermChg).toBeEnabled();
      expect(hasDependentsNoPermChg).not.toBeChecked();
      expect(hasDependentsNoPermChg).toBeEnabled();
    });

    // set order type to value that disables and defaults "has dependents"
    await userEvent.selectOptions(screen.getByLabelText(/Orders type/), ORDERS_TYPE.STUDENT_TRAVEL);

    // set order type to value the re-enables "has dependents"
    await userEvent.selectOptions(screen.getByLabelText(/Orders type/), 'LOCAL_MOVE');

    const hasDependentsYesLocalMove = screen.getByLabelText('Yes');
    const hasDependentsNoLocalMove = screen.getByLabelText('No');

    await waitFor(() => {
      expect(hasDependentsYesLocalMove).not.toBeChecked();
      expect(hasDependentsYesLocalMove).toBeEnabled();
      expect(hasDependentsNoLocalMove).not.toBeChecked();
      expect(hasDependentsNoLocalMove).toBeEnabled();
    });
  });

  it('has dependents becomes disabled and then re-enabled for order type early return', async () => {
    render(
      <Provider store={mockStore.store}>
        <OrdersInfoForm {...testProps} />
      </Provider>,
    );
    // set order type to perm change and verify the "has dependents" state
    await userEvent.selectOptions(screen.getByLabelText(/Orders type/), 'PERMANENT_CHANGE_OF_STATION');

    const hasDependentsYesPermChg = screen.getByLabelText('Yes');
    const hasDependentsNoPermChg = screen.getByLabelText('No');

    await waitFor(() => {
      expect(hasDependentsYesPermChg).not.toBeChecked();
      expect(hasDependentsYesPermChg).toBeEnabled();
      expect(hasDependentsNoPermChg).not.toBeChecked();
      expect(hasDependentsNoPermChg).toBeEnabled();
    });

    // set order type to value that disables and defaults "has dependents"
    await userEvent.selectOptions(screen.getByLabelText(/Orders type/), ORDERS_TYPE.EARLY_RETURN_OF_DEPENDENTS);

    // set order type to value the re-enables "has dependents"
    await userEvent.selectOptions(screen.getByLabelText(/Orders type/), 'LOCAL_MOVE');

    const hasDependentsYesLocalMove = screen.getByLabelText('Yes');
    const hasDependentsNoLocalMove = screen.getByLabelText('No');

    await waitFor(() => {
      expect(hasDependentsYesLocalMove).not.toBeChecked();
      expect(hasDependentsYesLocalMove).toBeEnabled();
      expect(hasDependentsNoLocalMove).not.toBeChecked();
      expect(hasDependentsNoLocalMove).toBeEnabled();
    });
  });

  it('does not render civilian TDY UB Allowance field if orders type is not TDY', async () => {
    isBooleanFlagEnabled.mockResolvedValue(true);

    render(
      <MockProviders>
        <OrdersInfoForm
          {...civilianTDYTestProps}
          initialValues={{
            ...civilianTDYTestProps.initialValues,
          }}
        />
      </MockProviders>,
    );

    await waitFor(() => {
      expect(screen.getByLabelText(/Orders type/)).toBeInTheDocument();
    });
    await userEvent.selectOptions(screen.getByLabelText(/Orders type/), ORDERS_TYPE_OPTIONS.LOCAL_MOVE);
    await userEvent.selectOptions(screen.getByLabelText(/Pay grade/), ORDERS_PAY_GRADE_TYPE.CIVILIAN_EMPLOYEE);
    await waitFor(() =>
      expect(
        screen.queryByText('If your orders specify a UB weight allowance, enter it here.'),
      ).not.toBeInTheDocument(),
    );
  });

  it('does not render civilian TDY UB Allowance field if grade is not CIVILIAN_EMPLOYEE', async () => {
    isBooleanFlagEnabled.mockResolvedValue(true);

    render(
      <MockProviders>
        <OrdersInfoForm
          {...civilianTDYTestProps}
          initialValues={{
            ...civilianTDYTestProps.initialValues,
          }}
        />
      </MockProviders>,
    );

    await waitFor(() => {
      expect(screen.getByLabelText(/Pay grade/)).toBeInTheDocument();
    });
    await userEvent.selectOptions(screen.getByLabelText(/Pay grade/), [ORDERS_PAY_GRADE_TYPE.E_5]);
    await waitFor(() =>
      expect(
        screen.queryByText('If your orders specify a UB weight allowance, enter it here.'),
      ).not.toBeInTheDocument(),
    );
  });

  it.each([[ORDERS_TYPE.RETIREMENT], [ORDERS_TYPE.SEPARATION]])(
    'renders correct DutyLocationInput label and hint for %s orders type',
    async (ordersType) => {
      render(
        <MockProviders>
          <OrdersInfoForm
            {...testProps}
            initialValues={{
              ...testProps.initialValues,
              orders_type: ordersType,
            }}
          />
        </MockProviders>,
      );

      await userEvent.selectOptions(screen.getByLabelText(/Orders type/), ordersType); // Select the orders type in the dropdown
      const destinationInput = await screen.findByLabelText(/Destination Location \(As Authorized on Orders\)/);
      expect(destinationInput).toBeInTheDocument();
      expect(
        screen.getByText(
          /Enter the option closest to your destination\. Your move counselor will identify if there might be a cost to you\./,
        ),
      ).toBeInTheDocument();
    },
  );

  afterEach(jest.restoreAllMocks);
});
