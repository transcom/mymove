import React from 'react';
// import { generatePath } from 'react-router';
// import { useHistory } from 'react-router-dom';
import { GridContainer } from '@trussworks/react-uswds';

import styles from './PrimeSimulatorAvailableMoves.module.scss';

import PrimeSimulatorListMoveCard from 'components/Office/PrimeSimulatorListMoveCard/PrimeSimulatorListMoveCard';
// import DateSelectFilter from 'components / Table / Filters / DateSelectFilter';
import { usePrimeSimulatorAvailableMovesQueries } from 'hooks/queries';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import SomethingWentWrong from 'shared/SomethingWentWrong';

const PrimeSimulatorAvailableMoves = () => {
  const { listMoves, isLoading, isError } = usePrimeSimulatorAvailableMovesQueries();
  if (isLoading) return <LoadingPlaceholder />;
  if (isError) return <SomethingWentWrong />;

  return (
    <GridContainer className={styles.gridContainer} data-testid="tio-payment-request-details">
      <h1>Available Moves</h1>
      <div className={styles.section} id="available-moves">
        {listMoves.map((listMove) => (
          <PrimeSimulatorListMoveCard listMove={listMove} key={listMove.id} />
        ))}
      </div>
    </GridContainer>
  );
};

export default PrimeSimulatorAvailableMoves;
