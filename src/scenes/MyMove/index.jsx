import React, { Component } from 'react';
import { get } from 'lodash';
import { LastLocationProvider } from 'react-router-last-location';

import ValidatedPrivateRoute from 'shared/User/ValidatedPrivateRoute';
import { Route, Switch } from 'react-router-dom';
import { ConnectedRouter, push, goBack } from 'react-router-redux';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';

import Alert from 'shared/Alert';
import StyleGuide from 'scenes/StyleGuide';
import Landing from 'scenes/Landing';
import Edit from 'scenes/Review/Edit';
import EditProfile from 'scenes/Review/EditProfile';
import EditBackupContact from 'scenes/Review/EditBackupContact';
import EditContactInfo from 'scenes/Review/EditContactInfo';
import EditOrders from 'scenes/Review/EditOrders';
import EditDateAndLocation from 'scenes/Review/EditDateAndLocation';
import EditWeight from 'scenes/Review/EditWeight';
import EditHHGDates from 'scenes/Review/EditShipment';
import Header from 'shared/Header/MyMove';
import PPMPaymentRequestIntro from 'scenes/Moves/Ppm/PPMPaymentRequestIntro';
import WeightTicket from 'scenes/Moves/Ppm/WeightTicket';
import ExpensesLanding from 'scenes/Moves/Ppm/ExpensesLanding';
import ExpensesUpload from 'scenes/Moves/Ppm/ExpensesUpload';
import AllowableExpenses from 'scenes/Moves/Ppm/AllowableExpenses';
import WeightTicketExamples from 'scenes/Moves/Ppm/WeightTicketExamples';
import PaymentRequest from 'scenes/Moves/Ppm/PaymentRequest';
import { history } from 'shared/store';
import Footer from 'shared/Footer';
import LogoutOnInactivity from 'shared/User/LogoutOnInactivity';
import PrivacyPolicyStatement from 'shared/Statements/PrivacyAndPolicyStatement';
import AccessibilityStatement from 'shared/Statements/AccessibilityStatement';
import { selectedMoveType, lastMoveIsCanceled } from 'scenes/Moves/ducks';
import { getWorkflowRoutes } from './getWorkflowRoutes';
import { getCurrentUserInfo } from 'shared/Data/users';
import { loadInternalSchema } from 'shared/Swagger/ducks';
import FailWhale from 'shared/FailWhale';
import { detectIE11, no_op } from 'shared/utils';
import DPSAuthCookie from 'scenes/DPSAuthCookie';
import TrailerCriteria from 'scenes/Moves/Ppm/TrailerCriteria';
import PaymentReview from 'scenes/Moves/Ppm/PaymentReview/index';
import CustomerAgreementLegalese from 'scenes/Moves/Ppm/CustomerAgreementLegalese';

export class AppWrapper extends Component {
  state = { hasError: false };
  componentDidMount() {
    this.props.loadInternalSchema();
    this.props.getCurrentUserInfo();
  }

  componentDidCatch(error, info) {
    this.setState({
      hasError: true,
    });
  }

  noMatch = () => (
    <div className="usa-grid">
      <h2>Page not found</h2>
      <p>Looks like you've followed a broken link or entered a URL that doesn't exist on this site.</p>
      <button onClick={this.props.goBack}>Go Back</button>
    </div>
  );

  render() {
    const props = this.props;
    const Tag = detectIE11() ? 'div' : 'main';

    return (
      <ConnectedRouter history={history}>
        <LastLocationProvider>
          <div className="my-move site">
            <Header />
            <Tag role="main" className="site__content my-move-container">
              <div className="usa-grid">
                <LogoutOnInactivity />
                {props.swaggerError && (
                  <Alert type="error" heading="An error occurred">
                    There was an error contacting the server.
                  </Alert>
                )}
              </div>
              {this.state.hasError && <FailWhale />}
              {!this.state.hasError &&
                !props.swaggerError && (
                  <Switch>
                    <Route exact path="/" component={Landing} />
                    <Route exact path="/sm_style_guide" component={StyleGuide} />
                    <Route path="/privacy-and-security-policy" component={PrivacyPolicyStatement} />
                    <Route path="/accessibility" component={AccessibilityStatement} />
                    {getWorkflowRoutes(props)}
                    <ValidatedPrivateRoute exact path="/moves/:moveId/edit" component={Edit} />
                    <ValidatedPrivateRoute exact path="/moves/review/edit-profile" component={EditProfile} />
                    <ValidatedPrivateRoute
                      exact
                      path="/moves/review/edit-backup-contact"
                      component={EditBackupContact}
                    />
                    <ValidatedPrivateRoute exact path="/moves/review/edit-contact-info" component={EditContactInfo} />

                    <ValidatedPrivateRoute path="/moves/:moveId/review/edit-orders" component={EditOrders} />
                    <ValidatedPrivateRoute
                      path="/moves/:moveId/review/edit-date-and-location"
                      component={EditDateAndLocation}
                    />
                    <ValidatedPrivateRoute path="/moves/:moveId/review/edit-weight" component={EditWeight} />

                    <ValidatedPrivateRoute
                      path="/shipments/:shipmentId/review/edit-hhg-dates"
                      component={EditHHGDates}
                    />
                    {/* <ValidatedPrivateRoute path="/moves/:moveId/review/edit-hhg-locations" component={EditHHGLocations} /> */}
                    {/* <ValidatedPrivateRoute path="/moves/:moveId/review/edit-hhg-weights" component={EditHHGWeights} /> */}

                    <ValidatedPrivateRoute path="/moves/:moveId/request-payment" component={PaymentRequest} />
                    <ValidatedPrivateRoute exact path="/weight-ticket-examples" component={WeightTicketExamples} />
                    <ValidatedPrivateRoute exact path="/trailer-criteria" component={TrailerCriteria} />
                    <ValidatedPrivateRoute exact path="/allowable-expenses" component={AllowableExpenses} />
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
        </LastLocationProvider>
      </ConnectedRouter>
    );
  }
}
AppWrapper.defaultProps = {
  loadInternalSchema: no_op,
  getCurrentUserInfo: no_op,
};

const mapStateToProps = state => {
  return {
    currentServiceMemberId: get(state, 'serviceMember.currentServiceMember.id'),
    lastMoveIsCanceled: lastMoveIsCanceled(state),
    latestMove: get(state, 'moves.latestMove'),
    moveId: get(state, 'moves.currentMove.id'),
    selectedMoveType: selectedMoveType(state),
    swaggerError: state.swaggerInternal.hasErrored,
  };
};
const mapDispatchToProps = dispatch =>
  bindActionCreators({ goBack, push, loadInternalSchema, getCurrentUserInfo }, dispatch);

export default connect(mapStateToProps, mapDispatchToProps)(AppWrapper);
