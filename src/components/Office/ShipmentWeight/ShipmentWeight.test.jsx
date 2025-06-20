import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
// import userEvent from '@testing-library/user-event';
import { Formik } from 'formik';

import ShipmentWeight from './ShipmentWeight';

describe('components/Office/ShipmentWeight', () => {
  it('defaults to customer not using Pro-gear', () => {
    render(
      <Formik initialValues={{}}>
        <ShipmentWeight />
      </Formik>,
    );

    expect(screen.getByLabelText('Yes')).not.toBeChecked();
    expect(screen.getByLabelText('No')).toBeChecked();

    expect(screen.queryByLabelText('Estimated pro-gear weight')).not.toBeInTheDocument();
    expect(screen.queryByLabelText('Estimated spouse pro-gear weight')).not.toBeInTheDocument();
  });

  it('displays estimated weight and pro-gear data', async () => {
    render(
      <Formik
        initialValues={{
          hasProGear: true,
          estimatedWeight: '4000',
          proGearWeight: '3000',
          spouseProGearWeight: '2000',
        }}
      >
        <ShipmentWeight />
      </Formik>,
    );
    await waitFor(() => {
      expect(screen.getByLabelText('Yes')).toBeChecked();
      expect(screen.getByLabelText('No')).not.toBeChecked();

      expect(screen.queryByLabelText('Estimated pro-gear weight')).toBeInTheDocument();
      expect(screen.queryByLabelText('Estimated spouse pro-gear weight')).toBeInTheDocument();
    });
  });
});
