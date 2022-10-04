import React from 'react';
import { render, waitFor, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { generatePath, useParams } from 'react-router-dom';
import { v4 } from 'uuid';

import { SHIPMENT_OPTIONS } from 'shared/constants';
import { MockProviders } from 'testUtils';
import { selectExpenseAndIndexById, selectMTOShipmentById } from 'store/entities/selectors';
import Expenses from 'pages/MyMove/PPM/Closeout/Expenses/Expenses';
import { customerRoutes, generalRoutes } from 'constants/routes';
import { createCompleteMovingExpense } from 'utils/test/factories/movingExpense';

const mockPush = jest.fn();
const mockReplace = jest.fn();
const mockMoveId = 'cc03c553-d317-46af-8b2d-3c9f899f6451';
const mockMTOShipmentId = '6b7a5769-4393-46fb-a4c4-d3f6ac7584c7';
const mockExpenseId = 'ba29f5f5-0a51-4161-adaa-c568f5d5eab0';
jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom'),
  useHistory: () => ({
    push: mockPush,
    replace: mockReplace,
  }),
  useParams: jest.fn(() => ({
    moveId: mockMoveId,
    mtoShipmentId: mockMTOShipmentId,
  })),
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

const mockExpenseAndIndex = {
  expense: createCompleteMovingExpense(),
  index: 0,
};

const homePath = generatePath(generalRoutes.HOME_PATH);

jest.mock('store/entities/selectors', () => ({
  ...jest.requireActual('store/entities/selectors'),
  selectMTOShipmentById: jest.fn(() => mockMTOShipment),
  selectExpenseAndIndexById: jest.fn(() => mockEmptyExpenseAndIndex),
}));

beforeEach(() => {
  jest.clearAllMocks();
});

describe('Expenses page', () => {
  it('loads the selected shipment from redux for a new expense', async () => {
    render(<Expenses />, { wrapper: MockProviders });

    await waitFor(() => {
      expect(selectMTOShipmentById).toHaveBeenCalledWith(expect.anything(), mockMTOShipmentId);
    });

    expect(selectExpenseAndIndexById).toHaveBeenCalledWith(expect.anything(), mockMTOShipmentId, undefined);
  });

  it('loads the selected shipment and existing expense from redux', async () => {
    useParams.mockImplementationOnce(() => ({
      moveId: mockMoveId,
      mtoShipmentId: mockMTOShipmentId,
      expenseId: mockExpenseId,
    }));

    selectExpenseAndIndexById.mockImplementationOnce(() => mockExpenseAndIndex);

    render(<Expenses />, { wrapper: MockProviders });

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

    const saveBtn = screen.getByRole('button', { name: 'Save & Continue' });
    expect(saveBtn).toBeEnabled();

    await userEvent.click(saveBtn);

    await waitFor(() => {
      expect(mockPush).toHaveBeenCalledWith(
        generatePath(customerRoutes.SHIPMENT_PPM_REVIEW_PATH, { moveId: mockMoveId, mtoShipmentId: mockMTOShipmentId }),
      );
    });
  });

  it('displays the creation form when adding a new expense', async () => {
    render(<Expenses />, { wrapper: MockProviders });

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

  it('routes to home when the return to homepage button is clicked', async () => {
    render(<Expenses />, { wrapper: MockProviders });

    await waitFor(() => {
      expect(screen.getByRole('button', { name: 'Return To Homepage' })).toBeInTheDocument();
    });
    await userEvent.click(screen.getByRole('button', { name: 'Return To Homepage' }));
    expect(mockPush).toHaveBeenCalledWith(homePath);
  });
});
