import React from 'react';
import { withRouter } from 'react-router-dom';

import { HistoryShape } from 'types/router';
import { createHeader } from 'components/Table/utils';
import { useMovesQueueQueries } from 'hooks/queries';
import { serviceMemberAgencyLabel } from 'shared/formatters';
import MultiSelectCheckBoxFilter from 'components/Table/Filters/MultiSelectCheckBoxFilter';
import SelectFilter from 'components/Table/Filters/SelectFilter';
import { BRANCH_OPTIONS, MOVE_STATUS_OPTIONS } from 'constants/queues';
import TableQueue from 'components/Table/TableQueue';

const moveStatusOptions = Object.keys(MOVE_STATUS_OPTIONS).map((key) => ({
  value: key,
  label: MOVE_STATUS_OPTIONS[`${key}`],
}));

const branchFilterOptions = [
  { value: '', label: 'All' },
  ...Object.keys(BRANCH_OPTIONS).map((key) => ({
    value: key,
    label: BRANCH_OPTIONS[`${key}`],
  })),
];

const columns = [
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
  createHeader('Status', 'status', {
    isFilterable: true,
    // eslint-disable-next-line react/jsx-props-no-spreading
    Filter: (props) => <MultiSelectCheckBoxFilter options={moveStatusOptions} {...props} />,
  }),
  createHeader('Move Code', 'locator', {
    id: 'moveID',
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
      // eslint-disable-next-line react/jsx-props-no-spreading
      Filter: (props) => <SelectFilter options={branchFilterOptions} {...props} />,
    },
  ),
  createHeader('# of shipments', 'shipmentsCount'),
  createHeader('Destination duty station', 'destinationDutyStation.name', {
    id: 'destinationDutyStation',
    isFilterable: true,
  }),
  createHeader('Origin GBLOC', 'originGBLOC'),
];

const MoveQueue = ({ history }) => {
  const handleClick = (values) => {
    history.push(`/moves/${values.id}/details`);
  };

  return (
    <TableQueue
      showFilters
      showPagination
      columns={columns}
      title="All moves"
      handleClick={handleClick}
      useQueries={useMovesQueueQueries}
    />
  );
};

MoveQueue.propTypes = {
  history: HistoryShape.isRequired,
};

export default withRouter(MoveQueue);
