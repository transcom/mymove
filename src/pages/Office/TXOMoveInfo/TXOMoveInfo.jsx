import React, { lazy, Suspense } from 'react';
import { matchPath, Navigate, Route, Routes, useLocation, useParams } from 'react-router-dom';
import { useSelector } from 'react-redux';

import 'styles/office.scss';

import { permissionTypes } from 'constants/permissions';
import { qaeCSRRoutes, tioRoutes } from 'constants/routes';
import TXOTabNav from 'components/Office/TXOTabNav/TXOTabNav';
import Restricted from 'components/Restricted/Restricted';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import CustomerHeader from 'components/CustomerHeader';
import SystemError from 'components/SystemError';
import { useTXOMoveInfoQueries } from 'hooks/queries';
import SomethingWentWrong from 'shared/SomethingWentWrong';

const MoveDetails = lazy(() => import('pages/Office/MoveDetails/MoveDetails'));
const MoveDocumentWrapper = lazy(() => import('pages/Office/MoveDocumentWrapper/MoveDocumentWrapper'));
const MoveTaskOrder = lazy(() => import('pages/Office/MoveTaskOrder/MoveTaskOrder'));
const PaymentRequestReview = lazy(() => import('pages/Office/PaymentRequestReview/PaymentRequestReview'));
const ReviewBillableWeight = lazy(() => import('pages/Office/ReviewBillableWeight/ReviewBillableWeight'));
const CustomerSupportRemarks = lazy(() => import('pages/Office/CustomerSupportRemarks/CustomerSupportRemarks'));
const EvaluationReports = lazy(() => import('pages/Office/EvaluationReports/EvaluationReports'));
const EvaluationReport = lazy(() => import('pages/Office/EvaluationReport/EvaluationReport'));
const EvaluationViolations = lazy(() => import('pages/Office/EvaluationViolations/EvaluationViolations'));
const MoveHistory = lazy(() => import('pages/Office/MoveHistory/MoveHistory'));
const MovePaymentRequests = lazy(() => import('pages/Office/MovePaymentRequests/MovePaymentRequests'));
const Forbidden = lazy(() => import('pages/Office/Forbidden/Forbidden'));

const TXOMoveInfo = () => {
  const [unapprovedShipmentCount, setUnapprovedShipmentCount] = React.useState(0);
  const [unapprovedServiceItemCount, setUnapprovedServiceItemCount] = React.useState(0);
  const [excessWeightRiskCount, setExcessWeightRiskCount] = React.useState(0);
  const [pendingPaymentRequestCount, setPendingPaymentRequestCount] = React.useState(0);
  const [unapprovedSITExtensionCount, setUnApprovedSITExtensionCount] = React.useState(0);

  const { hasRecentError, traceId } = useSelector((state) => state.interceptor);
  const { moveCode, reportId } = useParams();
  const { pathname } = useLocation();
  const { order, customerData, isLoading, isError } = useTXOMoveInfoQueries(moveCode);

  const hideNav =
    matchPath(
      {
        path: '/moves/:moveCode/payment-requests/:id',
        end: true,
      },
      pathname,
    ) ||
    matchPath(
      {
        path: '/moves/:moveCode/orders',
        end: true,
      },
      pathname,
    ) ||
    matchPath(
      {
        path: '/moves/:moveCode/allowances',
        end: true,
      },
      pathname,
    ) ||
    matchPath(
      {
        path: tioRoutes.BILLABLE_WEIGHT_PATH,
        end: true,
      },
      pathname,
    );

  if (isLoading) return <LoadingPlaceholder />;
  if (isError) return <SomethingWentWrong />;

  return (
    <>
      <CustomerHeader order={order} customer={customerData} moveCode={moveCode} />
      {hasRecentError && (
        <SystemError>
          Something isn&apos;t working, but we&apos;re not sure what. Wait a minute and try again.
          <br />
          If that doesn&apos;t fix it, contact the{' '}
          <a href="https://move.mil/customer-service#technical-help-desk">Technical Help Desk</a> and give them this
          code: <strong>{traceId}</strong>
        </SystemError>
      )}
      {!hideNav && (
        <TXOTabNav
          unapprovedShipmentCount={unapprovedShipmentCount}
          unapprovedServiceItemCount={unapprovedServiceItemCount}
          excessWeightRiskCount={excessWeightRiskCount}
          pendingPaymentRequestCount={pendingPaymentRequestCount}
          unapprovedSITExtensionCount={unapprovedSITExtensionCount}
          moveCode={moveCode}
          reportId={reportId}
          order={order}
        />
      )}

      <Suspense fallback={<LoadingPlaceholder />}>
        <Routes>
          <Route
            path="details"
            end
            element={
              <MoveDetails
                setUnapprovedShipmentCount={setUnapprovedShipmentCount}
                setUnapprovedServiceItemCount={setUnapprovedServiceItemCount}
                setExcessWeightRiskCount={setExcessWeightRiskCount}
                setUnapprovedSITExtensionCount={setUnApprovedSITExtensionCount}
              />
            }
          />
          http://officelocal:3000/moves/M4KJX4/details/allowances/orders
          <Route path="allowances" end element={<MoveDocumentWrapper />} />
          <Route path="orders" end element={<MoveDocumentWrapper />} />
          <Route
            path="mto"
            end
            element={
              <MoveTaskOrder
                setUnapprovedShipmentCount={setUnapprovedShipmentCount}
                setUnapprovedServiceItemCount={setUnapprovedServiceItemCount}
                setExcessWeightRiskCount={setExcessWeightRiskCount}
                setUnapprovedSITExtensionCount={setUnApprovedSITExtensionCount}
              />
            }
          />
          <Route path="payment-requests/:paymentRequestId" end element={<PaymentRequestReview order={order} />} />
          <Route
            path="payment-requests"
            end
            element={
              <MovePaymentRequests
                setUnapprovedShipmentCount={setUnapprovedShipmentCount}
                setUnapprovedServiceItemCount={setUnapprovedServiceItemCount}
                setPendingPaymentRequestCount={setPendingPaymentRequestCount}
              />
            }
          />
          <Route path="billable-weight" end element={<ReviewBillableWeight />} />
          <Route path={qaeCSRRoutes.CUSTOMER_SUPPORT_REMARKS_PATH} end element={<CustomerSupportRemarks />} />
          <Route
            path={qaeCSRRoutes.EVALUATION_REPORTS_PATH}
            end
            element={
              <EvaluationReports
                customerInfo={customerData}
                grade={order.grade}
                destinationDutyLocationPostalCode={order?.destinationDutyLocation?.address?.postalCode}
              />
            }
          />
          <Route
            path={qaeCSRRoutes.EVALUATION_REPORT_PATH}
            exact
            element={
              <Restricted to={permissionTypes.updateEvaluationReport} fallback={<Forbidden />}>
                <EvaluationReport
                  customerInfo={customerData}
                  grade={order.grade}
                  destinationDutyLocationPostalCode={order?.destinationDutyLocation?.address?.postalCode}
                />
              </Restricted>
            }
          />
          <Route
            path={qaeCSRRoutes.EVALUATION_VIOLATIONS_PATH}
            end
            element={
              <Restricted to={permissionTypes.updateEvaluationReport} fallback={<Forbidden />}>
                <EvaluationViolations
                  customerInfo={customerData}
                  grade={order.grade}
                  destinationDutyLocationPostalCode={order?.destinationDutyLocation?.address?.postalCode}
                />
              </Restricted>
            }
          />
          <Route path="history" end element={<MoveHistory moveCode={moveCode} />} />
          {/* TODO - clarify role/tab access */}
          <Route path="/moves/:moveCode" element={<Navigate to={`/moves/${moveCode}/details`} replace />} />
        </Routes>
      </Suspense>
    </>
  );
};

export default TXOMoveInfo;
