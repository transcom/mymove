import React, { lazy, Suspense } from 'react';
import { matchPath, Redirect, Route, Switch, useLocation, useParams } from 'react-router-dom';
import { useSelector } from 'react-redux';

import 'styles/office.scss';

import { tioRoutes } from 'constants/routes';
import TXOTabNav from 'components/Office/TXOTabNav/TXOTabNav';
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
const ShipmentEvaluationReport = lazy(() => import('pages/Office/ShipmentEvaluationReport/ShipmentEvaluationReport'));
const MoveHistory = lazy(() => import('pages/Office/MoveHistory/MoveHistory'));
const MovePaymentRequests = lazy(() => import('pages/Office/MovePaymentRequests/MovePaymentRequests'));

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
    matchPath(pathname, {
      path: '/moves/:moveCode/payment-requests/:id',
      exact: true,
    }) ||
    matchPath(pathname, {
      path: '/moves/:moveCode/orders',
      exact: true,
    }) ||
    matchPath(pathname, {
      path: '/moves/:moveCode/allowances',
      exact: true,
    }) ||
    matchPath(pathname, {
      path: tioRoutes.BILLABLE_WEIGHT_PATH,
      exact: true,
    });

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
        <Switch>
          <Route path="/moves/:moveCode/details" exact>
            <MoveDetails
              setUnapprovedShipmentCount={setUnapprovedShipmentCount}
              setUnapprovedServiceItemCount={setUnapprovedServiceItemCount}
              setExcessWeightRiskCount={setExcessWeightRiskCount}
              setUnapprovedSITExtensionCount={setUnApprovedSITExtensionCount}
            />
          </Route>

          <Route path={['/moves/:moveCode/allowances', '/moves/:moveCode/orders']} exact>
            <MoveDocumentWrapper />
          </Route>

          <Route path="/moves/:moveCode/mto" exact>
            <MoveTaskOrder
              setUnapprovedShipmentCount={setUnapprovedShipmentCount}
              setUnapprovedServiceItemCount={setUnapprovedServiceItemCount}
              setExcessWeightRiskCount={setExcessWeightRiskCount}
              setUnapprovedSITExtensionCount={setUnApprovedSITExtensionCount}
            />
          </Route>

          <Route path="/moves/:moveCode/payment-requests/:paymentRequestId" exact>
            <PaymentRequestReview order={order} />
          </Route>

          <Route path="/moves/:moveCode/payment-requests" exact>
            <MovePaymentRequests
              setUnapprovedShipmentCount={setUnapprovedShipmentCount}
              setUnapprovedServiceItemCount={setUnapprovedServiceItemCount}
              setPendingPaymentRequestCount={setPendingPaymentRequestCount}
            />
          </Route>

          <Route path="/moves/:moveCode/billable-weight" exact>
            <ReviewBillableWeight />
          </Route>

          <Route path="/moves/:moveCode/customer-support-remarks" exact>
            <CustomerSupportRemarks />
          </Route>

          <Route path="/moves/:moveCode/evaluation-reports" exact>
            <EvaluationReports />
          </Route>

          <Route path="/moves/:moveCode/evaluation-reports/:reportId" exact>
            <ShipmentEvaluationReport />
          </Route>

          <Route path="/moves/:moveCode/history" exact>
            <MoveHistory moveCode={moveCode} />
          </Route>

          {/* TODO - clarify role/tab access */}
          <Redirect from="/moves/:moveCode" to="/moves/:moveCode/details" />
        </Switch>
      </Suspense>
    </>
  );
};

export default TXOMoveInfo;
