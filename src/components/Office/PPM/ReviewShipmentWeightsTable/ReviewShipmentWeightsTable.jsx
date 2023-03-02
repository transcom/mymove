import React from 'react';
import { useTable, useFilters, usePagination, useSortBy } from 'react-table';

import styles from './ReviewShipmentWeightsTable.module.scss';

import Table from 'components/Table/Table';

const ReviewShipmentWeightsTable = (props) => {
  const { tableColumns, tableData, disableMultiSort, defaultCanSort, disableSortBy } = props;

  const { getTableProps, getTableBodyProps, headerGroups, rows, prepareRow } = useTable(
    {
      columns: tableColumns,
      data: tableData,
      manualFilters: false,
      manualPagination: false,
      manualSortBy: false,
      disableMultiSort,
      defaultCanSort,
      disableSortBy,
      autoResetSortBy: false,
      // If this option is true, the filters we get back from this hook
      // will not be memoized, which makes it easy to get into infinite render loops
      autoResetFilters: false,
    },
    useFilters,
    useSortBy,
    usePagination,
  );

  return (
    <div data-testid="table-queue" className={styles.ReviewShipmentWeightsTable}>
      {rows.length > 0 ? (
        <div className={styles.tableContainer}>
          <Table
            getTableProps={getTableProps}
            getTableBodyProps={getTableBodyProps}
            headerGroups={headerGroups}
            rows={rows}
            prepareRow={prepareRow}
            handleClick={() => {}}
          />
        </div>
      ) : (
        <p>No results found.</p>
      )}
    </div>
  );
};

export default ReviewShipmentWeightsTable;
