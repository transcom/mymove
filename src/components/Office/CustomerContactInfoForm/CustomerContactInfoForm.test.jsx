import React from 'react';
import { render, waitFor, screen } from '@testing-library/react';
// import userEvent from '@testing-library/user-event';

import CustomerContactInfoForm from './CustomerContactInfoForm';

describe('CustomerContactInfoForm Component', () => {
  const initialValues = {
    firstName: '',
    middleName: '',
    lastName: '',
    suffix: '',
    customerTelephone: '',
    customerEmail: '',
    customerAddress: {
      streetAddress1: '123 Happy St',
      streetAddress2: 'Unit 4',
      city: 'Missoula',
      state: 'MT',
      postalCode: '59802',
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
      expect(screen.getByLabelText('This is not the person named on the orders.')).toBeInstanceOf(HTMLInputElement);
      expect(screen.getByLabelText('First name')).toBeInstanceOf(HTMLInputElement);
      expect(screen.getByLabelText('First name')).toBeRequired();

      expect(screen.getByLabelText(/Middle name/)).toBeInstanceOf(HTMLInputElement);

      expect(screen.getByLabelText('Last name')).toBeInstanceOf(HTMLInputElement);
      expect(screen.getByLabelText('Last name')).toBeRequired();

      expect(screen.getByLabelText(/Suffix/)).toBeInstanceOf(HTMLInputElement);

      expect(screen.getAllByLabelText('Phone')[0]).toBeInstanceOf(HTMLInputElement);
      expect(screen.getAllByLabelText('Phone')[0]).toBeRequired();

      expect(screen.getAllByLabelText(/Alternate Phone/)[0]).toBeInstanceOf(HTMLInputElement);

      expect(screen.getAllByLabelText('Email')[0]).toBeInstanceOf(HTMLInputElement);
      expect(screen.getAllByLabelText('Email')[0]).toBeRequired();

      expect(screen.getByText('Current Address')).toBeInstanceOf(HTMLHeadingElement);
      expect(screen.getByDisplayValue('123 Happy St')).toBeInstanceOf(HTMLInputElement);
      expect(screen.getByDisplayValue('Unit 4')).toBeInstanceOf(HTMLInputElement);
      expect(screen.getByDisplayValue('Missoula')).toBeInstanceOf(HTMLInputElement);
      expect(screen.getByDisplayValue('MT')).toBeInstanceOf(HTMLSelectElement);
      expect(screen.getByDisplayValue('59802')).toBeInstanceOf(HTMLInputElement);

      expect(screen.getByLabelText('Name')).toBeInstanceOf(HTMLInputElement);
      expect(screen.getAllByLabelText('Phone')[1]).toBeInstanceOf(HTMLInputElement);
      expect(screen.getAllByLabelText('Email')[1]).toBeInstanceOf(HTMLInputElement);
    });
  });
});
