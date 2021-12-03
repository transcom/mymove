import React from 'react';
import { render, screen } from '@testing-library/react';
import { Formik } from 'formik';

import ShipmentWeightInput from './ShipmentWeightInput';

describe('components/Office/ShipmentWeightInput', () => {
  it('renders correctly', () => {
    render(
      <Formik initialValues={{}}>
        <ShipmentWeightInput />
      </Formik>,
    );

    expect(screen.getByText(/Shipment weight \(lbs\)/)).toBeInTheDocument();
  });

  it('populates Formik initialValues', () => {
    render(
      <Formik initialValues={{ primeActualWeight: '4500' }}>
        <ShipmentWeightInput />
      </Formik>,
    );

    expect(screen.getByLabelText(/Shipment weight \(lbs\)/)).toHaveValue('4500');
  });
});
