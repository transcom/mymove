import React, { useMemo } from 'react';
import { GridContainer } from '@trussworks/react-uswds';
import { useTable, useFilters, usePagination, useSortBy } from 'react-table';
import PropTypes from 'prop-types';

import styles from './SearchResultsTable.module.scss';

import Table from 'components/Table/Table';
import TextBoxFilter from 'components/Table/Filters/TextBoxFilter';
import { SortShape } from 'constants/queues';

// SearchResultsTable is a react-table that uses react-hooks to fetch, filter, sort and page data
const SearchResultsTable = (props) => {
  const {
    title,
    columns,
    manualSortBy,
    manualFilters,
    disableMultiSort,
    defaultCanSort,
    disableSortBy,
    defaultSortedColumns,
    defaultHiddenColumns,
    handleClick,
    showFilters,
    showPagination,
    data,
  } = props;

  const totalCount = data.length;

  const defaultColumn = useMemo(
    () => ({
      // Let's set up our default Filter UI
      Filter: TextBoxFilter,
    }),
    [],
  );

  const tableData = useMemo(() => {
    return data;
  }, [data]);
  const {
    getTableProps,
    getTableBodyProps,
    headerGroups,
    rows,
    prepareRow,
    canPreviousPage,
    canNextPage,
    pageOptions,
    gotoPage,
    nextPage,
    previousPage,
    setPageSize,
    state: { pageIndex, pageSize },
  } = useTable(
    {
      columns,
      data: tableData,
      initialState: {
        hiddenColumns: defaultHiddenColumns,
        sortBy: defaultSortedColumns,
      },
      defaultColumn, // Be sure to pass the defaultColumn option
      manualFilters,
      manualPagination: false,
      manualSortBy,
      disableMultiSort,
      defaultCanSort,
      disableSortBy,
      autoResetSortBy: false,
    },
    useFilters,
    useSortBy,
    usePagination,
  );

  return (
    <GridContainer data-testid="table-search" containerSize="widescreen" className={styles.SearchResultsTable}>
      <h2>{`${title} (${totalCount})`}</h2>
      {totalCount > 0 ? (
        <div className={styles.tableContainer}>
          <Table
            showFilters={showFilters}
            showPagination={showPagination}
            handleClick={handleClick}
            gotoPage={gotoPage}
            setPageSize={setPageSize}
            nextPage={nextPage}
            previousPage={previousPage}
            getTableProps={getTableProps}
            getTableBodyProps={getTableBodyProps}
            headerGroups={headerGroups}
            rows={rows}
            prepareRow={prepareRow}
            canPreviousPage={canPreviousPage}
            canNextPage={canNextPage}
            pageIndex={pageIndex}
            pageSize={pageSize}
            pageOptions={pageOptions}
          />
        </div>
      ) : (
        <p>No results found</p>
      )}
    </GridContainer>
  );
};

// TODO use an actual shape here
const SearchResultsShape = PropTypes.array;

SearchResultsTable.propTypes = {
  // handleClick is the handler to handle functionality to click on a row
  handleClick: PropTypes.func.isRequired,
  // title is the table title
  title: PropTypes.string.isRequired,
  // columns is the columns to show in the table
  columns: PropTypes.arrayOf(PropTypes.object).isRequired,
  // showFilters is bool value to show filters or not
  showFilters: PropTypes.bool,
  // showPagination is bool value to show pagination or not
  showPagination: PropTypes.bool,
  // manualSortBy should be enabled if doing sorting on the server side
  manualSortBy: PropTypes.bool,
  // manualFilters should be enabled if doing filtering on the server side
  manualFilters: PropTypes.bool,
  // disableMultiSort turns off keyboard selecting multiple columns to sort by
  disableMultiSort: PropTypes.bool,
  // defaultCanSort determines if all columns are by default sortable
  defaultCanSort: PropTypes.bool,
  // disableSortBy is bool flag to turn off sorting functionality
  disableSortBy: PropTypes.bool,
  // defaultSortedColumns is an object of the column id and sort direction
  defaultSortedColumns: SortShape,
  // defaultHiddenColumns is an array of columns to hide
  defaultHiddenColumns: PropTypes.arrayOf(PropTypes.string),
  data: SearchResultsShape,
};

SearchResultsTable.defaultProps = {
  showFilters: false,
  showPagination: false,
  manualSortBy: false,
  manualFilters: true,
  disableMultiSort: false,
  defaultCanSort: false,
  disableSortBy: true,
  defaultSortedColumns: [],
  defaultHiddenColumns: ['id'],
  data: [],
};

export default SearchResultsTable;
