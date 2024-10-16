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

  it('next button creates the orders and updates state', async () => {
    const testOrdersValues = {
      id: 'testOrdersId',
      orders_type: 'PERMANENT_CHANGE_OF_STATION',
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
      await userEvent.selectOptions(screen.getByLabelText(/Orders type/), 'PERMANENT_CHANGE_OF_STATION');
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

  it('redirects the user if canAddOrders is false', async () => {
    selectServiceMemberFromLoggedInUser.mockImplementation(() => serviceMember);
    selectCanAddOrders.mockImplementation(() => false);
    renderWithProviders(<AddOrders {...testPropsRedirect} />, {
      path: customerRoutes.ORDERS_ADD_PATH,
    });

    expect(mockNavigate).toHaveBeenCalledWith(generalRoutes.HOME_PATH);
  });
});
