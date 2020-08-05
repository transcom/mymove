import React from 'react';
import { Link } from 'react-router-dom';
import { useQuery } from 'react-query';

import { getPaymentRequestList } from 'services/ghcApi';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import SomethingWentWrong from 'shared/SomethingWentWrong';
import { mapObjectToArray } from 'utils/api';

const PaymentRequestIndex = () => {
  const { isLoading, isError, data, error } = useQuery('paymentRequests', getPaymentRequestList);

  // These values can be used to return the loading screen or error UI
  if (isLoading) return <LoadingPlaceholder />;
  if (isError) return <SomethingWentWrong error={error} />;

  const { paymentRequests } = data;
  const paymentRequestsArr = mapObjectToArray(paymentRequests);

  return (
    <>
      <h1>Payment Requests</h1>
      <table data-testid="PaymentRequestIndex">
        <thead>
          <tr>
            <th>ID</th>
            <th>Final</th>
            <th>Rejection Reason</th>
            <th>Service Item IDs</th>
            <th>Status</th>
          </tr>
        </thead>
        <tbody>
          {paymentRequestsArr.map((pr) => (
            <tr key={pr.id}>
              <td>
                <Link to={`/moves/MOVE_CODE/payment-requests/${pr.id}`}>{pr.id}</Link>
              </td>
              <td>{`${pr.isFinal}`}</td>
              <td>{pr.rejectionReason}</td>
              <td>{pr.serviceItemIDs}</td>
              <td>{pr.status}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </>
  );
};

export default PaymentRequestIndex;
