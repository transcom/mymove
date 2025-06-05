import React from 'react';
import { render, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import BackupContactForm from './index';

describe('BackupContactForm Component', () => {
  const initialValues = {
    firstName: '',
    lastName: '',
    telephone: '',
    email: '',
  };
  const testProps = {
    initialValues,
    onSubmit: jest.fn(),
    onBack: jest.fn(),
  };

  it('renders the form inputs', async () => {
    const { getByLabelText } = render(<BackupContactForm {...testProps} />);

    await waitFor(() => {
      expect(getByLabelText(/First Name/)).toBeInstanceOf(HTMLInputElement);
      expect(getByLabelText(/First Name/)).toBeRequired();
      expect(getByLabelText(/Last Name/)).toBeInstanceOf(HTMLInputElement);
      expect(getByLabelText(/Last Name/)).toBeRequired();
      expect(getByLabelText(/Phone/)).toBeInstanceOf(HTMLInputElement);
      expect(getByLabelText(/Phone/)).toBeRequired();
      expect(getByLabelText(/Email/)).toBeInstanceOf(HTMLInputElement);
      expect(getByLabelText(/Email/)).toBeRequired();
    });
  });

  it('validates the contact phone field', async () => {
    const { getByText, getByLabelText } = render(<BackupContactForm {...testProps} />);
    await userEvent.type(getByLabelText(/Phone/), '12345');
    await userEvent.tab();

    await waitFor(() => {
      expect(
        getByText('Please enter a valid phone number. Phone numbers must be entered as ###-###-####.'),
      ).toBeInTheDocument();
    });
  });

  it('validates the email field', async () => {
    const { getByText, getByLabelText } = render(<BackupContactForm {...testProps} />);
    await userEvent.type(getByLabelText(/Email/), 'sample@');
    await userEvent.tab();

    await waitFor(() => {
      expect(getByText('Must be a valid email address')).toBeInTheDocument();
    });
  });

  it('shows an error message when trying to submit an invalid form', async () => {
    const { getAllByTestId, getByRole, getByLabelText } = render(<BackupContactForm {...testProps} />);
    const submitBtn = getByRole('button', { name: 'Next' });

    // Touch all of the required fields so that they show error messages
    await userEvent.click(getByLabelText(/First Name/));
    await userEvent.click(getByLabelText(/Last Name/));
    await userEvent.click(getByLabelText(/Phone/));
    await userEvent.click(getByLabelText(/Email/));
    await userEvent.click(submitBtn);

    const numberOfContactFields = 4;

    await waitFor(() => {
      expect(getAllByTestId('errorMessage').length).toBe(numberOfContactFields);
    });

    expect(testProps.onSubmit).not.toHaveBeenCalled();
  });

  it('submits a form when it is valid', async () => {
    const { getByRole, getByLabelText } = render(<BackupContactForm {...testProps} />);
    const submitBtn = getByRole('button', { name: 'Next' });

    await userEvent.type(getByLabelText(/First Name/), 'Joe');
    await userEvent.type(getByLabelText(/Last Name/), 'Schmoe');
    await userEvent.type(getByLabelText(/Phone/), '555-555-5555');
    await userEvent.type(getByLabelText(/Email/), 'test@sample.com');
    await userEvent.click(submitBtn);

    await waitFor(() => {
      expect(testProps.onSubmit).toHaveBeenCalled();
    });
  });

  it('calls the back handler when back button is clicked', async () => {
    const { getByRole, getByLabelText } = render(<BackupContactForm {...testProps} />);
    const backBtn = getByRole('button', { name: 'Back' });

    await userEvent.type(getByLabelText(/First Name/), 'Janey');
    await userEvent.type(getByLabelText(/Last Name/), 'Profaney');
    await userEvent.type(getByLabelText(/Phone/), '555-555-1111');
    await userEvent.click(getByLabelText(/Email/));
    await userEvent.click(backBtn);

    await waitFor(() => {
      expect(testProps.onBack).toHaveBeenCalled();
    });
  });
});
