import React from 'react';
import { render, screen } from '@testing-library/react';
import { Formik } from 'formik';

import { CustomerAltContactInfoFields } from './index';

describe('CustomerAltContactInfoFields component', () => {
  it('renders a legend and all service member contact info inputs', () => {
    render(
      <Formik>
        <CustomerAltContactInfoFields legend="Contact info" />
      </Formik>,
    );
    expect(screen.getByText('Contact info')).toBeInstanceOf(HTMLLegendElement);
    expect(screen.getByLabelText('First name')).toBeInstanceOf(HTMLInputElement);
    expect(screen.getByLabelText('First name')).toBeRequired();

    expect(screen.getByLabelText(/Middle name/)).toBeInstanceOf(HTMLInputElement);

    expect(screen.getByLabelText('Last name')).toBeInstanceOf(HTMLInputElement);
    expect(screen.getByLabelText('Last name')).toBeRequired();

    expect(screen.getByLabelText(/Suffix/)).toBeInstanceOf(HTMLInputElement);

    expect(screen.getByLabelText('Phone')).toBeInstanceOf(HTMLInputElement);
    expect(screen.getByLabelText('Email')).toBeInstanceOf(HTMLInputElement);
  });

  describe('with pre-filled values', () => {
    it('renders a legend and all service member contact info inputs', () => {
      const initialValues = {
        firstName: 'Leo',
        middleName: 'Star',
        lastName: 'Spaceman',
        suffix: 'Mr.',
        customerTelephone: '555-555-5555',
        customerEmail: 'test@sample.com',
      };

      const { getByLabelText } = render(
        <Formik initialValues={initialValues}>
          <CustomerAltContactInfoFields legend="Contact info" name="contact" />
        </Formik>,
      );
      expect(getByLabelText('First name')).toHaveValue(initialValues.firstName);
      expect(getByLabelText(/Middle name/)).toHaveValue(initialValues.middleName);
      expect(getByLabelText('Last name')).toHaveValue(initialValues.lastName);
      expect(getByLabelText(/Suffix/)).toHaveValue(initialValues.suffix);
      expect(getByLabelText('Phone')).toHaveValue(initialValues.customerTelephone);
      expect(getByLabelText('Email')).toHaveValue(initialValues.customerEmail);
    });
  });
});
