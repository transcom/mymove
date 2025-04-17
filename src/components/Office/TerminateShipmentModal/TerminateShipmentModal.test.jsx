import React from 'react';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import TerminateShipmentModal from './TerminateShipmentModal';

describe('TerminateShipmentModal', () => {
  const mockOnClose = jest.fn();
  const mockOnSubmit = jest.fn();
  const shipmentID = 'test-shipment-id';
  const shipmentLocator = 'ABC123';

  beforeEach(() => {
    render(
      <TerminateShipmentModal
        isOpen
        onClose={mockOnClose}
        onSubmit={mockOnSubmit}
        shipmentID={shipmentID}
        shipmentLocator={shipmentLocator}
      />,
    );
  });

  afterEach(() => {
    jest.clearAllMocks();
  });

  it('renders the modal with title and shipment locator hint', () => {
    expect(screen.getByRole('heading', { name: 'Shipment termination' })).toBeInTheDocument();
    expect(screen.getByText(shipmentLocator)).toBeInTheDocument();
    expect(screen.getByTestId('terminationComments')).toBeInTheDocument();
  });

  it('renders the prefixed text field', () => {
    expect(screen.getByText('TERMINATED FOR CAUSE:')).toBeInTheDocument();
  });

  it('calls onClose when Cancel is clicked', async () => {
    const cancelButton = screen.getByTestId('modalBackBtn');
    await userEvent.click(cancelButton);
    expect(mockOnClose).toHaveBeenCalled();
  });

  it('calls onSubmit when Terminate is clicked and form is valid', async () => {
    const terminateButton = screen.getByTestId('modalSubmitBtn');
    await userEvent.type(screen.getByLabelText(/Termination reason/), 'get in the choppuh');
    expect(terminateButton).toBeEnabled();
    await userEvent.click(terminateButton);
    expect(mockOnSubmit).toHaveBeenCalledWith(shipmentID, {
      terminationComments: 'get in the choppuh',
    });
  });

  it('disables the Terminate button when form is invalid or submitting', () => {
    const terminateButton = screen.getByTestId('modalSubmitBtn');
    expect(terminateButton).toBeDisabled();
  });
});
