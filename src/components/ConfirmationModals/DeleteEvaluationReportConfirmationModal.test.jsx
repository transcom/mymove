import React from 'react';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import { DeleteEvaluationReportConfirmationModal } from 'components/ConfirmationModals/DeleteEvaluationReportConfirmationModal';

let onClose;
let onSubmit;
beforeEach(() => {
  onClose = jest.fn();
  onSubmit = jest.fn();
});

describe('DeleteEvaluationReportConfirmationModal', () => {
  it('renders the component', async () => {
    render(<DeleteEvaluationReportConfirmationModal submitModal={onSubmit} closeModal={onClose} />);

    expect(
      await screen.findByRole('heading', { level: 3, name: 'Are you sure you want to cancel this report?' }),
    ).toBeInTheDocument();
  });

  it('closes the modal when close icon is clicked', async () => {
    render(<DeleteEvaluationReportConfirmationModal submitModal={onSubmit} closeModal={onClose} />);

    const closeButton = await screen.findByTestId('modalCloseButton');

    await userEvent.click(closeButton);

    expect(onClose).toHaveBeenCalledTimes(1);
    expect(onSubmit).not.toHaveBeenCalled();
  });

  it('closes the modal when the keep button is clicked', async () => {
    render(<DeleteEvaluationReportConfirmationModal submitModal={onSubmit} closeModal={onClose} />);

    const keepButton = await screen.findByRole('button', { name: 'No, keep it' });

    await userEvent.click(keepButton);

    expect(onClose).toHaveBeenCalledTimes(1);
    expect(onSubmit).not.toHaveBeenCalled();
  });

  it('calls the submit function when delete button is clicked', async () => {
    render(<DeleteEvaluationReportConfirmationModal submitModal={onSubmit} closeModal={onClose} />);

    const deleteButton = await screen.findByRole('button', { name: 'Yes, Cancel' });

    await userEvent.click(deleteButton);

    expect(onSubmit).toHaveBeenCalledTimes(1);
    expect(onClose).not.toHaveBeenCalled();
  });
});
