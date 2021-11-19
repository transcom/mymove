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

    userEvent.click(flagForReview);
    userEvent.type(remarksInput, 'Because I said so...');
    userEvent.click(submitBtn);

    await waitFor(() => {
      expect(mockOnSubmit).toHaveBeenCalled();
    });
  });

  it('displays initial remarks', async () => {
    const remarks = 'Initial remarks';
    render(<FinancialReviewModal remarks={remarks} onSubmit={() => {}} onClose={() => {}} />);
    const flagForReview = screen.getByLabelText('Yes');
    const remarksInput = screen.getByLabelText('Remarks for financial office');

    await waitFor(() => {
      expect(flagForReview).toBeChecked();
      expect(remarksInput).toHaveValue(remarks);
    });
  });

  it('calls onclose prop on modal close', async () => {
    const mockClose = jest.fn();
    render(<FinancialReviewModal onSubmit={() => {}} onClose={mockClose} />);
    const closeBtn = screen.getByText('Cancel');

    userEvent.click(closeBtn);

    await waitFor(() => {
      expect(mockClose).toHaveBeenCalled();
    });
  });

  describe('Radio is yes', () => {
    it('does not allow submission without remarks', async () => {
      render(<FinancialReviewModal onSubmit={jest.fn()} onClose={() => {}} />);
      const flagForReview = screen.getByLabelText('Yes');
      const submitBtn = screen.getByText('Save');

      userEvent.click(flagForReview);

      await waitFor(() => {
        expect(submitBtn).toBeDisabled();
      });
    });
  });

  describe('Radio is no', () => {
    it('allows submission without remarks', async () => {
      render(<FinancialReviewModal onSubmit={jest.fn()} onClose={() => {}} />);
      const flagForReview = screen.getByLabelText('No');
      const submitBtn = screen.getByText('Save');

      userEvent.click(flagForReview);

      await waitFor(() => {
        expect(submitBtn).not.toBeDisabled();
      });
    });

    it('preserves remarks after No is selected', async () => {
      render(<FinancialReviewModal onSubmit={() => {}} onClose={() => {}} />);
      const remarksInput = screen.getByLabelText('Remarks for financial office');
      const doesNotHaveRemarks = screen.getByLabelText('No');

      userEvent.type(remarksInput, 'Test remark');
      userEvent.click(doesNotHaveRemarks);

      await waitFor(() => {
        expect(remarksInput).toHaveValue('Test remark');
      });
    });
  });
});
