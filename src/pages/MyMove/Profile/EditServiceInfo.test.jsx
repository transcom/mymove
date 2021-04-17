import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import { EditServiceInfo } from './EditServiceInfo';

import { patchServiceMember } from 'services/internalApi';

const mockPush = jest.fn();

jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom'),
  useLocation: () => ({
    pathname: 'localhost:3000/',
  }),
  useHistory: () => ({
    push: mockPush,
  }),
}));

jest.mock('services/internalApi', () => ({
  ...jest.requireActual('services/internalApi'),
  patchServiceMember: jest.fn(),
}));

describe('EditServiceInfo page', () => {
  const testProps = {
    updateServiceMember: jest.fn(),
    setFlashMessage: jest.fn(),
    serviceMember: {
      id: 'testServiceMemberId',
    },
    currentOrders: {},
    entitlement: {},
    moveIsInDraft: true,
  };

  it('renders the EditServiceInfo form', async () => {
    render(<EditServiceInfo {...testProps} />);

    expect(await screen.findByRole('heading', { name: 'Edit service info', level: 1 })).toBeInTheDocument();
  });

  it('the cancel button goes back to the profile page', async () => {
    render(<EditServiceInfo {...testProps} />);

    const cancelButton = await screen.findByText('Cancel');
    await waitFor(() => {
      expect(cancelButton).toBeInTheDocument();
    });

    userEvent.click(cancelButton);
    expect(mockPush).toHaveBeenCalledWith('/service-member/profile');
  });

  it('save button submits the form and goes to the profile page', async () => {
    const testServiceMemberValues = {
      id: 'testServiceMemberId',
      first_name: 'Leo',
      last_name: 'Spaceman',
      affiliation: 'NAVY',
      edipi: '1234567890',
      rank: 'E_5',
      current_station: {
        address: {
          city: 'Test City',
          id: '25be4d12-fe93-47f1-bbec-1db386dfa67f',
          postal_code: '12345',
          state: 'NY',
          street_address_1: '123 Main St',
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
      <EditServiceInfo
        {...testProps}
        serviceMember={testServiceMemberValues}
        currentOrders={{
          grade: testServiceMemberValues.rank,
          origin_duty_station: testServiceMemberValues.current_station,
        }}
      />,
    );

    const submitButton = await screen.findByText('Save');
    expect(submitButton).toBeInTheDocument();
    userEvent.click(submitButton);

    await waitFor(() => {
      expect(patchServiceMember).toHaveBeenCalled();
    });

    expect(testProps.updateServiceMember).toHaveBeenCalledWith(testServiceMemberValues);
    expect(testProps.setFlashMessage).toHaveBeenCalledWith(
      'EDIT_SERVICE_INFO_SUCCESS',
      'success',
      '',
      'Your changes have been saved.',
    );

    expect(mockPush).toHaveBeenCalledWith('/service-member/profile');
  });

  it('displays a flash message about entitlement when the rank changes', async () => {
    const testServiceMemberValues = {
      id: 'testServiceMemberId',
      first_name: 'Leo',
      last_name: 'Spaceman',
      affiliation: 'NAVY',
      edipi: '1234567890',
      rank: 'E_5',
      current_station: {
        address: {
          city: 'Test City',
          id: '25be4d12-fe93-47f1-bbec-1db386dfa67f',
          postal_code: '12345',
          state: 'NY',
          street_address_1: '123 Main St',
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
      <EditServiceInfo
        {...testProps}
        serviceMember={testServiceMemberValues}
        currentOrders={{
          grade: testServiceMemberValues.rank,
          origin_duty_station: testServiceMemberValues.current_station,
        }}
        entitlement={{ sum: 15000 }}
      />,
    );

    const rankInput = await screen.findByLabelText('Rank');
    userEvent.selectOptions(rankInput, ['E_2']);

    const submitButton = await screen.findByText('Save');
    expect(submitButton).toBeInTheDocument();
    userEvent.click(submitButton);

    await waitFor(() => {
      expect(patchServiceMember).toHaveBeenCalled();
    });

    expect(testProps.updateServiceMember).toHaveBeenCalledWith(testServiceMemberValues);
    expect(testProps.setFlashMessage).toHaveBeenCalledWith(
      'EDIT_SERVICE_INFO_SUCCESS',
      'info',
      `Your weight entitlement is now 15,000 lbs.`,
      'Your changes have been saved. Note that the entitlement has also changed.',
    );

    expect(mockPush).toHaveBeenCalledWith('/service-member/profile');
  });

  it('shows an error if the API returns an error', async () => {
    const testServiceMemberValues = {
      id: 'testServiceMemberId',
      first_name: 'Leo',
      last_name: 'Spaceman',
      affiliation: 'NAVY',
      edipi: '1234567890',
      rank: 'E_5',
      current_station: {
        address: {
          city: 'Test City',
          id: '25be4d12-fe93-47f1-bbec-1db386dfa67f',
          postal_code: '12345',
          state: 'NY',
          street_address_1: '123 Main St',
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
      <EditServiceInfo
        {...testProps}
        serviceMember={testServiceMemberValues}
        currentOrders={{
          grade: testServiceMemberValues.rank,
          origin_duty_station: testServiceMemberValues.current_station,
        }}
      />,
    );

    const submitButton = await screen.findByText('Save');
    expect(submitButton).toBeInTheDocument();
    userEvent.click(submitButton);

    await waitFor(() => {
      expect(patchServiceMember).toHaveBeenCalled();
    });

    expect(await screen.findByText('A server error occurred saving the service member')).toBeInTheDocument();
    expect(testProps.updateServiceMember).not.toHaveBeenCalled();
    expect(testProps.setFlashMessage).not.toHaveBeenCalled();
    expect(mockPush).not.toHaveBeenCalled();
  });

  describe('if the current move has been submitted', () => {
    it('redirects to the home page', async () => {
      render(<EditServiceInfo {...testProps} moveIsInDraft={false} />);

      await waitFor(() => {
        expect(mockPush).toHaveBeenCalledWith('/');
      });
    });
  });

  afterEach(jest.resetAllMocks);
});
