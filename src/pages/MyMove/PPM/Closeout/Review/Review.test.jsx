import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import { v4 } from 'uuid';
import userEvent from '@testing-library/user-event';
import { generatePath } from 'react-router-dom';

import { MockProviders } from 'testUtils';
import { isBooleanFlagEnabled } from 'utils/featureFlags';
import { selectMTOShipmentById } from 'store/entities/selectors';
import Review from 'pages/MyMove/PPM/Closeout/Review/Review';
import { PPM_TYPES, SHIPMENT_OPTIONS } from 'shared/constants';
import { customerRoutes } from 'constants/routes';
import {
  deleteWeightTicket,
  deleteProGearWeightTicket,
  deleteMovingExpense,
  deleteGunSafeWeightTicket,
  getMTOShipmentsForMove,
} from 'services/internalApi';
import { createBaseWeightTicket, createCompleteWeightTicket } from 'utils/test/factories/weightTicket';
import { createBaseProGearWeightTicket } from 'utils/test/factories/proGearWeightTicket';
import { createBaseGunSafeWeightTicket } from 'utils/test/factories/gunSafeWeightTicket';
import { createCompleteMovingExpense, createCompleteSITMovingExpense } from 'utils/test/factories/movingExpense';

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
    hasGunSafe: false,
    gunSafeWeight: null,
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
    hasGunSafe: false,
    gunSafeWeight: null,
    weightTickets: [weightTicketOne, weightTicketTwo],
    pickupAddress,
    destinationAddress,
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
        hasGunSafe: false,
        gunSafeWeight: null,
        weightTickets: [weightTicketTwo],
        pickupAddress,
        destinationAddress,
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
    hasGunSafe: false,
    gunSafeWeight: null,
    weightTickets: [createBaseWeightTicket()],
    pickupAddress,
    destinationAddress,
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
    hasGunSafe: false,
    gunSafeWeight: null,
    pickupAddress,
    destinationAddress,
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
        proGearWeightTickets: [],
        pickupAddress,
        destinationAddress,
      },
      eTag: 'dGVzdGluZzIzNDQzMjQ',
    },
  },
};

const gunSafeWeightOne = createBaseGunSafeWeightTicket();
const mockMTOShipmentWithGunSafe = {
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
    gunSafeWeightTickets: [gunSafeWeightOne],
    hasGunSafe: true,
    gunSafeWeight: null,
    pickupAddress,
    destinationAddress,
  },
  eTag: 'dGVzdGluZzIzNDQzMjQ',
};

const mockMTOShipmentWithGunSafeDeleted = {
  mtoShipments: {
    [mockMTOShipmentId]: {
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
        proGearWeight: 0,
        spouseProGearWeight: null,
        proGearWeightTickets: [],
        hasGunSafe: true,
        gunSafeWeight: 200,
        gunSafeWeightTickets: [],
        pickupAddress,
        destinationAddress,
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
    hasGunSafe: false,
    gunSafeWeight: 0,
    movingExpenses: [expenseOne, expenseTwo],
    pickupAddress,
    destinationAddress,
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
        hasGunSafe: false,
        gunSafeWeight: 0,
        movingExpenses: [expenseOne],
        pickupAddress,
        destinationAddress,
      },
      eTag: 'dGVzdGluZzIzNDQzMjQ',
    },
  },
};

const mockServiceMember = {
  id: 'testId',
};

const mockNavigate = jest.fn();
jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom'),
  useNavigate: () => mockNavigate,
}));

jest.mock('services/internalApi', () => ({
  ...jest.requireActual('services/internalApi'),
  deleteWeightTicket: jest.fn(() => {}),
  deleteProGearWeightTicket: jest.fn(() => {}),
  deleteGunSafeWeightTicket: jest.fn(() => {}),
  deleteMovingExpense: jest.fn(() => {}),
  getMTOShipmentsForMove: jest.fn(),
  getAllMoves: jest.fn().mockImplementation(() => Promise.resolve()),
}));

jest.mock('store/entities/selectors', () => ({
  ...jest.requireActual('store/entities/selectors'),
  selectMTOShipmentById: jest.fn(() => mockMTOShipment),
  selectServiceMemberFromLoggedInUser: jest.fn(() => mockServiceMember),
}));

jest.mock('utils/featureFlags', () => ({
  ...jest.requireActual('utils/featureFlags'),
  isBooleanFlagEnabled: jest.fn().mockImplementation(() => Promise.resolve(false)),
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
const newGunSafePath = generatePath(customerRoutes.SHIPMENT_PPM_GUN_SAFE_PATH, {
  moveId: mockMoveId,
  mtoShipmentId: mockMTOShipmentId,
});
const editGunSafeWeightTicket = generatePath(customerRoutes.SHIPMENT_PPM_GUN_SAFE_EDIT_PATH, {
  moveId: mockMoveId,
  mtoShipmentId: mockMTOShipmentId,
  gunSafeId: mockMTOShipmentWithGunSafe.ppmShipment.gunSafeWeightTickets[0].id,
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
    path: newGunSafePath,
    element: <div>New Gun Safe Page</div>,
  },
  {
    path: editGunSafeWeightTicket,
    element: <div>Edit Gun Safe Weight Ticket Page</div>,
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

const renderReviewPage = async (props) => {
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

  it('renders the page headings', async () => {
    isBooleanFlagEnabled.mockImplementation(() => Promise.resolve(true));
    await renderReviewPage();

    expect(screen.getByTestId('tag')).toHaveTextContent('PPM');
    expect(screen.getByRole('heading', { level: 1 })).toHaveTextContent('Review');
    expect(screen.getAllByRole('heading', { level: 2 })[0]).toHaveTextContent('About Your PPM');
    expect(screen.getAllByRole('heading', { level: 2 })[1]).toHaveTextContent('Documents');
    expect(screen.getAllByRole('heading', { level: 3 })[0]).toHaveTextContent('Weight moved');
    expect(screen.getAllByRole('heading', { level: 3 })[1]).toHaveTextContent('Pro-gear');
    expect(screen.getAllByRole('heading', { level: 3 })[2]).toHaveTextContent('Gun safe');
    expect(screen.getAllByRole('heading', { level: 3 })[3]).toHaveTextContent('Expenses');
  });

  it('excludes gun safe from the page headings when FF is toggled off', async () => {
    isBooleanFlagEnabled.mockImplementation(() => Promise.resolve(false));
    await renderReviewPage();

    expect(screen.getByTestId('tag')).toHaveTextContent('PPM');
    expect(screen.getByRole('heading', { level: 1 })).toHaveTextContent('Review');
    expect(screen.getAllByRole('heading', { level: 2 })[0]).toHaveTextContent('About Your PPM');
    expect(screen.getAllByRole('heading', { level: 2 })[1]).toHaveTextContent('Documents');
    expect(screen.getAllByRole('heading', { level: 3 })[0]).toHaveTextContent('Weight moved');
    expect(screen.getAllByRole('heading', { level: 3 })[1]).toHaveTextContent('Pro-gear');
    expect(screen.getAllByRole('heading', { level: 3 })[2]).toHaveTextContent('Expenses');
  });

  it('renders the empty message when there are no weight tickets', async () => {
    await renderReviewPage();

    expect(
      screen.getByText('No weight moved documented. At least one trip is required to continue.'),
    ).toBeInTheDocument();
  });

  it('routes to the edit about your ppm page when the edit link is clicked', async () => {
    await renderReviewPage();

    await userEvent.click(screen.getAllByText('Edit')[0]);

    await waitFor(() => {
      expect(screen.getByText('Edit About Your PPM Page')).toBeInTheDocument();
    });
  });

  it('routes to the add weight ticket page when the add link is clicked', async () => {
    await renderReviewPage();

    await userEvent.click(screen.getByText('Add More Weight'));

    await waitFor(() => {
      expect(screen.getByText('Add More Weight Page')).toBeInTheDocument();
    });
  });

  it('routes to the edit weight ticket page when the edit link is clicked', async () => {
    selectMTOShipmentById.mockImplementation(() => mockMTOShipmentWithWeightTicket);

    await renderReviewPage();

    await userEvent.click(screen.getAllByText('Edit')[1]);

    await waitFor(() => {
      expect(screen.getByText('Edit Weight Page')).toBeInTheDocument();
    });
  });

  it('routes to the add pro-gear page when the add link is clicked', async () => {
    await renderReviewPage();

    await userEvent.click(screen.getByText('Add Pro-gear Weight'));

    await waitFor(() => {
      expect(screen.getByText('New Pro Gear Page')).toBeInTheDocument();
    });
  });

  it('routes to the edit pro-gear page when the edit link is clicked', async () => {
    selectMTOShipmentById.mockImplementation(() => mockMTOShipmentWithProGear);

    await renderReviewPage();

    await userEvent.click(screen.getAllByText('Edit')[1]);

    await waitFor(() => {
      expect(screen.getByText('Edit Pro Gear Weight Ticket Page')).toBeInTheDocument();
    });
  });

  it('routes to the add gun safe page when the add link is clicked', async () => {
    isBooleanFlagEnabled.mockImplementation(() => Promise.resolve(true));
    await renderReviewPage();

    await userEvent.click(screen.getByText('Add Gun Safe Weight'));

    await waitFor(() => {
      expect(screen.getByText('New Gun Safe Page')).toBeInTheDocument();
    });
  });

  it('excludes the links to add a gun safe when the FF is toggled off', async () => {
    isBooleanFlagEnabled.mockImplementation(() => Promise.resolve(false));
    await renderReviewPage();
    expect(screen.queryByText('Add Gun Safe Weight', { exact: false })).toBeNull();
  });

  it('routes to the edit gun safe page when the edit link is clicked', async () => {
    isBooleanFlagEnabled.mockImplementation(() => Promise.resolve(true));
    selectMTOShipmentById.mockImplementation(() => mockMTOShipmentWithGunSafe);

    await renderReviewPage();

    await userEvent.click(screen.getAllByText('Edit')[1]);

    await waitFor(() => {
      expect(screen.getByText('Edit Gun Safe Weight Ticket Page')).toBeInTheDocument();
    });
  });

  it('excludes gun safe edit links when the FF is toggled off', async () => {
    isBooleanFlagEnabled.mockImplementation(() => Promise.resolve(false));
    selectMTOShipmentById.mockImplementation(() => mockMTOShipmentWithGunSafe);
    await renderReviewPage();
    expect(screen.queryByText('Edit Gun Safe Weight Ticket Page', { exact: false })).toBeNull();
  });

  it('routes to the add expenses page when the add link is clicked', async () => {
    await renderReviewPage();

    await userEvent.click(screen.getByText('Add Expenses'));

    await waitFor(() => {
      expect(screen.getByText('New Expense Page')).toBeInTheDocument();
    });
  });

  it('routes to the edit expense page when the edit link is clicked', async () => {
    selectMTOShipmentById.mockImplementation(() => mockMTOShipmentWithExpenses);

    await renderReviewPage();

    await userEvent.click(screen.getAllByText('Edit')[1]);

    await waitFor(() => {
      expect(screen.getByText('Edit Expense Page')).toBeInTheDocument();
    });
  });

  it('routes to the home page when the return to homepage link is clicked', async () => {
    await renderReviewPage();

    // await userEvent.click(screen.getByText('Return To Homepage'));
    await userEvent.click(screen.getByTestId('reviewReturnToHomepageLink'));

    // expect(mockNavigate).toHaveBeenCalledWith(generalRoutes.HOME_PATH);

    await waitFor(() => {
      expect(screen.getByText('Home Page')).toBeInTheDocument();
    });
  });

  it('routes to the complete page when the save and continue link is clicked', async () => {
    selectMTOShipmentById.mockImplementationOnce(() => mockMTOShipmentWithWeightTicket);

    await renderReviewPage();

    await userEvent.click(screen.getByText('Save & Continue'));

    await waitFor(() => {
      expect(screen.getByText('Complete Page')).toBeInTheDocument();
    });
  });

  it('disables the save and continue link when there are no weight tickets', async () => {
    selectMTOShipmentById.mockImplementationOnce(() => mockMTOShipment);
    await renderReviewPage();

    expect(screen.getByText('Save & Continue')).toHaveClass('usa-button--disabled');
    expect(screen.getByText('Save & Continue')).toHaveAttribute('aria-disabled', 'true');
  });

  it('disables the save and continue link when there is an incomplete weight ticket', async () => {
    selectMTOShipmentById.mockImplementationOnce(() => mockMTOShipmentWithIncompleteWeightTicket);
    await renderReviewPage();

    expect(screen.getByText('Save & Continue')).toHaveClass('usa-button--disabled');
    expect(screen.getByText('Save & Continue')).toHaveAttribute('aria-disabled', 'true');
  });

  it('error message is displayed when a PPM shipment is in an incomplete state', async () => {
    selectMTOShipmentById.mockImplementationOnce(() => mockMTOShipmentWithIncompleteWeightTicket);
    await renderReviewPage();

    expect(
      screen.getByText(
        'There are items below that are missing required information. Please select “Edit” to enter all required information or “Delete” to remove the item.',
      ),
    ).toBeInTheDocument();

    expect(screen.getByText('This trip is missing required information.')).toBeInTheDocument();
  });

  it('displays the delete confirmation modal when the delete button for Weight Moved/Trip 2 is clicked', async () => {
    selectMTOShipmentById.mockImplementation(() => mockMTOShipmentWithWeightTicket);
    await renderReviewPage();

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

    await renderReviewPage();

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
    await renderReviewPage();

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

  it('calls the delete gun safe weight ticket api when confirm is clicked', async () => {
    isBooleanFlagEnabled.mockImplementation(() => Promise.resolve(true));
    selectMTOShipmentById.mockImplementation(() => mockMTOShipmentWithGunSafe);
    const mockDeleteGunSafeWeightTicket = jest.fn().mockResolvedValue({});
    deleteGunSafeWeightTicket.mockImplementationOnce(mockDeleteGunSafeWeightTicket);
    getMTOShipmentsForMove.mockResolvedValue(mockMTOShipmentWithGunSafeDeleted);
    await renderReviewPage();

    await userEvent.click(screen.getAllByRole('button', { name: 'Delete' })[0]);

    await waitFor(() => {
      expect(screen.getByRole('heading', { level: 3, name: 'Delete this?' })).toBeInTheDocument();
    });

    await userEvent.click(screen.getByRole('button', { name: 'Yes, Delete' }));

    const gunSafeWeightTicket = mockMTOShipmentWithGunSafe.ppmShipment.gunSafeWeightTickets[0];
    await waitFor(() => {
      expect(mockDeleteGunSafeWeightTicket).toHaveBeenCalledWith(
        mockMTOShipmentWithWeightTicket.ppmShipment.id,
        gunSafeWeightTicket.id,
      );
    });

    await waitFor(() => {
      expect(screen.getByText('Set 1 successfully deleted.'));
    });
  });

  it('excludes links to delete gun safe tickets when the FF is toggled off', async () => {
    isBooleanFlagEnabled.mockImplementation(() => Promise.resolve(false));
    selectMTOShipmentById.mockImplementation(() => mockMTOShipmentWithGunSafe);
    const mockDeleteGunSafeWeightTicket = jest.fn().mockResolvedValue({});
    deleteGunSafeWeightTicket.mockImplementationOnce(mockDeleteGunSafeWeightTicket);
    getMTOShipmentsForMove.mockResolvedValue(mockMTOShipmentWithGunSafeDeleted);
    await renderReviewPage();

    expect(screen.queryAllByRole('button', { name: 'Delete' })).toHaveLength(0);
  });

  it('calls the delete moving expense api when confirm is clicked', async () => {
    selectMTOShipmentById.mockImplementation(() => mockMTOShipmentWithExpenses);
    const mockDeleteMovingExpense = jest.fn().mockResolvedValue({});
    deleteMovingExpense.mockImplementationOnce(mockDeleteMovingExpense);
    getMTOShipmentsForMove.mockResolvedValue(mockMTOShipmentWithExpensesDeleted);
    await renderReviewPage();

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

  it('disables the save and continue link for small package shipments with no expenses', async () => {
    const mockMTOShipmentWithSmallPackageNoExpense = {
      id: mockMTOShipmentId,
      shipmentType: SHIPMENT_OPTIONS.PPM,
      ppmShipment: {
        id: mockPPMShipmentId,
        ppmType: PPM_TYPES.SMALL_PACKAGE,
        weightTickets: [],
        movingExpenses: [],
        proGearWeightTickets: [],
        gunSafeWeightTickets: [],
        pickupAddress,
        destinationAddress,
      },
      eTag: 'dummyETag',
    };

    selectMTOShipmentById.mockImplementationOnce(() => mockMTOShipmentWithSmallPackageNoExpense);
    await renderReviewPage();

    const saveButton = screen.getByText('Save & Continue');
    expect(saveButton).toHaveClass('usa-button--disabled');
    expect(saveButton).toHaveAttribute('aria-disabled', 'true');
  });

  it('enables the save and continue link for small package shipments when at least one expense exists', async () => {
    const mockMTOShipmentWithSmallPackageExpense = {
      id: mockMTOShipmentId,
      shipmentType: SHIPMENT_OPTIONS.PPM,
      ppmShipment: {
        id: mockPPMShipmentId,
        ppmType: PPM_TYPES.SMALL_PACKAGE,
        weightTickets: [],
        movingExpenses: [expenseOne],
        proGearWeightTickets: [],
        gunSafeWeightTickets: [],
        pickupAddress,
        destinationAddress,
      },
      eTag: 'dummyETag',
    };

    selectMTOShipmentById.mockImplementationOnce(() => mockMTOShipmentWithSmallPackageExpense);
    await renderReviewPage();

    const saveButton = screen.getByText('Save & Continue');
    expect(saveButton).not.toHaveClass('usa-button--disabled');
    expect(saveButton).toHaveAttribute('aria-disabled', 'false');
  });
});
