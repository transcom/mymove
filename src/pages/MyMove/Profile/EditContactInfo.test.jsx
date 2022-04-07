import React from 'react';
import { faker } from '@faker-js/faker';
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { v4 as uuidv4 } from 'uuid';

import { EditContactInfo } from './EditContactInfo';

import { customerRoutes } from 'constants/routes';
import { patchBackupContact, patchServiceMember } from 'services/internalApi';
import { PHONE_FORMAT, serviceMemberBuilder } from 'utils/test/factories/serviceMember';

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
        id: uuidv4(),
        name: faker.name.findName(),
        email: faker.internet.exampleEmail(),
        telephone: faker.phone.phoneNumber(PHONE_FORMAT),
        permission: 'NONE',
      },
    ],
    serviceMember: serviceMemberBuilder(),
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

    userEvent.click(cancelButton);

    expect(mockPush).toHaveBeenCalledWith(customerRoutes.PROFILE_PATH);
  });

  it('saves backup contact info when it is updated and the save button is clicked', async () => {
    const newName = faker.name.findName();

    const expectedPayload = { ...testProps.currentBackupContacts[0], name: newName };

    const createdAt = faker.datatype.datetime();
    const patchResponse = {
      ...expectedPayload,
      serviceMemberId: testProps.serviceMember.id,
      created_at: createdAt,
      updated_at: faker.datatype.datetime({ min: createdAt.getTime() }),
    };

    patchBackupContact.mockImplementation(() => Promise.resolve(patchResponse));
    patchServiceMember.mockImplementation(() => Promise.resolve());

    render(<EditContactInfo {...testProps} />);

    const backupNameInput = await screen.findByLabelText('Name');

    userEvent.clear(backupNameInput);

    userEvent.type(backupNameInput, newName);

    const saveButton = screen.getByRole('button', { name: 'Save' });

    userEvent.click(saveButton);

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

    userEvent.clear(backupNameInput);

    userEvent.type(backupNameInput, faker.name.findName());

    const saveButton = screen.getByRole('button', { name: 'Save' });

    userEvent.click(saveButton);

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

    userEvent.click(saveButton);

    await waitFor(() => {
      expect(patchBackupContact).not.toHaveBeenCalled();
    });

    expect(testProps.updateBackupContact).not.toHaveBeenCalled();
  });

  it('saves service member info when the save button is clicked', async () => {
    const secondaryPhone = faker.phone.phoneNumber(PHONE_FORMAT);

    const expectedPayload = {
      id: testProps.serviceMember.id,
      telephone: testProps.serviceMember.telephone,
      secondary_telephone: secondaryPhone,
      personal_email: testProps.serviceMember.personal_email,
      phone_is_preferred: testProps.serviceMember.phone_is_preferred,
      email_is_preferred: testProps.serviceMember.email_is_preferred,
      residential_address: {
        streetAddress1: testProps.serviceMember.residential_address.streetAddress1,
        streetAddress2: testProps.serviceMember.residential_address.streetAddress2,
        city: testProps.serviceMember.residential_address.city,
        state: testProps.serviceMember.residential_address.state,
        postalCode: testProps.serviceMember.residential_address.postalCode,
      },
      backup_mailing_address: {
        streetAddress1: testProps.serviceMember.backup_mailing_address.streetAddress1,
        streetAddress2: testProps.serviceMember.backup_mailing_address.streetAddress2,
        city: testProps.serviceMember.backup_mailing_address.city,
        state: testProps.serviceMember.backup_mailing_address.state,
        postalCode: testProps.serviceMember.backup_mailing_address.postalCode,
      },
    };

    const createdAt = faker.datatype.datetime();
    const patchResponse = {
      ...expectedPayload,
      created_at: createdAt,
      updated_at: faker.datatype.datetime({ min: createdAt.getTime() }),
    };

    patchServiceMember.mockImplementation(() => Promise.resolve(patchResponse));

    render(<EditContactInfo {...testProps} />);

    const secondaryPhoneInput = await screen.findByLabelText(/Alt. phone/);

    userEvent.clear(secondaryPhoneInput);

    userEvent.type(secondaryPhoneInput, secondaryPhone);

    const saveButton = screen.getByRole('button', { name: 'Save' });

    userEvent.click(saveButton);

    await waitFor(() => {
      expect(patchServiceMember).toHaveBeenCalledWith(expectedPayload);
    });

    expect(testProps.updateServiceMember).toHaveBeenCalledWith(patchResponse);
  });

  it('sets a flash message when the save button is clicked', async () => {
    patchServiceMember.mockImplementation(() => Promise.resolve());

    render(<EditContactInfo {...testProps} />);

    const saveButton = screen.getByRole('button', { name: 'Save' });

    userEvent.click(saveButton);

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

    userEvent.click(saveButton);

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

    userEvent.click(saveButton);

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
