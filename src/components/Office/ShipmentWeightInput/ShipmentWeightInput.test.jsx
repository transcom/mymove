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

    expect(screen.getByText(/Previously recorded weight \(lbs\)/)).toBeInTheDocument();
  });

  it('populates Formik initialValues', () => {
    render(
      <Formik initialValues={{ ntsRecordedWeight: '4500' }}>
        <ShipmentWeightInput />
      </Formik>,
    );

    expect(screen.getByLabelText(/Previously recorded weight \(lbs\)/)).toHaveValue('4500');
  });
});
