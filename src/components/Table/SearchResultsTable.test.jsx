import React from 'react';
import { render, screen } from '@testing-library/react';

import SearchResultsTable from './SearchResultsTable';

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

describe('SearchResultsTable', () => {
  it('renders', () => {
    render(<SearchResultsTable handleClick={() => {}} title="Results" useQueries={mockQueries} />);
    const results = screen.queryByText('Results (1)');
    expect(results).toBeInTheDocument();
    const locator = screen.queryByText('P33YJB');
    expect(locator).toBeInTheDocument();
  });
});
