import React from 'react';
import { render, waitFor, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { generatePath } from 'react-router';
import { v4 } from 'uuid';

import About from 'pages/MyMove/PPM/Closeout/About/About';
import { selectMTOShipmentById } from 'store/entities/selectors';
import { customerRoutes, generalRoutes } from 'constants/routes';
import { getResponseError, patchMTOShipment } from 'services/internalApi';
import { updateMTOShipment } from 'store/entities/actions';
import { MockProviders } from 'testUtils';

const mockMoveId = v4();
const mockMTOShipmentId = v4();
const mockPPMShipmentId = v4();

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
  patchMTOShipment: jest.fn(),
  getResponseError: jest.fn(),
}));

const mockMTOShipment = {
  id: mockMTOShipmentId,
  shipmentType: 'PPM',
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

const partialPayload = {
  actualMoveDate: '2022-05-31',
  actualPickupPostalCode: '10001',
  actualDestinationPostalCode: '10002',
  hasReceivedAdvance: true,
  advanceAmountReceived: 750000,
};

const mockPayload = {
  shipmentType: 'PPM',
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
  userEvent.type(advanceAmountReceived, '7500');
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
    render(<About />, { wrapper: MockProviders });

    userEvent.click(screen.getByRole('button', { name: 'Finish Later' }));
    expect(mockPush).toHaveBeenCalledWith(homePath);
  });

  it('calls the patch shipment with the appropriate payload', async () => {
    patchMTOShipment.mockResolvedValueOnce(mockMTOShipmentResponse);

    render(<About />, { wrapper: MockProviders });

    fillOutBasicForm();
    fillOutAdvanceSections();

    userEvent.click(screen.getByRole('button', { name: 'Save & Continue' }));
    await waitFor(() => {
      expect(patchMTOShipment).toHaveBeenCalledWith(mockMTOShipmentId, mockPayload, mockMTOShipment.eTag);
    });

    expect(mockDispatch).toHaveBeenCalledWith(updateMTOShipment(mockMTOShipmentResponse));
    expect(mockPush).toHaveBeenCalledWith(weightTicketsPath);
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
});
