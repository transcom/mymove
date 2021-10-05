import React from 'react';
import { GridContainer } from '@trussworks/react-uswds';

import styles from './PrimeSimulatorAvailableMoves.module.scss';

import { createHeader } from 'components/Table/utils';
import TableQueue from 'components/Table/TableQueue';
import { usePrimeSimulatorAvailableMovesQueries } from 'hooks/queries';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import SomethingWentWrong from 'shared/SomethingWentWrong';

const columnHeaders = () => [
  createHeader('ID', 'id'),
  createHeader('Move code', 'moveCode'),
  createHeader('Created at', 'createdAt'),
  createHeader('Updated at', 'updatedAt'),
  createHeader('e-Tag', 'eTag'),
  createHeader('Order ID', 'orderID'),
  createHeader('Type', 'ppmType'),
  createHeader('Reference ID', 'referenceId'),
  createHeader('Available to Prime at', 'availableToPrimeAt'),
];

const PrimeSimulatorAvailableMoves = () => {
  const { isLoading, isError } = usePrimeSimulatorAvailableMovesQueries();
  if (isLoading) return <LoadingPlaceholder />;
  if (isError) return <SomethingWentWrong />;

  return (
    <GridContainer containerSize="widescreen" className={styles.gridContainer}>
      <TableQueue
        title="Moves available to Prime"
        columns={columnHeaders()}
        useQueries={usePrimeSimulatorAvailableMovesQueries}
        handleClick={() => {
          return null;
        }}
        defaultSortedColumns={[{ id: 'id', desc: false }]}
      />
    </GridContainer>
  );
};

export default PrimeSimulatorAvailableMoves;
