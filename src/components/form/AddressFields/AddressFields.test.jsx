import React from 'react';
import { render } from '@testing-library/react';
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
    expect(getByLabelText('Address 1')).toBeInstanceOf(HTMLInputElement);
    expect(getByLabelText(/Address 2/)).toBeInstanceOf(HTMLInputElement);
    expect(getByLabelText('City')).toBeInstanceOf(HTMLInputElement);
    expect(getByLabelText('State')).toBeInstanceOf(HTMLSelectElement);
    expect(getByLabelText('ZIP')).toBeInstanceOf(HTMLInputElement);
  });

  describe('with pre-filled values', () => {
    it('renders a legend and all address inputs', () => {
      const initialValues = {
        address: {
          street_address_1: '123 Main St',
          street_address_2: 'Apt 3A',
          city: 'New York',
          state: 'NY',
          postal_code: '10002',
        },
      };

      const { getByLabelText } = render(
        <Formik initialValues={initialValues}>
          <AddressFields legend="Address Form" name="address" />
        </Formik>,
      );
      expect(getByLabelText('Address 1')).toHaveValue(initialValues.address.street_address_1);
      expect(getByLabelText(/Address 2/)).toHaveValue(initialValues.address.street_address_2);
      expect(getByLabelText('City')).toHaveValue(initialValues.address.city);
      expect(getByLabelText('State')).toHaveValue(initialValues.address.state);
      expect(getByLabelText('ZIP')).toHaveValue(initialValues.address.postal_code);
    });
  });
});
