import React from 'react';
import { useNavigate } from 'react-router-dom';

import styles from './MoveQueue.module.scss';

import { createHeader } from 'components/Table/utils';
import { useMovesQueueQueries, useUserQueries } from 'hooks/queries';
import { formatDateFromIso, serviceMemberAgencyLabel } from 'utils/formatters';
import MultiSelectCheckBoxFilter from 'components/Table/Filters/MultiSelectCheckBoxFilter';
import SelectFilter from 'components/Table/Filters/SelectFilter';
import { BRANCH_OPTIONS, MOVE_STATUS_OPTIONS, GBLOC, MOVE_STATUS_LABELS } from 'constants/queues';
import TableQueue from 'components/Table/TableQueue';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import SomethingWentWrong from 'shared/SomethingWentWrong';
import DateSelectFilter from 'components/Table/Filters/DateSelectFilter';
import { DATE_FORMAT_STRING } from 'shared/constants';
import { SPECIAL_ORDERS_TYPES } from 'constants/orders';

const columns = (showBranchFilter = true) => [
  createHeader('ID', 'id'),
  createHeader(
    'Customer name',
    (row) => {
      return (
        <div>
          {['WOUNDED_WARRIOR', 'BLUEBARK'].includes(row.orderType) ? (
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
