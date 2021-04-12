import React from 'react';
import { useHistory } from 'react-router-dom';

import styles from './ServicesCounselingQueue.module.scss';

import { useServicesCounselingQueueQueries, useUserQueries } from 'hooks/queries';
import { createHeader } from 'components/Table/utils';
import {
  BRANCH_OPTIONS,
  SERVICE_COUNSELING_MOVE_STATUS_OPTIONS,
  GBLOC,
  SERVICE_COUNSELING_MOVE_STATUS_LABELS,
} from 'constants/queues';
import { formatDateFromIso, serviceMemberAgencyLabel } from 'shared/formatters';
import MultiSelectCheckBoxFilter from 'components/Table/Filters/MultiSelectCheckBoxFilter';
import SelectFilter from 'components/Table/Filters/SelectFilter';
import DateSelectFilter from 'components/Table/Filters/DateSelectFilter';
import TableQueue from 'components/Table/TableQueue';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import SomethingWentWrong from 'shared/SomethingWentWrong';

const columns = (isMarineCorpsUser = false) => [
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
      return formatDateFromIso(row.requestedMoveDate, 'DD MMM YYYY');
    },
    {
      id: 'requestedMoveDate',
      isFilterable: true,
      Filter: DateSelectFilter,
    },
  ),
  createHeader(
    'Date submitted',
    (row) => {
      return formatDateFromIso(row.submittedAt, 'DD MMM YYYY');
    },
    {
      id: 'submittedAt',
      isFilterable: true,
      Filter: DateSelectFilter,
    },
  ),
  createHeader(
    'Branch',
    (row) => {
      return serviceMemberAgencyLabel(row.customer.agency);
    },
    {
      id: 'branch',
      isFilterable: !isMarineCorpsUser,
      Filter: (props) => (
        // eslint-disable-next-line react/jsx-props-no-spreading
        <SelectFilter options={BRANCH_OPTIONS} {...props} />
      ),
      disableSortBy: isMarineCorpsUser,
    },
  ),
  createHeader('Origin GBLOC', 'originGBLOC', {
    isFilterable: isMarineCorpsUser,
    disableSortBy: !isMarineCorpsUser,
  }), // If the user is in the USMC GBLOC they will have many different GBLOCs and will want to sort and filter
  createHeader('Destination duty station', 'destinationDutyStation.name', {
    id: 'destinationDutyStation',
    isFilterable: true,
  }),
];

const ServicesCounselingQueue = () => {
  const {
    // eslint-disable-next-line camelcase
    data: { office_user },
    isLoading,
    isError,
  } = useUserQueries();

  const history = useHistory();

  const isMarineCorpsUser = office_user?.transportation_office?.gbloc === GBLOC.USMC;

  const handleClick = (values) => {
    history.push(`/counseling/moves/${values.locator}/details`);
  };

  if (isLoading) return <LoadingPlaceholder />;
  if (isError) return <SomethingWentWrong />;

  return (
    // TODO: Pull out header count and add new move button
    <div className={styles.ServicesCounselingQueue}>
      <TableQueue
        showFilters
        showPagination
        manualSortBy
        defaultCanSort
        defaultSortedColumns={[{ id: 'submittedAt', desc: false }]}
        disableMultiSort
        disableSortBy={false}
        columns={columns(isMarineCorpsUser)}
        title="Moves"
        handleClick={handleClick}
        useQueries={useServicesCounselingQueueQueries}
      />
    </div>
  );
};

export default ServicesCounselingQueue;
