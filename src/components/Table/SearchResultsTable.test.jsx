import React from 'react';
import { render, screen } from '@testing-library/react';

import SearchResultsTable from './SearchResultsTable';
import { createHeader } from './utils';
import MultiSelectCheckBoxFilter from './Filters/MultiSelectCheckBoxFilter';
import SelectFilter from './Filters/SelectFilter';

import { serviceMemberAgencyLabel } from 'utils/formatters';
import { BRANCH_OPTIONS, MOVE_STATUS_LABELS, MOVE_STATUS_OPTIONS } from 'constants/queues';

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

const columns = [
  createHeader('Move code', 'locator', {
    id: 'locator',
    isFilterable: false,
  }),
  createHeader('DOD ID', 'dodID', {
    id: 'dodID',
    isFilterable: false,
  }),
  createHeader(
    'Customer name',
    (row) => {
      return `${row.lastName}, ${row.firstName}`;
    },
    {
      id: 'customerName',
      isFilterable: false,
    },
  ),
  createHeader(
    'Status',
    (row) => {
      return MOVE_STATUS_LABELS[`${row.status}`];
    },
    {
      id: 'status',
      isFilterable: true,
      // eslint-disable-next-line react/jsx-props-no-spreading
      Filter: (props) => <MultiSelectCheckBoxFilter options={MOVE_STATUS_OPTIONS} {...props} />,
    },
  ),
  createHeader(
    'Origin ZIP',
    (row) => {
      return row.originDutyLocationPostalCode;
    },
    {
      id: 'originPostalCode',
      isFilterable: true,
    },
  ),
  createHeader(
    'Destination ZIP',
    (row) => {
      return row.destinationDutyLocationPostalCode;
    },
    {
      id: 'destinationPostalCode',
      isFilterable: true,
    },
  ),
  createHeader(
    'Branch',
    (row) => {
      return serviceMemberAgencyLabel(row.branch);
    },
    {
      id: 'branch',
      isFilterable: true,
      Filter: (props) => (
        // eslint-disable-next-line react/jsx-props-no-spreading
        <SelectFilter options={BRANCH_OPTIONS} {...props} />
      ),
    },
  ),
  createHeader(
    'Number of Shipments',
    (row) => {
      return Number(row.shipmentsCount);
    },
    { id: 'shipmentsCount', isFilterable: true },
  ),
];

describe('SearchResultsTable', () => {
  it('renders', () => {
    render(<SearchResultsTable handleClick={() => {}} title="Results" columns={columns} useQueries={mockQueries} />);
    const results = screen.queryByText('Results (1)');
    expect(results).toBeInTheDocument();
    const locator = screen.queryByText('P33YJB');
    expect(locator).toBeInTheDocument();
  });
});
