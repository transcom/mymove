import React from 'react';
import { capitalize } from 'lodash';
import { formatDate } from 'shared/formatters';
import SingleDatePicker from 'shared/JsonSchemaForm/SingleDatePicker';
import moment from 'moment';

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
});

const customerName = CreateReactTableColumn('Customer name', 'customer_name');

const dodId = CreateReactTableColumn('DoD ID', 'edipi');

const locator = CreateReactTableColumn('Locator #', 'locator', {
  Cell: row => <span data-cy="locator">{row.value}</span>,
});

const moveDate = CreateReactTableColumn('PPM start', 'move_date', {
  Cell: row => <span className="move_date">{formatDate(row.value)}</span>,
  filterable: true,
  filterMethod: (filter, row) => {
    // Filter dates that are same or before the filtered value
    if (filter.value === undefined) {
      return true;
    }

    const rowDate = moment(formatDate(row[filter.id]));
    const filterDate = moment(formatDate(filter.value));

    return rowDate.isSameOrBefore(filterDate);
  },
  Filter: ({ filter, onChange }) => {
    return SingleDatePicker({
      onChange: value => {
        return onChange(formatDate(value));
      },
      inputClassName: 'queue-date-picker-filter',
      value: filter ? filter.value : null,
      placeholder: 'DD-MMM-YY',
      format: 'DD-MMM-YY',
    });
  },
});

const origin = CreateReactTableColumn('Origin', 'origin_duty_station_name', {
  Cell: row => <span>{row.value}</span>,
});

const destination = CreateReactTableColumn('Destination', 'destination_duty_station_name', {
  Cell: row => <span>{row.value}</span>,
});

const branchOfService = CreateReactTableColumn('Branch', 'branch_of_service', {
  Cell: row => <span>{row.value}</span>,
  filterable: true,
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
      <option value="ARMY">Army</option>
      <option value="NAVY">Navy</option>
      <option value="MARINES">Marines</option>
      <option value="AIR_FORCE">Air Force</option>
      <option value="COAST_GUARD">Coast Guard</option>
    </select>
  ),
});

// Columns used to display in react table
export const defaultColumns = (props = {}) => {
  //const { queueType } = props;
  return [status, customerName, origin, destination, dodId, locator, moveDate, branchOfService];
};
