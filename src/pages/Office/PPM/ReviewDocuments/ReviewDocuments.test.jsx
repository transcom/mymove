import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import { ReviewDocuments } from './ReviewDocuments';

import { usePPMShipmentDocsQueries } from 'hooks/queries';

Element.prototype.scrollTo = jest.fn();

beforeEach(() => {
  jest.clearAllMocks();
});

const mockPush = jest.fn();
jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom'),
  useHistory: () => ({
    push: mockPush,
  }),
}));

const mockPatchWeightTicket = jest.fn();
jest.mock('services/ghcApi', () => ({
  ...jest.requireActual('services/ghcApi'),
  patchWeightTicket: (options) => mockPatchWeightTicket(options),
}));

// prevents react-fileviewer from throwing errors without mocking relevant DOM elements
jest.mock('components/DocumentViewer/Content/Content', () => {
  const MockContent = () => <div>Content</div>;
  return MockContent;
});

const mockPDFUpload = {
  bytes: 0,
  contentType: 'application/pdf',
  createdAt: '2020-09-17T16:00:48.099137Z',
  filename: 'test.pdf',
  id: '10',
  status: 'PROCESSING',
  updatedAt: '2020-09-17T16:00:48.099142Z',
  url: '/storage/prime/99/uploads/10?contentType=application%2Fpdf',
};

const mockXLSUpload = {
  bytes: 0,
  contentType: 'application/vnd.ms-excel',
  createdAt: '2020-09-17T16:00:48.099137Z',
  filename: 'test.xls',
  id: '11',
  status: 'PROCESSING',
  updatedAt: '11',
  url: '/storage/prime/99/uploads/10?contentType=image%2Fjpeg',
};

const mockJPGUpload = {
  bytes: 0,
  contentType: 'image/jpeg',
  createdAt: '2020-09-17T16:00:48.099137Z',
  filename: 'test.jpg',
  id: '12',
  status: 'PROCESSING',
  updatedAt: '2020-09-17T16:00:48.099142Z',
  url: '/storage/prime/99/uploads/10?contentType=image%2Fjpg',
};

jest.mock('hooks/queries', () => ({
  usePPMShipmentDocsQueries: jest.fn(),
}));

const testShipmentId = '4321';
const ticketDocuments = {
  emptyDocument: {
    uploads: [mockPDFUpload],
    createdAt: '3',
  },
  fullDocument: {
    uploads: [mockXLSUpload],
    createdAt: '2',
  },
  proofOfTrailerOwnershipDocument: {
    uploads: [mockJPGUpload],
    createdAt: '1',
  },
};

const usePPMShipmentDocsQueriesReturnValue = {
  mtoShipment: {
    ppmShipment: {
      actualPickupPostalCode: '90210',
      actualDestinationPostalCode: '11201',
      actualMoveDate: '2022-03-16',
      hasReceivedAdvance: true,
      advanceAmountReceived: 340000,
    },
  },
  weightTickets: [
    {
      ...ticketDocuments,
      id: '321',
      missingFullWeightTicket: false,
      missingEmptyWeightTicket: false,
      ppmShipmentId: '123',
      vehicleDescription: '2022 Honda CR-V Hybrid',
    },
  ],
};

const rejectedPayload = {
  payload: {
    ppmShipmentId: '123',
    vehicleDescription: '2022 Honda CR-V Hybrid',
    emptyWeight: 14500,
    missingEmptyWeightTicket: false,
    fullWeight: 18500,
    missingFullWeightTicket: false,
    ownsTrailer: false,
    trailerMeetsCriteria: false,
    reason: 'reason',
    status: 'REJECTED',
  },
  ppmShipmentId: '123',
  weightTicketId: '321',
};

const usePPMShipmentDocsQueriesReturnValueMultiple = {
  ...usePPMShipmentDocsQueriesReturnValue,
  weightTickets: [
    {
      ...ticketDocuments,
    },
    { ...ticketDocuments },
  ],
};

const requiredProps = {
  match: { params: { shipmentId: testShipmentId, moveCode: 'READY' } },
};

const loadingReturnValue = {
  isLoading: true,
  isError: false,
  isSuccess: false,
};

const errorReturnValue = {
  isLoading: false,
  isError: true,
  isSuccess: false,
};

describe('ReviewDocuments', () => {
  describe('check loading and error component states', () => {
    it('renders the Loading Placeholder when the query is still loading', async () => {
      usePPMShipmentDocsQueries.mockReturnValue(loadingReturnValue);
      render(<ReviewDocuments {...requiredProps} />);
      const h2 = await screen.getByRole('heading', { name: 'Loading, please wait...', level: 2 });
      expect(h2).toBeInTheDocument();
    });
    it('renders the Something Went Wrong component when the query errors', async () => {
      usePPMShipmentDocsQueries.mockReturnValue(errorReturnValue);
      render(<ReviewDocuments {...requiredProps} />);
      const errorMessage = await screen.getByText(/Something went wrong./);
      expect(errorMessage).toBeInTheDocument();
    });
  });
  describe('with data loaded', () => {
    it('renders the DocumentViewer', async () => {
      await waitFor(() => {
        usePPMShipmentDocsQueries.mockReturnValue(usePPMShipmentDocsQueriesReturnValue);
        render(<ReviewDocuments {...requiredProps} />);
      });
      const docs = screen.getByText(/Documents/);
      expect(docs).toBeInTheDocument();
      expect(screen.getAllByText('test.pdf').length).toBe(2);
      expect(screen.getByText('test.xls')).toBeInTheDocument();
      expect(screen.getByText('test.jpg')).toBeInTheDocument();
      expect(screen.getByRole('heading', { level: 2, name: '1 of 1 Document Sets' })).toBeInTheDocument();
      expect(screen.getByRole('button', { name: 'Continue' })).toBeInTheDocument();
      expect(screen.getByRole('button', { name: 'Back' })).toBeInTheDocument();
      expect(screen.getByTestId('closeSidebar')).toBeInTheDocument();
    });
    it('renders and handles the Continue button with the appropriate payload', async () => {
      usePPMShipmentDocsQueries.mockReturnValue(usePPMShipmentDocsQueriesReturnValue);
      const { getByRole, getByLabelText, queryByText } = render(<ReviewDocuments {...requiredProps} />);
      await waitFor(() => {
        expect(getByLabelText('Accept')).toBeInTheDocument();
        expect(getByLabelText('Reject')).toBeInTheDocument();
        expect(getByRole('button', { name: 'Continue' })).toBeInTheDocument();
      });
      await userEvent.type(getByRole('textbox', { name: 'Empty weight' }), '14500');
      await userEvent.type(getByRole('textbox', { name: 'Full weight' }), '18500');
      expect(queryByText(/4,000 lbs/)).toBeInTheDocument();
      await userEvent.click(getByLabelText('Reject'));
      await waitFor(() => {
        expect(getByLabelText('Reason')).toBeInTheDocument();
      });
      await userEvent.type(getByLabelText('Reason'), 'reason');
      await userEvent.click(getByRole('button', { name: 'Continue' }));
      expect(queryByText('Reviewing this weight ticket is required')).not.toBeInTheDocument();
      expect(mockPatchWeightTicket).toHaveBeenCalledWith(rejectedPayload);
      expect(mockPush).toHaveBeenCalled();
    });
    it('renders and handles the Close button', async () => {
      usePPMShipmentDocsQueries.mockReturnValue(usePPMShipmentDocsQueriesReturnValue);
      const { getByTestId } = render(<ReviewDocuments {...requiredProps} />);
      expect(getByTestId('closeSidebar')).toBeInTheDocument();
      await userEvent.click(getByTestId('closeSidebar'));
      expect(mockPush).toHaveBeenCalled();
    });
  });
  describe('with multiple document sets loaded', () => {
    it('renders and handles the Accept button', async () => {
      usePPMShipmentDocsQueries.mockReturnValue(usePPMShipmentDocsQueriesReturnValueMultiple);
      const { getByRole, getByLabelText } = render(<ReviewDocuments {...requiredProps} />);
      expect(getByRole('heading', { level: 2, text: '1 of 2 Document Sets' }));
      expect(getByLabelText('Accept')).toBeInTheDocument();
      expect(getByLabelText('Reject')).toBeInTheDocument();
      expect(getByRole('button', { name: 'Continue' })).toBeInTheDocument();
      expect(getByRole('button', { name: 'Back' })).toBeInTheDocument();
      await userEvent.type(getByRole('textbox', { name: 'Empty weight' }), '14500');
      await userEvent.type(getByRole('textbox', { name: 'Full weight' }), '18500');
      await userEvent.click(getByLabelText('Accept'));
      await userEvent.click(getByRole('button', { name: 'Continue' }));
      expect(getByRole('heading', { level: 2, text: '2 of 2 Document Sets' }));
    });
    it('renders and handles the Back button', async () => {
      usePPMShipmentDocsQueries.mockReturnValue(usePPMShipmentDocsQueriesReturnValueMultiple);
      const { getByRole, getByLabelText } = render(<ReviewDocuments {...requiredProps} />);
      expect(getByRole('heading', { level: 2, text: '1 of 2 Document Sets' }));
      expect(getByLabelText('Accept')).toBeInTheDocument();
      expect(getByLabelText('Reject')).toBeInTheDocument();
      expect(getByRole('button', { name: 'Continue' })).toBeInTheDocument();
      expect(getByRole('button', { name: 'Back' })).toBeInTheDocument();
      await userEvent.type(getByRole('textbox', { name: 'Empty weight' }), '14500');
      await userEvent.type(getByRole('textbox', { name: 'Full weight' }), '18500');
      await userEvent.click(getByLabelText('Accept'));
      await userEvent.click(getByRole('button', { name: 'Continue' }));
      expect(getByRole('heading', { level: 2, text: '2 of 2 Document Sets' }));
      await userEvent.click(getByRole('button', { name: 'Back' }));
      expect(getByRole('heading', { level: 2, text: '1 of 2 Document Sets' }));
    });
  });
});
