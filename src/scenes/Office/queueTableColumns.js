import React from 'react';
import { capitalize, memoize } from 'lodash';
import { formatDate } from 'shared/formatters';
import SingleDatePicker from 'shared/JsonSchemaForm/SingleDatePicker';
import moment from 'moment';

// Abstracting react table column creation
const CreateReactTableColumn = (header, accessor, options = {}) => ({
  Header: header,
  accessor: accessor,
  ...options,
});

// lodash memoize will prevent unnecessary rendering with the same state
// this will re-render if the state changes
const destination = memoize(destinationDutyStations =>
  CreateReactTableColumn('Destination', 'destination_duty_station_name', {
    Cell: row => <span>{row.value}</span>,
    Filter: ({ filter, onChange }) => (
      <select onChange={event => onChange(event.target.value)} value={filter ? filter.value : 'all'}>
        <option value="all">Show All</option>
        {destinationDutyStations.map(value => {
          return (
            <option key={value} value={value.toLowerCase()}>
              {value}
            </option>
          );
        })}
      </select>
    ),
    filterMethod: (filter, row) => {
      if (filter.value === 'all') {
        return true;
      } else if (row[filter.id] === undefined) {
        return false;
      }

      // filtered value should already be lowercase
      return row[filter.id].toLowerCase() === filter.value;
    },
    filterable: true,
  }),
);

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

const dateFormat = 'DD-MMM-YY';
const moveDate = CreateReactTableColumn('PPM start', 'move_date', {
  Cell: row => <span className="move_date">{formatDate(row.value)}</span>,
  Filter: ({ filter, onChange }) => {
    return (
      <div>
        <div>Before or on:</div>
        {SingleDatePicker({
          onChange: value => {
            return onChange(formatDate(value));
          },
          inputClassName: 'queue-date-picker-filter',
          value: filter ? filter.value : null,
          placeholder: dateFormat,
          format: dateFormat,
        })}
      </div>
    );
  },
  filterMethod: (filter, row) => {
    // Filter dates that are same or before the filtered value
    if (filter.value === undefined) {
      return true;
    } else if (row[filter.id] === undefined) {
      return false;
    }

    const rowDate = moment(row[filter.id]);
    const filterDate = moment(filter.value, dateFormat);

    return rowDate.isSameOrBefore(filterDate);
  },
  filterable: true,
});

const origin = CreateReactTableColumn('Origin', 'origin_duty_station_name', {
  Cell: row => <span>{row.value}</span>,
});

const branchOfService = CreateReactTableColumn('Branch', 'branch_of_service', {
  Cell: row => <span>{row.value}</span>,
  Filter: ({ filter, onChange }) => (
    <select onChange={event => onChange(event.target.value)} value={filter ? filter.value : 'all'}>
      <option value="all">Show All</option>
      <option value="ARMY">Army</option>
      <option value="NAVY">Navy</option>
      <option value="MARINES">Marines</option>
      <option value="AIR_FORCE">Air Force</option>
      <option value="COAST_GUARD">Coast Guard</option>
    </select>
  ),
  filterMethod: (filter, row) => {
    if (filter.value === 'all') {
      return true;
    }

    return row[filter.id] === filter.value;
  },
  filterable: true,
});

// Columns used to display in react table
export const defaultColumns = component => {
  return [
    status,
    customerName,
    origin,
    destination(component.getDestinationDutyStations()),
    dodId,
    locator,
    moveDate,
    branchOfService,
  ];
};
