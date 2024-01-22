import React, { useEffect, useState } from 'react';
import { Button } from '@trussworks/react-uswds';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';

import styles from './MultiMovesLandingPage.module.scss';
import MultiMovesMoveHeader from './MultiMovesMoveHeader/MultiMovesMoveHeader';
import MultiMovesMoveContainer from './MultiMovesMoveContainer/MultiMovesMoveContainer';
import { movesPCS, movesRetirement } from './MultiMovesTestData';

import { generatePageTitle } from 'hooks/custom';
import { milmoveLogger } from 'utils/milmoveLog';
import retryPageLoading from 'utils/retryPageLoading';
import { loadInternalSchema } from 'shared/Swagger/ducks';
import { loadUser } from 'store/auth/actions';
import { initOnboarding } from 'store/onboarding/actions';
import Helper from 'components/Customer/Home/Helper';

const MultiMovesLandingPage = () => {
  const [setErrorState] = useState({ hasError: false, error: undefined, info: undefined });
  useEffect(() => {
    const fetchData = async () => {
      try {
        loadInternalSchema();
        loadUser();
        initOnboarding();
        document.title = generatePageTitle('MilMove');

        const script = document.createElement('script');
        script.src = '//rum-static.pingdom.net/pa-6567b05deff3250012000426.js';
        script.async = true;
        document.body.appendChild(script);
      } catch (error) {
        const { message } = error;
        milmoveLogger.error({ message, info: null });
        setErrorState({
          hasError: true,
          error,
          info: null,
        });
        retryPageLoading(error);
      }
    };

    fetchData();
  }, [setErrorState]);

  // including test data to use - imported from MultiMovesTestData
  const moves = movesPCS;
  // const moves = movesSeparation;
  // const moves = movesRetirement;

  return (
    <div>
      <div className={styles.homeContainer}>
        <header data-testid="customer-header" className={styles.customerHeader}>
          <div className={`usa-prose grid-container ${styles['grid-container']}`}>
            <h2>First Last</h2>
          </div>
        </header>
        <div className={`usa-prose grid-container ${styles['grid-container']}`}>
          <Helper title="Welcome to MilMove!" className={styles['helper-paragraph-only']}>
            <p>
              We can put information at the top here - potentially important contact info or basic instructions on how
              to start a move?
            </p>
          </Helper>
          <div className={styles.centeredContainer}>
            <Button className={styles.createMoveBtn}>
              <span>Create a Move</span>
              <div>
                <FontAwesomeIcon icon="plus" />
              </div>
            </Button>
          </div>
          <div className={styles.movesContainer}>
            <MultiMovesMoveHeader data-testid="currentMoveHeader" title="Current Move" />
            <MultiMovesMoveContainer moves={moves.currentMove} />
            <MultiMovesMoveHeader title="Previous Moves" />
            <MultiMovesMoveContainer moves={moves.previousMoves} />
          </div>
        </div>
      </div>
    </div>
  );
};

export default MultiMovesLandingPage;
