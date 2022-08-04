/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import * as reactRedux from 'react-redux';
import { render, waitFor, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { push } from 'connected-react-router';

import ConnectedDutyLocation, { DutyLocation } from './DutyLocation';

import { MockProviders } from 'testUtils';
import { patchServiceMember } from 'services/internalApi';

jest.mock('services/internalApi', () => ({
  ...jest.requireActual('services/internalApi'),
  patchServiceMember: jest.fn(),
}));

describe('Duty Location page', () => {
  const testProps = {
    updateServiceMember: jest.fn(),
    push: jest.fn(),
    serviceMember: {
      id: 'testServiceMemberId',
      first_name: 'Leo',
      last_name: 'Spaceman',
    },
  };

  it('renders the CurrentDutyLocationForm', async () => {
    render(<DutyLocation {...testProps} />);
    expect(await screen.findByRole('heading', { level: 1 })).toHaveTextContent('Current duty location');
  });

  it('back button goes to the Contact Info step', async () => {
    render(<DutyLocation {...testProps} />);

    const backButton = await screen.findByText('Back');

    expect(backButton).toBeInTheDocument();

    await userEvent.click(backButton);
    expect(testProps.push).toHaveBeenCalledWith('/service-member/contact-info');
  });

  it('next button submits the form and goes to the Current Address step', async () => {
    const testServiceMemberValues = {
      id: 'testServiceMemberId',
      first_name: 'Leo',
      middle_name: 'Star',
      last_name: 'Spaceman',
      suffix: 'Mr.',
    };

    const testexistingDutyLocationValues = {
      address: {
        city: 'San Diego',
        state: 'CA',
        postalCode: '92104',
      },
      name: 'San Diego',
      id: 'testId',
    };

    patchServiceMember.mockImplementation(() => Promise.resolve(testServiceMemberValues));

    // Need to provide initial values because we aren't testing the form here, and just want to submit immediately
    render(
      <DutyLocation
        {...testProps}
        serviceMember={testServiceMemberValues}
        existingDutyLocation={testexistingDutyLocationValues}
      />,
    );

    const submitButton = screen.getByText('Next');
    expect(submitButton).toBeInTheDocument();
    await userEvent.click(submitButton);

    await waitFor(() => {
      expect(patchServiceMember).toHaveBeenCalled();
    });

    expect(testProps.updateServiceMember).toHaveBeenCalledWith(testServiceMemberValues);
    expect(testProps.push).toHaveBeenCalledWith('/service-member/current-address');
  });

  it('shows an error if the API returns an error', async () => {
    const testServiceMemberValues = {
      id: 'testServiceMemberId',
      first_name: 'Leo',
      middle_name: 'Star',
      last_name: 'Spaceman',
      suffix: 'Mr.',
    };

    const testexistingDutyLocationValues = {
      address: {
        city: 'San Diego',
        state: 'CA',
        postalCode: '92104',
      },
      name: 'San Diego',
      id: 'testId',
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
      <DutyLocation
        {...testProps}
        serviceMember={testServiceMemberValues}
        existingDutyLocation={testexistingDutyLocationValues}
      />,
    );

    const submitButton = screen.getByText('Next');
    expect(submitButton).toBeInTheDocument();
    await userEvent.click(submitButton);

    await waitFor(() => {
      expect(patchServiceMember).toHaveBeenCalled();
    });

    expect(screen.getByText('A server error occurred saving the service member')).toBeInTheDocument();
    expect(testProps.updateServiceMember).not.toHaveBeenCalled();
    expect(testProps.push).not.toHaveBeenCalled();
  });

  afterEach(jest.resetAllMocks);
});

describe('requireCustomerState DutyLocation', () => {
  const useDispatchMock = jest.spyOn(reactRedux, 'useDispatch');
  const mockDispatch = jest.fn();

  beforeEach(() => {
    useDispatchMock.mockClear();
    mockDispatch.mockClear();
    useDispatchMock.mockReturnValue(mockDispatch);
  });

  const props = {
    pages: ['first'],
    pageKey: '1',
    userEmail: 'my@email.com',
    schema: { my: 'schema' },
    updateServiceMember: jest.fn(),
    push: jest.fn(),
  };

  it('dispatches a redirect if the current state is earlier than the "CONTACT INFO COMPLETE" state', async () => {
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
        <ConnectedDutyLocation {...props} />
      </MockProviders>,
    );

    await waitFor(() => {
      expect(mockDispatch).toHaveBeenCalledWith(push('/service-member/contact-info'));
    });
  });

  it('does not redirect if the current state equals the "CONTACT INFO COMPLETE" state', async () => {
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
          },
        },
      },
    };

    render(
      <MockProviders initialState={mockState}>
        <ConnectedDutyLocation {...props} />
      </MockProviders>,
    );

    await waitFor(() => {
      expect(mockDispatch).not.toHaveBeenCalled();
    });
  });

  it('does not redirect if the current state is after the "CONTACT INFO COMPLETE" state and profile is not complete', async () => {
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
        <ConnectedDutyLocation {...props} />
      </MockProviders>,
    );

    await waitFor(() => {
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
        <ConnectedDutyLocation {...props} />
      </MockProviders>,
    );

    await waitFor(() => {
      expect(mockDispatch).toHaveBeenCalledWith(push('/'));
    });
  });
});
