import React from 'react';
import { fireEvent, render, screen, waitFor } from '@testing-library/react';

import { ReviewDocuments } from './ReviewDocuments';

import { usePPMShipmentDocsQueries } from 'hooks/queries';

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
  },
  fullDocument: {
    uploads: [mockXLSUpload],
  },
  proofOfTrailerOwnershipDocument: {
    uploads: [mockJPGUpload],
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
    },
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
      usePPMShipmentDocsQueries.mockReturnValue(usePPMShipmentDocsQueriesReturnValue);
      render(<ReviewDocuments {...requiredProps} />);

      const docs = await screen.getByText(/Documents/);
      expect(docs).toBeInTheDocument();
      expect(screen.getAllByText('test.pdf').length).toBe(2);
      expect(screen.getByText('test.xls')).toBeInTheDocument();
      expect(screen.getByText('test.jpg')).toBeInTheDocument();

      expect(screen.getByRole('heading', { level: 2, name: '1 of 1 Document Sets' })).toBeInTheDocument();
      expect(screen.getByRole('button', { name: 'Continue' })).toBeInTheDocument();
      expect(screen.getByRole('button', { name: 'Back' })).toBeInTheDocument();
    });
  });

  describe('returns to Review page when', () => {
    it('Close button is clicked', async () => {
      usePPMShipmentDocsQueries.mockReturnValue(usePPMShipmentDocsQueriesReturnValue);
      render(<ReviewDocuments {...requiredProps} />);

      const closeButton = await screen.getByTestId('closeSidebar');
      await waitFor(() => {
        expect(closeButton).toBeInTheDocument();
      });
      await fireEvent.click(closeButton);
      expect(mockPush).toHaveBeenCalled();
    });
  });

  describe('shows an error when review is invalid', () => {
    it('Continue button is clicked', async () => {
      usePPMShipmentDocsQueries.mockReturnValue(usePPMShipmentDocsQueriesReturnValue);
      render(<ReviewDocuments {...requiredProps} />);

      const continueButton = await screen.getByRole('button', { name: 'Continue' });
      await waitFor(() => {
        expect(continueButton).toBeInTheDocument();
      });
      await fireEvent.click(continueButton);
      expect(screen.getByText('Reviewing this weight ticket is required')).toBeInTheDocument();
    });
  });
});
