import React from 'react';
import { withRouter } from 'react-router-dom';

import { HistoryShape } from 'types/router';
import { createHeader } from 'components/Table/utils';
import { useMovesQueueQueries, useUserQueries } from 'hooks/queries';
import { serviceMemberAgencyLabel } from 'shared/formatters';
import MultiSelectCheckBoxFilter from 'components/Table/Filters/MultiSelectCheckBoxFilter';
import SelectFilter from 'components/Table/Filters/SelectFilter';
import { BRANCH_OPTIONS, BRANCH_OPTIONS_NO_MARINES, MOVE_STATUS_OPTIONS, GBLOC } from 'constants/queues';
import TableQueue from 'components/Table/TableQueue';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import SomethingWentWrong from 'shared/SomethingWentWrong';

const columns = (includeBranchOptionMarines = true) => [
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
    Filter: (props) => <MultiSelectCheckBoxFilter options={MOVE_STATUS_OPTIONS} {...props} />,
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
      Filter: (props) => (
        // eslint-disable-next-line react/jsx-props-no-spreading
        <SelectFilter options={includeBranchOptionMarines ? BRANCH_OPTIONS : BRANCH_OPTIONS_NO_MARINES} {...props} />
      ),
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
  const {
    // eslint-disable-next-line camelcase
    data: { office_user },
    isLoading,
    isError,
  } = useUserQueries();

  const includeBranchOptionMarines = office_user?.transportation_office?.gbloc === GBLOC.USMC;

  const handleClick = (values) => {
    history.push(`/moves/${values.id}/details`);
  };

  if (isLoading) return <LoadingPlaceholder />;
  if (isError) return <SomethingWentWrong />;

  return (
    <TableQueue
      showFilters
      showPagination
      columns={columns(includeBranchOptionMarines)}
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
