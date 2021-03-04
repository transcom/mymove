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
  getCustomer,
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
  CUSTOMER,
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

export const useTXOMoveInfoQueries = (moveCode) => {
  const { data: move, ...moveQuery } = useQuery([MOVES, moveCode], getMove);
  const orderId = move?.ordersId;

  // get move orders
  const { data: { moveOrders } = {}, ...moveOrderQuery } = useQuery([MOVE_ORDERS, orderId], getMoveOrder, {
    enabled: !!orderId,
  });

  // TODO - Need to refactor if we pass include customer in move order payload
  // get customer
  const order = moveOrders && Object.values(moveOrders)[0];
  const customerId = order?.customerID;
  const { data: { customer } = {}, ...customerQuery } = useQuery([CUSTOMER, customerId], getCustomer, {
    enabled: !!customerId,
  });
  const customerData = customer && Object.values(customer)[0];
  const { isLoading, isError, isSuccess } = getQueriesStatus([moveQuery, moveOrderQuery, customerQuery]);

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

  // get move orders
  const { data: { moveOrders } = {}, ...moveOrderQuery } = useQuery([MOVE_ORDERS, orderId], getMoveOrder, {
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
  const { data: { mtoShipments } = {}, ...mtoShipmentQuery } = useQuery([MTO_SHIPMENTS, mtoID, true], getMTOShipments, {
    enabled: !!mtoID,
  });

  // get MTO service items
  const { data: { mtoServiceItems } = {}, ...mtoServiceItemQuery } = useQuery(
    [MTO_SERVICE_ITEMS, mtoID, true],
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

export const useMoveOrderQueries = (orderId) => {
  // get move orders
  const { data: { moveOrders } = {}, ...moveOrderQuery } = useQuery([MOVE_ORDERS, orderId], getMoveOrder);

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

  const orderId = move?.ordersId;

  // get orders
  const { data: { moveOrders } = {}, ...moveOrderQuery } = useQuery([MOVE_ORDERS, orderId], getMoveOrder, {
    enabled: !!orderId,
  });

  const orders = moveOrders && moveOrders[`${orderId}`];
  // eslint-disable-next-line camelcase
  const documentId = orders?.uploaded_order_id;

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

export const useMovePaymentRequestsQueries = (moveCode) => {
  // This queries for the payment request
  const { data = [], ...movePaymentRequestsQuery } = useQuery(
    [MOVE_PAYMENT_REQUESTS, moveCode],
    getMovePaymentRequests,
  );
  const mtoID = data[0]?.moveTaskOrderID;
  const { data: mtoShipments = [], ...mtoShipmentQuery } = useQuery([MTO_SHIPMENTS, mtoID, false], getMTOShipments, {
    enabled: !!mtoID,
  });

  // This queries for the unapproved shipments count used in the Navbar
  const { data: move = {} } = useQuery([MOVES, moveCode], getMove);
  const moveId = move?.id;
  const { data: unapprovedShipments = [], ...unapprovedShipmentsQuery } = useQuery(
    [MTO_SHIPMENTS, moveId, false],
    getMTOShipments,
    {
      enabled: !!moveId,
    },
  );

  const { isLoading, isError, isSuccess } = getQueriesStatus([
    movePaymentRequestsQuery,
    mtoShipmentQuery,
    unapprovedShipmentsQuery,
  ]);

  return {
    paymentRequests: data,
    mtoShipments,
    unapprovedShipments,
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

  const { data: { moveOrders } = {}, ...moveOrderQuery } = useQuery([MOVE_ORDERS, orderId], getMoveOrder, {
    enabled: !!orderId,
  });

  const moveOrder = Object.values(moveOrders || {})?.[0];

  const { data: mtoShipments = [], ...mtoShipmentQuery } = useQuery([MTO_SHIPMENTS, moveId, false], getMTOShipments, {
    enabled: !!moveId,
  });

  // Must account for basic service items here not tied to a shipment
  const { data: mtoServiceItems = [], ...mtoServiceItemQuery } = useQuery(
    [MTO_SERVICE_ITEMS, moveId, false],
    getMTOServiceItems,
    { enabled: !!moveId },
  );

  const { isLoading, isError, isSuccess } = getQueriesStatus([
    moveQuery,
    moveOrderQuery,
    mtoShipmentQuery,
    mtoServiceItemQuery,
  ]);

  return {
    move,
    moveOrder,
    mtoShipments,
    mtoServiceItems,
    isLoading,
    isError,
    isSuccess,
  };
};
