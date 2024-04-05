import React from 'react';
import { render, waitFor, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { generatePath } from 'react-router-dom';
import { v4 } from 'uuid';

import { selectMTOShipmentById, selectWeightTicketAndIndexById } from 'store/entities/selectors';
import { customerRoutes } from 'constants/routes';
import { createWeightTicket, deleteUpload, patchWeightTicket } from 'services/internalApi';
import { MockProviders } from 'testUtils';
import { SHIPMENT_OPTIONS } from 'shared/constants';
import WeightTickets from 'pages/MyMove/PPM/Closeout/WeightTickets/WeightTickets';

const mockMoveId = 'cc03c553-d317-46af-8b2d-3c9f899f6451';
const mockMTOShipmentId = '6b7a5769-4393-46fb-a4c4-d3f6ac7584c7';
const mockPPMShipmentId = v4();
const mockWeightTicketId = v4();
const mockWeightTicketETag = window.btoa(new Date());

const mockNavigate = jest.fn();
jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom'),
  useNavigate: () => mockNavigate,
}));

jest.mock('services/internalApi', () => ({
  ...jest.requireActual('services/internalApi'),
  createWeightTicket: jest.fn(),
  createUploadForDocument: jest.fn(),
  deleteUpload: jest.fn(),
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
  emptyDocument: {
    uploads: [
      {
        id: '299e2fb4-432d-4261-bbed-d8280c6090af',
        createdAt: '2022-06-22T23:25:50.490Z',
        bytes: 819200,
        url: 'a/fake/path',
        filename: 'empty_weight.jpg',
        contentType: 'image/jpg',
      },
    ],
  },
  fullWeightDocumentId: mockFullWeightDocumentId,
  fullDocument: {
    uploads: [
      {
        id: 'f70af8a1-38e9-4ae2-a837-3c0c61069a0d',
        createdAt: '2022-06-23T23:25:50.490Z',
        bytes: 409600,
        url: 'a/fake/path',
        filename: 'full_weight.pdf',
        contentType: 'application/pdf',
      },
    ],
  },
  trailerOwnershipDocumentId: mockTrailerOwnershipWeightDocumentId,
  proofOfTrailerOwnershipDocument: {
    uploads: [
      {
        id: 'fd4e80f8-d025-44b2-8c33-15240fac51ab',
        createdAt: '2022-06-24T23:25:50.490Z',
        bytes: 204800,
        url: 'a/fake/path',
        filename: 'trailer_title.pdf',
        contentType: 'application/pdf',
      },
    ],
  },
  eTag: mockWeightTicketETag,
};

const mockEmptyWeightTicketAndIndex = {
  weightTicket: null,
  index: -1,
};

jest.mock('store/entities/selectors', () => ({
  ...jest.requireActual('store/entities/selectors'),
  selectMTOShipmentById: jest.fn(() => mockMTOShipment),
  selectWeightTicketAndIndexById: jest.fn(() => mockEmptyWeightTicketAndIndex),
}));

beforeEach(() => {
  jest.clearAllMocks();
});

const movePath = generatePath(customerRoutes.MOVE_HOME_PAGE);
const weightTicketsEditPath = generatePath(customerRoutes.SHIPMENT_PPM_WEIGHT_TICKETS_EDIT_PATH, {
  moveId: mockMoveId,
  mtoShipmentId: mockMTOShipmentId,
  weightTicketId: mockWeightTicketId,
});
const reviewPath = generatePath(customerRoutes.SHIPMENT_PPM_REVIEW_PATH, {
  moveId: mockMoveId,
  mtoShipmentId: mockMTOShipmentId,
});

const renderWeightTicketsPage = () => {
  const mockRoutingConfig = {
    path: customerRoutes.SHIPMENT_PPM_WEIGHT_TICKETS_PATH,
    params: { moveId: mockMoveId, mtoShipmentId: mockMTOShipmentId },
  };

  render(
    <MockProviders {...mockRoutingConfig}>
      <WeightTickets />
    </MockProviders>,
  );
};

const renderEditWeightTicketsPage = () => {
  const mockRoutingConfig = {
    path: customerRoutes.SHIPMENT_PPM_WEIGHT_TICKETS_EDIT_PATH,
    params: { moveId: mockMoveId, mtoShipmentId: mockMTOShipmentId, weightTicketId: mockWeightTicketId },
  };

  render(
    <MockProviders {...mockRoutingConfig}>
      <WeightTickets />
    </MockProviders>,
  );
};

describe('Weight Tickets page', () => {
  it('loads the selected shipment from redux', async () => {
    createWeightTicket.mockResolvedValue(mockWeightTicket);

    renderWeightTicketsPage();

    await waitFor(() => {
      expect(selectMTOShipmentById).toHaveBeenCalledWith(expect.anything(), mockMTOShipmentId);
    });
  });

  it('displays an error if the createWeightTicket request fails', async () => {
    createWeightTicket.mockRejectedValue('an error occurred');

    renderWeightTicketsPage();

    await waitFor(() => {
      expect(screen.getByText('Failed to create trip record')).toBeInTheDocument();
    });
  });

  it('does not make create weight ticket api request if id param exists', async () => {
    selectWeightTicketAndIndexById.mockReturnValue({ weightTicket: mockWeightTicket, index: 0 });

    renderEditWeightTicketsPage();

    await waitFor(() => {
      expect(createWeightTicket).not.toHaveBeenCalled();
    });
  });

  it('renders the page Content', async () => {
    createWeightTicket.mockResolvedValue(mockWeightTicket);
    selectWeightTicketAndIndexById.mockReturnValueOnce({ weightTicket: null, index: -1 });
    selectWeightTicketAndIndexById.mockReturnValue({ weightTicket: mockWeightTicket, index: 0 });

    renderEditWeightTicketsPage();

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
    expect(screen.getByRole('heading', { level: 2 })).toHaveTextContent('Trip 1');
  });

  it('replaces the router history with newly created weight ticket id', async () => {
    createWeightTicket.mockResolvedValue(mockWeightTicket);
    selectWeightTicketAndIndexById.mockReturnValueOnce({ weightTicket: null, index: -1 });
    selectWeightTicketAndIndexById.mockReturnValue({ weightTicket: mockWeightTicket, index: 0 });

    renderWeightTicketsPage();

    await waitFor(() => {
      expect(mockNavigate).toHaveBeenCalledWith(weightTicketsEditPath, { replace: true });
    });
  });

  it('routes back to home when return to homepage is clicked', async () => {
    createWeightTicket.mockResolvedValue(mockWeightTicket);
    selectWeightTicketAndIndexById.mockReturnValue({ weightTicket: mockWeightTicket, index: 0 });

    renderEditWeightTicketsPage();

    await waitFor(() => {
      expect(screen.getByRole('button', { name: 'Return To Homepage' })).toBeInTheDocument();
    });
    await userEvent.click(screen.getByRole('button', { name: 'Return To Homepage' }));
    expect(mockNavigate).toHaveBeenCalledWith(movePath);
  });

  it('calls patch weight ticket with the appropriate payload', async () => {
    createWeightTicket.mockResolvedValue(mockWeightTicketWithUploads);
    selectWeightTicketAndIndexById.mockReturnValue({ weightTicket: mockWeightTicketWithUploads, index: 1 });
    patchWeightTicket.mockResolvedValue({});

    renderWeightTicketsPage();

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
        mockPPMShipmentId,
        mockWeightTicketId,
        {
          ppmShipmentId: mockWeightTicketWithUploads.ppmShipmentId,
          vehicleDescription: 'DMC Delorean',
          emptyWeight: 4999,
          missingEmptyWeightTicket: false,
          fullWeight: 6999,
          missingFullWeightTicket: false,
          ownsTrailer: true,
          trailerMeetsCriteria: true,
        },
        mockWeightTicketETag,
      );
    });

    expect(mockNavigate).toHaveBeenCalledWith(reviewPath);
  });

  it('displays an error if patchWeightTicket fails', async () => {
    createWeightTicket.mockResolvedValue(mockWeightTicketWithUploads);
    selectWeightTicketAndIndexById.mockReturnValue({ weightTicket: mockWeightTicketWithUploads, index: 4 });
    patchWeightTicket.mockRejectedValueOnce('an error occurred');

    renderWeightTicketsPage();

    await waitFor(() => {
      expect(screen.getByRole('heading', { level: 2 })).toHaveTextContent('Trip 5');
    });
    await userEvent.type(screen.getByLabelText('Vehicle description'), 'DMC Delorean');
    await userEvent.type(screen.getByLabelText('Empty weight'), '4999');
    await userEvent.type(screen.getByLabelText('Full weight'), '6999');
    await userEvent.click(screen.getByLabelText('Yes'));
    await userEvent.click(screen.getAllByLabelText('Yes')[1]);

    await userEvent.click(screen.getByRole('button', { name: 'Save & Continue' }));

    await waitFor(() => {
      expect(screen.getByText('Failed to save updated trip record')).toBeInTheDocument();
    });
  });

  it('calls the delete handler when removing an existing upload', async () => {
    selectWeightTicketAndIndexById.mockReturnValue({ weightTicket: mockWeightTicketWithUploads, index: 0 });

    selectMTOShipmentById.mockReturnValue({
      ...mockMTOShipment,
      ppmShipment: {
        ...mockMTOShipment.ppmShipment,
        weightTickets: [mockWeightTicketWithUploads],
      },
    });
    deleteUpload.mockResolvedValue({});
    renderEditWeightTicketsPage();

    let deleteButtons;
    await waitFor(() => {
      deleteButtons = screen.getAllByRole('button', { name: 'Delete' });
      expect(deleteButtons).toHaveLength(2);
    });
    await userEvent.click(deleteButtons[0]);
    await waitFor(() => {
      expect(screen.queryByText('empty_weight.jpg')).not.toBeInTheDocument();
    });
  });

  it('expect loadingPlaceholder when mtoShipment is falsy', async () => {
    selectMTOShipmentById.mockReturnValueOnce(null);

    renderWeightTicketsPage();

    await waitFor(() => {
      expect(screen.getByRole('heading', { level: 2 })).toHaveTextContent('Loading, please wait...');
    });
  });
});
