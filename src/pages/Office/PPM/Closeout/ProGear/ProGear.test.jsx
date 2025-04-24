import React from 'react';
import { render, waitFor, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { generatePath } from 'react-router-dom';
import { v4 } from 'uuid';

import { reviewWeightsQuery } from '../../../MoveTaskOrder/moveTaskOrderUnitTestData';

import { MockProviders } from 'testUtils';
import { servicesCounselingRoutes } from 'constants/routes';
import ProGear from 'pages/Office/PPM/Closeout/ProGear/ProGear';
import { createProGearWeightTicket, patchProGearWeightTicket, deleteUploadForDocument } from 'services/ghcApi';
import { SHIPMENT_OPTIONS } from 'shared/constants';
import { usePPMShipmentAndDocsOnlyQueries, useReviewShipmentWeightsQuery } from 'hooks/queries';

const mockMoveId = 'cc03c553-d317-46af-8b2d-3c9f899f6451';
const mockMTOShipmentId = '6b7a5769-4393-46fb-a4c4-d3f6ac7584c7';
const mockPPMShipmentId = v4();
const mockProGearWeightTicketId = v4();
const mockProGearWeightTicketETag = window.btoa(new Date());

const mockNavigate = jest.fn();
jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom'),
  useNavigate: () => mockNavigate,
}));

jest.mock('services/ghcApi', () => ({
  ...jest.requireActual('services/ghcApi'),
  createProGearWeightTicket: jest.fn(),
  createUploadForPPMDocument: jest.fn(),
  deleteUploadForDocument: jest.fn(),
  patchProGearWeightTicket: jest.fn(),
  updateMTOShipment: jest.fn(),
}));

jest.mock('hooks/queries', () => ({
  usePPMShipmentAndDocsOnlyQueries: jest.fn(),
  useReviewShipmentWeightsQuery: jest.fn(),
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

const mockProGearWeightTicket = {
  id: mockProGearWeightTicketId,
  ppmShipmentId: mockPPMShipmentId,
  weight: 123,
  description: 'Professional items',
  belongsToSelf: true,
  hasWeightTickets: true,
  eTag: mockProGearWeightTicketETag,
};

const mockDocumentId = v4();

const mockUploads = [
  {
    id: '299e2fb4-432d-4261-bbed-d8280c6090af',
    createdAt: '2022-06-22T23:25:50.490Z',
    bytes: 819200,
    url: 'a/fake/path',
    filename: 'weight_ticket.jpg',
    contentType: 'image/jpg',
  },
  {
    id: 'fd4e80f8-d025-44b2-8c33-15240fac51ab',
    createdAt: '2022-06-24T23:25:50.490Z',
    bytes: 204800,
    url: 'a/fake/path',
    filename: 'weight_ticket.pdf',
    contentType: 'application/pdf',
  },
];

const mockProGearWeightTicketWithUploads = {
  id: mockProGearWeightTicketId,
  ppmShipmentId: mockPPMShipmentId,
  belongsToSelf: false,
  documentId: mockDocumentId,
  document: {
    uploads: mockUploads,
  },
  eTag: mockProGearWeightTicketETag,
};

beforeEach(() => {
  jest.clearAllMocks();
});

const reviewPath = generatePath(servicesCounselingRoutes.BASE_SHIPMENT_PPM_REVIEW_PATH, {
  moveCode: mockMoveId,
  shipmentId: mockMTOShipmentId,
});

const renderEditProGearPage = () => {
  const mockRoutingConfig = {
    path: servicesCounselingRoutes.BASE_SHIPMENT_PPM_PRO_GEAR_EDIT_PATH,
    params: { moveCode: mockMoveId, shipmentId: mockMTOShipmentId, proGearId: mockProGearWeightTicketId },
  };

  render(
    <MockProviders {...mockRoutingConfig}>
      <ProGear />
    </MockProviders>,
  );
};

const renderProGearPage = () => {
  const mockRoutingConfig = {
    path: servicesCounselingRoutes.BASE_SHIPMENT_PPM_PRO_GEAR_PATH,
    params: { moveCode: mockMoveId, shipmentId: mockMTOShipmentId },
  };

  render(
    <MockProviders {...mockRoutingConfig}>
      <ProGear />
    </MockProviders>,
  );
};

describe('test page', () => {
  it('displays an error if the createProGearWeightTicket request fails', async () => {
    usePPMShipmentAndDocsOnlyQueries.mockReturnValue({
      isLoading: false,
      mtoShipment: mockMTOShipment,
      documents: { ProGearWeightTickets: [mockProGearWeightTicket] },
      isError: null,
      refetchMTOShipment: jest.fn(), // Mock the refetch function
    });

    useReviewShipmentWeightsQuery.mockReturnValue(reviewWeightsQuery);

    createProGearWeightTicket.mockRejectedValue('an error occurred');

    renderProGearPage();

    await waitFor(() => {
      expect(screen.getByText('Failed to create trip record')).toBeInTheDocument();
    });
  });
});

describe('Pro-gear page', () => {
  it('displays an error if the createProGearWeightTicket request fails', async () => {
    createProGearWeightTicket.mockRejectedValue('an error occurred');

    renderProGearPage();

    await waitFor(() => {
      expect(screen.getByText('Failed to create trip record')).toBeInTheDocument();
    });
  });

  it('displays the page', async () => {
    createProGearWeightTicket.mockResolvedValue(mockProGearWeightTicket);

    usePPMShipmentAndDocsOnlyQueries.mockReturnValue({
      isLoading: false,
      mtoShipment: mockMTOShipment,
      documents: { ProGearWeightTickets: [mockProGearWeightTicket] },
      isError: null,
    });
    renderEditProGearPage();

    expect(screen.getByRole('heading', { level: 1 })).toHaveTextContent('Pro-gear');
  });

  it('routes back to home when cancel button is clicked', async () => {
    createProGearWeightTicket.mockResolvedValue(mockProGearWeightTicket);

    usePPMShipmentAndDocsOnlyQueries.mockReturnValue({
      isLoading: false,
      mtoShipment: mockMTOShipment,
      documents: { ProGearWeightTickets: [mockProGearWeightTicket] },
      isError: null,
    });

    useReviewShipmentWeightsQuery.mockReturnValue(reviewWeightsQuery);

    renderEditProGearPage();

    await waitFor(() => {
      expect(screen.getByRole('button', { name: 'Cancel' })).toBeInTheDocument();
    });
    await userEvent.click(screen.getByRole('button', { name: 'Cancel' }));
    expect(mockNavigate).toHaveBeenCalledWith(reviewPath);
  });

  it('calls patchProGearWeightTicket with the appropriate payload', async () => {
    usePPMShipmentAndDocsOnlyQueries.mockReturnValue({
      isLoading: false,
      mtoShipment: mockMTOShipment,
      documents: { ProGearWeightTickets: [mockProGearWeightTicketWithUploads] },
      isError: null,
      refetchMTOShipment: jest.fn().mockImplementation(() => Promise.resolve(mockMTOShipment)), // Mock the refetch function
    });

    useReviewShipmentWeightsQuery.mockReturnValue(reviewWeightsQuery);
    renderEditProGearPage();

    await userEvent.click(screen.getByLabelText('My spouse'));

    await waitFor(() => {
      expect(screen.getByRole('heading', { level: 2 })).toHaveTextContent('Set 1');
    });
    await waitFor(() => {
      expect(screen.getByLabelText(/^Brief description of the pro-gear/)).toBeInTheDocument();
      expect(screen.getByLabelText(/I don't have weight tickets/)).toBeInTheDocument();
    });
    await userEvent.type(screen.getByLabelText(/^Brief description of the pro-gear/), 'Professional gear');
    await userEvent.type(screen.getByLabelText(/^Shipment's pro-gear weight/), '100');

    await waitFor(() => {
      expect(screen.getByRole('button', { name: 'Save & Continue' })).toBeEnabled();
    });

    await userEvent.click(screen.getByRole('button', { name: 'Save & Continue' }));

    await waitFor(() => {
      expect(patchProGearWeightTicket).toHaveBeenCalledWith({
        ppmShipmentId: mockPPMShipmentId,
        proGearWeightTicketId: mockProGearWeightTicketId,
        eTag: mockProGearWeightTicketETag,
        payload: {
          hasWeightTickets: true,
          belongsToSelf: false,
          ppmShipmentId: mockPPMShipmentId,
          shipmentType: 'PPM',
          shipmentLocator: undefined,
          description: 'Professional gear',
          weight: 100,
        },
      });
    });
  });

  it('calls the delete handler when removing an existing upload', async () => {
    usePPMShipmentAndDocsOnlyQueries.mockReturnValue({
      isLoading: false,
      mtoShipment: mockMTOShipment,
      documents: { ProGearWeightTickets: [{}, {}, {}, {}, mockProGearWeightTicketWithUploads] },
      isError: null,
    });
    useReviewShipmentWeightsQuery.mockReturnValue(reviewWeightsQuery);

    deleteUploadForDocument.mockResolvedValue({});
    renderEditProGearPage();

    let deleteButtons;
    await waitFor(() => {
      deleteButtons = screen.getAllByText('Delete');
      expect(deleteButtons).toHaveLength(2);
    });
    await userEvent.click(deleteButtons[0]);
    await waitFor(() => {
      expect(screen.queryByText('weight_ticket.jpg')).not.toBeInTheDocument();
    });
  });

  it('expect loadingPlaceholder when mtoShipment is falsy', async () => {
    usePPMShipmentAndDocsOnlyQueries.mockReturnValue({
      isLoading: false,
      mtoShipment: null,
      documents: { ProGearWeightTickets: [{}, {}, {}, {}, mockProGearWeightTicketWithUploads] },
      isError: null,
    });
    useReviewShipmentWeightsQuery.mockReturnValue(reviewWeightsQuery);

    renderEditProGearPage();

    await waitFor(() => {
      expect(screen.getByRole('heading', { level: 2 })).toHaveTextContent('Loading, please wait...');
    });
  });
});
