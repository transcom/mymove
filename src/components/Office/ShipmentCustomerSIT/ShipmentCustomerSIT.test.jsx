import React from 'react';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { Formik } from 'formik';

import ShipmentCustomerSIT from './ShipmentCustomerSIT';

describe('components/Office/ShipmentCustomerSIT', () => {
  it('defaults to customer not using SIT', () => {
    render(
      <Formik initialValues={{}}>
        <ShipmentCustomerSIT />
      </Formik>,
    );

    expect(screen.getByLabelText('Yes')).not.toBeChecked();
    expect(screen.getByLabelText('No')).toBeChecked();

    expect(screen.queryByLabelText('Destination')).not.toBeInTheDocument();
  });

  it('defaults to customer using SIT at destination', async () => {
    render(
      <Formik initialValues={{}}>
        <ShipmentCustomerSIT />
      </Formik>,
    );

    await userEvent.click(screen.getByLabelText('Yes'));

    expect(await screen.findByLabelText('Destination')).toBeChecked();
  });
});
