import React from 'react';
import { render, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import ConnectedDodInfo, { DodInfo } from './DodInfo';

import { MockProviders } from 'testUtils';
import { patchServiceMember } from 'services/internalApi';

const mockNavigate = jest.fn();
jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom'),
  useNavigate: () => mockNavigate,
}));

jest.mock('services/internalApi', () => ({
  ...jest.requireActual('services/internalApi'),
  patchServiceMember: jest.fn(),
}));

beforeEach(() => {
  jest.resetAllMocks();
});

describe('DodInfo page', () => {
  const testProps = {
    updateServiceMember: jest.fn(),
    serviceMember: {
      id: 'testServiceMemberId',
    },
  };

  it('renders the DodInfoForm', async () => {
    const { queryByRole } = render(<DodInfo {...testProps} />);

    await waitFor(() => {
      expect(queryByRole('heading', { name: 'Create your profile', level: 1 })).toBeInTheDocument();
    });
  });

  it('back button goes to the CONUS/OCONUS step', async () => {
    const { queryByText } = render(<DodInfo {...testProps} />);

    const backButton = queryByText('Back');
    await waitFor(() => {
      expect(backButton).toBeInTheDocument();
    });

    await userEvent.click(backButton);
    expect(mockNavigate).toHaveBeenCalledWith('/service-member/conus-oconus');
  });

  it('next button submits the form and goes to the Name step', async () => {
    const testServiceMemberValues = {
      id: 'testServiceMemberId',
      affiliation: 'ARMY',
      edipi: '9999999999',
    };

    patchServiceMember.mockImplementation(() => Promise.resolve(testServiceMemberValues));

    // Need to provide initial values because we aren't testing the form here, and just want to submit immediately
    const { queryByText } = render(<DodInfo {...testProps} serviceMember={testServiceMemberValues} />);

    const submitButton = queryByText('Next');
    expect(submitButton).toBeInTheDocument();
    await userEvent.click(submitButton);

    await waitFor(() => {
      expect(patchServiceMember).toHaveBeenCalled();
    });

    expect(testProps.updateServiceMember).toHaveBeenCalledWith(testServiceMemberValues);
    expect(mockNavigate).toHaveBeenCalledWith('/service-member/name');
  });

  it('shows an error if the API returns an error', async () => {
    const testServiceMemberValues = {
      id: 'testServiceMemberId',
      affiliation: 'ARMY',
      edipi: '9999999999',
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
    const { queryByText } = render(<DodInfo {...testProps} serviceMember={testServiceMemberValues} />);

    const submitButton = queryByText('Next');
    expect(submitButton).toBeInTheDocument();
    await userEvent.click(submitButton);

    await waitFor(() => {
      expect(patchServiceMember).toHaveBeenCalled();
    });

    expect(queryByText('A server error occurred saving the service member')).toBeInTheDocument();
    expect(testProps.updateServiceMember).not.toHaveBeenCalled();
    expect(mockNavigate).not.toHaveBeenCalled();
  });

  afterEach(jest.resetAllMocks);
});

describe('requireCustomerState DodInfo', () => {
  const props = {
    updateServiceMember: jest.fn(),
  };

  it('does not redirect if the current state equals the "EMPTY PROFILE" state', async () => {
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

    render(
      <MockProviders initialState={mockState}>
        <ConnectedDodInfo {...props} />
      </MockProviders>,
    );

    await waitFor(() => {
      expect(mockNavigate).not.toHaveBeenCalled();
    });
  });

  it('does not redirect if the current state is after the "EMPTY PROFILE" state and profile is not complete', async () => {
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
            edipi: '1234567890',
            affiliation: 'ARMY',
            first_name: 'Tester',
            last_name: 'Testperson',
            telephone: '1234567890',
            personal_email: 'test@example.com',
            email_is_preferred: true,
          },
        },
      },
    };

    render(
      <MockProviders initialState={mockState}>
        <ConnectedDodInfo {...props} />
      </MockProviders>,
    );

    await waitFor(() => {
      expect(mockNavigate).not.toHaveBeenCalled();
    });
  });

  it('does redirect if the profile is complete', async () => {
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
            edipi: '1234567890',
            affiliation: 'ARMY',
            first_name: 'Tester',
            last_name: 'Testperson',
            telephone: '1234567890',
            personal_email: 'test@example.com',
            email_is_preferred: true,
            current_location: {
              id: 'testDutyLocationId',
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

    render(
      <MockProviders initialState={mockState}>
        <ConnectedDodInfo {...props} />
      </MockProviders>,
    );

    await waitFor(() => {
      expect(mockNavigate).toHaveBeenCalledWith('/');
    });
  });
});
