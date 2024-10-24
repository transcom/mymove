import React from 'react';
import { shallow } from 'enzyme';
import { render, screen } from '@testing-library/react';
import { MemoryRouter, useLocation } from 'react-router-dom';
import { QueryClientProvider, QueryClient } from '@tanstack/react-query';

import samplePDF from '../../../components/DocumentViewer/sample.pdf';
import sampleJPG from '../../../components/DocumentViewer/sample.jpg';
import samplePNG from '../../../components/DocumentViewer/sample2.png';
import sampleGIF from '../../../components/DocumentViewer/sample3.gif';

import MoveDocumentWrapper from './MoveDocumentWrapper';

import { useOrdersDocumentQueries, useAmendedDocumentQueries } from 'hooks/queries';

const mockFiles = [
  {
    id: 1,
    filename: 'Test File.pdf',
    contentType: 'application/pdf',
    url: samplePDF,
    createdAt: '2021-06-14T15:09:26.979879Z',
  },
  {
    id: 2,
    filename: 'Test File 2.jpg',
    contentType: 'image/jpeg',
    url: sampleJPG,
    createdAt: '2021-06-12T15:09:26.979879Z',
  },
  {
    id: 3,
    filename: 'Test File 3.png',
    contentType: 'image/png',
    url: samplePNG,
    createdAt: '2021-06-15T15:09:26.979879Z',
    rotation: 1,
  },
  {
    id: 4,
    filename: 'Test File 4.gif',
    contentType: 'image/gif',
    url: sampleGIF,
    createdAt: '2021-06-16T15:09:26.979879Z',
    rotation: 3,
  },
];

const mockOriginDutyLocation = {
  address: {
    city: 'Des Moines',
    country: 'US',
    eTag: 'MjAyMC0wOS0xNFQxNzo0MTozOC42OTg1OTha',
    id: '2e26b066-aaca-4563-b284-d7f3f978fb3c',
    postalCode: '50309',
    state: 'IA',
    streetAddress1: '987 Other Avenue',
    streetAddress2: 'P.O. Box 1234',
    streetAddress3: 'c/o Another Person',
  },
  address_id: '2e26b066-aaca-4563-b284-d7f3f978fb3c',
  eTag: 'MjAyMC0wOS0xNFQxNzo0MTozOC43MDcxOTVa',
  id: 'a3ec2bdd-aa0a-434a-ba58-34c85f047704',
  name: 'XBc1KNi3pA',
};

const mockDestinationDutyLocation = {
  address: {
    city: 'Augusta',
    country: 'United States',
    eTag: 'MjAyMC0wOS0xNFQxNzo0MDo0OC44OTM3MDVa',
    id: '5ac95be8-0230-47ea-90b4-b0f6f60de364',
    postalCode: '30813',
    state: 'GA',
    streetAddress1: 'Fort Gordon',
  },
  address_id: '5ac95be8-0230-47ea-90b4-b0f6f60de364',
  eTag: 'MjAyMC0wOS0xNFQxNzo0MDo0OC44OTM3MDVa',
  id: '2d5ada83-e09a-47f8-8de6-83ec51694a86',
  name: 'Fort Gordon',
};

jest.mock('hooks/queries', () => ({
  useOrdersDocumentQueries: jest.fn(),
  useAmendedDocumentQueries: jest.fn(),
}));

jest.mock('components/DocumentViewer/DocumentViewer', () => ({
  __esModule: true,
  default: ({ files, allowDownload }) => (
    <div>
      <div>
        <div data-testid="listOfFilesForViewer">
          {files ? JSON.stringify(files, null, 2).replace(/"/g, '') : 'No files available'}
        </div>
        <div data-testid="allowDownloadBool">Allow download: {allowDownload?.toString()}</div>
      </div>
    </div>
  ),
}));

const testMoveId = '10000';

const useOrdersDocumentQueriesReturnValue = {
  orders: {
    1: {
      agency: 'ARMY',
      customerID: '6ac40a00-e762-4f5f-b08d-3ea72a8e4b63',
      date_issued: '2018-03-15',
      department_indicator: 'AIR_AND_SPACE_FORCE',
      destinationDutyLocation: mockDestinationDutyLocation,
      eTag: 'MjAyMC0wOS0xNFQxNzo0MTozOC43MTE0Nlo=',
      entitlement: {
        authorizedWeight: 5000,
        dependentsAuthorized: true,
        eTag: 'MjAyMC0wOS0xNFQxNzo0MTozOC42ODAwOVo=',
        id: '0dbc9029-dfc5-4368-bc6b-dfc95f5fe317',
        nonTemporaryStorage: true,
        privatelyOwnedVehicle: true,
        proGearWeight: 2000,
        proGearWeightSpouse: 500,
        storageInTransit: 2,
        totalDependents: 1,
        totalWeight: 5000,
      },
      first_name: 'Leo',
      grade: 'E_1',
      id: '1',
      last_name: 'Spacemen',
      order_number: 'ORDER3',
      order_type: 'PERMANENT_CHANGE_OF_STATION',
      order_type_detail: 'HHG_PERMITTED',
      originDutyLocation: mockOriginDutyLocation,
      report_by_date: '2018-08-01',
      tac: 'F8E1',
      sac: 'E2P3',
    },
  },
  documents: {
    2: {
      id: '2',
      uploads: ['z'],
    },
  },
  upload: {
    z: mockFiles,
  },
  amendedDocuments: {
    3: {
      id: '3',
      uploads: ['x'],
    },
  },
  amendedUpload: {
    x: {
      id: 'z',
      filename: 'amended_test.pdf',
      contentType: 'application/pdf',
      url: '/storage/user/1/uploads/2?contentType=application%2Fpdf',
    },
  },
};
const useAmendedDocumentQueriesReturnValue = {
  amendedDocuments: {
    3: {
      id: '3',
      uploads: ['x'],
    },
  },
  amendedUpload: {
    x: {
      id: 'z',
      filename: 'amended_test.pdf',
      contentType: 'application/pdf',
      url: '/storage/user/1/uploads/2?contentType=application%2Fpdf',
    },
  },
};

const loadingReturnValue = {
  ...useOrdersDocumentQueriesReturnValue,
  isLoading: true,
  isError: false,
  isSuccess: false,
};

const errorReturnValue = {
  ...useOrdersDocumentQueriesReturnValue,
  isLoading: false,
  isError: true,
  isSuccess: false,
};

jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom'),
  useParams: jest.fn().mockReturnValue({ orderId: testMoveId }),
  useLocation: jest.fn(),
}));

describe('MoveDocumentWrapper', () => {
  describe('check loading and error component states', () => {
    it('renders the Loading Placeholder when the query is still loading', async () => {
      useLocation.mockReturnValue({ pathname: `/moves/${testMoveId}/orders` });
      useOrdersDocumentQueries.mockReturnValue(loadingReturnValue);
      useAmendedDocumentQueries.mockReturnValue(loadingReturnValue);

      render(<MoveDocumentWrapper />);

      const h2 = await screen.getByRole('heading', { name: 'Loading, please wait...', level: 2 });
      expect(h2).toBeInTheDocument();
    });

    it('renders the Something Went Wrong component when the query errors', async () => {
      useLocation.mockReturnValue({ pathname: `/moves/${testMoveId}/orders` });
      useOrdersDocumentQueries.mockReturnValue(errorReturnValue);
      useAmendedDocumentQueries.mockReturnValue(errorReturnValue);

      render(<MoveDocumentWrapper />);

      const errorMessage = await screen.getByText(/Something went wrong./);
      expect(errorMessage).toBeInTheDocument();
    });
  });

  describe('Basic rendering', () => {
    it('renders the document viewer', () => {
      useLocation.mockReturnValue({ pathname: `/moves/${testMoveId}/orders` });
      useOrdersDocumentQueries.mockReturnValue(useOrdersDocumentQueriesReturnValue);
      useAmendedDocumentQueries.mockReturnValue(useAmendedDocumentQueriesReturnValue);

      render(
        <MemoryRouter>
          <QueryClientProvider client={new QueryClient()}>
            <MoveDocumentWrapper allowDownload />
          </QueryClientProvider>
        </MemoryRouter>,
      );
      expect(screen.getByTestId('doc-wrapper')).toBeInTheDocument();
    });

    it('renders the sidebar Orders component', () => {
      useLocation.mockReturnValue({ pathname: `/moves/${testMoveId}/orders` });
      useOrdersDocumentQueries.mockReturnValue(useOrdersDocumentQueriesReturnValue);
      useAmendedDocumentQueries.mockReturnValue(useAmendedDocumentQueriesReturnValue);
      const wrapper = shallow(<MoveDocumentWrapper />);
      expect(wrapper.find('Orders').exists()).toBe(true);
    });

    it('renders the sidebar MoveAllowances component', () => {
      useLocation.mockReturnValue({ pathname: `/moves/${testMoveId}/allowances` });
      useOrdersDocumentQueries.mockReturnValue(useOrdersDocumentQueriesReturnValue);
      useAmendedDocumentQueries.mockReturnValue(useAmendedDocumentQueriesReturnValue);
      const wrapper = shallow(<MoveDocumentWrapper />);
      expect(wrapper.find('MoveAllowances').exists()).toBe(true);
    });

    it('combines orders and amended orders', () => {
      useLocation.mockReturnValue({ pathname: `/moves/${testMoveId}/orders` });
      useOrdersDocumentQueries.mockReturnValue(useOrdersDocumentQueriesReturnValue);
      useAmendedDocumentQueries.mockReturnValue(useAmendedDocumentQueriesReturnValue);

      render(
        <MemoryRouter>
          <QueryClientProvider client={new QueryClient()}>
            <MoveDocumentWrapper files={mockFiles} />
          </QueryClientProvider>
        </MemoryRouter>,
      );

      expect(screen.getByTestId('listOfFilesForViewer').textContent).toContain('Test File 3.png');
      expect(screen.getByTestId('listOfFilesForViewer').textContent).toContain('filename: Test File 4.gif');
      expect(screen.getByTestId('listOfFilesForViewer').textContent).toContain('filename: amended_test.pdf');
    });
  });
});
