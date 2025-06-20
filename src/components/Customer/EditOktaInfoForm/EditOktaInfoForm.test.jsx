import React from 'react';
import { screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import EditOktaInfoForm from './EditOktaInfoForm';

import { renderWithRouter } from 'testUtils';
import { isBooleanFlagEnabled } from 'utils/featureFlags';

jest.mock('utils/featureFlags', () => ({
  ...jest.requireActual('utils/featureFlags'),
  isBooleanFlagEnabled: jest.fn().mockImplementation(() => Promise.resolve(false)),
}));

describe('EditOktaInfoForm component', () => {
  const testProps = {
    initialValues: {
      oktaUsername: 'user@okta.mil',
      oktaEmail: 'user@okta.mil',
      oktaFirstName: 'Lucky',
      oktaLastName: 'Shamrock',
      oktaEdipi: '1112223334',
    },
    onSubmit: jest.fn().mockImplementation(() => Promise.resolve()),
    onCancel: jest.fn(),
  };

  const testPropsWithLock = {
    initialValues: {
      oktaUsername: 'user@okta.mil',
      oktaEmail: 'user@okta.mil',
      oktaFirstName: 'Lucky',
      oktaLastName: 'Shamrock',
      oktaEdipi: '1112223334',
    },
    onSubmit: jest.fn().mockImplementation(() => Promise.resolve()),
    onCancel: jest.fn(),
    isMoveLocked: true,
  };

  it('renders the form inputs', async () => {
    isBooleanFlagEnabled.mockImplementation(() => Promise.resolve(true));

    renderWithRouter(<EditOktaInfoForm {...testProps} />);

    const oktaUsername = await screen.findByLabelText('Okta Username *');
    expect(oktaUsername).toBeInstanceOf(HTMLInputElement);
    expect(oktaUsername).toHaveValue(testProps.initialValues.oktaUsername);

    const oktaEmail = await screen.findByLabelText('Okta Email *');
    expect(oktaEmail).toBeInstanceOf(HTMLInputElement);
    expect(oktaEmail).toHaveValue(testProps.initialValues.oktaEmail);

    const oktaFirstName = await screen.findByLabelText('First Name *');
    expect(oktaFirstName).toBeInstanceOf(HTMLInputElement);
    expect(oktaFirstName).toHaveValue(testProps.initialValues.oktaFirstName);

    const oktaLastName = await screen.findByLabelText('Last Name *');
    expect(oktaLastName).toBeInstanceOf(HTMLInputElement);
    expect(oktaLastName).toHaveValue(testProps.initialValues.oktaLastName);

    const oktaEdipi = await screen.findByLabelText('DoD ID number *');
    expect(oktaEdipi).toHaveValue(testProps.initialValues.oktaEdipi);
    expect(oktaEdipi).toBeDisabled();
  });

  it('shows an error message if Okta Email is not in email format', async () => {
    renderWithRouter(<EditOktaInfoForm {...testProps} />);

    const emailInput = await screen.findByLabelText('Okta Email *');
    await userEvent.clear(emailInput);
    await userEvent.type(emailInput, 'oktaUserWithoutEmail');
    await userEvent.tab();

    const alert = await screen.findByRole('alert');
    expect(alert).toBeInTheDocument();
    expect(alert).toHaveTextContent('Email address must end in a valid domain');
  });

  it('shows an error message if Okta Email is empty', async () => {
    renderWithRouter(<EditOktaInfoForm {...testProps} />);

    const saveButton = await screen.findByRole('button', { name: 'Save' });
    expect(saveButton).toBeEnabled();

    const emailInput = await screen.findByLabelText('Okta Email *');
    await userEvent.clear(emailInput);
    await userEvent.tab();

    const alert = await screen.findByRole('alert');
    expect(alert).toBeInTheDocument();
    expect(alert).toHaveTextContent('Required');

    expect(saveButton).toBeDisabled();
  });

  it('submits the form when its valid', async () => {
    renderWithRouter(<EditOktaInfoForm {...testProps} />);

    const saveButton = screen.getByRole('button', { name: 'Save' });

    await userEvent.click(saveButton);

    const expectedParams = {
      ...testProps.initialValues,
    };

    await waitFor(() => {
      expect(testProps.onSubmit).toHaveBeenCalledWith(expectedParams, expect.anything());
    });
  });

  it('converts the "Submit" button into the "Return to Home" button when the move has been locked by an office user', async () => {
    renderWithRouter(<EditOktaInfoForm {...testPropsWithLock} />);
    expect(screen.getByRole('button', { name: 'Return home' })).toBeInTheDocument();
  });

  it('implements the onCancel handler when the Cancel button is clicked', async () => {
    renderWithRouter(<EditOktaInfoForm {...testProps} />);

    const cancelButton = screen.getByRole('button', { name: 'Cancel' });

    await userEvent.click(cancelButton);

    await waitFor(() => {
      expect(testProps.onCancel).toHaveBeenCalled();
    });
  });

  afterEach(jest.resetAllMocks);
});
