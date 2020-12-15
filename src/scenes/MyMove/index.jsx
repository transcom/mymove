import React, { Component } from 'react';
import PropTypes from 'prop-types';
import { get } from 'lodash';
import { LastLocationProvider } from 'react-router-last-location';

import ValidatedPrivateRoute from 'shared/User/ValidatedPrivateRoute';
import { Route, Switch } from 'react-router-dom';
import { ConnectedRouter, push, goBack } from 'connected-react-router';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';

import 'uswds';
import '../../../node_modules/uswds/dist/css/uswds.css';
import 'styles/customer.scss';

import Alert from 'shared/Alert';
import InfectedUpload from 'shared/Uploader/InfectedUpload';
import ProcessingUpload from 'shared/Uploader/ProcessingUpload';
import StyleGuide from 'scenes/StyleGuide';
import PpmLanding from 'scenes/PpmLanding';
import Edit from 'scenes/Review/Edit';
import EditProfile from 'scenes/Review/EditProfile';
import EditBackupContact from 'scenes/Review/EditBackupContact';
import EditContactInfo from 'scenes/Review/EditContactInfo';
import EditOrders from 'scenes/Review/EditOrders';
import EditDateAndLocation from 'scenes/Review/EditDateAndLocation';
import EditWeight from 'scenes/Review/EditWeight';
import Header from 'shared/Header/MyMove';
import PPMPaymentRequestIntro from 'scenes/Moves/Ppm/PPMPaymentRequestIntro';
import WeightTicket from 'scenes/Moves/Ppm/WeightTicket';
import ExpensesLanding from 'scenes/Moves/Ppm/ExpensesLanding';
import ExpensesUpload from 'scenes/Moves/Ppm/ExpensesUpload';
import AllowableExpenses from 'scenes/Moves/Ppm/AllowableExpenses';
import WeightTicketExamples from 'scenes/Moves/Ppm/WeightTicketExamples';
import { history } from 'shared/store';
import Footer from 'shared/Footer';
import LogoutOnInactivity from 'shared/User/LogoutOnInactivity';
import PrivacyPolicyStatement from 'shared/Statements/PrivacyAndPolicyStatement';
import AccessibilityStatement from 'shared/Statements/AccessibilityStatement';
import { lastMoveIsCanceled, selectedConusStatus, selectedMoveType } from 'scenes/Moves/ducks';
import { getWorkflowRoutes } from './getWorkflowRoutes';
import { loadInternalSchema } from 'shared/Swagger/ducks';
import SomethingWentWrong from 'shared/SomethingWentWrong';
import { detectIE11, no_op } from 'shared/utils';
import DPSAuthCookie from 'scenes/DPSAuthCookie';
import TrailerCriteria from 'scenes/Moves/Ppm/TrailerCriteria';
import PaymentReview from 'scenes/Moves/Ppm/PaymentReview/index';
import CustomerAgreementLegalese from 'scenes/Moves/Ppm/CustomerAgreementLegalese';
import { withContext } from 'shared/AppContext';
import ConnectedCreateOrEditMtoShipment from 'pages/MyMove/CreateOrEditMtoShipment';
import Home from 'pages/MyMove/Home';
import { loadUser as loadUserAction } from 'store/auth/actions';
import { initOnboarding as initOnboardingAction } from 'store/onboarding/actions';
import { selectServiceMemberFromLoggedInUser, selectCurrentMove } from 'store/entities/selectors';

export class AppWrapper extends Component {
  state = { hasError: false };

  componentDidMount() {
    const { loadUser, loadInternalSchema, initOnboarding } = this.props;

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
    const Tag = detectIE11() ? 'div' : 'main';

    return (
      <ConnectedRouter history={history}>
        <LastLocationProvider>
          <div className="my-move site" id="app-root">
            <Header />
            <Tag role="main" className="site__content my-move-container">
              <div className="usa-grid">
                <LogoutOnInactivity />
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
              {this.state.hasError && <SomethingWentWrong />}
              {!this.state.hasError && !props.swaggerError && (
                <Switch>
                  <Route exact path="/" component={Home} />
                  <Route exact path="/ppm" component={PpmLanding} />
                  <Route exact path="/sm_style_guide" component={StyleGuide} />
                  <Route path="/privacy-and-security-policy" component={PrivacyPolicyStatement} />
                  <Route path="/accessibility" component={AccessibilityStatement} />
                  {getWorkflowRoutes(props)}
                  {props.context.flags.hhgFlow && <ValidatedPrivateRoute exact path="/" component={Home} />}
                  <ValidatedPrivateRoute exact path="/moves/:moveId/edit" component={Edit} />
                  <ValidatedPrivateRoute exact path="/moves/review/edit-profile" component={EditProfile} />
                  <ValidatedPrivateRoute
                    exact
                    path="/moves/:moveId/mto-shipments/:mtoShipmentId/edit-shipment"
                    component={ConnectedCreateOrEditMtoShipment}
                  />
                  <ValidatedPrivateRoute exact path="/moves/review/edit-backup-contact" component={EditBackupContact} />
                  <ValidatedPrivateRoute exact path="/moves/review/edit-contact-info" component={EditContactInfo} />
                  <ValidatedPrivateRoute path="/moves/:moveId/review/edit-orders" component={EditOrders} />
                  <ValidatedPrivateRoute
                    path="/moves/:moveId/review/edit-date-and-location"
                    component={EditDateAndLocation}
                  />
                  <ValidatedPrivateRoute path="/moves/:moveId/review/edit-weight" component={EditWeight} />
                  <ValidatedPrivateRoute exact path="/weight-ticket-examples" component={WeightTicketExamples} />
                  <ValidatedPrivateRoute exact path="/trailer-criteria" component={TrailerCriteria} />
                  <ValidatedPrivateRoute exact path="/allowable-expenses" component={AllowableExpenses} />
                  <ValidatedPrivateRoute exact path="/infected-upload" component={InfectedUpload} />
                  <ValidatedPrivateRoute exact path="/processing-upload" component={ProcessingUpload} />
                  <ValidatedPrivateRoute
                    path="/moves/:moveId/ppm-payment-request-intro"
                    component={PPMPaymentRequestIntro}
                  />
                  <ValidatedPrivateRoute path="/moves/:moveId/ppm-weight-ticket" component={WeightTicket} />
                  <ValidatedPrivateRoute path="/moves/:moveId/ppm-expenses-intro" component={ExpensesLanding} />
                  <ValidatedPrivateRoute path="/moves/:moveId/ppm-expenses" component={ExpensesUpload} />
                  <ValidatedPrivateRoute path="/moves/:moveId/ppm-payment-review" component={PaymentReview} />
                  <ValidatedPrivateRoute exact path="/ppm-customer-agreement" component={CustomerAgreementLegalese} />
                  <ValidatedPrivateRoute path="/dps_cookie" component={DPSAuthCookie} />
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
                  <Route component={this.noMatch} />
                </Switch>
              )}
            </Tag>
            <Footer />
          </div>
          <div id="modal-root"></div>
        </LastLocationProvider>
      </ConnectedRouter>
    );
  }
}

AppWrapper.propTypes = {
  loadInternalSchema: PropTypes.func,
  loadUser: PropTypes.func,
  initOnboarding: PropTypes.func,
  conusStatus: PropTypes.string.isRequired,
  context: PropTypes.shape({
    flags: PropTypes.shape({
      hhgFlow: PropTypes.bool,
      ghcFlow: PropTypes.bool,
    }),
  }).isRequired,
};

AppWrapper.defaultProps = {
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
  const move = selectCurrentMove(state);

  return {
    currentServiceMemberId: serviceMemberId,
    lastMoveIsCanceled: lastMoveIsCanceled(state),
    latestMove: get(state, 'moves.latestMove'),
    moveId: move.id,
    selectedMoveType: selectedMoveType(state),
    conusStatus: selectedConusStatus(state),
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

export default withContext(connect(mapStateToProps, mapDispatchToProps)(AppWrapper));
