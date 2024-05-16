import React, { lazy, Suspense, useEffect, useState } from 'react';
import { matchPath, Navigate, Route, Routes, useLocation, useParams } from 'react-router-dom';
import { useSelector } from 'react-redux';

import 'styles/office.scss';

import { permissionTypes } from 'constants/permissions';
import { qaeCSRRoutes, tioRoutes, tooRoutes } from 'constants/routes';
import TXOTabNav from 'components/Office/TXOTabNav/TXOTabNav';
import Restricted from 'components/Restricted/Restricted';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import CustomerHeader from 'components/CustomerHeader';
import SystemError from 'components/SystemError';
import { useTXOMoveInfoQueries, useUserQueries } from 'hooks/queries';
import SomethingWentWrong from 'shared/SomethingWentWrong';
import LockedMoveBanner from 'components/LockedMoveBanner/LockedMoveBanner';
import { isBooleanFlagEnabled } from 'utils/featureFlags';

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
const CustomerInfo = lazy(() => import('pages/Office/CustomerInfo/CustomerInfo'));

const TXOMoveInfo = () => {
  const [unapprovedShipmentCount, setUnapprovedShipmentCount] = React.useState(0);
  const [unapprovedServiceItemCount, setUnapprovedServiceItemCount] = React.useState(0);
  const [unapprovedSITAddressUpdateCount, setUnapprovedSITAddressUpdateCount] = React.useState(0);
  const [shipmentsWithDeliveryAddressUpdateRequestedCount, setShipmentsWithDeliveryAddressUpdateRequestedCount] =
    React.useState(0);
  const [excessWeightRiskCount, setExcessWeightRiskCount] = React.useState(0);
  const [pendingPaymentRequestCount, setPendingPaymentRequestCount] = React.useState(0);
  const [unapprovedSITExtensionCount, setUnApprovedSITExtensionCount] = React.useState(0);
  const [moveLockFlag, setMoveLockFlag] = useState(false);
  const [isMoveLocked, setIsMoveLocked] = useState(false);

  const { hasRecentError, traceId } = useSelector((state) => state.interceptor);
  const { moveCode, reportId } = useParams();
  const { pathname } = useLocation();
  const { move, order, customerData, isLoading, isError } = useTXOMoveInfoQueries(moveCode);
  const { data } = useUserQueries();
  const officeUser = data?.office_user;

  useEffect(() => {
    const fetchData = async () => {
      const lockedMoveFlag = await isBooleanFlagEnabled('move_lock');
      setMoveLockFlag(lockedMoveFlag);
      if (officeUser?.id !== move?.lockedByOfficeUserID && moveLockFlag) {
        setIsMoveLocked(true);
      }
    };

    fetchData();
  }, [move, officeUser, moveLockFlag]);

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
    ) ||
    matchPath(
      {
        path: tooRoutes.BASE_CUSTOMER_INFO_EDIT_PATH,
        end: true,
      },
      pathname,
    );

  if (isLoading) return <LoadingPlaceholder />;
  if (isError) return <SomethingWentWrong />;

  // this locked move banner will display if the current user is not the one who has it locked
  // if the current user is the one who has it locked, it will not display
  const renderLockedBanner = () => {
    if (move?.lockedByOfficeUserID && moveLockFlag) {
      if (move?.lockedByOfficeUserID !== officeUser?.id) {
        return (
          <LockedMoveBanner data-testid="locked-move-banner">
            This move is locked by {move.lockedByOfficeUser?.firstName} {move.lockedByOfficeUser?.lastName} at{' '}
            {move.lockedByOfficeUser?.transportationOffice?.name}
          </LockedMoveBanner>
        );
      }
      return null;
    }
    return null;
  };

  return (
    <>
      <CustomerHeader move={move} order={order} customer={customerData} moveCode={moveCode} />
      {renderLockedBanner()}
      {hasRecentError && (
        <SystemError>
          Something isn&apos;t working, but we&apos;re not sure what. Wait a minute and try again.
          <br />
          If that doesn&apos;t fix it, contact the{' '}
          <a href="mailto:usarmy.scott.sddc.mbx.G6-SRC-MilMove-HD@army.mil">Technical Help Desk</a>{' '}
          (usarmy.scott.sddc.mbx.G6-SRC-MilMove-HD@army.mil) and give them this code: <strong>{traceId}</strong>
        </SystemError>
      )}
      {!hideNav && (
        <TXOTabNav
          unapprovedShipmentCount={unapprovedShipmentCount}
          unapprovedServiceItemCount={unapprovedServiceItemCount}
          unapprovedSITAddressUpdateCount={unapprovedSITAddressUpdateCount}
          shipmentsWithDeliveryAddressUpdateRequestedCount={shipmentsWithDeliveryAddressUpdateRequestedCount}
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
                setUnapprovedSITAddressUpdateCount={setUnapprovedSITAddressUpdateCount}
                setShipmentsWithDeliveryAddressUpdateRequestedCount={
                  setShipmentsWithDeliveryAddressUpdateRequestedCount
                }
                setExcessWeightRiskCount={setExcessWeightRiskCount}
                setUnapprovedSITExtensionCount={setUnApprovedSITExtensionCount}
                isMoveLocked={isMoveLocked}
              />
            }
          />
          <Route path="allowances" end element={<MoveDocumentWrapper />} />
          <Route path="orders" end element={<MoveDocumentWrapper />} />
          <Route
            path="mto"
            end
            element={
              <MoveTaskOrder
                setUnapprovedShipmentCount={setUnapprovedShipmentCount}
                setUnapprovedServiceItemCount={setUnapprovedServiceItemCount}
                setUnapprovedSITAddressUpdateCount={setUnapprovedSITAddressUpdateCount}
                setExcessWeightRiskCount={setExcessWeightRiskCount}
                setUnapprovedSITExtensionCount={setUnApprovedSITExtensionCount}
                isMoveLocked={isMoveLocked}
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
                isMoveLocked={isMoveLocked}
              />
            }
          />
          <Route path="billable-weight" end element={<ReviewBillableWeight />} />
          <Route path={qaeCSRRoutes.CUSTOMER_SUPPORT_REMARKS_PATH} end element={<CustomerSupportRemarks />} />

          {/* WARN: MB-15562 captured this as a potential bug. An error was reported */}
          {/* that `order` was returned from `useTXOMoveInfoQueries` as a null value and */}
          {/* therefore broke the `EvaluationReport`, `EvaluationReports` and */}
          {/* `EvaluationViolations` components which expect to receive a `grade` */}
          {/* property from the `order.grade` lookup. */}
          {order.grade && (
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
          )}
          {order.grade && (
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
          )}
          {order.grade && (
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
          )}
          {order.grade && (
            <Route
              path={tooRoutes.CUSTOMER_INFO_EDIT_PATH}
              end
              element={
                <CustomerInfo ordersId={order.id} customer={customerData} isLoading={isLoading} isError={isError} />
              }
            />
          )}
          <Route path="history" end element={<MoveHistory moveCode={moveCode} />} />
          {/* TODO - clarify role/tab access */}
          <Route path="/" element={<Navigate to={`/moves/${moveCode}/details`} replace />} />
        </Routes>
      </Suspense>
    </>
  );
};

export default TXOMoveInfo;
