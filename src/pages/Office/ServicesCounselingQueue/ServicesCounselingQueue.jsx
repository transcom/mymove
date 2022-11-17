import React from 'react';
import { generatePath } from 'react-router';
import { useHistory, Switch, Route } from 'react-router-dom-old';

import styles from './ServicesCounselingQueue.module.scss';

import { createHeader } from 'components/Table/utils';
import MultiSelectCheckBoxFilter from 'components/Table/Filters/MultiSelectCheckBoxFilter';
import SelectFilter from 'components/Table/Filters/SelectFilter';
import DateSelectFilter from 'components/Table/Filters/DateSelectFilter';
import TableQueue from 'components/Table/TableQueue';
import {
  SERVICE_COUNSELING_BRANCH_OPTIONS,
  SERVICE_COUNSELING_MOVE_STATUS_OPTIONS,
  SERVICE_COUNSELING_MOVE_STATUS_LABELS,
  SERVICE_COUNSELING_PPM_TYPE_OPTIONS,
  SERVICE_COUNSELING_PPM_TYPE_LABELS,
} from 'constants/queues';
import { servicesCounselingRoutes } from 'constants/routes';
import { useServicesCounselingQueueQueries, useServicesCounselingQueuePPMQueries, useUserQueries } from 'hooks/queries';
import { DATE_FORMAT_STRING } from 'shared/constants';
import { formatDateFromIso, serviceMemberAgencyLabel } from 'utils/formatters';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import SomethingWentWrong from 'shared/SomethingWentWrong';

const counselingColumns = () => [
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
      // eslint-disable-next-line react/jsx-props-no-spreading
      Filter: (props) => <MultiSelectCheckBoxFilter options={SERVICE_COUNSELING_MOVE_STATUS_OPTIONS} {...props} />,
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
  const { data, isLoading, isError } = useUserQueries();

  const history = useHistory();

  const handleClick = (values) => {
    history.push(
      generatePath(servicesCounselingRoutes.MOVE_VIEW_PATH, {
        moveCode: values.locator,
      }),
    );
  };

  if (isLoading) return <LoadingPlaceholder />;
  if (isError) return <SomethingWentWrong />;

  // If the office user is in a closeout GBLOC and on the closeout tab, then we will want to disable
  // the column filter for the closeout location column because it will have no effect.
  const officeUserGBLOC = data.office_user.transportation_office.gbloc;
  const inPPMCloseoutGBLOC = officeUserGBLOC === 'TVCB' || officeUserGBLOC === 'NAVY' || officeUserGBLOC === 'USCG';

  return (
    // TODO: Pull out header count and add new move button
    <div className={styles.ServicesCounselingQueue}>
      <Switch>
        <Route path={servicesCounselingRoutes.QUEUE_CLOSEOUT_PATH} exact>
          <TableQueue
            showTabs
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
        </Route>
        <Route
          path={[
            servicesCounselingRoutes.QUEUE_COUNSELING_PATH,
            servicesCounselingRoutes.DEFAULT_QUEUE_PATH,
            servicesCounselingRoutes.QUEUE_VIEW_PATH,
          ]}
          exact
        >
          <div>
            <TableQueue
              className={styles.ServicesCounseling}
              showTabs
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
        </Route>
      </Switch>
    </div>
  );
};

export default ServicesCounselingQueue;
