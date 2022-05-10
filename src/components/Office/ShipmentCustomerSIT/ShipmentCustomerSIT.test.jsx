import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
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

    userEvent.click(screen.getByLabelText('Yes'));

    expect(await screen.findByLabelText('Destination')).toBeChecked();
  });

  it('disables Calculate SIT button until all fields filled out', async () => {
    const onCalculateClick = jest.fn();
    render(
      <Formik initialValues={{}}>
        <ShipmentCustomerSIT onCalculateClick={onCalculateClick} />
      </Formik>,
    );

    userEvent.click(screen.getByLabelText('Yes'));
    expect(await screen.findByRole('button', { name: 'Calculate SIT' })).toBeDisabled();

    userEvent.type(screen.getByLabelText('Estimated SIT weight'), '5725');
    userEvent.type(screen.getByLabelText('Estimated storage start'), '02 May 2022');
    userEvent.type(screen.getByLabelText('Estimated storage end'), '09 May 2022');

    await waitFor(() => expect(screen.getByRole('button', { name: 'Calculate SIT' })).toBeEnabled());

    userEvent.click(screen.getByRole('button', { name: 'Calculate SIT' }));

    await waitFor(() =>
      expect(onCalculateClick).toBeCalledWith({
        location: 'destination',
        weight: '5725',
        start: '02 May 2022',
        end: '09 May 2022',
      }),
    );
  });
});
