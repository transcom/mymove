import React, { Component, lazy } from 'react';
import PropTypes from 'prop-types';
import { LastLocationProvider } from 'react-router-last-location';

import { Route, Switch } from 'react-router-dom';
import { push, goBack } from 'connected-react-router';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';

import 'uswds';
import '../../../node_modules/uswds/dist/css/uswds.css';
import 'styles/customer.scss';

import Header from 'shared/Header/MyMove';
import Alert from 'shared/Alert';
import Footer from 'shared/Footer';
import ConnectedLogoutOnInactivity from 'layout/LogoutOnInactivity';
import SomethingWentWrong from 'shared/SomethingWentWrong';
import CustomerPrivateRoute from 'containers/CustomerPrivateRoute/CustomerPrivateRoute';
import { getWorkflowRoutes } from './getWorkflowRoutes';
import { loadInternalSchema } from 'shared/Swagger/ducks';
import { withContext } from 'shared/AppContext';
import { no_op } from 'shared/utils';
import { loadUser as loadUserAction } from 'store/auth/actions';
import { initOnboarding as initOnboardingAction } from 'store/onboarding/actions';
import { selectConusStatus } from 'store/onboarding/selectors';
import {
  selectServiceMemberFromLoggedInUser,
  selectCurrentMove,
  selectHasCanceledMove,
  selectMoveType,
} from 'store/entities/selectors';
import { generalRoutes, customerRoutes } from 'constants/routes';
/** Pages */
import InfectedUpload from 'shared/Uploader/InfectedUpload';
import ProcessingUpload from 'shared/Uploader/ProcessingUpload';
import PpmLanding from 'scenes/PpmLanding';
import Edit from 'scenes/Review/Edit';
import EditProfile from 'scenes/Review/EditProfile';
import EditBackupContact from 'scenes/Review/EditBackupContact';
import EditContactInfo from 'scenes/Review/EditContactInfo';
import EditOrders from 'scenes/Review/EditOrders';
import EditDateAndLocation from 'scenes/Review/EditDateAndLocation';
import EditWeight from 'scenes/Review/EditWeight';
import PPMPaymentRequestIntro from 'scenes/Moves/Ppm/PPMPaymentRequestIntro';
import WeightTicket from 'scenes/Moves/Ppm/WeightTicket';
import ExpensesLanding from 'scenes/Moves/Ppm/ExpensesLanding';
import ExpensesUpload from 'scenes/Moves/Ppm/ExpensesUpload';
import AllowableExpenses from 'scenes/Moves/Ppm/AllowableExpenses';
import WeightTicketExamples from 'scenes/Moves/Ppm/WeightTicketExamples';
import PrivacyPolicyStatement from 'shared/Statements/PrivacyAndPolicyStatement';
import AccessibilityStatement from 'shared/Statements/AccessibilityStatement';
import TrailerCriteria from 'scenes/Moves/Ppm/TrailerCriteria';
import PaymentReview from 'scenes/Moves/Ppm/PaymentReview/index';
import CustomerAgreementLegalese from 'scenes/Moves/Ppm/CustomerAgreementLegalese';
import ConnectedCreateOrEditMtoShipment from 'pages/MyMove/CreateOrEditMtoShipment';
import Home from 'pages/MyMove/Home';
// Pages should be lazy-loaded (they correspond to unique routes & only need to be loaded when that URL is accessed)
const SignIn = lazy(() => import('pages/SignIn/SignIn'));
const AccessCode = lazy(() => import('shared/User/AccessCode'));
const MovingInfo = lazy(() => import('pages/MyMove/MovingInfo'));

export class CustomerApp extends Component {
  constructor(props) {
    super(props);

    this.state = { hasError: false, error: undefined, info: undefined };
  }

  componentDidMount() {
    const { loadUser, initOnboarding, loadInternalSchema } = this.props;

    loadInternalSchema();
    loadUser();
    initOnboarding();
  }

  componentDidCatch(error, info) {
    this.setState({
      hasError: true,
      error,
      info,
    });
  }

  noMatch = () => (
    <div className="usa-grid">
      <div className="grid-container usa-prose">
        <h1>Page not found</h1>
        <p>Looks like you've followed a broken link or entered a URL that doesn't exist on this site.</p>
        <button className="usa-button" onClick={this.props.goBack}>
          Go Back
        </button>
      </div>
    </div>
  );

  render() {
    const props = this.props;
    const { hasError } = this.state;

    return (
      <>
        <LastLocationProvider>
          <div className="my-move site" id="app-root">
            <Header />

            <main role="main" className="site__content my-move-container" id="main">
              <ConnectedLogoutOnInactivity />

              <div className="usa-grid">
                {props.swaggerError && (
                  <div className="grid-container">
                    <div className="grid-row">
                      <div className="grid-col-12">
                        <Alert type="error" heading="An error occurred">
                          There was an error contacting the server.
                        </Alert>
                      </div>
                    </div>
                  </div>
                )}
              </div>

              {hasError && <SomethingWentWrong />}

              {!hasError && !props.swaggerError && (
                <Switch>
                  {/* no auth */}
                  <Route path={generalRoutes.SIGN_IN_PATH} component={SignIn} />
                  <Route path={customerRoutes.ACCESS_CODE_PATH} component={AccessCode} />
                  <Route path={generalRoutes.PRIVACY_SECURITY_POLICY_PATH} component={PrivacyPolicyStatement} />
                  <Route path={generalRoutes.ACCESSIBILITY_PATH} component={AccessibilityStatement} />

                  {/* auth required */}
                  <CustomerPrivateRoute exact path="/ppm" component={PpmLanding} />

                  {/* ROOT */}
                  <CustomerPrivateRoute path={generalRoutes.HOME_PATH} exact component={Home} />

                  {getWorkflowRoutes(props)}
                  <CustomerPrivateRoute exact path={customerRoutes.SHIPMENT_MOVING_INFO_PATH} component={MovingInfo} />
                  <CustomerPrivateRoute exact path="/moves/:moveId/edit" component={Edit} />
                  <CustomerPrivateRoute exact path="/moves/review/edit-profile" component={EditProfile} />
                  <CustomerPrivateRoute
                    path={customerRoutes.SHIPMENT_CREATE_PATH}
                    component={ConnectedCreateOrEditMtoShipment}
                  />
                  <CustomerPrivateRoute
                    exact
                    path={customerRoutes.SHIPMENT_EDIT_PATH}
                    component={ConnectedCreateOrEditMtoShipment}
                  />
                  <CustomerPrivateRoute exact path="/moves/review/edit-backup-contact" component={EditBackupContact} />
                  <CustomerPrivateRoute exact path="/moves/review/edit-contact-info" component={EditContactInfo} />
                  <CustomerPrivateRoute path="/moves/:moveId/review/edit-orders" component={EditOrders} />
                  <CustomerPrivateRoute
                    path="/moves/:moveId/review/edit-date-and-location"
                    component={EditDateAndLocation}
                  />
                  <CustomerPrivateRoute path="/moves/:moveId/review/edit-weight" component={EditWeight} />
                  <CustomerPrivateRoute exact path="/weight-ticket-examples" component={WeightTicketExamples} />
                  <CustomerPrivateRoute exact path="/trailer-criteria" component={TrailerCriteria} />
                  <CustomerPrivateRoute exact path="/allowable-expenses" component={AllowableExpenses} />
                  <CustomerPrivateRoute exact path="/infected-upload" component={InfectedUpload} />
                  <CustomerPrivateRoute exact path="/processing-upload" component={ProcessingUpload} />
                  <CustomerPrivateRoute
                    path="/moves/:moveId/ppm-payment-request-intro"
                    component={PPMPaymentRequestIntro}
                  />
                  <CustomerPrivateRoute path="/moves/:moveId/ppm-weight-ticket" component={WeightTicket} />
                  <CustomerPrivateRoute path="/moves/:moveId/ppm-expenses-intro" component={ExpensesLanding} />
                  <CustomerPrivateRoute path="/moves/:moveId/ppm-expenses" component={ExpensesUpload} />
                  <CustomerPrivateRoute path="/moves/:moveId/ppm-payment-review" component={PaymentReview} />
                  <CustomerPrivateRoute exact path="/ppm-customer-agreement" component={CustomerAgreementLegalese} />

                  {/* Errors */}
                  <Route exact path="/forbidden">
                    <div className="usa-grid">
                      <h2>You are forbidden to use this endpoint</h2>
                    </div>
                  </Route>
                  <Route exact path="/server_error">
                    <div className="usa-grid">
                      <h2>We are experiencing an internal server error</h2>
                    </div>
                  </Route>

                  {/* 404 */}
                  <Route component={this.noMatch} />
                </Switch>
              )}
            </main>
            <Footer />
          </div>
          <div id="modal-root"></div>
        </LastLocationProvider>
      </>
    );
  }
}

CustomerApp.propTypes = {
  loadInternalSchema: PropTypes.func,
  loadUser: PropTypes.func,
  initOnboarding: PropTypes.func,
  conusStatus: PropTypes.string,
  context: PropTypes.shape({
    flags: PropTypes.shape({
      hhgFlow: PropTypes.bool,
      ghcFlow: PropTypes.bool,
    }),
  }).isRequired,
};

CustomerApp.defaultProps = {
  loadInternalSchema: no_op,
  loadUser: no_op,
  initOnboarding: no_op,
  conusStatus: '',
  context: {
    flags: {
      hhgFlow: false,
      ghcFlow: false,
    },
  },
};

const mapStateToProps = (state) => {
  const serviceMember = selectServiceMemberFromLoggedInUser(state);
  const serviceMemberId = serviceMember?.id;
  const move = selectCurrentMove(state) || {};

  return {
    currentServiceMemberId: serviceMemberId,
    lastMoveIsCanceled: selectHasCanceledMove(state),
    moveId: move?.id,
    selectedMoveType: selectMoveType(state),
    conusStatus: selectConusStatus(state),
    swaggerError: state.swaggerInternal.hasErrored,
  };
};
const mapDispatchToProps = (dispatch) =>
  bindActionCreators(
    {
      goBack,
      push,
      loadInternalSchema,
      loadUser: loadUserAction,
      initOnboarding: initOnboardingAction,
    },
    dispatch,
  );

export default withContext(connect(mapStateToProps, mapDispatchToProps)(CustomerApp));
