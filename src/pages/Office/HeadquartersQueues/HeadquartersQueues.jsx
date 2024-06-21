import React, { useCallback, useEffect, useState } from 'react';
import { useNavigate, NavLink, useParams, Navigate } from 'react-router-dom';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';

import styles from './HeadquartersQueues.module.scss';

import { createHeader } from 'components/Table/utils';
import {
  useCustomerSearchQueries,
  useMovesQueueQueries,
  useMoveSearchQueries,
  usePaymentRequestQueueQueries,
  useServicesCounselingQueueQueries,
  useServicesCounselingQueuePPMQueries,
  useUserQueries,
} from 'hooks/queries';
import {
  getServicesCounselingQueue,
  getServicesCounselingPPMQueue,
  getPaymentRequestsQueue,
  getMovesQueue,
} from 'services/ghcApi';
import {
  formatAgeToDays,
  formatDateFromIso,
  paymentRequestStatusReadable,
  serviceMemberAgencyLabel,
} from 'utils/formatters';
import MultiSelectCheckBoxFilter from 'components/Table/Filters/MultiSelectCheckBoxFilter';
import SelectFilter from 'components/Table/Filters/SelectFilter';
import {
  BRANCH_OPTIONS,
  MOVE_STATUS_OPTIONS,
  GBLOC,
  MOVE_STATUS_LABELS,
  PAYMENT_REQUEST_STATUS_OPTIONS,
  SERVICE_COUNSELING_BRANCH_OPTIONS,
  SEARCH_QUEUE_STATUS_FILTER_OPTIONS,
  SERVICE_COUNSELING_MOVE_STATUS_LABELS,
  SERVICE_COUNSELING_PPM_TYPE_OPTIONS,
  SERVICE_COUNSELING_PPM_TYPE_LABELS,
} from 'constants/queues';
import TableQueue from 'components/Table/TableQueue';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import SomethingWentWrong from 'shared/SomethingWentWrong';
import DateSelectFilter from 'components/Table/Filters/DateSelectFilter';
import { DATE_FORMAT_STRING } from 'shared/constants';
import { CHECK_SPECIAL_ORDERS_TYPES, SPECIAL_ORDERS_TYPES } from 'constants/orders';
import MoveSearchForm from 'components/MoveSearchForm/MoveSearchForm';
import SearchResultsTable from 'components/Table/SearchResultsTable';
import TabNav from 'components/TabNav';
import { generalRoutes, hqRoutes } from 'constants/routes';
import { isNullUndefinedOrWhitespace } from 'shared/utils';
import NotFound from 'components/NotFound/NotFound';
import { isBooleanFlagEnabled } from 'utils/featureFlags';
import retryPageLoading from 'utils/retryPageLoading';
import { milmoveLogger } from 'utils/milmoveLog';
import ConnectedFlashMessage from 'containers/FlashMessage/FlashMessage';
import CustomerSearchForm from 'components/CustomerSearchForm/CustomerSearchForm';

const tooQueueColumns = (moveLockFlag, showBranchFilter = true) => [
  createHeader('ID', 'id', { id: 'id' }),
  createHeader(
    ' ',
    (row) => {
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
    },
    {
      id: 'lock',
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
          {`${row.customer.last_name}, ${row.customer.first_name}`}
        </div>
      );
    },
    {
      id: 'lastName',
      isFilterable: true,
      exportValue: (row) => {
        return `${row.customer.last_name}, ${row.customer.first_name}`;
      },
    },
  ),
  createHeader('DoD ID', 'customer.dodID', {
    id: 'dodID',
    isFilterable: true,
    exportValue: (row) => {
      return row.customer.dodID;
    },
  }),
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
  createHeader('Move code', 'locator', {
    id: 'locator',
    isFilterable: true,
  }),
  createHeader(
    'Requested move date',
    (row) => {
      return formatDateFromIso(row.requestedMoveDate, DATE_FORMAT_STRING);
    },
    {
      id: 'requestedMoveDate',
      isFilterable: true,
      // eslint-disable-next-line react/jsx-props-no-spreading
      Filter: (props) => <DateSelectFilter dateTime {...props} />,
    },
  ),
  createHeader(
    'Date submitted',
    (row) => {
      return formatDateFromIso(row.appearedInTooAt, DATE_FORMAT_STRING);
    },
    {
      id: 'appearedInTooAt',
      isFilterable: true,
      // eslint-disable-next-line react/jsx-props-no-spreading
      Filter: (props) => <DateSelectFilter dateTime {...props} />,
    },
  ),
  createHeader(
    'Branch',
    (row) => {
      return serviceMemberAgencyLabel(row.customer.agency);
    },
    {
      id: 'branch',
      isFilterable: showBranchFilter,
      Filter: (props) => (
        // eslint-disable-next-line react/jsx-props-no-spreading
        <SelectFilter options={BRANCH_OPTIONS} {...props} />
      ),
    },
  ),
  createHeader('# of shipments', 'shipmentsCount', { disableSortBy: true }),
  createHeader('Origin duty location', 'originDutyLocation.name', {
    id: 'originDutyLocation',
    isFilterable: true,
    exportValue: (row) => {
      return row.originDutyLocation?.name;
    },
  }),
  createHeader('Origin GBLOC', 'originGBLOC', { disableSortBy: true }),
];

const tioQueueColumns = (moveLockFlag, showBranchFilter = true) => [
  createHeader(
    ' ',
    (row) => {
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
    },
    {
      id: 'lock',
    },
  ),
  createHeader('ID', 'id', { id: 'id' }),
  createHeader(
    'Customer name',
    (row) => {
      return (
        <div>
          {CHECK_SPECIAL_ORDERS_TYPES(row.orderType) ? (
            <span className={styles.specialMoves}>{SPECIAL_ORDERS_TYPES[`${row.orderType}`]}</span>
          ) : null}
          {`${row.customer.last_name}, ${row.customer.first_name}`}
        </div>
      );
    },
    {
      id: 'lastName',
      isFilterable: true,
      exportValue: (row) => {
        return `${row.customer.last_name}, ${row.customer.first_name}`;
      },
    },
  ),
  createHeader('DoD ID', 'customer.dodID', {
    id: 'dodID',
    isFilterable: true,
    exportValue: (row) => {
      return row.customer.dodID;
    },
  }),
  createHeader(
    'Status',
    (row) => {
      return paymentRequestStatusReadable(row.status);
    },
    {
      id: 'status',
      isFilterable: true,
      Filter: (props) => (
        <div data-testid="statusFilter">
          {/* eslint-disable-next-line react/jsx-props-no-spreading */}
          <MultiSelectCheckBoxFilter options={PAYMENT_REQUEST_STATUS_OPTIONS} {...props} />
        </div>
      ),
    },
  ),
  createHeader(
    'Age',
    (row) => {
      return formatAgeToDays(row.age);
    },
    { id: 'age' },
  ),
  createHeader(
    'Submitted',
    (row) => {
      return formatDateFromIso(row.submittedAt, 'DD MMM YYYY');
    },
    {
      id: 'submittedAt',
      isFilterable: true,
      // eslint-disable-next-line react/jsx-props-no-spreading
      Filter: (props) => <DateSelectFilter dateTime {...props} />,
    },
  ),
  createHeader('Move Code', 'locator', {
    id: 'locator',
    isFilterable: true,
  }),
  createHeader(
    'Branch',
    (row) => {
      return serviceMemberAgencyLabel(row.customer.agency);
    },
    {
      id: 'branch',
      isFilterable: showBranchFilter,
      Filter: (props) => (
        // eslint-disable-next-line react/jsx-props-no-spreading
        <SelectFilter options={BRANCH_OPTIONS} {...props} />
      ),
    },
  ),
  createHeader('Origin GBLOC', 'originGBLOC', { disableSortBy: true }),
  createHeader('Origin Duty Location', 'originDutyLocation.name', {
    id: 'originDutyLocation',
    isFilterable: true,
    exportValue: (row) => {
      return row.originDutyLocation?.name;
    },
  }),
];

const counselingColumns = (moveLockFlag) => [
  createHeader(
    ' ',
    (row) => {
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
    },
    { id: 'lock' },
  ),
  createHeader('ID', 'id', { id: 'id' }),
  createHeader(
    'Customer name',
    (row) => {
      return (
        <div>
          {CHECK_SPECIAL_ORDERS_TYPES(row.orderType) ? (
            <span className={styles.specialMoves}>{SPECIAL_ORDERS_TYPES[`${row.orderType}`]}</span>
          ) : null}
          {`${row.customer.last_name}, ${row.customer.first_name}`}
        </div>
      );
    },
    {
      id: 'lastName',
      isFilterable: true,
      exportValue: (row) => {
        return `${row.customer.last_name}, ${row.customer.first_name}`;
      },
    },
  ),
  createHeader('DoD ID', 'customer.dodID', {
    id: 'dodID',
    isFilterable: true,
    exportValue: (row) => {
      return row.customer.dodID;
    },
  }),
  createHeader('Move code', 'locator', {
    id: 'locator',
    isFilterable: true,
  }),
  createHeader(
    'Status',
    (row) => {
      return SERVICE_COUNSELING_MOVE_STATUS_LABELS[`${row.status}`];
    },
    {
      id: 'status',
      isFilterable: true,
      Filter: (props) => (
        // eslint-disable-next-line react/jsx-props-no-spreading
        <MultiSelectCheckBoxFilter options={SEARCH_QUEUE_STATUS_FILTER_OPTIONS} {...props} />
      ),
    },
  ),
  createHeader(
    'Requested move date',
    (row) => {
      return formatDateFromIso(row.requestedMoveDate, DATE_FORMAT_STRING);
    },
    {
      id: 'requestedMoveDate',
      isFilterable: true,
      // eslint-disable-next-line react/jsx-props-no-spreading
      Filter: (props) => <DateSelectFilter {...props} />,
    },
  ),
  createHeader(
    'Date submitted',
    (row) => {
      return formatDateFromIso(row.submittedAt, DATE_FORMAT_STRING);
    },
    {
      id: 'submittedAt',
      isFilterable: true,
      // eslint-disable-next-line react/jsx-props-no-spreading
      Filter: (props) => <DateSelectFilter dateTime {...props} />,
    },
  ),
  createHeader(
    'Branch',
    (row) => {
      return serviceMemberAgencyLabel(row.customer.agency);
    },
    {
      id: 'branch',
      isFilterable: true,
      Filter: (props) => (
        // eslint-disable-next-line react/jsx-props-no-spreading
        <SelectFilter options={SERVICE_COUNSELING_BRANCH_OPTIONS} {...props} />
      ),
    },
  ),
  createHeader('Origin GBLOC', 'originGBLOC', {
    disableSortBy: true,
  }), // If the user is in the USMC GBLOC they will have many different GBLOCs and will want to sort and filter
  createHeader('Origin duty location', 'originDutyLocation.name', {
    id: 'originDutyLocation',
    isFilterable: true,
    exportValue: (row) => {
      return row.originDutyLocation?.name;
    },
  }),
];

const closeoutColumns = (moveLockFlag, ppmCloseoutGBLOC) => [
  createHeader(
    ' ',
    (row) => {
      const now = new Date();
      // this will render a lock icon if the move is locked & if the lockExpiresAt value is after right now
      if (row.lockedByOfficeUserID && row.lockExpiresAt && now < new Date(row.lockExpiresAt) && moveLockFlag) {
        return (
          <div id={row.id}>
            <FontAwesomeIcon icon="lock" />
          </div>
        );
      }
      return null;
    },
    { id: 'lock' },
  ),
  createHeader('ID', 'id', { id: 'id' }),
  createHeader(
    'Customer name',
    (row) => {
      return (
        <div>
          {CHECK_SPECIAL_ORDERS_TYPES(row.orderType) ? (
            <span className={styles.specialMoves}>{SPECIAL_ORDERS_TYPES[`${row.orderType}`]}</span>
          ) : null}
          {`${row.customer.last_name}, ${row.customer.first_name}`}
        </div>
      );
    },
    {
      id: 'lastName',
      isFilterable: true,
      exportValue: (row) => {
        return `${row.customer.last_name}, ${row.customer.first_name}`;
      },
    },
  ),
  createHeader('DoD ID', 'customer.dodID', {
    id: 'dodID',
    isFilterable: true,
    exportValue: (row) => {
      return row.customer.dodID;
    },
  }),
  createHeader('Move code', 'locator', {
    id: 'locator',
    isFilterable: true,
  }),
  createHeader(
    'Branch',
    (row) => {
      return serviceMemberAgencyLabel(row.customer.agency);
    },
    {
      id: 'branch',
      isFilterable: true,
      Filter: (props) => (
        // eslint-disable-next-line react/jsx-props-no-spreading
        <SelectFilter options={SERVICE_COUNSELING_BRANCH_OPTIONS} {...props} />
      ),
    },
  ),
  createHeader(
    'Closeout initiated',
    (row) => {
      return formatDateFromIso(row.closeoutInitiated, DATE_FORMAT_STRING);
    },
    {
      id: 'closeoutInitiated',
      isFilterable: true,
      // eslint-disable-next-line react/jsx-props-no-spreading
      Filter: (props) => <DateSelectFilter dateTime {...props} />,
    },
  ),
  createHeader(
    'Full or partial PPM',
    (row) => {
      return SERVICE_COUNSELING_PPM_TYPE_LABELS[`${row.ppmType}`];
    },
    {
      id: 'ppmType',
      isFilterable: true,
      Filter: (props) => (
        // eslint-disable-next-line react/jsx-props-no-spreading
        <SelectFilter options={SERVICE_COUNSELING_PPM_TYPE_OPTIONS} {...props} />
      ),
    },
  ),
  createHeader('Origin duty location', 'originDutyLocation.name', {
    id: 'originDutyLocation',
    isFilterable: true,
    exportValue: (row) => {
      return row.originDutyLocation?.name;
    },
  }),
  createHeader('Destination duty location', 'destinationDutyLocation.name', {
    id: 'destinationDutyLocation',
    isFilterable: true,
    exportValue: (row) => {
      return row.destinationDutyLocation?.name;
    },
  }),
  createHeader('PPM closeout location', 'closeoutLocation', {
    id: 'closeoutLocation',
    // This filter only makes sense if we're not in a closeout GBLOC. Users in a closeout GBLOC will
    // see the same value in this column for every move.
    isFilterable: !ppmCloseoutGBLOC,
  }),
];

const HeadquartersQueue = () => {
  const navigate = useNavigate();
  const { queueType } = useParams();
  const [search, setSearch] = useState({ moveCode: null, dodID: null, customerName: null });
  const [searchHappened, setSearchHappened] = useState(false);
  const [moveLockFlag, setMoveLockFlag] = useState(false);
  const [setErrorState] = useState({ hasError: false, error: undefined, info: undefined });

  useEffect(() => {
    const fetchData = async () => {
      const lockedMoveFlag = await isBooleanFlagEnabled('move_lock');
      setMoveLockFlag(lockedMoveFlag);
    };

    fetchData();
  }, []);

  // Feature Flag
  useEffect(() => {
    const fetchData = async () => {
      try {
        const lockedMoveFlag = await isBooleanFlagEnabled('move_lock');
        setMoveLockFlag(lockedMoveFlag);
      } catch (error) {
        const { message } = error;
        milmoveLogger.error({ message, info: null });
        setErrorState({
          hasError: true,
          error,
          info: null,
        });
        retryPageLoading(error);
      }
    };
    fetchData();
  }, [setErrorState]);

  const onSearchSubmit = useCallback((values) => {
    const payload = {
      moveCode: null,
      dodID: null,
      customerName: null,
    };
    if (!isNullUndefinedOrWhitespace(values.searchText)) {
      if (values.searchType === 'moveCode') {
        payload.moveCode = values.searchText.trim();
      } else if (values.searchType === 'dodID') {
        payload.dodID = values.searchText.trim();
      } else if (values.searchType === 'customerName') {
        payload.customerName = values.searchText.trim();
      }
    }
    setSearch(payload);
    setSearchHappened(true);
  }, []);

  const {
    // eslint-disable-next-line camelcase
    data: { office_user },
    isLoading,
    isError,
  } = useUserQueries();

  // eslint-disable-next-line camelcase
  const showBranchFilter = office_user?.transportation_office?.gbloc !== GBLOC.USMC;
  // eslint-disable-next-line camelcase
  const officeUserGBLOC = office_user?.transportation_office?.gbloc;
  const inPPMCloseoutGBLOC = officeUserGBLOC === 'TVCB' || officeUserGBLOC === 'NAVY' || officeUserGBLOC === 'USCG';

  const handleClickNavigateToDetails = (values) => {
    navigate(`/moves/${values.locator}/details`);
  };

  const handleClickNavigateToPaymentRequests = (values) => {
    navigate(`/moves/${values.locator}/payment-requests`);
  };

  if (isLoading) return <LoadingPlaceholder />;
  if (isError) return <SomethingWentWrong />;
  if (!queueType) {
    return <Navigate to={hqRoutes.BASE_MOVE_QUEUE} />;
  }
  const tabs = [
    <NavLink className={({ isActive }) => (isActive ? 'usa-current' : '')} to={hqRoutes.BASE_MOVE_QUEUE}>
      <span data-testid="task-order-queue-tab-link" className="tab-title" title="Move Queue">
        Task Order Queue
      </span>
    </NavLink>,
    <NavLink className={({ isActive }) => (isActive ? 'usa-current' : '')} to={hqRoutes.BASE_PAYMENT_REQUEST_QUEUE}>
      <span data-testid="payment-request-queue-tab-link" className="tab-title" title="Payment Request Queue">
        Payment Request Queue
      </span>
    </NavLink>,
    <NavLink end className={({ isActive }) => (isActive ? 'usa-current' : '')} to={hqRoutes.BASE_COUNSELING_QUEUE}>
      <span data-testid="counseling-queue-tab-link" className="tab-title">
        Counseling Queue
      </span>
    </NavLink>,
    <NavLink end className={({ isActive }) => (isActive ? 'usa-current' : '')} to={hqRoutes.BASE_CLOSEOUT_QUEUE}>
      <span data-testid="closeout-queue-tab-link" className="tab-title">
        PPM Closeout Queue
      </span>
    </NavLink>,
    <NavLink
      className={({ isActive }) => (isActive ? 'usa-current' : '')}
      to={generalRoutes.BASE_QUEUE_SEARCH_PATH}
      onClick={() => setSearchHappened(false)}
    >
      <span data-testid="move-search-tab-link" className="tab-title" title="Search">
        Move Search
      </span>
    </NavLink>,
    <NavLink
      end
      className={({ isActive }) => (isActive ? 'usa-current' : '')}
      to={hqRoutes.BASE_CUSTOMER_SEARCH}
      onClick={() => setSearchHappened(false)}
    >
      <span data-testid="customer-search-tab-link" className="tab-title">
        Customer Search
      </span>
    </NavLink>,
  ];

  const renderNavBar = () => {
    return <TabNav className={styles.tableTabs} items={tabs} />;
  };

  if (queueType === generalRoutes.QUEUE_SEARCH_PATH) {
    return (
      <div data-testid="move-search" className={styles.HeadquartersQueue}>
        {renderNavBar()}
        <h1>Search for a Move</h1>
        <MoveSearchForm onSubmit={onSearchSubmit} />
        {searchHappened && (
          <SearchResultsTable
            showFilters
            showPagination
            defaultCanSort
            disableMultiSort
            disableSortBy={false}
            title="Results"
            handleClick={handleClickNavigateToDetails}
            useQueries={useMoveSearchQueries}
            moveCode={search.moveCode}
            dodID={search.dodID}
            customerName={search.customerName}
            key="Move Search"
          />
        )}
      </div>
    );
  }
  if (queueType === hqRoutes.CUSTOMER_SEARCH) {
    return (
      <div data-testid="customer-search" className={styles.HeadquartersQueue}>
        {renderNavBar()}
        <ConnectedFlashMessage />
        <div className={styles.searchFormContainer}>
          <h1>Search for a Customer</h1>
        </div>
        <CustomerSearchForm onSubmit={onSearchSubmit} />
        {searchHappened && (
          <SearchResultsTable
            showFilters
            showPagination
            defaultCanSort
            disableMultiSort
            disableSortBy={false}
            title="Results"
            defaultHiddenColumns={['customerID', 'createMove']}
            handleClick={() => {}}
            useQueries={useCustomerSearchQueries}
            dodID={search.dodID}
            customerName={search.customerName}
            searchType="customer"
            key="Customer Search"
          />
        )}
      </div>
    );
  }
  if (queueType === hqRoutes.MOVE_QUEUE) {
    return (
      <div className={styles.HeadquartersQueue} data-testid="move-queue">
        {renderNavBar()}
        <TableQueue
          showFilters
          showPagination
          manualSortBy
          defaultCanSort
          defaultSortedColumns={[{ id: 'status', desc: false }]}
          disableMultiSort
          disableSortBy={false}
          columns={tooQueueColumns(moveLockFlag, showBranchFilter)}
          title="All moves"
          handleClick={handleClickNavigateToDetails}
          useQueries={useMovesQueueQueries}
          key="TOO Queue"
          showCSVExport
          csvExportFileNamePrefix="Task-Order-Queue"
          csvExportQueueFetcher={getMovesQueue}
          csvExportQueueFetcherKey="queueMoves"
        />
      </div>
    );
  }
  if (queueType === hqRoutes.PAYMENT_REQUEST_QUEUE) {
    return (
      <div className={styles.HeadquartersQueue} data-testid="payment-request-queue">
        {renderNavBar()}
        <TableQueue
          showFilters
          showPagination
          manualSortBy
          defaultCanSort
          defaultSortedColumns={[{ id: 'age', desc: true }]}
          disableMultiSort
          disableSortBy={false}
          columns={tioQueueColumns(moveLockFlag, showBranchFilter)}
          title="Payment requests"
          handleClick={handleClickNavigateToPaymentRequests}
          useQueries={usePaymentRequestQueueQueries}
          key="TIO Queue"
          showCSVExport
          csvExportFileNamePrefix="Payment-Request-Queue"
          csvExportQueueFetcher={getPaymentRequestsQueue}
          csvExportQueueFetcherKey="queuePaymentRequests"
        />
      </div>
    );
  }
  if (queueType === hqRoutes.CLOSEOUT_QUEUE) {
    return (
      <div className={styles.HeadquartersQueue}>
        {renderNavBar()}
        <TableQueue
          showFilters
          showPagination
          manualSortBy
          defaultCanSort
          defaultSortedColumns={[{ id: 'closeoutInitiated', desc: false }]}
          disableMultiSort
          disableSortBy={false}
          columns={closeoutColumns(moveLockFlag, inPPMCloseoutGBLOC)}
          title="Moves"
          handleClick={handleClickNavigateToDetails}
          useQueries={useServicesCounselingQueuePPMQueries}
          key="PPM Closeout Queue"
          showCSVExport
          csvExportFileNamePrefix="PPM-Closeout-Queue"
          csvExportQueueFetcher={getServicesCounselingPPMQueue}
          csvExportQueueFetcherKey="queueMoves"
        />
      </div>
    );
  }
  if (queueType === hqRoutes.COUNSELING_QUEUE) {
    return (
      <div className={styles.HeadquartersQueue}>
        {renderNavBar()}
        <TableQueue
          className={styles.ServicesCounseling}
          showFilters
          showPagination
          manualSortBy
          defaultCanSort
          defaultSortedColumns={[{ id: 'submittedAt', desc: false }]}
          disableMultiSort
          disableSortBy={false}
          columns={counselingColumns(moveLockFlag)}
          title="Moves"
          handleClick={handleClickNavigateToDetails}
          useQueries={useServicesCounselingQueueQueries}
          key="Counseling Queue"
          showCSVExport
          csvExportFileNamePrefix="Services-Counseling-Queue"
          csvExportQueueFetcher={getServicesCounselingQueue}
          csvExportQueueFetcherKey="queueMoves"
        />
      </div>
    );
  }
  return <NotFound />;
};

export default HeadquartersQueue;
