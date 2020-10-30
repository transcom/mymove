import React, { useMemo } from 'react';
import { withKnobs } from '@storybook/addon-knobs';
import { useFilters, useTable } from 'react-table';

import { createHeader } from './utils';
import Table from './Table';

import TextBoxFilter from 'components/Table/Filters/TextBoxFilter';

export default {
  title: 'TOO/TIO Components|Table',
  decorators: [
    withKnobs,
    (storyFn) => (
      <div style={{ margin: '10px', height: '80vh', display: 'flex', flexDirection: 'column', overflow: 'auto' }}>
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
  createHeader('Status', 'col3', { isFilterable }),
  createHeader('Move Code', 'col4', { isFilterable }),
  createHeader('Branch', 'col5', { isFilterable }),
  createHeader('# of shipments', 'col6', { isFilterable }),
  createHeader('Destination duty station', 'col7', { isFilterable }),
  createHeader('Origin GBLOC', 'col8', { isFilterable }),
  createHeader('Last modified by', 'col9', { isFilterable }),
];

// eslint-disable-next-line react/prop-types
const CreatedTable = () => {
  const defaultColumn = useMemo(
    () => ({
      // Let's set up our default Filter UI
      Filter: TextBoxFilter,
    }),
    [],
  );

  const tableData = useMemo(() => data, []);
  const tableColumns = useMemo(() => columns(), []);
  const propsWithFilters = useTable(
    {
      columns: tableColumns,
      data: tableData,
      initialState: { hiddenColumns: ['id'] },
      manualFilters: true,
      defaultColumn,
    },
    useFilters,
  );

  // eslint-disable-next-line react/jsx-props-no-spreading
  return <Table {...propsWithFilters} />;
};

const CreatedTableWithFilters = () => {
  const defaultColumn = useMemo(
    () => ({
      // Let's set up our default Filter UI
      Filter: TextBoxFilter,
    }),
    [],
  );

  const tableData = useMemo(() => data, []);
  const tableColumns = useMemo(() => columns(true), []);
  const propsWithFilters = useTable(
    {
      columns: tableColumns,
      data: tableData,
      initialState: { hiddenColumns: ['id'] },
      manualFilters: true,
      defaultColumn,
    },
    useFilters,
  );

  // eslint-disable-next-line react/jsx-props-no-spreading
  return <Table {...propsWithFilters} />;
};

export const TXOTable = () => <CreatedTable />;

export const TXOTableFilters = () => <CreatedTableWithFilters />;
