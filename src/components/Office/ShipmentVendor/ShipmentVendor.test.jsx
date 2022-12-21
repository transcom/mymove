import React from 'react';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { Formik } from 'formik';

import ShipmentVendor from './ShipmentVendor';

describe('components/Office/ShipmentVendor', () => {
  it('renders correctly', () => {
    render(
      <Formik initialValues={{}}>
        <ShipmentVendor />
      </Formik>,
    );

    expect(screen.getByLabelText('GHC prime contractor')).toBeChecked();
    expect(screen.getByLabelText('External vendor')).not.toBeChecked();
    expect(screen.queryByText('This shipment will be sent to the GHC prime contractor.')).not.toBeInTheDocument();
  });

  it('populates Formik initialValues', () => {
    render(
      <Formik initialValues={{ usesExternalVendor: true }}>
        <ShipmentVendor />
      </Formik>,
    );

    expect(screen.getByLabelText('GHC prime contractor')).not.toBeChecked();
    expect(screen.getByLabelText('External vendor')).toBeChecked();
    expect(screen.getByRole('list')).toBeInTheDocument();
  });

  it('changes messaging based on interaction', async () => {
    render(
      <Formik initialValues={{ usesExternalVendor: false }}>
        <ShipmentVendor />
      </Formik>,
    );

    expect(screen.queryByRole('list')).not.toBeInTheDocument();
    expect(screen.queryByText('This shipment will be sent to the GHC prime contractor.')).not.toBeInTheDocument();

    userEvent.click(screen.getByLabelText('External vendor'));
    expect(await screen.findByRole('list')).toBeInTheDocument();
    userEvent.click(screen.getByLabelText('GHC prime contractor'));
    expect(await screen.findByText('This shipment will be sent to the GHC prime contractor.')).toBeInTheDocument();
  });
});
