import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { useLocation } from 'react-router-dom';

import { EditServiceInfo } from './EditServiceInfo';

import { patchServiceMember, getAllMoves } from 'services/internalApi';
import { MockProviders } from 'testUtils';

const mockNavigate = jest.fn();
jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom'),
  useNavigate: () => mockNavigate,
  useLocation: jest.fn(),
}));

jest.mock('services/internalApi', () => ({
  ...jest.requireActual('services/internalApi'),
  patchServiceMember: jest.fn(),
  getAllMoves: jest.fn().mockImplementation(() => Promise.resolve()),
}));

beforeEach(() => {
  jest.resetAllMocks();
});

const serviceMemberMoves = {
  currentMove: [
    {
      lockExpiresAt: '2099-04-07T17:21:30.450Z',
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
};

describe('EditServiceInfo page', () => {
  const testProps = {
    updateServiceMember: jest.fn(),
    setFlashMessage: jest.fn(),
    serviceMember: {
      id: 'testServiceMemberId',
    },
    currentOrders: {},
    moveIsInDraft: true,
  };

  it('renders the EditServiceInfo form', async () => {
    useLocation.mockReturnValue({ state: { moveId: 'testMoveId' } });
    render(
      <MockProviders>
        <EditServiceInfo {...testProps} />
      </MockProviders>,
    );

    expect(await screen.findByRole('heading', { name: 'Edit service info', level: 1 })).toBeInTheDocument();
  });

  it('the cancel button goes back to the profile page', async () => {
    useLocation.mockReturnValue({ state: { moveId: 'testMoveId' } });
    render(
      <MockProviders>
        <EditServiceInfo {...testProps} />
      </MockProviders>,
    );
    const cancelButton = await screen.findByText('Cancel');
    await waitFor(() => {
      expect(cancelButton).toBeInTheDocument();
    });

    await userEvent.click(cancelButton);
    expect(mockNavigate).toHaveBeenCalled();
  });

  it('save button submits the form and goes to the profile page', async () => {
    useLocation.mockReturnValue({ state: { moveId: 'testMoveId' } });
    const testServiceMemberValues = {
      id: 'testServiceMemberId',
      first_name: 'Leo',
      last_name: 'Spaceman',
      affiliation: 'NAVY',
      edipi: '1234567890',
      current_location: {
        address: {
          city: 'Test City',
          id: '25be4d12-fe93-47f1-bbec-1db386dfa67f',
          postalCode: '12345',
          state: 'NY',
          streetAddress1: '123 Main St',
        },
        address_id: '25be4d12-fe93-47f1-bbec-1db386dfa67f',
        affiliation: 'AIR_FORCE',
        created_at: '2021-02-11T16:48:04.117Z',
        id: 'a8d6b33c-8370-4e92-8df2-356b8c9d0c1a',
        name: 'Luke AFB',
        updated_at: '2021-02-11T16:48:04.117Z',
      },
    };

    patchServiceMember.mockImplementation(() => Promise.resolve(testServiceMemberValues));

    // Need to provide initial values because we aren't testing the form here, and just want to submit immediately
    render(
      <MockProviders>
        <EditServiceInfo
          {...testProps}
          serviceMember={testServiceMemberValues}
          currentOrders={{
            has_dependents: false,
          }}
        />
      </MockProviders>,
    );

    const submitButton = await screen.findByText('Save');
    expect(submitButton).toBeInTheDocument();
    await userEvent.click(submitButton);

    await waitFor(() => {
      expect(patchServiceMember).toHaveBeenCalled();
    });

    expect(testProps.updateServiceMember).toHaveBeenCalledWith(testServiceMemberValues);

    expect(mockNavigate).toHaveBeenCalled();
  });

  it('converts the "Submit" button into the "Return to Home" button when the move has been locked by an office user', async () => {
    getAllMoves.mockResolvedValue(serviceMemberMoves);
    useLocation.mockReturnValue({ state: { moveId: 'testMoveId' } });
    const testServiceMemberValues = {
      id: 'testServiceMemberId',
      first_name: 'Leo',
      last_name: 'Spaceman',
      affiliation: 'NAVY',
      edipi: '1234567890',
      current_location: {
        address: {
          city: 'Test City',
          id: '25be4d12-fe93-47f1-bbec-1db386dfa67f',
          postalCode: '12345',
          state: 'NY',
          streetAddress1: '123 Main St',
        },
        address_id: '25be4d12-fe93-47f1-bbec-1db386dfa67f',
        affiliation: 'AIR_FORCE',
        created_at: '2021-02-11T16:48:04.117Z',
        id: 'a8d6b33c-8370-4e92-8df2-356b8c9d0c1a',
        name: 'Luke AFB',
        updated_at: '2021-02-11T16:48:04.117Z',
      },
    };

    // Need to provide initial values because we aren't testing the form here, and just want to submit immediately
    render(
      <MockProviders>
        <EditServiceInfo
          {...testProps}
          serviceMember={testServiceMemberValues}
          currentOrders={{
            has_dependents: false,
          }}
        />
      </MockProviders>,
    );

    expect(await screen.findByText('Return home')).toBeInTheDocument();
  });

  it('shows an error if the API returns an error', async () => {
    useLocation.mockReturnValue({ state: { moveId: 'testMoveId' } });
    const testServiceMemberValues = {
      id: 'testServiceMemberId',
      first_name: 'Leo',
      last_name: 'Spaceman',
      affiliation: 'NAVY',
      edipi: '1234567890',
      current_location: {
        address: {
          city: 'Test City',
          id: '25be4d12-fe93-47f1-bbec-1db386dfa67f',
          postalCode: '12345',
          state: 'NY',
          streetAddress1: '123 Main St',
        },
        address_id: '25be4d12-fe93-47f1-bbec-1db386dfa67f',
        affiliation: 'AIR_FORCE',
        created_at: '2021-02-11T16:48:04.117Z',
        id: 'a8d6b33c-8370-4e92-8df2-356b8c9d0c1a',
        name: 'Luke AFB',
        updated_at: '2021-02-11T16:48:04.117Z',
      },
    };

    patchServiceMember.mockImplementation(() =>
      // Disable this rule because makeSwaggerRequest does not throw an error if the API call fails
      // eslint-disable-next-line prefer-promise-reject-errors
      Promise.reject({
        message: 'A server error occurred saving the service member',
        response: {
          body: {
            detail: 'A server error occurred saving the service member',
          },
        },
      }),
    );

    // Need to provide complete & valid initial values because we aren't testing the form here, and just want to submit immediately
    render(
      <MockProviders>
        <EditServiceInfo {...testProps} serviceMember={testServiceMemberValues} currentOrders={{}} />
      </MockProviders>,
    );

    const submitButton = await screen.findByText('Save');
    expect(submitButton).toBeInTheDocument();
    await userEvent.click(submitButton);

    await waitFor(() => {
      expect(patchServiceMember).toHaveBeenCalled();
    });

    expect(await screen.findByText('A server error occurred saving the service member')).toBeInTheDocument();
    expect(testProps.updateServiceMember).not.toHaveBeenCalled();
    expect(testProps.setFlashMessage).not.toHaveBeenCalled();
    expect(mockNavigate).not.toHaveBeenCalled();
  });

  describe('if the current move has been submitted', () => {
    it('redirects to the home page', async () => {
      useLocation.mockReturnValue({ state: { moveId: 'testMoveId' } });
      render(
        <MockProviders>
          <EditServiceInfo {...testProps} moveIsInDraft={false} />
        </MockProviders>,
      );

      await waitFor(() => {
        expect(mockNavigate).toHaveBeenCalledWith('/');
      });
    });
  });

  afterEach(jest.resetAllMocks);
});
