import React from 'react';
import { render, waitFor, screen } from '@testing-library/react';
// import userEvent from '@testing-library/user-event';

import CustomerContactInfoForm from './CustomerContactInfoForm';

describe('CustomerContactInfoForm Component', () => {
  const initialValues = {
    first_name: '',
    middle_name: '',
    last_name: '',
    suffix: '',
    customer_telephone: '',
    customer_email: '',
    customer_address: {
      street_address_1: '',
      street_address_2: '',
      city: '',
      state: '',
      postal_code: '',
    },
    name: '',
    telephone: '',
    email: '',
  };
  const testProps = {
    initialValues,
    onSubmit: jest.fn(),
    onBack: jest.fn(),
  };

  it('renders the form inputs', async () => {
    render(<CustomerContactInfoForm {...testProps} />);

    await waitFor(() => {
      expect(screen.getByText('Contact info')).toBeInstanceOf(HTMLHeadingElement);
      expect(screen.getByLabelText('First name')).toBeInstanceOf(HTMLInputElement);
      expect(screen.getByLabelText('First name')).toBeRequired();

      expect(screen.getByLabelText(/Middle name/)).toBeInstanceOf(HTMLInputElement);

      expect(screen.getByLabelText('Last name')).toBeInstanceOf(HTMLInputElement);
      expect(screen.getByLabelText('Last name')).toBeRequired();

      expect(screen.getByLabelText(/Suffix/)).toBeInstanceOf(HTMLInputElement);

      expect(screen.getAllByLabelText('Phone')[0]).toBeInstanceOf(HTMLInputElement);
      expect(screen.getAllByLabelText('Phone')[0]).toBeRequired();
      expect(screen.getAllByLabelText('Email')[0]).toBeInstanceOf(HTMLInputElement);
      expect(screen.getAllByLabelText('Email')[0]).toBeRequired();

      expect(screen.getByText('Current Address')).toBeInstanceOf(HTMLHeadingElement);
      expect(screen.getByLabelText('Address 1')).toBeInstanceOf(HTMLInputElement);
      expect(screen.getByLabelText(/Address 2/)).toBeInstanceOf(HTMLInputElement);
      expect(screen.getByLabelText('City')).toBeInstanceOf(HTMLInputElement);
      expect(screen.getByLabelText('State')).toBeInstanceOf(HTMLSelectElement);
      expect(screen.getByLabelText('ZIP')).toBeInstanceOf(HTMLInputElement);

      expect(screen.getByLabelText('Name')).toBeInstanceOf(HTMLInputElement);
      expect(screen.getByLabelText('Name')).toBeRequired();
      expect(screen.getAllByLabelText('Phone')[1]).toBeInstanceOf(HTMLInputElement);
      expect(screen.getAllByLabelText('Phone')[1]).toBeRequired();
      expect(screen.getAllByLabelText('Email')[1]).toBeInstanceOf(HTMLInputElement);
      expect(screen.getAllByLabelText('Email')[1]).toBeRequired();
    });
  });

  // it('validates the customer phone field', async () => {
  //   render(<CustomerContactInfoForm {...testProps} />);
  //   userEvent.type(screen.getByLabelText('Best contact phone'), '12345');
  //   userEvent.tab();

  //   await waitFor(() => {
  //     expect(screen.getByText('Number must have 10 digits and a valid area code')).toBeInTheDocument();
  //   });
  // });

  // it('validates the customer email field', async () => {
  //   render(<CustomerContactInfoForm {...testProps} />);
  //   userEvent.type(screen.getByLabelText('Email'), 'sample@');
  //   userEvent.tab();

  //   await waitFor(() => {
  //     expect(screen.getByText('Must be a valid email address')).toBeInTheDocument();
  //   });
  // });

  // it('shows an error message when trying to submit an invalid form', async () => {
  //   render(<CustomerContactInfoForm {...testProps} />);
  //   const submitBtn = screen.getByRole('button', { name: 'Save' });

  //   userEvent.click(submitBtn);

  //   await waitFor(() => {
  //     expect(screen.getAllByText('Required').length).toBe(2);
  //   });

  //   expect(testProps.onSubmit).not.toHaveBeenCalled();
  // });

  // it('submits a form when it is valid', async () => {
  //   render(<CustomerContactInfoForm {...testProps} />);
  //   const submitBtn = screen.getByRole('button', { name: 'Next' });

  //   userEvent.type(screen.getByLabelText('Phone'), '555-555-5555');
  //   userEvent.type(screen.getByLabelText('Email'), 'test@sample.com');
  //   userEvent.click(screen.getByLabelText('Email'));
  //   userEvent.click(submitBtn);

  //   await waitFor(() => {
  //     expect(testProps.onSubmit).toHaveBeenCalled();
  //   });
  // });

  // it('calls the cancel handler when cancel button is clicked', async () => {
  //   render(<CustomerContactInfoForm {...testProps} />);
  //   const backBtn = screen.getByRole('button', { name: 'Cancel' });

  //   userEvent.type(screen.getByLabelText('Phone'), '555-555-1111');
  //   userEvent.type(screen.getByLabelText('Email'), 'test@sample.com');
  //   userEvent.click(screen.getByLabelText('Email'));
  //   userEvent.click(backBtn);

  //   await waitFor(() => {
  //     expect(testProps.onBack).toHaveBeenCalled();
  //   });
  // });
});
