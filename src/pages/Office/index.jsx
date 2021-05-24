/* eslint-disable react/jsx-props-no-spreading */
import React, { Component, lazy, Suspense } from 'react';
import PropTypes from 'prop-types';
import { Route, Switch, withRouter, matchPath, Link } from 'react-router-dom';
import { connect } from 'react-redux';
import classnames from 'classnames';

import '../../../node_modules/uswds/dist/css/uswds.css';
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
import OfficeLoggedInHeader from 'containers/Headers/OfficeLoggedInHeader';
import LoggedOutHeader from 'containers/Headers/LoggedOutHeader';
import { ConnectedSelectApplication } from 'pages/SelectApplication/SelectApplication';
import { roleTypes } from 'constants/userRoles';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import { withContext } from 'shared/AppContext';
import { LocationShape, UserRolesShape } from 'types/index';
import { servicesCounselingRoutes } from 'constants/routes';

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
const ServicesCounselingEditShipmentDetails = lazy(() =>
  import('pages/Office/ServicesCounselingEditShipmentDetails/ServicesCounselingEditShipmentDetails'),
);

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
      userIsLoggedIn,
      userRoles,
      location: { pathname },
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
        requiredRoles={[roleTypes.TOO, roleTypes.TIO]}
      />,
    ];

    const isFullscreenPage = matchPath(pathname, {
      path: '/moves/:moveCode/payment-requests/:id',
    });

    const siteClasses = classnames('site', {
      [`site--fullscreen`]: isFullscreenPage,
    });

    return (
      <>
        <div id="app-root">
          <div className={siteClasses}>
            <BypassBlock />
            <FOUOHeader />
            {displayChangeRole && <Link to="/select-application">Change user role</Link>}
            {!hideHeaderPPM && <>{userIsLoggedIn ? <OfficeLoggedInHeader /> : <LoggedOutHeader />}</>}
            <main id="main" role="main" className="site__content site-office__content">
              <ConnectedLogoutOnInactivity />

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
                      key="servicesCounselingEditShipmentDetailsRoute"
                      exact
                      path={servicesCounselingRoutes.EDIT_SHIPMENT_INFO_PATH}
                      component={ServicesCounselingEditShipmentDetails}
                      requiredRoles={[roleTypes.SERVICES_COUNSELOR]}
                    />
                    <PrivateRoute
                      path="/counseling/queue"
                      exact
                      component={ServicesCounselingQueue}
                      requiredRoles={[roleTypes.SERVICES_COUNSELOR]}
                    />
                    <PrivateRoute
                      key="servicesCounselingMoveInfoRoute"
                      path="/counseling/moves/:moveCode"
                      component={ServicesCounselingMoveInfo}
                      requiredRoles={[roleTypes.SERVICES_COUNSELOR]}
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
      </>
    );
  }
}

OfficeApp.propTypes = {
  loadInternalSchema: PropTypes.func.isRequired,
  loadPublicSchema: PropTypes.func.isRequired,
  loadUser: PropTypes.func.isRequired,
  location: LocationShape,
  userIsLoggedIn: PropTypes.bool,
  userRoles: UserRolesShape,
  activeRole: PropTypes.string,
};

OfficeApp.defaultProps = {
  location: { pathname: '' },
  userIsLoggedIn: false,
  userRoles: [],
  activeRole: null,
};

const mapStateToProps = (state) => {
  const user = selectLoggedInUser(state);

  return {
    swaggerError: state.swaggerInternal.hasErrored,
    userIsLoggedIn: selectIsLoggedIn(state),
    userRoles: user?.roles || [],
    activeRole: state.auth.activeRole,
  };
};

const mapDispatchToProps = {
  loadInternalSchema: loadInternalSchemaAction,
  loadPublicSchema: loadPublicSchemaAction,
  loadUser: loadUserAction,
};

export default withContext(withRouter(connect(mapStateToProps, mapDispatchToProps)(OfficeApp)));
