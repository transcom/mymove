import React from 'react';
import { render, fireEvent, waitFor, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import { CreateCustomerForm } from './CreateCustomerForm';

import { MockProviders } from 'testUtils';

const mockNavigate = jest.fn();
jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom'),
  useNavigate: () => mockNavigate,
}));

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
};

describe('CreateCustomerForm', () => {
  it('renders without crashing', async () => {
    render(
      <MockProviders>
        <CreateCustomerForm />
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

    const saveBtn = await screen.findByRole('button', { name: 'Save' });
    expect(saveBtn).toBeInTheDocument();
    expect(saveBtn).toBeDisabled();
    const cancelBtn = await screen.findByRole('button', { name: 'Cancel' });
    expect(cancelBtn).toBeInTheDocument();
  });

  it('navigates the user on cancel click', async () => {
    const { getByText } = render(
      <MockProviders>
        <CreateCustomerForm />
      </MockProviders>,
    );
    fireEvent.click(getByText('Cancel'));
    await waitFor(() => {
      expect(mockNavigate).toHaveBeenCalled();
    });
  });

  it('submits the form and navigates the user once all required fields are filled out', async () => {
    const { getByLabelText, getByTestId, getByRole } = render(
      <MockProviders>
        <CreateCustomerForm />
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

    const oktaRadioButton = getByLabelText('Yes');
    await userEvent.click(oktaRadioButton);

    await waitFor(() => {
      expect(saveBtn).toBeEnabled();
    });
    await userEvent.click(saveBtn);
    expect(mockNavigate).toHaveBeenCalled();
  });
});
