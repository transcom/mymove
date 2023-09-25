import React from 'react';
import { render, screen } from '@testing-library/react';
import { Formik } from 'formik';

import OktaInfoFields from './index';

describe('OktaInfoFields component', () => {
  it('renders a legend and all okta info inputs', () => {
    render(
      <Formik>
        <OktaInfoFields legend="Your contact info" />
      </Formik>,
    );
    expect(screen.getByLabelText('Okta Username')).toBeInstanceOf(HTMLInputElement);
    expect(screen.getByLabelText('Okta Email')).toBeInstanceOf(HTMLInputElement);
    expect(screen.getByLabelText('First Name')).toBeInstanceOf(HTMLInputElement);
    expect(screen.getByLabelText('Last Name')).toBeInstanceOf(HTMLInputElement);
    expect(screen.getByLabelText('DoD ID number | EDIPI')).toBeInstanceOf(HTMLInputElement);
  });

  describe('with pre-filled values', () => {
    it('renders a legend and all service member contact info inputs', async () => {
      const initialValues = {
        oktaUsername: 'user@okta.mil',
        oktaEmail: 'user@okta.mil',
        oktaFirstName: 'Okta',
        oktaLastName: 'User',
        oktaEdipi: '8888889990',
      };

      render(
        <Formik initialValues={initialValues}>
          <OktaInfoFields legend="Your Okta Profile" name="okta" />
        </Formik>,
      );
      expect(await screen.findByLabelText('Okta Username')).toHaveValue(initialValues.oktaUsername);
      expect(screen.getByLabelText('Okta Email')).toHaveValue(initialValues.oktaEmail);
      expect(screen.getByLabelText('First Name')).toHaveValue(initialValues.oktaFirstName);
      expect(screen.getByLabelText('Last Name')).toHaveValue(initialValues.oktaLastName);
      expect(screen.getByLabelText('DoD ID number | EDIPI')).toHaveValue(initialValues.oktaEdipi);
    });
  });
});
