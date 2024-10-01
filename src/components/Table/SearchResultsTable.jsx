import React, { useState, useEffect, useMemo } from 'react';
import { useTable, useFilters, usePagination, useSortBy } from 'react-table';
import { generatePath, useNavigate } from 'react-router';
import PropTypes from 'prop-types';
import { Button } from '@trussworks/react-uswds';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';

import styles from './SearchResultsTable.module.scss';
import { createHeader } from './utils';

import Table from 'components/Table/Table';
import DateSelectFilter from 'components/Table/Filters/DateSelectFilter';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import SomethingWentWrong from 'shared/SomethingWentWrong';
import TextBoxFilter from 'components/Table/Filters/TextBoxFilter';
import { BRANCH_OPTIONS, MOVE_STATUS_LABELS, SEARCH_QUEUE_STATUS_FILTER_OPTIONS, SortShape } from 'constants/queues';
import { DATE_FORMAT_STRING } from 'shared/constants';
import { formatDateFromIso, serviceMemberAgencyLabel } from 'utils/formatters';
import MultiSelectCheckBoxFilter from 'components/Table/Filters/MultiSelectCheckBoxFilter';
import SelectFilter from 'components/Table/Filters/SelectFilter';
import { servicesCounselingRoutes } from 'constants/routes';
import { CHECK_SPECIAL_ORDERS_TYPES, SPECIAL_ORDERS_TYPES } from 'constants/orders';
import { isBooleanFlagEnabled } from 'utils/featureFlags';

const moveSearchColumns = (moveLockFlag, handleEditProfileClick) => [
  createHeader(' ', (row) => {
    const now = new Date();
    // this will render a lock icon if the move is locked & if the lockExpiresAt value is after right now
    if (row.lockedByOfficeUserID && row.lockExpiresAt && now < new Date(row.lockExpiresAt) && moveLockFlag) {
      return (
        <div data-testid="lock-icon">
          <FontAwesomeIcon icon="lock" />
        </div>
      );
    }
    return null;
  }),
  createHeader('Move code', 'locator', {
    id: 'locator',
    isFilterable: false,
  }),
  createHeader('DOD ID', 'dodID', {
    id: 'dodID',
    isFilterable: false,
  }),
  createHeader('EMPLID', 'emplid', {
    id: 'emplid',
    isFilterable: false,
  }),
  createHeader('  ', (row) => {
    return (
      <div className={styles.editProfile} data-label="editProfile" data-testid="editProfileBtn">
        <Button unstyled type="button" onClick={() => handleEditProfileClick(row.locator)}>
          <FontAwesomeIcon icon={['far', 'user']} />
        </Button>
      </div>
    );
  }),
  createHeader(
    'Customer name',
    (row) => {
      return (
        <div>
          {CHECK_SPECIAL_ORDERS_TYPES(row.orderType) ? (
            <span className={styles.specialMoves}>{SPECIAL_ORDERS_TYPES[`${row.orderType}`]}</span>
          ) : null}
          {`${row.lastName}, ${row.firstName}`}
        </div>
      );
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
      Filter: (props) => {
        return (
          <MultiSelectCheckBoxFilter
            options={SEARCH_QUEUE_STATUS_FILTER_OPTIONS}
            // eslint-disable-next-line react/jsx-props-no-spreading
            {...props}
          />
        );
      },
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
  createHeader(
    'Scheduled Pickup Date',
    (row) => {
      return formatDateFromIso(row.requestedPickupDate, DATE_FORMAT_STRING);
    },
    {
      id: 'pickupDate',
      disableSortBy: true,
      isFilterable: true,
      // eslint-disable-next-line react/jsx-props-no-spreading
      Filter: (props) => <DateSelectFilter dateTime {...props} />,
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
    'Origin GBLOC',
    (row) => {
      return row.originGBLOC;
    },
    {
      id: 'originGBLOC',
      disableSortBy: true,
    },
  ),
  createHeader(
    'Scheduled Delivery Date',
    (row) => {
      return formatDateFromIso(row.requestedDeliveryDate, DATE_FORMAT_STRING);
    },
    {
      id: 'deliveryDate',
      disableSortBy: true,
      isFilterable: true,
      // eslint-disable-next-line react/jsx-props-no-spreading
      Filter: (props) => <DateSelectFilter dateTime {...props} />,
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
    'Destination GBLOC',
    (row) => {
      return row.destinationGBLOC;
    },
    {
      id: 'destinationGBLOC',
      disableSortBy: true,
    },
  ),
];

const customerSearchColumns = () => [
  createHeader(
    'Create Move',
    (row) => {
      return (
        <Button
          onClick={() =>
            useNavigate(generatePath(servicesCounselingRoutes.BASE_MOVE_VIEW_PATH, { moveCode: row.locator }))
          }
          type="button"
          className={styles.createNewMove}
          data-testid="searchCreateMoveButton"
        >
          Create New Move
        </Button>
      );
    },
    { id: 'createMove', isFilterable: false, disableSortBy: true },
  ),
  createHeader(
    'id',
    (row) => {
      return row.id;
    },
    {
      id: 'customerID',
      isFilterable: false,
    },
  ),
  createHeader(
    'Customer name',
    (row) => {
      return (
        <div>
          {CHECK_SPECIAL_ORDERS_TYPES(row.orderType) ? (
            <span className={styles.specialMoves}>{SPECIAL_ORDERS_TYPES[`${row.orderType}`]}</span>
          ) : null}
          {`${row.lastName}, ${row.firstName}`}
        </div>
      );
    },
    {
      id: 'customerName',
      isFilterable: false,
    },
  ),
  createHeader(
    'Branch',
    (row) => {
      return serviceMemberAgencyLabel(row.branch);
    },
    {
      id: 'branch',
      isFilterable: false,
    },
  ),
  createHeader('DOD ID', 'dodID', {
    id: 'dodID',
    isFilterable: false,
  }),
  createHeader('EMPLID', 'emplid', {
    id: 'emplid',
    isFilterable: false,
  }),
  createHeader('Email', 'personalEmail', {
    id: 'personalEmail',
    isFilterable: false,
  }),
  createHeader('Phone', 'telephone', {
    id: 'telephone',
    isFilterable: false,
  }),
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
    handleEditProfileClick,
    useQueries,
    showFilters,
    showPagination,
    dodID,
    moveCode,
    customerName,
    paymentRequestCode,
    searchType,
  } = props;
  const [paramSort, setParamSort] = useState(defaultSortedColumns);
  const [paramFilters, setParamFilters] = useState([]);
  const [currentPage, setCurrentPage] = useState(1);
  const [currentPageSize, setCurrentPageSize] = useState(20);
  const [pageCount, setPageCount] = useState(0);
  const [moveLockFlag, setMoveLockFlag] = useState(false);

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
  const tableColumns = useMemo(() => {
    return searchType === 'customer'
      ? customerSearchColumns()
      : moveSearchColumns(moveLockFlag, handleEditProfileClick);
  }, [searchType, moveLockFlag, handleEditProfileClick]);

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
    setParamFilters([]);
    const filtersToAdd = [];
    if (moveCode) {
      filtersToAdd.push({ id: 'locator', value: moveCode.trim() });
    }
    if (dodID) {
      filtersToAdd.push({ id: 'dodID', value: dodID.trim() });
    }
    if (customerName) {
      filtersToAdd.push({ id: 'customerName', value: customerName });
    }
    if (paymentRequestCode) {
      filtersToAdd.push({ id: 'paymentRequestCode', value: paymentRequestCode });
    }
    setParamFilters(filtersToAdd.concat(filters));
  }, [filters, moveCode, dodID, customerName, paymentRequestCode]);

  // this useEffect handles the fetching of feature flags
  useEffect(() => {
    const fetchData = async () => {
      const lockedMoveFlag = await isBooleanFlagEnabled('move_lock');
      setMoveLockFlag(lockedMoveFlag);
    };

    fetchData();
  }, []);

  if (isLoading) return <LoadingPlaceholder />;
  if (isError) return <SomethingWentWrong />;

  return (
    <div data-testid="table-queue" className={styles.SearchResultsTable}>
      <h2>
        {`${title} (${totalCount})`} {totalCount > 0 ? null : <p>No results found.</p>}
      </h2>
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
    </div>
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
  paymentRequestCode: PropTypes.string,
  searchType: PropTypes.string,
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
  paymentRequestCode: null,
  searchType: 'move',
};

export default SearchResultsTable;
