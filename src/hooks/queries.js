/* eslint-disable import/prefer-default-export */
import { useQuery } from 'react-query';

import {
  getPaymentRequest,
  getMTOShipments,
  getMTOServiceItems,
  getMoveOrder,
  getMoveTaskOrderList,
  getDocument,
  getMovesQueue,
  getPaymentRequestsQueue,
} from 'services/ghcApi';
import { getQueriesStatus } from 'utils/api';
import {
  PAYMENT_REQUESTS,
  MTO_SHIPMENTS,
  MTO_SERVICE_ITEMS,
  MOVE_ORDERS,
  MOVE_TASK_ORDERS,
  ORDERS_DOCUMENTS,
  MOVES_QUEUE,
  PAYMENT_REQUESTS_QUEUE,
} from 'constants/queryKeys';

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

export const useMoveOrderQueries = (moveOrderId) => {
  // get move orders
  const { data: { moveOrders } = {}, ...moveOrderQuery } = useQuery([MOVE_ORDERS, moveOrderId], getMoveOrder);

  const { isLoading, isError, isSuccess } = getQueriesStatus([moveOrderQuery]);

  return {
    moveOrders,
    isLoading,
    isError,
    isSuccess,
  };
};

export const useOrdersDocumentQueries = (moveOrderId) => {
  // Get the orders info so we can get the uploaded_orders_id (which is a document id)
  const { data: { moveOrders } = {}, ...moveOrderQuery } = useQuery([MOVE_ORDERS, moveOrderId], getMoveOrder);

  const orders = moveOrders && moveOrders[`${moveOrderId}`];
  // eslint-disable-next-line camelcase
  const documentId = orders?.uploaded_order_id;

  // Get a document
  // TODO - "upload" instead of "uploads" is because of the schema.js entity name. Change to "uploads"
  const { data: { documents, upload } = {}, ...ordersDocumentsQuery } = useQuery(
    [ORDERS_DOCUMENTS, documentId],
    getDocument,
    {
      enabled: !!documentId,
    },
  );

  const { isLoading, isError, isSuccess } = getQueriesStatus([moveOrderQuery, ordersDocumentsQuery]);

  return {
    moveOrders,
    documents,
    upload,
    isLoading,
    isError,
    isSuccess,
  };
};

// TODO skip normalizing of schema response and cleanup
export const useMovesQueueQueries = () => {
  const { data: { queueMovesResult } = {}, ...movesQueueQuery } = useQuery([MOVES_QUEUE], getMovesQueue);

  const { isLoading, isError, isSuccess } = getQueriesStatus([movesQueueQuery]);

  return {
    queueMovesResult,
    isLoading,
    isError,
    isSuccess,
  };
};

// TODO skip normalizing of schema response and cleanup
export const usePaymentRequestQueueQueries = () => {
  const { data = {}, ...paymentRequestsQueueQuery } = useQuery([PAYMENT_REQUESTS_QUEUE], getPaymentRequestsQueue);

  const { isLoading, isError, isSuccess } = getQueriesStatus([paymentRequestsQueueQuery]);

  return {
    queuePaymentRequestsResult: data,
    isLoading,
    isError,
    isSuccess,
  };
};
