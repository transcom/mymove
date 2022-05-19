import React, { useState, useEffect, useMemo } from 'react';
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
    searchKey,
    searchValue,
    data,
  } = props;
  // const [paramSort, setParamSort] = useState(defaultSortedColumns);
  // const [paramFilters, setParamFilters] = useState([]);
  // const [currentPage, setCurrentPage] = useState(1);
  // const [currentPageSize, setCurrentPageSize] = useState(20);
  const [pageCount, setPageCount] = useState(0);

  // const { id, desc } = paramSort.length ? paramSort[0] : {};

  // const { totalCount = 0, data = [], page = 1, perPage = 20 } = searchResult;
  const totalCount = data.length;
  const page = 1;
  const perPage = 20;

  const defaultColumn = useMemo(
    () => ({
      // Let's set up our default Filter UI
      Filter: TextBoxFilter,
    }),
    [],
  );
  // console.log(searchResult);
  // console.log('Data', data, isLoading, isError);
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
    state: { filters, pageIndex, pageSize, sortBy },
  } = useTable(
    {
      columns,
      data: tableData,
      initialState: {
        hiddenColumns: defaultHiddenColumns,
        pageSize: perPage,
        pageIndex: page - 1,
        sortBy: defaultSortedColumns,
      },
      defaultColumn, // Be sure to pass the defaultColumn option
      manualFilters,
      manualPagination: true,
      pageCount,
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

  // When these table states change, fetch new data!
  useEffect(() => {
    // console.log('useEffect', isLoading, isError);
    // if (!isLoading && !isError) {
    // console.log('useEffect not skipped');
    // setParamSort(sortBy);
    // setParamFilters(filters);
    // setCurrentPage(pageIndex + 1);
    // setCurrentPageSize(pageSize);
    setPageCount(Math.ceil(totalCount / pageSize));
    // }
  }, [sortBy, filters, pageIndex, pageSize, totalCount, searchKey, searchValue]);
  // }, [isLoading, isError, totalCount, searchKey, searchValue]);

  return (
    <GridContainer data-testid="table-search" containerSize="widescreen" className={styles.SearchResultsTable}>
      <h2>{`${title} (${totalCount})`}</h2>
      {totalCount > 0 ? (
        <div className={styles.tableContainer}>
          <p>{JSON.stringify(data)}</p>
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
            pageCount={pageCount}
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
  // TODO would need to handle move code and dod id
  searchKey: PropTypes.string,
  searchValue: PropTypes.string,
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
  searchKey: 'moveCode',
  searchValue: '',
  data: {},
};

export default SearchResultsTable;
