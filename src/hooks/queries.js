/* eslint-disable import/prefer-default-export */
import { useQueries, useQuery } from '@tanstack/react-query';
import { generatePath } from 'react-router-dom';

import { servicesCounselingRoutes } from '../constants/routes';

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
  getPrimeSimulatorAvailableMoves,
  getPPMCloseout,
  getPPMSITEstimatedCost,
  getPPMActualWeight,
  searchCustomers,
  getGBLOCs,
  getDestinationRequestsQueue,
  getBulkAssignmentData,
} from 'services/ghcApi';
import { getLoggedInUserQueries } from 'services/internalApi';
import { getPrimeSimulatorMove } from 'services/primeApi';
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
  PRIME_SIMULATOR_MOVE,
  CUSTOMER_SUPPORT_REMARKS,
  QAE_MOVE_SEARCH,
  SHIPMENT_EVALUATION_REPORTS,
  COUNSELING_EVALUATION_REPORTS,
  EVALUATION_REPORT,
  PWS_VIOLATIONS,
  REPORT_VIOLATIONS,
  MTO_SHIPMENT,
  DOCUMENTS,
  PRIME_SIMULATOR_AVAILABLE_MOVES,
  PPMCLOSEOUT,
  PPMACTUALWEIGHT,
  SC_CUSTOMER_SEARCH,
  PPMSIT_ESTIMATED_COST,
  GBLOCS,
} from 'constants/queryKeys';
import { PAGINATION_PAGE_DEFAULT, PAGINATION_PAGE_SIZE_DEFAULT } from 'constants/queues';

/**
 * Function that fetches and attaches weight tickets to corresponding ppmShipment objects on
 * each shipment in an array of MTO Shipments. This is used to incorporate the weight of PPM shipments
 * which is calculated from the net weights into various move-level weight calculations.
 *
 * @param {ShipmentShape[]} mtoShipments An array of MTO Shipments
 * @param {string} moveCode The move locator
 * @return {QueriesResults<any[]>} ppmDocsQueriesResults: an array of the documents queries for each PPM shipment in the mtoShipments array.
 */
const useAddExpensesToPPMShipments = (mtoShipments, moveCode) => {
  // Filter for ppm shipments to get their documents(including weight tickets)
  const shipmentIDs = mtoShipments?.filter((shipment) => shipment.ppmShipment).map((shipment) => shipment.id) ?? [];

  // get ppm documents
  const ppmDocsQueriesResults = useQueries({
    queries: shipmentIDs?.map((shipmentID) => {
      return {
        queryKey: [DOCUMENTS, shipmentID],
        queryFn: ({ queryKey }) => getPPMDocuments(...queryKey),
        enabled: !!shipmentID,
        select: (data) => {
          // Shove the weight tickets and other expenses into the corresponding ppmShipment object
          const shipment = mtoShipments.find((s) => s.id === shipmentID);
          shipment.ppmShipment.movingExpenses = data.MovingExpenses;
          shipment.ppmShipment.proGearWeightTickets = data.ProGearWeightTickets;
          shipment.ppmShipment.weightTickets = data.WeightTickets;
          // Attach the review url to each ppm shipment
          shipment.ppmShipment.reviewShipmentWeightsURL = generatePath(
            servicesCounselingRoutes.BASE_SHIPMENT_REVIEW_PATH,
            {
              moveCode,
              shipmentId: shipment.id,
            },
          );
        },
      };
    }),
  });
  return ppmDocsQueriesResults;
};

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
  const { isLoading, isError, isSuccess, errors } = getQueriesStatus([moveQuery, orderQuery, customerQuery]);

  return {
    move,
    order,
    customerData,
    isLoading,
    isError,
    isSuccess,
    errors,
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

export const useBulkAssignmentQueries = (queueType) => {
  const { data: bulkAssignmentData, ...bulkAssignmentDataQuery } = useQuery([queueType], ({ queryKey }) =>
    getBulkAssignmentData(queryKey),
  );
  const { isLoading, isError, isSuccess } = getQueriesStatus([bulkAssignmentDataQuery]);
  return {
    bulkAssignmentData,
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
  const {
    data: mtoShipment,
    refetch: refetchMTOShipment,
    ...mtoShipmentQuery
  } = useQuery([MTO_SHIPMENT, shipmentId], ({ queryKey }) => getMTOShipmentByID(...queryKey), {
    refetchOnMount: true,
    staleTime: 0,
  });

  const { data: documents, ...documentsQuery } = useQuery(
    [DOCUMENTS, shipmentId],
    ({ queryKey }) => getPPMDocuments(...queryKey),
    {
      enabled: !!shipmentId,
    },
  );

  const ppmShipmentId = mtoShipment?.ppmShipment?.id;
  const { data: ppmActualWeight, ...ppmActualWeightQuery } = useQuery(
    [PPMACTUALWEIGHT, ppmShipmentId],
    ({ queryKey }) => getPPMActualWeight(...queryKey),
    {
      enabled: !!ppmShipmentId,
    },
  );

  const { isLoading, isError, isSuccess, isFetching } = getQueriesStatus([
    mtoShipmentQuery,
    documentsQuery,
    ppmActualWeightQuery,
  ]);
  return {
    mtoShipment,
    documents,
    ppmActualWeight,
    refetchMTOShipment,
    isLoading,
    isError,
    isSuccess,
    isFetching,
  };
};

export const usePPMCloseoutQuery = (ppmShipmentId) => {
  const { data: ppmCloseout = {}, ...ppmCloseoutQuery } = useQuery([PPMCLOSEOUT, ppmShipmentId], ({ queryKey }) =>
    getPPMCloseout(...queryKey),
  );

  const { isLoading, isError, isSuccess, isFetching } = getQueriesStatus([ppmCloseoutQuery]);

  return {
    ppmCloseout,
    isLoading,
    isError,
    isSuccess,
    isFetching,
  };
};

export const useGetPPMSITEstimatedCostQuery = (
  ppmShipmentId,
  sitLocation,
  sitEntryDate,
  sitDepartureDate,
  weightStored,
) => {
  const { data: estimatedCost, ...ppmSITEstimatedCostQuery } = useQuery(
    [PPMSIT_ESTIMATED_COST, ppmShipmentId, sitLocation, sitEntryDate, sitDepartureDate, weightStored],
    ({ queryKey }) => getPPMSITEstimatedCost(...queryKey),
  );

  const { isLoading, isError, isSuccess } = getQueriesStatus([ppmSITEstimatedCostQuery]);

  return {
    estimatedCost,
    isLoading,
    isError,
    isSuccess,
  };
};

export const useReviewShipmentWeightsQuery = (moveCode) => {
  const { data: move, ...moveQuery } = useQuery({
    queryKey: [MOVES, moveCode],
    queryFn: ({ queryKey }) => getMove(...queryKey),
  });
  const orderId = move?.ordersId;

  // get orders
  const { data: { orders } = {}, ...orderQuery } = useQuery({
    queryKey: [ORDERS, orderId],
    queryFn: ({ queryKey }) => getOrder(...queryKey),
    options: {
      enabled: !!orderId,
    },
  });
  const mtoID = move?.id;

  // get MTO shipments
  const { data: mtoShipments, ...mtoShipmentQuery } = useQuery({
    queryKey: [MTO_SHIPMENTS, mtoID, false],
    queryFn: ({ queryKey }) => getMTOShipments(...queryKey),
    options: {
      enabled: !!mtoID,
    },
  });

  // attach ppm documents to their respective ppm shipments
  const ppmDocsQueriesResults = useAddExpensesToPPMShipments(mtoShipments, moveCode);

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

  // attach ppm documents to their respective ppm shipments
  const ppmDocsQueriesResults = useAddExpensesToPPMShipments(mtoShipments, moveCode);

  const { isLoading, isError, isSuccess } = getQueriesStatus([
    moveQuery,
    orderQuery,
    mtoShipmentQuery,
    mtoServiceItemQuery,
    ...ppmDocsQueriesResults,
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

export const useGetDocumentQuery = (documentId) => {
  const staleTime = 15 * 60000; // 15 * 60000 milliseconds = 15 mins
  const cacheTime = staleTime;
  const { data: { documents, upload } = {}, ...documentsQuery } = useQuery(
    [ORDERS_DOCUMENTS, documentId],
    ({ queryKey }) => getDocument(...queryKey),
    {
      enabled: !!documentId,
      staleTime,
      cacheTime,
      refetchOnWindowFocus: false,
    },
  );

  const { isLoading, isError, isSuccess } = getQueriesStatus([documentsQuery]);

  return {
    documents,
    upload,
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
    amendedOrderDocumentId,
  };
};

export const useAmendedDocumentQueries = (amendedOrderDocumentId) => {
  const staleTime = 0;
  const cacheTime = staleTime;

  const { data: { documents: amendedDocuments, upload: amendedUpload } = {}, ...amendedOrdersDocumentsQuery } =
    useQuery([ORDERS_DOCUMENTS, amendedOrderDocumentId], ({ queryKey }) => getDocument(...queryKey), {
      enabled: !!amendedOrderDocumentId,
      staleTime,
      cacheTime,
      refetchOnWindowFocus: false,
    });

  const { isLoading, isError, isSuccess } = getQueriesStatus([amendedOrdersDocumentsQuery]);

  return {
    amendedDocuments,
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
  viewAsGBLOC,
}) => {
  const { data = {}, ...movesQueueQuery } = useQuery(
    [MOVES_QUEUE, { sort, order, filters, currentPage, currentPageSize, viewAsGBLOC }],
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

export const useDestinationRequestsQueueQueries = ({
  sort,
  order,
  filters = [],
  currentPage = PAGINATION_PAGE_DEFAULT,
  currentPageSize = PAGINATION_PAGE_SIZE_DEFAULT,
  viewAsGBLOC,
}) => {
  const { data = {}, ...movesQueueQuery } = useQuery(
    [MOVES_QUEUE, { sort, order, filters, currentPage, currentPageSize, viewAsGBLOC }],
    ({ queryKey }) => getDestinationRequestsQueue(...queryKey),
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
  viewAsGBLOC,
}) => {
  const { data = {}, ...servicesCounselingQueueQuery } = useQuery(
    [
      SERVICES_COUNSELING_QUEUE,
      { sort, order, filters, currentPage, currentPageSize, needsPPMCloseout: true, viewAsGBLOC },
    ],
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
  viewAsGBLOC,
}) => {
  const { data = {}, ...servicesCounselingQueueQuery } = useQuery(
    [
      SERVICES_COUNSELING_QUEUE,
      { sort, order, filters, currentPage, currentPageSize, needsPPMCloseout: false, viewAsGBLOC },
    ],
    ({ queryKey }) => getServicesCounselingQueue(...queryKey),
  );

  const { isLoading, isError, isSuccess } = servicesCounselingQueueQuery;
  const { queueMoves, availableOfficeUsers, ...dataProps } = data;
  return {
    queueResult: { data: queueMoves, availableOfficeUsers, ...dataProps },
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
  viewAsGBLOC,
}) => {
  const { data = {}, ...paymentRequestsQueueQuery } = useQuery(
    [PAYMENT_REQUESTS_QUEUE, { sort, order, filters, currentPage, currentPageSize, viewAsGBLOC }],
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

  // attach all ppm documents to their respective ppm shipments
  const ppmDocsQueriesResults = useAddExpensesToPPMShipments(mtoShipments, moveCode);

  const orderId = move?.ordersId;
  const { data: { orders } = {}, ...orderQuery } = useQuery(
    [ORDERS, orderId],
    ({ queryKey }) => getOrder(...queryKey),
    {
      enabled: !!orderId,
    },
  );

  const order = Object.values(orders || {})?.[0];

  const { isLoading, isError, isSuccess } = getQueriesStatus([
    movePaymentRequestsQuery,
    mtoShipmentQuery,
    orderQuery,
    ...ppmDocsQueriesResults,
  ]);

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
  const { data: move = {}, ...moveQuery } = useQuery({
    queryKey: [MOVES, moveCode],
    queryFn: ({ queryKey }) => getMove(...queryKey),
  });

  const moveId = move?.id;
  const orderId = move?.ordersId;

  const { data: { orders } = {}, ...orderQuery } = useQuery({
    queryKey: [ORDERS, orderId],
    queryFn: ({ queryKey }) => getOrder(...queryKey),
    options: {
      enabled: !!orderId,
    },
  });

  const order = Object.values(orders || {})?.[0];

  const { upload: orderDocuments, ...documentQuery } = useGetDocumentQuery(order.uploaded_order_id);

  const { data: mtoShipments, ...mtoShipmentQuery } = useQuery({
    queryKey: [MTO_SHIPMENTS, moveId, false],
    queryFn: ({ queryKey }) => getMTOShipments(...queryKey),
    options: {
      enabled: !!moveId,
    },
  });

  // attach ppm documents to their respective ppm shipments
  const ppmDocsQueriesResults = useAddExpensesToPPMShipments(mtoShipments, moveCode);

  const customerId = order?.customerID;
  const { data: { customer } = {}, ...customerQuery } = useQuery({
    queryKey: [CUSTOMER, customerId],
    queryFn: ({ queryKey }) => getCustomer(...queryKey),
    options: {
      enabled: !!customerId,
    },
  });
  const customerData = customer && Object.values(customer)[0];
  const closeoutOffice = move.closeoutOffice && move.closeoutOffice.name;

  // Must account for basic service items here not tied to a shipment
  const { data: mtoServiceItems, ...mtoServiceItemQuery } = useQuery({
    queryKey: [MTO_SERVICE_ITEMS, moveId, false],
    queryFn: ({ queryKey }) => getMTOServiceItems(...queryKey),
    options: { enabled: !!moveId },
  });

  const { isLoading, isError, isSuccess, errors } = getQueriesStatus([
    moveQuery,
    orderQuery,
    documentQuery,
    customerQuery,
    mtoShipmentQuery,
    mtoServiceItemQuery,
    ...ppmDocsQueriesResults,
  ]);

  return {
    move,
    order,
    orderDocuments,
    customerData,
    closeoutOffice,
    mtoShipments,
    mtoServiceItems,
    isLoading,
    isError,
    isSuccess,
    errors,
  };
};

export const usePrimeSimulatorAvailableMovesQueries = ({
  filters = [],
  currentPage = PAGINATION_PAGE_DEFAULT,
  currentPageSize = PAGINATION_PAGE_SIZE_DEFAULT,
}) => {
  const { data = {}, ...primeSimulatorAvailableMovesQuery } = useQuery(
    [PRIME_SIMULATOR_AVAILABLE_MOVES, { filters, currentPage, currentPageSize }],
    ({ queryKey }) => getPrimeSimulatorAvailableMoves(...queryKey),
  );
  const { isLoading, isError, isSuccess } = primeSimulatorAvailableMovesQuery;
  const { queueMoves, ...dataProps } = data;

  return {
    queueResult: { data: queueMoves, ...dataProps },
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

  const { isLoading, isError, isSuccess, errors } = getQueriesStatus([primeSimulatorGetMoveQuery]);
  return {
    moveTaskOrder,
    isLoading,
    isError,
    isSuccess,
    errors,
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

export const useMoveSearchQueries = ({
  sort,
  order,
  filters = [],
  currentPage = PAGINATION_PAGE_DEFAULT,
  currentPageSize = PAGINATION_PAGE_SIZE_DEFAULT,
}) => {
  const queryResult = useQuery(
    [QAE_MOVE_SEARCH, { sort, order, filters, currentPage, currentPageSize }],
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

export const useCustomerSearchQueries = ({
  sort,
  order,
  filters = [],
  currentPage = PAGINATION_PAGE_DEFAULT,
  currentPageSize = PAGINATION_PAGE_SIZE_DEFAULT,
}) => {
  const queryResult = useQuery(
    [SC_CUSTOMER_SEARCH, { sort, order, filters, currentPage, currentPageSize }],
    ({ queryKey }) => searchCustomers(...queryKey),
    {
      enabled: filters.length > 0,
    },
  );
  const { data = {}, ...customerSearchQuery } = queryResult;
  const { isLoading, isError, isSuccess } = getQueriesStatus([customerSearchQuery]);
  const searchCustomersResult = data.searchCustomers;
  return {
    searchResult: { data: searchCustomersResult, page: data.page, perPage: data.perPage, totalCount: data.totalCount },
    isLoading,
    isError,
    isSuccess,
  };
};

export const useCustomerQuery = (customerId) => {
  const { data: { customer } = {}, ...customerQuery } = useQuery(
    [CUSTOMER, customerId],
    ({ queryKey }) => getCustomer(...queryKey),
    {
      enabled: !!customerId,
    },
  );
  const customerData = customer && Object.values(customer)[0];
  const { isLoading, isError, isSuccess } = getQueriesStatus([customerQuery]);
  return {
    customerData,
    isLoading,
    isError,
    isSuccess,
  };
};

export const useListGBLOCsQueries = () => {
  const { data = [], ...listGBLOCsQuery } = useQuery([GBLOCS, {}], ({ queryKey }) => getGBLOCs(...queryKey));
  const { isLoading, isError, isSuccess } = listGBLOCsQuery;
  const gblocs = data;
  return {
    result: gblocs,
    isLoading,
    isError,
    isSuccess,
  };
};
