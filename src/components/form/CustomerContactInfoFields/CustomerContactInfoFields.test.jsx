import React from 'react';
import { render, screen } from '@testing-library/react';
import { Formik } from 'formik';

import { CustomerContactInfoFields } from './index';

describe('ContactInfoFields component', () => {
  it('renders a legend and all service member contact info inputs during editing', () => {
    const editContact = true;
    render(
      <Formik>
        <CustomerContactInfoFields legend="Your contact info" signIn={editContact} />
      </Formik>,
    );
    expect(screen.getByText('Your contact info')).toBeInstanceOf(HTMLLegendElement);
    expect(screen.getByLabelText('Preferred Name')).toBeInstanceOf(HTMLInputElement);
    expect(screen.getByLabelText('Best contact phone')).toBeInstanceOf(HTMLInputElement);
    expect(screen.getByLabelText(/Alt. phone/)).toBeInstanceOf(HTMLInputElement);
    expect(screen.getByLabelText('Personal email')).toBeInstanceOf(HTMLInputElement);
    expect(screen.getByLabelText('Phone')).not.toBeChecked();
    expect(screen.getByLabelText('Email')).not.toBeChecked();
  });

  it('renders a legend and all service member contact info inputs when signing up.', () => {
    const editContact = false;
    render(
      <Formik>
        <CustomerContactInfoFields legend="Your contact info" signIn={editContact} />
      </Formik>,
    );
    expect(screen.getByText('Your contact info')).toBeInstanceOf(HTMLLegendElement);
    expect(screen.getByLabelText('Best contact phone')).toBeInstanceOf(HTMLInputElement);
    expect(screen.getByLabelText(/Alt. phone/)).toBeInstanceOf(HTMLInputElement);
    expect(screen.getByLabelText('Personal email')).toBeInstanceOf(HTMLInputElement);
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
        preferred_name: 'Alex',
      };
      const editContact = true;

      render(
        <Formik initialValues={initialValues}>
          <CustomerContactInfoFields legend="Your contact info" name="contact" signIn={editContact} />
        </Formik>,
      );
      expect(await screen.findByLabelText('Best contact phone')).toHaveValue(initialValues.telephone);
      expect(screen.getByLabelText(/Alt. phone/)).toHaveValue(initialValues.secondary_telephone);
      expect(screen.getByLabelText('Personal email')).toHaveValue(initialValues.personal_email);
      expect(screen.getByLabelText('Preferred Name')).toHaveValue(initialValues.preferred_name);
      expect(screen.getByLabelText('Phone')).toBeChecked();
      expect(screen.getByLabelText('Email')).toBeChecked();
    });
  });
});
