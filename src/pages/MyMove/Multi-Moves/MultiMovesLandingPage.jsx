import React, { useEffect, useState } from 'react';
import { Button } from '@trussworks/react-uswds';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { useNavigate } from 'react-router';
import { connect } from 'react-redux';

import styles from './MultiMovesLandingPage.module.scss';
import MultiMovesMoveHeader from './MultiMovesMoveHeader/MultiMovesMoveHeader';
import MultiMovesMoveContainer from './MultiMovesMoveContainer/MultiMovesMoveContainer';

import { detectFlags } from 'utils/featureFlags';
import { generatePageTitle } from 'hooks/custom';
import { milmoveLogger } from 'utils/milmoveLog';
import retryPageLoading from 'utils/retryPageLoading';
import { loadInternalSchema } from 'shared/Swagger/ducks';
import { loadUser } from 'store/auth/actions';
import { initOnboarding } from 'store/onboarding/actions';
import Helper from 'components/Customer/Home/Helper';
import { customerRoutes } from 'constants/routes';
import { withContext } from 'shared/AppContext';
import withRouter from 'utils/routing';
import requireCustomerState from 'containers/requireCustomerState/requireCustomerState';
import { selectAllMoves, selectIsProfileComplete, selectServiceMemberFromLoggedInUser } from 'store/entities/selectors';

const MultiMovesLandingPage = ({ serviceMember, serviceMemberMoves }) => {
  const [setErrorState] = useState({ hasError: false, error: undefined, info: undefined });
  const navigate = useNavigate();

  useEffect(() => {
    const fetchData = async () => {
      try {
        loadInternalSchema();
        loadUser();
        initOnboarding();
        document.title = generatePageTitle('MilMove');
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

  // handles logic when user clicks "Create a Move" button
  // if they have previous moves, they'll need to validate their profile
  // if they do not have previous moves, then they don't need to validate
  const handleCreateMoveBtnClick = () => {
    if (serviceMemberMoves && serviceMemberMoves.previousMoves && serviceMemberMoves.previousMoves.length !== 0) {
      const profileEditPath = customerRoutes.PROFILE_PATH;
      navigate(profileEditPath, { state: { needsToVerifyProfile: true } });
    } else {
      navigate(customerRoutes.MOVE_HOME_PAGE);
    }
  };

  // ! WILL ONLY SHOW IF MULTIMOVE FLAG IS TRUE
  return flags.multiMove ? (
    <div>
      <div className={styles.homeContainer}>
        <header data-testid="customerHeader" className={styles.customerHeader}>
          <div className={`usa-prose grid-container ${styles['grid-container']}`}>
            <h2>
              {serviceMember.first_name} {serviceMember.last_name}
            </h2>
          </div>
        </header>
        <div className={`usa-prose grid-container ${styles['grid-container']}`}>
          {serviceMemberMoves && serviceMemberMoves.previousMoves && serviceMemberMoves.previousMoves.length === 0 ? (
            <Helper title="Welcome to MilMove!" className={styles['helper-paragraph-only']}>
              <p data-testid="welcomeHeader">
                Select &quot;Create a Move&quot; to get started. <br />
                <br />
                If you encounter any issues please contact your local Transportation Office or the Help Desk for further
                assistance.
              </p>
            </Helper>
          ) : (
            <Helper title="Welcome to MilMove!" className={styles['helper-paragraph-only']}>
              <p data-testid="welcomeHeader">
                Select &quot;Create a Move&quot; to get started. <br />
                <br />
                Once you have validated your profile, pleasee click the &quot;Validate&quot; button and proceed to
                starting your move. <br />
                If you encounter any issues please contact your local Transportation Office or the Help Desk for further
                assistance.
              </p>
            </Helper>
          )}
          <div className={styles.centeredContainer}>
            <Button className={styles.createMoveBtn} onClick={handleCreateMoveBtnClick} data-testid="createMoveBtn">
              <span>Create a Move</span>
              <div>
                <FontAwesomeIcon icon="plus" />
              </div>
            </Button>
          </div>
          <div className={styles.movesContainer}>
            {serviceMemberMoves && serviceMemberMoves.currentMove && serviceMemberMoves.currentMove.length !== 0 ? (
              <>
                <div data-testid="currentMoveHeader">
                  <MultiMovesMoveHeader title="Current Move" />
                </div>
                <div data-testid="currentMoveContainer">
                  <MultiMovesMoveContainer moves={serviceMemberMoves.currentMove} />
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
            {serviceMemberMoves && serviceMemberMoves.previousMoves && serviceMemberMoves.previousMoves.length !== 0 ? (
              <>
                <div data-testid="prevMovesHeader">
                  <MultiMovesMoveHeader title="Previous Moves" />
                </div>
                <div data-testid="prevMovesContainer">
                  <MultiMovesMoveContainer moves={serviceMemberMoves.previousMoves} />
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

MultiMovesLandingPage.defaultProps = {
  serviceMember: null,
};

const mapStateToProps = (state) => {
  const serviceMember = selectServiceMemberFromLoggedInUser(state);
  const serviceMemberMoves = selectAllMoves(state);

  return {
    isProfileComplete: selectIsProfileComplete(state),
    serviceMember,
    serviceMemberMoves,
  };
};

// in order to avoid setting up proxy server only for storybook, pass in stub function so API requests don't fail
const mergeProps = (stateProps, dispatchProps, ownProps) => ({
  ...stateProps,
  ...dispatchProps,
  ...ownProps,
});

export default withContext(
  withRouter(connect(mapStateToProps, mergeProps)(requireCustomerState(MultiMovesLandingPage))),
);
