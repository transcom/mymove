/* eslint-disable react/jsx-props-no-spreading */
import React, { Component, lazy, Suspense } from 'react';
import PropTypes from 'prop-types';
import { Route, Routes, Link, matchPath, Navigate } from 'react-router-dom';
import { connect } from 'react-redux';
import classnames from 'classnames';

import styles from './Office.module.scss';
import 'styles/full_uswds.scss';
import 'scenes/Office/office.scss';

// API / Redux actions
import { selectGetCurrentUserIsLoading, selectIsLoggedIn } from 'store/auth/selectors';
import { loadUser as loadUserAction } from 'store/auth/actions';
import { selectLoggedInUser } from 'store/entities/selectors';
import {
  loadInternalSchema as loadInternalSchemaAction,
  loadPublicSchema as loadPublicSchemaAction,
} from 'shared/Swagger/ducks';
// Shared layout components
import ConnectedLogoutOnInactivity from 'layout/LogoutOnInactivity';
import SomethingWentWrong from 'shared/SomethingWentWrong';
import CUIHeader from 'components/CUIHeader/CUIHeader';
import BypassBlock from 'components/BypassBlock';
import SystemError from 'components/SystemError';
import NotFound from 'components/NotFound/NotFound';
import OfficeLoggedInHeader from 'containers/Headers/OfficeLoggedInHeader';
import LoggedOutHeader from 'containers/Headers/LoggedOutHeader';
import { ConnectedSelectApplication } from 'pages/SelectApplication/SelectApplication';
import { roleTypes } from 'constants/userRoles';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import { withContext } from 'shared/AppContext';
import { RouterShape, UserRolesShape } from 'types/index';
import { servicesCounselingRoutes, primeSimulatorRoutes, tooRoutes, qaeCSRRoutes } from 'constants/routes';
import PrimeBanner from 'pages/PrimeUI/PrimeBanner/PrimeBanner';
import PermissionProvider from 'components/Restricted/PermissionProvider';
import withRouter from 'utils/routing';
import ConnectedProtectedRoute from 'components/ProtectedRoute/ProtectedRoute';

// Lazy load these dependencies (they correspond to unique routes & only need to be loaded when that URL is accessed)
const SignIn = lazy(() => import('pages/SignIn/SignIn'));
const InvalidPermissions = lazy(() => import('pages/InvalidPermissions/InvalidPermissions'));
// PPM pages (TODO move into src/pages)
const MoveInfo = lazy(() => import('scenes/Office/MoveInfo'));
const Queues = lazy(() => import('scenes/Office/Queues'));
const OrdersInfo = lazy(() => import('scenes/Office/OrdersInfo'));
const DocumentViewer = lazy(() => import('scenes/Office/DocumentViewer'));
// TXO
const TXOMoveInfo = lazy(() => import('pages/Office/TXOMoveInfo/TXOMoveInfo'));
// TOO pages
const MoveQueue = lazy(() => import('pages/Office/MoveQueue/MoveQueue'));
// TIO pages
const PaymentRequestQueue = lazy(() => import('pages/Office/PaymentRequestQueue/PaymentRequestQueue'));
// Services Counselor pages
const ServicesCounselingMoveInfo = lazy(() =>
  import('pages/Office/ServicesCounselingMoveInfo/ServicesCounselingMoveInfo'),
);
const ServicesCounselingQueue = lazy(() => import('pages/Office/ServicesCounselingQueue/ServicesCounselingQueue'));
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
const PrimeSimulatorCreateServiceItem = lazy(() => import('pages/PrimeUI/CreateServiceItem/CreateServiceItem'));
const PrimeUIShipmentUpdateAddress = lazy(() => import('pages/PrimeUI/Shipment/PrimeUIShipmentUpdateAddress'));
const PrimeUIShipmentUpdateReweigh = lazy(() => import('pages/PrimeUI/Shipment/PrimeUIShipmentUpdateReweigh'));

const QAECSRMoveSearch = lazy(() => import('pages/Office/QAECSRMoveSearch/QAECSRMoveSearch'));

export class OfficeApp extends Component {
  constructor(props) {
    super(props);

    this.state = {
      hasError: false,
      error: undefined,
      info: undefined,
    };
  }

  componentDidMount() {
    document.title = 'Transcom PPP: Office';

    const { loadUser, loadInternalSchema, loadPublicSchema } = this.props;

    loadInternalSchema();
    loadPublicSchema();
    loadUser();
  }

  componentDidCatch(error, info) {
    this.setState({
      hasError: true,
      error,
      info,
    });
  }

  render() {
    const { hasError, error, info } = this.state;
    const {
      activeRole,
      officeUserId,
      loginIsLoading,
      userIsLoggedIn,
      userPermissions,
      userRoles,
      router: {
        location,
        location: { pathname },
      },
      hasRecentError,
      traceId,
    } = this.props;

    const roleSelected = userIsLoggedIn && activeRole;

    // TODO - test login page?

    // TODO - I don't love this solution but it will work for now. Ideally we can abstract the page layout into a separate file where each route can use it or not
    // Don't show Header on OrdersInfo or DocumentViewer pages (PPM only)
    const hideHeaderPPM =
      roleSelected &&
      activeRole === roleTypes.PPM &&
      (matchPath(
        {
          path: '/moves/:moveId/documents/:moveDocumentId?',
          end: true,
        },
        pathname,
      ) ||
        matchPath(
          {
            path: '/moves/:moveId/orders',
            end: true,
          },
          pathname,
        ));

    const displayChangeRole =
      userIsLoggedIn &&
      userRoles?.length > 1 &&
      !matchPath(
        {
          path: '/select-application',
          end: true,
        },
        pathname,
      );

    const ppmRoutes = [
      <Route
        key="ppmOrdersRoute"
        path="/moves/:moveId/orders"
        element={
          <ConnectedProtectedRoute requiredRoles={[roleTypes.PPM]}>
            <OrdersInfo />
          </ConnectedProtectedRoute>
        }
      />,
      <Route
        key="ppmMoveDocumentRoute"
        path="/moves/:moveId/documents/:moveDocumentId?"
        element={
          <ConnectedProtectedRoute requiredRoles={[roleTypes.PPM]}>
            <DocumentViewer />
          </ConnectedProtectedRoute>
        }
      />,
    ];

    // TODO - Services counseling routes not finalized, revisit
    const txoRoutes = [
      <Route
        key="txoMoveInfoRoute"
        path="/moves/:moveCode/*"
        element={
          <ConnectedProtectedRoute requiredRoles={[roleTypes.TOO, roleTypes.TIO, roleTypes.QAE_CSR]}>
            <TXOMoveInfo />
          </ConnectedProtectedRoute>
        }
      />,
    ];

    const isFullscreenPage = matchPath(
      {
        path: '/moves/:moveCode/payment-requests/:id',
      },
      pathname,
    );

    const siteClasses = classnames('site', {
      [`site--fullscreen`]: isFullscreenPage,
    });

    return (
      <PermissionProvider permissions={userPermissions} currentUserId={officeUserId}>
        <div id="app-root">
          <div className={siteClasses}>
            <BypassBlock />
            <CUIHeader />
            {roleSelected && activeRole === roleTypes.PRIME_SIMULATOR && <PrimeBanner />}
            {displayChangeRole && <Link to="/select-application">Change user role</Link>}
            {!hideHeaderPPM && userIsLoggedIn ? <OfficeLoggedInHeader /> : <LoggedOutHeader />}
            <main id="main" role="main" className="site__content site-office__content">
              <ConnectedLogoutOnInactivity />
              {hasRecentError && location.pathname === '/' && (
                <SystemError>
                  Something isn&apos;t working, but we&apos;re not sure what. Wait a minute and try again.
                  <br />
                  If that doesn&apos;t fix it, contact the{' '}
                  <a className={styles.link} href="https://move.mil/customer-service#technical-help-desk">
                    Technical Help Desk
                  </a>{' '}
                  and give them this code: <strong>{traceId}</strong>
                </SystemError>
              )}
              {hasError && <SomethingWentWrong error={error} info={info} />}

              <Suspense fallback={<LoadingPlaceholder />}>
                {!hasError && !loginIsLoading && (
                  <Routes>
                    {/* no auth */}
                    <Route path="/sign-in" element={<SignIn />} />
                    {/* no auth */}
                    <Route path="/invalid-permissions" element={<InvalidPermissions />} />
                    {/* PPM */}
                    <Route
                      path="/queues/:queueType/moves/:moveId"
                      element={
                        <ConnectedProtectedRoute requiredRoles={[roleTypes.PPM]}>
                          <MoveInfo />
                        </ConnectedProtectedRoute>
                      }
                    />
                    <Route
                      path="/queues/:queueType"
                      element={
                        <ConnectedProtectedRoute requiredRoles={[roleTypes.PPM]}>
                          <Queues />
                        </ConnectedProtectedRoute>
                      }
                    />
                    {/* TXO */}
                    <Route
                      path="/moves/queue"
                      end
                      element={
                        <ConnectedProtectedRoute requiredRoles={[roleTypes.TOO]}>
                          <MoveQueue />
                        </ConnectedProtectedRoute>
                      }
                    />
                    <Route
                      path="/invoicing/queue"
                      element={
                        <ConnectedProtectedRoute requiredRoles={[roleTypes.TIO]}>
                          <PaymentRequestQueue />
                        </ConnectedProtectedRoute>
                      }
                    />

                    {/* SERVICES_COUNSELOR */}
                    <Route
                      key="servicesCounselingMoveInfoRoute"
                      path={`${servicesCounselingRoutes.BASE_COUNSELING_MOVE_PATH}/*`}
                      element={
                        <ConnectedProtectedRoute requiredRoles={[roleTypes.SERVICES_COUNSELOR]}>
                          <ServicesCounselingMoveInfo />
                        </ConnectedProtectedRoute>
                      }
                    />
                    {/* TOO */}
                    <Route
                      key="tooEditShipmentDetailsRoute"
                      end
                      path={tooRoutes.BASE_SHIPMENT_EDIT_PATH}
                      element={
                        <ConnectedProtectedRoute requiredRoles={[roleTypes.TOO]}>
                          <EditShipmentDetails />
                        </ConnectedProtectedRoute>
                      }
                    />
                    {/* PRIME SIMULATOR */}
                    <Route
                      key="primeSimulatorMovePath"
                      path={primeSimulatorRoutes.VIEW_MOVE_PATH}
                      element={
                        <ConnectedProtectedRoute requiredRoles={[roleTypes.PRIME_SIMULATOR]}>
                          <PrimeSimulatorMoveDetails />
                        </ConnectedProtectedRoute>
                      }
                    />
                    <Route
                      key="primeSimulatorCreateShipmentPath"
                      path={primeSimulatorRoutes.CREATE_SHIPMENT_PATH}
                      element={
                        <ConnectedProtectedRoute requiredRoles={[roleTypes.PRIME_SIMULATOR]}>
                          <PrimeUIShipmentCreateForm />
                        </ConnectedProtectedRoute>
                      }
                    />
                    <Route
                      key="primeSimulatorShipmentUpdateAddressPath"
                      path={primeSimulatorRoutes.SHIPMENT_UPDATE_ADDRESS_PATH}
                      element={
                        <ConnectedProtectedRoute requiredRoles={[roleTypes.PRIME_SIMULATOR]}>
                          <PrimeUIShipmentUpdateAddress />
                        </ConnectedProtectedRoute>
                      }
                      end
                    />
                    <Route
                      key="primeSimulatorUpdateShipmentPath"
                      path={primeSimulatorRoutes.UPDATE_SHIPMENT_PATH}
                      end
                      element={
                        <ConnectedProtectedRoute requiredRoles={[roleTypes.PRIME_SIMULATOR]}>
                          <PrimeUIShipmentForm />
                        </ConnectedProtectedRoute>
                      }
                    />
                    <Route
                      key="primeSimulatorCreatePaymentRequestsPath"
                      path={primeSimulatorRoutes.CREATE_PAYMENT_REQUEST_PATH}
                      element={
                        <ConnectedProtectedRoute requiredRoles={[roleTypes.PRIME_SIMULATOR]}>
                          <PrimeSimulatorCreatePaymentRequest />
                        </ConnectedProtectedRoute>
                      }
                    />
                    <Route
                      key="primeSimulatorUploadPaymentRequestDocumentsPath"
                      path={primeSimulatorRoutes.UPLOAD_DOCUMENTS_PATH}
                      element={
                        <ConnectedProtectedRoute requiredRoles={[roleTypes.PRIME_SIMULATOR]}>
                          <PrimeSimulatorUploadPaymentRequestDocuments />
                        </ConnectedProtectedRoute>
                      }
                    />
                    <Route
                      key="primeSimulatorCreateServiceItem"
                      path={primeSimulatorRoutes.CREATE_SERVICE_ITEM_PATH}
                      element={
                        <ConnectedProtectedRoute requiredRoles={[roleTypes.PRIME_SIMULATOR]}>
                          <PrimeSimulatorCreateServiceItem />
                        </ConnectedProtectedRoute>
                      }
                    />
                    <Route
                      key="primeSimulatorUpdateReweighPath"
                      path={primeSimulatorRoutes.SHIPMENT_UPDATE_REWEIGH_PATH}
                      element={
                        <ConnectedProtectedRoute requiredRoles={[roleTypes.PRIME_SIMULATOR]}>
                          <PrimeUIShipmentUpdateReweigh />
                        </ConnectedProtectedRoute>
                      }
                    />

                    {/* QAE/CSR */}
                    <Route
                      key="qaeCSRMoveSearchPath"
                      path={qaeCSRRoutes.MOVE_SEARCH_PATH}
                      element={
                        <ConnectedProtectedRoute requiredRoles={[roleTypes.QAE_CSR]}>
                          <QAECSRMoveSearch />
                        </ConnectedProtectedRoute>
                      }
                    />

                    {/* PPM & TXO conflicting routes - select based on user role */}
                    {roleSelected && activeRole === roleTypes.PPM ? ppmRoutes : txoRoutes}

                    <Route end path="/select-application" element={<ConnectedSelectApplication />} />

                    {/* ROOT */}
                    {roleSelected && activeRole === roleTypes.PPM && (
                      <Route end path="/" element={<Queues queueType="new" />} />
                    )}
                    {roleSelected && activeRole === roleTypes.TIO && (
                      <Route end path="/" element={<PaymentRequestQueue />} />
                    )}
                    {roleSelected && activeRole === roleTypes.TOO && <Route end path="/" element={<MoveQueue />} />}
                    {roleSelected && activeRole === roleTypes.SERVICES_COUNSELOR && (
                      <Route end path="/*" element={<ServicesCounselingQueue />} />
                    )}
                    {roleSelected && activeRole === roleTypes.PRIME_SIMULATOR && (
                      <Route end path="/" element={<PrimeSimulatorAvailableMoves />} />
                    )}
                    {roleSelected && activeRole === roleTypes.QAE_CSR && (
                      <Route end path="/" element={<QAECSRMoveSearch />} />
                    )}
                    {/* 404 */}
                    <Route
                      path="*"
                      element={
                        (loginIsLoading && <LoadingPlaceholder />) ||
                        (!userIsLoggedIn && <Navigate to="/sign-in" replace />) || <NotFound />
                      }
                    />
                  </Routes>
                )}
              </Suspense>
            </main>
          </div>
        </div>
        <div id="modal-root" />
      </PermissionProvider>
    );
  }
}

OfficeApp.propTypes = {
  loadInternalSchema: PropTypes.func.isRequired,
  loadPublicSchema: PropTypes.func.isRequired,
  loadUser: PropTypes.func.isRequired,
  officeUserId: PropTypes.string,
  loginIsLoading: PropTypes.bool,
  userIsLoggedIn: PropTypes.bool,
  userPermissions: PropTypes.arrayOf(PropTypes.string),
  userRoles: UserRolesShape,
  activeRole: PropTypes.string,
  hasRecentError: PropTypes.bool.isRequired,
  traceId: PropTypes.string.isRequired,
  router: RouterShape.isRequired,
};

OfficeApp.defaultProps = {
  officeUserId: null,
  loginIsLoading: true,
  userIsLoggedIn: false,
  userPermissions: [],
  userRoles: [],
  activeRole: null,
};

const mapStateToProps = (state) => {
  const user = selectLoggedInUser(state);

  return {
    swaggerError: state.swaggerInternal.hasErrored,
    officeUserId: user?.office_user?.id,
    loginIsLoading: selectGetCurrentUserIsLoading(state),
    userIsLoggedIn: selectIsLoggedIn(state),
    userPermissions: user?.permissions || [],
    userRoles: user?.roles || [],
    activeRole: state.auth.activeRole,
    hasRecentError: state.interceptor.hasRecentError,
    traceId: state.interceptor.traceId,
  };
};

const mapDispatchToProps = {
  loadInternalSchema: loadInternalSchemaAction,
  loadPublicSchema: loadPublicSchemaAction,
  loadUser: loadUserAction,
};

export default withContext(withRouter(connect(mapStateToProps, mapDispatchToProps)(OfficeApp)));
