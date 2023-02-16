/* eslint-disable import/prefer-default-export */
import { useQueries, useQuery } from '@tanstack/react-query';

import {
  getPaymentRequest,
  getMTOShipments,
  getMTOServiceItems,
  getOrder,
  getMove,
  getMoveHistory,
  getDocument,
  getMovesQueue,
  getPaymentRequestsQueue,
  getPPMDocuments,
  getServicesCounselingQueue,
  getMovePaymentRequests,
  getCustomer,
  getShipmentsPaymentSITBalance,
  getCustomerSupportRemarksForMove,
  getShipmentEvaluationReports,
  getCounselingEvaluationReports,
  searchMoves,
  getEvaluationReportByID,
  getPWSViolations,
  getReportViolationsByReportID,
  getMTOShipmentByID,
  getServicesCounselingPPMQueue,
} from 'services/ghcApi';
import { getLoggedInUserQueries } from 'services/internalApi';
import { getPrimeSimulatorAvailableMoves, getPrimeSimulatorMove } from 'services/primeApi';
import { getQueriesStatus } from 'utils/api';
import {
  PAYMENT_REQUESTS,
  MTO_SHIPMENTS,
  MTO_SERVICE_ITEMS,
  MOVES,
  MOVE_HISTORY,
  ORDERS,
  MOVE_PAYMENT_REQUESTS,
  ORDERS_DOCUMENTS,
  MOVES_QUEUE,
  PAYMENT_REQUESTS_QUEUE,
  USER,
  CUSTOMER,
  SERVICES_COUNSELING_QUEUE,
  SHIPMENTS_PAYMENT_SIT_BALANCE,
  PRIME_SIMULATOR_AVAILABLE_MOVES,
  PRIME_SIMULATOR_MOVE,
  CUSTOMER_SUPPORT_REMARKS,
  QAE_CSR_MOVE_SEARCH,
  SHIPMENT_EVALUATION_REPORTS,
  COUNSELING_EVALUATION_REPORTS,
  EVALUATION_REPORT,
  PWS_VIOLATIONS,
  REPORT_VIOLATIONS,
  MTO_SHIPMENT,
  DOCUMENTS,
} from 'constants/queryKeys';
import { PAGINATION_PAGE_DEFAULT, PAGINATION_PAGE_SIZE_DEFAULT } from 'constants/queues';

export const useUserQueries = () => {
  const { data = {}, ...userQuery } = useQuery([USER, false], ({ queryKey }) => getLoggedInUserQueries(...queryKey));
  const { isLoading, isError, isSuccess } = userQuery;

  return {
    data,
    isLoading,
    isError,
    isSuccess,
  };
};

export const useTXOMoveInfoQueries = (moveCode) => {
  const { data: move, ...moveQuery } = useQuery([MOVES, moveCode], ({ queryKey }) => getMove(...queryKey));
  const orderId = move?.ordersId;

  // get orders
  const { data: { orders } = {}, ...orderQuery } = useQuery(
    [ORDERS, orderId],
    ({ queryKey }) => getOrder(...queryKey),
    {
      enabled: !!orderId,
    },
  );

  // TODO - Need to refactor if we pass include customer in order payload
  // get customer
  const order = orders && Object.values(orders)[0];
  const customerId = order?.customerID;
  const { data: { customer } = {}, ...customerQuery } = useQuery(
    [CUSTOMER, customerId],
    ({ queryKey }) => getCustomer(...queryKey),
    {
      enabled: !!customerId,
    },
  );
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
    ({ queryKey }) => getPaymentRequest(...queryKey),
  );

  const paymentRequest = paymentRequests && paymentRequests[`${paymentRequestId}`];
  const mtoID = paymentRequest?.moveTaskOrderID;

  const { data: mtoShipments, ...mtoShipmentQuery } = useQuery(
    [MTO_SHIPMENTS, mtoID, false],
    ({ queryKey }) => getMTOShipments(...queryKey),
    {
      enabled: !!mtoID,
    },
  );

  const { data: paymentSITBalances, ...shipmentsPaymentSITBalanceQuery } = useQuery(
    [SHIPMENTS_PAYMENT_SIT_BALANCE, paymentRequestId],
    ({ queryKey }) => getShipmentsPaymentSITBalance(...queryKey),
  );

  const shipmentsPaymentSITBalance = paymentSITBalances?.shipmentsPaymentSITBalance;

  const { isLoading, isError, isSuccess } = getQueriesStatus([
    paymentRequestQuery,
    mtoShipmentQuery,
    shipmentsPaymentSITBalanceQuery,
  ]);

  return {
    paymentRequest,
    paymentRequests,
    paymentServiceItems,
    mtoShipments,
    shipmentsPaymentSITBalance,
    isLoading,
    isError,
    isSuccess,
  };
};

export const useCustomerSupportRemarksQueries = (moveCode) => {
  const { data: customerSupportRemarks, ...customerSupportRemarksQuery } = useQuery(
    [CUSTOMER_SUPPORT_REMARKS, moveCode],
    ({ queryKey }) => getCustomerSupportRemarksForMove(...queryKey),
  );
  const { isLoading, isError, isSuccess } = getQueriesStatus([customerSupportRemarksQuery]);
  return {
    customerSupportRemarks,
    isLoading,
    isError,
    isSuccess,
  };
};

export const useEditShipmentQueries = (moveCode) => {
  // Get the orders info
  const { data: move = {}, ...moveQuery } = useQuery([MOVES, moveCode], ({ queryKey }) => getMove(...queryKey));

  const moveId = move?.id;
  const orderId = move?.ordersId;

  const { data: { orders } = {}, ...orderQuery } = useQuery(
    [ORDERS, orderId],
    ({ queryKey }) => getOrder(...queryKey),
    {
      enabled: !!orderId,
    },
  );

  const order = Object.values(orders || {})?.[0];

  const { data: mtoShipments, ...mtoShipmentQuery } = useQuery(
    [MTO_SHIPMENTS, moveId, false],
    ({ queryKey }) => getMTOShipments(...queryKey),
    {
      enabled: !!moveId,
    },
  );

  const { isLoading, isError, isSuccess } = getQueriesStatus([moveQuery, orderQuery, mtoShipmentQuery]);

  return {
    move,
    order,
    mtoShipments,
    isLoading,
    isError,
    isSuccess,
  };
};

export const usePPMShipmentDocsQueries = (shipmentId) => {
  const { data: mtoShipment, ...mtoShipmentQuery } = useQuery([MTO_SHIPMENT, shipmentId], ({ queryKey }) =>
    getMTOShipmentByID(...queryKey),
  );

  const { data: documents, ...documentsQuery } = useQuery(
    [DOCUMENTS, shipmentId],
    ({ queryKey }) => getPPMDocuments(...queryKey),
    {
      enabled: !!shipmentId,
    },
  );

  const { isLoading, isError, isSuccess } = getQueriesStatus([mtoShipmentQuery, documentsQuery]);
  return {
    mtoShipment,
    documents,
    isLoading,
    isError,
    isSuccess,
  };
};

export const useReviewShipmentWeightsQuery = (moveCode) => {
  const { data: move, ...moveQuery } = useQuery([MOVES, moveCode], ({ queryKey }) => getMove(...queryKey));
  const orderId = move?.ordersId;

  // get orders
  const { data: { orders } = {}, ...orderQuery } = useQuery(
    [ORDERS, orderId],
    ({ queryKey }) => getOrder(...queryKey),
    {
      enabled: !!orderId,
    },
  );
  const mtoID = move?.id;

  // get MTO shipments
  const { data: mtoShipments, ...mtoShipmentQuery } = useQuery(
    [MTO_SHIPMENTS, mtoID, false],
    ({ queryKey }) => getMTOShipments(...queryKey),
    {
      enabled: !!mtoID,
    },
  );

  // filter for ppm shipments to get their weight tickets
  const shipmentIDs = mtoShipments?.map((shipment) => shipment.id) ?? [];

  // get weight tickets
  const ppmDocsQueriesResults = useQueries({
    queries: shipmentIDs?.map((shipmentID) => {
      return {
        queryKey: [DOCUMENTS, shipmentID],
        queryFn: ({ queryKey }) => getPPMDocuments(...queryKey),
        enabled: !!shipmentID,
      };
    }),
  });

  const ppmDocs = ppmDocsQueriesResults.map((result) => result.data);

  const { isLoading, isError, isSuccess } = getQueriesStatus([
    moveQuery,
    orderQuery,
    mtoShipmentQuery,
    ...ppmDocsQueriesResults,
  ]);

  return {
    move,
    orders,
    mtoShipments,
    ppmDocs,
    isLoading,
    isError,
    isSuccess,
  };
};

export const useMoveTaskOrderQueries = (moveCode) => {
  const { data: move, ...moveQuery } = useQuery([MOVES, moveCode], ({ queryKey }) => getMove(...queryKey));
  const orderId = move?.ordersId;

  // get orders
  const { data: { orders } = {}, ...orderQuery } = useQuery(
    [ORDERS, orderId],
    ({ queryKey }) => getOrder(...queryKey),
    {
      enabled: !!orderId,
    },
  );

  const mtoID = move?.id;

  // get MTO shipments
  const { data: mtoShipments, ...mtoShipmentQuery } = useQuery(
    [MTO_SHIPMENTS, mtoID, false],
    ({ queryKey }) => getMTOShipments(...queryKey),
    {
      enabled: !!mtoID,
    },
  );

  // get MTO service items
  const { data: mtoServiceItems, ...mtoServiceItemQuery } = useQuery(
    [MTO_SERVICE_ITEMS, mtoID, false],
    ({ queryKey }) => getMTOServiceItems(...queryKey),
    { enabled: !!mtoID },
  );

  const { isLoading, isError, isSuccess } = getQueriesStatus([
    moveQuery,
    orderQuery,
    mtoShipmentQuery,
    mtoServiceItemQuery,
  ]);

  return {
    orders,
    move,
    mtoShipments,
    mtoServiceItems,
    isLoading,
    isError,
    isSuccess,
  };
};

export const useOrdersDocumentQueries = (moveCode) => {
  // Get the orders info so we can get the uploaded_orders_id (which is a document id)
  const { data: move, ...moveQuery } = useQuery([MOVES, moveCode], ({ queryKey }) => getMove(...queryKey));

  const orderId = move?.ordersId;

  // get orders
  const { data: { orders } = {}, ...orderQuery } = useQuery(
    [ORDERS, orderId],
    ({ queryKey }) => getOrder(...queryKey),
    {
      enabled: !!orderId,
    },
  );

  const order = orders && orders[`${orderId}`];
  // eslint-disable-next-line camelcase
  const documentId = order?.uploaded_order_id;
  const amendedOrderDocumentId = order?.uploadedAmendedOrderID;

  // Get a document
  // TODO - "upload" instead of "uploads" is because of the schema.js entity name. Change to "uploads"
  const staleTime = 15 * 60000; // 15 * 60000 milliseconds = 15 mins
  const cacheTime = staleTime;
  const { data: { documents, upload } = {}, ...ordersDocumentsQuery } = useQuery(
    [ORDERS_DOCUMENTS, documentId],
    ({ queryKey }) => getDocument(...queryKey),
    {
      enabled: !!documentId,
      staleTime,
      cacheTime,
      refetchOnWindowFocus: false,
    },
  );

  const { data: { documents: amendedDocuments, upload: amendedUpload } = {}, ...amendedOrdersDocumentsQuery } =
    useQuery([ORDERS_DOCUMENTS, amendedOrderDocumentId], ({ queryKey }) => getDocument(...queryKey), {
      enabled: !!amendedOrderDocumentId,
      staleTime,
      cacheTime,
      refetchOnWindowFocus: false,
    });

  const { isLoading, isError, isSuccess } = getQueriesStatus([
    moveQuery,
    orderQuery,
    ordersDocumentsQuery,
    amendedOrdersDocumentsQuery,
  ]);

  return {
    move,
    orders,
    documents,
    amendedDocuments,
    upload,
    amendedUpload,
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
    ({ queryKey }) => getMovesQueue(...queryKey),
  );
  const { isLoading, isError, isSuccess } = movesQueueQuery;
  const { queueMoves, ...dataProps } = data;
  return {
    queueResult: { data: queueMoves, ...dataProps },
    isLoading,
    isError,
    isSuccess,
  };
};

export const useServicesCounselingQueuePPMQueries = ({
  sort,
  order,
  filters = [],
  currentPage = PAGINATION_PAGE_DEFAULT,
  currentPageSize = PAGINATION_PAGE_SIZE_DEFAULT,
}) => {
  const { data = {}, ...servicesCounselingQueueQuery } = useQuery(
    [SERVICES_COUNSELING_QUEUE, { sort, order, filters, currentPage, currentPageSize, needsPPMCloseout: true }],
    ({ queryKey }) => getServicesCounselingPPMQueue(...queryKey),
  );

  const { isLoading, isError, isSuccess } = servicesCounselingQueueQuery;
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
    [SERVICES_COUNSELING_QUEUE, { sort, order, filters, currentPage, currentPageSize, needsPPMCloseout: false }],
    ({ queryKey }) => getServicesCounselingQueue(...queryKey),
  );

  const { isLoading, isError, isSuccess } = servicesCounselingQueueQuery;
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
    ({ queryKey }) => getPaymentRequestsQueue(...queryKey),
  );

  const { isLoading, isError, isSuccess } = paymentRequestsQueueQuery;
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
  const { data = [], ...movePaymentRequestsQuery } = useQuery([MOVE_PAYMENT_REQUESTS, moveCode], ({ queryKey }) =>
    getMovePaymentRequests(...queryKey),
  );
  const { data: move = {} } = useQuery([MOVES, moveCode], ({ queryKey }) => getMove(...queryKey));

  const mtoID = data[0]?.moveTaskOrderID || move?.id;

  const { data: mtoShipments, ...mtoShipmentQuery } = useQuery(
    [MTO_SHIPMENTS, mtoID, false],
    ({ queryKey }) => getMTOShipments(...queryKey),
    {
      enabled: !!mtoID,
    },
  );

  const orderId = move?.ordersId;
  const { data: { orders } = {}, ...orderQuery } = useQuery(
    [ORDERS, orderId],
    ({ queryKey }) => getOrder(...queryKey),
    {
      enabled: !!orderId,
    },
  );

  const order = Object.values(orders || {})?.[0];

  const { isLoading, isError, isSuccess } = getQueriesStatus([movePaymentRequestsQuery, mtoShipmentQuery, orderQuery]);

  return {
    paymentRequests: data,
    order,
    mtoShipments,
    isLoading,
    isError,
    isSuccess,
    move,
  };
};

// send in a single report ID and get all shipment information
export const useEvaluationReportShipmentListQueries = (reportID) => {
  const { data: evaluationReport = {}, ...viewEvaluationReportQuery } = useQuery(
    [EVALUATION_REPORT, reportID],
    ({ queryKey }) => getEvaluationReportByID(...queryKey),
  );
  const moveId = evaluationReport?.moveID;
  const { data: mtoShipments, ...mtoShipmentQuery } = useQuery(
    [MTO_SHIPMENTS, moveId, false],
    ({ queryKey }) => getMTOShipments(...queryKey),
    {
      enabled: !!moveId,
    },
  );
  const { data: reportViolations, ...reportViolationsQuery } = useQuery(
    [REPORT_VIOLATIONS, reportID],
    ({ queryKey }) => getReportViolationsByReportID(...queryKey),
    {
      enabled: !!reportID,
    },
  );
  const { isLoading, isError, isSuccess } = getQueriesStatus([
    viewEvaluationReportQuery,
    reportViolationsQuery,
    mtoShipmentQuery,
  ]);

  return {
    evaluationReport,
    mtoShipments,
    reportViolations,
    isLoading,
    isError,
    isSuccess,
  };
};

// lookup a single evaluation report, single shipment associated with that report
export const useEvaluationReportQueries = (reportID) => {
  const { data: evaluationReport = {}, ...shipmentEvaluationReportQuery } = useQuery(
    [EVALUATION_REPORT, reportID],
    getEvaluationReportByID,
  );

  const shipmentID = evaluationReport?.shipmentID;

  const { data: mtoShipment = {}, ...mtoShipmentQuery } = useQuery(
    [MTO_SHIPMENT, shipmentID],
    ({ queryKey }) => getMTOShipmentByID(...queryKey),
    {
      enabled: !!shipmentID,
    },
  );

  const { data: reportViolations = [], ...reportViolationsQuery } = useQuery(
    [REPORT_VIOLATIONS, reportID],
    ({ queryKey }) => getReportViolationsByReportID(...queryKey),
    {
      enabled: !!reportID,
    },
  );

  const { isLoading, isError, isSuccess } = getQueriesStatus([
    shipmentEvaluationReportQuery,
    mtoShipmentQuery,
    reportViolationsQuery,
  ]);
  return {
    evaluationReport,
    mtoShipment,
    reportViolations,
    isLoading,
    isError,
    isSuccess,
  };
};

// Lookup all Evaluation Reports and associated move/shipment data
export const useEvaluationReportsQueries = (moveCode) => {
  const { data: move = {}, ...moveQuery } = useQuery([MOVES, moveCode], ({ queryKey }) => getMove(...queryKey));
  const moveId = move?.id;

  const { data: shipments, ...shipmentQuery } = useQuery(
    [MTO_SHIPMENTS, moveId, false],
    ({ queryKey }) => getMTOShipments(...queryKey),
    {
      enabled: !!moveId,
    },
  );
  const { data: shipmentEvaluationReports, ...shipmentEvaluationReportsQuery } = useQuery(
    [SHIPMENT_EVALUATION_REPORTS, moveId],
    ({ queryKey }) => getShipmentEvaluationReports(...queryKey),
    {
      enabled: !!moveId,
    },
  );
  const { data: counselingEvaluationReports, ...counselingEvaluationReportsQuery } = useQuery(
    [COUNSELING_EVALUATION_REPORTS, moveId],
    ({ queryKey }) => getCounselingEvaluationReports(...queryKey),
    {
      enabled: !!moveId,
    },
  );

  const { isLoading, isError, isSuccess } = getQueriesStatus([
    moveQuery,
    shipmentQuery,
    shipmentEvaluationReportsQuery,
    counselingEvaluationReportsQuery,
  ]);
  return {
    move,
    shipments,
    counselingEvaluationReports,
    shipmentEvaluationReports,
    isLoading,
    isError,
    isSuccess,
  };
};

export const usePWSViolationsQueries = () => {
  const { data: violations = [], ...pwsViolationsQuery } = useQuery([PWS_VIOLATIONS], ({ queryKey }) =>
    getPWSViolations(...queryKey),
  );

  return {
    violations,
    ...pwsViolationsQuery,
  };
};

export const useMoveDetailsQueries = (moveCode) => {
  // Get the orders info so we can get the uploaded_orders_id (which is a document id)
  const { data: move = {}, ...moveQuery } = useQuery([MOVES, moveCode], ({ queryKey }) => getMove(...queryKey));

  const moveId = move?.id;
  const orderId = move?.ordersId;

  const { data: { orders } = {}, ...orderQuery } = useQuery(
    [ORDERS, orderId],
    ({ queryKey }) => getOrder(...queryKey),
    {
      enabled: !!orderId,
    },
  );

  const order = Object.values(orders || {})?.[0];

  const { data: mtoShipments, ...mtoShipmentQuery } = useQuery(
    [MTO_SHIPMENTS, moveId, false],
    ({ queryKey }) => getMTOShipments(...queryKey),
    {
      enabled: !!moveId,
    },
  );

  const customerId = order?.customerID;
  const { data: { customer } = {}, ...customerQuery } = useQuery(
    [CUSTOMER, customerId],
    ({ queryKey }) => getCustomer(...queryKey),
    {
      enabled: !!customerId,
    },
  );
  const customerData = customer && Object.values(customer)[0];
  const closeoutOffice = move.closeoutOffice && move.closeoutOffice.name;

  // Must account for basic service items here not tied to a shipment
  const { data: mtoServiceItems, ...mtoServiceItemQuery } = useQuery(
    [MTO_SERVICE_ITEMS, moveId, false],
    ({ queryKey }) => getMTOServiceItems(...queryKey),
    { enabled: !!moveId },
  );

  const { isLoading, isError, isSuccess } = getQueriesStatus([
    moveQuery,
    orderQuery,
    customerQuery,
    mtoShipmentQuery,
    mtoServiceItemQuery,
  ]);

  return {
    move,
    order,
    customerData,
    closeoutOffice,
    mtoShipments,
    mtoServiceItems,
    isLoading,
    isError,
    isSuccess,
  };
};

export const usePrimeSimulatorAvailableMovesQueries = () => {
  const { data = {}, ...primeSimulatorAvailableMovesQuery } = useQuery(
    [PRIME_SIMULATOR_AVAILABLE_MOVES, {}],
    ({ queryKey }) => getPrimeSimulatorAvailableMoves(...queryKey),
  );
  const { isLoading, isError, isSuccess } = getQueriesStatus([primeSimulatorAvailableMovesQuery]);
  // README: This queueResult is being artificially constructed rather than
  // created using the `..dataProp` destructering of other functions because
  // the Prime API does not return an Object that the TableQueue component can
  // consume. So the queueResult mimics that Objects properties since `data` in
  // this case is a simple Array of Prime Available Moves.
  const queueResult = {
    data,
    page: 1,
    perPage: data.length,
    totalCount: data.length,
  };

  return {
    queueResult,
    isLoading,
    isError,
    isSuccess,
  };
};

export const usePrimeSimulatorGetMove = (moveCode) => {
  const { data: moveTaskOrder, ...primeSimulatorGetMoveQuery } = useQuery(
    [PRIME_SIMULATOR_MOVE, moveCode],
    ({ queryKey }) => getPrimeSimulatorMove(...queryKey),
  );

  const { isLoading, isError, isSuccess } = getQueriesStatus([primeSimulatorGetMoveQuery]);

  return {
    moveTaskOrder,
    isLoading,
    isError,
    isSuccess,
  };
};

export const useGHCGetMoveHistory = ({
  moveCode,
  currentPage = PAGINATION_PAGE_DEFAULT,
  currentPageSize = PAGINATION_PAGE_SIZE_DEFAULT,
}) => {
  const { data = {}, ...getGHCMoveHistoryQuery } = useQuery(
    [MOVE_HISTORY, { moveCode, currentPage, currentPageSize }],
    ({ queryKey }) => getMoveHistory(...queryKey),
  );
  const { isLoading, isError, isSuccess } = getQueriesStatus([getGHCMoveHistoryQuery]);
  const { historyRecords, ...dataProps } = data;
  return {
    queueResult: { data: historyRecords, ...dataProps },
    isLoading,
    isError,
    isSuccess,
  };
};

export const useQAECSRMoveSearchQueries = ({
  sort,
  order,
  filters = [],
  currentPage = PAGINATION_PAGE_DEFAULT,
  currentPageSize = PAGINATION_PAGE_SIZE_DEFAULT,
}) => {
  const queryResult = useQuery(
    [QAE_CSR_MOVE_SEARCH, { sort, order, filters, currentPage, currentPageSize }],
    ({ queryKey }) => searchMoves(...queryKey),
    {
      enabled: filters.length > 0,
    },
  );
  const { data = {}, ...moveSearchQuery } = queryResult;
  const { isLoading, isError, isSuccess } = getQueriesStatus([moveSearchQuery]);
  const searchMovesResult = data.searchMoves;
  return {
    searchResult: { data: searchMovesResult, page: data.page, perPage: data.perPage, totalCount: data.totalCount },
    isLoading,
    isError,
    isSuccess,
  };
};
