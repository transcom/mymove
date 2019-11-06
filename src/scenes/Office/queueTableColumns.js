import React from 'react';
import { capitalize } from 'lodash';
import { formatDate } from 'shared/formatters';

// Abstracting react table column creation
const CreateReactTableColumn = (header, accessor, options = {}) => ({
  Header: header,
  accessor: accessor,
  ...options,
});

const status = CreateReactTableColumn('Status', 'synthetic_status', {
  Cell: row => (
    <span className="status" data-cy="status">
      {capitalize(row.value && row.value.replace('_', ' '))}
    </span>
  ),
  filterable: false,
});

const customerName = CreateReactTableColumn('Customer name', 'customer_name', { filterable: false });

const dodId = CreateReactTableColumn('DoD ID', 'edipi', { filterable: false });

const locator = CreateReactTableColumn('Locator #', 'locator', {
  Cell: row => <span data-cy="locator">{row.value}</span>,
  filterable: false,
});

const moveDate = CreateReactTableColumn('PPM start', 'move_date', {
  Cell: row => <span className="move_date">{formatDate(row.value)}</span>,
  filterable: false,
});

const origin = CreateReactTableColumn('Origin', 'origin_duty_station_name', {
  Cell: row => <span>{row.value}</span>,
  filterable: false,
});

const destination = CreateReactTableColumn('Destination', 'destination_duty_station_name', {
  Cell: row => <span>{row.value}</span>,
  filterable: false,
});

const branchOfService = CreateReactTableColumn('Branch', 'branch_of_service', {
  Cell: row => <span>{row.value}</span>,
  filterMethod: (filter, row) => {
    if (filter.value === 'all') {
      return true;
    }

    return row[filter.id] === filter.value;
  },
  Filter: ({ filter, onChange }) => (
    <select
      onChange={event => onChange(event.target.value)}
      style={{ width: '100%' }}
      value={filter ? filter.value : 'all'}
    >
      <option value="all">Show All</option>
      <option value="ARMY">ARMY</option>
      <option value="NAVY">NAVY</option>
      <option value="MARINES">MARINES</option>
      <option value="AIR_FORCE">AIR_FORCE</option>
      <option value="COAST_GUARD">COAST_GUARD</option>
    </select>
  ),
});

// Columns used to display in react table
export const defaultColumns = [status, customerName, origin, destination, dodId, locator, moveDate, branchOfService];
