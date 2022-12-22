import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import { EditContactInfo } from './EditContactInfo';

import { patchBackupContact, patchServiceMember } from 'services/internalApi';
import { customerRoutes } from 'constants/routes';

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
  patchBackupContact: jest.fn(),
  patchServiceMember: jest.fn(),
}));

describe('EditContactInfo page', () => {
  const testProps = {
    currentBackupContacts: [
      {
        id: 'backupContactID',
        name: 'Barbara St. Juste',
        email: 'bsj@example.com',
        telephone: '915-555-1234',
        permission: 'NONE',
      },
    ],
    serviceMember: {
      id: 'testServiceMemberID',
      telephone: '915-555-2945',
      secondary_telephone: '',
      personal_email: 'test@example.com',
      email_is_preferred: true,
      phone_is_preferred: false,
      residential_address: {
        streetAddress1: '148 S East St',
        streetAddress2: '',
        city: 'Fake City',
        state: 'TX',
        postalCode: '79936',
      },
      backup_mailing_address: {
        streetAddress1: '10642 N Second Ave',
        streetAddress2: '',
        city: 'Fake City',
        state: 'TX',
        postalCode: '79936',
      },
    },
    setFlashMessage: jest.fn(),
    updateBackupContact: jest.fn(),
    updateServiceMember: jest.fn(),
  };

  it('renders the EditContactInfo form', async () => {
    render(<EditContactInfo {...testProps} />);

    const h1 = await screen.findByRole('heading', { name: 'Edit contact info', level: 1 });
    expect(h1).toBeInTheDocument();

    const contactHeader = screen.getByRole('heading', { name: 'Your contact info', level: 2 });
    expect(contactHeader).toBeInTheDocument();

    const addressHeader = screen.getByRole('heading', { name: 'Current mailing address', level: 2 });
    expect(addressHeader).toBeInTheDocument();

    const backupAddressHeader = screen.getByRole('heading', { name: 'Backup mailing address', level: 2 });
    expect(backupAddressHeader).toBeInTheDocument();

    const backupContactHeader = screen.getByRole('heading', { name: 'Backup contact', level: 2 });
    expect(backupContactHeader).toBeInTheDocument();
  });

  it('goes back to the profile page when the cancel button is clicked', async () => {
    render(<EditContactInfo {...testProps} />);

    const cancelButton = await screen.findByRole('button', { name: 'Cancel' });

    await userEvent.click(cancelButton);

    expect(mockPush).toHaveBeenCalledWith(customerRoutes.PROFILE_PATH);
  });

  it('saves backup contact info when it is updated and the save button is clicked', async () => {
    const newName = 'Rosalie Wexler';

    const expectedPayload = { ...testProps.currentBackupContacts[0], name: newName };

    const patchResponse = {
      ...expectedPayload,
      serviceMemberId: testProps.serviceMember.id,
      created_at: '2021-02-08T16:48:04.117Z',
      updated_at: '2021-02-11T16:48:04.117Z',
    };

    patchBackupContact.mockImplementation(() => Promise.resolve(patchResponse));
    patchServiceMember.mockImplementation(() => Promise.resolve());

    render(<EditContactInfo {...testProps} />);

    const backupNameInput = await screen.findByLabelText('Name');

    await userEvent.clear(backupNameInput);

    await userEvent.type(backupNameInput, newName);

    const saveButton = screen.getByRole('button', { name: 'Save' });

    await userEvent.click(saveButton);

    await waitFor(() => {
      expect(patchBackupContact).toHaveBeenCalledWith(expectedPayload);
    });

    expect(testProps.updateBackupContact).toHaveBeenCalledWith(patchResponse);
  });

  it('shows an error if the patchBackupContact API returns an error', async () => {
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

    render(<EditContactInfo {...testProps} />);

    const backupNameInput = await screen.findByLabelText('Name');

    await userEvent.clear(backupNameInput);

    await userEvent.type(backupNameInput, 'Rosalie Wexler');

    const saveButton = screen.getByRole('button', { name: 'Save' });

    await userEvent.click(saveButton);

    await waitFor(() => {
      expect(patchBackupContact).toHaveBeenCalled();
    });

    expect(await screen.findByText('A server error occurred saving the backup contact')).toBeInTheDocument();
    expect(testProps.updateBackupContact).not.toHaveBeenCalled();
    expect(patchServiceMember).not.toHaveBeenCalled();
    expect(testProps.updateServiceMember).not.toHaveBeenCalled();
    expect(testProps.setFlashMessage).not.toHaveBeenCalled();
    expect(mockPush).not.toHaveBeenCalled();
  });

  it('does not save backup contact info if it is not updated and the save button is clicked', async () => {
    patchServiceMember.mockImplementation(() => Promise.resolve());

    render(<EditContactInfo {...testProps} />);

    const saveButton = screen.getByRole('button', { name: 'Save' });

    await userEvent.click(saveButton);

    await waitFor(() => {
      expect(patchBackupContact).not.toHaveBeenCalled();
    });

    expect(testProps.updateBackupContact).not.toHaveBeenCalled();
  });

  it('saves service member info when the save button is clicked', async () => {
    const secondaryPhone = '915-555-9753';

    const expectedPayload = { ...testProps.serviceMember, secondary_telephone: secondaryPhone };

    const patchResponse = {
      ...expectedPayload,
      created_at: '2021-02-08T16:48:04.117Z',
      updated_at: '2021-02-11T16:48:04.117Z',
    };

    patchServiceMember.mockImplementation(() => Promise.resolve(patchResponse));

    render(<EditContactInfo {...testProps} />);

    const secondaryPhoneInput = await screen.findByLabelText(/Alt. phone/);

    await userEvent.clear(secondaryPhoneInput);

    await userEvent.type(secondaryPhoneInput, secondaryPhone);

    const saveButton = screen.getByRole('button', { name: 'Save' });

    await userEvent.click(saveButton);

    await waitFor(() => {
      expect(patchServiceMember).toHaveBeenCalledWith(expectedPayload);
    });

    expect(testProps.updateServiceMember).toHaveBeenCalledWith(patchResponse);
  });

  it('sets a flash message when the save button is clicked', async () => {
    patchServiceMember.mockImplementation(() => Promise.resolve());

    render(<EditContactInfo {...testProps} />);

    const saveButton = screen.getByRole('button', { name: 'Save' });

    await userEvent.click(saveButton);

    await waitFor(() => {
      expect(testProps.setFlashMessage).toHaveBeenCalledWith(
        'EDIT_CONTACT_INFO_SUCCESS',
        'success',
        "You've updated your information.",
      );
    });
  });

  it('routes to the profile page when the save button is clicked', async () => {
    patchServiceMember.mockImplementation(() => Promise.resolve());

    render(<EditContactInfo {...testProps} />);

    const saveButton = screen.getByRole('button', { name: 'Save' });

    await userEvent.click(saveButton);

    await waitFor(() => {
      expect(mockPush).toHaveBeenCalledWith(customerRoutes.PROFILE_PATH);
    });
  });

  it('shows an error if the patchServiceMember API returns an error', async () => {
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

    render(<EditContactInfo {...testProps} />);

    const saveButton = screen.getByRole('button', { name: 'Save' });

    await userEvent.click(saveButton);

    await waitFor(() => {
      expect(patchServiceMember).toHaveBeenCalled();
    });

    expect(await screen.findByText('A server error occurred saving the service member')).toBeInTheDocument();
    expect(testProps.updateServiceMember).not.toHaveBeenCalled();
    expect(testProps.setFlashMessage).not.toHaveBeenCalled();
    expect(mockPush).not.toHaveBeenCalled();
  });

  afterEach(jest.resetAllMocks);
});
