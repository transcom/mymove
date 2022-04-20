import React from 'react';

import { HistoryShape } from 'types/router';
import { createHeader } from 'components/Table/utils';
import TableQueue from 'components/Table/TableQueue';
import { usePrimeSimulatorAvailableMovesQueries } from 'hooks/queries';
// TODO: This is very clunky. There are shared/formatters and util/formatters
// that determine dates. This way is a way to do it now, but this should be
// refactored as part of TRA work to be done differently across the app.
// For now though, I'm going to be using the `formatDateFromIso` function and
// then leverage a constant for how the date should be formatted.
import { formatDateFromIso } from 'utils/formatters';
import { DATE_TIME_FORMAT_STRING } from 'shared/constants';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import SomethingWentWrong from 'shared/SomethingWentWrong';

const columnHeaders = () => [
  createHeader('Move ID', 'id', {
    id: 'id',
    isFilterable: true,
  }),
  createHeader('Move code', 'moveCode', {
    id: 'moveCode',
    isFilterable: true,
  }),
  createHeader(
    'Created at',
    (row) => {
      return formatDateFromIso(row.createdAt, DATE_TIME_FORMAT_STRING);
    },
    { id: 'createdAt' },
  ),
  createHeader(
    'Updated at',
    (row) => {
      return formatDateFromIso(row.updatedAt, DATE_TIME_FORMAT_STRING);
    },
    { id: 'updatedAt' },
  ),
  createHeader('e-Tag', 'eTag'),
  createHeader('Order ID', 'orderID'),
  createHeader('Type', 'ppmType'),
  createHeader('Reference ID', 'referenceId'),
  createHeader(
    'Available to Prime at',
    (row) => {
      return formatDateFromIso(row.availableToPrimeAt, DATE_TIME_FORMAT_STRING);
    },
    { id: 'availableToPrimeAt' },
  ),
];

const PrimeSimulatorAvailableMoves = ({ history }) => {
  const { isLoading, isError } = usePrimeSimulatorAvailableMovesQueries();
  if (isLoading) return <LoadingPlaceholder />;
  if (isError) return <SomethingWentWrong />;

  return (
    <TableQueue
      title="Moves available to Prime"
      columns={columnHeaders()}
      useQueries={usePrimeSimulatorAvailableMovesQueries}
      handleClick={(row) => {
        history.push(`/simulator/moves/${row.id}/details`);
      }}
      defaultSortedColumns={[{ id: 'availableToPrimeAt', desc: false }]}
      defaultHiddenColumns={['eTag']}
      defaultCanSort
      disableSortBy={false}
      showFilters
      manualFilters={false}
    />
  );
};

PrimeSimulatorAvailableMoves.propTypes = {
  history: HistoryShape.isRequired,
};

export default PrimeSimulatorAvailableMoves;
