/* eslint-disable import/prefer-default-export */
import { useQuery } from 'react-query';

import {
  getPaymentRequest,
  getMTOShipments,
  getMTOServiceItems,
  getMoveOrder,
  getMoveTaskOrderList,
} from 'services/ghcApi';
import { getQueriesStatus } from 'utils/api';
import { PAYMENT_REQUESTS, MTO_SHIPMENTS, MTO_SERVICE_ITEMS, MOVE_ORDERS, MOVE_TASK_ORDERS } from 'constants/queryKeys';

export const usePaymentRequestQueries = (paymentRequestId) => {
  // get payment request by ID
  const { data: { paymentRequests, paymentServiceItems } = {}, ...paymentRequestQuery } = useQuery(
    [PAYMENT_REQUESTS, paymentRequestId],
    getPaymentRequest,
  );

  const paymentRequest = paymentRequests && paymentRequests[`${paymentRequestId}`];
  const mtoID = paymentRequest?.moveTaskOrderID;

  // get MTO shipments
  const { data: { mtoShipments } = {}, ...mtoShipmentQuery } = useQuery([MTO_SHIPMENTS, mtoID], getMTOShipments, {
    enabled: !!mtoID,
  });

  // get MTO service items
  const { data: { mtoServiceItems } = {}, ...mtoServiceItemQuery } = useQuery(
    [MTO_SERVICE_ITEMS, mtoID],
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

export const useMoveTaskOrderQueries = (moveOrderId) => {
  // get move orders
  const { data: { moveOrders } = {}, ...moveOrderQuery } = useQuery([MOVE_ORDERS, moveOrderId], getMoveOrder);

  // get move task orders
  const { data: { moveTaskOrders } = {}, ...moveTaskOrderQuery } = useQuery(
    [MOVE_TASK_ORDERS, moveOrderId],
    getMoveTaskOrderList,
  );

  const moveTaskOrder = moveTaskOrders && Object.values(moveTaskOrders)[0];
  const mtoID = moveTaskOrder?.id;

  // get MTO shipments
  const { data: { mtoShipments } = {}, ...mtoShipmentQuery } = useQuery([MTO_SHIPMENTS, mtoID], getMTOShipments, {
    enabled: !!mtoID,
  });

  // get MTO service items
  const { data: { mtoServiceItems } = {}, ...mtoServiceItemQuery } = useQuery(
    [MTO_SERVICE_ITEMS, mtoID],
    getMTOServiceItems,
    { enabled: !!mtoID },
  );

  const { isLoading, isError, isSuccess } = getQueriesStatus([
    moveOrderQuery,
    moveTaskOrderQuery,
    mtoShipmentQuery,
    mtoServiceItemQuery,
  ]);

  return {
    moveOrders,
    moveTaskOrders,
    mtoShipments,
    mtoServiceItems,
    isLoading,
    isError,
    isSuccess,
  };
};
