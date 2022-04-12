import React from 'react';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import { SubmitMoveConfirmationModal } from 'components/Office/SubmitMoveConfirmationModal/SubmitMoveConfirmationModal';

let onClose;
let onSubmit;
beforeEach(() => {
  onClose = jest.fn();
  onSubmit = jest.fn();
});

describe('SubmitMoveConfirmationModal', () => {
  it('renders the component', () => {
    render(<SubmitMoveConfirmationModal onSubmit={onSubmit} onClose={onClose} />);
    expect(screen.getByRole('heading', { level: 2, name: 'Are you sure?' })).toBeInTheDocument();
    expect(screen.getByText("You can't make changes after you submit the move.")).toBeInTheDocument();
  });

  it('closes the modal when close icon is clicked', () => {
    render(<SubmitMoveConfirmationModal onSubmit={onSubmit} onClose={onClose} />);
    userEvent.click(screen.getByTestId('modalCloseButton'));

    expect(onClose).toHaveBeenCalled();
  });

  it('closes the modal when the cancel button is clicked', () => {
    render(<SubmitMoveConfirmationModal onSubmit={onSubmit} onClose={onClose} />);
    userEvent.click(screen.getByTestId('modalCancelButton'));

    expect(onClose).toHaveBeenCalled();
  });

  it('calls the submit function when submit button is clicked', () => {
    render(<SubmitMoveConfirmationModal onSubmit={onSubmit} onClose={onClose} />);
    userEvent.click(screen.getByRole('button', { name: 'Yes, submit' }));

    expect(onSubmit).toHaveBeenCalled();
  });

  it('accepts an optional bodyText prop', () => {
    render(<SubmitMoveConfirmationModal onSubmit={onSubmit} onClose={onClose} bodyText="test text goes here" />);
    expect(screen.getByRole('heading', { level: 2, name: 'Are you sure?' })).toBeInTheDocument();
    expect(screen.queryByText("You can't make changes after you submit the move.")).not.toBeInTheDocument();
    expect(screen.getByText('test text goes here')).toBeInTheDocument();
  });
});
