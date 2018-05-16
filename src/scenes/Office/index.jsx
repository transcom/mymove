import React, { Component } from 'react';
import { Redirect, Switch } from 'react-router-dom';
import { ConnectedRouter } from 'react-router-redux';
import { history } from 'shared/store';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';

import QueueHeader from 'shared/Header/Office';
import QueueList from './QueueList';
import QueueTable from './QueueTable';
import MoveInfo from './MoveInfo';
import { loadLoggedInUser } from 'shared/User/ducks';
import { loadSchema } from 'shared/Swagger/ducks';
import { no_op } from 'shared/utils';
import LogoutOnInactivity from 'shared/User/LogoutOnInactivity';
import PrivateRoute from 'shared/User/PrivateRoute';


class Queues extends Component {
  render() {
    return (
      <div className="usa-grid grid-wide">
        <div className="usa-width-one-sixth">
          <QueueList />
        </div>
        <div className="usa-width-five-sixths">
          <QueueTable queueType={this.props.match.params.queueType} />
        </div>
      </div>
    );
  }
}

class OfficeWrapper extends Component {
  componentDidMount() {
    document.title = 'Transcom PPP: Office';
    this.props.loadSchema();
  }

  render() {
    return (
      <ConnectedRouter history={history}>
        <div className="Office site">
          <QueueHeader />
          <main className="site__content">
            <div>
              <LogoutOnInactivity />
              <Switch>
                <Redirect from="/" to="/queues/new" exact />
                <PrivateRoute
                  path="/queues/:queueType/moves/:moveId"
                  component={MoveInfo}
                />
                <PrivateRoute path="/queues/:queueType" component={Queues} />
              </Switch>
            </div>
          </main>
        </div>
      </ConnectedRouter>
    );
  }
}

OfficeWrapper.defaultProps = {
  loadSchema: no_op,
  loadLoggedInUser: no_op,
};

const mapStateToProps = state => ({
  swaggerError: state.swagger.hasErrored,
});

const mapDispatchToProps = dispatch =>
  bindActionCreators({ loadSchema, loadLoggedInUser }, dispatch);

export default connect(mapStateToProps, mapDispatchToProps)(OfficeWrapper);
