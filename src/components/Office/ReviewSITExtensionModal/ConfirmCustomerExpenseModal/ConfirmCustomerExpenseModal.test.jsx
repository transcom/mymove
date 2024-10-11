import React from 'react';
import { render, screen, fireEvent } from '@testing-library/react';

import '@testing-library/jest-dom/extend-expect';
import ConfirmCustomerExpenseModal from './ConfirmCustomerExpenseModal';

describe('ConfirmCustomerExpenseModal', () => {
  let setShowConfirmModal;
  let setValues;
  let values;

  beforeEach(() => {
    setShowConfirmModal = jest.fn();
    setValues = jest.fn();
    values = { convertToCustomerExpense: false };
  });

  it('renders the modal with correct content', () => {
    render(
      <ConfirmCustomerExpenseModal setShowConfirmModal={setShowConfirmModal} values={values} setValues={setValues} />,
    );

    expect(screen.getByText('Convert to Customer Expense')).toBeInTheDocument();
    expect(screen.getByText('Are you sure that you would like to convert to Customer Expense?')).toBeInTheDocument();
    // Check if both "Yes" and "No" buttons are rendered
    expect(screen.getByTestId('convertToCustomerExpenseConfirmationYes')).toBeInTheDocument();
    expect(screen.getByTestId('convertToCustomerExpenseConfirmationNo')).toBeInTheDocument();
  });

  it('handles "Yes" button click by setting convertToCustomerExpense to true and closing the modal', () => {
    render(
      <ConfirmCustomerExpenseModal setShowConfirmModal={setShowConfirmModal} values={values} setValues={setValues} />,
    );

    const yesButton = screen.getByTestId('convertToCustomerExpenseConfirmationYes');
    fireEvent.click(yesButton);
    expect(setValues).toHaveBeenCalledWith({ ...values, convertToCustomerExpense: true });
    expect(setShowConfirmModal).toHaveBeenCalledWith(false);
  });

  it('handles "No" button click by setting convertToCustomerExpense to false and closing the modal', () => {
    render(
      <ConfirmCustomerExpenseModal setShowConfirmModal={setShowConfirmModal} values={values} setValues={setValues} />,
    );

    const noButton = screen.getByTestId('convertToCustomerExpenseConfirmationNo');
    fireEvent.click(noButton);
    expect(setValues).toHaveBeenCalledWith({ ...values, convertToCustomerExpense: false });
    expect(setShowConfirmModal).toHaveBeenCalledWith(false);
  });

  it('handles closing the modal via the close button', () => {
    render(
      <ConfirmCustomerExpenseModal setShowConfirmModal={setShowConfirmModal} values={values} setValues={setValues} />,
    );

    const closeButton = screen.getByRole('button', { name: /close/i });
    fireEvent.click(closeButton);
    expect(setShowConfirmModal).toHaveBeenCalledWith(false);
  });
});
