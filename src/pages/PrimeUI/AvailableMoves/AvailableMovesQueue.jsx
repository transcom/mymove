import React from 'react';
import { useNavigate } from 'react-router-dom';

import styles from './AvailableMovesQueue.module.scss';

import TableQueue from 'components/Table/TableQueue';
import { createHeader } from 'components/Table/utils';
// TODO: This is very clunky. There are shared/formatters and util/formatters
// that determine dates. This way is a way to do it now, but this should be
// refactored as part of TRA work to be done differently across the app.
// For now though, I'm going to be using the `formatDateFromIso` function and
// then leverage a constant for how the date should be formatted.
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import SomethingWentWrong from 'shared/SomethingWentWrong';
import { DATE_TIME_FORMAT_STRING } from 'shared/constants';
import { formatDateFromIso } from 'utils/formatters';
import { usePrimeSimulatorAvailableMovesQueries, useUserQueries } from 'hooks/queries';
import { CHECK_SPECIAL_ORDERS_TYPES, SPECIAL_ORDERS_TYPES } from 'constants/orders';

const columnHeaders = () => [
  createHeader(
    'Move ID',
    (row) => (
      <div>
        {CHECK_SPECIAL_ORDERS_TYPES(row.orderType) ? (
          <span className={styles.specialMoves}>{SPECIAL_ORDERS_TYPES[`${row.orderType}`]}</span>
        ) : null}
        {`${row.id}`}
      </div>
    ),
    {
      id: 'id',
      isFilterable: true,
    },
  ),
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

const PrimeSimulatorAvailableMoves = () => {
  const navigate = useNavigate();

  const { isLoading, isError } = useUserQueries();
  if (isLoading) return <LoadingPlaceholder />;
  if (isError) return <SomethingWentWrong />;

  const handleClick = (values) => {
    navigate(`/simulator/moves/${values.id.props.children[1]}/details`);
  };

  return (
    <TableQueue
      title="Moves available to Prime"
      columns={columnHeaders()}
      useQueries={usePrimeSimulatorAvailableMovesQueries}
      handleClick={handleClick}
      defaultSortedColumns={[{ id: 'availableToPrimeAt', desc: false }]}
      defaultHiddenColumns={['eTag']}
      defaultCanSort
      disableSortBy={false}
      disableMultiSort
      showFilters
      showPagination
      sessionStorageKey="PrimeSimulatorAvailableMoves"
      key="PrimeSimulatorAvailableMoves"
    />
  );
};

export default PrimeSimulatorAvailableMoves;
