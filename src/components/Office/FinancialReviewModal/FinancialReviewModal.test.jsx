import React from 'react';
import { render, waitFor, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import FinancialReviewModal from './FinancialReviewModal';

describe('FinancialReviewModal', () => {
  it('calls onSubmit prop on approval with form values when validations pass', async () => {
    const mockOnSubmit = jest.fn();
    render(<FinancialReviewModal onSubmit={mockOnSubmit} onClose={() => {}} />);
    const reviewCheckbox = screen.getByTestId('reviewCheckbox');
    const remarksInput = screen.getByLabelText('Remarks');
    const submitBtn = screen.getByTestId('modalSaveButton');

    userEvent.click(reviewCheckbox);
    userEvent.type(remarksInput, 'Because I said so...');
    userEvent.click(submitBtn);

    await waitFor(() => {
      expect(mockOnSubmit).toHaveBeenCalled();
      expect(mockOnSubmit).toHaveBeenCalledWith({
        reviewCheckbox: true,
        officeRemarks: 'Because I said so...!',
      });
    });
  });

  it('does not allow submission without remarks', async () => {
    const mockOnSubmit = jest.fn();
    render(<FinancialReviewModal onSubmit={mockOnSubmit} onClose={() => {}} />);
    const reviewCheckbox = screen.getByTestId('reviewCheckbox');
    const submitBtn = screen.getByTestId('modalSaveButton');

    userEvent.click(reviewCheckbox);

    await waitFor(() => {
      expect(submitBtn).toBeDisabled();
    });
  });

  it('calls onclose prop on modal close', async () => {
    const mockClose = jest.fn();
    render(<FinancialReviewModal onSubmit={() => {}} onClose={mockClose} />);
    const closeBtn = screen.getByTestId('modalCancelButton');

    userEvent.click(closeBtn);

    await waitFor(() => {
      expect(mockClose).toHaveBeenCalled();
    });
  });
});
