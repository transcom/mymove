/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';

import { createHeader } from './utils';
import TableQueue from './TableQueue';

import MultiSelectCheckBoxFilter from 'components/Table/Filters/MultiSelectCheckBoxFilter';
import { BRANCH_OPTIONS, MOVE_STATUS_OPTIONS } from 'constants/queues';
import SelectFilter from 'components/Table/Filters/SelectFilter';
import DateSelectFilter from 'components/Table/Filters/DateSelectFilter';

export default {
  title: 'Office Components/Table',
  decorators: [
    (storyFn) => (
      <div style={{ margin: '10px', height: '100vh', display: 'flex', flexDirection: 'column', overflow: 'auto' }}>
        {storyFn()}
      </div>
    ),
  ],
};

const data = [
  {
    col1: 'Banks, Aaliyah',
    col2: '987654321',
    col3: 'New move',
    col4: 'LCKMAJ',
    col5: 'Navy',
    col6: '3',
    col7: 'NAS Jacksonville',
    col8: 'HAFC',
    col9: 'Garimundi, J (SW)',
  },
  {
    col1: 'Childers, Jamie',
    col2: '987654321',
    col3: 'New move',
    col4: 'XCQ5ZH',
    col5: 'Navy',
    col6: '3',
    col7: 'NAS Jacksonville',
    col8: 'HAFC',
    col9: 'Garimundi, J (SW)',
  },
  {
    col1: 'Clark-Nunez, Sofia',
    col2: '987654321',
    col3: 'New move',
    col4: 'UCAF8Q',
    col5: 'Navy',
    col6: '3',
    col7: 'NAS Jacksonville',
    col8: 'HAFC',
    col9: 'Garimundi, J (SW)',
  },
];

const columns = (isFilterable = false) => [
  createHeader('Customer name', 'col1', { isFilterable }),
  createHeader('DoD ID', 'col2', { isFilterable }),
  createHeader('Status', 'col3', {
    isFilterable,
    Filter: (props) => <MultiSelectCheckBoxFilter options={MOVE_STATUS_OPTIONS} {...props} />,
  }),
  createHeader('Move Code', 'col4', { isFilterable }),
  createHeader('Branch', 'col5', {
    isFilterable,
    Filter: (props) => <SelectFilter options={BRANCH_OPTIONS} {...props} />,
  }),
  createHeader('# of shipments', 'col6', { isFilterable }),
  createHeader('Destination duty location', 'col7', { isFilterable }),
  createHeader('Origin GBLOC', 'col8', { isFilterable }),
  createHeader('Last modified by', 'col9', { isFilterable, Filter: DateSelectFilter }),
];

const defaultProps = {
  title: 'Table queue',
  useQueries: () => ({ queueResult: { data, totalCount: data.length, perPage: 1 } }),
  handleClick: () => {},
  columns: columns(),
};

export const TXOTable = () => (
  <div className="officeApp">
    <TableQueue {...defaultProps} />
  </div>
);

export const TXOTableSortable = () => (
  <div className="officeApp">
    <TableQueue {...defaultProps} disableSortBy={false} defaultSortedColumns={[{ id: 'col1', desc: false }]} />
  </div>
);

export const TXOTableFilters = () => (
  <div className="officeApp">
    <TableQueue {...defaultProps} columns={columns(true)} showFilters />
  </div>
);

export const TXOTablePagination = () => (
  <div className="officeApp">
    {' '}
    <TableQueue {...defaultProps} showPagination />
  </div>
);
