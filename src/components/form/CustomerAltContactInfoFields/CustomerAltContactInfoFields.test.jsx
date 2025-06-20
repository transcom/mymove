import React from 'react';
import { render, screen } from '@testing-library/react';
import { Formik } from 'formik';

import { CustomerAltContactInfoFields } from './index';

describe('CustomerAltContactInfoFields component', () => {
  it('renders a legend and all service member contact info inputs and asterisks for required fields', () => {
    render(
      <Formik>
        <CustomerAltContactInfoFields legend="Contact info" />
      </Formik>,
    );

    expect(document.querySelector('#reqAsteriskMsg')).toHaveTextContent('Fields marked with * are required.');

    expect(screen.getByText('Contact info')).toBeInstanceOf(HTMLLegendElement);
    expect(screen.getByLabelText('First name *')).toBeInstanceOf(HTMLInputElement);
    expect(screen.getByLabelText('First name *')).toBeRequired();

    expect(screen.getByLabelText(/Middle name/)).toBeInstanceOf(HTMLInputElement);

    expect(screen.getByLabelText('Last name *')).toBeInstanceOf(HTMLInputElement);
    expect(screen.getByLabelText('Last name *')).toBeRequired();

    expect(screen.getByLabelText(/Suffix/)).toBeInstanceOf(HTMLInputElement);

    expect(screen.getAllByText('Phone')).toBeInstanceOf(Array);
    expect(screen.getAllByText('Email')).toBeInstanceOf(Array);
    expect(screen.getByLabelText(/Alternate Phone/)).toBeInstanceOf(HTMLInputElement);
    expect(screen.getByText('Preferred contact method')).toBeInstanceOf(HTMLLabelElement);
  });

  describe('with pre-filled values', () => {
    it('renders a legend and all service member contact info inputs', async () => {
      const initialValues = {
        firstName: 'Leo',
        middleName: 'Star',
        lastName: 'Spaceman',
        suffix: 'Mr.',
        customerTelephone: '555-555-5555',
        customerEmail: 'test@sample.com',
        secondaryPhone: '233-444-2222',
      };

      render(
        <Formik initialValues={initialValues}>
          <CustomerAltContactInfoFields legend="Contact info" name="contact" />
        </Formik>,
      );
      expect(await screen.findByLabelText('First name *')).toHaveValue(initialValues.firstName);
      expect(screen.getByLabelText(/Middle name/)).toHaveValue(initialValues.middleName);
      expect(screen.getByLabelText('Last name *')).toHaveValue(initialValues.lastName);
      expect(screen.getByLabelText(/Suffix/)).toHaveValue(initialValues.suffix);
      expect(screen.getByDisplayValue(initialValues.customerTelephone)).toBeInstanceOf(HTMLInputElement);
      expect(screen.getByDisplayValue(initialValues.customerEmail)).toBeInstanceOf(HTMLInputElement);
      expect(screen.getByLabelText(/Alternate Phone/)).toHaveValue(initialValues.secondaryPhone);
    });
  });
});
