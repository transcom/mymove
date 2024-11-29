import React, { useCallback, useEffect, useState } from 'react';
import { useNavigate, NavLink, useParams, Navigate } from 'react-router-dom';
import { Dropdown } from '@trussworks/react-uswds';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';

import styles from './PaymentRequestQueue.module.scss';

import SearchResultsTable from 'components/Table/SearchResultsTable';
import MoveSearchForm from 'components/MoveSearchForm/MoveSearchForm';
import { usePaymentRequestQueueQueries, useUserQueries, useMoveSearchQueries } from 'hooks/queries';
import { getPaymentRequestsQueue } from 'services/ghcApi';
import { createHeader } from 'components/Table/utils';
import {
  formatDateFromIso,
  serviceMemberAgencyLabel,
  paymentRequestStatusReadable,
  formatAgeToDays,
} from 'utils/formatters';
import SelectFilter from 'components/Table/Filters/SelectFilter';
import DateSelectFilter from 'components/Table/Filters/DateSelectFilter';
import { BRANCH_OPTIONS, GBLOC } from 'constants/queues';
import TableQueue from 'components/Table/TableQueue';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import SomethingWentWrong from 'shared/SomethingWentWrong';
import { CHECK_SPECIAL_ORDERS_TYPES, SPECIAL_ORDERS_TYPES } from 'constants/orders';
import TabNav from 'components/TabNav';
import { tioRoutes, generalRoutes } from 'constants/routes';
import { roleTypes } from 'constants/userRoles';
import { isNullUndefinedOrWhitespace } from 'shared/utils';
import NotFound from 'components/NotFound/NotFound';
import { isBooleanFlagEnabled } from 'utils/featureFlags';
import { DEFAULT_EMPTY_VALUE, PAYMENT_REQUEST_STATUS } from 'shared/constants';
import handleQueueAssignment from 'utils/queues';

export const columns = (moveLockFlag, isQueueManagementEnabled, showBranchFilter = true) => {
  const cols = [
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
        id: 'customerName',
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
    createHeader('EMPLID', 'customer.emplid', {
      id: 'emplid',
      isFilterable: true,
    }),
    createHeader(
      'Status',
      (row) => {
        return row.status !== PAYMENT_REQUEST_STATUS.PAYMENT_REQUESTED
          ? paymentRequestStatusReadable(row.status)
          : null;
      },
      {
        id: 'status',
        disableSortBy: true,
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

  if (isQueueManagementEnabled)
    cols.push(
      createHeader(
        'Assigned',
        (row) => {
          return !row?.assignable ? (
            <div data-testid="assigned-col">
              {row.assignedTo ? `${row.assignedTo?.lastName}, ${row.assignedTo?.firstName}` : ''}
            </div>
          ) : (
            <div data-label="assignedSelect" data-testid="assigned-col" className={styles.assignedToCol}>
              <Dropdown
                defaultValue={row.assignedTo?.officeUserId}
                onChange={(e) => {
                  handleQueueAssignment(row.moveID, e.target.value, roleTypes.TIO);
                }}
                title="Assigned dropdown"
              >
                <option value={null}>{DEFAULT_EMPTY_VALUE}</option>
                {row.availableOfficeUsers.map(({ lastName, firstName, officeUserId }) => {
                  return (
                    <option value={officeUserId} key={`filterOption_${officeUserId}`}>
                      {`${lastName}, ${firstName}`}
                    </option>
                  );
                })}
              </Dropdown>
            </div>
          );
        },
        {
          id: 'assignedTo',
          isFilterable: true,
          exportValue: (row) => {
            return row.assignedTo ? `${row.assignedTo?.lastName}, ${row.assignedTo?.firstName}` : '';
          },
        },
      ),
    );

  return cols;
};

const PaymentRequestQueue = ({ isQueueManagementFFEnabled }) => {
  const { queueType } = useParams();
  const navigate = useNavigate();
  const [search, setSearch] = useState({ moveCode: null, dodID: null, customerName: null, paymentRequestCode: null });
  const [searchHappened, setSearchHappened] = useState(false);
  const [moveLockFlag, setMoveLockFlag] = useState(false);

  useEffect(() => {
    const fetchData = async () => {
      const lockedMoveFlag = await isBooleanFlagEnabled('move_lock');
      setMoveLockFlag(lockedMoveFlag);
    };

    fetchData();
  }, []);

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
      paymentRequestCode: null,
    };
    if (!isNullUndefinedOrWhitespace(values.searchText)) {
      if (values.searchType === 'moveCode') {
        payload.moveCode = values.searchText;
      } else if (values.searchType === 'dodID') {
        payload.dodID = values.searchText;
      } else if (values.searchType === 'customerName') {
        payload.customerName = values.searchText;
      } else if (values.searchType === 'paymentRequestCode') {
        payload.paymentRequestCode = values.searchText.trim();
      }
    }
    setSearch(payload);
    setSearchHappened(true);
  }, []);

  // eslint-disable-next-line camelcase
  const showBranchFilter = office_user?.transportation_office?.gbloc !== GBLOC.USMC;

  const handleClick = (values, e) => {
    const assignedSelect = e.target.closest('div[data-label="assignedSelect"]');

    if (assignedSelect) {
      // do nothing when clicking on the assignedSelect column
    } else {
      navigate(`/moves/${values.locator}/payment-requests`);
    }
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
            paymentRequestCode={search.paymentRequestCode}
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
          columns={columns(moveLockFlag, isQueueManagementFFEnabled, showBranchFilter)}
          title="Payment requests"
          handleClick={handleClick}
          useQueries={usePaymentRequestQueueQueries}
          showCSVExport
          csvExportFileNamePrefix="Payment-Request-Queue"
          csvExportQueueFetcher={getPaymentRequestsQueue}
          csvExportQueueFetcherKey="queuePaymentRequests"
          sessionStorageKey={queueType}
          key={queueType}
        />
      </div>
    );
  }
  return <NotFound />;
};

export default PaymentRequestQueue;
