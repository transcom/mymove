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
      expect(getByLabelText('Best contact phone')).toBeInstanceOf(HTMLInputElement);
      expect(getByLabelText('Best contact phone')).toBeRequired();
      expect(getByLabelText(/Alt. phone/)).toBeInstanceOf(HTMLInputElement);
      expect(getByLabelText(/Alt. phone/)).not.toBeRequired();
      expect(getByLabelText('Personal email')).toBeInstanceOf(HTMLInputElement);
      expect(getByLabelText('Personal email')).toBeRequired();
      expect(getByLabelText('Phone')).toBeInstanceOf(HTMLInputElement);
      expect(getByLabelText('Email')).toBeInstanceOf(HTMLInputElement);
    });
  });

  it('validates the contact phone field', async () => {
    const { getByText, getByLabelText } = render(<ContactInfoForm {...testProps} />);
    userEvent.type(getByLabelText('Best contact phone'), '12345');
    userEvent.tab();

    await waitFor(() => {
      expect(getByText('Number must have 10 digits and a valid area code')).toBeInTheDocument();
    });
  });

  it('validates the alt phone field', async () => {
    const { getByText, getByLabelText } = render(<ContactInfoForm {...testProps} />);
    userEvent.type(getByLabelText(/Alt. phone/), '543');
    userEvent.tab();

    await waitFor(() => {
      expect(getByText('Number must have 10 digits and a valid area code')).toBeInTheDocument();
    });
  });

  it('validates the email field', async () => {
    const { getByText, getByLabelText } = render(<ContactInfoForm {...testProps} />);
    userEvent.type(getByLabelText('Personal email'), 'sample@');
    userEvent.tab();

    await waitFor(() => {
      expect(getByText('Must be a valid email address')).toBeInTheDocument();
    });
  });

  it('shows an error message when trying to submit an invalid form', async () => {
    const { getAllByText, getByRole } = render(<ContactInfoForm {...testProps} />);
    const submitBtn = getByRole('button', { name: 'Next' });

    userEvent.click(submitBtn);

    await waitFor(() => {
      expect(getAllByText('Required').length).toBe(2);
    });

    expect(testProps.onSubmit).not.toHaveBeenCalled();
  });

  it('is invalid if neither email nor phone is preferred is checked', async () => {
    const { getByRole, getByLabelText } = render(<ContactInfoForm {...testProps} />);
    const submitBtn = getByRole('button', { name: 'Next' });

    userEvent.type(getByLabelText('Best contact phone'), '555-555-5555');
    userEvent.type(getByLabelText('Personal email'), 'test@sample.com');
    userEvent.click(submitBtn);

    await waitFor(() => {
      expect(getByRole('button', { name: 'Next' })).toBeDisabled();
    });

    expect(testProps.onSubmit).not.toHaveBeenCalled();
  });

  it('submits a form when it is valid', async () => {
    const { getByRole, getByLabelText } = render(<ContactInfoForm {...testProps} />);
    const submitBtn = getByRole('button', { name: 'Next' });

    userEvent.type(getByLabelText('Best contact phone'), '555-555-5555');
    userEvent.type(getByLabelText('Personal email'), 'test@sample.com');
    userEvent.click(getByLabelText('Email'));
    userEvent.click(submitBtn);

    await waitFor(() => {
      expect(testProps.onSubmit).toHaveBeenCalled();
    });
  });

  it('calls the back handler when back button is clicked', async () => {
    const { getByRole, getByLabelText } = render(<ContactInfoForm {...testProps} />);
    const backBtn = getByRole('button', { name: 'Back' });

    userEvent.type(getByLabelText('Best contact phone'), '555-555-1111');
    userEvent.type(getByLabelText('Personal email'), 'test@sample.com');
    userEvent.click(getByLabelText('Email'));
    userEvent.click(backBtn);

    await waitFor(() => {
      expect(testProps.onBack).toHaveBeenCalled();
    });
  });
});
