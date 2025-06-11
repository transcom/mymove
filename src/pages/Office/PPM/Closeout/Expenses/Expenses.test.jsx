import React from 'react';
import { render, waitFor, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { generatePath } from 'react-router-dom';
import { v4 } from 'uuid';

import { SHIPMENT_OPTIONS } from 'shared/constants';
import { MockProviders } from 'testUtils';
import Expenses from 'pages/Office/PPM/Closeout/Expenses/Expenses';
import { expenseTypes } from 'constants/ppmExpenseTypes';
import { servicesCounselingRoutes } from 'constants/routes';
import { createMovingExpense, patchExpense, deleteUploadForDocument } from 'services/ghcApi';
import { usePPMShipmentAndDocsOnlyQueries } from 'hooks/queries';

// crete local alias for more descriptive name
const generateUUID = () => v4();

const mockMoveId = 'cc03c553-d317-46af-8b2d-3c9f899f6451';
const mockMTOShipmentId = '6b7a5769-4393-46fb-a4c4-d3f6ac7584c7';
const mockExpenseId = 'ba29f5f5-0a51-4161-adaa-c568f5d5eab0';
const mockExpenseEtag = window.btoa(new Date());
const mockExpenseDocumentId = generateUUID();
const mockExpenseDocumentUploadId = generateUUID();

const mockNavigate = jest.fn();
jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom'),
  useNavigate: () => mockNavigate,
}));

jest.mock('services/ghcApi', () => ({
  ...jest.requireActual('services/ghcApi'),
  createMovingExpense: jest.fn(),
  createUploadForPPMDocument: jest.fn(),
  deleteUploadForDocument: jest.fn(),
  patchExpense: jest.fn(),
}));

jest.mock('hooks/queries', () => ({
  usePPMShipmentAndDocsOnlyQueries: jest.fn(),
}));

const mockPPMShipmentId = generateUUID();
const mockMTOShipment = {
  id: mockMTOShipmentId,
  moveTaskOrderId: mockMoveId,
  shipmentType: SHIPMENT_OPTIONS.PPM,
  ppmShipment: {
    id: mockPPMShipmentId,
    expectedDepartureDate: '2022-04-30',
    advanceRequested: true,
    advance: 598700,
    estimatedWeight: 4000,
    estimatedIncentive: 1000000,
    sitExpected: false,
    hasProGear: false,
    proGearWeight: null,
    spouseProGearWeight: null,
  },
  eTag: 'dGVzdGluZzIzNDQzMjQ',
};

const mockEmptyExpenseAndIndex = {
  expense: null,
  index: -1,
};

const mockExpenseWithNoValues = {
  id: mockExpenseId,
  ppmShipmentId: mockPPMShipmentId,
  documentId: mockExpenseDocumentId,
  eTag: mockExpenseEtag,
};
const mockExpense = {
  id: mockExpenseId,
  ppmShipmentId: mockPPMShipmentId,
  amount: 7500,
  description: 'Medium and large boxes',
  movingExpenseType: expenseTypes.PACKING_MATERIALS,
  documentId: mockExpenseDocumentId,
  eTag: mockExpenseEtag,
};
const mockExpenseWithUpload = {
  id: mockExpenseId,
  ppmShipmentId: mockPPMShipmentId,
  amount: 8500,
  description: 'Peanuts and wrapping paper',
  movingExpenseType: expenseTypes.PACKING_MATERIALS,
  documentId: mockExpenseDocumentId,
  document: {
    id: mockExpenseDocumentId,
    uploads: [
      {
        id: mockExpenseDocumentUploadId,
        createdAt: '2022-06-22T23:25:50.490Z',
        bytes: 819200,
        url: 'a/fake/path',
        filename: 'an_expense.jpg',
        contentType: 'image/jpg',
      },
    ],
  },
  eTag: mockExpenseEtag,
};

const expensesEditPath = generatePath(servicesCounselingRoutes.BASE_SHIPMENT_PPM_EXPENSES_EDIT_PATH, {
  moveCode: mockMoveId,
  shipmentId: mockMTOShipmentId,
  expenseId: mockExpense.id,
});
const reviewPath = generatePath(servicesCounselingRoutes.BASE_SHIPMENT_PPM_REVIEW_PATH, {
  moveCode: mockMoveId,
  shipmentId: mockMTOShipmentId,
});

const renderExpensesPage = () => {
  const mockRoutingConfig = {
    path: servicesCounselingRoutes.BASE_SHIPMENT_PPM_EXPENSES_PATH,
    params: { moveCode: mockMoveId, shipmentId: mockMTOShipmentId },
  };

  render(
    <MockProviders {...mockRoutingConfig}>
      <Expenses />
    </MockProviders>,
  );
};

const renderEditExpensesPage = () => {
  const mockRoutingConfig = {
    path: servicesCounselingRoutes.BASE_SHIPMENT_PPM_EXPENSES_EDIT_PATH,
    params: { moveCode: mockMoveId, shipmentId: mockMTOShipmentId, expenseId: mockExpenseId },
  };

  render(
    <MockProviders {...mockRoutingConfig}>
      <Expenses />
    </MockProviders>,
  );
};

jest.mock('store/entities/selectors', () => ({
  ...jest.requireActual('store/entities/selectors'),
  selectMTOShipmentById: jest.fn(() => mockMTOShipment),
  selectExpenseAndIndexById: jest.fn(() => mockEmptyExpenseAndIndex),
}));

beforeEach(() => {
  jest.clearAllMocks();
});

describe('Expenses page', () => {
  it('displays an error if the createMovingExpense request fails', async () => {
    usePPMShipmentAndDocsOnlyQueries.mockReturnValue({
      isLoading: false,
      mtoShipment: mockMTOShipment,
      documents: { MovingExpenses: [mockExpense] },
      isError: null,
    });

    createMovingExpense.mockRejectedValue('an error occurred');

    renderExpensesPage();

    await waitFor(() => {
      expect(screen.getByText('Failed to create trip record')).toBeInTheDocument();
    });
  });

  it('does not make create moving expense api request if id param exists', async () => {
    usePPMShipmentAndDocsOnlyQueries.mockReturnValue({
      isLoading: false,
      mtoShipment: mockMTOShipment,
      documents: { MovingExpenses: [mockExpense] },
      isError: null,
    });

    renderEditExpensesPage();

    await waitFor(() => {
      expect(createMovingExpense).not.toHaveBeenCalled();
    });
  });

  it('renders the page content', async () => {
    createMovingExpense.mockResolvedValue(mockExpenseWithUpload);
    usePPMShipmentAndDocsOnlyQueries.mockReturnValue({
      isLoading: false,
      mtoShipment: mockMTOShipment,
      documents: { MovingExpenses: [mockExpenseWithUpload] },
      isError: null,
    });

    renderEditExpensesPage();

    await waitFor(() => {
      expect(screen.getByTestId('tag')).toHaveTextContent('PPM');
    });

    expect(screen.getByRole('heading', { level: 1 })).toHaveTextContent('Expenses');

    // renders form content
    expect(screen.getByRole('heading', { level: 2 })).toHaveTextContent('Receipt 1');
    expect(screen.getByLabelText('Select type *')).toHaveDisplayValue('Packing materials');
    expect(screen.getByLabelText('What did you buy or rent? *')).toHaveValue('Peanuts and wrapping paper');
    expect(screen.getByLabelText('No')).toBeChecked();
    expect(screen.getByLabelText('Amount *')).toHaveValue('85.00');
    expect(screen.getByLabelText("I don't have this receipt")).not.toBeChecked();
    expect(screen.getByText('an_expense.jpg')).toBeInTheDocument();

    expect(screen.getByRole('button', { name: 'Cancel' })).toBeInTheDocument();

    const saveBtn = screen.getByRole('button', { name: 'Save & Continue' });
    expect(saveBtn).toBeEnabled();
  });

  it('replaces the router history with newly created expense id', async () => {
    createMovingExpense.mockResolvedValueOnce(mockExpense);

    usePPMShipmentAndDocsOnlyQueries.mockReturnValue({
      isLoading: false,
      mtoShipment: mockMTOShipment,
      documents: { MovingExpenses: [mockExpense] },
      isError: null,
    });

    renderExpensesPage();

    await waitFor(() => {
      expect(mockNavigate).toHaveBeenCalledWith(expensesEditPath, { replace: true });
    });
  });

  it('displays the creation form when adding a new expense', async () => {
    createMovingExpense.mockResolvedValue(mockExpenseWithNoValues);
    usePPMShipmentAndDocsOnlyQueries.mockReturnValue({
      isLoading: false,
      mtoShipment: mockMTOShipment,
      documents: { MovingExpenses: [mockExpenseWithNoValues] },
      isError: null,
    });

    renderEditExpensesPage();

    await waitFor(() => {
      expect(screen.getByTestId('tag')).toHaveTextContent('PPM');
    });

    expect(screen.getByRole('heading', { level: 1 })).toHaveTextContent('Expenses');
    expect(
      screen.getByText(
        'Document your qualified expenses by uploading receipts. They should include a description of the item, the price you paid, the date of purchase, and the business name. All documents must be legible and unaltered.',
      ),
    ).toBeInTheDocument();

    expect(
      screen.getByText(
        'Your finance office will make the final decision about which expenses are deductible or reimbursable.',
      ),
    ).toBeInTheDocument();

    expect(
      screen.getByText('Upload one receipt at a time. Please do not put multiple receipts in one image.'),
    ).toBeInTheDocument();

    expect(screen.getByRole('heading', { level: 2 })).toHaveTextContent('Receipt 1');

    expect(screen.getByLabelText('Select type *')).toHaveDisplayValue('- Select -');

    expect(screen.getByRole('button', { name: 'Cancel' })).toBeInTheDocument();
    expect(screen.getByRole('button', { name: 'Save & Continue' })).toBeInTheDocument();
  });

  it('calls patchExpense with the appropriate payload', async () => {
    usePPMShipmentAndDocsOnlyQueries.mockReturnValue({
      isLoading: false,
      mtoShipment: mockMTOShipment,
      documents: { MovingExpenses: [mockExpenseWithUpload] },
      isError: null,
    });

    patchExpense.mockResolvedValue();

    renderEditExpensesPage();

    await waitFor(() => {
      expect(screen.getByRole('heading', { level: 2 })).toHaveTextContent('Receipt 1');
    });
    await userEvent.selectOptions(screen.getByLabelText('Select type *'), ['CONTRACTED_EXPENSE']);
    await userEvent.clear(screen.getByLabelText('What did you buy or rent? *'));
    await userEvent.type(screen.getByLabelText('What did you buy or rent? *'), 'Boxes and tape');
    await userEvent.click(screen.getByLabelText('Yes'));
    await userEvent.clear(screen.getByLabelText('Amount *'));
    await userEvent.type(screen.getByLabelText('Amount *'), '12.34');
    await userEvent.click(screen.getByLabelText("I don't have this receipt"));

    expect(screen.getByRole('button', { name: 'Save & Continue' })).toBeEnabled();

    await userEvent.click(screen.getByRole('button', { name: 'Save & Continue' }));

    await waitFor(() => {
      expect(patchExpense).toHaveBeenCalledWith({
        ppmShipmentId: mockPPMShipmentId,
        movingExpenseId: mockExpenseId,
        payload: {
          ppmShipmentId: mockPPMShipmentId,
          movingExpenseType: 'CONTRACTED_EXPENSE',
          description: 'Boxes and tape',
          missingReceipt: true,
          amount: 1234,
          SITEndDate: undefined,
          SITStartDate: undefined,
          paidWithGTCC: true,
          WeightStored: NaN,
          SITLocation: undefined,
          isProGear: false,
          trackingNumber: '',
          weightShipped: NaN,
        },
        eTag: mockExpenseEtag,
      });
    });

    expect(mockNavigate).toHaveBeenCalledWith(reviewPath);
  });

  it('has an appropriate amount payload when a whole dollar amount is entered', async () => {
    usePPMShipmentAndDocsOnlyQueries.mockReturnValue({
      isLoading: false,
      mtoShipment: mockMTOShipment,
      documents: { MovingExpenses: [mockExpenseWithUpload] },
      isError: null,
    });

    patchExpense.mockResolvedValue();

    renderEditExpensesPage();

    await waitFor(() => {
      expect(screen.getByRole('heading', { level: 2 })).toHaveTextContent('Receipt 1');
    });
    await userEvent.clear(screen.getByLabelText('Amount *'));
    await userEvent.type(screen.getByLabelText('Amount *'), '12');

    expect(screen.getByRole('button', { name: 'Save & Continue' })).toBeEnabled();

    await userEvent.click(screen.getByRole('button', { name: 'Save & Continue' }));

    await waitFor(() => {
      expect(patchExpense).toHaveBeenCalledWith({
        ppmShipmentId: mockPPMShipmentId,
        movingExpenseId: mockExpenseId,
        payload: {
          ppmShipmentId: mockPPMShipmentId,
          movingExpenseType: 'PACKING_MATERIALS',
          description: 'Peanuts and wrapping paper',
          missingReceipt: false,
          amount: 1200,
          isProGear: false,
          SITEndDate: undefined,
          SITStartDate: undefined,
          paidWithGTCC: false,
          WeightStored: NaN,
          SITLocation: undefined,
          trackingNumber: '',
          weightShipped: NaN,
        },
        eTag: mockExpenseEtag,
      });
    });

    expect(mockNavigate).toHaveBeenCalledWith(reviewPath);
  });

  it('has an appropriate payload when the type is Storage', async () => {
    usePPMShipmentAndDocsOnlyQueries.mockReturnValue({
      isLoading: false,
      mtoShipment: mockMTOShipment,
      documents: { MovingExpenses: [mockExpenseWithUpload] },
      isError: null,
    });

    patchExpense.mockResolvedValue();

    renderEditExpensesPage();

    await waitFor(() => {
      expect(screen.getByRole('heading', { level: 2 })).toHaveTextContent('Receipt 1');
    });

    await userEvent.selectOptions(screen.getByLabelText('Select type *'), ['STORAGE']);
    await userEvent.type(screen.getByLabelText('Start date *'), '10/10/2022');
    await userEvent.type(screen.getByLabelText('End date *'), '10/11/2022');
    await userEvent.click(screen.getByLabelText('Origin'));
    await userEvent.type(screen.getByLabelText('Weight Stored *'), '120');

    expect(screen.getByRole('button', { name: 'Save & Continue' })).toBeEnabled();

    await userEvent.click(screen.getByRole('button', { name: 'Save & Continue' }));

    await waitFor(() => {
      expect(patchExpense).toHaveBeenCalledWith({
        ppmShipmentId: mockPPMShipmentId,
        movingExpenseId: mockExpenseId,
        payload: {
          ppmShipmentId: mockPPMShipmentId,
          trackingNumber: '',
          weightShipped: NaN,
          movingExpenseType: 'STORAGE',
          description: 'Peanuts and wrapping paper',
          isProGear: false,
          missingReceipt: false,
          amount: 8500,
          SITEndDate: '2022-10-11',
          SITStartDate: '2022-10-10',
          paidWithGTCC: false,
          SITLocation: 'ORIGIN',
          WeightStored: 120,
        },
        eTag: mockExpenseEtag,
      });
    });

    expect(mockNavigate).toHaveBeenCalledWith(reviewPath);
  });

  it('displays an error if patchExpense fails', async () => {
    usePPMShipmentAndDocsOnlyQueries.mockReturnValue({
      isLoading: false,
      mtoShipment: mockMTOShipment,
      documents: { MovingExpenses: [{}, {}, {}, {}, mockExpenseWithUpload] },
      isError: null,
    });
    patchExpense.mockRejectedValueOnce('an error occurred');

    renderEditExpensesPage();

    await waitFor(() => {
      expect(screen.getByRole('heading', { level: 2 })).toHaveTextContent('Receipt 5');
    });
    await userEvent.selectOptions(screen.getByLabelText('Select type *'), ['CONTRACTED_EXPENSE']);
    await userEvent.type(screen.getByLabelText('What did you buy or rent? *'), 'Boxes and tape');
    await userEvent.click(screen.getByLabelText('Yes'));
    await userEvent.clear(screen.getByLabelText('Amount *'));
    await userEvent.type(screen.getByLabelText('Amount *'), '12.34');
    await userEvent.click(screen.getByLabelText("I don't have this receipt"));

    await userEvent.click(screen.getByRole('button', { name: 'Save & Continue' }));

    await waitFor(() => {
      expect(screen.getByText('Failed to save updated trip record')).toBeInTheDocument();
    });
  });

  it('routes to review when the cancel button is clicked', async () => {
    createMovingExpense.mockResolvedValue(mockExpense);
    usePPMShipmentAndDocsOnlyQueries.mockReturnValue({
      isLoading: false,
      mtoShipment: mockMTOShipment,
      documents: { MovingExpenses: [mockExpense] },
      isError: null,
    });

    renderEditExpensesPage();

    await waitFor(() => {
      expect(screen.getByRole('button', { name: 'Cancel' })).toBeInTheDocument();
    });
    await userEvent.click(screen.getByRole('button', { name: 'Cancel' }));
    expect(mockNavigate).toHaveBeenCalledWith(reviewPath);
  });

  it('calls the delete handler when removing an existing upload', async () => {
    usePPMShipmentAndDocsOnlyQueries.mockReturnValue({
      isLoading: false,
      mtoShipment: mockMTOShipment,
      documents: { MovingExpenses: [{}, {}, {}, {}, mockExpenseWithUpload] },
      isError: null,
    });
    deleteUploadForDocument.mockResolvedValue({});
    renderEditExpensesPage();

    let deleteButtons;
    await waitFor(() => {
      deleteButtons = screen.getAllByRole('button', { name: 'Delete' });
      expect(deleteButtons).toHaveLength(1);
    });
    await userEvent.click(deleteButtons[0]);
    await waitFor(() => {
      expect(screen.queryByText('an_expense.jpg')).not.toBeInTheDocument();
    });
  });

  it('expect loadingPlaceholder when mtoShipment is falsy', async () => {
    usePPMShipmentAndDocsOnlyQueries.mockReturnValue({
      isLoading: false,
      mtoShipment: null,
      documents: { MovingExpenses: [mockExpenseWithUpload] },
      isError: null,
    });

    renderExpensesPage();

    await waitFor(() => {
      expect(screen.getByRole('heading', { level: 2 })).toHaveTextContent('Loading, please wait...');
    });
  });
});
