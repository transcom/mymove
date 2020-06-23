import React, { Component, lazy, Suspense } from 'react';
import PropTypes from 'prop-types';
import { Route, Switch, withRouter, matchPath } from 'react-router-dom';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';

import 'uswds';
import '../../../node_modules/uswds/dist/css/uswds.css';
import 'scenes/Office/office.scss';

// API / Redux actions
import { getCurrentUserInfo as getCurrentUserInfoAction, selectCurrentUser } from 'shared/Data/users';
import {
  loadInternalSchema as loadInternalSchemaAction,
  loadPublicSchema as loadPublicSchemaAction,
} from 'shared/Swagger/ducks';
// Shared layout components
import ConnectedLogoutOnInactivity from 'shared/User/LogoutOnInactivity';
import PrivateRoute from 'shared/User/PrivateRoute';
import SomethingWentWrong from 'shared/SomethingWentWrong';
import { QueueHeader } from 'shared/Header/Office';
import FOUOHeader from 'components/FOUOHeader';
import { ConnectedSelectApplication } from 'pages/SelectApplication/SelectApplication';
import { roleTypes } from 'constants/userRoles';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import { withContext } from 'shared/AppContext';

// Lazy load these dependencies (they correspond to unique routes & only need to be loaded when that URL is accessed)
const ConnectedOfficeHome = lazy(() => import('pages/OfficeHome'));
// PPM pages (TODO move into src/pages)
const MoveInfo = lazy(() => import('scenes/Office/MoveInfo'));
const Queues = lazy(() => import('scenes/Office/Queues'));
const OrdersInfo = lazy(() => import('scenes/Office/OrdersInfo'));
const DocumentViewer = lazy(() => import('scenes/Office/DocumentViewer'));
// TXO
const TXOMoveInfo = lazy(() => import('pages/TXOMoveInfo'));
// TOO pages (TODO move into src/pages)
const TOO = lazy(() => import('scenes/Office/TOO/too'));
const CustomerDetails = lazy(() => import('scenes/Office/TOO/customerDetails'));
const TOOVerificationInProgress = lazy(() => import('scenes/Office/TOO/tooVerificationInProgress'));
// TIO pages (TODO move into src/pages)
const TIO = lazy(() => import('scenes/Office/TIO/tio'));
const PaymentRequestIndex = lazy(() => import('scenes/Office/TIO/paymentRequestIndex'));

export class OfficeWrapper extends Component {
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

    const { loadInternalSchema, loadPublicSchema, getCurrentUserInfo } = this.props;

    loadInternalSchema();
    loadPublicSchema();
    getCurrentUserInfo();
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
      context: {
        flags: { too, tio },
      },
      location: { pathname },
    } = this.props;

    // TODO - doesn't seem like any of these routes are accessible if not logged in, suggest refactor PrivateRoute into HOC required by this entrypoint
    // TODO - test login page?

    // TODO - I don't love this solution but it will work for now. Ideally we can abstract the page layout into a separate file where each route can use it or not
    // Don't show Header on OrdersInfo or DocumentViewer pages
    const hideHeader =
      matchPath(pathname, {
        path: '/moves/:moveId/documents/:moveDocumentId?',
        exact: true,
      }) ||
      matchPath(pathname, {
        path: '/moves/:moveId/orders',
        exact: true,
      });

    return (
      <div className="site">
        <FOUOHeader />
        {!hideHeader && <QueueHeader />}
        <main role="main" className="site__content">
          <ConnectedLogoutOnInactivity />

          {hasError && <SomethingWentWrong error={error} info={info} />}

          <Suspense fallback={<LoadingPlaceholder />}>
            {!hasError && (
              <Switch>
                {/* ROOT */}
                <PrivateRoute exact path="/" component={ConnectedOfficeHome} />
                <PrivateRoute exact path="/select-application" component={ConnectedSelectApplication} />

                {/* PPM routes */}
                <PrivateRoute
                  path="/queues/:queueType/moves/:moveId"
                  component={MoveInfo}
                  requiredRoles={[roleTypes.PPM]}
                />
                <PrivateRoute path="/queues/:queueType" component={Queues} requiredRoles={[roleTypes.PPM]} />
                <PrivateRoute path="/moves/:moveId/orders" component={OrdersInfo} requiredRoles={[roleTypes.PPM]} />
                <PrivateRoute
                  path="/moves/:moveId/documents/:moveDocumentId?"
                  component={DocumentViewer}
                  requiredRoles={[roleTypes.PPM]}
                />

                {/* TXO routes, depend on too/tio feature flags */}
                {too && (
                  <PrivateRoute
                    path="/too/:moveOrderId/customer/:customerId"
                    component={CustomerDetails}
                    requiredRoles={[roleTypes.TOO]}
                  />
                )}
                {too && <PrivateRoute path="/moves/queue" exact component={TOO} requiredRoles={[roleTypes.TOO]} />}
                {(too || tio) && (
                  <PrivateRoute
                    path="/moves/:moveOrderId"
                    component={TXOMoveInfo}
                    requiredRoles={[roleTypes.TOO, roleTypes.TIO]}
                  />
                )}
                {too && <Route path="/verification-in-progress" component={TOOVerificationInProgress} />}
                {tio && <PrivateRoute path="/invoicing/queue" component={TIO} requiredRoles={[roleTypes.TIO]} />}
                {tio && (
                  <PrivateRoute
                    path="/payment_requests"
                    component={PaymentRequestIndex}
                    requiredRoles={[roleTypes.TIO]}
                  />
                )}
              </Switch>
            )}
          </Suspense>
        </main>
      </div>
    );
  }
}

OfficeWrapper.propTypes = {
  loadInternalSchema: PropTypes.func.isRequired,
  loadPublicSchema: PropTypes.func.isRequired,
  getCurrentUserInfo: PropTypes.func.isRequired,
  context: PropTypes.shape({
    flags: PropTypes.shape({
      too: PropTypes.bool,
      tio: PropTypes.bool,
    }),
  }),
  location: PropTypes.shape({
    pathname: PropTypes.string,
  }),
};

OfficeWrapper.defaultProps = {
  context: {
    flags: {
      too: false,
      tio: false,
    },
  },
  location: { pathname: '' },
};

const mapStateToProps = (state) => {
  const user = selectCurrentUser(state);
  return {
    swaggerError: state.swaggerInternal.hasErrored,
    userIsLoggedIn: user.isLoggedIn,
  };
};

const mapDispatchToProps = (dispatch) =>
  bindActionCreators(
    {
      loadInternalSchema: loadInternalSchemaAction,
      loadPublicSchema: loadPublicSchemaAction,
      getCurrentUserInfo: getCurrentUserInfoAction,
    },
    dispatch,
  );

export default withContext(withRouter(connect(mapStateToProps, mapDispatchToProps)(OfficeWrapper)));
