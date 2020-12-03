import React from 'react';
import { useParams } from 'react-router-dom';
import classnames from 'classnames';

import styles from './MovePaymentRequests.module.scss';

import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import SomethingWentWrong from 'shared/SomethingWentWrong';
import { useMovePaymentRequestsQueries } from 'hooks/queries';

const MovePaymentRequests = () => {
  const { moveOrderId } = useParams();

  // eslint-disable-next-line no-unused-vars
  const { paymentRequests, isLoading, isError } = useMovePaymentRequestsQueries(moveOrderId);

  if (isLoading) return <LoadingPlaceholder />;
  if (isError) return <SomethingWentWrong />;

  return (
    <div className={classnames(styles.MovePaymentRequests, '.container')}>
      <h2>Payment Requests</h2>
    </div>
  );
};

export default MovePaymentRequests;
