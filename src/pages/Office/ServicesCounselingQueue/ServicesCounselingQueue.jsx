import React from 'react';
import { generatePath } from 'react-router';
import { useHistory, Switch, Route } from 'react-router-dom';

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
} from 'constants/queues';
import { servicesCounselingRoutes } from 'constants/routes';
import { useServicesCounselingQueueQueries, useServicesCounselingQueuePPMQueries, useUserQueries } from 'hooks/queries';
import { DATE_FORMAT_STRING } from 'shared/constants';
import { formatDateFromIso, serviceMemberAgencyLabel } from 'utils/formatters';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import SomethingWentWrong from 'shared/SomethingWentWrong';

const columns = () => [
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

const ServicesCounselingQueue = () => {
  const { isLoading, isError } = useUserQueries();

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
            defaultSortedColumns={[{ id: 'submittedAt', desc: false }]}
            disableMultiSort
            disableSortBy={false}
            columns={columns()}
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
            columns={columns()}
            title="Moves"
            handleClick={handleClick}
            useQueries={useServicesCounselingQueueQueries}
          />
        </Route>
      </Switch>
    </div>
  );
};

export default ServicesCounselingQueue;
