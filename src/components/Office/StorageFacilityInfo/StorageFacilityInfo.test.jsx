import React from 'react';
import { render, screen } from '@testing-library/react';
import { Formik } from 'formik';

import StorageFacilityInfo from './StorageFacilityInfo';

describe('components/Office/StorageFacilityInfo', () => {
  it('renders correctly', () => {
    render(
      <Formik initialValues={{}}>
        <StorageFacilityInfo />
      </Formik>,
    );

    expect(screen.getByRole('heading', { name: 'Storage facility info' })).toBeInTheDocument();
  });

  it('populates Formik initialValues', () => {
    render(
      <Formik
        initialValues={{
          facilityName: 'Most Excellent Storage',
        }}
      >
        <StorageFacilityInfo />
      </Formik>,
    );

    expect(screen.getByLabelText('Facility name')).toHaveValue('Most Excellent Storage');
  });
});
