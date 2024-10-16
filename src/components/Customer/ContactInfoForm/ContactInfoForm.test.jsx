import React from 'react';
import { render, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import ContactInfoForm from './index';

describe('ContactInfoForm Component', () => {
  const initialValues = {
    telephone: '',
    secondary_telephone: '',
    personal_email: '',
    phone_is_preferred: false,
    email_is_preferred: false,
  };
  const testProps = {
    initialValues,
    onSubmit: jest.fn(),
    onBack: jest.fn(),
  };

  it('renders the form inputs', async () => {
    const { getByLabelText } = render(<ContactInfoForm {...testProps} />);

    await waitFor(() => {
      expect(getByLabelText(/Best contact phone/)).toBeInstanceOf(HTMLInputElement);
      expect(getByLabelText(/Best contact phone/)).toBeRequired();
      expect(getByLabelText(/Alt. phone/)).toBeInstanceOf(HTMLInputElement);
      expect(getByLabelText(/Alt. phone/)).not.toBeRequired();
      expect(getByLabelText(/Personal email/)).toBeInstanceOf(HTMLInputElement);
      expect(getByLabelText(/Personal email/)).toBeRequired();
      expect(getByLabelText(/Phone/)).toBeInstanceOf(HTMLInputElement);
      expect(getByLabelText(/Email/)).toBeInstanceOf(HTMLInputElement);
    });
  });

  it('validates the contact phone field', async () => {
    const { getByText, getByLabelText } = render(<ContactInfoForm {...testProps} />);
    await userEvent.type(getByLabelText(/Best contact phone/), '12345');
    await userEvent.tab();

    await waitFor(() => {
      expect(
        getByText('Please enter a valid phone number. Phone numbers must be entered as ###-###-####.'),
      ).toBeInTheDocument();
    });
  });

  it('validates the alt phone field', async () => {
    const { getByText, getByLabelText } = render(<ContactInfoForm {...testProps} />);
    await userEvent.type(getByLabelText(/Alt. phone/), '543');
    await userEvent.tab();

    await waitFor(() => {
      expect(
        getByText('Please enter a valid phone number. Phone numbers must be entered as ###-###-####.'),
      ).toBeInTheDocument();
    });
  });

  it('validates the email field', async () => {
    const { getByText, getByLabelText } = render(<ContactInfoForm {...testProps} />);
    await userEvent.type(getByLabelText(/Personal email/), 'sample@');
    await userEvent.tab();

    await waitFor(() => {
      expect(getByText('Must be a valid email address')).toBeInTheDocument();
    });
  });

  it('shows an error message when trying to submit an invalid form', async () => {
    const { getAllByTestId, getByRole, getByLabelText } = render(<ContactInfoForm {...testProps} />);
    // Touch required fields to show validation errors
    await userEvent.click(getByLabelText(/Best contact phone/));
    await userEvent.click(getByLabelText(/Personal email/));

    const submitBtn = getByRole('button', { name: 'Next' });
    await userEvent.click(submitBtn);

    await waitFor(() => {
      expect(getAllByTestId('errorMessage').length).toBe(2);
    });

    expect(testProps.onSubmit).not.toHaveBeenCalled();
  });

  it('is invalid if neither email nor phone is preferred is checked', async () => {
    const { getByRole, getByLabelText } = render(<ContactInfoForm {...testProps} />);
    const submitBtn = getByRole('button', { name: 'Next' });

    await userEvent.type(getByLabelText(/Best contact phone/), '555-555-5555');
    await userEvent.type(getByLabelText(/Personal email/), 'test@sample.com');
    await userEvent.click(submitBtn);

    await waitFor(() => {
      expect(getByRole('button', { name: 'Next' })).toBeDisabled();
    });

    expect(testProps.onSubmit).not.toHaveBeenCalled();
  });

  it('submits a form when it is valid', async () => {
    const { getByRole, getByLabelText } = render(<ContactInfoForm {...testProps} />);
    const submitBtn = getByRole('button', { name: 'Next' });

    await userEvent.type(getByLabelText(/Best contact phone/), '555-555-5555');
    await userEvent.type(getByLabelText(/Personal email/), 'test@sample.com');
    await userEvent.click(getByLabelText(/Email/));
    await userEvent.click(submitBtn);

    await waitFor(() => {
      expect(testProps.onSubmit).toHaveBeenCalled();
    });
  });

  it('calls the back handler when back button is clicked', async () => {
    const { getByRole, getByLabelText } = render(<ContactInfoForm {...testProps} />);
    const backBtn = getByRole('button', { name: 'Back' });

    await userEvent.type(getByLabelText(/Best contact phone/), '555-555-1111');
    await userEvent.type(getByLabelText(/Personal email/), 'test@sample.com');
    await userEvent.click(getByLabelText(/Email/));
    await userEvent.click(backBtn);

    await waitFor(() => {
      expect(testProps.onBack).toHaveBeenCalled();
    });
  });
});
