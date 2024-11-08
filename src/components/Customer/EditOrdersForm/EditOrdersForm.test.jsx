import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import EditOrdersForm from './EditOrdersForm';

import { documentSizeLimitMsg } from 'shared/constants';
import { showCounselingOffices } from 'services/internalApi';
import { ORDERS_TYPE } from 'constants/orders';

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
        provides_services_counseling: true,
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
        provides_services_counseling: true,
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
    uploaded_orders: [],
    origin_duty_location: {
      provides_services_counseling: true,
    },
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
  ],
  currentDutyLocation: {},
  grade: '',
};

const initialValues = {
  orders_type: 'PERMANENT_CHANGE_OF_STATION',
  issue_date: '2020-11-08',
  report_by_date: '2020-11-26',
  has_dependents: 'no',
  origin_duty_location: {
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
    provides_services_counseling: true,
  },
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
  uploaded_orders: [
    {
      id: '123',
      createdAt: '2020-11-08',
      bytes: 1,
      url: 'url',
      filename: 'Test Upload',
      contentType: 'application/pdf',
    },
  ],
  grade: 'E_1',
  accompanied_tour: '',
  dependents_under_twelve: '',
  dependents_twelve_and_over: '',
};

jest.mock('utils/featureFlags', () => ({
  ...jest.requireActual('utils/featureFlags'),
  isBooleanFlagEnabled: jest.fn().mockImplementation(() => Promise.resolve(false)),
}));

describe('EditOrdersForm component', () => {
  describe('renders each input and checks if the field is required', () => {
    it.each([
      [/Orders type/, true, HTMLSelectElement],
      [/Orders date/, true, HTMLInputElement],
      [/Report by date/, true, HTMLInputElement],
      ['Yes', false, HTMLInputElement],
      ['No', false, HTMLInputElement],
      [/New duty location/, false, HTMLInputElement],
      [/Pay grade/, true, HTMLSelectElement],
      [/Current duty location/, false, HTMLInputElement],
    ])('rendering %s and is required is %s', async (formInput, required, inputType) => {
      render(<EditOrdersForm {...testProps} />);

      expect(await screen.findByLabelText(formInput)).toBeInstanceOf(inputType);
      if (required) {
        expect(await screen.findByLabelText(formInput)).toBeRequired();
      }
    });

    it('rendering the upload area', async () => {
      showCounselingOffices.mockImplementation(() => Promise.resolve({}));

      render(<EditOrdersForm {...testProps} />);

      expect(await screen.findByText(documentSizeLimitMsg)).toBeInTheDocument();
    });
  });

  describe('renders each option for the orders type dropdown', () => {
    it.each([
      ['PERMANENT_CHANGE_OF_STATION', 'PERMANENT_CHANGE_OF_STATION'],
      ['LOCAL_MOVE', 'LOCAL_MOVE'],
      ['RETIREMENT', 'RETIREMENT'],
      ['SEPARATION', 'SEPARATION'],
    ])('rendering the %s option', async (selectionOption, expectedValue) => {
      render(<EditOrdersForm {...testProps} />);

      const ordersTypeDropdown = await screen.findByLabelText(/Orders type/);
      expect(ordersTypeDropdown).toBeInstanceOf(HTMLSelectElement);

      await userEvent.selectOptions(ordersTypeDropdown, selectionOption);
      await waitFor(() => {
        expect(ordersTypeDropdown).toHaveValue(expectedValue);
      });
    });
  });

  it('allows new and current duty location to be the same', async () => {
    // Render the component
    render(
      <EditOrdersForm
        {...testProps}
        initialValues={{
          ...initialValues,
          origin_duty_location: {
            name: 'Luke AFB',
            provides_services_counseling: false,
            address: { isOconus: false },
          },
          new_duty_location: {
            name: 'Luke AFB',
            provides_services_counseling: false,
            address: { isOconus: false },
          },
          counseling_office_id: '3e937c1f-5539-4919-954d-017989130584',
          uploaded_orders: [
            {
              id: '123',
              createdAt: '2020-11-08',
              bytes: 1,
              url: 'url',
              filename: 'Test Upload',
              contentType: 'application/pdf',
            },
          ],
        }}
      />,
    );

    await waitFor(() => expect(screen.queryByText('Loading, please wait...')).not.toBeInTheDocument());

    const submitButton = screen.getByRole('button', { name: 'Save' });
    await waitFor(() => {
      expect(submitButton).not.toBeDisabled();
    });

    await userEvent.selectOptions(screen.getByLabelText(/Orders type/), 'PERMANENT_CHANGE_OF_STATION');
    await userEvent.type(screen.getByLabelText(/Orders date/), '08 Nov 2020');
    await userEvent.type(screen.getByLabelText(/Report by date/), '26 Nov 2020');
    await userEvent.click(screen.getByLabelText('No'));
    await userEvent.selectOptions(screen.getByLabelText(/Pay grade/), ['E_5']);
    await userEvent.click(screen.getByTestId('hasDependentsYes'));

    await waitFor(() => {
      expect(screen.getByRole('form')).toHaveFormValues({
        new_duty_location: 'Luke AFB',
        origin_duty_location: 'Luke AFB',
      });
    });
  });

  it('shows an error message if the form is invalid', async () => {
    render(<EditOrdersForm {...testProps} initialValues={initialValues} />);
    const submitButton = await screen.findByRole('button', { name: 'Save' });

    const ordersTypeDropdown = screen.getByLabelText(/Orders type/);
    await userEvent.selectOptions(ordersTypeDropdown, '');
    await userEvent.tab();

    await waitFor(() => {
      expect(submitButton).toBeDisabled();
    });

    const required = screen.getByTestId('errorMessage');
    expect(required).toBeInTheDocument();
  });

  it('submits the form when its valid', async () => {
    // Not testing the upload interaction, so give uploaded orders to the props.
    render(
      <EditOrdersForm
        {...testProps}
        initialValues={{
          origin_duty_location: {
            provides_services_counseling: true,
          },
          counseling_office_id: '3e937c1f-5539-4919-954d-017989130584',
          uploaded_orders: [
            {
              id: '123',
              createdAt: '2020-11-08',
              bytes: 1,
              url: 'url',
              filename: 'Test Upload',
              contentType: 'application/pdf',
            },
          ],
        }}
      />,
    );

    await waitFor(() => expect(screen.queryByText('Loading, please wait...')).not.toBeInTheDocument());

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

    await waitFor(() =>
      expect(screen.getByRole('form')).toHaveFormValues({
        new_duty_location: 'Luke AFB',
        origin_duty_location: 'Altus AFB',
      }),
    );

    const submitBtn = screen.getByRole('button', { name: 'Save' });
    expect(submitBtn).not.toBeDisabled();
    await userEvent.click(submitBtn);

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
          provides_services_counseling: true,
        },
        origin_duty_location: expect.any(Object),
        grade: 'E_5',
        counseling_office_id: '3e937c1f-5539-4919-954d-017989130584',
        uploaded_orders: expect.arrayContaining([
          expect.objectContaining({
            id: '123',
            createdAt: '2020-11-08',
            bytes: 1,
            url: 'url',
            filename: 'Test Upload',
            contentType: 'application/pdf',
          }),
        ]),
      }),
      expect.anything(),
    );
  });

  it('implements the onCancel handler when the Cancel button is clicked', async () => {
    render(<EditOrdersForm {...testProps} />);

    const cancelButton = await screen.findByRole('button', { name: 'Cancel' });

    await userEvent.click(cancelButton);

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
      uploaded_orders: [
        {
          id: '123',
          createdAt: '2020-11-08',
          bytes: 1,
          url: 'url',
          filename: 'Test Upload',
          contentType: 'application/pdf',
        },
      ],
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
      render(<EditOrdersForm {...testProps} initialValues={testInitialValues} />);

      expect(await screen.findByRole('form')).toHaveFormValues({
        new_duty_location: 'Yuma AFB',
        origin_duty_location: 'Altus AFB',
      });

      expect(screen.getByLabelText(/Orders type/)).toHaveValue(testInitialValues.orders_type);
      expect(screen.getByLabelText(/Orders date/)).toHaveValue('08 Nov 2020');
      expect(screen.getByLabelText(/Report by date/)).toHaveValue('26 Nov 2020');
      expect(screen.getByLabelText('Yes')).not.toBeChecked();
      expect(screen.getByLabelText('No')).toBeChecked();
      expect(screen.getByText('Yuma AFB')).toBeInTheDocument();
      expect(screen.getByLabelText(/Pay grade/)).toHaveValue(testInitialValues.grade);
      expect(screen.getByText('Altus AFB')).toBeInTheDocument();
    });

    it('renders the uploads table with an existing upload', async () => {
      render(<EditOrdersForm {...testProps} initialValues={testInitialValues} />);

      await waitFor(() => {
        expect(screen.queryByText('Test Upload')).toBeInTheDocument();
      });
    });
  });

  describe('disables the save button', () => {
    it.each([
      ['Orders Type', 'orders_type', ''],
      [/Orders date/, 'issue_date', ''],
      [/Report by date/, 'report_by_date', ''],
      ['Duty Location', 'new_duty_location', null],
      ['Uploaded Orders', 'uploaded_orders', []],
      [/Pay grade/, 'grade', ''],
    ])('when there is no %s', async (attributeNamePrettyPrint, attributeName, valueToReplaceIt) => {
      const modifiedProps = {
        onSubmit: jest.fn().mockImplementation(() => Promise.resolve()),
        initialValues: {
          orders_type: 'PERMANENT_CHANGE_OF_STATION',
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
          origin_duty_location: {
            provides_services_counseling: true,
          },
          uploaded_orders: [
            {
              id: '123',
              createdAt: '2020-11-08',
              bytes: 1,
              url: 'url',
              filename: 'Test Upload',
              contentType: 'application/pdf',
            },
          ],
          grade: 'E_1',
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
        ],
        currentDutyLocation: {},
      };

      modifiedProps.initialValues[attributeName] = valueToReplaceIt;

      render(<EditOrdersForm {...modifiedProps} />);

      const save = await screen.findByRole('button', { name: 'Save' });
      await waitFor(() => {
        expect(save).toBeInTheDocument();
      });

      expect(save).toBeDisabled();
    });
  });

  it('submits the form when temporary duty orders type is selected', async () => {
    // Not testing the upload interaction, so give uploaded orders to the props.
    render(
      <EditOrdersForm
        {...testProps}
        initialValues={{
          origin_duty_location: {
            provides_services_counseling: true,
          },
          uploaded_orders: [
            {
              id: '123',
              createdAt: '2020-11-08',
              bytes: 1,
              url: 'url',
              filename: 'Test Upload',
              contentType: 'application/pdf',
            },
          ],
          counseling_office_id: '3e937c1f-5539-4919-954d-017989130584',
        }}
      />,
    );

    await waitFor(() => expect(screen.queryByText('Loading, please wait...')).not.toBeInTheDocument());

    await userEvent.selectOptions(screen.getByLabelText(/Orders type/), 'TEMPORARY_DUTY');
    await userEvent.type(screen.getByLabelText(/Orders date/), '28 Oct 2024');
    await userEvent.type(screen.getByLabelText(/Report by date/), '28 Oct 2024');
    await userEvent.click(screen.getByLabelText('No'));
    await userEvent.selectOptions(screen.getByLabelText(/Pay grade/), ['E_8']);

    // Test Current Duty Location Search Box interaction
    await userEvent.type(screen.getByLabelText(/Current duty location/), 'AFB', { delay: 200 });
    const selectedOptionCurrent = await screen.findByText(/Altus/);
    await userEvent.click(selectedOptionCurrent);

    // Test New Duty Location Search Box interaction
    await userEvent.type(screen.getByLabelText(/New duty location/), 'AFB', { delay: 200 });
    const selectedOptionNew = await screen.findByText(/Luke/);
    await userEvent.click(selectedOptionNew);

    await waitFor(() =>
      expect(screen.getByRole('form')).toHaveFormValues({
        new_duty_location: 'Luke AFB',
        origin_duty_location: 'Altus AFB',
      }),
    );

    const submitBtn = screen.getByRole('button', { name: 'Save' });
    expect(submitBtn).not.toBeDisabled();
    await userEvent.click(submitBtn);

    await waitFor(() => {
      expect(testProps.onSubmit).toHaveBeenCalledWith(
        expect.objectContaining({
          orders_type: ORDERS_TYPE.TEMPORARY_DUTY,
        }),
        expect.anything(),
      );
    });
  });

  afterEach(jest.restoreAllMocks);
});
