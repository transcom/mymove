import React from 'react';
import { capitalize, memoize } from 'lodash';
import { formatDate } from 'utils/formatters';
import SingleDatePicker from 'shared/JsonSchemaForm/SingleDatePicker';
import moment from 'moment';

// testing
import Select from 'react-select';

// Abstracting react table column creation
const createReactTableColumn = (header, accessor, options = {}) => ({
  Header: header,
  accessor: accessor,
  ...options,
});

const getReactSelectFilterSettings = (data = []) => ({
  Filter: ({ filter, onChange }) => {
    const options = data.map((value) => ({ label: value, value: value }));
    return (
      <Select
        options={options}
        onChange={(value) => {
          // value example: {label: "Fort Gordon", value: "Fort Gordon"}
          return onChange(value ? value : undefined);
        }}
        defaultValue={filter ? filter.value : undefined}
        styles={{
          // overriding styles to match other table filters
          control: (baseStyles) => ({
            ...baseStyles,
            height: '1.5rem',
            minHeight: '1.5rem',
            border: '1px solid rgba(0,0,0,0.1)',
          }),
          indicatorsContainer: (baseStyles) => ({
            ...baseStyles,
            height: '1.5rem',
          }),
          clearIndicator: (baseStyles) => ({
            ...baseStyles,
            padding: '0.2rem',
          }),
          dropdownIndicator: (baseStyles) => ({
            ...baseStyles,
            padding: '0.2rem',
          }),
          input: (baseStyles) => ({
            ...baseStyles,
            margin: '0 2px',
            paddingTop: '0',
            paddingBottom: '0',
          }),
          valueContainer: (baseStyles) => ({
            ...baseStyles,
            padding: '0 8px',
          }),
        }}
        isClearable
      />
    );
  },
  filterMethod: (filter, row) => {
    if (filter.value === undefined) {
      return true;
    } else if (row[filter.id] === undefined) {
      return false;
    }

    return row[filter.id].toLowerCase() === filter.value.value.toLowerCase();
  },
});

// lodash memoize will prevent unnecessary rendering with the same state
// this will re-render if the state changes
const destination = memoize((destinationDutyLocations) =>
  createReactTableColumn('Destination', 'destination_duty_location_name', {
    Cell: (row) => <span>{row.value}</span>,
    filterable: true,
    ...getReactSelectFilterSettings(destinationDutyLocations),
  }),
);

const origin = memoize((originDutyLocations) =>
  createReactTableColumn('Origin', 'origin_destination_duty_location_name_name', {
    Cell: (row) => <span>{row.value}</span>,
    filterable: true,
    ...getReactSelectFilterSettings(originDutyLocations),
  }),
);

const status = createReactTableColumn('Status', 'synthetic_status', {
  Cell: (row) => (
    <span className="status" data-testid="status">
      {capitalize(row.value && row.value.replace('_', ' '))}
    </span>
  ),
});

const customerName = createReactTableColumn('Customer name', 'customer_name');

const dodId = createReactTableColumn('DoD ID', 'edipi');

const locator = createReactTableColumn('Locator #', 'locator', {
  Cell: (row) => <span data-testid="locator">{row.value}</span>,
});

const dateFormat = 'DD-MMM-YY';
const moveDate = createReactTableColumn('PPM start', 'move_date', {
  Cell: (row) => <span className="move_date">{formatDate(row.value)}</span>,
  Filter: ({ filter, onChange }) => {
    return (
      <div>
        <div>Before or on:</div>
        {SingleDatePicker({
          onChange: (value) => {
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

const branchOfService = createReactTableColumn('Branch', 'branch_of_service', {
  Cell: (row) => <span>{row.value}</span>,
  Filter: ({ filter, onChange }) => (
    <select onChange={(event) => onChange(event.target.value)} value={filter ? filter.value : 'all'}>
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
export const defaultColumns = (component) => {
  return [
    status,
    customerName,
    origin(component.getOriginDutyLocations()),
    destination(component.getDestinationDutyLocations()),
    dodId,
    locator,
    moveDate,
    branchOfService,
  ];
};
