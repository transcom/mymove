import React, { useEffect, useState, useMemo } from 'react';
import { queryCache, useMutation } from 'react-query';
import { useParams, useHistory } from 'react-router-dom';
import { generatePath } from 'react-router';
import { GridContainer, Tag } from '@trussworks/react-uswds';
import { func } from 'prop-types';
import classnames from 'classnames';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';

import txoStyles from '../TXOMoveInfo/TXOTab.module.scss';
import paymentRequestStatus from '../../../constants/paymentRequestStatus';

import styles from './MovePaymentRequests.module.scss';

import { MOVES } from 'constants/queryKeys';
import { shipmentIsOverweight } from 'utils/shipmentWeights';
import { tioRoutes } from 'constants/routes';
import handleScroll from 'utils/handleScroll';
import LeftNav from 'components/LeftNav';
import PaymentRequestCard from 'components/Office/PaymentRequestCard/PaymentRequestCard';
import BillableWeightCard from 'components/Office/BillableWeight/BillableWeightCard/BillableWeightCard';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import SomethingWentWrong from 'shared/SomethingWentWrong';
import { useMovePaymentRequestsQueries } from 'hooks/queries';
import { formatPaymentRequestAddressString, getShipmentModificationType } from 'utils/shipmentDisplay';
import { shipmentStatuses } from 'constants/shipments';
import SERVICE_ITEM_STATUSES from 'constants/serviceItems';
import {
  includedStatusesForCalculatingWeights,
  useCalculatedTotalBillableWeight,
  useCalculatedWeightRequested,
} from 'hooks/custom';
import { updateMTOReviewedBillableWeights } from 'services/ghcApi';
import { milmoveLog, MILMOVE_LOG_LEVEL } from 'utils/milmoveLog';

const sectionLabels = {
  'billable-weights': 'Billable weights',
  'payment-requests': 'Payment requests',
};

const MovePaymentRequests = ({
  setUnapprovedShipmentCount,
  setUnapprovedServiceItemCount,
  setPendingPaymentRequestCount,
}) => {
  const { moveCode } = useParams();
  const history = useHistory();

  const { move, paymentRequests, order, mtoShipments, isLoading, isError } = useMovePaymentRequestsQueries(moveCode);
  const [activeSection, setActiveSection] = useState('');
  const sections = useMemo(() => {
    return ['billable-weights', 'payment-requests'];
  }, []);
  const filteredShipments = mtoShipments?.filter((shipment) => {
    return includedStatusesForCalculatingWeights(shipment.status);
  });

  const [mutateMoves] = useMutation(updateMTOReviewedBillableWeights, {
    onSuccess: (data, variables) => {
      const updatedMove = data.moves[variables.moveTaskOrderID];
      queryCache.setQueryData([MOVES, move.locator], {
        moves: {
          [`${move.locator}`]: updatedMove,
        },
      });
      queryCache.invalidateQueries([MOVES, move.locator]);
    },
    onError: (error) => {
      const errorMsg = error?.response?.body;
      milmoveLog(MILMOVE_LOG_LEVEL.LOG, errorMsg);
    },
  });

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
        if (shipment.status === shipmentStatuses.APPROVED) {
          serviceItemCount += shipment.mtoServiceItems?.filter(
            (serviceItem) => serviceItem.status === SERVICE_ITEM_STATUSES.SUBMITTED,
          ).length;
        }
      });
    }
    setUnapprovedServiceItemCount(serviceItemCount);
  }, [mtoShipments, setUnapprovedServiceItemCount]);

  useEffect(() => {
    const pendingCount = paymentRequests?.filter((pr) => pr.status === paymentRequestStatus.PENDING).length;
    setPendingPaymentRequestCount(pendingCount);
  }, [paymentRequests, setPendingPaymentRequestCount]);

  useEffect(() => {
    // attach scroll listener
    window.addEventListener('scroll', handleScroll(sections, activeSection, setActiveSection));

    // remove scroll listener
    return () => {
      window.removeEventListener('scroll', handleScroll(sections, activeSection, setActiveSection));
    };
  }, [sections, activeSection]);

  const totalBillableWeight = useCalculatedTotalBillableWeight(mtoShipments);
  const weightRequested = useCalculatedWeightRequested(mtoShipments);
  const maxBillableWeight = order?.entitlement?.authorizedWeight;
  const billableWeightsReviewed = move?.billableWeightsReviewedAt;

  if (isLoading) return <LoadingPlaceholder />;
  if (isError) return <SomethingWentWrong />;

  const shipmentsInfo = [];

  if (paymentRequests.length) {
    mtoShipments.forEach((shipment) => {
      shipmentsInfo.push({
        mtoShipmentID: shipment.id,
        address: formatPaymentRequestAddressString(shipment.pickupAddress, shipment.destinationAddress),
        departureDate: shipment.actualPickupDate,
        modificationType: getShipmentModificationType(shipment),
        mtoServiceItems: shipment.mtoServiceItems,
      });
    });
  }

  const handleReviewWeightsClick = () => {
    history.push(generatePath(tioRoutes.BILLABLE_WEIGHT_PATH, { moveCode }));
    const payload = {
      moveTaskOrderID: move?.id,
      ifMatchETag: move?.eTag,
    };
    mutateMoves(payload);
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

  const maxBillableWeightExceeded = totalBillableWeight > maxBillableWeight;
  const noBillableWeightIssues =
    (billableWeightsReviewed && !maxBillableWeightExceeded) ||
    (!anyShipmentOverweight(filteredShipments) && !anyShipmentMissingWeight(filteredShipments));
  return (
    <div className={txoStyles.tabContent}>
      <div className={txoStyles.container} data-testid="MovePaymentRequests">
        <LeftNav className={txoStyles.sidebar}>
          {sections?.map((s) => {
            return (
              <a key={`sidenav_${s}`} href={`#${s}`} className={classnames({ active: s === activeSection })}>
                {sectionLabels[`${s}`]}
                {s === 'payment-requests' && paymentRequests?.length > 0 && (
                  <Tag className={txoStyles.tag} data-testid="numOfPaymentRequestsTag">
                    {paymentRequests.length}
                  </Tag>
                )}
                {s === 'billable-weights' && maxBillableWeightExceeded && filteredShipments?.length > 0 && (
                  <Tag
                    className={classnames('usa-tag usa-tag--alert', styles.errorTag)}
                    data-testid="maxBillableWeightErrorTag"
                  >
                    <FontAwesomeIcon icon="exclamation" />
                  </Tag>
                )}
                {s === 'billable-weights' &&
                  !maxBillableWeightExceeded &&
                  filteredShipments?.length > 0 &&
                  !billableWeightsReviewed &&
                  (anyShipmentOverweight(filteredShipments) || anyShipmentMissingWeight(filteredShipments)) && (
                    <FontAwesomeIcon
                      icon="exclamation-triangle"
                      data-testid="maxBillableWeightWarningTag"
                      className={classnames(styles.warning, styles.errorTag)}
                    />
                  )}
              </a>
            );
          })}
        </LeftNav>
        <GridContainer className={txoStyles.gridContainer} data-testid="tio-payment-request-details">
          <h1>Payment requests</h1>
          <div className={txoStyles.section} id="billable-weights">
            {/* Only show shipments in statuses of approved, diversion requested, or cancellation requested */}
            <BillableWeightCard
              maxBillableWeight={maxBillableWeight}
              totalBillableWeight={totalBillableWeight}
              weightRequested={weightRequested}
              weightAllowance={order?.entitlement?.totalWeight}
              onReviewWeights={handleReviewWeightsClick}
              shipments={filteredShipments}
              secondaryReviewWeightsBtn={noBillableWeightIssues}
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
