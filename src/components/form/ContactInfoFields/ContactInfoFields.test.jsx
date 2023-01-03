import React from 'react';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { Formik } from 'formik';

import { ContactInfoFields } from './ContactInfoFields';

describe('ContactInfoFields component', () => {
  it('renders a legend and all contact info inputs', () => {
    render(
      <Formik>
        <ContactInfoFields legend="Contact Info Form" name="contact" />
      </Formik>,
    );
    expect(screen.getByText('Contact Info Form')).toBeInstanceOf(HTMLLegendElement);
    expect(screen.getByLabelText('First name')).toBeInstanceOf(HTMLInputElement);
    expect(screen.getByLabelText('Last name')).toBeInstanceOf(HTMLInputElement);
    expect(screen.getByLabelText('Phone')).toBeInstanceOf(HTMLInputElement);
    expect(screen.getByLabelText('Email')).toBeInstanceOf(HTMLInputElement);
  });

  describe('with pre-filled values', () => {
    it('renders a legend and all address inputs', async () => {
      const initialValues = {
        contact: {
          firstName: 'Test',
          lastName: 'Person',
          phone: '555-123-4567',
          email: 'test@example.com',
        },
      };

      render(
        <Formik initialValues={initialValues}>
          <ContactInfoFields legend="Contact Info Form" name="contact" />
        </Formik>,
      );

      expect(await screen.findByLabelText('First name')).toHaveValue(initialValues.contact.firstName);
      expect(screen.getByLabelText('Last name')).toHaveValue(initialValues.contact.lastName);
      expect(screen.getByLabelText('Phone')).toHaveValue(initialValues.contact.phone);
      expect(screen.getByLabelText('Email')).toHaveValue(initialValues.contact.email);
    });
  });

  describe('with inserted values', () => {
    it('renders phone numbers formatted with dashes', async () => {
      render(
        <Formik initialValues={{}}>
          <ContactInfoFields legend="Contact Info Form" name="contact" />
        </Formik>,
      );

      expect(screen.getByLabelText('Phone')).toBeInTheDocument();
      await userEvent.type(screen.getByLabelText('Phone'), '5555555555');
      expect(await screen.findByLabelText('Phone')).toHaveValue('555-555-5555');
    });
  });
});
