import React from 'react';
import { render, screen } from '@testing-library/react';
import { Formik } from 'formik';

import ShipmentWeightInput from './ShipmentWeightInput';

describe('components/Office/ShipmentWeightInput', () => {
  it('renders correctly', () => {
    render(
      <Formik initialValues={{ ntsRecordedWeight: '' }}>
        <ShipmentWeightInput />
      </Formik>,
    );

    expect(screen.getByText(/Previous Recorded Weight \(lbs\)/)).toBeInTheDocument();
  });

  it('populates Formik initialValues', () => {
    render(
      <Formik initialValues={{ ntsRecordedWeight: '4500' }}>
        <ShipmentWeightInput />
      </Formik>,
    );

    expect(screen.getByLabelText(/Previous Recorded Weight \(lbs\)/)).toHaveValue('4500');
  });
});
