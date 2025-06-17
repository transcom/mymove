import React from 'react';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import CreateAccount from './CreateAccount';

import { MockProviders } from 'testUtils';
import { isBooleanFlagEnabledUnauthenticated } from 'utils/featureFlags';
import { validateCode } from 'services/internalApi';

const dummySetShowLoadingSpinner = jest.fn();

const mockNavigate = jest.fn();
jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom'),
  useNavigate: () => mockNavigate,
}));

jest.mock('utils/featureFlags', () => ({
  ...jest.requireActual('utils/featureFlags'),
  isBooleanFlagEnabledUnauthenticated: jest.fn().mockImplementation(() => Promise.resolve(false)),
}));

jest.mock('services/internalApi', () => ({
  ...jest.requireActual('services/internalApi'),
  validateCode: jest.fn().mockImplementation(() => Promise.resolve()),
}));

describe('CreateAccount Component', () => {
  const renderComponent = () =>
    render(
      <MockProviders>
        <CreateAccount setShowLoadingSpinner={dummySetShowLoadingSpinner} />
      </MockProviders>,
    );

  it('renders the form with expected fields and asterisks for required fields', async () => {
    renderComponent();
    await waitFor(() => {
      expect(screen.getByTestId('modal')).toBeInTheDocument();
    });
    fireEvent.click(screen.getByText('Yes'));
    expect(document.querySelector('#reqAsteriskMsg')).toHaveTextContent('Fields marked with * are required.');
    expect(screen.getByTestId('affiliationInput')).toBeInTheDocument();
    expect(screen.getByTestId('edipiInput')).toBeInTheDocument();
    expect(screen.getByTestId('edipiConfirmationInput')).toBeInTheDocument();
    expect(screen.getByTestId('firstName')).toBeInTheDocument();
    expect(screen.getByTestId('middleInitial')).toBeInTheDocument();
    expect(screen.getByTestId('lastName')).toBeInTheDocument();
    expect(screen.getByTestId('email')).toBeInTheDocument();
    expect(screen.getByTestId('emailConfirmation')).toBeInTheDocument();
    expect(screen.getByTestId('telephone')).toBeInTheDocument();
    expect(screen.getByTestId('secondaryTelephone')).toBeInTheDocument();
    expect(screen.getByTestId('phoneIsPreferred')).toBeInTheDocument();
    expect(screen.getByTestId('emailIsPreferred')).toBeInTheDocument();
    expect(screen.getByTestId('submitBtn')).toBeInTheDocument();
    expect(screen.getByTestId('submitBtn')).toBeDisabled();
    expect(screen.getByTestId('cancelBtn')).toBeInTheDocument();
    expect(screen.getByTestId('cancelBtn')).toBeEnabled();

    expect(screen.getByLabelText('Branch of service *')).toBeRequired();
    expect(screen.getByLabelText('DoD ID number *')).toBeRequired();
    expect(screen.getByLabelText('Confirm DoD ID number *')).toBeRequired();

    expect(screen.getByLabelText('First Name *')).toBeRequired();
    expect(screen.getByLabelText('Last Name *')).toBeRequired();
    expect(screen.getByLabelText('Email *')).toBeRequired();
    expect(screen.getByLabelText('Telephone *')).toBeRequired();

    await userEvent.selectOptions(screen.getByLabelText(/Branch of service/i), ['COAST_GUARD']);
    expect(screen.getByLabelText('EMPLID *')).toBeRequired();
    expect(screen.getByLabelText('Confirm EMPLID *')).toBeRequired();
  });

  it('shows the ValidCACModal on load', async () => {
    renderComponent();
    await waitFor(() => {
      expect(screen.getByTestId('modal')).toBeInTheDocument();
    });
  });

  it('calls navigate to /sign-in when user does not have valid CAC', async () => {
    renderComponent();
    await waitFor(() => {
      expect(screen.getByTestId('modal')).toBeInTheDocument();
    });
    fireEvent.click(screen.getByText('No'));
    expect(mockNavigate).toHaveBeenCalledWith('/sign-in', { state: { noValidCAC: true } });
  });

  it('Submit buttons stays disabled until form is validated', async () => {
    renderComponent();
    await waitFor(() => {
      expect(screen.getByTestId('modal')).toBeInTheDocument();
    });
    fireEvent.click(screen.getByText('Yes'));
    expect(screen.getByTestId('submitBtn')).toBeDisabled();
    await userEvent.selectOptions(screen.getByLabelText(/Branch of service/i), ['NAVY']);
    await userEvent.type(screen.getByTestId('edipiInput'), '1234567890');
    await userEvent.type(screen.getByTestId('edipiConfirmationInput'), '1234567890');
    await userEvent.type(screen.getByTestId('firstName'), 'Jim');
    await userEvent.type(screen.getByTestId('lastName'), 'Bob');
    await userEvent.type(screen.getByTestId('email'), 'jim@jim.com');
    await userEvent.type(screen.getByTestId('emailConfirmation'), 'jim@jim.com');
    await userEvent.type(screen.getByTestId('telephone'), '555-555-5555');

    expect(screen.getByTestId('submitBtn')).toBeEnabled();
  });

  it('Validations display when confirm fields do not match', async () => {
    renderComponent();
    await waitFor(() => {
      expect(screen.getByTestId('modal')).toBeInTheDocument();
    });
    fireEvent.click(screen.getByText('Yes'));
    expect(screen.getByTestId('submitBtn')).toBeDisabled();
    await userEvent.selectOptions(screen.getByLabelText(/Branch of service/i), ['COAST_GUARD']);
    await userEvent.type(screen.getByTestId('edipiInput'), '1234567890');
    await userEvent.type(screen.getByTestId('edipiConfirmationInput'), '1234567899');
    await userEvent.type(screen.getByTestId('emplidInput'), '123456');
    await userEvent.type(screen.getByTestId('emplidConfirmationInput'), '123455');
    await userEvent.type(screen.getByTestId('firstName'), 'Jim');
    await userEvent.type(screen.getByTestId('lastName'), 'Bob');
    await userEvent.type(screen.getByTestId('email'), 'jim@jim.com');
    await userEvent.type(screen.getByTestId('emailConfirmation'), 'jam@jim.com');
    await userEvent.type(screen.getByTestId('telephone'), '555-555-5555');

    expect(screen.getByText('DoD ID numbers must match')).toBeInTheDocument();
    expect(screen.getByText('EMPLID numbers must match')).toBeInTheDocument();
    expect(screen.getByText('Emails must match')).toBeInTheDocument();
  });

  it('renders the validation code screen when FF is on, able to submit and create account', async () => {
    isBooleanFlagEnabledUnauthenticated.mockImplementation(() => Promise.resolve(true));

    validateCode.mockImplementation(() =>
      Promise.resolve({
        body: {
          parameterValue: 'TestCode123123',
          parameterName: 'validation_code',
        },
      }),
    );

    renderComponent();

    await waitFor(() => {
      expect(screen.getByText('Please enter a validation code to begin creating a move')).toBeInTheDocument();
    });

    // submit the validation code
    const nextBtn = await screen.findByRole('button', { name: 'Next' });
    expect(nextBtn).toBeDisabled();
    await userEvent.type(screen.getByLabelText('Validation code *'), 'TestCode123123');
    expect(nextBtn).toBeEnabled();
    await userEvent.click(nextBtn);

    // should now see the CAC validated modal
    await waitFor(() => {
      expect(screen.getByTestId('modal')).toBeInTheDocument();
    });
  });
});
