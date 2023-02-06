import React from 'react';
import { render, waitFor, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import FinancialReviewModal from './FinancialReviewModal';

describe('FinancialReviewModal', () => {
  it('calls onSubmit prop on approval with form values when validations pass', async () => {
    const mockOnSubmit = jest.fn();
    render(<FinancialReviewModal onSubmit={mockOnSubmit} onClose={() => {}} />);
    const flagForReview = screen.getByLabelText('Yes');
    const remarksInput = screen.getByLabelText('Remarks for financial office');
    const submitBtn = screen.getByRole('button', { name: 'Save' });

    await userEvent.click(flagForReview);
    await userEvent.type(remarksInput, 'Because I said so...');
    await userEvent.click(submitBtn);

    await waitFor(() => {
      expect(mockOnSubmit).toHaveBeenCalled();
    });
  });

  it('does not allow submission without remarks', async () => {
    render(<FinancialReviewModal onSubmit={jest.fn()} onClose={() => {}} />);
    const flagForReview = screen.getByLabelText('Yes');
    const submitBtn = screen.getByText('Save');

    await userEvent.click(flagForReview);

    await waitFor(() => {
      expect(submitBtn).toBeDisabled();
    });
  });

  it('calls onclose prop on modal close', async () => {
    const mockClose = jest.fn();
    render(<FinancialReviewModal onSubmit={() => {}} onClose={mockClose} />);
    const closeBtn = screen.getByText('Cancel');

    await userEvent.click(closeBtn);

    await waitFor(() => {
      expect(mockClose).toHaveBeenCalled();
    });
  });
});
