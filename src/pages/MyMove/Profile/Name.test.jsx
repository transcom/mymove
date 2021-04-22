/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import { mount } from 'enzyme';
import * as reactRedux from 'react-redux';
import { push } from 'connected-react-router';
import { render, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import ConnectedName, { Name } from './Name';

import { MockProviders } from 'testUtils';
import { patchServiceMember } from 'services/internalApi';

jest.mock('services/internalApi', () => ({
  ...jest.requireActual('services/internalApi'),
  patchServiceMember: jest.fn(),
}));

describe('Name page', () => {
  const testProps = {
    updateServiceMember: jest.fn(),
    push: jest.fn(),
    serviceMember: {
      id: 'testServiceMemberId',
      first_name: 'Leo',
      last_name: 'Spaceman',
    },
  };

  it('renders the NameForm', async () => {
    await waitFor(() => {
      const wrapper = mount(<Name {...testProps} />);
      expect(wrapper.find('NameForm').exists()).toBe(true);
    });
  });

  it('back button goes to the DoD Info step', async () => {
    const { queryByText } = render(<Name {...testProps} />);

    const backButton = queryByText('Back');

    await waitFor(() => {
      expect(backButton).toBeInTheDocument();
    });

    userEvent.click(backButton);
    expect(testProps.push).toHaveBeenCalledWith('/service-member/dod-info');
  });

  it('next button submits the form and goes to the Name step', async () => {
    const testServiceMemberValues = {
      id: 'testServiceMemberId',
      first_name: 'Leo',
      middle_name: 'Star',
      last_name: 'Spaceman',
      suffix: 'Mr.',
    };

    patchServiceMember.mockImplementation(() => Promise.resolve(testServiceMemberValues));

    // Need to provide initial values because we aren't testing the form here, and just want to submit immediately
    const { queryByText } = render(<Name {...testProps} serviceMember={testServiceMemberValues} />);

    const submitButton = queryByText('Next');
    expect(submitButton).toBeInTheDocument();
    userEvent.click(submitButton);

    await waitFor(() => {
      expect(patchServiceMember).toHaveBeenCalled();
    });

    expect(testProps.updateServiceMember).toHaveBeenCalledWith(testServiceMemberValues);
    expect(testProps.push).toHaveBeenCalledWith('/service-member/contact-info');
  });

  it('shows an error if the API returns an error', async () => {
    const testServiceMemberValues = {
      id: 'testServiceMemberId',
      first_name: 'Leo',
      middle_name: 'Star',
      last_name: 'Spaceman',
      suffix: 'Mr.',
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
    const { queryByText } = render(<Name {...testProps} serviceMember={testServiceMemberValues} />);

    const submitButton = queryByText('Next');
    expect(submitButton).toBeInTheDocument();
    userEvent.click(submitButton);

    await waitFor(() => {
      expect(patchServiceMember).toHaveBeenCalled();
    });

    expect(queryByText('A server error occurred saving the service member')).toBeInTheDocument();
    expect(testProps.updateServiceMember).not.toHaveBeenCalled();
    expect(testProps.push).not.toHaveBeenCalled();
  });

  afterEach(jest.resetAllMocks);
});

describe('requireCustomerState Name', () => {
  const useDispatchMock = jest.spyOn(reactRedux, 'useDispatch');
  const mockDispatch = jest.fn();

  beforeEach(() => {
    useDispatchMock.mockClear();
    mockDispatch.mockClear();
    useDispatchMock.mockReturnValue(mockDispatch);
  });

  const props = {
    updateServiceMember: jest.fn(),
    push: jest.fn(),
  };

  it('dispatches a redirect if the current state is earlier than the "DOD INFO COMPLETE" state', () => {
    const mockState = {
      entities: {
        user: {
          testUserId: {
            id: 'testUserId',
            email: 'testuser@example.com',
            service_member: 'testServiceMemberId',
          },
        },
        serviceMembers: {
          testServiceMemberId: {
            id: 'testServiceMemberId',
          },
        },
      },
    };

    const wrapper = mount(
      <MockProviders initialState={mockState}>
        <ConnectedName {...props} />
      </MockProviders>,
    );

    expect(wrapper.exists()).toBe(true);
    expect(mockDispatch).toHaveBeenCalledWith(push('/service-member/conus-oconus'));
  });

  it('does not redirect if the current state equals the "DOD INFO COMPLETE" state', () => {
    const mockState = {
      entities: {
        user: {
          testUserId: {
            id: 'testUserId',
            email: 'testuser@example.com',
            service_member: 'testServiceMemberId',
          },
        },
        serviceMembers: {
          testServiceMemberId: {
            id: 'testServiceMemberId',
            rank: 'test rank',
            edipi: '1234567890',
            affiliation: 'ARMY',
          },
        },
      },
    };

    const wrapper = mount(
      <MockProviders initialState={mockState}>
        <ConnectedName {...props} />
      </MockProviders>,
    );

    expect(wrapper.exists()).toBe(true);
    expect(mockDispatch).not.toHaveBeenCalled();
  });
  it('does not redirect if the current state is after the "DOD INFO COMPLETE" state and profile is not complete', () => {
    const mockState = {
      entities: {
        user: {
          testUserId: {
            id: 'testUserId',
            email: 'testuser@example.com',
            service_member: 'testServiceMemberId',
          },
        },
        serviceMembers: {
          testServiceMemberId: {
            id: 'testServiceMemberId',
            rank: 'test rank',
            edipi: '1234567890',
            affiliation: 'ARMY',
            first_name: 'Tester',
            last_name: 'Testperson',
            telephone: '1234567890',
            personal_email: 'test@example.com',
            email_is_preferred: true,
            current_station: {
              id: 'testDutyStationId',
            },
          },
        },
      },
    };

    const wrapper = mount(
      <MockProviders initialState={mockState}>
        <ConnectedName {...props} />
      </MockProviders>,
    );

    expect(wrapper.exists()).toBe(true);
    expect(mockDispatch).not.toHaveBeenCalled();
  });

  it('does redirect if the profile is complete', () => {
    const mockState = {
      entities: {
        user: {
          testUserId: {
            id: 'testUserId',
            email: 'testuser@example.com',
            service_member: 'testServiceMemberId',
          },
        },
        serviceMembers: {
          testServiceMemberId: {
            id: 'testServiceMemberId',
            rank: 'test rank',
            edipi: '1234567890',
            affiliation: 'ARMY',
            first_name: 'Tester',
            last_name: 'Testperson',
            telephone: '1234567890',
            personal_email: 'test@example.com',
            email_is_preferred: true,
            current_station: {
              id: 'testDutyStationId',
            },
            residential_address: {
              street: '123 Main St',
            },
            backup_mailing_address: {
              street: '456 Main St',
            },
            backup_contacts: [
              {
                id: 'testBackupContact',
              },
            ],
          },
        },
      },
    };

    const wrapper = mount(
      <MockProviders initialState={mockState}>
        <ConnectedName {...props} />
      </MockProviders>,
    );

    expect(wrapper.exists()).toBe(true);
    expect(mockDispatch).toHaveBeenCalledWith(push('/'));
  });
});
