import React from 'react';
import { createHeader } from 'components/Table/utils';
import SelectFilter from 'components/Table/Filters/SelectFilter';
import DateSelectFilter from 'components/Table/Filters/DateSelectFilter';
import { BRANCH_OPTIONS } from 'constants/queues';
import { capitalize, memoize } from 'lodash';
import { formatDate } from 'shared/formatters';
// import SingleDatePicker from 'shared/JsonSchemaForm/SingleDatePicker';
import moment from 'moment';

// testing
import Select from 'react-select';

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
const destination = memoize((destinationDutyStations) =>
  createHeader('Destination', 'destination_duty_station_name', {
    Cell: (row) => <span>{row.value}</span>,
    isFilterable: true,
    ...getReactSelectFilterSettings(destinationDutyStations),
  }),
);

const origin = memoize((originDutyStations) =>
  createHeader('Origin', 'origin_duty_station_name', {
    Cell: (row) => <span>{row.value}</span>,
    isFilterable: true,
    ...getReactSelectFilterSettings(originDutyStations),
  }),
);

const status = createHeader('Status', 'synthetic_status', {
  Cell: (row) => (
    <span className="status" data-testid="status">
      {capitalize(row.value && row.value.replace('_', ' '))}
    </span>
  ),
});

const customerName = createHeader('Customer name', 'customer_name');

const dodId = createHeader('DoD ID', 'edipi');

const locator = createHeader('Locator #', 'locator', {
  Cell: (row) => <span data-testid="locator">{row.value}</span>,
});

const dateFormat = 'DD-MMM-YY';
const moveDate = createHeader('PPM start', 'move_date', {
  Cell: (row) => <span className="move_date">{formatDate(row.value)}</span>,
  Filter: DateSelectFilter,
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
  isFilterable: true,
});

const branchOfService = createHeader('Branch', 'branch_of_service', {
  Cell: (row) => <span>{row.value}</span>,
  Filter: (props) => (
    // eslint-disable-next-line react/jsx-props-no-spreading
    <SelectFilter options={BRANCH_OPTIONS} {...props} />
  ),
  filterMethod: (filter, row) => {
    if (filter.value === 'all') {
      return true;
    }

    return row[filter.id] === filter.value;
  },
  isFilterable: true,
});

// Columns used to display in react table
export const defaultColumns = (component) => {
  return [
    status,
    customerName,
    origin(component.getOriginDutyStations()),
    destination(component.getDestinationDutyStations()),
    dodId,
    locator,
    moveDate,
    branchOfService,
  ];
};
