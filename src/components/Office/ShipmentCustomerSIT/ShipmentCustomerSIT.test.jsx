import React from 'react';
import { render, screen, act } from '@testing-library/react';
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

  it('defaults to customer using SIT at destination and shows asterisks for required fields', async () => {
    render(
      <Formik initialValues={{}}>
        <ShipmentCustomerSIT />
      </Formik>,
    );

    await act(async () => {
      await userEvent.click(screen.getByLabelText('Yes'));
    });

    expect(await screen.findByLabelText('Destination')).toBeChecked();

    expect(document.querySelector('#reqAsteriskMsg')).toHaveTextContent('Fields marked with * are required.');

    expect(screen.getByLabelText('Estimated storage start *')).toBeInTheDocument();
    expect(screen.getByLabelText('Estimated storage end *')).toBeInTheDocument();
    expect(screen.getByLabelText('Estimated SIT weight *')).toBeInTheDocument();
  });
});
