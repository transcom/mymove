import React from 'react';
import { render, waitFor, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import ReviewExpense from './ReviewExpense';

import ppmDocumentStatus from 'constants/ppms';
import { expenseTypes } from 'constants/ppmExpenseTypes';
import { MockProviders } from 'testUtils';

beforeEach(() => {
  jest.clearAllMocks();
});

const defaultProps = {
  ppmShipmentInfo: {
    id: '32ecb311-edbe-4fd4-96ee-bd693113f3f3',
    expectedDepartureDate: '2022-12-02',
    actualMoveDate: '2022-12-06',
    actualPickupPostalCode: '90210',
    actualDestinationPostalCode: '94611',
    miles: 300,
    estimatedWeight: 3000,
    actualWeight: 3500,
  },
  tripNumber: 1,
  ppmNumber: 1,
  showAllFields: false,
  categoryIndex: 1,
};

const expenseRequiredProps = {
  expense: {
    id: '32ecb311-edbe-4fd4-96ee-bd693113f3f3',
    ppmShipmentId: '343bb456-63af-4f76-89bd-7403094a5c4d',
    movingExpenseType: expenseTypes.PACKING_MATERIALS,
    description: 'boxes, tape, bubble wrap',
    paidWithGtcc: false,
    amount: 123456,
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

const rejectedProps = {
  expense: {
    ...expenseRequiredProps.expense,
    reason: 'Rejection reason',
    status: ppmDocumentStatus.REJECTED,
  },
};

describe('ReviewExpenseForm component', () => {
  describe('displays form', () => {
    it('renders blank form on load with defaults', async () => {
      render(<ReviewExpense {...defaultProps} />, { wrapper: MockProviders });

      await waitFor(() => {
        expect(screen.getByRole('heading', { level: 3, name: 'Receipt 1' })).toBeInTheDocument();
      });

      expect(screen.getByText('Expense Type')).toBeInTheDocument();
      expect(screen.getByText('Description')).toBeInTheDocument();
      expect(screen.getByLabelText('Amount')).toBeInstanceOf(HTMLInputElement);

      expect(screen.getByRole('heading', { level: 3, name: `Review #1` })).toBeInTheDocument();

      expect(screen.getByText('Add a review for this')).toBeInTheDocument();

      expect(screen.getByLabelText('Accept')).toBeInstanceOf(HTMLInputElement);
      expect(screen.getByLabelText('Exclude')).toBeInstanceOf(HTMLInputElement);
      expect(screen.getByLabelText('Reject')).toBeInstanceOf(HTMLInputElement);
    });

    it('populates edit form with existing expense values', async () => {
      render(<ReviewExpense {...defaultProps} {...expenseRequiredProps} />, { wrapper: MockProviders });

      await waitFor(() => {
        expect(screen.getByText('Expense Type')).toBeInTheDocument();
      });
      expect(screen.getByText('Packing Materials #1')).toBeInTheDocument();
      expect(screen.getByText('boxes, tape, bubble wrap')).toBeInTheDocument();
      expect(screen.getByLabelText('Amount')).toHaveDisplayValue('1,234.56');
    });

    it('shows SIT fields when expense type is Storage', async () => {
      render(<ReviewExpense {...defaultProps} {...storageProps} />, { wrapper: MockProviders });
      await waitFor(() => {
        expect(screen.getByLabelText('Start date')).toBeInstanceOf(HTMLInputElement);
      });
      expect(screen.getByLabelText('End date')).toBeInstanceOf(HTMLInputElement);
      expect(screen.getByText('Total days in SIT')).toBeInTheDocument();
    });

    it('populates edit form with existing storage values', async () => {
      render(<ReviewExpense {...defaultProps} {...storageProps} />, { wrapper: MockProviders });
      await waitFor(() => {
        expect(screen.getByLabelText('Start date')).toHaveDisplayValue('15 Dec 2022');
      });
      expect(screen.getByLabelText('End date')).toHaveDisplayValue('25 Dec 2022');
    });

    it('correctly displays days in SIT', async () => {
      render(<ReviewExpense {...defaultProps} {...storageProps} />, { wrapper: MockProviders });
      await waitFor(() => {
        expect(screen.getByTestId('days-in-sit')).toHaveTextContent('11');
      });
    });

    it('correctly updates days in SIT', async () => {
      render(<ReviewExpense {...defaultProps} {...storageProps} />, { wrapper: MockProviders });
      const startDateInput = screen.getByLabelText('Start date');
      // clearing a date input throws an error about `children` being set to NaN
      // this replaces the '5' in '15 Dec 2022' with a '7' -> '17 Dec 2022'
      await userEvent.type(startDateInput, '7', {
        initialSelectionStart: 1,
        initialSelectionEnd: 2,
      });
      await waitFor(() => {
        expect(screen.getByTestId('days-in-sit')).toHaveTextContent('9');
      });
    });

    it('populates edit form with existing status and reason', async () => {
      render(<ReviewExpense {...defaultProps} {...rejectedProps} />, { wrapper: MockProviders });
      await waitFor(() => {
        expect(screen.getByLabelText('Reject')).toBeChecked();
      });
      expect(screen.getByText('Rejection reason')).toBeInTheDocument();
      expect(screen.getByText('484 characters')).toBeInTheDocument();
    });
  });
});
