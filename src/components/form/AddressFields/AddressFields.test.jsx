import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { Formik } from 'formik';
import { Provider } from 'react-redux';

import { AddressFields } from './AddressFields';

import * as ghcApi from 'services/ghcApi';
import { configureStore } from 'shared/store';
import { FEATURE_FLAG_KEYS } from 'shared/constants';
import { isBooleanFlagEnabled } from 'utils/featureFlags';

jest.mock('services/ghcApi');

jest.mock('utils/featureFlags', () => ({
  ...jest.requireActual('utils/featureFlags'),
  isBooleanFlagEnabled: jest.fn().mockImplementation(() => Promise.resolve(true)),
}));

const mockLoggedInUser = {
  id: '123',
  office_user: {
    id: '456',
  },
};

describe('AddressFields component', () => {
  const mockStore = configureStore({
    entities: {
      user: {
        [mockLoggedInUser.id]: mockLoggedInUser,
      },
    },
  });

  it('renders a legend and all address inputs', () => {
    const { getByText, getByLabelText, getByTestId } = render(
      <Provider store={mockStore.store}>
        <Formik>
          <AddressFields legend="Address Form" name="address" />
        </Formik>
      </Provider>,
    );
    expect(getByText('Address Form')).toBeInstanceOf(HTMLLegendElement);
    expect(getByLabelText(/Address 1/)).toBeInstanceOf(HTMLInputElement);
    expect(getByLabelText(/Address 2/)).toBeInstanceOf(HTMLInputElement);
    expect(getByTestId('City')).toBeInstanceOf(HTMLLabelElement);
    expect(getByTestId('State')).toBeInstanceOf(HTMLLabelElement);
    expect(getByTestId('ZIP')).toBeInstanceOf(HTMLLabelElement);
    expect(getByLabelText(/Location Lookup/)).toBeInstanceOf(HTMLInputElement);
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
          county: 'NEW YORK',
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
      expect(
        screen.getAllByText(
          `${initialValues.address.city}, ${initialValues.address.state} ${initialValues.address.postalCode} (${initialValues.address.county})`,
        ),
      );
    });
  });

  it('calls onCountryChange when a country value is entered', async () => {
    const mockCountryData = [
      {
        code: 'US',
        id: '791899e6-cd77-46f2-981b-176ecb8d7098',
        name: 'UNITED STATES',
      },
    ];

    ghcApi.searchCountry.mockResolvedValue(mockCountryData);

    const onCountryChange = jest.fn();

    render(
      <Provider store={mockStore.store}>
        <Formik initialValues={{}}>
          {(formikProps) => (
            <AddressFields name="address" formikProps={formikProps} onCountryChange={onCountryChange} />
          )}
        </Formik>
      </Provider>,
    );

    await waitFor(() => {
      expect(isBooleanFlagEnabled).toHaveBeenCalledWith(FEATURE_FLAG_KEYS.OCONUS_CITY_FINDER);
    });

    const countryInput = screen.getByRole('combobox', { name: /country/i });

    await userEvent.type(countryInput, 'US');

    await waitFor(() => {
      expect(ghcApi.searchCountry).toHaveBeenCalledWith('US');
    });

    await userEvent.type(countryInput, '{arrowdown}{enter}');
    await waitFor(() => {
      expect(onCountryChange).toHaveBeenCalledWith('US');
    });
  });
});
