import React from 'react';
import { screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { act } from 'react-dom/test-utils';

import AddOrders from './AddOrders';

import { createOrders, getServiceMember, showCounselingOffices } from 'services/internalApi';
import { renderWithProviders } from 'testUtils';
import { customerRoutes, generalRoutes } from 'constants/routes';
import { selectCanAddOrders, selectServiceMemberFromLoggedInUser } from 'store/entities/selectors';
import { setCanAddOrders, setMoveId } from 'store/general/actions';
import { isBooleanFlagEnabled } from 'utils/featureFlags';
import { ORDERS_TYPE } from 'constants/orders';

// Tests are timing out. High assumption it is due to service counseling office drop-down choice not being loaded on initial form load. It's another API call
jest.setTimeout(60000);

jest.mock('services/internalApi', () => ({
  ...jest.requireActual('services/internalApi'),
  getServiceMember: jest.fn().mockImplementation(() => Promise.resolve()),
  getResponseError: jest.fn().mockImplementation(() => Promise.resolve()),
  createOrders: jest.fn().mockImplementation(() => Promise.resolve()),
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

jest.mock('store/entities/selectors', () => ({
  ...jest.requireActual('store/entities/selectors'),
  selectServiceMemberFromLoggedInUser: jest.fn(),
  selectCanAddOrders: jest.fn(),
  selectMoveId: jest.fn(),
}));

jest.mock('store/general/actions', () => ({
  ...jest.requireActual('store/general/actions'),
  setCanAddOrders: jest.fn().mockImplementation(() => ({
    type: '',
    payload: '',
  })),
  setMoveId: jest.fn().mockImplementation(() => ({
    type: '',
    payload: '',
  })),
}));

const mockNavigate = jest.fn();
jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom'),
  useNavigate: () => mockNavigate,
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

const serviceMember = {
  id: 'id123',
};

jest.mock('utils/featureFlags', () => ({
  ...jest.requireActual('utils/featureFlags'),
  isBooleanFlagEnabled: jest.fn().mockImplementation(() => Promise.resolve(false)),
}));

describe('Add Orders page', () => {
  const testProps = {
    serviceMemberId: 'id123',
    context: { flags: { allOrdersTypes: true } },
    canAddOrders: true,
    moveId: '',
    updateOrders: jest.fn(),
    updateServiceMember: jest.fn(),
    setCanAddOrders: jest.fn(),
    setMoveId: jest.fn(),
  };

  const testPropsRedirect = {
    serviceMemberId: 'id123',
    context: { flags: { allOrdersTypes: true } },
    canAddOrders: false,
    moveId: '',
    updateOrders: jest.fn(),
    updateServiceMember: jest.fn(),
    setCanAddOrders: jest.fn(),
    setMoveId: jest.fn(),
  };

  it('renders all content of Orders component', async () => {
    showCounselingOffices.mockImplementation(() => Promise.resolve({}));
    selectServiceMemberFromLoggedInUser.mockImplementation(() => serviceMember);
    renderWithProviders(<AddOrders {...testProps} />, {
      path: customerRoutes.ORDERS_ADD_PATH,
    });

    await screen.findByRole('heading', { level: 1, name: 'Tell us about your move orders' });
    expect(screen.getByTestId('main-container')).toBeInTheDocument();
    expect(screen.getByTestId('orders-form-container')).toBeInTheDocument();
    const saveBtn = await screen.findByRole('button', { name: 'Back' });
    expect(saveBtn).toBeInTheDocument();
    const cancelBtn = await screen.findByRole('button', { name: 'Next' });
    expect(cancelBtn).toBeInTheDocument();
  });

  it('renders all fields on load', async () => {
    selectServiceMemberFromLoggedInUser.mockImplementation(() => serviceMember);
    renderWithProviders(<AddOrders {...testProps} />, {
      path: customerRoutes.ORDERS_ADD_PATH,
    });

    await screen.findByRole('heading', { level: 1, name: 'Tell us about your move orders' });
    expect(screen.getByLabelText(/Orders type/)).toBeInTheDocument();
    expect(screen.getByLabelText(/Orders date/)).toBeInTheDocument();
    expect(screen.getByLabelText(/Report by date/)).toBeInTheDocument();
    expect(screen.getByText('Are dependents included in your orders?')).toBeInTheDocument();
    expect(screen.getByLabelText(/Current duty location/)).toBeInTheDocument();
    expect(screen.getByLabelText(/New duty location/)).toBeInTheDocument();
    expect(screen.getByLabelText(/Pay grade/)).toBeInTheDocument();

    const backBtn = await screen.findByRole('button', { name: 'Back' });
    expect(backBtn).toBeInTheDocument();
    expect(backBtn).toBeEnabled();

    const nextBtn = await screen.findByRole('button', { name: 'Next' });
    expect(nextBtn).toBeInTheDocument();
    expect(nextBtn).toBeDisabled();
  });

  it('does not render conditional dependent fields on load', async () => {
    selectServiceMemberFromLoggedInUser.mockImplementation(() => serviceMember);
    renderWithProviders(<AddOrders {...testProps} />, {
      path: customerRoutes.ORDERS_ADD_PATH,
    });

    await screen.findByRole('heading', { level: 1, name: 'Tell us about your move orders' });
    expect(screen.queryByText('Is this an accompanied tour?')).not.toBeInTheDocument();
    expect(screen.queryByLabelText(/Number of dependents under the age of 12/)).not.toBeInTheDocument();
    expect(screen.queryByLabelText(/Number of dependents of the age 12 or over/)).not.toBeInTheDocument();
    expect(
      screen.queryByText(
        'Unaccompanied Tour: An authorized order (assignment or tour) that DOES NOT allow dependents to travel to the new Permanent Duty Station (PDS)',
      ),
    ).not.toBeInTheDocument();
    expect(
      screen.queryByText(
        'Accompanied Tour: An authorized order (assignment or tour) that allows dependents to travel to the new Permanent Duty Station (PDS)',
      ),
    ).not.toBeInTheDocument();
  });

  it('does not render the input boxes for number of dependents over or under 12 if both locations are CONUS', async () => {
    selectServiceMemberFromLoggedInUser.mockImplementation(() => serviceMember);
    renderWithProviders(<AddOrders {...testProps} />, {
      path: customerRoutes.ORDERS_ADD_PATH,
    });

    await screen.findByRole('heading', { level: 1, name: 'Tell us about your move orders' });
    // Select a CONUS current duty location and new duty location
    await userEvent.type(screen.getByLabelText(/Current duty location/), 'AFB', { delay: 100 });
    const selectedOptionCurrent = await screen.findByText(/Altus/);
    await userEvent.click(selectedOptionCurrent);
    await userEvent.type(screen.getByLabelText(/New duty location/), 'AFB', { delay: 100 });
    const selectedOptionNew = await screen.findByText(/Luke/);
    await userEvent.click(selectedOptionNew);

    // Select that dependents are present
    await userEvent.click(screen.getByTestId('hasDependentsYes'));

    // With both addresses being CONUS, the number of dependents input boxes should be missing
    expect(screen.queryByLabelText(/Number of dependents under the age of 12/)).not.toBeInTheDocument();
    expect(screen.queryByLabelText(/Number of dependents of the age 12 or over/)).not.toBeInTheDocument();
  });

  it('does render the input boxes for number of dependents over or under 12 if one of the locations are OCONUS', async () => {
    isBooleanFlagEnabled.mockImplementation(() => Promise.resolve(true));
    selectServiceMemberFromLoggedInUser.mockImplementation(() => serviceMember);
    renderWithProviders(<AddOrders {...testProps} />, {
      path: customerRoutes.ORDERS_ADD_PATH,
    });

    await screen.findByRole('heading', { level: 1, name: 'Tell us about your move orders' });
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
    expect(screen.getByLabelText(/Number of dependents under the age of 12/)).toBeInTheDocument();
    expect(screen.getByLabelText(/Number of dependents of the age 12 or over/)).toBeInTheDocument();
  });

  it('only renders dependents age groupings and accompanied tour if dependents are present', async () => {
    isBooleanFlagEnabled.mockImplementation(() => Promise.resolve(true));
    selectServiceMemberFromLoggedInUser.mockImplementation(() => serviceMember);
    renderWithProviders(<AddOrders {...testProps} />, {
      path: customerRoutes.ORDERS_ADD_PATH,
    });

    await screen.findByRole('heading', { level: 1, name: 'Tell us about your move orders' });
    // Select a CONUS current duty location
    await userEvent.type(screen.getByLabelText(/Current duty location/), 'AFB', { delay: 100 });
    const selectedOptionCurrent = await screen.findByText(/Altus/);
    await userEvent.click(selectedOptionCurrent);
    // Select an OCONUS new duty location
    await userEvent.type(screen.getByLabelText(/New duty location/), 'AFB', { delay: 100 });
    const selectedOptionNew = await screen.findByText(/Elmendorf/);
    await userEvent.click(selectedOptionNew);
    // Select that dependents are present
    await userEvent.click(screen.getByTestId('hasDependentsNo'));
    // With one of the duty locations being OCONUS, the number of dependents input boxes should be present
    expect(screen.queryByLabelText(/Number of dependents under the age of 12/)).not.toBeInTheDocument();
    expect(screen.queryByLabelText(/Number of dependents of the age 12 or over/)).not.toBeInTheDocument();
    expect(screen.queryByLabelText(/Is this an accompanied tour?/)).not.toBeInTheDocument();
  });

  it('next button creates the orders and updates state', async () => {
    const testOrdersValues = {
      id: 'testOrdersId',
      orders_type: ORDERS_TYPE.PERMANENT_CHANGE_OF_STATION,
      issue_date: '2020-11-08',
      report_by_date: '2020-11-26',
      has_dependents: false,
      moves: ['testMovId'],
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
    };

    selectServiceMemberFromLoggedInUser.mockImplementation(() => serviceMember);
    createOrders.mockImplementation(() => Promise.resolve(testOrdersValues));
    getServiceMember.mockImplementation(() => Promise.resolve());

    await act(async () => {
      renderWithProviders(<AddOrders {...testProps} />, {
        path: customerRoutes.ORDERS_ADD_PATH,
      });
    });

    const nextBtn = await screen.findByRole('button', { name: 'Next' });
    expect(nextBtn).toBeInTheDocument();

    await act(async () => {
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
    });

    await waitFor(() => {
      expect(screen.getByRole('form')).toHaveFormValues({
        new_duty_location: 'Luke AFB',
        origin_duty_location: 'Altus AFB',
      });
    });

    await waitFor(() => expect(nextBtn).toBeEnabled());

    await act(async () => {
      await userEvent.click(nextBtn);
    });

    await waitFor(() => {
      expect(createOrders).toHaveBeenCalled();
      expect(setMoveId).toHaveBeenCalled();
      expect(setCanAddOrders).toHaveBeenCalled();
      expect(getServiceMember).toHaveBeenCalledWith(testProps.serviceMemberId);
    });
  });

  it('submits OCONUS fields correctly on form submit', async () => {
    const testOrdersValues = {
      orders_type: 'PERMANENT_CHANGE_OF_STATION',
      issue_date: '2020-11-08',
      report_by_date: '2020-11-26',
      has_dependents: true,
      accompanied_tour: true,
      dependents_under_twelve: 1,
      dependents_twelve_and_over: 2,
      counseling_office_id: null,
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
      new_duty_location_id: 'a8d6b33c-8370-4e92-8df2-356b8c9d0c1a',
      new_duty_location: {
        address: {
          city: 'Elmendorf AFB',
          country: 'US',
          isOconus: true,
          id: 'fa51dab0-4553-4732-b843-1f33407f11bc',
          postalCode: '78112',
          state: 'AK',
          streetAddress1: 'n/a',
        },
        affiliation: 'AIR_FORCE',
        created_at: '2021-02-11T16:48:04.117Z',
        id: 'a8d6b33c-8370-4e92-8df2-356b8c9d0c1a',
        name: 'Elmendorf AFB',
        updated_at: '2021-02-11T16:48:04.117Z',
        address_id: 'fa51dab0-4553-4732-b843-1f33407f11bc',
      },
      grade: 'E_5',
      origin_duty_location_id: '93f0755f-6f35-478b-9a75-35a69211da1c',
      service_member_id: 'id123',
      spouse_has_pro_gear: false,
    };
    isBooleanFlagEnabled.mockImplementation(() => Promise.resolve(true));
    selectServiceMemberFromLoggedInUser.mockImplementation(() => serviceMember);
    renderWithProviders(<AddOrders {...testProps} />, {
      path: customerRoutes.ORDERS_ADD_PATH,
    });
    await screen.findByRole('heading', { level: 1, name: 'Tell us about your move orders' });

    const nextBtn = await screen.findByRole('button', { name: 'Next' });
    expect(nextBtn).toBeInTheDocument();

    // Set standard form fields
    await act(async () => {
      await userEvent.selectOptions(screen.getByLabelText(/Orders type/), 'PERMANENT_CHANGE_OF_STATION');
      await userEvent.type(screen.getByLabelText(/Orders date/), '08 Nov 2020');
      await userEvent.type(screen.getByLabelText(/Report by date/), '26 Nov 2020');
      await userEvent.click(screen.getByLabelText('No'));
      await userEvent.selectOptions(screen.getByLabelText(/Pay grade/), ['E_5']);

      // Select a CONUS current duty location
      await userEvent.type(screen.getByLabelText(/Current duty location/), 'AFB', { delay: 100 });
      const selectedOptionCurrent = await screen.findByText(/Altus/);
      await userEvent.click(selectedOptionCurrent);
      // Select an OCONUS new duty location
      await userEvent.type(screen.getByLabelText(/New duty location/), 'AFB', { delay: 100 });
      const selectedOptionNew = await screen.findByText(/Elmendorf/);
      await userEvent.click(selectedOptionNew);
    });

    await waitFor(() => {
      expect(screen.getByRole('form')).toHaveFormValues({
        new_duty_location: 'Elmendorf AFB',
        origin_duty_location: 'Altus AFB',
      });
    });

    // Set dependents and accompanied tour
    await userEvent.click(screen.getByTestId('hasDependentsYes'));
    await userEvent.click(screen.getByTestId('isAnAccompaniedTourYes'));
    await userEvent.type(screen.getByTestId('dependentsUnderTwelve'), '1');
    await userEvent.type(screen.getByTestId('dependentsTwelveAndOver'), '2');

    await waitFor(() => expect(nextBtn).toBeEnabled());

    await act(async () => {
      await userEvent.click(nextBtn);
    });

    await waitFor(() => {
      expect(createOrders).toHaveBeenCalledWith(testOrdersValues);
    });
  });

  it('properly does not pass in OCONUS fields when is a CONUS move', async () => {
    const testOrdersValues = {
      orders_type: 'PERMANENT_CHANGE_OF_STATION',
      issue_date: '2020-11-08',
      report_by_date: '2020-11-26',
      has_dependents: false,
      counseling_office_id: null,
      dependents_twelve_and_over: null,
      dependents_under_twelve: null,
      accompanied_tour: null,
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
      new_duty_location_id: 'a8d6b33c-8370-4e92-8df2-356b8c9d0c1a',
      new_duty_location: {
        address: {
          city: 'Glendale Luke AFB',
          country: 'United States',
          id: 'fa51dab0-4553-4732-b843-1f33407f77bc',
          postalCode: '85309',
          state: 'AZ',
          streetAddress1: 'n/a',
        },
        affiliation: 'AIR_FORCE',
        created_at: '2021-02-11T16:48:04.117Z',
        id: 'a8d6b33c-8370-4e92-8df2-356b8c9d0c1a',
        name: 'Luke AFB',
        updated_at: '2021-02-11T16:48:04.117Z',
        address_id: '25be4d12-fe93-47f1-bbec-1db386dfa67f',
      },
      grade: 'E_5',
      origin_duty_location_id: '93f0755f-6f35-478b-9a75-35a69211da1c',
      service_member_id: 'id123',
      spouse_has_pro_gear: false,
    };
    isBooleanFlagEnabled.mockImplementation(() => Promise.resolve(true));
    selectServiceMemberFromLoggedInUser.mockImplementation(() => serviceMember);
    renderWithProviders(<AddOrders {...testProps} />, {
      path: customerRoutes.ORDERS_ADD_PATH,
    });
    await screen.findByRole('heading', { level: 1, name: 'Tell us about your move orders' });

    const nextBtn = await screen.findByRole('button', { name: 'Next' });
    expect(nextBtn).toBeInTheDocument();

    // Set standard form fields
    await act(async () => {
      await userEvent.selectOptions(screen.getByLabelText(/Orders type/), 'PERMANENT_CHANGE_OF_STATION');
      await userEvent.type(screen.getByLabelText(/Orders date/), '08 Nov 2020');
      await userEvent.type(screen.getByLabelText(/Report by date/), '26 Nov 2020');
      await userEvent.click(screen.getByLabelText('No'));
      await userEvent.selectOptions(screen.getByLabelText(/Pay grade/), ['E_5']);

      // Select a CONUS current duty location
      await userEvent.type(screen.getByLabelText(/Current duty location/), 'AFB', { delay: 100 });
      const selectedOptionCurrent = await screen.findByText(/Altus/);
      await userEvent.click(selectedOptionCurrent);
      // Select an CONUS new duty location
      await userEvent.type(screen.getByLabelText(/New duty location/), 'AFB', { delay: 100 });
      const selectedOptionNew = await screen.findByText(/Luke/);
      await userEvent.click(selectedOptionNew);
    });

    await waitFor(() => {
      expect(screen.getByRole('form')).toHaveFormValues({
        new_duty_location: 'Luke AFB',
        origin_duty_location: 'Altus AFB',
      });
    });

    await waitFor(() => expect(nextBtn).toBeEnabled());

    await act(async () => {
      await userEvent.click(nextBtn);
    });

    await waitFor(() => {
      expect(createOrders).toHaveBeenCalledWith(testOrdersValues);
    });
  });

  it('redirects the user if canAddOrders is false', async () => {
    selectServiceMemberFromLoggedInUser.mockImplementation(() => serviceMember);
    selectCanAddOrders.mockImplementation(() => false);
    renderWithProviders(<AddOrders {...testPropsRedirect} />, {
      path: customerRoutes.ORDERS_ADD_PATH,
    });

    expect(mockNavigate).toHaveBeenCalledWith(generalRoutes.HOME_PATH);
  });
});
