import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { Button, Dropdown } from '@trussworks/react-uswds';
import React, { useCallback, useContext, useEffect, useState } from 'react';
import { connect } from 'react-redux';
import { generatePath, Navigate, NavLink, useNavigate, useParams } from 'react-router-dom';

import styles from './ServicesCounselingQueue.module.scss';

import CustomerSearchForm from 'components/CustomerSearchForm/CustomerSearchForm';
import MoveSearchForm from 'components/MoveSearchForm/MoveSearchForm';
import NotFound from 'components/NotFound/NotFound';
import SelectedGblocContext from 'components/Office/GblocSwitcher/SelectedGblocContext';
import DateSelectFilter from 'components/Table/Filters/DateSelectFilter';
import MultiSelectTypeAheadCheckBoxFilter from 'components/Table/Filters/MutliSelectTypeAheadCheckboxFilter';
import SelectFilter from 'components/Table/Filters/SelectFilter';
import SearchResultsTable from 'components/Table/SearchResultsTable';
import TableQueue from 'components/Table/TableQueue';
import { createHeader } from 'components/Table/utils';
import TabNav from 'components/TabNav';
import { CHECK_SPECIAL_ORDERS_TYPES, SPECIAL_ORDERS_TYPES } from 'constants/orders';
import {
  BRANCH_OPTIONS,
  QUEUE_TYPES,
  SERVICE_COUNSELING_MOVE_STATUS_LABELS,
  SERVICE_COUNSELING_PPM_STATUS_LABELS,
  SERVICE_COUNSELING_PPM_STATUS_OPTIONS,
  SERVICE_COUNSELING_PPM_TYPE_LABELS,
  SERVICE_COUNSELING_PPM_TYPE_OPTIONS,
} from 'constants/queues';
import { generalRoutes, servicesCounselingRoutes } from 'constants/routes';
import { elevatedPrivilegeTypes } from 'constants/userPrivileges';
import { roleTypes } from 'constants/userRoles';
import ConnectedFlashMessage from 'containers/FlashMessage/FlashMessage';
import {
  useCustomerSearchQueries,
  useMoveSearchQueries,
  useServicesCounselingQueuePPMQueries,
  useServicesCounselingQueueQueries,
  useUserQueries,
} from 'hooks/queries';
import {
  getServicesCounselingOriginLocations,
  getServicesCounselingPPMQueue,
  getServicesCounselingQueue,
} from 'services/ghcApi';
import { DATE_FORMAT_STRING, DEFAULT_EMPTY_VALUE, MOVE_STATUSES } from 'shared/constants';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import SomethingWentWrong from 'shared/SomethingWentWrong';
import { isNullUndefinedOrWhitespace } from 'shared/utils';
import { selectLoggedInUser } from 'store/entities/selectors';
import { setRefetchQueue as setRefetchQueueAction } from 'store/general/actions';
import { isBooleanFlagEnabled, isCounselorMoveCreateEnabled } from 'utils/featureFlags';
import { formatDateFromIso, serviceMemberAgencyLabel } from 'utils/formatters';
import { milmoveLogger } from 'utils/milmoveLog';
import { handleQueueAssignment, getQueue } from 'utils/queues';
import retryPageLoading from 'utils/retryPageLoading';

export const counselingColumns = (
  moveLockFlag,
  originLocationList,
  supervisor,
  queueType,
  isQueueManagementEnabled,
  setRefetchQueue,
) => {
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
    createHeader('Move code', 'locator', {
      id: 'locator',
      isFilterable: true,
    }),
    createHeader(
      'Status',
      (row) => {
        return row.status !== MOVE_STATUSES.SERVICE_COUNSELING_COMPLETED
          ? SERVICE_COUNSELING_MOVE_STATUS_LABELS[`${row.status}`]
          : null;
      },
      {
        id: 'status',
        disableSortBy: true,
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
          <SelectFilter options={BRANCH_OPTIONS} {...props} />
        ),
      },
    ),
    createHeader('Origin GBLOC', 'originGBLOC', {
      disableSortBy: true,
    }), // If the user is in the USMC GBLOC they will have many different GBLOCs and will want to sort and filter
    supervisor
      ? createHeader(
          'Origin duty location',
          (row) => {
            return `${row.originDutyLocation.name}`;
          },
          {
            id: 'originDutyLocation',
            isFilterable: true,
            exportValue: (row) => {
              return row.originDutyLocation?.name;
            },
            Filter: (props) => (
              <MultiSelectTypeAheadCheckBoxFilter
                options={originLocationList}
                placeholder="Start typing a duty location..."
                // eslint-disable-next-line react/jsx-props-no-spreading
                {...props}
              />
            ),
          },
        )
      : createHeader('Origin duty location', 'originDutyLocation.name', {
          id: 'originDutyLocation',
          isFilterable: true,
          exportValue: (row) => {
            return row.originDutyLocation?.name;
          },
        }),
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
            <div>{row.assignedTo ? `${row.assignedTo?.lastName}, ${row.assignedTo?.firstName}` : ''}</div>
          ) : (
            <div data-label="assignedSelect" className={styles.assignedToCol}>
              <Dropdown
                key={row.id}
                onChange={(e) => {
                  handleQueueAssignment(row.id, e.target.value, getQueue(queueType));
                  setRefetchQueue(true);
                }}
                title="Assigned dropdown"
              >
                <option value={null}>{DEFAULT_EMPTY_VALUE}</option>
                {row.availableOfficeUsers.map(({ lastName, firstName, officeUserId }) => (
                  <option
                    value={officeUserId}
                    key={officeUserId}
                    selected={row.assignedTo?.officeUserId === officeUserId}
                  >
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
export const closeoutColumns = (
  moveLockFlag,
  ppmCloseoutGBLOC,
  ppmCloseoutOriginLocationList,
  supervisor,
  queueType,
  isQueueManagementEnabled,
  setRefetchQueue,
) => {
  const cols = [
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
          <SelectFilter options={BRANCH_OPTIONS} {...props} />
        ),
      },
    ),
    createHeader(
      'Status',
      (row) => {
        return SERVICE_COUNSELING_PPM_STATUS_LABELS[`${row.ppmStatus}`];
      },
      {
        id: 'ppmStatus',
        isFilterable: true,
        Filter: (props) => (
          // eslint-disable-next-line react/jsx-props-no-spreading
          <SelectFilter options={SERVICE_COUNSELING_PPM_STATUS_OPTIONS} {...props} />
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
    supervisor
      ? createHeader(
          'Origin duty location',
          (row) => {
            return `${row.originDutyLocation.name}`;
          },
          {
            id: 'originDutyLocation',
            isFilterable: true,
            exportValue: (row) => {
              return row.originDutyLocation?.name;
            },
            Filter: (props) => (
              <MultiSelectTypeAheadCheckBoxFilter
                options={ppmCloseoutOriginLocationList}
                placeholder="Start typing a duty location..."
                // eslint-disable-next-line react/jsx-props-no-spreading
                {...props}
              />
            ),
          },
        )
      : createHeader('Origin duty location', 'originDutyLocation.name', {
          id: 'originDutyLocation',
          isFilterable: true,
          exportValue: (row) => {
            return row.originDutyLocation?.name;
          },
        }),
    createHeader('Counseling office', 'counselingOffice', {
      id: 'counselingOffice',
      isFilterable: true,
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
  if (isQueueManagementEnabled)
    cols.push(
      createHeader(
        'Assigned',
        (row) => {
          return !row?.assignable ? (
            <div>{row.assignedTo ? `${row.assignedTo?.lastName}, ${row.assignedTo?.firstName}` : ''}</div>
          ) : (
            <div data-label="assignedSelect" className={styles.assignedToCol} key={row.id}>
              <Dropdown
                onChange={(e) => {
                  handleQueueAssignment(row.id, e.target.value, getQueue(queueType));
                  setRefetchQueue(true);
                }}
                title="Assigned dropdown"
              >
                <option value={null}>{DEFAULT_EMPTY_VALUE}</option>
                {row.availableOfficeUsers.map(({ lastName, firstName, officeUserId }) => (
                  <option
                    value={officeUserId}
                    key={officeUserId}
                    selected={row.assignedTo?.officeUserId === officeUserId}
                  >
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

const ServicesCounselingQueue = ({
  userPrivileges,
  isQueueManagementFFEnabled,
  officeUser,
  isBulkAssignmentFFEnabled,
  setRefetchQueue,
  activeRole,
}) => {
  const { queueType } = useParams();
  const { data, isLoading, isError } = useUserQueries();

  const navigate = useNavigate();

  const [isCounselorMoveCreateFFEnabled, setisCounselorMoveCreateFFEnabled] = useState(false);
  const [moveLockFlag, setMoveLockFlag] = useState(false);
  const [setErrorState] = useState({ hasError: false, error: undefined, info: undefined });
  const [originLocationList, setOriginLocationList] = useState([]);
  const [ppmCloseoutOriginLocationList, setPpmCloseoutOriginLocationList] = useState([]);
  const supervisor = userPrivileges
    ? userPrivileges.some((p) => p.privilegeType === elevatedPrivilegeTypes.SUPERVISOR)
    : false;

  const gblocContext = useContext(SelectedGblocContext);
  const { selectedGbloc } =
    officeUser?.transportation_office_assignments?.length > 1 && gblocContext
      ? gblocContext
      : { selectedGbloc: undefined };

  // Feature Flag
  useEffect(() => {
    const getOriginLocationList = (needsPPMCloseout, gbloc) => {
      if (supervisor) {
        getServicesCounselingOriginLocations(needsPPMCloseout, gbloc).then((response) => {
          if (needsPPMCloseout) {
            setPpmCloseoutOriginLocationList(response);
          } else {
            setOriginLocationList(response);
          }
        });
      }
    };

    getOriginLocationList(true, selectedGbloc);
    getOriginLocationList(false, selectedGbloc);

    const fetchData = async () => {
      try {
        const isEnabled = await isCounselorMoveCreateEnabled();
        setisCounselorMoveCreateFFEnabled(isEnabled);
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
  }, [setErrorState, supervisor, selectedGbloc]);

  const handleEditProfileClick = (locator) => {
    navigate(generatePath(servicesCounselingRoutes.BASE_CUSTOMER_INFO_EDIT_PATH, { moveCode: locator }));
  };

  const handleClick = (values, e) => {
    // if the user clicked the profile icon to edit, we want to route them elsewhere
    // since we don't have innerText, we are using the data-label property
    const editProfileDiv = e.target.closest('div[data-label="editProfile"]');
    const assignedSelect = e.target.closest('div[data-label="assignedSelect"]');
    if (editProfileDiv) {
      navigate(generatePath(servicesCounselingRoutes.BASE_CUSTOMER_INFO_EDIT_PATH, { moveCode: values.locator }));
    } else if (assignedSelect) {
      // do nothing
    } else {
      navigate(generatePath(servicesCounselingRoutes.BASE_MOVE_VIEW_PATH, { moveCode: values.locator }));
    }
  };

  const handleCustomerSearchClick = (values) => {
    navigate(
      generatePath(servicesCounselingRoutes.BASE_CUSTOMERS_CUSTOMER_INFO_PATH, { customerId: values.customerID }),
    );
  };

  const handleAddCustomerClick = () => {
    navigate(generatePath(servicesCounselingRoutes.CREATE_CUSTOMER_PATH));
  };

  const [search, setSearch] = useState({ moveCode: null, dodID: null, customerName: null });
  const [searchHappened, setSearchHappened] = useState(false);

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

  const tabs = [
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
      onClick={() => setSearchHappened(false)}
    >
      <span data-testid="search-tab-link" className="tab-title">
        Move Search
      </span>
    </NavLink>,
  ];

  // when FEATURE_FLAG_COUNSELOR_MOVE_CREATE is removed,
  // this can simply be the tabs for this component
  const ffTabs = [
    ...tabs,
    <NavLink
      end
      className={({ isActive }) => (isActive ? 'usa-current' : '')}
      to={servicesCounselingRoutes.BASE_CUSTOMER_SEARCH_PATH}
      onClick={() => setSearchHappened(false)}
    >
      <span data-testid="search-tab-link" className="tab-title">
        Customer Search
      </span>
    </NavLink>,
  ];

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
  const navTabs = () => (isCounselorMoveCreateFFEnabled ? ffTabs : tabs);

  const renderNavBar = () => {
    return <TabNav className={styles.tableTabs} items={navTabs()} />;
  };

  if (queueType === generalRoutes.QUEUE_SEARCH_PATH) {
    return (
      <div data-testid="move-search" className={styles.ServicesCounselingQueue}>
        {renderNavBar()}
        <ConnectedFlashMessage />
        <div className={styles.searchFormContainer}>
          <h1>Search for a move</h1>
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
            handleEditProfileClick={handleEditProfileClick}
            useQueries={useMoveSearchQueries}
            moveCode={search.moveCode}
            dodID={search.dodID}
            customerName={search.customerName}
            roleType={roleTypes.SERVICES_COUNSELOR}
            searchType="move"
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
          columns={closeoutColumns(
            moveLockFlag,
            inPPMCloseoutGBLOC,
            ppmCloseoutOriginLocationList,
            supervisor,
            queueType,
            isQueueManagementFFEnabled,
            setRefetchQueue,
          )}
          title="Moves"
          handleClick={handleClick}
          useQueries={useServicesCounselingQueuePPMQueries}
          showCSVExport
          csvExportFileNamePrefix="PPM-Closeout-Queue"
          csvExportQueueFetcher={getServicesCounselingPPMQueue}
          csvExportQueueFetcherKey="queueMoves"
          sessionStorageKey={queueType}
          key={queueType}
          isSupervisor={supervisor}
          isBulkAssignmentFFEnabled={isBulkAssignmentFFEnabled}
          queueType={QUEUE_TYPES.CLOSEOUT}
          activeRole={activeRole}
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
          columns={counselingColumns(
            moveLockFlag,
            originLocationList,
            supervisor,
            queueType,
            isQueueManagementFFEnabled,
            setRefetchQueue,
          )}
          title="Moves"
          handleClick={handleClick}
          useQueries={useServicesCounselingQueueQueries}
          showCSVExport
          csvExportFileNamePrefix="Services-Counseling-Queue"
          csvExportQueueFetcher={getServicesCounselingQueue}
          csvExportQueueFetcherKey="queueMoves"
          sessionStorageKey={queueType}
          key={queueType}
          isSupervisor={supervisor}
          isBulkAssignmentFFEnabled={isBulkAssignmentFFEnabled}
          queueType={QUEUE_TYPES.COUNSELING}
          activeRole={activeRole}
        />
      </div>
    );
  }
  if (queueType === 'customer-search') {
    return (
      <div data-testid="customer-search" className={styles.ServicesCounselingQueue}>
        {renderNavBar()}
        <ConnectedFlashMessage />
        <div className={styles.searchFormContainer}>
          <h1>Search for a customer</h1>
          {searchHappened && isCounselorMoveCreateFFEnabled && (
            <Button type="submit" onClick={handleAddCustomerClick} className={styles.addCustomerBtn}>
              Add Customer
            </Button>
          )}
        </div>
        <CustomerSearchForm onSubmit={onSubmit} role={roleTypes.SERVICES_COUNSELOR} />
        {searchHappened && (
          <SearchResultsTable
            showFilters
            showPagination
            defaultCanSort
            disableMultiSort
            disableSortBy={false}
            title="Results"
            defaultHiddenColumns={['customerID']}
            handleClick={handleCustomerSearchClick}
            useQueries={useCustomerSearchQueries}
            dodID={search.dodID}
            customerName={search.customerName}
            roleType={roleTypes.SERVICES_COUNSELOR}
            searchType="customer"
          />
        )}
      </div>
    );
  }

  return <NotFound />;
};
const mapStateToProps = (state) => {
  const user = selectLoggedInUser(state);

  return {
    officeUser: user?.office_user || {},
    setRefetchQueue: state.generalState.setRefetchQueue,
  };
};

const mapDispatchToProps = { setRefetchQueue: setRefetchQueueAction };

export default connect(mapStateToProps, mapDispatchToProps)(ServicesCounselingQueue);
