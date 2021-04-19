import React from 'react';
import { bool, func } from 'prop-types';
import { connect } from 'react-redux';
import { GridContainer, Grid } from '@trussworks/react-uswds';

import styles from './Review.module.scss';

import ScrollToTop from 'components/ScrollToTop';
import { hasShortHaulError } from 'utils/incentives';
import ConnectedFlashMessage from 'containers/FlashMessage/FlashMessage';
import ConnectedSummary from 'components/Customer/Review/Summary/index';
import 'scenes/Review/Review.css';
import formStyles from 'styles/form.module.scss';
import WizardNavigation from 'components/Customer/WizardNavigation/WizardNavigation';
import { selectServiceMemberFromLoggedInUser } from 'store/entities/selectors';
import { selectPPMEstimateError } from 'store/onboarding/selectors';
import { customerRoutes, generalRoutes } from 'constants/routes';

const Review = ({ push, canMoveNext }) => {
  const handleCancel = () => {
    push(generalRoutes.HOME_PATH);
  };

  const handleNext = () => {
    push(customerRoutes.MOVE_AGREEMENT_PATH);
  };

  return (
    <GridContainer>
      <ScrollToTop />
      <ConnectedFlashMessage />
      <Grid row>
        <Grid col desktop={{ col: 8, offset: 2 }}>
          <div className={styles.reviewMoveContainer}>
            <div className={styles.reviewMoveHeaderContainer}>
              <h1 data-testid="review-move-header">Review your details</h1>
              <p>
                You’re almost done setting up your move. Double&#8209;check that your information is accurate, then move
                on to the final step.
              </p>
            </div>
            <ConnectedSummary />
            <div className={formStyles.formActions}>
              <WizardNavigation
                onNextClick={handleNext}
                disableNext={!canMoveNext}
                onCancelClick={handleCancel}
                isFirstPage
                showFinishLater
              />
            </div>
          </div>
        </Grid>
      </Grid>
    </GridContainer>
  );
};

Review.propTypes = {
  canMoveNext: bool.isRequired,
  push: func.isRequired,
};

const mapStateToProps = (state, ownProps) => {
  const serviceMember = selectServiceMemberFromLoggedInUser(state);
  const ppmEstimateError = selectPPMEstimateError(state);

  const ppmEstimate = {
    hasEstimateError: !!ppmEstimateError,
    rateEngineError: ppmEstimateError,
    originDutyStationZip: serviceMember?.current_station?.address?.postal_code,
  };

  return {
    ...ownProps,
    ppmEstimate,
    canMoveNext: !hasShortHaulError(ppmEstimateError),
  };
};

export default connect(mapStateToProps)(Review);
