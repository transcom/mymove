import React from 'react';
import { render, fireEvent } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { Formik } from 'formik';

import { AddressFields } from './AddressFields';

describe('AddressFields component', () => {
  it('renders a legend and all address inputs', () => {
    const { getByText, getByLabelText } = render(
      <Formik>
        <AddressFields legend="Address Form" name="address" />
      </Formik>,
    );
    expect(getByText('Address Form')).toBeInstanceOf(HTMLLegendElement);
    expect(getByLabelText(/Address 1/)).toBeInstanceOf(HTMLInputElement);
    expect(getByLabelText(/Address 2/)).toBeInstanceOf(HTMLInputElement);
    expect(getByLabelText(/City/)).toBeInstanceOf(HTMLInputElement);
    expect(getByLabelText(/State/)).toBeInstanceOf(HTMLSelectElement);
    expect(getByLabelText(/ZIP/)).toBeInstanceOf(HTMLInputElement);
  });

  describe('with pre-filled values', () => {
    it('renders a legend and all address inputs', () => {
      const initialValues = {
        address: {
          streetAddress1: '123 Main St',
          streetAddress2: 'Apt 3A',
          city: 'New York',
          state: 'NY',
          postalCode: '10002',
        },
      };

      const { getByLabelText } = render(
        <Formik initialValues={initialValues}>
          <AddressFields legend="Address Form" name="address" />
        </Formik>,
      );
      expect(getByLabelText(/Address 1/)).toHaveValue(initialValues.address.streetAddress1);
      expect(getByLabelText(/Address 2/)).toHaveValue(initialValues.address.streetAddress2);
      expect(getByLabelText(/City/)).toHaveValue(initialValues.address.city);
      expect(getByLabelText(/State/)).toHaveValue(initialValues.address.state);
      expect(getByLabelText(/ZIP/)).toHaveValue(initialValues.address.postalCode);
    });
  });

  describe('with validators', () => {
    it('puts the validator on the expected field', async () => {
      const initialValues = {
        address: {
          streetAddress1: '',
          streetAddress2: '',
          city: '',
          state: '',
          postalCode: '',
        },
      };

      const postalCodeErrorText = 'ZIP code must be 99999';

      const { getByLabelText, findByRole } = render(
        <Formik initialValues={initialValues}>
          {() => (
            <AddressFields
              legend="Address Form"
              name="address"
              validators={{
                postalCode: (value) => (value !== '99999' ? postalCodeErrorText : ''),
              }}
            />
          )}
        </Formik>,
      );

      const postalCodeInput = getByLabelText(/ZIP/);
      await userEvent.type(postalCodeInput, '12345');
      fireEvent.blur(postalCodeInput);

      const postalCodeError = await findByRole('alert');

      expect(postalCodeError).toHaveTextContent(postalCodeErrorText);
    });
  });

  describe('with on change postal code validator', () => {
    it('puts the validator on the expected field', async () => {
      const initialValues = {
        address: {
          streetAddress1: '',
          streetAddress2: '',
          city: '',
          state: '',
          postalCode: '',
        },
      };

      const postalCodeErrorText = 'ZIP code must be 99999';

      const { getByLabelText, findByRole } = render(
        <Formik initialValues={initialValues} validateOnChange>
          {({ handleChange, setFieldTouched }) => (
            <AddressFields
              legend="Address Form"
              name="address"
              formikFunctionsToValidatePostalCodeOnChange={{ handleChange, setFieldTouched }}
              validators={{
                postalCode: (value) => (value !== '99999' ? postalCodeErrorText : ''),
              }}
            />
          )}
        </Formik>,
      );

      const postalCodeInput = getByLabelText(/ZIP/);
      await userEvent.type(postalCodeInput, '12345');

      const postalCodeError = await findByRole('alert');

      expect(postalCodeError).toHaveTextContent(postalCodeErrorText);
    });
  });
});
