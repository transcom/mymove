import React, { useState, useEffect, useMemo, useContext } from 'react';
import { connect } from 'react-redux';
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
import { selectLoggedInUser } from 'store/entities/selectors';
import SelectedGblocContext from 'components/Office/GblocSwitcher/SelectedGblocContext';
import {
  setTableQueueFilterSessionStorageValue,
  getTableQueueFilterSessionStorageValue,
  setTableQueuePageSizeSessionStorageValue,
  getTableQueuePageSizeSessionStorageValue,
  setTableQueuePageSessionStorageValue,
  getTableQueuePageSessionStorageValue,
  setTableQueueSortParamSessionStorageValue,
  getTableQueueSortParamSessionStorageValue,
  getSelectionOptionLabel,
} from 'components/Table/utils';
import { roleTypes } from 'constants/userRoles';

const defaultPageSize = 20;
const defaultPage = 1;

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
  sessionStorageKey,
  officeUser,
  activeRole,
}) => {
  const [isPageReload, setIsPageReload] = useState(true);
  useEffect(() => {
    // Component is mounted. Set flag to tell component
    // subsequent effects are post mount.
    setTimeout(() => {
      setIsPageReload(false);
    }, 500);
  }, []);

  const [paramSort, setParamSort] = useState(
    getTableQueueSortParamSessionStorageValue(sessionStorageKey) || defaultSortedColumns,
  );
  useEffect(() => {
    setTableQueueSortParamSessionStorageValue(sessionStorageKey, paramSort);
  }, [paramSort, sessionStorageKey]);

  // Pull table filters directly from cache. Updates are done in general table useEffect below.
  const paramFilters = getTableQueueFilterSessionStorageValue(sessionStorageKey) || [];

  const [currentPage, setCurrentPage] = useState(
    getTableQueuePageSessionStorageValue(sessionStorageKey) || defaultPage,
  );
  useEffect(() => {
    setTableQueuePageSessionStorageValue(sessionStorageKey, currentPage);
  }, [currentPage, sessionStorageKey]);

  const [currentPageSize, setCurrentPageSize] = useState(
    getTableQueuePageSizeSessionStorageValue(sessionStorageKey) || defaultPageSize,
  );
  useEffect(() => {
    setTableQueuePageSizeSessionStorageValue(sessionStorageKey, currentPageSize);
  }, [currentPageSize, sessionStorageKey]);

  const [pageCount, setPageCount] = useState(0);

  const { id, desc } = paramSort.length ? paramSort[0] : {};

  const gblocContext = useContext(SelectedGblocContext);
  const { selectedGbloc } =
    (activeRole === roleTypes.HQ || officeUser?.transportation_office_assignments?.length > 1) && gblocContext
      ? gblocContext
      : { selectedGbloc: undefined };

  const multiSelectValueDelimiter = ',';

  const {
    queueResult: {
      totalCount = 0,
      data = [],
      page = getTableQueuePageSessionStorageValue(sessionStorageKey) || defaultPage,
      perPage = getTableQueuePageSizeSessionStorageValue(sessionStorageKey) || defaultPageSize,
    },
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
    setAllFilters,
    state: { filters, pageIndex, pageSize, sortBy },
  } = useTable(
    {
      columns: tableColumns,
      data: tableData,
      initialState: {
        hiddenColumns: defaultHiddenColumns,
        pageSize: perPage,
        pageIndex: page - 1,
        sortBy: getTableQueueSortParamSessionStorageValue(sessionStorageKey) || defaultSortedColumns,
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

      if (filters.length === 0 && isPageReload) {
        // This is executed once. This is to ensure filters
        // is set with cached values during page reload use case.
        const filterStorage = getTableQueueFilterSessionStorageValue(sessionStorageKey) || [];
        filterStorage.forEach((item) => {
          // add cached filters to current prop filters var
          filters.push(item);
        });
      }

      // Save to cache.
      setTableQueueFilterSessionStorageValue(sessionStorageKey, filters);

      setCurrentPage(pageIndex + 1);
      setCurrentPageSize(pageSize);
      setPageCount(Math.ceil(totalCount / pageSize));
    }
  }, [sortBy, filters, pageIndex, pageSize, isLoading, isError, totalCount, isPageReload, sessionStorageKey]);

  if (isLoading || (title === 'Move history' && data.length <= 0 && !isError)) return <LoadingPlaceholder />;
  if (isError) return <SomethingWentWrong />;

  const isDateFilterValue = (value) => {
    return !Number.isNaN(Date.parse(value));
  };

  const handleRemoveFilterClick = (index) => {
    if (index === null) {
      filters.length = 0;
    } else {
      filters.splice(index, 1);
    }
    setAllFilters(filters);
  };

  const handleRemoveMultiSelectFilterClick = (index, valueToDelete) => {
    const filter = filters[index];
    const isObjectBasedArrayItem = Array.isArray(filter.value);
    const filterValues = !isObjectBasedArrayItem ? filter.value.split(multiSelectValueDelimiter) : filter.value;
    if (filterValues.length === 1) {
      filters.splice(index, 1);
    } else {
      const indexToDelete = filterValues.indexOf(valueToDelete);
      if (indexToDelete !== -1) {
        filterValues.splice(indexToDelete, 1);
      }
      filters[index].value = isObjectBasedArrayItem ? filterValues : filterValues.join(multiSelectValueDelimiter);
    }
    setAllFilters(filters);
  };

  const renderFilterPillButton = (index, value, buttonTitle, label, dataTestId) => {
    return (
      <button
        type="button"
        title={buttonTitle}
        data-testid={dataTestId}
        className={styles.pillButton}
        onClick={() => (value ? handleRemoveMultiSelectFilterClick(index, value) : handleRemoveFilterClick(index))}
      >
        {label} <span aria-hidden="true">&times;</span>
      </button>
    );
  };

  const renderRemoveAllPillButton = () => {
    let totalFilterValues = 0;
    // Loop through all filters to ensure there are really more than one filter values.
    // There is a chance filter.value that is object based array is empty. We can't totally
    // rely on filters.length.
    filters.forEach((filter) => {
      if (Array.isArray(filter.value)) {
        totalFilterValues += filter.value.length;
      } else {
        // legacy column filter control uses commas to represent array in one single string value
        totalFilterValues += filter.value.split(multiSelectValueDelimiter).length;
      }
    });
    if (totalFilterValues > 1) {
      return renderFilterPillButton(null, null, 'Remove all filters', 'All', 'remove-filters-all');
    }
    return null;
  };

  const renderFilterPillButtonList = () => {
    if (filters?.length > 0) {
      const filterPillButtons = [];
      const removeAllPillButton = renderRemoveAllPillButton();
      if (removeAllPillButton !== null) {
        filterPillButtons.push(removeAllPillButton);
      }
      const buttonTitle = 'Remove filter';
      const prefixDataTestId = 'remove-filters-';
      filters.forEach(function callback(filter, index) {
        columns.forEach((col) => {
          if (col.id === filter.id) {
            if ('Filter' in col) {
              if (isDateFilterValue(filter.value)) {
                filterPillButtons.push(
                  renderFilterPillButton(index, null, buttonTitle, col.Header, `${prefixDataTestId}${filter.id}`),
                );
              } else if (Array.isArray(filter.value)) {
                // value as real array
                filter.value.forEach((val) => {
                  const label = filter.value.length > 1 ? `${col.Header} (${val})` : col.Header;
                  filterPillButtons.push(
                    renderFilterPillButton(index, val, buttonTitle, label, `${prefixDataTestId}${filter.id}-${val}`),
                  );
                });
              } else {
                // value as string representing array using comma delimiter
                const values = filter.value.split(multiSelectValueDelimiter);
                values.forEach((val) => {
                  const label = values.length > 1 ? `${col.Header} (${getSelectionOptionLabel(val)})` : col.Header;
                  filterPillButtons.push(
                    renderFilterPillButton(index, val, buttonTitle, label, `${prefixDataTestId}${filter.id}-${val}`),
                  );
                });
              }
            } else {
              // default filter TextInput
              filterPillButtons.push(
                renderFilterPillButton(index, null, buttonTitle, col.Header, `${prefixDataTestId}${filter.id}`),
              );
            }
          }
        });
      });
      return <div className={styles.pillButtonRow}>Filters: {filterPillButtons}</div>;
    }
    return '';
  };

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
            isHeadquartersUser={activeRole === roleTypes.HQ}
          />
        )}
      </div>
      {renderFilterPillButtonList()}
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
  // session storage key to store search filters
  sessionStorageKey: PropTypes.string,
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
  sessionStorageKey: 'default',
};

const mapStateToProps = (state) => {
  const user = selectLoggedInUser(state);

  return {
    officeUser: user?.office_user || {},
    activeRole: state.auth.activeRole,
  };
};

export default connect(mapStateToProps)(TableQueue);
