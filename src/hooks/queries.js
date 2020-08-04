/* eslint-disable import/prefer-default-export */
import { useQuery } from 'react-query';

import { getPaymentRequest, getMTOShipments, getMTOServiceItems } from 'services/ghcApi';
import { getQueriesStatus } from 'utils/api';

export const usePaymentRequestQueries = (paymentRequestId) => {
  // get payment request by ID
  const { data: { paymentRequests, paymentServiceItems } = {}, ...paymentRequestQuery } = useQuery(
    ['paymentRequest', paymentRequestId],
    getPaymentRequest,
  );

  const paymentRequest = paymentRequests && paymentRequests[`${paymentRequestId}`];
  const mtoID = paymentRequest?.moveTaskOrderID;

  // get MTO shipments
  const { data: { mtoShipments } = {}, ...mtoShipmentQuery } = useQuery(['mtoShipment', mtoID], getMTOShipments, {
    enabled: !!mtoID,
  });

  // get MTO service items
  const { data: { mtoServiceItems } = {}, ...mtoServiceItemQuery } = useQuery(
    ['mtoServiceItem', mtoID],
    getMTOServiceItems,
    {
      enabled: !!mtoID,
    },
  );

  const { isLoading, isError, isSuccess } = getQueriesStatus([
    paymentRequestQuery,
    mtoShipmentQuery,
    mtoServiceItemQuery,
  ]);

  return {
    paymentRequest,
    paymentRequests,
    paymentServiceItems,
    mtoShipments,
    mtoServiceItems,
    isLoading,
    isError,
    isSuccess,
  };
};
