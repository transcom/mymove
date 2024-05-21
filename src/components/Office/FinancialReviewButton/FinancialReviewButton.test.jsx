import React from 'react';
import { render, waitFor, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import FinancialReviewButton from './FinancialReviewButton';

describe('FinancialReviewButton', () => {
  it('calls the onClick function when clicked', async () => {
    const mockOnClick = jest.fn();
    render(<FinancialReviewButton onClick={mockOnClick} />);
    const submitBtn = screen.getByText('Flag move for financial review');

    await userEvent.click(submitBtn);

    await waitFor(() => {
      expect(mockOnClick).toHaveBeenCalled();
    });
  });

  it('displays a tag when a review has been requested', async () => {
    const mockOnClick = jest.fn();
    render(<FinancialReviewButton onClick={mockOnClick} reviewRequested />);
    const tag = screen.getByTestId('tag');

    expect(tag).toHaveTextContent('Flagged for financial review');
  });

  it('disables the button when the move is locked', async () => {
    const mockOnClick = jest.fn();
    const isMoveLocked = true;
    render(<FinancialReviewButton onClick={mockOnClick} isMoveLocked={isMoveLocked} />);
    const submitBtn = screen.getByText('Flag move for financial review');

    expect(submitBtn).toBeDisabled();
  });
});
