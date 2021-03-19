import React, { useEffect } from 'react';
import { useParams } from 'react-router-dom';
import { func } from 'prop-types';

import txoStyles from '../TXOMoveInfo/TXOTab.module.scss';
import paymentRequestStatus from '../../../constants/paymentRequestStatus';

import PaymentRequestCard from 'components/Office/PaymentRequestCard/PaymentRequestCard';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import SomethingWentWrong from 'shared/SomethingWentWrong';
import { useMovePaymentRequestsQueries } from 'hooks/queries';
import { formatPaymentRequestAddressString } from 'utils/shipmentDisplay';
import { shipmentStatuses } from 'constants/shipments';
import SERVICE_ITEM_STATUSES from 'constants/serviceItems';

const MovePaymentRequests = ({
  setUnapprovedShipmentCount,
  setUnapprovedServiceItemCount,
  setPendingPaymentRequestCount,
}) => {
  const { moveCode } = useParams();

  const { paymentRequests, mtoShipments, isLoading, isError } = useMovePaymentRequestsQueries(moveCode);

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
    const pendingCount = paymentRequests.filter((pr) => pr.status === paymentRequestStatus.PENDING).length;
    setPendingPaymentRequestCount(pendingCount);
  }, [paymentRequests, setPendingPaymentRequestCount]);

  if (isLoading) return <LoadingPlaceholder />;
  if (isError) return <SomethingWentWrong />;

  const shipmentAddresses = [];

  if (paymentRequests.length) {
    mtoShipments.forEach((shipment) => {
      shipmentAddresses.push({
        mtoShipmentID: shipment.id,
        shipmentAddress: formatPaymentRequestAddressString(shipment.pickupAddress, shipment.destinationAddress),
      });
    });
  }

  return (
    <div className={txoStyles.tabContent}>
      <div className="grid-container-widescreen" data-testid="MovePaymentRequests">
        <h1>Payment requests</h1>

        {paymentRequests.length ? (
          paymentRequests.map((paymentRequest) => (
            <PaymentRequestCard
              paymentRequest={paymentRequest}
              shipmentAddresses={shipmentAddresses}
              key={paymentRequest.id}
            />
          ))
        ) : (
          <div className={txoStyles.emptyMessage}>
            <p>No payment requests have been submitted for this move yet.</p>
          </div>
        )}
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
