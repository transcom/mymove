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

const MovePaymentRequests = ({ setUnapprovedShipmentCount, setPendingPaymentRequestCount }) => {
  const { moveCode } = useParams();

  const { paymentRequests, mtoShipments, isLoading, isError } = useMovePaymentRequestsQueries(moveCode);

  const mtoShipmentsArr = Object.values(mtoShipments);

  useEffect(() => {
    const shipmentCount = mtoShipments
      ? mtoShipmentsArr.filter((shipment) => shipment.status === 'SUBMITTED').length
      : 0;
    setUnapprovedShipmentCount(shipmentCount);
  }, [mtoShipments, mtoShipmentsArr, setUnapprovedShipmentCount]);

  useEffect(() => {
    const pendingCount = paymentRequests.filter((pr) => pr.status === paymentRequestStatus.PENDING).length;
    setPendingPaymentRequestCount(pendingCount);
  }, [paymentRequests, setPendingPaymentRequestCount]);

  if (isLoading) return <LoadingPlaceholder />;
  if (isError) return <SomethingWentWrong />;

  const shipmentAddresses = [];

  if (paymentRequests.length) {
    mtoShipmentsArr.forEach((shipment) => {
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
  setPendingPaymentRequestCount: func.isRequired,
};

export default MovePaymentRequests;
