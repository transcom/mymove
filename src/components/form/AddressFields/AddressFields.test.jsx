import React from 'react';
import { render, fireEvent } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { Formik } from 'formik';
import { Provider } from 'react-redux';

import { AddressFields } from './AddressFields';

import { configureStore } from 'shared/store';

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

  describe('zip city enabled with pre-filled values', () => {
    it('renders zip city lookup with info', () => {
      const initialValues = {
        address: {
          streetAddress1: '123 Main St',
          streetAddress2: 'Apt 3A',
          city: 'New York',
          state: 'NY',
          postalCode: '10002',
          county: 'NEW YORK',
        },
      };
      const mockStore = configureStore({});

      const { getByLabelText, getByTestId } = render(
        <Provider store={mockStore.store}>
          <Formik initialValues={initialValues}>
            {({ ...formikProps }) => {
              return <AddressFields legend="Address Form" name="address" locationLookup formikProps={formikProps} />;
            }}
          </Formik>
        </Provider>,
      );
      expect(getByLabelText('Address 1')).toHaveValue(initialValues.address.streetAddress1);
      expect(getByLabelText(/Address 2/)).toHaveValue(initialValues.address.streetAddress2);
      expect(getByTestId('City')).toHaveTextContent(initialValues.address.city);
      expect(getByTestId('State')).toHaveTextContent(initialValues.address.state);
      expect(getByTestId('ZIP')).toHaveTextContent(initialValues.address.postalCode);
    });
  });
});
