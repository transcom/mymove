import React, { useEffect, useState } from 'react';
import { connect } from 'react-redux';
import { GridContainer, Grid } from '@trussworks/react-uswds';
import { generatePath, useNavigate, useParams } from 'react-router-dom';

import { isBooleanFlagEnabled } from '../../../utils/featureFlags';

import styles from './Review.module.scss';

import ConnectedFlashMessage from 'containers/FlashMessage/FlashMessage';
import ConnectedSummary from 'components/Customer/Review/Summary/Summary';
import WizardNavigation from 'components/Customer/WizardNavigation/WizardNavigation';
import MOVE_STATUSES from 'constants/moves';
import { customerRoutes } from 'constants/routes';
import 'scenes/Review/Review.css';
import { selectAllMoves, selectServiceMemberFromLoggedInUser } from 'store/entities/selectors';
import formStyles from 'styles/form.module.scss';
import { SHIPMENT_TYPES } from 'shared/constants';
import { isPPMShipmentComplete, isBoatShipmentComplete, isMobileHomeShipmentComplete } from 'utils/shipments';
import { useTitle } from 'hooks/custom';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import { getAllMoves } from 'services/internalApi';
import { updateAllMoves as updateAllMovesAction } from 'store/entities/actions';

const Review = ({ serviceMemberId, serviceMemberMoves, updateAllMoves }) => {
  useTitle('Move review');
  const navigate = useNavigate();
  const [multiMove, setMultiMove] = useState(false);
  const { moveId } = useParams();
  const handleCancel = () => {
    if (multiMove) {
      navigate(generatePath(customerRoutes.MOVE_HOME_PATH, { moveId }));
    } else {
      navigate(customerRoutes.MOVE_HOME_PAGE);
    }
  };

  // fetching all move data on load since this component is dependent on that data
  // this will run each time the component is loaded/accessed
  useEffect(() => {
    getAllMoves(serviceMemberId).then((response) => {
      updateAllMoves(response);
    });
    isBooleanFlagEnabled('multi_move').then((enabled) => {
      setMultiMove(enabled);
    });
  }, [updateAllMoves, serviceMemberId]);

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

  // PPM, boat, and mobile home shipments can be left in an incomplete state, disable proceeding to the signature move
  // submission page to force them to complete or delete the shipment.
  const hasIncompleteShipment = () => {
    if (!mtoShipments) return false;
    const shipmentValidators = {
      [SHIPMENT_TYPES.PPM]: isPPMShipmentComplete,
      [SHIPMENT_TYPES.BOAT_HAUL_AWAY]: isBoatShipmentComplete,
      [SHIPMENT_TYPES.BOAT_TOW_AWAY]: isBoatShipmentComplete,
      [SHIPMENT_TYPES.MOBILE_HOME]: isMobileHomeShipmentComplete,
    };

    return mtoShipments.some((shipment) => {
      const validateShipment = shipmentValidators[shipment.shipmentType];
      return validateShipment && !validateShipment(shipment);
    });
  };

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
                disableNext={hasIncompleteShipment() || !mtoShipments?.length}
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
  const serviceMember = selectServiceMemberFromLoggedInUser(state);
  const serviceMemberMoves = selectAllMoves(state);
  return {
    ...ownProps,
    serviceMemberId: serviceMember.id,
    serviceMemberMoves,
  };
};

const mapDispatchToProps = {
  updateAllMoves: updateAllMovesAction,
};

// in order to avoid setting up proxy server only for storybook, pass in stub function so API requests don't fail
const mergeProps = (stateProps, dispatchProps, ownProps) => ({
  ...stateProps,
  ...dispatchProps,
  ...ownProps,
});

export default connect(mapStateToProps, mapDispatchToProps, mergeProps)(Review);
