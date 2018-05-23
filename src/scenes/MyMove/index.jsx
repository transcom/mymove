import React, { Component } from 'react';
import { get } from 'lodash';
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
import EditProfile from 'scenes/Review/EditProfile';
import EditBackupContact from 'scenes/Review/EditBackupContact';
import EditOrders from 'scenes/Review/EditOrders';
import Header from 'shared/Header/MyMove';
import { history } from 'shared/store';
import Footer from 'shared/Footer';
import LogoutOnInactivity from 'shared/User/LogoutOnInactivity';

import { getWorkflowRoutes } from './getWorkflowRoutes';
import { loadLoggedInUser } from 'shared/User/ducks';
import { loadSchema } from 'shared/Swagger/ducks';
import { no_op } from 'shared/utils';

const NoMatch = ({ location }) => (
  <div className="usa-grid">
    <h3>
      No match for <code>{location.pathname}</code>
    </h3>
  </div>
);
export class AppWrapper extends Component {
  componentDidMount() {
    this.props.loadLoggedInUser();
    this.props.loadSchema();
  }
  render() {
    const props = this.props;
    return (
      <ConnectedRouter history={history}>
        <div className="App site">
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
            {!props.swaggerError && (
              <Switch>
                <Route exact path="/" component={Landing} />
                <Route path="/submitted" component={SubmittedFeedback} />
                <Route
                  path="/shipments/:shipmentsStatus"
                  component={Shipments}
                />
                <Route path="/feedback" component={Feedback} />
                {getWorkflowRoutes(props)}
                <PrivateRoute
                  exact
                  path="/moves/:moveId/review/edit-profile"
                  component={EditProfile}
                />
                <PrivateRoute
                  exact
                  path="/moves/:moveId/review/edit-backup-contact"
                  component={EditBackupContact}
                />
                <PrivateRoute
                  exact
                  path="/moves/:moveId/review/edit-orders"
                  component={EditOrders}
                />
                <Route component={NoMatch} />
              </Switch>
            )}
          </main>
          <Footer />
        </div>
      </ConnectedRouter>
    );
  }
}
AppWrapper.defaultProps = {
  loadSchema: no_op,
  loadLoggedInUser: no_op,
};

const mapStateToProps = state => {
  return {
    hasCompleteProfile: get(
      state.loggedInUser,
      'loggedInUser.service_member.is_profile_complete',
    ),
    swaggerError: state.swagger.hasErrored,
    selectedMoveType: state.submittedMoves.currentMove
      ? state.submittedMoves.currentMove.selected_move_type
      : 'PPM', // hack: this makes development easier when an eng has to reload a page in the ppm flow over and over but there must be a better way.
    hasMove: Boolean(state.submittedMoves.currentMove),
    moveId: state.submittedMoves.currentMove
      ? state.submittedMoves.currentMove.id
      : null,
    currentOrdersId:
      get(state.loggedInUser, 'loggedInUser.service_member.orders[0].id') ||
      get(state.orders, 'currentOrders.id'), // should we get the latest or the first?
  };
};
const mapDispatchToProps = dispatch =>
  bindActionCreators({ push, loadSchema, loadLoggedInUser }, dispatch);

export default connect(mapStateToProps, mapDispatchToProps)(AppWrapper);
