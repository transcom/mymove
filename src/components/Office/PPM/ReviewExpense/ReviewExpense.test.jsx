import React from 'react';
import { render, waitFor, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import ReviewExpense from './ReviewExpense';

import { expenseTypes } from 'constants/ppmExpenseTypes';

beforeEach(() => {
  jest.clearAllMocks();
});

const defaultProps = {
  ppmShipment: {
    id: '32ecb311-edbe-4fd4-96ee-bd693113f3f3',
    actualPickupPostalCode: '90210',
    actualMoveDate: '2022-04-30',
    actualDestinationPostalCode: '94611',
    hasReceivedAdvance: true,
    advanceAmountReceived: 60000,
  },
  expenseNumber: 1,
  ppmNumber: 1,
};

const expenseRequiredProps = {
  expense: {
    id: '32ecb311-edbe-4fd4-96ee-bd693113f3f3',
    ppmShipmentId: '343bb456-63af-4f76-89bd-7403094a5c4d',
    movingExpenseType: expenseTypes.PACKING_MATERIALS,
    description: 'boxes, tape, bubble wrap',
    paidWithGtcc: false,
    amount: 12345,
  },
};

const storageProps = {
  expense: {
    ...expenseRequiredProps.expense,
    movingExpenseType: expenseTypes.STORAGE,
    sitStartDate: '2022-12-15',
    sitEndDate: '2022-12-25',
  },
};

describe('ReviewExpenseForm component', () => {
  describe('displays form', () => {
    it('renders blank form on load with defaults', async () => {
      render(<ReviewExpense {...defaultProps} />);

      await waitFor(() => {
        expect(screen.getByRole('heading', { level: 3, name: 'Receipt 1' })).toBeInTheDocument();
      });

      expect(screen.getByText('Expense type')).toBeInTheDocument();
      expect(screen.getByText('Description')).toBeInTheDocument();
      expect(screen.getByLabelText('Amount')).toBeInstanceOf(HTMLInputElement);

      expect(screen.getByRole('heading', { level: 3, name: 'Review receipt 1' })).toBeInTheDocument();

      expect(screen.getByText('Add a review for this receipt')).toBeInTheDocument();

      expect(screen.getByLabelText('Accept')).toBeInstanceOf(HTMLInputElement);
      expect(screen.getByLabelText('Exclude')).toBeInstanceOf(HTMLInputElement);
      expect(screen.getByLabelText('Reject')).toBeInstanceOf(HTMLInputElement);
    });

    it('populates edit form with existing expense values', async () => {
      render(<ReviewExpense {...defaultProps} {...expenseRequiredProps} />);

      await waitFor(() => {
        expect(screen.getByText('Packing materials')).toBeInTheDocument();
      });
      expect(screen.getByText('boxes, tape, bubble wrap')).toBeInTheDocument();

      expect(screen.getByLabelText('Amount')).toHaveDisplayValue('123.45');
    });

    it('shows SIT fields when expense type is Storage', async () => {
      render(<ReviewExpense {...defaultProps} {...storageProps} />);
      await waitFor(() => {
        expect(screen.getByLabelText('Start date')).toBeInstanceOf(HTMLInputElement);
      });
      expect(screen.getByLabelText('End date')).toBeInstanceOf(HTMLInputElement);
      expect(screen.getByText('Total days in SIT')).toBeInTheDocument();
    });

    it('populates edit form with existing storage values', async () => {
      render(<ReviewExpense {...defaultProps} {...storageProps} />);
      await waitFor(() => {
        expect(screen.getByLabelText('Start date')).toHaveDisplayValue('15 Dec 2022');
      });
      expect(screen.getByLabelText('End date')).toHaveDisplayValue('25 Dec 2022');
    });

    it('correctly displays days in SIT', async () => {
      render(<ReviewExpense {...defaultProps} {...storageProps} />);
      await waitFor(() => {
        expect(screen.getByTestId('days-in-sit')).toHaveTextContent('10');
      });
    });

    it('correctly updates days in SIT', async () => {
      render(<ReviewExpense {...defaultProps} {...storageProps} />);
      const startDateInput = screen.getByLabelText('Start date');
      // clearing a date input throws an error about `children` being set to NaN
      // this replaces the '5' in '15 Dec 2022' with a '7' -> '17 Dec 2022'
      await userEvent.type(startDateInput, '7', {
        initialSelectionStart: 1,
        initialSelectionEnd: 2,
      });
      await waitFor(() => {
        expect(screen.getByTestId('days-in-sit')).toHaveTextContent('8');
      });
    });
  });
});
