import React, { useCallback, useEffect, useState } from 'react';
import { useNavigate, NavLink, useParams, Navigate, generatePath } from 'react-router-dom';
import { Dropdown } from '@trussworks/react-uswds';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';

import styles from './MoveQueue.module.scss';

import { createHeader } from 'components/Table/utils';
import {
  useMovesQueueQueries,
  useUserQueries,
  useMoveSearchQueries,
  useDestinationRequestsQueueQueries,
} from 'hooks/queries';
import { getDestinationRequestsQueue, getMovesQueue } from 'services/ghcApi';
import { formatDateFromIso, serviceMemberAgencyLabel } from 'utils/formatters';
import MultiSelectCheckBoxFilter from 'components/Table/Filters/MultiSelectCheckBoxFilter';
import SelectFilter from 'components/Table/Filters/SelectFilter';
import { MOVE_STATUS_OPTIONS, GBLOC, MOVE_STATUS_LABELS, BRANCH_OPTIONS } from 'constants/queues';
import TableQueue from 'components/Table/TableQueue';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import SomethingWentWrong from 'shared/SomethingWentWrong';
import DateSelectFilter from 'components/Table/Filters/DateSelectFilter';
import { DATE_FORMAT_STRING, DEFAULT_EMPTY_VALUE } from 'shared/constants';
import { CHECK_SPECIAL_ORDERS_TYPES, SPECIAL_ORDERS_TYPES } from 'constants/orders';
import MoveSearchForm from 'components/MoveSearchForm/MoveSearchForm';
import { roleTypes } from 'constants/userRoles';
import SearchResultsTable from 'components/Table/SearchResultsTable';
import TabNav from 'components/TabNav';
import { generalRoutes, tooRoutes } from 'constants/routes';
import { isNullUndefinedOrWhitespace } from 'shared/utils';
import NotFound from 'components/NotFound/NotFound';
import { isBooleanFlagEnabled } from 'utils/featureFlags';
import handleQueueAssignment from 'utils/queues';

export const columns = (moveLockFlag, isQueueManagementEnabled, showBranchFilter = true) => {
  const cols = [
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
        id: 'customerName',
        isFilterable: true,
        exportValue: (row) => {
          return `${row.customer.last_name}, ${row.customer.first_name}`;
        },
      },
    ),
    createHeader('DoD ID', 'customer.edipi', {
      id: 'edipi',
      isFilterable: true,
      exportValue: (row) => {
        return row.customer.edipi;
      },
    }),
    createHeader('EMPLID', 'customer.emplid', {
      id: 'emplid',
      isFilterable: true,
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
    createHeader('Counseling office', 'counselingOffice', {
      id: 'counselingOffice',
      isFilterable: true,
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
            <div data-label="assignedSelect" data-testid="assigned-col" className={styles.assignedToCol} key={row.id}>
              <Dropdown
                defaultValue={row.assignedTo?.officeUserId}
                onChange={(e) => handleQueueAssignment(row.id, e.target.value, roleTypes.TOO)}
                title="Assigned dropdown"
              >
                <option value={null}>{DEFAULT_EMPTY_VALUE}</option>
                {row.availableOfficeUsers?.map(({ lastName, firstName, officeUserId }) => (
                  <option value={officeUserId} key={`filterOption_${officeUserId}`}>
                    {`${lastName}, ${firstName}`}
                  </option>
                ))}
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

const MoveQueue = ({ isQueueManagementFFEnabled }) => {
  const navigate = useNavigate();
  const { queueType } = useParams();
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

  const onSubmit = useCallback((values) => {
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

  const handleEditProfileClick = (locator) => {
    navigate(generatePath(tooRoutes.BASE_CUSTOMER_INFO_EDIT_PATH, { moveCode: locator }));
  };

  const handleClick = (values, e) => {
    // if the user clicked the profile icon to edit, we want to route them elsewhere
    // since we don't have innerText, we are using the data-label property
    const editProfileDiv = e.target.closest('div[data-label="editProfile"]');
    const assignedSelect = e.target.closest('div[data-label="assignedSelect"]');
    if (editProfileDiv) {
      navigate(generatePath(tooRoutes.BASE_CUSTOMER_INFO_EDIT_PATH, { moveCode: values.locator }));
    } else if (assignedSelect) {
      // don't want to page redirect if clicking in the assigned dropdown
    } else {
      navigate(generatePath(tooRoutes.BASE_MOVE_VIEW_PATH, { moveCode: values.locator }));
    }
  };

  if (isLoading) return <LoadingPlaceholder />;
  if (isError) return <SomethingWentWrong />;
  if (!queueType) {
    return <Navigate to={tooRoutes.BASE_MOVE_QUEUE} />;
  }
  const renderNavBar = () => {
    return (
      <TabNav
        className={styles.tableTabs}
        items={[
          <NavLink end className={({ isActive }) => (isActive ? 'usa-current' : '')} to={tooRoutes.BASE_MOVE_QUEUE}>
            <span data-testid="closeout-tab-link" className="tab-title" title="Move Queue">
              Task Order Queue
            </span>
          </NavLink>,
          <NavLink
            end
            className={({ isActive }) => (isActive ? 'usa-current' : '')}
            to={tooRoutes.BASE_DESTINATION_REQUESTS_QUEUE}
          >
            <span className="tab-title" title="Destination Requests Queue">
              Destination Requests Queue
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
      <div data-testid="move-search" className={styles.ServicesCounselingQueue}>
        {renderNavBar()}
        <h1>Search for a move</h1>
        <MoveSearchForm onSubmit={onSubmit} role={roleTypes.TOO} />
        {searchHappened && (
          <SearchResultsTable
            showFilters
            showPagination
            defaultCanSort
            disableMultiSort
            disableSortBy={false}
            title="Results"
            handleClick={handleClick}
            handleEditProfileClick={handleEditProfileClick}
            useQueries={useMoveSearchQueries}
            moveCode={search.moveCode}
            dodID={search.dodID}
            customerName={search.customerName}
            paymentRequestCode={search.paymentRequestCode}
            roleType={roleTypes.TOO}
          />
        )}
      </div>
    );
  }
  if (queueType === tooRoutes.MOVE_QUEUE) {
    return (
      <div className={styles.MoveQueue} data-testid="move-queue">
        {renderNavBar()}
        <TableQueue
          showFilters
          showPagination
          manualSortBy
          defaultCanSort
          defaultSortedColumns={[{ id: 'status', desc: false }]}
          disableMultiSort
          disableSortBy={false}
          columns={columns(moveLockFlag, isQueueManagementFFEnabled, showBranchFilter)}
          title="All moves"
          handleClick={handleClick}
          useQueries={useMovesQueueQueries}
          showCSVExport
          csvExportFileNamePrefix="Task-Order-Queue"
          csvExportQueueFetcher={getMovesQueue}
          csvExportQueueFetcherKey="queueMoves"
          sessionStorageKey={queueType}
          key={queueType}
        />
      </div>
    );
  }
  if (queueType === tooRoutes.DESTINATION_REQUESTS_QUEUE) {
    return (
      <div className={styles.MoveQueue} data-testid="destination-requests-queue">
        {renderNavBar()}
        <TableQueue
          showFilters
          showPagination
          manualSortBy
          defaultCanSort
          defaultSortedColumns={[{ id: 'status', desc: false }]}
          disableMultiSort
          disableSortBy={false}
          columns={columns(moveLockFlag, isQueueManagementFFEnabled, showBranchFilter)}
          title="Destination requests"
          handleClick={handleClick}
          useQueries={useDestinationRequestsQueueQueries}
          showCSVExport
          csvExportFileNamePrefix="Destination-Requests-Queue"
          csvExportQueueFetcher={getDestinationRequestsQueue}
          csvExportQueueFetcherKey="queueMoves"
          sessionStorageKey={queueType}
          key={queueType}
        />
      </div>
    );
  }
  return <NotFound />;
};

export default MoveQueue;
