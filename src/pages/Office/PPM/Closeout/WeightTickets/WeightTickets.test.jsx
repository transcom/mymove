import React from 'react';
import { render, waitFor, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { generatePath } from 'react-router-dom';
import { v4 } from 'uuid';

import { servicesCounselingRoutes } from 'constants/routes';
import { createWeightTicket, deleteUploadForDocument, patchWeightTicket } from 'services/ghcApi';
import { MockProviders } from 'testUtils';
import { SHIPMENT_OPTIONS } from 'shared/constants';
import WeightTickets from 'pages/Office/PPM/Closeout/WeightTickets/WeightTickets';
import { usePPMShipmentAndDocsOnlyQueries } from 'hooks/queries';

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

jest.mock('services/ghcApi', () => ({
  ...jest.requireActual('services/ghcApi'),
  createWeightTicket: jest.fn(),
  createUploadForPPMDocument: jest.fn(),
  deleteUploadForDocument: jest.fn(),
  patchWeightTicket: jest.fn(),
}));

jest.mock('hooks/queries', () => ({
  usePPMShipmentAndDocsOnlyQueries: jest.fn(),
}));

const mockMTOShipment = {
  id: mockMTOShipmentId,
  shipmentType: SHIPMENT_OPTIONS.PPM,
  ppmShipment: {
    id: mockPPMShipmentId,
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

beforeEach(() => {
  jest.clearAllMocks();
});

const weightTicketsEditPath = generatePath(servicesCounselingRoutes.BASE_SHIPMENT_PPM_WEIGHT_TICKETS_EDIT_PATH, {
  moveCode: mockMoveId,
  shipmentId: mockMTOShipmentId,
  weightTicketId: mockWeightTicketId,
});
const reviewPath = generatePath(servicesCounselingRoutes.BASE_SHIPMENT_PPM_REVIEW_PATH, {
  moveCode: mockMoveId,
  shipmentId: mockMTOShipmentId,
});

const renderWeightTicketsPage = () => {
  const mockRoutingConfig = {
    path: servicesCounselingRoutes.BASE_SHIPMENT_PPM_WEIGHT_TICKETS_PATH,
    params: { moveCode: mockMoveId, shipmentId: mockMTOShipmentId },
  };

  render(
    <MockProviders {...mockRoutingConfig}>
      <WeightTickets />
    </MockProviders>,
  );
};

const renderEditWeightTicketsPage = () => {
  const mockRoutingConfig = {
    path: servicesCounselingRoutes.BASE_SHIPMENT_PPM_WEIGHT_TICKETS_EDIT_PATH,
    params: { moveCode: mockMoveId, shipmentId: mockMTOShipmentId, weightTicketId: mockWeightTicketId },
  };

  render(
    <MockProviders {...mockRoutingConfig}>
      <WeightTickets />
    </MockProviders>,
  );
};

describe('Weight Tickets page', () => {
  it('displays an error if the createWeightTicket request fails', async () => {
    usePPMShipmentAndDocsOnlyQueries.mockReturnValue({
      isLoading: false,
      mtoShipment: mockMTOShipment,
      documents: { WeightTickets: [mockWeightTicket] },
      isError: null,
    });

    createWeightTicket.mockRejectedValue('an error occurred');

    renderWeightTicketsPage();

    await waitFor(() => {
      expect(screen.getByText('Failed to create trip record')).toBeInTheDocument();
    });
  });

  it('does not make create weight ticket api request if id param exists', async () => {
    usePPMShipmentAndDocsOnlyQueries.mockReturnValue({
      isLoading: false,
      mtoShipment: mockMTOShipment,
      documents: { WeightTickets: [mockWeightTicket] },
      isError: null,
    });

    renderEditWeightTicketsPage();

    await waitFor(() => {
      expect(createWeightTicket).not.toHaveBeenCalled();
    });
  });

  it('renders the page Content', async () => {
    createWeightTicket.mockResolvedValue(mockWeightTicket);
    usePPMShipmentAndDocsOnlyQueries.mockReturnValue({
      isLoading: false,
      mtoShipment: mockMTOShipment,
      documents: { WeightTickets: [mockWeightTicket] },
      isError: null,
    });

    renderEditWeightTicketsPage();

    await waitFor(() => {
      expect(screen.getByTestId('tag')).toHaveTextContent('PPM');
    });

    expect(screen.getByRole('heading', { level: 1 })).toHaveTextContent('Weight Tickets');

    // renders form content
    expect(screen.getByRole('heading', { level: 2 })).toHaveTextContent('Trip 1');
  });

  it('replaces the router history with newly created weight ticket id', async () => {
    createWeightTicket.mockResolvedValue(mockWeightTicket);
    usePPMShipmentAndDocsOnlyQueries.mockReturnValue({
      isLoading: false,
      mtoShipment: mockMTOShipment,
      documents: { WeightTickets: [mockWeightTicket] },
      isError: null,
    });

    renderWeightTicketsPage();

    await waitFor(() => {
      expect(mockNavigate).toHaveBeenCalledWith(weightTicketsEditPath, { replace: true });
    });
  });

  it('routes back to review page when cancel is clicked', async () => {
    createWeightTicket.mockResolvedValue(mockWeightTicket);
    usePPMShipmentAndDocsOnlyQueries.mockReturnValue({
      isLoading: false,
      mtoShipment: mockMTOShipment,
      documents: { WeightTickets: [mockWeightTicket] },
      isError: null,
    });

    renderEditWeightTicketsPage();

    await waitFor(() => {
      expect(screen.getByRole('button', { name: 'Cancel' })).toBeInTheDocument();
    });
    await userEvent.click(screen.getByRole('button', { name: 'Cancel' }));
    expect(mockNavigate).toHaveBeenCalledWith(reviewPath);
  });

  it('calls patch weight ticket with the appropriate payload', async () => {
    usePPMShipmentAndDocsOnlyQueries.mockReturnValue({
      isLoading: false,
      mtoShipment: mockMTOShipment,
      documents: { WeightTickets: [mockWeightTicketWithUploads] },
      isError: null,
    });

    patchWeightTicket.mockResolvedValue();

    renderEditWeightTicketsPage();

    await waitFor(() => {
      expect(screen.getByRole('heading', { level: 2 })).toHaveTextContent('Trip 1');
    });
    await userEvent.type(screen.getByLabelText('Vehicle description *'), 'DMC Delorean');
    await userEvent.type(screen.getByLabelText('Empty weight *'), '4999');
    await userEvent.type(screen.getByLabelText('Full weight *'), '6999');
    await userEvent.click(screen.getByLabelText('Yes'));
    await userEvent.click(screen.getAllByLabelText('Yes')[1]);

    await userEvent.click(screen.getByRole('button', { name: 'Save & Continue' }));

    await waitFor(() => {
      expect(patchWeightTicket).toHaveBeenCalledWith({
        ppmShipmentId: mockPPMShipmentId,
        weightTicketId: mockWeightTicketId,
        payload: {
          ppmShipmentId: mockWeightTicketWithUploads.ppmShipmentId,
          vehicleDescription: 'DMC Delorean',
          emptyWeight: 4999,
          missingEmptyWeightTicket: false,
          fullWeight: 6999,
          missingFullWeightTicket: false,
          ownsTrailer: true,
          trailerMeetsCriteria: true,
        },
        eTag: mockWeightTicketETag,
      });
    });

    expect(mockNavigate).toHaveBeenCalledWith(reviewPath);
  });

  it('displays an error if patchWeightTicket fails', async () => {
    usePPMShipmentAndDocsOnlyQueries.mockReturnValue({
      isLoading: false,
      mtoShipment: mockMTOShipment,
      documents: { WeightTickets: [{}, {}, {}, {}, mockWeightTicketWithUploads] },
      isError: null,
    });
    patchWeightTicket.mockRejectedValueOnce('an error occurred');

    renderEditWeightTicketsPage();

    await waitFor(() => {
      expect(screen.getByRole('heading', { level: 2 })).toHaveTextContent('Trip 5');
    });
    await userEvent.type(screen.getByLabelText('Vehicle description *'), 'DMC Delorean');
    await userEvent.type(screen.getByLabelText('Empty weight *'), '4999');
    await userEvent.type(screen.getByLabelText('Full weight *'), '6999');
    await userEvent.click(screen.getByLabelText('Yes'));
    await userEvent.click(screen.getAllByLabelText('Yes')[1]);

    await userEvent.click(screen.getByRole('button', { name: 'Save & Continue' }));

    await waitFor(() => {
      expect(screen.getByText('Failed to save updated trip record')).toBeInTheDocument();
    });
  });

  it('calls the delete handler when removing an existing upload', async () => {
    usePPMShipmentAndDocsOnlyQueries.mockReturnValue({
      isLoading: false,
      mtoShipment: mockMTOShipment,
      documents: { WeightTickets: [{}, {}, {}, {}, mockWeightTicketWithUploads] },
      isError: null,
    });
    deleteUploadForDocument.mockResolvedValue({});
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
    usePPMShipmentAndDocsOnlyQueries.mockReturnValue({
      isLoading: false,
      mtoShipment: null,
      documents: { WeightTickets: [mockWeightTicketWithUploads] },
      isError: null,
    });

    renderWeightTicketsPage();

    await waitFor(() => {
      expect(screen.getByRole('heading', { level: 2 })).toHaveTextContent('Loading, please wait...');
    });
  });
});
