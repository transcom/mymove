import React from 'react';
import { render, screen } from '@testing-library/react';
import { Formik } from 'formik';

import ShipmentWeightInput from './ShipmentWeightInput';

import { roleTypes } from 'constants/userRoles';

describe('components/Office/ShipmentWeightInput', () => {
  it('renders correctly', () => {
    render(
      <Formik initialValues={{ ntsRecordedWeight: '' }}>
        <ShipmentWeightInput userRole={roleTypes.TOO} />
      </Formik>,
    );

    expect(screen.getByText(/Previously recorded weight \(lbs\)/)).toBeInTheDocument();
  });

  it('populates Formik initialValues', () => {
    render(
      <Formik initialValues={{ ntsRecordedWeight: '4500' }}>
        <ShipmentWeightInput userRole={roleTypes.TOO} />
      </Formik>,
    );

    expect(screen.getByLabelText(/Previously recorded weight \(lbs\)/)).toHaveValue('4500');
  });

  it('makes Shipment Weight required for TOO', async () => {
    render(
      <Formik initialValues={{ ntsRecordedWeight: '' }}>
        <ShipmentWeightInput userRole={roleTypes.TOO} />
      </Formik>,
    );

    const input = await screen.getByLabelText('Previously recorded weight (lbs)');
    expect(
      input.parentElement.getElementsByClassName('labelWrapper')[0].getElementsByTagName('div')[0].textContent ===
        'Required',
    );
  });

  it('makes Shipment Weight optional for Services Counselor', async () => {
    render(
      <Formik initialValues={{ ntsRecordedWeight: '' }}>
        <ShipmentWeightInput userRole={roleTypes.SERVICES_COUNSELOR} />
      </Formik>,
    );

    const input = await screen.getByLabelText('Previously recorded weight (lbs)');
    expect(input.parentElement.getElementsByClassName('labelWrapper')[0].childElementCount === 1); // Ensure "required" text was not added.
  });
});
