import React, { useState, useEffect, useMemo } from 'react';
import { GridContainer } from '@trussworks/react-uswds';
import { useTable, useFilters, usePagination, useSortBy } from 'react-table';
import PropTypes from 'prop-types';

import styles from './SearchResultsTable.module.scss';
import { createHeader } from './utils';

import Table from 'components/Table/Table';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import SomethingWentWrong from 'shared/SomethingWentWrong';
import TextBoxFilter from 'components/Table/Filters/TextBoxFilter';
import { BRANCH_OPTIONS, MOVE_STATUS_LABELS, MOVE_STATUS_OPTIONS, SortShape } from 'constants/queues';
import { serviceMemberAgencyLabel } from 'utils/formatters';
import MultiSelectCheckBoxFilter from 'components/Table/Filters/MultiSelectCheckBoxFilter';
import SelectFilter from 'components/Table/Filters/SelectFilter';

const columns = [
  createHeader('Move code', 'locator', {
    id: 'locator',
    isFilterable: false,
  }),
  createHeader('DOD ID', 'dodID', {
    id: 'dodID',
    isFilterable: false,
  }),
  createHeader(
    'Customer name',
    (row) => {
      return `${row.lastName}, ${row.firstName}`;
    },
    {
      id: 'customerName',
      isFilterable: false,
    },
  ),
  createHeader(
    'Status',
    (row) => {
      return MOVE_STATUS_LABELS[`${row.status}`];
    },
    {
      id: 'status',
      isFilterable: true,
      // eslint-disable-next-line react/jsx-props-no-spreading
      Filter: (props) => <MultiSelectCheckBoxFilter options={MOVE_STATUS_OPTIONS} {...props} />,
    },
  ),
  createHeader(
    'Origin ZIP',
    (row) => {
      return row.originDutyLocationPostalCode;
    },
    {
      id: 'originPostalCode',
      isFilterable: true,
    },
  ),
  createHeader(
    'Destination ZIP',
    (row) => {
      return row.destinationDutyLocationPostalCode;
    },
    {
      id: 'destinationPostalCode',
      isFilterable: true,
    },
  ),
  createHeader(
    'Branch',
    (row) => {
      return serviceMemberAgencyLabel(row.branch);
    },
    {
      id: 'branch',
      isFilterable: true,
      Filter: (props) => (
        // eslint-disable-next-line react/jsx-props-no-spreading
        <SelectFilter options={BRANCH_OPTIONS} {...props} />
      ),
    },
  ),
  createHeader(
    'Number of Shipments',
    (row) => {
      return Number(row.shipmentsCount || 0);
    },
    { id: 'shipmentsCount', isFilterable: true },
  ),
];

// SearchResultsTable is a react-table that uses react-hooks to fetch, filter, sort and page data
const SearchResultsTable = (props) => {
  const {
    title,
    disableMultiSort,
    defaultCanSort,
    disableSortBy,
    defaultSortedColumns,
    defaultHiddenColumns,
    handleClick,
    useQueries,
    showFilters,
    showPagination,
    dodID,
    moveCode,
    customerName,
  } = props;
  const [paramSort, setParamSort] = useState(defaultSortedColumns);
  const [paramFilters, setParamFilters] = useState([]);
  const [currentPage, setCurrentPage] = useState(1);
  const [currentPageSize, setCurrentPageSize] = useState(20);
  const [pageCount, setPageCount] = useState(0);

  const { id, desc } = paramSort.length ? paramSort[0] : {};

  let order;
  if (desc !== undefined) {
    order = desc ? 'desc' : 'asc';
  }

  const {
    searchResult: { totalCount = 0, data = [], page = 1, perPage = 20 },
    isLoading,
    isError,
  } = useQueries({
    sort: id,
    order,
    filters: paramFilters,
    currentPage,
    currentPageSize,
  });

  // react-table setup below
  const defaultColumn = useMemo(
    () => ({
      // Let's set up our default Filter UI
      Filter: TextBoxFilter,
    }),
    [],
  );
  const tableData = useMemo(() => data, [data]);
  const tableColumns = useMemo(() => columns, []);
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
      columns: tableColumns,
      data: tableData,
      initialState: {
        hiddenColumns: defaultHiddenColumns,
        pageSize: perPage,
        pageIndex: page - 1,
        sortBy: defaultSortedColumns,
      },
      defaultColumn, // Be sure to pass the defaultColumn option
      manualFilters: true,
      manualPagination: true,
      pageCount,
      manualSortBy: true,
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

  // When these table states change, fetch new data!
  useEffect(() => {
    if (!isLoading && !isError) {
      setParamSort(sortBy);
      setCurrentPage(pageIndex + 1);
      setCurrentPageSize(pageSize);
      setPageCount(Math.ceil(totalCount / pageSize));
    }
  }, [sortBy, pageIndex, pageSize, isLoading, isError, totalCount]);

  // Update filters when we get a new search or a column filter is edited
  useEffect(() => {
    const filtersToAdd = [];
    if (moveCode) {
      filtersToAdd.push({ id: 'locator', value: moveCode });
    }
    if (dodID) {
      filtersToAdd.push({ id: 'dodID', value: dodID });
    }
    if (customerName) {
      filtersToAdd.push({ id: 'customerName', value: customerName });
    }
    setParamFilters(filtersToAdd.concat(filters));
  }, [filters, moveCode, dodID, customerName]);

  if (isLoading) return <LoadingPlaceholder />;
  if (isError) return <SomethingWentWrong />;

  return (
    <GridContainer data-testid="table-queue" containerSize="widescreen" className={styles.SearchResultsTable}>
      <h1>{`${title} (${totalCount})`}</h1>
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
          pageCount={pageCount}
          pageOptions={pageOptions}
        />
      </div>
    </GridContainer>
  );
};

SearchResultsTable.propTypes = {
  // handleClick is the handler to handle functionality to click on a row
  handleClick: PropTypes.func.isRequired,
  // useQueries is the react-query hook call to handle data fetching
  useQueries: PropTypes.func.isRequired,
  // title is the table title
  title: PropTypes.string.isRequired,
  // showFilters is bool value to show filters or not
  showFilters: PropTypes.bool,
  // showPagination is bool value to show pagination or not
  showPagination: PropTypes.bool,
  // manualSortBy should be enabled if doing sorting on the server side
  disableMultiSort: PropTypes.bool,
  // defaultCanSort determines if all columns are by default sortable
  defaultCanSort: PropTypes.bool,
  // disableSortBy is bool flag to turn off sorting functionality
  disableSortBy: PropTypes.bool,
  // defaultSortedColumns is an object of the column id and sort direction
  defaultSortedColumns: SortShape,
  // defaultHiddenColumns is an array of columns to hide
  defaultHiddenColumns: PropTypes.arrayOf(PropTypes.string),
  // dodID is the DOD ID that is being searched for
  dodID: PropTypes.string,
  // moveCode is the move code that is being searched for
  moveCode: PropTypes.string,
  // customerName is the customer name search text
  customerName: PropTypes.string,
};

SearchResultsTable.defaultProps = {
  showFilters: false,
  showPagination: false,
  disableMultiSort: false,
  defaultCanSort: false,
  disableSortBy: true,
  defaultSortedColumns: [],
  defaultHiddenColumns: ['id'],
  dodID: null,
  moveCode: null,
  customerName: null,
};

export default SearchResultsTable;
