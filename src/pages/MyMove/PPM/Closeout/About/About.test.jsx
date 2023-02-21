import React from 'react';
import { render, waitFor, screen, within } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { generatePath } from 'react-router';
import { v4 } from 'uuid';
import moment from 'moment';

import About from 'pages/MyMove/PPM/Closeout/About/About';
import { selectMTOShipmentById } from 'store/entities/selectors';
import { customerRoutes, generalRoutes } from 'constants/routes';
import { getResponseError, patchMTOShipment } from 'services/internalApi';
import { updateMTOShipment } from 'store/entities/actions';
import { MockProviders } from 'testUtils';
import { SHIPMENT_OPTIONS } from 'shared/constants';
import { ppmShipmentStatuses, shipmentStatuses } from 'constants/shipments';

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
  advanceAmountReceived: 750000,
  w2Address: {
    streetAddress1: '10642 N Second Ave',
    streetAddress2: '',
    city: 'Goldsboro',
    state: 'NC',
    postalCode: '27534',
  },
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

jest.mock('services/internalApi', () => ({
  ...jest.requireActual('services/internalApi'),
  getMTOShipmentsForMove: jest.fn(() => mockMTOShipment),
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

const fillOutBasicForm = async (form) => {
  within(form).getByLabelText('When did you leave your origin?').focus();
  await userEvent.paste('31 May 2022');

  within(form).getByLabelText('Starting ZIP').focus();
  await userEvent.paste('10001', {
    initialSelectionStart: 0,
    initialSelectionEnd: 5,
  });

  within(form).getByLabelText('Ending ZIP').focus();
  await userEvent.paste('10002', {
    initialSelectionStart: 0,
    initialSelectionEnd: 5,
  });

  within(form).getByLabelText('Address 1').focus();
  await userEvent.paste('10642 N Second Ave');

  within(form).getByLabelText('City').focus();
  await userEvent.paste('Goldsboro');

  await userEvent.selectOptions(within(form).getByLabelText('State'), 'NC');

  within(form).getByLabelText('ZIP').focus();
  await userEvent.paste('27534');
};

const fillOutAdvanceSections = async (form) => {
  await userEvent.click(within(form).getByLabelText('Yes'));

  within(form).getByLabelText('How much did you receive?').focus();
  await userEvent.paste('7500');
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
    const headings = screen.getAllByRole('heading', { level: 2 });
    expect(headings[0]).toHaveTextContent('How to complete your PPM');
    expect(headings[1]).toHaveTextContent('About your final payment');

    // renders form content
    expect(headings[2]).toHaveTextContent('Departure date');
    expect(headings[3]).toHaveTextContent('Locations');
    expect(headings[4]).toHaveTextContent('Advance (AOA)');
    expect(headings[5]).toHaveTextContent('W-2 address');
  });

  it('routes back to home when return to homepage is clicked', async () => {
    render(<About />, { wrapper: MockProviders });

    await userEvent.click(screen.getByRole('button', { name: 'Return To Homepage' }));
    expect(mockPush).toHaveBeenCalledWith(homePath);
  });

  it('calls the patch shipment with the appropriate payload', async () => {
    patchMTOShipment.mockResolvedValueOnce(mockMTOShipmentResponse);

    render(<About />, { wrapper: MockProviders });
    const form = screen.getByTestId('aboutForm');

    await fillOutBasicForm(form);
    await fillOutAdvanceSections(form);

    await userEvent.click(within(form).getByRole('button', { name: 'Save & Continue' }));
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

    const form = screen.getByTestId('aboutForm');
    await fillOutBasicForm(form);

    await userEvent.click(within(form).getByRole('button', { name: 'Save & Continue' }));
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
});
