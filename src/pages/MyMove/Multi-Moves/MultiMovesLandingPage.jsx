import React, { useEffect, useMemo, useState } from 'react';
import { Button } from '@trussworks/react-uswds';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { useLocation, useNavigate } from 'react-router';

import styles from './MultiMovesLandingPage.module.scss';
import MultiMovesMoveHeader from './MultiMovesMoveHeader/MultiMovesMoveHeader';
import MultiMovesMoveContainer from './MultiMovesMoveContainer/MultiMovesMoveContainer';
import {
  mockMovesPCS,
  mockMovesSeparation,
  mockMovesRetirement,
  mockMovesNoPreviousMoves,
  mockMovesNoCurrentMoveWithPreviousMoves,
  mockMovesNoCurrentOrPreviousMoves,
} from './MultiMovesTestData';

import { detectFlags } from 'utils/featureFlags';
import { generatePageTitle } from 'hooks/custom';
import { milmoveLogger } from 'utils/milmoveLog';
import retryPageLoading from 'utils/retryPageLoading';
import { loadInternalSchema } from 'shared/Swagger/ducks';
import { loadUser } from 'store/auth/actions';
import { initOnboarding } from 'store/onboarding/actions';
import Helper from 'components/Customer/Home/Helper';
import { customerRoutes, generalRoutes } from 'constants/routes';

const MultiMovesLandingPage = () => {
  const [setErrorState] = useState({ hasError: false, error: undefined, info: undefined });
  const { search } = useLocation();
  const navigate = useNavigate();
  const searchParams = useMemo(() => new URLSearchParams(search), [search]);

  const handleNewPathClick = (path, paramKey, paramValue) => {
    if (!paramKey || !paramValue) {
      navigate({
        pathname: path,
      });
    } else {
      searchParams.set(paramKey, paramValue);
      navigate({
        pathname: path,
        search: searchParams.toString(),
      });
    }
  };

  // ! This is just used for testing and viewing different variations of data that MilMove will use
  // user can add params of ?moveData=PCS, etc to view different views
  let moves;
  const currentUrl = new URL(window.location.href);
  const moveDataSource = currentUrl.searchParams.get('moveData');
  switch (moveDataSource) {
    case 'PCS':
      moves = mockMovesPCS;
      break;
    case 'retirement':
      moves = mockMovesRetirement;
      break;
    case 'separation':
      moves = mockMovesSeparation;
      break;
    case 'noPreviousMoves':
      moves = mockMovesNoPreviousMoves;
      break;
    case 'noCurrentMove':
      moves = mockMovesNoCurrentMoveWithPreviousMoves;
      break;
    case 'noMoves':
      moves = mockMovesNoCurrentOrPreviousMoves;
      break;
    default:
      moves = mockMovesPCS;
      break;
  }
  // ! end of test data
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

  const handleCreateMoveBtnClick = () => {
    if (moves.previousMoves.length > 0) {
      const profileEditPath = customerRoutes.PROFILE_PATH;
      handleNewPathClick(profileEditPath, 'verifyProfile', 'true');
    } else {
      handleNewPathClick(generalRoutes.HOME_PATH);
    }
  };

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
            <Button className={styles.createMoveBtn} onClick={handleCreateMoveBtnClick}>
              <span>Create a Move</span>
              <div>
                <FontAwesomeIcon icon="plus" />
              </div>
            </Button>
          </div>
          <div className={styles.movesContainer}>
            {moves.currentMove.length > 0 ? (
              <>
                <div data-testid="currentMoveHeader">
                  <MultiMovesMoveHeader title="Current Move" />
                </div>
                <div data-testid="currentMoveContainer">
                  <MultiMovesMoveContainer moves={moves.currentMove} />
                </div>
              </>
            ) : (
              <>
                <div data-testid="currentMoveHeader">
                  <MultiMovesMoveHeader title="Current Moves" />
                </div>
                <div>You do not have a current move.</div>
              </>
            )}
            {moves.previousMoves.length > 0 ? (
              <>
                <div data-testid="prevMovesHeader">
                  <MultiMovesMoveHeader title="Previous Moves" />
                </div>
                <div data-testid="prevMovesContainer">
                  <MultiMovesMoveContainer moves={moves.previousMoves} />
                </div>
              </>
            ) : (
              <>
                <div data-testid="prevMovesHeader">
                  <MultiMovesMoveHeader title="Previous Moves" />
                </div>
                <div>You have no previous moves.</div>
              </>
            )}
          </div>
        </div>
      </div>
    </div>
  ) : null;
};

export default MultiMovesLandingPage;
