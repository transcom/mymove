import React from 'react';
import { capitalize } from 'lodash';
import { formatDate, formatDate4DigitYear } from 'shared/formatters';
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
      formattedValue: filter ? formatDate4DigitYear(filter.value) : null,
      placeholder: 'DD-MMM-YYYY',
    });
  },
});

const origin = CreateReactTableColumn('Origin', 'origin_duty_station_name', {
  Cell: row => <span>{row.value}</span>,
});

const destination = CreateReactTableColumn('Destination', 'destination_duty_station_name', {
  Cell: row => <span>{row.value}</span>,
});

// Columns used to display in react table
export const defaultColumns = [status, customerName, origin, destination, dodId, locator, moveDate];
