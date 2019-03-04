import React, { Component } from 'react';
import { get } from 'lodash';
import { LastLocationProvider } from 'react-router-last-location';

import PrivateRoute from 'shared/User/PrivateRoute';
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
import { no_op } from 'shared/utils';
import DPSAuthCookie from 'scenes/DPSAuthCookie';

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
    return (
      <ConnectedRouter history={history}>
        <LastLocationProvider>
          <div className="my-move site">
            <Header />
            <main className="site__content my-move-container">
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
                    <PrivateRoute exact path="/moves/:moveId/edit" component={Edit} />
                    <PrivateRoute exact path="/moves/review/edit-profile" component={EditProfile} />
                    <PrivateRoute exact path="/moves/review/edit-backup-contact" component={EditBackupContact} />
                    <PrivateRoute exact path="/moves/review/edit-contact-info" component={EditContactInfo} />

                    <PrivateRoute path="/moves/:moveId/review/edit-orders" component={EditOrders} />
                    <PrivateRoute path="/moves/:moveId/review/edit-date-and-location" component={EditDateAndLocation} />
                    <PrivateRoute path="/moves/:moveId/review/edit-weight" component={EditWeight} />

                    <PrivateRoute path="/shipments/:shipmentId/review/edit-hhg-dates" component={EditHHGDates} />
                    {/* <PrivateRoute path="/moves/:moveId/review/edit-hhg-locations" component={EditHHGLocations} /> */}
                    {/* <PrivateRoute path="/moves/:moveId/review/edit-hhg-weights" component={EditHHGWeights} /> */}

                    <PrivateRoute path="/moves/:moveId/request-payment" component={PaymentRequest} />
                    <PrivateRoute path="/dps_cookie" component={DPSAuthCookie} />
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
            </main>
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
