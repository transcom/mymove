import React from 'react';
import { useParams } from 'react-router-dom';
import classnames from 'classnames';

import styles from './MovePaymentRequests.module.scss';

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

  Object.values(mtoShipments).forEach((shipment) => {
    shipmentAddresses.push({
      mtoShipmentID: shipment.id,
      shipmentAddress: formatPaymentRequestAddressString(shipment.pickupAddress, shipment.destinationAddress),
    });
  });

  return (
    <div
      className={classnames(styles.MovePaymentRequests, 'grid-container-widescreen')}
      data-testid="MovePaymentRequests"
    >
      <h2>Payment Requests</h2>
      {paymentRequests.map((paymentRequest) => (
        <PaymentRequestCard
          paymentRequest={paymentRequest}
          shipmentAddresses={shipmentAddresses}
          key={paymentRequest.id}
        />
      ))}
    </div>
  );
};

export default MovePaymentRequests;
