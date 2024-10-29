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

describe('CancelMoveConfirmationModal', () => {
  const moveID = '123456';

  it('renders the component', async () => {
    render(<CancelMoveConfirmationModal onSubmit={onSubmit} onClose={onClose} moveId={moveID} />);

    expect(await screen.findByRole('heading', { level: 3, name: 'Are you sure?' })).toBeInTheDocument();
  });

  it('closes the modal when close icon is clicked', async () => {
    render(<CancelMoveConfirmationModal onSubmit={onSubmit} onClose={onClose} moveId={moveID} />);

    const closeButton = await screen.findByTestId('modalCloseButton');

    await userEvent.click(closeButton);

    expect(onClose).toHaveBeenCalledTimes(1);
  });

  it('closes the modal when the keep button is clicked', async () => {
    render(<CancelMoveConfirmationModal onSubmit={onSubmit} onClose={onClose} moveId={moveID} />);

    const keepButton = await screen.findByRole('button', { name: 'Keep move' });

    await userEvent.click(keepButton);

    expect(onClose).toHaveBeenCalledTimes(1);
  });

  it('calls the submit function when cancel button is clicked', async () => {
    render(<CancelMoveConfirmationModal onSubmit={onSubmit} onClose={onClose} moveId={moveID} />);

    const cancelButton = await screen.findByRole('button', { name: 'Cancel move' });

    await userEvent.click(cancelButton);

    expect(onSubmit).toHaveBeenCalledWith(moveID);
    expect(onSubmit).toHaveBeenCalledTimes(1);
  });
});
