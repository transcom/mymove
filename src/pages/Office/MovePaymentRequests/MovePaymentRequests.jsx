import React, { useEffect, useState, useMemo } from 'react';
import { useParams, useHistory } from 'react-router-dom';
import { generatePath } from 'react-router';
import { GridContainer } from '@trussworks/react-uswds';
import { func } from 'prop-types';
import classnames from 'classnames';

import txoStyles from '../TXOMoveInfo/TXOTab.module.scss';
import paymentRequestStatus from '../../../constants/paymentRequestStatus';

import { tioRoutes } from 'constants/routes';
import handleScroll from 'utils/handleScroll';
import LeftNav from 'components/LeftNav';
import PaymentRequestCard from 'components/Office/PaymentRequestCard/PaymentRequestCard';
import BillableWeightCard from 'components/Office/BillableWeightCard/BillableWeightCard';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import SomethingWentWrong from 'shared/SomethingWentWrong';
import { useMovePaymentRequestsQueries, useMoveDetailsQueries } from 'hooks/queries';
import { formatPaymentRequestAddressString, getShipmentModificationType } from 'utils/shipmentDisplay';
import { shipmentStatuses } from 'constants/shipments';
import SERVICE_ITEM_STATUSES from 'constants/serviceItems';

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

  const { paymentRequests, mtoShipments, isLoading, isError } = useMovePaymentRequestsQueries(moveCode);
  const { order, moveIsLoading, moveIsError } = useMoveDetailsQueries(moveCode);
  const [activeSection, setActiveSection] = useState('');
  const sections = useMemo(() => {
    return ['billable-weights', 'payment-requests'];
  }, []);

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

  if (isLoading || moveIsLoading) return <LoadingPlaceholder />;
  if (isError || moveIsError) return <SomethingWentWrong />;

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
  };

  return (
    <div className={txoStyles.tabContent}>
      <div className={txoStyles.container} data-testid="MovePaymentRequests">
        <LeftNav className={txoStyles.sidebar}>
          {paymentRequests.length &&
            sections?.map((s) => {
              return (
                <a key={`sidenav_${s}`} href={`#${s}`} className={classnames({ active: s === activeSection })}>
                  {sectionLabels[`${s}`]}
                </a>
              );
            })}
        </LeftNav>
        <GridContainer className={txoStyles.gridContainer} data-testid="tio-payment-request-details">
          <h1>Payment requests</h1>
          <div className={txoStyles.section} id="billable-weights">
            {/* TODO
                totalBillableWeights needs to be calculated using serviceObject from MB-9278
                weightRequested needs to be calculated using the serviceObject from MB-9278
                weightAllowance needs to be updated to use the default weight allowance from authorizedWeight
              */}
            <BillableWeightCard
              maxBillableWeight={order?.entitlement?.authorizedWeight}
              totalBillableWeight={0}
              weightRequested={0}
              weightAllowance={order?.entitlement?.authorizedWeight}
              reviewWeights={handleReviewWeightsClick}
              shipments={mtoShipments}
            />
          </div>
          <h2>Payment requests</h2>
          <div className={txoStyles.section} id="payment-requests">
            {paymentRequests.length ? (
              paymentRequests.map((paymentRequest) => (
                <PaymentRequestCard
                  paymentRequest={paymentRequest}
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
