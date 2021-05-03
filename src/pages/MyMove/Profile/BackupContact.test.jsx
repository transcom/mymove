import React from 'react';
import { mount } from 'enzyme';
import * as reactRedux from 'react-redux';
import { push } from 'connected-react-router';
import { render, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import ConnectedBackupContact, { BackupContact } from './BackupContact';

import { MockProviders } from 'testUtils';
import { createBackupContactForServiceMember, patchBackupContact, getServiceMember } from 'services/internalApi';

jest.mock('services/internalApi', () => ({
  ...jest.requireActual('services/internalApi'),
  createBackupContactForServiceMember: jest.fn(),
  patchBackupContact: jest.fn(),
  getServiceMember: jest.fn(),
}));

describe('BackupContact page', () => {
  const testProps = {
    updateServiceMember: jest.fn(),
    updateBackupContact: jest.fn(),
    push: jest.fn(),
    serviceMember: {
      id: 'testServiceMemberId',
    },
    currentBackupContacts: [],
  };

  const testBackupContactValues = {
    name: 'Ima Goddess',
    telephone: '555-555-5555',
    email: 'test@example.com',
    // permission: 'NONE',
  };

  const testBackupContacts = [testBackupContactValues];

  it('renders the BackupContactForm', async () => {
    const { queryByRole } = render(<BackupContact {...testProps} />);

    await waitFor(() => {
      expect(queryByRole('heading', { name: 'Backup contact', level: 1 })).toBeInTheDocument();
    });
  });

  it('back button goes to the BACKUP ADDRESS step', async () => {
    // Need to provide initial values because we aren't testing the form here, and just want to submit immediately
    const { queryByText } = render(<BackupContact {...testProps} currentBackupContacts={testBackupContacts} />);

    const backButton = queryByText('Back');
    expect(backButton).toBeInTheDocument();
    userEvent.click(backButton);

    expect(testProps.push).toHaveBeenCalledWith('/service-member/backup-address');
  });
  describe('if there is an existing backup contact', () => {
    it('next button submits the form and goes to the Home step', async () => {
      patchBackupContact.mockImplementation(() => Promise.resolve(testBackupContactValues));
      getServiceMember.mockImplementation(() => Promise.resolve(testProps.serviceMember));
      testProps.updateServiceMember.mockImplementation(() => Promise.resolve({}));

      // Need to provide initial values because we aren't testing the form here, and just want to submit immediately
      const { queryByText } = render(<BackupContact {...testProps} currentBackupContacts={testBackupContacts} />);

      const submitButton = queryByText('Next');
      expect(submitButton).toBeInTheDocument();
      userEvent.click(submitButton);

      await waitFor(() => {
        expect(patchBackupContact).toHaveBeenCalled();
      });

      expect(testProps.updateBackupContact).toHaveBeenCalledWith(testBackupContactValues);
      expect(getServiceMember).toHaveBeenCalledWith(testProps.serviceMember.id);
      expect(testProps.updateServiceMember).toHaveBeenCalledWith(testProps.serviceMember);
      expect(testProps.push).toHaveBeenCalledWith('/');
    });

    it('shows an error if the API returns an error', async () => {
      patchBackupContact.mockImplementation(() =>
        // Disable this rule because makeSwaggerRequest does not throw an error if the API call fails
        // eslint-disable-next-line prefer-promise-reject-errors
        Promise.reject({
          message: 'A server error occurred saving the backup contact',
          response: {
            body: {
              detail: 'A server error occurred saving the backup contact',
            },
          },
        }),
      );

      // Need to provide complete & valid initial values because we aren't testing the form here, and just want to submit immediately
      const { queryByText } = render(<BackupContact {...testProps} currentBackupContacts={testBackupContacts} />);

      const submitButton = queryByText('Next');
      expect(submitButton).toBeInTheDocument();
      userEvent.click(submitButton);

      await waitFor(() => {
        expect(patchBackupContact).toHaveBeenCalled();
      });

      expect(queryByText('A server error occurred saving the backup contact')).toBeInTheDocument();
      expect(testProps.updateBackupContact).not.toHaveBeenCalled();
      expect(testProps.push).not.toHaveBeenCalled();
    });
  });
  describe('if there is no existing backup contact', () => {
    it('next button submits the form and goes to the Home step', async () => {
      createBackupContactForServiceMember.mockImplementation(() => Promise.resolve(testBackupContactValues));
      getServiceMember.mockImplementation(() => Promise.resolve(testProps.serviceMember));
      testProps.updateServiceMember.mockImplementation(() => Promise.resolve({}));

      // Need to provide initial values because we aren't testing the form here, and just want to submit immediately
      const { queryByText, getByLabelText } = render(<BackupContact {...testProps} />);

      const submitButton = queryByText('Next');
      expect(submitButton).toBeInTheDocument();
      userEvent.type(getByLabelText('Name'), 'Joe Schmoe');
      userEvent.type(getByLabelText('Phone'), '555-555-5555');
      userEvent.type(getByLabelText('Email'), 'test@sample.com');
      userEvent.click(submitButton);

      await waitFor(() => {
        expect(createBackupContactForServiceMember).toHaveBeenCalled();
      });

      expect(testProps.updateBackupContact).toHaveBeenCalledWith(testBackupContactValues);
      expect(getServiceMember).toHaveBeenCalledWith(testProps.serviceMember.id);
      expect(testProps.updateServiceMember).toHaveBeenCalledWith(testProps.serviceMember);
      expect(testProps.push).toHaveBeenCalledWith('/');
    });

    it('shows an error if the API returns an error', async () => {
      createBackupContactForServiceMember.mockImplementation(() =>
        // Disable this rule because makeSwaggerRequest does not throw an error if the API call fails
        // eslint-disable-next-line prefer-promise-reject-errors
        Promise.reject({
          message: 'A server error occurred saving the backup contact',
          response: {
            body: {
              detail: 'A server error occurred saving the backup contact',
            },
          },
        }),
      );

      // Need to provide complete & valid initial values because we aren't testing the form here, and just want to submit immediately
      const { queryByText, getByLabelText } = render(<BackupContact {...testProps} />);

      const submitButton = queryByText('Next');
      expect(submitButton).toBeInTheDocument();
      userEvent.type(getByLabelText('Name'), 'Joe Schmitty');
      userEvent.type(getByLabelText('Phone'), '555-555-5555');
      userEvent.type(getByLabelText('Email'), 'test@sample.com');
      userEvent.click(submitButton);

      await waitFor(() => {
        expect(createBackupContactForServiceMember).toHaveBeenCalled();
      });

      expect(queryByText('A server error occurred saving the backup contact')).toBeInTheDocument();
      expect(testProps.updateBackupContact).not.toHaveBeenCalled();
      expect(testProps.push).not.toHaveBeenCalled();
    });
  });
  afterEach(jest.resetAllMocks);
});

describe('requireCustomerState BackupContact', () => {
  const useDispatchMock = jest.spyOn(reactRedux, 'useDispatch');
  const mockDispatch = jest.fn();

  beforeEach(() => {
    useDispatchMock.mockClear();
    mockDispatch.mockClear();
    useDispatchMock.mockReturnValue(mockDispatch);
  });

  const props = {
    updateServiceMember: jest.fn(),
    updateBackupContact: jest.fn(),
    push: jest.fn(),
    serviceMember: {
      id: 'testServiceMemberId',
    },
    currentBackupContacts: [],
  };

  it('dispatches a redirect if the current state is earlier than the "BACKUP MAILING ADDRESS COMPLETE" state', () => {
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
        <ConnectedBackupContact {...props} />
      </MockProviders>,
    );

    expect(wrapper.exists()).toBe(true);
    expect(mockDispatch).toHaveBeenCalledWith(push('/service-member/backup-address'));
  });

  it('does not redirect if the current state equals the "BACKUP MAILING ADDRESS COMPLETE" state', () => {
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
        <ConnectedBackupContact {...props} />
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
        <ConnectedBackupContact {...props} />
      </MockProviders>,
    );

    expect(wrapper.exists()).toBe(true);
    expect(mockDispatch).toHaveBeenCalledWith(push('/'));
  });
});
