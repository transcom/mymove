import React from 'react';
import { Link } from 'react-router-dom';
import { useQuery } from 'react-query';

// Need functions that actually make the API call(s)
import { getPaymentRequestList } from 'services/ghcApi';

const PaymentRequestIndex = () => {
  const { isLoading, isError, data } = useQuery('paymentRequestList', getPaymentRequestList);

  // These values can be used to return the loading screen or error UI
  if (isLoading || isError) return '';

  const { paymentRequests } = data;

  // eslint-disable-next-line security/detect-object-injection
  const paymentRequestsArr = Object.keys(paymentRequests).map((i) => paymentRequests[i]);

  return (
    <>
      <h1>Payment Requests</h1>
      <table>
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
