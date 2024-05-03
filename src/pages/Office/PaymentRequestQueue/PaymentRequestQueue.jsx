import React, { useCallback, useState } from 'react';
import { useNavigate, NavLink, useParams, Navigate } from 'react-router-dom';

import styles from './PaymentRequestQueue.module.scss';

import SearchResultsTable from 'components/Table/SearchResultsTable';
import MoveSearchForm from 'components/MoveSearchForm/MoveSearchForm';
import { usePaymentRequestQueueQueries, useUserQueries, useMoveSearchQueries } from 'hooks/queries';
import { createHeader } from 'components/Table/utils';
import {
  formatDateFromIso,
  serviceMemberAgencyLabel,
  paymentRequestStatusReadable,
  formatAgeToDays,
} from 'utils/formatters';
import MultiSelectCheckBoxFilter from 'components/Table/Filters/MultiSelectCheckBoxFilter';
import SelectFilter from 'components/Table/Filters/SelectFilter';
import DateSelectFilter from 'components/Table/Filters/DateSelectFilter';
import { BRANCH_OPTIONS, GBLOC, PAYMENT_REQUEST_STATUS_OPTIONS } from 'constants/queues';
import TableQueue from 'components/Table/TableQueue';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import SomethingWentWrong from 'shared/SomethingWentWrong';
import { CHECK_SPECIAL_ORDERS_TYPES, SPECIAL_ORDERS_TYPES } from 'constants/orders';
import TabNav from 'components/TabNav';
import { tioRoutes, generalRoutes } from 'constants/routes';
import { roleTypes } from 'constants/userRoles';
import { isNullUndefinedOrWhitespace } from 'shared/utils';
import NotFound from 'components/NotFound/NotFound';

const columns = (showBranchFilter = true) => [
  createHeader('ID', 'id'),
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
    },
  ),
  createHeader('DoD ID', 'customer.dodID', {
    id: 'dodID',
    isFilterable: true,
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
  createHeader(
    'Origin Duty Location',
    (row) => {
      return row.originDutyLocation.name;
    },
    {
      id: 'originDutyLocation',
      isFilterable: true,
    },
  ),
];

const PaymentRequestQueue = () => {
  const { queueType } = useParams();
  const navigate = useNavigate();
  const [search, setSearch] = useState({ moveCode: null, dodID: null, customerName: null });
  const [searchHappened, setSearchHappened] = useState(false);
  const {
    // eslint-disable-next-line camelcase
    data: { office_user },
    isLoading,
    isError,
  } = useUserQueries();
  const onSubmit = useCallback((values) => {
    const payload = {
      moveCode: null,
      dodID: null,
      customerName: null,
    };
    if (!isNullUndefinedOrWhitespace(values.searchText)) {
      if (values.searchType === 'moveCode') {
        payload.moveCode = values.searchText;
      } else if (values.searchType === 'dodID') {
        payload.dodID = values.searchText;
      } else if (values.searchType === 'customerName') {
        payload.customerName = values.searchText;
      }
    }
    setSearch(payload);
    setSearchHappened(true);
  }, []);

  // eslint-disable-next-line camelcase
  const showBranchFilter = office_user?.transportation_office?.gbloc !== GBLOC.USMC;

  const handleClick = (values) => {
    navigate(`/moves/${values.locator}/payment-requests`);
  };

  if (isLoading) return <LoadingPlaceholder />;
  if (isError) return <SomethingWentWrong />;
  if (!queueType) {
    return <Navigate to={tioRoutes.BASE_PAYMENT_REQUEST_QUEUE} />;
  }
  const renderNavBar = () => {
    return (
      <TabNav
        className={styles.tableTabs}
        items={[
          <NavLink
            end
            className={({ isActive }) => (isActive ? 'usa-current' : '')}
            to={tioRoutes.BASE_PAYMENT_REQUEST_QUEUE}
          >
            <span data-testid="payment-request-queue-tab-link" className="tab-title" title="Payment Request Queue">
              Payment Request Queue
            </span>
          </NavLink>,
          <NavLink
            end
            className={({ isActive }) => (isActive ? 'usa-current' : '')}
            to={generalRoutes.BASE_QUEUE_SEARCH_PATH}
          >
            <span data-testid="search-tab-link" className="tab-title" title="Search">
              Search
            </span>
          </NavLink>,
        ]}
      />
    );
  };

  if (queueType === generalRoutes.QUEUE_SEARCH_PATH) {
    return (
      <div data-testid="move-search" className={styles.PaymentRequestQueue}>
        {renderNavBar()}
        <h1>Search for a move</h1>
        <MoveSearchForm onSubmit={onSubmit} role={roleTypes.TIO} />
        {searchHappened && (
          <SearchResultsTable
            showFilters
            showPagination
            defaultCanSort
            disableMultiSort
            disableSortBy={false}
            title="Results"
            handleClick={handleClick}
            useQueries={useMoveSearchQueries}
            moveCode={search.moveCode}
            dodID={search.dodID}
            customerName={search.customerName}
            roleType={roleTypes.TIO}
          />
        )}
      </div>
    );
  }
  if (queueType === tioRoutes.PAYMENT_REQUEST_QUEUE) {
    return (
      <div className={styles.PaymentRequestQueue} data-testid="payment-request-queue">
        {renderNavBar()}
        <TableQueue
          showFilters
          showPagination
          manualSortBy
          defaultCanSort
          defaultSortedColumns={[{ id: 'age', desc: true }]}
          disableMultiSort
          disableSortBy={false}
          columns={columns(showBranchFilter)}
          title="Payment requests"
          handleClick={handleClick}
          useQueries={usePaymentRequestQueueQueries}
        />
      </div>
    );
  }
  return <NotFound />;
};

export default PaymentRequestQueue;
