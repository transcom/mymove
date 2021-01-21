import React from 'react';
import { useParams } from 'react-router-dom';
import classnames from 'classnames';

import styles from './MovePaymentRequests.module.scss';

import PaymentRequestCard from 'components/Office/PaymentRequestCard/PaymentRequestCard';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import SomethingWentWrong from 'shared/SomethingWentWrong';
import { useMovePaymentRequestsQueries, useMoveTaskOrderQueries } from 'hooks/queries';

const MovePaymentRequests = () => {
  const { moveCode } = useParams();
  const { paymentRequests, isLoading, isError } = useMovePaymentRequestsQueries(moveCode);

  const { mtoShipments, isLoading: mtoShipmentsLoading, isError: mtoShipmentsError } = useMoveTaskOrderQueries(
    moveCode,
  );

  if (isLoading || mtoShipmentsLoading) return <LoadingPlaceholder />;
  if (isError || mtoShipmentsError) return <SomethingWentWrong />;

  return (
    <div
      className={classnames(styles.MovePaymentRequests, 'grid-container-widescreen')}
      data-testid="MovePaymentRequests"
    >
      <h2>Payment Requests</h2>
      {paymentRequests.map((paymentRequest) => (
        <PaymentRequestCard paymentRequest={paymentRequest} mtoShipments={mtoShipments} key={paymentRequest.id} />
      ))}
    </div>
  );
};

export default MovePaymentRequests;
