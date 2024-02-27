import React from 'react';
import { connect } from 'react-redux';
import { GridContainer, Grid } from '@trussworks/react-uswds';
import { generatePath, useNavigate, useParams } from 'react-router-dom';

import styles from './Review.module.scss';

import ConnectedFlashMessage from 'containers/FlashMessage/FlashMessage';
import ConnectedSummary from 'components/Customer/Review/Summary/Summary';
import WizardNavigation from 'components/Customer/WizardNavigation/WizardNavigation';
import MOVE_STATUSES from 'constants/moves';
import { customerRoutes, generalRoutes } from 'constants/routes';
import 'scenes/Review/Review.css';
import { selectAllMoves } from 'store/entities/selectors';
import formStyles from 'styles/form.module.scss';
import { SHIPMENT_OPTIONS } from 'shared/constants';
import { isPPMShipmentComplete } from 'utils/shipments';
import { useTitle } from 'hooks/custom';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';

const Review = ({ serviceMemberMoves }) => {
  useTitle('Move review');
  const navigate = useNavigate();
  const { moveId } = useParams();
  const handleCancel = () => {
    navigate(generalRoutes.HOME_PATH);
  };

  // loading placeholder while data loads - this handles any async issues
  if (!serviceMemberMoves || !serviceMemberMoves.currentMove || !serviceMemberMoves.previousMoves) {
    return (
      <div className={styles.homeContainer}>
        <div className={`usa-prose grid-container ${styles['grid-container']}`}>
          <LoadingPlaceholder />
        </div>
      </div>
    );
  }

  const currentMove = serviceMemberMoves.currentMove.find((m) => m.id === moveId);
  const previousMove = serviceMemberMoves.previousMoves.find((m) => m.id === moveId);
  const move = currentMove || previousMove;
  const { mtoShipments } = move;

  const handleNext = () => {
    const nextPath = generatePath(customerRoutes.MOVE_AGREEMENT_PATH, {
      moveId,
    });
    navigate(nextPath);
  };

  const inDraftStatus = move.status === MOVE_STATUSES.DRAFT;

  // PPM shipments can be left in an incomplete state, disable proceeding to the signature move
  // submission page to force them to complete or delete the shipment.
  const hasCompletedPPMShipments = mtoShipments
    ?.filter((s) => s.shipmentType === SHIPMENT_OPTIONS.PPM)
    ?.every((s) => isPPMShipmentComplete(s));

  return (
    <GridContainer>
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
                You are almost done setting up your move. Double&#8209;check that your information is accurate, add more
                shipments if needed, then move on to the final step.
              </p>
            </div>
            <ConnectedSummary />
            <div className={formStyles.formActions}>
              <WizardNavigation
                onNextClick={handleNext}
                disableNext={!hasCompletedPPMShipments || !mtoShipments.length}
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

const mapStateToProps = (state, ownProps) => {
  const serviceMemberMoves = selectAllMoves(state);
  return {
    ...ownProps,
    serviceMemberMoves,
  };
};

export default connect(mapStateToProps)(Review);
