import React, { useCallback, useState } from 'react';
import { useNavigate, NavLink, useParams } from 'react-router-dom';

import styles from './MoveQueue.module.scss';

import { createHeader } from 'components/Table/utils';
import { useMovesQueueQueries, useUserQueries, useMoveSearchQueries } from 'hooks/queries';
import { formatDateFromIso, serviceMemberAgencyLabel } from 'utils/formatters';
import MultiSelectCheckBoxFilter from 'components/Table/Filters/MultiSelectCheckBoxFilter';
import SelectFilter from 'components/Table/Filters/SelectFilter';
import { BRANCH_OPTIONS, MOVE_STATUS_OPTIONS, GBLOC, MOVE_STATUS_LABELS } from 'constants/queues';
import TableQueue from 'components/Table/TableQueue';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import SomethingWentWrong from 'shared/SomethingWentWrong';
import DateSelectFilter from 'components/Table/Filters/DateSelectFilter';
import { DATE_FORMAT_STRING } from 'shared/constants';
import MoveSearchForm from 'components/MoveSearchForm/MoveSearchForm';
import { roleTypes } from 'constants/userRoles';
import SearchResultsTable from 'components/Table/SearchResultsTable';
import TabNav from 'components/TabNav';
import { generalRoutes, tooRoutes } from 'constants/routes';

const columns = (showBranchFilter = true) => [
  createHeader('ID', 'id'),
  createHeader(
    'Customer name',
    (row) => {
      return `${row.customer.last_name}, ${row.customer.first_name}`;
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
  }),
  createHeader('Origin GBLOC', 'originGBLOC', { disableSortBy: true }),
];

const MoveQueue = () => {
  const navigate = useNavigate();
  const queueType = useParams();
  const [search, setSearch] = useState({ moveCode: null, dodID: null, customerName: null });
  const [searchHappened, setSearchHappened] = useState(false);

  const onSubmit = useCallback((values) => {
    const payload = {
      moveCode: null,
      dodID: null,
      customerName: null,
    };

    if (values.searchType === 'moveCode') {
      payload.moveCode = values.searchText;
    } else if (values.searchType === 'dodID') {
      payload.dodID = values.searchText;
    } else if (values.searchType === 'customerName') {
      payload.customerName = values.searchText;
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

  const handleClick = (values) => {
    navigate(`/moves/${values.locator}/details`);
  };

  if (isLoading) return <LoadingPlaceholder />;
  if (isError) return <SomethingWentWrong />;

  const renderNavBar = () => {
    return (
      <TabNav
        className={styles.tableTabs}
        items={[
          <NavLink
            end
            className={({ isActive }) => (isActive ? 'usa-current' : '')}
            to={tooRoutes.BASE_QUEUE_COUNSELING_PATH}
          >
            <span data-testid="closeout-tab-link" className="tab-title">
              Move Queue
            </span>
          </NavLink>,
          <NavLink
            end
            className={({ isActive }) => (isActive ? 'usa-current' : '')}
            to={generalRoutes.BASE_QUEUE_SEARCH_PATH}
          >
            <span data-testid="search-tab-link" className="tab-title">
              Search
            </span>
          </NavLink>,
        ]}
      />
    );
  };

  if (queueType === 'Search') {
    return (
      <div data-testid="move-search" className={styles.ServicesCounselingQueue}>
        {renderNavBar()}
        <h1>Search for a move</h1>
        <MoveSearchForm onSubmit={onSubmit} role={roleTypes.SERVICES_COUNSELOR} />
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
            roleType={roleTypes.SERVICES_COUNSELOR}
          />
        )}
      </div>
    );
  }

  return (
    <div className={styles.MoveQueue}>
      <TableQueue
        showFilters
        showPagination
        manualSortBy
        defaultCanSort
        defaultSortedColumns={[{ id: 'status', desc: false }]}
        disableMultiSort
        disableSortBy={false}
        columns={columns(showBranchFilter)}
        title="All moves"
        handleClick={handleClick}
        useQueries={useMovesQueueQueries}
      />
    </div>
  );
};

export default MoveQueue;
