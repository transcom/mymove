import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import { v4 } from 'uuid';
import userEvent from '@testing-library/user-event';
import { generatePath } from 'react-router-dom';

import { MockProviders } from 'testUtils';
import { selectMTOShipmentById } from 'store/entities/selectors';
import Review from 'pages/MyMove/PPM/Closeout/Review/Review';
import { SHIPMENT_OPTIONS } from 'shared/constants';
import { customerRoutes } from 'constants/routes';
import { deleteWeightTicket } from 'services/internalApi';
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
    weightTickets: [createCompleteWeightTicket()],
  },
  eTag: 'dGVzdGluZzIzNDQzMjQ',
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
    proGearWeightTickets: [createBaseProGearWeightTicket()],
  },
  eTag: 'dGVzdGluZzIzNDQzMjQ',
};

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
    movingExpenses: [createCompleteMovingExpense(), createCompleteSITMovingExpense()],
  },
  eTag: 'dGVzdGluZzIzNDQzMjQ',
};

const mockNavigate = jest.fn();
jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom'),
  useNavigate: () => mockNavigate,
}));

jest.mock('services/internalApi', () => ({
  ...jest.requireActual('services/internalApi'),
  deleteWeightTicket: jest.fn(() => {}),
}));

jest.mock('store/entities/selectors', () => ({
  ...jest.requireActual('store/entities/selectors'),
  selectMTOShipmentById: jest.fn(() => mockMTOShipment),
}));

beforeEach(() => {
  jest.clearAllMocks();
});

const newWeightPath = generatePath(customerRoutes.SHIPMENT_PPM_WEIGHT_TICKETS_PATH, {
  moveId: mockMoveId,
  mtoShipmentId: mockMTOShipmentId,
});
const editAboutYourPPMPath = generatePath(customerRoutes.SHIPMENT_PPM_ABOUT_PATH, {
  moveId: mockMoveId,
  mtoShipmentId: mockMTOShipmentId,
});
const editWeightPath = generatePath(customerRoutes.SHIPMENT_PPM_WEIGHT_TICKETS_EDIT_PATH, {
  moveId: mockMoveId,
  mtoShipmentId: mockMTOShipmentId,
  weightTicketId: mockMTOShipmentWithWeightTicket.ppmShipment.weightTickets[0].id,
});
const newProGearPath = generatePath(customerRoutes.SHIPMENT_PPM_PRO_GEAR_PATH, {
  moveId: mockMoveId,
  mtoShipmentId: mockMTOShipmentId,
});
const editProGearWeightTicket = generatePath(customerRoutes.SHIPMENT_PPM_PRO_GEAR_EDIT_PATH, {
  moveId: mockMoveId,
  mtoShipmentId: mockMTOShipmentId,
  proGearId: mockMTOShipmentWithProGear.ppmShipment.proGearWeightTickets[0].id,
});
const newExpensePath = generatePath(customerRoutes.SHIPMENT_PPM_EXPENSES_PATH, {
  moveId: mockMoveId,
  mtoShipmentId: mockMTOShipmentId,
});
const editExpensePath = generatePath(customerRoutes.SHIPMENT_PPM_EXPENSES_EDIT_PATH, {
  moveId: mockMoveId,
  mtoShipmentId: mockMTOShipmentId,
  expenseId: mockMTOShipmentWithExpenses.ppmShipment.movingExpenses[0].id,
});
const completePath = generatePath(customerRoutes.SHIPMENT_PPM_COMPLETE_PATH, {
  moveId: mockMoveId,
  mtoShipmentId: mockMTOShipmentId,
});

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
  { path: '/', element: <div>Home Page</div> },
];

const renderReviewPage = (props) => {
  return render(
    <MockProviders
      path={customerRoutes.SHIPMENT_PPM_REVIEW_PATH}
      params={{
        moveId: mockMoveId,
        mtoShipmentId: mockMTOShipmentId,
      }}
      routes={mockRoutes}
    >
      <Review {...props} />
    </MockProviders>,
  );
};

describe('Review page', () => {
  it('loads the selected shipment from redux', () => {
    renderReviewPage();
    expect(selectMTOShipmentById).toHaveBeenCalledWith(expect.anything(), mockMTOShipmentId);
  });

  it('renders the page headings', () => {
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
    renderReviewPage();

    expect(
      screen.getByText('No weight tickets uploaded. Add at least one set of weight tickets to request payment.'),
    ).toBeInTheDocument();
  });

  it('routes to the edit about your ppm page when the edit link is clicked', async () => {
    renderReviewPage();

    await userEvent.click(screen.getAllByText('Edit')[0]);

    await waitFor(() => {
      expect(screen.getByText('Edit About Your PPM Page')).toBeInTheDocument();
    });
  });

  it('routes to the add weight ticket page when the add link is clicked', async () => {
    renderReviewPage();

    await userEvent.click(screen.getByText('Add More Weight'));

    await waitFor(() => {
      expect(screen.getByText('Add More Weight Page')).toBeInTheDocument();
    });
  });

  it('routes to the edit weight ticket page when the edit link is clicked', async () => {
    selectMTOShipmentById.mockImplementation(() => mockMTOShipmentWithWeightTicket);

    renderReviewPage();

    await userEvent.click(screen.getAllByText('Edit')[1]);

    await waitFor(() => {
      expect(screen.getByText('Edit Weight Page')).toBeInTheDocument();
    });
  });

  it('routes to the add pro-gear page when the add link is clicked', async () => {
    renderReviewPage();

    await userEvent.click(screen.getByText('Add Pro-gear Weight'));

    await waitFor(() => {
      expect(screen.getByText('New Pro Gear Page')).toBeInTheDocument();
    });
  });

  it('routes to the edit pro-gear page when the edit link is clicked', async () => {
    selectMTOShipmentById.mockImplementation(() => mockMTOShipmentWithProGear);

    renderReviewPage();

    await userEvent.click(screen.getAllByText('Edit')[1]);

    await waitFor(() => {
      expect(screen.getByText('Edit Pro Gear Weight Ticket Page')).toBeInTheDocument();
    });
  });

  it('routes to the add expenses page when the add link is clicked', async () => {
    renderReviewPage();

    await userEvent.click(screen.getByText('Add Expenses'));

    await waitFor(() => {
      expect(screen.getByText('New Expense Page')).toBeInTheDocument();
    });
  });

  it('routes to the edit expense page when the edit link is clicked', async () => {
    selectMTOShipmentById.mockImplementation(() => mockMTOShipmentWithExpenses);

    renderReviewPage();

    await userEvent.click(screen.getAllByText('Edit')[1]);

    await waitFor(() => {
      expect(screen.getByText('Edit Expense Page')).toBeInTheDocument();
    });
  });

  it('routes to the home page when the return to homepage link is clicked', async () => {
    renderReviewPage();

    await userEvent.click(screen.getByText('Return To Homepage'));

    await waitFor(() => {
      expect(screen.getByText('Home Page')).toBeInTheDocument();
    });
  });

  it('routes to the complete page when the save and continue link is clicked', async () => {
    selectMTOShipmentById.mockImplementationOnce(() => mockMTOShipmentWithWeightTicket);

    renderReviewPage();

    await userEvent.click(screen.getByText('Save & Continue'));

    await waitFor(() => {
      expect(screen.getByText('Complete Page')).toBeInTheDocument();
    });
  });

  it('disables the save and continue link when there are no weight tickets', async () => {
    selectMTOShipmentById.mockImplementationOnce(() => mockMTOShipment);
    renderReviewPage();

    expect(screen.getByText('Save & Continue')).toHaveClass('usa-button--disabled');
    expect(screen.getByText('Save & Continue')).toHaveAttribute('aria-disabled', 'true');
  });

  it('disables the save and continue link when there is an incomplete weight ticket', async () => {
    selectMTOShipmentById.mockImplementationOnce(() => mockMTOShipmentWithIncompleteWeightTicket);
    renderReviewPage();

    expect(screen.getByText('Save & Continue')).toHaveClass('usa-button--disabled');
    expect(screen.getByText('Save & Continue')).toHaveAttribute('aria-disabled', 'true');
  });

  it('displays the delete confirmation modal when the delete button is clicked', async () => {
    selectMTOShipmentById.mockImplementation(() => mockMTOShipmentWithWeightTicket);
    renderReviewPage();

    await userEvent.click(screen.getAllByRole('button', { name: 'Delete' })[0]);

    await waitFor(() => {
      expect(screen.getByRole('heading', { level: 3, name: 'Delete this?' })).toBeInTheDocument();
    });

    await userEvent.click(screen.getByRole('button', { name: 'No, Keep It' }));

    expect(screen.queryByRole('heading', { level: 3, name: 'Delete this?' })).not.toBeInTheDocument();
  });

  it('calls the delete weight ticket api when confirm is clicked', async () => {
    selectMTOShipmentById.mockImplementation(() => mockMTOShipmentWithWeightTicket);
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
      expect(mockDeleteWeightTicket).toHaveBeenCalledWith(weightTicket.id, weightTicket.eTag);
    });
  });
});
