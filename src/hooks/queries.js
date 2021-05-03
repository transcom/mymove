/* eslint-disable import/prefer-default-export */
import { useQuery } from 'react-query';

import {
  getPaymentRequest,
  getMTOShipments,
  getMTOServiceItems,
  getOrder,
  getMove,
  getMoveTaskOrderList,
  getDocument,
  getMovesQueue,
  getPaymentRequestsQueue,
  getServicesCounselingQueue,
  getMovePaymentRequests,
  getCustomer,
} from 'services/ghcApi';
import { getLoggedInUserQueries } from 'services/internalApi';
import { getQueriesStatus } from 'utils/api';
import {
  PAYMENT_REQUESTS,
  MTO_SHIPMENTS,
  MTO_SERVICE_ITEMS,
  MOVES,
  ORDERS,
  MOVE_PAYMENT_REQUESTS,
  MOVE_TASK_ORDERS,
  ORDERS_DOCUMENTS,
  MOVES_QUEUE,
  PAYMENT_REQUESTS_QUEUE,
  USER,
  CUSTOMER,
  SERVICES_COUNSELING_QUEUE,
} from 'constants/queryKeys';
import { PAGINATION_PAGE_DEFAULT, PAGINATION_PAGE_SIZE_DEFAULT } from 'constants/queues';

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

export const useTXOMoveInfoQueries = (moveCode) => {
  const { data: move, ...moveQuery } = useQuery([MOVES, moveCode], getMove);
  const orderId = move?.ordersId;

  // get orders
  const { data: { orders } = {}, ...orderQuery } = useQuery([ORDERS, orderId], getOrder, {
    enabled: !!orderId,
  });

  // TODO - Need to refactor if we pass include customer in order payload
  // get customer
  const order = orders && Object.values(orders)[0];
  const customerId = order?.customerID;
  const { data: { customer } = {}, ...customerQuery } = useQuery([CUSTOMER, customerId], getCustomer, {
    enabled: !!customerId,
  });
  const customerData = customer && Object.values(customer)[0];
  const { isLoading, isError, isSuccess } = getQueriesStatus([moveQuery, orderQuery, customerQuery]);

  return {
    order,
    customerData,
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

  const { isLoading, isError, isSuccess } = getQueriesStatus([paymentRequestQuery]);

  return {
    paymentRequest,
    paymentRequests,
    paymentServiceItems,
    isLoading,
    isError,
    isSuccess,
  };
};

export const useMoveTaskOrderQueries = (moveCode) => {
  const { data: move, ...moveQuery } = useQuery([MOVES, moveCode], getMove);
  const orderId = move?.ordersId;

  // get orders
  const { data: { orders } = {}, ...orderQuery } = useQuery([ORDERS, orderId], getOrder, {
    enabled: !!orderId,
  });

  // get move task orders
  const { data: { moveTaskOrders } = {}, ...moveTaskOrderQuery } = useQuery(
    [MOVE_TASK_ORDERS, orderId],
    getMoveTaskOrderList,
    { enabled: !!orderId },
  );

  const moveTaskOrder = moveTaskOrders && Object.values(moveTaskOrders)[0];
  const mtoID = moveTaskOrder?.id;

  // get MTO shipments
  const { data: mtoShipments, ...mtoShipmentQuery } = useQuery([MTO_SHIPMENTS, mtoID, false], getMTOShipments, {
    enabled: !!mtoID,
  });

  // get MTO service items
  const { data: mtoServiceItems, ...mtoServiceItemQuery } = useQuery(
    [MTO_SERVICE_ITEMS, mtoID, false],
    getMTOServiceItems,
    { enabled: !!mtoID },
  );

  const { isLoading, isError, isSuccess } = getQueriesStatus([
    moveQuery,
    orderQuery,
    moveTaskOrderQuery,
    mtoShipmentQuery,
    mtoServiceItemQuery,
  ]);

  return {
    orders,
    moveTaskOrders,
    mtoShipments,
    mtoServiceItems,
    isLoading,
    isError,
    isSuccess,
  };
};

export const useOrdersDocumentQueries = (moveCode) => {
  // Get the orders info so we can get the uploaded_orders_id (which is a document id)
  const { data: move, ...moveQuery } = useQuery([MOVES, moveCode], getMove);

  const orderId = move?.ordersId;

  // get orders
  const { data: { orders } = {}, ...orderQuery } = useQuery([ORDERS, orderId], getOrder, {
    enabled: !!orderId,
  });

  const order = orders && orders[`${orderId}`];
  // eslint-disable-next-line camelcase
  const documentId = order?.uploaded_order_id;

  // Get a document
  // TODO - "upload" instead of "uploads" is because of the schema.js entity name. Change to "uploads"
  const staleTime = 15 * 60000; // 15 * 60000 milliseconds = 15 mins
  const cacheTime = staleTime;
  const { data: { documents, upload } = {}, ...ordersDocumentsQuery } = useQuery(
    [ORDERS_DOCUMENTS, documentId],
    getDocument,
    {
      enabled: !!documentId,
      staleTime,
      cacheTime,
      refetchOnWindowFocus: false,
    },
  );

  const { isLoading, isError, isSuccess } = getQueriesStatus([moveQuery, orderQuery, ordersDocumentsQuery]);

  return {
    move,
    orders,
    documents,
    upload,
    isLoading,
    isError,
    isSuccess,
  };
};

export const useMovesQueueQueries = ({
  sort,
  order,
  filters = [],
  currentPage = PAGINATION_PAGE_DEFAULT,
  currentPageSize = PAGINATION_PAGE_SIZE_DEFAULT,
}) => {
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

export const useServicesCounselingQueueQueries = ({
  sort,
  order,
  filters = [],
  currentPage = PAGINATION_PAGE_DEFAULT,
  currentPageSize = PAGINATION_PAGE_SIZE_DEFAULT,
}) => {
  const { data = {}, ...servicesCounselingQueueQuery } = useQuery(
    [SERVICES_COUNSELING_QUEUE, { sort, order, filters, currentPage, currentPageSize }],
    getServicesCounselingQueue,
  );
  const { isLoading, isError, isSuccess } = getQueriesStatus([servicesCounselingQueueQuery]);
  const { queueMoves, ...dataProps } = data;
  return {
    queueResult: { data: queueMoves, ...dataProps },
    isLoading,
    isError,
    isSuccess,
  };
};

export const usePaymentRequestQueueQueries = ({
  sort,
  order,
  filters = [],
  currentPage = PAGINATION_PAGE_DEFAULT,
  currentPageSize = PAGINATION_PAGE_SIZE_DEFAULT,
}) => {
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

export const useMovePaymentRequestsQueries = (moveCode) => {
  // This queries for the payment request
  const { data = [], ...movePaymentRequestsQuery } = useQuery(
    [MOVE_PAYMENT_REQUESTS, moveCode],
    getMovePaymentRequests,
  );
  const { data: move = {} } = useQuery([MOVES, moveCode], getMove);

  const mtoID = data[0]?.moveTaskOrderID || move?.id;

  const { data: mtoShipments, ...mtoShipmentQuery } = useQuery([MTO_SHIPMENTS, mtoID, false], getMTOShipments, {
    enabled: !!mtoID,
  });

  const { isLoading, isError, isSuccess } = getQueriesStatus([movePaymentRequestsQuery, mtoShipmentQuery]);

  return {
    paymentRequests: data,
    mtoShipments,
    isLoading,
    isError,
    isSuccess,
  };
};

export const useMoveDetailsQueries = (moveCode) => {
  // Get the orders info so we can get the uploaded_orders_id (which is a document id)
  const { data: move = {}, ...moveQuery } = useQuery([MOVES, moveCode], getMove);

  const moveId = move?.id;
  const orderId = move?.ordersId;

  const { data: { orders } = {}, ...orderQuery } = useQuery([ORDERS, orderId], getOrder, {
    enabled: !!orderId,
  });

  const order = Object.values(orders || {})?.[0];

  const { data: mtoShipments, ...mtoShipmentQuery } = useQuery([MTO_SHIPMENTS, moveId, false], getMTOShipments, {
    enabled: !!moveId,
  });

  // Must account for basic service items here not tied to a shipment
  const { data: mtoServiceItems, ...mtoServiceItemQuery } = useQuery(
    [MTO_SERVICE_ITEMS, moveId, false],
    getMTOServiceItems,
    { enabled: !!moveId },
  );

  const { isLoading, isError, isSuccess } = getQueriesStatus([
    moveQuery,
    orderQuery,
    mtoShipmentQuery,
    mtoServiceItemQuery,
  ]);

  return {
    move,
    order,
    mtoShipments,
    mtoServiceItems,
    isLoading,
    isError,
    isSuccess,
  };
};
