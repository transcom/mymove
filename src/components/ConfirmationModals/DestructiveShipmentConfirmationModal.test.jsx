import React from 'react';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import { DestructiveShipmentConfirmationModal } from 'components/ConfirmationModals/DestructiveShipmentConfirmationModal';

let onClose;
let onSubmit;
beforeEach(() => {
  onClose = jest.fn();
  onSubmit = jest.fn();
});

describe('DestructiveShipmentConfirmationModal', () => {
  const shipmentID = '123456';

  it('renders the component', async () => {
    render(<DestructiveShipmentConfirmationModal onSubmit={onSubmit} onClose={onClose} shipmentID={shipmentID} />);

    expect(await screen.findByRole('heading', { level: 3, name: 'Are you sure?' })).toBeInTheDocument();
  });

  it('closes the modal when close icon is clicked', async () => {
    render(<DestructiveShipmentConfirmationModal onSubmit={onSubmit} onClose={onClose} shipmentID={shipmentID} />);

    const closeButton = await screen.findByTestId('modalCloseButton');

    userEvent.click(closeButton);

    expect(onClose).toHaveBeenCalledTimes(1);
  });

  it('closes the modal when the keep button is clicked', async () => {
    render(<DestructiveShipmentConfirmationModal onSubmit={onSubmit} onClose={onClose} shipmentID={shipmentID} />);

    const keepButton = await screen.findByRole('button', { name: 'Keep shipment' });

    userEvent.click(keepButton);

    expect(onClose).toHaveBeenCalledTimes(1);
  });

  it('calls the submit function when delete button is clicked', async () => {
    render(<DestructiveShipmentConfirmationModal onSubmit={onSubmit} onClose={onClose} shipmentID={shipmentID} />);

    const deleteButton = await screen.findByRole('button', { name: 'Delete shipment' });

    userEvent.click(deleteButton);

    expect(onSubmit).toHaveBeenCalledWith(shipmentID);
    expect(onSubmit).toHaveBeenCalledTimes(1);
  });
});
