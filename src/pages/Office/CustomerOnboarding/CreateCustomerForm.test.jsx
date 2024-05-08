import React from 'react';
import { render, fireEvent, waitFor, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import { CreateCustomerForm } from './CreateCustomerForm';

import { MockProviders } from 'testUtils';
import { createCustomerWithOktaOption } from 'services/ghcApi';

const mockNavigate = jest.fn();
jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom'),
  useNavigate: () => mockNavigate,
}));

jest.mock('services/ghcApi', () => ({
  ...jest.requireActual('services/ghcApi'),
  createCustomerWithOktaOption: jest.fn(),
}));

jest.mock('store/flash/actions', () => ({
  ...jest.requireActual('store/flash/actions'),
  setFlashMessage: jest.fn(),
}));

beforeEach(jest.resetAllMocks);

const fakePayload = {
  affiliation: 'ARMY',
  edipi: '1234567890',
  first_name: 'Shish',
  middle_name: 'Ka',
  last_name: 'Bob',
  suffix: 'Mr.',
  telephone: '555-555-5555',
  secondary_telephone: '999-867-5309',
  personal_email: 'tastyAndDelicious@mail.mil',
  phone_is_preferred: true,
  email_is_preferred: '',
  residential_address: {
    streetAddress1: '8711 S Hungry Ave.',
    streetAddress2: '',
    streetAddress3: '',
    city: 'Starving',
    state: 'OK',
    postalCode: '74133',
  },
  backup_mailing_address: {
    streetAddress1: '420 S. Munchies Lane',
    streetAddress2: '',
    streetAddress3: '',
    city: 'Mustang',
    state: 'KS',
    postalCode: '73064',
  },
  backup_contact: {
    name: 'Silly String',
    telephone: '666-666-6666',
    email: 'allOverDaPlace@mail.com',
  },
  create_okta_account: 'true',
  cac_user: 'false',
};

const fakeResponse = {
  affiliation: 'string',
  firstName: 'John',
  lastName: 'Doe',
  telephone: '216-421-1392',
  personalEmail: '73sGJ6jq7cS%6@PqElR.WUzkqFNvtduyyA',
  suffix: 'Jr.',
  middleName: 'David',
  residentialAddress: {
    id: 'c56a4180-65aa-42ec-a945-5fd21dec0538',
    streetAddress1: '123 Main Ave',
    streetAddress2: 'Apartment 9000',
    streetAddress3: 'Montmârtre',
    city: 'Anytown',
    eTag: 'string',
    state: 'AL',
    postalCode: '90210',
    country: 'USA',
  },
  backupContact: {
    name: 'string',
    email: 'backupContact@mail.com',
    phone: '381-100-5880',
  },
  id: 'c56a4180-65aa-42ec-a945-5fd21dec0538',
  edipi: 'string',
  userID: 'c56a4180-65aa-42ec-a945-5fd21dec0538',
  oktaID: 'string',
  oktaEmail: 'string',
  phoneIsPreferred: true,
  emailIsPreferred: true,
  secondaryTelephone: '499-793-2722',
  backupAddress: {
    id: 'c56a4180-65aa-42ec-a945-5fd21dec0538',
    streetAddress1: '123 Main Ave',
    streetAddress2: 'Apartment 9000',
    streetAddress3: 'Montmârtre',
    city: 'Anytown',
    eTag: 'string',
    state: 'AL',
    postalCode: '90210',
    country: 'USA',
  },
};

const testProps = {
  setFlashMessage: jest.fn(),
};

describe('CreateCustomerForm', () => {
  it('renders without crashing', async () => {
    render(
      <MockProviders>
        <CreateCustomerForm {...testProps} />
      </MockProviders>,
    );

    // checking that all headers exist
    expect(screen.getByText('Create Customer Profile')).toBeInTheDocument();
    expect(screen.getByText('Customer Affiliation')).toBeInTheDocument();
    expect(screen.getByText('Customer Name')).toBeInTheDocument();
    expect(screen.getByText('Contact Info')).toBeInTheDocument();
    expect(screen.getByText('Current Address')).toBeInTheDocument();
    expect(screen.getByText('Backup Address')).toBeInTheDocument();
    expect(screen.getByText('Backup Contact')).toBeInTheDocument();
    expect(screen.getByText('Okta Account')).toBeInTheDocument();
    expect(screen.getByText('Non-CAC Users')).toBeInTheDocument();

    const saveBtn = await screen.findByRole('button', { name: 'Save' });
    expect(saveBtn).toBeInTheDocument();
    expect(saveBtn).toBeDisabled();
    const cancelBtn = await screen.findByRole('button', { name: 'Cancel' });
    expect(cancelBtn).toBeInTheDocument();
  });

  it('navigates the user on cancel click', async () => {
    const { getByText } = render(
      <MockProviders>
        <CreateCustomerForm {...testProps} />
      </MockProviders>,
    );
    fireEvent.click(getByText('Cancel'));
    await waitFor(() => {
      expect(mockNavigate).toHaveBeenCalled();
    });
  });

  it('submits the form and navigates the user once all required fields are filled out', async () => {
    createCustomerWithOktaOption.mockImplementation(() => Promise.resolve(fakeResponse));

    const { getByLabelText, getByTestId, getByRole } = render(
      <MockProviders>
        <CreateCustomerForm {...testProps} />
      </MockProviders>,
    );

    const user = userEvent.setup();

    const saveBtn = await screen.findByRole('button', { name: 'Save' });
    expect(saveBtn).toBeInTheDocument();

    await user.selectOptions(getByLabelText('Branch of service'), [fakePayload.affiliation]);

    await user.type(getByLabelText('First name'), fakePayload.first_name);
    await user.type(getByLabelText('Last name'), fakePayload.last_name);

    await user.type(getByLabelText('Best contact phone'), fakePayload.telephone);
    await user.type(getByLabelText('Personal email'), fakePayload.personal_email);

    await user.type(getByTestId('res-add-street1'), fakePayload.residential_address.streetAddress1);
    await user.type(getByTestId('res-add-city'), fakePayload.residential_address.city);
    await user.selectOptions(getByTestId('res-add-state'), [fakePayload.residential_address.state]);
    await user.type(getByTestId('res-add-zip'), fakePayload.residential_address.postalCode);

    await user.type(getByTestId('backup-add-street1'), fakePayload.backup_mailing_address.streetAddress1);
    await user.type(getByTestId('backup-add-city'), fakePayload.backup_mailing_address.city);
    await user.selectOptions(getByTestId('backup-add-state'), [fakePayload.backup_mailing_address.state]);
    await user.type(getByTestId('backup-add-zip'), fakePayload.backup_mailing_address.postalCode);

    await user.type(getByLabelText('Name'), fakePayload.backup_contact.name);
    await user.type(getByRole('textbox', { name: 'Email' }), fakePayload.backup_contact.email);
    await user.type(getByRole('textbox', { name: 'Phone' }), fakePayload.backup_contact.telephone);

    await userEvent.type(getByTestId('create-okta-account-yes'), fakePayload.create_okta_account);

    await userEvent.type(getByTestId('cac-user-no'), fakePayload.cac_user);

    await waitFor(() => {
      expect(saveBtn).toBeEnabled();
    });

    const waiter = waitFor(() => {
      expect(createCustomerWithOktaOption).toHaveBeenCalled();
    });

    await user.click(saveBtn);
    await waiter;
    expect(mockNavigate).toHaveBeenCalled();
  }, 10000);

  it('submits the form and tests for unsupported state validation', async () => {
    createCustomerWithOktaOption.mockImplementation(() => Promise.resolve(fakeResponse));

    const { getByLabelText, getByTestId, getByRole, getByText } = render(
      <MockProviders>
        <CreateCustomerForm {...testProps} />
      </MockProviders>,
    );

    const user = userEvent.setup();

    const saveBtn = await screen.findByRole('button', { name: 'Save' });
    expect(saveBtn).toBeInTheDocument();

    await user.selectOptions(getByLabelText('Branch of service'), [fakePayload.affiliation]);

    await user.type(getByLabelText('First name'), fakePayload.first_name);
    await user.type(getByLabelText('Last name'), fakePayload.last_name);

    await user.type(getByLabelText('Best contact phone'), fakePayload.telephone);
    await user.type(getByLabelText('Personal email'), fakePayload.personal_email);

    await userEvent.type(getByTestId('res-add-street1'), fakePayload.residential_address.streetAddress1);
    await userEvent.type(getByTestId('res-add-city'), fakePayload.residential_address.city);
    await userEvent.selectOptions(getByTestId('res-add-state'), [fakePayload.residential_address.state]);
    await userEvent.type(getByTestId('res-add-zip'), fakePayload.residential_address.postalCode);

    await userEvent.type(getByTestId('backup-add-street1'), fakePayload.backup_mailing_address.streetAddress1);
    await userEvent.type(getByTestId('backup-add-city'), fakePayload.backup_mailing_address.city);
    await userEvent.selectOptions(getByTestId('backup-add-state'), [fakePayload.backup_mailing_address.state]);
    await userEvent.type(getByTestId('backup-add-zip'), fakePayload.backup_mailing_address.postalCode);

    await userEvent.type(getByLabelText('Name'), fakePayload.backup_contact.name);
    await userEvent.type(getByRole('textbox', { name: 'Email' }), fakePayload.backup_contact.email);
    await userEvent.type(getByRole('textbox', { name: 'Phone' }), fakePayload.backup_contact.telephone);

    await userEvent.type(getByTestId('create-okta-account-yes'), fakePayload.create_okta_account);

    await userEvent.type(getByTestId('cac-user-no'), fakePayload.cac_user);

    await waitFor(() => {
      expect(saveBtn).toBeEnabled();
    });

    await userEvent.selectOptions(getByTestId('backup-add-state'), 'AK');
    await userEvent.tab();

    const msg = getByText('Moves to this state are not supported at this time.');
    expect(msg).toBeVisible();

    await userEvent.selectOptions(getByTestId('backup-add-state'), [fakePayload.residential_address.state]);
    await userEvent.tab();
    expect(msg).not.toBeVisible();

    await waitFor(() => {
      expect(saveBtn).toBeEnabled();
    });

    await userEvent.click(saveBtn);

    expect(createCustomerWithOktaOption).toHaveBeenCalled();
    expect(mockNavigate).toHaveBeenCalled();
  }, 10000);
});
