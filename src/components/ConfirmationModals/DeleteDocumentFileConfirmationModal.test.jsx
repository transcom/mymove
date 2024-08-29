import React from 'react';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import { DeleteDocumentFileConfirmationModal } from 'components/ConfirmationModals/DeleteDocumentFileConfirmationModal';

let onClose;
let onSubmit;

const fileInfo = {
  filename: 'test-file',
  bytes: '1212',
  createdAt: '12/01/2024',
};

beforeEach(() => {
  onClose = jest.fn();
  onSubmit = jest.fn();
});

describe('DeleteDocumentFileConfirmationModal', () => {
  it('renders the component', async () => {
    render(<DeleteDocumentFileConfirmationModal submitModal={onSubmit} closeModal={onClose} fileInfo={fileInfo} />);

    expect(
      await screen.findByRole('heading', { level: 3, name: 'Are you sure you want to delete this file?' }),
    ).toBeInTheDocument();
  });

  it('closes the modal when close icon is clicked', async () => {
    render(<DeleteDocumentFileConfirmationModal submitModal={onSubmit} closeModal={onClose} fileInfo={fileInfo} />);

    const closeButton = await screen.findByTestId('modalCloseButton');

    await userEvent.click(closeButton);

    expect(onClose).toHaveBeenCalledTimes(1);
    expect(onSubmit).not.toHaveBeenCalled();
  });

  it('closes the modal when the keep button is clicked', async () => {
    render(<DeleteDocumentFileConfirmationModal submitModal={onSubmit} closeModal={onClose} fileInfo={fileInfo} />);

    const keepButton = await screen.findByRole('button', { name: 'No, keep it' });

    await userEvent.click(keepButton);

    expect(onClose).toHaveBeenCalledTimes(1);
    expect(onSubmit).not.toHaveBeenCalled();
  });

  it('calls the submit function when delete button is clicked', async () => {
    render(<DeleteDocumentFileConfirmationModal submitModal={onSubmit} closeModal={onClose} fileInfo={fileInfo} />);

    const deleteButton = await screen.findByRole('button', { name: 'Yes, delete' });

    await userEvent.click(deleteButton);

    expect(onSubmit).toHaveBeenCalledTimes(1);
    expect(onClose).not.toHaveBeenCalled();
  });
});
