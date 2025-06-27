import React from 'react';
import { render, waitFor, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { Provider } from 'react-redux';

import AddOrdersForm from './AddOrdersForm';

import { MockProviders } from 'testUtils';
import { dropdownInputOptions } from 'utils/formatters';
import { ORDERS_PAY_GRADE_TYPE, ORDERS_TYPE, ORDERS_TYPE_OPTIONS } from 'constants/orders';
import { configureStore } from 'shared/store';
import { isBooleanFlagEnabled } from 'utils/featureFlags';
import { servicesCounselingRoutes } from 'constants/routes';

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
      {
        address: {
          city: '',
          id: '25be4d12-fe93-47f1-bbec-1db386dfa67e',
          postalCode: '',
          state: '',
          streetAddress1: '',
        },
        address_id: '4334640b-c35e-4293-a2f1-36c7b629f904',
        affiliation: 'AIR_FORCE',
        created_at: '2021-02-11T16:48:04.117Z',
        id: '22f0755f-6f35-478b-9a75-35a69211da1d',
        name: 'Scott AFB',
        updated_at: '2021-02-11T16:48:04.117Z',
        provides_services_counseling: true,
      },
    ]),
  ),
}));

jest.mock('services/ghcApi', () => ({
  ...jest.requireActual('services/ghcApi'),
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
      ],
    });
  }),
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
  counselingOfficeId: '',
};
const testProps = {
  initialValues,
  ordersTypeOptions: dropdownInputOptions(ORDERS_TYPE_OPTIONS),
  onSubmit: jest.fn(),
  onBack: jest.fn(),
};
const mockParams = { customerId: 'ea51dab0-4553-4732-b843-1f33407f77bd' };
const mockPath = servicesCounselingRoutes.BASE_CUSTOMERS_ORDERS_ADD_PATH;

describe('CreateMoveCustomerInfo Component', () => {
  it('renders the form inputs and asterisks for required fields', async () => {
    render(
      <Provider store={mockStore.store}>
        <AddOrdersForm {...testProps} />
      </Provider>,
    );

    await waitFor(() => {
      expect(screen.getByText('Tell us about the orders')).toBeInTheDocument();
      expect(screen.getByLabelText(/Orders type/)).toBeInTheDocument();
      expect(screen.getByLabelText(/Orders date/)).toBeInTheDocument();
      expect(screen.getByLabelText(/Report by date/)).toBeInTheDocument();
      expect(screen.getByText(/Are dependents included in the orders?/)).toBeInTheDocument();
      expect(screen.getByTestId('hasDependentsYes')).toBeInTheDocument();
      expect(screen.getByTestId('hasDependentsNo')).toBeInTheDocument();
      expect(screen.getByLabelText(/Current duty location/)).toBeInTheDocument();
      expect(screen.getByLabelText(/New duty location/)).toBeInTheDocument();
      expect(screen.getByLabelText(/Pay grade/)).toBeInTheDocument();

      expect(screen.getByTestId('reqAsteriskMsg')).toBeInTheDocument();

      // check for asterisks on required fields
      expect(document.querySelector('#reqAsteriskMsg')).toHaveTextContent('Fields marked with * are required.');

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

    const { getByLabelText } = render(
      <Provider store={mockStore.store}>
        <AddOrdersForm {...testProps} />
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

    await userEvent.selectOptions(ordersTypeDropdown, ORDERS_TYPE.WOUNDED_WARRIOR);
    expect(ordersTypeDropdown).toHaveValue(ORDERS_TYPE.WOUNDED_WARRIOR);

    await userEvent.selectOptions(ordersTypeDropdown, ORDERS_TYPE.BLUEBARK);
    expect(ordersTypeDropdown).toHaveValue(ORDERS_TYPE.BLUEBARK);

    // Saftey option should not be available for non safety moves
    const options = ordersTypeDropdown.querySelectorAll('option');
    const isSafetyOptionPresent = Array.from(options).some((option) => option.value === ORDERS_TYPE.SAFETY);
    expect(isSafetyOptionPresent).toBe(false);
  });

  it('shows an error message if trying to submit an invalid form', async () => {
    const { getByRole, findAllByRole, getByLabelText } = render(
      <Provider store={mockStore.store}>
        <AddOrdersForm {...testProps} />
      </Provider>,
    );
    await userEvent.click(getByLabelText(/Orders type/));
    await userEvent.click(getByLabelText(/Orders date/));
    await userEvent.click(getByLabelText(/Report by date/));
    await userEvent.click(getByLabelText(/Current duty location/));
    await userEvent.click(getByLabelText(/New duty location/));
    await userEvent.click(getByLabelText(/Pay grade/));

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
      <Provider params={mockParams} store={mockStore.store}>
        <AddOrdersForm {...testProps} />
      </Provider>,
    );

    await userEvent.selectOptions(await screen.findByLabelText(/Orders type/), 'PERMANENT_CHANGE_OF_STATION');
    await userEvent.type(screen.getByLabelText(/Orders date/), '08 Nov 2020');
    await userEvent.type(screen.getByLabelText(/Report by date/), '26 Nov 2020');
    await userEvent.click(screen.getByLabelText('No'));
    await userEvent.selectOptions(screen.getByLabelText(/Pay grade/), [ORDERS_PAY_GRADE_TYPE.E_5]);

    await userEvent.type(screen.getByLabelText(/Current duty location/), 'AFB');
    await userEvent.click(await screen.findByText(/Elmendorf/));

    const counselingOfficeLabel = await screen.queryByText(/Counseling office/);
    expect(counselingOfficeLabel).toBeFalsy();

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

describe('AddOrdersForm - Student Travel, Early Return of Dependents Test', () => {
  it('has dependents is yes and disabled when order type is student travel', async () => {
    isBooleanFlagEnabled.mockImplementation(() => Promise.resolve(true));

    render(
      <Provider store={mockStore.store}>
        <AddOrdersForm {...testProps} />
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
    isBooleanFlagEnabled.mockImplementation(() => Promise.resolve(true));

    render(
      <Provider store={mockStore.store}>
        <AddOrdersForm {...testProps} />
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
    isBooleanFlagEnabled.mockImplementation(() => Promise.resolve(true));

    render(
      <Provider store={mockStore.store}>
        <AddOrdersForm {...testProps} />
      </Provider>,
    );

    await userEvent.selectOptions(screen.getByLabelText(/Orders type/), ORDERS_TYPE.PERMANENT_CHANGE_OF_STATION);

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

    const hasDependentsYesStudent = screen.getByLabelText('Yes');
    const hasDependentsNoStudent = screen.getByLabelText('No');

    await waitFor(() => {
      expect(hasDependentsYesStudent).toBeChecked();
      expect(hasDependentsYesStudent).toBeDisabled();
      expect(hasDependentsNoStudent).toBeDisabled();
    });

    // set order type to value the re-enables "has dependents"
    await userEvent.selectOptions(screen.getByLabelText(/Orders type/), ORDERS_TYPE.LOCAL_MOVE);

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
    isBooleanFlagEnabled.mockImplementation(() => Promise.resolve(true));

    render(
      <Provider store={mockStore.store}>
        <AddOrdersForm {...testProps} />
      </Provider>,
    );

    await userEvent.selectOptions(screen.getByLabelText(/Orders type/), ORDERS_TYPE.PERMANENT_CHANGE_OF_STATION);

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

    const hasDependentsYesEarly = screen.getByLabelText('Yes');
    const hasDependentsNoEarly = screen.getByLabelText('No');

    await waitFor(() => {
      expect(hasDependentsYesEarly).toBeChecked();
      expect(hasDependentsYesEarly).toBeDisabled();
      expect(hasDependentsNoEarly).toBeDisabled();
    });

    // set order type to value the re-enables "has dependents"
    await userEvent.selectOptions(screen.getByLabelText(/Orders type/), ORDERS_TYPE.LOCAL_MOVE);

    const hasDependentsYesLocalMove = screen.getByLabelText('Yes');
    const hasDependentsNoLocalMove = screen.getByLabelText('No');

    await waitFor(() => {
      expect(hasDependentsYesLocalMove).not.toBeChecked();
      expect(hasDependentsYesLocalMove).toBeEnabled();
      expect(hasDependentsNoLocalMove).not.toBeChecked();
      expect(hasDependentsNoLocalMove).toBeEnabled();
    });
  });
});

describe('AddOrdersForm - Edge Cases and Additional Scenarios', () => {
  it('disables orders type when safety move is selected', async () => {
    render(
      <Provider store={mockStore.store}>
        <AddOrdersForm {...testProps} isSafetyMoveSelected />
      </Provider>,
    );

    expect(screen.getByLabelText(/Orders type/)).toBeDisabled();
  });

  it('disables orders type when bluebark move is selected', async () => {
    render(
      <Provider store={mockStore.store}>
        <AddOrdersForm {...testProps} isBluebarkMoveSelected />
      </Provider>,
    );
    expect(screen.getByLabelText(/Orders type/)).toBeDisabled();
  });
});

describe('AddOrdersForm - With Counseling Office', () => {
  it('displays the counseling office dropdown', async () => {
    isBooleanFlagEnabled.mockImplementation(() => Promise.resolve(true));
    render(
      <MockProviders params={mockParams} path={mockPath} store={mockStore.store}>
        <AddOrdersForm {...testProps} />
      </MockProviders>,
    );

    await userEvent.selectOptions(await screen.findByLabelText(/Orders type/), 'PERMANENT_CHANGE_OF_STATION');
    await userEvent.paste(screen.getByLabelText(/Orders date/), '08 Nov 2020');
    await userEvent.paste(screen.getByLabelText(/Report by date/), '26 Nov 2020');
    await userEvent.click(screen.getByLabelText('No'));
    await userEvent.selectOptions(screen.getByLabelText(/Pay grade/), [ORDERS_PAY_GRADE_TYPE.E_5]);

    // Test Current Duty Location Search Box interaction
    await userEvent.type(screen.getByLabelText(/Current duty location/), 'AFB', { delay: 100 });
    const selectedOptionCurrent = await screen.findByText(/Scott/);
    await userEvent.click(selectedOptionCurrent);

    // Test New Duty Location Search Box interaction
    await userEvent.type(screen.getByLabelText(/New duty location/), 'AFB', { delay: 100 });
    const selectedOptionNew = await screen.findByText(/Luke/);
    await userEvent.click(selectedOptionNew);

    const counselingOfficeLabel = await screen.queryByText(/Counseling office/);
    expect(counselingOfficeLabel).toBeTruthy();

    await userEvent.selectOptions(screen.getByLabelText(/Counseling office/), ['Albuquerque AFB']);

    const nextBtn = screen.getByRole('button', { name: 'Next' });
    expect(nextBtn.getAttribute('disabled')).toBeFalsy();
  });

  it('disabled submit if counseling office is required and blank', async () => {
    isBooleanFlagEnabled.mockImplementation(() => Promise.resolve(true));
    render(
      <MockProviders params={mockParams} path={mockPath} store={mockStore.store}>
        <AddOrdersForm {...testProps} />
      </MockProviders>,
    );

    await userEvent.selectOptions(await screen.findByLabelText(/Orders type/), 'PERMANENT_CHANGE_OF_STATION');
    await userEvent.type(screen.getByLabelText(/Orders date/), '08 Nov 2024');
    await userEvent.type(screen.getByLabelText(/Report by date/), '26 Nov 2024');

    // Test Current Duty Location Search Box interaction
    await userEvent.type(screen.getByLabelText(/Current duty location/), 'AFB', { delay: 100 });
    const selectedOptionCurrent = await screen.findByText(/Scott/);
    await userEvent.click(selectedOptionCurrent);

    // Test New Duty Location Search Box interaction
    await userEvent.type(screen.getByLabelText(/New duty location/), 'AFB', { delay: 100 });
    const selectedOptionNew = await screen.findByText(/Luke/);
    await userEvent.click(selectedOptionNew);

    const counselingOfficeLabel = await screen.queryByText(/Counseling office/);
    expect(counselingOfficeLabel).toBeTruthy(); // If the field is visible then it it required

    await userEvent.selectOptions(screen.getByLabelText(/Pay grade/), [ORDERS_PAY_GRADE_TYPE.E_5]);
    await userEvent.click(screen.getByLabelText('No'));

    const nextBtn = await screen.getByRole('button', { name: 'Next' }, { delay: 100 });
    expect(nextBtn).toBeDisabled();
  });
});
