import React from 'react';
import { useParams } from 'react-router-dom';
import classnames from 'classnames';

import styles from './MovePaymentRequests.module.scss';

import PaymentRequestCard from 'components/Office/PaymentRequestCard/PaymentRequestCard';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import SomethingWentWrong from 'shared/SomethingWentWrong';
import { useMovePaymentRequestsQueries } from 'hooks/queries';

const MovePaymentRequests = () => {
  const { locator } = useParams();
  const { paymentRequests, isLoading, isError } = useMovePaymentRequestsQueries(locator);

  if (isLoading) return <LoadingPlaceholder />;
  if (isError) return <SomethingWentWrong />;

  return (
    <div className={classnames(styles.MovePaymentRequests, '.container')} data-testid="MovePaymentRequests">
      <h2>Payment Requests</h2>
      {paymentRequests.map((paymentRequest) => (
        <PaymentRequestCard paymentRequest={paymentRequest} key={paymentRequest.id} locator={locator} />
      ))}
    </div>
  );
};

export default MovePaymentRequests;
