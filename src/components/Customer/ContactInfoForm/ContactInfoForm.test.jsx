/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import { render, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import ContactInfoForm from './index';

describe('ContactInfoForm Component', () => {
  const initialValues = {
    telephone: '555-555-5555',
    secondary_phone: '555-444-5555',
    personal_email: 'test@sample.com',
    phone_is_preferred: false,
    email_is_preferred: true,
  };
  const testProps = {
    initialValues,
    onSubmit: jest.fn,
    onBack: jest.fn,
  };

  //  renders the input fields
  it('renders the form inputs', async () => {
    const { getByTestId } = render(<ContactInfoForm {...testProps} />);

    await waitFor(() => {
      expect(getByTestId('contactPhone')).toBeInstanceOf(HTMLInputElement);
      expect(getByTestId('contactPhone')).toBeRequired();
      expect(getByTestId('secondaryPhone')).toBeInstanceOf(HTMLInputElement);
      expect(getByTestId('secondaryPhone')).not.toBeRequired();
      expect(getByTestId('personalEmail')).toBeInstanceOf(HTMLInputElement);
      expect(getByTestId('personalEmail')).toBeRequired();
      expect(getByTestId('contactInfoPhonePreferred')).toBeInstanceOf(HTMLInputElement);
      expect(getByTestId('contactInfoEmailPreferred')).toBeInstanceOf(HTMLInputElement);
    });
  });
  //  validates the phone fields
  it('validates the contact phone field', async () => {
    const { getByText, getByTestId } = render(<ContactInfoForm {...testProps} />);
    userEvent.type(getByTestId('contactPhone'), '12345');
    userEvent.tab();

    await waitFor(() => {
      expect(getByTestId('contactPhone')).toBeInvalid();
      expect(getByText('Number must have 10 digits and a valid area code')).toBeInTheDocument();
    });
  });
  it('validates the alt phone field', async () => {
    const { getByText, getByTestId } = render(<ContactInfoForm {...testProps} />);
    userEvent.type(getByTestId('secondaryPhone'), '543');
    userEvent.tab();

    await waitFor(() => {
      expect(getByTestId('secondaryPhone')).toBeInvalid();
      expect(getByText('Number must have 10 digits and a valid area code')).toBeInTheDocument();
    });
  });
  //  validates the email field
  it('validates the email field', async () => {
    const { getByText, getByTestId } = render(<ContactInfoForm {...testProps} />);
    userEvent.type(getByTestId('personalEmail'), 'sample@');
    userEvent.tab();

    await waitFor(() => {
      expect(getByTestId('personalEmail')).toBeInvalid();
      expect(getByText('Must be a valid email address')).toBeInTheDocument();
    });
  });
  //  shows an error message when trying to submit an invalid form
  it('shows an error message when trying to submit an invalid form', async () => {
    const { getAllByText, getByRole } = render(<ContactInfoForm {...testProps} />);
    const submitBtn = getByRole('button', { name: 'Next' });

    userEvent.click(submitBtn);

    await waitFor(() => {
      expect(getAllByText('Required').length).toBe(3);
    });

    expect(testProps.onSubmit).not.toHaveBeenCalled();
  });

  it('is invalid if neither email nor phone is preferred is checked', async () => {
    const { getByText, getByRole, getByDataId } = render(<ContactInfoForm {...testProps} />);
    const submitBtn = getByRole('button', { name: 'Next' });

    userEvent.type(getByDataId('contactPhone'), '555-555-5555');
    userEvent.type(getByDataId('personalEmail'), 'test@sample.com');
    userEvent.click(submitBtn);

    await waitFor(() => {
      expect(getByText('msg to check one or the other')).toBeInTheDocument();
    });

    expect(testProps.onSubmit).not.toHaveBeenCalled();
  });

  //  submits a form when it is valid
  it('submits a form when it is valid', async () => {
    const { getAllByText, getByRole, getByDataId } = render(<ContactInfoForm {...testProps} />);
    const submitBtn = getByRole('button', { name: 'Next' });

    userEvent.type(getByDataId('contactPhone'), '555-555-5555');
    userEvent.type(getByDataId('personalEmail'), 'test@sample.com');
    userEvent.click(getByDataId('contactInfoEmailPreferred'));
    userEvent.click(submitBtn);

    await waitFor(() => {
      expect(getAllByText('Required').length).toBe(3);
    });

    expect(testProps.onSubmit).not.toHaveBeenCalled();
  });
});
