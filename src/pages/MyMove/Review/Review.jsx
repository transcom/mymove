import React from 'react';
import { func } from 'prop-types';
import { connect } from 'react-redux';
import { GridContainer, Grid } from '@trussworks/react-uswds';
import { generatePath } from 'react-router';

import styles from './Review.module.scss';

import { MatchShape } from 'types/router';
import ScrollToTop from 'components/ScrollToTop';
import ConnectedFlashMessage from 'containers/FlashMessage/FlashMessage';
import ConnectedSummary from 'components/Customer/Review/Summary/Summary';
import 'scenes/Review/Review.css';
import formStyles from 'styles/form.module.scss';
import WizardNavigation from 'components/Customer/WizardNavigation/WizardNavigation';
import { customerRoutes, generalRoutes } from 'constants/routes';
import { showLoggedInUser as showLoggedInUserAction } from 'shared/Entities/modules/user';
import { selectCurrentMove } from 'store/entities/selectors';
import { MoveShape } from 'types/customerShapes';
import MOVE_STATUSES from 'constants/moves';

const Review = ({ push, currentMove, match }) => {
  const handleCancel = () => {
    push(generalRoutes.HOME_PATH);
  };

  const handleNext = () => {
    const nextPath = generatePath(customerRoutes.MOVE_AGREEMENT_PATH, {
      moveId: match.params.moveId,
    });
    push(nextPath);
  };

  const inDraftStatus = currentMove.status === MOVE_STATUSES.DRAFT;

  return (
    <GridContainer>
      <ScrollToTop />
      <Grid row>
        <Grid col desktop={{ col: 8, offset: 2 }}>
          <ConnectedFlashMessage />
        </Grid>
      </Grid>
      <Grid row>
        <Grid col desktop={{ col: 8, offset: 2 }}>
          <div className={styles.reviewMoveContainer}>
            <div className={styles.reviewMoveHeaderContainer}>
              <h1 data-testid="review-move-header">Review your details</h1>
              <p>
                You’re almost done setting up your move. Double&#8209;check that your information is accurate, add more
                shipments if needed, then move on to the final step.
              </p>
            </div>
            <ConnectedSummary />
            <div className={formStyles.formActions}>
              <WizardNavigation
                onNextClick={handleNext}
                disableNext={false}
                onCancelClick={handleCancel}
                isFirstPage
                showFinishLater
                readOnly={!inDraftStatus}
              />
            </div>
          </div>
        </Grid>
      </Grid>
    </GridContainer>
  );
};

Review.propTypes = {
  currentMove: MoveShape.isRequired,
  push: func.isRequired,
  match: MatchShape.isRequired,
};

const mapStateToProps = (state, ownProps) => {
  return {
    ...ownProps,
    currentMove: selectCurrentMove(state) || {},
  };
};

const mapDispatchToProps = {
  showLoggedInUser: showLoggedInUserAction,
};

export default connect(mapStateToProps, mapDispatchToProps)(Review);
