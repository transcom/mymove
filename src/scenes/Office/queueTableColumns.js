import React from 'react';
import { createHeader } from 'components/Table/utils';
import SelectFilter from 'components/Table/Filters/SelectFilter';
import DateSelectFilter from 'components/Table/Filters/DateSelectFilter';
import MultiSelectCheckBoxFilter from 'components/Table/Filters/MultiSelectCheckBoxFilter';
import { BRANCH_OPTIONS } from 'constants/queues';
import { capitalize, memoize } from 'lodash';
import { formatDate } from 'shared/formatters';
import moment from 'moment';

const getReactSelectFilterSettings = (data = []) => ({
  Filter: ({ filter }) => {
    const options = data.map((value) => ({ label: value, value: value }));
    return (
      <MultiSelectCheckBoxFilter
        options={options}
        column={{
          filterValue: filter ? filter.value : undefined,
          setFilter: (value) => {
            return value ? value : undefined;
          },
        }}
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
  filter: (filter, row) => {
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
  filter: (rows, id, filterValue) => {
    if (filterValue === 'all') {
      return true;
    }

    return rows[`${id}`] === filterValue;
  },
  isFilterable: true,
});

// Columns used to display in react table
export const defaultColumns = (origDutyStationData, desDutyStationData) => {
  return [
    status,
    customerName,
    origin(origDutyStationData),
    destination(desDutyStationData),
    dodId,
    locator,
    moveDate,
    branchOfService,
  ];
};
