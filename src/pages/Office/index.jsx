/* eslint-disable react/jsx-props-no-spreading */
import React, { Component, lazy, Suspense } from 'react';
import PropTypes from 'prop-types';
import { Route, Switch, withRouter, matchPath, Link } from 'react-router-dom';
import { connect } from 'react-redux';
import classnames from 'classnames';

import styles from './Office.module.scss';
import 'uswds/dist/css/uswds.css';
import 'scenes/Office/office.scss';

// API / Redux actions
import { selectIsLoggedIn } from 'store/auth/selectors';
import { loadUser as loadUserAction } from 'store/auth/actions';
import { selectLoggedInUser } from 'store/entities/selectors';
import {
  loadInternalSchema as loadInternalSchemaAction,
  loadPublicSchema as loadPublicSchemaAction,
} from 'shared/Swagger/ducks';
// Shared layout components
import ConnectedLogoutOnInactivity from 'layout/LogoutOnInactivity';
import PrivateRoute from 'containers/PrivateRoute';
import SomethingWentWrong from 'shared/SomethingWentWrong';
import FOUOHeader from 'components/FOUOHeader';
import BypassBlock from 'components/BypassBlock';
import SystemError from 'components/SystemError';
import OfficeLoggedInHeader from 'containers/Headers/OfficeLoggedInHeader';
import LoggedOutHeader from 'containers/Headers/LoggedOutHeader';
import { ConnectedSelectApplication } from 'pages/SelectApplication/SelectApplication';
import { roleTypes } from 'constants/userRoles';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import { withContext } from 'shared/AppContext';
import { LocationShape, UserRolesShape } from 'types/index';
import { servicesCounselingRoutes, primeSimulatorRoutes, tooRoutes, qaeCSRRoutes } from 'constants/routes';
import PrimeBanner from 'pages/PrimeUI/PrimeBanner/PrimeBanner';
import PermissionProvider from 'components/Restricted/PermissionProvider';

// Lazy load these dependencies (they correspond to unique routes & only need to be loaded when that URL is accessed)
const SignIn = lazy(() => import('pages/SignIn/SignIn'));
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
const ServicesCounselingAddShipment = lazy(() =>
  import('pages/Office/ServicesCounselingAddShipment/ServicesCounselingAddShipment'),
);
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
      userIsLoggedIn,
      userPermissions,
      userRoles,
      location: { pathname },
      hasRecentError,
      traceId,
      history,
    } = this.props;
    const selectedRole = userIsLoggedIn && activeRole;

    // TODO - test login page?

    // TODO - I don't love this solution but it will work for now. Ideally we can abstract the page layout into a separate file where each route can use it or not
    // Don't show Header on OrdersInfo or DocumentViewer pages (PPM only)
    const hideHeaderPPM =
      selectedRole === roleTypes.PPM &&
      (matchPath(pathname, {
        path: '/moves/:moveId/documents/:moveDocumentId?',
        exact: true,
      }) ||
        matchPath(pathname, {
          path: '/moves/:moveId/orders',
          exact: true,
        }));

    const displayChangeRole =
      userIsLoggedIn &&
      userRoles?.length > 1 &&
      !matchPath(pathname, {
        path: '/select-application',
        exact: true,
      });

    const ppmRoutes = [
      <PrivateRoute
        key="ppmOrdersRoute"
        path="/moves/:moveId/orders"
        component={OrdersInfo}
        requiredRoles={[roleTypes.PPM]}
      />,
      <PrivateRoute
        key="ppmMoveDocumentRoute"
        path="/moves/:moveId/documents/:moveDocumentId?"
        component={DocumentViewer}
        requiredRoles={[roleTypes.PPM]}
      />,
    ];

    // TODO - Services counseling routes not finalized, revisit
    const txoRoutes = [
      <PrivateRoute
        key="txoMoveInfoRoute"
        path="/moves/:moveCode"
        component={TXOMoveInfo}
        requiredRoles={[roleTypes.TOO, roleTypes.TIO, roleTypes.QAE_CSR]}
      />,
    ];

    const isFullscreenPage = matchPath(pathname, {
      path: '/moves/:moveCode/payment-requests/:id',
    });

    const siteClasses = classnames('site', {
      [`site--fullscreen`]: isFullscreenPage,
    });

    return (
      <PermissionProvider permissions={userPermissions} currentUserId={officeUserId}>
        <div id="app-root">
          <div className={siteClasses}>
            <BypassBlock />
            <FOUOHeader />
            {selectedRole === roleTypes.PRIME_SIMULATOR && <PrimeBanner />}
            {displayChangeRole && <Link to="/select-application">Change user role</Link>}
            {!hideHeaderPPM && userIsLoggedIn ? <OfficeLoggedInHeader /> : <LoggedOutHeader />}
            <main id="main" role="main" className="site__content site-office__content">
              <ConnectedLogoutOnInactivity />
              {hasRecentError && history.location.pathname === '/' && (
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
                {!hasError && (
                  <Switch>
                    {/* no auth */}
                    <Route path="/sign-in" component={SignIn} />

                    {/* PPM */}
                    <PrivateRoute
                      path="/queues/:queueType/moves/:moveId"
                      component={MoveInfo}
                      requiredRoles={[roleTypes.PPM]}
                    />
                    <PrivateRoute path="/queues/:queueType" component={Queues} requiredRoles={[roleTypes.PPM]} />

                    {/* TXO */}
                    <PrivateRoute path="/moves/queue" exact component={MoveQueue} requiredRoles={[roleTypes.TOO]} />
                    <PrivateRoute
                      path="/invoicing/queue"
                      component={PaymentRequestQueue}
                      requiredRoles={[roleTypes.TIO]}
                    />

                    {/* SERVICES_COUNSELOR */}
                    <PrivateRoute
                      key="servicesCounselingAddShipment"
                      exact
                      path={servicesCounselingRoutes.SHIPMENT_ADD_PATH}
                      component={ServicesCounselingAddShipment}
                      requiredRoles={[roleTypes.SERVICES_COUNSELOR]}
                    />
                    <PrivateRoute
                      path={servicesCounselingRoutes.QUEUE_VIEW_PATH}
                      exact
                      component={ServicesCounselingQueue}
                      requiredRoles={[roleTypes.SERVICES_COUNSELOR]}
                    />
                    <PrivateRoute
                      key="servicesCounselingMoveInfoRoute"
                      path={servicesCounselingRoutes.BASE_MOVE_PATH}
                      component={ServicesCounselingMoveInfo}
                      requiredRoles={[roleTypes.SERVICES_COUNSELOR]}
                    />

                    {/* TOO */}
                    <PrivateRoute
                      key="tooEditShipmentDetailsRoute"
                      exact
                      path={tooRoutes.SHIPMENT_EDIT_PATH}
                      component={EditShipmentDetails}
                      requiredRoles={[roleTypes.TOO]}
                    />

                    {/* PRIME SIMULATOR */}
                    <PrivateRoute
                      key="primeSimulatorMovePath"
                      path={primeSimulatorRoutes.VIEW_MOVE_PATH}
                      component={PrimeSimulatorMoveDetails}
                      requiredRoles={[roleTypes.PRIME_SIMULATOR]}
                    />

                    <PrivateRoute
                      key="primeSimulatorCreateShipmentPath"
                      path={primeSimulatorRoutes.CREATE_SHIPMENT_PATH}
                      component={PrimeUIShipmentCreateForm}
                      requiredRoles={[roleTypes.PRIME_SIMULATOR]}
                    />

                    <PrivateRoute
                      key="primeSimulatorShipmentUpdateAddressPath"
                      path={primeSimulatorRoutes.SHIPMENT_UPDATE_ADDRESS_PATH}
                      component={PrimeUIShipmentUpdateAddress}
                      requiredRoles={[roleTypes.PRIME_SIMULATOR]}
                      exact
                    />

                    <PrivateRoute
                      key="primeSimulatorUpdateShipmentPath"
                      path={primeSimulatorRoutes.UPDATE_SHIPMENT_PATH}
                      exact
                      component={PrimeUIShipmentForm}
                      requiredRoles={[roleTypes.PRIME_SIMULATOR]}
                    />

                    <PrivateRoute
                      key="primeSimulatorCreatePaymentRequestsPath"
                      path={primeSimulatorRoutes.CREATE_PAYMENT_REQUEST_PATH}
                      component={PrimeSimulatorCreatePaymentRequest}
                      requiredRoles={[roleTypes.PRIME_SIMULATOR]}
                    />

                    <PrivateRoute
                      key="primeSimulatorUploadPaymentRequestDocumentsPath"
                      path={primeSimulatorRoutes.UPLOAD_DOCUMENTS_PATH}
                      component={PrimeSimulatorUploadPaymentRequestDocuments}
                      requiredRoles={[roleTypes.PRIME_SIMULATOR]}
                    />

                    <PrivateRoute
                      key="primeSimulatorCreateServiceItem"
                      path={primeSimulatorRoutes.CREATE_SERVICE_ITEM_PATH}
                      component={PrimeSimulatorCreateServiceItem}
                      requiredRoles={[roleTypes.PRIME_SIMULATOR]}
                    />

                    <PrivateRoute
                      key="primeSimulatorUpdateReweighPath"
                      path={primeSimulatorRoutes.SHIPMENT_UPDATE_REWEIGH_PATH}
                      component={PrimeUIShipmentUpdateReweigh}
                      requiredRoles={[roleTypes.PRIME_SIMULATOR]}
                    />

                    {/* QAE/CSR */}
                    <PrivateRoute
                      key="qaeCSRMoveSearchPath"
                      path={qaeCSRRoutes.MOVE_SEARCH_PATH}
                      component={QAECSRMoveSearch}
                      requiredRoles={[roleTypes.QAE_CSR]}
                    />

                    {/* PPM & TXO conflicting routes - select based on user role */}
                    {selectedRole === roleTypes.PPM ? ppmRoutes : txoRoutes}

                    <PrivateRoute exact path="/select-application" component={ConnectedSelectApplication} />
                    {/* ROOT */}
                    <PrivateRoute
                      exact
                      path="/"
                      render={(routeProps) => {
                        switch (selectedRole) {
                          case roleTypes.PPM:
                            return <Queues queueType="new" {...routeProps} />;
                          case roleTypes.TIO:
                            return <PaymentRequestQueue {...routeProps} />;
                          case roleTypes.TOO:
                            return <MoveQueue {...routeProps} />;
                          case roleTypes.SERVICES_COUNSELOR:
                            return <ServicesCounselingQueue {...routeProps} />;
                          case roleTypes.PRIME_SIMULATOR:
                            return <PrimeSimulatorAvailableMoves {...routeProps} />;
                          case roleTypes.QAE_CSR:
                            return <QAECSRMoveSearch {...routeProps} />;
                          default:
                            // User has unknown role or shouldn't have access
                            return <div />;
                        }
                      }}
                    />
                  </Switch>
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
  location: LocationShape,
  officeUserId: PropTypes.string,
  userIsLoggedIn: PropTypes.bool,
  userPermissions: PropTypes.arrayOf(PropTypes.string),
  userRoles: UserRolesShape,
  activeRole: PropTypes.string,
  hasRecentError: PropTypes.bool.isRequired,
  traceId: PropTypes.string.isRequired,
  history: PropTypes.shape({
    location: PropTypes.shape({
      pathname: PropTypes.string,
    }),
  }),
};

OfficeApp.defaultProps = {
  location: { pathname: '' },
  officeUserId: null,
  userIsLoggedIn: false,
  userPermissions: [],
  userRoles: [],
  activeRole: null,
  history: {
    location: { pathname: '' },
  },
};

const mapStateToProps = (state) => {
  const user = selectLoggedInUser(state);

  return {
    swaggerError: state.swaggerInternal.hasErrored,
    officeUserId: user?.office_user?.id,
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
