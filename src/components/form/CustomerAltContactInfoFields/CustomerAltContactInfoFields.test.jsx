import React from 'react';
import { render, screen } from '@testing-library/react';
import { Formik } from 'formik';

import { CustomerAltContactInfoFields } from './index';

describe('ContactInfoFields component', () => {
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
        first_name: 'Leo',
        middle_name: 'Star',
        last_name: 'Spaceman',
        suffix: 'Mr.',
        customer_telephone: '555-555-5555',
        customer_email: 'test@sample.com',
      };

      const { getByLabelText } = render(
        <Formik initialValues={initialValues}>
          <CustomerAltContactInfoFields legend="Contact info" name="contact" />
        </Formik>,
      );
      expect(getByLabelText('First name')).toHaveValue(initialValues.first_name);
      expect(getByLabelText(/Middle name/)).toHaveValue(initialValues.middle_name);
      expect(getByLabelText('Last name')).toHaveValue(initialValues.last_name);
      expect(getByLabelText(/Suffix/)).toHaveValue(initialValues.suffix);
      expect(getByLabelText('Phone')).toHaveValue(initialValues.customer_telephone);
      expect(getByLabelText('Email')).toHaveValue(initialValues.customer_email);
    });
  });
});
