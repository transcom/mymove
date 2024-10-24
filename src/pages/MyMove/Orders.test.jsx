import React from 'react';
import { screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import Orders from './Orders';

import { getOrders, patchOrders } from 'services/internalApi';
import { renderWithProviders } from 'testUtils';
import { customerRoutes } from 'constants/routes';
import {
  selectAllMoves,
  selectOrdersForLoggedInUser,
  selectServiceMemberFromLoggedInUser,
} from 'store/entities/selectors';

jest.mock('services/internalApi', () => ({
  ...jest.requireActual('services/internalApi'),
  patchOrders: jest.fn().mockImplementation(() => Promise.resolve()),
  getOrders: jest.fn().mockImplementation(() => Promise.resolve()),
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
          city: '',
          id: '00000000-0000-0000-0000-000000000000',
          postalCode: '',
          state: '',
          streetAddress1: '',
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

jest.mock('store/entities/selectors', () => ({
  ...jest.requireActual('store/entities/selectors'),
  selectServiceMemberFromLoggedInUser: jest.fn(),
  selectOrdersForLoggedInUser: jest.fn(),
  selectAllMoves: jest.fn(),
}));

const mockNavigate = jest.fn();
jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom'),
  useNavigate: () => mockNavigate,
}));

afterEach(() => {
  jest.resetAllMocks();
});

const testPropsWithUploads = {
  id: 'testOrderId',
  orders_type: 'PERMANENT_CHANGE_OF_STATION',
  issue_date: '2020-11-08',
  report_by_date: '2020-11-26',
  has_dependents: false,
  grade: 'E_8',
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
    address: {
      city: 'Altus AFB',
      country: 'United States',
      id: 'fa51dab0-4553-4732-b843-1f33407f77bd',
      postalCode: '73523',
      state: 'OK',
      streetAddress1: 'n/a',
    },
    address_id: 'fa51dab0-4553-4732-b843-1f33407f77bd',
    affiliation: 'AIR_FORCE',
    created_at: '2021-02-11T16:48:04.117Z',
    id: '93f0755f-6f35-478b-9a75-35a69211da1c',
    name: 'Altus AFB',
    updated_at: '2021-02-11T16:48:04.117Z',
  },
  uploaded_orders: {
    id: 'testId',
    uploads: [
      {
        bytes: 1578588,
        contentType: 'image/png',
        createdAt: '2024-02-23T16:51:45.504Z',
        filename: 'Screenshot 2024-02-15 at 12.22.53 PM (2).png',
        id: 'fd88b0e6-ff6d-4a99-be6f-49458a244209',
        status: 'PROCESSING',
        updatedAt: '2024-02-23T16:51:45.504Z',
        url: '/storage/user/5fe4d948-aa1c-4823-8967-b1fb40cf6679/uploads/fd88b0e6-ff6d-4a99-be6f-49458a244209?contentType=image%2Fpng',
      },
    ],
  },
  moves: ['testMoveId'],
};

const testPropsNoUploads = {
  id: 'testOrderId2',
  orders_type: 'PERMANENT_CHANGE_OF_STATION',
  issue_date: '2020-11-08',
  report_by_date: '2020-11-26',
  has_dependents: false,
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
    address: {
      city: 'Altus AFB',
      country: 'United States',
      id: 'fa51dab0-4553-4732-b843-1f33407f77bd',
      postalCode: '73523',
      state: 'OK',
      streetAddress1: 'n/a',
    },
    address_id: 'fa51dab0-4553-4732-b843-1f33407f77bd',
    affiliation: 'AIR_FORCE',
    created_at: '2021-02-11T16:48:04.117Z',
    id: '93f0755f-6f35-478b-9a75-35a69211da1c',
    name: 'Altus AFB',
    updated_at: '2021-02-11T16:48:04.117Z',
  },
  uploaded_orders: {
    id: 'testId',
    service_member_id: 'testId',
    uploads: [],
  },
  moves: ['testMoveId'],
};

const testOrders = [testPropsWithUploads, testPropsNoUploads];

const serviceMember = {
  id: 'id123',
};

describe('Orders page', () => {
  const testProps = {
    serviceMemberId: 'id123',
    context: { flags: { allOrdersTypes: true } },
    orders: testOrders,
    serviceMemberMoves: {
      currentMove: [
        {
          createdAt: '2024-02-23T19:30:11.374Z',
          eTag: 'MjAyNC0wMi0yM1QxOTozMDoxMS4zNzQxN1o=',
          id: 'testMoveId',
          moveCode: '44649B',
          orders: {
            authorizedWeight: 11000,
            created_at: '2024-02-23T19:30:11.369Z',
            entitlement: {
              proGear: 2000,
              proGearSpouse: 500,
            },
            grade: 'E_7',
            has_dependents: false,
            id: 'testOrders1',
            issue_date: '2024-02-29',
            new_duty_location: {
              address: {
                city: 'Fort Irwin',
                country: 'United States',
                id: '77dca457-d0d6-4718-9ca4-a630b4614cf8',
                postalCode: '92310',
                state: 'CA',
                streetAddress1: 'n/a',
              },
              address_id: '77dca457-d0d6-4718-9ca4-a630b4614cf8',
              affiliation: 'ARMY',
              created_at: '2024-02-22T21:34:21.449Z',
              id: '12421bcb-2ded-4165-b0ac-05f76301082a',
              name: 'Fort Irwin, CA 92310',
              transportation_office: {
                address: {
                  city: 'Fort Irwin',
                  country: 'United States',
                  id: '65a97b21-cf6a-47c1-a4b6-e3f885dacba5',
                  postalCode: '92310',
                  state: 'CA',
                  streetAddress1: 'Langford Lake Rd',
                  streetAddress2: 'Bldg 105',
                },
                created_at: '2018-05-28T14:27:37.312Z',
                gbloc: 'LKNQ',
                id: 'd00e3ee8-baba-4991-8f3b-86c2e370d1be',
                name: 'PPPO Fort Irwin - USA',
                phone_lines: [],
                updated_at: '2018-05-28T14:27:37.312Z',
              },
              transportation_office_id: 'd00e3ee8-baba-4991-8f3b-86c2e370d1be',
              updated_at: '2024-02-22T21:34:21.449Z',
            },
            orders_type: 'PERMANENT_CHANGE_OF_STATION',
            originDutyLocationGbloc: 'BGAC',
            origin_duty_location: {
              address: {
                city: 'Fort Gregg-Adams',
                country: 'United States',
                id: '12270b68-01cf-4416-8b19-125d11bc8340',
                postalCode: '23801',
                state: 'VA',
                streetAddress1: 'n/a',
              },
              address_id: '12270b68-01cf-4416-8b19-125d11bc8340',
              affiliation: 'ARMY',
              created_at: '2024-02-22T21:34:26.430Z',
              id: '9cf15b8d-985b-4ca3-9f27-4ba32a263908',
              name: 'Fort Gregg-Adams, VA 23801',
              transportation_office: {
                address: {
                  city: 'Fort Gregg-Adams',
                  country: 'United States',
                  id: '10dc88f5-d76a-427f-89a0-bf85587b0570',
                  postalCode: '23801',
                  state: 'VA',
                  streetAddress1: '1401 B Ave',
                  streetAddress2: 'Bldg 3400, Room 119',
                },
                created_at: '2018-05-28T14:27:42.125Z',
                gbloc: 'BGAC',
                id: '4cc26e01-f0ea-4048-8081-1d179426a6d9',
                name: 'PPPO Fort Gregg-Adams - USA',
                phone_lines: [],
                updated_at: '2018-05-28T14:27:42.125Z',
              },
              transportation_office_id: '4cc26e01-f0ea-4048-8081-1d179426a6d9',
              updated_at: '2024-02-22T21:34:26.430Z',
            },
            report_by_date: '2024-02-29',
            service_member_id: '81aeac60-80f3-44d1-9b74-ba6d405ee2da',
            spouse_has_pro_gear: false,
            status: 'DRAFT',
            updated_at: '2024-02-23T19:30:11.369Z',
            uploaded_orders: {
              id: 'bd35c4c2-41c6-44a1-bf54-9098c68d87cc',
              service_member_id: '81aeac60-80f3-44d1-9b74-ba6d405ee2da',
              uploads: [
                {
                  bytes: 92797,
                  contentType: 'image/png',
                  createdAt: '2024-02-26T18:43:58.515Z',
                  filename: 'Screenshot 2024-02-08 at 12.57.43 PM.png',
                  id: '786237dc-c240-449d-8859-3f37583b3406',
                  status: 'PROCESSING',
                  updatedAt: '2024-02-26T18:43:58.515Z',
                  url: '/storage/user/5fe4d948-aa1c-4823-8967-b1fb40cf6679/uploads/786237dc-c240-449d-8859-3f37583b3406?contentType=image%2Fpng',
                },
              ],
            },
          },
          status: 'DRAFT',
          submittedAt: '0001-01-01T00:00:00.000Z',
          updatedAt: '0001-01-01T00:00:00.000Z',
        },
      ],
      previousMoves: [],
    },
    updateOrders: jest.fn(),
  };

  it('renders all content of Orders component', async () => {
    selectServiceMemberFromLoggedInUser.mockImplementation(() => serviceMember);
    selectOrdersForLoggedInUser.mockImplementation(() => testProps.orders);
    selectAllMoves.mockImplementation(() => testProps.serviceMemberMoves);
    renderWithProviders(<Orders {...testProps} />, {
      path: customerRoutes.ORDERS_INFO_PATH,
      params: { orderId: 'testOrderId' },
    });

    await screen.findByRole('heading', { level: 1, name: 'Tell us about your move orders' });
    expect(screen.getByTestId('main-container')).toBeInTheDocument();
    expect(screen.getByTestId('orders-form-container')).toBeInTheDocument();
    const saveBtn = await screen.findByRole('button', { name: 'Back' });
    expect(saveBtn).toBeInTheDocument();
    const cancelBtn = await screen.findByRole('button', { name: 'Next' });
    expect(cancelBtn).toBeInTheDocument();
  });

  it('renders appropriate order data on load', async () => {
    selectServiceMemberFromLoggedInUser.mockImplementation(() => serviceMember);
    selectOrdersForLoggedInUser.mockImplementation(() => testProps.orders);
    selectAllMoves.mockImplementation(() => testProps.serviceMemberMoves);
    renderWithProviders(<Orders {...testProps} />, {
      path: customerRoutes.ORDERS_INFO_PATH,
      params: { orderId: 'testOrderId' },
    });

    await screen.findByRole('heading', { level: 1, name: 'Tell us about your move orders' });
    expect(screen.getByLabelText(/Orders type/)).toHaveValue('PERMANENT_CHANGE_OF_STATION');
    expect(screen.getByLabelText(/Orders date/)).toHaveValue('08 Nov 2020');
    expect(screen.getByLabelText(/Report by date/)).toHaveValue('26 Nov 2020');
    expect(screen.getByLabelText('Yes')).not.toBeChecked();
    expect(screen.getByLabelText('No')).toBeChecked();
    expect(screen.queryByText('Yuma AFB')).toBeInTheDocument();
    expect(screen.getByLabelText(/Pay grade/)).toHaveValue('E_8');
    expect(screen.queryByText('Altus AFB')).toBeInTheDocument();
  });

  it('next button patches the orders updates state', async () => {
    selectServiceMemberFromLoggedInUser.mockImplementation(() => serviceMember);
    selectOrdersForLoggedInUser.mockImplementation(() => testProps.orders);
    selectAllMoves.mockImplementation(() => testProps.serviceMemberMoves);
    const testOrdersValues = {
      id: 'testOrdersId',
      orders_type: 'PERMANENT_CHANGE_OF_STATION',
      issue_date: '2020-11-08',
      report_by_date: '2020-11-26',
      has_dependents: false,
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
    patchOrders.mockImplementation(() => Promise.resolve(testOrdersValues));
    getOrders.mockImplementation(() => Promise.resolve());

    renderWithProviders(<Orders {...testProps} />, {
      path: customerRoutes.ORDERS_INFO_PATH,
      params: { orderId: 'testOrderId' },
    });

    const nextBtn = await screen.findByRole('button', { name: 'Next' });
    expect(nextBtn).toBeInTheDocument();

    await userEvent.click(nextBtn);

    await waitFor(() => {
      expect(patchOrders).toHaveBeenCalled();
      expect(getOrders).toHaveBeenCalledWith('testOrderId');
    });
  });

  it('shows an error if the API returns an error', async () => {
    selectServiceMemberFromLoggedInUser.mockImplementation(() => serviceMember);
    selectOrdersForLoggedInUser.mockImplementation(() => testProps.orders);
    selectAllMoves.mockImplementation(() => testProps.serviceMemberMoves);
    patchOrders.mockImplementation(() =>
      // Disable this rule because makeSwaggerRequest does not throw an error if the API call fails
      // eslint-disable-next-line prefer-promise-reject-errors
      Promise.reject({
        message: 'A server error occurred saving the orders',
        response: {
          body: {
            detail: 'A server error occurred saving the orders',
          },
        },
      }),
    );

    renderWithProviders(<Orders {...testProps} />, {
      path: customerRoutes.ORDERS_INFO_PATH,
      params: { orderId: 'testOrderId' },
    });

    await waitFor(() => {
      expect(screen.queryByText('Next')).toBeInTheDocument();
    });
    await userEvent.click(screen.queryByText('Next'));

    await waitFor(() => {
      expect(patchOrders).toHaveBeenCalled();
    });

    expect(screen.queryByText('A server error occurred saving the orders')).toBeInTheDocument();
    expect(mockNavigate).not.toHaveBeenCalled();
  });
});
