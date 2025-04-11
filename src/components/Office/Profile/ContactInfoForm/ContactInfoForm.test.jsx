import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { Provider } from 'react-redux';

import ContactInfoForm from './ContactInfoForm';

import { configureStore } from 'shared/store';

describe('ContactInfoForm component', () => {
  const testProps = {
    initialValues: {
      firstName: 'John',
      middleName: 'F',
      lastName: 'Doe',
      telephone: '915-555-1248',
      email: 'mm@example.com',
    },
    onSubmit: jest.fn().mockImplementation(() => Promise.resolve()),
    onCancel: jest.fn(),
  };

  it('renders the form fields', async () => {
    const mockStore = configureStore({});

    render(
      <Provider store={mockStore.store}>
        <ContactInfoForm {...testProps} />
      </Provider>,
    );

    expect(await screen.findByRole('heading', { name: 'Your contact info', level: 2 })).toBeInTheDocument();

    const telephoneInput = await screen.findByLabelText('Phone *');
    const emailInput = await screen.findByLabelText('Email');
    const firstNameInput = await screen.findByLabelText('First name');
    const middleNameInput = await screen.findByLabelText('Middle name');
    const lastNameInput = await screen.findByLabelText('Last name');

    expect(telephoneInput).toBeInstanceOf(HTMLInputElement);
    expect(telephoneInput).toHaveValue(testProps.initialValues.telephone);

    expect(emailInput).toBeInstanceOf(HTMLInputElement);
    expect(emailInput).toHaveValue(testProps.initialValues.email);

    expect(firstNameInput).toBeInstanceOf(HTMLInputElement);
    expect(firstNameInput).toHaveValue(testProps.initialValues.firstName);

    expect(middleNameInput).toBeInstanceOf(HTMLInputElement);
    expect(middleNameInput).toHaveValue(testProps.initialValues.middleName);

    expect(lastNameInput).toBeInstanceOf(HTMLInputElement);
    expect(lastNameInput).toHaveValue(testProps.initialValues.lastName);
  });

  it('shows an error message if trying to submit an invalid form', async () => {
    const mockStore = configureStore({});

    render(
      <Provider store={mockStore.store}>
        <ContactInfoForm {...testProps} />
      </Provider>,
    );

    const saveButton = await screen.findByRole('button', { name: 'Save' });

    expect(saveButton).toBeEnabled();

    const telephoneInput = await screen.findByLabelText('Phone *');

    await userEvent.clear(telephoneInput);
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
        <ContactInfoForm {...testProps} />
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
        <ContactInfoForm {...testProps} />
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
