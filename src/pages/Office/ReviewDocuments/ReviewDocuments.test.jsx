import React from 'react';
import { mount } from 'enzyme';
import { render, screen } from '@testing-library/react';

import { ReviewDocuments } from './ReviewDocuments';

import { useMoveDetailsQueries } from 'hooks/queries';

const mockPDFUpload = {
  contentType: 'application/pdf',
  createdAt: '2020-09-17T16:00:48.099137Z',
  filename: 'test.pdf',
  id: '10',
  status: 'PROCESSING',
  updatedAt: '2020-09-17T16:00:48.099142Z',
  url: '/storage/prime/99/uploads/10?contentType=application%2Fpdf',
};

const mockXLSUpload = {
  contentType: 'application/vnd.ms-excel',
  createdAt: '2020-09-17T16:00:48.099137Z',
  filename: 'test.xls',
  id: '11',
  status: 'PROCESSING',
  updatedAt: '11',
  url: '/storage/prime/99/uploads/10?contentType=image%2Fjpeg',
};

jest.mock('hooks/queries', () => ({
  useMoveDetailsQueries: jest.fn(),
}));

const testShipmentId = '4321';
const useMoveDetailsQueriesReturnValue = {
  moveCode: 'READY',
  mtoShipments: [
    {
      id: testShipmentId,
      status: 'SUBMITTED',
      moveTaskOrderID: '123',
      ppmShipment: {
        city: 'Beverly Hills',
        id: '0cf43b1f-04e8-4c56-a6a1-06aec192ca07',
        weightTickets: [
          {
            emptyDocument: {
              uploads: [mockPDFUpload],
            },
            fullDocument: {
              uploads: [mockXLSUpload],
            },
          },
        ],
      },
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
      useMoveDetailsQueries.mockReturnValue(loadingReturnValue);

      render(<ReviewDocuments {...requiredProps} />);

      const h2 = await screen.getByRole('heading', { name: 'Loading, please wait...', level: 2 });
      expect(h2).toBeInTheDocument();
    });

    it('renders the Something Went Wrong component when the query errors', async () => {
      useMoveDetailsQueries.mockReturnValue(errorReturnValue);

      render(<ReviewDocuments {...requiredProps} />);

      const errorMessage = await screen.getByText(/Something went wrong./);
      expect(errorMessage).toBeInTheDocument();
    });
  });

  describe('with data loaded', () => {
    useMoveDetailsQueries.mockReturnValue(useMoveDetailsQueriesReturnValue);
    const wrapper = mount(<ReviewDocuments {...requiredProps} />);

    it('renders without errors', () => {
      expect(wrapper.find('[data-testid="ReviewDocuments"]').exists()).toBe(true);
    });

    it('renders the DocumentViewer', () => {
      const documentViewer = wrapper.find('DocumentViewer');
      expect(documentViewer.exists()).toBe(true);
      expect(documentViewer.prop('files')).toEqual([mockPDFUpload, mockXLSUpload]);
    });
  });
});
