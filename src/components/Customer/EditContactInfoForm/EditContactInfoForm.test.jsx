import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { Provider } from 'react-redux';

import EditContactInfoForm from './EditContactInfoForm';

import { configureStore } from 'shared/store';

describe('EditContactInfoForm component', () => {
  const testProps = {
    initialValues: {
      telephone: '915-555-1248',
      secondary_telephone: '512-555-9285',
      personal_email: 'mm@example.com',
      phone_is_preferred: false,
      email_is_preferred: true,
      residential_address: {
        streetAddress1: '235 Prospect Valley Road SE',
        streetAddress2: '#125',
        city: 'El Paso',
        state: 'TX',
        postalCode: '79912',
        country: {
          code: 'US',
          name: 'UNITED STATES',
          id: '791899e6-cd77-46f2-981b-176ecb8d7098',
        },
        countryID: '791899e6-cd77-46f2-981b-176ecb8d7098',
      },
      backup_mailing_address: {
        streetAddress1: '9 W 2nd Ave',
        streetAddress2: '',
        city: 'El Paso',
        state: 'TX',
        postalCode: '79936',
        country: {
          code: 'US',
          name: 'UNITED STATES',
          id: '791899e6-cd77-46f2-981b-176ecb8d7098',
        },
        countryID: '791899e6-cd77-46f2-981b-176ecb8d7098',
      },
      backup_contact: {
        name: 'Peyton Wing',
        email: 'pw@example.com',
        telephone: '915-555-8761',
      },
    },
    onSubmit: jest.fn().mockImplementation(() => Promise.resolve()),
    onCancel: jest.fn(),
  };

  it('renders the form inputs', async () => {
    const mockStore = configureStore({});

    render(
      <Provider store={mockStore.store}>
        <EditContactInfoForm {...testProps} />
      </Provider>,
    );

    const telephoneInput = await screen.findByLabelText(/Best contact phone/);

    expect(telephoneInput).toBeInstanceOf(HTMLInputElement);

    expect(telephoneInput).toHaveValue(testProps.initialValues.telephone);

    const secondaryPhoneInput = await screen.findByLabelText(/Alt. phone/);

    expect(secondaryPhoneInput).toBeInstanceOf(HTMLInputElement);

    expect(secondaryPhoneInput).toHaveValue(testProps.initialValues.secondary_telephone);

    const personalEmailInput = await screen.findByLabelText(/Personal email/);

    expect(personalEmailInput).toBeInstanceOf(HTMLInputElement);

    expect(personalEmailInput).toHaveValue(testProps.initialValues.personal_email);

    const nameInput = await screen.findByLabelText(/Name/);

    expect(nameInput).toBeInstanceOf(HTMLInputElement);

    expect(nameInput).toHaveValue(testProps.initialValues.backup_contact.name);

    // We have two sets of addresses and the labels are the same across both
    const address1Inputs = await screen.findAllByLabelText(/Address 1/);

    expect(address1Inputs.length).toBe(2);

    const [residentialAddress1, backupAddress1] = address1Inputs;

    expect(residentialAddress1).toBeInstanceOf(HTMLInputElement);
    expect(residentialAddress1).toHaveValue(testProps.initialValues.residential_address.streetAddress1);

    expect(backupAddress1).toBeInstanceOf(HTMLInputElement);
    expect(backupAddress1).toHaveValue(testProps.initialValues.backup_mailing_address.streetAddress1);

    const address2Inputs = await screen.findAllByLabelText(/Address 2/);

    expect(address2Inputs.length).toBe(2);

    const [residentialAddress2, backupAddress2] = address2Inputs;

    expect(residentialAddress2).toBeInstanceOf(HTMLInputElement);
    expect(residentialAddress2).toHaveValue(testProps.initialValues.residential_address.streetAddress2);

    expect(backupAddress2).toBeInstanceOf(HTMLInputElement);
    expect(backupAddress2).toHaveValue(testProps.initialValues.backup_mailing_address.streetAddress2);

    const cityInputs = screen.getAllByTestId(/City/);

    expect(cityInputs.length).toBe(2);

    const [residentialCity, backupCity] = cityInputs;

    expect(residentialCity).toBeInstanceOf(HTMLLabelElement);
    expect(residentialCity).toHaveTextContent(testProps.initialValues.residential_address.city);

    expect(backupCity).toBeInstanceOf(HTMLLabelElement);
    expect(backupCity).toHaveTextContent(testProps.initialValues.backup_mailing_address.city);

    const stateInputs = screen.getAllByTestId(/State/);

    expect(stateInputs.length).toBe(2);

    const [residentialState, backupState] = stateInputs;

    expect(residentialState).toBeInstanceOf(HTMLLabelElement);
    expect(residentialState).toHaveTextContent(testProps.initialValues.residential_address.state);

    expect(backupState).toBeInstanceOf(HTMLLabelElement);
    expect(backupState).toHaveTextContent(testProps.initialValues.backup_mailing_address.state);

    const zipInputs = screen.getAllByTestId(/ZIP/);

    expect(zipInputs.length).toBe(2);

    const [residentialZIP, backupZIP] = zipInputs;

    expect(residentialZIP).toBeInstanceOf(HTMLLabelElement);
    expect(residentialZIP).toHaveTextContent(testProps.initialValues.residential_address.postalCode);

    expect(backupZIP).toBeInstanceOf(HTMLLabelElement);
    expect(backupZIP).toHaveTextContent(testProps.initialValues.backup_mailing_address.postalCode);

    expect(
      screen.getAllByText(
        `${testProps.initialValues.residential_address.city}, ${testProps.initialValues.residential_address.state} ${testProps.initialValues.residential_address.postalCode} ()`,
      ),
    );
    expect(
      screen.getAllByText(
        `${testProps.initialValues.backup_mailing_address.city}, ${testProps.initialValues.backup_mailing_address.state} ${testProps.initialValues.backup_mailing_address.postalCode} ()`,
      ),
    );

    // These next few have the same label for different field types
    const phoneInputs = await screen.findAllByLabelText(/Phone/);

    expect(phoneInputs.length).toBe(2);

    const [phoneCheckbox, phoneTextField] = phoneInputs;

    expect(phoneCheckbox).toBeInstanceOf(HTMLInputElement);
    if (testProps.initialValues.phone_is_preferred) {
      expect(phoneCheckbox).toBeChecked();
    } else {
      expect(phoneCheckbox).not.toBeChecked();
    }

    expect(phoneTextField).toBeInstanceOf(HTMLInputElement);
    expect(phoneTextField).toHaveValue(testProps.initialValues.backup_contact.telephone);

    const emailInputs = await screen.findAllByLabelText(/Email/);

    expect(emailInputs.length).toBe(2);

    const [emailCheckbox, emailTextField] = emailInputs;

    expect(emailCheckbox).toBeInstanceOf(HTMLInputElement);
    if (testProps.initialValues.email_is_preferred) {
      expect(emailCheckbox).toBeChecked();
    } else {
      expect(emailCheckbox).not.toBeChecked();
    }

    expect(emailTextField).toBeInstanceOf(HTMLInputElement);
    expect(emailTextField).toHaveValue(testProps.initialValues.backup_contact.email);
  });

  it('shows an error message if trying to submit an invalid form', async () => {
    const mockStore = configureStore({});

    render(
      <Provider store={mockStore.store}>
        <EditContactInfoForm {...testProps} />
      </Provider>,
    );

    const saveButton = await screen.findByRole('button', { name: 'Save' });

    expect(saveButton).toBeEnabled();

    const emailInput = await screen.findByLabelText(/Personal email/);

    await userEvent.clear(emailInput);
    await userEvent.tab();

    const alert = await screen.findByRole('alert');

    expect(alert).toBeInTheDocument();

    expect(alert).toHaveTextContent('Required');

    expect(saveButton).toBeDisabled();
  });

  it('submits the form when its valid', async () => {
    const mockStore = configureStore({});

    render(
      <Provider store={mockStore.store}>
        <EditContactInfoForm {...testProps} />
      </Provider>,
    );

    const saveButton = screen.getByRole('button', { name: 'Save' });

    await userEvent.click(saveButton);

    const expectedParams = {
      ...testProps.initialValues,
    };

    await waitFor(() => {
      expect(testProps.onSubmit).toHaveBeenCalledWith(expectedParams, expect.anything());
    });
  });

  it('implements the onCancel handler when the Cancel button is clicked', async () => {
    const mockStore = configureStore({});

    render(
      <Provider store={mockStore.store}>
        <EditContactInfoForm {...testProps} />
      </Provider>,
    );

    const cancelButton = screen.getByRole('button', { name: 'Cancel' });

    await userEvent.click(cancelButton);

    await waitFor(() => {
      expect(testProps.onCancel).toHaveBeenCalled();
    });
  });

  afterEach(jest.resetAllMocks);
});
