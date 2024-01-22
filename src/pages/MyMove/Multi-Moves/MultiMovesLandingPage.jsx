import React, { useEffect, useState } from 'react';
import { Button } from '@trussworks/react-uswds';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';

import styles from './MultiMovesLandingPage.module.scss';
import MultiMovesMoveHeader from './MultiMovesMoveHeader/MultiMovesMoveHeader';
import MultiMovesMoveContainer from './MultiMovesMoveContainer/MultiMovesMoveContainer';
import { mockMovesPCS } from './MultiMovesTestData';

import { detectFlags } from 'utils/featureFlags';
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

  const flags = detectFlags(process.env.NODE_ENV, window.location.host, window.location.search);

  // including test data to use - imported from MultiMovesTestData
  const moves = mockMovesPCS;
  // const moves = mockMovesSeparation;
  // const moves = mockMovesRetirement;

  // ! WILL ONLY SHOW IF MULTIMOVE FLAG IS TRUE
  return flags.multiMove ? (
    <div>
      <div className={styles.homeContainer}>
        <header data-testid="customerHeader" className={styles.customerHeader}>
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
            <div data-testid="currentMoveHeader">
              <MultiMovesMoveHeader title="Current Move" />
            </div>
            <div data-testid="currentMoveContainer">
              <MultiMovesMoveContainer moves={moves.currentMove} />
            </div>
            <div data-testid="prevMovesHeader">
              <MultiMovesMoveHeader title="Previous Moves" />
            </div>
            <div data-testid="prevMovesContainer">
              <MultiMovesMoveContainer moves={moves.previousMoves} />
            </div>
          </div>
        </div>
      </div>
    </div>
  ) : null;
};

export default MultiMovesLandingPage;
