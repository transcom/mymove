import React, { useState, useEffect, useMemo, useContext } from 'react';
import { GridContainer } from '@trussworks/react-uswds';
import { useTable, useFilters, usePagination, useSortBy } from 'react-table';
import PropTypes from 'prop-types';

import styles from './TableQueue.module.scss';
import TableCSVExportButton from './TableCSVExportButton';

import Table from 'components/Table/Table';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import SomethingWentWrong from 'shared/SomethingWentWrong';
import TextBoxFilter from 'components/Table/Filters/TextBoxFilter';
import { SortShape } from 'constants/queues';
import SelectedGblocContext from 'components/Office/GblocSwitcher/SelectedGblocContext';

// TableQueue is a react-table that uses react-hooks to fetch, filter, sort and page data
const TableQueue = ({
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
  useQueries,
  showFilters,
  showPagination,
  showCSVExport,
  csvExportFileNamePrefix,
  csvExportHiddenColumns,
  csvExportQueueFetcher,
  csvExportQueueFetcherKey,
}) => {
  const [paramSort, setParamSort] = useState(defaultSortedColumns);
  const [paramFilters, setParamFilters] = useState([]);
  const [currentPage, setCurrentPage] = useState(1);
  const [currentPageSize, setCurrentPageSize] = useState(20);
  const [pageCount, setPageCount] = useState(0);

  const { id, desc } = paramSort.length ? paramSort[0] : {};

  const gblocContext = useContext(SelectedGblocContext);
  const { selectedGbloc } = gblocContext || { selectedGbloc: undefined };

  const {
    queueResult: { totalCount = 0, data = [], page = 1, perPage = 20 },
    isInitialLoading: isLoading,
    isError,
  } = useQueries({
    sort: id,
    order: desc ? 'desc' : 'asc',
    filters: paramFilters,
    currentPage,
    currentPageSize,
    viewAsGBLOC: selectedGbloc,
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
  const tableColumns = useMemo(() => columns, [columns]);
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
    if (!isLoading && !isError) {
      setParamSort(sortBy);
      setParamFilters(filters);
      setCurrentPage(pageIndex + 1);
      setCurrentPageSize(pageSize);
      setPageCount(Math.ceil(totalCount / pageSize));
    }
  }, [sortBy, filters, pageIndex, pageSize, isLoading, isError, totalCount]);

  if (isLoading || (title === 'Move history' && data.length <= 0 && !isError)) return <LoadingPlaceholder />;
  if (isError) return <SomethingWentWrong />;

  return (
    <GridContainer data-testid="table-queue" containerSize="widescreen" className={styles.TableQueue}>
      <div className={styles.queueHeader}>
        <h1>{`${title} (${totalCount})`}</h1>
        {showCSVExport && (
          <TableCSVExportButton
            className={styles.csvDownloadLink}
            tableColumns={columns}
            hiddenColumns={csvExportHiddenColumns}
            filePrefix={csvExportFileNamePrefix}
            queueFetcher={csvExportQueueFetcher}
            queueFetcherKey={csvExportQueueFetcherKey}
            totalCount={totalCount}
            paramSort={paramSort}
            paramFilters={paramFilters}
          />
        )}
      </div>
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

TableQueue.propTypes = {
  // handleClick is the handler to handle functionality to click on a row
  handleClick: PropTypes.func.isRequired,
  // useQueries is the react-query hook call to handle data fetching
  useQueries: PropTypes.func.isRequired,
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
  // showCSVExport shows the CSV export button
  showCSVExport: PropTypes.bool,
  // csvExportFileNamePrefix is the prefix used when this queue is exported to a CSV
  csvExportFileNamePrefix: PropTypes.string,
  // csvExportHiddenColumns is a array of the column ids to not use in a CSV export of the queue
  csvExportHiddenColumns: PropTypes.arrayOf(PropTypes.string),
  // csvExportQueueFetcher is the function to handle refetching non-paginated queue data
  csvExportQueueFetcher: PropTypes.func,
  // csvExportQueueFetcherKey is the key the queue data is stored under in the retrun value of csvExportQueueFetcher
  csvExportQueueFetcherKey: PropTypes.string,
};

TableQueue.defaultProps = {
  showFilters: false,
  showPagination: false,
  manualSortBy: false,
  manualFilters: true,
  disableMultiSort: false,
  defaultCanSort: false,
  disableSortBy: true,
  defaultSortedColumns: [],
  defaultHiddenColumns: ['id'],
  showCSVExport: false,
  csvExportFileNamePrefix: 'Moves',
  csvExportHiddenColumns: ['id', 'lock'],
  csvExportQueueFetcher: null,
  csvExportQueueFetcherKey: null,
};
export default TableQueue;
