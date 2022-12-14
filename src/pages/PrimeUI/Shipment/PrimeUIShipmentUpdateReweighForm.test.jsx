import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import PrimeUIShipmentUpdateReweighForm from './PrimeUIShipmentUpdateReweighForm';

describe('PrimeUIShipmentUpdateReweighForm', () => {
  const testProps = {
    initialValues: {
      reweighWeight: '0',
      reweighRemarks: '',
    },
    onSubmit: jest.fn().mockImplementation(() => Promise.resolve()),
    handleClose: jest.fn(),
  };

  it('renders the form', async () => {
    render(<PrimeUIShipmentUpdateReweighForm {...testProps} />);

    expect(await screen.findByLabelText('Reweigh Weight (lbs)')).toBeInstanceOf(HTMLInputElement);
    expect(screen.getByTestId('remarks')).toBeInstanceOf(HTMLTextAreaElement);
  });

  it('submits the form when valid', async () => {
    render(<PrimeUIShipmentUpdateReweighForm {...testProps} />);
    const submitBtn = await screen.findByRole('button', { name: 'Save' });

    expect(submitBtn).toBeDisabled();

    const reweighInput = screen.getByLabelText('Reweigh Weight (lbs)');
    await userEvent.clear(reweighInput);
    await userEvent.type(reweighInput, '123');
    await userEvent.type(screen.getByTestId('remarks'), 'test');

    await waitFor(() => {
      expect(submitBtn).not.toBeDisabled();
    });
    await userEvent.click(submitBtn);

    await waitFor(() => {
      expect(testProps.onSubmit).toHaveBeenCalled();
    });
  });

  it('shows an error message when the reweigh weight is 0', async () => {
    render(<PrimeUIShipmentUpdateReweighForm {...testProps} />);

    const alert = await screen.findByRole('alert');

    expect(alert).toHaveTextContent('Authorized weight must be greater than or equal to 1');
  });

  it('implements the handleClose handler when the Cancel button is clicked', async () => {
    render(<PrimeUIShipmentUpdateReweighForm {...testProps} />);
    const cancelButton = await screen.findByRole('button', { name: 'Cancel' });

    await userEvent.click(cancelButton);

    await waitFor(() => {
      expect(testProps.handleClose).toHaveBeenCalled();
    });
  });
});
