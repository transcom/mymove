import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import '@testing-library/jest-dom/extend-expect';
import { ViolationAppealModal } from './ViolationAppealModal';

describe('ViolationAppealModal', () => {
  const mockOnClose = jest.fn();
  const mockOnSubmit = jest.fn();

  beforeEach(() => {
    render(<ViolationAppealModal onClose={mockOnClose} onSubmit={mockOnSubmit} />);
  });

  afterEach(() => {
    jest.clearAllMocks();
  });

  test('renders correctly with all fields and buttons', () => {
    expect(screen.getByText('Leave Appeal Decision')).toBeInTheDocument();
    expect(screen.getByLabelText('Remarks')).toBeInTheDocument();
    expect(screen.getByLabelText('Sustained')).toBeInTheDocument();
    expect(screen.getByLabelText('Rejected')).toBeInTheDocument();
    expect(screen.getByRole('button', { name: 'Save' })).toBeDisabled();
    expect(screen.getByRole('button', { name: 'Cancel' })).toBeInTheDocument();
  });

  test('displays validation messages when form is submitted with empty fields', async () => {
    expect(screen.getByRole('button', { name: 'Save' })).toBeDisabled();

    // click the input but don't do anything
    userEvent.click(screen.getByLabelText('Remarks'));
    userEvent.click(screen.getByTestId('sustainedRadio'));

    await waitFor(() => {
      expect(screen.getByText('Remarks are required')).toBeInTheDocument();
    });

    expect(screen.getByRole('button', { name: 'Save' })).toBeDisabled();
  });

  test('the form successfully submits when all required fields are filled out', async () => {
    await userEvent.type(screen.getByLabelText('Remarks'), 'These are my remarks');
    await userEvent.click(screen.getByTestId('sustainedRadio'));

    expect(screen.getByRole('button', { name: 'Save' })).toBeEnabled();

    userEvent.click(screen.getByRole('button', { name: /Save/i }));

    await waitFor(() => {
      expect(mockOnSubmit).toHaveBeenCalledWith(
        {
          remarks: 'These are my remarks',
          appealStatus: 'sustained',
        },
        expect.anything(),
      );
    });
  });

  test('Cancel button triggers onClose callback', async () => {
    await userEvent.click(screen.getByTestId('modalCancelButton'));
    expect(mockOnClose).toHaveBeenCalled();
  });
});
