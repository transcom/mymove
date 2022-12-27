import React from 'react';
import * as reactRedux from 'react-redux';
import { push } from 'connected-react-router';
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import ConnectedContactInfo, { ContactInfo } from './ContactInfo';

import { MockProviders } from 'testUtils';
import { patchServiceMember } from 'services/internalApi';

jest.mock('services/internalApi', () => ({
  ...jest.requireActual('services/internalApi'),
  patchServiceMember: jest.fn(),
}));

describe('ContactInfo page', () => {
  const testProps = {
    updateServiceMember: jest.fn(),
    push: jest.fn(),
    serviceMember: {
      id: 'testServiceMemberId',
    },
    userEmail: 'sm@example.com',
  };

  const testServiceMemberValues = {
    id: 'testServiceMemberId',
    telephone: '555-555-5555',
    personal_email: 'test@example.com',
    email_is_preferred: true,
  };

  it('renders the ContactInfoForm', async () => {
    render(<ContactInfo {...testProps} />);
    expect(await screen.findByRole('heading', { name: 'Your contact info', level: 1 })).toBeInTheDocument();
  });

  it('back button goes to the NAME step', async () => {
    // Need to provide initial values because we aren't testing the form here, and just want to submit immediately
    render(<ContactInfo {...testProps} serviceMember={testServiceMemberValues} />);

    const backButton = screen.getByRole('button', { name: 'Back' });
    expect(backButton).toBeInTheDocument();
    await userEvent.click(backButton);

    await waitFor(async () => {
      expect(testProps.push).toHaveBeenCalledWith('/service-member/name');
    });
  });

  it('next button submits the form and goes to the Name step', async () => {
    patchServiceMember.mockImplementation(() => Promise.resolve(testServiceMemberValues));

    // Need to provide initial values because we aren't testing the form here, and just want to submit immediately
    render(<ContactInfo {...testProps} serviceMember={testServiceMemberValues} />);

    const submitButton = screen.getByRole('button', { name: 'Next' });
    expect(submitButton).toBeInTheDocument();
    await userEvent.click(submitButton);

    await waitFor(() => {
      expect(patchServiceMember).toHaveBeenCalled();
    });

    expect(testProps.updateServiceMember).toHaveBeenCalledWith(testServiceMemberValues);
    expect(testProps.push).toHaveBeenCalledWith('/service-member/current-duty');
  });

  it('shows an error if the API returns an error', async () => {
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
    render(<ContactInfo {...testProps} serviceMember={testServiceMemberValues} />);

    const submitButton = screen.getByRole('button', { name: 'Next' });
    expect(submitButton).toBeInTheDocument();
    await userEvent.click(submitButton);

    await waitFor(() => {
      expect(patchServiceMember).toHaveBeenCalled();
    });

    expect(screen.queryByText('A server error occurred saving the service member')).toBeInTheDocument();
    expect(testProps.updateServiceMember).not.toHaveBeenCalled();
    expect(testProps.push).not.toHaveBeenCalled();
  });

  afterEach(jest.resetAllMocks);
});

describe('requireCustomerState ContactInfo', () => {
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
    serviceMember: {
      id: 'testServiceMemberId',
    },
    userEmail: 'sm@example.com',
  };

  it('dispatches a redirect if the current state is earlier than the "NAME COMPLETE" state', async () => {
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

    render(
      <MockProviders initialState={mockState}>
        <ConnectedContactInfo {...props} />
      </MockProviders>,
    );

    const h1 = screen.getByRole('heading', { name: 'Your contact info', level: 1 });
    expect(h1).toBeInTheDocument();

    await waitFor(async () => {
      expect(mockDispatch).toHaveBeenCalledWith(push('/service-member/name'));
    });
  });

  it('does not redirect if the current state equals the "NAME COMPLETE" state', async () => {
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
          },
        },
      },
    };

    render(
      <MockProviders initialState={mockState}>
        <ConnectedContactInfo {...props} />
      </MockProviders>,
    );
    const h1 = screen.getByRole('heading', { name: 'Your contact info', level: 1 });
    expect(h1).toBeInTheDocument();

    await waitFor(async () => {
      expect(mockDispatch).not.toHaveBeenCalled();
    });
  });
  it('does not redirect if the current state is after the "NAME COMPLETE" state and profile is not complete', async () => {
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
            current_location: {
              id: 'testDutyLocationId',
            },
            residential_address: {
              street: '123 Main St',
            },
            backup_mailing_address: {
              street: '456 Main St',
            },
          },
        },
      },
    };

    render(
      <MockProviders initialState={mockState}>
        <ConnectedContactInfo {...props} />
      </MockProviders>,
    );

    const h1 = screen.getByRole('heading', { name: 'Your contact info', level: 1 });
    expect(h1).toBeInTheDocument();

    await waitFor(async () => {
      expect(mockDispatch).not.toHaveBeenCalled();
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
            rank: 'test rank',
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
        <ConnectedContactInfo {...props} />
      </MockProviders>,
    );

    const h1 = screen.getByRole('heading', { name: 'Your contact info', level: 1 });
    expect(h1).toBeInTheDocument();

    await waitFor(async () => {
      expect(mockDispatch).toHaveBeenCalledWith(push('/'));
    });
  });
});
