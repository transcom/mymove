import React from 'react';
import { render } from '@testing-library/react';
import { Formik } from 'formik';

import { CustomerContactInfoFields } from './index';

describe('ContactInfoFields component', () => {
  it('renders a legend and all service member contact info inputs', () => {
    const { getByText, getByLabelText } = render(
      <Formik>
        <CustomerContactInfoFields legend="Your contact info" />
      </Formik>,
    );
    expect(getByText('Your contact info')).toBeInstanceOf(HTMLLegendElement);
    expect(getByLabelText('Best contact phone')).toBeInstanceOf(HTMLInputElement);
    expect(getByLabelText(/Alt. phone/)).toBeInstanceOf(HTMLInputElement);
    expect(getByLabelText('Personal email')).toBeInstanceOf(HTMLInputElement);
    expect(getByLabelText('Phone')).not.toBeChecked();
    expect(getByLabelText('Email')).not.toBeChecked();
  });

  describe('with pre-filled values', () => {
    it('renders a legend and all service member contact info inputs', () => {
      const initialValues = {
        telephone: '555-123-4567',
        secondary_telephone: '555-890-1234',
        personal_email: 'test@example.com',
        phone_is_preferred: true,
        email_is_preferred: true,
      };

      const { getByLabelText } = render(
        <Formik initialValues={initialValues}>
          <CustomerContactInfoFields legend="Your contact info" name="contact" />
        </Formik>,
      );
      expect(getByLabelText('Best contact phone')).toHaveValue(initialValues.telephone);
      expect(getByLabelText(/Alt. phone/)).toHaveValue(initialValues.secondary_telephone);
      expect(getByLabelText('Personal email')).toHaveValue(initialValues.personal_email);
      expect(getByLabelText('Phone')).toBeChecked();
      expect(getByLabelText('Email')).toBeChecked();
    });
  });
});
