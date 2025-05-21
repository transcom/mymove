import React, { useEffect, useState, useMemo } from 'react';
import { useQueryClient, useMutation } from '@tanstack/react-query';
import { useParams, useNavigate } from 'react-router-dom';
import { Alert, Grid, GridContainer } from '@trussworks/react-uswds';
import { func } from 'prop-types';
import classnames from 'classnames';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';

import txoStyles from '../TXOMoveInfo/TXOTab.module.scss';

import styles from './MovePaymentRequests.module.scss';

import paymentRequestStatus from 'constants/paymentRequestStatus';
import { MOVES, MTO_SHIPMENTS } from 'constants/queryKeys';
import { shipmentIsOverweight } from 'utils/shipmentWeights';
import { tioRoutes } from 'constants/routes';
import PaymentRequestCard from 'components/Office/PaymentRequestCard/PaymentRequestCard';
import BillableWeightCard from 'components/Office/BillableWeight/BillableWeightCard/BillableWeightCard';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import SomethingWentWrong from 'shared/SomethingWentWrong';
import { SHIPMENT_OPTIONS, LOA_TYPE } from 'shared/constants';
import { useMovePaymentRequestsQueries } from 'hooks/queries';
import { formatPaymentRequestAddressString, getShipmentModificationType } from 'utils/shipmentDisplay';
import { shipmentStatuses, WEIGHT_ADJUSTMENT } from 'constants/shipments';
import SERVICE_ITEM_STATUSES from 'constants/serviceItems';
import {
  calculateWeightRequested,
  includedStatusesForCalculatingWeights,
  useCalculatedTotalBillableWeight,
} from 'hooks/custom';
import { updateFinancialFlag, updateMTOReviewedBillableWeights, updateMTOShipment } from 'services/ghcApi';
import { milmoveLogger } from 'utils/milmoveLog';
import FinancialReviewButton from 'components/Office/FinancialReviewButton/FinancialReviewButton';
import FinancialReviewModal from 'components/Office/FinancialReviewModal/FinancialReviewModal';
import LeftNav from 'components/LeftNav/LeftNav';
import LeftNavTag from 'components/LeftNavTag/LeftNavTag';
import Restricted from 'components/Restricted/Restricted';
import { permissionTypes } from 'constants/permissions';

const MovePaymentRequests = ({
  setUnapprovedShipmentCount,
  setUnapprovedServiceItemCount,
  setPendingPaymentRequestCount,
  isMoveLocked,
}) => {
  const { moveCode } = useParams();
  const navigate = useNavigate();

  const { move, paymentRequests, order, mtoShipments, isLoading, isError } = useMovePaymentRequestsQueries(moveCode);
  const [alertMessage, setAlertMessage] = useState(null);
  const [alertType, setAlertType] = useState('success');
  const sections = useMemo(() => {
    return ['billable-weights', 'payment-requests'];
  }, []);
  const [isFinancialModalVisible, setIsFinancialModalVisible] = useState(false);
  const filteredShipments = mtoShipments?.filter((shipment) => {
    return includedStatusesForCalculatingWeights(shipment.status);
  });
  const queryClient = useQueryClient();
  const { mutate: mutateMoves } = useMutation(updateMTOReviewedBillableWeights, {
    onSuccess: (data, variables) => {
      const updatedMove = data.moves[variables.moveTaskOrderID];
      queryClient.setQueryData([MOVES, move.locator], updatedMove);
      queryClient.invalidateQueries([MOVES, move.locator]);
    },
    onError: (error) => {
      const errorMsg = error?.response?.body;
      milmoveLogger.error(errorMsg);
    },
  });

  const { mutate: mutateFinancialReview } = useMutation(updateFinancialFlag, {
    onSuccess: (data) => {
      queryClient.setQueryData([MOVES, data.locator], data);
      queryClient.invalidateQueries([MOVES, data.locator]);
      if (data.financialReviewFlag) {
        setAlertMessage('Move flagged for financial review.');
        setAlertType('success');
      } else {
        setAlertMessage('Move unflagged for financial review.');
        setAlertType('success');
      }
    },
    onError: () => {
      setAlertMessage('There was a problem flagging the move for financial review. Please try again later.');
      setAlertType('error');
    },
  });

  const { mutate: mutateMTOhipment } = useMutation(updateMTOShipment, {
    onSuccess(_, variables) {
      queryClient.setQueryData([MTO_SHIPMENTS, variables.moveTaskOrderID, false], mtoShipments);
      queryClient.invalidateQueries([MTO_SHIPMENTS, variables.moveTaskOrderID]);
    },
  });

  const handleShowFinancialReviewModal = () => {
    setIsFinancialModalVisible(true);
  };

  const handleSubmitFinancialReviewModal = (remarks, flagForReview) => {
    // if it's set to yes let's send a true to the backend. If not we'll send false.
    const flagForReviewBool = flagForReview === 'yes';
    mutateFinancialReview({
      moveID: move.id,
      ifMatchETag: move.eTag,
      body: { remarks, flagForReview: flagForReviewBool },
    });
    setIsFinancialModalVisible(false);
  };

  const handleCancelFinancialReviewModal = () => {
    setIsFinancialModalVisible(false);
  };

  useEffect(() => {
    const shipmentCount = mtoShipments
      ? mtoShipments.filter((shipment) => shipment.status === shipmentStatuses.SUBMITTED).length
      : 0;
    setUnapprovedShipmentCount(shipmentCount);
  }, [mtoShipments, setUnapprovedShipmentCount]);

  useEffect(() => {
    let serviceItemCount = 0;
    if (mtoShipments) {
      mtoShipments.forEach((shipment) => {
        if (shipment.status === shipmentStatuses.APPROVED && shipment.mtoServiceItems) {
          serviceItemCount += shipment.mtoServiceItems.filter(
            (serviceItem) => serviceItem.status === SERVICE_ITEM_STATUSES.SUBMITTED,
          ).length;
        }
      });
    }
    setUnapprovedServiceItemCount(serviceItemCount);
  }, [mtoShipments, setUnapprovedServiceItemCount]);

  const pendingCount = paymentRequests?.filter((pr) => pr.status === paymentRequestStatus.PENDING).length;
  useEffect(() => {
    setPendingPaymentRequestCount(pendingCount);
  }, [pendingCount, setPendingPaymentRequestCount]);

  const excludePPMShipments = mtoShipments?.filter((shipment) => shipment.shipmentType !== 'PPM');

  const actualBillableWeight = useCalculatedTotalBillableWeight(excludePPMShipments, WEIGHT_ADJUSTMENT);
  const weightRequested = calculateWeightRequested(excludePPMShipments);
  const maxBillableWeight = order?.entitlement?.authorizedWeight;
  const billableWeightsReviewed = move?.billableWeightsReviewedAt;

  if (isLoading) return <LoadingPlaceholder />;
  if (isError) return <SomethingWentWrong />;

  const shipmentsInfo = [];

  if (paymentRequests.length) {
    // NOTE: We are attempting a try-catch here for the `mtoShipments`
    // iteration because it has previously been flagged as a bug in our
    // backlog, MB-15562.
    try {
      mtoShipments.forEach((shipment) => {
        const tacType = shipment.shipmentType === SHIPMENT_OPTIONS.HHG ? LOA_TYPE.HHG : shipment.tacType;
        const sacType = shipment.shipmentType === SHIPMENT_OPTIONS.HHG ? LOA_TYPE.HHG : shipment.sacType;

        shipmentsInfo.push({
          mtoShipmentID: shipment.id,
          address: formatPaymentRequestAddressString(shipment.pickupAddress, shipment.destinationAddress),
          departureDate: shipment.actualPickupDate,
          modificationType: getShipmentModificationType(shipment),
          mtoServiceItems: shipment.mtoServiceItems,
          tacType,
          sacType,
        });
      });
    } catch (error) {
      milmoveLogger.warn(
        'MovePaymentRequests.jsx: Move Task Orders should always have Shipments associated with them.',
        error,
      );
    }
  }

  const handleReviewWeightsClick = () => {
    navigate(`../${tioRoutes.BILLABLE_WEIGHT_PATH}`);
    const payload = {
      moveTaskOrderID: move?.id,
      ifMatchETag: move?.eTag,
    };
    mutateMoves(payload);
  };

  const handleEditAccountingCodes = (shipmentID, body) => {
    const shipment = mtoShipments.find((s) => s.id === shipmentID);

    if (shipment) {
      mutateMTOhipment({
        shipmentID,
        moveTaskOrderID: shipment.moveTaskOrderID,
        ifMatchETag: shipment.eTag,
        body,
      });
    }
  };

  const anyShipmentOverweight = (shipments) => {
    return shipments.some((shipment) => {
      return shipmentIsOverweight(shipment.primeEstimatedWeight, shipment.calculatedBillableWeight);
    });
  };

  const anyShipmentMissingWeight = (shipments) => {
    return shipments.some((shipment) => {
      return !shipment.primeEstimatedWeight || (shipment.reweigh?.id && !shipment.reweigh?.weight);
    });
  };

  const maxBillableWeightExceeded = actualBillableWeight > maxBillableWeight;
  const noBillableWeightIssues =
    (billableWeightsReviewed && !maxBillableWeightExceeded) ||
    (!anyShipmentOverweight(filteredShipments) && !anyShipmentMissingWeight(filteredShipments));
  return (
    <div className={txoStyles.tabContent}>
      <div className={txoStyles.container} data-testid="MovePaymentRequests">
        <LeftNav sections={sections}>
          <LeftNavTag
            background="#e34b11"
            associatedSectionName="payment-requests"
            showTag={pendingCount > 0}
            testID="numOfPendingPaymentRequestsTag"
          >
            {pendingCount}
          </LeftNavTag>
          <LeftNavTag
            className={classnames('usa-tag usa-tag--alert', styles.errorTag)}
            associatedSectionName="billable-weights"
            showTag={maxBillableWeightExceeded && filteredShipments?.length > 0}
            testID="maxBillableWeightErrorTag"
          >
            <FontAwesomeIcon icon="exclamation" />
          </LeftNavTag>
          <LeftNavTag
            className={styles.warningTag}
            background="none"
            associatedSectionName="billable-weights"
            showTag={
              !maxBillableWeightExceeded &&
              filteredShipments?.length > 0 &&
              !billableWeightsReviewed &&
              (anyShipmentOverweight(filteredShipments) || anyShipmentMissingWeight(filteredShipments))
            }
            testID="maxBillableWeightWarningTag"
          >
            <FontAwesomeIcon icon="exclamation-triangle" className={classnames(styles.warning, styles.errorTag)} />
          </LeftNavTag>
        </LeftNav>
        <GridContainer className={txoStyles.gridContainer} data-testid="tio-payment-request-details">
          <Grid row className={txoStyles.pageHeader}>
            {alertMessage && (
              <Grid col={12} className={txoStyles.alertContainer}>
                <Alert headingLevel="h4" slim type={alertType}>
                  {alertMessage}
                </Alert>
              </Grid>
            )}
          </Grid>
          <div className={styles.tioPaymentRequestsHeadingFlexbox}>
            <h1>Payment Requests</h1>
            <Restricted to={permissionTypes.updateFinancialReviewFlag}>
              <FinancialReviewButton
                onClick={handleShowFinancialReviewModal}
                reviewRequested={move?.financialReviewFlag}
                isMoveLocked={isMoveLocked}
              />
            </Restricted>
          </div>
          {isFinancialModalVisible && (
            <FinancialReviewModal
              onClose={handleCancelFinancialReviewModal}
              onSubmit={handleSubmitFinancialReviewModal}
              initialRemarks={move?.financialReviewRemarks}
              initialSelection={move?.financialReviewFlag}
            />
          )}
          <div className={txoStyles.section} id="billable-weights">
            {/* Only show shipments in statuses of approved, diversion requested, or cancellation requested */}
            <BillableWeightCard
              maxBillableWeight={maxBillableWeight}
              actualBillableWeight={actualBillableWeight}
              weightRequested={weightRequested}
              weightAllowance={order?.entitlement?.totalWeight}
              onReviewWeights={handleReviewWeightsClick}
              shipments={filteredShipments}
              secondaryReviewWeightsBtn={noBillableWeightIssues}
              isMoveLocked={isMoveLocked}
            />
          </div>
          <h2>Payment requests</h2>
          <div className={txoStyles.section} id="payment-requests">
            {paymentRequests?.length > 0 ? (
              paymentRequests.map((paymentRequest) => (
                <PaymentRequestCard
                  paymentRequest={paymentRequest}
                  hasBillableWeightIssues={!noBillableWeightIssues}
                  shipmentsInfo={shipmentsInfo}
                  key={paymentRequest.id}
                  onEditAccountingCodes={handleEditAccountingCodes}
                  isMoveLocked={isMoveLocked}
                />
              ))
            ) : (
              <div className={txoStyles.emptyMessage}>
                <p>No payment requests have been submitted for this move yet.</p>
              </div>
            )}
          </div>
        </GridContainer>
      </div>
    </div>
  );
};

MovePaymentRequests.propTypes = {
  setUnapprovedShipmentCount: func.isRequired,
  setUnapprovedServiceItemCount: func.isRequired,
  setPendingPaymentRequestCount: func.isRequired,
};

export default MovePaymentRequests;
