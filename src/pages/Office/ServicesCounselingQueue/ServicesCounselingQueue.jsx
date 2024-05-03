import React, { useCallback, useEffect, useState } from 'react';
import { generatePath, useNavigate, Navigate, useParams, NavLink } from 'react-router-dom';
import { Button } from '@trussworks/react-uswds';

import styles from './ServicesCounselingQueue.module.scss';

import { createHeader } from 'components/Table/utils';
import MultiSelectCheckBoxFilter from 'components/Table/Filters/MultiSelectCheckBoxFilter';
import SelectFilter from 'components/Table/Filters/SelectFilter';
import DateSelectFilter from 'components/Table/Filters/DateSelectFilter';
import TableQueue from 'components/Table/TableQueue';
import {
  SERVICE_COUNSELING_BRANCH_OPTIONS,
  SERVICE_COUNSELING_QUEUE_MOVE_STATUS_FILTER_OPTIONS,
  SERVICE_COUNSELING_MOVE_STATUS_LABELS,
  SERVICE_COUNSELING_PPM_TYPE_OPTIONS,
  SERVICE_COUNSELING_PPM_TYPE_LABELS,
} from 'constants/queues';
import { generalRoutes, servicesCounselingRoutes } from 'constants/routes';
import {
  useServicesCounselingQueueQueries,
  useServicesCounselingQueuePPMQueries,
  useUserQueries,
  useMoveSearchQueries,
} from 'hooks/queries';
import { DATE_FORMAT_STRING } from 'shared/constants';
import { formatDateFromIso, serviceMemberAgencyLabel } from 'utils/formatters';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import SomethingWentWrong from 'shared/SomethingWentWrong';
import NotFound from 'components/NotFound/NotFound';
import MoveSearchForm from 'components/MoveSearchForm/MoveSearchForm';
import { roleTypes } from 'constants/userRoles';
import SearchResultsTable from 'components/Table/SearchResultsTable';
import TabNav from 'components/TabNav';
import { isBooleanFlagEnabled, isCounselorMoveCreateEnabled } from 'utils/featureFlags';
import retryPageLoading from 'utils/retryPageLoading';
import { milmoveLogger } from 'utils/milmoveLog';
import { CHECK_SPECIAL_ORDERS_TYPES, SPECIAL_ORDERS_TYPES } from 'constants/orders';
import ConnectedFlashMessage from 'containers/FlashMessage/FlashMessage';
import { isNullUndefinedOrWhitespace } from 'shared/utils';

const counselingColumns = () => [
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
        <MultiSelectCheckBoxFilter options={SERVICE_COUNSELING_QUEUE_MOVE_STATUS_FILTER_OPTIONS} {...props} />
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
  }),
];
const closeoutColumns = (ppmCloseoutGBLOC) => [
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
  }),
  createHeader('Destination duty location', 'destinationDutyLocation.name', {
    id: 'destinationDutyLocation',
    isFilterable: true,
  }),
  createHeader('PPM closeout location', 'closeoutLocation', {
    id: 'closeoutLocation',
    // This filter only makes sense if we're not in a closeout GBLOC. Users in a closeout GBLOC will
    // see the same value in this column for every move.
    isFilterable: !ppmCloseoutGBLOC,
  }),
];

const ServicesCounselingQueue = () => {
  const { queueType } = useParams();
  const { data, isLoading, isError } = useUserQueries();

  const navigate = useNavigate();

  const [isCounselorMoveCreateFFEnabled, setisCounselorMoveCreateFFEnabled] = useState(false);
  const [setErrorState] = useState({ hasError: false, error: undefined, info: undefined });

  // Feature Flag
  useEffect(() => {
    const fetchData = async () => {
      try {
        const isEnabled = await isCounselorMoveCreateEnabled();
        setisCounselorMoveCreateFFEnabled(isEnabled);
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

  const handleClick = (values, e) => {
    if (e?.target?.innerHTML === 'Create New Move') {
      navigate(
        generatePath(servicesCounselingRoutes.BASE_CREATE_MOVE_EDIT_CUSTOMER_PATH, { moveCode: values.locator }),
      );
    } else {
      navigate(generatePath(servicesCounselingRoutes.BASE_MOVE_VIEW_PATH, { moveCode: values.locator }));
    }
  };

  const handleAddCustomerClick = () => {
    navigate(generatePath(servicesCounselingRoutes.CREATE_CUSTOMER_PATH));
  };

  const [search, setSearch] = useState({ moveCode: null, dodID: null, customerName: null });
  const [searchHappened, setSearchHappened] = useState(false);
  const counselorMoveCreateFeatureFlag = isBooleanFlagEnabled('counselor_move_create');

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

  // If the office user is in a closeout GBLOC and on the closeout tab, then we will want to disable
  // the column filter for the closeout location column because it will have no effect.
  const officeUserGBLOC = data?.office_user?.transportation_office?.gbloc;
  const inPPMCloseoutGBLOC = officeUserGBLOC === 'TVCB' || officeUserGBLOC === 'NAVY' || officeUserGBLOC === 'USCG';
  if (isLoading) return <LoadingPlaceholder />;
  if (isError) return <SomethingWentWrong />;
  if (!queueType) {
    return inPPMCloseoutGBLOC ? (
      <Navigate to={servicesCounselingRoutes.BASE_QUEUE_CLOSEOUT_PATH} />
    ) : (
      <Navigate to={servicesCounselingRoutes.BASE_QUEUE_COUNSELING_PATH} />
    );
  }

  const renderNavBar = () => {
    return (
      <TabNav
        className={styles.tableTabs}
        items={[
          <NavLink
            end
            className={({ isActive }) => (isActive ? 'usa-current' : '')}
            to={servicesCounselingRoutes.BASE_QUEUE_COUNSELING_PATH}
          >
            <span data-testid="counseling-tab-link" className="tab-title">
              Counseling Queue
            </span>
          </NavLink>,
          <NavLink
            end
            className={({ isActive }) => (isActive ? 'usa-current' : '')}
            to={servicesCounselingRoutes.BASE_QUEUE_CLOSEOUT_PATH}
          >
            <span data-testid="closeout-tab-link" className="tab-title">
              PPM Closeout Queue
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
        <ConnectedFlashMessage />
        <div className={styles.searchFormContainer}>
          <h1>Search for a move</h1>
          {searchHappened && counselorMoveCreateFeatureFlag && (
            <Button type="submit" onClick={handleAddCustomerClick} className={styles.addCustomerBtn}>
              Add Customer
            </Button>
          )}
        </div>
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
            isCounselorMoveCreateFFEnabled={isCounselorMoveCreateFFEnabled}
          />
        )}
      </div>
    );
  }

  if (queueType === 'PPM-closeout') {
    return (
      <div className={styles.ServicesCounselingQueue}>
        {renderNavBar()}
        <TableQueue
          showFilters
          showPagination
          manualSortBy
          defaultCanSort
          defaultSortedColumns={[{ id: 'closeoutInitiated', desc: false }]}
          disableMultiSort
          disableSortBy={false}
          columns={closeoutColumns(inPPMCloseoutGBLOC)}
          title="Moves"
          handleClick={handleClick}
          useQueries={useServicesCounselingQueuePPMQueries}
        />
      </div>
    );
  }
  if (queueType === 'counseling') {
    return (
      <div className={styles.ServicesCounselingQueue}>
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
          columns={counselingColumns()}
          title="Moves"
          handleClick={handleClick}
          useQueries={useServicesCounselingQueueQueries}
        />
      </div>
    );
  }

  return <NotFound />;
};

export default ServicesCounselingQueue;
