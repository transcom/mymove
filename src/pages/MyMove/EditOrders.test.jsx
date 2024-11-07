import React from 'react';
import { screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { act } from 'react-dom/test-utils';

import EditOrders from './EditOrders';

import { getOrders, patchOrders, showCounselingOffices } from 'services/internalApi';
import { renderWithProviders } from 'testUtils';
import { customerRoutes } from 'constants/routes';
import {
  selectAllMoves,
  selectOrdersForLoggedInUser,
  selectServiceMemberFromLoggedInUser,
} from 'store/entities/selectors';
import { ORDERS_TYPE } from 'constants/orders';
import { isBooleanFlagEnabled } from 'utils/featureFlags';

const mockNavigate = jest.fn();
jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom'),
  useNavigate: () => mockNavigate,
}));

jest.mock('store/entities/selectors', () => ({
  ...jest.requireActual('store/entities/selectors'),
  selectServiceMemberFromLoggedInUser: jest.fn(),
  selectOrdersForLoggedInUser: jest.fn(),
  selectAllMoves: jest.fn(),
}));

jest.mock('services/internalApi', () => ({
  ...jest.requireActual('services/internalApi'),
  getOrders: jest.fn().mockImplementation(() => Promise.resolve()),
  patchOrders: jest.fn().mockImplementation(() => Promise.resolve()),
  getAllMoves: jest.fn().mockImplementation(() => Promise.resolve()),
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

jest.mock('utils/featureFlags', () => ({
  ...jest.requireActual('utils/featureFlags'),
  isBooleanFlagEnabled: jest.fn().mockImplementation(() => Promise.resolve(false)),
}));

describe('EditOrders Page', () => {
  const testProps = {
    serviceMember: {
      id: 'id123',
      current_location: {
        address: {
          city: 'Fort Bragg',
          country: 'United States',
          id: 'f1ee4cea-6b23-4971-9947-efb51294ed32',
          postalCode: '29310',
          state: 'NC',
          streetAddress1: '',
        },
        address_id: 'f1ee4cea-6b23-4971-9947-efb51294ed32',
        affiliation: 'ARMY',
        created_at: '2020-10-19T17:01:16.114Z',
        id: 'dca78766-e76b-4c6d-ba82-81b50ca824b9"',
        name: 'Fort Bragg',
        updated_at: '2020-10-19T17:01:16.114Z',
      },
    },
    serviceMemberId: 'id123',
    orders: [
      {
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
        moves: ['testMoveId'],
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
            phone_lines: ['760-380-3823', '470-3823'],
            updated_at: '2018-05-28T14:27:37.312Z',
          },
          transportation_office_id: 'd00e3ee8-baba-4991-8f3b-86c2e370d1be',
          updated_at: '2024-02-22T21:34:21.449Z',
        },
        orders_type: ORDERS_TYPE.PERMANENT_CHANGE_OF_STATION,
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
            address: null,
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
    ],
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
            orders_type: ORDERS_TYPE.PERMANENT_CHANGE_OF_STATION,
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
    setFlashMessage: jest.fn(),
    updateOrders: jest.fn(),
    updateAllMoves: jest.fn(),
    context: { flags: { allOrdersTypes: true } },
  };
  selectServiceMemberFromLoggedInUser.mockImplementation(() => testProps.serviceMember);
  selectOrdersForLoggedInUser.mockImplementation(() => testProps.orders);
  getOrders.mockResolvedValue(() => testProps.orders[0]);
  selectAllMoves.mockImplementation(() => testProps.serviceMemberMoves);

  it('renders the edit orders form', async () => {
    renderWithProviders(<EditOrders {...testProps} />, {
      path: customerRoutes.ORDERS_EDIT_PATH,
      params: { moveId: 'testMoveId', orderId: 'testOrders1' },
    });

    const h1 = await screen.findByRole('heading', { name: 'Orders', level: 1 });
    expect(h1).toBeInTheDocument();

    const editOrdersHeader = await screen.findByRole('heading', { name: 'Edit Orders:', level: 2 });
    expect(editOrdersHeader).toBeInTheDocument();
  });

  it('delete button visible for orders when move is in draft state', async () => {
    renderWithProviders(<EditOrders {...testProps} />, {
      path: customerRoutes.ORDERS_EDIT_PATH,
      params: { moveId: 'testMoveId', orderId: 'testOrders1' },
    });
    const deleteBtn = await screen.findByRole('button', { name: 'Delete' });
    expect(deleteBtn).toBeInTheDocument();
  });

  it('no option to delete uploaded orders when move is submitted', async () => {
    testProps.serviceMemberMoves.currentMove[0].status = 'SUBMITTED';
    renderWithProviders(<EditOrders {...testProps} />, {
      path: customerRoutes.ORDERS_EDIT_PATH,
      params: { moveId: 'testMoveId', orderId: 'testOrders1' },
    });
    expect(screen.queryByRole('button', { name: 'Delete' })).toBeNull();
  });

  it('goes back to the previous page when the cancel button is clicked', async () => {
    showCounselingOffices.mockImplementation(() => Promise.resolve({}));
    renderWithProviders(<EditOrders {...testProps} />, {
      path: customerRoutes.ORDERS_EDIT_PATH,
      params: { moveId: 'testMoveId', orderId: 'testOrders1' },
    });
    const deleteBtn = await screen.findByRole('button', { name: 'Cancel' });
    expect(deleteBtn).toBeInTheDocument();

    await userEvent.click(deleteBtn);

    await waitFor(() => {
      expect(mockNavigate).toHaveBeenCalledWith(-1);
    });
  });

  it('shows an error if the API returns an error', async () => {
    renderWithProviders(<EditOrders {...testProps} />, {
      path: customerRoutes.ORDERS_EDIT_PATH,
      params: { moveId: 'testMoveId', orderId: 'testOrders1' },
    });
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

    const submitButton = await screen.findByRole('button', { name: 'Save' });
    expect(submitButton).not.toBeDisabled();

    await userEvent.click(submitButton);

    await waitFor(() => {
      expect(patchOrders).toHaveBeenCalledTimes(1);
    });

    expect(screen.queryByText('A server error occurred saving the orders')).toBeInTheDocument();
    expect(mockNavigate).not.toHaveBeenCalled();
  });

  it('next button patches the orders and goes to the previous page', async () => {
    renderWithProviders(<EditOrders {...testProps} />, {
      path: customerRoutes.ORDERS_EDIT_PATH,
      params: { moveId: 'testMoveId', orderId: 'testOrders1' },
    });
    patchOrders.mockImplementation(() => Promise.resolve(testProps.currentOrders));

    const submitButton = await screen.findByRole('button', { name: 'Save' });
    expect(submitButton).not.toBeDisabled();

    await userEvent.click(submitButton);

    await waitFor(() => {
      expect(patchOrders).toHaveBeenCalledTimes(1);
    });

    expect(mockNavigate).toHaveBeenCalledTimes(1);
    expect(mockNavigate).toHaveBeenCalledWith(-1);
  });

  it('submits OCONUS fields correctly on form submit', async () => {
    isBooleanFlagEnabled.mockImplementation(() => Promise.resolve(true));
    testProps.orders[0].origin_duty_location.address = {
      ...testProps.orders[0].origin_duty_location.address,
      isOconus: true,
    };

    renderWithProviders(<EditOrders {...testProps} />, {
      path: customerRoutes.ORDERS_EDIT_PATH,
      params: { moveId: 'testMoveId', orderId: 'testOrders1' },
    });
    await waitFor(() => {
      expect(screen.getByRole('form')).toHaveFormValues({
        new_duty_location: 'Fort Irwin, CA 92310',
        origin_duty_location: 'Fort Gregg-Adams, VA 23801',
      });
    });
    await userEvent.click(screen.getByTestId('hasDependentsYes'));
    await userEvent.click(screen.getByTestId('isAnAccompaniedTourYes'));
    await userEvent.type(screen.getByTestId('dependentsUnderTwelve'), '1');
    await userEvent.type(screen.getByTestId('dependentsTwelveAndOver'), '2');

    const submitButton = await screen.findByRole('button', { name: 'Save' });
    expect(submitButton).not.toBeDisabled();

    await act(async () => {
      userEvent.click(submitButton);
    });

    await waitFor(() => {
      expect(patchOrders).toHaveBeenCalledWith(
        expect.objectContaining({
          accompanied_tour: true,
          dependents_under_twelve: 1,
          dependents_twelve_and_over: 2,
        }),
      );
    });
  });

  afterEach(jest.clearAllMocks);
});
