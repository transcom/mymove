import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import { v4 } from 'uuid';
import userEvent from '@testing-library/user-event';
import { generatePath } from 'react-router-dom';

import { MockProviders } from 'testUtils';
import Review from 'pages/Office/PPM/Closeout/Review/Review';
import { SHIPMENT_OPTIONS } from 'shared/constants';
import { servicesCounselingRoutes } from 'constants/routes';
import { deleteWeightTicket, deleteMovingExpense } from 'services/ghcApi';
import { createBaseWeightTicket, createCompleteWeightTicket } from 'utils/test/factories/weightTicket';
import { createBaseProGearWeightTicket } from 'utils/test/factories/proGearWeightTicket';
import { createCompleteMovingExpense, createCompleteSITMovingExpense } from 'utils/test/factories/movingExpense';
import { usePPMShipmentAndDocsOnlyQueries } from 'hooks/queries';

const mockMoveId = v4();
const mockMTOShipmentId = v4();
const mockPPMShipmentId = v4();

const pickupAddress = {
  id: 'test1',
  streetAddress1: 'Pickup Road',
  city: 'PPM City',
  state: 'CA',
  postalCode: '90210',
};

const destinationAddress = {
  id: 'test1',
  streetAddress1: 'Destination Road',
  city: 'PPM City',
  state: 'CA',
  postalCode: '90210',
};

const mockMTOShipment = {
  id: mockMTOShipmentId,
  shipmentType: SHIPMENT_OPTIONS.PPM,
  ppmShipment: {
    id: mockPPMShipmentId,
    actualMoveDate: '2022-05-01',
    advanceReceived: true,
    advanceAmountReceived: '6000000',
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
    pickupAddress,
    destinationAddress,
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
    advanceReceived: true,
    advanceAmountReceived: '6000000',
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
    pickupAddress,
    destinationAddress,
  },
  eTag: 'dGVzdGluZzIzNDQzMjQ',
};

const mockDocumentsWithWeightTickets = {
  WeightTickets: [weightTicketOne, weightTicketTwo],
  ProGearWeightTickets: [],
  MovingExpenses: [],
};

const mockMTOShipmentWithIncompleteWeightTicket = {
  id: mockMTOShipmentId,
  shipmentType: SHIPMENT_OPTIONS.PPM,
  ppmShipment: {
    id: mockPPMShipmentId,
    actualMoveDate: '2022-05-01',
    advanceReceived: true,
    advanceAmountReceived: '6000000',
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
    pickupAddress,
    destinationAddress,
  },
  eTag: 'dGVzdGluZzIzNDQzMjQ',
};
const mockDocumentsWithIncompleteWeightTicket = {
  WeightTickets: [createBaseWeightTicket()],
  ProGearWeightTickets: [],
  MovingExpenses: [],
};

const proGearWeightOne = createBaseProGearWeightTicket();
const mockMTOShipmentWithProGear = {
  id: mockMTOShipmentId,
  shipmentType: SHIPMENT_OPTIONS.PPM,
  ppmShipment: {
    id: mockPPMShipmentId,
    actualMoveDate: '2022-05-01',
    advanceReceived: true,
    advanceAmountReceived: '6000000',
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
    pickupAddress,
    destinationAddress,
  },
  eTag: 'dGVzdGluZzIzNDQzMjQ',
};

const expenseOne = createCompleteMovingExpense();
const expenseTwo = createCompleteSITMovingExpense();
const mockMTOShipmentWithExpenses = {
  id: mockMTOShipmentId,
  shipmentType: SHIPMENT_OPTIONS.PPM,
  ppmShipment: {
    id: mockPPMShipmentId,
    actualMoveDate: '2022-05-01',
    advanceReceived: true,
    advanceAmountReceived: '6000000',
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
    pickupAddress,
    destinationAddress,
  },
  eTag: 'dGVzdGluZzIzNDQzMjQ',
};

const mockDocumentsWithExpenses = {
  WeightTickets: [weightTicketOne, weightTicketTwo],
  ProGearWeightTickets: [],
  MovingExpenses: [expenseOne, expenseTwo],
};

const mockNavigate = jest.fn();
jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom'),
  useNavigate: () => mockNavigate,
}));

jest.mock('hooks/queries', () => ({
  usePPMShipmentAndDocsOnlyQueries: jest.fn(),
}));

jest.mock('services/ghcApi', () => ({
  ...jest.requireActual('services/ghcApi'),
  deleteWeightTicket: jest.fn(),
  deleteMovingExpense: jest.fn(),
}));

beforeEach(() => {
  jest.clearAllMocks();
});

const newWeightPath = generatePath(servicesCounselingRoutes.BASE_SHIPMENT_PPM_WEIGHT_TICKETS_PATH, {
  moveCode: mockMoveId,
  shipmentId: mockMTOShipmentId,
});
const editAboutYourPPMPath = generatePath(servicesCounselingRoutes.BASE_SHIPMENT_PPM_ABOUT_PATH, {
  moveCode: mockMoveId,
  shipmentId: mockMTOShipmentId,
});
const editWeightPath = generatePath(servicesCounselingRoutes.BASE_SHIPMENT_PPM_WEIGHT_TICKETS_EDIT_PATH, {
  moveCode: mockMoveId,
  shipmentId: mockMTOShipmentId,
  weightTicketId: mockMTOShipmentWithWeightTicket.ppmShipment.weightTickets[0].id,
});
const newProGearPath = generatePath(servicesCounselingRoutes.BASE_SHIPMENT_PPM_PRO_GEAR_PATH, {
  moveCode: mockMoveId,
  shipmentId: mockMTOShipmentId,
});
const editProGearWeightTicket = generatePath(servicesCounselingRoutes.BASE_SHIPMENT_PPM_PRO_GEAR_EDIT_PATH, {
  moveCode: mockMoveId,
  shipmentId: mockMTOShipmentId,
  proGearId: mockMTOShipmentWithProGear.ppmShipment.proGearWeightTickets[0].id,
});
const newExpensePath = generatePath(servicesCounselingRoutes.BASE_SHIPMENT_PPM_EXPENSES_PATH, {
  moveCode: mockMoveId,
  shipmentId: mockMTOShipmentId,
});
const editExpensePath = generatePath(servicesCounselingRoutes.BASE_SHIPMENT_PPM_EXPENSES_EDIT_PATH, {
  moveCode: mockMoveId,
  shipmentId: mockMTOShipmentId,
  expenseId: mockMTOShipmentWithExpenses.ppmShipment.movingExpenses[0].id,
});

const completePath = generatePath(servicesCounselingRoutes.BASE_SHIPMENT_PPM_COMPLETE_PATH, {
  moveCode: mockMoveId,
  shipmentId: mockMTOShipmentId,
});

const moveDetailsPath = generatePath(servicesCounselingRoutes.BASE_MOVE_VIEW_PATH, { moveCode: mockMoveId });

const mockRoutes = [
  {
    path: newWeightPath,
    element: <div>Add More Weight Page</div>,
  },
  {
    path: editAboutYourPPMPath,
    element: <div>Edit About Your PPM Page</div>,
  },
  {
    path: editWeightPath,
    element: <div>Edit Weight Page</div>,
  },
  {
    path: newProGearPath,
    element: <div>New Pro Gear Page</div>,
  },
  {
    path: editProGearWeightTicket,
    element: <div>Edit Pro Gear Weight Ticket Page</div>,
  },
  {
    path: newExpensePath,
    element: <div>New Expense Page</div>,
  },
  {
    path: editExpensePath,
    element: <div>Edit Expense Page</div>,
  },
  {
    path: completePath,
    element: <div>Complete Page</div>,
  },
];

const renderReviewPage = (props) => {
  return render(
    <MockProviders
      path={servicesCounselingRoutes.BASE_SHIPMENT_PPM_REVIEW_PATH}
      params={{
        moveCode: mockMoveId,
        shipmentId: mockMTOShipmentId,
      }}
      routes={mockRoutes}
    >
      <Review {...props} />
    </MockProviders>,
  );
};

describe('Review page', () => {
  it('renders the page headings', () => {
    usePPMShipmentAndDocsOnlyQueries.mockReturnValue({
      isLoading: false,
      mtoShipment: mockMTOShipment,
      error: null,
    });

    renderReviewPage();

    expect(screen.getByTestId('tag')).toHaveTextContent('PPM');
    expect(screen.getByRole('heading', { level: 1 })).toHaveTextContent('Review');
    expect(screen.getAllByRole('heading', { level: 2 })[0]).toHaveTextContent('About Your PPM');
    expect(screen.getAllByRole('heading', { level: 2 })[1]).toHaveTextContent('Documents');
    expect(screen.getAllByRole('heading', { level: 3 })[0]).toHaveTextContent('Weight moved');
    expect(screen.getAllByRole('heading', { level: 3 })[1]).toHaveTextContent('Pro-gear');
    expect(screen.getAllByRole('heading', { level: 3 })[2]).toHaveTextContent('Expenses');
  });

  it('renders the empty message when there are no weight tickets', () => {
    usePPMShipmentAndDocsOnlyQueries.mockReturnValue({
      isLoading: false,
      mtoShipment: mockMTOShipmentWithIncompleteWeightTicket,
      error: null,
    });
    renderReviewPage();

    expect(
      screen.getByText('No weight moved documented. At least one trip is required to continue.'),
    ).toBeInTheDocument();
  });

  it('routes to the edit about your ppm page when the edit link is clicked', async () => {
    usePPMShipmentAndDocsOnlyQueries.mockReturnValue({
      isLoading: false,
      mtoShipment: mockMTOShipment,
      error: null,
    });
    renderReviewPage();

    await userEvent.click(screen.getAllByText('Edit')[0]);

    await waitFor(() => {
      expect(screen.getByText('Edit About Your PPM Page')).toBeInTheDocument();
    });
  });

  it('routes to the Move Details page when the Back button is clicked', async () => {
    usePPMShipmentAndDocsOnlyQueries.mockReturnValue({
      isLoading: false,
      mtoShipment: mockMTOShipment,
      error: null,
    });
    renderReviewPage();

    await userEvent.click(screen.getByTestId('formBackButton'));

    await waitFor(() => {
      expect(mockNavigate).toHaveBeenCalledWith(moveDetailsPath);
    });
  });

  it('disables the save and continue link when there are no weight tickets', async () => {
    usePPMShipmentAndDocsOnlyQueries.mockReturnValue({
      isLoading: false,
      mtoShipment: mockMTOShipmentWithIncompleteWeightTicket,
      documents: mockDocumentsWithIncompleteWeightTicket,
      error: null,
    });
    renderReviewPage();

    expect(screen.getByTestId('saveAndContinueButton')).toBeDisabled();
  });

  it('disables the save and continue link when there is an incomplete weight ticket', async () => {
    usePPMShipmentAndDocsOnlyQueries.mockReturnValue({
      isLoading: false,
      mtoShipment: mockMTOShipmentWithIncompleteWeightTicket,
      documents: mockDocumentsWithIncompleteWeightTicket,
      error: null,
    });
    renderReviewPage();

    expect(screen.getByTestId('saveAndContinueButton')).toBeDisabled();
  });

  it('error message is displayed when a PPM shipment is in an incomplete state', async () => {
    usePPMShipmentAndDocsOnlyQueries.mockReturnValue({
      isLoading: false,
      mtoShipment: mockMTOShipmentWithIncompleteWeightTicket,
      documents: mockDocumentsWithIncompleteWeightTicket,
      error: null,
    });
    renderReviewPage();

    expect(
      screen.getByText(
        'There are items below that are missing required information. Please select “Edit” to enter all required information or “Delete” to remove the item.',
      ),
    ).toBeInTheDocument();

    expect(screen.getByText('This trip is missing required information.')).toBeInTheDocument();
  });

  it('displays the delete confirmation modal when the delete button for Weight Moved/Trip 2 is clicked', async () => {
    usePPMShipmentAndDocsOnlyQueries.mockReturnValue({
      isLoading: false,
      mtoShipment: mockMTOShipmentWithWeightTicket,
      documents: mockDocumentsWithWeightTickets,
      error: null,
    });
    renderReviewPage();

    await userEvent.click(screen.getAllByRole('button', { name: 'Delete' })[1]);

    await waitFor(() => {
      expect(screen.getByRole('heading', { level: 3, name: 'Delete this?' })).toBeInTheDocument();
      expect(screen.getByText('You are about to delete Trip 2. This cannot be undone.')).toBeInTheDocument();
    });

    await userEvent.click(screen.getByRole('button', { name: 'No, Keep It' }));

    expect(screen.queryByRole('heading', { level: 3, name: 'Delete this?' })).not.toBeInTheDocument();
  });

  it('calls the delete weight ticket api when confirm is clicked', async () => {
    const mockDeleteWeightTicket = jest.fn().mockResolvedValue({});
    deleteWeightTicket.mockImplementationOnce(mockDeleteWeightTicket);

    renderReviewPage();

    await userEvent.click(screen.getAllByRole('button', { name: 'Delete' })[0]);

    await waitFor(() => {
      expect(screen.getByRole('heading', { level: 3, name: 'Delete this?' })).toBeInTheDocument();
    });

    await userEvent.click(screen.getByRole('button', { name: 'Yes, Delete' }));

    const weightTicket = mockMTOShipmentWithWeightTicket.ppmShipment.weightTickets[0];
    await waitFor(() => {
      expect(mockDeleteWeightTicket).toHaveBeenCalledWith({
        ppmShipmentId: mockMTOShipmentWithWeightTicket.ppmShipment.id,
        weightTicketId: weightTicket.id,
      });
    });
    await waitFor(() => {
      expect(screen.getByText('Trip 1 successfully deleted.'));
    });
  });

  it('displays the delete confirmation modal when the delete button for Weight Expenses is clicked', async () => {
    usePPMShipmentAndDocsOnlyQueries.mockReturnValue({
      isLoading: false,
      mtoShipment: mockMTOShipmentWithExpenses,
      documents: mockDocumentsWithExpenses,
      error: null,
    });
    renderReviewPage();

    await userEvent.click(screen.getAllByRole('button', { name: 'Delete' })[3]);

    await waitFor(() => {
      expect(screen.getByRole('heading', { level: 3, name: 'Delete this?' })).toBeInTheDocument();
      expect(screen.getByText('You are about to delete Receipt 2. This cannot be undone.')).toBeInTheDocument();
    });

    await userEvent.click(screen.getByRole('button', { name: 'No, Keep It' }));

    expect(screen.queryByRole('heading', { level: 3, name: 'Delete this?' })).not.toBeInTheDocument();
  });

  it('calls the delete expenses api when confirm is clicked', async () => {
    const mockDeleteMovingExpense = jest.fn().mockResolvedValue({});
    deleteMovingExpense.mockImplementationOnce(mockDeleteMovingExpense);

    renderReviewPage();

    await userEvent.click(screen.getAllByRole('button', { name: 'Delete' })[2]);

    await waitFor(() => {
      expect(screen.getByRole('heading', { level: 3, name: 'Delete this?' })).toBeInTheDocument();
    });

    await userEvent.click(screen.getByRole('button', { name: 'Yes, Delete' }));

    const movingExpense = mockMTOShipmentWithExpenses.ppmShipment.movingExpenses[0];
    await waitFor(() => {
      expect(mockDeleteMovingExpense).toHaveBeenCalledWith({
        ppmShipmentId: mockMTOShipmentWithExpenses.ppmShipment.id,
        movingExpenseId: movingExpense.id,
      });
    });
    await waitFor(() => {
      expect(screen.getByText('Receipt 1 successfully deleted.'));
    });
  });
});
