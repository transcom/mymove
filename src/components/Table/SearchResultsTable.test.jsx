import React from 'react';
import { render, screen } from '@testing-library/react';

import { BRANCH_OPTIONS, MOVE_STATUS_LABELS, MOVE_STATUS_OPTIONS } from '../../constants/queues';
import { serviceMemberAgencyLabel } from '../../utils/formatters';

import SearchResultsTable from './SearchResultsTable';
import { createHeader } from './utils';
import MultiSelectCheckBoxFilter from './Filters/MultiSelectCheckBoxFilter';
import SelectFilter from './Filters/SelectFilter';

const mockTableData = [
  {
    customer: {
      agency: 'ARMY',
      dodID: '5177210523',
      first_name: 'Felicia',
      last_name: 'Arnold',
    },
    departmentIndicator: 'ARMY',
    destinationDutyLocation: {
      address: {
        postalCode: '30813',
      },
      name: 'Fort Gordon',
    },
    locator: 'P33YJB',
    originDutyLocation: {
      address: {
        postalCode: '40475',
      },
      name: 'Blue Grass Army Depot',
    },
    requestedMoveDate: '2022-04-27',
    shipmentsCount: 1,
    status: 'APPROVALS REQUESTED',
  },
];

const columns = [
  createHeader('Move code', 'locator', {
    id: 'locator',
    isFilterable: true,
  }),
  createHeader('DOD ID', 'customer.dodID', {
    id: 'dodID',
    isFilterable: true,
  }),
  createHeader(
    'Customer name',
    (row) => {
      return `${row.customer.last_name}, ${row.customer.first_name}`;
    },
    {
      id: 'lastName',
      isFilterable: true,
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
      return `${row.originDutyLocation.address.postalCode}`;
    },
    {
      id: 'originZIP',
      isFilterable: true,
    },
  ),
  createHeader(
    'Destination ZIP',
    (row) => {
      return `${row.destinationDutyLocation.address.postalCode}`;
    },
    {
      id: 'destinationZIP',
      isFilterable: true,
    },
  ),
  createHeader(
    'Branch',
    (row) => {
      return serviceMemberAgencyLabel(row.customer.agency);
    },
    {
      id: 'branch',
      isFilterable: false,
      Filter: (props) => (
        // eslint-disable-next-line react/jsx-props-no-spreading
        <SelectFilter options={BRANCH_OPTIONS} {...props} />
      ),
    },
  ),
  createHeader('Number of shipments', 'shipmentsCount', { disableSortBy: true }),
];

describe('SearchResultsTable', () => {
  it('renders', () => {
    render(<SearchResultsTable handleClick={() => {}} title="Results" columns={columns} data={mockTableData} />);
    const results = screen.queryByText('Results (1)');
    expect(results).toBeInTheDocument();
    const locator = screen.queryByText('P33YJB');
    expect(locator).toBeInTheDocument();
  });
});
