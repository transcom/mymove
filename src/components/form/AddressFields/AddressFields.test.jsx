import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import { Formik } from 'formik';
import { Provider } from 'react-redux';

import { AddressFields } from './AddressFields';

import { configureStore } from 'shared/store';
import { isBooleanFlagEnabled } from 'utils/featureFlags';

jest.mock('utils/featureFlags', () => ({
  ...jest.requireActual('utils/featureFlags'),
  isBooleanFlagEnabled: jest.fn().mockImplementation(() => Promise.resolve(true)),
}));

describe('AddressFields component', () => {
  const mockStore = configureStore({});

  it('renders a legend and all address inputs and asterisks for required fields', () => {
    isBooleanFlagEnabled.mockImplementation(() => Promise.resolve(true));
    const initialValues = {
      address: {
        streetAddress1: '123 Main St',
        streetAddress2: 'Apt 3A',
        city: 'New York',
        state: 'NY',
        postalCode: '10002',
        county: 'NEW YORK',
        country: {
          code: 'US',
          name: 'UNITED STATES',
          id: '791899e6-cd77-46f2-981b-176ecb8d7098',
        },
        countryID: '791899e6-cd77-46f2-981b-176ecb8d7098',
      },
    };
    const { getByText, getByLabelText, getByTestId } = render(
      <Provider store={mockStore.store}>
        <Formik initialValues={initialValues}>
          <AddressFields legend="Address Form" name="address" />
        </Formik>
      </Provider>,
    );
    expect(getByTestId('reqAsteriskMsg')).toBeInTheDocument();

    expect(getByText('Address Form')).toBeInstanceOf(HTMLLegendElement);
    expect(getByLabelText(/Address 1/)).toBeInstanceOf(HTMLInputElement);
    expect(getByLabelText(/Address 1 */)).toBeInTheDocument();
    expect(getByLabelText(/Address 2/)).toBeInstanceOf(HTMLInputElement);
    expect(getByTestId('City')).toBeInstanceOf(HTMLLabelElement);
    expect(getByTestId('State')).toBeInstanceOf(HTMLLabelElement);
    expect(getByTestId('ZIP')).toBeInstanceOf(HTMLLabelElement);
    expect(getByLabelText(/Country Lookup/)).toBeInstanceOf(HTMLInputElement);
    expect(getByLabelText(/Location Lookup/)).toBeInstanceOf(HTMLInputElement);
    expect(getByLabelText(/Location Lookup */)).toBeInTheDocument();
  });

  describe('with pre-filled values', () => {
    it('renders a legend and all address inputs', async () => {
      isBooleanFlagEnabled.mockImplementation(() => Promise.resolve(true));
      const initialValues = {
        address: {
          streetAddress1: '123 Main St',
          streetAddress2: 'Apt 3A',
          city: 'New York',
          state: 'NY',
          postalCode: '10002',
          county: 'NEW YORK',
          country: {
            code: 'US',
            name: 'UNITED STATES',
            id: '791899e6-cd77-46f2-981b-176ecb8d7098',
          },
          countryID: '791899e6-cd77-46f2-981b-176ecb8d7098',
        },
      };

      const { getByLabelText, getByTestId } = render(
        <Provider store={mockStore.store}>
          <Formik initialValues={initialValues}>
            <AddressFields legend="Address Form" name="address" />
          </Formik>
        </Provider>,
      );
      expect(getByLabelText(/Address 1/)).toHaveValue(initialValues.address.streetAddress1);
      expect(getByLabelText(/Address 2/)).toHaveValue(initialValues.address.streetAddress2);
      expect(getByTestId('City')).toHaveTextContent(initialValues.address.city);
      expect(getByTestId('State')).toHaveTextContent(initialValues.address.state);
      expect(getByTestId('ZIP')).toHaveTextContent(initialValues.address.postalCode);
      expect(
        screen.getAllByText(
          `${initialValues.address.city}, ${initialValues.address.state} ${initialValues.address.postalCode} (${initialValues.address.county})`,
        ),
      );

      await waitFor(() => {
        const elements = screen.queryAllByText(/(US)/);
        expect(elements).toHaveLength(1);
      });
    });
  });

  describe('zip city enabled with pre-filled values', () => {
    it('renders zip city lookup with info', async () => {
      isBooleanFlagEnabled.mockImplementation(() => Promise.resolve(true));
      const initialValues = {
        address: {
          streetAddress1: '123 Main St',
          streetAddress2: 'Apt 3A',
          city: 'New York',
          state: 'NY',
          postalCode: '10002',
          county: 'NEW YORK',
          country: {
            code: 'US',
            name: 'UNITED STATES',
            id: '791899e6-cd77-46f2-981b-176ecb8d7098',
          },
          countryID: '791899e6-cd77-46f2-981b-176ecb8d7098',
        },
      };

      const { getByLabelText, getByTestId } = render(
        <Provider store={mockStore.store}>
          <Formik initialValues={initialValues}>
            {({ ...formikProps }) => {
              return <AddressFields legend="Address Form" name="address" locationLookup formikProps={formikProps} />;
            }}
          </Formik>
        </Provider>,
      );
      expect(getByLabelText(/Address 1/)).toHaveValue(initialValues.address.streetAddress1);
      expect(getByLabelText(/Address 2/)).toHaveValue(initialValues.address.streetAddress2);
      expect(getByTestId('City')).toHaveTextContent(initialValues.address.city);
      expect(getByTestId('State')).toHaveTextContent(initialValues.address.state);
      expect(getByTestId('ZIP')).toHaveTextContent(initialValues.address.postalCode);
      expect(
        screen.getAllByText(
          `${initialValues.address.city}, ${initialValues.address.state} ${initialValues.address.postalCode} (${initialValues.address.county})`,
        ),
      );

      await waitFor(() => {
        const elements = screen.queryAllByText(/(US)/);
        expect(elements).toHaveLength(1);
      });
    });
  });
});
