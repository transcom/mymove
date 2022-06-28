import React from 'react';
import { render, waitFor, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { generatePath } from 'react-router';
import { v4 } from 'uuid';

import { selectMTOShipmentById } from 'store/entities/selectors';
import { customerRoutes, generalRoutes } from 'constants/routes';
import { createWeightTicket, patchWeightTicket } from 'services/internalApi';
import { MockProviders } from 'testUtils';
import { SHIPMENT_OPTIONS } from 'shared/constants';
import WeightTickets from 'pages/MyMove/PPM/Closeout/WeightTickets/WeightTickets';

const mockMoveId = v4();
const mockMTOShipmentId = v4();
const mockPPMShipmentId = v4();
const mockWeightTicketId = v4();
const mockWeightTicketETag = window.btoa(new Date());

const mockPush = jest.fn();
const mockReplace = jest.fn();
jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom'),
  useHistory: () => ({
    push: mockPush,
    replace: mockReplace,
  }),
  useParams: () => ({
    moveId: mockMoveId,
    mtoShipmentId: mockMTOShipmentId,
  }),
  useLocation: () => ({
    search: 'tripNumber=2',
  }),
}));

jest.mock('services/internalApi', () => ({
  ...jest.requireActual('services/internalApi'),
  createWeightTicket: jest.fn(),
  createUploadForDocument: jest.fn(),
  patchWeightTicket: jest.fn(),
  getResponseError: jest.fn(),
}));

const mockMTOShipment = {
  id: mockMTOShipmentId,
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

const mockEmptyWeightDocumentId = v4();
const mockFullWeightDocumentId = v4();
const mockTrailerOwnershipWeightDocumentId = v4();

const mockWeightTicket = {
  id: mockWeightTicketId,
  ppmShipmentId: mockPPMShipmentId,
  emptyWeightDocumentId: mockEmptyWeightDocumentId,
  fullWeightDocumentId: mockFullWeightDocumentId,
  trailerOwnershipDocumentId: mockTrailerOwnershipWeightDocumentId,
  eTag: mockWeightTicketETag,
};

const mockWeightTicketWithUploads = {
  id: mockWeightTicketId,
  ppmShipmentId: mockPPMShipmentId,
  emptyWeightDocumentId: mockEmptyWeightDocumentId,
  emptyWeightTickets: [
    {
      id: '299e2fb4-432d-4261-bbed-d8280c6090af',
      created_at: '2022-06-22T23:25:50.490Z',
      bytes: 819200,
      url: 'a/fake/path',
      filename: 'empty_weight.jpg',
      content_type: 'image/jpg',
    },
  ],
  fullWeightDocumentId: mockFullWeightDocumentId,
  fullWeightTickets: [
    {
      id: 'f70af8a1-38e9-4ae2-a837-3c0c61069a0d',
      created_at: '2022-06-23T23:25:50.490Z',
      bytes: 409600,
      url: 'a/fake/path',
      filename: 'full_weight.pdf',
      content_type: 'application/pdf',
    },
  ],
  trailerOwnershipDocumentId: mockTrailerOwnershipWeightDocumentId,
  trailerOwnershipDocs: [
    {
      id: 'fd4e80f8-d025-44b2-8c33-15240fac51ab',
      created_at: '2022-06-24T23:25:50.490Z',
      bytes: 204800,
      url: 'a/fake/path',
      filename: 'trailer_title.pdf',
      content_type: 'application/pdf',
    },
  ],
  eTag: mockWeightTicketETag,
};

jest.mock('store/entities/selectors', () => ({
  ...jest.requireActual('store/entities/selectors'),
  selectMTOShipmentById: jest.fn(() => mockMTOShipment),
}));

beforeEach(() => {
  jest.clearAllMocks();
});

const homePath = generatePath(generalRoutes.HOME_PATH);
const weightTicketsEditPath = generatePath(customerRoutes.SHIPMENT_PPM_WEIGHT_TICKETS_EDIT_PATH, {
  moveId: mockMoveId,
  mtoShipmentId: mockMTOShipmentId,
  weightTicketId: mockWeightTicketId,
});
const reviewPath = generatePath(customerRoutes.SHIPMENT_PPM_REVIEW_PATH, {
  moveId: mockMoveId,
  mtoShipmentId: mockMTOShipmentId,
});

describe('Weight Tickets page', () => {
  it('loads the selected shipment from redux', async () => {
    createWeightTicket.mockResolvedValue(mockWeightTicket);

    render(<WeightTickets />, { wrapper: MockProviders });

    await waitFor(() => {
      expect(selectMTOShipmentById).toHaveBeenCalledWith(expect.anything(), mockMTOShipmentId);
    });
  });

  it('renders the page Content', async () => {
    createWeightTicket.mockResolvedValue(mockWeightTicket);

    render(<WeightTickets />, { wrapper: MockProviders });

    await waitFor(() => {
      expect(screen.getByTestId('tag')).toHaveTextContent('PPM');
    });

    expect(screen.getByRole('heading', { level: 1 })).toHaveTextContent('Weight Tickets');
    expect(
      screen.getByText(
        'Weight tickets should include both an empty or full weight ticket for each segment or trip. If you’re missing a weight ticket, you’ll be able to use a government-created spreadsheet to estimate the weight.',
      ),
    ).toBeInTheDocument();
    expect(
      screen.getByText('Weight tickets must be certified, legible, and unaltered. Files must be 25MB or smaller.'),
    ).toBeInTheDocument();
    expect(screen.getByText('You must upload at least one set of weight tickets to get paid for your PPM.'));

    // renders form content
    expect(screen.getByRole('heading', { level: 2 })).toHaveTextContent('Trip 2');
  });

  it('replaces the router history with newly created weight ticket id', async () => {
    createWeightTicket.mockResolvedValue(mockWeightTicket);

    render(<WeightTickets />, { wrapper: MockProviders });

    await waitFor(() => {
      expect(mockReplace).toHaveBeenCalledWith(weightTicketsEditPath);
    });
  });

  it('routes back to home when finish later is clicked', async () => {
    createWeightTicket.mockResolvedValue(mockWeightTicket);

    render(<WeightTickets />, { wrapper: MockProviders });

    await waitFor(() => {
      expect(screen.getByRole('button', { name: 'Finish Later' })).toBeInTheDocument();
    });
    await userEvent.click(screen.getByRole('button', { name: 'Finish Later' }));
    expect(mockPush).toHaveBeenCalledWith(homePath);
  });

  it('calls patch weight ticket with the appropriate payload', async () => {
    createWeightTicket.mockResolvedValue(mockWeightTicketWithUploads);
    patchWeightTicket.mockResolvedValue({});

    render(<WeightTickets />, { wrapper: MockProviders });

    await waitFor(() => {
      expect(screen.getByRole('heading', { level: 2 })).toHaveTextContent('Trip 2');
    });
    await userEvent.type(screen.getByLabelText('Vehicle description'), 'DMC Delorean');
    await userEvent.type(screen.getByLabelText('Empty weight'), '4999');
    await userEvent.type(screen.getByLabelText('Full weight'), '6999');
    await userEvent.click(screen.getByLabelText('Yes'));
    await userEvent.click(screen.getAllByLabelText('Yes')[1]);

    await userEvent.click(screen.getByRole('button', { name: 'Save & Continue' }));

    await waitFor(() => {
      expect(patchWeightTicket).toHaveBeenCalledWith(
        mockMTOShipmentId,
        mockWeightTicketId,
        {
          shipmentId: mockWeightTicketWithUploads.ppmShipmentId,
          vehicleDescription: 'DMC Delorean',
          emptyWeight: 4999,
          missingEmptyWeightTicket: false,
          fullWeight: 6999,
          missingFullWeightTicket: false,
          hasOwnTrailer: true,
          hasClaimedTrailer: true,
        },
        mockWeightTicketETag,
      );
    });

    expect(mockPush).toHaveBeenCalledWith(reviewPath);
  });

  it('expect loadingPlaceholder when mtoShipment is falsy', async () => {
    selectMTOShipmentById.mockReturnValue(null);

    render(<WeightTickets />, { wrapper: MockProviders });

    await waitFor(() => {
      expect(screen.getByRole('heading', { level: 2 })).toHaveTextContent('Loading, please wait...');
    });
  });
});
