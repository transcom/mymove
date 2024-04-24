/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import ConnectedName, { Name } from './Name';

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

jest.mock('utils/featureFlags', () => ({
  ...jest.requireActual('utils/featureFlags'),
  isBooleanFlagEnabled: jest.fn().mockImplementation(() => Promise.resolve(false)),
}));

beforeEach(() => {
  jest.resetAllMocks();
});

describe('Name page', () => {
  const testProps = {
    updateServiceMember: jest.fn(),
    serviceMember: {
      id: 'testServiceMemberId',
      first_name: 'Leo',
      last_name: 'Spaceman',
    },
  };

  it('renders the NameForm', async () => {
    render(<Name {...testProps} />);

    expect(await screen.findByRole('heading', { name: 'Name', level: 1 })).toBeInTheDocument();
  });

  it('back button goes to the DoD Info step', async () => {
    render(<Name {...testProps} />);

    const backButton = await screen.findByRole('button', { name: 'Back' });

    expect(backButton).toBeInTheDocument();

    await userEvent.click(backButton);
    expect(mockNavigate).toHaveBeenCalledWith('/service-member/dod-info');
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
    render(<Name {...testProps} serviceMember={testServiceMemberValues} />);

    const submitButton = await screen.findByRole('button', { name: 'Next' });

    expect(submitButton).toBeInTheDocument();
    await userEvent.click(submitButton);

    await waitFor(() => {
      expect(patchServiceMember).toHaveBeenCalled();
    });

    expect(testProps.updateServiceMember).toHaveBeenCalledWith(testServiceMemberValues);
    expect(mockNavigate).toHaveBeenCalledWith('/service-member/contact-info');
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
    render(<Name {...testProps} serviceMember={testServiceMemberValues} />);

    const submitButton = await screen.findByRole('button', { name: 'Next' });

    expect(submitButton).toBeInTheDocument();
    await userEvent.click(submitButton);

    await waitFor(() => {
      expect(patchServiceMember).toHaveBeenCalled();
    });

    expect(screen.getByText('A server error occurred saving the service member')).toBeInTheDocument();
    expect(testProps.updateServiceMember).not.toHaveBeenCalled();
    expect(mockNavigate).not.toHaveBeenCalled();
  });

  afterEach(jest.resetAllMocks);
});

describe('requireCustomerState Name', () => {
  const props = {
    updateServiceMember: jest.fn(),
  };

  it('dispatches a redirect if the current state is earlier than the "DOD INFO COMPLETE" state', async () => {
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
        <ConnectedName {...props} />
      </MockProviders>,
    );

    expect(await screen.findByRole('heading', { name: 'Name', level: 1 })).toBeInTheDocument();

    expect(mockNavigate).toHaveBeenCalledWith('/service-member/conus-oconus');
  });

  it('does not redirect if the current state equals the "DOD INFO COMPLETE" state', async () => {
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
          },
        },
      },
    };

    render(
      <MockProviders initialState={mockState}>
        <ConnectedName {...props} />
      </MockProviders>,
    );

    expect(await screen.findByRole('heading', { name: 'Name', level: 1 })).toBeInTheDocument();

    expect(mockNavigate).not.toHaveBeenCalled();
  });

  it('does not redirect if the current state is after the "DOD INFO COMPLETE" state and profile is not complete', async () => {
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
          },
        },
      },
    };

    render(
      <MockProviders initialState={mockState}>
        <ConnectedName {...props} />
      </MockProviders>,
    );

    expect(await screen.findByRole('heading', { name: 'Name', level: 1 })).toBeInTheDocument();

    expect(mockNavigate).not.toHaveBeenCalled();
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
        <ConnectedName {...props} />
      </MockProviders>,
    );

    expect(await screen.findByRole('heading', { name: 'Name', level: 1 })).toBeInTheDocument();

    expect(mockNavigate).toHaveBeenCalledWith('/');
  });
});
