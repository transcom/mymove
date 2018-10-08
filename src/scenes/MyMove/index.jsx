import React, { Component } from 'react';
import { get } from 'lodash';
import { LastLocationProvider } from 'react-router-last-location';

import PrivateRoute from 'shared/User/PrivateRoute';
import { Route, Switch } from 'react-router-dom';
import { ConnectedRouter, push } from 'react-router-redux';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';

import Alert from 'shared/Alert';
import Feedback from 'scenes/Feedback';
import Landing from 'scenes/Landing';
import Shipments from 'scenes/Shipments';
import SubmittedFeedback from 'scenes/SubmittedFeedback';
import Edit from 'scenes/Review/Edit';
import EditProfile from 'scenes/Review/EditProfile';
import EditBackupContact from 'scenes/Review/EditBackupContact';
import EditContactInfo from 'scenes/Review/EditContactInfo';
import EditOrders from 'scenes/Review/EditOrders';
import EditDateAndLocation from 'scenes/Review/EditDateAndLocation';
import EditWeight from 'scenes/Review/EditWeight';
import Header from 'shared/Header/MyMove';
import PaymentRequest from 'scenes/Moves/Ppm/PaymentRequest';
import { history } from 'shared/store';
import Footer from 'shared/Footer';
import LogoutOnInactivity from 'shared/User/LogoutOnInactivity';
import PrivacyPolicyStatement from 'shared/Statements/PrivacyAndPolicyStatement';
import AccessibilityStatement from 'shared/Statements/AccessibilityStatement';
import { selectedMoveType, lastMoveIsCanceled } from 'scenes/Moves/ducks';
import { getWorkflowRoutes } from './getWorkflowRoutes';
import { loadLoggedInUser } from 'shared/User/ducks';
import { loadInternalSchema } from 'shared/Swagger/ducks';
import FailWhale from 'shared/FailWhale';
import { no_op } from 'shared/utils';

const NoMatch = ({ location }) => (
  <div className="usa-grid">
    <h3>
      No match for <code>{location.pathname}</code>
    </h3>
  </div>
);
export class AppWrapper extends Component {
  state = { hasError: false };
  componentDidMount() {
    this.props.loadLoggedInUser();
    this.props.loadInternalSchema();
  }

  componentDidCatch(error, info) {
    this.setState({
      hasError: true,
    });
  }
  render() {
    const props = this.props;
    return (
      <ConnectedRouter history={history}>
        <LastLocationProvider>
          <div className="my-move site">
            <Header />
            <main className="site__content">
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
                    <Route path="/submitted" component={SubmittedFeedback} />
                    <Route
                      path="/shipments/:shipmentsStatus"
                      component={Shipments}
                    />
                    <Route path="/feedback" component={Feedback} />
                    <Route
                      path="/privacy-and-security-policy"
                      component={PrivacyPolicyStatement}
                    />
                    <Route
                      path="/accessibility"
                      component={AccessibilityStatement}
                    />
                    {getWorkflowRoutes(props)}
                    <PrivateRoute
                      exact
                      path="/moves/:moveId/edit"
                      component={Edit}
                    />
                    <PrivateRoute
                      exact
                      path="/moves/review/edit-profile"
                      component={EditProfile}
                    />
                    <PrivateRoute
                      exact
                      path="/moves/review/edit-backup-contact"
                      component={EditBackupContact}
                    />
                    <PrivateRoute
                      exact
                      path="/moves/review/edit-contact-info"
                      component={EditContactInfo}
                    />
                    <PrivateRoute
                      path="/moves/:moveId/review/edit-orders"
                      component={EditOrders}
                    />
                    <PrivateRoute
                      path="/moves/:moveId/review/edit-date-and-location"
                      component={EditDateAndLocation}
                    />
                    <PrivateRoute
                      path="/moves/:moveId/review/edit-weight"
                      component={EditWeight}
                    />
                    <PrivateRoute
                      path="/moves/:moveId/request-payment"
                      component={PaymentRequest}
                    />
                    <Route component={NoMatch} />
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
  loadLoggedInUser: no_op,
};

const mapStateToProps = state => {
  return {
    swaggerError: state.swagger.hasErrored,
    currentServiceMemberId: get(state, 'serviceMember.currentServiceMember.id'),
    selectedMoveType: selectedMoveType(state),
    moveId: get(state, 'moves.currentMove.id'),
    lastMoveIsCanceled: lastMoveIsCanceled(state),
    latestMove: get(state, 'moves.latestMove'),
  };
};
const mapDispatchToProps = dispatch =>
  bindActionCreators({ push, loadInternalSchema, loadLoggedInUser }, dispatch);

export default connect(mapStateToProps, mapDispatchToProps)(AppWrapper);
