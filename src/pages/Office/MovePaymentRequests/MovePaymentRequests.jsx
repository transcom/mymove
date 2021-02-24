import React from 'react';
import { useParams } from 'react-router-dom';

import txoStyles from '../TXOMoveInfo/TXOTab.module.scss';

import PaymentRequestCard from 'components/Office/PaymentRequestCard/PaymentRequestCard';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import SomethingWentWrong from 'shared/SomethingWentWrong';
import { useMovePaymentRequestsQueries } from 'hooks/queries';
import { formatPaymentRequestAddressString } from 'utils/shipmentDisplay';

const MovePaymentRequests = () => {
  const { moveCode } = useParams();

  const { paymentRequests, mtoShipments, isLoading, isError } = useMovePaymentRequestsQueries(moveCode);

  if (isLoading) return <LoadingPlaceholder />;
  if (isError) return <SomethingWentWrong />;

  const shipmentAddresses = [];

  if (paymentRequests.length) {
    Object.values(mtoShipments).forEach((shipment) => {
      shipmentAddresses.push({
        mtoShipmentID: shipment.id,
        shipmentAddress: formatPaymentRequestAddressString(shipment.pickupAddress, shipment.destinationAddress),
      });
    });
  }

  return (
    <div className={txoStyles.tabContent}>
      <div className="grid-container-widescreen" data-testid="MovePaymentRequests">
        <h1>Payment Requests</h1>

        {paymentRequests.length ? (
          paymentRequests.map((paymentRequest) => (
            <PaymentRequestCard
              paymentRequest={paymentRequest}
              shipmentAddresses={shipmentAddresses}
              key={paymentRequest.id}
            />
          ))
        ) : (
          <div>
            <p>No payment requests have been submitted for this move yet.</p>
          </div>
        )}
      </div>
    </div>
  );
};

export default MovePaymentRequests;
