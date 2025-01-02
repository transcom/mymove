import React from 'react';
import { render, screen } from '@testing-library/react';
import { Formik } from 'formik';
import { Provider } from 'react-redux';

import StorageFacilityAddress from './StorageFacilityAddress';

import { configureStore } from 'shared/store';

const mockStore = configureStore({});

describe('components/Office/StorageFacilityAddress', () => {
  it('renders correctly', () => {
    render(
      <Provider store={mockStore.store}>
        <Formik initialValues={{}}>
          <StorageFacilityAddress />
        </Formik>
      </Provider>,
    );

    expect(screen.getByRole('heading', { name: 'Storage facility address' })).toBeInTheDocument();
  });

  it('populates Formik initialValues', () => {
    render(
      <Provider store={mockStore.store}>
        <Formik
          initialValues={{
            storageFacility: {
              lotNumber: '42',
              address: {
                streetAddress1: '3373 NW Martin Luther King Blvd',
                city: 'San Antonio',
                state: 'TX',
                postalCode: '78234',
              },
            },
          }}
        >
          <StorageFacilityAddress />
        </Formik>
      </Provider>,
    );

    expect(screen.getByLabelText(/Address 1/)).toHaveValue('3373 NW Martin Luther King Blvd');
    expect(screen.getByTestId('State')).toHaveTextContent('TX');
    expect(screen.getByLabelText(/Lot number/)).toHaveValue('42');
  });
});
