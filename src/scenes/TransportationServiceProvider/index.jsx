import React, { Component } from 'react';
import { Redirect, Switch } from 'react-router-dom';
import { ConnectedRouter } from 'react-router-redux';
import { history } from 'shared/store';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';

import TspHeader from 'shared/Header/Tsp';
import QueueList from './QueueList';
import QueueTable from './QueueTable';
import { loadLoggedInUser } from 'shared/User/ducks';
import { loadSchema } from 'shared/Swagger/ducks';
import { no_op } from 'shared/utils';
import LogoutOnInactivity from 'shared/User/LogoutOnInactivity';
import PrivateRoute from 'shared/User/PrivateRoute';

import './tsp.css';

class Queues extends Component {
  render() {
    return (
      <div className="usa-grid grid-wide queue-columns">
        <div className="queue-menu-column">
          <QueueList />
        </div>
        <div className="queue-list-column">
          <div className="queue-table-scrollable">
            <QueueTable queueType={this.props.match.params.queueType} />
          </div>
        </div>
      </div>
    );
  }
}

class TspWrapper extends Component {
  componentDidMount() {
    document.title = 'Transcom PPP: TSP';
    this.props.loadSchema();
  }

  render() {
    return (
      <ConnectedRouter history={history}>
        <div className="TSP site">
          <TspHeader />
          <main className="site__content">
            <div>
              <LogoutOnInactivity />
              <Switch>
                <Redirect from="/" to="/queues/all" exact />
                <PrivateRoute path="/queues/:queueType" component={Queues} />
              </Switch>
            </div>
          </main>
        </div>
      </ConnectedRouter>
    );
  }
}

TspWrapper.defaultProps = {
  loadSchema: no_op,
  loadLoggedInUser: no_op,
};

const mapStateToProps = state => ({
  swaggerError: state.swagger.hasErrored,
});

const mapDispatchToProps = dispatch =>
  bindActionCreators({ loadSchema, loadLoggedInUser }, dispatch);

export default connect(mapStateToProps, mapDispatchToProps)(TspWrapper);
