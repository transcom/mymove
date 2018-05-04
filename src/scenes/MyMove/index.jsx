import React, { Component } from 'react';
import { Route, Switch } from 'react-router-dom';
import { ConnectedRouter, push } from 'react-router-redux';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';

import Feedback from 'scenes/Feedback';
import Landing from 'scenes/Landing';
import Shipments from 'scenes/Shipments';
import SubmittedFeedback from 'scenes/SubmittedFeedback';
import Header from 'shared/Header/MyMove';
import { history } from 'shared/store';
import Footer from 'shared/Footer';
import Uploader from 'shared/Uploader';
import PrivateRoute from 'shared/User/PrivateRoute';

import { getWorkflowRoutes } from './getWorkflowRoutes';
import { createMove } from 'scenes/Moves/ducks';
import { loadUserAndToken } from 'shared/User/ducks';
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
    this.props.loadUserAndToken();
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
            <Switch>
              <Route exact path="/" component={Landing} />
              <Route path="/submitted" component={SubmittedFeedback} />
              <Route path="/shipments/:shipmentsStatus" component={Shipments} />
              <Route path="/feedback" component={Feedback} />
              <PrivateRoute path="/upload" component={Uploader} />
              {getWorkflowRoutes(props)}
              <Route component={NoMatch} />
            </Switch>
          </main>
          <Footer />
        </div>
      </ConnectedRouter>
    );
  }
}
AppWrapper.defaultProps = {
  loadSchema: no_op,
  loadUserAndToken: no_op,
  loadLoggedInUser: no_op,
};

const mapStateToProps = state => ({
  hasCompleteProfile: false, //todo update this when user service is ready
  selectedMoveType: state.submittedMoves.currentMove
    ? state.submittedMoves.currentMove.selected_move_type
    : 'PPM', // hack: this makes development easier when an eng has to reload a page in the ppm flow over and over but there must be a better way.
  hasMove: Boolean(state.submittedMoves.currentMove),
  moveId: state.submittedMoves.currentMove
    ? state.submittedMoves.currentMove.id
    : null,
});
const mapDispatchToProps = dispatch =>
  bindActionCreators(
    { push, loadSchema, loadLoggedInUser, loadUserAndToken, createMove },
    dispatch,
  );

export default connect(mapStateToProps, mapDispatchToProps)(AppWrapper);
