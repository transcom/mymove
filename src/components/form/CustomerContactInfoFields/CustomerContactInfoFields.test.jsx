import React from 'react';
import { render, screen } from '@testing-library/react';
import { Formik } from 'formik';

import { CustomerContactInfoFields } from './index';

describe('ContactInfoFields component', () => {
  it('renders a legend and all service member contact info inputs and asterisks for required fields', () => {
    render(
      <Formik>
        <CustomerContactInfoFields legend="Your contact info" />
      </Formik>,
    );
    expect(screen.getByText('Your contact info')).toBeInstanceOf(HTMLLegendElement);
    expect(screen.getByLabelText('Best contact phone *')).toBeInTheDocument();
    expect(screen.getByLabelText('Best contact phone *')).toBeInstanceOf(HTMLInputElement);
    expect(screen.getByLabelText('Best contact phone *')).toBeRequired();

    expect(screen.getByLabelText(/Alt. phone/)).toBeInstanceOf(HTMLInputElement);
    expect(screen.getByLabelText('Personal email *')).toBeInTheDocument();
    expect(screen.getByLabelText('Personal email *')).toBeInstanceOf(HTMLInputElement);
    expect(screen.getByLabelText('Personal email *')).toBeRequired();
    expect(screen.getByLabelText('Phone')).not.toBeChecked();
    expect(screen.getByLabelText('Email')).not.toBeChecked();
  });

  describe('with pre-filled values', () => {
    it('renders a legend and all service member contact info inputs', async () => {
      const initialValues = {
        telephone: '555-123-4567',
        secondary_telephone: '555-890-1234',
        personal_email: 'test@example.com',
        phone_is_preferred: true,
        email_is_preferred: true,
      };

      render(
        <Formik initialValues={initialValues}>
          <CustomerContactInfoFields legend="Your contact info" name="contact" />
        </Formik>,
      );
      expect(screen.getByLabelText('Best contact phone *')).toHaveValue(initialValues.telephone);
      expect(screen.getByLabelText(/Alt. phone/)).toHaveValue(initialValues.secondary_telephone);
      expect(screen.getByLabelText('Personal email *')).toHaveValue(initialValues.personal_email);
      expect(screen.getByLabelText('Phone')).toBeChecked();
      expect(screen.getByLabelText('Email')).toBeChecked();
    });
  });
});
