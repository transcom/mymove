import React from 'react';
import { render, screen, fireEvent } from '@testing-library/react';

import '@testing-library/jest-dom';
import ConfirmCustomerExpenseModal from './ConfirmCustomerExpenseModal';

describe('ConfirmCustomerExpenseModal', () => {
  const setShowConfirmModalMock = jest.fn();
  const setValuesMock = jest.fn();

  const defaultProps = {
    setShowConfirmModal: setShowConfirmModalMock,
    values: { convertToCustomerExpense: false },
    setValues: setValuesMock,
  };

  it('renders without crashing', () => {
    render(<ConfirmCustomerExpenseModal {...defaultProps} />);
    expect(screen.getByText('Convert to Customer Expense')).toBeInTheDocument();
  });

  it('calls handleConfirmYes on "Yes" button click', () => {
    render(<ConfirmCustomerExpenseModal {...defaultProps} />);
    fireEvent.click(screen.getByTestId('convertToCustomerExpenseConfirmationYes'));
    expect(setValuesMock).toHaveBeenCalledWith({
      ...defaultProps.values,
      convertToCustomerExpense: true,
    });
    expect(setShowConfirmModalMock).toHaveBeenCalledWith(false);
  });

  it('calls handleConfirmNo on "No" button click', () => {
    render(<ConfirmCustomerExpenseModal {...defaultProps} />);
    fireEvent.click(screen.getByTestId('convertToCustomerExpenseConfirmationNo'));
    expect(setValuesMock).toHaveBeenCalledWith({
      ...defaultProps.values,
      convertToCustomerExpense: false,
    });
    expect(setShowConfirmModalMock).toHaveBeenCalledWith(false);
  });
});
