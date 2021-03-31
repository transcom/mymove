import React from 'react';
import { mount } from 'enzyme';
import * as reactRedux from 'react-redux';
import { push } from 'connected-react-router';
import { render, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import { MockProviders } from 'testUtils';
import ConnectedBackupMailingAddress, { BackupMailingAddress } from 'pages/MyMove/Profile/BackupMailingAddress';
import { customerRoutes } from 'constants/routes';
import { patchServiceMember } from 'services/internalApi';

jest.mock('services/internalApi', () => ({
  ...jest.requireActual('services/internalApi'),
  patchServiceMember: jest.fn(),
}));

describe('BackupMailingAddress page', () => {
  const fakeAddress = {
    street_address_1: '235 Prospect Valley Road SE',
    street_address_2: '#125',
    city: 'El Paso',
    state: 'TX',
    postal_code: '79912',
  };

  const blankAddress = Object.fromEntries(Object.keys(fakeAddress).map((k) => [k, '']));

  const generateTestProps = (address) => ({
    updateServiceMember: jest.fn(),
    push: jest.fn(),
    serviceMember: {
      id: 'testServiceMemberId',
      backup_mailing_address: address,
    },
  });

  it('renders the BackupMailingAddressForm', async () => {
    const testProps = generateTestProps(blankAddress);

    const { queryByRole } = render(<BackupMailingAddress {...testProps} />);

    await waitFor(() => {
      expect(queryByRole('heading', { name: 'Backup mailing address', level: 1 })).toBeInTheDocument();
    });
  });

  it('back button goes to the Residential address step', async () => {
    const testProps = generateTestProps(blankAddress);

    const { findByRole } = render(<BackupMailingAddress {...testProps} />);

    const backButton = await findByRole('button', { name: 'Back' });
    expect(backButton).toBeInTheDocument();
    userEvent.click(backButton);

    expect(testProps.push).toHaveBeenCalledWith(customerRoutes.CURRENT_ADDRESS_PATH);
  });

  it('next button submits the form and goes to the Backup contact step', async () => {
    const testProps = generateTestProps(blankAddress);

    const expectedServiceMemberPayload = { ...testProps.serviceMember, backup_mailing_address: fakeAddress };

    patchServiceMember.mockImplementation(() => Promise.resolve(expectedServiceMemberPayload));

    const { getByRole, getByLabelText } = render(<BackupMailingAddress {...testProps} />);

    userEvent.type(getByLabelText('Address 1'), fakeAddress.street_address_1);
    userEvent.type(getByLabelText(/Address 2/), fakeAddress.street_address_2);
    userEvent.type(getByLabelText('City'), fakeAddress.city);
    userEvent.selectOptions(getByLabelText('State'), [fakeAddress.state]);
    userEvent.type(getByLabelText('ZIP'), fakeAddress.postal_code);

    const submitButton = getByRole('button', { name: 'Next' });
    expect(submitButton).toBeInTheDocument();
    userEvent.click(submitButton);

    await waitFor(() => {
      expect(patchServiceMember).toHaveBeenCalledWith(expectedServiceMemberPayload);
    });

    expect(testProps.updateServiceMember).toHaveBeenCalledWith(expectedServiceMemberPayload);
    expect(testProps.push).toHaveBeenCalledWith(customerRoutes.BACKUP_CONTACTS_PATH);
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

    const { getByRole, queryByText } = render(<BackupMailingAddress {...testProps} />);

    const submitButton = getByRole('button', { name: 'Next' });
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

describe('requireCustomerState BackupMailingAddress', () => {
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

  it('dispatches a redirect if the current state is earlier than the "ADDRESS COMPLETE" state', () => {
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
        <ConnectedBackupMailingAddress {...props} />
      </MockProviders>,
    );

    expect(wrapper.exists()).toBe(true);
    expect(mockDispatch).toHaveBeenCalledWith(push('/service-member/current-address'));
  });

  it('does not redirect if the current state equals the "ADDRESS COMPLETE" state', () => {
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
          },
        },
      },
    };

    const wrapper = mount(
      <MockProviders initialState={mockState}>
        <ConnectedBackupMailingAddress {...props} />
      </MockProviders>,
    );

    expect(wrapper.exists()).toBe(true);
    expect(mockDispatch).not.toHaveBeenCalled();
  });

  it('does not redirect if the current state is after the "ADDRESS COMPLETE" state and profile is not complete', () => {
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
          },
        },
      },
    };

    const wrapper = mount(
      <MockProviders initialState={mockState}>
        <ConnectedBackupMailingAddress {...props} />
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
        <ConnectedBackupMailingAddress {...props} />
      </MockProviders>,
    );

    expect(wrapper.exists()).toBe(true);
    expect(mockDispatch).toHaveBeenCalledWith(push('/'));
  });
});
