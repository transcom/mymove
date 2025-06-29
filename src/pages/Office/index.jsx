import React, { lazy, Suspense, useEffect, useState } from 'react';
import PropTypes from 'prop-types';
import { Route, Routes, Link, matchPath, Navigate, useLocation } from 'react-router-dom';
import { connect } from 'react-redux';
import classnames from 'classnames';

import styles from './Office.module.scss';

import 'styles/full_uswds.scss';
import 'scenes/Office/office.scss';
// API / Redux actions
import {
  selectGetCurrentUserIsLoading,
  selectIsLoggedIn,
  selectUnderMaintenance,
  selectLoadingSpinnerMessage,
  selectShowLoadingSpinner,
} from 'store/auth/selectors';
import { loadUser as loadUserAction } from 'store/auth/actions';
import { selectLoggedInUser } from 'store/entities/selectors';
import {
  loadInternalSchema as loadInternalSchemaAction,
  loadPublicSchema as loadPublicSchemaAction,
} from 'shared/Swagger/ducks';
// Feature Flags
import { isBooleanFlagEnabled } from 'utils/featureFlags';
// Shared layout components
import ConnectedLogoutOnInactivity from 'layout/LogoutOnInactivity';
import PrivateRoute from 'containers/PrivateRoute';
import CUIHeader from 'components/CUIHeader/CUIHeader';
import BypassBlock from 'components/BypassBlock';
import SystemError from 'components/SystemError';
import NotFound from 'components/NotFound/NotFound';
import OfficeLoggedInHeader from 'containers/Headers/OfficeLoggedInHeader';
import LoggedOutHeader from 'containers/Headers/LoggedOutHeader';
import { ConnectedSelectApplication } from 'pages/SelectApplication/SelectApplication';
import { roleTypes } from 'constants/userRoles';
import { pageNames } from 'constants/signInPageNames';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import { withContext } from 'shared/AppContext';
import { RouterShape, UserRolesShape } from 'types/index';
import {
  servicesCounselingRoutes,
  primeSimulatorRoutes,
  tooRoutes,
  qaeCSRRoutes,
  contractingOfficerRoutes,
  officeRoutes,
} from 'constants/routes';
import PrimeBanner from 'pages/PrimeUI/PrimeBanner/PrimeBanner';
import PermissionProvider from 'components/Restricted/PermissionProvider';
import withRouter from 'utils/routing';
import { OktaLoggedOutBanner, OktaNeedsLoggedOutBanner } from 'components/OktaLogoutBanner';
import SelectedGblocProvider from 'components/Office/GblocSwitcher/SelectedGblocProvider';
import MaintenancePage from 'pages/Maintenance/MaintenancePage';
import { FEATURE_FLAG_KEYS } from 'shared/constants';
import LoadingSpinner from 'components/LoadingSpinner/LoadingSpinner';

// Lazy load these dependencies (they correspond to unique routes & only need to be loaded when that URL is accessed)
const SignIn = lazy(() => import('pages/SignIn/SignIn'));
const RequestAccount = lazy(() => import('pages/Office/RequestAccount/RequestAccount'));
const InvalidPermissions = lazy(() => import('pages/InvalidPermissions/InvalidPermissions'));
// TXO
const TXOMoveInfo = lazy(() => import('pages/Office/TXOMoveInfo/TXOMoveInfo'));
// TOO pages
const MoveQueue = lazy(() => import('pages/Office/MoveQueue/MoveQueue'));
// TIO pages
const PaymentRequestQueue = lazy(() => import('pages/Office/PaymentRequestQueue/PaymentRequestQueue'));
// HQ pages
const HeadquartersQueues = lazy(() => import('pages/Office/HeadquartersQueues/HeadquartersQueues'));
// Services Counselor pages
const ServicesCounselingMoveInfo = lazy(() =>
  import('pages/Office/ServicesCounselingMoveInfo/ServicesCounselingMoveInfo'),
);
const ServicesCounselingQueue = lazy(() => import('pages/Office/ServicesCounselingQueue/ServicesCounselingQueue'));
const ServicesCounselingAddShipment = lazy(() =>
  import('pages/Office/ServicesCounselingAddShipment/ServicesCounselingAddShipment'),
);
const AddShipment = lazy(() => import('pages/Office/AddShipment/AddShipment'));
const EditShipmentDetails = lazy(() => import('pages/Office/EditShipmentDetails/EditShipmentDetails'));
const PrimeSimulatorAvailableMoves = lazy(() => import('pages/PrimeUI/AvailableMoves/AvailableMovesQueue'));
const PrimeSimulatorMoveDetails = lazy(() => import('pages/PrimeUI/MoveTaskOrder/MoveDetails'));
const PrimeSimulatorCreatePaymentRequest = lazy(() =>
  import('pages/PrimeUI/CreatePaymentRequest/CreatePaymentRequest'),
);
const PrimeUIShipmentCreateForm = lazy(() => import('pages/PrimeUI/Shipment/PrimeUIShipmentCreate'));
const PrimeUIShipmentForm = lazy(() => import('pages/PrimeUI/Shipment/PrimeUIShipmentUpdate'));

const PrimeSimulatorUploadPaymentRequestDocuments = lazy(() =>
  import('pages/PrimeUI/UploadPaymentRequestDocuments/UploadPaymentRequestDocuments'),
);
const PrimeSimulatorUploadServiceRequestDocuments = lazy(() =>
  import('pages/PrimeUI/UploadServiceRequestDocuments/UploadServiceRequestDocuments'),
);
const PrimeSimulatorCreateServiceItem = lazy(() => import('pages/PrimeUI/CreateServiceItem/CreateServiceItem'));
const PrimeSimulatorUpdateServiceItem = lazy(() => import('pages/PrimeUI/UpdateServiceItems/PrimeUIUpdateServiceItem'));
const PrimeUIShipmentUpdateAddress = lazy(() => import('pages/PrimeUI/Shipment/PrimeUIShipmentUpdateAddress'));
const PrimeUIShipmentUpdateReweigh = lazy(() => import('pages/PrimeUI/Shipment/PrimeUIShipmentUpdateReweigh'));
const PrimeSimulatorCreateSITExtensionRequest = lazy(() =>
  import('pages/PrimeUI/CreateSITExtensionRequest/CreateSITExtensionRequest'),
);
const PrimeUIShipmentUpdateDestinationAddress = lazy(() =>
  import('pages/PrimeUI/Shipment/PrimeUIShipmentUpdateDestinationAddress'),
);
const PrimeUIAcknowledgeMove = lazy(() => import('pages/PrimeUI/MoveTaskOrder/AcknowledgeMove'));
const PrimeUIAcknowledgeShipment = lazy(() => import('pages/PrimeUI/Shipment/AcknowledgeShipment'));

const MoveSearch = lazy(() => import('pages/Office/MoveSearch/MoveSearch'));
const CreateCustomerForm = lazy(() => import('pages/Office/CustomerOnboarding/CreateCustomerForm'));
const CreateMoveCustomerInfo = lazy(() => import('pages/Office/CreateMoveCustomerInfo/CreateMoveCustomerInfo'));
const CustomerInfo = lazy(() => import('pages/Office/CustomerInfo/CustomerInfo'));
const ServicesCounselingAddOrders = lazy(() =>
  import('pages/Office/ServicesCounselingAddOrders/ServicesCounselingAddOrders'),
);
const OfficeUserProfile = lazy(() => import('pages/Office/Profile/Profile'));
const OfficeUserEditContactInfo = lazy(() => import('pages/Office/Profile/ContactInfo'));

const OfficeApp = ({ loadUser, loadInternalSchema, loadPublicSchema, ...props }) => {
  // Local state for feature flags and Okta booleans
  const [oktaLoggedOut, setOktaLoggedOut] = useState(false);
  const [oktaNeedsLoggedOut, setOktaNeedsLoggedOut] = useState(false);
  const [gsrRoleFlag, setGsrRoleFlag] = useState(false);
  const [queueManagementFlag, setQueueManagementFlag] = useState(false);
  const [bulkAssignmentFlag, setBulkAssignmentFlag] = useState(false);
  const [approvalRequestTypeFlag, setApprovalRequestTypeFlag] = useState(false);

  const location = useLocation();
  const displayChangeRole =
    props.userIsLoggedIn &&
    !!props.userInactiveRoles?.length &&
    !matchPath(
      {
        path: '/select-application',
        end: true,
      },
      location.pathname,
    );
  const isFullscreenPage = matchPath(
    {
      path: '/moves/:moveCode/payment-requests/:id',
    },
    location.pathname,
  );
  const siteClasses = classnames('site', {
    [`site--fullscreen`]: isFullscreenPage,
  });

  useEffect(() => {
    loadInternalSchema();
    loadPublicSchema();
    loadUser();

    isBooleanFlagEnabled('gsr_role').then(setGsrRoleFlag);
    isBooleanFlagEnabled('queue_management').then(setQueueManagementFlag);
    isBooleanFlagEnabled(FEATURE_FLAG_KEYS.BULK_ASSIGNMENT).then(setBulkAssignmentFlag);
    isBooleanFlagEnabled(FEATURE_FLAG_KEYS.APPROVAL_REQUEST_TYPE_COLUMN).then(setApprovalRequestTypeFlag);

    // We need to check if the user was redirected back from Okta after logging out
    // This can occur when they click "sign out" or if they try to access MM
    // while still logged into Okta which will force a redirect to logout
    const urlParams = new URLSearchParams(window.location.search);
    // okta_error=true params are added when the user is still logged into Okta elsewhere and Okta denies access
    // due to authentication method limitations
    // If the params "okta_logged_out=true" or "okta_error=true" are in the url, a banner will display
    if (urlParams.get('okta_logged_out') === 'true') {
      setOktaLoggedOut(true);
    } else if (urlParams.get('okta_error') === 'true') {
      setOktaNeedsLoggedOut(true);
    }
  }, [loadInternalSchema, loadPublicSchema, loadUser]);

  if (props.underMaintenance) {
    return <MaintenancePage />;
  }
  return (
    <PermissionProvider permissions={props.userPermissions} currentUserId={props.officeUserId}>
      <SelectedGblocProvider>
        <div id="app-root">
          <div className={siteClasses}>
            <BypassBlock />
            <CUIHeader />
            {props.userIsLoggedIn && props.activeRole === roleTypes.PRIME_SIMULATOR && <PrimeBanner />}
            {displayChangeRole && <Link to="/select-application">Change user role</Link>}
            {props.userIsLoggedIn ? <OfficeLoggedInHeader /> : <LoggedOutHeader app={pageNames.OFFICE} />}
            <main id="main" role="main" className="site__content site-office__content">
              <ConnectedLogoutOnInactivity />
              {props.hasRecentError && location.pathname === '/' && (
                <SystemError>
                  Something isn&apos;t working, but we&apos;re not sure what. Wait a minute and try again.
                  <br />
                  If that doesn&apos;t fix it, contact the{' '}
                  <a className={styles.link} href="mailto:usarmy.scott.sddc.mbx.G6-SRC-MilMove-HD@army.mil">
                    Technical Help Desk
                  </a>{' '}
                  (usarmy.scott.sddc.mbx.G6-SRC-MilMove-HD@army.mil)and give them this code:
                  <strong>{props.traceId}</strong>
                </SystemError>
              )}
              {oktaLoggedOut && <OktaLoggedOutBanner />}
              {oktaNeedsLoggedOut && <OktaNeedsLoggedOutBanner />}

              {props.showLoadingSpinner && <LoadingSpinner message={props.loadingSpinnerMessage} />}

              <Suspense fallback={<LoadingPlaceholder />}>
                {!props.userIsLoggedIn && (
                  // No Auth Routes
                  <Routes>
                    <Route path="/sign-in" element={<SignIn />} />
                    <Route path="/request-account" element={<RequestAccount />} />
                    <Route path="/invalid-permissions" element={<InvalidPermissions />} />

                    {/* 404 */}
                    <Route
                      path="*"
                      element={(props.loginIsLoading && <LoadingPlaceholder />) || <Navigate to="/sign-in" replace />}
                    />
                  </Routes>
                )}
                {!props.hasError && props.userIsLoggedIn && (
                  // Auth Routes
                  <Routes>
                    <Route path="/invalid-permissions" element={<InvalidPermissions />} />
                    {/* TOO */}
                    <Route
                      path="/moves/queue"
                      end
                      element={
                        <PrivateRoute requiredRoles={[roleTypes.TOO]}>
                          <MoveQueue
                            isQueueManagementFFEnabled={queueManagementFlag}
                            activeRole={props.activeRole}
                            isApprovalRequestTypeFFEnabled={approvalRequestTypeFlag}
                          />
                        </PrivateRoute>
                      }
                    />
                    {/* TIO */}
                    <Route
                      path="/invoicing/queue"
                      element={
                        <PrivateRoute requiredRoles={[roleTypes.TIO]}>
                          <PaymentRequestQueue
                            isQueueManagementFFEnabled={queueManagementFlag}
                            activeRole={props.activeRole}
                          />
                        </PrivateRoute>
                      }
                    />
                    {/* HQ */}
                    <Route
                      path="/hq/queues"
                      end
                      element={
                        <PrivateRoute requiredRoles={[roleTypes.HQ]}>
                          <HeadquartersQueues
                            isQueueManagementFFEnabled={queueManagementFlag}
                            activeRole={props.activeRole}
                            isApprovalRequestTypeFFEnabled={approvalRequestTypeFlag}
                          />
                        </PrivateRoute>
                      }
                    />
                    {/* SERVICES_COUNSELOR */}
                    <Route
                      key="servicesCounselingAddShipment"
                      end
                      path={servicesCounselingRoutes.SHIPMENT_ADD_PATH}
                      element={
                        <PrivateRoute requiredRoles={[roleTypes.SERVICES_COUNSELOR]}>
                          <ServicesCounselingAddShipment />
                        </PrivateRoute>
                      }
                    />
                    {props.activeRole === roleTypes.SERVICES_COUNSELOR && (
                      <Route
                        path="/:queueType/*"
                        end
                        element={
                          <PrivateRoute requiredRoles={[roleTypes.SERVICES_COUNSELOR]}>
                            <ServicesCounselingQueue
                              userPrivileges={props.userPrivileges}
                              isQueueManagementFFEnabled={queueManagementFlag}
                              isBulkAssignmentFFEnabled={bulkAssignmentFlag}
                              activeRole={props.activeRole}
                            />
                          </PrivateRoute>
                        }
                      />
                    )}
                    <Route
                      path={`${servicesCounselingRoutes.BASE_CUSTOMERS_CUSTOMER_INFO_PATH}/*`}
                      element={
                        <PrivateRoute requiredRoles={[roleTypes.SERVICES_COUNSELOR]}>
                          <CreateMoveCustomerInfo />
                        </PrivateRoute>
                      }
                    />
                    <Route
                      path={`${servicesCounselingRoutes.BASE_CUSTOMERS_ORDERS_ADD_PATH}/*`}
                      element={
                        <PrivateRoute requiredRoles={[roleTypes.SERVICES_COUNSELOR]}>
                          <ServicesCounselingAddOrders userPrivileges={props.userPrivileges} />
                        </PrivateRoute>
                      }
                    />
                    {props.activeRole === roleTypes.TIO && (
                      <Route
                        path="/:queueType/*"
                        end
                        element={
                          <PrivateRoute requiredRoles={[roleTypes.TIO]}>
                            <PaymentRequestQueue
                              isQueueManagementFFEnabled={queueManagementFlag}
                              userPrivileges={props.userPrivileges}
                              isBulkAssignmentFFEnabled={bulkAssignmentFlag}
                              activeRole={props.activeRole}
                            />
                          </PrivateRoute>
                        }
                      />
                    )}
                    {props.activeRole === roleTypes.TOO && (
                      <Route
                        path="/:queueType/*"
                        end
                        element={
                          <PrivateRoute requiredRoles={[roleTypes.TOO]}>
                            <MoveQueue
                              isQueueManagementFFEnabled={queueManagementFlag}
                              userPrivileges={props.userPrivileges}
                              isBulkAssignmentFFEnabled={bulkAssignmentFlag}
                              activeRole={props.activeRole}
                              isApprovalRequestTypeFFEnabled={approvalRequestTypeFlag}
                            />
                          </PrivateRoute>
                        }
                      />
                    )}

                    {props.activeRole === roleTypes.HQ && (
                      <Route
                        path="/:queueType/*"
                        end
                        element={
                          <PrivateRoute requiredRoles={[roleTypes.HQ]}>
                            <HeadquartersQueues
                              isQueueManagementFFEnabled={queueManagementFlag}
                              activeRole={props.activeRole}
                              isApprovalRequestTypeFFEnabled={approvalRequestTypeFlag}
                            />
                          </PrivateRoute>
                        }
                      />
                    )}
                    <Route
                      path={servicesCounselingRoutes.CREATE_CUSTOMER_PATH}
                      element={
                        <PrivateRoute requiredRoles={[roleTypes.SERVICES_COUNSELOR]}>
                          <CreateCustomerForm userPrivileges={props.userPrivileges} />
                        </PrivateRoute>
                      }
                    />
                    <Route
                      key="servicesCounselingMoveInfoRoute"
                      path={`${servicesCounselingRoutes.BASE_COUNSELING_MOVE_PATH}/*`}
                      element={
                        <PrivateRoute requiredRoles={[roleTypes.SERVICES_COUNSELOR]}>
                          <ServicesCounselingMoveInfo />
                        </PrivateRoute>
                      }
                    />

                    {/* TOO */}
                    <Route
                      path={`${tooRoutes.BASE_CUSTOMERS_CUSTOMER_INFO_PATH}`}
                      element={
                        <PrivateRoute requiredRoles={[roleTypes.TOO]}>
                          <CustomerInfo />
                        </PrivateRoute>
                      }
                    />
                    <Route
                      key="tooAddShipmentRoute"
                      end
                      path={tooRoutes.SHIPMENT_ADD_PATH}
                      element={
                        <PrivateRoute requiredRoles={[roleTypes.TOO]}>
                          <AddShipment />
                        </PrivateRoute>
                      }
                    />
                    <Route
                      key="tooEditShipmentDetailsRoute"
                      end
                      path={tooRoutes.BASE_SHIPMENT_EDIT_PATH}
                      element={
                        <PrivateRoute requiredRoles={[roleTypes.TOO]}>
                          <EditShipmentDetails />
                        </PrivateRoute>
                      }
                    />
                    <Route
                      key="tooCounselingMoveInfoRoute"
                      path={`${tooRoutes.BASE_SHIPMENT_ADVANCE_PATH_TOO}/*`}
                      element={
                        <PrivateRoute requiredRoles={[roleTypes.TOO]}>
                          <ServicesCounselingMoveInfo />
                        </PrivateRoute>
                      }
                    />

                    {/* PRIME SIMULATOR */}
                    <Route
                      key="primeSimulatorMovePath"
                      path={primeSimulatorRoutes.VIEW_MOVE_PATH}
                      element={
                        <PrivateRoute requiredRoles={[roleTypes.PRIME_SIMULATOR]}>
                          <PrimeSimulatorMoveDetails />
                        </PrivateRoute>
                      }
                    />
                    <Route
                      key="primeSimulatorCreateShipmentPath"
                      path={primeSimulatorRoutes.CREATE_SHIPMENT_PATH}
                      element={
                        <PrivateRoute requiredRoles={[roleTypes.PRIME_SIMULATOR]}>
                          <PrimeUIShipmentCreateForm />
                        </PrivateRoute>
                      }
                    />
                    <Route
                      key="primeSimulatorShipmentUpdateAddressPath"
                      path={primeSimulatorRoutes.SHIPMENT_UPDATE_ADDRESS_PATH}
                      element={
                        <PrivateRoute requiredRoles={[roleTypes.PRIME_SIMULATOR]}>
                          <PrimeUIShipmentUpdateAddress />
                        </PrivateRoute>
                      }
                      end
                    />
                    <Route
                      key="primeSimulatorUpdateShipmentPath"
                      path={primeSimulatorRoutes.UPDATE_SHIPMENT_PATH}
                      end
                      element={
                        <PrivateRoute requiredRoles={[roleTypes.PRIME_SIMULATOR]}>
                          <PrimeUIShipmentForm />
                        </PrivateRoute>
                      }
                    />
                    <Route
                      key="primeSimulatorCreatePaymentRequestsPath"
                      path={primeSimulatorRoutes.CREATE_PAYMENT_REQUEST_PATH}
                      element={
                        <PrivateRoute requiredRoles={[roleTypes.PRIME_SIMULATOR]}>
                          <PrimeSimulatorCreatePaymentRequest />
                        </PrivateRoute>
                      }
                    />
                    <Route
                      key="primeSimulatorUploadPaymentRequestDocumentsPath"
                      path={primeSimulatorRoutes.UPLOAD_DOCUMENTS_PATH}
                      element={
                        <PrivateRoute requiredRoles={[roleTypes.PRIME_SIMULATOR]}>
                          <PrimeSimulatorUploadPaymentRequestDocuments />
                        </PrivateRoute>
                      }
                    />
                    <Route
                      key="primeSimulatorUploadServiceRequestDocumentsPath"
                      path={primeSimulatorRoutes.UPLOAD_SERVICE_REQUEST_DOCUMENTS_PATH}
                      element={
                        <PrivateRoute requiredRoles={[roleTypes.PRIME_SIMULATOR]}>
                          <PrimeSimulatorUploadServiceRequestDocuments />
                        </PrivateRoute>
                      }
                    />
                    <Route
                      key="primeSimulatorCreateServiceItem"
                      path={primeSimulatorRoutes.CREATE_SERVICE_ITEM_PATH}
                      element={
                        <PrivateRoute requiredRoles={[roleTypes.PRIME_SIMULATOR]}>
                          <PrimeSimulatorCreateServiceItem />
                        </PrivateRoute>
                      }
                    />
                    <Route
                      key="primeSimulatorUpdateServiceItems"
                      path={primeSimulatorRoutes.UPDATE_SIT_SERVICE_ITEM_PATH}
                      element={
                        <PrivateRoute requiredRoles={[roleTypes.PRIME_SIMULATOR]}>
                          <PrimeSimulatorUpdateServiceItem />
                        </PrivateRoute>
                      }
                    />
                    <Route
                      key="primeSimulatorUpdateReweighPath"
                      path={primeSimulatorRoutes.SHIPMENT_UPDATE_REWEIGH_PATH}
                      element={
                        <PrivateRoute requiredRoles={[roleTypes.PRIME_SIMULATOR]}>
                          <PrimeUIShipmentUpdateReweigh />
                        </PrivateRoute>
                      }
                    />
                    <Route
                      key="primeSimulatorCreateSITExtensionRequestsPath"
                      path={primeSimulatorRoutes.CREATE_SIT_EXTENSION_REQUEST_PATH}
                      element={
                        <PrivateRoute requiredRoles={[roleTypes.PRIME_SIMULATOR]}>
                          <PrimeSimulatorCreateSITExtensionRequest />
                        </PrivateRoute>
                      }
                    />
                    <Route
                      key="primeSimulatorUpdateDestinationAddressPath"
                      path={primeSimulatorRoutes.SHIPMENT_UPDATE_DESTINATION_ADDRESS_PATH}
                      element={
                        <PrivateRoute requiredRoles={[roleTypes.PRIME_SIMULATOR]}>
                          <PrimeUIShipmentUpdateDestinationAddress />
                        </PrivateRoute>
                      }
                    />
                    <Route
                      key="primeSimulatorAcknowledgeMovePath"
                      path={primeSimulatorRoutes.ACKNOWLEDGE_MOVE_PATH}
                      element={
                        <PrivateRoute requiredRoles={[roleTypes.PRIME_SIMULATOR]}>
                          <PrimeUIAcknowledgeMove />
                        </PrivateRoute>
                      }
                    />
                    <Route
                      key="primeSimulatorAcknowledgeShipmentPath"
                      path={primeSimulatorRoutes.ACKNOWLEDGE_SHIPMENT_PATH}
                      element={
                        <PrivateRoute requiredRoles={[roleTypes.PRIME_SIMULATOR]}>
                          <PrimeUIAcknowledgeShipment />
                        </PrivateRoute>
                      }
                    />

                    {/* QAE/CSR/GSR */}
                    <Route
                      key="moveSearchPath"
                      path={qaeCSRRoutes.MOVE_SEARCH_PATH}
                      element={
                        <PrivateRoute
                          requiredRoles={[roleTypes.QAE, roleTypes.CUSTOMER_SERVICE_REPRESENTATIVE, roleTypes.GSR]}
                        >
                          <MoveSearch />
                        </PrivateRoute>
                      }
                    />
                    {/* COR */}
                    <Route
                      key="corMoveSearchPath"
                      path={contractingOfficerRoutes.MOVE_SEARCH_PATH}
                      element={
                        <PrivateRoute requiredRoles={[roleTypes.CONTRACTING_OFFICER]}>
                          <MoveSearch landingPath="mto" />
                        </PrivateRoute>
                      }
                    />

                    <Route
                      key="txoMoveInfoRoute"
                      path="/moves/:moveCode/*"
                      element={
                        <PrivateRoute
                          requiredRoles={[
                            roleTypes.TOO,
                            roleTypes.TIO,
                            roleTypes.QAE,
                            roleTypes.CONTRACTING_OFFICER,
                            roleTypes.CUSTOMER_SERVICE_REPRESENTATIVE,
                            roleTypes.GSR,
                            roleTypes.HQ,
                          ]}
                        >
                          <TXOMoveInfo />
                        </PrivateRoute>
                      }
                    />

                    <Route end path="/select-application" element={<ConnectedSelectApplication />} />

                    {/* ROOT */}
                    {props.activeRole === roleTypes.TIO && (
                      <Route
                        end
                        path="/*"
                        element={<PaymentRequestQueue isQueueManagementFFEnabled={queueManagementFlag} />}
                      />
                    )}
                    {props.activeRole === roleTypes.TOO && <Route end path="/*" element={<MoveQueue />} />}
                    {props.activeRole === roleTypes.HQ && <Route end path="/*" element={<HeadquartersQueues />} />}
                    {props.activeRole === roleTypes.SERVICES_COUNSELOR && (
                      <Route end path="/*" element={<ServicesCounselingQueue />} />
                    )}
                    {props.activeRole === roleTypes.PRIME_SIMULATOR && (
                      <Route end path="/" element={<PrimeSimulatorAvailableMoves />} />
                    )}
                    {(props.activeRole === roleTypes.QAE ||
                      props.activeRole === roleTypes.CUSTOMER_SERVICE_REPRESENTATIVE ||
                      (props.activeRole === roleTypes.GSR && gsrRoleFlag)) && (
                      <Route end path="/" element={<MoveSearch />} />
                    )}
                    {props.activeRole === roleTypes.CONTRACTING_OFFICER && (
                      <Route end path="/" element={<MoveSearch landingPath="mto" />} />
                    )}
                    {props.activeRole === roleTypes.GSR && !gsrRoleFlag && (
                      <Route end path="/*" element={<InvalidPermissions />} />
                    )}

                    <Route end path={officeRoutes.PROFILE_PATH} element={<OfficeUserProfile />} />
                    <Route end path={officeRoutes.CONTACT_INFO_EDIT_PATH} element={<OfficeUserEditContactInfo />} />
                    {/* 404 */}
                    <Route path="*" element={<NotFound />} />
                  </Routes>
                )}
              </Suspense>
            </main>
          </div>
        </div>
        <div id="modal-root" />
      </SelectedGblocProvider>
    </PermissionProvider>
  );
};

OfficeApp.propTypes = {
  loadInternalSchema: PropTypes.func.isRequired,
  loadPublicSchema: PropTypes.func.isRequired,
  loadUser: PropTypes.func.isRequired,
  officeUserId: PropTypes.string,
  loginIsLoading: PropTypes.bool,
  userIsLoggedIn: PropTypes.bool,
  userPermissions: PropTypes.arrayOf(PropTypes.string),
  userInactiveRoles: UserRolesShape,
  activeRole: PropTypes.string,
  hasRecentError: PropTypes.bool.isRequired,
  traceId: PropTypes.string.isRequired,
  router: RouterShape.isRequired,
  userPrivileges: PropTypes.arrayOf(PropTypes.object),
  underMaintenance: PropTypes.bool,
  showLoadingSpinner: PropTypes.bool,
  loadingSpinnerMessage: PropTypes.string,
  hasError: PropTypes.bool,
};

OfficeApp.defaultProps = {
  officeUserId: null,
  loginIsLoading: false,
  userIsLoggedIn: false,
  userPermissions: [],
  userInactiveRoles: [],
  activeRole: null,
  userPrivileges: [],
  underMaintenance: false,
  showLoadingSpinner: false,
  loadingSpinnerMessage: null,
  hasError: false,
};

const mapStateToProps = (state) => {
  const user = selectLoggedInUser(state);
  return {
    swaggerError: state.swaggerInternal?.hasErrored,
    officeUserId: user?.office_user?.id,
    loginIsLoading: selectGetCurrentUserIsLoading(state),
    userIsLoggedIn: selectIsLoggedIn(state),
    userPermissions: user?.permissions || [],
    userInactiveRoles: user?.inactiveRoles ?? null,
    activeRole: state.auth.activeRole,
    hasRecentError: state.interceptor.hasRecentError,
    traceId: state.interceptor.traceId,
    userPrivileges: user?.privileges || null,
    underMaintenance: selectUnderMaintenance(state),
    showLoadingSpinner: selectShowLoadingSpinner(state),
    loadingSpinnerMessage: selectLoadingSpinnerMessage(state),
    hasError: state.auth.hasError,
  };
};

const mapDispatchToProps = {
  loadInternalSchema: loadInternalSchemaAction,
  loadPublicSchema: loadPublicSchemaAction,
  loadUser: loadUserAction,
};

export default withContext(withRouter(connect(mapStateToProps, mapDispatchToProps)(OfficeApp)));
