import React from 'react';
import { render, screen } from '@testing-library/react';
import { Formik } from 'formik';

import StorageFacilityInfo from './StorageFacilityInfo';

import { roleTypes } from 'constants/userRoles';

describe('components/Office/StorageFacilityInfo', () => {
  it('renders correctly', () => {
    render(
      <Formik initialValues={{}}>
        <StorageFacilityInfo userRole={roleTypes.SERVICES_COUNSELOR} />
      </Formik>,
    );

    expect(screen.getByRole('heading', { name: 'Storage facility info' })).toBeInTheDocument();
    expect(screen.getAllByText(/Optional/)).toHaveLength(3);
  });

  it('populates Formik initialValues', async () => {
    render(
      <Formik
        initialValues={{
          serviceOrderNumber: '12341234',
          storageFacility: { facilityName: 'Most Excellent Storage', phone: '555-456-4567' },
        }}
      >
        <StorageFacilityInfo userRole={roleTypes.SERVICES_COUNSELOR} />
      </Formik>,
    );

    expect(screen.getByLabelText('Facility name')).toHaveValue('Most Excellent Storage');
    expect(screen.getByLabelText(/Service order number/)).toHaveValue('12341234');
    expect(await screen.findByLabelText(/Phone/)).toHaveValue('555-456-4567');
  });

  it('makes Service Order Number required for TOO', () => {
    render(
      <Formik initialValues={{}}>
        <StorageFacilityInfo userRole={roleTypes.TOO} />
      </Formik>,
    );

    expect(screen.getAllByText(/Optional/)).toHaveLength(2);
  });
});
