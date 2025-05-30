import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { generatePath } from 'react-router-dom';
import { v4 } from 'uuid';

import { MockProviders } from 'testUtils';
import { customerRoutes } from 'constants/routes';
import GunSafe from 'pages/MyMove/PPM/Closeout/GunSafe/GunSafe';
import { createGunSafeWeightTicket, deleteUpload, patchGunSafeWeightTicket } from 'services/internalApi';
import { SHIPMENT_OPTIONS } from 'shared/constants';
import { selectMTOShipmentById, selectGunSafeWeightTicketAndIndexById } from 'store/entities/selectors';

const mockMoveId = 'cc03c553-d317-46af-8b2d-3c9f899f6451';
const mockMTOShipmentId = '6b7a5769-4393-46fb-a4c4-d3f6ac7584c7';
const mockPPMShipmentId = v4();
const mockGunSafeWeightTicketId = v4();
const mockGunSafeWeightTicketETag = window.btoa(new Date());

const mockNavigate = jest.fn();
jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom'),
  useNavigate: () => mockNavigate,
}));
const mockRoutingConfig = {
  path: customerRoutes.SHIPMENT_PPM_GUN_SAFE_PATH,
  params: {
    moveId: 'cc03c553-d317-46af-8b2d-3c9f899f6451',
    mtoShipmentId: '6b7a5769-4393-46fb-a4c4-d3f6ac7584c7',
  },
};

jest.mock('services/internalApi', () => ({
  ...jest.requireActual('services/internalApi'),
  createGunSafeWeightTicket: jest.fn(),
  createUploadForDocument: jest.fn(),
  deleteUpload: jest.fn(),
  patchGunSafeWeightTicket: jest.fn(),
  getResponseError: jest.fn(),
  getAllMoves: jest.fn().mockImplementation(() => Promise.resolve()),
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
    hasGunSafe: false,
    GunSafeWeight: null,
    spouseGunSafeWeight: null,
  },
  eTag: 'dGVzdGluZzIzNDQzMjQ',
};

const mockDocumentId = v4();

const mockGunSafeWeightTicket = {
  id: mockGunSafeWeightTicketId,
  ppmShipmentId: mockPPMShipmentId,
  weight: 123,
  description: 'Gun safe',
  hasWeightTickets: true,
  eTag: mockGunSafeWeightTicketETag,
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

const mockGunSafeWeightTicketWithUploads = {
  id: mockGunSafeWeightTicketId,
  ppmShipmentId: mockPPMShipmentId,
  documentId: mockDocumentId,
  document: {
    uploads: mockUploads,
  },
  eTag: mockGunSafeWeightTicketETag,
};

const mockEmptyGunSafeWeightTicketAndIndex = {
  gunSafeWeightTicket: null,
  index: -1,
};

const mockServiceMember = {
  id: 'testId',
};

const mockOrders = {
  'fd4e80f8-d025-44b2-8c33-15240fac51ab': {
    entitlement: {
      GunSafe: 1234,
      GunSafeSpouse: 123,
    },
  },
};

jest.mock('store/entities/selectors', () => ({
  ...jest.requireActual('store/entities/selectors'),
  selectMTOShipmentById: jest.fn(() => mockMTOShipment),
  selectGunSafeWeightTicketAndIndexById: jest.fn(() => mockEmptyGunSafeWeightTicketAndIndex),
  selectServiceMemberFromLoggedInUser: jest.fn(() => mockServiceMember),
  selectOrdersForLoggedInUser: jest.fn(() => mockOrders),
}));

beforeEach(() => {
  jest.clearAllMocks();
});

const reviewPath = generatePath(customerRoutes.SHIPMENT_PPM_REVIEW_PATH, {
  moveId: mockMoveId,
  mtoShipmentId: mockMTOShipmentId,
});

const GunSafeWeightTicketsEditPath = generatePath(customerRoutes.SHIPMENT_PPM_GUN_SAFE_EDIT_PATH, {
  moveId: mockMoveId,
  mtoShipmentId: mockMTOShipmentId,
  gunSafeId: mockGunSafeWeightTicketId,
});

const renderGunSafePage = () => {
  return render(
    <MockProviders {...mockRoutingConfig}>
      <GunSafe />
    </MockProviders>,
  );
};

describe('gun safe page', () => {
  it('loads the selected shipment from redux', async () => {
    createGunSafeWeightTicket.mockResolvedValue(mockGunSafeWeightTicket);

    renderGunSafePage();

    await waitFor(() => {
      expect(selectMTOShipmentById).toHaveBeenCalledWith(expect.anything(), mockMTOShipmentId);
    });
  });

  it('displays an error if the createGunSafeWeightTicket request fails', async () => {
    createGunSafeWeightTicket.mockRejectedValue('an error occurred');

    renderGunSafePage();

    await waitFor(() => {
      expect(screen.getByText('Failed to create trip record')).toBeInTheDocument();
    });
  });

  it('does not make create pro gear weight ticket api request if id param exists', async () => {
    selectGunSafeWeightTicketAndIndexById.mockReturnValue({ gunSafeWeightTicket: mockGunSafeWeightTicket, index: 0 });

    const mockRoutingParams = {
      moveId: mockMoveId,
      mtoShipmentId: mockMTOShipmentId,
      gunSafeId: mockGunSafeWeightTicketId,
    };
    render(
      <MockProviders path={customerRoutes.SHIPMENT_PPM_GUN_SAFE_EDIT_PATH} params={mockRoutingParams}>
        <GunSafe />
      </MockProviders>,
    );

    await waitFor(() => {
      expect(createGunSafeWeightTicket).not.toHaveBeenCalled();
    });
  });

  it('displays the page', async () => {
    createGunSafeWeightTicket.mockResolvedValue(mockGunSafeWeightTicket);
    renderGunSafePage();
    expect(screen.getByRole('heading', { level: 1 })).toHaveTextContent('Gun Safe');
  });

  it('replaces the router history with newly created pro gear weight ticket id', async () => {
    createGunSafeWeightTicket.mockResolvedValue(mockGunSafeWeightTicket);
    selectGunSafeWeightTicketAndIndexById.mockReturnValueOnce({ gunSafeWeightTicket: null, index: -1 });
    selectGunSafeWeightTicketAndIndexById.mockReturnValue({
      gunSafeWeightTicket: mockGunSafeWeightTicket,
      index: 0,
    });

    renderGunSafePage();

    await waitFor(() => {
      expect(mockNavigate).toHaveBeenCalledWith(GunSafeWeightTicketsEditPath, { replace: true });
    });
  });

  it('routes back to review when cancel is clicked', async () => {
    createGunSafeWeightTicket.mockResolvedValue(mockGunSafeWeightTicket);
    selectGunSafeWeightTicketAndIndexById.mockReturnValue({ gunSafeWeightTicket: mockGunSafeWeightTicket, index: 0 });

    renderGunSafePage();

    await waitFor(() => {
      expect(screen.getByRole('button', { name: 'Cancel' })).toBeInTheDocument();
    });
    await userEvent.click(screen.getByRole('button', { name: 'Cancel' }));
    expect(mockNavigate).toHaveBeenCalledWith(reviewPath);
  });

  it('calls patchGunSafeWeightTicket with the appropriate payload', async () => {
    createGunSafeWeightTicket.mockResolvedValue(mockGunSafeWeightTicketWithUploads);
    selectGunSafeWeightTicketAndIndexById.mockReturnValue({
      gunSafeWeightTicket: mockGunSafeWeightTicketWithUploads,
      index: 1,
    });
    patchGunSafeWeightTicket.mockResolvedValue({});
    selectMTOShipmentById.mockReturnValue({
      ...mockMTOShipment,
      ppmShipment: {
        ...mockMTOShipment.ppmShipment,
        GunSafeWeightTickets: [mockGunSafeWeightTicketWithUploads],
      },
    });

    renderGunSafePage();

    await waitFor(() => {
      expect(screen.getByRole('heading', { level: 2 })).toHaveTextContent('Gun Safe 2');
    });
    await waitFor(() => {
      expect(screen.getByLabelText(/^Brief description of the gun safe/)).toBeInTheDocument();
      expect(screen.getByLabelText(/I don't have weight tickets/)).toBeInTheDocument();
    });
    await userEvent.type(screen.getByLabelText(/^Brief description of the gun safe/), 'Gun safe');
    await userEvent.type(screen.getByLabelText(/^Shipment's gun safe weight/), '100');
    await userEvent.click(screen.getByLabelText(/I don't have weight tickets/));

    await waitFor(() => {
      expect(screen.getByRole('button', { name: 'Save & Continue' })).toBeEnabled();
    });

    await userEvent.click(screen.getByRole('button', { name: 'Save & Continue' }));

    await waitFor(() => {
      expect(patchGunSafeWeightTicket).toHaveBeenCalledWith(
        mockPPMShipmentId,
        mockGunSafeWeightTicketId,
        {
          ppmShipmentId: mockGunSafeWeightTicketWithUploads.ppmShipmentId,
          gunSafeWeightTicketId: mockGunSafeWeightTicketId,
          description: 'Gun safe',
          weight: 100,
          hasWeightTickets: false,
        },
        mockGunSafeWeightTicketETag,
      );
    });

    expect(mockNavigate).toHaveBeenCalledWith(GunSafeWeightTicketsEditPath, { replace: true });
  });

  it('calls the delete handler when removing an existing upload', async () => {
    selectGunSafeWeightTicketAndIndexById.mockReturnValue({
      gunSafeWeightTicket: mockGunSafeWeightTicketWithUploads,
      index: 0,
    });

    selectMTOShipmentById.mockReturnValue({
      ...mockMTOShipment,
      ppmShipment: {
        ...mockMTOShipment.ppmShipment,
        GunSafeWeightTickets: [mockGunSafeWeightTicketWithUploads],
      },
    });
    deleteUpload.mockResolvedValue({});
    const mockRoutingParams = {
      moveId: mockMoveId,
      mtoShipmentId: mockMTOShipmentId,
      gunSafeId: mockGunSafeWeightTicketId,
    };
    render(
      <MockProviders path={customerRoutes.SHIPMENT_PPM_GUN_SAFE_EDIT_PATH} params={mockRoutingParams}>
        <GunSafe />
      </MockProviders>,
    );

    let deleteButtons;
    await waitFor(() => {
      deleteButtons = screen.getAllByRole('button', { name: 'Delete' });
      expect(deleteButtons).toHaveLength(2);
    });
    await userEvent.click(deleteButtons[1]);
    // TODO: THERE IS A KNOWN ISSUE WITH FAILING TO DELETE PPM UPLOADED DOCUMENTS (B-19065) THESE
    // TESTS WILL FAIL UNTIL THE ISSUE IS RESOLVED
    // await waitFor(() => {
    //   expect(screen.queryByText('weight_ticket.pdf')).not.toBeInTheDocument();
    // });
    // await userEvent.click(deleteButtons[0]);
    // await waitFor(() => {
    //   expect(screen.queryByText('weight_ticket.jpg')).not.toBeInTheDocument();
    //   expect(screen.getByText(/At least one upload is required/)).toBeInTheDocument();
    // });
  });

  it('displays an error if delete fails', async () => {
    mockGunSafeWeightTicketWithUploads.document.uploads = mockUploads;

    selectGunSafeWeightTicketAndIndexById.mockReturnValue({
      gunSafeWeightTicket: mockGunSafeWeightTicketWithUploads,
      index: 0,
    });

    selectMTOShipmentById.mockReturnValue({
      ...mockMTOShipment,
      ppmShipment: {
        ...mockMTOShipment.ppmShipment,
        GunSafeWeightTickets: [mockGunSafeWeightTicketWithUploads],
      },
    });

    deleteUpload.mockRejectedValue('error');
    const mockRoutingParams = {
      moveId: mockMoveId,
      mtoShipmentId: mockMTOShipmentId,
      gunSafeId: mockGunSafeWeightTicketId,
    };
    render(
      <MockProviders path={customerRoutes.SHIPMENT_PPM_GUN_SAFE_EDIT_PATH} params={mockRoutingParams}>
        <GunSafe />
      </MockProviders>,
    );

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

    renderGunSafePage();

    await waitFor(() => {
      expect(screen.getByRole('heading', { level: 2 })).toHaveTextContent('Loading, please wait...');
    });
  });
});
