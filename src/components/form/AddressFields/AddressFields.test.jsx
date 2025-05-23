import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { Formik } from 'formik';
import { Provider } from 'react-redux';

import { AddressFields } from './AddressFields';

import * as ghcApi from 'services/ghcApi';
import { configureStore } from 'shared/store';

jest.mock('services/ghcApi');

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

  it('calls onCountryChanged when a country value is entered', async () => {
    const mockLocationData = [
      {
        city: 'New York',
        state: 'NY',
        county: 'New York County',
        postalCode: '10001',
        usPostRegionCitiesID: '123456',
      },
    ];

    ghcApi.searchLocationByZipCityState.mockResolvedValue(mockLocationData);

    const onCountryChanged = jest.fn();

    render(
      <Provider store={mockStore.store}>
        <Formik initialValues={{}}>
          {(formikProps) => (
            <AddressFields name="address" formikProps={formikProps} onCountryChanged={onCountryChanged} />
          )}
        </Formik>
      </Provider>,
    );

    const countryInput = screen.getByRole('combobox', { name: /country/i });

    await userEvent.type(countryInput, 'US');

    await waitFor(() => {
      expect(ghcApi.searchLocationByZipCityState).toHaveBeenCalledWith('New York');
    });

    // Simulate selecting the option
    await userEvent.type(countryInput, '{arrowdown}{enter}');

    await waitFor(() => {
      expect(onCountryChanged).toHaveBeenCalledWith(
        expect.objectContaining({
          city: 'New York',
          state: 'NY',
          county: 'New York County',
          postalCode: '10001',
          usPostRegionCitiesID: '123456',
        }),
      );
    });
  });
});
