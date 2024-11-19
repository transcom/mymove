import React, { useCallback, useEffect, useState } from 'react';
import { useNavigate, NavLink, useParams, Navigate } from 'react-router-dom';

import { counselingColumns, closeoutColumns } from '../ServicesCounselingQueue/ServicesCounselingQueue';
import { columns as tooQueueColumns } from '../MoveQueue/MoveQueue';
import { columns as tioQueueColumns } from '../PaymentRequestQueue/PaymentRequestQueue';

import styles from './HeadquartersQueues.module.scss';

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
import { GBLOC } from 'constants/queues';
import TableQueue from 'components/Table/TableQueue';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import SomethingWentWrong from 'shared/SomethingWentWrong';
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

const HeadquartersQueue = ({ isQueueManagementFFEnabled }) => {
  const navigate = useNavigate();
  const { queueType } = useParams();
  const [search, setSearch] = useState({ moveCode: null, dodID: null, customerName: null, paymentRequestCode: null });
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
      paymentRequestCode: null,
    };
    if (!isNullUndefinedOrWhitespace(values.searchText)) {
      if (values.searchType === 'moveCode') {
        payload.moveCode = values.searchText.trim();
      } else if (values.searchType === 'dodID') {
        payload.dodID = values.searchText.trim();
      } else if (values.searchType === 'customerName') {
        payload.customerName = values.searchText.trim();
      } else if (values.searchType === 'paymentRequestCode') {
        payload.paymentRequestCode = values.searchText.trim();
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
            paymentRequestCode={search.paymentRequestCode}
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
            paymentRequestCode={search.paymentRequestCode}
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
          columns={tooQueueColumns(moveLockFlag, isQueueManagementFFEnabled, showBranchFilter)}
          title="All moves"
          handleClick={handleClickNavigateToDetails}
          useQueries={useMovesQueueQueries}
          key="TOO Queue"
          showCSVExport
          csvExportFileNamePrefix="Task-Order-Queue"
          csvExportQueueFetcher={getMovesQueue}
          csvExportQueueFetcherKey="queueMoves"
          sessionStorageKey={queueType}
          // isHeadquartersUser
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
          columns={tioQueueColumns(moveLockFlag, isQueueManagementFFEnabled, showBranchFilter)}
          title="Payment requests"
          handleClick={handleClickNavigateToPaymentRequests}
          useQueries={usePaymentRequestQueueQueries}
          key="TIO Queue"
          showCSVExport
          csvExportFileNamePrefix="Payment-Request-Queue"
          csvExportQueueFetcher={getPaymentRequestsQueue}
          csvExportQueueFetcherKey="queuePaymentRequests"
          sessionStorageKey={queueType}
          // isHeadquartersUser
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
          columns={closeoutColumns(moveLockFlag, inPPMCloseoutGBLOC, null, null, isQueueManagementFFEnabled)}
          title="Moves"
          handleClick={handleClickNavigateToDetails}
          useQueries={useServicesCounselingQueuePPMQueries}
          key="PPM Closeout Queue"
          showCSVExport
          csvExportFileNamePrefix="PPM-Closeout-Queue"
          csvExportQueueFetcher={getServicesCounselingPPMQueue}
          csvExportQueueFetcherKey="queueMoves"
          sessionStorageKey={queueType}
          // isHeadquartersUser
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
          columns={counselingColumns(moveLockFlag, null, null, isQueueManagementFFEnabled)}
          title="Moves"
          handleClick={handleClickNavigateToDetails}
          useQueries={useServicesCounselingQueueQueries}
          key="Counseling Queue"
          showCSVExport
          csvExportFileNamePrefix="Services-Counseling-Queue"
          csvExportQueueFetcher={getServicesCounselingQueue}
          csvExportQueueFetcherKey="queueMoves"
          sessionStorageKey={queueType}
          // isHeadquartersUser
        />
      </div>
    );
  }
  return <NotFound />;
};

export default HeadquartersQueue;
