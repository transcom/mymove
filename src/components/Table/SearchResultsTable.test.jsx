import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';

import SearchResultsTable from './SearchResultsTable';

import { isBooleanFlagEnabled } from 'utils/featureFlags';

jest.mock('utils/featureFlags', () => ({
  ...jest.requireActual('utils/featureFlags'),
  isBooleanFlagEnabled: jest.fn().mockImplementation(() => Promise.resolve()),
}));

const mockTableData = [
  {
    branch: 'ARMY',
    destinationDutyLocationPostalCode: '30813',
    dodID: '5177210523',
    firstName: 'Felicia',
    id: '630519ab-f0ee-40ea-8414-ed5524df0386',
    lastName: 'Arnold',
    locator: 'P33YJB',
    originDutyLocationPostalCode: '40475',
    shipmentsCount: 1,
    status: 'APPROVALS REQUESTED',
    requestedPickupDate: '2024-04-05',
    requestedDeliveryDate: '2024-04-10',
    originGBLOC: 'KKFA',
    destinationGBLOC: 'CNNQ',
    lockExpiresAt: '2099-10-15T23:48:35.420Z',
    lockedByOfficeUserID: '2744435d-7ba8-4cc5-bae5-f302c72c966e',
  },
];

const mockCustomerTableData = [
  {
    branch: 'MARINES',
    dodID: '6585626513',
    firstName: 'Ted',
    id: '8604447b-cbfc-4d59-a9a1-dec219eb2046',
    lastName: 'Marine',
    personalEmail: 'leo_spaceman_sm@example.com',
    telephone: '212-123-4567',
  },
];

function mockQueries() {
  return {
    searchResult: {
      data: mockTableData,
      totalCount: mockTableData.length,
    },
    isLoading: false,
    isError: false,
    isSuccess: true,
  };
}
function mockCustomerQueries() {
  return {
    searchResult: {
      data: mockCustomerTableData,
      totalCount: mockCustomerTableData.length,
    },
    isLoading: false,
    isError: false,
    isSuccess: true,
  };
}
function mockLoadingQuery() {
  return {
    searchResult: {
      data: [],
      totalCount: 0,
    },
    isLoading: true,
    isError: false,
    isSuccess: false,
  };
}
function mockErrorQuery() {
  return {
    searchResult: {
      data: [],
      totalCount: 0,
    },
    isLoading: false,
    isError: true,
    isSuccess: false,
  };
}

describe('SearchResultsTable', () => {
  afterEach(() => {
    jest.restoreAllMocks();
  });

  it('renders a move search', () => {
    render(<SearchResultsTable handleClick={() => {}} title="Results" useQueries={mockQueries} />);
    const results = screen.queryByText('Results (1)');
    expect(results).toBeInTheDocument();
    const locator = screen.queryByText('P33YJB');
    expect(locator).toBeInTheDocument();
    const pickupDate = screen.queryByText('05 Apr 2024');
    expect(pickupDate).toBeInTheDocument();
    const deliveryDate = screen.queryByText('10 Apr 2024');
    expect(deliveryDate).toBeInTheDocument();
    const originGBLOC = screen.queryByText('KKFA');
    expect(originGBLOC).toBeInTheDocument();
    const destinationGBLOC = screen.queryByText('CNNQ');
    expect(destinationGBLOC).toBeInTheDocument();
  });
  it('renders a customer search', () => {
    render(
      <SearchResultsTable
        handleClick={() => {}}
        title="Results"
        useQueries={mockCustomerQueries}
        searchType="customer"
      />,
    );
    const results = screen.queryByText('Results (1)');
    expect(results).toBeInTheDocument();
    const branch = screen.queryByText('Marine Corps');
    expect(branch).toBeInTheDocument();
    const dodID = screen.queryByText('6585626513');
    expect(dodID).toBeInTheDocument();
    const name = screen.queryByText('Marine, Ted');
    expect(name).toBeInTheDocument();
    const email = screen.queryByText('leo_spaceman_sm@example.com');
    expect(email).toBeInTheDocument();
    const phone = screen.queryByText('212-123-4567');
    expect(phone).toBeInTheDocument();
  });
  it('renders a lock icon when move lock flag is on', async () => {
    isBooleanFlagEnabled.mockResolvedValue(true);

    render(<SearchResultsTable handleClick={() => {}} title="Results" useQueries={mockQueries} searchType="move" />);

    await waitFor(() => {
      const lockIcon = screen.queryByTestId('lock-icon');
      expect(lockIcon).toBeInTheDocument();
    });
  });
  it('does NOT render a lock icon when move lock flag is off', async () => {
    isBooleanFlagEnabled.mockResolvedValue(false);

    render(<SearchResultsTable handleClick={() => {}} title="Results" useQueries={mockQueries} searchType="move" />);

    await waitFor(() => {
      const lockIcon = screen.queryByTestId('lock-icon');
      expect(lockIcon).not.toBeInTheDocument();
    });
  });
  it('renders create move button on customer search', () => {
    render(
      <SearchResultsTable handleClick={() => {}} title="Results" useQueries={mockQueries} searchType="customer" />,
    );

    const createMoveButton = screen.queryByTestId('searchCreateMoveButton');
    expect(createMoveButton).toBeInTheDocument();
  });
  it('does not render create move button on move search', () => {
    render(<SearchResultsTable handleClick={() => {}} title="Results" useQueries={mockQueries} searchType="move" />);

    const createMoveButton = screen.queryByTestId('searchCreateMoveButton');
    expect(createMoveButton).not.toBeInTheDocument();
  });
  it('loading', () => {
    render(
      <SearchResultsTable handleClick={() => {}} title="Results" useQueries={mockLoadingQuery} dodID="1234567890" />,
    );
    expect(screen.getByRole('heading', { name: 'Loading, please wait...', level: 2 })).toBeInTheDocument();
  });
  it('error', () => {
    render(
      <SearchResultsTable
        handleClick={() => {}}
        title="Results"
        useQueries={mockErrorQuery}
        customerName="leo spacemen"
      />,
    );
    expect(screen.getByRole('heading', { name: /something went wrong/i, level: 2 })).toBeInTheDocument();
  });
});
