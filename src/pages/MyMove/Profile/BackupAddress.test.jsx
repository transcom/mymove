import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import { MockProviders } from 'testUtils';
import ConnectedBackupAddress, { BackupAddress } from 'pages/MyMove/Profile/BackupAddress';
import { customerRoutes } from 'constants/routes';
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

describe('BackupAddress page', () => {
  const fakeAddress = {
    streetAddress1: '235 Prospect Valley Road SE',
    streetAddress2: '#125',
    city: 'El Paso',
    state: 'TX',
    postalCode: '79912',
  };

  const blankAddress = Object.fromEntries(Object.keys(fakeAddress).map((k) => [k, '']));

  const generateTestProps = (address) => ({
    updateServiceMember: jest.fn(),
    serviceMember: {
      id: 'testServiceMemberId',
      backup_mailing_address: address,
    },
  });

  it('renders the BackupAddressForm', async () => {
    const testProps = generateTestProps(blankAddress);

    const { queryByRole } = render(<BackupAddress {...testProps} />);

    await waitFor(() => {
      expect(queryByRole('heading', { name: 'Backup address', level: 1 })).toBeInTheDocument();
    });
  });

  it('back button goes to the Residential address step', async () => {
    const testProps = generateTestProps(blankAddress);

    const { findByRole } = render(<BackupAddress {...testProps} />);

    const backButton = await findByRole('button', { name: 'Back' });
    expect(backButton).toBeInTheDocument();
    await userEvent.click(backButton);

    expect(mockNavigate).toHaveBeenCalledWith(customerRoutes.CURRENT_ADDRESS_PATH);
  });

  it('next button submits the form and goes to the Backup contact step', async () => {
    const testProps = generateTestProps(blankAddress);

    const expectedServiceMemberPayload = { ...testProps.serviceMember, backup_mailing_address: fakeAddress };

    patchServiceMember.mockImplementation(() => Promise.resolve(expectedServiceMemberPayload));

    const { getByRole, getByLabelText } = render(<BackupAddress {...testProps} />);

    await userEvent.type(getByLabelText(/Address 1/), fakeAddress.streetAddress1);
    await userEvent.type(getByLabelText(/Address 2/), fakeAddress.streetAddress2);
    await userEvent.type(getByLabelText(/City/), fakeAddress.city);
    await userEvent.selectOptions(getByLabelText(/State/), [fakeAddress.state]);
    await userEvent.type(getByLabelText(/ZIP/), fakeAddress.postalCode);
    await userEvent.tab();

    const submitButton = getByRole('button', { name: 'Next' });
    expect(submitButton).toBeInTheDocument();
    await waitFor(() => {
      expect(submitButton).toBeEnabled();
    });

    await userEvent.click(submitButton);

    await waitFor(() => {
      expect(patchServiceMember).toHaveBeenCalledWith(expectedServiceMemberPayload);
    });

    expect(testProps.updateServiceMember).toHaveBeenCalledWith(expectedServiceMemberPayload);
    expect(mockNavigate).toHaveBeenCalledWith(customerRoutes.BACKUP_CONTACTS_PATH);
  });

  it('Selecting an unsupported state should display an unsupported state message', async () => {
    const testProps = generateTestProps(blankAddress);

    const expectedServiceMemberPayload = { ...testProps.serviceMember, backup_mailing_address: fakeAddress };

    patchServiceMember.mockImplementation(() => Promise.resolve(expectedServiceMemberPayload));

    const { getByLabelText } = render(<BackupAddress {...testProps} />);

    await userEvent.type(getByLabelText(/Address 1/), fakeAddress.streetAddress1);
    await userEvent.type(getByLabelText(/Address 2/), fakeAddress.streetAddress2);
    await userEvent.type(getByLabelText(/City/), fakeAddress.city);
    await userEvent.selectOptions(getByLabelText(/State/), 'HI');
    await userEvent.type(getByLabelText(/ZIP/), fakeAddress.postalCode);
    await userEvent.tab();

    let msg = screen.getByText('Moves to this state are not supported at this time.');
    expect(msg).toBeVisible();

    await userEvent.selectOptions(getByLabelText(/State/), 'AL');
    await userEvent.type(getByLabelText(/ZIP/), fakeAddress.postalCode);
    await userEvent.tab();
    expect(msg).not.toBeVisible();

    await userEvent.selectOptions(getByLabelText(/State/), 'HI');
    await userEvent.type(getByLabelText(/ZIP/), fakeAddress.postalCode);
    await userEvent.tab();
    msg = screen.getByText('Moves to this state are not supported at this time.');
    expect(msg).toBeVisible();
  });

  it('shows an error if the patchServiceMember API returns an error', async () => {
    const testProps = generateTestProps(fakeAddress);

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

    const { getByRole, queryByText } = render(<BackupAddress {...testProps} />);

    const submitButton = getByRole('button', { name: 'Next' });
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

describe('requireCustomerState BackupAddress', () => {
  const props = {
    updateServiceMember: jest.fn(),
  };

  it('dispatches a redirect if the current state is earlier than the "ADDRESS COMPLETE" state', async () => {
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
        <ConnectedBackupAddress {...props} />
      </MockProviders>,
    );

    await waitFor(() => {
      expect(mockNavigate).toHaveBeenCalledWith('/service-member/current-address');
    });
  });

  it('does not redirect if the current state equals the "ADDRESS COMPLETE" state', async () => {
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
          },
        },
      },
    };

    render(
      <MockProviders initialState={mockState}>
        <ConnectedBackupAddress {...props} />
      </MockProviders>,
    );

    await waitFor(() => {
      expect(mockNavigate).not.toHaveBeenCalled();
    });
  });

  it('does not redirect if the current state is after the "ADDRESS COMPLETE" state and profile is not complete', async () => {
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
          },
        },
      },
    };

    render(
      <MockProviders initialState={mockState}>
        <ConnectedBackupAddress {...props} />
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
        <ConnectedBackupAddress {...props} />
      </MockProviders>,
    );

    await waitFor(() => {
      expect(mockNavigate).toHaveBeenCalledWith('/');
    });
  });
});
