import React from 'react';
import { Route, Switch } from 'react-router-dom';
import { ConnectedRouter, push } from 'react-router-redux';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';

import DD1299 from 'scenes/DD1299';
import Feedback from 'scenes/Feedback';
import Landing from 'scenes/Landing';
import Shipments from 'scenes/Shipments';
import SubmittedFeedback from 'scenes/SubmittedFeedback';
import Header from 'shared/Header';
import { history } from 'shared/store';
import Footer from 'shared/Footer';
import Uploader from 'shared/Uploader';
import PrivateRoute from 'shared/User/PrivateRoute';
import { getWorkflowRoutes } from './getWorkflowRoutes';

const NoMatch = ({ location }) => (
  <div className="usa-grid">
    <h3>
      No match for <code>{location.pathname}</code>
    </h3>
  </div>
);
export const AppWrapper = props => (
  <ConnectedRouter history={history}>
    <div className="App site">
      <Header />
      <main className="site__content">
        <Switch>
          <Route exact path="/" component={Landing} />
          <Route path="/submitted" component={SubmittedFeedback} />
          <Route path="/shipments/:shipmentsStatus" component={Shipments} />
          <Route path="/DD1299" component={DD1299} />
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

const mapStateToProps = state => ({
  hasCompleteProfile: false, //todo update this when user service is ready
  selectedMoveType: state.submittedMoves.currentMove
    ? state.submittedMoves.currentMove.selected_move_type
    : null,
  hasMove: Boolean(state.submittedMoves.currentMove),
});
const mapDispatchToProps = dispatch => bindActionCreators({ push }, dispatch);

export default connect(mapStateToProps, mapDispatchToProps)(AppWrapper);
