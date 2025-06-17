import React from 'react';
import { render, waitFor, screen, act } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import moment from 'moment';

import ConvertSITToCustomerExpenseModal from './ConvertSITToCustomerExpenseModal';

import { utcDateFormat } from 'shared/dates';

const defaultValues = {
  sitStatus: {
    totalDaysRemaining: 210,
    totalSITDaysUsed: 60,
    currentSIT: {
      daysInSIT: 60,
      sitEntryDate: moment().subtract(60, 'days').format(utcDateFormat),
    },
  },
  shipment: {
    sitDaysAllowance: 270,
  },
};

describe('ConvertSITToCustomerExpenseModal', () => {
  it('calls onSubmit prop on approval with form values when validations pass', async () => {
    const mockOnSubmit = jest.fn();
    await render(<ConvertSITToCustomerExpenseModal onSubmit={mockOnSubmit} onClose={() => {}} {...defaultValues} />);
    const remarksInput = screen.getByLabelText('Remarks *');
    const submitBtn = screen.getByRole('button', { name: 'Save' });

    await act(() => userEvent.type(remarksInput, 'Approved!'));
    await act(() => userEvent.click(submitBtn));

    await waitFor(() => {
      expect(mockOnSubmit).toHaveBeenCalled();
      expect(mockOnSubmit).toHaveBeenCalledWith(true, 'Approved!');
    });
  });

  it('does not allow submission when office remarks is empty', async () => {
    const mockOnSubmit = jest.fn();
    await render(<ConvertSITToCustomerExpenseModal onSubmit={mockOnSubmit} onClose={() => {}} {...defaultValues} />);
    const remarksInput = screen.getByLabelText('Remarks *');
    const submitBtn = screen.getByRole('button', { name: 'Save' });

    await act(() => userEvent.clear(remarksInput));
    await waitFor(() => {
      expect(submitBtn).toBeDisabled();
    });
  });

  it('calls onclose prop on modal close', async () => {
    const mockClose = jest.fn();
    await render(<ConvertSITToCustomerExpenseModal onSubmit={() => {}} onClose={mockClose} {...defaultValues} />);
    const closeBtn = screen.getByRole('button', { name: 'Cancel' });

    await act(() => userEvent.click(closeBtn));

    await waitFor(() => {
      expect(mockClose).toHaveBeenCalled();
    });
  });

  it('renders the summary SIT component', async () => {
    await render(<ConvertSITToCustomerExpenseModal onSubmit={jest.fn()} onClose={jest.fn()} {...defaultValues} />);

    await waitFor(() => {
      expect(screen.getByText('SIT (STORAGE IN TRANSIT)')).toBeInTheDocument();
    });
  });

  it('renders asterisks for required fields', async () => {
    await render(<ConvertSITToCustomerExpenseModal onSubmit={jest.fn()} onClose={jest.fn()} {...defaultValues} />);

    await waitFor(() => {
      expect(screen.getByText('SIT (STORAGE IN TRANSIT)')).toBeInTheDocument();
    });

    expect(document.querySelector('#reqAsteriskMsg')).toHaveTextContent('Fields marked with * are required.');
    expect(screen.getByLabelText('Remarks *')).toBeInTheDocument();
  });
});
