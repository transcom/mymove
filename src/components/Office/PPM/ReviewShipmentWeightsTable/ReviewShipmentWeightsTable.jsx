import React from 'react';
import { useFilters, usePagination, useSortBy, useTable } from 'react-table';

import styles from './ReviewShipmentWeightsTable.module.scss';
import { addShipmentNumbersToTableData } from './helpers';

import Table from 'components/Table/Table';

const ReviewShipmentWeightsTable = ({ tableData, tableConfig }) => {
  const { tableColumns, noRowsMsg, determineShipmentNumbers } = tableConfig;

  const reviewWeightsData = React.useMemo(
    () => addShipmentNumbersToTableData(tableData, determineShipmentNumbers),
    [tableData, determineShipmentNumbers],
  );
  const columns = React.useMemo(() => tableColumns, [tableColumns]);

  const { getTableProps, getTableBodyProps, headerGroups, rows, prepareRow } = useTable(
    {
      columns,
      data: reviewWeightsData,
      manualFilters: false,
      manualPagination: false,
      manualSortBy: false,
      disableMultiSort: true,
      defaultCanSort: false,
      disableSortBy: true,
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
    <div data-testid="reviewShipmentWeightsTable" className={styles.ReviewShipmentWeightsTable}>
      {rows.length > 0 ? (
        <div className={styles.tableContainer}>
          <Table
            getTableProps={getTableProps}
            getTableBodyProps={getTableBodyProps}
            headerGroups={headerGroups}
            rows={rows}
            prepareRow={prepareRow}
            handleClick={() => {}}
            tableStyle={styles.Table}
          />
        </div>
      ) : (
        <p>{noRowsMsg || 'No results found.'}</p>
      )}
    </div>
  );
};

export default ReviewShipmentWeightsTable;
