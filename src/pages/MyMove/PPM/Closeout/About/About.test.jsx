import React from 'react';
import { render, waitFor, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import moment from 'moment';
import { generatePath } from 'react-router-dom';
import { v4 } from 'uuid';

import About from 'pages/MyMove/PPM/Closeout/About/About';
import { customerRoutes, generalRoutes } from 'constants/routes';
import { ppmShipmentStatuses, shipmentStatuses } from 'constants/shipments';
import { getResponseError, patchMTOShipment } from 'services/internalApi';
import { SHIPMENT_OPTIONS } from 'shared/constants';
import { updateMTOShipment } from 'store/entities/actions';
import { selectMTOShipmentById } from 'store/entities/selectors';
import { MockProviders, setUpProvidersWithHistory } from 'testUtils';
import { createBaseWeightTicket, createCompleteWeightTicket } from 'utils/test/factories/weightTicket';

const mockMoveId = v4();
const mockMTOShipmentId = v4();
const mockPPMShipmentId = v4();

jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom'),
  useParams: () => ({
    moveId: mockMoveId,
    mtoShipmentId: mockMTOShipmentId,
  }),
}));

jest.mock('services/internalApi', () => ({
  ...jest.requireActual('services/internalApi'),
  patchMTOShipment: jest.fn(),
  getResponseError: jest.fn(),
}));

const mtoShipmentCreatedDate = new Date();
const ppmShipmentCreatedDate = moment(mtoShipmentCreatedDate).add(5, 'seconds');
const approvedDate = moment(ppmShipmentCreatedDate).add(2, 'days');

const mockMTOShipment = {
  id: mockMTOShipmentId,
  shipmentType: SHIPMENT_OPTIONS.PPM,
  status: shipmentStatuses.APPROVED,
  moveTaskOrderId: mockMoveId,
  ppmShipment: {
    id: mockPPMShipmentId,
    shipmentId: mockMTOShipmentId,
    status: ppmShipmentStatuses.WAITING_ON_CUSTOMER,
    pickupPostalCode: '10001',
    destinationPostalCode: '10002',
    expectedDepartureDate: '2022-04-30',
    hasRequestedAdvance: true,
    advanceAmountRequested: 598700,
    estimatedWeight: 4000,
    estimatedIncentive: 1000000,
    sitExpected: false,
    hasProGear: false,
    proGearWeight: null,
    spouseProGearWeight: null,
    actualMoveDate: null,
    actualPickupPostalCode: null,
    actualDestinationPostalCode: null,
    hasReceivedAdvance: null,
    advanceAmountReceived: null,
    weightTickets: [],
    createdAt: ppmShipmentCreatedDate.toISOString(),
    updatedAt: approvedDate.toISOString(),
    eTag: window.btoa(approvedDate.toISOString()),
  },
  createdAt: mtoShipmentCreatedDate.toISOString(),
  updatedAt: approvedDate.toISOString(),
  eTag: window.btoa(approvedDate.toISOString()),
};

const partialPayload = {
  actualMoveDate: '2022-05-31',
  actualPickupPostalCode: '10001',
  actualDestinationPostalCode: '10002',
  hasReceivedAdvance: true,
  advanceAmountReceived: 598700,
};

const mockPayload = {
  shipmentType: SHIPMENT_OPTIONS.PPM,
  ppmShipment: {
    id: mockPPMShipmentId,
    ...partialPayload,
  },
};

const mockMTOShipmentResponse = {
  ...mockMTOShipment,
  ppmShipment: {
    ...mockMTOShipment.ppmShipment,
    ...partialPayload,
  },
};

jest.mock('store/entities/selectors', () => ({
  ...jest.requireActual('store/entities/selectors'),
  selectMTOShipmentById: jest.fn(() => mockMTOShipment),
}));

jest.mock('utils/validation', () => ({
  ...jest.requireActual('utils/validation'),
  validatePostalCode: jest.fn(),
}));

const mockDispatch = jest.fn();
jest.mock('react-redux', () => ({
  ...jest.requireActual('react-redux'),
  useDispatch: jest.fn(() => mockDispatch),
}));

beforeEach(() => {
  jest.clearAllMocks();
});

const homePath = generatePath(generalRoutes.HOME_PATH);
const weightTicketsPath = generatePath(customerRoutes.SHIPMENT_PPM_WEIGHT_TICKETS_PATH, {
  moveId: mockMoveId,
  mtoShipmentId: mockMTOShipmentId,
});

const fillOutBasicForm = () => {
  const actualMoveDate = screen.getByLabelText('When did you leave your origin?');
  userEvent.clear(actualMoveDate);
  userEvent.type(actualMoveDate, '31 May 2022');

  const actualStartingZip = screen.getByLabelText('Starting ZIP');
  userEvent.clear(actualStartingZip);
  userEvent.type(actualStartingZip, '10001');

  const actualDestinationZip = screen.getByLabelText('Ending ZIP');
  userEvent.clear(actualDestinationZip);
  userEvent.type(actualDestinationZip, '10002');
};

const fillOutAdvanceSections = () => {
  const hasReceivedAdvance = screen.getByLabelText('Yes');
  userEvent.click(hasReceivedAdvance);

  const advanceAmountReceived = screen.getByLabelText('How much did you receive?');
  userEvent.clear(advanceAmountReceived);
  userEvent.type(advanceAmountReceived, '5987');
};

describe('About page', () => {
  it('loads the selected shipment from redux', () => {
    render(<About />, { wrapper: MockProviders });

    expect(selectMTOShipmentById).toHaveBeenCalledWith(expect.anything(), mockMTOShipmentId);
  });

  it('renders the page Content', () => {
    render(<About />, { wrapper: MockProviders });

    expect(screen.getByTestId('tag')).toHaveTextContent('PPM');
    expect(screen.getByRole('heading', { level: 1 })).toHaveTextContent('About your PPM');
    expect(screen.getByText('Finish moving this PPM before you start documenting it.')).toBeInTheDocument();

    // renders form content
    expect(screen.getAllByRole('heading', { level: 2 })[0]).toHaveTextContent('Departure date');
    expect(screen.getAllByRole('heading', { level: 2 })[1]).toHaveTextContent('Locations');
    expect(screen.getAllByRole('heading', { level: 2 })[2]).toHaveTextContent('Advance (AOA)');
  });

  it('routes back to home when finish later is clicked', () => {
    const { memoryHistory, mockProviderWithHistory } = setUpProvidersWithHistory();

    render(<About />, { wrapper: mockProviderWithHistory });

    userEvent.click(screen.getByRole('button', { name: 'Finish Later' }));

    expect(memoryHistory.location.pathname).toBe(homePath);
  });

  it('calls the patch shipment with the appropriate payload', async () => {
    patchMTOShipment.mockResolvedValueOnce(mockMTOShipmentResponse);

    const { memoryHistory, mockProviderWithHistory } = setUpProvidersWithHistory();

    render(<About />, { wrapper: mockProviderWithHistory });

    fillOutBasicForm();
    fillOutAdvanceSections();

    userEvent.click(screen.getByRole('button', { name: 'Save & Continue' }));
    await waitFor(() => {
      expect(patchMTOShipment).toHaveBeenCalledWith(mockMTOShipmentId, mockPayload, mockMTOShipment.eTag);
    });

    expect(mockDispatch).toHaveBeenCalledWith(updateMTOShipment(mockMTOShipmentResponse));

    expect(memoryHistory.location.pathname).toBe(weightTicketsPath);
  });

  it('displays an error when the patch shipment API fails', async () => {
    const mockErrorMsg = 'Error Updating';
    patchMTOShipment.mockRejectedValue(mockErrorMsg);
    getResponseError.mockReturnValue(mockErrorMsg);

    render(<About />, { wrapper: MockProviders });

    fillOutBasicForm();

    userEvent.click(screen.getByRole('button', { name: 'Save & Continue' }));
    const payload = {
      ...mockPayload,
      ppmShipment: {
        ...mockPayload.ppmShipment,
        hasReceivedAdvance: false,
        advanceAmountReceived: null,
      },
    };
    await waitFor(() => {
      expect(patchMTOShipment).toHaveBeenCalledWith(mockMTOShipmentId, payload, mockMTOShipment.eTag);
    });

    expect(screen.getByText(mockErrorMsg)).toBeInTheDocument();
  });

  it('expect loadingPlaceholder when mtoShipment is falsy', () => {
    selectMTOShipmentById.mockReturnValue(null);

    render(<About />, { wrapper: MockProviders });
    expect(screen.getByRole('heading', { level: 2 })).toHaveTextContent('Loading, please wait...');
  });

  describe('routes to the', () => {
    const ppmShipmentWithActualShipmentInfo = {
      ...mockMTOShipment,
      ppmShipment: {
        ...mockMTOShipment.ppmShipment,
        actualMoveDate: mockMTOShipment.ppmShipment.expectedDepartureDate,
        actualPickupPostalCode: mockMTOShipment.ppmShipment.pickupPostalCode,
        actualDestinationPostalCode: mockMTOShipment.ppmShipment.destinationPostalCode,
        hasReceivedAdvance: mockMTOShipment.ppmShipment.hasRequestedAdvance,
        advanceAmountReceived: mockMTOShipment.ppmShipment.advanceAmountRequested,
      },
    };

    const serviceMemberId = v4();

    const ppmShipmentWithIncompleteWeightTicket = {
      ...ppmShipmentWithActualShipmentInfo,
      ppmShipment: {
        ...ppmShipmentWithActualShipmentInfo.ppmShipment,
        weightTickets: [
          createBaseWeightTicket({ serviceMemberId }, { ppmShipmentId: ppmShipmentWithActualShipmentInfo.id }),
        ],
      },
    };

    const ppmShipmentWithCompleteWeightTicket = {
      ...ppmShipmentWithIncompleteWeightTicket,
      ppmShipment: {
        ...ppmShipmentWithIncompleteWeightTicket.ppmShipment,
        weightTickets: [
          createCompleteWeightTicket({ serviceMemberId }, { ppmShipmentId: ppmShipmentWithActualShipmentInfo.id }),
        ],
      },
    };

    const ppmShipmentWithMultipleIncompleteWeightTickets = {
      ...ppmShipmentWithIncompleteWeightTicket,
      ppmShipment: {
        ...ppmShipmentWithIncompleteWeightTicket.ppmShipment,
        weightTickets: [
          ...ppmShipmentWithIncompleteWeightTicket.ppmShipment.weightTickets,
          createBaseWeightTicket({ serviceMemberId }, { ppmShipmentId: ppmShipmentWithIncompleteWeightTicket.id }),
        ],
      },
    };

    const ppmShipmentWithMultipleWeightTickets = {
      ...ppmShipmentWithCompleteWeightTicket,
      ppmShipment: {
        ...ppmShipmentWithCompleteWeightTicket.ppmShipment,
        weightTickets: [
          ...ppmShipmentWithCompleteWeightTicket.ppmShipment.weightTickets,
          createBaseWeightTicket({ serviceMemberId }, { ppmShipmentId: ppmShipmentWithActualShipmentInfo.id }),
        ],
      },
    };

    it.each([
      [
        'new Weight Ticket page if weight ticket info is missing',
        ppmShipmentWithActualShipmentInfo,
        generatePath(customerRoutes.SHIPMENT_PPM_WEIGHT_TICKETS_PATH, {
          moveId: mockMoveId,
          mtoShipmentId: ppmShipmentWithActualShipmentInfo.id,
        }),
      ],
      [
        'edit Weight Ticket page if weight ticket info is incomplete',
        ppmShipmentWithIncompleteWeightTicket,
        generatePath(customerRoutes.SHIPMENT_PPM_WEIGHT_TICKETS_EDIT_PATH, {
          moveId: mockMoveId,
          mtoShipmentId: ppmShipmentWithIncompleteWeightTicket.id,
          weightTicketId: ppmShipmentWithIncompleteWeightTicket.ppmShipment.weightTickets[0].id,
        }),
      ],
      [
        'edit Weight Ticket page for the first weight ticket if there are multiple but none are complete',
        ppmShipmentWithMultipleIncompleteWeightTickets,
        generatePath(customerRoutes.SHIPMENT_PPM_WEIGHT_TICKETS_EDIT_PATH, {
          moveId: mockMoveId,
          mtoShipmentId: ppmShipmentWithMultipleIncompleteWeightTickets.id,
          weightTicketId: ppmShipmentWithMultipleIncompleteWeightTickets.ppmShipment.weightTickets[0].id,
        }),
      ],
      [
        'Review page if weight ticket info is complete',
        ppmShipmentWithCompleteWeightTicket,
        generatePath(customerRoutes.SHIPMENT_PPM_REVIEW_PATH, {
          moveId: mockMoveId,
          mtoShipmentId: ppmShipmentWithCompleteWeightTicket.id,
        }),
      ],
      [
        'Review page if at least one weight ticket is completely filled out',
        ppmShipmentWithMultipleWeightTickets,
        generatePath(customerRoutes.SHIPMENT_PPM_REVIEW_PATH, {
          moveId: mockMoveId,
          mtoShipmentId: ppmShipmentWithMultipleWeightTickets.id,
        }),
      ],
    ])('%s', async (scenarioDescription, shipment, expectedRoute) => {
      selectMTOShipmentById.mockReturnValue(shipment);
      patchMTOShipment.mockResolvedValueOnce(shipment);

      const { memoryHistory, mockProviderWithHistory } = setUpProvidersWithHistory();

      render(<About />, { wrapper: mockProviderWithHistory });

      userEvent.click(screen.getByRole('button', { name: 'Save & Continue' }));

      const expectedPayload = {
        shipmentType: SHIPMENT_OPTIONS.PPM,
        ppmShipment: {
          id: mockPPMShipmentId,
          actualMoveDate: shipment.ppmShipment.actualMoveDate,
          actualPickupPostalCode: shipment.ppmShipment.actualPickupPostalCode,
          actualDestinationPostalCode: shipment.ppmShipment.actualDestinationPostalCode,
          hasReceivedAdvance: shipment.ppmShipment.hasReceivedAdvance,
          advanceAmountReceived: shipment.ppmShipment.advanceAmountReceived,
        },
      };

      await waitFor(() => {
        expect(patchMTOShipment).toHaveBeenCalledWith(mockMTOShipmentId, expectedPayload, shipment.eTag);
      });

      expect(mockDispatch).toHaveBeenCalledWith(updateMTOShipment(shipment));

      expect(memoryHistory.location.pathname).toEqual(expectedRoute);
    });
  });
});
