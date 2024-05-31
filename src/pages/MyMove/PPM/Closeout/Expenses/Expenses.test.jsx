import React from 'react';
import { render, waitFor, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { generatePath } from 'react-router-dom';
import { v4 } from 'uuid';

import { SHIPMENT_OPTIONS } from 'shared/constants';
import { MockProviders } from 'testUtils';
import { selectExpenseAndIndexById, selectMTOShipmentById } from 'store/entities/selectors';
import Expenses from 'pages/MyMove/PPM/Closeout/Expenses/Expenses';
import { customerRoutes } from 'constants/routes';
import { createBaseMovingExpense, createCompleteMovingExpense } from 'utils/test/factories/movingExpense';
import { createMovingExpense, patchMovingExpense, deleteUpload } from 'services/internalApi';

const mockNavigate = jest.fn();
const mockMoveId = 'cc03c553-d317-46af-8b2d-3c9f899f6451';
const mockMTOShipmentId = '6b7a5769-4393-46fb-a4c4-d3f6ac7584c7';
const mockExpenseId = 'ba29f5f5-0a51-4161-adaa-c568f5d5eab0';
jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom'),
  useNavigate: () => mockNavigate,
}));
jest.mock('services/internalApi', () => ({
  ...jest.requireActual('services/internalApi'),
  createMovingExpense: jest.fn(),
  createUploadForPPMDocument: jest.fn(),
  deleteUpload: jest.fn(),
  patchMovingExpense: jest.fn(),
}));

const mockPPMShipmentId = v4();
const mockMTOShipment = {
  id: mockMTOShipmentId,
  moveTaskOrderId: mockMoveId,
  shipmentType: SHIPMENT_OPTIONS.PPM,
  ppmShipment: {
    id: mockPPMShipmentId,
    pickupPostalCode: '10001',
    destinationPostalCode: '10002',
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

const mockExpense = createCompleteMovingExpense();
const mockNewExpense = createBaseMovingExpense();

const mockExpenseAndIndex = {
  expense: mockExpense,
  index: 0,
};

const mockNewExpenseAndIndex = {
  expense: mockNewExpense,
  index: 0,
};

const movePath = generatePath(customerRoutes.MOVE_HOME_PAGE);

const expensesEditPath = generatePath(customerRoutes.SHIPMENT_PPM_EXPENSES_EDIT_PATH, {
  moveId: mockMoveId,
  mtoShipmentId: mockMTOShipmentId,
  expenseId: mockExpense.id,
});
const reviewPath = generatePath(customerRoutes.SHIPMENT_PPM_REVIEW_PATH, {
  moveId: mockMoveId,
  mtoShipmentId: mockMTOShipmentId,
});

jest.mock('store/entities/selectors', () => ({
  ...jest.requireActual('store/entities/selectors'),
  selectMTOShipmentById: jest.fn(() => mockMTOShipment),
  selectExpenseAndIndexById: jest.fn(() => mockEmptyExpenseAndIndex),
}));

beforeEach(() => {
  jest.clearAllMocks();
});

const params = {
  moveId: mockMoveId,
  mtoShipmentId: mockMTOShipmentId,
};

const renderWithMocks = () => {
  render(
    <MockProviders path={customerRoutes.SHIPMENT_PPM_EXPENSES_PATH} params={params}>
      <Expenses />
    </MockProviders>,
  );
};

describe('Expenses page', () => {
  it('loads the selected shipment from redux for a new expense', async () => {
    createMovingExpense.mockResolvedValue(mockExpense);

    renderWithMocks();

    await waitFor(() => {
      expect(selectMTOShipmentById).toHaveBeenCalledWith(expect.anything(), mockMTOShipmentId);
    });

    expect(selectExpenseAndIndexById).toHaveBeenCalledWith(expect.anything(), mockMTOShipmentId, undefined);
  });

  it('displays an error if the createMovingExpense request fails', async () => {
    createMovingExpense.mockRejectedValue('an error occurred');

    renderWithMocks();

    await waitFor(() => {
      expect(screen.getByText('Failed to create trip record')).toBeInTheDocument();
    });
  });

  it('does not make create moving expense api request if id param exists', async () => {
    const editParams = {
      moveId: mockMoveId,
      mtoShipmentId: mockMTOShipmentId,
      expenseId: mockExpenseId,
    };

    selectExpenseAndIndexById.mockReturnValue(mockExpenseAndIndex);

    render(
      <MockProviders path={customerRoutes.SHIPMENT_PPM_EXPENSES_EDIT_PATH} params={editParams}>
        <Expenses />
      </MockProviders>,
    );

    await waitFor(() => {
      expect(createMovingExpense).not.toHaveBeenCalled();
    });
  });

  it('renders the page content', async () => {
    createMovingExpense.mockResolvedValue(mockExpense);
    selectExpenseAndIndexById.mockReturnValueOnce(mockEmptyExpenseAndIndex);
    selectExpenseAndIndexById.mockReturnValue(mockExpenseAndIndex);

    renderWithMocks();

    await waitFor(() => {
      expect(screen.getByTestId('tag')).toHaveTextContent('PPM');
    });

    expect(screen.getByRole('heading', { level: 1 })).toHaveTextContent('Expenses');

    // renders form content
    expect(screen.getByRole('heading', { level: 2 })).toHaveTextContent('Receipt 1');
    expect(screen.getByLabelText('Select type')).toHaveDisplayValue('Packing materials');
    expect(screen.getByLabelText('What did you buy?')).toHaveValue('Medium and large boxes');
    expect(screen.getByLabelText('No')).toBeChecked();
    expect(screen.getByLabelText('Amount')).toHaveValue('75.00');
    expect(screen.getByLabelText("I don't have this receipt")).not.toBeChecked();
    expect(screen.getByText('expense.pdf')).toBeInTheDocument();

    expect(screen.getByRole('button', { name: 'Return To Homepage' })).toBeInTheDocument();

    const saveBtn = screen.getByRole('button', { name: 'Save & Continue' });
    expect(saveBtn).toBeEnabled();
  });

  it('replaces the router history with newly created weight ticket id', async () => {
    createMovingExpense.mockResolvedValueOnce(mockExpense);
    selectExpenseAndIndexById.mockReturnValueOnce(mockEmptyExpenseAndIndex);
    selectExpenseAndIndexById.mockReturnValue(mockExpenseAndIndex);

    renderWithMocks();

    await waitFor(() => {
      expect(mockNavigate).toHaveBeenCalledWith(expensesEditPath, { replace: true });
    });
  });

  it('loads the selected shipment and existing expense from redux', async () => {
    const editParams = {
      moveId: mockMoveId,
      mtoShipmentId: mockMTOShipmentId,
      expenseId: mockExpenseId,
    };

    selectExpenseAndIndexById.mockImplementationOnce(() => mockExpenseAndIndex);

    render(
      <MockProviders path={customerRoutes.SHIPMENT_PPM_EXPENSES_EDIT_PATH} params={editParams}>
        <Expenses />
      </MockProviders>,
    );

    await waitFor(() => {
      expect(selectMTOShipmentById).toHaveBeenCalledWith(expect.anything(), mockMTOShipmentId);
    });
    expect(selectExpenseAndIndexById).toHaveBeenCalledWith(expect.anything(), mockMTOShipmentId, mockExpenseId);

    expect(screen.getByRole('heading', { level: 2 })).toHaveTextContent('Receipt 1');
    expect(screen.getByLabelText('Select type')).toHaveDisplayValue('Packing materials');
    expect(screen.getByLabelText('What did you buy?')).toHaveValue('Medium and large boxes');
    expect(screen.getByLabelText('No')).toBeChecked();
    expect(screen.getByLabelText('Amount')).toHaveValue('75.00');
    expect(screen.getByLabelText("I don't have this receipt")).not.toBeChecked();
    expect(screen.getByText('expense.pdf')).toBeInTheDocument();

    expect(screen.getByRole('button', { name: 'Return To Homepage' })).toBeInTheDocument();
    expect(screen.getByRole('button', { name: 'Save & Continue' })).toBeInTheDocument();
  });

  it('displays the creation form when adding a new expense', async () => {
    selectExpenseAndIndexById.mockReturnValueOnce(mockNewExpenseAndIndex);
    renderWithMocks();

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

    expect(screen.getByLabelText('Select type')).toHaveDisplayValue('- Select -');

    expect(screen.getByRole('button', { name: 'Return To Homepage' })).toBeInTheDocument();
    expect(screen.getByRole('button', { name: 'Save & Continue' })).toBeInTheDocument();
  });

  it('calls patchMovingExpense with the appropriate payload', async () => {
    createMovingExpense.mockResolvedValue(mockExpense);
    selectExpenseAndIndexById.mockReturnValue({ expense: mockExpense, index: 1 });
    patchMovingExpense.mockResolvedValue({});

    renderWithMocks();

    await waitFor(() => {
      expect(screen.getByRole('heading', { level: 2 })).toHaveTextContent('Receipt 2');
    });
    await userEvent.selectOptions(screen.getByLabelText('Select type'), ['CONTRACTED_EXPENSE']);
    await userEvent.clear(screen.getByLabelText('What did you buy?'));
    await userEvent.type(screen.getByLabelText('What did you buy?'), 'Boxes and tape');
    await userEvent.click(screen.getByLabelText('Yes'));
    await userEvent.clear(screen.getByLabelText('Amount'));
    await userEvent.type(screen.getByLabelText('Amount'), '12.34');
    await userEvent.click(screen.getByLabelText("I don't have this receipt"));

    expect(screen.getByRole('button', { name: 'Save & Continue' })).toBeEnabled();

    await userEvent.click(screen.getByRole('button', { name: 'Save & Continue' }));

    await waitFor(() => {
      expect(patchMovingExpense).toHaveBeenCalledWith(
        mockPPMShipmentId,
        mockExpense.id,
        {
          ppmShipmentId: mockPPMShipmentId,
          movingExpenseType: 'CONTRACTED_EXPENSE',
          description: 'Boxes and tape',
          missingReceipt: true,
          amount: 1234,
          SITEndDate: undefined,
          SITStartDate: undefined,
          paidWithGTCC: true,
          WeightStored: NaN,
          SITLocation: '',
        },
        mockExpense.eTag,
      );
    });

    expect(mockNavigate).toHaveBeenCalledWith(reviewPath);
  });

  it('has an appropriate amount payload when a whole dollar amount is entered', async () => {
    createMovingExpense.mockResolvedValue(mockExpense);
    selectExpenseAndIndexById.mockReturnValue({ expense: mockExpense, index: 1 });
    patchMovingExpense.mockResolvedValue({});

    renderWithMocks();

    await waitFor(() => {
      expect(screen.getByRole('heading', { level: 2 })).toHaveTextContent('Receipt 2');
    });
    await userEvent.clear(screen.getByLabelText('Amount'));
    await userEvent.type(screen.getByLabelText('Amount'), '12');

    expect(screen.getByRole('button', { name: 'Save & Continue' })).toBeEnabled();

    await userEvent.click(screen.getByRole('button', { name: 'Save & Continue' }));

    await waitFor(() => {
      expect(patchMovingExpense).toHaveBeenCalledWith(
        mockPPMShipmentId,
        mockExpense.id,
        {
          ppmShipmentId: mockPPMShipmentId,
          movingExpenseType: 'PACKING_MATERIALS',
          description: 'Medium and large boxes',
          missingReceipt: false,
          amount: 1200,
          SITEndDate: undefined,
          SITStartDate: undefined,
          paidWithGTCC: false,
          WeightStored: NaN,
          SITLocation: '',
        },
        mockExpense.eTag,
      );
    });

    expect(mockNavigate).toHaveBeenCalledWith(reviewPath);
  });

  it('has an appropriate payload when the type is Storage', async () => {
    createMovingExpense.mockResolvedValue(mockExpense);
    selectExpenseAndIndexById.mockReturnValue({ expense: mockExpense, index: 1 });
    patchMovingExpense.mockResolvedValue({});

    renderWithMocks();

    await waitFor(() => {
      expect(screen.getByRole('heading', { level: 2 })).toHaveTextContent('Receipt 2');
    });
    await userEvent.selectOptions(screen.getByLabelText('Select type'), ['STORAGE']);
    await userEvent.type(screen.getByLabelText('Start date'), '10/10/2022');
    await userEvent.type(screen.getByLabelText('End date'), '10/11/2022');
    await userEvent.click(screen.getByLabelText('Origin'));
    await userEvent.type(screen.getByLabelText('Weight Stored'), '120');

    expect(screen.getByRole('button', { name: 'Save & Continue' })).toBeEnabled();

    await userEvent.click(screen.getByRole('button', { name: 'Save & Continue' }));

    await waitFor(() => {
      expect(patchMovingExpense).toHaveBeenCalledWith(
        mockPPMShipmentId,
        mockExpense.id,
        {
          ppmShipmentId: mockPPMShipmentId,
          movingExpenseType: 'STORAGE',
          description: 'Medium and large boxes',
          missingReceipt: false,
          amount: 7500,
          SITEndDate: '2022-10-11',
          SITStartDate: '2022-10-10',
          paidWithGTCC: false,
          SITLocation: 'ORIGIN',
          WeightStored: 120,
        },
        mockExpense.eTag,
      );
    });

    expect(mockNavigate).toHaveBeenCalledWith(reviewPath);
  });

  it('displays an error if patchMovingExpense fails', async () => {
    createMovingExpense.mockResolvedValue(mockExpense);
    selectExpenseAndIndexById.mockReturnValue({ expense: mockExpense, index: 4 });
    patchMovingExpense.mockRejectedValueOnce('an error occurred');

    renderWithMocks();

    await waitFor(() => {
      expect(screen.getByRole('heading', { level: 2 })).toHaveTextContent('Receipt 5');
    });
    await userEvent.selectOptions(screen.getByLabelText('Select type'), ['CONTRACTED_EXPENSE']);
    await userEvent.type(screen.getByLabelText('What did you buy?'), 'Boxes and tape');
    await userEvent.click(screen.getByLabelText('Yes'));
    await userEvent.clear(screen.getByLabelText('Amount'));
    await userEvent.type(screen.getByLabelText('Amount'), '12.34');
    await userEvent.click(screen.getByLabelText("I don't have this receipt"));

    await userEvent.click(screen.getByRole('button', { name: 'Save & Continue' }));

    await waitFor(() => {
      expect(screen.getByText('Failed to save updated trip record')).toBeInTheDocument();
    });
  });

  it('routes to home when the return to homepage button is clicked', async () => {
    createMovingExpense.mockResolvedValue(mockExpense);

    renderWithMocks();

    await waitFor(() => {
      expect(selectMTOShipmentById).toHaveBeenCalledWith(expect.anything(), mockMTOShipmentId);
    });

    await waitFor(() => {
      expect(screen.getByRole('button', { name: 'Return To Homepage' })).toBeInTheDocument();
    });
    await userEvent.click(screen.getByRole('button', { name: 'Return To Homepage' }));
    expect(mockNavigate).toHaveBeenCalledWith(movePath);
  });

  it('calls the delete handler when removing an existing upload', async () => {
    const editParams = {
      moveId: mockMoveId,
      mtoShipmentId: mockMTOShipmentId,
      expenseId: mockExpense.id,
    };

    selectExpenseAndIndexById.mockReturnValue({ expense: mockExpense, index: 0 });

    selectMTOShipmentById.mockReturnValue({
      ...mockMTOShipment,
      ppmShipment: {
        ...mockMTOShipment.ppmShipment,
        expenses: [mockExpense],
      },
    });
    deleteUpload.mockResolvedValue({});
    render(
      <MockProviders path={customerRoutes.SHIPMENT_PPM_EXPENSES_EDIT_PATH} params={editParams}>
        <Expenses />
      </MockProviders>,
    );

    let deleteButtons;
    await waitFor(() => {
      deleteButtons = screen.getAllByRole('button', { name: 'Delete' });
      expect(deleteButtons).toHaveLength(1);
    });
    await userEvent.click(deleteButtons[0]);
    await waitFor(() => {
      expect(screen.queryByText('empty_weight.jpg')).not.toBeInTheDocument();
    });
  });

  it('expect loadingPlaceholder when mtoShipment is falsy', async () => {
    selectMTOShipmentById.mockReturnValueOnce(null);

    renderWithMocks();

    await waitFor(() => {
      expect(screen.getByRole('heading', { level: 2 })).toHaveTextContent('Loading, please wait...');
    });
  });
});
