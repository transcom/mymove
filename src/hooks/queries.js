/* eslint-disable import/prefer-default-export */
import { useQuery } from 'react-query';

import {
  getPaymentRequest,
  getMTOShipments,
  getMTOServiceItems,
  getMoveOrder,
  getMove,
  getMoveTaskOrderList,
  getDocument,
  getMovesQueue,
  getPaymentRequestsQueue,
  getMovePaymentRequests,
} from 'services/ghcApi';
import { getLoggedInUserQueries } from 'services/internalApi';
import { getQueriesStatus } from 'utils/api';
import {
  PAYMENT_REQUESTS,
  MTO_SHIPMENTS,
  MTO_SERVICE_ITEMS,
  MOVES,
  MOVE_ORDERS,
  MOVE_PAYMENT_REQUESTS,
  MOVE_TASK_ORDERS,
  ORDERS_DOCUMENTS,
  MOVES_QUEUE,
  PAYMENT_REQUESTS_QUEUE,
  USER,
} from 'constants/queryKeys';

export const useUserQueries = () => {
  const { data = {}, ...userQuery } = useQuery([USER, false], getLoggedInUserQueries);
  const { isLoading, isError, isSuccess } = getQueriesStatus([userQuery]);

  return {
    data,
    isLoading,
    isError,
    isSuccess,
  };
};

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

export const useMoveTaskOrderQueries = (moveCode) => {
  const { data: move, ...moveQuery } = useQuery([MOVES, moveCode], getMove);
  const moveOrderId = move?.ordersId;

  // get move orders
  const { data: { moveOrders } = {}, ...moveOrderQuery } = useQuery([MOVE_ORDERS, moveOrderId], getMoveOrder, {
    enabled: !!moveOrderId,
  });

  // get move task orders
  const { data: { moveTaskOrders } = {}, ...moveTaskOrderQuery } = useQuery(
    [MOVE_TASK_ORDERS, moveOrderId],
    getMoveTaskOrderList,
    { enabled: !!moveOrderId },
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
    moveQuery,
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

export const useOrdersDocumentQueries = (moveCode) => {
  // Get the orders info so we can get the uploaded_orders_id (which is a document id)
  const { data: move, ...moveQuery } = useQuery([MOVES, moveCode], getMove);

  const moveOrderId = move?.ordersId;

  // get orders
  const { data: { moveOrders } = {}, ...moveOrderQuery } = useQuery([MOVE_ORDERS, moveOrderId], getMoveOrder, {
    enabled: !!moveOrderId,
  });

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
      refetchOnWindowFocus: false,
    },
  );

  const { isLoading, isError, isSuccess } = getQueriesStatus([moveQuery, moveOrderQuery, ordersDocumentsQuery]);

  return {
    move,
    moveOrders,
    documents,
    upload,
    isLoading,
    isError,
    isSuccess,
  };
};

export const useMovesQueueQueries = ({ sort, order, filters = [], currentPage = 1, currentPageSize = 20 }) => {
  const { data = {}, ...movesQueueQuery } = useQuery(
    [MOVES_QUEUE, { sort, order, filters, currentPage, currentPageSize }],
    getMovesQueue,
  );
  const { isLoading, isError, isSuccess } = getQueriesStatus([movesQueueQuery]);
  const { queueMoves, ...dataProps } = data;
  return {
    queueResult: { data: queueMoves, ...dataProps },
    isLoading,
    isError,
    isSuccess,
  };
};

export const usePaymentRequestQueueQueries = ({ sort, order, filters = [], currentPage = 1, currentPageSize = 20 }) => {
  const { data = {}, ...paymentRequestsQueueQuery } = useQuery(
    [PAYMENT_REQUESTS_QUEUE, { sort, order, filters, currentPage, currentPageSize }],
    getPaymentRequestsQueue,
  );

  const { isLoading, isError, isSuccess } = getQueriesStatus([paymentRequestsQueueQuery]);
  const { queuePaymentRequests, ...dataProps } = data;
  return {
    queueResult: { data: queuePaymentRequests, ...dataProps },
    isLoading,
    isError,
    isSuccess,
  };
};

export const useMovePaymentRequestsQueries = (locator) => {
  const { data = {}, ...movePaymentRequestsQuery } = useQuery([MOVE_PAYMENT_REQUESTS, locator], getMovePaymentRequests);

  const { isLoading, isError, isSuccess } = getQueriesStatus([movePaymentRequestsQuery]);

  return {
    paymentRequests: data,
    isLoading,
    isError,
    isSuccess,
  };
};
