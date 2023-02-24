import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import { v4 } from 'uuid';
import userEvent from '@testing-library/user-event';
import { generatePath } from 'react-router-dom';

import { MockProviders, setUpProvidersWithHistory } from 'testUtils';
import { selectMTOShipmentById } from 'store/entities/selectors';
import Review from 'pages/MyMove/PPM/Closeout/Review/Review';
import { SHIPMENT_OPTIONS } from 'shared/constants';
import { customerRoutes } from 'constants/routes';
import {
  deleteWeightTicket,
  deleteProGearWeightTicket,
  deleteMovingExpense,
  getMTOShipmentsForMove,
} from 'services/internalApi';
import { createBaseWeightTicket, createCompleteWeightTicket } from 'utils/test/factories/weightTicket';
import { createBaseProGearWeightTicket } from 'utils/test/factories/proGearWeightTicket';
import { createCompleteMovingExpense, createCompleteSITMovingExpense } from 'utils/test/factories/movingExpense';

const mockMoveId = v4();
const mockMTOShipmentId = v4();
const mockPPMShipmentId = v4();

const mockMTOShipment = {
  id: mockMTOShipmentId,
  shipmentType: SHIPMENT_OPTIONS.PPM,
  ppmShipment: {
    id: mockPPMShipmentId,
    actualMoveDate: '2022-05-01',
    actualPickupPostalCode: '10003',
    actualDestinationPostalCode: '10004',
    advanceReceived: true,
    advanceAmountReceived: '6000000',
    pickupPostalCode: '10001',
    destinationPostalCode: '10002',
    expectedDepartureDate: '2022-04-30',
    advanceRequested: true,
    advanceAmountRequested: 598700,
    estimatedWeight: 4000,
    estimatedIncentive: 1000000,
    sitExpected: false,
    hasProGear: false,
    proGearWeight: null,
    spouseProGearWeight: null,
    weightTickets: [],
  },
  eTag: 'dGVzdGluZzIzNDQzMjQ',
};

const weightTicketOne = createCompleteWeightTicket();
const weightTicketTwo = createCompleteWeightTicket();
const mockMTOShipmentWithWeightTicket = {
  id: mockMTOShipmentId,
  shipmentType: SHIPMENT_OPTIONS.PPM,
  ppmShipment: {
    id: mockPPMShipmentId,
    actualMoveDate: '2022-05-01',
    actualPickupPostalCode: '10003',
    actualDestinationPostalCode: '10004',
    advanceReceived: true,
    advanceAmountReceived: '6000000',
    pickupPostalCode: '10001',
    destinationPostalCode: '10002',
    expectedDepartureDate: '2022-04-30',
    advanceRequested: true,
    advanceAmountRequested: 598700,
    estimatedWeight: 4000,
    estimatedIncentive: 1000000,
    sitExpected: false,
    hasProGear: false,
    proGearWeight: null,
    spouseProGearWeight: null,
    weightTickets: [weightTicketOne, weightTicketTwo],
  },
  eTag: 'dGVzdGluZzIzNDQzMjQ',
};

const mockMTOShipmentWithWeightTicketDeleted = {
  mtoShipments: {
    [mockMTOShipmentId]: {
      id: mockMTOShipmentId,
      shipmentType: SHIPMENT_OPTIONS.PPM,
      ppmShipment: {
        id: mockPPMShipmentId,
        actualMoveDate: '2022-05-01',
        actualPickupPostalCode: '10003',
        actualDestinationPostalCode: '10004',
        advanceReceived: true,
        advanceAmountReceived: '6000000',
        pickupPostalCode: '10001',
        destinationPostalCode: '10002',
        expectedDepartureDate: '2022-04-30',
        advanceRequested: true,
        advanceAmountRequested: 598700,
        estimatedWeight: 4000,
        estimatedIncentive: 1000000,
        sitExpected: false,
        hasProGear: false,
        proGearWeight: null,
        spouseProGearWeight: null,
        weightTickets: [weightTicketTwo],
      },
      eTag: 'dGVzdGluZzIzNDQzMjQ',
    },
  },
};

const mockMTOShipmentWithIncompleteWeightTicket = {
  id: mockMTOShipmentId,
  shipmentType: SHIPMENT_OPTIONS.PPM,
  ppmShipment: {
    id: mockPPMShipmentId,
    actualMoveDate: '2022-05-01',
    actualPickupPostalCode: '10003',
    actualDestinationPostalCode: '10004',
    advanceReceived: true,
    advanceAmountReceived: '6000000',
    pickupPostalCode: '10001',
    destinationPostalCode: '10002',
    expectedDepartureDate: '2022-04-30',
    advanceRequested: true,
    advanceAmountRequested: 598700,
    estimatedWeight: 4000,
    estimatedIncentive: 1000000,
    sitExpected: false,
    hasProGear: false,
    proGearWeight: null,
    spouseProGearWeight: null,
    weightTickets: [createBaseWeightTicket()],
  },
  eTag: 'dGVzdGluZzIzNDQzMjQ',
};

const proGearWeightOne = createBaseProGearWeightTicket();
const mockMTOShipmentWithProGear = {
  id: mockMTOShipmentId,
  shipmentType: SHIPMENT_OPTIONS.PPM,
  ppmShipment: {
    id: mockPPMShipmentId,
    actualMoveDate: '2022-05-01',
    actualPickupPostalCode: '10003',
    actualDestinationPostalCode: '10004',
    advanceReceived: true,
    advanceAmountReceived: '6000000',
    pickupPostalCode: '10001',
    destinationPostalCode: '10002',
    expectedDepartureDate: '2022-04-30',
    advanceRequested: true,
    advanceAmountRequested: 598700,
    estimatedWeight: 4000,
    estimatedIncentive: 1000000,
    sitExpected: false,
    hasProGear: true,
    proGearWeight: 100,
    spouseProGearWeight: null,
    proGearWeightTickets: [proGearWeightOne],
  },
  eTag: 'dGVzdGluZzIzNDQzMjQ',
};

const mockMTOShipmentWithProGearDeleted = {
  mtoShipments: {
    [mockMTOShipmentId]: {
      id: mockMTOShipmentId,
      shipmentType: SHIPMENT_OPTIONS.PPM,
      ppmShipment: {
        id: mockPPMShipmentId,
        actualMoveDate: '2022-05-01',
        actualPickupPostalCode: '10003',
        actualDestinationPostalCode: '10004',
        advanceReceived: true,
        advanceAmountReceived: '6000000',
        pickupPostalCode: '10001',
        destinationPostalCode: '10002',
        expectedDepartureDate: '2022-04-30',
        advanceRequested: true,
        advanceAmountRequested: 598700,
        estimatedWeight: 4000,
        estimatedIncentive: 1000000,
        sitExpected: false,
        hasProGear: true,
        proGearWeight: 100,
        spouseProGearWeight: null,
        proGearWeightTickets: [],
      },
      eTag: 'dGVzdGluZzIzNDQzMjQ',
    },
  },
};

const expenseOne = createCompleteMovingExpense();
const expenseTwo = createCompleteSITMovingExpense();
const mockMTOShipmentWithExpenses = {
  id: mockMTOShipmentId,
  shipmentType: SHIPMENT_OPTIONS.PPM,
  ppmShipment: {
    id: mockPPMShipmentId,
    actualMoveDate: '2022-05-01',
    actualPickupPostalCode: '10003',
    actualDestinationPostalCode: '10004',
    advanceReceived: true,
    advanceAmountReceived: '6000000',
    pickupPostalCode: '10001',
    destinationPostalCode: '10002',
    expectedDepartureDate: '2022-04-30',
    advanceRequested: true,
    advanceAmountRequested: 598700,
    estimatedWeight: 4000,
    estimatedIncentive: 1000000,
    sitExpected: false,
    hasProGear: true,
    proGearWeight: 100,
    spouseProGearWeight: null,
    movingExpenses: [expenseOne, expenseTwo],
  },
  eTag: 'dGVzdGluZzIzNDQzMjQ',
};

const mockMTOShipmentWithExpensesDeleted = {
  mtoShipments: {
    [mockMTOShipmentId]: {
      id: mockMTOShipmentId,
      shipmentType: SHIPMENT_OPTIONS.PPM,
      ppmShipment: {
        id: mockPPMShipmentId,
        actualMoveDate: '2022-05-01',
        actualPickupPostalCode: '10003',
        actualDestinationPostalCode: '10004',
        advanceReceived: true,
        advanceAmountReceived: '6000000',
        pickupPostalCode: '10001',
        destinationPostalCode: '10002',
        expectedDepartureDate: '2022-04-30',
        advanceRequested: true,
        advanceAmountRequested: 598700,
        estimatedWeight: 4000,
        estimatedIncentive: 1000000,
        sitExpected: false,
        hasProGear: true,
        proGearWeight: 100,
        spouseProGearWeight: null,
        movingExpenses: [expenseOne],
      },
      eTag: 'dGVzdGluZzIzNDQzMjQ',
    },
  },
};
const mockPush = jest.fn();
jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom'),
  useHistory: () => ({
    push: mockPush,
  }),
  useParams: () => ({
    moveId: mockMoveId,
    mtoShipmentId: mockMTOShipmentId,
  }),
}));

jest.mock('services/internalApi', () => ({
  ...jest.requireActual('services/internalApi'),
  deleteWeightTicket: jest.fn(() => {}),
  deleteProGearWeightTicket: jest.fn(() => {}),
  deleteMovingExpense: jest.fn(() => {}),
  getMTOShipmentsForMove: jest.fn(),
}));

jest.mock('store/entities/selectors', () => ({
  ...jest.requireActual('store/entities/selectors'),
  selectMTOShipmentById: jest.fn(() => mockMTOShipment),
}));

beforeEach(() => {
  jest.clearAllMocks();
});

describe('Review page', () => {
  it('loads the selected shipment from redux', () => {
    render(<Review />, { wrapper: MockProviders });

    expect(selectMTOShipmentById).toHaveBeenCalledWith(expect.anything(), mockMTOShipmentId);
  });

  it('renders the page headings', () => {
    render(<Review />, { wrapper: MockProviders });

    expect(screen.getByTestId('tag')).toHaveTextContent('PPM');
    expect(screen.getByRole('heading', { level: 1 })).toHaveTextContent('Review');
    expect(screen.getAllByRole('heading', { level: 2 })[0]).toHaveTextContent('About Your PPM');
    expect(screen.getAllByRole('heading', { level: 2 })[1]).toHaveTextContent('Documents');
    expect(screen.getAllByRole('heading', { level: 3 })[0]).toHaveTextContent('Weight moved');
    expect(screen.getAllByRole('heading', { level: 3 })[1]).toHaveTextContent('Pro-gear');
    expect(screen.getAllByRole('heading', { level: 3 })[2]).toHaveTextContent('Expenses');
  });

  it('renders the empty message when there are no weight tickets', () => {
    render(<Review />, { wrapper: MockProviders });

    expect(
      screen.getByText('No weight moved documented. At least one trip is required to continue.'),
    ).toBeInTheDocument();
  });

  it('routes to the edit about your ppm page when the edit link is clicked', async () => {
    const editAboutYourPPM = generatePath(customerRoutes.SHIPMENT_PPM_ABOUT_PATH, {
      moveId: mockMoveId,
      mtoShipmentId: mockMTOShipmentId,
    });

    const { memoryHistory, mockProviderWithHistory } = setUpProvidersWithHistory();

    render(<Review />, { wrapper: mockProviderWithHistory });

    await userEvent.click(screen.getAllByText('Edit')[0]);

    await waitFor(() => {
      expect(memoryHistory.location.pathname).toEqual(editAboutYourPPM);
    });
  });

  it('routes to the add weight ticket page when the add link is clicked', async () => {
    const newWeightTicket = generatePath(customerRoutes.SHIPMENT_PPM_WEIGHT_TICKETS_PATH, {
      moveId: mockMoveId,
      mtoShipmentId: mockMTOShipmentId,
    });

    const { memoryHistory, mockProviderWithHistory } = setUpProvidersWithHistory();

    render(<Review />, { wrapper: mockProviderWithHistory });

    await userEvent.click(screen.getByText('Add More Weight'));

    await waitFor(() => {
      expect(memoryHistory.location.pathname).toEqual(newWeightTicket);
    });
  });

  it('routes to the edit weight ticket page when the edit link is clicked', async () => {
    selectMTOShipmentById.mockImplementation(() => mockMTOShipmentWithWeightTicket);
    const editWeightTicket = generatePath(customerRoutes.SHIPMENT_PPM_WEIGHT_TICKETS_EDIT_PATH, {
      moveId: mockMoveId,
      mtoShipmentId: mockMTOShipmentId,
      weightTicketId: mockMTOShipmentWithWeightTicket.ppmShipment.weightTickets[0].id,
    });

    const { memoryHistory, mockProviderWithHistory } = setUpProvidersWithHistory();

    render(<Review />, { wrapper: mockProviderWithHistory });

    await userEvent.click(screen.getAllByText('Edit')[1]);

    await waitFor(() => {
      expect(memoryHistory.location.pathname).toEqual(editWeightTicket);
    });
  });

  it('routes to the add pro-gear page when the add link is clicked', async () => {
    const newProGear = generatePath(customerRoutes.SHIPMENT_PPM_PRO_GEAR_PATH, {
      moveId: mockMoveId,
      mtoShipmentId: mockMTOShipmentId,
    });

    const { memoryHistory, mockProviderWithHistory } = setUpProvidersWithHistory();

    render(<Review />, { wrapper: mockProviderWithHistory });

    await userEvent.click(screen.getByText('Add Pro-gear Weight'));

    await waitFor(() => {
      expect(memoryHistory.location.pathname).toEqual(newProGear);
    });
  });

  it('routes to the edit pro-gear page when the edit link is clicked', async () => {
    selectMTOShipmentById.mockImplementation(() => mockMTOShipmentWithProGear);
    const editProGearWeightTicket = generatePath(customerRoutes.SHIPMENT_PPM_PRO_GEAR_EDIT_PATH, {
      moveId: mockMoveId,
      mtoShipmentId: mockMTOShipmentId,
      proGearId: mockMTOShipmentWithProGear.ppmShipment.proGearWeightTickets[0].id,
    });

    const { memoryHistory, mockProviderWithHistory } = setUpProvidersWithHistory();

    render(<Review />, { wrapper: mockProviderWithHistory });

    await userEvent.click(screen.getAllByText('Edit')[1]);

    await waitFor(() => {
      expect(memoryHistory.location.pathname).toEqual(editProGearWeightTicket);
    });
  });

  it('routes to the add expenses page when the add link is clicked', async () => {
    const newExpensePath = generatePath(customerRoutes.SHIPMENT_PPM_EXPENSES_PATH, {
      moveId: mockMoveId,
      mtoShipmentId: mockMTOShipmentId,
    });

    const { memoryHistory, mockProviderWithHistory } = setUpProvidersWithHistory();

    render(<Review />, { wrapper: mockProviderWithHistory });

    await userEvent.click(screen.getByText('Add Expenses'));

    await waitFor(() => {
      expect(memoryHistory.location.pathname).toEqual(newExpensePath);
    });
  });

  it('routes to the edit expense page when the edit link is clicked', async () => {
    selectMTOShipmentById.mockImplementation(() => mockMTOShipmentWithExpenses);
    const editExpensePath = generatePath(customerRoutes.SHIPMENT_PPM_EXPENSES_EDIT_PATH, {
      moveId: mockMoveId,
      mtoShipmentId: mockMTOShipmentId,
      expenseId: mockMTOShipmentWithExpenses.ppmShipment.movingExpenses[0].id,
    });

    const { memoryHistory, mockProviderWithHistory } = setUpProvidersWithHistory();

    render(<Review />, { wrapper: mockProviderWithHistory });

    await userEvent.click(screen.getAllByText('Edit')[1]);

    await waitFor(() => {
      expect(memoryHistory.location.pathname).toEqual(editExpensePath);
    });
  });

  it('routes to the home page when the return to homepage link is clicked', async () => {
    const { memoryHistory, mockProviderWithHistory } = setUpProvidersWithHistory();

    render(<Review />, { wrapper: mockProviderWithHistory });

    await userEvent.click(screen.getByText('Return To Homepage'));

    await waitFor(() => {
      expect(memoryHistory.location.pathname).toEqual('/');
    });
  });

  it('routes to the complete page when the save and continue link is clicked', async () => {
    selectMTOShipmentById.mockImplementationOnce(() => mockMTOShipmentWithWeightTicket);
    const completePath = generatePath(customerRoutes.SHIPMENT_PPM_COMPLETE_PATH, {
      moveId: mockMoveId,
      mtoShipmentId: mockMTOShipmentId,
    });

    const { memoryHistory, mockProviderWithHistory } = setUpProvidersWithHistory();

    render(<Review />, { wrapper: mockProviderWithHistory });

    await userEvent.click(screen.getByText('Save & Continue'));

    await waitFor(() => {
      expect(memoryHistory.location.pathname).toEqual(completePath);
    });
  });

  it('disables the save and continue link when there are no weight tickets', async () => {
    selectMTOShipmentById.mockImplementationOnce(() => mockMTOShipment);
    render(<Review />, { wrapper: MockProviders });

    expect(screen.getByText('Save & Continue')).toHaveClass('usa-button--disabled');
    expect(screen.getByText('Save & Continue')).toHaveAttribute('aria-disabled', 'true');
  });

  it('disables the save and continue link when there is an incomplete weight ticket', async () => {
    selectMTOShipmentById.mockImplementationOnce(() => mockMTOShipmentWithIncompleteWeightTicket);
    render(<Review />, { wrapper: MockProviders });

    expect(screen.getByText('Save & Continue')).toHaveClass('usa-button--disabled');
    expect(screen.getByText('Save & Continue')).toHaveAttribute('aria-disabled', 'true');
  });

  it('error message is displayed when a PPM shipment is in an incomplete state', async () => {
    selectMTOShipmentById.mockImplementationOnce(() => mockMTOShipmentWithIncompleteWeightTicket);
    render(<Review />, { wrapper: MockProviders });

    expect(
      screen.getByText(
        'There are items below that are missing required information. Please select “Edit” to enter all required information or “Delete” to remove the item.',
      ),
    ).toBeInTheDocument();

    expect(screen.getByText('This trip is missing required information.')).toBeInTheDocument();
  });

  it('displays the delete confirmation modal when the delete button for Weight Moved/Trip 2 is clicked', async () => {
    selectMTOShipmentById.mockImplementation(() => mockMTOShipmentWithWeightTicket);
    render(<Review />, { wrapper: MockProviders });

    await userEvent.click(screen.getAllByRole('button', { name: 'Delete' })[1]);

    await waitFor(() => {
      expect(screen.getByRole('heading', { level: 3, name: 'Delete this?' })).toBeInTheDocument();
      expect(screen.getByText('You are about to delete Trip 2. This cannot be undone.')).toBeInTheDocument();
    });

    await userEvent.click(screen.getByRole('button', { name: 'No, Keep It' }));

    expect(screen.queryByRole('heading', { level: 3, name: 'Delete this?' })).not.toBeInTheDocument();
  });

  it('calls the delete weight ticket api when confirm is clicked', async () => {
    selectMTOShipmentById.mockImplementation(() => mockMTOShipmentWithWeightTicket);
    const mockDeleteWeightTicket = jest.fn().mockResolvedValue({});
    deleteWeightTicket.mockImplementationOnce(mockDeleteWeightTicket);
    getMTOShipmentsForMove.mockResolvedValue(mockMTOShipmentWithWeightTicketDeleted);
    render(<Review />, { wrapper: MockProviders });

    await userEvent.click(screen.getAllByRole('button', { name: 'Delete' })[0]);

    await waitFor(() => {
      expect(screen.getByRole('heading', { level: 3, name: 'Delete this?' })).toBeInTheDocument();
    });

    await userEvent.click(screen.getByRole('button', { name: 'Yes, Delete' }));

    const weightTicket = mockMTOShipmentWithWeightTicket.ppmShipment.weightTickets[0];
    await waitFor(() => {
      expect(mockDeleteWeightTicket).toHaveBeenCalledWith(
        mockMTOShipmentWithWeightTicket.ppmShipment.id,
        weightTicket.id,
      );
    });
    await waitFor(() => {
      expect(screen.getByText('Trip 1 successfully deleted.'));
    });
  });

  it('calls the delete progear weight ticket api when confirm is clicked', async () => {
    selectMTOShipmentById.mockImplementation(() => mockMTOShipmentWithProGear);
    const mockDeleteProGearWeightTicket = jest.fn().mockResolvedValue({});
    deleteProGearWeightTicket.mockImplementationOnce(mockDeleteProGearWeightTicket);
    getMTOShipmentsForMove.mockResolvedValue(mockMTOShipmentWithProGearDeleted);
    render(<Review />, { wrapper: MockProviders });

    await userEvent.click(screen.getAllByRole('button', { name: 'Delete' })[0]);

    await waitFor(() => {
      expect(screen.getByRole('heading', { level: 3, name: 'Delete this?' })).toBeInTheDocument();
    });

    await userEvent.click(screen.getByRole('button', { name: 'Yes, Delete' }));

    const proGearWeightTicket = mockMTOShipmentWithProGear.ppmShipment.proGearWeightTickets[0];
    await waitFor(() => {
      expect(mockDeleteProGearWeightTicket).toHaveBeenCalledWith(
        mockMTOShipmentWithWeightTicket.ppmShipment.id,
        proGearWeightTicket.id,
      );
    });

    await waitFor(() => {
      expect(screen.getByText('Set 1 successfully deleted.'));
    });
  });

  it('calls the delete moving expense api when confirm is clicked', async () => {
    selectMTOShipmentById.mockImplementation(() => mockMTOShipmentWithExpenses);
    const mockDeleteMovingExpense = jest.fn().mockResolvedValue({});
    deleteMovingExpense.mockImplementationOnce(mockDeleteMovingExpense);
    getMTOShipmentsForMove.mockResolvedValue(mockMTOShipmentWithExpensesDeleted);
    render(<Review />, { wrapper: MockProviders });

    await userEvent.click(screen.getAllByRole('button', { name: 'Delete' })[0]);

    await waitFor(() => {
      expect(screen.getByRole('heading', { level: 3, name: 'Delete this?' })).toBeInTheDocument();
    });

    await userEvent.click(screen.getByRole('button', { name: 'Yes, Delete' }));

    const movingExpense = mockMTOShipmentWithExpenses.ppmShipment.movingExpenses[0];
    await waitFor(() => {
      expect(mockDeleteMovingExpense).toHaveBeenCalledWith(
        mockMTOShipmentWithWeightTicket.ppmShipment.id,
        movingExpense.id,
      );
    });
    await waitFor(() => {
      expect(screen.getByText('Receipt 1 successfully deleted.'));
    });
  });
});
