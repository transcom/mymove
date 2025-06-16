import React, { useEffect, useState } from 'react';
import { connect } from 'react-redux';
import { GridContainer, Grid, Alert } from '@trussworks/react-uswds';
import { generatePath, useNavigate, useParams } from 'react-router-dom';

import styles from './Review.module.scss';

import ConnectedFlashMessage from 'containers/FlashMessage/FlashMessage';
import ConnectedSummary from 'components/Customer/Review/Summary/Summary';
import WizardNavigation from 'components/Customer/WizardNavigation/WizardNavigation';
import MOVE_STATUSES from 'constants/moves';
import { customerRoutes } from 'constants/routes';
import 'scenes/Review/Review.css';
import { selectAllMoves, selectServiceMemberFromLoggedInUser } from 'store/entities/selectors';
import formStyles from 'styles/form.module.scss';
import { checkIfMoveIsLocked, MOVE_LOCKED_WARNING, SHIPMENT_TYPES } from 'shared/constants';
import { isPPMShipmentComplete, isBoatShipmentComplete, isMobileHomeShipmentComplete } from 'utils/shipments';
import { useTitle } from 'hooks/custom';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import { getAllMoves } from 'services/internalApi';
import { updateAllMoves as updateAllMovesAction } from 'store/entities/actions';

const Review = ({ serviceMemberId, serviceMemberMoves, updateAllMoves }) => {
  useTitle('Move review');
  const navigate = useNavigate();
  const { moveId } = useParams();
  const [isMoveLocked, setIsMoveLocked] = useState(false);
  const handleCancel = () => {
    navigate(generatePath(customerRoutes.MOVE_HOME_PATH, { moveId }));
  };

  // fetching all move data on load since this component is dependent on that data
  // this will run each time the component is loaded/accessed
  useEffect(() => {
    getAllMoves(serviceMemberId).then((response) => {
      updateAllMoves(response);
    });
  }, [updateAllMoves, serviceMemberId]);

  let mtoShipments;
  let move;

  if (serviceMemberMoves && serviceMemberMoves.currentMove && serviceMemberMoves.previousMoves) {
    // Find the move in the currentMove array
    const currentMove = serviceMemberMoves.currentMove.find((thisMove) => thisMove.id === moveId);
    // Find the move in the previousMoves array if not found in currentMove
    const previousMove = serviceMemberMoves.previousMoves.find((thisMove) => thisMove.id === moveId);
    // the move will either be in the currentMove or previousMove object
    move = currentMove || previousMove;
    if (!move.mtoShipments) {
      mtoShipments = [];
    } else {
      mtoShipments = move.mtoShipments;
    }
  }

  useEffect(() => {
    if (checkIfMoveIsLocked(move)) {
      setIsMoveLocked(true);
    }
  }, [move]);

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

  const handleNext = () => {
    const nextPath = generatePath(customerRoutes.MOVE_AGREEMENT_PATH, {
      moveId,
    });
    navigate(nextPath);
  };

  const handleAddShipment = () => {
    const addShipmentPath = generatePath(customerRoutes.SHIPMENT_SELECT_TYPE_PATH, {
      moveId,
    });
    navigate(addShipmentPath);
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
    <>
      {isMoveLocked && (
        <Alert headingLevel="h4" type="warning">
          {MOVE_LOCKED_WARNING}
        </Alert>
      )}
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
                  You are almost done setting up your move. Double&#8209;check that your information is accurate, add
                  more shipments if needed, then move on to the final step.
                </p>
              </div>
              <ConnectedSummary isMoveLocked={isMoveLocked} />
              <div className={formStyles.formActions}>
                <WizardNavigation
                  isReviewPage
                  onNextClick={handleNext}
                  onAddShipment={handleAddShipment}
                  disableNext={hasIncompleteShipment() || !mtoShipments?.length}
                  onCancelClick={handleCancel}
                  isFirstPage
                  showFinishLater
                  readOnly={!inDraftStatus || isMoveLocked}
                />
              </div>
            </div>
          </Grid>
        </Grid>
      </GridContainer>
    </>
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
