// SmallPackageForm.test.jsx
import React from 'react';
import { render, screen } from '@testing-library/react';
import { Formik } from 'formik';

import SmallPackageForm from './SmallPackageForm';

// Define two sets of initial values for Formik:
const initialValuesProGearTrue = {
  amount: '',
  trackingNumber: '',
  isProGear: 'true',
  proGearBelongsToSelf: 'true',
  proGearDescription: '',
  weightShipped: '',
};

const initialValuesProGearFalse = {
  amount: '',
  trackingNumber: '',
  isProGear: 'false',
  proGearBelongsToSelf: '',
  proGearDescription: '',
  weightShipped: '',
};

describe('SmallPackageForm', () => {
  test('displays pro gear fields when isProGear is true and asterisks for required fields', () => {
    render(
      <Formik initialValues={initialValuesProGearTrue} onSubmit={() => {}}>
        <SmallPackageForm />
      </Formik>,
    );

    expect(screen.getByLabelText(/Package shipment cost/i)).toBeInTheDocument();
    expect(screen.getByLabelText(/Tracking number/i)).toBeInTheDocument();

    expect(screen.getByLabelText(/Yes/i)).toBeInTheDocument();
    expect(screen.getByLabelText(/No/i)).toBeInTheDocument();

    // since isProGear is "true", additional pro gear fields should appear
    expect(screen.getByText(/Who does this pro-gear belong to/i)).toBeInTheDocument();
    expect(screen.getByLabelText(/^Me$/i)).toBeInTheDocument();
    expect(screen.getByLabelText(/My Spouse/i)).toBeInTheDocument();
    expect(screen.getByLabelText(/Brief description of the pro-gear/i)).toBeInTheDocument();
    expect(screen.getByTestId('proGearWeight')).toBeInTheDocument();

    // the standard weightShipped field should not be rendered
    expect(screen.queryByTestId('weightShipped')).not.toBeInTheDocument();
  });

  test('displays non-pro gear field when isProGear is false', () => {
    render(
      <Formik initialValues={initialValuesProGearFalse} onSubmit={() => {}}>
        <SmallPackageForm />
      </Formik>,
    );

    expect(screen.getByLabelText(/Package shipment cost/i)).toBeInTheDocument();
    expect(screen.getByLabelText(/Tracking number/i)).toBeInTheDocument();

    expect(screen.getByLabelText(/Yes/i)).toBeInTheDocument();
    expect(screen.getByLabelText(/No/i)).toBeInTheDocument();

    // isProGear is "false", extra pro gear fields should NOT be rendered
    expect(screen.queryByText(/Who does this pro-gear belong to/i)).not.toBeInTheDocument();
    expect(screen.queryByLabelText(/^Me$/i)).not.toBeInTheDocument();
    expect(screen.queryByLabelText(/My Spouse/i)).not.toBeInTheDocument();
    expect(screen.queryByLabelText(/Brief description of the pro-gear/i)).not.toBeInTheDocument();
    expect(screen.queryByTestId('proGearWeight')).not.toBeInTheDocument();

    // instead, the default weight field should be rendered
    expect(screen.getByTestId('weightShipped')).toBeInTheDocument();
  });
});
