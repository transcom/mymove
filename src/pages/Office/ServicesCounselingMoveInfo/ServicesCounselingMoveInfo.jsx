import React, { Suspense, lazy, useState, useEffect } from 'react';
import { Routes, useParams, Route, Navigate, useLocation, matchPath } from 'react-router-dom';
import { useSelector } from 'react-redux';

import 'styles/office.scss';

import ServicesCounselorTabNav from 'components/Office/ServicesCounselingTabNav/ServicesCounselingTabNav';
import CustomerHeader from 'components/CustomerHeader';
import SystemError from 'components/SystemError';
import { servicesCounselingRoutes } from 'constants/routes';
import { useTXOMoveInfoQueries, useUserQueries } from 'hooks/queries';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import SomethingWentWrong from 'shared/SomethingWentWrong';
import { roleTypes } from 'constants/userRoles';
import LockedMoveBanner from 'components/LockedMoveBanner/LockedMoveBanner';
import { isBooleanFlagEnabled } from 'utils/featureFlags';

const ServicesCounselingMoveDocumentWrapper = lazy(() =>
  import('pages/Office/ServicesCounselingMoveDocumentWrapper/ServicesCounselingMoveDocumentWrapper'),
);
const ServicesCounselingMoveDetails = lazy(() =>
  import('pages/Office/ServicesCounselingMoveDetails/ServicesCounselingMoveDetails'),
);
const ServicesCounselingAddShipment = lazy(() =>
  import('pages/Office/ServicesCounselingAddShipment/ServicesCounselingAddShipment'),
);
const ServicesCounselingEditShipmentDetails = lazy(() =>
  import('pages/Office/ServicesCounselingEditShipmentDetails/ServicesCounselingEditShipmentDetails'),
);
const CustomerInfo = lazy(() => import('pages/Office/CustomerInfo/CustomerInfo'));
const MoveTaskOrder = lazy(() => import('pages/Office/MoveTaskOrder/MoveTaskOrder'));
const CustomerSupportRemarks = lazy(() => import('pages/Office/CustomerSupportRemarks/CustomerSupportRemarks'));
const MoveHistory = lazy(() => import('pages/Office/MoveHistory/MoveHistory'));
const ReviewDocuments = lazy(() => import('pages/Office/PPM/ReviewDocuments/ReviewDocuments'));
const ServicesCounselingReviewShipmentWeights = lazy(() =>
  import('pages/Office/ServicesCounselingReviewShipmentWeights/ServicesCounselingReviewShipmentWeights'),
);
const SupportingDocuments = lazy(() => import('../SupportingDocuments/SupportingDocuments'));

const ServicesCounselingMoveInfo = () => {
  const [unapprovedShipmentCount, setUnapprovedShipmentCount] = React.useState(0);
  const [unapprovedServiceItemCount, setUnapprovedServiceItemCount] = React.useState(0);
  const [excessWeightRiskCount, setExcessWeightRiskCount] = React.useState(0);
  const [unapprovedSITExtensionCount, setUnApprovedSITExtensionCount] = React.useState(0);
  const [shipmentWarnConcernCount, setShipmentWarnConcernCount] = useState(0);
  const [shipmentErrorConcernCount, setShipmentErrorConcernCount] = useState(0);
  const [missingOrdersInfoCount, setMissingOrdersInfoCount] = useState(0);
  const [infoSavedAlert, setInfoSavedAlert] = useState(null);
  const { hasRecentError, traceId } = useSelector((state) => state.interceptor);
  const [moveLockFlag, setMoveLockFlag] = useState(false);
  const [isMoveLocked, setIsMoveLocked] = useState(false);
  const onInfoSavedUpdate = (alertType) => {
    if (alertType === 'error') {
      setInfoSavedAlert({
        alertType,
        message: 'Something went wrong, and your changes were not saved. Please try again later.',
      });
    } else {
      setInfoSavedAlert({
        alertType,
        message: 'Your changes were saved.',
      });
    }
  };

  // Clear the alert when route changes
  const location = useLocation();
  const { moveCode } = useParams();
  const { move, order, customerData, isLoading, isError } = useTXOMoveInfoQueries(moveCode);
  const { data } = useUserQueries();
  const officeUserID = data?.office_user?.id;

  const [supportingDocsFF, setSupportingDocsFF] = useState(false);

  useEffect(() => {
    const fetchData = async () => {
      setSupportingDocsFF(await isBooleanFlagEnabled('manage_supporting_docs'));
    };
    fetchData();
  }, []);

  useEffect(() => {
    const fetchData = async () => {
      if (
        infoSavedAlert &&
        !matchPath(
          {
            path: servicesCounselingRoutes.BASE_MOVE_VIEW_PATH,
            end: true,
          },
          location.pathname,
        )
      ) {
        setInfoSavedAlert(null);
      }
      const lockedMoveFlag = await isBooleanFlagEnabled('move_lock');
      setMoveLockFlag(lockedMoveFlag);
      const now = new Date();
      if (officeUserID !== move?.lockedByOfficeUserID && now < new Date(move?.lockExpiresAt) && moveLockFlag) {
        setIsMoveLocked(true);
      }
    };

    fetchData();
  }, [infoSavedAlert, location, move, officeUserID, moveLockFlag]);

  const { pathname } = useLocation();
  const hideNav =
    matchPath(
      {
        path: servicesCounselingRoutes.BASE_SHIPMENT_ADD_PATH,
        end: true,
      },
      pathname,
    ) ||
    matchPath(
      {
        path: servicesCounselingRoutes.BASE_SHIPMENT_EDIT_PATH,
        end: true,
      },
      pathname,
    ) ||
    matchPath(
      {
        path: servicesCounselingRoutes.BASE_ORDERS_EDIT_PATH,
        end: true,
      },
      pathname,
    ) ||
    matchPath(
      {
        path: servicesCounselingRoutes.BASE_ALLOWANCES_EDIT_PATH,
        end: true,
      },
      pathname,
    ) ||
    matchPath(
      {
        path: servicesCounselingRoutes.BASE_SHIPMENT_REVIEW_PATH,
        end: true,
      },
      pathname,
    ) ||
    matchPath(
      {
        path: servicesCounselingRoutes.BASE_CUSTOMER_INFO_EDIT_PATH,
        end: true,
      },
      pathname,
    );

  if (isLoading) return <LoadingPlaceholder />;
  if (isError) return <SomethingWentWrong />;

  // this locked move banner will display if the current user is not the one who has it locked
  // if the current user is the one who has it locked, it will not display
  const renderLockedBanner = () => {
    const now = new Date();
    if (move?.lockedByOfficeUserID && move?.lockExpiresAt && moveLockFlag) {
      if (move?.lockedByOfficeUserID !== officeUserID && now < new Date(move?.lockExpiresAt)) {
        return (
          <LockedMoveBanner data-testid="locked-move-banner">
            This move is locked by {move?.lockedByOfficeUser?.firstName} {move?.lockedByOfficeUser?.lastName} at{' '}
            {move?.lockedByOfficeUser?.transportationOffice?.name}
          </LockedMoveBanner>
        );
      }
      return null;
    }
    return null;
  };

  return (
    <>
      <CustomerHeader
        move={move}
        order={order}
        customer={customerData}
        moveCode={moveCode}
        userRole={roleTypes.SERVICES_COUNSELOR}
      />
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
        <ServicesCounselorTabNav
          unapprovedShipmentCount={unapprovedShipmentCount}
          moveCode={moveCode}
          unapprovedServiceItemCount={unapprovedServiceItemCount}
          excessWeightRiskCount={excessWeightRiskCount}
          unapprovedSITExtensionCount={unapprovedSITExtensionCount}
          shipmentWarnConcernCount={shipmentWarnConcernCount}
          shipmentErrorConcernCount={shipmentErrorConcernCount}
          missingOrdersInfoCount={missingOrdersInfoCount}
        />
      )}

      <Suspense fallback={<LoadingPlaceholder />}>
        <Routes>
          {/* TODO - Routes not finalized, revisit */}
          <Route path={servicesCounselingRoutes.SHIPMENT_REVIEW_PATH} end element={<ReviewDocuments />} />
          <Route
            path={servicesCounselingRoutes.SHIPMENT_VIEW_DOCUMENT_PATH}
            end
            element={<ReviewDocuments readOnly />}
          />
          <Route
            path={servicesCounselingRoutes.MOVE_VIEW_PATH}
            end
            element={
              <ServicesCounselingMoveDetails
                infoSavedAlert={infoSavedAlert}
                setUnapprovedShipmentCount={setUnapprovedShipmentCount}
                setShipmentWarnConcernCount={setShipmentWarnConcernCount}
                setShipmentErrorConcernCount={setShipmentErrorConcernCount}
                shipmentWarnConcernCount={shipmentWarnConcernCount}
                shipmentErrorConcernCount={shipmentErrorConcernCount}
                missingOrdersInfoCount={missingOrdersInfoCount}
                setMissingOrdersInfoCount={setMissingOrdersInfoCount}
                isMoveLocked={isMoveLocked}
              />
            }
          />
          <Route
            key="servicesCounselingAddShipment"
            end
            path={servicesCounselingRoutes.SHIPMENT_ADD_PATH}
            element={<ServicesCounselingAddShipment />}
          />
          <Route
            path={servicesCounselingRoutes.CUSTOMER_SUPPORT_REMARKS_PATH}
            end
            element={<CustomerSupportRemarks isMoveLocked={isMoveLocked} />}
          />
          <Route
            path={servicesCounselingRoutes.MTO_PATH}
            end
            element={
              <MoveTaskOrder
                setUnapprovedShipmentCount={setUnapprovedShipmentCount}
                setUnapprovedServiceItemCount={setUnapprovedServiceItemCount}
                setExcessWeightRiskCount={setExcessWeightRiskCount}
                setUnapprovedSITExtensionCount={setUnApprovedSITExtensionCount}
                isMoveLocked={isMoveLocked}
              />
            }
          />
          <Route path={servicesCounselingRoutes.MOVE_HISTORY_PATH} end element={<MoveHistory moveCode={moveCode} />} />
          {supportingDocsFF && (
            <Route
              path={servicesCounselingRoutes.SUPPORTING_DOCUMENTS_PATH}
              end
              element={<SupportingDocuments move={move} uploads={move?.additionalDocuments?.uploads} />}
            />
          )}
          <Route
            path={servicesCounselingRoutes.ALLOWANCES_EDIT_PATH}
            end
            element={<ServicesCounselingMoveDocumentWrapper />}
          />
          <Route
            path={servicesCounselingRoutes.ORDERS_EDIT_PATH}
            end
            element={<ServicesCounselingMoveDocumentWrapper />}
          />

          {/*  WARN: MB-15562 captured this as a potential bug. An error was reported */}
          {/* that `order` was without an `id` field. Therefore this broke the */}
          {/* `CustomerInfo` component because it is expecting an `ordersId` to come */}
          {/* from the `order.id` property returned by `useTXOMoveInfoQueries`. */}
          {order.id && (
            <Route
              path={servicesCounselingRoutes.CUSTOMER_INFO_EDIT_PATH}
              end
              element={
                <CustomerInfo
                  ordersId={order.id}
                  customer={customerData}
                  isLoading={isLoading}
                  isError={isError}
                  onUpdate={onInfoSavedUpdate}
                />
              }
            />
          )}
          <Route
            path={servicesCounselingRoutes.SHIPMENT_EDIT_PATH}
            end
            element={<ServicesCounselingEditShipmentDetails onUpdate={onInfoSavedUpdate} />}
          />
          <Route
            path={servicesCounselingRoutes.SHIPMENT_ADVANCE_PATH}
            end
            element={<ServicesCounselingEditShipmentDetails onUpdate={onInfoSavedUpdate} isAdvancePage />}
          />
          <Route
            path={servicesCounselingRoutes.SHIPMENT_REVIEW_PATH}
            exact
            element={<ReviewDocuments onUpdate={onInfoSavedUpdate} />}
          />
          <Route
            path={servicesCounselingRoutes.REVIEW_SHIPMENT_WEIGHTS_PATH}
            exact
            element={<ServicesCounselingReviewShipmentWeights moveCode={moveCode} />}
          />
          {/* TODO - clarify role/tab access */}
          <Route path="/" element={<Navigate to={servicesCounselingRoutes.MOVE_VIEW_PATH} replace />} />
        </Routes>
      </Suspense>
    </>
  );
};

export default ServicesCounselingMoveInfo;
