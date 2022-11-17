import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { useParams, generatePath } from 'react-router-dom-old';
import { v4 } from 'uuid';

import { MockProviders } from 'testUtils';
import { customerRoutes, generalRoutes } from 'constants/routes';
import ProGear from 'pages/MyMove/PPM/Closeout/ProGear/ProGear';
import { createProGearWeightTicket, deleteUpload, patchProGearWeightTicket } from 'services/internalApi';
import { SHIPMENT_OPTIONS } from 'shared/constants';
import { selectMTOShipmentById, selectProGearWeightTicketAndIndexById } from 'store/entities/selectors';

const mockMoveId = 'cc03c553-d317-46af-8b2d-3c9f899f6451';
const mockMTOShipmentId = '6b7a5769-4393-46fb-a4c4-d3f6ac7584c7';
const mockPPMShipmentId = v4();
const mockProGearWeightTicketId = v4();
const mockProGearWeightTicketETag = window.btoa(new Date());

const mockPush = jest.fn();
const mockReplace = jest.fn();
jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom'),
  useHistory: () => ({
    push: mockPush,
    replace: mockReplace,
  }),
  useParams: jest.fn(() => ({
    moveId: 'cc03c553-d317-46af-8b2d-3c9f899f6451',
    mtoShipmentId: '6b7a5769-4393-46fb-a4c4-d3f6ac7584c7',
  })),
}));

jest.mock('services/internalApi', () => ({
  ...jest.requireActual('services/internalApi'),
  createProGearWeightTicket: jest.fn(),
  createUploadForDocument: jest.fn(),
  deleteUpload: jest.fn(),
  patchProGearWeightTicket: jest.fn(),
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

const mockDocumentId = v4();

const mockProGearWeightTicket = {
  id: mockProGearWeightTicketId,
  ppmShipmentId: mockPPMShipmentId,
  weight: 123,
  description: 'Professional items',
  belongsToSelf: true,
  hasWeightTickets: true,
  eTag: mockProGearWeightTicketETag,
};

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

const mockEmptyProGearWeightTicketAndIndex = {
  proGearWeightTicket: null,
  index: -1,
};

const mockEntitlement = {
  proGear: 1234,
  proGearSpouse: 123,
};

jest.mock('store/entities/selectors', () => ({
  ...jest.requireActual('store/entities/selectors'),
  selectMTOShipmentById: jest.fn(() => mockMTOShipment),
  selectProGearWeightTicketAndIndexById: jest.fn(() => mockEmptyProGearWeightTicketAndIndex),
  selectProGearEntitlements: jest.fn(() => mockEntitlement),
}));

beforeEach(() => {
  jest.clearAllMocks();
});

const homePath = generatePath(generalRoutes.HOME_PATH);
const proGearWeightTicketsEditPath = generatePath(customerRoutes.SHIPMENT_PPM_PRO_GEAR_EDIT_PATH, {
  moveId: mockMoveId,
  mtoShipmentId: mockMTOShipmentId,
  proGearId: mockProGearWeightTicketId,
});
const reviewPath = generatePath(customerRoutes.SHIPMENT_PPM_REVIEW_PATH, {
  moveId: mockMoveId,
  mtoShipmentId: mockMTOShipmentId,
});

describe('Pro-gear page', () => {
  it('loads the selected shipment from redux', async () => {
    createProGearWeightTicket.mockResolvedValue(mockProGearWeightTicket);

    render(<ProGear />, { wrapper: MockProviders });

    await waitFor(() => {
      expect(selectMTOShipmentById).toHaveBeenCalledWith(expect.anything(), mockMTOShipmentId);
    });
  });

  it('displays an error if the createProGearWeightTicket request fails', async () => {
    createProGearWeightTicket.mockRejectedValue('an error occurred');

    render(<ProGear />, { wrapper: MockProviders });

    await waitFor(() => {
      expect(screen.getByText('Failed to create trip record')).toBeInTheDocument();
    });
  });

  it('does not make create pro gear weight ticket api request if id param exists', async () => {
    useParams.mockImplementationOnce(() => ({
      moveId: mockMoveId,
      mtoShipmentId: mockMTOShipmentId,
      proGearId: mockProGearWeightTicketId,
    }));
    selectProGearWeightTicketAndIndexById.mockReturnValue({ proGearWeightTicket: mockProGearWeightTicket, index: 0 });

    render(<ProGear />, { wrapper: MockProviders });

    await waitFor(() => {
      expect(createProGearWeightTicket).not.toHaveBeenCalled();
    });
  });

  it('displays the page', async () => {
    createProGearWeightTicket.mockResolvedValue(mockProGearWeightTicket);
    render(<ProGear />, { wrapper: MockProviders });
    expect(screen.getByRole('heading', { level: 1 })).toHaveTextContent('Pro-gear');
  });
  it('displays reminder to include pro-gear weight in total', () => {
    render(<ProGear />, { wrapper: MockProviders });
    expect(screen.getByText(/This pro-gear should be included in your total weight moved./)).toBeInTheDocument();
  });

  it('replaces the router history with newly created pro gear weight ticket id', async () => {
    createProGearWeightTicket.mockResolvedValue(mockProGearWeightTicket);
    selectProGearWeightTicketAndIndexById.mockReturnValueOnce({ proGearWeightTicket: null, index: -1 });
    selectProGearWeightTicketAndIndexById.mockReturnValue({
      proGearWeightTicket: mockProGearWeightTicket,
      index: 0,
    });

    render(<ProGear />, { wrapper: MockProviders });

    await waitFor(() => {
      expect(mockReplace).toHaveBeenCalledWith(proGearWeightTicketsEditPath);
    });
  });

  it('routes back to home when return to homepage is clicked', async () => {
    createProGearWeightTicket.mockResolvedValue(mockProGearWeightTicket);
    selectProGearWeightTicketAndIndexById.mockReturnValue({ proGearWeightTicket: mockProGearWeightTicket, index: 0 });

    render(<ProGear />, { wrapper: MockProviders });

    await waitFor(() => {
      expect(screen.getByRole('button', { name: 'Return To Homepage' })).toBeInTheDocument();
    });
    await userEvent.click(screen.getByRole('button', { name: 'Return To Homepage' }));
    expect(mockPush).toHaveBeenCalledWith(homePath);
  });

  it('calls patchProGearWeightTicket with the appropriate payload', async () => {
    createProGearWeightTicket.mockResolvedValue(mockProGearWeightTicketWithUploads);
    selectProGearWeightTicketAndIndexById.mockReturnValue({
      proGearWeightTicket: mockProGearWeightTicketWithUploads,
      index: 1,
    });
    patchProGearWeightTicket.mockResolvedValue({});

    render(<ProGear />, { wrapper: MockProviders });
    await userEvent.click(screen.getByLabelText('My spouse'));
    await waitFor(() => {
      expect(screen.getByRole('heading', { level: 2 })).toHaveTextContent('Set 2');
    });
    await waitFor(() => {
      expect(screen.getByLabelText(/^Brief description of the pro-gear/)).toBeInTheDocument();
      expect(screen.getByLabelText(/I don't have weight tickets/)).toBeInTheDocument();
    });
    await userEvent.type(screen.getByLabelText(/^Brief description of the pro-gear/), 'Professional gear');
    await userEvent.type(screen.getByLabelText(/^Shipment's pro-gear weight/), '100');
    await userEvent.click(screen.getByLabelText(/I don't have weight tickets/));

    await waitFor(() => {
      expect(screen.getByRole('button', { name: 'Save & Continue' })).toBeEnabled();
    });

    await userEvent.click(screen.getByRole('button', { name: 'Save & Continue' }));

    await waitFor(() => {
      expect(patchProGearWeightTicket).toHaveBeenCalledWith(
        mockPPMShipmentId,
        mockProGearWeightTicketId,
        {
          ppmShipmentId: mockProGearWeightTicketWithUploads.ppmShipmentId,
          proGearWeightTicketId: mockProGearWeightTicketId,
          description: 'Professional gear',
          belongsToSelf: false,
          weight: 100,
          hasWeightTickets: false,
        },
        mockProGearWeightTicketETag,
      );
    });

    expect(mockPush).toHaveBeenCalledWith(reviewPath);
  });

  it('calls the delete handler when removing an existing upload', async () => {
    useParams.mockImplementation(() => ({
      moveId: mockMoveId,
      mtoShipmentId: mockMTOShipmentId,
      proGearId: mockProGearWeightTicketId,
    }));
    selectProGearWeightTicketAndIndexById.mockReturnValue({
      proGearWeightTicket: mockProGearWeightTicketWithUploads,
      index: 0,
    });

    selectMTOShipmentById.mockReturnValue({
      ...mockMTOShipment,
      ppmShipment: {
        ...mockMTOShipment.ppmShipment,
        proGearWeightTickets: [mockProGearWeightTicketWithUploads],
      },
    });
    deleteUpload.mockResolvedValue({});
    render(<ProGear />, { wrapper: MockProviders });

    let deleteButtons;
    await waitFor(() => {
      deleteButtons = screen.getAllByRole('button', { name: 'Delete' });
      expect(deleteButtons).toHaveLength(2);
    });
    await userEvent.click(deleteButtons[1]);
    await waitFor(() => {
      expect(screen.queryByText('weight_ticket.pdf')).not.toBeInTheDocument();
    });
    await userEvent.click(deleteButtons[0]);
    await waitFor(() => {
      expect(screen.queryByText('weight_ticket.jpg')).not.toBeInTheDocument();
      expect(screen.getByText(/At least one upload is required/)).toBeInTheDocument();
    });
  });

  it('displays an error if delete fails', async () => {
    mockProGearWeightTicketWithUploads.document.uploads = mockUploads;
    useParams.mockImplementation(() => ({
      moveId: mockMoveId,
      mtoShipmentId: mockMTOShipmentId,
      proGearId: mockProGearWeightTicketId,
    }));

    selectProGearWeightTicketAndIndexById.mockReturnValue({
      proGearWeightTicket: mockProGearWeightTicketWithUploads,
      index: 0,
    });

    selectMTOShipmentById.mockReturnValue({
      ...mockMTOShipment,
      ppmShipment: {
        ...mockMTOShipment.ppmShipment,
        proGearWeightTickets: [mockProGearWeightTicketWithUploads],
      },
    });

    deleteUpload.mockRejectedValue('error');
    render(<ProGear />, { wrapper: MockProviders });

    let deleteButtons;
    await waitFor(() => {
      deleteButtons = screen.getAllByRole('button', { name: 'Delete' });
      expect(deleteButtons).toHaveLength(2);
    });
    await userEvent.click(deleteButtons[1]);
    await waitFor(() => {
      expect(screen.getByText(/Failed to delete the file upload/)).toBeInTheDocument();
    });
  });

  it('expect loadingPlaceholder when mtoShipment is falsy', async () => {
    selectMTOShipmentById.mockReturnValueOnce(null);

    render(<ProGear />, { wrapper: MockProviders });

    await waitFor(() => {
      expect(screen.getByRole('heading', { level: 2 })).toHaveTextContent('Loading, please wait...');
    });
  });
});
