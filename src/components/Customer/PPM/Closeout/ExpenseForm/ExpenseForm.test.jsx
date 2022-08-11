import React from 'react';
import { render, waitFor, screen, within } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import ExpenseForm from 'components/Customer/PPM/Closeout/ExpenseForm/ExpenseForm';
import { DocumentAndImageUploadInstructions } from 'content/uploads';

beforeEach(() => {
  jest.clearAllMocks();
});

const defaultProps = {
  expense: {
    id: '32ecb311-edbe-4fd4-96ee-bd693113f3f3',
    ppmShipmentId: '343bb456-63af-4f76-89bd-7403094a5c4d',
    expenseType: 'packing_materials',
  },
  receiptNumber: '1',
  onCreateUpload: jest.fn(),
  onUploadComplete: jest.fn(),
  onUploadDelete: jest.fn(),
  onBack: jest.fn(),
  onSubmit: jest.fn(),
};

const missingReceiptProps = {
  expense: {
    id: '32ecb311-edbe-4fd4-96ee-bd693113f3f3',
    ppmShipmentId: '343bb456-63af-4f76-89bd-7403094a5c4d',
    expenseType: 'packing_materials',
    description: 'bubble wrap',
    missingReceipt: true,
  },
  receiptNumber: '1',
  onCreateUpload: jest.fn(),
  onUploadComplete: jest.fn(),
  onUploadDelete: jest.fn(),
  onBack: jest.fn(),
  onSubmit: jest.fn(),
};

const expenseRequiredProps = {
  expense: {
    id: '32ecb311-edbe-4fd4-96ee-bd693113f3f3',
    ppmShipmentId: '343bb456-63af-4f76-89bd-7403094a5c4d',
    expenseType: 'packing_materials',
    description: 'bubble wrap',
    missingReceipt: false,
    paidWithGTCC: false,
    amount: 600,
    receiptDocument: {
      uploads: [
        {
          id: 'db4713ae-6087-4330-8b0d-926b3d65c454',
          created_at: '2022-06-10T12:59:30.000Z',
          bytes: 204800,
          url: 'some/path/to/',
          filename: 'expenseReceipt.pdf',
          content_type: 'application/pdf',
        },
      ],
    },
  },
};

const sitExpenseProps = {
  expense: {
    id: '32ecb311-edbe-4fd4-96ee-bd693113f3f3',
    ppmShipmentId: '343bb456-63af-4f76-89bd-7403094a5c4d',
    expenseType: 'storage',
    description: '10x10 storage pod',
    missingReceipt: false,
    paidWithGTCC: false,
    amount: 1600,
    sitStartDate: '2022-09-24',
    sitEndDate: '2022-12-26',
    receiptDocument: {
      uploads: [
        {
          id: 'db4713ae-6087-4330-8b0d-926b3d65c454',
          created_at: '2022-08-10T12:59:30.000Z',
          bytes: 204800,
          url: 'some/path/to/',
          filename: 'uhaulReceipt.pdf',
          content_type: 'application/pdf',
        },
      ],
    },
  },
};

describe('ExpenseForm component', () => {
  describe('displays form', () => {
    it('renders blank form on load with defaults', async () => {
      render(<ExpenseForm {...defaultProps} />);

      await waitFor(() => {
        expect(screen.getByRole('heading', { level: 2, name: 'Receipt 1' })).toBeInTheDocument();
      });

      expect(screen.getByLabelText('Select type')).toBeInstanceOf(HTMLSelectElement);
      expect(screen.getByRole('heading', { level: 3, name: 'Description' })).toBeInTheDocument();
      expect(screen.getByLabelText('What did you buy?')).toBeInstanceOf(HTMLInputElement);
      expect(screen.getByText('Add a brief description of the expense.')).toBeInTheDocument();
      expect(screen.getByText('Did you pay with your GTCC (Government Travel Charge Card)?')).toBeInTheDocument();
      expect(screen.getByLabelText('Yes')).toBeInstanceOf(HTMLInputElement);
      expect(screen.getByLabelText('No')).toBeInstanceOf(HTMLInputElement);

      expect(screen.getByRole('heading', { level: 3, name: 'Amount' })).toBeInTheDocument();
      expect(screen.getByLabelText('Amount')).toBeInstanceOf(HTMLInputElement);
      expect(
        screen.getByText(
          "Enter the total unit price for all items on the receipt that you're claiming as part of your PPM moving expenses.",
        ),
      ).toBeInTheDocument();
      const missingReceipt = screen.getByLabelText("I don't have this receipt");
      expect(missingReceipt).toBeInstanceOf(HTMLInputElement);
      expect(missingReceipt).not.toBeChecked();
      expect(screen.getByText('Upload receipt')).toBeInstanceOf(HTMLLabelElement);
      const uploadFileTypeHints = screen.getAllByText(DocumentAndImageUploadInstructions);
      expect(uploadFileTypeHints[0]).toBeInTheDocument();
      expect(screen.queryByRole('heading', { level: 3, name: 'Dates' })).not.toBeInTheDocument();

      expect(screen.getByRole('button', { name: 'Finish Later' })).toBeInTheDocument();
      expect(screen.getByRole('button', { name: 'Save & Continue' })).toBeInTheDocument();
    });

    it('populates edit form with existing expense values', async () => {
      render(<ExpenseForm {...defaultProps} {...expenseRequiredProps} />);

      await waitFor(() => {
        expect(screen.getByLabelText('What did you buy?')).toHaveDisplayValue('bubble wrap');
      });
      expect(screen.getByLabelText('Select type')).toBeInstanceOf(HTMLSelectElement);
      expect(screen.getByRole('option', { name: 'Packing materials' }).selected).toBe(true);
      expect(screen.getByText('expenseReceipt.pdf')).toBeInTheDocument();
      const deleteButton = screen.getByRole('button', { name: 'Delete' });
      expect(deleteButton).toBeInTheDocument();
      expect(screen.getByText('200KB')).toBeInTheDocument();

      expect(screen.getByLabelText('No')).toBeChecked();
      expect(screen.queryByRole('heading', { level: 3, name: 'Dates' })).not.toBeInTheDocument();

      expect(screen.getByRole('button', { name: 'Finish Later' })).toBeInTheDocument();
      expect(screen.getByRole('button', { name: 'Save & Continue' })).toBeEnabled();
    });

    it('populates edit form when reciept is missing', async () => {
      render(<ExpenseForm {...defaultProps} {...missingReceiptProps} />);
      await waitFor(() => {
        expect(screen.getByLabelText('What did you buy?')).toHaveDisplayValue('bubble wrap');
      });
      expect(
        screen.getByText(
          'If you can, get a replacement copy of your receipt and upload that. If that is not possible, write and sign a statement that explains why this receipt is missing. Include details about where and when you purchased this item. Upload that statement. Your reimbursement for this expense will be based on the information you provide.',
        ),
      ).toBeInTheDocument();
    });

    it('populates edit form with SIT values', async () => {
      render(<ExpenseForm {...defaultProps} {...sitExpenseProps} />);
      await waitFor(() => {
        expect(screen.getByLabelText('What did you buy?')).toHaveDisplayValue('10x10 storage pod');
      });
      expect(screen.getByLabelText('Select type')).toBeInstanceOf(HTMLSelectElement);
      expect(screen.getByRole('option', { name: 'Storage' }).selected).toBe(true);
      expect(screen.getByText('uhaulReceipt.pdf')).toBeInTheDocument();
      const deleteButton = screen.getByRole('button', { name: 'Delete' });
      expect(deleteButton).toBeInTheDocument();
      expect(screen.getByText('200KB')).toBeInTheDocument();
      expect(screen.getByText('Uploaded 10 Aug 2022 12:59 PM')).toBeInTheDocument();

      expect(screen.getByLabelText('No')).toBeChecked();
      expect(screen.getByRole('heading', { level: 3, name: 'Dates' })).toBeInTheDocument();
      expect(screen.getByLabelText('Start date')).toBeInstanceOf(HTMLInputElement);
      expect(screen.getByLabelText('Start date')).toHaveDisplayValue('24 Sep 2022');
      expect(screen.getByLabelText('End date')).toBeInstanceOf(HTMLInputElement);
      expect(screen.getByLabelText('End date')).toHaveDisplayValue('26 Dec 2022');

      expect(screen.getByRole('button', { name: 'Finish Later' })).toBeInTheDocument();
      expect(screen.getByRole('button', { name: 'Save & Continue' })).toBeEnabled();
    });
  });

  describe('validates the form', () => {
    it('marks required fields of empty form', async () => {
      render(<ExpenseForm {...defaultProps} />);

      await userEvent.click(screen.getByRole('button', { name: 'Save & Continue' }));

      let invalidAlerts;
      await waitFor(() => {
        invalidAlerts = screen.getAllByRole('alert');
      });

      expect(invalidAlerts).toHaveLength(2);
      expect(invalidAlerts[0].nextSibling).toHaveAttribute('name', 'description');
      expect(within(invalidAlerts[1].previousSibling).getByText('Amount')).toBeInTheDocument();
    });
  });

  describe('attaches button handler callbacks', () => {
    it('calls the onSubmit callback', async () => {
      render(<ExpenseForm {...defaultProps} {...expenseRequiredProps} />);

      await userEvent.click(screen.getByRole('button', { name: 'Save & Continue' }));

      await waitFor(() => {
        expect(defaultProps.onSubmit).toHaveBeenCalled();
      });
    });
    it('calls the onBack prop when the Finish Later button is clicked', async () => {
      render(<ExpenseForm {...defaultProps} />);

      await userEvent.click(screen.getByRole('button', { name: 'Finish Later' }));

      await waitFor(() => {
        expect(defaultProps.onBack).toHaveBeenCalled();
      });
    });
  });
});
