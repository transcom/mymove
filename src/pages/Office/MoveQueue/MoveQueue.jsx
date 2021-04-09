import React from 'react';
import { withRouter } from 'react-router-dom';

import styles from './MoveQueue.module.scss';

import { HistoryShape } from 'types/router';
import { createHeader } from 'components/Table/utils';
import { useMovesQueueQueries, useUserQueries } from 'hooks/queries';
import { serviceMemberAgencyLabel } from 'shared/formatters';
import MultiSelectCheckBoxFilter from 'components/Table/Filters/MultiSelectCheckBoxFilter';
import SelectFilter from 'components/Table/Filters/SelectFilter';
import { BRANCH_OPTIONS, MOVE_STATUS_OPTIONS, GBLOC, MOVE_STATUS_LABELS } from 'constants/queues';
import TableQueue from 'components/Table/TableQueue';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import SomethingWentWrong from 'shared/SomethingWentWrong';

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
  createHeader('Destination duty station', 'destinationDutyStation.name', {
    id: 'destinationDutyStation',
    isFilterable: true,
  }),
  createHeader('Origin GBLOC', 'originGBLOC', { disableSortBy: true }),
];

const MoveQueue = ({ history }) => {
  const {
    // eslint-disable-next-line camelcase
    data: { office_user },
    isLoading,
    isError,
  } = useUserQueries();

  const showBranchFilter = office_user?.transportation_office?.gbloc !== GBLOC.USMC;

  const handleClick = (values) => {
    history.push(`/moves/${values.locator}/details`);
  };

  if (isLoading) return <LoadingPlaceholder />;
  if (isError) return <SomethingWentWrong />;

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

MoveQueue.propTypes = {
  history: HistoryShape.isRequired,
};

export default withRouter(MoveQueue);
