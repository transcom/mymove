import React, { Component } from 'react';
import { Redirect, Route, Switch } from 'react-router-dom';
import { ConnectedRouter } from 'react-router-redux';
import { history } from 'shared/store';

import QueueHeader from 'shared/Header/Office';
import QueueList from './QueueList';
import MoveInfo from './MoveInfo';

class QueueTable extends Component {
  render() {
    return (
      <div style={{ background: 'rgb(200,255,200)' }}>
        <h3>QueueTable</h3>
        <p>
          Now showing the <strong>{this.props.queueType}</strong> queue!
        </p>
      </div>
    );
  }
}

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
  }

  render() {
    return (
      <ConnectedRouter history={history}>
        <div className="Office site">
          <QueueHeader />
          <main className="site__content">
            <div>
              <Switch>
                <Redirect from="/" to="/queues/new_moves" exact />
                <Route
                  path="/queues/:queueType/moves/:moveID"
                  component={MoveInfo}
                />
                <Route path="/queues/:queueType" component={Queues} />
              </Switch>
            </div>
          </main>
        </div>
      </ConnectedRouter>
    );
  }
}

export default OfficeWrapper;
