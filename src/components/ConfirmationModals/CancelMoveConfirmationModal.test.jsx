import React from 'react';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import { CancelMoveConfirmationModal } from 'components/ConfirmationModals/CancelMoveConfirmationModal';

let onClose;
let onSubmit;
beforeEach(() => {
  onClose = jest.fn();
  onSubmit = jest.fn();
});

describe('DeleteCustomerSupportRemarkConfirmationModal', () => {
  const moveId = '123456';

  it('renders the component', async () => {
    render(<CancelMoveConfirmationModal onSubmit={onSubmit} onClose={onClose} moveId={moveId} />);

    expect(await screen.findByRole('heading', { level: 3, name: 'Cancel this move?' })).toBeInTheDocument();
  });

  it('closes the modal when close icon is clicked', async () => {
    render(<CancelMoveConfirmationModal onSubmit={onSubmit} onClose={onClose} moveId={moveId} />);

    const closeButton = await screen.findByTestId('modalCloseButton');

    await userEvent.click(closeButton);

    expect(onClose).toHaveBeenCalledTimes(1);
  });

  it('closes the modal when the keep button is clicked', async () => {
    render(<CancelMoveConfirmationModal onSubmit={onSubmit} onClose={onClose} moveId={moveId} />);

    const keepButton = await screen.findByRole('button', { name: 'No, Keep it' });

    await userEvent.click(keepButton);

    expect(onClose).toHaveBeenCalledTimes(1);
  });

  it('calls the submit function when delete button is clicked', async () => {
    render(<CancelMoveConfirmationModal onSubmit={onSubmit} onClose={onClose} moveId={moveId} />);

    const deleteButton = await screen.findByRole('button', { name: 'Yes, Cancel' });

    await userEvent.click(deleteButton);

    expect(onSubmit).toHaveBeenCalledWith(moveId);
    expect(onSubmit).toHaveBeenCalledTimes(1);
  });
});
