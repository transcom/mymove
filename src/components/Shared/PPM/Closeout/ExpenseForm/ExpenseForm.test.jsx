import React from 'react';
import { render, waitFor, screen, within } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import ExpenseForm from 'components/Shared/PPM/Closeout/ExpenseForm/ExpenseForm';
import { DocumentAndImageUploadInstructions } from 'content/uploads';
import { expenseTypes } from 'constants/ppmExpenseTypes';
import { PPM_TYPES } from 'shared/constants';
import { APP_NAME } from 'constants/apps';

beforeEach(() => {
  jest.clearAllMocks();
});

const defaultProps = {
  expense: {
    id: '32ecb311-edbe-4fd4-96ee-bd693113f3f3',
    ppmShipmentId: '343bb456-63af-4f76-89bd-7403094a5c4d',
    movingExpenseType: expenseTypes.PACKING_MATERIALS,
  },
  receiptNumber: 1,
  onCreateUpload: jest.fn(),
  onUploadComplete: jest.fn(),
  onUploadDelete: jest.fn(),
  onBack: jest.fn(),
  onSubmit: jest.fn(),
};

const smallPackageProps = {
  expense: {
    id: '32ecb311-edbe-4fd4-96ee-bd693113f3f3',
    ppmShipmentId: '343bb456-63af-4f76-89bd-7403094a5c4d',
  },
  receiptNumber: 1,
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
    movingExpenseType: expenseTypes.PACKING_MATERIALS,
    description: 'bubble wrap',
    missingReceipt: true,
  },
  receiptNumber: 1,
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
    movingExpenseType: expenseTypes.PACKING_MATERIALS,
    description: 'bubble wrap',
    missingReceipt: false,
    paidWithGtcc: false,
    amount: 60000,
    document: {
      uploads: [
        {
          id: 'db4713ae-6087-4330-8b0d-926b3d65c454',
          createdAt: '2022-06-10T12:59:30.000Z',
          bytes: 204800,
          url: 'some/path/to/',
          filename: 'expenseReceipt.pdf',
          contentType: 'application/pdf',
        },
      ],
    },
  },
};

const sitExpenseProps = {
  expense: {
    id: '32ecb311-edbe-4fd4-96ee-bd693113f3f3',
    ppmShipmentId: '343bb456-63af-4f76-89bd-7403094a5c4d',
    movingExpenseType: expenseTypes.STORAGE,
    description: '10x10 storage pod',
    missingReceipt: false,
    paidWithGtcc: false,
    amount: 16099,
    sitStartDate: '2022-09-24',
    sitEndDate: '2022-12-26',
    sitLocation: 'ORIGIN',
    weightStored: '120',
    document: {
      uploads: [
        {
          id: 'db4713ae-6087-4330-8b0d-926b3d65c454',
          createdAt: '2022-08-10T12:59:30.000Z',
          bytes: 204800,
          url: 'some/path/to/',
          filename: 'uhaulReceipt.pdf',
          contentType: 'application/pdf',
        },
      ],
    },
  },
};

const smallPackageExpenseProps = {
  expense: {
    paidWithGtcc: false,
    amount: 5309,
    missingReceipt: false,
    document: { uploads: [] },
    trackingNumber: 'track THIS!',
    weightShipped: '500',
    isProGear: true,
    proGearBelongsToSelf: false,
    proGearDescription: 'describte THAT',
  },
};

describe('ExpenseForm component', () => {
  describe('displays form', () => {
    it('renders blank form on load with defaults - Customer page', async () => {
      render(<ExpenseForm {...defaultProps} appName={APP_NAME.MYMOVE} />);

      await waitFor(() => {
        expect(screen.getByRole('heading', { level: 2, name: 'Receipt 1' })).toBeInTheDocument();
      });

      expect(
        screen.getByText(
          'Document your qualified expenses by uploading receipts. They should include a description of the item, the price you paid, the date of purchase, and the business name. All documents must be legible and unaltered.',
        ),
      ).toBeInTheDocument();

      expect(screen.getByLabelText('Select type')).toBeInstanceOf(HTMLSelectElement);
      expect(screen.getByRole('heading', { level: 3, name: 'Description' })).toBeInTheDocument();
      expect(screen.getByLabelText('What did you buy or rent?')).toBeInstanceOf(HTMLInputElement);
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

      expect(screen.getByRole('button', { name: 'Cancel' })).toBeInTheDocument();
      expect(screen.getByRole('button', { name: 'Save & Continue' })).toBeInTheDocument();
    });

    it('renders blank form on load with defaults - Office page', async () => {
      render(<ExpenseForm {...defaultProps} appName={APP_NAME.OFFICE} />);

      await waitFor(() => {
        expect(screen.getByRole('heading', { level: 2, name: 'Receipt 1' })).toBeInTheDocument();
      });

      expect(screen.getByLabelText('Select type')).toBeInstanceOf(HTMLSelectElement);
      expect(screen.getByRole('heading', { level: 3, name: 'Description' })).toBeInTheDocument();
      expect(screen.getByLabelText('What did you buy or rent?')).toBeInstanceOf(HTMLInputElement);
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

      expect(screen.getByRole('button', { name: 'Cancel' })).toBeInTheDocument();
      expect(screen.getByRole('button', { name: 'Save & Continue' })).toBeInTheDocument();
    });

    it('populates edit form with existing expense values', async () => {
      render(<ExpenseForm {...defaultProps} {...expenseRequiredProps} appName={APP_NAME.MYMOVE} />);

      await waitFor(() => {
        expect(screen.getByLabelText('What did you buy or rent?')).toHaveDisplayValue('bubble wrap');
      });
      expect(screen.getByLabelText('Select type')).toBeInstanceOf(HTMLSelectElement);
      expect(screen.getAllByRole('option')[3].selected).toBe(true);
      expect(screen.getByText('expenseReceipt.pdf')).toBeInTheDocument();
      const deleteButton = screen.getByRole('button', { name: 'Delete' });
      expect(deleteButton).toBeInTheDocument();
      expect(screen.getByText('200KB')).toBeInTheDocument();

      expect(screen.getByLabelText('No')).toBeChecked();
      expect(screen.queryByRole('heading', { level: 3, name: 'Dates' })).not.toBeInTheDocument();

      expect(screen.getByRole('button', { name: 'Cancel' })).toBeInTheDocument();
      expect(screen.getByRole('button', { name: 'Save & Continue' })).toBeEnabled();
    });

    it('populates edit form when reciept is missing', async () => {
      render(<ExpenseForm {...defaultProps} {...missingReceiptProps} />);
      await waitFor(() => {
        expect(screen.getByLabelText('What did you buy or rent?')).toHaveDisplayValue('bubble wrap');
      });
      expect(
        screen.getByText(
          'If you can, get a replacement copy of your receipt and upload that. If that is not possible, write and sign a statement that explains why this receipt is missing. Include details about where and when you purchased this item. Upload that statement. Your reimbursement for this expense will be based on the information you provide.',
        ),
      ).toBeInTheDocument();
    });

    it('populates edit form with SIT values', async () => {
      render(<ExpenseForm {...defaultProps} {...sitExpenseProps} appName={APP_NAME.MYMOVE} />);
      await waitFor(() => {
        expect(screen.getByLabelText('What did you buy or rent?')).toHaveDisplayValue('10x10 storage pod');
      });
      expect(screen.getByLabelText('Select type')).toBeInstanceOf(HTMLSelectElement);
      expect(screen.getByRole('option', { name: 'Storage' }).selected).toBe(true);
      expect(screen.getByText('uhaulReceipt.pdf')).toBeInTheDocument();
      const deleteButton = screen.getByRole('button', { name: 'Delete' });
      expect(deleteButton).toBeInTheDocument();
      expect(screen.getByText('200KB')).toBeInTheDocument();
      expect(screen.getByText('Uploaded 10 Aug 2022 12:59 PM')).toBeInTheDocument();

      expect(screen.getByLabelText('Origin')).toBeChecked();
      expect(screen.getByLabelText('Destination')).not.toBeChecked();
      expect(screen.getByLabelText('Weight Stored')).toHaveDisplayValue('120');
      expect(screen.getByLabelText('No')).toBeChecked();
      expect(screen.getByRole('heading', { level: 3, name: 'Dates' })).toBeInTheDocument();
      expect(screen.getByLabelText('Start date')).toBeInstanceOf(HTMLInputElement);
      expect(screen.getByLabelText('Start date')).toHaveDisplayValue('24 Sep 2022');
      expect(screen.getByLabelText('End date')).toBeInstanceOf(HTMLInputElement);
      expect(screen.getByLabelText('End date')).toHaveDisplayValue('26 Dec 2022');

      expect(screen.getByRole('button', { name: 'Cancel' })).toBeInTheDocument();
      expect(screen.getByRole('button', { name: 'Save & Continue' })).toBeEnabled();
    });

    it('renders base expense form for small package', async () => {
      render(<ExpenseForm ppmType={PPM_TYPES.SMALL_PACKAGE} {...smallPackageProps} />);

      expect(screen.getByLabelText('Select type')).toBeInstanceOf(HTMLSelectElement);
      expect(screen.getByRole('option', { name: 'Small package reimbursement' }).selected).toBe(true);
      expect(screen.getByRole('option', { name: 'Small package reimbursement' })).toBeDisabled();
      expect(screen.getByTestId('smallPackageInfo')).toBeInTheDocument();

      await waitFor(() => {
        expect(screen.getByTestId('weightShipped')).toBeInTheDocument();
      });
      // the extra pro gear fields should not be rendered until the user selects that the expense is pro gear
      expect(screen.queryByText(/Who does this pro-gear belong to/i)).toBeNull();
      expect(screen.queryByLabelText(/Brief description of the pro-gear/i)).toBeNull();
    });

    it('renders pro gear values on existing expense form for small package', async () => {
      render(<ExpenseForm ppmType={PPM_TYPES.SMALL_PACKAGE} {...smallPackageProps} {...smallPackageExpenseProps} />);

      expect(screen.getByLabelText('Select type')).toBeInstanceOf(HTMLSelectElement);
      expect(screen.getByRole('option', { name: 'Small package reimbursement' }).selected).toBe(true);
      expect(screen.getByRole('option', { name: 'Small package reimbursement' })).toBeDisabled();

      // these should be visible now because smallPackageExpense props has a isProGear value of true
      await waitFor(() => {
        expect(screen.getByTestId('proGearWeight')).toBeInTheDocument();
      });
      expect(screen.queryByText(/Who does this pro-gear belong to/i)).toBeInTheDocument();
      expect(screen.queryByText(/Brief description of the pro-gear/i)).toBeInTheDocument();
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

  describe('attaches button handler callbacks - Customer page', () => {
    it('calls the onSubmit callback', async () => {
      render(<ExpenseForm {...defaultProps} {...expenseRequiredProps} appName={APP_NAME.MYMOVE} />);

      await userEvent.click(screen.getByRole('button', { name: 'Save & Continue' }));

      await waitFor(() => {
        expect(defaultProps.onSubmit).toHaveBeenCalled();
      });
    });
    it('calls the onBack prop when the back button is clicked', async () => {
      render(<ExpenseForm {...defaultProps} appName={APP_NAME.MYMOVE} />);

      await userEvent.click(screen.getByRole('button', { name: 'Cancel' }));

      await waitFor(() => {
        expect(defaultProps.onBack).toHaveBeenCalled();
      });
    });
  });

  describe('attaches button handler callbacks - Office page', () => {
    it('calls the onSubmit callback', async () => {
      render(<ExpenseForm {...defaultProps} {...expenseRequiredProps} appName={APP_NAME.OFFICE} />);

      await userEvent.click(screen.getByRole('button', { name: 'Save & Continue' }));

      await waitFor(() => {
        expect(defaultProps.onSubmit).toHaveBeenCalled();
      });
    });
    it('calls the onBack prop when the Cancel button is clicked', async () => {
      render(<ExpenseForm {...defaultProps} appName={APP_NAME.OFFICE} />);

      await userEvent.click(screen.getByRole('button', { name: 'Cancel' }));

      await waitFor(() => {
        expect(defaultProps.onBack).toHaveBeenCalled();
      });
    });
  });
});
