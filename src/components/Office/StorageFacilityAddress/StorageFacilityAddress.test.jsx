import React from 'react';
import { render, screen } from '@testing-library/react';
import { Formik } from 'formik';

import StorageFacilityAddress from './StorageFacilityAddress';

describe('components/Office/StorageFacilityAddress', () => {
  it('renders correctly', () => {
    render(
      <Formik initialValues={{}}>
        <StorageFacilityAddress />
      </Formik>,
    );

    expect(screen.getByRole('heading', { name: 'Storage facility address' })).toBeInTheDocument();
  });

  it('populates Formik initialValues', () => {
    render(
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
      </Formik>,
    );

    expect(screen.getByLabelText(/Address 1/)).toHaveValue('3373 NW Martin Luther King Blvd');
    expect(screen.getByLabelText('State')).toHaveValue('TX');
    expect(screen.getByLabelText(/Lot number/)).toHaveValue('42');
  });
});
