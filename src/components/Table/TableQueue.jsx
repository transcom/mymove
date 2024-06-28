import React, { useState, useEffect, useMemo } from 'react';
import { GridContainer } from '@trussworks/react-uswds';
import { useTable, useFilters, usePagination, useSortBy } from 'react-table';
import PropTypes from 'prop-types';

import styles from './TableQueue.module.scss';

import Table from 'components/Table/Table';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import SomethingWentWrong from 'shared/SomethingWentWrong';
import TextBoxFilter from 'components/Table/Filters/TextBoxFilter';
import { SortShape } from 'constants/queues';
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
  sessionStorageKey,
}) => {
  const [isPageReload, setIsPageReload] = useState(true);
  useEffect(() => {
    // Component is mounted. Set flag to tell component
    // subsequent effects are post mount.
    setTimeout(() => {
      setIsPageReload(false);
    }, 1000);
  }, []);

  const [paramSort, setParamSort] = useState(
    getTableQueueSortParamSessionStorageValue(sessionStorageKey) || defaultSortedColumns,
  );
  useEffect(() => {
    setTableQueueSortParamSessionStorageValue(sessionStorageKey, paramSort);
  }, [paramSort, sessionStorageKey]);

  // Pull table filters directly from cache. Updates are done in general table useEffect below.
  let paramFilters = getTableQueueFilterSessionStorageValue(sessionStorageKey) || [];

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

      if (filters.length === 0 && paramFilters.length > 0 && isPageReload) {
        // This is executed once. This is to ensure paramFilters
        // is set with cached values during page reload use case.
        paramFilters.forEach((item) => {
          // add cached filters to current prop filters var
          filters.push(item);
        });
      }

      // eslint-disable-next-line react-hooks/exhaustive-deps
      paramFilters = filters;

      // Save to cache.
      setTableQueueFilterSessionStorageValue(sessionStorageKey, paramFilters);

      setCurrentPage(pageIndex + 1);
      setCurrentPageSize(pageSize);
      setPageCount(Math.ceil(totalCount / pageSize));
    }
  }, [sortBy, filters, pageIndex, pageSize, isLoading, isError, totalCount]);

  if (isLoading || (title === 'Move history' && data.length <= 0 && !isError)) return <LoadingPlaceholder />;
  if (isError) return <SomethingWentWrong />;

  const isDateFilterParam = (filterParam) => {
    return !Number.isNaN(Date.parse(filterParam.value));
  };

  const handleRemoveFilterClick = (index) => {
    if (index === null) {
      paramFilters.length = 0;
    } else {
      paramFilters.splice(index, 1);
    }
    setAllFilters(paramFilters);
  };

  const handleRemoveMultiSelectFilterClick = (index, valueToDelete) => {
    const filter = paramFilters[index];
    const filterValues = filter.value.split(multiSelectValueDelimiter);
    if (filterValues.length === 1) {
      paramFilters.splice(index, 1);
    } else {
      const indexToDelete = filterValues.indexOf(valueToDelete);
      if (indexToDelete !== -1) {
        filterValues.splice(indexToDelete, 1);
      }
      paramFilters[index].value = filterValues.join(multiSelectValueDelimiter);
    }
    setAllFilters(paramFilters);
  };

  const renderFilterPillBtn = (index, value, useMultiSelectHandler, titleAttributeText, label, dataTestId) => {
    if (useMultiSelectHandler) {
      return (
        <button
          type="button"
          title={titleAttributeText}
          data-testid={dataTestId}
          className={styles.pillButton}
          onClick={() => handleRemoveMultiSelectFilterClick(index, value)}
        >
          {label} <span aria-hidden="true">&times;</span>
        </button>
      );
    }

    return (
      <button
        type="button"
        title={titleAttributeText}
        data-testid={dataTestId}
        className={styles.pillButton}
        onClick={() => handleRemoveFilterClick(index)}
      >
        {label} <span aria-hidden="true">&times;</span>
      </button>
    );
  };

  const renderRemoveAllPillBtn = () => {
    let isVisible = paramFilters?.length > 1;
    if (paramFilters?.length === 1) {
      if (!isDateFilterParam(paramFilters[0])) {
        isVisible = paramFilters[0].value.split(multiSelectValueDelimiter).length > 1;
      }
    }
    if (isVisible) {
      return renderFilterPillBtn(null, null, false, 'Remove all filters', 'All', 'remove-filters-all');
    }
    return null;
  };

  const renderFilterPillBtnList = () => {
    if (paramFilters?.length > 0) {
      const filterPillBtns = [];
      const removeAllPill = renderRemoveAllPillBtn();
      if (removeAllPill !== null) {
        filterPillBtns.push(removeAllPill);
      }
      // index, value, useMultiSelectHandler, titleAttributeText, label, dataTestId
      paramFilters.forEach(function callback(filter, index) {
        columns.forEach((col) => {
          if (col.id === filter.id) {
            if ('Filter' in col) {
              // MultiSelect filter column  can will contain Filter property.
              if (!isDateFilterParam(filter)) {
                const valueArray = filter.value.split(multiSelectValueDelimiter);
                if (valueArray.length > 1) {
                  // Multiselect filter type containing multiple filter values. Render
                  // pill button for each filter value item ..ex: Status (Approved request).
                  valueArray.forEach((val) => {
                    filterPillBtns.push(
                      renderFilterPillBtn(
                        index,
                        val,
                        true,
                        'Remove filter',
                        `${col.Header} (${getSelectionOptionLabel(val)})`,
                        `remove-filters-${filter.id}-${val}`,
                      ),
                    );
                  });
                } else {
                  // Multiselect filter type containing one value.
                  // In this case just display generic label using
                  // column header text.
                  filterPillBtns.push(
                    renderFilterPillBtn(
                      index,
                      null,
                      false,
                      'Remove filter',
                      `${col.Header}`,
                      `remove-filters-${filter.id}`,
                    ),
                  );
                }
              } else {
                // It's a datePicker filter.
                filterPillBtns.push(
                  renderFilterPillBtn(
                    index,
                    null,
                    false,
                    'Remove filter',
                    `${col.Header}`,
                    `remove-filters-${filter.id}`,
                  ),
                );
              }
            } else {
              // For filters column using default filter control.
              filterPillBtns.push(
                renderFilterPillBtn(
                  index,
                  null,
                  false,
                  'Remove filter',
                  `${col.Header}`,
                  `remove-filters-${filter.id}`,
                ),
              );
            }
          }
        });
      });
      return <div className={styles.pillButtonRow}>Filters: {filterPillBtns}</div>;
    }
    return '';
  };

  return (
    <GridContainer data-testid="table-queue" containerSize="widescreen" className={styles.TableQueue}>
      <h1>{`${title} (${totalCount})`}</h1>
      {renderFilterPillBtnList()}
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
  sessionStorageKey: 'default',
};
export default TableQueue;
